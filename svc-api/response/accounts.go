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

//Package response ...
package response

import (
	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
)

//User struct is used to ommit password for display purposes
type User struct {
	OdataContext string  `json:"@odata.context"`
	Etag         string  `json:"@odata.etag,omitempty"`
	OdataID      string  `json:"@odata.id"`
	OdataType    string  `json:"@odata.type"`
	UserName     string  `json:"UserName"`
	RoleID       string  `json:"RoleID"`
	Password     *string `json:"Password"`
	ID           string  `json:"ID"`
	Name         string  `json:"Name"`
	Description  string  `json:"Description"`
	Links        Links   `json:"Links"`
	OEM          OEM     `json:"Oem"`
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
	OdataType                          string           `json:"@odata.type"`
	ID                                 string           `json:"Id"`
	Name                               string           `json:"Name"`
	Description                        string           `json:"Description"`
	Status                             Status           `json:"Status"`
	ServiceEnabled                     bool             `json:"ServiceEnabled"`
	AuthFailureLoggingThreshold        int              `json:"AuthFailureLoggingThreshold"`
	MinPasswordLength                  int              `json:"MinPasswordLength"`
	AccountLockoutThreshold            int              `json:"AccountLockoutThreshold"`
	AccountLockoutDuration             int              `json:"AccountLockoutDuration"`
	AccountLockoutCounterResetAfter    int              `json:"AccountLockoutCounterResetAfter"`
	Accounts                           Accounts         `json:"Accounts"`
	Roles                              Accounts         `json:"Roles"`
	OdataContext                       string           `json:"@odata.context"`
	OdataID                            string           `json:"@odata.id"`
	AccountLockoutCounterResetEnabled  bool             `json:"AccountLockoutCounterResetEnabled,omitempty"`
	Actions                            *dmtf.Actions    `json:"Actions,omitempty"`
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

//OAuth2 struct definition
type OAuth2 struct {
}

//ActiveDirectory struct definition
type ActiveDirectory struct {
}

//LDAP struct definition
type LDAP struct {
}

//TACACSplus struct definition
type TACACSplus struct {
}
