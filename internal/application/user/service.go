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
	VerifyNumber(verifyNumberInput VerifyNumberInput, deviceName string, deviceIP string) (step int, responseDTO datatypes.ResponseDTO)
	Signup(userInput SignupUserInput, deviceName string, deviceIP string) (responseDTO datatypes.ResponseDTO)
	GetUserInfo(user domain_user.User) datatypes.ResponseDTO
	CheckNumber(numberInput NumberInput) datatypes.ResponseDTO
	Login(loginInput LoginUserInput, deviceName string, deviceIP string) (responseDTO datatypes.ResponseDTO)
	GetAccessFromRefresh(refresh string) (responseDTO datatypes.ResponseDTO)
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

func (s *service) VerifyNumber(verifyNumberInput VerifyNumberInput, deviceName string, deviceIP string) (int, datatypes.ResponseDTO) {

	var responseDTO datatypes.ResponseDTO
	responseDTO.Data = make(map[string]any)

	// isBlocked will checked in step 2

	userErr, serverErr := s.domainService.VerifyNumber(verifyNumberInput.Number, verifyNumberInput.Code, verifyNumberInput.Token, false)
	if serverErr != nil {
		responseDTO.ServerErr = serverErr
		return 0, responseDTO
	}
	if userErr != nil {
		responseDTO.UserErr = userErr
		responseDTO.ResponseCode = rcodes.InvalidField
		return 0, responseDTO
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
					responseDTO.ServerErr = err
					return 0, responseDTO
				}

				err = s.cacheRepo.Save(verifyNumberInput.Number, verifyInfoString, 10*time.Minute)
				if err != nil {
					responseDTO.ServerErr = err
					return 0, responseDTO
				}
				fmt.Println("OTP: ", otp)

				responseDTO.ResponseCode = rcodes.CodeSendToNumber
				responseDTO.Data["token"] = token.String()
				return 1, responseDTO

			} else {
				responseDTO.ServerErr = err
				return 0, responseDTO
			}
		}

		// otp code sent and not expired

		responseDTO.ResponseCode = rcodes.NumberDelay
		responseDTO.UserErr = errors.New("number: otp not expired. wait some minutes")
		responseDTO.ServerErr = err
		return 1, responseDTO

	} else {
		// step 2: check otp code
		verifyInfoString, err := s.cacheRepo.Get(verifyNumberInput.Number)

		if err != nil {
			if err == database_errors.ErrRecordNotFound || err == database_errors.ErrExpired {

				responseDTO.ResponseCode = rcodes.ZeroCodeFirst
				responseDTO.UserErr = errors.New("code: zero code in first step")
				return 0, responseDTO

			}
			// else if err == database_errors.ErrExpired {
			// 	errMap["code"] = "code expired"
			// 	return 0, rcodes.OTPExpired, tokens, errMap, nil
			// }

			responseDTO.ServerErr = err
			return 0, responseDTO
		}

		verifyInfo, err := utils.ConvertStringToMap(verifyInfoString)
		if err != nil {
			responseDTO.ServerErr = err
			return 0, responseDTO
		}

		// get token from verifyInfo
		token, ok := verifyInfo["token"]
		if !ok {
			responseDTO.ServerErr = errors.New("cannot convert token to string")
			return 0, responseDTO
		}

		// check token
		if token != verifyNumberInput.Token {

			responseDTO.ResponseCode = rcodes.InvalidField
			responseDTO.UserErr = errors.New("token: invalid token")
			return 2, responseDTO
		}

		// get otp
		otp, ok := verifyInfo["otp"]
		if !ok {
			responseDTO.ServerErr = errors.New("cannot convert otp to int")
			return 0, responseDTO
		}

		// check otp
		if otp == strconv.Itoa(int(verifyNumberInput.Code)) {
			// get user
			user, err := s.repo.GetByNumber(verifyNumberInput.Number)

			otpInt, err2 := strconv.Atoi(otp)
			if err2 != nil {
				responseDTO.ServerErr = err
				return 0, responseDTO
			}

			// call domain service
			userErr, serverErr = s.domainService.VerifyNumber(user.Number, uint(otpInt), token, user.IsBlocked)
			responseDTO.ServerErr = serverErr
			responseDTO.UserErr = userErr
			if serverErr != nil || userErr != nil {
				return 0, responseDTO
			}

			if err != nil {
				if err == database_errors.ErrRecordNotFound {

					// user not exist. redirect to signup

					verifyInfo := make(map[string]string)
					verifyInfo["token"] = token
					verifyInfo["verify"] = "signup"

					verifyInfoString, err = utils.ConvertMapToString(verifyInfo)
					if err != nil {
						responseDTO.ServerErr = err
						return 0, responseDTO
					}

					// save number and token to cache for signup
					err = s.cacheRepo.Save(verifyNumberInput.Number, verifyInfoString, 10*time.Minute)
					if err != nil {
						responseDTO.ServerErr = err
						return 0, responseDTO
					}

					responseDTO.ResponseCode = rcodes.GoSignup
					return 2, responseDTO
				}

				responseDTO.ServerErr = err
				return 0, responseDTO
			}

			// user found in database
			tokens, err := jwt.CreateRefreshAndAccessFromUserWithMap(config.JWtRefreshExpire, config.JWTAccessExpire, user.ID, user.Name, user.Number, user.IsRegistered)
			if err != nil {
				responseDTO.ServerErr = err
				return 0, responseDTO
			}

			// create device

			err = s.deviceAppService.CreateOrUpdate(app_device.DeviceInput{
				Name:         deviceName,
				IP:           deviceIP,
				RefreshToken: tokens["refresh"],
				UserID:       user.ID,
			})
			if err != nil {
				responseDTO.ServerErr = err
				return 0, responseDTO
			}

			responseDTO.Data = utils.ConvertMapStringStringToMapStringAny(tokens)
			return 3, responseDTO

		} else {
			responseDTO.ResponseCode = rcodes.WrongCode
			responseDTO.UserErr = errors.New("code: wrong code")
			return 0, responseDTO
		}

	}
}

