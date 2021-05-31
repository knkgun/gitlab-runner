package runner

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"gitlab.com/ajwalker/splitic/internal/reports/cover"
	"gitlab.com/ajwalker/splitic/internal/reports/junit"
	"gitlab.com/ajwalker/splitic/internal/runner/flags"
	"gitlab.com/ajwalker/splitic/internal/timings"
)

var (
	ErrRunHadFailures  = errors.New("failures/errors occurred during tests")
	ErrRunWasTruncated = errors.New("truncated test output: panic occurred?")

	malformedTestOutputWarning = "The following tests had either no initial ===RUN or final status terminator.\n" +
		"This can occur if logged test output is interfering with the Go test format."
)

type runner struct {
	report  timings.Report
	options flags.Options
	stdout  io.Writer
	stderr  io.Writer
	dir     string
}

type event struct {
	Time    time.Time
	Action  string
	Package string
	Test    string
	Elapsed float64
	Output  string
}

type Config struct {
	Stdout io.Writer
	Stderr io.Writer
}

func Run(report timings.Report, options flags.Options, config *Config) error {
	if config == nil {
		config = new(Config)
		config.Stdout = os.Stdout
		config.Stderr = os.Stderr
	}

	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return fmt.Errorf("creating temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	return (&runner{
		report:  report,
		options: options,
		stdout:  config.Stdout,
		stderr:  config.Stderr,
		dir:     tmpDir,
	}).run()
}

func (r *runner) run() error {
	idx := r.options.NodeIndex - 1
	if idx < 0 || r.options.NodeIndex > r.options.NodeTotal {
		return fmt.Errorf("invalid node index/total (%d/%d)", r.options.NodeIndex, r.options.NodeTotal)
	}

	var buildFlags []string
	if r.options.Tags != "" {
		buildFlags = append(buildFlags, "-tags", r.options.Tags)
	}

	tests, err := list(r.options.WorkingDirectory, buildFlags, r.options.PkgList)
	if err != nil {
		return fmt.Errorf("extracting test names: %w", err)
	}

	buckets := make(Buckets, r.options.NodeTotal)

	for _, test := range tests {
		buckets.Add(r.report, test)
	}

	for idx, bucket := range buckets {
		fmt.Fprintf(r.stderr, "%d tests for index %d/%d, ~%.2f seconds.\n", len(bucket.items), idx+1, r.options.NodeTotal, bucket.time)
	}

	fmt.Fprintf(r.stderr, "Running tests for index %d/%d:\n", r.options.NodeIndex, r.options.NodeTotal)
	suites := make(map[string]*junit.TestSuite)
	failures := make(map[string]struct{})

	runGroupErr := r.runGroups(buckets[idx].RunGroups(), suites, failures)

	// merge cover profiles
	if err := r.mergeCover(); err != nil {
		return errorPrecedence(runGroupErr, fmt.Errorf("merging cover profiles: %w", err))
	}

	var report junit.Report
	for _, suite := range suites {
		report.Suites = append(report.Suites, *suite)
	}

	// save junit test report
	if err := report.Save(filepath.Join(r.options.OutputDirectory, r.options.JUnitReport)); err != nil {
		return errorPrecedence(runGroupErr, fmt.Errorf("saving junit test report: %w", err))
	}

	if len(failures) > 0 {
		return errorPrecedence(runGroupErr, ErrRunHadFailures)
	}

	return runGroupErr
}

func (r *runner) runGroups(groups []RunGroup, suites map[string]*junit.TestSuite, failures map[string]struct{}) error {
	for idx, group := range groups {
		run := group.Run

		attempts := r.options.FlakyRetries
		if attempts < 1 {
			attempts = 1
		}
		for attempt := 1; attempt <= attempts; attempt++ {
			testcases, err := r.runTests(idx, attempt, run, group.Packages, r.options)
			if err != nil {
				return err
			}

			for _, tc := range testcases {
				if _, ok := suites[tc.Classname]; !ok {
					suites[tc.Classname] = &junit.TestSuite{
						Name: tc.Classname,
						Properties: &junit.Properties{
							Property: []junit.Property{
								{Name: "go.version", Value: runtime.Version()},
								{Name: "go.os", Value: runtime.GOOS},
								{Name: "go.arch", Value: runtime.GOARCH},
							},
						},
					}
				}

				suites[tc.Classname].TestCases = append(suites[tc.Classname].TestCases, tc)
			}

			run = run[:0]
			for _, tc := range testcases {
				if r.options.Quarantined.Has(tc.Classname+" "+tc.Name) || r.options.Quarantined.Has(tc.Name) {
					continue
				}

				if len(tc.Error) == 0 && len(tc.Failure) == 0 {
					delete(failures, tc.Classname+" "+tc.Name)
					continue
				}

				if r.options.Flaky.Has(tc.Classname+" "+tc.Name) || r.options.Flaky.Has(tc.Name) {
					run = append(run, tc.Name)
					continue
				}

				failures[tc.Classname+" "+tc.Name] = struct{}{}
			}

			if len(run) == 0 {
				break
			}
		}
	}

	return nil
}

