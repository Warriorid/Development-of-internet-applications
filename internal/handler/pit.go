package handler

import (
	"DIA/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)


func (h *Handler) GetDraftPit(c *gin.Context) {
	creatorID := 2 

	pitID, itemsCount, err := h.service.GetDraftPitWithItemsCount(creatorID)
	if err != nil {
		if err.Error() == "pit not found" {
			newErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"pit_id": pitID,
		"pits_count": itemsCount,
	})
}

func (h *Handler) GetPits(c *gin.Context) {
	statusFilter := c.Query("status")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

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


func (h *Handler) GetPitWithMaterials(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
        return
    }
    creatorId := 2

    pit, err := h.service.GetPitWithMaterials(id, creatorId)
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

func (h *Handler) UpdatePit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
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

func (h *Handler) FormPit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}
	creatorId := 2
	if err := h.service.FormPit(id, creatorId); err != nil {
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

func (h *Handler) CompletePit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}
	var status string
	if err := c.BindJSON(&status); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid status")
		return
	}
	if status != "completed" && status != "rejected" {
		newErrorResponse(c, http.StatusBadRequest, "status must be 'completed' or 'rejected'")
		return
	}
	moderatorID := 1
	if err := h.service.CompletePit(id, moderatorID, status); err != nil {
		switch {
		case err.Error() == "not found":
			newErrorResponse(c, http.StatusNotFound, "pit not found")
		case err.Error() == "can only complete or reject pits with formed status":
			newErrorResponse(c, http.StatusConflict, err.Error())
		case err.Error() == "invalid status: must be 'completed' or 'rejected'":
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": status})
}

func (h *Handler) DeletePit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id parameter")
		return
	}
	if err := h.service.DeletePit(id); err != nil {
		switch {
		case err.Error() == "not found":
			newErrorResponse(c, http.StatusNotFound, "pit not found")
		default:
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}