func (s *service) Signup(userInput SignupUserInput, deviceName string, deviceIP string) (responseDTO datatypes.ResponseDTO) {

	responseDTO.Data = make(map[string]any)

	var user domain_user.User
	user.Number = userInput.Number
	user.Name = userInput.Name
	user.Password = userInput.Password

	userErr, serverErr := s.domainService.Signup(&user, userInput.Token)

	if serverErr != nil {
		responseDTO.ServerErr = serverErr
		return responseDTO
	}
	if userErr != nil {
		responseDTO.UserErr = userErr
		responseDTO.ResponseCode = rcodes.InvalidField
		return responseDTO
	}

	verifyInfoString, err := s.cacheRepo.Get(userInput.Number)
	if err != nil {
		if err == database_errors.ErrRecordNotFound || err == database_errors.ErrExpired {

			responseDTO.ResponseCode = rcodes.VerifyNumberFirst
			responseDTO.UserErr = errors.New("verify number first")
			return responseDTO

		} else {
			responseDTO.ServerErr = err
			return responseDTO
		}
	}

	verifyInfo, err := utils.ConvertStringToMap(verifyInfoString)
	if err != nil {
		responseDTO.ServerErr = err
		return responseDTO
	}

	// check token
	if verifyInfo["token"] != userInput.Token {

		responseDTO.ResponseCode = rcodes.InvalidField
		responseDTO.UserErr = errors.New("token: invalid token")
		return responseDTO
	}

	verify, ok := verifyInfo["verify"]
	if !ok || verify != "signup" {

		responseDTO.ResponseCode = rcodes.VerifyNumberFirst
		responseDTO.UserErr = errors.New("verify number first")
		return responseDTO
	}

	// delete cache
	if err := s.cacheRepo.Delete(userInput.Number + userInput.Token); err != nil {
		responseDTO.ServerErr = err
		return responseDTO
	}

	tokens, err := jwt.CreateRefreshAndAccessFromUserWithMap(config.JWtRefreshExpire, config.JWTAccessExpire, user.ID, user.Name, user.Number, user.IsRegistered)
	if err != nil {
		responseDTO.ServerErr = err
		return responseDTO
	}

	// create device
	err = s.deviceAppService.CreateOrUpdate(app_device.DeviceInput{
		UserID:       user.ID,
		Name:         deviceName,
		IP:           deviceIP,
		RefreshToken: tokens["refresh"],
	})
	if err != nil {
		responseDTO.ServerErr = err
		return responseDTO
	}

	err = s.repo.Create(user)
	responseDTO.ServerErr = err
	return responseDTO
}

