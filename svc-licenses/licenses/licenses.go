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
	"encoding/json"
	"net/http"
	"strings"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	licenseproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/licenses"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	log "github.com/sirupsen/logrus"
)

// GetLicenseService to get license service details
func (e *ExternalInterface) GetLicenseService(req *licenseproto.GetLicenseServiceRequest) response.RPC {
	var resp response.RPC
	license := dmtf.LicenseService{
		OdataContext:   "/redfish/v1/$metadata#LicenseService.LicenseService",
		OdataID:        "/redfish/v1/LicenseService",
		OdataType:      "#LicenseService.v1_0_0.LicenseService",
		Description:    "License Service",
		Name:           "License Service",
		ServiceEnabled: true,
	}
	license.Licenses = &dmtf.Link{Oid: "/redfish/v1/LicenseService/Licenses"}

	resp.Body = license
	resp.StatusCode = http.StatusOK
	return resp
}

// GetLicenseCollection to get license collection details
func (e *ExternalInterface) GetLicenseCollection(req *licenseproto.GetLicenseRequest) response.RPC {
	var resp response.RPC
	licenseCollection := dmtf.LicenseCollection{
		OdataContext: "/redfish/v1/$metadata#LicenseCollection.LicenseCollection",
		OdataID:      "/redfish/v1/LicenseService/Licenses",
		OdataType:    "#LicenseCollection.v1_0_0.LicenseCollection",
		Description:  "License Collection",
		Name:         "License Collection",
	}
	var members []*dmtf.Link

	licenseCollectionKeysArray, err := e.DB.GetAllKeysFromTable("Licenses", persistencemgr.InMemory)
	if err != nil || len(licenseCollectionKeysArray) == 0 {
		log.Error("odimra doesnt have Licenses")
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, err.Error(), nil, nil)
	}

	for _, key := range licenseCollectionKeysArray {
		members = append(members, &dmtf.Link{Oid: key})
	}
	licenseCollection.Members = members
	licenseCollection.MembersCount = len(members)
	resp.Body = licenseCollection
	resp.StatusCode = http.StatusOK
	return resp
}

// GetLicenseResource to get individual license resource
func (e *ExternalInterface) GetLicenseResource(req *licenseproto.GetLicenseResourceRequest) response.RPC {
	var resp response.RPC
	licenseResp := dmtf.License{}
	uri := req.URL
	ID := strings.Split(uri, "/")

	data, dbErr := e.DB.GetResource("Licenses", uri, persistencemgr.InMemory)
	if dbErr != nil {
		log.Error("Unable to get license data : " + dbErr.Error())
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, dbErr.Error(), nil, nil)
	}

	if data != "" {
		err := json.Unmarshal([]byte(data), &licenseResp)
		if err != nil {
			log.Error("Unable to unmarshall  the data: " + err.Error())
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		}
	}
	licenseResp.OdataContext = "/redfish/v1/$metadata#License.License"
	licenseResp.OdataType = "#License.v1_0_0.License"
	licenseResp.Description = "License"
	licenseResp.OdataID = uri
	licenseResp.ID = ID[len(ID)-1]

	resp.Body = licenseResp
	resp.StatusCode = http.StatusOK
	return resp
}
