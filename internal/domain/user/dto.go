package domain_user

import shared_dto "github.com/yaghoubi-mn/pedarkharj/internal/shared/dto"

type SendOTPInput struct {
	shared_dto.SendOTPInput
}

func NewSendOTPInput(number string) SendOTPInput {
	return SendOTPInput{
		SendOTPInput: shared_dto.SendOTPInput{
			PhoneNumber: number,
		},
	}
}

type VerifyOTPInput struct {
	shared_dto.VerifyOTPInput
}

func NewVerifyOTPInput(number string, otp uint, token, mode string) VerifyOTPInput {
	return VerifyOTPInput{
		shared_dto.VerifyOTPInput{
			PhoneNumber: number,
			OTP:         otp,
			Token:       token,
			Mode:        mode,
		},
	}
}

type NumberInput struct {
	shared_dto.NumberInput
}

func NewNumberInput(number string) NumberInput {
	return NumberInput{
		shared_dto.NumberInput{
			PhoneNumber: number,
		},
	}
}

type SignupUserInput struct {
	shared_dto.SignupUserInput
}

func NewSignupUserInput(number, name, password, token string) SignupUserInput {
	return SignupUserInput{
		shared_dto.SignupUserInput{
			PhoneNumber: number,
			Name:        name,
			Password:    password,
			Token:       token,
		},
	}
}

func (v SignupUserInput) GetUser() User {
	return User{
		Number:   v.PhoneNumber,
		Name:     v.Name,
		Password: v.Password,
	}
}

type LoginUserInput struct {
	shared_dto.LoginUserInput

	StoredPassword string // password that stored in database
	Salt           string
	IsBlocked      bool
	IsRegistered   bool
}

func NewLoginUserInput(number, password, storedPassword, salt string, isBlocked, isRegistered bool) LoginUserInput {
	return LoginUserInput{
		LoginUserInput: shared_dto.LoginUserInput{
			PhoneNumber:   number,
			InputPassword: password,
		},
		StoredPassword: storedPassword,
		Salt:           salt,
		IsBlocked:      isBlocked,
		IsRegistered:   isRegistered,
	}
}

type RefreshInput struct {
	shared_dto.RefreshInput
}

func NewRefreshInput(refresh string) RefreshInput {
	return RefreshInput{
		shared_dto.RefreshInput{
			Refresh: refresh,
		},
	}
}

type AvatarChooseInput struct {
	shared_dto.AvatarChooseInput
}

func NewAvatarChooseInput(avatar string) AvatarChooseInput {
	return AvatarChooseInput{
		shared_dto.AvatarChooseInput{
			Avatar: avatar,
		},
	}
}

type ResetPasswordInput struct {
	shared_dto.ResetPasswordInput
}

func NewResetPasswordInput(number, password, token string) ResetPasswordInput {
	return ResetPasswordInput{
		shared_dto.ResetPasswordInput{
			PhoneNumber: number,
			Password:    password,
			Token:       token,
		},
	}
}
