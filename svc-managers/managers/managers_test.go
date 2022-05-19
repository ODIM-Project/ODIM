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
package managers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	managersproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/managers"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrcommon"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrmodel"
	"github.com/ODIM-Project/ODIM/svc-managers/mgrresponse"
	"github.com/stretchr/testify/assert"
)

func TestGetManagersCollection(t *testing.T) {
	req := &managersproto.ManagerRequest{}
	e := mockGetExternalInterface()
	response, err := e.GetManagersCollection(req)
	assert.Nil(t, err, "There should be no error")

	manager := response.Body.(mgrresponse.ManagersCollection)
	assert.Equal(t, int(response.StatusCode), http.StatusOK, "Status code should be StatusOK.")
	assert.Equal(t, manager.MembersCount, 1, fmt.Sprintf("Managers count is expected to be 1 but got %v", manager.MembersCount))
}

func TestGetManagerRootUUIDNotFound(t *testing.T) {
	config.SetUpMockConfig(t)
	config.Data.RootServiceUUID = "nonExistingUUID"
	req := &managersproto.ManagerRequest{
		ManagerID: config.Data.RootServiceUUID,
	}
	e := mockGetExternalInterface()
	response := e.GetManagers(req)

	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound")
}

func TestGetManager(t *testing.T) {
	config.SetUpMockConfig(t)
	req := &managersproto.ManagerRequest{
		ManagerID: config.Data.RootServiceUUID,
	}
	e := mockGetExternalInterface()
	response := e.GetManagers(req)

	var manager mgrmodel.Manager
	data, _ := json.Marshal(response.Body)
	json.Unmarshal(data, &manager)

	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK")
	assert.Equal(t, "odimra", manager.Name, "Manager name should be odimra")
	assert.Equal(t, "Service", manager.ManagerType, "Manager type should be Service")
	assert.Equal(t, req.ManagerID, manager.ID, "Unexpected manager ID, should be equal to the ID in request")
	assert.Equal(t, "1.0", manager.FirmwareVersion, "Manager firmware version should be 1.0")
	assert.Equal(t, time.Now().Format(time.RFC3339), manager.DateTime, "Invalid DateTime format")
}

func TestGetManagerWithDeviceAbsent(t *testing.T) {
	req := &managersproto.ManagerRequest{
		ManagerID: "noDeviceManager.1",
		URL:       "/redfish/v1/Managers/deviceAbsent.1",
	}
	e := mockGetExternalInterface()
	response := e.GetManagers(req)

	var manager mgrmodel.Manager
	data, _ := json.Marshal(response.Body)
	json.Unmarshal(data, &manager)

	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")
	assert.Equal(t, "Absent", manager.Status.State, "Status state should be Absent.")

}

func TestGetManagerwithInvalidURL(t *testing.T) {
	req := &managersproto.ManagerRequest{
		ManagerID: "uuid.1",
		URL:       "/redfish/v1/Managers/invalidURL.1",
	}
	e := mockGetExternalInterface()
	response := e.GetManagers(req)
	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusOK.")

}

