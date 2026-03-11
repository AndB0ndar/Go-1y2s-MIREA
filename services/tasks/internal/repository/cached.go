package repository

import (
    "context"
    "encoding/json"
    "fmt"

    "app/services/tasks/internal/cache"
    "app/services/tasks/internal/service"
    "github.com/sirupsen/logrus"
)

type cachedTaskRepository struct {
    repo      TaskRepository
    cache     *cache.RedisClient
    log       *logrus.Logger
    baseTTL   int
    jitterTTL int
}

func NewCachedTaskRepository(
	repo TaskRepository,
	cache *cache.RedisClient,
	log *logrus.Logger,
	baseTTL, jitterTTL int,
) TaskRepository {
    return &cachedTaskRepository{
        repo:      repo,
        cache:     cache,
        log:       log,
        baseTTL:   baseTTL,
        jitterTTL: jitterTTL,
    }
}

func (c *cachedTaskRepository) taskKey(id string) string {
    return fmt.Sprintf("tasks:task:%s", id)
}

func (c *cachedTaskRepository) Create(
	task service.Task,
) (service.Task, error) {
	createdTask, err := c.repo.Create(task)
    if err != nil {
        return service.Task{}, err
    }
    return createdTask, err
}

func (c *cachedTaskRepository) GetAll() ([]service.Task, error) {
    return c.repo.GetAll()
}

func (c *cachedTaskRepository) GetByID(id string) (service.Task, error) {
    ctx := context.Background()
    key := c.taskKey(id)

    data, hit, err := c.cache.GetTask(ctx, key)
    if err == nil && hit {
        var task service.Task
        if err := json.Unmarshal(data, &task); err == nil {
            c.log.WithField("key", key).Info("Cache hit")
            return task, nil
        } else {
            c.log.WithError(err).Warn("Failed to unmarshal cached task")
        }
    }

    c.log.WithField("key", key).Info("Cache miss")
    task, err := c.repo.GetByID(id)
    if err != nil {
        return task, err
    }

    go func() {
        ctxBg := context.Background()
        data, _ := json.Marshal(task)
        c.cache.SetTask(ctxBg, key, data, c.baseTTL, c.jitterTTL)
    }()

    return task, nil
}

func (c *cachedTaskRepository) Update(task service.Task) error {
    if err := c.repo.Update(task); err != nil {
        return err
    }
    go func() {
        ctxBg := context.Background()
        c.cache.Delete(ctxBg, c.taskKey(task.ID))
    }()
    return nil
}

func (c *cachedTaskRepository) Delete(id string) error {
    if err := c.repo.Delete(id); err != nil {
        return err
    }
    go func() {
        ctxBg := context.Background()
        c.cache.Delete(ctxBg, c.taskKey(id))
    }()
    return nil
}

func (c *cachedTaskRepository) SearchByTitle(title string) ([]service.Task, error) {
    return c.repo.SearchByTitle(title)
}

func (c *cachedTaskRepository) SearchByTitleVulnerable(title string) ([]service.Task, error) {
    return c.repo.SearchByTitleVulnerable(title)
}
