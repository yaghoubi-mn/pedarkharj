package device

import (
	"time"

	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
)

type DeviceService interface {
	Create(deviceInput DeviceInput)
	CreateOrUpdate(DeviceInput DeviceInput)
}

type service struct {
	repo      DeviceRepository
	validator datatypes.Validator
}

// func NewDeviceService(repo DeviceRepository, validator datatypes.Validator) DeviceService {
// 	return &service{
// 		repo:      repo,
// 		validator: validator,
// 	}
// }

func (s *service) Create(deviceInput DeviceInput) (map[string]string, error) {
	errMap := s.validator.Struct(deviceInput)
	if errMap != nil {
		return errMap, nil
	}

	var device Device
	device.Name = deviceInput.Name
	device.LastIP = deviceInput.IP
	device.FirstLogin = time.Now()
	device.LastLogin = time.Now()
	device.RefreshToken = deviceInput.RefreshToken

	err := s.repo.Create(device)

	return nil, err

}
