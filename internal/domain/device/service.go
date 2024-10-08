package domain_device

import (
	"errors"
	"time"

	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
)

type DeviceDomainService interface {
	Create(device *Device) error
	CreateOrUpdate(device *Device) error
	Logout(userID uint64, deviceName string) error
	LogoutAllUserDevices(userID uint64) error
}

type service struct {
	validator datatypes.Validator
}

func NewDeviceService(validator datatypes.Validator) DeviceDomainService {
	return &service{
		validator: validator,
	}
}

func (s *service) Create(device *Device) error {
	if err := s.validator.ValidateFieldByFieldName("Name", device.Name, Device{}); err != nil {
		return errors.New("name: " + err.Error())
	}

	if err := s.validator.ValidateFieldByFieldName("LastIP", device.LastIP, Device{}); err != nil {
		return errors.New("lastIP: " + err.Error())
	}

	if err := s.validator.ValidateFieldByFieldName("RefreshToken", device.RefreshToken, Device{}); err != nil {
		return errors.New("refresh: invalid refresh token: " + err.Error())
	}

	if device.UserID == 0 {
		return errors.New("invalid user id")
	}

	device.LastLogin = time.Now()
	device.FirstLogin = time.Now()

	return nil

}

func (s *service) CreateOrUpdate(device *Device) error {
	if err := s.validator.ValidateFieldByFieldName("Name", device.Name, Device{}); err != nil {
		return errors.New("name: " + err.Error())
	}

	if err := s.validator.ValidateFieldByFieldName("LastIP", device.LastIP, Device{}); err != nil {
		return errors.New("lastIP: " + err.Error())
	}

	if err := s.validator.ValidateFieldByFieldName("RefreshToken", device.RefreshToken, Device{}); err != nil {
		return errors.New("refresh: invalid refresh token: " + err.Error())
	}

	device.LastLogin = time.Now()

	// for create device
	if device.ID == 0 {
		device.FirstLogin = time.Now()
	}

	return nil
}

func (s *service) Logout(userID uint64, deviceName string) error {
	if err := s.validator.ValidateFieldByFieldName("Name", deviceName, Device{}); err != nil {
		return err
	}

	if userID == 0 {
		return errors.New("invalid userID")
	}

	return nil
}

func (s *service) LogoutAllUserDevices(userID uint64) error {

	if userID == 0 {
		return errors.New("ivnalid userID")
	}

	return nil
}
