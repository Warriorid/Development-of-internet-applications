package model

type Users struct {
	ID       int    `gorm:"primaryKey;column:id;autoIncrement"`
	Username string `gorm:"column:username;type:varchar(150);not null;unique"`
	Password string `gorm:"column:password;type:varchar(128);not null"`
	IsStaff  bool   `gorm:"column:is_staff;not null;default:false"`
}



type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
    Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}