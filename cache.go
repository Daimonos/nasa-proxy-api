package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis"
)

type Cache struct {
	Client *redis.Client
}

func (c *Cache) InitializeClient(addr string, poolsize int, pooltimeout int) {
	if poolsize == 0 {
		poolsize = 10
	}
	if pooltimeout == 0 {
		pooltimeout = 30
	}
	c.Client = redis.NewClient(&redis.Options{
		Addr:        addr,
		PoolSize:    poolsize,
		PoolTimeout: time.Duration(pooltimeout) * time.Second,
	})
}

func (c *Cache) Set(key string, value interface{}, expiration int) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		return err
	}
	return c.Client.Set(key, string(bytes), time.Duration(expiration)*time.Second).Err()
}

func (c *Cache) Get(key string) ([]byte, error) {
	s, err := c.Client.Get(key).Result()
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
	}
	return []byte(s), nil
}
