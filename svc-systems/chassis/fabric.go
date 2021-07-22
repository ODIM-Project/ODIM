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

package chassis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-rest-client/pmbhandle"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/ODIM-Project/ODIM/svc-systems/sresponse"
	log "github.com/sirupsen/logrus"
)

type fabricFactory struct {
	collection        *sresponse.Collection
	chassisMap        map[string]bool
	wg                *sync.WaitGroup
	mu                *sync.RWMutex
	getFabricManagers func() ([]smodel.Plugin, error)
	contactClient     func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
}

func getFabricFactory(collection *sresponse.Collection) *fabricFactory {
	chassisMap := make(map[string]bool)
	return &fabricFactory{
		collection:        collection,
		chassisMap:        chassisMap,
		wg:                &sync.WaitGroup{},
		mu:                &sync.RWMutex{},
		getFabricManagers: smodel.GetFabricManagers,
		contactClient:     pmbhandle.ContactPlugin,
	}
}

type pluginContactRequest struct {
	URL             string
	HTTPMethodType  string
	ContactClient   func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	PostBody        interface{}
	LoginCredential map[string]string
	Plugin          smodel.Plugin
	Token           string
}

// PluginToken interface to hold the token
type PluginToken struct {
	Tokens map[string]string
	lock   sync.RWMutex
}

// Token variable hold the all the XAuthToken  against the plguin ID
var Token PluginToken

func (c *sourceProviderImpl) findFabricChassis(collection *sresponse.Collection) {
	f := c.getFabricFactory(collection)
	managers, err := f.getFabricManagers()
	if err != nil {
		log.Warn("while trying to collect fabric managers details from DB, got " + err.Error())
		return
	}
	for _, manager := range managers {
		f.wg.Add(1)
		go f.getFabricManagerChassis(manager)
	}
	f.wg.Wait()
}

// getFabricManagerChassis will send a request to the plugin for the chassis collection,
// and add them to the existing chassis collection.
func (f *fabricFactory) getFabricManagerChassis(plugin smodel.Plugin) {
	defer f.wg.Done()
	req, errResp, err := f.createChassisRequest(plugin, collectionURL, http.MethodGet, nil)
	if errResp != nil {
		log.Warn("while trying to create fabric plugin request for " + plugin.ID + ", got " + err.Error())
		return
	}
	links, err := collectChassisCollection(f, req)
	if err != nil {
		log.Warn("while trying to create fabric plugin request for " + plugin.ID + ", got " + err.Error())
		return
	}
	for _, link := range links {
		f.mu.Lock()
		if !f.chassisMap[link.Oid] { // uniqueness check for the chassis URI
			f.chassisMap[link.Oid] = true
			f.collection.AddMember(link)
		}
		f.mu.Unlock()
	}

}

// createChassisRequest creates the parameters ready for the plugin communication
func (f *fabricFactory) createChassisRequest(plugin smodel.Plugin, url, method string, body *json.RawMessage) (pReq *pluginContactRequest, errResp *response.RPC, err error) {
	var token string
	cred := make(map[string]string)

	if strings.EqualFold(plugin.PreferredAuthType, "XAuthToken") {
		token = f.getPluginToken(plugin)
		if token == "" {
			*errResp = common.GeneralError(http.StatusUnauthorized, response.ResourceAtURIUnauthorized, "unable to create session for plugin "+plugin.ID, []interface{}{url}, nil)
			return nil, errResp, fmt.Errorf("unable to create session for plugin " + plugin.ID)
		}
	} else {
		cred["UserName"] = plugin.Username
		cred["Password"] = string(plugin.Password)
	}

	// validating Patch request properties are in uppercamelcase or not
	if strings.EqualFold(method, http.MethodPatch) {
		errResp = validateReqParamsCase(body)
		if errResp != nil {
			return nil, errResp, fmt.Errorf("validation of request body failed")
		}
	}

	for key, value := range config.Data.URLTranslation.SouthBoundURL {
		if body != nil {
			*body = json.RawMessage(strings.Replace(string(*body), key, value, -1))
		}
		url = strings.Replace(url, key, value, -1)
	}

	pReq = &pluginContactRequest{
		Token:           token,
		LoginCredential: cred,
		ContactClient:   f.contactClient,
		Plugin:          plugin,
		HTTPMethodType:  method,
		URL:             url,
		PostBody:        body,
	}
	return pReq, nil, nil
}

// collectChassisCollection contacts the plugin and collect the chassis response
func collectChassisCollection(f *fabricFactory, pluginRequest *pluginContactRequest) ([]dmtf.Link, error) {
	body, _, statusCode, _, err := contactPlugin(pluginRequest)
	if statusCode == http.StatusUnauthorized && strings.EqualFold(pluginRequest.Plugin.PreferredAuthType, "XAuthToken") {
		body, _, statusCode, _, err = retryFabricsOperation(f, pluginRequest)
	}
	if err != nil {
		return []dmtf.Link{}, fmt.Errorf("while trying contact plugin " + pluginRequest.Plugin.ID + ", got " + err.Error())
	}
	if !is2xx(statusCode) {
		return []dmtf.Link{}, fmt.Errorf("while trying contact plugin " + pluginRequest.Plugin.ID + ", got " + strconv.Itoa(statusCode))
	}
	return extractChassisCollection(body)
}

