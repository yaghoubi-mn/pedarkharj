package domain_expense

import (
	domain_shared "github.com/yaghoubi-mn/pedarkharj/internal/domain/shared"
	"github.com/yaghoubi-mn/pedarkharj/pkg/service_errors"
	"slices"
)

type ExpenseDomainService interface {
	Create(input ExpenseInputWithPhoneNumber) (expense Expense, userErr error)
	Update(input ExpenseUpdateInput) (userErr error)
	Delete(expenseID uint64) (userEr error)
	Get(expenseID uint64) (userErr error)
	GetLimited(page, limit uint) (userErr error)
	GetMyExpenseDebtLimited(userID uint64, page, limit uint) (userErr error)
}

type service struct {
	validator domain_shared.Validator
}

func NewExpenseService(validator domain_shared.Validator) ExpenseDomainService {
	return &service{
		validator: validator,
	}
}

func (s service) Create(input ExpenseInputWithPhoneNumber) (expense Expense, userErr error) {

	if err := s.validator.ValidateFieldByFieldName("Name", input.Name, Expense{}); err != nil {
		return expense, service_errors.ErrInvalidName
	}

	if err := s.validator.ValidateFieldByFieldName("Description", input.Description, Expense{}); err != nil {
		return expense, service_errors.ErrInvalidDescription
	}

	if len(input.Creditors) == 0 {
		return expense, service_errors.ErrEmptyCreditors
	}

	isCreatorFound := false

	for phoneNumber, credit := range input.Creditors {
		expense.TotalAmount += credit

		if err := s.validator.ValidateField(phoneNumber, "phone_number"); err != nil {
			return expense, service_errors.ErrInvalidNumber
		}

		if credit == 0 {
			return expense, service_errors.ErrInvalidCredit
		}

		// check creditor isn't in debtors
		if slices.Contains(input.Debtors, phoneNumber) {
			return expense, service_errors.ErrCommonCreditorAndDebtor
		}

		if phoneNumber == input.CreatorPhoneNumber {
			isCreatorFound = true
		}
	}

	for _, debtorPhoneNumber := range input.Debtors {

		if err := s.validator.ValidateField(debtorPhoneNumber, "phone_number"); err != nil {
			return expense, service_errors.ErrInvalidNumber
		}

		if debtorPhoneNumber == input.CreatorPhoneNumber {
			isCreatorFound = true
		}
	}

	if !isCreatorFound {
		return expense, service_errors.ErrCreatorNustBeInCreditorsOrDebtors
	}

	expense = input.GetExpense(expense.TotalAmount)
	return expense, nil

}

func (s service) Update(input ExpenseUpdateInput) (userErr error) {

	if err := s.validator.ValidateFieldByFieldName("Name", input.Name, Expense{}); err != nil {
		return service_errors.ErrInvalidName
	}

	if err := s.validator.ValidateFieldByFieldName("Description", input.Description, Expense{}); err != nil {
		return service_errors.ErrInvalidDescription
	}
	return nil
}

func (s service) Delete(expenseID uint64) (userEr error) {

	if expenseID == 0 {
		return service_errors.ErrInvalidID
	}

	return nil
}

func (s service) Get(expenseID uint64) (userErr error) {

	if expenseID == 0 {
		return service_errors.ErrInvalidID
	}

	return nil
}

func (s service) GetLimited(page uint, limit uint) error {

	if page == 0 {
		return service_errors.ErrInvalidPage
	}

	if limit < 1 {
		return service_errors.ErrInvalidLimit
	}

	return nil
}

func (s service) GetMyExpenseDebtLimited(userID uint64, page uint, limit uint) error {

	if page == 0 {
		return service_errors.ErrInvalidPage
	}

	if limit < 1 {
		return service_errors.ErrInvalidLimit
	}

	if userID == 0 {
		return service_errors.ErrInvalidID
	}

	return nil
}
