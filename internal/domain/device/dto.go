package domain_device

import shared_dto "github.com/yaghoubi-mn/pedarkharj/internal/shared/dto"

type DeviceInput struct {
	shared_dto.DeviceInput
}

func NewDeviceInput(name, ip, refreshToken string, userID uint64) DeviceInput {
	return DeviceInput{
		shared_dto.DeviceInput{
			Name:         name,
			IP:           ip,
			RefreshToken: refreshToken,
			UserID:       userID,
		},
	}
}

func (d *DeviceInput) CreateDevice() (device Device) {
	device.Name = d.Name
	device.LastIP = d.IP
	device.RefreshToken = d.RefreshToken
	device.UserID = d.UserID

	return
}
