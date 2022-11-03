//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
//Licensed under the Apache License, Version 2.0 (the "License"); you may
//not use this file except in compliance with the License. You may obtain
//a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//License for the specific language governing permissions and limitations
// under the License.

//Package persistencemgr provides an  interfaces for database communication
package persistencemgr

import (
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

// MockDBConnection provides a mock db for unit testing
func MockDBConnection(t *testing.T) (*ConnPool, *errors.Error) {
	config.SetUpMockConfig(t)
	cfg, err := GetMockDBConfig()
	if err != nil {
		return nil, errors.PackError(errors.UndefinedErrorType, "error while trying to initiate mock db: ", err)
	}
	return cfg.Connection()
}

// MockDBWriteConnection provides a mock db write connection for unit testing
func MockDBWriteConnection(t *testing.T) (*Conn, *errors.Error) {
	connPool, err := MockDBConnection(t)
	if err != nil {
		return nil, errors.PackError(errors.UndefinedErrorType, "error while trying to initiate mock write db connection: ", err)
	}
	return connPool.GetWriteConnection()
}

// GetMockDBConfig will initiate mock db and will provide the config file
func GetMockDBConfig() (*Config, *errors.Error) {
	//Need to discuss more on this
	config.Data.DBConf = &config.DBConf{
		Protocol:              config.DefaultDBProtocol,
		OnDiskPort:            "6380",
		OnDiskHost:            "localhost",
		InMemoryHost:          "localhost",
		InMemoryPort:          "6379",
		RedisHAEnabled:        false,
		InMemorySentinelPort:  "26379",
		OnDiskSentinelPort:    "26379",
		InMemoryPrimarySet:    "redisSentinel",
		OnDiskPrimarySet:      "redisSentinel",
		MaxIdleConns:          config.DefaultDBMaxIdleConns,
		MaxActiveConns:        config.DefaultDBMaxActiveConns,
		RedisInMemoryPassword: []byte("redis_password"),
	}
	config := &Config{
		Port:     config.Data.DBConf.InMemoryPort,
		Protocol: config.Data.DBConf.Protocol,
		Host:     config.Data.DBConf.InMemoryHost,
		Password: string(config.Data.DBConf.RedisInMemoryPassword),
	}

	return config, nil
}
