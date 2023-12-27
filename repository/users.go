package repository

import (
	"database/sql"
	"github.com/daniel-vuky/golang-todo-list-v2/model"
)

type UsersRepository struct {
	Db *sql.DB
}

// CreateNewUser register an new user
func (usersRepository UsersRepository) CreateNewUser(user *model.User) error {
	result, err := usersRepository.Db.Exec(
		"INSERT INTO users (username, email, password) values (?, ?, ?)",
		user.Username,
		user.Email,
		user.Password,
	)
	if err != nil {
		return err
	}
	lastInsertId, insertErr := result.LastInsertId()
	if insertErr != nil {
		return insertErr
	}
	user.UserId = uint64(lastInsertId)
	return nil
}

// GetUser get existed user
func (usersRepository UsersRepository) GetUser(user *model.User) error {
	exec := "SELECT user_id, username, email, password FROM users WHERE username = ?"
	queryErr := usersRepository.Db.QueryRow(exec, user.Username).Scan(
		&user.UserId,
		&user.Username,
		&user.Email,
		&user.Password,
	)

	return queryErr
}
