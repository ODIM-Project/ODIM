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
	stdCtx "context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"
	"github.com/go-redis/redis/v8"
)

var ErrAlreadyExists = errors.New("already exists")

func NewConnectionManager(redisAddress, sentinelMasterName string) *ConnectionManager {
	if sentinelMasterName == "" {
		return &ConnectionManager{
			redis.NewClient(&redis.Options{
				Addr: redisAddress,
			}),
		}
	} else {
		return &ConnectionManager{
			redis.NewFailoverClient(&redis.FailoverOptions{
				MasterName:    sentinelMasterName,
				SentinelAddrs: []string{redisAddress},
			}),
		}
	}
}

type ConnectionManager struct {
	*redis.Client
}

func (c *ConnectionManager) FindChassis(chassisOid string) (*redfish.Chassis, error) {
	v, err := c.Client.Get(stdCtx.TODO(), CreateKey("Chassis", chassisOid).String()).Result()
	if err != nil && err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	chassis := new(redfish.Chassis)
	err = json.Unmarshal([]byte(v), chassis)
	if err != nil {
		return nil, err
	}

	chassisContainsKey := CreateContainsKey("Chassis", chassisOid)
	members, err := c.Client.SMembers(stdCtx.TODO(), chassisContainsKey.String()).Result()
	if err != nil {
		return nil, err
	}

	for _, m := range members {
		chassis.Links.Contains = append(chassis.Links.Contains, redfish.Link{Oid: m})
	}

	return chassis, nil
}

func (c *ConnectionManager) DAO() *redis.Client {
	return c.Client
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
