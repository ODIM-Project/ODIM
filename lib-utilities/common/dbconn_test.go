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
package common

import (
	"fmt"
	"math"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
)

const (
	undefinedErr int = iota
	getRedisConnErr
	emptyRedisConnObjErr
	getRedisSearchConnErr
	emptyRedisSearchConnObjErr
	redisDataWriteErr
	redisDataReadErr
	redisEmptyDBErr
	redisDBFlushErr
	redisFailedFlushErr
	reuseConnPoolErr
)

func reportError(t *testing.T, errType int, errMsg ...interface{}) {
	switch errType {
	case getRedisConnErr:
		t.Errorf("expected error to be nil while using GetDBConnection but got: %v", errMsg)
	case emptyRedisConnObjErr:
		t.Error("expected conn not to be nil while using GetDBConnection but got nil")
	case getRedisSearchConnErr:
		t.Errorf("expected error to be nil while using GetRediSearchConnection but got: %v", errMsg)
	case emptyRedisSearchConnObjErr:
		t.Error("expected conn not to be nil while using GetRediSearchConnection but got nil")
	case redisDataWriteErr:
		t.Errorf("Error while creating data: %v", errMsg)
	case redisDataReadErr:
		t.Errorf("error while fetching data: %v", errMsg)
	case redisEmptyDBErr:
		t.Error("no data found before cleaning the db")
	case redisDBFlushErr:
		t.Errorf("error while trying to flush db: %v", errMsg)
	case redisFailedFlushErr:
		t.Error("database was not fully cleaned")
	case reuseConnPoolErr:
		t.Errorf("expected second connection pool, to be equal to first connection pool %v", errMsg)
	default:
		t.Error(errMsg...)
	}
}

func TestGetDBConnectionInMemory(t *testing.T) {
	config.SetUpMockConfig(t)
	// Get In-Memory db connection pool
	conn, err := GetDBConnection(InMemory)
	if err != nil {
		reportError(t, getRedisConnErr, err)
	}
	if conn == nil {
		reportError(t, emptyRedisConnObjErr)
	}
}
func TestGetDBConnectionOnDisk(t *testing.T) {
	config.SetUpMockConfig(t)
	// Get In-Disk db connection pool
	conn, err := GetDBConnection(OnDisk)
	if err != nil {
		reportError(t, getRedisConnErr, err)
	}
	if conn == nil {
		reportError(t, emptyRedisConnObjErr)
	}
}
func TestGetDBConnectionDefaultCase(t *testing.T) {
	config.SetUpMockConfig(t)
	// Get In-Disk db connection pool
	conn, err := GetDBConnection(math.MaxInt32)
	if err == nil {
		reportError(t, getRedisConnErr, err)
	}
	if conn != nil {
		reportError(t, emptyRedisConnObjErr)
	}
}

func TestGetDBConnectionForExistingConnInMemory(t *testing.T) {
	config.SetUpMockConfig(t)
	// Get In-Memory db connection pool
	conn, err := GetDBConnection(InMemory)
	if err != nil {
		reportError(t, getRedisConnErr, err)
	}
	if conn == nil {
		reportError(t, emptyRedisConnObjErr)
	}

	// now in-memory connection pool already exists, expects to be served by same connection pool
	sConn, err := GetDBConnection(InMemory)
	if err != nil {
		reportError(t, getRedisConnErr, err)
	}
	if sConn != conn {
		reportError(t, reuseConnPoolErr, "got:", sConn, "want:", conn)
	}
}
func TestGetDBConnectionForExistingConnOnDisk(t *testing.T) {
	config.SetUpMockConfig(t)
	// Get In-Disk db connection pool
	conn, err := GetDBConnection(OnDisk)
	if err != nil {
		reportError(t, getRedisConnErr, err)
	}
	if conn == nil {
		reportError(t, emptyRedisConnObjErr)
	}

	// now In-Disk connection pool already exists, expects to be served by same connection pool
	sConn, err := GetDBConnection(OnDisk)
	if err != nil {
		reportError(t, getRedisConnErr, err)
	}
	if sConn != conn {
		reportError(t, reuseConnPoolErr, "got:", sConn, "want:", conn)
	}
}

func TestTruncateDBInMemory(t *testing.T) {
	config.SetUpMockConfig(t)
	// Get In-Memory db connection pool
	conn, err := GetDBConnection(InMemory)
	if err != nil {
		reportError(t, getRedisConnErr, err)
	}
	if conn == nil {
		reportError(t, emptyRedisConnObjErr)
	}

	if errs := conn.Create("table", "key1", "data1"); errs != nil {
		reportError(t, redisDataWriteErr, errs)
	}
	keys, errs := conn.GetAllDetails("*")
	if errs != nil {
		reportError(t, redisDataReadErr, errs)
	}
	if len(keys) == 0 {
		reportError(t, redisEmptyDBErr)
	}
	err = TruncateDB(InMemory)
	if err != nil {
		reportError(t, redisDBFlushErr, err)
	}
	keys, errs = conn.GetAllDetails("*")
	if errs != nil {
		reportError(t, redisDataReadErr, errs)
	}
	if len(keys) != 0 {
		reportError(t, redisFailedFlushErr)
	}
}

func TestTruncateDBOnDisk(t *testing.T) {
	config.SetUpMockConfig(t)
	// Get In-Disk db connection pool
	conn, err := GetDBConnection(OnDisk)
	if err != nil {
		reportError(t, getRedisConnErr, err)
	}
	if conn == nil {
		reportError(t, emptyRedisConnObjErr)
	}

	if errs := conn.Create("table", "key1", "data1"); errs != nil {
		reportError(t, redisDataWriteErr, errs)
	}
	keys, errs := conn.GetAllDetails("*")
	if errs != nil {
		reportError(t, redisDataReadErr, errs)
	}
	if len(keys) == 0 {
		reportError(t, redisEmptyDBErr)
	}
	err = TruncateDB(OnDisk)
	if err != nil {
		reportError(t, redisDBFlushErr, err)
	}
	keys, errs = conn.GetAllDetails("*")
	if errs != nil {
		reportError(t, redisDataReadErr, errs)
	}
	if len(keys) != 0 {
		reportError(t, redisFailedFlushErr)
	}
}
func TestTruncateDBNotExist(t *testing.T) {
	config.SetUpMockConfig(t)
	err := TruncateDB(math.MaxInt32)
	if err == nil {
		reportError(t, undefinedErr, "expected err to be non nil while trying to flush db:", err)
	}
}

func TestCheckDBConnection(t *testing.T) {
	config.SetUpMockConfig(t)
	if err := CheckDBConnection(); err != nil {
		reportError(t, undefinedErr, fmt.Sprintf("expected err to be nil while using CheckDBConnection but got: %v", err))
	}
}
