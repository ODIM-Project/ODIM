package db

import (
	"errors"
	"strings"

	"github.com/go-redis/redis/v8"
)

var ErrAlreadyExists = errors.New("already exists")

func NewConnectionManager(redisAddress, sentinelMasterName string) *ConnectionManager {
	if sentinelMasterName == "" {
		return &ConnectionManager{
			c: redis.NewClient(&redis.Options{
				Addr: redisAddress,
			}),
		}
	} else {
		return &ConnectionManager{
			c: redis.NewFailoverClient(&redis.FailoverOptions{
				MasterName:    sentinelMasterName,
				SentinelAddrs: []string{redisAddress},
			}),
		}
	}
}

type ConnectionManager struct {
	c *redis.Client
}

func (c *ConnectionManager) DAO() *redis.Client {
	return c.c
}

type Key string

func (k Key) String() string {
	return string(k)
}

func (k Key) Prefix() string {
	return string(k) + ":"
}

func (k Key) TrimWildcard() Key {
	return Key(strings.TrimSuffix(string(k), "*"))
}

func (k Key) WithWildcard() Key {
	return k + "*"
}

func CreateContainsKey(tokens ...string) Key {
	return CreateKey(append([]string{"CONTAINS"}, tokens...)...)
}

func CreateContainedInKey(tokens ...string) Key {
	return CreateKey(append([]string{"CONTAINEDIN"}, tokens...)...)
}

func CreateKey(keys ...string) Key {
	return Key(strings.Join(keys, ":"))
}
