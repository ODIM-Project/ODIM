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
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

var user = User{
	UserName: "successID",
	Password: "SomePassword",
	RoleID:   "someRole",
}

func mockData(dbType common.DbType, table, id string, data interface{}) {
	connPool, _ := common.GetDBConnection(dbType)
	connPool.Create(table, id, data)
}

func TestCreate(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	err := CreateUser(user)
	assert.Nil(t, err, "There should be no error")
}

func TestGetAllUsers(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(common.OnDisk, "User", "successID", User{UserName: "successID"})
	_, err := GetAllUsers()
	assert.Nil(t, err, "There should be no error")
}

func TestGetUserDetails(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(common.OnDisk, "User", "successID", User{UserName: "successID"})
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    User
		wantErr bool
	}{
		{
			name: "success case",
			args: args{
				key: "successID",
			},
			want:    User{UserName: "successID"},
			wantErr: false,
		},
		{
			name: "not found case",
			args: args{
				key: "InvalidID",
			},
			want:    User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserDetails(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(common.OnDisk, "User", "successID", User{UserName: "successID"})
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want *errors.Error
	}{
		{
			name: "success case",
			args: args{
				key: "successID",
			},
			want: nil,
		},
		{
			name: "not found case",
			args: args{
				key: "InvalidID",
			},
			want: errors.PackError(errors.DBKeyNotFound, "no data with the with key InvalidID found"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeleteUser(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteUser() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestUpdateUserDetails(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	user := User{
		UserName: "successID",
		Password: "SomePassword",
		RoleID:   "someRole",
	}
	user1 := User{
		UserName: "user1",
		Password: "password1",
		RoleID:   common.RoleAdmin,
	}
	user2 := User{
		UserName: "user2",
		Password: "password",
		RoleID:   common.RoleMonitor,
	}
	user3 := User{
		UserName: "user3",
		Password: "password",
		RoleID:   common.RoleClient,
	}
	user4 := User{
		UserName: "user4",
		Password: "password",
		RoleID:   "testRole",
	}
	mockData(common.OnDisk, "User", "successID", user)
	type args struct {
		userData User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "positive case",
			args:    args{userData: User{UserName: "successID"}},
			wantErr: false,
		},
		{
			name:    "positive case1",
			args:    args{userData: user1},
			wantErr: false,
		},
		{
			name:    "positive case2",
			args:    args{userData: user2},
			wantErr: false,
		},
		{
			name:    "positive case3",
			args:    args{userData: user3},
			wantErr: false,
		},
		{
			name:    "positive case4",
			args:    args{userData: user4},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateUserDetails(user, tt.args.userData); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserDetails() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateUserDetailsNegativeTestCase(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	user := User{}
	mockData(common.OnDisk, "User", "successID", "user")
	userData := User{UserName: "successID"}
	err := UpdateUserDetails(user, userData)
	assert.NotNil(t, err, "There should be an error")
}
