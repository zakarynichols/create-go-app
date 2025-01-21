package main

import "errors"

// TODO: strike a balance between swallowing errors and showing our own vs simply returning the error that propogates.
var (
	ErrDirExists     = errors.New("create-go-app: directory already exists")
	ErrMkdir         = errors.New("create-go-app: failed to create directory")
	ErrChdir         = errors.New("create-go-app: failed to change directory")
	ErrWkdir         = errors.New("create-go-app: failed to get working directory")
	ErrInitMod       = errors.New("create-go-app: failed to init a module")
	ErrFmt           = errors.New("create-go-app: failed to format code")
	ErrWriteFiles    = errors.New("create-go-app: failed to write files")
	ErrReadModule    = errors.New("create-go-app: failed to read module name")
	ErrEmptyModule   = errors.New("create-go-app: module name cannot be empty")
	ErrCreateDir     = errors.New("create-go-app: failed to create directory")
	ErrLongModule    = errors.New("create-go-app: module name is too long")
	ErrInvalidModule = errors.New("create-go-app: invalid module name")
	ErrNamedFlag     = errors.New("create-go-app: invalid named flag parameters")
	ErrPosArgs       = errors.New("create-go-app: invalid positional arguments")
	ErrCreateFile    = errors.New("create-go-app: failed to create file")
	ErrNamedFlags    = errors.New("create-go-app: please provide only a single named flag")
)
