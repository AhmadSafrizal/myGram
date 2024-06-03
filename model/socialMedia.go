package model

import (
	"errors"
	"time"
)

type SocialMedia struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name"`
	SocialMediaURL string    `json:"social_media_url"`
	UserID         uint      `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (s *SocialMedia) Validate() error {
	if s.Name == "" {
		return errors.New("name must not be empty")
	}

	if s.SocialMediaURL == "" {
		return errors.New("social_media_url must not be empty")
	}

	return nil
}
