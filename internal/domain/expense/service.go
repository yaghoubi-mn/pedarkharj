package domain_expense

import "github.com/yaghoubi-mn/pedarkharj/internal/domain/shared"

type ExpenseDomainService interface {
	Create(input ExpenseInput) error
}

type service struct {
	validator domain_shared.Validator
}
