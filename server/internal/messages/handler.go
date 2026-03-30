package messages

import (
	"errors"
	"net/http"
	"strconv"

	"manumental-effort/server/internal/auth"
	"manumental-effort/server/internal/channels"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateMessage(c *gin.Context) {
	channelID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid channel id"})
		return
	}

	userIDHex, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authenticated user id"})
		return
	}

	var input CreateMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	message, err := h.service.CreateMessage(c.Request.Context(), channelID, userID, input.Body)
	if err != nil {
		switch {
		case errors.Is(err, ErrMessageBodyRequired):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, ErrMessageBodyTooLong):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, channels.ErrChannelNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, ErrMembershipRequired):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create message"})
		}
		return
	}

	c.JSON(http.StatusCreated, message)
}

func (h *Handler) ListMessages(c *gin.Context) {
	channelID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid channel id"})
		return
	}

	userIDHex, err := auth.GetUserID(c)

	userID, err := primitive.ObjectIDFromHex(userIDHex)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	limit := DefaultMessageListLimit
	if limitValue := c.Query("limit"); limitValue != "" {
		parsedLimit, err := strconv.Atoi(limitValue)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}
		limit = parsedLimit
	}

	var beforeID *primitive.ObjectID
	if beforeValue := c.Query("before"); beforeValue != "" {
		parsedBeforeID, err := primitive.ObjectIDFromHex(beforeValue)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid before cursor"})
			return
		}
		beforeID = &parsedBeforeID
	}

	result, err := h.service.ListMessages(c.Request.Context(), channelID, userID, beforeID, limit)
	if err != nil {
		switch {
		case errors.Is(err, channels.ErrChannelNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, ErrMembershipRequired):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list messages"})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}
