package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrEmptyPayload = errors.New("empty payload")
var ErrPasswordMismatch = errors.New("the passwords don't match")
var ErrInvalidCredentials = errors.New("incorrect email or password")
var ErrEmailAlreadyExists = errors.New("such user is already registered")
