package utils

import (
	"log"
	"sync"
	"time"

	"github.com/TheAlpha16/isolet/api/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/memory/v2"
	"github.com/gofiber/storage/redis/v3"
)

var (
	RedisStore  *redis.Storage
	MemoryStore *memory.Storage
	storeLock   sync.RWMutex
	IsRedisDown bool
)

func InitRedis() {
	MemoryStore = memory.New()

	defer func() {
		if r := recover(); r != nil {
			log.Println("[ERROR] Redis initialization panicked, using in-memory storage")
			IsRedisDown = true
		}
	}()

	store := redis.New(redis.Config{
		URL: config.REDIS_URL,
	})

	if err := store.Set("rate:test", []byte("ok"), 5*time.Second); err != nil {
		log.Println("[WARN] Redis is unreachable, using in-memory storage")
		IsRedisDown = true
	} else {
		RedisStore = store
		IsRedisDown = false
		log.Println("[INFO] Redis storage initialized successfully")
	}

	go PingRedis()
}

func PingRedis() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if RedisStore == nil {
			continue
		}

		err := RedisStore.Set("rate:health", []byte("ok"), 5*time.Second)

		storeLock.Lock()
		if err != nil {
			if !IsRedisDown {
				log.Println("[ERROR] Redis is down! Switching to in-memory storage")
				IsRedisDown = true
			}
		} else {
			if IsRedisDown {
				log.Println("[INFO] Redis is back online! Switching to Redis storage")
				IsRedisDown = false
			}
		}
		storeLock.Unlock()
	}
}

func GetActiveStore() fiber.Storage {
	storeLock.RLock()
	defer storeLock.RUnlock()
	if IsRedisDown {
		return MemoryStore
	}
	return RedisStore
}
