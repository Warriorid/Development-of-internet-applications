package service

import (
	"DIA/internal/model"
	"DIA/internal/repository"
	"fmt"
)

type UserService struct {
	repo repository.User
}

func NewUserService (repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(user *model.Users) (*model.Users, error) {
	err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(username, password string) (*model.Users, error) {
	user, err := s.repo.Authentication(username, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserProfile(id int) (*model.Users, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}


func (s *UserService) UpdateUserProfile(id int, username, password string) (*model.Users, error) {
    _, err := s.repo.GetUserByID(id)
    if err != nil {
        return nil, fmt.Errorf("user not found")
    }
    updateData := &model.Users{
        Username: username,
        Password: password,
    }
    err = s.repo.UpdateUser(id, updateData)
    if err != nil {
        return nil, err
    }
    return s.repo.GetUserByID(id)
}

func (s *UserService) Logout() error {
	return nil
}