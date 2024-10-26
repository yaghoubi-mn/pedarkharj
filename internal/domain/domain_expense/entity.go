package domain_expense

import (
	"time"

	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
)

type Expense struct {
	ID          uint64
	CreatorID   uint64 `gorm:"not null"`
	Creator     domain_user.User
	Name        string    `gorm:"not null,size:100"`
	Description string    `gorm:"size:400"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
	TotalAmount uint64    `gorm:"not null"`
}
