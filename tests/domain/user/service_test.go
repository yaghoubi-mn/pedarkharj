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

func TestVerifyNumber(t *testing.T) {

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
		userErr, serverErr := userService.VerifyNumber(domain_user.VerifyNumberInput{
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
		Number  string
		WantErr error
	}{
		{ // test success
			Number:  "+989123456789",
			WantErr: nil,
		},
		{ // test
			Number:  "+989123456789",
			WantErr: nil,
		},
		{ // test success
			Number:  "+989123456789",
			WantErr: nil,
		},
		{ // test success
			Number:  "+989123456789",
			WantErr: nil,
		},
		{ // test success
			Number:  "+989123456789",
			WantErr: nil,
		},
		{ // test success
			Number:  "+989123456789",
			WantErr: nil,
		},
		{ // test success
			Number:  "+989123456789",
			WantErr: nil,
		},
	}
}
