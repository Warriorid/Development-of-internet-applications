package handler

import (
	"DIA/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)



func (h *Handler) RegisterUser(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user := &model.Users{
		Username: req.Username,
		Password: req.Password,
		IsStaff:  false,
	}
	createdUser, err := h.service.User.Register(user)
	if err != nil {
		if err.Error() == "username already exists" {
			newErrorResponse(c, http.StatusConflict, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusCreated, createdUser)
}

func (h *Handler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.service.User.Login(req.Username, req.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			newErrorResponse(c, http.StatusUnauthorized, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetUserProfile(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.User.GetUserProfile(userID)
	if err != nil {
		if err.Error() == "user not found" {
			newErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUserProfile(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	var req model.UpdateProfileRequest
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.service.User.UpdateUserProfile(userID, req.Username, req.Password)
	if err != nil {
		if err.Error() == "user not found" {
			newErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) Logout(c *gin.Context) {
	err := h.service.User.Logout()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}