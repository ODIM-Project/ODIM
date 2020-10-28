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

// Package common ...
package common

import (
	"fmt"
	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

// DbType is a alias name for int32
type DbType int32

const (
	//InMemory - To select in-memory db connection pool
	InMemory DbType = iota
	// OnDisk - To select in-disk db connection pool
	OnDisk
)

// GetDBConnection is for maintaining and supplying DB connection pool for InMemory and OnDisk DB's
// Takes dbFlag of type DbType/int32
// dbFlag:
//	InMemory:	returns In-Memory DB connection pool
//	OnDsik:  	returns On-Disk DB connection pool
func GetDBConnection(dbFlag DbType) (*persistencemgr.ConnPool, *errors.Error) {
	switch dbFlag {
	case InMemory:
		pool, err := persistencemgr.GetDBConnection(persistencemgr.InMemory)
		return pool, err
	case OnDisk:
		pool, err := persistencemgr.GetDBConnection(persistencemgr.OnDisk)
		return pool, err
	default:
		return nil, errors.PackError(errors.UndefinedErrorType, "error invalid db type selection")
	}
}

// TruncateDB will clear DB. It will be useful for test cases
// Takes DbFlag of type DbType/int32 to choose Inmemory or OnDisk db to truncate
//dbFlag:
//    InMemory: Truncates InMemory DB
//    OnDisk: Truncates OnDisk DB
func TruncateDB(dbFlag DbType) *errors.Error {
	conn, err := GetDBConnection(dbFlag)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to connect to DB: ", err.Error())
	}
	err = conn.CleanUpDB()
	if err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to flush out DB: ", err.Error())
	}
	return nil
}

// CheckDBConnection will check both inMemory and onDisk DB connections
// This function is expected to be called at each service startup
func CheckDBConnection() error {
	inMemConn, err := GetDBConnection(InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to create InMemory DB connection: %v", err)
	}
	onDiskConn, err := GetDBConnection(OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to create OnDisk DB connection: %v", err)
	}

	if err := inMemConn.Ping(); err != nil {
		return fmt.Errorf("error while trying to ping InMemory DB: %v", err)
	}
	if err := onDiskConn.Ping(); err != nil {
		return fmt.Errorf("error while trying to ping OnDisk DB: %v", err)
	}

	return nil
}
