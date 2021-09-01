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

//Package fabrics ...
package fabrics

import (
	"encoding/json"
	"fmt"
	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"github.com/ODIM-Project/ODIM/svc-fabrics/fabresponse"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-fabrics/fabmodel"
)

// Fabrics struct helps to hold the behaviours
type Fabrics struct {
	Auth          func(sessionToken string, privileges []string, oemPrivileges []string) response.RPC
	ContactClient func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
}

type pluginContactRequest struct {
	URL             string
	HTTPMethodType  string
	ContactClient   func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	PostBody        interface{}
	LoginCredential map[string]string
	Plugin          fabmodel.Plugin
	Token           string
}
type responseStatus struct {
	StatusCode    int32
	StatusMessage string
	Location      string
}

// PluginToken interface to hold the token
type PluginToken struct {
	Tokens map[string]string
	lock   sync.Mutex
}

// Zones struct to check request body cases
type Zones struct {
	Name     string     `json:"Name"`
	ZoneType string     `json:"ZoneType"`
	Links    dmtf.Links `json:"Links"`
}

//Endpoints struct to check request body cases
type Endpoints struct {
	Name        string       `json:"Name"`
	Description string       `json:"Description"`
	Redundancy  []Redundancy `json:"Redundancy"`
	Links       dmtf.Links   `json:"Links"`
}

//Redundancy struct to check request body cases
type Redundancy struct {
	Mode          string      `json:"Mode"`
	RedundencySet []dmtf.Link `json:"RedundencySet"`
}

// Token variable hold the all the XAuthToken  against the plguin ID
var Token PluginToken

func (p *PluginToken) storeToken(plguinID, token string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Tokens[plguinID] = token
}

func (p *PluginToken) getToken(pluginID string) string {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.Tokens[pluginID]
}

func contactPlugin(req pluginContactRequest, errorMessage string) ([]byte, string, responseStatus, error) {
	var resp responseStatus

	pluginResponse, err := callPlugin(req)
	if err != nil {
		if getPluginStatus(req.Plugin) {
			pluginResponse, err = callPlugin(req)
		}
		if err != nil {
			errorMessage = errorMessage + err.Error()
			resp.StatusCode = http.StatusInternalServerError
			resp.StatusMessage = response.InternalError
			log.Error(errorMessage)
			return nil, "", resp, fmt.Errorf(errorMessage)
		}
	}
	defer pluginResponse.Body.Close()
	if !(pluginResponse.StatusCode == http.StatusCreated || pluginResponse.StatusCode == http.StatusOK) {
		body, err := ioutil.ReadAll(pluginResponse.Body)
		if err != nil {
			errorMessage := "error while trying to read response body: " + err.Error()
			resp.StatusCode = http.StatusInternalServerError
			resp.StatusMessage = response.InternalError
			log.Error(errorMessage)
			return nil, "", resp, fmt.Errorf(errorMessage)
		}
		resp.StatusCode = int32(pluginResponse.StatusCode)
		log.Info("Read response successfully")
		return body, "", resp, fmt.Errorf("Read response successfully")
	}
	body, err := ioutil.ReadAll(pluginResponse.Body)
	resp.StatusCode = int32(pluginResponse.StatusCode)
	if err != nil {
		errorMessage := "error while trying to read response body: " + err.Error()
		resp.StatusCode = http.StatusInternalServerError
		resp.StatusMessage = response.InternalError
		return nil, "", resp, fmt.Errorf(errorMessage)
	}
	resp.Location = pluginResponse.Header.Get("Location")
	return body, pluginResponse.Header.Get("X-Auth-Token"), resp, nil
}

func callPlugin(req pluginContactRequest) (*http.Response, error) {
	var reqURL = "https://" + req.Plugin.IP + ":" + req.Plugin.Port + req.URL
	if strings.EqualFold(req.Plugin.PreferredAuthType, "BasicAuth") {
		req.ContactClient(reqURL, req.HTTPMethodType, "", "", req.PostBody, req.LoginCredential)
	}
	return req.ContactClient(reqURL, req.HTTPMethodType, req.Token, "", req.PostBody, nil)
}

// getPluginStatus checks the status of given plugin in configured interval
func getPluginStatus(plugin fabmodel.Plugin) bool {
	var pluginStatus = common.PluginStatus{
		Method: http.MethodGet,
		RequestBody: common.StatusRequest{
			Comment: "",
		},
		PluginIP:         plugin.IP,
		PluginPort:       plugin.Port,
		ResponseWaitTime: config.Data.PluginStatusPolling.ResponseTimeoutInSecs,
		Count:            config.Data.PluginStatusPolling.MaxRetryAttempt,
		RetryInterval:    config.Data.PluginStatusPolling.RetryIntervalInMins,
		CACertificate:    &config.Data.KeyCertConf.RootCACertificate,
	}
	status, _, _, err := pluginStatus.CheckStatus()
	if err != nil && !status {
		log.Error("Error While getting the status for plugin " + plugin.ID + err.Error())
		return status
	}
	log.Info("Status of plugin" + plugin.ID + strconv.FormatBool(status))
	return status
}

