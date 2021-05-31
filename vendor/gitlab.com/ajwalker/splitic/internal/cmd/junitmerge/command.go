package junitmerge

import (
	"flag"
	"fmt"
	"os"

	_ "gitlab.com/ajwalker/splitic/internal/timings/gitlab"
	_ "gitlab.com/ajwalker/splitic/internal/timings/junit"

	"gitlab.com/ajwalker/splitic/internal/reports/junit"
)

type junitMergeCmd struct {
}

func New() *junitMergeCmd {
	return &junitMergeCmd{}
}

func (cmd *junitMergeCmd) Name() string {
	return "junit-merge"
}

func (cmd *junitMergeCmd) Execute() error {
	fs := flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s:\n\n", fs.Name())
		fmt.Fprintf(fs.Output(), "%s junit1 junitN\n\n", fs.Name())
		fs.PrintDefaults()
	}
	fs.Parse(os.Args[2:])

	return junit.Merge(fs.Args(), os.Stdout)
}
