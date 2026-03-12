package repository

import "app/services/graphql/internal/models"

type TaskRepository interface {
	Create(task models.Task) error
	GetAll() ([]models.Task, error)
	GetByID(id string) (models.Task, error)
	Update(task models.Task) error
	Delete(id string) error
}
