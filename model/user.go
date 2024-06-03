package model

import (
	"errors"
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email must not be empty")
	}

	if u.Username == "" {
		return errors.New("username must not be empty")
	}

	if len(u.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	if u.Age < 8 {
		return errors.New("age must be at least 8")
	}

	return nil
}
