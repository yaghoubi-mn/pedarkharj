package user

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/jwt"
	"github.com/yaghoubi-mn/pedarkharj/pkg/rcodes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/utils"
)

var (
	JwtRefreshExpireMinutes time.Duration = 30 * 24 * 60
	JwtAccessExpireMinutes  time.Duration = 15
)

type UserService interface {
	Signup(signupUserInput SignupUserInput) (tokens map[string]string, code rcodes.ResponseCode, errMap map[string]string, err error)
	VerifyNumber(v VerifyNumberInput) (step int, code rcodes.ResponseCode, token map[string]string, errMap map[string]string, err error)
}

type DeviceRepository interface {
	CreateWithParam(name string, lastIP string, firstLogin time.Time, lastLogin time.Time, refreshToken string) error
}

type service struct {
	repo       UserRepository
	cacheRepo  datatypes.CacheRepository
	validator  datatypes.Validator
	deviceRepo DeviceRepository
}

func NewUserService(repo UserRepository, cacheRepo datatypes.CacheRepository, validator datatypes.Validator) UserService {
	return &service{
		repo:      repo,
		cacheRepo: cacheRepo,
		validator: validator,
		// deviceRepo: deviceRepo,
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

		// check for number delay
		_, err := s.cacheRepo.Get(verifyNumberInput.Number)

		if err != nil {
			if err == database_errors.ErrExpired || err == database_errors.ErrRecordNotFound {

				// generate random otp code between 10000 and 99999
				otp := rand.Intn(90000) + 10000
				token := uuid.New()
				// TODO: send code to number

				verifyInfo := make(map[string]string)
				verifyInfo["token"] = token.String()
				verifyInfo["otp"] = strconv.Itoa(otp)

				// convert map to bytes
				verifyInfoJsonBytes, err := json.Marshal(verifyInfo)
				if err != nil {
					return 0, "", nil, nil, err
				}

				// convert bytes to string
				verifyInfoString := hex.EncodeToString(verifyInfoJsonBytes)

				err = s.cacheRepo.Save(verifyNumberInput.Number, verifyInfoString, 10*time.Minute)
				if err != nil {
					return 0, "", nil, nil, err
				}
				fmt.Println("OTP: ", otp)

				tokens["token"] = token.String()
				return 1, rcodes.CodeSendToNumber, tokens, nil, nil

			} else {
				return 0, "", nil, nil, err
			}
		}

		// otp code sent and not expired
		errMap["number"] = "wait some seconds"
		return 1, rcodes.NumberDelay, nil, errMap, err
	} else {
		// step 2: check otp code
		verifyInfoString, err := s.cacheRepo.Get(verifyNumberInput.Number)

		if err != nil {
			if err == database_errors.ErrRecordNotFound || err == database_errors.ErrExpired {
				errMap["non-field"] = "zero code first"
				return 0, rcodes.ZeroCodeFirst, tokens, errMap, nil
			} else if err == database_errors.ErrExpired {
				errMap["code"] = "code expired"
				return 0, rcodes.OTPExpired, tokens, errMap, nil
			}

			return 0, "", nil, nil, err
		}

		// convert verfiyInfo to map
		verifyInfoBytes, err := hex.DecodeString(verifyInfoString)
		if err != nil {
			return 0, "", nil, nil, err
		}

		verifyInfo := make(map[string]string)

		err = json.Unmarshal(verifyInfoBytes, &verifyInfo)
		if err != nil {
			return 0, "", nil, nil, err
		}

		// get token from verifyInfo
		token, ok := verifyInfo["token"]
		if !ok {
			return 0, "", nil, nil, errors.New("cannot convert token to string")
		}

		// check token
		if token != verifyNumberInput.Token {
			errMap["token"] = "invalid token"
			return 2, rcodes.InvalidField, nil, errMap, nil
		}

		// get otp
		otp, ok := verifyInfo["otp"]
		if !ok {
			return 0, "", nil, nil, errors.New("cannot convert otp to int")
		}

		// check otp
		if otp == strconv.Itoa(int(verifyNumberInput.Code)) {
			user, err := s.repo.GetByNumber(verifyNumberInput.Number)

			if err != nil {
				if err == database_errors.ErrRecordNotFound {

					verifyInfo := make(map[string]string)
					verifyInfo["token"] = token
					verifyInfo["verify"] = "signup"

					// convert map to bytes
					verifyInfoJsonBytes, err := json.Marshal(verifyInfo)
					if err != nil {
						return 0, "", nil, nil, err
					}

					// convert bytes to string
					verifyInfoString := hex.EncodeToString(verifyInfoJsonBytes)

					// save number and token to cache for signup
					err = s.cacheRepo.Save(verifyNumberInput.Number, verifyInfoString, 10*time.Minute)
					if err != nil {
						return 0, "", tokens, nil, err
					}

					return 2, rcodes.GoSignup, tokens, nil, nil
				}

				return 0, "", tokens, nil, err
			}

			// user found in database
			refresh, access, err := createTokensAndUserDevice(user)
			if err != nil {
				return 0, "", nil, nil, err
			}

			tokens["refresh"] = refresh
			tokens["access"] = access
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

	verifyInfoString, err := s.cacheRepo.Get(userInput.Number)
	if err != nil {
		if err == database_errors.ErrRecordNotFound || err == database_errors.ErrExpired {
			errMap["non-field"] = "verify number first"
			return nil, rcodes.VerifyNumberFirst, errMap, nil
		} else {
			return nil, "", nil, err
		}
	}

	// convert verfiyInfo to map
	verifyInfoBytes, err := hex.DecodeString(verifyInfoString)
	if err != nil {
		return nil, "", nil, err
	}

	verifyInfo := make(map[string]string)

	err = json.Unmarshal(verifyInfoBytes, &verifyInfo)
	if err != nil {
		return nil, "", nil, err
	}

	// check token
	if verifyInfo["token"] != userInput.Token {
		errMap["token"] = "invalid token"
		return nil, rcodes.InvalidField, errMap, nil
	}

	verify, ok := verifyInfo["verify"]
	if !ok || verify != "signup" {
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

	refresh, access, err := createTokensAndUserDevice(user)
	if err != nil {
		return tokens, "", nil, err
	}

	tokens["refresh"] = refresh
	tokens["access"] = access

	err = s.repo.Create(user)
	return tokens, "", nil, err
}

func createTokensAndUserDevice(user User) (refresh string, access string, err error) {
	refresh, access, err = createRefreshAndAccessFromUser(user)

	// create device
	// deviceRepo.CreateWithParam()

	return refresh, access, err
}

func createRefreshAndAccessFromUser(user User) (refresh string, access string, err error) {
	refresh, err = jwt.CreateJwt(map[string]any{
		"exp": JwtRefreshExpireMinutes,
	})

	if err != nil {
		return "", "", err
	}

	access, err = jwt.CreateJwt(map[string]any{
		"exp":          time.Now().Add(JwtAccessExpireMinutes * time.Minute),
		"name":         user.Name,
		"number":       user.Number,
		"isRegistered": user.IsRegistered,
	})

	return refresh, access, err

}

func getUserFromAccess(access string) (User, error) {
	var user User
	mapClaims, err := jwt.VerifyJwt(access)
	if err != nil {
		return user, err
	}

	user.Name = mapClaims["name"].(string)
	user.Number = mapClaims["number"].(string)
	user.IsRegistered = mapClaims["isRegistered"].(bool)

	return user, nil
}
