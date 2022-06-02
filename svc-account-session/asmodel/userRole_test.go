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
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

var role = Role{
	ID:                 "someID",
	AssignedPrivileges: []string{"somePrivilege"},
}

var invalidRole = Role{}

func TestCreateRole(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	err := role.Create()
	assert.Nil(t, err, "There should be no error")
}

func TestGetAllRoles(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	mockData(common.OnDisk, "role", role.ID, role)
	_, err := GetAllRoles()
	assert.Nil(t, err, "There should be no error")
}

func TestGetRole(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(common.OnDisk, "role", role.ID, role)
	type args struct {
		key string
	}
	tests := []struct {
		name                string
		args                args
		GetDBConnectionFunc func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error)
		want                Role
		wantErr             bool
	}{
		{
			name: "Db conn error",
			args: args{
				key: role.ID,
			},
			GetDBConnectionFunc: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return nil, &errors.Error{}
			},
			want:    Role{},
			wantErr: true,
		},
		{
			name: "success case",
			args: args{
				key: role.ID,
			},
			GetDBConnectionFunc: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
			want:    role,
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
			want:    Role{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		GetDBConnectionFunc = tt.GetDBConnectionFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRoleDetailsByID(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRoleDetailsByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRoleDetailsByID() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestDeleteRole(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(common.OnDisk, "role", role.ID, role)
	tests := []struct {
		name                string
		GetDBConnectionFunc func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error)
		want                *errors.Error
	}{
		{
			name: "Db conn error",
			GetDBConnectionFunc: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return nil, &errors.Error{}
			},
			want: &errors.Error{},
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
			want: errors.PackError(errors.DBKeyNotFound, "no data with the with key someID found"),
		},
	}
	for _, tt := range tests {
		GetDBConnectionFunc = tt.GetDBConnectionFunc
		t.Run(tt.name, func(t *testing.T) {
			err := role.Delete()
			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("Delete() = %v, want %v", err, tt.want)
			}
		})
	}

}

func TestUpdateRoleDetails(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	mockData(common.OnDisk, "role", role.ID, role)
	err := role.UpdateRoleDetails()
	assert.Nil(t, err, "There should be no error")
}

func TestUpdateRoleNegativeTestCase(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	mockData(common.OnDisk, "role", role.ID, "role")
	err := invalidRole.UpdateRoleDetails()
	assert.NotNil(t, err, "There should be an error")
}

func TestUpdateRoleDetailsDBError(t *testing.T) {
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	err := role.UpdateRoleDetails()
	assert.Equalf(t, &errors.Error{}, err, "UpdateRoleDetails() ")
}

func TestGetAllRolesDBError(t *testing.T) {
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	role, err := GetAllRoles()
	assert.Equalf(t, []Role(nil), role, "GetAllRoles() ")
	assert.Equalf(t, &errors.Error{}, err, "GetAllRoles() ")
}

func TestCreateRoleDBError(t *testing.T) {
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	err := role.Create()
	assert.Equalf(t, &errors.Error{}, err, "Create() ")
}