func (r *runner) runTests(idx, attempt int, run, pkgs []string, options flags.Options) ([]junit.TestCase, error) {
	// build run pattern, eg: ^TestOne$|^TestTwo$
	var runPattern strings.Builder
	for idx, name := range run {
		runPattern.WriteString("^")
		runPattern.WriteString(name)
		runPattern.WriteString("$")
		if idx != len(run)-1 {
			runPattern.WriteString("|")
		}
	}

	args := options.GoTestFlags(fmt.Sprintf("%d_%d", idx, attempt))
	args = append(args, "-outputdir", r.dir)
	args = append(args, "-run")
	args = append(args, runPattern.String())
	args = append(args, pkgs...)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if options.Debug {
		fmt.Fprintf(r.stderr, "go %v\n", strings.Join(args, " "))
	}

	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = options.WorkingDirectory
	cmd.Stderr = os.Stderr

	// filter environment variables passed to go test
	if len(options.EnvPassthrough) == 0 {
		cmd.Env = os.Environ()
	} else {
		for _, key := range options.EnvPassthrough {
			val, ok := os.LookupEnv(key)
			if ok {
				cmd.Env = append(cmd.Env, key+"="+val)
			}
		}
	}

	rc, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	testcases, err := r.handle(rc, attempt)
	if err != nil {
		cmd.Wait()
		return testcases, err
	}

	err = cmd.Wait()
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return testcases, nil
	}

	return testcases, err
}

func (r *runner) handle(eventReader io.Reader, attempt int) ([]junit.TestCase, error) {
	var (
		outputs = make(map[string]*strings.Builder)
		bufPool = sync.Pool{
			New: func() interface{} {
				return new(strings.Builder)
			},
		}
		testcases []junit.TestCase
	)

	scanner := bufio.NewScanner(eventReader)

	var panicked bool
	for scanner.Scan() {
		line := scanner.Bytes()

		if r.options.Debug {
			fmt.Fprintln(r.stderr, string(line))
		}

		var e event
		err := json.Unmarshal(line, &e)
		if err != nil {
			r.stderr.Write(line)
			r.stderr.Write([]byte{'\n'})
			continue
		}

		key := e.Package + "/" + e.Test
		if _, ok := outputs[key]; !ok {
			outputs[key] = bufPool.Get().(*strings.Builder)
			outputs[key].Reset()
		}

		switch e.Action {
		case "pass", "fail", "skip":
			if attempt <= 1 || e.Test == "" {
				fmt.Fprintf(r.stdout, "%s %.2fs %s %s\n", e.Action, e.Elapsed, e.Package, e.Test)
			} else {
				fmt.Fprintf(r.stdout, "%s %.2fs %s %s (#%d)\n", e.Action, e.Elapsed, e.Package, e.Test, attempt)
			}

			output := outputs[key].String()
			bufPool.Put(outputs[key])
			delete(outputs, key)

			if e.Test == "" {
				switch {
				case
					strings.HasPrefix(output, "FAIL\n"),
					strings.HasPrefix(output, "PASS\n"),
					strings.HasPrefix(output, "SKIP\n"):
				default:
					fmt.Fprintln(r.stdout, output)
					panicked = true
				}
				continue
			}

			tc := junit.TestCase{
				Classname: e.Package,
				Name:      e.Test,
				Status:    strings.ToUpper(e.Action),
				Time:      e.Elapsed,
			}

			switch e.Action {
			case "fail":
				tc.Failure = []junit.Failure{{
					Message:  "Failed",
					Contents: output,
				}}

				// for failures, we always output the test contents
				fmt.Fprintln(r.stdout, output)

			case "skip":
				tc.Skipped = output
			}

			testcases = append(testcases, tc)

		case "output":
			outputs[key].WriteString(e.Output)
		}
	}

	if len(outputs) > 0 {
		fmt.Fprintln(r.stderr, malformedTestOutputWarning)
		for key := range outputs {
			fmt.Fprintln(r.stderr, "\t", key)
		}
	}

	if panicked {
		return testcases, ErrRunWasTruncated
	}

	return testcases, nil
}

func (r *runner) mergeCover() error {
	if !r.options.Cover {
		return nil
	}

	filename := r.options.CoverReport
	if !filepath.IsAbs(r.options.CoverReport) {
		filename = filepath.Join(r.options.OutputDirectory, r.options.CoverReport)
	}

	if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
		return fmt.Errorf("creating output directory for cover profile: %w", err)
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	matches, _ := filepath.Glob(filepath.Join(r.dir, "cover_*.profile"))
	if err := cover.Merge(matches, f); err != nil {
		return err
	}

	return f.Close()
}

func errorPrecedence(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
