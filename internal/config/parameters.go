package config

import (
	"github.com/muonsoft/errors"
	"github.com/strider2038/digen/internal/di"
)

var errInvalidPath = errors.New("invalid path")

type Parameters struct {
	Version       string        `json:"version" yaml:"version"`
	Container     Container     `json:"container" yaml:"container"`
	ErrorHandling ErrorHandling `json:"errorHandling,omitempty" yaml:"errorHandling,omitempty"`
}

type Container struct {
	Dir string `json:"dir" yaml:"dir"`
}

type ErrorHandling struct {
	New  ErrorOptions     `json:"new,omitempty" yaml:"new,omitempty"`
	Join ErrorOptions     `json:"join,omitempty" yaml:"join,omitempty"`
	Wrap WrapErrorOptions `json:"wrap,omitempty" yaml:"wrap,omitempty"`
}

func (h ErrorHandling) MapToOptions() di.ErrorHandling {
	return di.ErrorHandling{
		New:  h.New.mapToOptions(),
		Join: h.Join.mapToOptions(),
		Wrap: h.Wrap.mapToOptions(),
	}
}

type ErrorOptions struct {
	Pkg  string `json:"pkg,omitempty" yaml:"pkg,omitempty"`
	Func string `json:"func,omitempty" yaml:"func,omitempty"`
}

func (o ErrorOptions) mapToOptions() di.ErrorOptions {
	return di.ErrorOptions{
		Package:  o.Pkg,
		Function: o.Func,
	}
}

type WrapErrorOptions struct {
	Pkg  string `json:"pkg,omitempty" yaml:"pkg,omitempty"`
	Func string `json:"func,omitempty" yaml:"func,omitempty"`
	Verb string `json:"verb,omitempty" yaml:"verb,omitempty"`
}

func (o WrapErrorOptions) mapToOptions() di.ErrorOptions {
	return di.ErrorOptions{
		Package:  o.Pkg,
		Function: o.Func,
		Verb:     o.Verb,
	}
}
