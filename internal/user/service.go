package user

import "github.com/go-playground/validator/v10"

type UserService interface {
	CreateUser(user User) error
}

type service struct {
	repo      UserRepository
	validator validator.Validate
}

func NewUserService(repo UserRepository, validator validator.Validate) UserService {
	return &service{
		repo:      repo,
		validator: validator,
	}
}

func (s *service) CreateUser(user User) error {
	err := s.validator.Struct(user)
	if err != nil {
		return err
	}

	return s.repo.Create(user)
}
