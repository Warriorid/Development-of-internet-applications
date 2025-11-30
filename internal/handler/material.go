package handler

import (
	"DIA/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetMaterials godoc
// @Summary Получение списка материалов
// @Description Получение списка всех материалов с возможностью фильтрации по названию
// @Tags materials
// @Accept json
// @Produce json
// @Param materialTitle query string false "Фильтр по названию материала"
// @Success 200 {array} model.Material
// @Failure 500 {object} errorResponse
// @Router /materials [get]
func (h *Handler) GetMaterials(c *gin.Context) {
	materialTitle := c.Query("materialTitle")
	
	materials, err := h.service.Material.GetMaterials(materialTitle)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	c.JSON(http.StatusOK, materials)
}

// GetMaterial godoc
// @Summary Получение материала по ID
// @Description Получение информации о конкретном материале
// @Tags materials
// @Accept json
// @Produce json
// @Param id path int true "ID материала"
// @Success 200 {object} model.Material
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Router /materials/{id} [get]
func (h *Handler) GetMaterial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid material ID")
		return
	}
	
	material, err := h.service.GetMaterial(id)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, "Material not found")
		return
	}
	
	c.JSON(http.StatusOK, material)
}

// CreateMaterial godoc
// @Summary Создание нового материала
// @Description Создание материала (только для модераторов)
// @Tags materials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body model.Material true "Данные материала"
// @Success 201 {object} model.Material
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /materials [post]
func (h *Handler) CreateMaterial(c *gin.Context) {
	role, err := getUserRole(c)
	if err != nil || role != 1 { 
		newErrorResponse(c, http.StatusForbidden, "Access denied: moderator role required")
		return
	}
	var material model.Material
	if err := c.BindJSON(&material); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}
	
	createdMaterial, err := h.service.Material.CreateMaterial(&material)
    if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	c.JSON(http.StatusCreated, createdMaterial)
}

// UpdateMaterial godoc
// @Summary Обновление материала
// @Description Обновление информации о материале (только для модераторов)
// @Tags materials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID материала"
// @Param input body model.Material true "Данные для обновления"
// @Success 200 {object} model.Material
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /materials/{id} [put]
func (h *Handler) UpdateMaterial(c *gin.Context) {
	role, err := getUserRole(c)
	if err != nil || role != 1 {
		newErrorResponse(c, http.StatusForbidden, "Access denied: moderator role required")
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid material ID")
		return
	}
	
	var material model.Material
	if err := c.BindJSON(&material); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid input data")
		return
	}
	
	if err := h.service.Material.UpdateMaterial(id, &material); err != nil {
		if err.Error() == "material not found" {
            newErrorResponse(c, http.StatusNotFound, err.Error())
        } else {
            newErrorResponse(c, http.StatusInternalServerError, err.Error())
        }
        return
	}
	
	c.JSON(http.StatusOK, material)
}

// DeleteMaterial godoc
// @Summary Удаление материала
// @Description Удаление материала (только для модераторов)
// @Tags materials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID материала"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /materials/{id} [delete]
func (h *Handler) DeleteMaterial(c *gin.Context) {
	role, err := getUserRole(c)
	if err != nil || role != 1 {
		newErrorResponse(c, http.StatusForbidden, "Access denied: moderator role required")
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid material ID")
		return
	}

	
	if err := h.service.Material.DeleteMaterial(id); err != nil {
		if err.Error() == "material not found" {
            newErrorResponse(c, http.StatusNotFound, err.Error())
        } else {
            newErrorResponse(c, http.StatusInternalServerError, err.Error())
        }
        return
	}
	
	c.Status(http.StatusOK)
}

// AddMaterialToPit godoc
// @Summary Добавление материала в заявку
// @Description Добавление материала в текущую заявку пользователя
// @Tags materials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID материала"
// @Success 201
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /materials/{id}/pit [post]
func (h *Handler) AddMaterialToPit(c *gin.Context) {
	userId, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	materialId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid material ID")
		return
	}
	
	if err := h.service.AddMaterialToPit(userId, materialId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	c.Status(http.StatusCreated)
}

// UploadMaterialImage godoc
// @Summary Загрузка изображения материала
// @Description Загрузка изображения для материала (только для модераторов)
// @Tags materials
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID материала"
// @Param image formData file true "Изображение материала"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /materials/{id}/image [post]
func (h *Handler) UploadMaterialImage(c *gin.Context) {
	role, err := getUserRole(c)
	if err != nil || role != 1 {
		newErrorResponse(c, http.StatusForbidden, "Access denied: moderator role required")
		return
	}
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        newErrorResponse(c, http.StatusBadRequest, "Invalid material ID")
        return
    }
    
    file, header, err := c.Request.FormFile("image")
    if err != nil {
        newErrorResponse(c, http.StatusBadRequest, "No image file provided")
        return
    }
    defer file.Close()
    
    fileBytes := make([]byte, header.Size)
    _, err = file.Read(fileBytes)
    if err != nil {
        newErrorResponse(c, http.StatusInternalServerError, "Failed to read file")
        return
    }
    
    if err := h.service.UploadMaterialImage(id, fileBytes, header.Filename); err != nil {
        newErrorResponse(c, http.StatusInternalServerError, err.Error())
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Image uploaded successfully",
    })
}