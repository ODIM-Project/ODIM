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

package account

import (
	"encoding/base64"
	"fmt"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/svc-account-session/asmodel"
	"golang.org/x/crypto/sha3"
	"testing"
)

func TestGetExternalInterface(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetExternalInterface(); got == nil {
				t.Errorf("Result of GetExternalInterface() should not be equal to nil")
			}
		})
	}
}

func getMockExternalInterface() *ExternalInterface {
	return &ExternalInterface{
		CreateUser:         mockCreateUser,
		GetUserDetails:     mockGetUserDetails,
		GetRoleDetailsByID: mockGetRoleDetailsByID,
		UpdateUserDetails:  mockUpdateUserDetails,
	}
}

func mockCreateUser(user asmodel.User) *errors.Error {
	if user.UserName == "existingUser" {
		return errors.PackError(errors.DBKeyAlreadyExist, "error: data with key existingUser already exists")
	}
	return nil
}

func mockGetUserDetails(userName string) (asmodel.User, *errors.Error) {
	hash := sha3.New512()
	hash.Write([]byte("P@$$w0rd"))
	hashSum := hash.Sum(nil)
	hashedPassword := base64.URLEncoding.EncodeToString(hashSum)

	user := asmodel.User{
		UserName: userName,
		Password: hashedPassword,
	}

	if userName == "testUser1" || userName == "testUser2" {
		user.RoleID = common.RoleAdmin
	} else if userName == "testUser3" {
		user.RoleID = "PrivilegeLogin"
	} else if userName == "operatorUser" {
		user.RoleID = common.RoleMonitor
	} else {
		return user, errors.PackError(errors.DBKeyNotFound, "error while trying to get user: ", fmt.Sprintf("no data with the with key %v found", userName))
	}
	return user, nil
}

func mockUpdateUserDetails(user, newData asmodel.User) *errors.Error {
	return nil
}

func mockGetRoleDetailsByID(roleID string) (asmodel.Role, *errors.Error) {
	if roleID == "xyz" {
		return asmodel.Role{}, errors.PackError(errors.DBKeyNotFound, "error while trying to get role details: ", fmt.Sprintf("error: Invalid RoleID %v present", roleID))
	}
	return asmodel.Role{}, nil
}
