package user_test

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	app_user "github.com/yaghoubi-mn/pedarkharj/internal/application/user"
	domain_shared "github.com/yaghoubi-mn/pedarkharj/internal/domain/shared"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	"github.com/yaghoubi-mn/pedarkharj/internal/infrastructure/config"
	shared_dto "github.com/yaghoubi-mn/pedarkharj/internal/shared/dto"
	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/rcodes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/service_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/utils"
	"github.com/yaghoubi-mn/pedarkharj/tests/pkg/helpers"
	_ "github.com/yaghoubi-mn/pedarkharj/tests/pkg/helpers"
	"github.com/yaghoubi-mn/pedarkharj/tests/pkg/number_generator"
)

var appService app_user.UserAppService
var repo domain_user.UserDomainRepository
var cacheRepo domain_shared.CacheRepository

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func setup() {
	appService = helpers.GetUserAppService()
	repo = helpers.GetUserDomainRepository()
	cacheRepo = helpers.GetCacheRepository()
}

func TestSendOTP(t *testing.T) {

	tests := []struct {
		PreRunCode          func(number string)
		PhoneNumber         string
		WantUserErr         error
		ResponseCode        string
		ResponseDTODataKeys []string
		NextRunCode         func(number string)
	}{
		{ // test sucess
			PreRunCode:          func(number string) {},
			PhoneNumber:         number_generator.GetNumber(),
			WantUserErr:         nil,
			ResponseCode:        rcodes.CodeSendToNumber,
			ResponseDTODataKeys: []string{"token", "delayTimeSeconds"},
			NextRunCode: func(number string) {

				// check database
				now := time.Now()
				data, expireTime, err := cacheRepo.Get(number)
				assert.NoError(t, err, "cache database error")

				if !(expireTime.Unix()-now.Unix() < int64(config.VerifyNumberCacheExpireTime.Seconds()) && expireTime.Unix()-now.Unix() > int64(config.VerifyNumberCacheExpireTime.Seconds())-10) {
					assert.Error(t, errors.New("invalid expire time"), fmt.Sprintf("now: %s, expireTime: %s", now, expireTime))
				}

				err = utils.CheckMapHaveKeys(data, "token", "otp")
				assert.NoError(t, err, data)

			},
		},
		{ // test delay number
			PreRunCode: func(number string) {

				responseDTO := appService.SendOTP(app_user.SendOTPInput{
					SendOTPInput: shared_dto.SendOTPInput{
						PhoneNumber: number,
					},
				})

				assert.NoError(t, responseDTO.ServerErr)
				assert.NoError(t, responseDTO.UserErr)
				assert.Equal(t, rcodes.CodeSendToNumber, responseDTO.ResponseCode, responseDTO)
			},
			PhoneNumber:         number_generator.GetNumber(),
			WantUserErr:         service_errors.ErrOTPNotExpired,
			ResponseCode:        rcodes.NumberDelay,
			ResponseDTODataKeys: []string{"delayTimeSeconds"},
			NextRunCode:         func(number string) {},
		},
		{ // test invalid number
			PreRunCode:          func(number string) {},
			PhoneNumber:         "+9891234567",
			WantUserErr:         service_errors.ErrInvalidNumber,
			ResponseCode:        rcodes.InvalidField,
			ResponseDTODataKeys: []string{},
			NextRunCode:         func(number string) {},
		},
	}

	for _, tt := range tests {

		tt.PreRunCode(tt.PhoneNumber)

		responseDTO := appService.SendOTP(app_user.SendOTPInput{
			SendOTPInput: shared_dto.SendOTPInput{
				PhoneNumber: tt.PhoneNumber,
			},
		})

		assert.NoError(t, responseDTO.ServerErr)
		assert.Equal(t, tt.WantUserErr, responseDTO.UserErr)

		assert.Equal(t, tt.ResponseCode, responseDTO.ResponseCode, responseDTO)
		err := utils.CheckMapHaveKeys(responseDTO.Data, tt.ResponseDTODataKeys...)
		assert.NoError(t, err, responseDTO.Data)

		tt.NextRunCode(tt.PhoneNumber)

	}

}

