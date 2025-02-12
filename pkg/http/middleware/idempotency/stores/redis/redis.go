package stores

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/FotiadisM/service-template/pkg/http/middleware/idempotency"
)

type RedisStore struct {
	client        *redis.Client
	log           *slog.Logger
	keysKey       string
	dataPrefixKey string
}

func NewRedisStore(client *redis.Client, log *slog.Logger, keysKey, dataPrefixKey string) *RedisStore {
	if log == nil {
		log = slog.Default()
	}

	return &RedisStore{
		client:        client,
		log:           log,
		keysKey:       keysKey,
		dataPrefixKey: dataPrefixKey,
	}
}

func (s *RedisStore) SetKey(ctx context.Context, key string) bool {
	items, err := s.client.SAdd(ctx, s.keysKey, key).Result()
	if err != nil {
		s.log.ErrorContext(ctx, "http-idempotency: redis-store: failed to set key", "error", err)
		return false
	}

	return items != 1
}

func (s *RedisStore) DelKey(ctx context.Context, key string) {
	err := s.client.SRem(ctx, s.keysKey, key).Err()
	if err != nil {
		s.log.ErrorContext(ctx, "http-idempotency: redis-store: failed to delete key", "error", err)
	}
}

func (s *RedisStore) GetData(ctx context.Context, key string) *idempotency.Data {
	res, err := s.client.Get(ctx, s.dataPrefixKey+key).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	if err != nil {
		s.log.ErrorContext(ctx, "http-idempotency: redis-store: failed to get data", "error", err)
		return nil
	}

	data := &idempotency.Data{}
	err = json.Unmarshal([]byte(res), data)
	if err != nil {
		s.log.ErrorContext(ctx, "http-idempotency: redis-store: failed to unmarshal data", "error", err)
	}

	return data
}

func (s *RedisStore) SetData(ctx context.Context, key string, data *idempotency.Data, exp time.Duration) {
	b, err := json.Marshal(data)
	if err != nil {
		s.log.ErrorContext(ctx, "http-idempotency: redis-store: failed to marshal data", "error", err)
		return
	}

	err = s.client.Set(ctx, s.dataPrefixKey+key, string(b), exp).Err()
	if err != nil {
		s.log.ErrorContext(ctx, "http-idempotency: redis-store: failed to set data", "error", err)
	}
}
