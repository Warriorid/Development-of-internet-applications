package handler

import (
	"DIA/internal/model"
	"DIA/internal/role"
	"DIA/pkg/config"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterUser godoc
// @Summary Регистрация пользователя
// @Description Создание нового пользователя
// @Tags users
// @Accept json
// @Produce json
// @Param input body model.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} model.Users
// @Failure 400 {object} errorResponse
// @Failure 409 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users [post]
func (h *Handler) RegisterUser(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userRole := 0
	if req.Role != nil {
		userRole = *req.Role
	}
	user := &model.Users{
		Username: req.Username,
		Password: req.Password,
		Role:  role.Role(userRole),
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

// Login godoc
// @Summary Авторизация пользователя
// @Description Вход в систему и получение JWT токена
// @Tags users
// @Accept json
// @Produce json
// @Param input body model.LoginRequest true "Данные для входа"
// @Success 200 {object} model.LoginResp
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/login [post]
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

// GetUserProfile godoc
// @Summary Получение профиля пользователя
// @Description Получение информации о пользователе по ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID пользователя"
// @Success 200 {object} model.Users
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/{id} [get]
func (h *Handler) GetUserProfile(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	authUserID, err := getUserID(c)
	if err != nil || authUserID != userID {
		newErrorResponse(c, http.StatusForbidden, "Access denied")
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

// UpdateUserProfile godoc
// @Summary Обновление профиля пользователя
// @Description Обновление данных пользователя
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID пользователя"
// @Param input body model.UpdateProfileRequest true "Данные для обновления"
// @Success 200 {object} model.Users
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/{id} [put]
func (h *Handler) UpdateUserProfile(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	authUserID, err := getUserID(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	if authUserID != userID {
		newErrorResponse(c, http.StatusForbidden, "You can only update your own profile")
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

// Logout godoc
// @Summary Выход из системы
// @Description Выход пользователя и добавление токена в blacklist
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Example: {\"message\": \"logged out successfully\"}"
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	token, err := GetToken(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Token not found")
		return
	}
	claims, err := config.ParseToken(token)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid token")
		return
	}
	var expirationTime time.Duration
	
	if claims.ExpiresAt != nil {
		expirationTime = time.Until(claims.ExpiresAt.Time)
	} else {
		expirationTime = 24 * time.Hour
	}
	if expirationTime <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "logged out successfully (token already expired)",
		})
		return
	}
	err = h.addTokenToBlacklist(token, expirationTime)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to logout")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}
