package shared_dto

type SendOTPInput struct {
	PhoneNumber string `json:"number" validate:"required,phone_number"`
}

type VerifyOTPInput struct {
	PhoneNumber string `json:"number" validate:"required,phone_number"`
	OTP         uint   `json:"otp" validate:"required"`
	Token       string `json:"token" validate:"uuid,omitempty"` // temporary token
	Mode        string `json:"mode" validate:"required"`        // signup or reset password
}

type NumberInput struct {
	PhoneNumber string `json:"number" validate:"required,phone_number"`
}

type SignupUserInput struct {
	PhoneNumber string `json:"number" validate:"required,phone_number"` // TODO: size 13
	Name        string `json:"name" validate:"required,name"`
	Password    string `json:"password" validate:"required"`
	Token       string `json:"token" validate:"required,uuid"`
}

type LoginUserInput struct {
	PhoneNumber   string `json:"number" validator:"required,phone_number"`
	InputPassword string `json:"password" validator:"size:30,required"`
}

type RefreshInput struct {
	Refresh string `json:"refresh" validate:"required"`
}

type AvatarChooseInput struct {
	Avatar string `json:"avatar" validate:"requied,url"`
}

type ResetPasswordInput struct {
	PhoneNumber string `json:"number"`
	Password    string `json:"password"`
	Token       string `json:"token"`
}
