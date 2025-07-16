package domain_notification

import (
	domain_shared "github.com/yaghoubi-mn/pedarkharj/internal/domain/shared"
)

type NotificationDomainService interface {
	Create(notif Notification) domain_shared.ErrorsMap
	Update(notif Notification) domain_shared.ErrorsMap
	Delete(notifID uint64) domain_shared.ErrorsMap
	GetAll(page int, limit int) ([]Notification, domain_shared.ErrorsMap)
	Get(notifID uint64) (Notification, domain_shared.ErrorsMap)
}

type service struct {
	validator domain_shared.Validator
}

func NewNotificationDomainService(validator domain_shared.Validator) NotificationDomainService {
	return service{
		validator: validator,
	}
}
func (s service) Create(notif Notification) domain_shared.ErrorsMap {

	return s.validator.Struct(notif)
}

func (s service) Update(notif Notification) domain_shared.ErrorsMap {

	errs := s.validator.Struct(notif)
	if notif.ID == 0 {

		if errs == nil {
			errs = make(map[string]string)
		}

		errs["id"] = "ID can't be zero"

	}

	return errs
}

func (s service) Delete(notifID uint64) domain_shared.ErrorsMap {
	panic("unimplemented")
}

func (s service) Get(notifID uint64) (Notification, domain_shared.ErrorsMap) {
	panic("unimplemented")
}

func (s service) GetAll(page int, limit int) ([]Notification, domain_shared.ErrorsMap) {
	panic("unimplemented")
}
