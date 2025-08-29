package app_debt

import (
	app_shared "github.com/yaghoubi-mn/pedarkharj/internal/application/shared"
	domain_debt "github.com/yaghoubi-mn/pedarkharj/internal/domain/debt"
)

type DebtAppService interface {
	Create(input ExpenseDebtInputWithID) app_shared.ResponseDTO
}

type service struct {
	repo          domain_debt.DebtDomainRepository
	domainService domain_debt.DebtDomainService
}

func NewDebtAppService(repo domain_debt.DebtDomainRepository, domainService domain_debt.DebtDomainService) DebtAppService {
	return service{
		repo:          repo,
		domainService: domainService,
	}
}

func (s service) Create(input ExpenseDebtInputWithID) (responseDTO app_shared.ResponseDTO) {
	responseDTO.Data = make(map[string]any)

	debts, userErr := s.domainService.Create(domain_debt.NewExpenseDebtInput(
		input.Name,
		input.Description,
		input.Creditors,
		input.Debtors,
		input.ExpenseID,
	))
	if userErr != nil {
		responseDTO.UserErr = userErr
		return
	}

	err := s.repo.CreateMultipleWithTransaction(debts)
	if err != nil {
		responseDTO.ServerErr = err
		return
	}

	responseDTO.Data["msg"] = "Done"
	return
}
