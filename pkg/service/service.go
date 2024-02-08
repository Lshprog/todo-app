package service

import (
	"todo-app/pkg/entities"
	"todo-app/pkg/repository"
)

type Authorization interface {
	CreateUser(user entities.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, list entities.TodoList) (int, error)
	GetAll(userId int) ([]entities.TodoList, error)
	GetById(userId, listId int) (entities.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input entities.UpdateListInput) error
}

type TodoItem interface {
	Create(userId, listId int, item entities.TodoItem) (int, error)
	GetAll(userId, listId int) ([]entities.TodoItem, error)
	GetById(userId, itemId int) (entities.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, input entities.UpdateItemInput) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
