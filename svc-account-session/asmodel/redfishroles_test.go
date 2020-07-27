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
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
)

var roles = RedfishRoles{
	List: []string{
		"SomeRole1",
		"SomeRole2",
	},
}

func TestCreateRedfishRoles(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	err := roles.Create()
	assert.Nil(t, err, "There should be no error")
}

func TestGetRedfishRoles(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockData(common.OnDisk, "roles", "redfishdefined", roles)
	_, err := GetRedfishRoles()
	assert.Nil(t, err, "There should be no error")
}

func TestGetRedfishRolesNegativeTestCase(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	_, err := GetRedfishRoles()
	assert.NotNil(t, err, "There should be an error")
	mockData(common.OnDisk, "roles", "redfishdefined", "roles")
	_, err = GetRedfishRoles()
	assert.NotNil(t, err, "There should be an error")
}
