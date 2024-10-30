package domain_user

type SendOTPInput struct {
	Number string `validate:"required,phone_number"`
}

type VerifyOTPInput struct {
	Number string `validate:"required,phone_number"`
	OTP    uint
	Token  string `validate:"uuid,omitempty"` // temporary token
	Mode   string // signup or reset password
}

type NumberInput struct {
	Number string
}

type SignupUserInput struct {
	Number   string `validate:"required,phone_number"` // TODO: size 13
	Name     string `validate:"required,name"`
	Password string `validate:"required"`
	Token    string `validate:"required,uuid"`
}

type LoginUserInput struct {
	Number        string `validator:"required,phone_number"`
	InputPassword string `validator:"size:20"`

	RealPassword string // password that stored in database
	Salt         string
	IsBlocked    bool
	IsRegistered bool
}

type RefreshInput struct {
	Refresh string
}

type AvatarChooseInput struct {
	Avatar string
}

type RestPasswordInput struct {
	Number   string
	Password string
	Token    string
}

func (v SignupUserInput) GetUser() User {
	return User{
		Number:   v.Number,
		Name:     v.Name,
		Password: v.Password,
	}
}
