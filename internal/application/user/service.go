package app_user

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	app_device "github.com/yaghoubi-mn/pedarkharj/internal/application/device"

	// app_user "github.com/yaghoubi-mn/pedarkharj/internal/application/user"
	domain_device "github.com/yaghoubi-mn/pedarkharj/internal/domain/device"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	"github.com/yaghoubi-mn/pedarkharj/internal/infrastructure/config"
	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/jwt"
	"github.com/yaghoubi-mn/pedarkharj/pkg/rcodes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/s3"
	"github.com/yaghoubi-mn/pedarkharj/pkg/service_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/sms"
	"github.com/yaghoubi-mn/pedarkharj/pkg/utils"
)

type UserAppService interface {
	VerifyNumber(verifyNumberInput VerifyNumberInput, deviceName string, deviceIP string) (step int, responseDTO datatypes.ResponseDTO)
	Signup(userInput SignupUserInput, deviceName string, deviceIP string) (responseDTO datatypes.ResponseDTO)
	GetUserInfo(userID uint64) datatypes.ResponseDTO
	CheckNumber(numberInput NumberInput) datatypes.ResponseDTO
	Login(loginInput LoginUserInput, deviceName string, deviceIP string) (responseDTO datatypes.ResponseDTO)
	GetAccessFromRefresh(refresh string) (responseDTO datatypes.ResponseDTO)
	ChooseUserAvatar(avatarName string, userID uint64) datatypes.ResponseDTO
	GetAvatars() datatypes.ResponseDTO
	ResetPassword(input RestPasswordInput) datatypes.ResponseDTO
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

	userErr, serverErr := s.domainService.VerifyNumber(domain_user.VerifyNumberInput{
		Number: verifyNumberInput.Number,
		OTP:    verifyNumberInput.OTP,
		Token:  verifyNumberInput.Token,
		Mode:   verifyNumberInput.Mode,
	})
	if serverErr != nil {
		responseDTO.ServerErr = serverErr
		return 0, responseDTO
	}
	if userErr != nil {
		responseDTO.UserErr = userErr
		responseDTO.ResponseCode = rcodes.InvalidField
		return 0, responseDTO
	}

	if verifyNumberInput.OTP == 0 {
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

				if config.Debug == false {

					// send code to number
					err = sms.SendOTPSMS(verifyNumberInput.Number[3:], otp)
					if err != nil {
						responseDTO.ServerErr = err
						return 1, responseDTO
					}
				}

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
		if otp == strconv.Itoa(int(verifyNumberInput.OTP)) {

			// get user
			user, databaseErr := s.repo.GetByNumber(verifyNumberInput.Number)

			var isUserRegistered bool
			isUserExist := true

			isUserRegistered = user.IsRegistered

			if databaseErr != nil {
				if databaseErr == database_errors.ErrRecordNotFound {
					isUserRegistered = false
					isUserExist = false

				} else {
					responseDTO.ServerErr = err
					return 0, responseDTO
				}
			}

			if !isUserRegistered {
				// user not registered

				// becuase user not found only signup mode is allowed
				if verifyNumberInput.Mode != "signup" {
					responseDTO.ResponseCode = rcodes.UserNotRegistered
					responseDTO.UserErr = errors.New("user not exist. reset_password not allowed")
					return 0, responseDTO
				}

				// user not exist. redirect to signup

				verifyInfo := make(map[string]string)
				verifyInfo["token"] = token
				verifyInfo["mode"] = "signup"
				verifyInfo["is_user_exist"] = strconv.FormatBool(isUserExist)

				verifyInfoString, err = utils.ConvertMapToString(verifyInfo)
				if err != nil {
					responseDTO.ServerErr = err
					return 0, responseDTO
				}

				// save number and token to cache for signup
				err = s.cacheRepo.Save(verifyNumberInput.Number, verifyInfoString, config.VerifyNumberCacheExpireTime)
				if err != nil {
					responseDTO.ServerErr = err
					return 0, responseDTO
				}

				responseDTO.ResponseCode = rcodes.GoSignup
				return 2, responseDTO

			} else {
				// user registered

				// user exist so just reset password allowed
				if verifyNumberInput.Mode != "reset_password" {
					responseDTO.ResponseCode = rcodes.UserAlreadyRegistered
					responseDTO.UserErr = errors.New("user already exist. signup mode not allowed")
					return 0, responseDTO
				}

				// check is blocked
				if user.IsBlocked {
					responseDTO.UserErr = errors.New("you are blocked")
					return 0, responseDTO
				}

				verifyInfo := make(map[string]string)
				verifyInfo["token"] = token
				verifyInfo["mode"] = "reset_password"

				verifyInfoString, err = utils.ConvertMapToString(verifyInfo)
				if err != nil {
					responseDTO.ServerErr = err
					return 0, responseDTO
				}

				err := s.cacheRepo.Save(verifyNumberInput.Number, verifyInfoString, config.VerifyNumberCacheExpireTime)
				if err != nil {
					responseDTO.ServerErr = err
					return 0, responseDTO
				}

				responseDTO.Data["msg"] = "go to reset password"
				responseDTO.ResponseCode = rcodes.GoRestPassword
				return 3, responseDTO
			}
		} else {
			responseDTO.ResponseCode = rcodes.WrongOTP
			responseDTO.UserErr = errors.New("code: wrong code")
			return 0, responseDTO
		}

	}
}

