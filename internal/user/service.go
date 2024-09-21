package user

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/rcodes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/utils"
)

type UserService interface {
	Signup(signupUserInput SignupUserInput) (tokens map[string]string, code rcodes.ResponseCode, errMap map[string]string, err error)
	VerifyNumber(v VerifyNumberInput) (step int, code rcodes.ResponseCode, token map[string]string, errMap map[string]string, err error)
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
func (s *service) VerifyNumber(verifyNumberInput VerifyNumberInput) (int, rcodes.ResponseCode, map[string]string, map[string]string, error) {

	errMap := s.validator.Struct(verifyNumberInput)
	tokens := make(map[string]string)
	if errMap != nil {
		return 0, rcodes.InvalidField, tokens, errMap, nil
	}

	errMap = make(map[string]string)

	if verifyNumberInput.Code == 0 {
		// step 1: sent otp code to number

		// generate random otp code between 10000 and 99999
		otp := rand.Intn(90000) + 10000
		token := uuid.New()
		// TODO: send code to number
		otpString := strconv.Itoa(otp)

		err := s.cacheRepo.Save(verifyNumberInput.Number+token.String(), otpString, 10*time.Minute)
		if err != nil {
			return 0, "", tokens, nil, err
		}
		fmt.Println("OTP: ", otp)

		tokens["token"] = token.String()
		return 1, rcodes.CodeSendToNumber, tokens, nil, nil
	} else {
		// step 2: check otp code
		otp, err := s.cacheRepo.Get(verifyNumberInput.Number + verifyNumberInput.Token)

		if err != nil {
			if err == database_errors.ErrRecordNotFound || err == database_errors.ErrExpired {
				errMap["non-field"] = "zero code first"
				return 0, rcodes.ZeroCodeFirst, tokens, errMap, nil
			} else if err == database_errors.ErrExpired {
				errMap["code"] = "code expired"
				return 0, rcodes.OTPExpired, tokens, errMap, nil
			}

			return 0, "", tokens, nil, err
		}

		if otp == strconv.Itoa(int(verifyNumberInput.Code)) {
			user, err := s.repo.GetByNumber(verifyNumberInput.Number)
			_ = user
			if err != nil {
				if err == database_errors.ErrRecordNotFound {

					// save
					err := s.cacheRepo.Save(verifyNumberInput.Number+verifyNumberInput.Token, "signup", 10*time.Minute)
					if err != nil {
						return 0, "", tokens, nil, err
					}

					return 2, rcodes.GoSignup, tokens, nil, nil
				}

				return 0, "", tokens, nil, err
			}

			// user found in database
			//TODO: get refresh token
			return 3, "", tokens, nil, nil

		} else {
			errMap["code"] = "wrong code"
			return 0, rcodes.WrongCode, tokens, errMap, nil
		}

	}
}

func (s *service) Signup(userInput SignupUserInput) (map[string]string, rcodes.ResponseCode, map[string]string, error) {
	errMap := s.validator.Struct(userInput)
	tokens := make(map[string]string)

	if errMap != nil {
		return tokens, rcodes.InvalidField, errMap, nil
	}

	errMap = make(map[string]string)

	verify, err := s.cacheRepo.Get(userInput.Number + userInput.Token)
	if err != nil {
		if err == database_errors.ErrRecordNotFound {
			errMap["non-field"] = "verify number first"
			return tokens, rcodes.VerifyNumberFirst, errMap, nil
		} else {
			return tokens, "", nil, err
		}
	}

	if verify != "signup" {
		errMap["non-field"] = "verify number first"
		return tokens, rcodes.VerifyNumberFirst, errMap, nil
	}

	// delete cache
	if err := s.cacheRepo.Delete(userInput.Number + userInput.Token); err != nil {
		return tokens, "", nil, err
	}

	var user User

	user.Name = userInput.Name
	user.Number = userInput.Number

	user.IsRegistered = true
	user.LastLogin = time.Now()
	user.RegisteredAt = time.Now()

	user.Salt, err = utils.GenerateRandomSalt()
	if err != nil {
		return tokens, "", nil, err
	}

	user.Password, err = utils.HashPasswordWithSalt(userInput.Password, user.Salt)
	if err != nil {
		return tokens, "", nil, err
	}

	// TODO: create device

	// TODO: get refresh and access

	err = s.repo.Create(user)
	return tokens, "", nil, err
}
