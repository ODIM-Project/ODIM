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

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
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

	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")
	assert.Equal(t, "odimra", manager.Name, "Status code should be StatusOK.")
	assert.Equal(t, "Service", manager.ManagerType, "Status code should be StatusOK.")
	assert.Equal(t, req.ManagerID, manager.ID, "Status code should be StatusOK.")
	assert.Equal(t, "1.0", manager.FirmwareVersion, "Status code should be StatusOK.")

}

func TestGetManagerWithDeviceAbsent(t *testing.T) {
	req := &managersproto.ManagerRequest{
		ManagerID: "noDeviceManager:1",
		URL:       "/redfish/v1/Managers/deviceAbsent:1",
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
		ManagerID: "uuid:1",
		URL:       "/redfish/v1/Managers/invalidURL:1",
	}
	e := mockGetExternalInterface()
	response := e.GetManagers(req)
	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusOK.")

}

func TestGetManagerwithValidURL(t *testing.T) {
	req := &managersproto.ManagerRequest{
		ManagerID: "uuid:1",
		URL:       "/redfish/v1/Managers/uuid:1",
	}
	e := mockGetExternalInterface()
	response := e.GetManagers(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

}

func TestGetManagerInvalidID(t *testing.T) {
	req := &managersproto.ManagerRequest{
		ManagerID: "invalidID",
	}
	e := mockGetExternalInterface()
	response := e.GetManagers(req)

	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound")
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
		ManagerID: "uuid:1",
		URL:       "/redfish/v1/Managers/uuid:1/EthernetInterfaces",
	}
	e := mockGetExternalInterface()
	response := e.GetManagersResource(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid:1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid:1/EthernetInterfaces/1",
	}
	response = e.GetManagersResource(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid:1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid:1/VirtualMedia",
	}
	response = e.GetManagersResource(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid:1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid:1/VirtualMedia/1",
	}
	response = e.GetManagersResource(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")
}

func TestGetManagerResourcewithInvalidURL(t *testing.T) {
	config.SetUpMockConfig(t)
	req := &managersproto.ManagerRequest{
		ManagerID: "uuid1:1",
		URL:       "/redfish/v1/Managers/uuid1:1/Ethernet",
	}
	e := mockGetExternalInterface()
	response := e.GetManagersResource(req)
	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound.")

	req = &managersproto.ManagerRequest{
		ManagerID: "uuid1:1",
		URL:       "/redfish/v1/Managers/uuid1:1/Virtual",
	}
	response = e.GetManagersResource(req)
	assert.Equal(t, http.StatusNotFound, int(response.StatusCode), "Status code should be StatusNotFound.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid1:1",
		ResourceID: "4",
		URL:        "/redfish/v1/Managers/uuid1:1/VirtualMedia/4",
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
		ManagerID:  "uuid:1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid:1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia",
		RequestBody: []byte(`{"Image":"http://10.1.1.1/ISO",
							"WriteProtected":true,
							"Inserted":true}`),
	}
	e := mockGetExternalInterface()
	response := e.VirtualMediaActions(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")

	req = &managersproto.ManagerRequest{
		ManagerID:  "uuid1:1",
		ResourceID: "1",
		URL:        "/redfish/v1/Managers/uuid:1/VirtualMedia/1/Actions/VirtualMedia.EjectMedia",
	}
	response = e.VirtualMediaActions(req)
	assert.Equal(t, http.StatusOK, int(response.StatusCode), "Status code should be StatusOK.")
}
