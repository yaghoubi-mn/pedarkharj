package domain_device

type DeviceDomainRepository interface {
	Create(device Device) error
	Update(device Device) error
	CreateOrUpdate(device Device) error
}
