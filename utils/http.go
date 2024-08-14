package utils

import "errors"

var (
	ErrInternalServer = errors.New("internal server")
	ErrBadRequest     = errors.New("bad request")
	ErrForbidden      = errors.New("forbidden")
)
