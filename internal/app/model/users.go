package model

type Users struct {
    ID       int    `gorm:"primaryKey;column:id"`
    Username string `gorm:"column:username"`
    Password string `gorm:"column:password"`
    IsStaff  bool   `gorm:"column:is_staff"`
}