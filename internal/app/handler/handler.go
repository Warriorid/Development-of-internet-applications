package handler

import (
	"DIA/internal/app/repository"

	"github.com/gin-gonic/gin"
)

type Handler struct{
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}


func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/materials", h.GetMaterials)
	router.GET("/materials/:id", h.GetMaterial)
	router.GET("/pits-calculations/:id", h.GetPit)
	router.POST("/pits-calculations/:id/delete", h.DeletePit)
	router.POST("/materials/:id/add-to-pit", h.AddMaterialToPit)
}


func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}
