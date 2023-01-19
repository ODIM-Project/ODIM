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

package ucommon

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-update/umodel"
	"github.com/stretchr/testify/assert"
)

func mockGetTarget(id string) (*umodel.Target, *errors.Error) {
	var target umodel.Target
	target.PluginID = id
	target.DeviceUUID = "uuid"
	target.UserName = "admin"
	target.Password = []byte("password")
	target.ManagerAddress = "ip"
	return &target, nil
}

func mockGetPluginData(id string) (umodel.Plugin, *errors.Error) {
	var plugin umodel.Plugin
	plugin.IP = "ip"
	plugin.Port = "port"
	plugin.Username = "plugin"
	plugin.Password = []byte("password")
	plugin.PluginType = "Redfish"
	plugin.PreferredAuthType = "XAuthToken"
	return plugin, nil
}
func mockGetPluginBasicData(id string) (umodel.Plugin, *errors.Error) {
	var plugin umodel.Plugin
	plugin.IP = "ip"
	plugin.Port = "port"
	plugin.Username = "plugin"
	plugin.Password = []byte("password")
	plugin.PluginType = "Redfish"
	plugin.PreferredAuthType = "basic"
	return plugin, nil
}

func mockContactPlugin(ctx context.Context, req PluginContactRequest, errorMessage string) ([]byte, string, ResponseStatus, error) {
	var resp ResponseStatus
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success

	reqBody := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString("{\"abc\":\"abc\"}")),
	}
	body, _ := ioutil.ReadAll(reqBody.Body)
	return body, "1234", resp, nil
}

func mockInterface() *CommonInterface {
	return &CommonInterface{
		GetTarget:     mockGetTarget,
		GetPluginData: mockGetPluginData,
		ContactPlugin: mockContactPlugin,
	}
}

func stubDevicePassword(password []byte) ([]byte, error) {
	return password, nil
}
func mockDevicePassword(password []byte) ([]byte, error) {
	return password, &errors.Error{}
}

func mockContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, common.TransactionID, "xyz")
	ctx = context.WithValue(ctx, common.ActionID, "001")
	ctx = context.WithValue(ctx, common.ActionName, "xyz")
	ctx = context.WithValue(ctx, common.ThreadID, "0")
	ctx = context.WithValue(ctx, common.ThreadName, "xyz")
	ctx = context.WithValue(ctx, common.ProcessName, "xyz")
	return ctx
}

