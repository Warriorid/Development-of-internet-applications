package handler

import (
	"DIA/internal/app/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct{
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetMaterials(ctx *gin.Context) {
    var materials []repository.Material
    var err error

    searchQuery := ctx.Query("query")
    if searchQuery == "" {          
        materials, err = h.Repository.GetMaterials()
        if err != nil {
            logrus.Error(err)
            materials = []repository.Material{}
        }
    } else {
        materials, err = h.Repository.GetMaterialsByTitle(searchQuery)
        if err != nil {
            logrus.Error(err)
            materials = []repository.Material{}
        }
    }

	exampleIdForCart := 1
	cart, cartSize, err := h.Repository.GetCart(exampleIdForCart)
	if err != nil {
		logrus.Error(err)
	}
	cartId := cart.ID


    ctx.HTML(http.StatusOK, "index.html", gin.H{
        "time":      time.Now().Format("15:04:05"),
        "materials": materials,
        "query":     searchQuery,
        "cartCount": cartSize,
		"cartId": cartId,
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
	c.HTML(http.StatusOK, "material.html", gin.H{
		"material": material,
	})
}

func (h *Handler) GetCart(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Error(err)
	}
	cart, _, err := h.Repository.GetCart(id)
	if err != nil {
		logrus.Error(err)
	}

	c.HTML(http.StatusOK, "cart.html", gin.H{
		"cartItems": cart.CartItems,
		"volume": 0.0,
	})
}