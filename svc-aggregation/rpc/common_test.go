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

package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
	"github.com/ODIM-Project/ODIM/svc-aggregation/system"
	uuid "github.com/satori/go.uuid"
)

var connector = &system.ExternalInterface{
	ContactClient:            mockContactClient,
	Auth:                     mockIsAuthorized,
	CreateTask:               createTaskForTesting,
	CreateChildTask:          mockCreateChildTask,
	UpdateTask:               mockUpdateTask,
	DecryptPassword:          stubDevicePassword,
	GetPluginStatus:          GetPluginStatusForTesting,
	CreateSubcription:        EventFunctionsForTesting,
	PublishEvent:             PostEventFunctionForTesting,
	EncryptPassword:          stubDevicePassword,
	DeleteComputeSystem:      deleteComputeforTest,
	DeleteSystem:             deleteSystemforTest,
	DeleteEventSubscription:  mockDeleteSubscription,
	EventNotification:        mockEventNotification,
	SubscribeToEMB:           mockSubscribeEMB,
	GetSessionUserName:       getSessionUserNameForTesting,
	GetAllKeysFromTable:      mockGetAllKeysFromTable,
	GetConnectionMethod:      mockGetConnectionMethod,
	UpdateConnectionMethod:   mockUpdateConnectionMethod,
	GetAggregationSourceInfo: mockGetAggregationSourceInfo,
	GenericSave:              mockGenericSave,
	CheckActiveRequest:       mockCheckActiveRequest,
	DeleteActiveRequest:      mockDeleteActiveRequest,
}

func mockGetAggregationSourceInfo(reqURI string) (agmodel.AggregationSource, *errors.Error) {
	var aggSource agmodel.AggregationSource
	if reqURI == "/redfish/v1/AggregationService/AggregationSources/36474ba4-a201-46aa-badf-d8104da418e8" {
		aggSource = agmodel.AggregationSource{
			HostName: "9.9.9.0",
			UserName: "admin",
			Password: []byte("admin12345"),
			Links: map[string]interface{}{
				"ConnectionMethod": map[string]interface{}{
					"OdataID": "/redfish/v1/AggregationService/ConnectionMethods/c41cbd97-937d-1b73-c41c-1b7385d3906",
				},
			},
		}
		return aggSource, nil
	}
	return aggSource, errors.PackError(errors.DBKeyNotFound, "error while trying to get compute details: no data with the with key "+reqURI+" found")
}

func mockGetAllKeysFromTable(table string) ([]string, error) {
	if table == "ConnectionMethod" {
		return []string{"/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73"}, nil
	}
	return []string{}, fmt.Errorf("Table not found")
}

func mockUpdateConnectionMethod(connectionMethod agmodel.ConnectionMethod, cmURI string) *errors.Error {
	return nil
}
func mockGetConnectionMethod(ConnectionMethodURI string) (agmodel.ConnectionMethod, *errors.Error) {
	var connMethod agmodel.ConnectionMethod
	if ConnectionMethodURI == "/redfish/v1/AggregationService/ConnectionMethods/7ff3bd97-c41c-5de0-937d-85d390691b73" {
		connMethod.ConnectionMethodType = "Redfish"
		connMethod.ConnectionMethodVariant = "iLO_v1.0.0"
		return connMethod, nil
	}
	return connMethod, errors.PackError(errors.DBKeyNotFound, "error while trying to get compute details: no data with the with key "+ConnectionMethodURI+" found")
}

func deleteComputeforTest(index int, key string) *errors.Error {
	if key == "/redfish/v1/systems/del-comp-internal-err:1" {
		return errors.PackError(errors.UndefinedErrorType, "some internal error happed")
	}
	if key != "/redfish/v1/systems/ef83e569-7336-492a-aaee-31c02d9db831:1" && key != "/redfish/v1/systems/" &&
		key != "/redfish/v1/systems/del-sys-internal-err:1" && key != "/redfish/v1/systems/sys-not-found:1" {
		return errors.PackError(errors.DBKeyNotFound, "error while trying to get compute details: no data with the with key "+key+" found")
	}
	return nil
}

