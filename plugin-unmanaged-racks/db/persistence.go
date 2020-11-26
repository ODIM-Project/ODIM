/*
 * Copyright (c) 2020 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

func (k Key) Id() string {
	return k.TrimWildcard().TrimPrefix().String()
}

func (k Key) TrimPrefix() Key {
	return k[strings.LastIndex(k.String(), ":")+1:]
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
