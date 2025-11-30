package model

import (
	"DIA/internal/role"
)

type Users struct {
    ID       int       `gorm:"primaryKey;column:id;autoIncrement" json:"id"`
    Username string    `gorm:"column:username;type:varchar(150);not null;unique" json:"username"`
    Password string    `gorm:"column:password;type:varchar(128);not null" json:"-"`
    Role     role.Role `gorm:"column:role;type:integer;not null;default:0" json:"role"`
}

type JWTClaims struct {
    UserID   int      `json:"user_id"`
    Username string   `json:"username"`
    Role     role.Role `json:"role"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     *int `json:"role"`
}


type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginResp struct {
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type UpdateProfileRequest struct {
    Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}