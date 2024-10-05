package user_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	"github.com/yaghoubi-mn/pedarkharj/pkg/service_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/validator"
)

func TestVerifyNumber(t *testing.T) {
	userService := domain_user.NewUserService(validator.NewValidator())

	tests := []struct {
		name      string
		number    string
		code      uint
		token     string
		isBlocked bool
		mockErr   error
		wantErr   error
	}{
		{
			name:      "Success",
			number:    "+989123456789",
			code:      0,
			token:     uuid.NewString(),
			isBlocked: false,
			wantErr:   nil,
		},
		{
			name:      "Success",
			number:    "+989123456789",
			code:      12345,
			token:     uuid.NewString(),
			isBlocked: false,
			wantErr:   nil,
		},
		{
			name:      "Success",
			number:    "+989123456789",
			code:      0,
			token:     "",
			isBlocked: false,
			wantErr:   nil,
		},
		{
			name:      "Blocked user",
			number:    "+989123456789",
			code:      12345,
			token:     uuid.NewString(),
			isBlocked: true,
			wantErr:   service_errors.ErrBlockedUser,
		},
		{
			name:      "invalid number",
			number:    "+9891234567893",
			code:      0,
			token:     uuid.NewString(),
			isBlocked: false,
			wantErr:   service_errors.ErrInvalidNumber,
		},
	}

	for _, tt := range tests {
		userErr, serverErr := userService.VerifyNumber(tt.number, tt.code, tt.token, tt.isBlocked)
		assert.NoError(t, serverErr, tt)

		assert.Equal(t, tt.wantErr, userErr, tt)
	}
}

func TestSignup(t *testing.T) {
	userService := domain_user.NewUserService(validator.NewValidator())

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
				Name:     "long long long long long name",
				Number:   "+989123456789",
				Password: "12345678",
			},
			token:   uuid.NewString(),
			wantErr: service_errors.ErrLongName,
		},
		{ // small name
			user: domain_user.User{
				Name:     "a",
				Number:   "+98123456789",
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
				Password: "12345678901234567",
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
		userErr, serverErr := userService.Signup(&tt.user, tt.token)
		assert.NoError(t, serverErr, tt)

		assert.Equal(t, tt.wantErr, userErr, tt)

		if userErr == nil {

			assert.Equal(t, true, tt.user.IsRegistered, tt)
			now := time.Now()
			if now.Unix()-tt.user.RegisteredAt.Unix() < 10 {
				assert.Error(t, errors.New("invalid RegisteredAt"), "actual:", tt.user.RegisteredAt)
			}
		}
	}
}