func contactPlugin(req *pluginContactRequest) ([]byte, string, int, string, error) {
	pluginResponse, err := callPlugin(req)
	if err != nil {
		if getPluginStatus(req.Plugin) {
			pluginResponse, err = callPlugin(req)
		}
		if err != nil {
			return nil, "", http.StatusInternalServerError, response.InternalError, fmt.Errorf(err.Error())
		}
	}
	defer pluginResponse.Body.Close()
	body, err := ioutil.ReadAll(pluginResponse.Body)
	if err != nil {
		return nil, "", http.StatusInternalServerError, response.InternalError, fmt.Errorf(err.Error())
	}
	var statusMessage string
	switch pluginResponse.StatusCode {
	case http.StatusOK:
		statusMessage = response.Success
	case http.StatusUnauthorized:
		statusMessage = response.ResourceAtURIUnauthorized
	case http.StatusNotFound:
		statusMessage = response.ResourceNotFound
	default:
		statusMessage = response.CouldNotEstablishConnection
	}
	return body, pluginResponse.Header.Get("X-Auth-Token"), pluginResponse.StatusCode, statusMessage, nil
}

// retryFabricsOperation will be called whenever  the unauthorized status code during the plugin call
// This function will create a new session token reexcutes the plugin call
func retryFabricsOperation(f *fabricFactory, req *pluginContactRequest) ([]byte, string, int, string, error) {
	var resp response.RPC
	var token = f.createToken(req.Plugin)
	if token == "" {
		resp = common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "error: Unable to create session with plugin "+req.Plugin.ID,
			[]interface{}{}, nil)
		data, _ := json.Marshal(resp.Body)
		return data, "", int(resp.StatusCode), response.NoValidSession, fmt.Errorf("error: Unable to create session with plugin")
	}
	req.Token = token
	return contactPlugin(req)

}

func callPlugin(req *pluginContactRequest) (*http.Response, error) {
	var reqURL = "https://" + req.Plugin.IP + ":" + req.Plugin.Port + req.URL
	if strings.EqualFold(req.Plugin.PreferredAuthType, "BasicAuth") {
		return req.ContactClient(reqURL, req.HTTPMethodType, "", "", req.PostBody, req.LoginCredential)
	}
	return req.ContactClient(reqURL, req.HTTPMethodType, req.Token, "", req.PostBody, nil)
}

// getPluginStatus checks the status of given plugin in configured interval
func getPluginStatus(plugin smodel.Plugin) bool {
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
		log.Warn("while getting the status for plugin " + plugin.ID + err.Error())
		return status
	}
	log.Info("Status of plugin" + plugin.ID + strconv.FormatBool(status))
	return status
}

// getPluginToken will verify the if any token present to the plugin else it will create token for the new plugin
func (f *fabricFactory) getPluginToken(plugin smodel.Plugin) string {
	authToken := Token.getToken(plugin.ID)
	if authToken == "" {
		return f.createToken(plugin)
	}
	return authToken
}

func (f *fabricFactory) createToken(plugin smodel.Plugin) string {
	var contactRequest pluginContactRequest
	contactRequest.ContactClient = f.contactClient
	contactRequest.Plugin = plugin
	contactRequest.HTTPMethodType = http.MethodPost
	contactRequest.PostBody = map[string]interface{}{
		"Username": plugin.Username,
		"Password": string(plugin.Password),
	}
	contactRequest.URL = "/ODIM/v1/Sessions"
	_, token, _, _, err := contactPlugin(&contactRequest)
	if err != nil {
		log.Error(err.Error())
	}
	if token != "" {
		Token.storeToken(plugin.ID, token)
	}
	return token
}

func (p *PluginToken) storeToken(plguinID, token string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.Tokens[plguinID] = token
}

func (p *PluginToken) getToken(pluginID string) string {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.Tokens[pluginID]
}

// extractChassisCollection unmarshals the plugin response and returns the collection members
func extractChassisCollection(body []byte) ([]dmtf.Link, error) {
	var resp sresponse.Collection
	data := string(body)
	//replacing the resposne with north bound translation URL
	for key, value := range config.Data.URLTranslation.NorthBoundURL {
		data = strings.Replace(data, key, value, -1)
	}
	err := json.Unmarshal([]byte(data), &resp)
	if err != nil {
		return resp.Members, fmt.Errorf("while unmarshalling the chassis fabric collection, got: %v", err)
	}

	return resp.Members, nil

}

func is2xx(status int) bool {
	return status/100 == 2
}
