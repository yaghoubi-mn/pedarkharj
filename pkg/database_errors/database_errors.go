package database_errors

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrExpired        = errors.New("expired")
)
