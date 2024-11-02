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
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	"github.com/yaghoubi-mn/pedarkharj/internal/infrastructure/config"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/rcodes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/service_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/utils"
	"github.com/yaghoubi-mn/pedarkharj/tests/pkg/helpers"
	_ "github.com/yaghoubi-mn/pedarkharj/tests/pkg/helpers"
	"github.com/yaghoubi-mn/pedarkharj/tests/pkg/number_generator"
)

var appService app_user.UserAppService
var repo domain_user.UserDomainRepository
var cacheRepo datatypes.CacheRepository

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
		Number              string
		WantUserErr         error
		ResponseCode        string
		ResponseDTODataKeys []string
		NextRunCode         func(number string)
	}{
		{ // test sucess
			PreRunCode:          func(number string) {},
			Number:              number_generator.GetNumber(),
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
					Number: number,
				})

				assert.NoError(t, responseDTO.ServerErr)
				assert.NoError(t, responseDTO.UserErr)
				assert.Equal(t, rcodes.CodeSendToNumber, responseDTO.ResponseCode, responseDTO)
			},
			Number:              number_generator.GetNumber(),
			WantUserErr:         service_errors.ErrOTPNotExpired,
			ResponseCode:        rcodes.NumberDelay,
			ResponseDTODataKeys: []string{"delayTimeSeconds"},
			NextRunCode:         func(number string) {},
		},
		{ // test invalid number
			PreRunCode:          func(number string) {},
			Number:              "+9891234567",
			WantUserErr:         service_errors.ErrInvalidNumber,
			ResponseCode:        rcodes.InvalidField,
			ResponseDTODataKeys: []string{},
			NextRunCode:         func(number string) {},
		},
	}

	for _, tt := range tests {

		tt.PreRunCode(tt.Number)

		responseDTO := appService.SendOTP(app_user.SendOTPInput{
			Number: tt.Number,
		})

		assert.NoError(t, responseDTO.ServerErr)
		assert.Equal(t, tt.WantUserErr, responseDTO.UserErr)

		assert.Equal(t, tt.ResponseCode, responseDTO.ResponseCode, responseDTO)
		err := utils.CheckMapHaveKeys(responseDTO.Data, tt.ResponseDTODataKeys...)
		assert.NoError(t, err, responseDTO.Data)

		tt.NextRunCode(tt.Number)

	}

}

func TestVerifyOTP(t *testing.T) {

	tests := []struct {
		TestID              int
		PreRunCode          func(number string)
		Number              string
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
			PreRunCode:          func(number string) {},
			Number:              number_generator.GetNumber(),
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
			Number:              number_generator.GetNumber(),
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
			Number:              number_generator.GetNumber(),
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
			Number:              number_generator.GetNumber(),
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
			Number:              number_generator.GetNumber(),
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
			Number:              number_generator.GetNumber(),
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
			Number:              number_generator.GetNumber(),
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
			Number:              number_generator.GetNumber(),
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

		tt.PreRunCode(tt.Number)

		// call send otp
		responseDTO := appService.SendOTP(app_user.SendOTPInput{
			Number: tt.Number,
		})
		assert.NoError(t, responseDTO.ServerErr)
		assert.NoError(t, responseDTO.UserErr)

		// get otp from cache
		data, _, err := cacheRepo.Get(tt.Number)
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
			Number: tt.Number,
			OTP:    uint(otp),
			Token:  token,
			Mode:   tt.Mode,
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

		tt.NextRunCode(tt.Number)

	}

}