func (s *service) Signup(userInput SignupUserInput, deviceName string, deviceIP string) (responseDTO datatypes.ResponseDTO) {

	responseDTO.Data = make(map[string]any)

	// call domain service
	user, userErr, serverErr := s.domainService.Signup(domain_user.SignupUserInput{
		Number:   userInput.Number,
		Name:     userInput.Name,
		Password: userInput.Password,
		Token:    userInput.Token,
	})

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
		responseDTO.UserErr = service_errors.ErrInvalidToken
		return responseDTO
	}

	verify, ok := verifyInfo["mode"]
	if !ok || verify != "signup" {

		responseDTO.ResponseCode = rcodes.VerifyNumberFirst
		responseDTO.UserErr = errors.New("verify number first")
		return responseDTO
	}

	// delete cache
	if err := s.cacheRepo.Delete(userInput.Number); err != nil {
		responseDTO.ServerErr = err
		return responseDTO
	}

	// select random avatar for user
	// get list of avatars
	avatars, err := s3.GetListObjects(config.AvatarPath)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}

	// select random avatar
	randomIndex := rand.Intn(len(avatars))
	user.Avatar = avatars[randomIndex]

	isUserExist, ok := verifyInfo["is_user_exist"]
	if !ok {
		responseDTO.ServerErr = errors.New("is_user_exist not found in verifyInfo map in signup")
		return
	}

	if isUserExist == "true" {
		// update user

		err = s.repo.UpdateColumns(user)
		if err != nil {
			responseDTO.ServerErr = err
			return
		}
	} else {
		// insert user into database
		err = s.repo.Create(&user)
		if err != nil {
			responseDTO.ServerErr = err
			return responseDTO
		}
	}

	tokens, err := jwt.CreateRefreshAndAccessFromUserWithMap(config.JWtRefreshExpire, config.JWTAccessExpire, user.ID, user.Name, user.Number, user.IsRegistered)
	if err != nil {
		responseDTO.ServerErr = err
		return responseDTO
	}

	// create device
	err = s.deviceAppService.CreateOrUpdate(domain_device.DeviceInput{
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

func (s *service) ResetPassword(input RestPasswordInput) (responseDTO datatypes.ResponseDTO) {
	responseDTO.Data = make(map[string]any)

	userErr, serverErr, salt, hashedPassword := s.domainService.ResetPassword(domain_user.RestPasswordInput{
		Number:   input.Number,
		Password: input.Password,
		Token:    input.Token,
	})
	if serverErr != nil {
		responseDTO.ServerErr = serverErr
		return
	}
	if userErr != nil {
		responseDTO.UserErr = userErr
		return
	}

	verifyInfoString, err := s.cacheRepo.Get(input.Number)
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

	if token, ok := verifyInfo["token"]; !ok && token != input.Token {
		responseDTO.UserErr = errors.New("token: invlid token")
		return
	}

	if mode, ok := verifyInfo["mode"]; !ok && mode != "reset_password" {
		responseDTO.Data["msg"] = "invalid mode"
		return
	}

	user, err := s.repo.GetByNumber(input.Number)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}

	user.Password = hashedPassword
	user.Salt = salt

	err = s.repo.Update(user)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}

	s.cacheRepo.Delete(input.Number)

	tokens, err := jwt.CreateRefreshAndAccessFromUserWithMap(config.JWtRefreshExpire, config.JWTAccessExpire, user.ID, user.Name, user.Number, user.IsRegistered)
	if err != nil {
		responseDTO.ServerErr = err
		return responseDTO
	}

	responseDTO.Data = utils.ConvertMapStringStringToMapStringAny(tokens)
	responseDTO.Data["msg"] = "password updated"
	return

}

