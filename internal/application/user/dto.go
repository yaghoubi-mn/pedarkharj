package app_user

import (
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	shared_dto "github.com/yaghoubi-mn/pedarkharj/internal/shared/dto"
)

// user fields that saved in jwt token
type JWTUser struct {
	ID           uint64
	Name         string
	Number       string
	IsRegistered bool
}

type SendOTPInput struct {
	shared_dto.SendOTPInput
}

// func (s SendOTPInput) GetDomainStruct() domain_user.SendOTPInput {
// 	return domain_user.SendOTPInput{
// 		SendOTPInput: s.SendOTPInput,
// 	}
// }

type VerifyOTPInput struct {
	shared_dto.VerifyOTPInput
}

// func (v VerifyOTPInput) GetDomainStruct() domain_user.VerifyOTPInput {
// 	return domain_user.VerifyOTPInput{
// 		VerifyOTPInput: v.VerifyOTPInput,
// 	}
// }

type NumberInput struct {
	shared_dto.NumberInput
}

type SignupUserInput struct {
	shared_dto.SignupUserInput
}

// func (s SignupUserInput) GetDomainStruct() domain_user.SignupUserInput {
// 	return domain_user.SignupUserInput{
// 		SignupUserInput: s.SignupUserInput,
// 	}
// }

type LoginUserInput struct {
	shared_dto.LoginUserInput
}

// func (l LoginUserInput) (storedPassword, salt string, isBlocked, isRegistered bool) domain_user.LoginUserInput {
// 	return domain_user.LoginUserInput{
// 		LoginUserInput: l.LoginUserInput,
// 		RealPassword:   storedPassword,
// 		Salt:           salt,
// 		IsBlocked:      isBlocked,
// 		IsRegistered:   isRegistered,
// 	}
// }

type RefreshInput struct {
	shared_dto.RefreshInput
}

type AvatarChooseInput struct {
	shared_dto.AvatarChooseInput
}

type RestPasswordInput struct {
	shared_dto.RestPasswordInput
}

// func (r RestPasswordInput) GetDomainStruct() domain_user.RestPasswordInput {
// 	return domain_user.RestPasswordInput{
// 		RestPasswordInput: r.RestPasswordInput,
// 	}
// }

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
