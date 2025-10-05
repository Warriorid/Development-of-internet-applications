package service

import "DIA/internal/repository"


type PitsService struct{
	repo repository.Pits
}

func NewPitsService(repo repository.Pits) *PitsService{
	return &PitsService{repo: repo}
}