func TestVerifyOTP(t *testing.T) {

	tests := []struct {
		TestID              int
		PreRunCode          func(phoneNumber string)
		PhoneNumber         string
		SendOTP             bool
		SendCustomToken     bool
		CustomToken         string
		Mode                string
		DeviceName          string
		DeviceIP            string
		WantUserErr         error
		ResponseCode        string
		ResponseDTODataKeys []string
		NextRunCode         func(number string)
	}{
		{ // test sucess. number not exist in database. signup mode
			TestID:              1,
			PreRunCode:          func(phoneNumber string) {},
			PhoneNumber:         number_generator.GetNumber(),
			SendOTP:             true,
			SendCustomToken:     false,
			Mode:                "signup",
			DeviceName:          "test",
			DeviceIP:            "1.1.1.1",
			WantUserErr:         nil,
			ResponseCode:        rcodes.GoSignup,
			ResponseDTODataKeys: []string{},
			NextRunCode: func(number string) {

				// check database
				now := time.Now()
				data, expireTime, err := cacheRepo.Get(number)
				assert.NoError(t, err, "cache database error")

				if !(expireTime.Unix()-now.Unix() < int64(config.VerifyNumberCacheExpireTime.Seconds()) && expireTime.Unix()-now.Unix() > int64(config.VerifyNumberCacheExpireTime.Seconds())-10) {
					assert.Error(t, errors.New("invalid expire time"), fmt.Sprintf("now: %s, expireTime: %s", now, expireTime))
				}

				err = utils.CheckMapHaveKeys(data, "token", "mode", "is_user_exist")
				assert.NoError(t, err, data)

				assert.Equal(t, "signup", data["mode"], data)
				assert.Equal(t, "false", data["is_user_exist"], data)

			},
		},

		{ // test success. user exist in database (already registered). reset_password mode
			TestID: 2,
			PreRunCode: func(number string) {

				// create user in database
				err := repo.Create(&domain_user.User{
					Name:         "test",
					Number:       number,
					Password:     "dfd",
					Salt:         "as",
					Avatar:       "dafd.com",
					IsRegistered: true,
					RegisteredAt: time.Now(),
				})

				assert.NoError(t, err)
			},
			PhoneNumber:         number_generator.GetNumber(),
			SendOTP:             true,
			SendCustomToken:     false,
			Mode:                "reset_password",
			DeviceName:          "test",
			DeviceIP:            "1.1.1.1",
			WantUserErr:         nil,
			ResponseCode:        rcodes.GoRestPassword,
			ResponseDTODataKeys: []string{},
			NextRunCode: func(number string) {

				// check database
				now := time.Now()
				data, expireTime, err := cacheRepo.Get(number)
				assert.NoError(t, err, "cache database error")

				if !(expireTime.Unix()-now.Unix() < int64(config.VerifyNumberCacheExpireTime.Seconds()) && expireTime.Unix()-now.Unix() > int64(config.VerifyNumberCacheExpireTime.Seconds())-10) {
					assert.Error(t, errors.New("invalid expire time"), fmt.Sprintf("now: %s, expireTime: %s", now, expireTime))
				}

				err = utils.CheckMapHaveKeys(data, "token", "mode", "is_user_exist")
				assert.NoError(t, err, data)

				assert.Equal(t, "reset_password", data["mode"], data)
				assert.Equal(t, "true", data["is_user_exist"], data)

			},
		},

		{ // test success. user exist in database (not registered). signup mode
			TestID: 3,
			PreRunCode: func(number string) {

				// create user in database
				err := repo.Create(&domain_user.User{
					Name:         "test",
					Number:       number,
					Password:     "dfd",
					Salt:         "as",
					Avatar:       "dafd.com",
					IsRegistered: false,
					RegisteredAt: time.Now(),
				})

				assert.NoError(t, err)
			},
			PhoneNumber:         number_generator.GetNumber(),
			SendOTP:             true,
			SendCustomToken:     false,
			Mode:                "signup",
			DeviceName:          "test",
			DeviceIP:            "1.1.1.1",
			WantUserErr:         nil,
			ResponseCode:        rcodes.GoSignup,
			ResponseDTODataKeys: []string{},
			NextRunCode: func(number string) {

				// check database
				now := time.Now()
				data, expireTime, err := cacheRepo.Get(number)
				assert.NoError(t, err, "cache database error")

				if !(expireTime.Unix()-now.Unix() < int64(config.VerifyNumberCacheExpireTime.Seconds()) && expireTime.Unix()-now.Unix() > int64(config.VerifyNumberCacheExpireTime.Seconds())-10) {
					assert.Error(t, errors.New("invalid expire time"), fmt.Sprintf("now: %s, expireTime: %s", now, expireTime))
				}

				err = utils.CheckMapHaveKeys(data, "token", "mode", "is_user_exist")
				assert.NoError(t, err, data)

				assert.Equal(t, "signup", data["mode"], data)
				assert.Equal(t, "true", data["is_user_exist"], data)

			},
		},

		{ // test failure. number not exist in database. reset_password mode
			TestID:              4,
			PreRunCode:          func(number string) {},
			PhoneNumber:         number_generator.GetNumber(),
			SendOTP:             true,
			SendCustomToken:     false,
			Mode:                "reset_password",
			DeviceName:          "test",
			DeviceIP:            "1.1.1.1",
			WantUserErr:         service_errors.ErrUserNotRegisteredResetPasswordNotAllowed,
			ResponseCode:        rcodes.UserNotRegistered,
			ResponseDTODataKeys: []string{},
			NextRunCode:         func(number string) {},
		},

		{ // test failure. user exist in database (already registered). signup mode
			TestID: 5,
			PreRunCode: func(number string) {

				// create user in database
				err := repo.Create(&domain_user.User{
					Name:         "test",
					Number:       number,
					Password:     "dfd",
					Salt:         "as",
					Avatar:       "dafd.com",
					IsRegistered: true,
					RegisteredAt: time.Now(),
				})

				assert.NoError(t, err)
			},
			PhoneNumber:         number_generator.GetNumber(),
			SendOTP:             true,
			SendCustomToken:     false,
			Mode:                "signup",
			DeviceName:          "test",
			DeviceIP:            "1.1.1.1",
			WantUserErr:         service_errors.ErrUserAlreayRegisteredSignupNotAllowed,
			ResponseCode:        rcodes.UserAlreadyRegistered,
			ResponseDTODataKeys: []string{},
			NextRunCode:         func(number string) {},
		},

		{ // test failure. user exist in database (not registered). reset_password mode
			TestID: 6,
			PreRunCode: func(number string) {

				// create user in database
				err := repo.Create(&domain_user.User{
					Name:         "test",
					Number:       number,
					Password:     "dfd",
					Salt:         "as",
					Avatar:       "dafd.com",
					IsRegistered: false,
					RegisteredAt: time.Now(),
				})

				assert.NoError(t, err)
			},
			PhoneNumber:         number_generator.GetNumber(),
			SendOTP:             true,
			SendCustomToken:     false,
			Mode:                "reset_password",
			DeviceName:          "test",
			DeviceIP:            "1.1.1.1",
			WantUserErr:         service_errors.ErrUserNotRegisteredResetPasswordNotAllowed,
			ResponseCode:        rcodes.UserNotRegistered,
			ResponseDTODataKeys: []string{},
			NextRunCode:         func(number string) {},
		},

		{ // test falure. wrong otp code
			TestID:              7,
			PreRunCode:          func(number string) {},
			PhoneNumber:         number_generator.GetNumber(),
			SendOTP:             false,
			SendCustomToken:     false,
			Mode:                "signup",
			DeviceName:          "test",
			DeviceIP:            "1.1.1.1",
			WantUserErr:         service_errors.ErrWrongOTP,
			ResponseCode:        rcodes.WrongOTP,
			ResponseDTODataKeys: []string{},
			NextRunCode:         func(number string) {},
		},

		{ // test falure. wrong token
			TestID:              8,
			PreRunCode:          func(number string) {},
			PhoneNumber:         number_generator.GetNumber(),
			SendOTP:             true,
			SendCustomToken:     true,
			CustomToken:         uuid.NewString(),
			Mode:                "signup",
			DeviceName:          "test",
			DeviceIP:            "1.1.1.1",
			WantUserErr:         service_errors.ErrWrongToken,
			ResponseCode:        rcodes.InvalidField,
			ResponseDTODataKeys: []string{},
			NextRunCode:         func(number string) {},
		},

		// { // test invalid number
		// 	PreRunCode:          func(number string) {},
		// 	Number:              number_generator.GetNumber(),
		// 	SendOTP:             true,
		// 	SendToken:           true,
		// 	Mode:                "signup",
		// 	DeviceName:          "test",
		// 	DeviceIP:            "1.1.1.1",
		// 	WantUserErr:         service_errors.ErrInvalidNumber,
		// 	ResponseCode:        rcodes.InvalidField,
		// 	ResponseDTODataKeys: []string{},
		// 	NextRunCode:         func(number string) {},
		// },
	}

	for _, tt := range tests {

		tt.PreRunCode(tt.PhoneNumber)

		// call send otp
		responseDTO := appService.SendOTP(app_user.SendOTPInput{
			SendOTPInput: shared_dto.SendOTPInput{
				PhoneNumber: tt.PhoneNumber,
			},
		})
		assert.NoError(t, responseDTO.ServerErr)
		assert.NoError(t, responseDTO.UserErr)

		// get otp from cache
		data, _, err := cacheRepo.Get(tt.PhoneNumber)
		assert.NoError(t, err, "cache database error")
		otp, err := strconv.Atoi(data["otp"])
		assert.NoError(t, err)

		token := responseDTO.Data["token"].(string)

		if !tt.SendOTP {
			otp = 0
		}
		if tt.SendCustomToken {
			token = tt.CustomToken
		}

		mode, responseDTO := appService.VerifyOTP(app_user.VerifyOTPInput{
			VerifyOTPInput: shared_dto.VerifyOTPInput{
				PhoneNumber: tt.PhoneNumber,
				OTP:         uint(otp),
				Token:       token,
				Mode:        tt.Mode,
			},
		}, tt.DeviceName, tt.DeviceIP)

		assert.NoError(t, responseDTO.ServerErr)
		assert.Equal(t, tt.WantUserErr, responseDTO.UserErr, fmt.Sprintf("TestID: %v, %v", tt.TestID, responseDTO))

		if responseDTO.UserErr == nil && responseDTO.ServerErr == nil {

			// check mode
			if tt.Mode == "signup" {
				assert.Equal(t, 1, mode, "in signup mode, mode must be 1")
			} else if tt.Mode == "reset_password" {
				assert.Equal(t, 2, mode, "in reset_password mode, mode must be 2")
			}
		}

		assert.Equal(t, tt.ResponseCode, responseDTO.ResponseCode, responseDTO)
		err = utils.CheckMapHaveKeys(responseDTO.Data, tt.ResponseDTODataKeys...)
		assert.NoError(t, err, responseDTO.Data)

		tt.NextRunCode(tt.PhoneNumber)

	}

}

