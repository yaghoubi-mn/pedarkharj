package domain_user

import (
	"time"

	"github.com/yaghoubi-mn/pedarkharj/internal/infrastructure/config"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/service_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/utils"
)

type UserDomainService interface {
	Signup(user *User, token string) (userError error, serverError error)
	VerifyNumber(number string, code uint, token string, isBlocked bool) (userError error, serverError error)
	CheckNumber(number string) error
	Login(number, inputPassword, realPassword, salt string, isBlocked bool) (userError, serverError error)
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
		return service_errors.ErrBlockedUser, nil
	}

	if err := s.validator.ValidateFieldByFieldName("Number", number, User{}); err != nil {
		return service_errors.ErrInvalidNumber, nil
	}

	if (code > 99999 || code < 10000) && code != 0 {
		return service_errors.ErrInvalidCode, nil
	}

	if err := s.validator.ValidateField(token, "uuid,allowempty"); err != nil {
		return service_errors.ErrInvalidToken, nil
	}

	return nil, nil

}

func (s *service) Signup(user *User, token string) (error, error) {
	if err := s.validator.ValidateFieldByFieldName("Name", user.Name, User{}); err != nil {
		return service_errors.ErrInvalidName, nil
	}

	if err := s.validator.ValidateFieldByFieldName("Number", user.Number, User{}); err != nil {
		return service_errors.ErrInvalidNumber, nil
	}

	if err := s.validator.ValidateField(token, "uuid,required"); err != nil {
		return service_errors.ErrInvalidToken, nil
	}

	if len(user.Name) < 2 {
		return service_errors.ErrSmallName, nil
	}

	if len(user.Password) < 8 {
		return service_errors.ErrSmallPassword, nil
	}

	user.RegisteredAt = time.Now()
	user.IsRegistered = true

	var err error
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

	if err := s.validator.ValidateFieldByFieldName("Number", number, User{}); err != nil {
		return service_errors.ErrInvalidNumber
	}

	return nil
}

// realPassword is hashed (stored password in database)
func (s *service) Login(number, inputPassword, hashedRealPassword, salt string, isBlocked bool) (error, error) {

	if isBlocked {
		return service_errors.ErrBlockedUser, nil
	}

	if err := utils.CompareHashAndPassword(hashedRealPassword, inputPassword, salt); err != nil {
		return service_errors.ErrWrongPassword, nil
	}

	if err := s.validator.ValidateFieldByFieldName("Number", number, User{}); err != nil {
		return service_errors.ErrInvalidNumber, nil
	}

	return nil, nil
}
