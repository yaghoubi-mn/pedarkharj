package domain_expense

import (
	"time"

	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
)

type Expense struct {
	ID          uint64
	CreatorID   uint64 `gorm:"not null"`
	Creator     domain_user.User
	Name        string `gorm:"not null,size:100" validate:"name,required"`
	Description string `gorm:"size:400" validate:"description"`
	// Type        ExpenseType    `gorm:"not null,type:enum('group','contact','single')"` // group, contact (pay all a contact debts), single (a debt or credit to another person)
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	TotalAmount uint64    `gorm:"not null"`
}

// type ExpenseType string

// const (
// 	ExpenseTypeGroup = "group"
// 	ExpenseTypeContact = "contact"
// 	ExpenseTypeSingle = "single"
// )
