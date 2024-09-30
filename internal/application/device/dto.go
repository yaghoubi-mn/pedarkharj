package app_device

import domain_device "github.com/yaghoubi-mn/pedarkharj/internal/domain/device"

type DeviceInput struct {
	Name         string `validate:"required,name"`
	IP           string `validate:"ipv4"`
	RefreshToken string `validate:"jwt"`
	UserID       uint64
}

func (d *DeviceInput) CreateDevice() (device domain_device.Device) {
	device.Name = d.Name
	device.LastIP = d.IP
	device.RefreshToken = d.RefreshToken
	device.UserID = d.UserID

	return
}
