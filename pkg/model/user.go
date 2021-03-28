package model

import "time"

// User model, also known as trader
type User struct {
	ID           int64     `json:"id" gorm:"column:id"`
	DisplayName  string    `json:"display_name" gorm:"column:display_name"`
	EMail        string    `json:"e_mail" gorm:"column:e_mail"`
	PasswordHash string    `json:"password_hash" gorm:"column:password_hash"`
	LastLoginAt  time.Time `json:"last_login_at" gorm:"column:last_login_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
}
