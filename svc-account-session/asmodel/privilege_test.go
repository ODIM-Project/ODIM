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
	"testing"

	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/stretchr/testify/assert"
)

var list = Privileges{
	List: []string{
		"Login",
		"ConfigureManager",
		"ConfigureUsers",
		"ConfigureSelf",
		"ConfigureComponents",
	},
}

var OEMList = OEMPrivileges{
	List: []string{},
}

func TestCreatePrivilegeRegistry(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	err := list.Create()
	assert.Nil(t, err, "There should be no error")
}

func TestGetPrivilegeRegistry(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	mockData(common.OnDisk, "registry", "assignedprivileges", list)
	_, err := GetPrivilegeRegistry()
	assert.Nil(t, err, "There should be no error")
}

func TestGetPrivilegeRegistryNegativeTestCase(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	_, err := GetPrivilegeRegistry()
	assert.NotNil(t, err, "There should be an error")
	mockData(common.OnDisk, "registry", "assignedprivileges", "list")
	_, err = GetPrivilegeRegistry()
	assert.NotNil(t, err, "There should be an error")
}

func TestCreateOEMPrivilegeRegistry(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	err := OEMList.Create()
	assert.Nil(t, err, "There should be no error")
}

func TestGetOEMPrivileges(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	mockData(common.OnDisk, "registry", "oemprivileges", OEMList)
	_, err := GetOEMPrivileges()
	assert.Nil(t, err, "There should be no error")
}

func TestGetOEMPrivilegesNegativeTestCase(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	_, err := GetOEMPrivileges()
	assert.NotNil(t, err, "There should be an error")
	mockData(common.OnDisk, "registry", "oemprivileges", "OEMList")
	_, err = GetOEMPrivileges()
	assert.NotNil(t, err, "There should be an error")
}

func TestCreateDBError(t *testing.T) {
	common.SetUpMockConfig()
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	err := list.Create()
	assert.Equalf(t, &errors.Error{}, err, " Privilege Create() ")
}

func TestGetOEMPrivilegesDBError(t *testing.T) {
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	oemPriv, err := GetOEMPrivileges()
	assert.Equalf(t, OEMPrivileges{}, oemPriv, "GetOEMPrivileges() ")
	assert.Equalf(t, &errors.Error{}, err, "GetOEMPrivileges() ")
}

func TestCreateOEMPrivilegeRegistryDBError(t *testing.T) {
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	err := OEMList.Create()
	assert.Equalf(t, &errors.Error{}, err, " OEM Create() ")
}

func TestGetPrivilegeRegistryDBError(t *testing.T) {
	GetDBConnectionFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	priv, err := GetPrivilegeRegistry()
	assert.Equalf(t, Privileges{}, priv, "GetPrivilegeRegistry() ")
	assert.Equalf(t, &errors.Error{}, err, "GetPrivilegeRegistry() ")
}
