package resolvers

import (
	"context"
	"errors"

	"app/services/graphql/graph/generated"
	"app/services/graphql/internal/models"
)

func (r *MutationResolver) CreateTask(ctx context.Context, input generated.CreateTaskInput) (*models.Task, error) {
	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
		Done:        false,
	}
	if err := r.Repo.Create(task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *MutationResolver) UpdateTask(ctx context.Context, id string, input generated.UpdateTaskInput) (*models.Task, error) {
	existing, err := r.Repo.GetByID(id)
	if err != nil {
		if err.Error() == "task not found" {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	if input.Title != nil {
		existing.Title = *input.Title
	}
	if input.Description != nil {
		existing.Description = input.Description
	}
	if input.Done != nil {
		existing.Done = *input.Done
	}
	if input.DueDate != nil {
		existing.DueDate = input.DueDate
	}
	if err := r.Repo.Update(existing); err != nil {
		return nil, err
	}
	return &existing, nil
}

func (r *MutationResolver) DeleteTask(ctx context.Context, id string) (bool, error) {
	err := r.Repo.Delete(id)
	if err != nil {
		if err.Error() == "task not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
