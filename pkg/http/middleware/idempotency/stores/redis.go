package stores

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/FotiadisM/mock-microservice/pkg/http/middleware/idempotency"
)

type RedisStore struct {
	client        *redis.Client
	keysKey       string
	dataPrefixKey string
}

func NewRedisStore(client *redis.Client, keysKey, dataPrefixKey string) *RedisStore {
	return &RedisStore{
		client:        client,
		keysKey:       keysKey,
		dataPrefixKey: dataPrefixKey,
	}
}

func (s *RedisStore) SetKey(ctx context.Context, key string) (bool, error) {
	items, err := s.client.SAdd(ctx, s.keysKey, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set key %s: %w", key, err)
	}

	return items != 1, nil
}

func (s *RedisStore) DelKey(ctx context.Context, key string) error {
	err := s.client.SRem(ctx, s.keysKey, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}

	return nil
}

func (s *RedisStore) GetData(ctx context.Context, key string) (*idempotency.Data, error) {
	res, err := s.client.Get(ctx, s.dataPrefixKey+key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, idempotency.ErrNoDataFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get data for key %s: %w", key, err)
	}

	data := &idempotency.Data{}
	err = json.Unmarshal([]byte(res), data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return data, nil
}

func (s *RedisStore) SetData(ctx context.Context, key string, data *idempotency.Data, exp time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	err = s.client.Set(ctx, s.dataPrefixKey+key, string(b), exp).Err()
	if err != nil {
		return fmt.Errorf("failed to set data for key %s: %w", key, err)
	}

	return nil
}