func (s *service) GetUserInfo(user domain_user.User) (responseDTO datatypes.ResponseDTO) {
	responseDTO.Data = make(map[string]any)

	var userOutput UserOutput
	userOutput.Fill(user)

	responseDTO.Data["data"] = userOutput
	return responseDTO
}

func (s *service) CheckNumber(numberInput NumberInput) (responseDTO datatypes.ResponseDTO) {
	responseDTO.Data = make(map[string]any)

	err := s.domainService.CheckNumber(numberInput.Number)
	if err != nil {
		responseDTO.ServerErr = err
		return responseDTO
	}

	_, err = s.repo.GetByNumber(numberInput.Number)

	if err != nil {
		if err == database_errors.ErrRecordNotFound {
			responseDTO.Data["isExist"] = false
			return responseDTO
		}
		responseDTO.ServerErr = err
		return responseDTO
	} else {
		responseDTO.Data["isExist"] = true
		return responseDTO
	}

}

func (s *service) Login(loginInput LoginUserInput, deviceName string, deviceIP string) (responseDTO datatypes.ResponseDTO) {
	responseDTO.Data = make(map[string]any)

	user, err := s.repo.GetByNumber(loginInput.Number)
	if err != nil {
		responseDTO.ServerErr = err
		return responseDTO
	}

	userErr, serverErr := s.domainService.Login(user.Number, loginInput.Password, user.Password, user.Salt)
	if serverErr != nil {
		responseDTO.ServerErr = serverErr
		return responseDTO
	}
	if userErr != nil {
		responseDTO.ResponseCode = rcodes.InvalidField
		responseDTO.UserErr = userErr
		return responseDTO
	}

	tokens, err := jwt.CreateRefreshAndAccessFromUserWithMap(config.JWtRefreshExpire, config.JWTAccessExpire, user.ID, user.Name, user.Number, user.IsRegistered)
	if err != nil {
		responseDTO.ServerErr = err
		return responseDTO
	}

	// create device
	err = s.deviceAppService.CreateOrUpdate(app_device.DeviceInput{
		UserID:       user.ID,
		Name:         deviceName,
		IP:           deviceIP,
		RefreshToken: tokens["refresh"],
	})
	if err != nil {
		responseDTO.ServerErr = err
		return responseDTO
	}
	responseDTO.Data = utils.ConvertMapStringStringToMapStringAny(tokens)
	return responseDTO
}

func (s *service) GetAccessFromRefresh(refresh string) (responseDTO datatypes.ResponseDTO) {
	responseDTO.Data = make(map[string]any)

	user, userErr, serverErr := s.deviceAppService.GetDeviceUserByRefreshToken(refresh)
	if serverErr != nil {
		responseDTO.ServerErr = serverErr
		return responseDTO
	}
	if userErr != nil {
		responseDTO.UserErr = userErr
		return responseDTO
	}

	access, err := jwt.CreateAccessFromUser(config.JWTAccessExpire, user.ID, user.Name, user.Number, user.IsRegistered)
	responseDTO.ServerErr = err
	responseDTO.Data["access"] = access

	return responseDTO

}
