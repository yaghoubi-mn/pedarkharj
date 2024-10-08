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
	err := s.validator.ValidateField(device.Name, "useragent,size:300,required")
	if err != nil {
		return errors.New("name: " + err.Error())
	}

	err = s.validator.ValidateField(device.LastIP, "size:15,required")
	if err != nil {
		return errors.New("lastIP: " + err.Error())
	}

	err = s.validator.ValidateField(device.RefreshToken, "size:200,required")
	if err != nil {
		return errors.New("refresh: invalid refresh token: " + err.Error())
	}

	if device.UserID == 0 {
		return errors.New("invalid user id")
	}

	device.LastLogin = time.Now()
	device.FirstLogin = time.Now()

	return err

}

func (s *service) CreateOrUpdate(device *Device) error {
	err := s.validator.ValidateField(device.Name, "useragent,size:300,required")
	if err != nil {
		return errors.New("name: " + err.Error())
	}

	err = s.validator.ValidateField(device.LastIP, "size:15,required")
	if err != nil {
		return errors.New("lastIP: " + err.Error())
	}

	err = s.validator.ValidateField(device.RefreshToken, "size:200,required")
	if err != nil {
		return errors.New("refresh: invalid refresh token: " + err.Error())
	}

	device.LastLogin = time.Now()

	// for create device
	if device.ID == 0 {
		device.FirstLogin = time.Now()
	}

	return err
}

func (s *service) Logout(userID uint64, deviceName string) error {
	err := s.validator.ValidateField(deviceName, "useragent")
	if err != nil {
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
