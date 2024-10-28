package domain_user

import (
	"time"

	"github.com/yaghoubi-mn/pedarkharj/internal/infrastructure/config"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/service_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/utils"
)

type UserDomainService interface {
	Signup(input SignupUserInput) (user User, userError error, serverError error)
	VerifyNumber(input VerifyNumberInput) (userError error, serverError error)
	CheckNumber(number string) error
	Login(input LoginUserInput) (userError, serverError error)
	ResetPassword(input RestPasswordInput) (userErr, serverErr error, salt, outPassword string)
}

type service struct {
	validator datatypes.Validator
}

func NewUserService(validator datatypes.Validator) UserDomainService {
	return &service{
		validator: validator,
	}
}

// step: int, code: string, token: string, errors: []error, err: error
func (s *service) VerifyNumber(input VerifyNumberInput) (error, error) {

	// if input.IsBlocked {
	// 	return service_errors.ErrBlockedUser, nil
	// }

	if err := s.validator.ValidateFieldByFieldName("Number", input.Number, User{}); err != nil {
		return service_errors.ErrInvalidNumber, nil
	}

	if (input.OTP > 99999 || input.OTP < 10000) && input.OTP != 0 {
		return service_errors.ErrInvalidCode, nil
	}

	if err := s.validator.ValidateField(input.Token, "uuid,allowempty"); err != nil {
		return service_errors.ErrInvalidToken, nil
	}

	if input.Mode != "signup" && input.Mode != "reset_password" {
		return service_errors.ErrInvalidMode, nil
	}

	return nil, nil

}

func (s *service) Signup(input SignupUserInput) (User, error, error) {
	var user User
	if err := s.validator.ValidateFieldByFieldName("Name", input.Name, User{}); err != nil {
		return user, service_errors.ErrInvalidName, nil
	}

	if err := s.validator.ValidateFieldByFieldName("Number", input.Number, User{}); err != nil {
		return user, service_errors.ErrInvalidNumber, nil
	}

	if err := s.validator.ValidateField(input.Token, "uuid,required"); err != nil {
		return user, service_errors.ErrInvalidToken, nil
	}

	if len(input.Name) < 2 {
		return user, service_errors.ErrSmallName, nil
	}

	if len(input.Password) < 8 {
		return user, service_errors.ErrSmallPassword, nil
	}

	if len(input.Password) > 30 {
		return user, service_errors.ErrLongPassword, nil
	}

	user = input.GetUser()
	user.RegisteredAt = time.Now()
	user.IsRegistered = true

	var err error
	user.Salt, err = utils.GenerateRandomSalt()
	if err != nil {
		return user, nil, err
	}

	user.Password, err = utils.HashPasswordWithSalt(user.Password, user.Salt, config.BcryptCost)
	if err != nil {
		return user, nil, err
	}

	return user, nil, nil
}

func (s *service) ResetPassword(input RestPasswordInput) (userErr, serverErr error, salt, outPassword string) {

	if err := s.validator.ValidateFieldByFieldName("Number", input.Number, User{}); err != nil {
		return service_errors.ErrInvalidNumber, nil, "", ""
	}

	// if err := s.validator.ValidateFieldByFieldName("Password", password, User{}); err!=nil{
	// 	return service_errors.
	// }

	if len(input.Password) < 8 {
		return service_errors.ErrSmallPassword, nil, "", ""
	}

	if len(input.Password) > 30 {
		return service_errors.ErrLongPassword, nil, "", ""
	}

	salt, err := utils.GenerateRandomSalt()
	if err != nil {
		return nil, err, "", ""
	}

	outPassword, err = utils.HashPasswordWithSalt(input.Password, salt, config.BcryptCost)
	if err != nil {
		return nil, err, "", ""
	}

	return nil, nil, salt, outPassword
}

func (s *service) CheckNumber(number string) error {

	if err := s.validator.ValidateFieldByFieldName("Number", number, User{}); err != nil {
		return service_errors.ErrInvalidNumber
	}

	return nil
}

// realPassword is hashed (stored password in database)
func (s *service) Login(input LoginUserInput) (error, error) {

	if input.IsBlocked {
		return service_errors.ErrBlockedUser, nil
	}

	if err := utils.CompareHashAndPassword(input.RealPassword, input.InputPassword, input.Salt); err != nil {
		return service_errors.ErrWrongPassword, nil
	}

	if err := s.validator.ValidateFieldByFieldName("Number", input.Number, User{}); err != nil {
		return service_errors.ErrInvalidNumber, nil
	}

	if !input.IsRegistered {
		return service_errors.ErrUserNotRegistered, nil
	}

	return nil, nil
}