func deleteSystemforTest(key string) *errors.Error {
	if key == "del-sys-internal-err" {
		return errors.PackError(errors.UndefinedErrorType, "some internal error happed")
	}
	if key != "ef83e569-7336-492a-aaee-31c02d9db831" {
		return errors.PackError(errors.DBKeyNotFound, "error while trying to get compute details: no data with the with key "+key+" found")
	}
	return nil
}

func mockDeleteSubscription(uuid string) (*eventsproto.EventSubResponse, error) {
	if uuid == "/redfish/v1/systems/delete-subscription-error:1" {
		return nil, fmt.Errorf("error while trying to delete event subcription")
	} else if uuid == "/redfish/v1/systems/unexpected-statuscode:1" {
		return &eventsproto.EventSubResponse{
			StatusCode: http.StatusCreated,
		}, nil
	}
	return &eventsproto.EventSubResponse{
		StatusCode: http.StatusNoContent,
	}, nil
}

func mockEventNotification(systemID, eventType, collectionType string) {
	return
}

func mockManagersData(id string, data map[string]interface{}) error {
	reqData, _ := json.Marshal(data)

	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("Managers", id, string(reqData)); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Managaers", err.Error())
	}
	return nil
}

func mockContactClientForDelete(url, method, token string, odataID string, body interface{}, credentials map[string]string) (*http.Response, error) {
	if url == "https://localhost:9092/ODIM/v1/Status" || (strings.Contains(url, "/ODIM/v1/Status") && credentials["UserName"] == "noStatusUser") {
		body := `{"MessageId": "` + response.Success + `"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func EventFunctionsForTesting(s []string) {}

func PostEventFunctionForTesting(s []string, name string) {}

func GetPluginStatusForTesting(plugin agmodel.Plugin) bool {
	return true
}

func mockIsAuthorized(sessionToken string, privileges, oemPrivileges []string) response.RPC {
	if sessionToken == "invalidToken" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "", nil, nil)
	}
	return common.GeneralError(http.StatusOK, response.Success, "", nil, nil)
}

func getSessionUserNameForTesting(sessionToken string) (string, error) {
	if sessionToken == "noDetailsToken" {
		return "", fmt.Errorf("no details")
	} else if sessionToken == "noTaskToken" {
		return "noTaskUser", nil
	} else if sessionToken == "taskWithSlashToken" {
		return "taskWithSlashUser", nil
	}
	return "someUserName", nil
}

func createTaskForTesting(sessionUserName string) (string, error) {
	if sessionUserName == "noTaskUser" {
		return "", fmt.Errorf("no details")
	} else if sessionUserName == "taskWithSlashUser" {
		return "some/Task/", nil
	}
	return "some/Task", nil
}

func mockSubscribeEMB(pluginID string, list []string) {
	return
}

func mockCreateChildTask(sessionID, taskID string) (string, error) {
	switch taskID {
	case "taskWithoutChild":
		return "", fmt.Errorf("subtask cannot created")
	case "subTaskWithSlash":
		return "someSubTaskID/", nil
	default:
		return "someSubTaskID", nil
	}
}

func mockSystemData(systemID string) error {
	reqData, _ := json.Marshal(&map[string]interface{}{
		"Id": "1",
	})

	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("ComputerSystem", systemID, string(reqData)); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "System", err.Error())
	}
	return nil
}

func mockUpdateTask(task common.TaskData) error {
	if task.TaskID == "invalid" {
		return fmt.Errorf(common.Cancelling)
	}
	return nil
}

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockPluginData(t *testing.T, pluginID string) error {
	password := getEncryptedKey(t, []byte("$2a$10$OgSUYvuYdI/7dLL5KkYNp.RCXISefftdj.MjbBTr6vWyNwAvht6ci"))
	plugin := agmodel.Plugin{
		IP:                "localhost",
		Port:              "9091",
		Username:          "admin",
		Password:          password,
		ID:                pluginID,
		PreferredAuthType: "BasicAuth",
	}
	switch pluginID {
	case "XAuthPlugin":
		plugin.PreferredAuthType = "XAuthToken"
	case "XAuthPluginFail":
		plugin.PreferredAuthType = "XAuthToken"
		plugin.Username = "incorrectusername"
	case "NoStatusPlugin":
		plugin.Username = "noStatusUser"
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
func mockDeviceData(uuid string, device agmodel.Target) error {

	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("System", uuid, device); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "System", err.Error())
	}
	return nil
}

func mockContactClient(url, method, token string, odataID string, body interface{}, credentials map[string]string) (*http.Response, error) {
	var bData agmodel.SaveSystem
	bBytes, _ := json.Marshal(body)
	json.Unmarshal(bBytes, &bData)
	host := strings.Split(url, "/ODIM")[0]
	uid := uuid.NewV4().String()
	if url == "https://localhost:9091/ODIM/v1/Systems/1/Actions/ComputerSystem.Reset" || url == "https://localhost:9091/ODIM/v1/Systems/1/Actions/ComputerSystem.Add" ||
		url == "https://localhost:9091/ODIM/v1/Systems/1/Actions/ComputerSystem.SetDefaultBootOrder" {
		body := `{"MessageId": "` + response.Success + `"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9091/ODIM/v1/Systems" {
		body := `{"Members":[{"@odata.id":"/ODIM/v1/Systems/1"}]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == "https://localhost:9091/ODIM/v1/Systems/1" {
		body := `{"@odata.id":"/ODIM/v1/Systems/1", "UUID": "` + uid + `", "Id": "1"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == "https://localhost:9091/ODIM/v1/Chassis" {
		body := `{"Members":[{"@odata.id":"/ODIM/v1/Chassis/1"}]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == "https://localhost:9091/ODIM/v1/Chassis/1" {
		body := `{"@odata.id":"/ODIM/v1/Chassis/1", "UUID": "` + uid + `", "Id": "1"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == host+"/ODIM/v1/Managers" {
		body := `{"Members":[{"@odata.id":"/ODIM/v1/Managers/1"}]}`
		if host == "https://100.0.0.5:9091" {
			return nil, fmt.Errorf("manager data not available not reachable")
		}
		if host == "https://100.0.0.6:9091" {
			body = "incorrectResponse"
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == host+"/ODIM/v1/Managers/1" {
		body := `{"@odata.id":"/ODIM/v1/Managers/1", "UUID": "1s7sda8asd-asdas8as0", "Id": "1"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == host+"/ODIM/v1/Status" {
		body := `{"EventMessageBus":{"EmbQueue":[{"EmbQueueName":"GRF"}]}}`
		if host == "https://100.0.0.3:9091" {
			return nil, fmt.Errorf("plugin not reachable")
		}
		if host == "https://100.0.0.4:9091" {
			body = "incorrectResponse"
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	} else if url == "https://localhost:9091/ODIM/v1/validate" || url == "https://localhost:9091/ODIM/v1/Sessions" || url == host+"/ODIM/v1/Sessions" {
		body := `{"MessageId": "` + response.Success + `"}`
		if bData.UserName == "incorrectusername" || bytes.Compare(bData.Password, []byte("incorrectPassword")) == 0 {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			}, nil
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil

	}

	return nil, fmt.Errorf("InvalidRequest")
}

func stubDevicePassword(password []byte) ([]byte, error) {
	if bytes.Compare(password, []byte("passwordWithInvalidEncryption")) == 0 {
		return []byte{}, fmt.Errorf("password decryption failed")
	}
	return password, nil
}

func TestGetAggregator(t *testing.T) {
	tests := []struct {
		name string
		want *Aggregator
	}{
		{
			name: "positive case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAggregator(); got == nil {
				t.Errorf("GetAggregator() = %v, want %v", got, tt.want)
			}
		})
	}
}

var activeReqFlag bool

func mockGenericSave(data []byte, table, key string) error {
	common.MuxLock.Lock()
	activeReqFlag = true
	common.MuxLock.Unlock()
	return nil
}

func mockCheckActiveRequest(managerAddress string) (bool, *errors.Error) {
	return activeReqFlag, nil
}

func mockDeleteActiveRequest(managerAddress string) *errors.Error {
	activeReqFlag = false
	return nil
}
