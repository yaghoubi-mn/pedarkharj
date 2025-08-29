package domain_debt

import (
	"github.com/yaghoubi-mn/pedarkharj/internal/domain/expense"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
)

type Debt struct {
	ID uint64

	ExpenseID uint64 `gorm:"not null"`
	Expense   domain_expense.Expense

	CreditorID uint64 `gorm:"not null"`
	Creditor   domain_user.User
	DebtorID   uint64 `gorm:"not null"`
	Debtor     domain_user.User

	Amount uint64 `gorm:"not null"`

	IsCreditorAccepted           bool
	IsDebtorAccepted             bool
	IsCreditorRejected           bool
	IsDebtorRejected             bool
	IsPaid                       bool
	IsPaymentAccepted            bool
	IsDebtorRequestedForDelete   bool
	IsCreditorRequestedForDelete bool
}
