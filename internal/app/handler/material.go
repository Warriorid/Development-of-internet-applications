package handler

import (
	"DIA/internal/app/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


func (h *Handler) GetMaterials(ctx *gin.Context) {
    var materials []model.Material
    var err error

    searchMaterialTitle := ctx.Query("materialTitle")
    if searchMaterialTitle == "" {          
        materials, err = h.Repository.GetMaterials()
        if err != nil {
            logrus.Error(err)
            materials = []model.Material{}
        }
    } else {
        materials, err = h.Repository.GetMaterialsByTitle(searchMaterialTitle)
        if err != nil {
            logrus.Error(err)
            materials = []model.Material{}
        }
    }

    pitID, err := h.Repository.GetDraftPitID()
	hasDraftPit := err == nil
    if err != nil {
        logrus.Println("Error getting draft pit ID:", err)
    }
    pitCount := h.Repository.GetPitCount(pitID)

    ctx.HTML(http.StatusOK, "materials.html", gin.H{
        "time":          time.Now().Format("15:04:05"),
        "materials":     materials,
        "materialTitle": searchMaterialTitle,
        "pitCount":      pitCount,
        "pitId":         pitID, 
		"hasActivePit": hasDraftPit,
    })
}

func (h *Handler) GetMaterial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error("Invalid material ID:", err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Неверный ID материала",
			"code":  400,
		})
		return
	}
	material, err := h.Repository.GetMaterial(id)
	if err != nil {
		logrus.Info("Material not found, ID:", id)
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Материал не найден",
			"code":  404,
		})
		return
	}
	c.HTML(http.StatusOK, "material_info.html", gin.H{
		"material": material,
	})
}

func (h *Handler) AddMaterialToPit(c *gin.Context) {
	materialId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error("uncorrenct materialId,", err)
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "неверный id материала",
			"code": 400,
		})
		return
	}
    
	err = h.Repository.AddMaterialToPit(materialId)
    if err != nil {
        logrus.Error("Error adding material to pit:", err)
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "error": "Ошибка при добавлении в заявку",
            "code":  500,
        })
        return
    }
	c.Redirect(http.StatusFound, "/materials")
}