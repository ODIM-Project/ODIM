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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/config"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/logging"
	"github.com/ODIM-Project/ODIM/plugin-unmanaged-racks/redfish"

	"github.com/go-redis/redis/v8"
)

type TLSConfig func(*config.PluginConfig) (*tls.Config, error)

var GetTLSConfig TLSConfig = func(c *config.PluginConfig) (*tls.Config, error) {
	caCert, err := ioutil.ReadFile(c.PKIRootCAPath)
	if err != nil {
		return &tls.Config{}, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCert)
	cert, err := tls.LoadX509KeyPair(c.PKICertificatePath, c.PKIPrivateKeyPath)
	if err != nil {
		return &tls.Config{}, err
	}
	cfg := &tls.Config{
		RootCAs:      pool,
		MinVersion:   c.TLSConf.MinVersion,
		Certificates: []tls.Certificate{cert},
	}
	return cfg, nil
}

// CreateDAO creates new instance of DAO
func CreateDAO(c *config.PluginConfig, sentinelMasterName string, getTLSConfig TLSConfig) *DAO {
	tlsConfig, err := getTLSConfig(c)
	if err != nil {
		logging.Fatalf("error while getting tls configuration: %s", err.Error())
	}
	if sentinelMasterName == "" {
		return &DAO{
			redis.NewClient(&redis.Options{
				Addr:      c.RedisAddress,
				TLSConfig: tlsConfig,
				Password:  string(c.RedisOnDiskPassword),
			}),
		}
	}

	return &DAO{
		redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       sentinelMasterName,
			SentinelAddrs:    []string{c.RedisAddress},
			TLSConfig:        tlsConfig,
			SentinelPassword: string(c.RedisOnDiskPassword),
			Password:         string(c.RedisOnDiskPassword),
		}),
	}
}

// DAO struct provides access to Redis
type DAO struct {
	*redis.Client
}

// FindChassis finds requested chassis or returns nil if chassis doesn't exists
func (c *DAO) FindChassis(chassisOid string) (*redfish.Chassis, error) {
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

// Key is an string alias representing Redis key. Implementation of Key wrapper comes with number of utility functions.
// Key wraps number of tokens separated with a ":" separator, for example: "foo", "foo:bar", "foo:bar:foobar", etc.
// Last token is called ID, all other tokens taken together are called PREFIX.
type Key string

// String returns unwrapped key
func (k Key) String() string {
	return string(k)
}

// WithWildcard returns key which contains value of current key(k) concatenated with wildcard("*")
func (k Key) WithWildcard() Key {
	return k + "*"
}

// ID returns unwrapped key with trimmed ending wildcard
func (k Key) ID() string {
	return k.trimWildcard().trimPrefix().String()
}

func (k Key) trimPrefix() Key {
	return k[strings.LastIndex(k.String(), ":")+1:]
}

func (k Key) trimWildcard() Key {
	return Key(strings.TrimSuffix(string(k), "*"))
}

// CreateContainsKey utility function which produces key for CONTAINS relation, for example: "CONTAINS:CHASSIS"
func CreateContainsKey(tokens ...string) Key {
	return CreateKey(append([]string{"CONTAINS"}, tokens...)...)
}

// CreateContainedInKey utility function which produces key for CONTAINEDIN relation, for example: "CONTAINEDIN:CHASSIS"
func CreateContainedInKey(tokens ...string) Key {
	return CreateKey(append([]string{"CONTAINEDIN"}, tokens...)...)
}

// CreateKey creates new instance of key
func CreateKey(tokens ...string) Key {
	return Key(strings.Join(tokens, ":"))
}
