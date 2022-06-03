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

type LicenseCollection struct {
	OdataContext string  `json:"@odata.context,omitempty"`
	Etag         string  `json:"@odata.etag,omitempty"`
	OdataID      string  `json:"@odata.id"`
	OdataType    string  `json:"@odata.type"`
	Description  string  `json:"Description,omitempty"`
	Name         string  `json:"Name"`
	Members      []*Link `json:"Members"`
	MembersCount int     `json:"Members@odata.count"`
}

type License struct {
	OdataContext         string       `json:"@odata.context,omitempty"`
	OdataID              string       `json:"@odata.id"`
	OdataType            string       `json:"@odata.type"`
	ID                   string       `json:"Id"`
	Name                 string       `json:"Name"`
	Description          string       `json:"Description,omitempty"`
	AuthorizationScope   string       `json:"AuthorizationScope,omitempty"`
	Contact              *Contact     `json:"Contact,omitempty"`
	DownloadURI          string       `json:"DownloadURI,omitempty"`
	EntitlementId        string       `json:"EntitlementId,omitempty"`
	ExpirationDate       string       `json:"ExpirationDate,omitempty"`
	GracePeriodDays      int32        `json:"GracePeriodDays,omitempty"`
	InstallDate          string       `json:"InstallDate,omitempty"`
	LicenseInfoURI       string       `json:"LicenseInfoURI,omitempty"`
	LicenseOrigin        string       `json:"LicenseOrigin,omitempty"`
	LicenseString        string       `json:"LicenseString,omitempty"`
	LicenseType          string       `json:"LicenseType,omitempty"`
	Links                *LicenseLink `json:"Links,omitempty"`
	Manufacturer         string       `json:"Manufacturer,omitempty"`
	MaxAuthorizedDevices int32        `json:"MaxAuthorizedDevices,omitempty"`
	PartNumber           string       `json:"PartNumber,omitempty"`
	RemainingDuration    string       `json:"RemainingDuration,omitempty"`
	RemainingUseCount    int32        `json:"RemainingUseCount,omitempty"`
	Removable            bool         `json:"Removable,omitempty"`
	SerialNumber         string       `json:"SerialNumber,omitempty"`
	SKU                  string       `json:"SKU,omitempty"`
	Status               *Status      `json:"Status,omitempty"`
}

type Contact struct {
	ContactName  string `json:"ContactName,omitempty"`
	EmailAddress string `json:"EmailAddress,omitempty"`
	PhoneNumber  string `json:"PhoneNumber,omitempty"`
}

type LicenseLink struct {
	AuthorizedDevices []*Link `json:"AuthorizedDevices,omitempty"`
	Oem               *Oem    `json:"Oem,omitempty"`
}

type LicenseService struct {
	OdataContext                 string      `json:"@odata.context,omitempty"`
	Etag                         string      `json:"@odata.etag,omitempty"`
	OdataID                      string      `json:"@odata.id"`
	OdataType                    string      `json:"@odata.type"`
	Description                  string      `json:"Description,omitempty"`
	ID                           string      `json:"Id"`
	Name                         string      `json:"Name"`
	Actions                      *OemActions `json:"Actions,omitempty"`
	LicenseExpirationWarningDays int32       `json:"LicenseExpirationWarningDays,omitempty"`
	Licenses                     *Link       `json:"Licenses,omitempty"`
	ServiceEnabled               bool        `json:"ServiceEnabled,omitempty"`
}

type LicenseInstallRequest struct {
	LicenseString string             `json:"LicenseString,omitempty"`
	Links         *AuthorizedDevices `json:"Links,omitempty"`
}

type AuthorizedDevices struct {
	Link []*Link `json:"AuthorizedDevices,omitempty"`
}
