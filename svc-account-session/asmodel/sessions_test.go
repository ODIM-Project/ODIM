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

// Package asmodel ...
package asmodel

import (
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"

	"github.com/stretchr/testify/assert"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

var session = Session{
	ID:       "someID",
	Token:    "token",
	UserName: "successID",
	RoleID:   "someRole",
}
var invalidSession = Session{
	ID:       "invalidID",
	UserName: "invalidName",
	RoleID:   "invalidRole",
}

func TestPersist(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	err := session.Persist()
	assert.Nil(t, err, "There should be no error")
}
func TestGetAllSessionKeys(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	mockData(common.InMemory, "session", session.Token, session)
	_, err := GetAllSessionKeys()
	assert.Nil(t, err, "There should be no error")
}

func TestGetSession(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(common.InMemory, "session", session.Token, session)
	type args struct {
		key string
	}
	tests := []struct {
		name                string
		args                args
		GetDBConnectionFunc func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error)
		want                Session
		wantErr             bool
	}{
		{
			name: "DB Error",
			args: args{
				key: "token",
			},
			GetDBConnectionFunc: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return nil, &errors.Error{}
			},
			want:    Session{},
			wantErr: true,
		},
		{
			name: "success case",
			args: args{
				key: "token",
			},
			GetDBConnectionFunc: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
			want:    session,
			wantErr: false,
		},
		{
			name: "not found case",
			args: args{
				key: "InvalidID",
			},
			GetDBConnectionFunc: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
			want:    Session{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		GetDBConnectionFunc = tt.GetDBConnectionFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSession(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSession() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestDelete(t *testing.T) {
	config.SetUpMockConfig(t)
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	mockData(common.InMemory, "session", session.Token, session)
	tests := []struct {
		name                string
		GetDBConnectionFunc func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error)
		want                *errors.Error
	}{
		{
			name: "DB error",
			GetDBConnectionFunc: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return nil, errors.PackError(0, "fakeError : ", " fakeErr")
			},
			want: errors.PackError(0, "error while trying to connecting to DB: ", "fakeError :  fakeErr"),
		},
		{
			name: "success case",
			GetDBConnectionFunc: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
			want: nil,
		},
		{
			name: "not found case",
			GetDBConnectionFunc: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
			want: errors.PackError(errors.DBKeyNotFound, "error while trying to delete session: no data with the with key token found"),
		},
	}
	for _, tt := range tests {
		GetDBConnectionFunc = tt.GetDBConnectionFunc
		t.Run(tt.name, func(t *testing.T) {
			err := session.Delete()
			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("Delete() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	mockData(common.InMemory, "session", session.Token, session)
	err := session.Update()
	assert.Nil(t, err, "There should be no error")
}

func TestUpdateNegativeTestCase(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	mockData(common.InMemory, "session", session.Token, "session")
	err := invalidSession.Update()
	assert.NotNil(t, err, "There should be an error")
}

func TestPersistDBError(t *testing.T) {
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	err := session.Persist()
	assert.Equalf(t, errors.PackError(0, "error while trying to connecting to DB: ", ""), err, "Persist()")
}

func TestUpdateDBError(t *testing.T) {
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	err := session.Update()
	assert.Equalf(t, errors.PackError(0, "error while trying to connecting to DB: ", ""), err, "Persist()")

}

func TestGetAllSessionKeysDBError(t *testing.T) {
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	_, err := GetAllSessionKeys()
	assert.Equalf(t, errors.PackError(0, "error while trying to connecting to DB: ", ""), err, "Persist()")
}
