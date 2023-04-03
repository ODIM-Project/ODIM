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

// LicenseCollection - This resource shall represent a resource collection of
// License instances for a Redfish implementation.
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

// License - This resource shall represent a license for a Redfish implementation.
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
	EntitlementID        string       `json:"EntitlementId,omitempty"`
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

// Contact - This property shall contain an object containing information
// about the contact of the license.
type Contact struct {
	ContactName  string `json:"ContactName,omitempty"`
	EmailAddress string `json:"EmailAddress,omitempty"`
	PhoneNumber  string `json:"PhoneNumber,omitempty"`
}

// LicenseLink - This property shall contain links to resources that are
// related to but are not contained by, or subordinate to, this resource.
type LicenseLink struct {
	AuthorizedDevices []*Link `json:"AuthorizedDevices,omitempty"`
	Oem               *Oem    `json:"Oem,omitempty"`
}

// LicenseService - The LicenseService schema describes the license service and the properties for the service
// itself with a link to the collection of licenses.
// The license service also provides methods for installing licenses in a Redfish service.
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

// LicenseInstallRequest - Request for install License
type LicenseInstallRequest struct {
	LicenseString string             `json:"LicenseString"`
	Links         *AuthorizedDevices `json:"Links"`
}

// AuthorizedDevices - This property shall contain an array of links to devices that are authorized by the license.
// Clients can provide this property when installing a license to apply the license to specific devices.
// If not provided when installing a license, the service may determine the devices to which the license applies.
// This property shall not be present if the AuthorizationScope property contains the value `Service`.
type AuthorizedDevices struct {
	Link []*Link `json:"AuthorizedDevices"`
}
