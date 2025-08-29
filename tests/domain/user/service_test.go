package user_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	"github.com/yaghoubi-mn/pedarkharj/pkg/service_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/utils"
	"github.com/yaghoubi-mn/pedarkharj/pkg/validator"
)

var userService domain_user.UserDomainService

func setup() {
	userService = domain_user.NewUserService(validator.NewValidator())

}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestSendOTP(t *testing.T) {

	tests := []struct {
		ID      int
		Number  string
		WantErr error
	}{
		{ // test success
			ID:      1,
			Number:  "+989123456789",
			WantErr: nil,
		},
		{ // test invalid number
			ID:      2,
			Number:  "+98912345678",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      3,
			Number:  "09123456789",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      4,
			Number:  "+99123456789",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      5,
			Number:  "+9893456789",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      6,
			Number:  "",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      7,
			Number:  "+",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      8,
			Number:  "+9891234567891",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      9,
			Number:  "+981234567899",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      10,
			Number:  "+989123a45678",
			WantErr: service_errors.ErrInvalidNumber,
		},
	}

	for _, tt := range tests {
		userErr := userService.SendOTP(domain_user.SendOTPInput{
			Number: tt.Number,
		})

		assert.Equal(t, tt.WantErr, userErr, tt)
	}
}

func TestVerifyOTP(t *testing.T) {

	tests := []struct {
		name    string
		number  string
		code    uint
		token   string
		mode    string
		mockErr error
		wantErr error
	}{
		{
			name:    "Success",
			number:  "+989123456789",
			code:    0,
			token:   uuid.NewString(),
			mode:    "signup",
			wantErr: nil,
		},
		{
			name:    "Success",
			number:  "+989123456789",
			code:    12345,
			token:   uuid.NewString(),
			mode:    "signup",
			wantErr: nil,
		},
		{
			name:    "Success",
			number:  "+989123456789",
			code:    0,
			token:   "",
			mode:    "signup",
			wantErr: nil,
		},
		{
			name:    "Blocked user",
			number:  "+989123456789",
			code:    12345,
			token:   uuid.NewString(),
			mode:    "signup",
			wantErr: service_errors.ErrBlockedUser,
		},
		{
			name:    "invalid number",
			number:  "+9891234567893",
			code:    0,
			token:   uuid.NewString(),
			mode:    "signup",
			wantErr: service_errors.ErrInvalidNumber,
		},
		{
			name:    "invalid mode",
			number:  "+989123456789",
			code:    12345,
			token:   uuid.NewString(),
			mode:    "a",
			wantErr: service_errors.ErrInvalidMode,
		},
	}

	for _, tt := range tests {
		userErr, serverErr := userService.VerifyOTP(domain_user.VerifyOTPInput{
			Number: tt.number,
			OTP:    tt.code,
			Token:  tt.token,
			Mode:   tt.mode,
		})
		assert.NoError(t, serverErr, tt)

		assert.Equal(t, tt.wantErr, userErr, tt)
	}
}

