package domain_device

import (
	"time"

	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
)

type Device struct {
	ID           uint64    `gorm:"primaryKey"`
	Name         string    `gorm:"size:300,not null"` // user agent
	LastIP       string    `gorm:"size:15,not null"`
	FirstLogin   time.Time `gorm:"not null"`
	LastLogin    time.Time `gorm:"not null"`
	RefreshToken string    `gorm:"size:200"`
	UserID       uint64    `gorm:"index,not null"`
	User         domain_user.User
}
