package app_expense

import (
	app_debt "github.com/yaghoubi-mn/pedarkharj/internal/application/debt"
	app_shared "github.com/yaghoubi-mn/pedarkharj/internal/application/shared"
	domain_expense "github.com/yaghoubi-mn/pedarkharj/internal/domain/expense"
	"github.com/yaghoubi-mn/pedarkharj/pkg/rcodes"
)

type ExpenseAppService interface {
	Create(input ExpenseInputWithPhoneNumber, userID uint64, userPhoneNumber string) app_shared.ResponseDTO
	Update(input ExpenseUpdateInput)
	Delete(expenseID, userID uint64) app_shared.ResponseDTO
	Get(expenseID, userID uint64) app_shared.ResponseDTO
	GetLimited(userID uint64, page, limit uint) app_shared.ResponseDTO
}

type service struct {
	domainService  domain_expense.ExpenseDomainService
	repo           domain_expense.ExpenseDomainRepository
	debtAppService app_debt.DebtAppService
}

func NewExpenseAppService(repo domain_expense.ExpenseDomainRepository, domainService domain_expense.ExpenseDomainService, debtAppService app_debt.DebtAppService) ExpenseAppService {
	return service{
		repo:           repo,
		debtAppService: debtAppService,
		domainService:  domainService,
	}
}

func (s service) Create(input ExpenseInputWithPhoneNumber, userID uint64, userPhoneNumber string) (responseDTO app_shared.ResponseDTO) {
	responseDTO.Data = make(map[string]any)

	expense, userErr := s.domainService.Create(domain_expense.NewExpenseInputWithPhoneNumber(
		input.Name,
		input.Description,
		input.Creditors,
		input.Debtors,
		userID,
		userPhoneNumber,
	))

	if userErr != nil {
		responseDTO.UserErr = userErr
		responseDTO.ResponseCode = rcodes.InvalidField
		return
	}

	err := s.repo.Create(&expense)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}
	numbers := make([]string, 0, len(input.Creditors))
	for k := range input.Creditors {
		numbers = append(numbers, k)
	}

	numbers = append(numbers, input.Debtors...)

	err = s.repo.CreateUsersWithNumbers(numbers)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}

	idPhoneMap, err := s.repo.GetUserIDOfPhoneNumbers(numbers)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}

	// create expense input with userID instead of phone number
	var expenseInput ExpenseInputWithID
	expenseInput.Fill(input, idPhoneMap, expense.ID)

	responseDTO2 := s.debtAppService.Create(app_debt.NewExpenseDebtInputWithID(
		expenseInput.Name,
		expenseInput.Description,
		expenseInput.Creditors,
		expenseInput.Debtors,
		expenseInput.ExpenseID,
	))

	if responseDTO2.UserErr != nil || responseDTO2.ServerErr != nil {
		return responseDTO2
	}

	responseDTO.Data["msg"] = "Done"
	return
}

func (s service) Delete(expenseID uint64, userID uint64) (responseDTO app_shared.ResponseDTO) {

	responseDTO.Data = make(map[string]any)

	userErr := s.domainService.Delete(expenseID)
	if userErr != nil {
		responseDTO.UserErr = userErr
		return
	}

	err := s.repo.Delete(expenseID, userID)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}

	responseDTO.Data["msg"] = "Done"
	return
}

func (s service) Get(expenseID uint64, userID uint64) (responseDTO app_shared.ResponseDTO) {

	responseDTO.Data = make(map[string]any)

	userErr := s.domainService.Get(expenseID)
	if userErr != nil {
		responseDTO.UserErr = userErr
		return
	}

	expense, err := s.repo.GetByID(expenseID, userID)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}

	responseDTO.Data["expense"] = expense
	return
}

func (s service) GetLimited(userID uint64, page, limit uint) (responseDTO app_shared.ResponseDTO) {

	responseDTO.Data = make(map[string]any)

	userErr := s.domainService.GetLimited(page, limit)
	if userErr != nil {
		responseDTO.UserErr = userErr
		return
	}

	expenses, err := s.repo.GetLimitedExpenseDebtByUserID(userID, int(page), int(limit))
	if err != nil {
		responseDTO.ServerErr = err
		return
	}

	responseDTO.Data["data"] = expenses
	return
}

func (s service) Update(input ExpenseUpdateInput) {
	panic("unimplemented")
}
