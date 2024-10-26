package domain_expense

type ExpenseDomainRepository interface {
	GetByID(id uint64) (Expense, error)
	GetLimitedByUserID(userId uint64, limit int) ([]Expense, error)
	Create(expense *Expense) error
	Update(expense Expense) error
	Delete(id uint64) error
}
