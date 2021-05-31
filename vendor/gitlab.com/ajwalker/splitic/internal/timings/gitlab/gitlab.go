package gitlab

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"time"

	"gitlab.com/ajwalker/splitic/internal/timings"
)

var defaultClient = &http.Client{
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	},
}

type results struct {
	Suites []testsuite `json:"test_suites"`
}

type testsuite struct {
	Name  string     `json:"name"`
	Cases []testcase `json:"test_cases"`
}

type testcase struct {
	Class  string  `json:"classname"`
	Name   string  `json:"name"`
	Timing float64 `json:"execution_time"`
}

func init() {
	timings.Register(&gitlab{})
}

type gitlab struct {
	endpoint string
	project  string
	branch   string
	date     string
	pattern  regexpFlag
}

type regexpFlag struct {
	re *regexp.Regexp
}

func (rf *regexpFlag) Set(value string) error {
	re, err := regexp.Compile(value)
	rf.re = re
	return err
}

func (rf *regexpFlag) String() string {
	if rf.re == nil {
		return ""
	}
	return rf.re.String()
}

func (p *gitlab) Name() string {
	return "gitlab"
}

func (p *gitlab) IsDefault() bool {
	return os.Getenv("GITLAB_CI") != ""
}

func (p *gitlab) Flags(f *flag.FlagSet) {
	endpoint := os.Getenv("CI_API_V4_URL")
	if endpoint == "" {
		endpoint = "https://gitlab.com/api/v4/"
	}

	branch := os.Getenv("CI_DEFAULT_BRANCH")
	if branch == "" {
		branch = "master"
	}

	date := os.Getenv("CI_PIPELINE_CREATED_AT")
	if date == "" {
		date = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	}

	f.StringVar(&p.endpoint, "gitlab-endpoint", endpoint, "api endpoint")
	f.StringVar(&p.project, "gitlab-project", os.Getenv("CI_PROJECT_ID"), "project")
	f.StringVar(&p.branch, "gitlab-branch", branch, "branch")
	f.StringVar(&p.date, "gitlab-date", date, "only consider finished pipelines before this date")
	f.Var(&p.pattern, "gitlab-suite-pattern", "regex pattern to select test suites")
}

func (p *gitlab) Get() (timings.Report, error) {
	var report timings.Report

	pID, err := p.getSuccessfulLatestPipelineID()
	if err != nil {
		return report, fmt.Errorf("fetching latest pipeline: %w", err)
	}

	var results results

	err = fetch(fmt.Sprintf("%s/projects/%s/pipelines/%d/test_report", p.endpoint, p.project, pID), &results)
	if err != nil {
		return report, err
	}

	for _, suite := range results.Suites {
		if p.pattern.re != nil && !p.pattern.re.MatchString(suite.Name) {
			continue
		}

		for _, testcase := range suite.Cases {
			report = append(report, timings.Timing{
				Package: testcase.Class,
				Method:  testcase.Name,
				Timing:  testcase.Timing,
			})
		}
	}

	return report, nil
}

func (p *gitlab) getSuccessfulLatestPipelineID() (uint64, error) {
	var results []struct {
		ID uint64 `json:"id"`
	}

	err := fetch(fmt.Sprintf("%s/projects/%s/pipelines?ref=%s&updated_before=%s&status=success", p.endpoint, p.project, p.branch, p.date), &results)
	if err != nil {
		return 0, err
	}

	if len(results) == 0 {
		return 0, fmt.Errorf("no results found")
	}

	return results[0].ID, nil
}

func fetch(url string, results interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating timing request: %w", err)
	}

	switch {
	case os.Getenv("SPLITIC_GITLAB_TOKEN") != "":
		req.Header.Set("PRIVATE-TOKEN", os.Getenv("SPLITIC_GITLAB_TOKEN"))

	case os.Getenv("SPLITIC_GITLAB_TOKEN_PATH") != "":
		token, _ := ioutil.ReadFile(os.Getenv("SPLITIC_GITLAB_TOKEN_PATH"))
		req.Header.Set("PRIVATE-TOKEN", string(token))

	case os.Getenv("CI_JOB_TOKEN") != "":
		// attempt to use CI_JOB_TOKEN, although, this token doesn't (yet?) allow access
		// to these endpoints.
		req.Header.Set("JOB-TOKEN", os.Getenv("CI_JOB_TOKEN"))
	}

	resp, err := defaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("performing timing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		snippet, _ := ioutil.ReadAll(io.LimitReader(resp.Body, 200))

		return fmt.Errorf("non-200 response (%s): %s", url, string(snippet))
	}

	if err := json.NewDecoder(resp.Body).Decode(results); err != nil {
		return fmt.Errorf("decoding results (%s): %w", url, err)
	}

	return nil
}
