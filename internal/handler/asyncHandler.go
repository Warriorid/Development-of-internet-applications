package handler

import (
	"DIA/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CompleteAsyncCalculation godoc
// @Summary Завершение асинхронного расчета
// @Description Прием результатов расчета от асинхронного сервиса
// @Tags async
// @Accept json
// @Produce json
// @Param input body model.AsyncCalculationResult true "Результат расчета"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /calculations/complete [post]
func (h *Handler) CompleteAsyncCalculation(c *gin.Context) {
    var input model.AsyncCalculationResult
    
    if err := c.BindJSON(&input); err != nil {
        newErrorResponse(c, http.StatusBadRequest, err.Error())
        return
    }
    
    if err := h.service.Pits.CompleteAsyncCalculation(input); err != nil {
        if err.Error() == "unauthorized: invalid token" {
            newErrorResponse(c, http.StatusUnauthorized, err.Error())
        } else {
            newErrorResponse(c, http.StatusInternalServerError, err.Error())
        }
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"status": "result updated"})
}