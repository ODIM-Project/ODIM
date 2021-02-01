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

// Package asresponse ...
package asresponse

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// Account struct is used to ommit password for display purposes
type Account struct {
	response.Response
	UserName     string   `json:"UserName"`
	RoleID       string   `json:"RoleId"`
	AccountTypes []string `json:"AccountTypes"`
	Password     *string  `json:"Password"`
	Links        Links    `json:"Links"`
	OEM          *OEM     `json:"Oem,omitempty"`
}

//OEM struct definition
type OEM struct {
}

//Links struct definition
type Links struct {
	Role Role `json:"Role"`
}

//Role struct definition
type Role struct {
	OdataID string `json:"@odata.id"`
}

//AccountService struct definition
type AccountService struct {
	response.Response
	Status                          Status   `json:"Status"`
	ServiceEnabled                  bool     `json:"ServiceEnabled"`
	AuthFailureLoggingThreshold     int      `json:"AuthFailureLoggingThreshold"`
	MinPasswordLength               int      `json:"MinPasswordLength"`
	AccountLockoutThreshold         int      `json:"AccountLockoutThreshold"`
	AccountLockoutDuration          int      `json:"AccountLockoutDuration"`
	AccountLockoutCounterResetAfter int      `json:"AccountLockoutCounterResetAfter"`
	Accounts                        Accounts `json:"Accounts"`
	Roles                           Accounts `json:"Roles"`
}

//Accounts struct definition
type Accounts struct {
	OdataID string `json:"@odata.id"`
}
