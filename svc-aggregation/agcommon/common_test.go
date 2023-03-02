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

package agcommon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	uuid "github.com/satori/go.uuid"
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
func TestAddConnectionMethods(t *testing.T) {
	var e = DBInterface{
		GetAllKeysFromTableInterface: stubGetAllkeys,
		GetConnectionMethodInterface: stubGetConnectionMethod,
		AddConnectionMethodInterface: stubAddConnectionMethod,
		DeleteInterface:              stubDeleteConnectionMethod,
	}
	config.SetUpMockConfig(t)
	err := e.AddConnectionMethods(config.Data.ConnectionMethodConf)
	assert.Nil(t, err, "err should be nil")
}

func TestAddConnectionMethods_failGetAllKeys(t *testing.T) {
	var e = DBInterface{
		GetAllKeysFromTableInterface: invalidGetAllkeys,
		GetConnectionMethodInterface: stubGetConnectionMethod,
		AddConnectionMethodInterface: stubAddConnectionMethod,
		DeleteInterface:              stubDeleteConnectionMethod,
	}
	config.SetUpMockConfig(t)
	err := e.AddConnectionMethods(config.Data.ConnectionMethodConf)

	assert.NotNil(t, err, "error should be not nil")
}

func TestAddConnectionMethods_failGetConnectionMethod(t *testing.T) {
	var e = DBInterface{
		GetAllKeysFromTableInterface: stubGetAllkeys,
		GetConnectionMethodInterface: invalidGetConnectionMethod,
		AddConnectionMethodInterface: stubAddConnectionMethod,
		DeleteInterface:              stubDeleteConnectionMethod,
	}
	config.SetUpMockConfig(t)
	err := e.AddConnectionMethods(config.Data.ConnectionMethodConf)

	assert.NotNil(t, err, "error should be not nil")
}

func TestAddConnectionMethods_failConnectionMethodInterface(t *testing.T) {
	var e = DBInterface{
		GetAllKeysFromTableInterface: stubGetAllkeys,
		GetConnectionMethodInterface: stubGetConnectionMethod,
		AddConnectionMethodInterface: invalidConnectionMethod,
		DeleteInterface:              stubDeleteConnectionMethod,
	}
	config.SetUpMockConfig(t)
	err := e.AddConnectionMethods(config.Data.ConnectionMethodConf)

	assert.NotNil(t, err, "error should be not nil")
}

func TestAddConnectionMethods_failDeleteConnectionMethod(t *testing.T) {
	var e = DBInterface{
		GetAllKeysFromTableInterface: stubGetAllkeys,
		GetConnectionMethodInterface: stubGetConnectionMethod,
		AddConnectionMethodInterface: stubAddConnectionMethod,
		DeleteInterface:              invalidDeleteConnectionMethod,
	}
	config.SetUpMockConfig(t)
	err := e.AddConnectionMethods(config.Data.ConnectionMethodConf)

	assert.NotNil(t, err, "error should be not nil")
}

func TestAddConnectionMethods_failConnection(t *testing.T) {
	var e = DBInterface{
		GetAllKeysFromTableInterface: stubGetAllkeys,
		GetConnectionMethodInterface: stubGetConnectionMethod,
		AddConnectionMethodInterface: stubAddConnectionMethod,
		DeleteInterface:              stubDeleteConnectionMethod,
	}
	config.SetUpMockConfig(t)
	var ConnectionMethodConf = []config.ConnectionMethodConf{}
	config.Data.ConnectionMethodConf = ConnectionMethodConf
	err := e.AddConnectionMethods(config.Data.ConnectionMethodConf)
	assert.Nil(t, err, "err should be nil")
}

var connectionMethod = []string{"/redfish/v1/AggregationService/ConnectionMethods/1234567545691234f",
	"/redfish/v1/AggregationService/ConnectionMethods/1234567545691234g",
	"/redfish/v1/AggregationService/ConnectionMethods/1234567545691234h"}

