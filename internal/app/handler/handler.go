package handler

import (
	"DIA/internal/app/repository"
)

type Handler struct{
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}