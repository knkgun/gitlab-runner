package main

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/ajwalker/splitic/internal/cmd/covermerge"
	"gitlab.com/ajwalker/splitic/internal/cmd/junitcheck"
	"gitlab.com/ajwalker/splitic/internal/cmd/junitmerge"
	"gitlab.com/ajwalker/splitic/internal/cmd/test"
)

type Command interface {
	Name() string
	Execute() error
}

type commands []Command

func main() {
	cmds := commands{
		test.New(),
		junitmerge.New(),
		junitcheck.New(),
		covermerge.New(),
	}

	defaults := func() {
		var names []string
		for _, cmd := range cmds {
			names = append(names, cmd.Name())
		}

		fmt.Fprintln(os.Stderr, "subcommands:", strings.Join(names, ", "))
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		defaults()
	}

	for _, cmd := range cmds {
		if os.Args[1] == cmd.Name() {
			if err := cmd.Execute(); err != nil {
				panic(err)
			}
			return
		}
	}

	defaults()
}
