package redis

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/pkg/errors"
	redisClient "github.com/redis/go-redis/v9"
)

type ItoaCache[T interface{}] struct {
	client *redisClient.Client
	name   string
}

func NewItoaCache[T interface{}](client *redisClient.Client) *ItoaCache[T] {
	return &ItoaCache[T]{client: client, name: reflect.TypeOf(new(T)).Name()}
}

func (c *ItoaCache[T]) Get(ctx context.Context, id int) (*T, error) {
	bytes, err := c.client.Get(ctx, c.stringifyKey(id)).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "ItoaCache.GetByID.redisClient.Get")
	}

	data := new(T)
	if err = json.Unmarshal(bytes, data); err != nil {
		return nil, errors.Wrap(err, "ItoaCache.GetByID.UnmarshalJSON")
	}

	return data, nil
}

func (c *ItoaCache[T]) Set(ctx context.Context, id int, data *T, ttl int) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "ItoaCache.SetUser.json.Unmarshal")
	}

	exp := time.Duration(ttl) * time.Second
	err = c.client.Set(ctx, c.stringifyKey(id), bytes, exp).Err()
	if err != nil {
		return errors.Wrap(err, "ItoaCache.SetUser.redisClient.Set")
	}
	return nil
}

func (c *ItoaCache[T]) Del(ctx context.Context, id int) error {
	if err := c.client.Del(ctx, c.stringifyKey(id)).Err(); err != nil {
		return errors.Wrap(err, "ItoaCache.DeleteUser.redisClient.Del")
	}
	return nil
}

func (c *ItoaCache[T]) stringifyKey(id int) string {
	return c.name + "_" + strconv.Itoa(id)
}
