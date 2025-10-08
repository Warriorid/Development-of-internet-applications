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
        logrus.Error("Invalid pit ID:", err)
        return
    }
    pit, err := h.Repository.GetPit(id)
    if err != nil {
        logrus.Info("Draft pit not found, ID:", id)
    }

    c.HTML(http.StatusOK, "pits.html", gin.H{
        "pit": pit,
    })
}

func (h *Handler) DeletePit(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error("Invalid pit ID:", err)
        return
	}
	if err := h.Repository.DeletePit(id); err != nil {
		logrus.Error("error of deleting pit:", err)
		return
	}
	c.Redirect(http.StatusFound, "/materials")
}