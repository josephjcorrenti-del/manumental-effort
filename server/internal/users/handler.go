package users

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

func (h *Handler) CreateUser(c *gin.Context) {
	var input CreateUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, ErrDuplicateUsername):
			c.JSON(http.StatusConflict, gin.H{
				"error": "username is already taken",
			})
			return
		case errors.Is(err, ErrDuplicateEmail):
			c.JSON(http.StatusConflict, gin.H{
				"error": "email is already in use",
			})
			return
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.service.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
