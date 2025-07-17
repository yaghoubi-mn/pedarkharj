package app_device_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	app_device "github.com/yaghoubi-mn/pedarkharj/internal/application/device"
	domain_device "github.com/yaghoubi-mn/pedarkharj/internal/domain/device"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
)

// Mock implementations
type mockDeviceRepo struct {
	createOrUpdateErr     error
	getUserByRefreshToken domain_user.User
	getUserErr            error
	logoutErr             error
	logoutAllErr          error
}

func (m *mockDeviceRepo) Create(device domain_device.Device) error {
	return m.createOrUpdateErr
}

func (m *mockDeviceRepo) Update(device domain_device.Device) error {
	return m.createOrUpdateErr
}

func (m *mockDeviceRepo) CreateOrUpdate(device domain_device.Device) error {
	return m.createOrUpdateErr
}

func (m *mockDeviceRepo) GetUserByRefreshToken(refresh string) (domain_user.User, error) {
	return m.getUserByRefreshToken, m.getUserErr
}

func (m *mockDeviceRepo) Logout(userID uint64, deviceName string) error {
	return m.logoutErr
}

func (m *mockDeviceRepo) LogoutAllUserDevices(userID uint64) error {
	return m.logoutAllErr
}

type mockDeviceDomainService struct {
	createOrUpdateErr error
	logoutErr         error
	logoutAllErr      error
}

func (m *mockDeviceDomainService) Create(device *domain_device.Device) error {
	return m.createOrUpdateErr
}

func (m *mockDeviceDomainService) Update(device *domain_device.Device) error {
	return m.createOrUpdateErr
}

func (m *mockDeviceDomainService) CreateOrUpdate(device *domain_device.Device) error {
	return m.createOrUpdateErr
}

func (m *mockDeviceDomainService) Logout(userID uint64, deviceName string) error {
	return m.logoutErr
}

func (m *mockDeviceDomainService) LogoutAllUserDevices(userID uint64) error {
	return m.logoutAllErr
}

func TestCreateOrUpdate(t *testing.T) {
	tests := []struct {
		name          string
		repoErr       error
		serviceErr    error
		expectedError error
	}{
		{
			name:          "success",
			repoErr:       nil,
			serviceErr:    nil,
			expectedError: nil,
		},
		{
			name:          "domain service error",
			serviceErr:    errors.New("domain service error"),
			expectedError: errors.New("domain service error"),
		},
		{
			name:          "repository error",
			repoErr:       errors.New("repository error"),
			expectedError: errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockDeviceRepo{createOrUpdateErr: tt.repoErr}
			domainService := &mockDeviceDomainService{createOrUpdateErr: tt.serviceErr}
			service := app_device.NewDeviceAppService(repo, domainService)

			err := service.CreateOrUpdate(domain_device.DeviceInput{})
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestGetDeviceUserByRefreshToken(t *testing.T) {
	tests := []struct {
		name            string
		user            domain_user.User
		repoErr         error
		expectedUser    domain_user.User
		expectedErr     error
		expectedUserErr error
	}{
		{
			name:            "success",
			user:            domain_user.User{ID: 1, Name: "Test User"},
			expectedUser:    domain_user.User{ID: 1, Name: "Test User"},
			expectedErr:     nil,
			expectedUserErr: nil,
		},
		{
			name:            "repository error",
			repoErr:         errors.New("repository error"),
			expectedErr:     errors.New("repository error"),
			expectedUserErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockDeviceRepo{
				getUserByRefreshToken: tt.user,
				getUserErr:            tt.repoErr,
			}
			service := app_device.NewDeviceAppService(repo, &mockDeviceDomainService{})

			user, userErr, err := service.GetDeviceUserByRefreshToken("refresh-token")
			assert.Equal(t, tt.expectedUser, user)
			assert.Equal(t, tt.expectedUserErr, userErr)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name       string
		serviceErr error
		repoErr    error
		expectErr  bool
	}{
		{
			name:       "success",
			serviceErr: nil,
			repoErr:    nil,
			expectErr:  false,
		},
		{
			name:       "domain service error",
			serviceErr: errors.New("service error"),
			expectErr:  true,
		},
		{
			name:      "repository error",
			repoErr:   errors.New("repo error"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockDeviceRepo{logoutErr: tt.repoErr}
			domainService := &mockDeviceDomainService{logoutErr: tt.serviceErr}
			service := app_device.NewDeviceAppService(repo, domainService)

			response := service.Logout(1, "device-name")
			if tt.expectErr {
				assert.NotNil(t, response.ServerErr)
			} else {
				assert.Nil(t, response.ServerErr)
			}
		})
	}
}

func TestLogoutAllUserDevices(t *testing.T) {
	tests := []struct {
		name       string
		serviceErr error
		repoErr    error
		expectErr  bool
	}{
		{
			name:       "success",
			serviceErr: nil,
			repoErr:    nil,
			expectErr:  false,
		},
		{
			name:       "domain service error",
			serviceErr: errors.New("service error"),
			expectErr:  true,
		},
		{
			name:      "repository error",
			repoErr:   errors.New("repo error"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockDeviceRepo{logoutAllErr: tt.repoErr}
			domainService := &mockDeviceDomainService{logoutAllErr: tt.serviceErr}
			service := app_device.NewDeviceAppService(repo, domainService)

			response := service.LogoutAllUserDevices(1)
			if tt.expectErr {
				assert.NotNil(t, response.ServerErr)
			} else {
				assert.Nil(t, response.ServerErr)
			}
		})
	}
}
