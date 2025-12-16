package handler

import (
	_ "DIA/docs"
	"DIA/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct{
	service *service.Service
	redis   *redis.Client
}

func NewHandler(service *service.Service, redis *redis.Client) *Handler {
	return &Handler{service: service, redis: redis}
}


func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"https://localhost:3000", "http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")

	materials := api.Group("/materials")
{
    materials.GET("", h.GetMaterials)     
    materials.GET("/:id", h.GetMaterial)  
}

	protectedMaterials := api.Group("/materials", h.authMiddleware())
	{
		protectedMaterials.POST("", h.CreateMaterial) 
		protectedMaterials.PUT("/:id", h.UpdateMaterial)     
		protectedMaterials.DELETE("/:id", h.DeleteMaterial)  
		protectedMaterials.POST("/:id/pit", h.AddMaterialToPit)  
		protectedMaterials.POST("/:id/image", h.UploadMaterialImage) 
	}

	pits := api.Group("/pits", h.authMiddleware())
	{
		pits.GET("/draft", h.GetDraftPit)
		pits.GET("", h.GetPits)
		pits.GET("/:id", h.GetPitWithMaterials)
		pits.PUT("/:id", h.UpdatePit)
		pits.PUT("/:id/form", h.FormPit)
		pits.PUT("/:id/complete", h.CompletePit)
		pits.DELETE("/:id", h.DeletePit)
	}
	
	calculationMaterials := api.Group("/calculation-materials", h.authMiddleware())
	{
		calculationMaterials.DELETE("/:calculation_id/:material_id", h.removeMaterialFromPit)
		calculationMaterials.PUT("/:calculation_id/:material_id", h.UpdateCalculationMaterial)
	}
	
	
	users := api.Group("/users")
	{
		users.POST("", h.RegisterUser)
    	users.POST("/login", h.Login)
	}
	protectUser := api.Group("/users", h.authMiddleware())
	{
		protectUser.GET("/:id", h.GetUserProfile)
		protectUser.PUT("/:id", h.UpdateUserProfile)
		protectUser.POST("/logout", h.Logout)
	}

}

