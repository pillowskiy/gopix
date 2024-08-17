package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pkg/errors"
	redisClient "github.com/redis/go-redis/v9"
)

type userRedisCache struct {
	client *redisClient.Client
}

func NewUserCache(client *redisClient.Client) *userRedisCache {
	return &userRedisCache{client: client}
}

func (c *userRedisCache) GetByID(ctx context.Context, id int) (*domain.User, error) {
	userBytes, err := c.client.Get(ctx, c.stringifyKey(id)).Bytes()
	if err != nil {
		return nil, fmt.Errorf("userRedisCache.GetByID: %v", err)
	}

	user := new(domain.User)
	if err = json.Unmarshal(userBytes, user); err != nil {
		return nil, errors.Wrap(err, "userRedisCache.GetByID.UnmarshalJSON")
	}

	return user, nil
}

func (c *userRedisCache) SetUser(ctx context.Context, id int, user *domain.User, ttl int) error {
	userBytes, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "authRedisRepo.SetUserCtx.json.Unmarshal")
	}

	exp := time.Duration(ttl) * time.Second
	err = c.client.Set(ctx, c.stringifyKey(id), userBytes, exp).Err()
	if err != nil {
		return errors.Wrap(err, "authRedisRepo.SetUserCtx.redisClient.Set")
	}
	return nil
}

func (c *userRedisCache) DeleteUser(ctx context.Context, id int) error {
	if err := c.client.Del(ctx, c.stringifyKey(id)).Err(); err != nil {
		return errors.Wrap(err, "authRedisRepo.DeleteUserCtx.redisClient.Del")
	}
	return nil
}

func (c *userRedisCache) stringifyKey(id int) string {
	return strconv.Itoa(id)
}
