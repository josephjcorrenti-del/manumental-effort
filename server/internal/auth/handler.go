package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	result, err := h.service.Login(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid email or password",
			})
			return
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) Me(c *gin.Context) {
	userIDValue, exists := c.Get(ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "authenticated user context missing",
		})
		return
	}

	userID, ok := userIDValue.(string)
	if !ok || userID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid authenticated user context",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
	})
}