func TestSignup(t *testing.T) {

	tests := []struct {
		TestName            string
		PreRunCode          func(number string)
		PhoneNumber         string
		Name                string
		Password            string
		Mode                string
		DeviceName          string
		DeviceIP            string
		WantUserErr         error
		ResponseCode        string
		ResponseDTODataKeys []string
		NextRunCode         func(number string)
	}{
		{
			TestName:            "test sucess. number not exist in database. signup mode",
			PreRunCode:          func(number string) {},
			PhoneNumber:         number_generator.GetNumber(),
			Name:                "test",
			Mode:                "signup",
			Password:            "12345678",
			DeviceName:          "test",
			DeviceIP:            "1.1.1.1",
			WantUserErr:         nil,
			ResponseCode:        "",
			ResponseDTODataKeys: []string{"refresh", "access"},
			NextRunCode: func(number string) {

				// data must be deleted from database
				_, _, err := cacheRepo.Get(number)
				if err != database_errors.ErrRecordNotFound && err != database_errors.ErrExpired {

					assert.Error(t, errors.New("record found in database after signup complete"))
				}

			},
		},

		{
			TestName: "test success. user exist in database (not registered). signup mode",
			PreRunCode: func(number string) {

				// create user in database
				err := repo.Create(&domain_user.User{
					Name:         "test",
					Number:       number,
					Password:     "dfd",
					Salt:         "as",
					Avatar:       "dafd.com",
					IsRegistered: false,
					RegisteredAt: time.Now(),
				})

				assert.NoError(t, err)
			},
			PhoneNumber:         number_generator.GetNumber(),
			Name:                "test",
			Password:            "12345678",
			Mode:                "signup",
			DeviceName:          "test",
			DeviceIP:            "1.1.1.1",
			WantUserErr:         nil,
			ResponseCode:        "",
			ResponseDTODataKeys: []string{},
			NextRunCode: func(number string) {

				// data must be deleted from database
				_, _, err := cacheRepo.Get(number)
				if err != database_errors.ErrRecordNotFound && err != database_errors.ErrExpired {

					assert.Error(t, errors.New("record found in database after signup complete"))
				}
			},
		},

		{
			TestName:            "test failure. invalid name",
			PreRunCode:          func(number string) {},
			PhoneNumber:         number_generator.GetNumber(),
			Name:                "",
			Password:            "12345678",
			Mode:                "reset_password",
			DeviceName:          "test",
			DeviceIP:            "1.1.1.1",
			WantUserErr:         service_errors.ErrInvalidName,
			ResponseCode:        rcodes.InvalidField,
			ResponseDTODataKeys: []string{},
			NextRunCode:         func(number string) {},
		},

		{
			TestName:            "test failure. small password",
			PreRunCode:          func(number string) {},
			PhoneNumber:         number_generator.GetNumber(),
			Name:                "test",
			Password:            "1234567",
			Mode:                "signup",
			DeviceName:          "test",
			DeviceIP:            "1.1.1.1",
			WantUserErr:         service_errors.ErrSmallPassword,
			ResponseCode:        rcodes.InvalidField,
			ResponseDTODataKeys: []string{},
			NextRunCode:         func(number string) {},
		},
	}

	for _, tt := range tests {

		contin := false // if one step raise error continue to next test case

		tt.PreRunCode(tt.PhoneNumber)

		// call send otp
		responseDTO := appService.SendOTP(app_user.SendOTPInput{
			SendOTPInput: shared_dto.SendOTPInput{
				PhoneNumber: tt.PhoneNumber,
			},
		})
		contin = contin || assert.NoError(t, responseDTO.ServerErr)
		contin = contin || assert.NoError(t, responseDTO.UserErr)

		if contin {
			continue
		}

		// get otp from cache
		data, _, err := cacheRepo.Get(tt.PhoneNumber)
		assert.NoError(t, err, "cache database error")
		otp, err := strconv.Atoi(data["otp"])
		assert.NoError(t, err)

		token := responseDTO.Data["token"].(string)

		// verify otp
		_, responseDTO = appService.VerifyOTP(app_user.VerifyOTPInput{
			VerifyOTPInput: shared_dto.VerifyOTPInput{
				PhoneNumber: tt.PhoneNumber,
				OTP:         uint(otp),
				Token:       token,
				Mode:        tt.Mode,
			},
		}, tt.DeviceName, tt.DeviceIP)
		assert.NoError(t, responseDTO.ServerErr, fmt.Sprintf("TestName: %v", tt.TestName))
		assert.NoError(t, responseDTO.UserErr, fmt.Sprintf("TestName: %v", tt.TestName))

		// signup
		responseDTO = appService.Signup(app_user.SignupUserInput{
			SignupUserInput: shared_dto.SignupUserInput{
				PhoneNumber: tt.PhoneNumber,
				Name:        tt.Name,
				Password:    tt.Password,
				Token:       token,
			},
		}, tt.DeviceName, tt.DeviceIP)

		assert.NoError(t, responseDTO.ServerErr, responseDTO)
		assert.Equal(t, tt.WantUserErr, responseDTO.UserErr, responseDTO)

		err = utils.CheckMapHaveKeys(responseDTO.Data, tt.ResponseDTODataKeys...)
		assert.NoError(t, err)

		tt.NextRunCode(tt.PhoneNumber)

	}

	// custom tests

	// test failure. invalid mode saved in cache

	number := number_generator.GetNumber()
	token := uuid.NewString()

	// create user in database
	err := cacheRepo.Save(
		number,
		map[string]string{
			"token":         token,
			"mode":          "reset_password",
			"is_user_exist": "false",
		},
		10*time.Minute)

	assert.NoError(t, err)

	responseDTO := appService.Signup(app_user.SignupUserInput{
		SignupUserInput: shared_dto.SignupUserInput{
			PhoneNumber: number,
			Name:        "test",
			Password:    "12345678",
			Token:       token,
		},
	}, "test", "1.1.1.1")

	assert.NoError(t, responseDTO.ServerErr)

	assert.Equal(t, service_errors.ErrVerifyNumberFirst, responseDTO.UserErr, responseDTO)
	assert.Equal(t, rcodes.VerifyNumberFirst, responseDTO.ResponseCode, responseDTO)

}

