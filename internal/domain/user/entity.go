package domain_user

import "time"

type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Name     string `gorm:"size:30,not null" validate:"name,required,size:30"`
	Number   string `gorm:"size:13,not null,index,unique" validate:"e164,required,size:15"`
	Password string `gorm:"size:30,not null" validate:"size:30"`
	Salt     string `gorm:"size:32,not null"`

	RegisteredAt time.Time `gorm:"not null"`
	IsRegistered bool      `gorm:"not null"`
	IsBlocked    bool
}
