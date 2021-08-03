package digen

import "errors"

var (
	ErrContainerNotFound = errors.New("container not found")
	ErrUnexpectedType    = errors.New("unexpected type")
	ErrNotSupported      = errors.New("not supported")
	ErrParsing           = errors.New("parsing error")

	errFileIgnored = errors.New("file ignored")
)
