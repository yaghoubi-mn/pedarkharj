package app_expense

import (
	shared_dto "github.com/yaghoubi-mn/pedarkharj/internal/shared/dto"
)

type ExpenseDebtOuput struct {
	shared_dto.ExpenseDebtOuput
}

type ExpenseInputWithID struct {
	shared_dto.ExpenseDebtInputWithID
}

func (e *ExpenseInputWithID) Fill(input ExpenseInputWithPhoneNumber, idAndNumberMap map[string]uint64, expenseID uint64) {
	e.Name = input.Name
	e.Description = input.Description
	e.ExpenseID = expenseID
	e.Creditors = make(map[uint64]uint64)
	e.Debtors = make([]uint64, 0, len(input.Debtors))
	for creditorPhoneNumber, creditorAmount := range input.Creditors {

		e.Creditors[idAndNumberMap[creditorPhoneNumber]] = creditorAmount
	}

	for _, debtorPhoneNumber := range input.Debtors {

		e.Debtors = append(e.Debtors, idAndNumberMap[debtorPhoneNumber])
	}
}

type ExpenseInputWithPhoneNumber struct {
	shared_dto.ExpenseInputWithPhoneNumber
}

type ExpenseUpdateInput struct {
	shared_dto.ExpenseUpdateInput
}
