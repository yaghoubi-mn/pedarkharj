package domain_expense_comment

type DebtDomainRepository interface {
	GetByID(id uint64) (ExpenseComment, error)
	GetLimitedByExpenseID(expenseID uint64, offset int, limit uint) ([]ExpenseComment, error)
	Create(expenseComment *ExpenseComment) error
	Update(expenseComment ExpenseComment) error
	Delete(id uint64) error
}
