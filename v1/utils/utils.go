package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetIdFromContext(context *gin.Context) (int8, error) {
	idStr := context.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 8)
	if err != nil {
		return 0, fmt.Errorf("invalid ID: %s", idStr)
	}
	return int8(id), nil
}

func GetCurrentISODate() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func RemoveProtectedFields(itemMap map[string]interface{}, protectedFields []string) {
	for _, field := range protectedFields {
		delete(itemMap, field)
	}
}

func CheckRequiredFields(itemMap map[string]interface{}, requiredFields []string) error {
	for _, field := range requiredFields {
		if _, exists := itemMap[field]; !exists {
			return fmt.Errorf("missing required field: %s", field)
		}
	}
	return nil
}

func InRange(value, min, max int) bool {
	return value >= min && value <= max
}
