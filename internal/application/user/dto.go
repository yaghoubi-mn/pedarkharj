package app_user

import domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"

type VerifyNumberInput struct {
	Number string `json:"number" validate:"required,e164"`
	Code   uint   `json:"code"`
	Token  string `json:"token" validate:"uuid,omitempty"` // temporary token
}

type NumberInput struct {
	Number string `json:"number"`
}

type SignupUserInput struct {
	Number   string `json:"number" validate:"required,e164"` // TODO: size 13
	Name     string `json:"name" validate:"required,name"`
	Password string `json:"password" validate:"required"`
	Token    string `json:"token" validate:"required,uuid"`
}

type LoginUserInput struct {
	Number   string `json:"number" validator:"required,e164"`
	Password string `json:"password" validator:"size:20"`
}

type RefreshInput struct {
	Refresh string `json:"refresh"`
}

type UserOutput struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

func (u *UserOutput) Fill(user domain_user.User) {
	u.Name = user.Name
	u.Number = user.Number
}
