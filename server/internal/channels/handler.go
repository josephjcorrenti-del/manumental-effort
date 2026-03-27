package channels

import (
	"net/http"

	"manumental-effort/server/internal/auth"

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

func (h *Handler) CreateChannel(c *gin.Context) {
	var input CreateChannelInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	spaceID := c.Param("id")
	if spaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "space id is required",
		})
		return
	}

	channel, err := h.service.CreateChannel(c.Request.Context(), userID, spaceID, input)
	if err != nil {
		if err.Error() == "membership required" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "membership required",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, channel)
}

func (h *Handler) ListChannels(c *gin.Context) {
	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	spaceID := c.Param("id")
	if spaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "space id is required",
		})
		return
	}

	channels, err := h.service.ListChannels(c.Request.Context(), userID, spaceID)
	if err != nil {
		if err.Error() == "membership required" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "membership required",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, channels)
}
