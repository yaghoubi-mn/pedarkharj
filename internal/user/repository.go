package user

import (
	"errors"

	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetByID(id uint64) (User, error)
	GetByNumber(number string) (User, error)
	Create(user User) error
	Update(user User) error
	Delete(id uint64) error
}

type GormUserRepository struct {
	DB *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{DB: db}
}

func (repo *GormUserRepository) GetByID(id uint64) (User, error) {
	var user User
	if err := repo.DB.Where(User{ID: id}).Find(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, errors.New("Record Not Found")
		}

		return user, err
	}

	return user, nil
}

func (repo *GormUserRepository) GetByNumber(number string) (User, error) {
	var user User
	if err := repo.DB.First(&user, User{Number: number}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, database_errors.ErrRecordNotFound
		}

		return user, err
	}

	return user, nil
}

func (repo *GormUserRepository) Create(user User) error {

	if err := repo.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (repo *GormUserRepository) Update(user User) error {

	if err := repo.DB.Updates(&user).Error; err != nil {
		return err
	}

	return nil
}

func (repo *GormUserRepository) Delete(id uint64) error {

	if err := repo.DB.Delete(&User{ID: id}).Error; err != nil {
		return err
	}

	return nil
}
