package items

import (
	"encoding/json"
	"fmt"

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
		return schemas.Item{}, err
	}

	var item schemas.Item
	err = json.Unmarshal(data, &item)
	if err != nil {
		return schemas.Item{}, err
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
		return schemas.Item{}, err
	}

	var item schemas.Item
	err = json.Unmarshal(data, &item)

	if err != nil {
		return schemas.Item{}, err
	}

	return item, nil
}

func CreateItem(item schemas.Item) (schemas.Item, error) {
	client := db.Connect()

	item.CreatedAt = utils.GetCurrentISODate()
	item.UpdatedAt = utils.GetCurrentISODate()

	data, _, err := client.From("items").Insert(item, false, "", "", "").Single().Execute()

	if err != nil {
		return schemas.Item{}, err
	}

	var createdItem schemas.Item
	err = json.Unmarshal(data, &createdItem)

	if err != nil {
		return schemas.Item{}, err
	}

	return createdItem, nil
}

func DeleteItem(id int8) error {
	client := db.Connect()
	idStr := fmt.Sprintf("%d", id)

	_, _, err := client.From("items").Delete("", "").Eq("id", idStr).Execute()

	if err != nil {
		return err
	}

	return nil
}
