package domain_expense_comment

import (
	"time"

	"github.com/yaghoubi-mn/pedarkharj/internal/domain/expense"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
)

type ExpenseComment struct {
	ID uint64

	UserID    uint64 `gorm:"not null"`
	User      domain_user.User
	ExpenseID uint64 `gorm:"not null"`
	Expense   domain_expense.Expense

	Content string `gorm:"not null,size:400" validate:"description"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
