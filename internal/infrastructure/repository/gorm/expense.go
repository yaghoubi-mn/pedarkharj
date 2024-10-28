package repository

import (
	domain_expense "github.com/yaghoubi-mn/pedarkharj/internal/domain/expense"
	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"gorm.io/gorm"
)

type GormExpenseRepository struct {
	DB *gorm.DB
}

func NewGormExpenseRepository(db *gorm.DB) domain_expense.ExpenseDomainRepository {
	return &GormExpenseRepository{DB: db}
}

func (repo *GormExpenseRepository) GetByID(id uint64) (domain_expense.Expense, error) {
	var user domain_expense.Expense
	if err := repo.DB.Where(domain_expense.Expense{ID: id}).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, database_errors.ErrRecordNotFound
		}

		return user, err
	}

	return user, nil
}

func (repo *GormExpenseRepository) GetLimitedExpenseDebtByUserID(userID uint64, offset int, limit int) ([]domain_expense.ExpenseDebtOuput, error) {
	var expenses []domain_expense.ExpenseDebtOuput
	if err := repo.DB.Joins("LEFT JOIN debts ON debts.expense_id = expenses.id").
		Joins("LEFT JOIN users creditor_users ON creditor_users.id = debts.creditor_id"). // join users for his/her name and avatar
		// Joins("LEFT JOIN users debtor_users ON debtor_users.id = debts.debtor_id").
		Select("").
		Where("debts.creditor_id=? OR debts.debtor_id=?", userID, userID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&expenses).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return expenses, database_errors.ErrRecordNotFound
		}

		return nil, err
	}

	var debtors []struct {
		UserName   string
		UserAvatar string
	}

	// for user name and avatar
	if err := repo.DB.Joins("LEFT JOIN debts ON debts.expense_id = expenses.id").
		Joins("LEFT JOIN users debtor_users ON debtor_users.id = debts.debtor_id").
		Select("debtor_users.name as user_name, debtor_users.avatar as user_avatar").
		Where("debts.creditor_id=? OR debts.debtor_id=?", userID, userID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&debtors).Error; err != nil {

		return nil, err
	}

	for i := range expenses {
		if userID == expenses[i].CreditorID {
			expenses[i].UserName = debtors[i].UserName
			expenses[i].UserAvatar = debtors[i].UserAvatar
		}
	}
	return expenses, nil
}

func (repo *GormExpenseRepository) Create(user *domain_expense.Expense) error {

	if err := repo.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (repo *GormExpenseRepository) Update(user domain_expense.Expense) error {

	if err := repo.DB.Updates(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return database_errors.ErrRecordNotFound
		}

		return err
	}

	return nil
}

func (repo *GormExpenseRepository) Delete(id uint64) error {

	if err := repo.DB.Delete(&domain_expense.Expense{ID: id}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return database_errors.ErrRecordNotFound
		}

		return err
	}

	return nil
}
