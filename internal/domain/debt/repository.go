package domain_debt

type DebtDomainRepository interface {
	GetByID(id uint64) (Debt, error)
	GetLimitedByUserID(userId uint64, offset int, limit int) ([]Debt, error)
	Create(debt *Debt) error
	Update(debt Debt) error
	Delete(id uint64) error
}
