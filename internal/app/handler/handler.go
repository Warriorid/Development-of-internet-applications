package handler

import (
	"DIA/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	router.GET("/materials/pits/:id", h.GetPit)
	router.POST("/materials/pits/:id/delete", h.DeletePit)
	router.POST("/materials/:id/add-to-pit", h.AddMaterialToPit)
}


func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./resources")
}


func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}