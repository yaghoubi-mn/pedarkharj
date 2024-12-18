package domain_expense

import "time"

type ExpenseInput struct {
	Name        string
	Description string
	Creditors   map[string]uint64 // {"<ID>": <Amount>, ...}
	Debtors     []uint64          // list if debtors IDs

}

type ExpenseDebtOuput struct {
	// expense
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// debt
	CreditorID                   uint64 `json:"creditor_id"`
	DebtorID                     uint64 `json:"debtor_id"`
	Amount                       uint64 `json:"amount"`
	Type                         string `json:"type"`
	IsCreditorAccepted           bool   `json:"is_creditor_accepted"`
	IsDebtorAccepted             bool   `json:"is_debtor_accepted"`
	IsCreditorRejected           bool   `json:"is_creditor_rejected"`
	IsDebtorRejected             bool   `json:"is_debtor_rejected"`
	IsPaid                       bool   `json:"is_paid"`
	IsPaymentAccpeted            bool   `json:"is_payment_accpeted"`
	IsCreditorRequestedForDelete bool   `json:"is_creditor_requested_for_delete"`
	IsDebtorRequestedForDelete   bool   `json:"is_debtor_requested_for_delete"`

	// user
	UserAvatar string `json:"user_avatar"`
	UserName   string `json:"user_name"`
}
