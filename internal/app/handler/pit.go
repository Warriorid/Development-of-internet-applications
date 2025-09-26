package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetPit(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)
	}
	pit, _, err := h.Repository.GetPit(id)
	if err != nil {
		logrus.Error(err)
	}

	c.HTML(http.StatusOK, "pits.html", gin.H{
		"pit": pit,
	})
}