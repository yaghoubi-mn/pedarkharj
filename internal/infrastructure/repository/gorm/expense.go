package repository

import (
	domain_expense "github.com/yaghoubi-mn/pedarkharj/internal/domain/expense"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	"github.com/yaghoubi-mn/pedarkharj/pkg/database_errors"
	"gorm.io/gorm"
)

type GormExpenseRepository struct {
	DB *gorm.DB
}

func NewGormExpenseRepository(db *gorm.DB) domain_expense.ExpenseDomainRepository {
	return &GormExpenseRepository{DB: db}
}

// todo: this should be done in users service
func (repo *GormExpenseRepository) CreateUsersWithNumbers(numbers []string) error {
	var existingUsers []domain_user.User
	if err := repo.DB.Select("number").Where("number IN ?", numbers).Find(&existingUsers).Error; err != nil {
		return err
	}

	existingPhones := make(map[string]bool)
	for _, user := range existingUsers {
		existingPhones[user.Number] = true
	}

	if len(numbers) == len(existingPhones) {
		// all users alredy exist in database
		return nil
	}

	newUsers := make([]domain_user.User, 0, 2)
	for _, number := range numbers {
		_, ok := existingPhones[number]
		if ok {
			continue
		}
		var user domain_user.User
		user.IsRegistered = false
		user.Name = number
		user.Number = number
		user.IsBlocked = false
		user.Avatar = "default"
		user.Password = "No Password"
		user.Salt = "No Salt"

		newUsers = append(newUsers, user)
	}

	return repo.DB.Transaction(func(tx *gorm.DB) error {

		return tx.Create(&newUsers).Error
	})
}

func (repo *GormExpenseRepository) GetUserIDOfPhoneNumbers(numbers []string) (map[string]uint64, error) {

	var users []domain_user.User
	if err := repo.DB.Where("number In ?", numbers).Find(&users).Error; err != nil {
		return nil, err
	}

	idNumberMap := make(map[string]uint64)
	for _, user := range users {
		idNumberMap[user.Number] = user.ID
	}

	return idNumberMap, nil
}

func (repo *GormExpenseRepository) GetByID(id uint64, userID uint64) (domain_expense.Expense, error) {
	var expense domain_expense.Expense
	if err := repo.DB.Model(&domain_expense.Expense{}).
		Joins("JOIN debts ON debts.expense_id = expense.id").
		Where("expense.id = ? AND (debts.creditor_id=? OR debts.debtor_id=?)", id, userID, userID).
		First(&expense).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return expense, database_errors.ErrRecordNotFound
		}

		return expense, err
	}

	return expense, nil
}

func (repo *GormExpenseRepository) GetLimitedExpenseDebtByUserID(userID uint64, offset int, limit int) ([]domain_expense.ExpenseDebtOuput, error) {
	var expenses []domain_expense.ExpenseDebtOuput
	if err := repo.DB.Model(&domain_expense.Expense{}).
		Select(`expenses.id,
			expenses.name,
			expenses.description,
			expenses.created_at,
			expenses.updated_at,
			debts.creditor_id,
			debts.debtor_id,
			debts.amount,
			debts.is_creditor_accepted,
			debts.is_debtor_accepted,
			debts.is_creditor_rejected,
			debts.is_debtor_rejected,
			debts.is_paid,
			debts.is_payment_accepted,
			debts.is_creditor_requested_for_delete,
			debts.is_debtor_requested_for_delete,
			CASE
				WHEN debts.creditor_id = ? then debtor_user.name
				ELSE creditor_user.name
			END as contact_user_name
			CASE
				WHEN debts.creditor_id = ? then debtor_user.avatar
				ELSE creditor_user.avatar
			END as contact_user_avatar
			CASE
				WHEN debts.creditor_id = ? then 'debtor'
				ELSE 'creditor'
			END as user_role
			`, userID, userID, userID).
		Joins("JOIN debts ON debts.expense_id = expenses.id").
		Joins("JOIN users as creditor_user ON creditor_user.id = debts.creditor_id"). // join users for contact name and avatar
		Joins("JOIN users as debtor_user ON debtor_user.id = debts.debtor_id").
		Where("debts.creditor_id=? OR debts.debtor_id=?", userID, userID).Order("created_at DESC").Offset(offset).Limit(limit).Find(&expenses).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return expenses, database_errors.ErrRecordNotFound
		}

		return nil, err
	}

	return expenses, nil
}

func (repo *GormExpenseRepository) Create(user *domain_expense.Expense) error {

	if err := repo.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (repo *GormExpenseRepository) Update(expense domain_expense.Expense) error {

	if err := repo.DB.Updates(&expense).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return database_errors.ErrRecordNotFound
		}

		return err
	}

	return nil
}

func (repo *GormExpenseRepository) Delete(id uint64, userID uint64) error {

	if err := repo.DB.Delete(&domain_expense.Expense{ID: id, CreatorID: userID}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return database_errors.ErrRecordNotFound
		}

		return err
	}

	return nil
}
