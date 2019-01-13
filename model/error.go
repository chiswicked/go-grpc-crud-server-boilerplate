package model

import "errors"

var (
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrNotFound            = errors.New("Your requested Item is not found")
	ErrInvalidParameter    = errors.New("Invalid parameter")
)
