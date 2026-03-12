package resolvers

import (
	"app/services/graphql/graph/generated"
	"app/services/graphql/internal/repository"
)

type Resolver struct {
	Repo repository.TaskRepository
}

func (r *Resolver) Task() generated.TaskResolver {
	return &TaskResolver{r}
}

func (r *Resolver) Query() generated.QueryResolver {
	return &QueryResolver{r}
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &MutationResolver{r}
}

type QueryResolver struct{ *Resolver }
type MutationResolver struct{ *Resolver }
type TaskResolver struct{ *Resolver }
