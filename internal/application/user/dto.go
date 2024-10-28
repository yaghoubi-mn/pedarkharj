package app_user

import domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"

// user fields that saved in jwt token
type JWTUser struct {
	ID           uint64
	Name         string
	Number       string
	IsRegistered bool
}

type VerifyNumberInput struct {
	Number string `json:"number" validate:"required,phone_number"`
	OTP    uint   `json:"otp"`
	Token  string `json:"token" validate:"uuid,omitempty"` // temporary token
	Mode   string `json:"mode"`                            // signup or reset password
}

type NumberInput struct {
	Number string `json:"number"`
}

type SignupUserInput struct {
	Number   string `json:"number" validate:"required,phone_number"` // TODO: size 13
	Name     string `json:"name" validate:"required,name"`
	Password string `json:"password" validate:"required"`
	Token    string `json:"token" validate:"required,uuid"`
}

type LoginUserInput struct {
	Number   string `json:"number" validator:"required,phone_number"`
	Password string `json:"password" validator:"size:20"`
}

type RefreshInput struct {
	Refresh string `json:"refresh"`
}

type AvatarChooseInput struct {
	Avatar string `json:"avatar"`
}

type RestPasswordInput struct {
	Number   string `json:"number"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type UserOutput struct {
	Name   string `json:"name"`
	Number string `json:"number"` // number must be fill from user contact for security. user contact may be empty for adding unknown user to user contact
	Avatar string `json:"avatar"`
}

func (u *UserOutput) Fill(user domain_user.User) {
	u.Name = user.Name
	u.Number = user.Number
	u.Avatar = user.Avatar
}
