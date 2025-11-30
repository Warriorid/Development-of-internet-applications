package handler

import (
	"DIA/pkg/config"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		
		if authHeader == "" {
			c.Set("userRole", 2) 
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			newErrorResponse(c, http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
			return
		}

		tokenString := parts[1]

		if h.isTokenBlacklisted(tokenString) {
			newErrorResponse(c, http.StatusUnauthorized, "Token has been invalidated")
			return
		}
		
		claims, err := config.ParseToken(tokenString)
		if err != nil {
			newErrorResponse(c, http.StatusUnauthorized, "Invalid token")
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Set("token", tokenString)

		c.Next()
	}
}

func (h *Handler) isTokenBlacklisted(token string) bool {
	if h.redis == nil {
		return false
	}

	ctx := context.Background()
	exists, err := h.redis.Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		return false
	}
	return exists > 0
}

func (h *Handler) addTokenToBlacklist(token string, expirationTime time.Duration) error {
	if h.redis == nil {
		return nil
	}

	ctx := context.Background()
	err := h.redis.Set(ctx, "blacklist:"+token, "1", expirationTime).Err()
	return err
}

func getUserID(c *gin.Context) (int, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, gin.Error{}
	}

	id, ok := userID.(int)
	if !ok {
		return 0, gin.Error{}
	}

	return id, nil
}

func getUserRole(c *gin.Context) (int, error) {
	role, exists := c.Get("userRole")
	if !exists {
		return 0, gin.Error{}
	}

	r, ok := role.(int)
	if !ok {
		return 0, gin.Error{}
	}

	return r, nil
}

func GetToken(c *gin.Context) (string, error) {
	token, exists := c.Get("token")
	if !exists {
		return "", gin.Error{}
	}

	t, ok := token.(string)
	if !ok {
		return "", gin.Error{}
	}

	return t, nil
}