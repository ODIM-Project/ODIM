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

//Package rfputilities ...
package rfputilities

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	lutilconf "github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/plugin-redfish/config"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

//RedfishDeviceCollection struct definition
type RedfishDeviceCollection struct {
	RedfishDevices            []*RedfishDevice           `json:"targetHosts"`
	UnreachableRedfishDevices []UnreachableRedfishDevice `json:"failedHosts,omitempty"`
}

//RedfishDevice struct definition
type RedfishDevice struct {
	Host            string            `json:"hostAddress"`
	Username        string            `json:"username,omitempty"`
	Password        string            `json:"password,omitempty"`
	Token           string            `json:"token,omitempty"`
	Tags            []string          `json:"Tags"`
	RootNode        *dmtf.ServiceRoot `json:"rootNode,omitempty"`
	ComputerSystems []*Identifier
	PostBody        []byte `json:"PostBody,omitempty"`
	Location        string `json:"Location"`
}

//Identifier struct definition
type Identifier struct {
	UUID uuid.UUID
	URI  string
}

//UnreachableRedfishDevice struct definition
type UnreachableRedfishDevice struct {
	Host         string `json:"hostAddress"`
	ErrorMessage string `json:"errorMessage"`
}

//RedfishClient struct definition
type RedfishClient struct {
	httpClient *http.Client
}

// MarshalJSON Custom marshalling code used to prevent the display of password or authtoken
func (rfd RedfishDevice) MarshalJSON() ([]byte, error) {
	type redfishdevice RedfishDevice

	sanitizedRedfishDevice := redfishdevice(rfd)
	sanitizedRedfishDevice.Password = ""
	sanitizedRedfishDevice.Token = ""

	return json.Marshal(sanitizedRedfishDevice)
}

var redfishServiceRootURI = "/redfish/v1"

// GetRedfishClient : Returns a new RedfishClient with insecure flag set.
func GetRedfishClient() (*RedfishClient, error) {
	var err error
	newClient := RedfishClient{}
	httpConf := &lutilconf.HTTPConfig{
		CACertificate: &config.Data.KeyCertConf.RootCACertificate,
	}
	if newClient.httpClient, err = httpConf.GetHTTPClientObj(); err != nil {
		return nil, err
	}
	return &newClient, nil
}

// Get : Executes the REST call with the specified host and URI, then returns the response object.
func (client *RedfishClient) Get(device *RedfishDevice, requestURI string) (*http.Response, error) {
	endpoint := fmt.Sprintf("https://%s%s", device.Host, requestURI)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Close = true
	req.Header.Set("Accept", "application/json")
	if device.Token != "" {
		req.Header.Set("X-Auth-Token", device.Token)
	}
	req.Close = true
	lutilconf.TLSConfMutex.RLock()
	resp, err := client.httpClient.Do(req)
	lutilconf.TLSConfMutex.RUnlock()
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetRootService : Retrieves the ServiceRoot endpoint for the device and saves the return in the device object
func (client *RedfishClient) GetRootService(device *RedfishDevice) error {
	resp, err := client.Get(device, redfishServiceRootURI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("While trying to read response body, got: ", err.Error())
		return err
	}
	if resp.StatusCode >= 300 {
		errMessage := "Could not retrieve ServiceRoot for " + device.Host + ":" + string(body)
		log.Error(errMessage)
		return fmt.Errorf(errMessage)
	}
	serviceRoot := &dmtf.ServiceRoot{}
	json.Unmarshal(body, serviceRoot)
	device.RootNode = serviceRoot
	return nil
}

// AuthWithDevice : Performs authentication with the given device and saves the token
func (client *RedfishClient) AuthWithDevice(device *RedfishDevice) error {
	if device.RootNode == nil {
		return fmt.Errorf("No ServiceRoot found for device")
	}

	// TODO auth (Issue #22)
	endpoint := fmt.Sprintf("https://%s%s", device.Host, "/redfish/v1/SessionService/Sessions")

	var jsonStr = []byte(`{"UserName":"` + device.Username + `","Password":"` + string(device.Password) + `"}`)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonStr))
	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OData-Version", "4.0")

	req.Close = true

	lutilconf.TLSConfMutex.RLock()
	resp, err := client.httpClient.Do(req)
	lutilconf.TLSConfMutex.RUnlock()
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	device.Token = resp.Header["X-Auth-Token"][0]
	log.Debug("Token: " + device.Token)

	return nil
}

