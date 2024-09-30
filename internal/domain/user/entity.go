package domain_user

import "time"

type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Name     string `gorm:"size:20,not null"`
	Number   string `gorm:"size:13,not null,index,unique"`
	Password string `gorm:"size:16,not null"`
	Salt     string `gorm:"size:32,not null"`

	RegisteredAt time.Time `gorm:"not null"`
	IsRegistered bool      `gorm:"not null"`
	IsBlocked    bool
}
