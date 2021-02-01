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
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

// Config is the configuration for db which is set by the wrapper package.
/*
Port is the port number for the database connection
Protocol is the type of protocol with which the connection takes place
Host is hostname/IP on which the database is running
*/
type Config struct {
	Port         string
	Protocol     string
	Host         string
	SentinelPort string
	MasterSet    string
}

// ConnPool is the established connection
type ConnPool struct {
	ReadPool        *redis.Pool
	WritePool       *redis.Pool
	MasterIP        string
	PoolUpdatedTime time.Time
	Mux             sync.Mutex
}
