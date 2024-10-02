package app_user

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	app_device "github.com/yaghoubi-mn/pedarkharj/internal/application/device"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	"github.com/yaghoubi-mn/pedarkharj/internal/infrastructure/config"
	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/jwt"
	"github.com/yaghoubi-mn/pedarkharj/pkg/rcodes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/utils"
)

type UserAppService interface {
	VerifyNumber(verifyNumberInput VerifyNumberInput, deviceName string, deviceIP string) (step int, responseCode rcodes.ResponseCode, tokens map[string]string, userError error, serverError error)
	Signup(userInput SignupUserInput, deviceName string, deviceIP string) (tokens map[string]string, responseCode rcodes.ResponseCode, userError error, serverError error)
	GetUserInfo(user domain_user.User) UserOutput
	CheckNumber(numberInput NumberInput) (isExist bool, err error)
	Login(loginInput LoginUserInput, deviceName string, deviceIP string) (tokens map[string]string, responseCode rcodes.ResponseCode, userError error, serverError error)
}

type deviceRepository interface {
	CreateWithParam(name string, lastIP string, firstLogin time.Time, lastLogin time.Time, refreshToken string) error
}

type service struct {
	repo             domain_user.UserDomainRepository
	cacheRepo        datatypes.CacheRepository
	domainService    domain_user.UserDomainService
	deviceAppService app_device.DeviceAppService
}

func NewUserService(repo domain_user.UserDomainRepository, cacheRepo datatypes.CacheRepository, deviceAppService app_device.DeviceAppService, domainService domain_user.UserDomainService) UserAppService {
	return &service{
		repo:             repo,
		cacheRepo:        cacheRepo,
		domainService:    domainService,
		deviceAppService: deviceAppService,
	}
}

