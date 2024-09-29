package redis

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/pkg/errors"
	redisClient "github.com/redis/go-redis/v9"
)

type Cache[T interface{}] struct {
	client *redisClient.Client
	name   string
}

func NewCache[T interface{}](client *redisClient.Client) *Cache[T] {
	return &Cache[T]{client: client, name: reflect.TypeOf(new(T)).Name()}
}

func (c *Cache[T]) Get(ctx context.Context, id string) (*T, error) {
	bytes, err := c.client.Get(ctx, id).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "ItoaCache.GetByID.redisClient.Get")
	}

	data := new(T)
	if err = json.Unmarshal(bytes, data); err != nil {
		return nil, errors.Wrap(err, "ItoaCache.GetByID.UnmarshalJSON")
	}

	return data, nil
}

func (c *Cache[T]) Set(ctx context.Context, id string, data *T, ttl int) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "ItoaCache.SetUser.json.Unmarshal")
	}

	exp := time.Duration(ttl) * time.Second
	err = c.client.Set(ctx, id, bytes, exp).Err()
	if err != nil {
		return errors.Wrap(err, "ItoaCache.SetUser.redisClient.Set")
	}
	return nil
}

func (c *Cache[T]) Del(ctx context.Context, id string) error {
	if err := c.client.Del(ctx, id).Err(); err != nil {
		return errors.Wrap(err, "ItoaCache.DeleteUser.redisClient.Del")
	}
	return nil
}
