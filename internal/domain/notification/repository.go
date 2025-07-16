package domain_notification

type NotificationDomainRepository interface {
	Create(notif Notification) error
	Update(notif Notification) error
	Delete(notifID uint64) error
	GetAll(page int, limit int, sort string) ([]Notification, error)
	Get(notifID uint64) (Notification, error)
}
