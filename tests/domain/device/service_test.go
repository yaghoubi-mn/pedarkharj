package device_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	domain_device "github.com/yaghoubi-mn/pedarkharj/internal/domain/device"
	"github.com/yaghoubi-mn/pedarkharj/pkg/jwt"
	"github.com/yaghoubi-mn/pedarkharj/pkg/service_errors"
	"github.com/yaghoubi-mn/pedarkharj/pkg/validator"
)

var deviceService domain_device.DeviceDomainService

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	validator := validator.NewValidator()
	deviceService = domain_device.NewDeviceService(validator)
}

func TestCreate(t *testing.T) {

	jwtToken, err := jwt.CreateJwt(map[string]any{"exp": time.Now()})
	assert.NoError(t, err)

	tests := []struct {
		TestID int
		domain_device.Device
		WantErr error
	}{
		{ // test success
			TestID: 1,
			Device: domain_device.Device{
				ID:           0,
				Name:         "test",
				LastIP:       "127.0.0.1",
				RefreshToken: jwtToken,
				UserID:       1,
			},
			WantErr: nil,
		},
		{ // test invalid UserID
			TestID: 2,
			Device: domain_device.Device{
				ID:           0,
				Name:         "test",
				LastIP:       "127.0.0.1",
				RefreshToken: jwtToken,
				UserID:       0,
			},
			WantErr: service_errors.ErrInvalidUserID,
		},
		{ // test invalid device name
			TestID: 3,
			Device: domain_device.Device{
				ID:           0,
				Name:         "",
				LastIP:       "127.0.0.1",
				RefreshToken: jwtToken,
				UserID:       1,
			},
			WantErr: service_errors.ErrInvalidUserAgent,
		},
		{ // invalid device ip
			TestID: 4,
			Device: domain_device.Device{
				ID:           0,
				Name:         "test",
				LastIP:       "127.0.0",
				RefreshToken: jwtToken,
				UserID:       1,
			},
			WantErr: service_errors.ErrInvalidIP,
		},
		{ // test invalid ip
			TestID: 5,
			Device: domain_device.Device{
				ID:           0,
				Name:         "test",
				LastIP:       "",
				RefreshToken: jwtToken,
				UserID:       1,
			},
			WantErr: service_errors.ErrInvalidIP,
		},
		{ // test invalid refresh token
			TestID: 6,
			Device: domain_device.Device{
				ID:           0,
				Name:         "test",
				LastIP:       "127.0.0.1",
				RefreshToken: "eydkfdkf",
				UserID:       1,
			},
			WantErr: service_errors.ErrInvalidRefreshToken,
		},
		{ // test invalid refresh token
			TestID: 1,
			Device: domain_device.Device{
				ID:           0,
				Name:         "test",
				LastIP:       "127.0.0.1",
				RefreshToken: "",
				UserID:       1,
			},
			WantErr: service_errors.ErrInvalidRefreshToken,
		},
	}

	for _, tt := range tests {

		err := deviceService.Create(&tt.Device)

		assert.Equal(t, tt.WantErr, err, tt)

		assert.GreaterOrEqual(t, time.Now().Unix()+5, tt.Device.FirstLogin.Unix(), tt)
		assert.GreaterOrEqual(t, time.Now().Unix()+5, tt.Device.LastLogin.Unix(), tt)
	}
}

