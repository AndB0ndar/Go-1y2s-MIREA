package resolvers

import (
	"context"
	"time"

	"app/services/graphql/internal/models"
)

func (r *TaskResolver) ID(ctx context.Context, obj *models.Task) (string, error) {
	return obj.ID, nil
}

func (r *TaskResolver) Title(ctx context.Context, obj *models.Task) (string, error) {
	return obj.Title, nil
}

func (r *TaskResolver) Description(ctx context.Context, obj *models.Task) (*string, error) {
	return obj.Description, nil
}

func (r *TaskResolver) Done(ctx context.Context, obj *models.Task) (bool, error) {
	return obj.Done, nil
}

func (r *TaskResolver) DueDate(ctx context.Context, obj *models.Task) (*string, error) {
	return obj.DueDate, nil
}

func (r *TaskResolver) CreatedAt(ctx context.Context, obj *models.Task) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

func (r *TaskResolver) UpdatedAt(ctx context.Context, obj *models.Task) (string, error) {
	return obj.UpdatedAt.Format(time.RFC3339), nil
}
