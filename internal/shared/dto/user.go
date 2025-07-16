package shared_dto

type SendOTPInput struct {
	Number string `json:"number" validate:"required,phone_number"`
}

type VerifyOTPInput struct {
	Number string `json:"number" validate:"required,phone_number"`
	OTP    uint   `json:"otp" validate:"required"`
	Token  string `json:"token" validate:"uuid,omitempty"` // temporary token
	Mode   string `json:"mode" validate:"required"`        // signup or reset password
}

type NumberInput struct {
	Number string `json:"number" validate:"required,phone_number"`
}

type SignupUserInput struct {
	Number   string `json:"number" validate:"required,phone_number"` // TODO: size 13
	Name     string `json:"name" validate:"required,name"`
	Password string `json:"password" validate:"required"`
	Token    string `json:"token" validate:"required,uuid"`
}

type LoginUserInput struct {
	Number        string `json:"number" validator:"required,phone_number"`
	InputPassword string `json:"password" validator:"size:30,required"`
}

type RefreshInput struct {
	Refresh string `json:"refresh" validate:"required"`
}

type AvatarChooseInput struct {
	Avatar string `json:"avatar" validate:"requied,url"`
}

type RestPasswordInput struct {
	Number   string `json:"number"`
	Password string `json:"password"`
	Token    string `json:"token"`
}
