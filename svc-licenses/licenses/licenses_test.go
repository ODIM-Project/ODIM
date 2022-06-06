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

package licenses

import (
	"fmt"
	"net/http"
	"testing"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	licenseproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/licenses"
	"github.com/stretchr/testify/assert"
)

var licenseServiceResponse = dmtf.LicenseService{
	OdataContext:   "/redfish/v1/$metadata#LicenseService.LicenseService",
	OdataID:        "/redfish/v1/LicenseService",
	OdataType:      "#LicenseService.v1_0_0.LicenseService",
	Description:    "License Service",
	Name:           "License Service",
	ServiceEnabled: true,
	ID:             "LicenseService",
	Licenses:       &dmtf.Link{Oid: "/redfish/v1/LicenseService/Licenses"},
}

var licenseCollectionResponse = dmtf.LicenseCollection{
	OdataContext: "/redfish/v1/$metadata#LicenseCollection.LicenseCollection",
	OdataID:      "/redfish/v1/LicenseService/Licenses",
	OdataType:    "#LicenseCollection.v1_0_0.LicenseCollection",
	Description:  "License Collection",
	Name:         "License Collection",
	Members: []*dmtf.Link{
		&dmtf.Link{
			Oid: "/redfish/v1/LicenseService/Licenses/uuid.1.1",
		},
		&dmtf.Link{
			Oid: "/redfish/v1/LicenseService/Licenses/uuid.1.2",
		},
	},
	MembersCount: 2,
}

var licenseResourceResponse = dmtf.License{
	OdataContext: "/redfish/v1/$metadata#License.License",
	OdataID:      "/redfish/v1/LicenseService/Licenses/uuid.1.1",
	OdataType:    "#License.v1_0_0.License",
	ID:           "uuid.1.1",
	Name:         "iLO License",
	LicenseType:  "Perpetual",
}

func TestGetLicenseService(t *testing.T) {
	req := &licenseproto.GetLicenseServiceRequest{}
	e := mockGetExternalInterface()
	response := e.GetLicenseService(req)

	license := response.Body.(dmtf.LicenseService)
	assert.Equal(t, int(response.StatusCode), http.StatusOK, "Status code should be StatusOK.")
	assert.Equal(t, licenseServiceResponse, license, "Valid License service response is expected")
}

func TestGetLicenseCollection(t *testing.T) {
	req := &licenseproto.GetLicenseRequest{}
	e := mockGetExternalInterface()
	response := e.GetLicenseCollection(req)

	license := response.Body.(dmtf.LicenseCollection)
	assert.Equal(t, int(response.StatusCode), http.StatusOK, "Status code should be StatusOK.")
	assert.Equal(t, licenseCollectionResponse, license, "Valid License collection response is expected")
}

func TestGetLicenseResource(t *testing.T) {
	req := &licenseproto.GetLicenseResourceRequest{
		URL: "/redfish/v1/LicenseService/Licenses/uuid.1.1",
	}
	e := mockGetExternalInterface()
	response := e.GetLicenseResource(req)
	fmt.Println(response)
	license := response.Body.(dmtf.License)
	assert.Equal(t, int(response.StatusCode), http.StatusOK, "Status code should be StatusOK.")
	assert.Equal(t, licenseResourceResponse, license, "Valid License resource response is expected")
}

func TestInstallLicenseService(t *testing.T) {
	req := &licenseproto.InstallLicenseRequest{
		RequestBody: []byte(`{
			"LicenseString": "333K2-8QLV4-Y8R3G-LJQQX-7BK6M",
			"Links": {
				"AuthorizedDevices": [{
					"@odata.id": "/redfish/v1/Systems/uuid.1"
				}]
			}
		}`)}
	e := mockGetExternalInterface()
	response := e.InstallLicenseService(req)

	assert.Equal(t, http.StatusNoContent, int(response.StatusCode), "Status code should be StatusNoContent.")
}

func TestInstallLicenseService_InvalidRequest(t *testing.T) {
	req := &licenseproto.InstallLicenseRequest{}
	e := mockGetExternalInterface()
	response := e.InstallLicenseService(req)

	assert.Equal(t, http.StatusBadRequest, int(response.StatusCode), "Status code should be StatusBadRequest.")
}

func TestInstallLicenseService_EmptyLinks(t *testing.T) {
	req := &licenseproto.InstallLicenseRequest{
		RequestBody: []byte(`{
			"LicenseString": "XYZ"
		}`)}
	e := mockGetExternalInterface()
	response := e.InstallLicenseService(req)

	assert.Equal(t, http.StatusBadRequest, int(response.StatusCode), "Status code should be StatusBadRequest.")
}

func TestInstallLicenseService_InvalidManager(t *testing.T) {
	req := &licenseproto.InstallLicenseRequest{
		RequestBody: []byte(`{
			"LicenseString": "XYZ",
			"Links": {
				"AuthorizedDevices": [{
					"@odata.id": "/redfish/v1/Systems/uuid.1.1"
				}]
			}
		}`)}
	e := mockGetExternalInterface()
	response := e.InstallLicenseService(req)

	assert.Equal(t, http.StatusInternalServerError, int(response.StatusCode), "Status code should be StatusInternalServerError.")
}

func TestInstallLicenseService_InvalidAuthorizedDevices(t *testing.T) {
	req := &licenseproto.InstallLicenseRequest{
		RequestBody: []byte(`{
			"LicenseString": "XYZ",
			"Links": {
				"AuthorizedDevices": [{
					"@odata.id": "/redfish/v1/invalid/uuid.1"
				}]
			}
		}`)}
	e := mockGetExternalInterface()
	response := e.InstallLicenseService(req)

	assert.Equal(t, http.StatusBadRequest, int(response.StatusCode), "Status code should be StatusBadRequest.")
}

func TestInstallLicenseService_ManagerURL(t *testing.T) {
	req := &licenseproto.InstallLicenseRequest{
		RequestBody: []byte(`{
			"LicenseString": "XYZ",
			"Links": {
				"AuthorizedDevices": [{
					"@odata.id": "/redfish/v1/Managers/uuid.1"
				}]
			}
		}`)}
	e := mockGetExternalInterface()
	response := e.InstallLicenseService(req)

	assert.Equal(t, http.StatusNoContent, int(response.StatusCode), "Status code should be StatusNoContent.")
}

