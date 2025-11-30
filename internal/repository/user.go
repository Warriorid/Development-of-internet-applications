package repository

import (
	"DIA/internal/model"
	"fmt"

	"gorm.io/gorm"
)

type UserPostgres struct {
	db *gorm.DB
}

func NewUserPostgres(db *gorm.DB) *UserPostgres {
	return &UserPostgres{db: db}
}


func (r *UserPostgres) CreateUser(user *model.Users) error {
	var count int64
	r.db.Model(&model.Users{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return fmt.Errorf("username already exists")
	}
	return r.db.Create(user).Error
}

func (r *UserPostgres) GetUserByID(id int) (*model.Users, error) {
	var user model.Users
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *UserPostgres) UpdateUser(id int, user *model.Users) error {
    if user.Username != "" {
        var count int64
        r.db.Model(&model.Users{}).Where("username = ? AND id != ?", user.Username, id).Count(&count)
        if count > 0 {
            return fmt.Errorf("username already exists")
        }
    }
    return r.db.Model(&model.Users{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "username": user.Username,
            "password": user.Password,
        }).Error
}

func (r *UserPostgres) GetUserByUsername(username string) (*model.Users, error) {
	var user model.Users
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *UserPostgres) Authentication(username, password string) (*model.Users, error) {
	user, err := r.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("invalid authentication")
	}

	if user.Password != password {
		return nil, fmt.Errorf("invalid authentication")
	}

	return user, nil
}