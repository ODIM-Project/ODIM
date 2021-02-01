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
	"encoding/json"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
)

// Account is the model for creating/updating an Account
type Account struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"`
	RoleID   string `json:"RoleId"`
}

// User is the model for User Account
type User struct {
	UserName     string   `json:"UserName"`
	Password     string   `json:"Password"`
	RoleID       string   `json:"RoleId"`
	AccountTypes []string `json:"AccountTypes"`
}

// CreateUser connects to the persistencemgr and creates a user in db
func CreateUser(user User) *errors.Error {

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	//Create a header for data entry
	const table string = "User"
	//Save data into Database
	return conn.Create(table, user.UserName, user)
}

//GetAllUsers gets all the accounts from the db
func GetAllUsers() ([]User, *errors.Error) {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return nil, err
	}
	keys, err := conn.GetAllDetails("User")
	if err != nil {
		return nil, err
	}
	var users []User
	//users := make(map[string]User)
	for _, key := range keys {
		var user User
		userdata, err := conn.Read("User", key)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(userdata), &user); err != nil {
			return nil, errors.PackError(errors.UndefinedErrorType, err)
		}
		users = append(users, user)

	}
	return users, nil
}

// GetUserDetails will fetch details of specific user from the db
func GetUserDetails(userName string) (User, *errors.Error) {
	var user User

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return user, err
	}
	userdata, err := conn.Read("User", userName)
	if err != nil {
		return user, errors.PackError(err.ErrNo(), "error while trying to get user: ", err.Error())
	}
	if jerr := json.Unmarshal([]byte(userdata), &user); jerr != nil {
		return user, errors.PackError(errors.UndefinedErrorType, jerr)
	}
	return user, nil

}

//DeleteUser will delete the user entry from the database based on the uuid
func DeleteUser(key string) *errors.Error {
	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	if err = conn.Delete("User", key); err != nil {
		return err
	}
	return nil
}

// UpdateUserDetails will modify the current details to given changes
func UpdateUserDetails(user, newData User) *errors.Error {

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	//Create a header for data entry
	const table string = "User"
	//Save data into Database

	if newData.Password != "" {
		user.Password = newData.Password
	}
	if newData.RoleID != "" {
		user.RoleID = newData.RoleID
	}
	if _, err = conn.Update(table, user.UserName, user); err != nil {
		return err
	}
	return nil
}
