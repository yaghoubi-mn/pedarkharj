package app_debt

import shared_dto "github.com/yaghoubi-mn/pedarkharj/internal/shared/dto"

type ExpenseDebtInputWithID struct {
	shared_dto.ExpenseDebtInputWithID
}

func NewExpenseDebtInputWithID(name, description string, creditors map[uint64]uint64, debtors []uint64, expenseID uint64) ExpenseDebtInputWithID {
	return ExpenseDebtInputWithID{
		ExpenseDebtInputWithID: shared_dto.ExpenseDebtInputWithID{
			Name:        name,
			Description: description,
			Creditors:   creditors,
			Debtors:     debtors,
			ExpenseID:   expenseID,
		},
	}

}
