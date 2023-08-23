package entity

import "time"

type Account struct {
	ID          uint64     `json:"id" db:"id"`
	UserName    string     `json:"user_name" db:"user_name"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at" db:"-"`
}
