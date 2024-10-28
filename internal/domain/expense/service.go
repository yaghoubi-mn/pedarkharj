package domain_expense

import "github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"

type ExpenseDomainService interface {
	Create(input ExpenseInput) error
}

type service struct {
	validator datatypes.Validator
}
