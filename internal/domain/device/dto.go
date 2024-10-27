package domain_device

type DeviceInput struct {
	Name         string `validate:"required,name"`
	IP           string `validate:"ipv4"`
	RefreshToken string `validate:"jwt"`
	UserID       uint64
}

func (d *DeviceInput) CreateDevice() (device Device) {
	device.Name = d.Name
	device.LastIP = d.IP
	device.RefreshToken = d.RefreshToken
	device.UserID = d.UserID

	return
}
