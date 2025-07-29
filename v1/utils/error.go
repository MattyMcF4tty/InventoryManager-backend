package utils

import (
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/MattyMcF4tty/InventoryManager-backend/v1/schemas"
)

func IsCustomError(err error) bool {
	if _, ok := err.(schemas.CustomError); ok {
		return true
	}
	if _, ok := err.(*schemas.CustomError); ok {
		return true
	}
	return false
}

func PostgresToHTTPError(err error) *int {
	// Check if the error is a Postgres error
	if strings.Contains(err.Error(), "PGRST") {
		// Extract the error code from the error message
		regex := regexp.MustCompile(`\(([^)]+)\)`)
		codeMatches := regex.FindStringSubmatch(err.Error())
		var errCode string
		if len(codeMatches) > 1 {
			errCode = codeMatches[1]
		}

		// The code handling is based of the PostgREST documentation
		// https://docs.postgrest.org/en/v12/references/errors.html

		postgresErrorMap := map[string]int{
			"PGRST100": http.StatusBadRequest,
			"PGRST102": http.StatusBadRequest,
			"PGRST108": http.StatusBadRequest,
			"PGRST116": http.StatusNotFound,
			"PGRST121": http.StatusNotFound,
		}

		// Check if the error code exists in the map
		if errCode != "" {
			if status, exists := postgresErrorMap[errCode]; exists {
				return &status
			} else {
				slog.Error("Postgres error code not found in map", "error_code", errCode, "error", err.Error())
			}
		} else {
			slog.Error("Postgres error code not found in postgres error message", "error", err.Error())
		}

		status := http.StatusInternalServerError
		return &status
	}

	return nil
}
