package model

import (
	"time"

	"github.com/google/uuid"
)

type Alert struct {
	AlertID   uint64      `json:"alert_id"`
	UserId    uuid.UUID   `json:"user_id"`
	AlertData []AlertType `json:"alert_data"`
	CreatedAt *time.Time  `json:"created_at"`
}

type AlertType struct {
	AlertTypeId uuid.UUID `json:"alert_type_id"`
	AlertName   string    `json:"alert_name"`
}
