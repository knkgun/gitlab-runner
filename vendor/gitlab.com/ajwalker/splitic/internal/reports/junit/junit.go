package junit

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Report struct {
	XMLName xml.Name `xml:"testsuites"`

	Name     string  `xml:"name,attr"`
	Time     float64 `xml:"time,attr,omitempty"`
	Tests    int     `xml:"tests,attr,omitempty"`
	Failures int     `xml:"failures,attr,omitempty"`
	Disabled int     `xml:"disabled,attr,omitempty"`
	Errors   int     `xml:"errors,attr,omitempty"`

	Suites []TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	Name      string  `xml:"name,attr,omitempty"`
	Tests     int     `xml:"tests,attr,omitempty"`
	Failures  int     `xml:"failures,attr,omitempty"`
	Errors    int     `xml:"errors,attr,omitempty"`
	Time      float64 `xml:"time,attr,omitempty"`
	Disabled  int     `xml:"disabled,attr,omitempty"`
	Skipped   int     `xml:"skipped,attr,omitempty"`
	Timestamp string  `xml:"timestamp,attr,omitempty"`
	Hostname  string  `xml:"hostname,attr,omitempty"`
	ID        int     `xml:"id,attr"`
	Package   string  `xml:"package,attr,omitempty"`

	Properties *Properties `xml:"properties,omitempty"`
	TestCases  []TestCase  `xml:"testcase"`
	SystemOut  string      `xml:"system-out,omitempty"`
	SystemErr  string      `xml:"system-err,omitempty"`
}

type Properties struct {
	Property []Property `xml:"property,omitempty"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type TestCase struct {
	Name       string  `xml:"name,attr"`
	Assertions int     `xml:"assertions,attr,omitempty"`
	Time       float64 `xml:"time,attr"`
	Classname  string  `xml:"classname,attr,omitempty"`
	Status     string  `xml:"status,attr,omitempty"`

	Skipped   string    `xml:"skipped,omitempty"`
	Error     []Failure `xml:"error,omitempty"`
	Failure   []Failure `xml:"failure,omitempty"`
	SystemOut []string  `xml:"system-out,omitempty"`
	SystemErr []string  `xml:"system-err,omitempty"`
}

type Failure struct {
	Type     string `xml:"type,attr,omitempty"`
	Message  string `xml:"message,attr,omitempty"`
	Contents string `xml:",chardata"`
}

func (s *Report) write(w io.Writer) error {
	for idx := range s.Suites {
		suite := &s.Suites[idx]

		for _, testcase := range suite.TestCases {
			suite.Tests++

			if len(testcase.Skipped) > 0 {
				suite.Skipped++
			}
			if len(testcase.Error) > 0 {
				suite.Errors++
			}
			if len(testcase.Failure) > 0 {
				suite.Failures++
			}
		}

		s.Tests += suite.Tests
		s.Failures += suite.Failures
		s.Disabled += suite.Disabled
		s.Errors += suite.Errors
	}

	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return fmt.Errorf("writing junit xml header: %w", err)
	}

	enc := xml.NewEncoder(w)
	enc.Indent("", " ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("encoding junit: %w", err)
	}

	return nil
}

func Merge(inputs []string, w io.Writer) error {
	var merged Report

	for _, input := range inputs {
		suites, err := Load(input)
		if err != nil {
			return err
		}

		merged.Suites = append(merged.Suites, suites.Suites...)
	}

	return merged.write(w)
}

func Load(filename string) (*Report, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening junit: %w", err)
	}
	defer f.Close()

	var report Report
	if err := xml.NewDecoder(f).Decode(&report); err != nil {
		return nil, fmt.Errorf("decoding junit %q: %w", filename, err)
	}

	return &report, nil
}

func (s *Report) Save(filename string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
		return fmt.Errorf("creating output directory for junit: %w", err)
	}

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating junit: %w", err)
	}
	defer f.Close()

	if err := s.write(f); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("saving junit %q: %w", filename, err)
	}

	return nil
}
