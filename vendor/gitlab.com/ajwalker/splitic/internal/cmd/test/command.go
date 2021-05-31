package test

import (
	"errors"
	"fmt"
	"os"

	_ "gitlab.com/ajwalker/splitic/internal/timings/gitlab"
	_ "gitlab.com/ajwalker/splitic/internal/timings/junit"

	"gitlab.com/ajwalker/splitic/internal/runner"
	"gitlab.com/ajwalker/splitic/internal/runner/flags"
)

type testCmd struct {
}

func New() *testCmd {
	c := &testCmd{}
	return c
}

func (cmd *testCmd) Name() string {
	return "test"
}

func (cmd *testCmd) Execute() error {
	provider, options := flags.Parse(os.Args[1], os.Args[2:], os.Stderr)

	// get timings
	report, err := provider.Get()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error collecting timing information:", err)
		os.Exit(1)
	}

	// run tests
	err = runner.Run(report, options, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running tests:", err)

		if errors.Is(err, runner.ErrRunHadFailures) {
			os.Exit(options.TestFailExitCode)
		}
		os.Exit(1)
	}

	return nil
}