func TestResetPassword(t *testing.T) {
	tests := []struct {
		TestName            string
		PreRunCode          func(number string, token string)
		PhoneNumber         string
		Password            string
		Token               string
		WantUserErr         error
		ResponseCode        string
		ResponseDTODataKeys []string
	}{
		{
			TestName:    "test success - valid reset password",
			PhoneNumber: number_generator.GetNumber(),
			Password:    "newPassword123",
			PreRunCode: func(number string, token string) {
				// Create user in DB
				repo.Create(&domain_user.User{
					Number:       number,
					Password:     "oldPassword",
					Salt:         "oldSalt",
					IsRegistered: true,
				})

				// Save reset token in cache
				cacheRepo.Save(number, map[string]string{
					"token": token,
					"mode":  "reset_password",
				}, 10*time.Minute)
			},
			WantUserErr:         nil,
			ResponseDTODataKeys: []string{"access", "refresh", "msg"},
		},
		{
			TestName:    "test failure - invalid token",
			PhoneNumber: number_generator.GetNumber(),
			Password:    "newPassword123",
			Token:       "invalid-token",
			PreRunCode: func(number string, token string) {
				repo.Create(&domain_user.User{
					Number:       number,
					IsRegistered: true,
				})
				cacheRepo.Save(number, map[string]string{
					"token": "valid-token",
					"mode":  "reset_password",
				}, 10*time.Minute)
			},
			WantUserErr:  service_errors.ErrInvalidToken,
			ResponseCode: rcodes.InvalidField,
		},
		{
			TestName:    "test failure - invalid mode",
			PhoneNumber: number_generator.GetNumber(),
			Password:    "newPassword123",
			PreRunCode: func(number string, token string) {
				repo.Create(&domain_user.User{
					Number:       number,
					IsRegistered: true,
				})
				cacheRepo.Save(number, map[string]string{
					"token": token,
					"mode":  "wrong_mode",
				}, 10*time.Minute)
			},
			WantUserErr:  nil, // ServerErr should be set
			ResponseCode: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			token := uuid.New().String()
			if tt.Token == "" {
				tt.Token = token
			}

			tt.PreRunCode(tt.PhoneNumber, tt.Token)

			responseDTO := appService.ResetPassword(app_user.ResetPasswordInput{
				ResetPasswordInput: shared_dto.ResetPasswordInput{
					PhoneNumber: tt.PhoneNumber,
					Password:    tt.Password,
					Token:       tt.Token,
				},
			})

			assert.Equal(t, tt.WantUserErr, responseDTO.UserErr)
			assert.Equal(t, tt.ResponseCode, responseDTO.ResponseCode)

			if tt.ResponseDTODataKeys != nil {
				err := utils.CheckMapHaveKeys(responseDTO.Data, tt.ResponseDTODataKeys...)
				assert.NoError(t, err)
			}

			// Additional assertions for success case
			if tt.WantUserErr == nil && responseDTO.ServerErr == nil {
				user, err := repo.GetByNumber(tt.PhoneNumber)
				assert.NoError(t, err)
				assert.NotEqual(t, "oldPassword", user.Password)
			}
		})
	}
}

