package redis

import (
	"backend_golang_codeing_test/config"
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Manager struct {
	read  *redis.Client
	write *redis.Client
	pub   *redis.Client
}

var (
	redisInstance *Manager
	onceRedis     sync.Once
)

func NewRedisManager(conf *config.Redis) *Manager {
	onceRedis.Do(func() {
		redisInstance = &Manager{
			read:  newRedisClient(&conf.Read, "read"),
			write: newRedisClient(&conf.Write, "write"),
			pub:   newRedisClient(&conf.Pub, "pub"),
		}
	})

	return redisInstance
}

func newRedisClient(conf *config.RedisInstance, label string) *redis.Client {
	addr := fmt.Sprintf("%s:%s", conf.Host, conf.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Password,
		DB:       conf.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("❌ Redis %s connection failed (%s): %v", label, addr, err)
	}

	log.Printf("✅ Redis %s connected at %s", label, addr)
	return rdb
}

func (m *Manager) GetRead() *redis.Client  { return m.read }
func (m *Manager) GetWrite() *redis.Client { return m.write }
func (m *Manager) GetPub() *redis.Client   { return m.pub }
