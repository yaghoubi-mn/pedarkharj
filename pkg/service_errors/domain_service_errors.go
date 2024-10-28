package service_errors

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrBlockedUser         = errors.New("you are blocked")
	ErrInvalidNumber       = errors.New("number: invalid number")
	ErrInvalidCode         = errors.New("code: invalid code")
	ErrInvalidToken        = errors.New("token: invalid token")
	ErrSmallPassword       = errors.New("password: small password")
	ErrLongPassword        = errors.New("password: long password")
	ErrInvalidName         = errors.New("name: invalid name")
	ErrLongName            = errors.New("name: long name")
	ErrSmallName           = errors.New("name: small name")
	ErrWrongPassword       = errors.New("password: wrong password")
	ErrUserNotRegistered   = errors.New("user not registered")
	ErrInvalidMode         = errors.New("mode: invalid mode")
)
