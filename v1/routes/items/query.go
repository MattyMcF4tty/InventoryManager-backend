package items

// ---- WILL NOT BE USED ----
// This file was intended to handle querying items based on various parameters.
// However, it will be replaced by a more robust query system later on.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MattyMcF4tty/InventoryManager-backend/v1/database"
	"github.com/MattyMcF4tty/InventoryManager-backend/v1/schemas"
	"github.com/MattyMcF4tty/InventoryManager-backend/v1/utils"
	"github.com/supabase-community/postgrest-go"
)

// We define the possible query parameters
var allowedQueryParameters = map[string]struct{}{
	"name":        {},
	"description": {},
	"price":       {},
	"quantity":    {},
	"category":    {},
	"location":    {},
	"status":      {},

	"min_price":    {},
	"max_price":    {},
	"min_quantity": {},
	"max_quantity": {},
	"sort_by":      {},
	"sort_order":   {},
}

func queryItems(queryParams map[string][]string) ([]schemas.Item, *int64, error) {

	client := database.Connect()
	query := client.From("items").Select("*", "exact", false)

	for key, values := range queryParams {
		if _, isAllowed := allowedQueryParameters[key]; !isAllowed {
			// We skip unknown parameters
			continue
		}

		if len(values) == 0 {
			// We skip keys with no values
			continue
		}

		// Handle special cases for min_ and max_ filters
		if column, found := strings.CutPrefix(key, "min_"); found {
			// We validate that the value is a valid float64
			value := values[0]
			_, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, nil, &schemas.CustomError{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("Invalid value for %s: %s", key, value),
					Details: fmt.Sprintf("Item query failed. Expected a float64 number for %s, got %s", key, value),
				}
			}

			query = query.Gte(column, value)
			continue
		} else if column, found := strings.CutPrefix(key, "max_"); found {
			// We validate that the value is a valid float64
			value := values[0]
			_, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, nil, &schemas.CustomError{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("Invalid value for %s: %s", key, value),
					Details: fmt.Sprintf("Item query failed. Expected a float64 number for %s, got %s", key, value),
				}
			}

			query = query.Lte(column, values[0])
			continue
		}

		if key == "sort_by" {
			sortOrders := queryParams["sort_order"]
			for index, value := range values {
				// We validate that the value is a valid field name
				if _, isAllowed := allowedQueryParameters[value]; !isAllowed {
					return nil, nil, &schemas.CustomError{
						Code:    http.StatusBadRequest,
						Message: fmt.Sprintf("Invalid sort field: %s", value),
						Details: fmt.Sprintf("Item query failed. Expected a valid field name for sorting, got %s", value),
					}
				}

				// Check if there is a corresponding sort order. If not, we default to ascending order
				asc := true
				if index < len(sortOrders) && strings.ToLower(sortOrders[index]) == "desc" {
					asc = false
				}
				query = query.Order(value, &postgrest.OrderOpts{Ascending: asc})
			}
			continue
		}

		// For other parameters, we assume they are equality filters
		for _, value := range values {
			// we apply the filter to the query
			query = query.Eq(key, value)
		}
	}

	itemsData, count, err := query.Execute()
	if err != nil {
		// Set the default error code and message
		code := http.StatusInternalServerError
		message := "An error occurred while querying items"
		if status := utils.PostgresToHTTPError(err); status != nil {
			code = *status

			if code == http.StatusNotFound {
				message = "No items found matching query"
			}
		}

		return []schemas.Item{}, nil, &schemas.CustomError{
			Code:    code,
			Message: message,
			Details: fmt.Sprintf("Item query failed: %v", err),
		}
	}

	var items []schemas.Item
	err = json.Unmarshal(itemsData, &items)
	if err != nil {
		return []schemas.Item{}, nil, &schemas.CustomError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse items data",
			Details: fmt.Sprintf("Failed to parse data from item query failed: %v", err),
		}
	}

	return items, &count, nil
}