func mockContactClient(ctx context.Context, url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {

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
	if url == "https://localhost:9091/ODIM/v1/UpdateService/FirmwareInventory/1" && token == "12345" {
		body := `{"data": "/ODIM/v1/UpdateService/FirmwareInventory/1"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9093/ODIM/v1/UpdateService/FirmwareInventory/1" {
		body := `{"data": "/ODIM/v1/UpdateService/FirmwareInventory/1"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9092/ODIM/v1/UpdateService/FirmwareInventory/1" && token == "23456" {
		body := `{"data": "/ODIM/v1/UpdateService/FirmwareInventory/uuid"}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9091/ODIM/v1/UpdateService/FirmwareInventory/1" {
		body := `{"@odata.id": "/ODIM/v1/UpdateService/FirmwareInventory/1"}`
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

func TestGetResourceInfoFromDevice(t *testing.T) {
	config.SetUpMockConfig(t)
	ctx := mockContext()
	i := mockInterface()
	var req = ResourceInfoRequest{
		URL:            "/redfish/v1/UpdateService/FirmwareInventory/uuid.1",
		UUID:           "uuid",
		SystemID:       "1",
		ContactClient:  mockContactClient,
		DevicePassword: stubDevicePassword,
	}
	_, err := i.GetResourceInfoFromDevice(ctx, req)
	assert.Nil(t, err, "There should be an error")
	i.GetPluginData = mockGetPluginBasicData
	req = ResourceInfoRequest{
		URL:            "/redfish/v1/UpdateService/FirmwareInventory/uuid.1",
		UUID:           "uuid",
		SystemID:       "1",
		ContactClient:  mockContactClient,
		DevicePassword: stubDevicePassword,
	}
	_, err = i.GetResourceInfoFromDevice(ctx, req)
	assert.Nil(t, err, "There should be an error")

	req = ResourceInfoRequest{
		URL:            "/redfish/v1/UpdateService/FirmwareInventory/uuid.1",
		UUID:           "uuid",
		SystemID:       "1",
		ContactClient:  mockContactClient,
		DevicePassword: mockDevicePassword,
	}
	_, err = i.GetResourceInfoFromDevice(ctx, req)
	assert.NotNil(t, err, "There should be an error")

}

func TestContactPlugin(t *testing.T) {
	config.SetUpMockConfig(t)
	ctx := mockContext()
	CallPluginFunc = func(ctx context.Context, req PluginContactRequest) (*http.Response, error) {
		return &http.Response{Status: "Save", StatusCode: 204, Body: ioutil.NopCloser(bytes.NewBufferString("Dummy"))}, nil
	}

	_, _, _, err := ContactPlugin(ctx, PluginContactRequest{}, "Dumyy")
	assert.NotNil(t, err, "There should be error")
	CallPluginFunc = func(ctx context.Context, req PluginContactRequest) (*http.Response, error) {
		return &http.Response{Status: "Save", StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBufferString("Dummy"))}, nil
	}
	_, _, _, err = ContactPlugin(ctx, PluginContactRequest{}, "Dumyy")
	assert.Nil(t, err, "There should be no error")
	CallPluginFunc = func(ctx context.Context, req PluginContactRequest) (*http.Response, error) {
		return &http.Response{Status: "Save", StatusCode: 202, Body: ioutil.NopCloser(bytes.NewBufferString("Dummy"))}, nil
	}

	_, _, _, err = ContactPlugin(ctx, PluginContactRequest{}, "Dumyy")
	assert.Nil(t, err, "There should be error")
	CallPluginFunc = func(ctx context.Context, req PluginContactRequest) (*http.Response, error) {
		return &http.Response{Status: "Save", StatusCode: 401, Body: ioutil.NopCloser(bytes.NewBufferString("Dummy"))}, nil
	}
	_, _, _, err = ContactPlugin(ctx, PluginContactRequest{}, "Dumyy")
	assert.NotNil(t, err, "There should be error")
	IOUtilReadAllFunc = func(r io.Reader) ([]byte, error) {
		return nil, &errors.Error{}
	}
	CallPluginFunc = func(ctx context.Context, req PluginContactRequest) (*http.Response, error) {
		return &http.Response{Status: "Save", StatusCode: 401, Body: ioutil.NopCloser(bytes.NewBufferString("Dummy"))}, nil
	}

	_, _, _, err = ContactPlugin(ctx, PluginContactRequest{}, "Dumyy")
	assert.NotNil(t, err, "There should be error")
	CallPluginFunc = func(ctx context.Context, req PluginContactRequest) (*http.Response, error) {
		return &http.Response{Status: "Save", StatusCode: 401, Body: ioutil.NopCloser(bytes.NewBufferString("Dummy"))}, &errors.Error{}
	}
	_, _, _, err = ContactPlugin(ctx, PluginContactRequest{}, "Dumyy")
	assert.NotNil(t, err, "There should be error")
}

func Test_callPlugin(t *testing.T) {
	config.SetUpMockConfig(t)
	_, err := callPlugin(context.TODO(), PluginContactRequest{ContactClient: mockContactClient})
	assert.NotNil(t, err, "There should be an error ")
	_, err = callPlugin(context.TODO(), PluginContactRequest{ContactClient: mockContactClient, Plugin: umodel.Plugin{PreferredAuthType: "BasicAuth"}})
	assert.NotNil(t, err, "There should be an error ")
}
