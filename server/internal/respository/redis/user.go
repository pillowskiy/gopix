package redis

import (
	"github.com/pillowskiy/gopix/internal/domain"
	redisClient "github.com/redis/go-redis/v9"
)

func NewUserCache(client *redisClient.Client) *ItoaCache[domain.User] {
	return NewItoaCache[domain.User](client)
}
