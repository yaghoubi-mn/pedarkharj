package domain_notification

import (
	"time"

	domain_debt "github.com/yaghoubi-mn/pedarkharj/internal/domain/debt"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
)

type Notification struct {
	ID          uint64
	Title       string `gorm:"size:50; not null"`
	Image       string `gorm:"size:1000; not null"`
	Description string `gorm:"size:300; not null"`
	UserID      uint64 `gorm:"not null"`
	User        domain_user.User
	DebtID      uint64 `gorm:"not null"`
	Debt        domain_debt.Debt
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	Type        string    `gorm:"size:30; not null"`
	Amount      uint64    `grom:"not null"`
	IsCreditor  bool      `gorm:"not null"`
}
