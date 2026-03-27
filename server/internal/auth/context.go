package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (string, error) {
	value, exists := c.Get(ContextUserIDKey)
	if !exists {
		return "", fmt.Errorf("authenticated user context missing")
	}

	userID, ok := value.(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("invalid authenticated user context")
	}

	return userID, nil
}
