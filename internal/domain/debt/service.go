package domain_debt

import (
	domain_shared "github.com/yaghoubi-mn/pedarkharj/internal/domain/shared"
	"github.com/yaghoubi-mn/pedarkharj/pkg/service_errors"
)

type DebtDomainService interface {
	Create(input ExpenseDebtInput) (debts []Debt, userErr error)
	Delete(debt Debt, requesterUserID uint64) (procceedDeletation bool, outDebt Debt, userErr error)
	Get(debtID uint64) (userErr error)
	GetLimited(page, limit uint, userID uint64) (userErr error)
	Accept(debt Debt, acceptorUserID uint64, isCreditorRegistered, isDebtorRegistered bool) (outDebt Debt, userErr error)
	Reject(debt Debt, rejectorUserID uint64, isCreditorRegistered, isDebtorRegistered bool) (outDebt Debt, userErr error)
	Pay(debt Debt, payerUserID uint64, isCreditorRegistered bool) (outDebt Debt, userErr error)
	AcceptPayment(debt Debt, acceptorUserID uint64, isDebtorRegistered bool) (outDebt Debt, userErr error)
}

type service struct {
	validator domain_shared.Validator
}

func NewDebtDomainService(validator domain_shared.Validator) DebtDomainService {
	return service{
		validator: validator,
	}
}

func (s service) Create(input ExpenseDebtInput) (debts []Debt, userErr error) {

	// calculate debts

	var average uint64
	for _, creditAmount := range input.Creditors {
		average += creditAmount
	}
	average /= uint64(len(input.Creditors) + len(input.Debtors))

	if average == 0 {
		return nil, service_errors.ErrLowCredit
	}

	debtorsMap := make(map[uint64]uint64)

	// find creditors that paid lower than average
	for creditorID, creditAmount := range input.Creditors {

		if creditAmount >= average {

			debtorsMap[creditorID] = 0
		} else if creditAmount > 0 {

			debtorsMap[creditorID] = average - creditAmount
		}
	}

	for _, debtorID := range input.Debtors {
		debtorsMap[debtorID] = average
	}

	debts = make([]Debt, 0, len(debtorsMap))

	for creditorID, creditAmount := range input.Creditors {
		if creditAmount < average {
			// creditor acutally is a debtor
			continue
		}

		for debtorID, debtAmount := range debtorsMap {

			if debtAmount == 0 {
				continue
			}

			var calAmount uint64
			if creditAmount >= debtAmount {
				debtorsMap[debtorID] -= debtAmount
				calAmount = debtAmount
				creditAmount -= debtAmount
			} else if creditAmount > 0 {
				calAmount = debtAmount - creditAmount
				debtorsMap[debtorID] -= creditAmount
			} else {
				break
			}

			debts = append(debts, Debt{
				ExpenseID:  input.ExpenseID,
				CreditorID: creditorID,
				DebtorID:   debtorID,
				Amount:     calAmount,
			})

		}
	}

	return debts, nil
}

func (s service) Delete(debt Debt, requesterUserID uint64) (bool, Debt, error) {

	if requesterUserID == debt.CreditorID {

		if debt.IsDebtorRequestedForDelete {
			return true, debt, nil
		} else {
			debt.IsCreditorRequestedForDelete = true
			return false, debt, nil
		}
	} else if requesterUserID == debt.DebtorID {

		if debt.IsCreditorRequestedForDelete {
			return true, debt, nil
		} else {
			debt.IsDebtorRequestedForDelete = true
			return false, debt, nil
		}
	} else {
		return false, debt, service_errors.ErrPermissionDenied
	}

}

func (s service) Get(debtID uint64) (userErr error) {
	if debtID == 0 {
		return service_errors.ErrInvalidID
	}

	return nil
}

func (s service) GetLimited(page uint, limit uint, userID uint64) (userErr error) {
	panic("unimplemented")
}

func (s service) Accept(debt Debt, acceptorUserID uint64, isCreditorRegistered, isDebtorRegistered bool) (Debt, error) {

	if debt.CreditorID == acceptorUserID {
		debt.IsCreditorAccepted = true
		debt.IsCreditorRejected = false

		if !isDebtorRegistered {
			debt.IsDebtorAccepted = true
			debt.IsDebtorRejected = false
		}
		return debt, nil

	} else if debt.DebtorID == acceptorUserID {
		debt.IsDebtorAccepted = true
		debt.IsDebtorRejected = false

		if !isCreditorRegistered {
			debt.IsCreditorAccepted = true
			debt.IsCreditorRejected = false
		}
		return debt, nil

	} else {
		return debt, service_errors.ErrPermissionDenied
	}
}

func (s service) Reject(debt Debt, rejectorUserID uint64, isCreditorRegistered, isDebtorRegistered bool) (Debt, error) {

	if debt.CreditorID == rejectorUserID {
		debt.IsCreditorRejected = true
		debt.IsCreditorAccepted = false

		if !isDebtorRegistered {
			debt.IsDebtorRejected = true
			debt.IsDebtorAccepted = false
		}
		return debt, nil

	} else if debt.DebtorID == rejectorUserID {
		debt.IsDebtorRejected = true
		debt.IsDebtorAccepted = false

		if !isCreditorRegistered {
			debt.IsCreditorRejected = true
			debt.IsCreditorAccepted = false
		}
		return debt, nil

	} else {
		return debt, service_errors.ErrPermissionDenied
	}
}

func (s service) Pay(debt Debt, payerUserID uint64, isCreditorRegistered bool) (Debt, error) {

	if debt.DebtorID == payerUserID {
		debt.IsPaid = true

		if !isCreditorRegistered {
			debt.IsPaymentAccepted = true
		}
		return debt, nil
	} else {
		return debt, service_errors.ErrPermissionDenied
	}
}

func (s service) AcceptPayment(debt Debt, acceptorUserID uint64, isDebtorRegistered bool) (Debt, error) {

	if !debt.IsPaid && isDebtorRegistered {
		return debt, service_errors.ErrDebtIsNotPaid
	}

	if debt.CreditorID == acceptorUserID {
		debt.IsPaymentAccepted = true

		if !isDebtorRegistered {
			debt.IsPaid = true
		}
		return debt, nil
	} else {
		return debt, service_errors.ErrPermissionDenied
	}
}