// getPluginToken will verify the if any token present to the plugin else it will create token for the new plugin
func (f *Fabrics) getPluginToken(plugin fabmodel.Plugin) string {
	authToken := Token.getToken(plugin.ID)
	if authToken == "" {
		return f.createToken(plugin)
	}
	return authToken
}

func (f *Fabrics) createToken(plugin fabmodel.Plugin) string {
	var contactRequest pluginContactRequest

	contactRequest.ContactClient = f.ContactClient
	contactRequest.Plugin = plugin
	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.PostBody = map[string]interface{}{
		"Username": plugin.Username,
		"Password": string(plugin.Password),
	}
	contactRequest.URL = "/ODIM/v1/Sessions"
	_, token, _, err := contactPlugin(contactRequest, "error while logging in to plugin: ")
	if err != nil {
		log.Error(err.Error())
	}
	if token != "" {
		Token.storeToken(plugin.ID, token)
	}
	return token
}

// retryFabricsOperation will be called whenever  the unauthorized status code during the plugin call
// This function will create a new session token reexcutes the plugin call
func (f *Fabrics) retryFabricsOperation(req pluginContactRequest, errorMessage string) ([]byte, string, responseStatus, error) {
	var resp response.RPC
	var token = f.createToken(req.Plugin)
	if token == "" {
		resp = common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "error: Unable to create session with plugin "+req.Plugin.ID,
			[]interface{}{}, nil)
		data, _ := json.Marshal(resp.Body)
		return data, "", responseStatus{
			StatusCode: resp.StatusCode,
		}, fmt.Errorf("error: Unable to create session with plugin")
	}
	req.Token = token
	return contactPlugin(req, errorMessage)

}

func (f *Fabrics) parseFabricsRequest(req *fabricsproto.FabricRequest) (pluginContactRequest, response.RPC, error) {
	var contactRequest pluginContactRequest
	var resp response.RPC
	sessionToken := req.SessionToken
	authResp := f.Auth(sessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		errMsg := "error while trying to authenticate session"
		log.Error(errMsg)
		return contactRequest, authResp, fmt.Errorf(errMsg)
	}

	if req.URL == "/redfish/v1/Fabrics" {
		resp = getFabricCollection()
		return contactRequest, resp, nil
	}
	log.Info("Request url" + req.URL)
	fabID := getFabricID(req.URL)
	log.Info("Fabric UUID" + fabID)
	fabric, err := fabmodel.GetManagingPluginIDForFabricID(fabID)
	if err != nil {
		errMsg := fmt.Sprintf("error while trying to get fabric Data: %v", err.Error())
		log.Error(errMsg)
		resp = common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg,
			[]interface{}{"Plugin", "Fabric"}, nil)
		return contactRequest, resp, err
	}
	// Get the Plugin info
	plugin, errs := fabmodel.GetPluginData(fabric.PluginID)
	if errs != nil {
		errMsg := fmt.Sprintf("error while trying to get plugin Data: %v", errs.Error())
		log.Error(errMsg)
		resp = common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errMsg,
			[]interface{}{"Plugin", "Fabric"}, nil)
		return contactRequest, resp, errs
	}

	contactRequest.ContactClient = f.ContactClient
	contactRequest.Plugin = plugin
	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token := f.getPluginToken(plugin)
		if token == "" {
			var errorMessage = "error: Unable to create session with plugin " + plugin.ID
			log.Error(errorMessage)
			resp = common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage,
				[]interface{}{}, nil)
			return contactRequest, resp, fmt.Errorf(errorMessage)
		}
		contactRequest.Token = token
	} else {
		contactRequest.LoginCredential = map[string]string{
			"Username": plugin.Username,
			"Password": string(plugin.Password),
		}

	}

	// Validating Post/Patch request properties are in uppercamelcase or not
	if strings.EqualFold(req.Method, "POST") || strings.EqualFold(req.Method, "PATCH") {
		valResp, err := validateReqParamsCase(req)
		if err != nil {
			return contactRequest, valResp, err
		}
	}

	var reqURL string
	var reqData string
	//replacing the reruest url with south bound translation URL
	for key, value := range config.Data.URLTranslation.SouthBoundURL {
		reqURL = strings.Replace(req.URL, key, value, -1)
		reqData = strings.Replace(string(req.RequestBody), key, value, -1)
	}

	contactRequest.URL = reqURL
	contactRequest.HTTPMethodType = req.Method
	if !(req.Method == http.MethodGet || req.Method == http.MethodDelete) {
		err := json.Unmarshal([]byte(reqData), &contactRequest.PostBody)
		if err != nil {
			log.Error("error while trying to get JSON request body: " + err.Error())
			resp = common.GeneralError(http.StatusBadRequest, response.MalformedJSON,
				"error while trying to get JSON request body: "+err.Error(),
				[]interface{}{}, nil)
			return contactRequest, resp, fmt.Errorf("error while trying to get JSON request body: %v", err)
		}
	}
	return contactRequest, resp, nil
}

