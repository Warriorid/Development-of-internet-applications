package handler

import (
	"DIA/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetDraftPit godoc
// @Summary Получение черновика заявки
// @Description Получение ID черновика заявки пользователя и количества материалов в ней. Для гостей возвращает -1.
// @Tags pits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object "Для авторизованных: {pit_id, pits_count}, для гостей: -1"
// @Failure 401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /pits/draft [get]
func (h *Handler) GetDraftPit(c *gin.Context) {
	userId, err := getUserID(c)
    if err != nil {
        c.JSON(http.StatusOK, -1)
        return
    }
    role, roleErr := getUserRole(c)
    if roleErr != nil || role == 2 {
        c.JSON(http.StatusOK, -1)
        return
    }
	pitID, itemsCount, err := h.service.GetDraftPitWithItemsCount(userId)
	if err != nil {
		if err.Error() == "pit not found" {
			newErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"pit_id":     pitID,
		"pits_count": itemsCount,
	})
}

// GetPits godoc
// @Summary Получение списка заявок
// @Description Получение списка всех заявок с фильтрацией (только для модераторов)
// @Tags pits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Фильтр по статусу"
// @Param start_date query string false "Начальная дата (формат: YYYY-MM-DD)"
// @Param end_date query string false "Конечная дата (формат: YYYY-MM-DD)"
// @Success 200 {array} model.PitsCalculationListItem
// @Failure 403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /pits [get]
func (h *Handler) GetPits(c *gin.Context) {
	statusFilter := c.Query("status")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	role, err := getUserRole(c)
    if err != nil || role != 1 {
        newErrorResponse(c, http.StatusForbidden, "Access denied: moderator role required")
        return
    }
	var startDate, endDate *time.Time

	if startDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = &parsed
		}
	}

	if endDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = &parsed
		}
	}

	pits, err := h.service.GetPits(statusFilter, startDate, endDate)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, pits)
}

// GetPitWithMaterials godoc
// @Summary Получение заявки с материалами
// @Description Получение детальной информации о заявке включая материалы
// @Tags pits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Success 200 {object} model.PitsCalculationWithMaterials
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /pits/{id} [get]
func (h *Handler) GetPitWithMaterials(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
        return
    }
    userId, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

    pit, err := h.service.GetPitWithMaterials(id, userId)
    if err != nil {
        if err.Error() == "not found" {
            newErrorResponse(c, http.StatusNotFound, "pit not found")
        } else {
            newErrorResponse(c, http.StatusInternalServerError, err.Error())
        }
        return
    }

    c.JSON(http.StatusOK, pit)
}

// UpdatePit godoc
// @Summary Обновление заявки
// @Description Обновление данных заявки (только владелец заявки)
// @Tags pits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Param input body model.UpdatePitParam true "Данные для обновления"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /pits/{id} [put]
func (h *Handler) UpdatePit(c *gin.Context) {
	userID, err := getUserID(c)
    if err != nil {
        newErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}
	isOwner, err := h.service.IsPitOwner(id, userID)
	if err != nil || !isOwner {
		newErrorResponse(c, http.StatusForbidden, "You can only update your own pits")
		return
	}
	var input model.UpdatePitParam
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdatePit(id, input); err != nil {
		if err.Error() == "not found" {
			newErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// FormPit godoc
// @Summary Формирование заявки
// @Description Перевод заявки из черновика в статус "formed" (только владелец заявки)
// @Tags pits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /pits/{id}/form [put]
func (h *Handler) FormPit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}
	userId, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	if err := h.service.FormPit(id, userId); err != nil {
		switch {
		case err.Error() == "not found":
			newErrorResponse(c, http.StatusNotFound, "pit not found")
		case err.Error() == "pit must have at least one material":
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		case err.Error() == "all materials must have slope angle specified":
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		case err.Error() == "pit dimensions must be greater than default values":
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		case err.Error() == "pit cannot be formed":
			newErrorResponse(c, http.StatusConflict, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "formed"})
}

// CompletePit godoc
// @Summary Завершение/отклонение заявки
// @Description Перевод заявки в статус "completed" или "rejected" (только для модераторов)
// @Tags pits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Param input body string true "Статус: 'completed' или 'rejected'"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /pits/{id}/complete [put]
func (h *Handler) CompletePit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}
	var request struct {
		Status string `json:"status"`
	}
	if err := c.BindJSON(&request); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid status")
		return
	}
	if request.Status != "completed" && request.Status != "rejected" {
		newErrorResponse(c, http.StatusBadRequest, "status must be 'completed' or 'rejected'")
		return
	}
	moderatorID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	role, err := getUserRole(c)
	if err != nil || role != 1 {
		newErrorResponse(c, http.StatusForbidden, "Access denied: moderator role required")
		return
	}
	if err := h.service.CompletePit(id, moderatorID, request.Status); err != nil {
		switch {
		case err.Error() == "not found":
			newErrorResponse(c, http.StatusNotFound, "pit not found")
		case err.Error() == "can only complete or reject pits with formed status":
			newErrorResponse(c, http.StatusConflict, err.Error())
		case err.Error() == "invalid status: must be 'completed' или 'rejected'":
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": request})
}

// DeletePit godoc
// @Summary Удаление заявки
// @Description Удаление заявки (только для черновиков и только владелец заявки или модератор)
// @Tags pits
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID заявки"
// @Success 200
// @Failure 400 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /pits/{id} [delete]
func (h *Handler) DeletePit(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
        return
    }
    
    userID, err := getUserID(c)
    if err != nil {
        newErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }
    
    role, err := getUserRole(c)
    if err != nil {
        newErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
        return
    }

    canDelete, err := h.service.CanUserDeletePit(id, userID, role)
    if err != nil || !canDelete {
        newErrorResponse(c, http.StatusForbidden, "Access denied: you can only delete your own draft pits")
        return
    }
    
    if err := h.service.DeletePit(id); err != nil {
        switch {
        case err.Error() == "pit not found":
            newErrorResponse(c, http.StatusNotFound, "pit not found")
        case err.Error() == "only draft pits can be deleted":
            newErrorResponse(c, http.StatusForbidden, err.Error())
        default:
            newErrorResponse(c, http.StatusInternalServerError, err.Error())
        }
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}