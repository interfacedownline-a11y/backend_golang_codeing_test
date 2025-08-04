package redis

import "github.com/redis/go-redis/v9"

type Redis interface {
	GetRead() *redis.Client
	GetWrite() *redis.Client
	GetPub() *redis.Client
}
