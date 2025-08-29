package domain_debt

import shared_dto "github.com/yaghoubi-mn/pedarkharj/internal/shared/dto"

type ExpenseDebtInput struct {
	shared_dto.ExpenseDebtInputWithID
}

func NewExpenseDebtInput(name, description string, creditors map[uint64]uint64, debtors []uint64, expenseID uint64) ExpenseDebtInput {
	return ExpenseDebtInput{
		shared_dto.ExpenseDebtInputWithID{
			Name:        name,
			Description: description,
			Creditors:   creditors,
			Debtors:     debtors,
			ExpenseID:   expenseID,
		},
	}
}
