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
	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
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
	Status                             Status           `json:"Status,omitempty"`
	ServiceEnabled                     bool             `json:"ServiceEnabled,omitempty"`
	AuthFailureLoggingThreshold        int              `json:"AuthFailureLoggingThreshold,omitempty"`
	MinPasswordLength                  int              `json:"MinPasswordLength,omitempty"`
	AccountLockoutThreshold            int              `json:"AccountLockoutThreshold,omitempty"`
	AccountLockoutDuration             int              `json:"AccountLockoutDuration,omitempty"`
	AccountLockoutCounterResetAfter    int              `json:"AccountLockoutCounterResetAfter,omitempty"`
	Accounts                           Accounts         `json:"Accounts,omitempty"`
	Roles                              Accounts         `json:"Roles,omitempty"`
	AccountLockoutCounterResetEnabled  bool             `json:"AccountLockoutCounterResetEnabled,omitempty"`
	Actions                            *dmtf.OemActions `json:"Actions,omitempty"`
	ActiveDirectory                    *ActiveDirectory `json:"ActiveDirectory,omitempty"`
	AdditionalExternalAccountProviders *dmtf.Link       `json:"AdditionalExternalAccountProviders,omitempty"`
	LDAP                               *LDAP            `json:"LDAP,omitempty"`
	LocalAccountAuth                   string           `json:"LocalAccountAuth,omitempty"`
	MaxPasswordLength                  int              `json:"MaxPasswordLength,omitempty"`
	OAuth2                             *OAuth2          `json:"OAuth2,omitempty"`
	Oem                                *OEM             `json:"Oem,omitempty"`
	PasswordExpirationDays             int              `json:"PasswordExpirationDays,omitempty"`
	PrivilegeMap                       *dmtf.Link       `json:"PrivilegeMap,omitempty"`
	RestrictedOemPrivileges            []string         `json:"RestrictedOemPrivileges,omitempty"`
	RestrictedPrivileges               []string         `json:"RestrictedPrivileges,omitempty"`
	SupportedAccountTypes              []string         `json:"SupportedAccountTypes,omitempty"`
	SupportedOEMAccountTypes           []string         `json:"SupportedOEMAccountTypes,omitempty"`
	TACACSplus                         *TACACSplus      `json:"TACACSplus,omitempty"`
}

//Accounts struct definition
type Accounts struct {
	OdataID string `json:"@odata.id"`
}

// OAuth2 struct definition
type OAuth2 struct {
}

// ActiveDirectory struct definition
type ActiveDirectory struct {
}

// LDAP struct definition
type LDAP struct {
}

// TACACSplus struct definition
type TACACSplus struct {
}
