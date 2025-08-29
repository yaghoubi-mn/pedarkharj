package domain_expense

import (
	shared_dto "github.com/yaghoubi-mn/pedarkharj/internal/shared/dto"
)

type ExpenseUpdateInput struct {
	shared_dto.ExpenseUpdateInput
}

func NewExpenseUpdateInput(name, description string) ExpenseUpdateInput {
	return ExpenseUpdateInput{
		ExpenseUpdateInput: shared_dto.ExpenseUpdateInput{
			Name:        name,
			Description: description,
		},
	}
}

type ExpenseInputWithPhoneNumber struct {
	shared_dto.ExpenseInputWithPhoneNumber
	CreatorID          uint64
	CreatorPhoneNumber string
}

func NewExpenseInputWithPhoneNumber(name, description string, creditors map[string]uint64, debtors []string, creatorID uint64, creatorPhoneNumber string) ExpenseInputWithPhoneNumber {
	return ExpenseInputWithPhoneNumber{
		ExpenseInputWithPhoneNumber: shared_dto.ExpenseInputWithPhoneNumber{
			Name:        name,
			Description: description,
			Creditors:   creditors,
			Debtors:     debtors,
		},
		CreatorID:          creatorID,
		CreatorPhoneNumber: creatorPhoneNumber,
	}
}

func (e ExpenseInputWithPhoneNumber) GetExpense(totalAmount uint64) Expense {
	return Expense{
		CreatorID:   e.CreatorID,
		Name:        e.Name,
		Description: e.Description,
		TotalAmount: totalAmount,
	}
}

type ExpenseDebtOuput struct {
	shared_dto.ExpenseDebtOuput
}
