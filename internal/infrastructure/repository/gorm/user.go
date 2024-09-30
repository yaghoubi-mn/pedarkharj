package repository

import (
	"github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	DB *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) domain_user.UserDomainRepository {
	return &GormUserRepository{DB: db}
}

func (repo *GormUserRepository) GetByID(id uint64) (domain_user.User, error) {
	var user domain_user.User
	if err := repo.DB.Where(domain_user.User{ID: id}).Find(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, database_errors.ErrRecordNotFound
		}

		return user, err
	}

	return user, nil
}

func (repo *GormUserRepository) GetByNumber(number string) (domain_user.User, error) {
	var u domain_user.User
	if err := repo.DB.First(&u, domain_user.User{Number: number}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return u, database_errors.ErrRecordNotFound
		}

		return u, err
	}

	return u, nil
}

func (repo *GormUserRepository) Create(user domain_user.User) error {

	if err := repo.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (repo *GormUserRepository) Update(user domain_user.User) error {

	if err := repo.DB.Updates(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return database_errors.ErrRecordNotFound
		}

		return err
	}

	return nil
}

func (repo *GormUserRepository) Delete(id uint64) error {

	if err := repo.DB.Delete(&domain_user.User{ID: id}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return database_errors.ErrRecordNotFound
		}

		return err
	}

	return nil
}
