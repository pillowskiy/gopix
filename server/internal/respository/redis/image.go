package redis

import (
	"github.com/pillowskiy/gopix/internal/domain"
	redisClient "github.com/redis/go-redis/v9"
)

func NewImageCache(client *redisClient.Client) *ItoaCache[domain.Image] {
	return NewItoaCache[domain.Image](client)
}
