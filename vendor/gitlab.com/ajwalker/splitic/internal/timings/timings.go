package timings

import (
	"flag"
)

var providers []Provider

type Report []Timing

type Timing struct {
	Package string
	Method  string
	Timing  float64
}

type Provider interface {
	IsDefault() bool
	Name() string
	Flags(flags *flag.FlagSet)
	Get() (Report, error)
}

func Register(provider Provider) {
	providers = append(providers, provider)
}

func Default() Provider {
	for _, provider := range providers {
		if provider.IsDefault() {
			return provider
		}
	}

	return noneProvider
}

func Get(name string) Provider {
	if name == noneProvider.Name() {
		return noneProvider
	}

	for _, provider := range providers {
		if provider.Name() == name {
			return provider
		}
	}

	return nil
}

var noneProvider = &none{}

type none struct {
}

func (none) IsDefault() bool {
	return false
}

func (none) Name() string {
	return "none"
}

func (none) Flags(flags *flag.FlagSet) {
}

func (none) Get() (Report, error) {
	return Report{}, nil
}