// BasicAuthWithDevice : Performs authentication with the given device and saves the token
func (client *RedfishClient) BasicAuthWithDevice(device *RedfishDevice, requestURI string) (*http.Response, error) {
	// if device.RootNode == nil {
	// 	return errors.New("no ServiceRoot found for device")
	// }
	endpoint := fmt.Sprintf("https://%s%s", device.Host, requestURI)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Close = true
	req.Header.Set("Accept", "application/json")
	auth := device.Username + ":" + string(device.Password)
	Basicauth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", Basicauth)
	req.Header.Add("Content-Type", "application/json")
	req.Close = true

	lutilconf.TLSConfMutex.RLock()
	resp, err := client.httpClient.Do(req)
	lutilconf.TLSConfMutex.RUnlock()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetWithBasicAuth : Performs authentication with the given device and saves the token
func (client *RedfishClient) GetWithBasicAuth(device *RedfishDevice, requestURI string) (*http.Response, error) {

	endpoint := fmt.Sprintf("https://%s%s", device.Host, requestURI)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Close = true
	req.Header.Set("Accept", "application/json")
	auth := device.Username + ":" + string(device.Password)
	Basicauth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", Basicauth)
	req.Header.Add("Content-Type", "application/json")
	req.Close = true

	lutilconf.TLSConfMutex.RLock()
	resp, err := client.httpClient.Do(req)
	lutilconf.TLSConfMutex.RUnlock()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SubscribeForEvents :Subscribes for events with Basic Auth
func (client *RedfishClient) SubscribeForEvents(device *RedfishDevice) (*http.Response, error) {
	endpoint := fmt.Sprintf("https://%s%s", device.Host, "/redfish/v1/EventService/Subscriptions")
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(device.PostBody))
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Accept", "application/json")
	auth := device.Username + ":" + string(device.Password)
	Basicauth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", Basicauth)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	var resp *http.Response
	lutilconf.TLSConfMutex.RLock()
	resp, err = client.httpClient.Do(req)
	lutilconf.TLSConfMutex.RUnlock()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// ResetComputerSystem :Reset the computer system with given ResetType
func (client *RedfishClient) ResetComputerSystem(device *RedfishDevice, uri string) (*http.Response, error) {

	endpoint := "https://" + device.Host + uri

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(device.PostBody))
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Accept", "application/json")
	auth := device.Username + ":" + string(device.Password)
	Basicauth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", Basicauth)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	var resp *http.Response
	lutilconf.TLSConfMutex.RLock()
	resp, err = client.httpClient.Do(req)
	lutilconf.TLSConfMutex.RUnlock()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SetDefaultBootOrder : sets default boot order
func (client *RedfishClient) SetDefaultBootOrder(device *RedfishDevice, uri string) (*http.Response, error) {

	endpoint := "https://" + device.Host + uri

	req, err := http.NewRequest(http.MethodPost, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Accept", "application/json")
	auth := device.Username + ":" + string(device.Password)
	Basicauth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", Basicauth)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	var resp *http.Response
	lutilconf.TLSConfMutex.RLock()
	resp, err = client.httpClient.Do(req)
	lutilconf.TLSConfMutex.RUnlock()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DeleteSubscriptionDetail will accepts device struct
// and it will delete the subscription detail
func (client *RedfishClient) DeleteSubscriptionDetail(device *RedfishDevice) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", device.Location, nil)
	if err != nil {
		return nil, err
	}

	req.Close = true
	auth := device.Username + ":" + string(device.Password)
	Basicauth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", Basicauth)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Close = true

	lutilconf.TLSConfMutex.RLock()
	resp, err := client.httpClient.Do(req)
	lutilconf.TLSConfMutex.RUnlock()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DeviceCall will call device with the given device details on the url given
// TODO: use same method to all other calls in this file
func (client *RedfishClient) DeviceCall(device *RedfishDevice, url, method string) (*http.Response, error) {
	endpoint := fmt.Sprintf("https://%s%s", device.Host, url)
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(device.PostBody))
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Accept", "application/json")
	auth := device.Username + ":" + string(device.Password)
	Basicauth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", Basicauth)
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	var resp *http.Response
	lutilconf.TLSConfMutex.RLock()
	resp, err = client.httpClient.Do(req)
	lutilconf.TLSConfMutex.RUnlock()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetSubscriptionDetail will accepts device struct
// and it will get the subscription detail
func (client *RedfishClient) GetSubscriptionDetail(device *RedfishDevice) (*http.Response, error) {
	req, err := http.NewRequest("GET", device.Location, nil)
	if err != nil {
		return nil, err
	}

	req.Close = true
	req.Header.Set("Accept", "application/json")
	auth := device.Username + ":" + string(device.Password)
	Basicauth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", Basicauth)
	req.Header.Add("Content-Type", "application/json")
	req.Close = true

	lutilconf.TLSConfMutex.RLock()
	resp, err := client.httpClient.Do(req)
	lutilconf.TLSConfMutex.RUnlock()
	if err != nil {
		return nil, err
	}

	return resp, nil
}
