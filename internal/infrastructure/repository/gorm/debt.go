package repository

import (
	domain_debt "github.com/yaghoubi-mn/pedarkharj/internal/domain/debt"
	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"gorm.io/gorm"
)

type GormDebtRepository struct {
	DB *gorm.DB
}

func NewGormDebtRepository(db *gorm.DB) domain_debt.DebtDomainRepository {
	return &GormDebtRepository{DB: db}
}

func (repo *GormDebtRepository) GetByID(id uint64) (domain_debt.Debt, error) {
	var user domain_debt.Debt
	if err := repo.DB.Where(domain_debt.Debt{ID: id}).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, database_errors.ErrRecordNotFound
		}

		return user, err
	}

	return user, nil
}

func (repo *GormDebtRepository) GetLimitedByUserID(userID uint64, offset int, limit int) ([]domain_debt.Debt, error) {
	var debts []domain_debt.Debt
	if err := repo.DB.Where("creditor_id=? or debtor_id=?", userID, userID).Limit(limit).Find(&debts).Error; err != nil {
		return nil, err
	}
	return debts, nil
}

// the pointer for debt is for returning id
func (repo *GormDebtRepository) Create(debt *domain_debt.Debt) error {

	if err := repo.DB.Create(&debt).Error; err != nil {
		return err
	}

	return nil
}

func (repo *GormDebtRepository) Update(debt domain_debt.Debt) error {

	if err := repo.DB.Updates(&debt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return database_errors.ErrRecordNotFound
		}

		return err
	}

	return nil
}

func (repo *GormDebtRepository) Delete(id uint64) error {

	if err := repo.DB.Delete(&domain_debt.Debt{ID: id}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return database_errors.ErrRecordNotFound
		}

		return err
	}

	return nil
}
