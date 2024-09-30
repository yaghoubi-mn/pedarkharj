package repository

import (
	"github.com/yaghoubi-mn/pedarkharj/internal/domain/device"
	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"gorm.io/gorm"
)

type GormDeviceRepository struct {
	DB *gorm.DB
}

func NewGormDeviceRepository(db *gorm.DB) domain_device.DeviceDomainRepository {
	return &GormDeviceRepository{DB: db}
}

func (repo *GormDeviceRepository) Create(device domain_device.Device) error {

	if err := repo.DB.Create(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return database_errors.ErrRecordNotFound
		}

		return err
	}

	return nil
}

func (repo *GormDeviceRepository) Update(device domain_device.Device) error {

	if err := repo.DB.Updates(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return database_errors.ErrRecordNotFound
		}

		return err
	}

	return nil
}

// distinguish devices with device.Name and device.UserID
// device.ID can be zero
func (repo *GormDeviceRepository) CreateOrUpdate(device domain_device.Device) error {

	// check device exist or not
	var d domain_device.Device
	if err := repo.DB.First(&d, &domain_device.Device{Name: device.Name, UserID: device.UserID}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {

			// device not exist. insert device
			if err = repo.DB.Create(&device).Error; err != nil {
				return err
			}

			return nil
		}

		return err
	}

	// device found. update it
	device.ID = d.ID
	if err := repo.DB.Updates(&device).Error; err != nil {
		return err
	}

	return nil
}
