package domain_expense

type ExpenseDomainRepository interface {
	GetByID(id uint64) (Expense, error)
	GetLimitedExpenseDebtByUserID(userId uint64, offset int, limit int) ([]ExpenseDebtOuput, error)
	Create(expense *Expense) error
	Update(expense Expense) error
	Delete(id uint64) error
}
