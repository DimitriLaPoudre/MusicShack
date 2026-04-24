package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetFromContext[T any](c *gin.Context, key string) (T, error) {
	var empty T

	value, exists := c.Get(key)
	if !exists {
		return empty, fmt.Errorf("GetFromContext: key %q not found in context", key)
	}

	typedValue, ok := value.(T)
	if !ok {
		return empty, fmt.Errorf("GetFromContext: key %q not found in context", key)
	}

	return typedValue, nil
}
