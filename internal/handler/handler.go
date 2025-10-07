package handler

import (
	"DIA/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct{
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}


func (h *Handler) RegisterHandler(router *gin.Engine) {
	api := router.Group("/api")

	materials := api.Group("/materials")
	{
		materials.GET("", h.GetMaterials)
		materials.GET("/:id", h.GetMaterial)
		materials.POST("", h.CreateMaterial)
		materials.PUT("/:id", h.UpdateMaterial)
		materials.DELETE("/:id", h.DeleteMaterial)
		materials.POST("/:id/pit", h.AddMaterialToPit)
		materials.POST("/:id/image", h.UploadMaterialImage)
	}

	pits := api.Group("/pits")
	{
		pits.GET("/draft", h.GetDraftPit)
		pits.GET("", h.GetPits)
		pits.GET("/:id", h.GetPitWithMaterials)
		pits.PUT("/:id", h.UpdatePit)
		pits.PUT("/:id/form", h.FormPit)
		pits.PUT("/:id/complete", h.CompletePit)
		pits.DELETE("/:id", h.DeletePit)
	}
	
	calculationMaterials := api.Group("/calculation-materials")
	{
		calculationMaterials.DELETE("", h.removeMaterialFromPit)
		calculationMaterials.PUT("", h.UpdateCalculationMaterial)
	}
	
	
	users := api.Group("/users")
	{
		users.POST("/register", h.RegisterUser)
    	users.POST("/login", h.Login)
    	users.POST("/logout", h.Logout)
    	users.GET("/profile/:id", h.GetUserProfile)
    	users.PUT("/profile/:id", h.UpdateUserProfile)
	}

}

