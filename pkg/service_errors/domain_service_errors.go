package service_errors

import "errors"

var (
	// global
	ErrPermissionDenied   = errors.New("permission denied")
	ErrInvalidDescription = errors.New("description: invalid description")
	ErrInvalidPage        = errors.New("page: invalid page")
	ErrInvalidLimit       = errors.New("limit: invalid limit")
	ErrInvalidID          = errors.New("id: invalid id")

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
	ErrInvalidUserAgent    = errors.New("useragent: invalid user agent")

	// expense
	ErrInvalidCredit                     = errors.New("credit: invalid credit")
	ErrDebtIsNotPaid                     = errors.New("debt is not paid yet")
	ErrEmptyCreditors                    = errors.New("creditors: creditors cannot be empty")
	ErrCommonCreditorAndDebtor           = errors.New("debtors: list of creditors and debtors cannot overlap")
	ErrLowCredit                         = errors.New("creditors: credit is too low")
	ErrCreatorNustBeInCreditorsOrDebtors = errors.New("creator must be in creditors or debtors")
)
