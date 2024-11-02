package service_errors

import "errors"

var (
	// global

	// user
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidName         = errors.New("name: invalid name")
	ErrBlockedUser         = errors.New("you are blocked")
	ErrInvalidNumber       = errors.New("number: invalid number")
	ErrInvalidCode         = errors.New("code: invalid code")
	ErrInvalidToken        = errors.New("token: invalid token")
	ErrSmallPassword       = errors.New("password: small password")
	ErrLongPassword        = errors.New("password: long password")
	ErrLongName            = errors.New("name: long name")
	ErrSmallName           = errors.New("name: small name")
	ErrWrongPassword       = errors.New("password: wrong password")
	ErrUserNotRegistered   = errors.New("user not registered")
	ErrInvalidMode         = errors.New("mode: invalid mode")

	// device
	ErrInvalidIP           = errors.New("lastIP: invalid last ip")
	ErrInvalidRefreshToken = errors.New("refresh: invalid refresh token")
	ErrInvalidUserID       = errors.New("userID: invalid userID")
	ErrInvalidUserAgent    = errors.New("useragent: invalid user agent")
)