func TestGetManagerwithValidURL(t *testing.T) {
	req := &managersproto.ManagerRequest{
		ManagerID: "uuid.1",
		URL:       "/redfish/v1/Managers/uuid.1",
	}
	e := mockGetExternalInterface()
	response := e.GetManagers(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

}

func TestGetManagerLinkDetails(t *testing.T) {
	e := mockGetExternalInterface()
	var chassisLink, serverLink, managerLink []*dmtf.Link

	chassisLink = append(chassisLink, &dmtf.Link{Oid: "/redfish/v1/Managers/uuid.1"})
	serverLink = append(serverLink, &dmtf.Link{Oid: "/redfish/v1/Managers/uuid.1"})
	managerLink = append(managerLink, &dmtf.Link{Oid: "/redfish/v1/Managers/uuid.1"})
	response, _ := e.getManagerDetails("/redfish/v1/Managers/uuid.1")

	assert.Equal(t, chassisLink, response.Links.ManagerForChassis, "ManagerForChassis should be returned.")
	assert.Equal(t, serverLink, response.Links.ManagerForServers, "ManagerForServers should be returned.")
	assert.Equal(t, managerLink, response.Links.ManagerForManagers, "ManagerForManagers should be returned.")

	JsonUnMarshalFunc = func(data []byte, v interface{}) error {
		return &errors.Error{}
	}
	response, err := e.getManagerDetails("/redfish/v1/Managers/uuid.1")
	assert.NotNil(t, err, "Invalid Json .")

}

func TestGetManagerInvalidID(t *testing.T) {
	req := &managersproto.ManagerRequest{
		ManagerID: "invalidID",
	}
	e := mockGetExternalInterface()
	response := e.GetManagers(req)

	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound")
	JsonUnMarshalFunc = func(data []byte, v interface{}) error {
		return &errors.Error{}
	}
	response = e.GetManagers(req)
	assert.NotNil(t, response, "Status code should be StatusNotFound")

}

func TestGetManagerResourcewithBadManagerID(t *testing.T) {
	config.SetUpMockConfig(t)
	req := &managersproto.ManagerRequest{
		ManagerID: "invalidURL",
		URL:       "/redfish/v1/Managers/uuid",
	}
	e := mockGetExternalInterface()
	response := e.GetManagersResource(req)
	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusBadRequest.")
}

func TestGetManagerResourcewithValidURL(t *testing.T) {
	config.SetUpMockConfig(t)
	req := &managersproto.ManagerRequest{
		ManagerID: "uuid.1",
		URL:       "/redfish/v1/Managers/uuid.1/EthernetInterfaces",
	}
	e := mockGetExternalInterface()
	response := e.GetManagersResource(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid.1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid.1/EthernetInterfaces/1",
	}
	response = e.GetManagersResource(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid.1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid.1/VirtualMedia",
	}
	response = e.GetManagersResource(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid.1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid.1/VirtualMedia/1",
	}
	response = e.GetManagersResource(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")
	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid.1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid.1/LogServices",
	}
	response = e.GetManagersResource(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

}

func TestGetManagerResourcewithInvalidURL(t *testing.T) {
	config.SetUpMockConfig(t)
	req := &managersproto.ManagerRequest{
		ManagerID: "uuid1.1",
		URL:       "/redfish/v1/Managers/uuid1.1/Ethernet",
	}
	e := mockGetExternalInterface()
	response := e.GetManagersResource(req)
	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound.")

	req = &managersproto.ManagerRequest{
		ManagerID: "uuid1.1",
		URL:       "/redfish/v1/Managers/uuid1.1/Virtual",
	}
	response = e.GetManagersResource(req)
	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound.")

	req = &managersproto.ManagerRequest{
		ManagerID: "uuid1.1",
		URL:       "/redfish/v1/Managers/uuid1.1/Logservice",
	}
	response = e.GetManagersResource(req)
	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid1.1",
		ResourceID: "4",
		URL:        "/redfish/v1/Managers/uuid1.1/VirtualMedia/4",
	}
	response = e.GetManagersResource(req)
	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound.")
}

func TestGetPluginManagerResourceSuccess(t *testing.T) {
	mgrcommon.Token.Tokens = make(map[string]string)

	config.SetUpMockConfig(t)
	req := &managersproto.ManagerRequest{
		ManagerID: "uuid",
		URL:       "/redfish/v1/Managers/uuid/EthernetInterfaces",
	}
	e := mockGetExternalInterface()
	response := e.GetManagersResource(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid1/EthernetInterfaces",
	}
	response = e.GetManagersResource(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

}

func TestGetPluginManagerResourceInvalidPluginFail(t *testing.T) {
	mgrcommon.Token.Tokens = make(map[string]string)

	config.SetUpMockConfig(t)
	req := &managersproto.ManagerRequest{
		ManagerID: "noPlugin",
		URL:       "/redfish/v1/Managers/noPlugin/EthernetInterfaces",
	}
	e := mockGetExternalInterface()
	response := e.GetManagersResource(req)
	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound.")
}

func TestGetPluginManagerResourceInvalidPluginSessions(t *testing.T) {
	mgrcommon.Token.Tokens = make(map[string]string)

	config.SetUpMockConfig(t)
	req := &managersproto.ManagerRequest{
		ManagerID: "noToken",
		URL:       "/redfish/v1/Managers/uuid/EthernetInterfaces",
	}
	e := mockGetExternalInterface()
	response := e.GetManagersResource(req)
	assert.Equal(t, http.StatusUnauthorized, int(response.StatusCode), "Status code should be StatusUnauthorized.")
	mgrcommon.Token.Tokens = map[string]string{
		"CFM": "23456",
	}
	response = e.GetManagersResource(req)
	assert.Equal(t, http.StatusUnauthorized, int(response.StatusCode), "Status code should be StatusUnauthorized.")

}

func TestVirtualMediaActionsResource(t *testing.T) {
	mgrcommon.Token.Tokens = make(map[string]string)

	config.SetUpMockConfig(t)
	req := &managersproto.ManagerRequest{
		ManagerID:  "uuid.1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid.1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia",
		RequestBody: []byte(`{"Image":"http://10.1.1.1/ISO",
							"WriteProtected":true,
							"Inserted":true}`),
	}
	e := mockGetExternalInterface()
	response := e.VirtualMediaActions(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	RequestParamsCaseValidatorFunc = func(rawRequestBody []byte, reqStruct interface{}) (string, error) {
		return "test", &errors.Error{}
	}
	response = e.VirtualMediaActions(req)
	assert.Equal(t, http.StatusInternalServerError, int(response.StatusCode), "Status code should be StatusOK.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid1.1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid.1/VirtualMedia/1/Actions/VirtualMedia.EjectMedia",
	}
	response = e.VirtualMediaActions(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	req = &managersproto.ManagerRequest{
		ManagerID:   "uuid1.1",
		ResourceID:  "1",
		URL:         "/redfish/v1/Managers/uuid.1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia",
		RequestBody: []byte("dummydata"),
	}
	response = e.VirtualMediaActions(req)
	assert.Equal(t, http.StatusBadRequest, int(response.StatusCode), "Status code should be StatusBadRequest.")
	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid1.1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid.1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia",
		RequestBody: []byte(`{"image":"http://10.1.1.1/ISO",
		"writeProtected":true,
		"inserted":true}`),
	}
	RequestParamsCaseValidatorFunc = func(rawRequestBody []byte, reqStruct interface{}) (string, error) {
		return common.RequestParamsCaseValidator(rawRequestBody, reqStruct)
	}
	response = e.VirtualMediaActions(req)
	assert.Equal(t, http.StatusBadRequest, int(response.StatusCode), "Status code should be StatusBadRequest.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid1.1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid.1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia",
		RequestBody: []byte(`{
		"WriteProtected":true
		}`),
	}
	response = e.VirtualMediaActions(req)
	assert.Equal(t, http.StatusBadRequest, int(response.StatusCode), "Status code should be StatusBadRequest.")

}

func TestGetRemoteAccountService(t *testing.T) {
	mgrcommon.Token.Tokens = make(map[string]string)

	config.SetUpMockConfig(t)

	req := &managersproto.ManagerRequest{
		ManagerID: "uuid.1",
		URL:       "/redfish/v1/Managers/uuid.1/RemoteAccountService",
	}
	e := mockGetExternalInterface()
	response := e.GetRemoteAccountService(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid1.1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts/1",
	}
	response = e.GetRemoteAccountService(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid1.1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid.1/RemoteAccountService/Roles/1",
	}
	response = e.GetRemoteAccountService(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")
}

func TestCreateRemoteAccountService(t *testing.T) {
	mgrcommon.Token.Tokens = make(map[string]string)
	e := mockGetExternalInterface()
	config.SetUpMockConfig(t)
	req := &managersproto.ManagerRequest{
		ManagerID: "uuid.1",
		URL:       "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts",
		RequestBody: []byte(`{"UserName":"UserName",
                                 "Password":"Password",
                                 "RoleId":"Administrator"}`),
	}
	response := e.CreateRemoteAccountService(req)
	assert.Equal(t, http.StatusCreated, int(response.StatusCode), "Status code should be StatusCreated.")

	RequestParamsCaseValidatorFunc = func(rawRequestBody []byte, reqStruct interface{}) (string, error) {
		return "", &errors.Error{}
	}
	response = e.CreateRemoteAccountService(req)
	assert.Equal(t, http.StatusInternalServerError, int(response.StatusCode), "Status code should be StatusInternalServerError.")

	req = &managersproto.ManagerRequest{
		ManagerID: "uuid.1",
		URL:       "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts",
	}
	response = e.CreateRemoteAccountService(req)
	assert.Equal(t, http.StatusBadRequest, int(response.StatusCode), "Status code should be StatusBadRequest.")

	req = &managersproto.ManagerRequest{
		ManagerID: "uuid.1",
		URL:       "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts",
		RequestBody: []byte(`{"userName":"UserName",
                                 "password":"Password",
                                 "roleId":"Administrator"}`),
	}
	RequestParamsCaseValidatorFunc = func(rawRequestBody []byte, reqStruct interface{}) (string, error) {
		return common.RequestParamsCaseValidator(rawRequestBody, reqStruct)
	}

	response = e.CreateRemoteAccountService(req)
	assert.Equal(t, http.StatusBadRequest, int(response.StatusCode), "Status code should be StatusBadRequest.")

	req = &managersproto.ManagerRequest{
		ManagerID:   "uuid.1",
		URL:         "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts",
		RequestBody: []byte(`{}`),
	}
	response = e.CreateRemoteAccountService(req)
	assert.Equal(t, http.StatusBadRequest, int(response.StatusCode), "Status code should be StatusBadRequest.")

}

func TestDeleteRemoteAccountService(t *testing.T) {
	mgrcommon.Token.Tokens = make(map[string]string)
	e := mockGetExternalInterface()
	config.SetUpMockConfig(t)
	req := &managersproto.ManagerRequest{
		ManagerID: "uuid.1",
		URL:       "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts/5",
	}
	response := e.DeleteRemoteAccountService(req)
	assert.Equal(t, http.StatusNoContent, int(response.StatusCode), "Status code should be StatusNoContent.")
}

func TestUpdateRemoteAccountService(t *testing.T) {
	mgrcommon.Token.Tokens = make(map[string]string)
	e := mockGetExternalInterface()
	config.SetUpMockConfig(t)
	req := &managersproto.ManagerRequest{
		ManagerID: "uuid.1",
		URL:       "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts",
	}
	response := e.UpdateRemoteAccountService(req)
	assert.Equal(t, http.StatusBadRequest, int(response.StatusCode), "Status code should be StatusBadRequest")

	req = &managersproto.ManagerRequest{
		ManagerID: "uuid.1",
		URL:       "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts",
		RequestBody: []byte(`{"password":"Password",
                                 "roleId":"Administrator"}`),
	}
	response = e.UpdateRemoteAccountService(req)
	assert.Equal(t, http.StatusBadRequest, int(response.StatusCode), "Status code should be StatusBadRequest.")

	req = &managersproto.ManagerRequest{
		ManagerID: "uuid.1",
		URL:       "/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts",
		RequestBody: []byte(`{"Password":"Password",
	                           "RoleId":"Administrator"}`),
	}
	response = e.UpdateRemoteAccountService(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	RequestParamsCaseValidatorFunc = func(rawRequestBody []byte, reqStruct interface{}) (string, error) {
		return "", &errors.Error{}
	}
	response = e.UpdateRemoteAccountService(req)
	assert.Equal(t, http.StatusInternalServerError, int(response.StatusCode), "Status code should be StatusInternalServerError")

}

func TestHandleRemoteAccountServiceError(t *testing.T) {

	resp := handleRemoteAccountServiceError("/redfish/v1/Managers/uuid.1/RemoteAccountService/Accounts/1", "uuid1.1", fmt.Errorf("Dummy Error "))
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status code should be StatusNotFound.")

	resp = handleRemoteAccountServiceError("/redfish/v1/Managers/uuid.1/RemoteAccountService/Test/1", "uuid1.1", fmt.Errorf("Dummy Error "))
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status code should be StatusNotFound.")

	resp = handleRemoteAccountServiceError("/redfish/v1/Managers/uuid.1/RemoteAccountService/Roles/1", "uuid1.1", fmt.Errorf("Dummy Error "))
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status code should be StatusNotFound.")

}
