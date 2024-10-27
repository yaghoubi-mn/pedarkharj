package app_device

import (
	"errors"

	domain_device "github.com/yaghoubi-mn/pedarkharj/internal/domain/device"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/jwt"
)

type DeviceAppService interface {
	CreateOrUpdate(deviceInput domain_device.DeviceInput) error
	GetDeviceUserByRefreshToken(refresh string) (user domain_user.User, userErr error, serverErr error)
	Logout(userID uint64, deviceName string) datatypes.ResponseDTO
	LogoutAllUserDevices(userID uint64) datatypes.ResponseDTO
}

type service struct {
	repo          domain_device.DeviceDomainRepository
	domainService domain_device.DeviceDomainService
}

func NewDeviceAppService(repo domain_device.DeviceDomainRepository, domainService domain_device.DeviceDomainService) DeviceAppService {
	return &service{
		repo:          repo,
		domainService: domainService,
	}
}

func (s *service) CreateOrUpdate(deviceInput domain_device.DeviceInput) error {

	device := deviceInput.CreateDevice()

	err := s.domainService.CreateOrUpdate(&device)
	if err != nil {
		return err
	}

	err = s.repo.CreateOrUpdate(device)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetDeviceUserByRefreshToken(refresh string) (user domain_user.User, userErr error, err error) {

	_, err = jwt.VerifyJwt(refresh)
	if err != nil {
		return user, errors.New("refresh: invalid token"), nil
	}
	user, err = s.repo.GetUserByRefreshToken(refresh)
	return user, nil, err

}

func (s *service) Logout(userID uint64, deviceName string) (responseDTO datatypes.ResponseDTO) {
	responseDTO.Data = make(map[string]any)

	err := s.domainService.Logout(userID, deviceName)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}

	err = s.repo.Logout(userID, deviceName)
	responseDTO.ServerErr = err

	return

}

func (s *service) LogoutAllUserDevices(userID uint64) (responseDTO datatypes.ResponseDTO) {
	responseDTO.Data = make(map[string]any)

	err := s.domainService.LogoutAllUserDevices(userID)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}

	err = s.repo.LogoutAllUserDevices(userID)
	responseDTO.ServerErr = err
	return
}