func (f *Fabrics) parseFabricsResponse(pluginRequest pluginContactRequest, reqURI string) response.RPC {
	var resp response.RPC
	var errorMessage = fmt.Sprintf("error while performing %s operation on %s: ", pluginRequest.HTTPMethodType, reqURI)
	var header = map[string]string{"Content-type": "application/json; charset=utf-8"}
	//contactPlugin
	body, _, getResponse, err := contactPlugin(pluginRequest, errorMessage)
	if err != nil {
		if getResponse.StatusCode == http.StatusUnauthorized && strings.EqualFold(pluginRequest.Plugin.PreferredAuthType, "XAuthToken") {
			if body, _, getResponse, err = f.retryFabricsOperation(pluginRequest, errorMessage); err != nil {
				data := string(body)
				//replacing the resposne with north bound translation URL
				for key, value := range config.Data.URLTranslation.NorthBoundURL {
					data = strings.Replace(data, key, value, -1)
				}
				resp.StatusCode = getResponse.StatusCode
				json.Unmarshal([]byte(data), &resp.Body)
				resp.Header = header
				return resp
			}
		} else {
			data := string(body)
			//replacing the resposne with north bound translation URL
			for key, value := range config.Data.URLTranslation.NorthBoundURL {
				data = strings.Replace(data, key, value, -1)
			}
			resp.StatusCode = getResponse.StatusCode
			json.Unmarshal([]byte(data), &resp.Body)
			resp.Header = header
			return resp
		}
	}
	return fillResponse(body, getResponse.Location, pluginRequest.HTTPMethodType, getResponse.StatusCode)
}

func fillResponse(body []byte, location string, method string, statusCode int32) response.RPC {
	var resp response.RPC
	data := string(body)
	//replacing the resposne with north bound translation URL
	for key, value := range config.Data.URLTranslation.NorthBoundURL {
		location = strings.Replace(location, key, value, -1)
		data = strings.Replace(data, key, value, -1)
	}
	if method != http.MethodDelete {
		var respData map[string]interface{}
		err := json.Unmarshal([]byte(data), &respData)
		if err != nil {
			log.Printf(err.Error())
			return common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(),
				[]interface{}{}, nil)
		}
		resp.Body = respData
	}

	resp.Header = map[string]string{
		"Allow":             `"GET", "PUT", "POST", "PATCH", "DELETE"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	if location != "" {
		resp.Header["Location"] = location
	}

	resp.StatusCode = statusCode
	resp.StatusMessage = response.Success
	return resp

}

func getFabricID(url string) string {
	data := strings.Split(url, "/redfish/v1/Fabrics/")
	if len(data) > 1 {
		fabricData := strings.Split(data[1], "/")
		return fabricData[0]
	}
	return ""
}

func getFabricCollection() response.RPC {
	var resp response.RPC
	// ignoring error since we are trying to get collection of fabrics
	// So even its errored out we have to return empty collection
	fabrics, _ := fabmodel.GetAllTheFabrics()
	fabricCollection := fabresponse.FabricCollection{
		OdataContext: "/redfish/v1/$metadata#FabricCollection.FabricCollection",
		OdataID:      "/redfish/v1/Fabrics",
		OdataType:    "#FabricCollection.FabricCollection",
		Description:  "Fabric Collection view",
		Name:         "Fabric Collection",
	}
	members := []dmtf.Link{}
	for _, fab := range fabrics {
		members = append(members, dmtf.Link{Oid: fmt.Sprintf("/redfish/v1/Fabrics/%s", fab.FabricUUID)})
	}
	fabricCollection.Members = members
	fabricCollection.MembersCount = len(members)
	resp.Header = map[string]string{
		"Allow":             `"GET"`,
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	resp.Body = fabricCollection
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}

// Validating if request properties are in uppercamelcase or not
func validateReqParamsCase(req *fabricsproto.FabricRequest) (response.RPC, error) {
	var resp response.RPC
	var fabricRequest interface{}
	//Checking the request type, whether it is Zones,AddressPool or Endpoints request
	if strings.Contains(req.URL, "/Zones") {
		fabricRequest = &Zones{}
	} else if strings.Contains(req.URL, "/AddressPools") {
		fabricRequest = &dmtf.AddressPool{}
	} else if strings.Contains(req.URL, "/Endpoints") {
		fabricRequest = &Endpoints{}
	}

	// parsing the fabricRequest
	err := json.Unmarshal(req.RequestBody, &fabricRequest)
	if err != nil {
		errMsg := "unable to parse the fabrics request" + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), fmt.Errorf(errMsg)
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.RequestBody, fabricRequest)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		log.Error(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil), fmt.Errorf(errMsg)
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
		return response, fmt.Errorf(errorMessage)
	}

	return resp, nil
}
