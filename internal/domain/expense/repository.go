package domain_expense

type ExpenseDomainRepository interface {
	GetByID(id uint64, userID uint64) (Expense, error)
	GetLimitedExpenseDebtByUserID(userId uint64, offset int, limit int) ([]ExpenseDebtOuput, error)
	Create(expense *Expense) error
	Update(expense Expense) error
	Delete(id uint64, userID uint64) error
	CreateUsersWithNumbers(numbers []string) error
	GetUserIDOfPhoneNumbers(numbers []string) (map[string]uint64, error)
}