func TestSignup(t *testing.T) {

	tests := []struct {
		user    domain_user.User
		token   string
		wantErr error
	}{
		{
			user: domain_user.User{
				Name:     "success",
				Number:   "+989123456789",
				Password: "12345678",
			},
			token:   uuid.NewString(),
			wantErr: nil,
		},
		{
			user: domain_user.User{
				Name:     "invalid number",
				Number:   "+98912",
				Password: "12345678",
			},
			token:   uuid.NewString(),
			wantErr: service_errors.ErrInvalidNumber,
		},
		{
			user: domain_user.User{
				Name:     "long long long long long long name",
				Number:   "+989123456789",
				Password: "12345678",
			},
			token:   uuid.NewString(),
			wantErr: service_errors.ErrInvalidName,
		},
		{ // small name
			user: domain_user.User{
				Name:     "a",
				Number:   "+989123456789",
				Password: "12345678",
			},
			token:   uuid.NewString(),
			wantErr: service_errors.ErrSmallName,
		},
		{
			user: domain_user.User{
				Name:     "small password",
				Number:   "+989123456789",
				Password: "12345",
			},
			token:   uuid.NewString(),
			wantErr: service_errors.ErrSmallPassword,
		},
		{
			user: domain_user.User{
				Name:     "long password",
				Number:   "+989123456789",
				Password: "12345678901234567890012345678901",
			},
			token:   uuid.NewString(),
			wantErr: service_errors.ErrLongPassword,
		},
		{
			user: domain_user.User{
				Name:     "invalid token",
				Number:   "+989123456789",
				Password: "12345678",
			},
			token:   "dfdf",
			wantErr: service_errors.ErrInvalidToken,
		},
		{
			user: domain_user.User{
				Name:     "invalid token",
				Number:   "+989123456789",
				Password: "12345678",
			},
			token:   "",
			wantErr: service_errors.ErrInvalidToken,
		},
	}

	for _, tt := range tests {

		user, userErr, serverErr := userService.Signup(domain_user.SignupUserInput{
			Number:   tt.user.Number,
			Name:     tt.user.Name,
			Password: tt.user.Password,
			Token:    tt.token,
		})
		assert.NoError(t, serverErr, tt)

		assert.Equal(t, tt.wantErr, userErr, tt)

		if userErr == nil {

			assert.Equal(t, true, user.IsRegistered, "isRegistered must be true", tt)
			now := time.Now()
			if now.Unix()-user.RegisteredAt.Unix() < 10 {
				assert.Error(t, errors.New("invalid RegisteredAt"), "actual:", user.RegisteredAt)
			}
		}
	}
}

