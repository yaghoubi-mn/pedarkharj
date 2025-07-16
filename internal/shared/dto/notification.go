package shared_dto

import (
	"time"
)

type NotificationInput struct {
	Title       string `validate:"name,required"`
	Image       string `validate:"url,required"`
	Description string `validate:"description,required"`
	UserID      uint64 `validate:"gte=1"`
	DebtID      uint64 `validate:"gte=1"`
	Type        string `validate:"name,required"`
	Amount      uint64 `validate:"gte=1"`
	IsCreditor  bool   `validate:"required"`
}

type NotificationOutput struct {
	Title       string    `json:"title"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	UserID      uint64    `json:"user_id"`
	DebtID      uint64    `json:"debt_id"`
	Type        string    `json:"type"`
	Amount      uint64    `json:"amount"`
	IsCreditor  bool      `json:"is_creditor"`
	CreatedAt   time.Time `json:"created_at"`
}
