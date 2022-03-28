//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

package model

// AccountService the supported properties,
// this structure should be updated once ODIMRA supports more properties
type AccountService struct {
	ODataContext           string `json:"@odata.context,omitempty"`
	ODataEtag              string `json:"@odata.etag,omitempty"`
	ODataID                string `json:"@odata.id"`
	ODataType              string `json:"@odata.type"`
	ID                     string `json:"Id"`
	Name                   string `json:"Name"`
	Description            string `json:"Description,omitempty"`
	Status                 Status `json:"Status,omitempty"`
	Accounts               Link   `json:"Accounts,omitempty"`
	Roles                  Link   `json:"Roles,omitempty"`
	MinPasswordLength      int    `json:"MinPasswordLength,omitempty"`
	MaxPasswordLength      int    `json:"MaxPasswordLength,omitempty"`
	PasswordExpirationDays int    `json:"PasswordExpirationDays,omitempty"`
	ServiceEnabled         bool   `json:"ServiceEnabled,omitempty"`
	LocalAccountAuth       string `json:"LocalAccountAuth,omitempty"`
}

// ManagerAccount the supported properties of manager account schema,
// this structure should be updated once ODIMRA supports more properties
type ManagerAccount struct {
	ODataContext           string       `json:"@odata.context,omitempty"`
	ODataEtag              string       `json:"@odata.etag,omitempty"`
	ODataID                string       `json:"@odata.id"`
	ODataType              string       `json:"@odata.type"`
	ID                     string       `json:"Id"`
	Name                   string       `json:"Name"`
	Description            string       `json:"Description,omitempty"`
	UserName               string       `json:"UserName,omitempty"`
	Password               string       `json:"Password,omitempty"`
	RoleID                 string       `json:"RoleId,omitempty"`
	Enabled                bool         `json:"Enabled,omitempty"`
	Locked                 bool         `json:"Locked,omitempty"`
	PasswordChangeRequired bool         `json:"PasswordChangeRequired,omitempty"`
	PasswordExpiration     string       `json:"PasswordExpiration,omitempty"`
	AccountExpiration      string       `json:"AccountExpiration,omitempty"`
	Links                  AccountLinks `json:"Links,omitempty"`
	AccountTypes           string       `json:"AccountTypes,omitempty"`
}

//AccountLinks struct definition
type AccountLinks struct {
	Role Link `json:"Role"`
}

// Role the supported properties of role schema,
// this structure should be updated once ODIMRA supports more properties
type Role struct {
	ODataContext       string   `json:"@odata.context,omitempty"`
	ODataEtag          string   `json:"@odata.etag,omitempty"`
	ODataID            string   `json:"@odata.id"`
	ODataType          string   `json:"@odata.type"`
	ID                 string   `json:"Id"`
	Name               string   `json:"Name"`
	Description        string   `json:"Description,omitempty"`
	AlternateRoleID    string   `json:"AlternateRoleId,omitempty"`
	AssignedPrivileges []string `json:"AssignedPrivileges,omitempty"`
	IsPredefined       bool     `json:"IsPredefined,omitempty"`
	Restricted         bool     `json:"Restricted,omitempty"`
	RoleID             string   `json:"RoleId,omitempty"`
}
