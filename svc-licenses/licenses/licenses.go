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
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	licenseproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/licenses"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	lcommon "github.com/ODIM-Project/ODIM/svc-licenses/lcommon"
	"github.com/ODIM-Project/ODIM/svc-licenses/model"

	log "github.com/sirupsen/logrus"
)

var (
	JsonUnMarshalFunc = json.Unmarshal
	JsonMarshalFunc   = json.Marshal
)

// GetLicenseService to get license service details
func (e *ExternalInterface) GetLicenseService(req *licenseproto.GetLicenseServiceRequest) response.RPC {
	var resp response.RPC
	license := dmtf.LicenseService{
		OdataContext:   "/redfish/v1/$metadata#LicenseService.LicenseService",
		OdataID:        "/redfish/v1/LicenseService",
		OdataType:      "#LicenseService.v1_0_0.LicenseService",
		ID:             "LicenseService",
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
		err := json.Unmarshal([]byte(data.(string)), &licenseResp)
		if err != nil {
			log.Error("Unable to unmarshall  the data: " + err.Error())
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		}
	}
	licenseResp.OdataContext = "/redfish/v1/$metadata#License.License"
	licenseResp.OdataType = "#License.v1_0_0.License"
	licenseResp.OdataID = uri
	licenseResp.ID = ID[len(ID)-1]

	resp.Body = licenseResp
	resp.StatusCode = http.StatusOK
	return resp
}

// InstallLicenseService to install license
func (e *ExternalInterface) InstallLicenseService(req *licenseproto.InstallLicenseRequest) response.RPC {
	var resp response.RPC
	var contactRequest model.PluginContactRequest
	var installreq dmtf.LicenseInstallRequest

	genErr := JsonUnMarshalFunc(req.RequestBody, &installreq)
	if genErr != nil {
		errMsg := "Unable to unmarshal the install license request" + genErr.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.InternalError, errMsg, nil, nil)
	}

	if installreq.Links == nil || len(installreq.Links.Link) == 0 || installreq.LicenseString == "" {
		errMsg := "Invalid request, AuthorizedDevices links missing"
		log.Error(errMsg)
		return common.GeneralError(http.StatusBadRequest, response.InternalError, errMsg, nil, nil)
	}
	var serverURI string
	var err error
	var managerLink []string
	linksMap := make(map[string]bool)
	for _, serverIDs := range installreq.Links.Link {
		serverURI = serverIDs.Oid
		switch {
		case strings.Contains(serverURI, "Systems"):
			managerLink, err = e.getManagerURL(serverURI)
			if err != nil {
				errMsg := "Unable to get manager link"
				log.Error(errMsg)
				return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
			}
			for _, link := range managerLink {
				linksMap[link] = true
			}
		case strings.Contains(serverURI, "Managers"):
			linksMap[serverURI] = true
		case strings.Contains(serverURI, "Aggregates"):
			managerLink, err = e.getDetailsFromAggregate(serverURI)
			if err != nil {
				errMsg := "Unable to get manager link from aggregates"
				log.Error(errMsg)
				return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
			}
			for _, link := range managerLink {
				linksMap[link] = true
			}
		default:
			errMsg := "Invalid AuthorizedDevices links"
			log.Error(errMsg)
			return common.GeneralError(http.StatusBadRequest, response.InternalError, errMsg, nil, nil)
		}
	}
	log.Info("Map with manager Links: ", linksMap)

	for serverURI := range linksMap {
		uuid, managerID, err := lcommon.GetIDsFromURI(serverURI)
		if err != nil {
			errMsg := "error while trying to get system ID from " + serverURI + ": " + err.Error()
			log.Error(errMsg)
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"SystemID", serverURI}, nil)
		}
		// Get target device Credentials from using device UUID
		target, targetErr := e.External.GetTarget(uuid)
		if targetErr != nil {
			errMsg := targetErr.Error()
			log.Error(errMsg)
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"target", uuid}, nil)
		}

		decryptedPasswordByte, err := e.External.DevicePassword(target.Password)
		if err != nil {
			errMsg := "error while trying to decrypt device password: " + err.Error()
			log.Error(errMsg)
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
		}
		target.Password = decryptedPasswordByte

		// Get the Plugin info
		plugin, errs := e.External.GetPluginData(target.PluginID)
		if errs != nil {
			errMsg := "error while getting plugin data: " + errs.Error()
			log.Error(errMsg)
			return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg, []interface{}{"PluginData", target.PluginID}, nil)
		}
		log.Info("Plugin info: ", plugin)

		encodedKey := base64.StdEncoding.EncodeToString([]byte(installreq.LicenseString))
		managerURI := "/redfish/v1/Managers/" + managerID
		reqPostBody := map[string]interface{}{"LicenseString": encodedKey, "AuthorizedDevices": managerURI}
		reqBody, _ := json.Marshal(reqPostBody)

		contactRequest.Plugin = *plugin
		contactRequest.ContactClient = e.External.ContactClient
		contactRequest.Plugin.ID = target.PluginID
		contactRequest.HTTPMethodType = http.MethodPost

		if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
			contactRequest.DeviceInfo = map[string]interface{}{
				"UserName": plugin.Username,
				"Password": string(plugin.Password),
			}
			contactRequest.OID = "/ODIM/v1/Sessions"
			_, token, getResponse, err := lcommon.ContactPlugin(contactRequest, "error while logging in to plugin: ")
			if err != nil {
				errMsg := err.Error()
				log.Error(errMsg)
				return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, nil)
			}
			contactRequest.Token = token
		} else {
			contactRequest.LoginCredentials = map[string]string{
				"UserName": plugin.Username,
				"Password": string(plugin.Password),
			}

		}
		target.PostBody = []byte(reqBody)
		contactRequest.DeviceInfo = target
		contactRequest.OID = "/ODIM/v1/LicenseService/Licenses"
		contactRequest.PostBody = reqBody
		_, _, getResponse, err := e.External.ContactPlugin(contactRequest, "error while installing license: ")
		if err != nil {
			errMsg := err.Error()
			log.Error(errMsg)
			return common.GeneralError(getResponse.StatusCode, getResponse.StatusMessage, errMsg, getResponse.MsgArgs, nil)
		}
		log.Info("Install license response: ", getResponse)
	}

	resp.StatusCode = http.StatusNoContent
	return resp
}

