package models

import (
	"time"

	"github.com/google/uuid"
)

const TimestampFormat = "2006-01-02T15:04:05.000Z"

type Base struct {
	ID        	uuid.UUID	`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	IsEnabled	bool      	`json:"is_enabled"`
}

func (b *Base) Init() error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	b.ID = id
	b.CreatedAt = now
	b.UpdatedAt = now
	b.IsEnabled = true

	return nil
}

func (b *Base) Touch() {
	b.UpdatedAt = time.Now().UTC()
}