func TestResetPassword(t *testing.T) {

	tests := []struct {
		ID       int
		Number   string
		Password string
		Token    string
		WantErr  error
	}{
		{ // test success
			ID:       1,
			Number:   "+989123456789",
			Password: "12345678",
			Token:    uuid.NewString(),
			WantErr:  nil,
		},
		{ // test invalid number
			ID:       2,
			Number:   "+98123456789",
			Password: "12345678",
			Token:    uuid.NewString(),
			WantErr:  service_errors.ErrInvalidNumber,
		},
		{ // invalid number
			ID:       3,
			Number:   "09123456789",
			Password: "12345678",
			Token:    uuid.NewString(),
			WantErr:  service_errors.ErrInvalidNumber,
		},
		{ // small password
			ID:       4,
			Number:   "+989123456789",
			Password: "1234",
			Token:    uuid.NewString(),
			WantErr:  service_errors.ErrSmallPassword,
		},
		{
			ID:       5,
			Number:   "+989123456789",
			Password: "1234567890123456789012345678900",
			Token:    uuid.NewString(),
			WantErr:  service_errors.ErrLongPassword,
		},
		{ // invalid token
			ID:       6,
			Number:   "+989123456789",
			Password: "12345678",
			Token:    "121",
			WantErr:  service_errors.ErrInvalidToken,
		},
		{ // invalid token
			ID:       6,
			Number:   "+989123456789",
			Password: "12345678",
			Token:    "",
			WantErr:  service_errors.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		userErr, serverErr, salt, outPassword := userService.ResetPassword(domain_user.RestPasswordInput{
			Number:   tt.Number,
			Password: tt.Password,
			Token:    tt.Token,
		})

		assert.NoError(t, serverErr, tt)

		assert.Equal(t, tt.WantErr, userErr, tt)

		assert.NotEqual(t, 0, len(salt), tt)
		assert.NotEqual(t, 0, len(outPassword), tt)

	}
}

func TestCheckNumber(t *testing.T) {

	tests := []struct {
		ID      int
		Number  string
		WantErr error
	}{
		{ // test success
			ID:      1,
			Number:  "+989123456789",
			WantErr: nil,
		},
		{ // test invalid number
			ID:      2,
			Number:  "+98912345678",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      3,
			Number:  "09123456789",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      4,
			Number:  "+99123456789",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      5,
			Number:  "+9893456789",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      6,
			Number:  "",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      7,
			Number:  "+",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      8,
			Number:  "+9891234567891",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      9,
			Number:  "+981234567899",
			WantErr: service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:      10,
			Number:  "+989123a45678",
			WantErr: service_errors.ErrInvalidNumber,
		},
	}

	for _, tt := range tests {
		userErr := userService.CheckNumber(tt.Number)

		assert.Equal(t, tt.WantErr, userErr, tt)
	}
}

func TestLogin(t *testing.T) {
	salt, err := utils.GenerateRandomSalt()
	assert.NoError(t, err)

	getPasswordHash := func(password string) string {
		hashedPassword, err := utils.HashPasswordWithSalt(password, salt, 4)
		assert.NoError(t, err)

		return hashedPassword
	}

	tests := []struct {
		ID             int
		Number         string
		Password       string
		HashedPassword string
		IsRegistered   bool
		IsBlocked      bool
		WantErr        error
	}{
		{ // test success
			ID:             1,
			Number:         "+989123456789",
			Password:       "12345678",
			HashedPassword: getPasswordHash("12345678"),
			IsRegistered:   true,
			IsBlocked:      false,
			WantErr:        nil,
		},
		{ // test invalid number
			ID:             2,
			Number:         "+98912345678",
			Password:       "12345678",
			HashedPassword: getPasswordHash("12345678"),
			IsRegistered:   true,
			IsBlocked:      false,
			WantErr:        service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:             3,
			Number:         "09123456789",
			Password:       "12345678",
			HashedPassword: getPasswordHash("12345678"),
			IsRegistered:   true,
			IsBlocked:      false,
			WantErr:        service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:             4,
			Number:         "+99123456789",
			Password:       "12345678",
			HashedPassword: getPasswordHash("12345678"),
			IsRegistered:   true,
			IsBlocked:      false,
			WantErr:        service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:             5,
			Number:         "+9893456789",
			Password:       "12345678",
			HashedPassword: getPasswordHash("12345678"),
			IsRegistered:   true,
			IsBlocked:      false,
			WantErr:        service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:             6,
			Number:         "",
			Password:       "12345678",
			HashedPassword: getPasswordHash("12345678"),
			IsRegistered:   true,
			IsBlocked:      false,
			WantErr:        service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:             7,
			Number:         "+",
			Password:       "12345678",
			HashedPassword: getPasswordHash("12345678"),
			IsRegistered:   true,
			IsBlocked:      false,
			WantErr:        service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:             8,
			Number:         "+9891234567891",
			Password:       "12345678",
			HashedPassword: getPasswordHash("12345678"),
			IsRegistered:   true,
			IsBlocked:      false,
			WantErr:        service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:             9,
			Number:         "+981234567899",
			Password:       "12345678",
			HashedPassword: getPasswordHash("12345678"),
			IsRegistered:   true,
			IsBlocked:      false,
			WantErr:        service_errors.ErrInvalidNumber,
		},
		{ // test invalid number
			ID:             10,
			Number:         "+989123a45678",
			Password:       "12345678",
			HashedPassword: getPasswordHash("12345678"),
			IsRegistered:   true,
			IsBlocked:      false,
			WantErr:        service_errors.ErrInvalidNumber,
		},
		{ // test wrong password
			ID:             11,
			Number:         "+989123456789",
			Password:       "12345678",
			HashedPassword: getPasswordHash("12345678d"),
			IsRegistered:   true,
			IsBlocked:      false,
			WantErr:        service_errors.ErrWrongPassword,
		},
		{ // test user not registered
			ID:             11,
			Number:         "+989123456789",
			Password:       "12345678",
			HashedPassword: getPasswordHash("12345678"),
			IsRegistered:   false,
			IsBlocked:      false,
			WantErr:        service_errors.ErrUserNotRegistered,
		},
		{ // test blocked user
			ID:             11,
			Number:         "+989123456789",
			Password:       "12345678",
			HashedPassword: getPasswordHash("12345678"),
			IsRegistered:   true,
			IsBlocked:      true,
			WantErr:        service_errors.ErrBlockedUser,
		},
	}

	for _, tt := range tests {

		userErr, serverErr := userService.Login(domain_user.LoginUserInput{
			Number:        tt.Number,
			InputPassword: tt.Password,
			RealPassword:  tt.HashedPassword,
			Salt:          salt,
			IsRegistered:  tt.IsRegistered,
			IsBlocked:     tt.IsBlocked,
		})

		assert.NoError(t, serverErr)
		assert.Equal(t, tt.WantErr, userErr, tt)
	}
}
