package repository

import (
	"database/sql"
	"github.com/daniel-vuky/golang-todo-list-and-chat/model"
)

type ItemsRepository struct {
	Db *sql.DB
}

// Insert method of ItemsRepository
// @param item
// @throw error
func (itemsRepository ItemsRepository) Insert(item model.Item) (int64, error) {
	result, err := itemsRepository.Db.Exec(
		"INSERT INTO items (user_id, title, description, status) values (?, ?, ?, ?)",
		item.UserId,
		item.Title,
		item.Description,
		item.Status,
	)
	if err != nil {
		return 0, err
	}
	lastInsertId, insertErr := result.LastInsertId()
	return lastInsertId, insertErr
}

// Find method of ItemsRepository
// @param id
// @return item
// @throw error
func (itemsRepository ItemsRepository) Find(id int) (model.Item, error) {
	var item model.Item
	queryErr := itemsRepository.Db.QueryRow("SELECT * FROM items WHERE item_id = ?", id).Scan(
		&item.ItemId,
		&item.UserId,
		&item.Title,
		&item.Description,
		&item.Status,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if queryErr != nil {
		return model.Item{}, queryErr
	}

	return item, nil
}

// FindAll method of ItemsRepository
// @param limit
// @param offset
// @return list item
// @throw error
func (itemsRepository ItemsRepository) FindAll(limit, offset int) ([]model.Item, error) {
	items, err := itemsRepository.Db.Query(
		"SELECT * FROM items order by status, item_id LIMIT ? OFFSET ?",
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer items.Close()

	listItems := []model.Item{}
	for items.Next() {
		var item model.Item
		if scanErr := items.Scan(
			&item.ItemId,
			&item.UserId,
			&item.Title,
			&item.Description,
			&item.Status,
			&item.CreatedAt,
			&item.UpdatedAt,
		); scanErr != nil {
			return nil, scanErr
		}
		listItems = append(listItems, item)
	}

	return listItems, nil
}

// Update method of ItemsRepository
// @param item
// @param itemInput
// @throw error
func (itemsRepository ItemsRepository) Update(item model.Item, itemInput model.ItemInput) error {
	_, updatedError := itemsRepository.Db.Exec(
		"UPDATE items set title = ?, description = ?, status = ? WHERE item_id = ?",
		itemInput.Title,
		itemInput.Description,
		itemInput.Status,
		item.ItemId,
	)
	return updatedError
}

// Delete method of ItemsRepository
// @param itemId
// @throw error
func (itemsRepository ItemsRepository) Delete(itemId int) error {
	_, updatedError := itemsRepository.Db.Exec(
		"DELETE FROM items WHERE item_id = ?",
		itemId,
	)
	return updatedError
}