func TestCreateOrUpdate(t *testing.T) {

	jwtToken, err := jwt.CreateJwt(map[string]any{"exp": time.Now()})
	assert.NoError(t, err)

	tests := []struct {
		TestID int
		domain_device.Device
		WantErr error
	}{
		{ // test success
			TestID: 1,
			Device: domain_device.Device{
				ID:           0,
				Name:         "test",
				LastIP:       "127.0.0.1",
				RefreshToken: jwtToken,
				UserID:       1,
			},
			WantErr: nil,
		},
		{ // test invalid UserID
			TestID: 2,
			Device: domain_device.Device{
				ID:           0,
				Name:         "test",
				LastIP:       "127.0.0.1",
				RefreshToken: jwtToken,
				UserID:       0,
			},
			WantErr: service_errors.ErrInvalidUserID,
		},
		{ // test invalid device name
			TestID: 3,
			Device: domain_device.Device{
				ID:           0,
				Name:         "",
				LastIP:       "127.0.0.1",
				RefreshToken: jwtToken,
				UserID:       1,
			},
			WantErr: service_errors.ErrInvalidUserAgent,
		},
		{ // invalid device ip
			TestID: 4,
			Device: domain_device.Device{
				ID:           0,
				Name:         "test",
				LastIP:       "127.0.0",
				RefreshToken: jwtToken,
				UserID:       1,
			},
			WantErr: service_errors.ErrInvalidIP,
		},
		{ // test invalid ip
			TestID: 5,
			Device: domain_device.Device{
				ID:           0,
				Name:         "test",
				LastIP:       "",
				RefreshToken: jwtToken,
				UserID:       1,
			},
			WantErr: service_errors.ErrInvalidIP,
		},
		{ // test invalid refresh token
			TestID: 6,
			Device: domain_device.Device{
				ID:           0,
				Name:         "test",
				LastIP:       "127.0.0.1",
				RefreshToken: "eydkfdkf",
				UserID:       1,
			},
			WantErr: service_errors.ErrInvalidRefreshToken,
		},
		{ // test invalid refresh token
			TestID: 1,
			Device: domain_device.Device{
				ID:           0,
				Name:         "test",
				LastIP:       "127.0.0.1",
				RefreshToken: "",
				UserID:       1,
			},
			WantErr: service_errors.ErrInvalidRefreshToken,
		},
		{ // test devicea already exist
			TestID: 1,
			Device: domain_device.Device{
				ID:           1,
				Name:         "test",
				LastIP:       "127.0.0.1",
				RefreshToken: "",
				UserID:       1,
			},
			WantErr: service_errors.ErrInvalidRefreshToken,
		},
	}

	for _, tt := range tests {

		err := deviceService.Create(&tt.Device)

		assert.Equal(t, tt.WantErr, err, tt)

		assert.GreaterOrEqual(t, time.Now().Unix()+5, tt.Device.LastLogin.Unix(), tt)
		if tt.Device.ID == 0 {
			assert.GreaterOrEqual(t, time.Now().Unix()+5, tt.Device.FirstLogin.Unix(), tt)
		}
	}
}

func TestLogout(t *testing.T) {

	tests := []struct {
		TestID     int
		UserID     uint64
		DeviceName string
		WantErr    error
	}{
		{ // test success
			TestID:     1,
			UserID:     1,
			DeviceName: "test",
			WantErr:    nil,
		},
		{ // test invalid user id
			TestID:     2,
			UserID:     0,
			DeviceName: "test",
			WantErr:    service_errors.ErrInvalidUserID,
		},
		{ // test invalid device name (user agent)
			TestID:     3,
			UserID:     1,
			DeviceName: "test'",
			WantErr:    service_errors.ErrInvalidUserAgent,
		},
	}

	for _, tt := range tests {

		err := deviceService.Logout(tt.UserID, tt.DeviceName)

		assert.Equal(t, tt.WantErr, err, tt)
	}
}

func TestLogoutAllUserDevices(t *testing.T) {

	tests := []struct {
		TestID  int
		UserID  uint64
		WantErr error
	}{
		{ // test success
			TestID:  1,
			UserID:  1,
			WantErr: nil,
		},
		{ // test invalid user id
			TestID:  2,
			UserID:  0,
			WantErr: service_errors.ErrInvalidUserID,
		},
	}

	for _, tt := range tests {

		err := deviceService.LogoutAllUserDevices(tt.UserID)

		assert.Equal(t, tt.WantErr, err, tt)
	}
}
