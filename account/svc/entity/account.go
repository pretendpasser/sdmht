package entity

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Account struct {
	ID          uint64    `json:"id" db:"id"`
	UserName    string    `json:"user_name" db:"user_name"`
	WeChatID    string    `json:"wechat_id" db:"wechat_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	LastLoginAt time.Time `json:"last_login_at" db:"last_login_at"`
}

type Claims struct {
	jwt.StandardClaims
	AccountID uint64
	WechatID  string
}

type RegisterReq struct {
	WechatID string `json:"wechat_id"`
	UserName string `json:"user_name"`
}

type LoginReq struct {
	WechatID string `json:"wechat_id"`
}

type LoginRes struct {
	Token   string   `json:"token"`
	Account *Account `json:"account"`
}
