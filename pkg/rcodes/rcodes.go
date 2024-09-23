package rcodes

// response codes

const (
	// global
	InvalidField      = "invalid_field"
	InvalidQueryParam = "invalid_query_param"

	// user
	CodeSendToNumber  = "code_sent_to_number"
	VerifyNumberFirst = "verify_number_first"
	WrongCode         = "wrong_code"
	GoSignup          = "go_signup"
	OTPExpired        = "otp_expired"
	ZeroCodeFirst     = "zero_code_first"
	NumberDelay       = "number_delay"
)

type ResponseCode string
