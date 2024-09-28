package device

type DeviceInput struct {
	Name         string `validate:"required,name"`
	IP           string `validate:"ipv4"`
	RefreshToken string `validate:"jwt"`
}
