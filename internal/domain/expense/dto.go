package domain_expense

import "time"

type ExpenseDebtOuput struct {
	// expense
	ID          uint64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// debt
	CreditorID                   uint64
	DebtorID                     uint64
	Amount                       uint64
	Type                         string
	IsCreditorAccepted           bool
	IsDebtorAccepted             bool
	IsCreditorRejected           bool
	IsDebtorRejected             bool
	IsPaid                       bool
	IsPaymentAccpeted            bool
	IsCreditorRequestedForDelete bool
	IsDebtorRequestedForDelete   bool

	// user
	UserAvatar string
	UserName   string
}
