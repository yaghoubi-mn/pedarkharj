package domain_notification

import (
	shared_dto "github.com/yaghoubi-mn/pedarkharj/internal/shared/dto"
)

type NotificationInput struct {
	shared_dto.NotificationInput
}

func NewNotificationInput(title, image, description, typ string, debtID, userID uint64, amount uint64, isCreditor bool) NotificationInput {
	return NotificationInput{
		shared_dto.NotificationInput{
			Title:       title,
			Image:       image,
			Description: description,
			UserID:      userID,
			DebtID:      debtID,
			Type:        typ,
			Amount:      amount,
			IsCreditor:  isCreditor,
		},
	}
}

func (n NotificationInput) GetNotification() Notification {
	return Notification{
		Title:       n.Title,
		Image:       n.Image,
		Description: n.Description,
		UserID:      n.UserID,
		DebtID:      n.DebtID,
		Type:        n.Type,
		Amount:      n.Amount,
		IsCreditor:  n.IsCreditor,
	}
}

type NotificationOutput struct {
	shared_dto.NotificationOutput
}
