package repository

import (
	"todo-app/pkg/entities"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user entities.User) (int, error)
	GetUser(username, password string) (entities.User, error)
}

type TodoList interface {
	Create(userId int, list entities.TodoList) (int, error)
	GetAll(userId int) ([]entities.TodoList, error)
	GetById(userId, listId int) (entities.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input entities.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, item entities.TodoItem) (int, error)
	GetAll(userId, listId int) ([]entities.TodoItem, error)
	GetById(userId, itemId int) (entities.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input entities.UpdateItemInput) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
