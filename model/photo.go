package model

import (
	"errors"
	"time"
)

type Photo struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Caption   string    `json:"caption"`
	PhotoURL  string    `json:"photo_url"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Photo) Validate() error {
	if p.Title == "" {
		return errors.New("title must not be empty")
	}

	if p.PhotoURL == "" {
		return errors.New("photo_url must not be empty")
	}

	return nil
}
