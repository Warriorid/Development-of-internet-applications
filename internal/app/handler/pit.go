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
        c.HTML(http.StatusBadRequest, "error.html", gin.H{
            "error": "Неверный ID котлована",
            "code":  400,
        })
        return
    }
    pit, err := h.Repository.GetPit(id)
    if err != nil {
        logrus.Info("Draft pit not found, ID:", id)
        c.HTML(http.StatusNotFound, "error.html", gin.H{
            "error": "Котлован не найден",
            "code":  404,
        })
        return
    }

    c.HTML(http.StatusOK, "pits.html", gin.H{
        "pit": pit,
    })
}

func (h *Handler) DeletePit(c *gin.Context){
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error("Invalid pit ID:", err)
        c.HTML(http.StatusBadRequest, "error.html", gin.H{
            "error": "Неверный ID котлована",
            "code":  400,
        })
        return
	}
	if err := h.Repository.DeletePit(id); err != nil {
		logrus.Error("error of deleting pit:", err)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err,
			"code": 500,
		})
		return
	}
	c.Redirect(http.StatusFound, "/materials")
}