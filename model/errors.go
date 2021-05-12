package model

import "errors"

var (
	ErrAlreadyExists  = errors.New("already exists")
	ErrNotImplemented = errors.New("not implemented")
	ErrNotFound = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")
	ErrNoAccess = errors.New("no access")
)