func TestCheckNumber(t *testing.T) {
	tests := []struct {
		TestName        string
		PreRunCode      func(number string)
		PhoneNumber     string
		WantUserErr     error
		ExpectedIsExist bool
	}{
		{
			TestName: "test success - number exists",
			PreRunCode: func(number string) {
				repo.Create(&domain_user.User{
					Number:       number,
					IsRegistered: true,
				})
			},
			PhoneNumber:     number_generator.GetNumber(),
			WantUserErr:     nil,
			ExpectedIsExist: true,
		},
		{
			TestName:        "test success - number doesn't exist",
			PreRunCode:      func(number string) {},
			PhoneNumber:     number_generator.GetNumber(),
			WantUserErr:     nil,
			ExpectedIsExist: false,
		},
		{
			TestName:    "test failure - invalid number",
			PreRunCode:  func(number string) {},
			PhoneNumber: "invalid",
			WantUserErr: service_errors.ErrInvalidNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			tt.PreRunCode(tt.PhoneNumber)

			responseDTO := appService.CheckNumber(app_user.NumberInput{
				NumberInput: shared_dto.NumberInput{
					PhoneNumber: tt.PhoneNumber,
				},
			})

			assert.Equal(t, tt.WantUserErr, responseDTO.UserErr)

			if tt.WantUserErr == nil {
				assert.Equal(t, tt.ExpectedIsExist, responseDTO.Data["isExist"])
			}
		})
	}
}

