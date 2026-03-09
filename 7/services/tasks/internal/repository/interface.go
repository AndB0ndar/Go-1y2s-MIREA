package repository

import "app/services/tasks/internal/service"

type TaskRepository interface {
	Create(task service.Task) (service.Task, error)
	GetAll() ([]service.Task, error)
	GetByID(id string) (service.Task, error)
	Update(task service.Task) error
	Delete(id string) error
	SearchByTitle(title string) ([]service.Task, error)
	SearchByTitleVulnerable(title string) ([]service.Task, error)
}
