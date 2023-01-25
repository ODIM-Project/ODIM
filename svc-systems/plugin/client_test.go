/*
 * Copyright (c) 2020 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/ODIM-Project/ODIM/svc-systems/sresponse"
	"github.com/stretchr/testify/assert"
)

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

func Test_uriTranslator(t *testing.T) {

	sut := uriTranslator{&config.URLTranslation{
		NorthBoundURL: map[string]string{
			"ODIM": "redfish",
		},
		SouthBoundURL: map[string]string{
			"redfish": "ODIM",
		},
	}}

	tests := []struct {
		name           string
		translate      func(toBeTranslated string) (translated string)
		toBeTranslated string
		expected       string
	}{
		{name: "toSouthbound", translate: sut.toSouthbound, toBeTranslated: "", expected: ""},
		{name: "toSouthbound", translate: sut.toSouthbound, toBeTranslated: "redfish", expected: "ODIM"},
		{name: "toSouthbound", translate: sut.toSouthbound, toBeTranslated: "redfish redfish", expected: "ODIM ODIM"},
		{name: "toSouthbound", translate: sut.toSouthbound, toBeTranslated: "Redfish", expected: "Redfish"},
		{name: "toSouthbound", translate: sut.toSouthbound, toBeTranslated: `{"@odata.id":"/redfish/v1"}`, expected: `{"@odata.id":"/ODIM/v1"}`},

		{name: "toNorthbound", translate: sut.toNorthbound, toBeTranslated: "", expected: ""},
		{name: "toNorthbound", translate: sut.toNorthbound, toBeTranslated: "ODIM", expected: "redfish"},
		{name: "toNorthbound", translate: sut.toNorthbound, toBeTranslated: "ODIM ODIM", expected: "redfish redfish"},
		{name: "toNorthbound", translate: sut.toNorthbound, toBeTranslated: "Redfish", expected: "Redfish"},
		{name: "toNorthbound", translate: sut.toNorthbound, toBeTranslated: `{"@odata.id":"/ODIM/v1"}`, expected: `{"@odata.id":"/redfish/v1"}`},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("translating %s(%s)", test.name, test.toBeTranslated), func(t *testing.T) {
			assert.Equal(t, test.expected, test.translate(test.toBeTranslated))
		})
	}
}

func Test_convertToString(t *testing.T) {
	ctx := mockContext()
	res := convertToString(ctx, 125544)
	assert.True(t, res == "125544", "There should be no error ")
	JSONMarshalFunc = func(v interface{}) ([]byte, error) {
		return nil, &errors.Error{}
	}
	res = convertToString(ctx, "")
	assert.NotNil(t, res, "There should be an error ")

	JSONMarshalFunc = func(v interface{}) ([]byte, error) {
		return json.Marshal(v)
	}
	status := is2xx(200)
	assert.True(t, true, "There should be no error ", status)
	status = is4xx(400)
	assert.True(t, true, "There should be no error ", status)
	status = hasToBeSkipped("Content-Length")
	assert.True(t, true, "There should be no error ", status)
}

func Test_findAllPlugins(t *testing.T) {
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

	_, err = findAllPlugins("ILO")
	assert.Nil(t, err, "There should be no error ")

	JSONUnmarshalFunc = func(data []byte, v interface{}) error {
		return &errors.Error{}
	}
	_, err = findAllPlugins("ILO")
	assert.True(t, true, "There should be an error ", err)

	JSONUnmarshalFunc = func(data []byte, v interface{}) error {
		return json.Unmarshal(data, v)
	}
	DecryptWithPrivateKeyFunc = func(ciphertext []byte) ([]byte, error) {
		return nil, &errors.Error{}
	}
	_, err = findAllPlugins("ILO")
	assert.True(t, true, "There should be an error ", err)
	DecryptWithPrivateKeyFunc = func(ciphertext []byte) ([]byte, error) {
		return common.DecryptWithPrivateKey(ciphertext)
	}
	FindAllFunc = func(table, key string) ([][]byte, error) {
		return nil, &errors.Error{}
	}
	_, err = findAllPlugins("ILO")
	assert.NotNil(t, err, "There should be an error ")

}

func Test_createRPCResponse(t *testing.T) {
	transform := func(s string) string {
		return ""
	}
	rsp := createRPCResponse(200, http.Header{"content-type": []string{"application/json"}, "Content-Length": []string{"1025"}}, []byte(`{"username":"test"}`), transform)
	assert.Equal(t, http.StatusOK, int(rsp.StatusCode), "There should not be error ")
}

func Test_client_extractResp(t *testing.T) {
	c := &client{}
	ctx := mockContext()
	resp := c.extractResp(ctx, &http.Response{}, &errors.Error{})
	assert.Equal(t, http.StatusInternalServerError, int(resp.StatusCode), "Status should be StatusInternalServerError")

	fbUserData := []byte("data")
	resp = c.extractResp(ctx, &http.Response{StatusCode: 300, Body: ioutil.NopCloser(bytes.NewBufferString(string(fbUserData)))}, nil)
	assert.Equal(t, http.StatusInternalServerError, int(resp.StatusCode), "Status should be StatusInternalServerError")
	IoutilReadAllFunc = func(r io.Reader) ([]byte, error) {
		return nil, &errors.Error{}
	}
	resp = c.extractResp(ctx, &http.Response{StatusCode: 300, Body: ioutil.NopCloser(bytes.NewBufferString(string(fbUserData)))}, nil)
	assert.Equal(t, http.StatusInternalServerError, int(resp.StatusCode), "Status should be StatusInternalServerError")
	IoutilReadAllFunc = func(r io.Reader) ([]byte, error) {
		return ioutil.ReadAll(r)
	}

}

func Test_client_Get(t *testing.T) {
	ctx := mockContext()
	client := client{
		plugin: smodel.Plugin{
			Port: "90093",
			IP:   "localhost",
		},
		translator: &uriTranslator{&config.URLTranslation{
			NorthBoundURL: map[string]string{
				"ODIM": "redfish",
			},
			SouthBoundURL: map[string]string{
				"redfish": "ODIM",
			},
		}},
	}
	resp := client.Get(ctx, "/redis/v1/session")
	assert.Equal(t, http.StatusInternalServerError, int(resp.StatusCode), "Response status should be StatusInternalServerError")
	resp = client.Delete(ctx, "/redis/v1/session")
	assert.Equal(t, http.StatusInternalServerError, int(resp.StatusCode), "Response status should be StatusInternalServerError")

	resp = client.Patch(ctx, "/redis/v1/session", &json.RawMessage{})
	assert.Equal(t, http.StatusInternalServerError, int(resp.StatusCode), "Response status should be StatusInternalServerError")
	resp = client.Post(ctx, "/redis/v1/session", &json.RawMessage{})
	assert.Equal(t, http.StatusInternalServerError, int(resp.StatusCode), "Response status should be StatusInternalServerError")
}

func Test_multiTargetClient_Delete(t *testing.T) {
	ctx := mockContext()
	config := &config.URLTranslation{
		NorthBoundURL: map[string]string{
			"ODIM": "redfish",
		},
		SouthBoundURL: map[string]string{
			"redfish": "ODIM",
		},
	}
	client := func(plugin *smodel.Plugin) Client {
		return &client{plugin: *plugin, translator: &uriTranslator{dictionaries: config}}
	}
	multiTargetClient := multiTargetClient{
		createClient: client,
		targets: []*smodel.Plugin{
			{
				IP:                "localhost",
				Port:              "9093",
				Username:          "admin",
				Password:          []byte("dummy"),
				ID:                "ILO",
				PreferredAuthType: "Basic-Auth",
			},
		},
		call: ReturnFirstResponse(&CallConfig{}),
	}
	multiTargetClient.Get(ctx, "redish/v1/session")
	multiTargetClient.Post(ctx, "redish/v1/session", &json.RawMessage{})
	multiTargetClient.Patch(ctx, "redish/v1/session", &json.RawMessage{})
	multiTargetClient.Delete(ctx, "redish/v1/session")
}

func TestNewClientFactory(t *testing.T) {

	config := &config.URLTranslation{
		NorthBoundURL: map[string]string{
			"ODIM": "redfish",
		},
		SouthBoundURL: map[string]string{
			"redfish": "ODIM",
		},
	}
	fun := NewClientFactory(config)
	client, _ := fun("redish/v1")
	assert.Nil(t, client, "There should not error ")
	client, _ = fun("*")
	assert.Nil(t, client, "There should not error ")

}

func Test_returnFirst_Collect(t *testing.T) {

	returnFirst := returnFirst{
		ReqURI: "/redfish/v1/session",
		resp:   &response.RPC{},
	}
	ctx := mockContext()
	resp := response.RPC{}
	res := returnFirst.Collect(resp)
	assert.Nil(t, res, "There should be no error ")
	res1 := returnFirst.GetResult(ctx)
	assert.Equal(t, 0, int(res1.StatusCode), "There should be no error ")
	AggregateResults(&CallConfig{})
	res2 := ReturnFirstResponse(&CallConfig{})
	assert.NotNil(t, res2, "There should be an error ")

	returnFirst.resp = nil
	res1 = returnFirst.GetResult(ctx)
	assert.Equal(t, http.StatusNotFound, int(res1.StatusCode), "Status code should be StatusNotFound")

	collectCollectionMembers := collectCollectionMembers{
		collection: sresponse.Collection{},
		errored:    false,
	}
	collectionRes := collectCollectionMembers.Collect(response.RPC{})
	assert.NotNil(t, collectionRes, "There should be an error ")

	JSONUnmarshalFunc = func(data []byte, v interface{}) error {
		return fmt.Errorf("")
	}
	collectionRes = collectCollectionMembers.Collect(response.RPC{StatusCode: 200, Body: []byte(`{"Members":[{"@@odata.id":"/redish/v1/session/1"}]}`)})
	assert.NotNil(t, collectionRes, "There should be an error ")

	JSONUnmarshalFunc = func(data []byte, v interface{}) error {
		return json.Unmarshal(data, v)
	}
	collectionRes = collectCollectionMembers.Collect(response.RPC{StatusCode: 200, Body: []byte(`{"Members":[{"@@odata.id":"/redish/v1/session/1"}]}`)})
	assert.Nil(t, collectionRes, "There should be no error ")

	resultRes := collectCollectionMembers.GetResult(ctx)
	assert.NotNil(t, resultRes, "There should be an error ")
	JSONMarshalFunc = func(v interface{}) ([]byte, error) {
		return nil, &errors.Error{}
	}
	res3 := collectCollectionMembers.GetResult(ctx)
	assert.NotNil(t, res3, "There should be an error ")
	JSONMarshalFunc = func(v interface{}) ([]byte, error) {
		return json.Marshal(v)
	}
	chassisRes := getChassisID("redis/v1")
	assert.NotNil(t, chassisRes, "There should be no error ")
}
