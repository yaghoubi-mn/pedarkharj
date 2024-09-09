package rcodes

// response codes

const (
	// global
	INVALID_FIELD       = "invalid_field"
	INVALID_QUERY_PARAM = "invalid_query_param"
	// user
	CODE_SENT_TO_NUMBER = "code_sent_to_number"
	VERIFY_NUMBER_FIRST = "verify_number_first"
	WRONG_CODE          = "wrong_code"
	GO_SIGNUP           = "go_signup"
	OTP_EXPIRED         = "otp_expired"
	ZERO_CODE_FIRST     = "zero_code_first"
)
