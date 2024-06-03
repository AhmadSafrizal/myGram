package model

import (
	"errors"
	"time"
)

type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	PhotoID   uint      `json:"photo_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Comment) Validate() error {
	if c.Message == "" {
		return errors.New("message must not be empty")
	}

	return nil
}
