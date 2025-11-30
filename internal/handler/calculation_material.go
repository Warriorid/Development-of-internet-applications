package handler

import (
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
)

// RemoveMaterialFromPit godoc
// @Summary Удаление материала из заявки
// @Description Удаление материала из заявки расчета (только владелец заявки)
// @Tags calculation-materials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param calculation_id path int true "ID заявки расчета"
// @Param material_id path int true "ID материала"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /calculation-materials/{calculation_id}/{material_id} [delete]
func (h *Handler) removeMaterialFromPit(c *gin.Context) {
	calculationID, err := strconv.Atoi(c.Param("calculation_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid calculation_id")
		return
	}
	materialID, err := strconv.Atoi(c.Param("material_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid material_id")
		return
	}
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	isOwner, err := h.service.IsPitOwner(calculationID, userID)
	if err != nil || !isOwner {
		newErrorResponse(c, http.StatusForbidden, "You can only modify your own pits")
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

// UpdateCalculationMaterial godoc
// @Summary Обновление материала в заявке
// @Description Обновление угла откоса материала в заявке расчета (только владелец заявки)
// @Tags calculation-materials
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param calculation_id path int true "ID заявки расчета"
// @Param material_id path int true "ID материала"
// @Param input body object true "Угол откоса"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /calculation-materials/{calculation_id}/{material_id} [put]
func (h *Handler) UpdateCalculationMaterial(c *gin.Context) {
	var input struct {
		SlopeAngle    int `json:"slope_angle"`
	}
	calculationID, err := strconv.Atoi(c.Param("calculation_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid calculation_id")
		return
	}
	materialID, err := strconv.Atoi(c.Param("material_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid material_id")
		return
	}
	userID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	isOwner, err := h.service.IsPitOwner(calculationID, userID)
	if err != nil || !isOwner {
		newErrorResponse(c, http.StatusForbidden, "You can only modify your own pits")
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid input data")
		return
	}
	if calculationID == 0 || materialID == 0 {
		newErrorResponse(c, http.StatusBadRequest, "id is required")
		return
	}	

	if err := h.service.CalculationMaterial.UpdateCalculationMaterial(calculationID, materialID, input.SlopeAngle); err != nil {
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