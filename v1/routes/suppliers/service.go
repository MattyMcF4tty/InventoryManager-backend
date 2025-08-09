package suppliers

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/MattyMcF4tty/InventoryManager-backend/v1/database"
	"github.com/MattyMcF4tty/InventoryManager-backend/v1/schemas"
	"github.com/MattyMcF4tty/InventoryManager-backend/v1/utils"
)

func GetSupplier(id int8) (schemas.Supplier, error) {
	client := db.Connect()
	idStr := fmt.Sprintf("%d", id)

	data, _, err := client.
		From("suppliers").
		Select("*", "", false).
		Eq("id", idStr).
		Is("deleted_at", "null").
		Single().
		Execute()

	if err != nil {
		// Set the default error code and message
		code := http.StatusInternalServerError
		message := "An error occurred while retrieving the supplier"

		// Check if the error is a Postgres error
		// If true we update the code and message accordingly
		if status := utils.PostgresToHTTPError(err); status != nil {
			code = *status

			if code == http.StatusNotFound {
				message = "Supplier not found"
			}
		}

		return schemas.Supplier{}, &schemas.CustomError{
			Code:    code,
			Message: message,
			Details: fmt.Sprintf("Error retrieving supplier with ID %d: %v", id, err),
		}
	}

	var supplier schemas.Supplier
	err = json.Unmarshal(data, &supplier)
	if err != nil {
		return schemas.Supplier{}, &schemas.CustomError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse supplier data",
			Details: fmt.Sprintf("Error parsing supplier data for ID %d: %v", id, err),
		}
	}

	// We make sure that the contact info is an empty array before returning it
	if supplier.ContactInfo == nil {
		supplier.ContactInfo = []*schemas.SupplierContactInfo{}
	}

	return supplier, nil
}
