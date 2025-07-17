package domain_user

import "time"

type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Name     string `gorm:"size:30,not null" validate:"name,required,max=30"`
	Number   string `gorm:"size:13,not null,index,unique" validate:"phone_number,required,max=15"`
	Password string `gorm:"size:30,not null" validate:"max=30"`
	Salt     string `gorm:"size:32,not null"`
	Avatar   string `gorm:"size:500,not null"`

	RegisteredAt time.Time `gorm:"not null"`
	IsRegistered bool      `gorm:"not null"`
	IsBlocked    bool
}
