package device

import "time"

type Device struct {
	ID           uint64 `gorm:"primaryKey"`
	Name         string `gorm:"size:300"`
	LastIP       string `gorm:"size:15"`
	FirstLogin   time.Time
	LastLogin    time.Time
	RefreshToken string `gorm:"size:200"`
	UserID       uint64 `gorm:"index,constraint:OnDelete:CASCADE"` //`gorm:"foreginKey:users,references:ID,OnDelete:CASECADE"`
}
