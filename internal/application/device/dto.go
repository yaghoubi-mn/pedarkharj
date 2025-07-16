package app_device

import (
	domain_device "github.com/yaghoubi-mn/pedarkharj/internal/domain/device"
	shared_dto "github.com/yaghoubi-mn/pedarkharj/internal/shared/dto"
)

type DeviceInput struct {
	shared_dto.DeviceInput
}

func (d *DeviceInput) CreateDevice() (device domain_device.Device) {
	device.Name = d.Name
	device.LastIP = d.IP
	device.RefreshToken = d.RefreshToken
	device.UserID = d.UserID

	return
}
