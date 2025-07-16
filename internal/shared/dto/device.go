package shared_dto

type DeviceInput struct {
	Name         string `validate:"required,name"`
	IP           string `validate:"ipv4"`
	RefreshToken string `validate:"jwt"`
	UserID       uint64
}
