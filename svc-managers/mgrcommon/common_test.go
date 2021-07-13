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

package mgrcommon

import (
	"bytes"
	"fmt"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrmodel"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockTarget() error {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	target := &mgrmodel.DeviceTarget{
		ManagerAddress: "10.24.0.14",
		Password:       []byte("Password"),
		UserName:       "admin",
		DeviceUUID:     "uuid",
		PluginID:       "GRF",
	}
	target1 := &mgrmodel.DeviceTarget{
		ManagerAddress: "10.24.0.13",
		Password:       []byte("Password"),
		UserName:       "admin",
		DeviceUUID:     "uuid1",
		PluginID:       "ILO",
	}
	target2 := &mgrmodel.DeviceTarget{
		ManagerAddress: "localhost",
		Password:       []byte("Password"),
		UserName:       "admin",
		DeviceUUID:     "uuid2",
		PluginID:       "INVALID",
	}
	const table string = "System"
	//Save data into Database
	if err = connPool.Create(table, target.DeviceUUID, target); err != nil {
		return err
	}
	if err = connPool.Create(table, target1.DeviceUUID, target1); err != nil {
		return err
	}
	if err = connPool.Create(table, target2.DeviceUUID, target2); err != nil {
		return err
	}
	return nil
}
func stubDevicePassword(password []byte) ([]byte, error) {
	return password, nil
}

func mockContactClient(url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {
	baseURI := "/redfish/v1"
	baseURI = TranslateToSouthBoundURL(baseURI)

	if url == "https://localhost:9091"+baseURI+"/Sessions" {
		body := `{"Token": "12345"}`
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			Header: http.Header{
				"X-Auth-Token": []string{"12345"},
			},
		}, nil
	} else if url == "https://localhost:9092"+baseURI+"/Sessions" {
		body := `{"Token": ""}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091"+baseURI+"/Managers/1/EthernetInterfaces" && token == "12345" {
		body := `{"data": "/ODIM/v1/Managers/1/EthernetInterfaces"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9093"+baseURI+"/Managers/1/EthernetInterfaces" {
		body := `{"data": "/ODIM/v1/Managers/1/EthernetInterfaces"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9092"+baseURI+"/Managers/1/EthernetInterfaces" && token == "23456" {
		body := `{"data": "/ODIM/v1/Managers/uuid/EthernetInterfaces"}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9091"+baseURI+"/Managers/1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia" {
		body := `{"data": "Success"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9091"+baseURI+"/Managers/1/VirtualMedia/1/Actions/VirtualMedia.EjectMedia" {
		body := `{"data": "Success"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9091"+baseURI+"/Managers/1/VirtualMedia/1" {
		body := `{"data": "Success"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func mockPluginData(t *testing.T, pluginID, PreferredAuthType, port string) error {
	password := getEncryptedKey(t, []byte("$2a$10$OgSUYvuYdI/7dLL5KkYNp.RCXISefftdj.MjbBTr6vWyNwAvht6ci"))
	plugin := mgrmodel.Plugin{
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

func TestGetResourceInfoFromDevice(t *testing.T) {
	Token.Tokens = make(map[string]string)

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
		URL:                   "/redfish/v1/Managers/uuid:1/EthernetInterfaces",
		UUID:                  "uuid",
		SystemID:              "1",
		ContactClient:         mockContactClient,
		DecryptDevicePassword: stubDevicePassword,
	}
	_, err = GetResourceInfoFromDevice(req)
	assert.Nil(t, err, "There should be no error getting data")
	req.UUID = "uuid1"
	req.URL = "/redfish/v1/Managers/uuid1:1/EthernetInterfaces"
	_, err = GetResourceInfoFromDevice(req)
	assert.Nil(t, err, "There should be no error getting data")

	req.UUID = "uuid"
	req.URL = "/redfish/v1/Managers/uuid:1/VirtualMedia/1"
	_, err = GetResourceInfoFromDevice(req)
	assert.Nil(t, err, "There should be no error getting data")
}

func TestGetResourceInfoFromDeviceInvalidPlugin(t *testing.T) {
	Token.Tokens = make(map[string]string)

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

	var err error
	err = mockPluginData(t, "INVALID", "BasicAuth", "2325")
	if err != nil {
		t.Fatalf("Error in creating mock PluginData :%v", err)
	}
	err = mockTarget()
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}

	// for test purpose will set lower values
	// default values will be set back again with config.SetUpMockConfig call
	config.Data.PluginStatusPolling.ResponseTimeoutInSecs = 1
	config.Data.PluginStatusPolling.MaxRetryAttempt = 1
	config.Data.PluginStatusPolling.RetryIntervalInMins = 0

	var req = ResourceInfoRequest{
		URL:                   "/redfish/v1/Managers/uuid:1/EthernetInterfaces",
		UUID:                  "uuid2",
		SystemID:              "1",
		ContactClient:         mockContactClient,
		DecryptDevicePassword: stubDevicePassword,
	}
	_, err = GetResourceInfoFromDevice(req)
	assert.NotNil(t, err, "There should be an error")
}

func TestGetResourceInfoFromDeviceWithInvalidPluginSession(t *testing.T) {
	Token.Tokens = make(map[string]string)

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
		URL:                   "/redfish/v1/Managers/uuid:1/EthernetInterfaces",
		UUID:                  "uuid",
		SystemID:              "1",
		ContactClient:         mockContactClient,
		DecryptDevicePassword: stubDevicePassword,
	}
	_, err = GetResourceInfoFromDevice(req)

	assert.NotNil(t, err, "There should be an error")
	Token.Tokens = map[string]string{
		"GRF": "23456",
	}
	_, err = GetResourceInfoFromDevice(req)
	assert.NotNil(t, err, "There should be an error")
}

func TestDeviceCommunication(t *testing.T) {
	Token.Tokens = make(map[string]string)

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
		URL:                   "/redfish/v1/Managers/uuid:1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia",
		UUID:                  "uuid",
		SystemID:              "1",
		ContactClient:         mockContactClient,
		DecryptDevicePassword: stubDevicePassword,
		HTTPMethod:            http.MethodPost,
		RequestBody:           []byte(`{"Image":"http://10.1.1.1/ISO"}`),
	}
	response := DeviceCommunication(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	req = ResourceInfoRequest{
		URL:                   "/redfish/v1/Managers/uuid:1/VirtualMedia/1/Actions/VirtualMedia.EjectMedia",
		UUID:                  "uuid",
		SystemID:              "1",
		ContactClient:         mockContactClient,
		DecryptDevicePassword: stubDevicePassword,
		HTTPMethod:            http.MethodPost,
		RequestBody:           []byte(`{"Image":"http://10.1.1.1/ISO"}`),
	}
	response = DeviceCommunication(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")
}
