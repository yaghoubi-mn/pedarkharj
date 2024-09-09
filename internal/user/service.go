package user

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/rcodes"
)

type UserService interface {
	Signup(signupUserInput SignupUserInput) (code string, errMap map[string]string, err error)
	VerifyNumber(v VerifyNumberInput) (step int, code string, token string, errMap map[string]string, err error)
}

type service struct {
	repo      UserRepository
	cacheRepo datatypes.CacheRepository
	validator datatypes.Validator
}

func NewUserService(repo UserRepository, cacheRepo datatypes.CacheRepository, validator datatypes.Validator) UserService {
	return &service{
		repo:      repo,
		cacheRepo: cacheRepo,
		validator: validator,
	}
}

// step: int, code: string, token: string, errors: []error, err: error
func (s *service) VerifyNumber(verifyNumberInput VerifyNumberInput) (int, string, string, map[string]string, error) {

	errMap := s.validator.Struct(verifyNumberInput)
	if errMap != nil {
		return 0, rcodes.INVALID_FIELD, "", errMap, nil
	}

	errMap = make(map[string]string)

	if verifyNumberInput.Code == 0 {
		// step 1: sent otp code to number

		// generate random otp code between 10000 and 99999
		otp := rand.Intn(90000) + 10000
		token := uuid.New()
		// TODO: send code to number
		err := s.cacheRepo.Save(verifyNumberInput.Number+token.String(), otp, 10*time.Minute)
		if err != nil {
			return 0, "", "", nil, err
		}
		fmt.Println("OTP: ", otp)

		return 1, rcodes.CODE_SENT_TO_NUMBER, token.String(), nil, nil
	} else {
		// step 2: check otp code
		otp, err := s.cacheRepo.Get(verifyNumberInput.Number + verifyNumberInput.Token)
		if err != nil {
			if err == database_errors.ErrRecordNotFound {
				errMap["code"] = "code expired or invalid token"
				return 0, rcodes.ZERO_CODE_FIRST, "", errMap, nil
			}

			return 0, "", "", nil, err
		}

		if otp == verifyNumberInput.Code {
			// save
			err := s.cacheRepo.Save(verifyNumberInput.Number+verifyNumberInput.Token, "signup", 10*time.Minute)
			if err != nil {
				return 0, "", "", nil, err
			}

			return 0, rcodes.GO_SIGNUP, "", nil, nil

		} else {
			errMap["code"] = "wrong code"
			return 0, rcodes.WRONG_CODE, "", errMap, nil
		}

	}
}

func (s *service) Signup(userInput SignupUserInput) (string, map[string]string, error) {
	errMap := s.validator.Struct(userInput)
	if errMap != nil {
		return "", errMap, nil
	}

	errMap = make(map[string]string)

	verify, err := s.cacheRepo.Get(userInput.Number + userInput.Token)
	if err != nil {
		if err == database_errors.ErrRecordNotFound {
			errMap["non-field"] = "verify number first"
			return rcodes.VERIFY_NUMBER_FIRST, errMap, nil
		} else {
			return "", nil, err
		}
	}

	if verify != "signup" {
		errMap["non-field"] = "verify number first"
		return rcodes.VERIFY_NUMBER_FIRST, errMap, nil
	}

	var user User

	user.IsRegistered = true
	user.Name = userInput.Name
	user.Number = userInput.Number

	// TODO: create device

	err = s.repo.Create(user)
	return "", nil, err
}
