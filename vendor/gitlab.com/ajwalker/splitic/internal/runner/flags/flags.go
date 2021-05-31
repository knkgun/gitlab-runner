package flags

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gitlab.com/ajwalker/splitic/internal/timings"
)

type Options struct {
	WorkingDirectory string
	OutputDirectory  string
	EnvPassthrough   FileEntries
	Quarantined      FileEntries
	Flaky            FileEntries
	FlakyRetries     int
	TestFailExitCode int

	PkgList   []string
	NodeIndex int
	NodeTotal int

	Cover    bool
	CoverPkg string
	Race     bool
	Tags     string
	Debug    bool

	CoverReport string
	JUnitReport string

	goTestFlags []string
}

type FileEntries []string

var testingMode bool

func (e *FileEntries) String() string {
	return strings.Join(*e, ",")
}

func (e *FileEntries) Has(value string) bool {
	for _, entry := range *e {
		if entry == value {
			return true
		}
	}

	return false
}

func (e *FileEntries) Set(value string) error {
	f, err := os.Open(value)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(f)
	for s.Scan() {
		*e = append(*e, s.Text())
	}
	f.Close()

	return s.Err()
}

func Parse(name string, args []string, output io.Writer) (timings.Provider, Options) {
	var providerName string
	var options Options

	goTestArgsIndex := -1
	for idx, arg := range args {
		if arg == "--" {
			goTestArgsIndex = idx
			break
		}
	}
	if goTestArgsIndex > -1 {
		options.goTestFlags = args[goTestArgsIndex+1:]
		args = args[:goTestArgsIndex]
	}

	// parse flags just to determine provider first
	pfs := flag.NewFlagSet("", flag.ContinueOnError)
	{
		pfs.StringVar(&providerName, "provider", timings.Default().Name(), "provider of timing data")
		pfs.Usage = func() {}
		pfs.SetOutput(ioutil.Discard)

		pfs.Parse(args)
	}

	// load provider
	provider := timings.Get(providerName)

	// parse all tags
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	{
		fs.SetOutput(output)

		options.NodeIndex, _ = strconv.Atoi(os.Getenv("CI_NODE_INDEX"))
		options.NodeTotal, _ = strconv.Atoi(os.Getenv("CI_NODE_TOTAL"))

		fs.StringVar(&providerName, "provider", provider.Name(), "provider of timing data")
		fs.StringVar(&options.WorkingDirectory, "dir", ".", "working directory")
		fs.StringVar(&options.OutputDirectory, "outputdir", ".", "output directory")
		fs.Var(&options.EnvPassthrough, "env-passthrough", "environment variable passthrough file")
		fs.Var(&options.Quarantined, "quarantined", "a file of quarantined test entries that are allowed to fail")
		fs.Var(&options.Flaky, "flaky", "a file of flaky tests that will be retried")
		fs.IntVar(&options.FlakyRetries, "flaky-retries", 3, "number of times to retry defined flaky tests")
		fs.IntVar(&options.TestFailExitCode, "fail-exit-code", 1, "exit code used specifically for test failures")

		fs.IntVar(&options.NodeIndex, "node-index", options.NodeIndex, "node index determines which test bucket to use")
		fs.IntVar(&options.NodeTotal, "node-total", options.NodeTotal, "node total determines how many tests buckets there are")

		fs.BoolVar(&options.Cover, "cover", false, "output coverage report")
		fs.StringVar(&options.CoverPkg, "coverpkg", "", "cover package")
		fs.BoolVar(&options.Race, "race", false, "enable race detection")
		fs.StringVar(&options.Tags, "tags", "", "build tags")
		fs.BoolVar(&options.Debug, "debug", false, "debug output")

		fs.StringVar(&options.CoverReport, "cover-report", "cover.profile", "cover report name")
		fs.StringVar(&options.JUnitReport, "junit-report", "junit.xml", "junit report name")

		provider.Flags(fs)

		err := fs.Parse(args)
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		if err != nil && !testingMode {
			os.Exit(2)
		}
	}

	if options.NodeIndex <= 0 {
		options.NodeIndex = 1
	}

	if options.NodeTotal <= 0 {
		options.NodeTotal = 1
	}

	options.PkgList = fs.Args()

	return provider, options
}

func (o Options) CoverMode() string {
	if o.Race {
		return "atomic"
	}
	return "count"
}

func (o Options) GoTestFlags(id string) []string {
	args := []string{"test"}
	args = append(args, "-json")

	if o.Tags != "" {
		args = append(args, "-tags", o.Tags)
	}

	if o.Race {
		args = append(args, "-race")
	}

	if o.Cover {
		args = append(args, "-covermode", o.CoverMode())
		args = append(args, "-coverprofile", fmt.Sprintf("cover_%s.profile", id))
	}

	if o.CoverPkg != "" {
		args = append(args, "-coverpkg", o.CoverPkg)
	}

	return append(args, o.goTestFlags...)
}
