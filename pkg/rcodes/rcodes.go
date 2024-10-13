package rcodes

// response codes

const (
	// global
	InvalidField      = "invalid_field"
	InvalidQueryParam = "invalid_query_param"
	InvalidHeader     = "invalid_header"
	InvalidToken      = "invalid_token"
	Unauthenticated   = "unauthenticated"

	// user
	CodeSendToNumber  = "code_sent_to_number"
	VerifyNumberFirst = "verify_number_first"
	WrongOTP          = "wrong_otp"
	GoSignup          = "go_signup"
	OTPExpired        = "otp_expired"
	ZeroCodeFirst     = "zero_code_first"
	NumberDelay       = "number_delay"
	NumberNotExist    = "number_not_exist"
	AvatarNotFound = "avatar_not_found"
)

type ResponseCode string
