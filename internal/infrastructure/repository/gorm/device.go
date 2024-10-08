package repository

import (
	"github.com/yaghoubi-mn/pedarkharj/internal/domain/device"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
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

func (repo *GormDeviceRepository) GetUserByRefreshToken(refresh string) (user domain_user.User, err error) {

	// // TODO: fix preload
	// var device domain_device.Device
	// if err = repo.DB.Preload("User").Where(domain_device.Device{RefreshToken: refresh}).Find(&device).Error; err != nil {
	// 	if err == gorm.ErrRecordNotFound {
	// 		return user, database_errors.ErrRecordNotFound
	// 	}

	// 	return user, err
	// }

	var userID uint64
	if err = repo.DB.Model(&domain_device.Device{}).Select("user_id").Where(domain_device.Device{RefreshToken: refresh}).Find(&userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, database_errors.ErrRecordNotFound
		}

		return user, err
	}

	if err = repo.DB.First(&user, &domain_user.User{ID: userID}).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (repo *GormDeviceRepository) Logout(userID uint64, deviceName string) error {

	if err := repo.DB.Model(&domain_device.Device{}).Where(domain_device.Device{UserID: userID, Name: deviceName}).Update("refresh_token", "").Error; err != nil {
		return err
	}

	return nil
}

func (repo *GormDeviceRepository) LogoutAllUserDevices(userID uint64) error {

	if err := repo.DB.Model(&domain_device.Device{}).Where(domain_device.Device{UserID: userID}).Update("refresh_token", "").Error; err != nil {
		return err
	}

	return nil
}
