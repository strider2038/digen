package di

import "errors"

var (
	ErrContainerNotFound = errors.New("container not found")
	ErrUnexpectedType    = errors.New("unexpected type")
	ErrNotSupported      = errors.New("not supported")
	ErrParsing           = errors.New("parsing error")
	ErrFileAlreadyExists = errors.New("file already exists")
	ErrInvalidDefinition = errors.New("invalid definition")

	errMissingModule = errors.New("cannot detect module from go.mod")
)
