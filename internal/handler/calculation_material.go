package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) removeMaterialFromPit(c *gin.Context) {
	calculationID, err := strconv.Atoi(c.Query("calculation_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid calculation_id")
		return
	}
	materialID, err := strconv.Atoi(c.Query("material_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid material_id")
		return
	}

	if err := h.service.CalculationMaterial.RemoveMaterialFromPit(calculationID, materialID); err != nil {
		if err.Error() == "material not found in pit" {
			newErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "material removed from pit successfully",
	})
}


func (h *Handler) UpdateCalculationMaterial(c *gin.Context) {
	var input struct {
		CalculationID int `json:"calculation_id"`
		MaterialID    int `json:"material_id"`
		SlopeAngle    int `json:"slope_angle"`
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input data")
		return
	}
	if input.CalculationID == 0 || input.MaterialID == 0 {
		newErrorResponse(c, http.StatusBadRequest, "id is required")
		return
	}	

	if err := h.service.CalculationMaterial.UpdateCalculationMaterial(input.CalculationID, input.MaterialID, input.SlopeAngle,); err != nil {
		switch {
		case err.Error() == "material not found in pit":
			newErrorResponse(c, http.StatusNotFound, err.Error())
		case err.Error() == "slope angle must be between 0 and 45 degrees":
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "slope angle updated successfully",
	})
}