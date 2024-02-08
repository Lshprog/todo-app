package repository

import (
	"fmt"
	"strings"
	"todo-app/pkg/entities"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, item entities.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)
	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	if err := row.Scan(&itemId); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersItemQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createUsersItemQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]entities.TodoItem, error) {
	var items []entities.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON li.item_id = ti.id
		INNER JOIN %s ul ON ul.list_id = li.list_id WHERE ul.user_id = $1 AND li.list_id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)

	if err := r.db.Select(&items, query, userId, listId); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemPostgres) GetById(userId, itemId int) (entities.TodoItem, error) {
	var item entities.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON li.item_id = ti.id
		INNER JOIN %s ul ON ul.list_id = li.list_id WHERE ul.user_id = $1 AND ti.id = $2`,
		todoItemsTable, listsItemsTable, usersListsTable)
	err := r.db.Get(&item, query, userId, itemId)
	return item, err
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul
		WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id=$1 AND ti.id=$2`,
		todoItemsTable, listsItemsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, itemId)

	return err
}

func (r *TodoItemPostgres) Update(userId, itemId int, input entities.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul
		WHERE ti.id = li.item_id AND ul.list_id = li.list_id AND ti.id=$%d AND ul.user_id=$%d`,
		todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)

	args = append(args, itemId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)

	return err
}
