package items

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	db "github.com/MattyMcF4tty/InventoryManager-backend/v1/database"
	"github.com/MattyMcF4tty/InventoryManager-backend/v1/schemas"
	"github.com/MattyMcF4tty/InventoryManager-backend/v1/utils"
	"github.com/supabase-community/postgrest-go"
)

func GetItem(id int8) (schemas.Item, error) {
	client := db.Connect()
	idStr := fmt.Sprintf("%d", id)

	data, _, err := client.
		From("items").
		Select("*", "", false).
		Eq("id", idStr).
		Is("deleted_at", "null").
		Single().
		Execute()

	if err != nil {
		// Set the default error code and message
		code := http.StatusInternalServerError
		message := "An error occurred while retrieving the item"

		// Check if the error is a Postgres error
		// If true we update the code and message accordingly
		if status := utils.PostgresToHTTPError(err); status != nil {
			code = *status

			if code == http.StatusNotFound {
				message = "Item not found"
			}
		}

		return schemas.Item{}, &schemas.CustomError{
			Code:    code,
			Message: message,
			Details: fmt.Sprintf("Error retrieving item with ID %d: %v", id, err),
		}
	}

	var item schemas.Item
	err = json.Unmarshal(data, &item)
	if err != nil {
		return schemas.Item{}, &schemas.CustomError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse item data",
			Details: fmt.Sprintf("Error parsing item data for ID %d: %v", id, err),
		}
	}

	return item, nil
}

func UpdateItem(id int8, updates map[string]interface{}) (schemas.Item, error) {
	client := db.Connect()
	idStr := fmt.Sprintf("%d", id)

	// Add updated_at field
	updates["updated_at"] = utils.GetCurrentISODate()

	data, _, err := client.
		From("items").
		Update(updates, "", "").
		Eq("id", idStr).
		Is("deleted_at", "null").
		Single().Execute()

	if err != nil {
		// Set the default error code and message
		code := http.StatusInternalServerError
		message := "An error occurred while updating the item"

		// Check if the error is a Postgres error
		// If true we update the code and message accordingly
		if status := utils.PostgresToHTTPError(err); status != nil {
			code = *status

			if code == http.StatusNotFound {
				message = "Item not found"
			}
		}

		return schemas.Item{}, &schemas.CustomError{
			Code:    code,
			Message: message,
			Details: fmt.Sprintf("Error updating item with ID %d: %v", id, err),
		}
	}

	var item schemas.Item
	err = json.Unmarshal(data, &item)

	if err != nil {
		return schemas.Item{}, &schemas.CustomError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse item data",
			Details: fmt.Sprintf("Error parsing item data for ID %d: %v", id, err),
		}
	}

	return item, nil
}

func CreateItem(item schemas.Item) (schemas.Item, error) {
	client := db.Connect()

	item.CreatedAt = utils.GetCurrentISODate()
	item.UpdatedAt = utils.GetCurrentISODate()

	data, _, err := client.
		From("items").
		Insert(item, false, "", "", "").
		Single().
		Execute()

	if err != nil {
		// Set the default error code and message
		code := http.StatusInternalServerError
		message := "An error occurred while creating the item"

		// Check if the error is a Postgres error
		// If true we update the code and message accordingly
		if status := utils.PostgresToHTTPError(err); status != nil {
			code = *status

			if code == http.StatusNotFound {
				message = "Item not found"
			}
		}

		return schemas.Item{}, &schemas.CustomError{
			Code:    code,
			Message: message,
			Details: fmt.Sprintf("Error creating item: %v", err),
		}
	}

	var createdItem schemas.Item
	err = json.Unmarshal(data, &createdItem)

	if err != nil {
		return schemas.Item{}, &schemas.CustomError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse item data",
			Details: fmt.Sprintf("Error parsing item data while creating item: %v", err),
		}
	}

	return createdItem, nil
}

func DeleteItem(id int8) error {
	client := db.Connect()
	idStr := fmt.Sprintf("%d", id)

	_, _, err := client.
		From("items").
		Delete("", "").
		Eq("id", idStr).
		Is("deleted_at", "null").
		Execute()

	if err != nil {
		// Set the default error code and message
		code := http.StatusInternalServerError
		message := "An error occurred while updating the item"

		// Check if the error is a Postgres error
		// If true we update the code and message accordingly
		if status := utils.PostgresToHTTPError(err); status != nil {
			code = *status

			if code == http.StatusNotFound {
				message = "Item not found"
			}
		}

		return &schemas.CustomError{
			Code:    code,
			Message: message,
			Details: fmt.Sprintf("Error updating item with ID %d: %v", id, err),
		}
	}

	return nil
}

func getPagedItems(page int, pageSize int) ([]schemas.Item, *int64, error) {
	client := db.Connect()

	_, count, err := client.
		From("items").
		Select("", "exact", false).
		Is("deleted_at", "null").
		Execute()

	if err != nil {
		return nil, nil, &schemas.CustomError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve items",
			Details: fmt.Sprintf("Failed to retrieve item count: %v", err),
		}
	}

	// If count is zero, return an empty slice to save time and resources;
	if count == 0 {
		return []schemas.Item{}, &count, nil
	}

	lastPage := (int(count) + pageSize - 1) / pageSize
	if page > lastPage {
		return nil, nil, &schemas.CustomError{
			Code:    http.StatusBadRequest,
			Message: "Page out of range",
			Details: fmt.Sprintf("Requested page %d with page size %d is out of range for total items %d", page, pageSize, count),
		}
	}

	// The end index must not exceed the total count of items
	pageEndIndex := min(page*pageSize, int(count))

	pageStartIndex := (page - 1) * pageSize

	slog.Info("Indexes", "start", pageStartIndex, "end", pageEndIndex)

	data, _, err := client.
		From("items").
		Select("*", "", false).Order("name", &postgrest.OrderOpts{Ascending: true}).
		Range(pageStartIndex, pageEndIndex-1, "").
		Execute()

	if err != nil {
		// Set the default error code and message
		code := http.StatusInternalServerError
		message := fmt.Sprintf("An error occurred while getting item page %v", page)

		// Check if the error is a Postgres error
		// If true we update the code and message accordingly
		if status := utils.PostgresToHTTPError(err); status != nil {
			code = *status

			if code == http.StatusNotFound {
				message = fmt.Sprintf("No items found for page %d", page)
			}
		}

		return nil, nil, &schemas.CustomError{
			Code:    code,
			Message: message,
			Details: fmt.Sprintf("Error retrieving items for page %d: %v", page, err),
		}
	}

	var items []schemas.Item
	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, nil, &schemas.CustomError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse items data",
			Details: fmt.Sprintf("Error parsing items data for page %d: %v", page, err),
		}
	}

	return items, &count, nil
}