func (s *service) VerifyNumber(verifyNumberInput VerifyNumberInput, deviceName string, deviceIP string) (int, rcodes.ResponseCode, map[string]string, error, error) {

	// isBlocked will checked in step 2
	tokens := make(map[string]string)

	userErr, serverErr := s.domainService.VerifyNumber(verifyNumberInput.Number, verifyNumberInput.Code, verifyNumberInput.Token, false)
	if serverErr != nil {
		return 0, "", nil, nil, serverErr
	}
	if userErr != nil {
		return 0, rcodes.InvalidField, tokens, userErr, nil
	}

	if verifyNumberInput.Code == 0 {
		// step 1: sent otp code to number

		// check for number delay
		_, err := s.cacheRepo.Get(verifyNumberInput.Number)

		if config.Debug {
			err = database_errors.ErrExpired
		}
		if err != nil {
			if err == database_errors.ErrExpired || err == database_errors.ErrRecordNotFound {

				// generate random otp code between 10000 and 99999
				otp := rand.Intn(90000) + 10000
				token := uuid.New()
				// TODO: send code to number

				verifyInfo := make(map[string]string)
				verifyInfo["token"] = token.String()
				verifyInfo["otp"] = strconv.Itoa(otp)

				verifyInfoString, err := utils.ConvertMapToString(verifyInfo)
				if err != nil {
					return 0, "", nil, nil, err
				}

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
		return 1, rcodes.NumberDelay, nil, errors.New("number: otp not expired"), err

	} else {
		// step 2: check otp code
		verifyInfoString, err := s.cacheRepo.Get(verifyNumberInput.Number)

		if err != nil {
			if err == database_errors.ErrRecordNotFound || err == database_errors.ErrExpired {

				return 0, rcodes.ZeroCodeFirst, tokens, errors.New("code: zero code in first step"), nil

			}
			// else if err == database_errors.ErrExpired {
			// 	errMap["code"] = "code expired"
			// 	return 0, rcodes.OTPExpired, tokens, errMap, nil
			// }

			return 0, "", nil, nil, err
		}

		verifyInfo, err := utils.ConvertStringToMap(verifyInfoString)
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

			return 2, rcodes.InvalidField, nil, errors.New("token: invalid token"), nil
		}

		// get otp
		otp, ok := verifyInfo["otp"]
		if !ok {
			return 0, "", nil, nil, errors.New("cannot convert otp to int")
		}

		// check otp
		if otp == strconv.Itoa(int(verifyNumberInput.Code)) {
			// get user
			user, err := s.repo.GetByNumber(verifyNumberInput.Number)

			otpInt, err2 := strconv.Atoi(otp)
			if err2 != nil {
				return 0, "", nil, nil, err
			}

			// call domain service
			userErr, serverErr = s.domainService.VerifyNumber(user.Number, uint(otpInt), token, user.IsBlocked)
			if serverErr != nil {
				return 0, "", nil, nil, serverErr
			}
			if userErr != nil {
				return 0, "", nil, userErr, nil
			}

			if err != nil {
				if err == database_errors.ErrRecordNotFound {

					// user not exist. redirect to signup

					verifyInfo := make(map[string]string)
					verifyInfo["token"] = token
					verifyInfo["verify"] = "signup"

					verifyInfoString, err = utils.ConvertMapToString(verifyInfo)
					if err != nil {
						return 0, "", nil, nil, err
					}

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
			tokens, err := jwt.CreateRefreshAndAccessFromUserWithMap(config.JWtRefreshExpireMinutes, config.JWTAccessExpireMinutes, user.ID, user.Name, user.Number, user.IsRegistered)
			if err != nil {
				return 0, "", nil, nil, err
			}

			// create device

			err = s.deviceAppService.CreateOrUpdate(app_device.DeviceInput{
				Name:         deviceName,
				IP:           deviceIP,
				RefreshToken: tokens["refresh"],
				UserID:       user.ID,
			})
			if err != nil {
				return 0, "", nil, nil, err
			}

			return 3, "", tokens, nil, nil

		} else {
			return 0, rcodes.WrongCode, tokens, errors.New("code: wrong code"), nil
		}

	}
}

func (s *service) Signup(userInput SignupUserInput, deviceName string, deviceIP string) (map[string]string, rcodes.ResponseCode, error, error) {
	var user domain_user.User
	user.Number = userInput.Number
	user.Name = userInput.Name
	user.Password = userInput.Password

	userErr, serverErr := s.domainService.Signup(&user, userInput.Token)
	if serverErr != nil {
		return nil, "", nil, serverErr
	}
	if userErr != nil {
		return nil, rcodes.InvalidField, userErr, nil
	}

	verifyInfoString, err := s.cacheRepo.Get(userInput.Number)
	if err != nil {
		if err == database_errors.ErrRecordNotFound || err == database_errors.ErrExpired {

			return nil, rcodes.VerifyNumberFirst, errors.New("verify number first"), nil

		} else {
			return nil, "", nil, err
		}
	}

	verifyInfo, err := utils.ConvertStringToMap(verifyInfoString)
	if err != nil {
		return nil, "", nil, err
	}

	// check token
	if verifyInfo["token"] != userInput.Token {

		return nil, rcodes.InvalidField, errors.New("token: invalid token"), nil
	}

	verify, ok := verifyInfo["verify"]
	if !ok || verify != "signup" {

		return nil, rcodes.VerifyNumberFirst, errors.New("verify number first"), nil
	}

	// delete cache
	if err := s.cacheRepo.Delete(userInput.Number + userInput.Token); err != nil {
		return nil, "", nil, err
	}

	tokens, err := jwt.CreateRefreshAndAccessFromUserWithMap(config.JWtRefreshExpireMinutes, config.JWTAccessExpireMinutes, user.ID, user.Name, user.Number, user.IsRegistered)
	if err != nil {
		return nil, "", nil, err
	}

	// create device
	err = s.deviceAppService.CreateOrUpdate(app_device.DeviceInput{
		UserID:       user.ID,
		Name:         deviceName,
		IP:           deviceIP,
		RefreshToken: tokens["refresh"],
	})
	if err != nil {
		return nil, "", nil, err
	}

	err = s.repo.Create(user)
	return tokens, "", nil, err
}

func (s *service) GetUserInfo(user domain_user.User) UserOutput {

	var userOutput UserOutput
	userOutput.Fill(user)

	return userOutput
}

func (s *service) CheckNumber(numberInput NumberInput) (bool, error) {

	err := s.domainService.CheckNumber(numberInput.Number)
	if err != nil {
		return false, err
	}

	_, err = s.repo.GetByNumber(numberInput.Number)

	if err != nil {
		if err == database_errors.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	} else {
		return true, nil
	}

}

func (s *service) Login(loginInput LoginUserInput, deviceName string, deviceIP string) (map[string]string, rcodes.ResponseCode, error, error) {

	user, err := s.repo.GetByNumber(loginInput.Number)
	if err != nil {
		return nil, "", nil, err
	}

	userErr, serverErr := s.domainService.Login(user.Number, loginInput.Password, user.Password, user.Salt)
	if serverErr != nil {
		return nil, "", nil, serverErr
	}
	if userErr != nil {
		return nil, rcodes.InvalidField, userErr, nil
	}

	tokens, err := jwt.CreateRefreshAndAccessFromUserWithMap(config.JWtRefreshExpireMinutes, config.JWTAccessExpireMinutes, user.ID, user.Name, user.Number, user.IsRegistered)
	if err != nil {
		return nil, "", nil, err
	}

	// create device
	err = s.deviceAppService.CreateOrUpdate(app_device.DeviceInput{
		UserID:       user.ID,
		Name:         deviceName,
		IP:           deviceIP,
		RefreshToken: tokens["refresh"],
	})
	if err != nil {
		return nil, "", nil, err
	}

	return tokens, "", nil, nil
}
