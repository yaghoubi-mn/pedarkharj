package domain_user

import (
	"errors"
	"time"

	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/utils"
)

type UserDomainService interface {
	Signup(user *User, token string) error
	VerifyNumber(number string, code uint, token string, isBlocked bool) error
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
func (s *service) VerifyNumber(number string, code uint, token string, isBlocked bool) error {

	if isBlocked {
		return errors.New("you are blocked")
	}

	err := s.validator.ValidateField(number, "e164,required")
	if err != nil {
		return err
	}

	if code > 99999 || code < 10000 {
		return errors.New("code: invalid code")
	}

	err = s.validator.ValidateField(token, "uuid,required")
	if err != nil {
		return err
	}

	return nil

}

func (s *service) Signup(user *User, token string) error {
	err := s.validator.ValidateField(user.Name, "name,required")
	if err != nil {
		return errors.New("name: " + err.Error())
	}

	err = s.validator.ValidateField(user.Number, "e164,required")
	if err != nil {
		return errors.New("number: " + err.Error())
	}

	if len(user.Password) < 8 {
		return errors.New("password: small password")
	}

	user.RegisteredAt = time.Now()
	user.IsRegistered = true

	user.Salt, err = utils.GenerateRandomSalt()
	if err != nil {
		return err
	}

	user.Password, err = utils.HashPasswordWithSalt(user.Password, user.Salt)
	if err != nil {
		return err
	}

	return nil
}