func stubGetAllkeys(tableName string) ([]string, error) {
	return connectionMethod, nil
}

func invalidGetAllkeys(tableName string) ([]string, error) {
	return nil, errors.PackError(0, "error while trying to connecting to DB: ")
}

func invalidGetConnectionMethod(key string) (agmodel.ConnectionMethod, *errors.Error) {
	return agmodel.ConnectionMethod{}, errors.PackError(0, "error while trying to connecting to DB: ")
}

func invalidConnectionMethod(data agmodel.ConnectionMethod, key string) *errors.Error {
	return errors.PackError(0, "error while trying to connecting to DB: ")
}

func invalidDeleteConnectionMethod(table, key string, dbtype common.DbType) *errors.Error {
	return errors.PackError(0, "error while trying to connecting to DB: ")
}

func stubGetConnectionMethod(key string) (agmodel.ConnectionMethod, *errors.Error) {
	if key == "/redfish/v1/AggregationService/ConnectionMethods/1234567545691234f" {
		return agmodel.ConnectionMethod{
			ConnectionMethodType:    "Redfish",
			ConnectionMethodVariant: "Compute:BasicAuth:GRF:1.0.0",
			Links: agmodel.Links{
				AggregationSources: []agmodel.OdataID{},
			},
		}, nil
	}

	if key == "/redfish/v1/AggregationService/ConnectionMethods/1234567545691234g" {
		return agmodel.ConnectionMethod{
			ConnectionMethodType:    "Redfish",
			ConnectionMethodVariant: "Fabric:XAuthToken:FabricPlugin:1.0.0",
			Links: agmodel.Links{
				AggregationSources: []agmodel.OdataID{
					agmodel.OdataID{OdataID: "/redfish/v1/AggregationService/AggregationSources/1234656881231fg1"},
				},
			},
		}, nil
	}
	return agmodel.ConnectionMethod{
		ConnectionMethodType:    "Redfish",
		ConnectionMethodVariant: "Storage:BasicAuth:Stg1:1.0.0",
		Links: agmodel.Links{
			AggregationSources: []agmodel.OdataID{},
		},
	}, nil
}

func stubAddConnectionMethod(data agmodel.ConnectionMethod, key string) *errors.Error {
	ConnectionMethod := agmodel.ConnectionMethod{
		ConnectionMethodType:    "Redfish",
		ConnectionMethodVariant: "Compute:BasicAuth:GRF:1.0.0",
		Links: agmodel.Links{
			AggregationSources: []agmodel.OdataID{},
		},
	}
	connectionMethodURI := "/redfish/v1/AggregationService/ConnectionMethods/" + uuid.NewV4().String()

	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to connecting to DB: ", err.Error())

	}
	if err = connPool.Create("ConnectionMethod", connectionMethodURI, ConnectionMethod); err != nil {
		return errors.PackError(err.ErrNo(), "error while trying to create new  resource :ConnectionMethod  ", err.Error())
	}
	return nil
}

func stubDeleteConnectionMethod(table, key string, dbtype common.DbType) *errors.Error {

	return nil

}

func mockData(t *testing.T, dbType common.DbType, table, id string, data interface{}) {
	connPool, err := common.GetDBConnection(dbType)
	if err != nil {
		t.Fatalf("error: mockData() failed to DB connection: %v", err)
	}
	if err = connPool.Create(table, id, data); err != nil {
		t.Fatalf("error: mockData() failed to create entry %s-%s: %v", table, id, err)
	}
}

func TestGetStorageResources(t *testing.T) {
	config.SetUpMockConfig(t)
	storageURI := "/redfish/v1/Systems/12345677651245-12341/Storage"
	GetResourceDetailsFunc = func(key string) (string, *errors.Error) {
		return "", errors.PackError(0, "error while trying to connecting to DB: ")
	}
	ctx := mockContext()
	resp := GetStorageResources(ctx, storageURI)
	assert.NotNil(t, resp, "There should be an error ")
	GetResourceDetailsFunc = func(key string) (string, *errors.Error) {
		return string([]byte(`{"user":"name"}`)), nil
	}
	resp = GetStorageResources(ctx, storageURI)
	assert.NotNil(t, resp, "There should be no error ")
}

