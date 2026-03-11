package cache

import (
    "context"
    "math/rand"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/sirupsen/logrus"
)

type RedisClient struct {
    client *redis.Client
    log    *logrus.Logger
}

func NewRedisClient(addr, password string, db int, log *logrus.Logger) (*RedisClient, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, err
    }
    return &RedisClient{client: client, log: log}, nil
}

func (r *RedisClient) Close() error {
    return r.client.Close()
}

func (r *RedisClient) GetTask(ctx context.Context, key string) ([]byte, bool, error) {
    data, err := r.client.Get(ctx, key).Bytes()
    if err == redis.Nil {
        return nil, false, nil // miss
    }
    if err != nil {
        r.log.WithError(err).Warn("Redis get error")
        return nil, false, err // ошибка, считаем как miss и идём в БД
    }
    return data, true, nil // hit
}

func (r *RedisClient) SetTask(ctx context.Context, key string, data []byte, baseTTL, jitterSeconds int) {
    ttl := time.Duration(baseTTL) * time.Second
    if jitterSeconds > 0 {
        jitter := time.Duration(rand.Intn(jitterSeconds)) * time.Second
        ttl += jitter
    }
    err := r.client.Set(ctx, key, data, ttl).Err()
    if err != nil {
        r.log.WithError(err).Warn("Redis set error")
    }
}

func (r *RedisClient) Delete(ctx context.Context, key string) {
    err := r.client.Del(ctx, key).Err()
    if err != nil {
        r.log.WithError(err).Warn("Redis delete error")
    }
}
