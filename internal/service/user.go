package service

import (
	"DIA/internal/model"
	"DIA/internal/repository"
	"DIA/pkg/config"
	"fmt"
)

type UserService struct {
	repo repository.User
}

func NewUserService (repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(user *model.Users) (*model.Users, error) {
	user.Password = config.GeneratePasswordHash(user.Password)
	return user, s.repo.CreateUser(user)
}

func (s *UserService) Login(username, password string) (*model.LoginResp, error) {
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	if !config.CheckPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := config.GenerateJWTToken(user.ID, user.Username, int(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token")
	}
	
	loginResp := &model.LoginResp{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   24 * 60 * 60, 
	}
	
	return loginResp, nil
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
        Password: config.GeneratePasswordHash(password),
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