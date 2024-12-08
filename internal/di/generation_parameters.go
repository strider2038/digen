package di

import (
	"strconv"
	"strings"

	"github.com/dave/jennifer/jen"
)

type GenerationParameters struct {
	RootPackage   string
	ErrorHandling ErrorHandling
	Factories     FactoriesParameters
	Version       string
}

func (params GenerationParameters) Defaults() GenerationParameters {
	params.ErrorHandling = params.ErrorHandling.Defaults()
	if params.Version == "" {
		params.Version = "(unknown version)"
	}

	return params
}

type FactoriesParameters struct {
	SkipError bool
}

func (p FactoriesParameters) ReturnError() bool {
	return !p.SkipError
}

type ErrorHandling struct {
	New  ErrorOptions
	Join ErrorOptions
	Wrap ErrorOptions
}

type ErrorOptions struct {
	Package  string
	Function string
	Verb     string
}

func (w ErrorHandling) Defaults() ErrorHandling {
	if w.New.Package == "" {
		w.New.Package = "fmt"
	}
	if w.New.Function == "" {
		w.New.Function = "Errorf"
	}
	if w.Join.Package == "" {
		w.Join.Package = "errors"
	}
	if w.Join.Function == "" {
		w.Join.Function = "Join"
	}
	if w.Wrap.Package == "" {
		w.Wrap.Package = "fmt"
	}
	if w.Wrap.Function == "" {
		w.Wrap.Function = "Errorf"
	}
	if w.Wrap.Verb == "" {
		w.Wrap.Verb = "%w"
	}

	return w
}

func (params GenerationParameters) rootPackageName() string {
	path := strings.Split(params.RootPackage, "/")
	if len(path) == 0 {
		return ""
	}
	return path[len(path)-1]
}

func (params GenerationParameters) packageName(packageType PackageType) string {
	return strings.Trim(strconv.Quote(params.RootPackage+"/"+packageDirs[packageType]), `"`)
}

func (params GenerationParameters) wrapError(message string, errorIdentifier jen.Code) *jen.Statement {
	path := params.ErrorHandling.Wrap.Package
	funcName := params.ErrorHandling.Wrap.Function
	verb := params.ErrorHandling.Wrap.Verb

	return jen.Qual(path, funcName).Call(jen.Lit(message+": "+verb), errorIdentifier)
}

func (params GenerationParameters) joinErrors(errs ...jen.Code) *jen.Statement {
	path := params.ErrorHandling.Join.Package
	funcName := params.ErrorHandling.Join.Function

	return jen.Qual(path, funcName).Call(errs...)
}