func (s *service) GetUserInfo(userID uint64) (responseDTO datatypes.ResponseDTO) {
	responseDTO.Data = make(map[string]any)

	var userOutput UserOutput
	// get user from database for full information
	user, err := s.repo.GetByID(userID)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}
	userOutput.Fill(user)

	responseDTO.Data["data"] = userOutput
	return responseDTO
}

func (s *service) CheckNumber(numberInput NumberInput) (responseDTO datatypes.ResponseDTO) {
	responseDTO.Data = make(map[string]any)

	err := s.domainService.CheckNumber(numberInput.Number)
	if err != nil {
		responseDTO.UserErr = err
		responseDTO.ResponseCode = rcodes.InvalidField
		return responseDTO
	}

	user, err := s.repo.GetByNumber(numberInput.Number)
	isExist := false
	if err != nil {
		if err == database_errors.ErrRecordNotFound {
			isExist = false
		} else {
			responseDTO.ServerErr = err
			return responseDTO
		}

	} else {

		if user.IsRegistered {
			isExist = true
		} else {
			isExist = false
		}
	}

	responseDTO.Data["isExist"] = isExist
	return responseDTO

}

func (s *service) Login(loginInput LoginUserInput, deviceName string, deviceIP string) (responseDTO datatypes.ResponseDTO) {
	responseDTO.Data = make(map[string]any)

	if err := s.domainService.CheckNumber(loginInput.Number); err != nil {
		responseDTO.UserErr = err
		return
	}

	user, err := s.repo.GetByNumber(loginInput.Number)
	if err != nil {
		if err == database_errors.ErrRecordNotFound {
			responseDTO.ResponseCode = rcodes.NumberNotExist
			responseDTO.Data["msg"] = "number not exist."
			responseDTO.UserErr = errors.New("number: number not exist")
			return
		}
		responseDTO.ServerErr = err
		return responseDTO
	}

	userErr, serverErr := s.domainService.Login(domain_user.LoginUserInput{
		Number:        loginInput.Number,
		InputPassword: loginInput.Password,
		RealPassword:  user.Password,
		Salt:          user.Salt,
		IsBlocked:     user.IsBlocked,
		IsRegistered:  user.IsRegistered,
	})
	if serverErr != nil {
		responseDTO.ServerErr = serverErr
		return responseDTO
	}
	if userErr != nil {
		responseDTO.ResponseCode = rcodes.InvalidField
		if userErr == service_errors.ErrUserNotRegistered {
			responseDTO.ResponseCode = rcodes.UserNotRegistered
		}
		responseDTO.UserErr = userErr
		return responseDTO
	}

	tokens, err := jwt.CreateRefreshAndAccessFromUserWithMap(config.JWtRefreshExpire, config.JWTAccessExpire, user.ID, user.Name, user.Number, user.IsRegistered)
	if err != nil {
		responseDTO.ServerErr = err
		return responseDTO
	}

	// create device
	err = s.deviceAppService.CreateOrUpdate(domain_device.DeviceInput{
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

func (s *service) ChooseUserAvatar(avatarName string, userID uint64) (responseDTO datatypes.ResponseDTO) {
	responseDTO.Data = make(map[string]any)

	avatars, err := s3.GetListObjects(config.AvatarPath)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}

	found := false
	for _, avatar := range avatars {
		if avatar == avatarName {
			found = true
			break
		}
	}

	if !found {
		responseDTO.ResponseCode = rcodes.AvatarNotFound
		responseDTO.UserErr = errors.New("avatar: avatar not found")
		return
	}

	user, err := s.repo.GetByID(userID)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}

	user.Avatar = avatarName
	err = s.repo.Update(user)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}

	responseDTO.Data["msg"] = "avatar saved"
	return
}

func (s *service) GetAvatars() (responseDTO datatypes.ResponseDTO) {
	responseDTO.Data = make(map[string]any)

	avatars, err := s3.GetListObjects(config.AvatarPath)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}

	responseDTO.Data["data"] = avatars
	return

}
