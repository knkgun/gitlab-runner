package runner

import (
	"math"
	"sort"
	"strings"

	"gitlab.com/ajwalker/splitic/internal/timings"
)

type Buckets []Bucket

// many tests are reported as taking zero time, however in practice
// there's an overhead for every test run, no matter how small so
// we assume each test takes at least 200ms.
const minTime = 0.2

type RunGroup struct {
	Packages []string
	Run      []string
}

func (bs *Buckets) Add(r timings.Report, item testcase) {
	timed := getTiming(r, item.pkg, item.name)

	chosen := 0
	lowest := math.MaxFloat64
	for idx, b := range *bs {
		score := b.time + timed.Timing
		if score < lowest {
			lowest = score
			chosen = idx
		}
	}

	(*bs)[chosen].time += timed.Timing
	(*bs)[chosen].items = append((*bs)[chosen].items, timed)
}

type Bucket struct {
	items []timings.Timing
	time  float64
}

func (b *Bucket) RunGroups() []RunGroup {
	var (
		seen      = make(map[string]struct{})
		ambigious = make(map[string]struct{})
		tests     = make(map[string][]string)
	)

	for _, item := range b.items {
		if _, ok := seen[item.Method]; ok {
			ambigious[item.Package] = struct{}{}
		}

		seen[item.Method] = struct{}{}
		tests[item.Package] = append(tests[item.Package], item.Method)
	}

	groups := make(map[string]RunGroup)
	var sorted []string
	for pkg, methods := range tests {
		if _, ok := ambigious[pkg]; ok {
			group := groups[pkg]
			group.Packages = []string{pkg}
			group.Run = methods
			groups[pkg] = group

			sorted = append(sorted, pkg)
			continue
		}

		group, ok := groups[""]
		if !ok {
			sorted = append(sorted, "")
		}
		group.Packages = append(group.Packages, pkg)
		group.Run = append(group.Run, methods...)
		groups[""] = group
	}

	var result []RunGroup

	sort.Strings(sorted)
	for _, name := range sorted {
		group := groups[name]
		sort.Strings(group.Run)
		sort.Strings(group.Packages)
		result = append(result, group)
	}

	return result
}

func getTiming(r timings.Report, pkg, method string) timings.Timing {
	var length int
	var timing float64
	for _, tc := range r {
		if tc.Method != method {
			continue
		}

		if strings.HasSuffix(pkg, tc.Package) && tc.Method == method {
			if len(tc.Package) == length && tc.Timing > timing {
				timing = tc.Timing
				continue
			}

			if len(tc.Package) > length {
				length = len(tc.Package)
				timing = tc.Timing
			}
		}
	}

	return timings.Timing{
		Package: pkg,
		Method:  method,
		Timing:  math.Max(minTime, timing),
	}
}