func (e *ExternalInterface) getDetailsFromAggregate(aggregateURI string) ([]string, error) {
	var resource model.Elements
	var links []string
	respData, err := e.DB.GetResource("Aggregate", aggregateURI, persistencemgr.OnDisk)
	if err != nil {
		return nil, err
	}
	jsonStr, jerr := JsonMarshalFunc(respData)
	if jerr != nil {
		return nil, jerr
	}
	jerr = JsonUnMarshalFunc([]byte(jsonStr), &resource)
	if jerr != nil {
		return nil, jerr
	}
	log.Info("System URL's from agrregate: ", resource)

	for _, key := range resource.Elements {
		res, err := e.getManagerURL(key)
		if err != nil {
			errMsg := "Unable to get manager link"
			log.Error(errMsg)
			return nil, err
		}
		links = append(links, res...)
	}
	log.Info("manager links: ", links)
	return links, nil
}

func (e *ExternalInterface) getManagerURL(systemURI string) ([]string, error) {
	var resource dmtf.ComputerSystem
	var managerLink string
	var links []string
	respData, err := e.DB.GetResource("ComputerSystem", systemURI, persistencemgr.InMemory)
	if err != nil {
		return nil, err
	}
	jerr := JsonUnMarshalFunc([]byte(respData.(string)), &resource)
	if jerr != nil {
		return nil, jerr
	}
	members := resource.Links.ManagedBy
	for _, member := range members {
		managerLink = member.Oid
	}
	links = append(links, managerLink)

	return links, nil
}
