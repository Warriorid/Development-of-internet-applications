package handler

import (
	"DIA/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)



func (h *Handler) GetMaterials(c *gin.Context) {
	materialTitle := c.Query("materialTitle")
	
	materials, err := h.service.Material.GetMaterials(materialTitle)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	c.JSON(http.StatusOK, materials)
}

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

func (h *Handler) CreateMaterial(c *gin.Context) {
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

func (h *Handler) UpdateMaterial(c *gin.Context) {
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

func (h *Handler) DeleteMaterial(c *gin.Context) {
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

func (h *Handler) AddMaterialToPit(c *gin.Context) {
	userId := 5
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

func (h *Handler) UploadMaterialImage(c *gin.Context) {
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