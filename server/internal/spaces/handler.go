package spaces

import (
	"errors"
	"net/http"

	"manumental-effort/server/internal/auth"
	"manumental-effort/server/internal/memberships"

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

func (h *Handler) CreateSpace(c *gin.Context) {
	var input CreateSpaceInput

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

	space, err := h.service.CreateSpace(c.Request.Context(), userID, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, space)
}

func (h *Handler) JoinSpace(c *gin.Context) {
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

	err = h.service.JoinSpace(c.Request.Context(), userID, spaceID)
	if err != nil {
		switch {
		case errors.Is(err, memberships.ErrDuplicateMembership):
			c.JSON(http.StatusConflict, gin.H{
				"error": "already a member of this space",
			})
			return
		case err.Error() == "space not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error": "space not found",
			})
			return
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "joined",
	})
}

func (h *Handler) ListSpaces(c *gin.Context) {
	userID, err := auth.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	spaces, err := h.service.ListSpaces(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to list spaces",
		})
		return
	}

	c.JSON(http.StatusOK, spaces)
}