func TestGetStorageResources_invalidJson(t *testing.T) {
	config.SetUpMockConfig(t)
	storageURI := "/redfish/v1/Systems/12345677651245-12341/Storage"
	GetResourceDetailsFunc = func(key string) (string, *errors.Error) {
		return "", errors.PackError(0, "error while trying to connecting to DB: ")
	}
	ctx := mockContext()
	resp := GetStorageResources(ctx, storageURI)
	assert.NotNil(t, resp, "There should be an error ")
	GetResourceDetailsFunc = func(key string) (string, *errors.Error) {
		return string([]byte(`{"user":"name"}`)), nil
	}
	JSONUnMarshalFunc = func(data []byte, v interface{}) error {
		return &errors.Error{}
	}
	resp = GetStorageResources(ctx, storageURI)
	assert.NotNil(t, resp, "Invalid Json")
}

func stubDevicePassword(password []byte) ([]byte, error) {
	if bytes.Compare(password, []byte("passwordWithInvalidEncryption")) == 0 {
		return []byte{}, fmt.Errorf("password decryption failed")
	}
	return password, nil
}

func TestPluginHealthCheckInterface_GetPluginStatus(t *testing.T) {
	type args struct {
		plugin agmodel.Plugin
	}
	PluginHealthCheck := &PluginHealthCheckInterface{
		DecryptPassword: common.DecryptWithPrivateKey,
	}
	password, _ := stubDevicePassword([]byte("password"))
	pluginData := agmodel.Plugin{
		IP:                "duphost",
		Port:              "9091",
		Username:          "admin",
		Password:          password,
		PreferredAuthType: "BasicAuth",
		ManagerUUID:       "mgr-addr",
	}
	var p = []string{}
	tests := []struct {
		name  string
		phc   *PluginHealthCheckInterface
		args  args
		want  bool
		want1 []string
	}{
		{
			name: "test1",
			phc:  PluginHealthCheck,
			args: args{
				plugin: pluginData,
			},
			want:  false,
			want1: p,
		},
	}
	ctx := mockContext()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.phc.GetPluginStatus(ctx, tt.args.plugin)
			if got != tt.want {
				t.Errorf("PluginHealthCheckInterface.GetPluginStatus() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestDupPluginConf(t *testing.T) {
	config.SetUpMockConfig(t)
	type args struct {
		plugin agmodel.Plugin
	}
	PluginHealthCheck := &PluginHealthCheckInterface{
		DecryptPassword: common.DecryptWithPrivateKey,
	}
	password, _ := stubDevicePassword([]byte("password"))
	pluginData := agmodel.Plugin{
		IP:                "duphost",
		Port:              "9091",
		Username:          "admin",
		Password:          password,
		PreferredAuthType: "BasicAuth",
		ManagerUUID:       "mgr-addr",
	}
	var p = []string{}
	tests := []struct {
		name  string
		phc   *PluginHealthCheckInterface
		args  args
		want  bool
		want1 []string
	}{
		{
			name: "test1",
			phc:  PluginHealthCheck,
			args: args{
				plugin: pluginData,
			},
			want:  false,
			want1: p,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.phc.DupPluginConf()
		})
	}
}

func TestGetPluginManagedServers(t *testing.T) {
	config.SetUpMockConfig(t)
	type args struct {
		plugin agmodel.Plugin
	}
	PluginHealthCheck := &PluginHealthCheckInterface{
		DecryptPassword: common.DecryptWithPrivateKey,
	}
	password, _ := stubDevicePassword([]byte("password"))
	pluginData := agmodel.Plugin{
		IP:                "duphost",
		Port:              "9091",
		Username:          "admin",
		Password:          password,
		PreferredAuthType: "BasicAuth",
		ManagerUUID:       "mgr-addr",
	}
	tests := []struct {
		name string
		phc  *PluginHealthCheckInterface
		args args
		want []agmodel.Target
	}{
		{
			name: "test1",
			phc:  PluginHealthCheck,
			args: args{
				plugin: pluginData,
			},
			want: []agmodel.Target{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.phc.GetPluginManagedServers(pluginData)
			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("PluginManagedServers() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestLookupHost(t *testing.T) {
	config.SetUpMockConfig(t)

	ip, _, _, _ := LookupHost("10.0.0.0")
	assert.Equal(t, "10.0.0.0", ip, "Ip should be same")

	LookupIPfunc = func(host string) (ip []net.IP, err error) {
		err = fmt.Errorf("error")
		return
	}
	ip, _, _, err := LookupHost("10.0.0")
	assert.NotNil(t, err, "There should be an error")
	LookupIPfunc = func(host string) (ip []net.IP, err error) {
		return
	}
	ip, _, _, err = LookupHost("10.0.0.1")
	assert.NotNil(t, err, "There should be an error")

}
func TestLookupPlugin(t *testing.T) {
	config.SetUpMockConfig(t)
	ctx := mockContext()
	var data = agmodel.Plugin{IP: "", Port: "", Username: "", Password: []uint8(nil), ID: "", PluginType: "", PreferredAuthType: "", ManagerUUID: ""}
	res, _ := LookupPlugin(ctx, "10.0.0.0")
	assert.Equal(t, res, data, "It should be same")

}

func TestGetDeviceSubscriptionDetails(t *testing.T) {
	config.SetUpMockConfig(t)
	var data = ""
	res, _, _ := GetDeviceSubscriptionDetails("10.0.0.0")
	assert.Equal(t, res, data, "It should be same")

}

func mockGetSearchKey(key, index string) string {
	return "10.0.0.0"
}

func mockLookupHost(serverAddress string) (ip string, host string, port string, err error) {
	return "10.0.0.0", "", "", nil
}

func TestRemoveDuplicates(t *testing.T) {
	config.SetUpMockConfig(t)
	var data = []string{"example", "example", "example2"}
	res := removeDuplicates(data)
	assert.Equal(t, []string{"example", "example2"}, res, "It should remove duplicates")

}

func TestGetSearchKey(t *testing.T) {
	config.SetUpMockConfig(t)
	res := GetSearchKey("100.100.100.100", "0")
	assert.Equal(t, string("100.100.100.100"), res, "It should remove duplicates")

}

func TestGetSubscribedEvtTypes(t *testing.T) {
	config.SetUpMockConfig(t)
	res, _ := GetSubscribedEvtTypes("100.100.100.100")
	assert.Equal(t, []string([]string{}), res, "It should be same")

}

func TestGetSubscribedEvtTypes_fail(t *testing.T) {
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
	mockGetEventSubscriptionsFunc("*" + "100.100.100.100" + "*")
	res, _ := GetSubscribedEvtTypes("100.100.100.100")
	assert.Equal(t, []string([]string{}), res, "It should be same")

}

func mockGetEventSubscriptionsFunc(key string) []string {
	return []string{"100.100.100.100"}
}

func TestUpdateDeviceSubscriptionDetails(t *testing.T) {
	config.SetUpMockConfig(t)
	ctx := mockContext()
	var data = make(map[string]string)
	data["100.100.100.100"] = "location"
	UpdateDeviceSubscriptionDetails(ctx, data)
}

func TestGetPluginStatusRecord(t *testing.T) {
	config.SetUpMockConfig(t)
	_, got := GetPluginStatusRecord("ILO")
	assert.Equal(t, false, got, "It should be same")

}

func TestSetPluginStatusRecord(t *testing.T) {
	config.SetUpMockConfig(t)

	SetPluginStatusRecord("ILO", 1)
}

func mockGetPluginStatus(plugin agmodel.Plugin) bool {
	return true
}

func TestGetAllPlugins(t *testing.T) {
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
	mockPlugins(t)
	ctx := mockContext()
	plugins, err := GetAllPlugins(ctx)
	assert.Nil(t, err, "Error Should be nil")
	assert.Equal(t, 3, len(plugins), "should be only 3 plugins")
}

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockPlugins(t *testing.T) {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		t.Errorf("error while trying to connecting to DB: %v", err.Error())
	}

	password := getEncryptedKey(t, []byte("Password"))
	pluginArr := []agmodel.Plugin{
		{
			IP:                "localhost",
			Port:              "1234",
			Password:          password,
			Username:          "admin",
			ID:                "GRF",
			PreferredAuthType: "BasicAuth",
			PluginType:        "GRF",
		},
		{
			IP:                "localhost",
			Port:              "1234",
			Password:          password,
			Username:          "admin",
			ID:                "ILO",
			PreferredAuthType: "XAuthToken",
			PluginType:        "ILO",
		},
		{
			IP:                "localhost",
			Port:              "1234",
			Password:          password,
			Username:          "admin",
			ID:                "CFM",
			PreferredAuthType: "XAuthToken",
			PluginType:        "CFM",
		},
	}
	for _, plugin := range pluginArr {
		pl := "Plugin"
		//Save data into Database
		if err := connPool.Create(pl, plugin.ID, &plugin); err != nil {
			t.Fatalf("error: %v", err)
		}
	}
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
	plugin := agmodel.Plugin{}
	var contactRequest agmodel.PluginContactRequest

	contactRequest.ContactClient = mockContactClient
	contactRequest.Plugin = plugin
	_, err = ContactPlugin(context.TODO(), contactRequest, "")
	assert.NotNil(t, err, "There should be an error")
}

func TestContactPlugin_XAuth(t *testing.T) {
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
	err := mockPluginData(t, "falseData", "XAuthToken", "InvalidPort")
	if err != nil {
		t.Fatalf("Error in creating mock PluginData :%v", err)
	}
	plugin := agmodel.Plugin{PreferredAuthType: "XAuthToken"}
	var contactRequest agmodel.PluginContactRequest

	contactRequest.ContactClient = mockContactClient
	contactRequest.Plugin = plugin
	_, err = ContactPlugin(context.TODO(), contactRequest, "")
	assert.NotNil(t, err, "There should be an error")
}

type PluginContactRequest struct {
	Token           string
	OID             string
	DeviceInfo      interface{}
	BasicAuth       map[string]string
	ContactClient   func(string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	GetPluginStatus func(agmodel.Plugin) bool
	Plugin          agmodel.Plugin
	HTTPMethodType  string
}

func (phc PluginHealthCheckInterface) mockGetPluginStatus(plugin agmodel.Plugin) bool {
	return true
}

func mockPluginData(t *testing.T, pluginID, PreferredAuthType, port string) error {
	password := getEncryptedKey(t, []byte("$2a$10$OgSUYvuYdI/7dLL5KkYNp.RCXISefftdj.MjbBTr6vWyNwAvht6ci"))
	plugin := agmodel.Plugin{
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

func TestGetAllServers(t *testing.T) {
	config.SetUpMockConfig(t)
	st := StartUpInteraface{
		DecryptPassword: stubDevicePassword,
		GetAllSystems:   MockGetAllSystems,
		GetSingleSystem: MockGetSingleSystem,
	}
	servers, err := st.getAllServers("ILO")
	assert.Nil(t, err, "Error Should be nil")
	assert.Equal(t, 2, len(servers), "there should be 2 server")

}

func TestGetAllServers_fail(t *testing.T) {
	config.SetUpMockConfig(t)
	st := StartUpInteraface{
		DecryptPassword: stubDevicePassword,
		GetAllSystems:   MockInvalidGetAllSystems,
		GetSingleSystem: MockGetSingleSystem,
	}
	_, err := st.getAllServers("ILO")
	assert.Equal(t, err, &errors.Error{})
}

type SavedSystems struct {
	ManagerAddress string
	Password       []byte
	UserName       string
	DeviceUUID     string
	PluginID       string
}

func (st *StartUpInteraface) getAllServers(pluginID string) ([]SavedSystems, error) {
	var matchedServers []SavedSystems
	allServers, err := st.GetAllSystems()
	if err != nil {
		return matchedServers, err
	}
	for i := 0; i < len(allServers); i++ {
		var s SavedSystems
		singleServer, err := st.GetSingleSystem(allServers[i])
		if err != nil {
			continue
		}
		json.Unmarshal([]byte(singleServer), &s)
		if s.PluginID == pluginID {
			decryptedPasswordByte, err := st.DecryptPassword(s.Password)
			if err != nil {
				errorMessage := "error while trying to decrypt device password for the host: " + s.ManagerAddress + ":" + err.Error()
				l.Log.Error(errorMessage)
				continue
			}
			s.Password = decryptedPasswordByte
			matchedServers = append(matchedServers, s)
		}
	}
	return matchedServers, err
}

// MockGetSingleSystem is for mocking up of get system info
func MockGetSingleSystem(id string) (string, error) {
	var systemData SavedSystems
	switch id {
	case "6d4a0a66-7efa-578e-83cf-44dc68d2874e":
		systemData = SavedSystems{
			ManagerAddress: "100.100.100.100",
			Password:       []byte("Password"),
			UserName:       "admin",
			DeviceUUID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
			PluginID:       "ILO",
		}
	case "11081de0-4859-984c-c35a-6c50732d72da":
		systemData = SavedSystems{
			ManagerAddress: "10.10.1.3",
			Password:       []byte("Password"),
			UserName:       "admin",
			DeviceUUID:     "11081de0-4859-984c-c35a-6c50732d72da",
			PluginID:       "ILO",
		}
	case "d72dade0-c35a-984c-4859-1108132d72da":
		systemData = SavedSystems{
			ManagerAddress: "odim.system.com",
			Password:       []byte("Password"),
			UserName:       "admin",
			DeviceUUID:     "d72dade0-c35a-984c-4859-1108132d72da",
			PluginID:       "GRF",
		}
	default:
		return "", fmt.Errorf("No Data found for the id")
	}
	marshalData, _ := json.Marshal(systemData)
	return string(marshalData), nil
}

type StartUpInteraface struct {
	DecryptPassword func([]byte) ([]byte, error)
	EMBConsume      func(string)
	GetAllPlugins   func() ([]agmodel.Plugin, *errors.Error)
	GetAllSystems   func() ([]string, error)
	GetSingleSystem func(string) (string, error)
	GetPluginData   func(string) (*agmodel.Plugin, *errors.Error)
}

// MockGetAllSystems is for mocking up of get all system info
func MockGetAllSystems() ([]string, error) {
	return []string{
		"6d4a0a66-7efa-578e-83cf-44dc68d2874e",
		"11081de0-4859-984c-c35a-6c50732d72da",
		"d72dade0-c35a-984c-4859-1108132d72da",
	}, nil
}

func MockInvalidGetAllSystems() ([]string, error) {
	return []string{}, &errors.Error{}
}

func TestCreateContext(t *testing.T) {
	ctx := CreateContext("123", "001", "TestAction", "0", "Test-svc-aggregation", "TestCreateContext")
	assert.Equal(t, ctx.Value("transactionid"), "123", "Context id  is not the same")
	assert.Equal(t, ctx.Value("actionid"), "001", "Context  actionId is not the same")
	assert.Equal(t, ctx.Value("actionname"), "TestAction", "Context actionName is not the same")
	assert.Equal(t, ctx.Value("threadid"), "0", "Context threadId is not the same")
	assert.Equal(t, ctx.Value("threadname"), "Test-svc-aggregation", "Context threadName is not the same")
	assert.Equal(t, ctx.Value("processname"), "TestCreateContext", "Context processName is not the same")
}
