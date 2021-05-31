package junitcheck

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"gitlab.com/ajwalker/splitic/internal/runner/flags"
	_ "gitlab.com/ajwalker/splitic/internal/timings/gitlab"
	_ "gitlab.com/ajwalker/splitic/internal/timings/junit"

	"gitlab.com/ajwalker/splitic/internal/reports/junit"
)

type junitCheckCmd struct {
	quarantined flags.FileEntries
	flaky       flags.FileEntries
}

func New() *junitCheckCmd {
	return &junitCheckCmd{}
}

func (cmd *junitCheckCmd) Name() string {
	return "junit-check"
}

func (cmd *junitCheckCmd) Execute() error {
	fs := flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	fs.Var(&cmd.quarantined, "quarantined", "a file of quarantined test entries that are allowed to fail")
	fs.Var(&cmd.flaky, "flaky", "tests that are allowed to fail")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s:\n\n", fs.Name())
		fmt.Fprintf(fs.Output(), "%s junit1 junitN\n\n", fs.Name())
		fs.PrintDefaults()
	}
	fs.Parse(os.Args[2:])

	f, err := ioutil.TempFile("", "")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()

	if err := junit.Merge(fs.Args(), f); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	report, err := junit.Load(f.Name())
	os.Remove(f.Name())
	if err != nil {
		return err
	}

	var unknown []string
	var passing []string
	var failing []string
	for _, name := range cmd.quarantined {
		tc := find(report.Suites, name)
		if tc == nil {
			unknown = append(unknown, name)
			continue
		}

		if len(tc.Failure) == 0 && len(tc.Error) == 0 {
			passing = append(passing, name)
		} else {
			failing = append(failing, name)
		}
	}

	for _, flaky := range cmd.flaky {
		tc := find(report.Suites, flaky)
		if tc == nil {
			unknown = append(unknown, flaky)
			continue
		}
	}

	if len(passing) > 0 || len(failing) > 0 {
		fmt.Printf("%d quarantined tests passed, %d failed\n", len(passing), len(failing))
		for _, name := range passing {
			fmt.Printf("\tpass: %s\n", name)
		}
		for _, name := range failing {
			fmt.Printf("\tfail: %s\n", name)
		}
	}

	if len(failing) > 0 {
		fmt.Printf("found %d defined but unknown/unrun tests\n", len(unknown))
		for _, name := range unknown {
			fmt.Printf("\tunknown: %s\n", name)
		}

		os.Exit(1)
	}

	return nil
}

func find(report []junit.TestSuite, name string) *junit.TestCase {
	for i, suites := range report {
		for j, testcase := range suites.TestCases {
			if name == testcase.Classname+" "+testcase.Name || name == testcase.Name {
				return &report[i].TestCases[j]
			}
		}
	}

	return nil
}
