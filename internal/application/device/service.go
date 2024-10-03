package app_device

import (
	"errors"

	domain_device "github.com/yaghoubi-mn/pedarkharj/internal/domain/device"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	"github.com/yaghoubi-mn/pedarkharj/pkg/jwt"
)

type DeviceAppService interface {
	CreateOrUpdate(deviceInput DeviceInput) error
	GetDeviceUserByRefreshToken(refresh string) (user domain_user.User, userErr error, serverErr error)
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

func (s *service) CreateOrUpdate(deviceInput DeviceInput) error {

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
