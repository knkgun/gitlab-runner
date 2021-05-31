package junit

import (
	"flag"

	"gitlab.com/ajwalker/splitic/internal/reports/junit"
	"gitlab.com/ajwalker/splitic/internal/timings"
)

func init() {
	timings.Register(&junitfile{})
}

type junitfile struct {
	filename string
}

func (p *junitfile) Name() string {
	return "junit"
}

func (p *junitfile) IsDefault() bool {
	return false
}

func (p *junitfile) Flags(f *flag.FlagSet) {
	f.StringVar(&p.filename, "junit-filename", "junit.xml", "junit filename")
}

func (p *junitfile) Get() (timings.Report, error) {
	suites, err := junit.Load(p.filename)
	if err != nil {
		return nil, err
	}

	var report timings.Report
	for _, suite := range suites.Suites {
		for _, testcase := range suite.TestCases {
			report = append(report, timings.Timing{
				Package: suite.Package,
				Method:  testcase.Name,
				Timing:  testcase.Time,
			})
		}
	}

	return report, nil
}
