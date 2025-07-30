package items

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/MattyMcF4tty/InventoryManager-backend/v1/database"
	"github.com/MattyMcF4tty/InventoryManager-backend/v1/schemas"
	"github.com/MattyMcF4tty/InventoryManager-backend/v1/utils"
)

func GetItem(id int8) (schemas.Item, error) {
	client := db.Connect()
	idStr := fmt.Sprintf("%d", id)

	data, _, err := client.
		From("items").
		Select("*", "", false).
		Eq("id", idStr).
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

	data, _, err := client.From("items").Update(updates, "", "").Eq("id", idStr).Single().Execute()

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

	data, _, err := client.From("items").Insert(item, false, "", "", "").Single().Execute()

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

	_, _, err := client.From("items").Delete("", "").Eq("id", idStr).Execute()

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
