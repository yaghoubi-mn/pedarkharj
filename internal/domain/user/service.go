package domain_user

import (
	"errors"
	"time"

	"github.com/yaghoubi-mn/pedarkharj/internal/infrastructure/config"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/utils"
)

type UserDomainService interface {
	Signup(user *User, token string) (userError error, serverError error)
	VerifyNumber(number string, code uint, token string, isBlocked bool) (userError error, serverError error)
	CheckNumber(number string) error
	Login(number, inputPassword, realPassword, salt string) (userError, serverError error)
}

type DeviceRepository interface {
	CreateWithParam(name string, lastIP string, firstLogin time.Time, lastLogin time.Time, refreshToken string) error
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
func (s *service) VerifyNumber(number string, code uint, token string, isBlocked bool) (error, error) {

	if isBlocked {
		return errors.New("you are blocked"), nil
	}

	err := s.validator.ValidateField(number, "e164,required")
	if err != nil {
		return errors.New("number: invalid number"), nil
	}

	if (code > 99999 || code < 10000) && code != 0 {
		return errors.New("code: invalid code"), nil
	}

	err = s.validator.ValidateField(token, "uuid,allowempty")
	if err != nil {
		return errors.New("token: invalid token"), nil
	}

	return nil, nil

}

func (s *service) Signup(user *User, token string) (error, error) {
	err := s.validator.ValidateField(user.Name, "name,required")
	if err != nil {
		return errors.New("name: " + err.Error()), nil
	}

	err = s.validator.ValidateField(user.Number, "e164,required")
	if err != nil {
		return errors.New("number: " + err.Error()), nil
	}

	if len(user.Password) < 8 {
		return errors.New("password: small password"), nil
	}

	user.RegisteredAt = time.Now()
	user.IsRegistered = true

	user.Salt, err = utils.GenerateRandomSalt()
	if err != nil {
		return nil, err
	}

	user.Password, err = utils.HashPasswordWithSalt(user.Password, user.Salt, config.BcryptCost)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *service) CheckNumber(number string) error {

	err := s.validator.ValidateField(number, "e164,required")
	return err
}

// realPassword is hashed (stored password in database)
func (s *service) Login(number, inputPassword, hashedRealPassword, salt string) (error, error) {

	err := utils.CompareHashAndPassword(hashedRealPassword, inputPassword, salt)
	if err != nil {
		return errors.New("invalid number or password"), nil
	}

	err = s.validator.ValidateField(number, "e164,required")
	if err != nil {
		return errors.New("number: invalid number"), nil
	}

	return nil, nil
}
