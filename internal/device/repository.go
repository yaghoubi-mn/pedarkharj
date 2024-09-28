package device

import (
	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	Create(device Device) error
	Update(device Device) error
}

type GormDeviceRepository struct {
	DB *gorm.DB
}

func NewGormDeviceRepository(db *gorm.DB) DeviceRepository {
	return &GormDeviceRepository{DB: db}
}

func (repo *GormDeviceRepository) Create(device Device) error {

	if err := repo.DB.Create(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return database_errors.ErrRecordNotFound
		}

		return err
	}

	return nil
}

func (repo *GormDeviceRepository) Update(device Device) error {

	if err := repo.DB.Updates(&device).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return database_errors.ErrRecordNotFound
		}

		return err
	}

	return nil
}
