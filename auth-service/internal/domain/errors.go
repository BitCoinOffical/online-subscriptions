package domain

import "errors"

var ErrNotFound = errors.New("not found")
var ErrEmptyPayload = errors.New("empty payload")
var ErrPasswordMismatch = errors.New("the passwords don't match")
