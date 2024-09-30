package app_device

import (
	domain_device "github.com/yaghoubi-mn/pedarkharj/internal/domain/device"
)

type DeviceAppService interface {
	CreateOrUpdate(deviceInput DeviceInput) error
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
