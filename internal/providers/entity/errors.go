package entity

import "errors"

var (
	ErrAlreadyExists       = errors.New("already exists")
	ErrNotFound            = errors.New("not found")
	ErrInternalServerError = errors.New("internal server error")
)
