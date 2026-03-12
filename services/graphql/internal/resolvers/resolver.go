package resolvers

import (
	"app/services/graphql/internal/repository"
)

type Resolver struct {
	Repo repository.TaskRepository
}

type QueryResolver struct{ *Resolver }
type MutationResolver struct{ *Resolver }
