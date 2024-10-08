package domain_device

import domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"

type DeviceDomainRepository interface {
	Create(device Device) error
	Update(device Device) error
	CreateOrUpdate(device Device) error
	GetUserByRefreshToken(refresh string) (domain_user.User, error)
	Logout(userID uint64, deviceName string) error
	LogoutAllUserDevices(userID uint64) error
}
