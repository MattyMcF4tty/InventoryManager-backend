package suppliercontactinfo

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/MattyMcF4tty/InventoryManager-backend/v1/database"
	"github.com/MattyMcF4tty/InventoryManager-backend/v1/schemas"
	"github.com/MattyMcF4tty/InventoryManager-backend/v1/utils"
)

func GetSupplierContactInfo(supplierId int8) ([]schemas.SupplierContactInfo, error) {
	client := db.Connect()
	idStr := fmt.Sprintf("%d", supplierId)

	data, _, err := client.
		From("supplier_contact_information").
		Select("*", "", false).
		Eq("supplier_id", idStr).
		Execute()

	if err != nil {
		// Set the default error code and message
		code := http.StatusInternalServerError
		message := "An error occurred while retrieving the supplier contact info"

		// Check if the error is a Postgres error
		// If true we update the code and message accordingly
		if status := utils.PostgresToHTTPError(err); status != nil {
			code = *status

			if code == http.StatusNotFound {
				return []schemas.SupplierContactInfo{}, nil
			}
		}

		return []schemas.SupplierContactInfo{}, &schemas.CustomError{
			Code:    code,
			Message: message,
			Details: fmt.Sprintf("Error retrieving supplier contact info for supplier with ID %d: %v", supplierId, err),
		}
	}

	var supplierContactInfo []schemas.SupplierContactInfo
	err = json.Unmarshal(data, &supplierContactInfo)
	if err != nil {
		return []schemas.SupplierContactInfo{}, &schemas.CustomError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to parse supplier contact info data",
			Details: fmt.Sprintf("Error parsing supplier contact info data for supplier ID %d: %v", supplierId, err),
		}
	}

	return supplierContactInfo, nil
}
