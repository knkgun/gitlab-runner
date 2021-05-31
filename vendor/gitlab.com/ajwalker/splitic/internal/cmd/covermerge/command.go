package covermerge

import (
	"flag"
	"fmt"
	"os"

	_ "gitlab.com/ajwalker/splitic/internal/timings/gitlab"
	_ "gitlab.com/ajwalker/splitic/internal/timings/junit"

	"gitlab.com/ajwalker/splitic/internal/reports/cover"
)

type coverMergeCmd struct {
}

func New() *coverMergeCmd {
	return &coverMergeCmd{}
}

func (cmd *coverMergeCmd) Name() string {
	return "cover-merge"
}

func (cmd *coverMergeCmd) Execute() error {
	fs := flag.NewFlagSet("cover-merge", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage of %s:\n\n", fs.Name())
		fmt.Fprintf(fs.Output(), "%s cover1 coverN\n\n", fs.Name())
		fs.PrintDefaults()
	}
	fs.Parse(os.Args[2:])

	return cover.Merge(fs.Args(), os.Stdout)
}