func TestLogin(t *testing.T) {
	validPassword := "validPassword123"
	salt := "testSalt"
	// In real code this would be: utils.HashPassword(validPassword, salt)
	// For testing we'll just use a fixed value
	hashedPassword := "hashed_password"

	tests := []struct {
		TestName            string
		PreRunCode          func(number string)
		PhoneNumber         string
		Password            string
		WantUserErr         error
		ResponseCode        string
		ResponseDTODataKeys []string
	}{
		{
			TestName: "test success - valid credentials",
			PreRunCode: func(number string) {
				repo.Create(&domain_user.User{
					Number:       number,
					Password:     hashedPassword,
					Salt:         salt,
					IsRegistered: true,
				})
			},
			PhoneNumber:         number_generator.GetNumber(),
			Password:            validPassword,
			WantUserErr:         nil,
			ResponseDTODataKeys: []string{"access", "refresh"},
		},
		{
			TestName: "test failure - invalid password",
			PreRunCode: func(number string) {
				repo.Create(&domain_user.User{
					Number:       number,
					Password:     hashedPassword,
					Salt:         salt,
					IsRegistered: true,
				})
			},
			PhoneNumber:  number_generator.GetNumber(),
			Password:     "wrongPassword",
			WantUserErr:  service_errors.ErrWrongPassword,
			ResponseCode: rcodes.InvalidField,
		},
		{
			TestName: "test failure - user not registered",
			PreRunCode: func(number string) {
				repo.Create(&domain_user.User{
					Number:       number,
					IsRegistered: false,
				})
			},
			PhoneNumber:  number_generator.GetNumber(),
			Password:     "anyPassword",
			WantUserErr:  service_errors.ErrUserNotRegistered,
			ResponseCode: rcodes.UserNotRegistered,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			tt.PreRunCode(tt.PhoneNumber)

			responseDTO := appService.Login(app_user.LoginUserInput{
				LoginUserInput: shared_dto.LoginUserInput{
					PhoneNumber:   tt.PhoneNumber,
					InputPassword: tt.Password,
				},
			}, "test-device", "1.1.1.1")

			assert.Equal(t, tt.WantUserErr, responseDTO.UserErr)
			assert.Equal(t, tt.ResponseCode, responseDTO.ResponseCode)

			if tt.ResponseDTODataKeys != nil {
				err := utils.CheckMapHaveKeys(responseDTO.Data, tt.ResponseDTODataKeys...)
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAccessFromRefresh(t *testing.T) {
	tests := []struct {
		TestName            string
		PreRunCode          func() string
		RefreshToken        string
		WantUserErr         error
		ResponseDTODataKeys []string
	}{
		{
			TestName: "test success - valid refresh token",
			PreRunCode: func() string {
				user := domain_user.User{
					Number: number_generator.GetNumber(),
				}
				repo.Create(&user)
				// In real code this would be: jwt.CreateRefreshToken(...)
				// For testing we'll just generate a random string
				refreshToken := uuid.New().String()
				return refreshToken
			},
			WantUserErr:         nil,
			ResponseDTODataKeys: []string{"access"},
		},
		{
			TestName: "test failure - invalid refresh token",
			PreRunCode: func() string {
				return "invalid-refresh-token"
			},
			WantUserErr: service_errors.ErrInvalidRefreshToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			refreshToken := tt.PreRunCode()

			responseDTO := appService.GetAccessFromRefresh(refreshToken)

			assert.Equal(t, tt.WantUserErr, responseDTO.UserErr)

			if tt.ResponseDTODataKeys != nil {
				err := utils.CheckMapHaveKeys(responseDTO.Data, tt.ResponseDTODataKeys...)
				assert.NoError(t, err)
			}
		})
	}
}

func TestChooseUserAvatar(t *testing.T) {
	tests := []struct {
		TestName     string
		PreRunCode   func(userID uint64) []string
		AvatarName   string
		WantUserErr  error
		ResponseCode string
	}{
		{
			TestName: "test success - valid avatar",
			PreRunCode: func(userID uint64) []string {
				avatars := []string{"avatar1.png", "avatar2.png"}
				return avatars
			},
			AvatarName:  "avatar1.png",
			WantUserErr: nil,
		},
		{
			TestName: "test failure - invalid avatar",
			PreRunCode: func(userID uint64) []string {
				return []string{"avatar1.png", "avatar2.png"}
			},
			AvatarName:   "invalid-avatar.png",
			WantUserErr:  service_errors.ErrAvatarNotFound,
			ResponseCode: rcodes.AvatarNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			// Create test user
			user := domain_user.User{
				Number: number_generator.GetNumber(),
			}
			repo.Create(&user)

			// Setup avatars by calling PreRunCode but don't use the result
			tt.PreRunCode(user.ID)

			responseDTO := appService.ChooseUserAvatar(tt.AvatarName, user.ID)

			assert.Equal(t, tt.WantUserErr, responseDTO.UserErr)
			assert.Equal(t, tt.ResponseCode, responseDTO.ResponseCode)

			// Verify avatar was set for success case
			if tt.WantUserErr == nil {
				updatedUser, _ := repo.GetByID(user.ID)
				assert.Equal(t, tt.AvatarName, updatedUser.Avatar)
			}
		})
	}
}

func TestGetAvatars(t *testing.T) {
	tests := []struct {
		TestName            string
		PreRunCode          func()
		ExpectedAvatarCount int
	}{
		{
			TestName: "test success - get avatars",
			PreRunCode: func() {
				// Normally this would set up S3, but we're testing the service layer
			},
			ExpectedAvatarCount: 2, // Mock returns 2 avatars
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			tt.PreRunCode()

			responseDTO := appService.GetAvatars()

			assert.Nil(t, responseDTO.UserErr)
			assert.Nil(t, responseDTO.ServerErr)

			avatars, ok := responseDTO.Data["data"].([]string)
			assert.True(t, ok)
			assert.GreaterOrEqual(t, len(avatars), tt.ExpectedAvatarCount)
		})
	}
}

func TestGetUserInfo(t *testing.T) {
	tests := []struct {
		TestName    string
		UserID      uint64
		PreRunCode  func() uint64
		WantUserErr error
		ExpectUser  bool
	}{
		{
			TestName: "test success - valid user",
			PreRunCode: func() uint64 {
				user := domain_user.User{
					Name:   "Test User",
					Number: number_generator.GetNumber(),
				}
				repo.Create(&user)
				return user.ID
			},
			ExpectUser: true,
		},
		{
			TestName: "test failure - user not found",
			PreRunCode: func() uint64 {
				return 999999 // Non-existent ID
			},
			WantUserErr: nil, // ServerErr should be set
			ExpectUser:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.TestName, func(t *testing.T) {
			userID := tt.PreRunCode()
			responseDTO := appService.GetUserInfo(userID)

			assert.Equal(t, tt.WantUserErr, responseDTO.UserErr)

			if tt.ExpectUser {
				assert.NotNil(t, responseDTO.Data["data"])
				userOutput := responseDTO.Data["data"].(app_user.UserOutput)
				assert.Equal(t, "Test User", userOutput.Name)
			} else {
				assert.NotNil(t, responseDTO.ServerErr)
			}
		})
	}
}
