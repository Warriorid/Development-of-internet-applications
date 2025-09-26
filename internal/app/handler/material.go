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

	exampleIdForPit := 1
	pit, pitCount, err := h.Repository.GetPit(exampleIdForPit)
	if err != nil {
		logrus.Error(err)
	}
	pitId := pit.ID


    ctx.HTML(http.StatusOK, "materials.html", gin.H{
        "time":      time.Now().Format("15:04:05"),
        "materials": materials,
        "materialTitle":     searchMaterialTitle,
        "pitCount": pitCount,
		"pitId": pitId,
    })
}

func (h *Handler) GetMaterial(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)
	}
	material, err := h.Repository.GetMaterial(id)
	if err != nil {
		logrus.Error(err)
	}
	c.HTML(http.StatusOK, "material_info.html", gin.H{
		"material": material,
	})
}