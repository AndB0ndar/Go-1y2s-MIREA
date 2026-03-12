package resolvers

import (
	"context"

	"app/services/graphql/internal/models"
)

func (r *QueryResolver) Tasks(ctx context.Context) ([]*models.Task, error) {
	tasks, err := r.Repo.GetAll()
	if err != nil {
		return nil, err
	}
	// Convert to slice of pointers
	result := make([]*models.Task, len(tasks))
	for i := range tasks {
		result[i] = &tasks[i]
	}
	return result, nil
}

func (r *QueryResolver) Task(ctx context.Context, id string) (*models.Task, error) {
	task, err := r.Repo.GetByID(id)
	if err != nil {
		if err.Error() == "task not found" {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}
