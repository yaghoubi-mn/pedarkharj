package service_errors

import "errors"

var (
	ErrOTPNotExpired                            = errors.New("otp not expired. wait some minutes")
	ErrOTPNotSend                               = errors.New("code: OTP wasn't sent. go send-otp first")
	ErrUserNotRegisteredResetPasswordNotAllowed = errors.New("user not exist. reset_password not allowed")
	ErrUserAlreayRegisteredSignupNotAllowed     = errors.New("user already exist. signup mode not allowed")
	ErrVerifyNumberFirst                        = errors.New("verify number first")
	ErrNumberNotExist                           = errors.New("number: number not exist")
	ErrAvatarNotFound                           = errors.New("avatar: avatar not found")
	ErrWrongOTP                                 = errors.New("otp: wrong otp")
	ErrWrongToken                               = errors.New("token: wrong token")
)
