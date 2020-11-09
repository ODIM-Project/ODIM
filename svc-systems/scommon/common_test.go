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

package scommon

import (
	"bytes"
	"fmt"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func mockTarget() error {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	target := &smodel.Target{
		ManagerAddress: "10.24.0.14",
		Password:       []byte("Password"),
		UserName:       "admin",
		DeviceUUID:     "uuid",
		PluginID:       "GRF",
	}
	target1 := &smodel.Target{
		ManagerAddress: "10.24.0.13",
		Password:       []byte("Password"),
		UserName:       "admin",
		DeviceUUID:     "uuid1",
		PluginID:       "ILO",
	}
	const table string = "System"
	//Save data into Database
	if err = connPool.Create(table, target.DeviceUUID, target); err != nil {
		return err
	}
	if err = connPool.Create(table, target1.DeviceUUID, target1); err != nil {
		return err
	}
	return nil
}

func stubDevicePassword(password []byte) ([]byte, error) {
	return password, nil
}

func stubDeviceInvalidPassword(password []byte) ([]byte, error) {
	return []byte(""), fmt.Errorf("error decrypting device password")
}

func mockInvalidContactClient(url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {
	return nil, fmt.Errorf("InvalidRequest")
}

func mockContactClient(url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {

	if url == "https://localhost:9091/ODIM/v1/Sessions" {
		body := `{"Token": "12345"}`
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			Header: http.Header{
				"X-Auth-Token": []string{"12345"},
			},
		}, nil
	} else if url == "https://localhost:9092/ODIM/v1/Sessions" {
		body := `{"Token": ""}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091/ODIM/v1/Systems/1/EthernetInterfaces" && token == "12345" {
		body := `{"data": "/ODIM/v1/Systems/1/EthernetInterfaces"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9093/ODIM/v1/Systems/1/EthernetInterfaces" {
		body := `{"data": "/ODIM/v1/Systems/1/EthernetInterfaces"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9092/ODIM/v1/Systems/1/EthernetInterfaces" && token == "23456" {
		body := `{"data": "/ODIM/v1/Systems/uuid/EthernetInterfaces"}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9091/ODIM/v1/Systems/1/LogServices" {
		body := `{"@odata.id": "/ODIM/v1/Systems/1/LogServices"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockPluginData(t *testing.T, pluginID, PreferredAuthType, port string) error {
	password := getEncryptedKey(t, []byte("$2a$10$OgSUYvuYdI/7dLL5KkYNp.RCXISefftdj.MjbBTr6vWyNwAvht6ci"))
	plugin := smodel.Plugin{
		IP:                "localhost",
		Port:              port,
		Username:          "admin",
		Password:          password,
		ID:                pluginID,
		PreferredAuthType: PreferredAuthType,
	}
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("Plugin", pluginID, plugin); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Plugin", err.Error())
	}
	return nil
}

func mockPluginStatus(plugin smodel.Plugin) bool {
	return true
}

func TestGetResourceInfoFromDevice(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t, "GRF", "XAuthToken", "9091")
	if err != nil {
		t.Fatalf("Error in creating mock PluginData :%v", err)
	}

	err = mockPluginData(t, "ILO", "BasicAuth", "9093")
	if err != nil {
		t.Fatalf("Error in creating mock PluginData :%v", err)
	}
	err = mockTarget()
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	var req = ResourceInfoRequest{
		URL:            "/redfish/v1/Systems/uuid:1/EthernetInterfaces",
		UUID:           "uuid",
		SystemID:       "1",
		ContactClient:  mockContactClient,
		DevicePassword: stubDevicePassword,
	}
	_, err = GetResourceInfoFromDevice(req, true)
	assert.Nil(t, err, "There should be no error getting data")
	req.UUID = "uuid1"
	req.URL = "/redfish/v1/Systems/uuid1:1/EthernetInterfaces"
	_, err = GetResourceInfoFromDevice(req, true)
	assert.Nil(t, err, "There should be no error getting data")
}

func TestGetResourceInfoFromDeviceWithInvalidPluginSession(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t, "GRF", "XAuthToken", "9092")
	if err != nil {
		t.Fatalf("Error in creating mock PluginData :%v", err)
	}

	err = mockTarget()
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	var req = ResourceInfoRequest{
		URL:            "/redfish/v1/Systems/uuid:1/EthernetInterfaces",
		UUID:           "uuid",
		SystemID:       "1",
		ContactClient:  mockContactClient,
		DevicePassword: stubDevicePassword,
	}
	_, err = GetResourceInfoFromDevice(req, true)

	assert.NotNil(t, err, "There should be an error")
	//PluginContactRequest.Token = "23456"
	_, err = GetResourceInfoFromDevice(req, true)
	assert.NotNil(t, err, "There should be an error")
}

func TestGetResourceInfoFromDeviceWithInvalidPluginData(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t, "falseData", "InvalidToken", "InvalidPort")
	if err != nil {
		t.Fatalf("Error in creating mock PluginData :%v", err)
	}

	err = mockTarget()
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	var req = ResourceInfoRequest{
		URL:            "/redfish/v1/Systems/uuid:1/EthernetInterfaces",
		UUID:           "uuid",
		SystemID:       "1",
		ContactClient:  mockContactClient,
		DevicePassword: stubDevicePassword,
	}
	_, err = GetResourceInfoFromDevice(req, true)
	assert.NotNil(t, err, "There should be an error")
}

func TestGetResourceInfoFromDeviceWithNoTarget(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t, "GRF", "BasicAuth", "9093")
	if err != nil {
		t.Fatalf("Error in creating mock PluginData :%v", err)
	}
	var req = ResourceInfoRequest{
		URL:            "/redfish/v1/Systems/uuid:1/EthernetInterfaces",
		UUID:           "uuid",
		SystemID:       "1",
		ContactClient:  mockContactClient,
		DevicePassword: stubDevicePassword,
	}
	_, err = GetResourceInfoFromDevice(req, true)
	assert.NotNil(t, err, "There should be an error")
}

func TestGetResourceInfoFromDeviceWithInvalidDevicePassword(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t, "GRF", "BasicAuth", "9093")
	if err != nil {
		t.Fatalf("Error in creating mock PluginData :%v", err)
	}

	err = mockTarget()
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	var req = ResourceInfoRequest{
		URL:            "/redfish/v1/Systems/uuid:1/EthernetInterfaces",
		UUID:           "uuid",
		SystemID:       "1",
		ContactClient:  mockContactClient,
		DevicePassword: stubDeviceInvalidPassword,
	}
	_, err = GetResourceInfoFromDevice(req, true)
	assert.NotNil(t, err, "There should be an error")
}

func TestContactPlugin(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t, "falseData", "BasicAuth", "InvalidPort")
	if err != nil {
		t.Fatalf("Error in creating mock PluginData :%v", err)
	}
	plugin := smodel.Plugin{}
	var contactRequest PluginContactRequest

	contactRequest.ContactClient = mockContactClient
	contactRequest.Plugin = plugin
	contactRequest.GetPluginStatus = mockPluginStatus
	_, _, _, err = ContactPlugin(contactRequest, "")
	assert.NotNil(t, err, "There should be an error")
}
