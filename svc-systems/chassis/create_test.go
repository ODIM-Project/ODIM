//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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
	"strconv"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	chassisproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/chassis"
	"github.com/ODIM-Project/ODIM/svc-systems/plugin"
	"github.com/ODIM-Project/ODIM/svc-systems/smodel"
	"github.com/stretchr/testify/assert"
)

func TestNewCreateHandler(t *testing.T) {
	config.SetUpMockConfig(t)
	managerData := []byte(`{"@odata.context":"/redfish/v1/$metadata#Manager.Manager","@odata.etag":"W/\"6C220104\"","@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1","@odata.type":"#Manager.v1_13_0.Manager","Actions":{"#Manager.Reset":{"ResetType@Redfish.AllowableValues":["ForceRestart","GracefulRestart"],"target":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/Actions/Manager.Reset"}},"CommandShell":{"ConnectTypesSupported":["SSH","Oem"],"MaxConcurrentSessions":9,"ServiceEnabled":true},"DateTime":"2022-05-12T19:20:27Z","DateTimeLocalOffset":"+00:00","EthernetInterfaces":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/EthernetInterfaces"},"FirmwareVersion":"iLO 5 v2.60","GraphicalConsole":{"ConnectTypesSupported":["KVMIP"],"MaxConcurrentSessions":10,"ServiceEnabled":true},"HostInterfaces":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/HostInterfaces"},"Id":"1","Links":{"ManagerForChassis":[{"@odata.id":"/redfish/v1/Chassis/ba06875a-a292-445d-89ef-90e984e806ed.1"}],"ManagerForServers":[{"@odata.id":"/redfish/v1/Systems/ba06875a-a292-445d-89ef-90e984e806ed.1"}],"ManagerInChassis":{"@odata.id":"/redfish/v1/Chassis/ba06875a-a292-445d-89ef-90e984e806ed.1"}},"LogServices":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/LogServices"},"ManagerType":"BMC","Model":"iLO 5","Name":"Manager","NetworkProtocol":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/NetworkProtocol"},"Oem":{"Hpe":{"@odata.context":"/redfish/v1/$metadata#HpeiLO.HpeiLO","@odata.type":"#HpeiLO.v2_8_0.HpeiLO","Actions":{"#HpeiLO.ClearHotKeys":{"target":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/Actions/Oem/Hpe/HpeiLO.ClearHotKeys"},"#HpeiLO.ClearRestApiState":{"target":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/Actions/Oem/Hpe/HpeiLO.ClearRestApiState"},"#HpeiLO.DisableiLOFunctionality":{"target":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/Actions/Oem/Hpe/HpeiLO.DisableiLOFunctionality"},"#HpeiLO.RequestFirmwareAndOsRecovery":{"target":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/Actions/Oem/Hpe/HpeiLO.RequestFirmwareAndOsRecovery"},"#HpeiLO.ResetToFactoryDefaults":{"ResetType@Redfish.AllowableValues":["Default"],"target":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/Actions/Oem/Hpe/HpeiLO.ResetToFactoryDefaults"}},"ClearRestApiStatus":"DataPresent","ConfigurationLimitations":"None","ConfigurationSettings":"Current","FederationConfig":{"IPv6MulticastScope":"Site","MulticastAnnouncementInterval":600,"MulticastDiscovery":"Enabled","MulticastTimeToLive":5,"iLOFederationManagement":"Enabled"},"Firmware":{"Current":{"Date":"Dec 06 2021","DebugBuild":false,"MajorVersion":2,"MinorVersion":60,"VersionString":"iLO 5 v2.60"}},"FrontPanelUSB":{"State":"Ready"},"IdleConnectionTimeoutMinutes":30,"IntegratedRemoteConsole":{"HotKeys":[{"KeySequence":["NONE","NONE","NONE","NONE","NONE"],"Name":"Ctrl-T"},{"KeySequence":["NONE","NONE","NONE","NONE","NONE"],"Name":"Ctrl-U"},{"KeySequence":["NONE","NONE","NONE","NONE","NONE"],"Name":"Ctrl-V"},{"KeySequence":["NONE","NONE","NONE","NONE","NONE"],"Name":"Ctrl-W"},{"KeySequence":["NONE","NONE","NONE","NONE","NONE"],"Name":"Ctrl-X"},{"KeySequence":["NONE","NONE","NONE","NONE","NONE"],"Name":"Ctrl-Y"}],"LockKey":{"CustomKeySequence":["NONE","NONE","NONE","NONE","NONE"],"LockOption":"Disabled"},"TrustedCertificateRequired":false},"License":{"LicenseKey":"XXXXX-XXXXX-XXXXX-XXXXX-7BK6M","LicenseString":"iLO Advanced limited-distribution test","LicenseType":"Internal"},"Links":{"ActiveHealthSystem":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/ActiveHealthSystem"},"BackupRestoreService":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/BackupRestoreService"},"DateTimeService":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/DateTime"},"EmbeddedMediaService":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/EmbeddedMedia"},"FederationDispatch":{"extref":"/dispatch"},"FederationGroups":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/FederationGroups"},"FederationPeers":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/FederationPeers"},"GUIService":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/GUIService"},"LicenseService":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/LicenseService"},"RemoteSupport":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/RemoteSupportService"},"SNMPService":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/SnmpService"},"SecurityService":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/SecurityService"},"Thumbnail":{"extref":"/images/thumbnail.bmp"},"VSPLogLocation":{"extref":"/sol.log.gz"}},"PersistentMouseKeyboardEnabled":false,"PhysicalMonitorHealthStatusEnabled":true,"RIBCLEnabled":true,"RemoteConsoleThumbnailEnabled":true,"RequireHostAuthentication":false,"RequiredLoginForiLORBSU":false,"SerialCLISpeed":9600,"SerialCLIStatus":"EnabledAuthReq","SerialCLIUART":"Present","VSPDlLoggingEnabled":false,"VSPLogDownloadEnabled":false,"VideoPresenceDetectOverride":true,"VideoPresenceDetectOverrideSupported":true,"VirtualNICEnabled":false,"WebGuiEnabled":true,"iLOFunctionalityEnabled":true,"iLOFunctionalityRequired":false,"iLOIPduringPOSTEnabled":true,"iLORBSUEnabled":true,"iLOSelfTestResults":[{"Notes":"","SelfTestName":"NVRAMData","Status":"OK"},{"Notes":"Controller firmware revision  2.11.00  ","SelfTestName":"EmbeddedFlash","Status":"OK"},{"Notes":"","SelfTestName":"EEPROM","Status":"OK"},{"Notes":"","SelfTestName":"HostRom","Status":"OK"},{"Notes":"","SelfTestName":"SupportedHost","Status":"OK"},{"Notes":"Version 1.0.7","SelfTestName":"PowerManagementController","Status":"Informational"},{"Notes":"ProLiant DL360 Gen10 System Programmable Logic Device 0x2A","SelfTestName":"CPLDPAL0","Status":"Informational"},{"Notes":"","SelfTestName":"ASICFuses","Status":"OK"}],"iLOServicePort":{"MassStorageAuthenticationRequired":false,"USBEthernetAdaptersEnabled":true,"USBFlashDriveEnabled":true,"iLOServicePortEnabled":true}}},"SerialConsole":{"ConnectTypesSupported":["SSH","IPMI","Oem"],"MaxConcurrentSessions":13,"ServiceEnabled":true},"SerialInterfaces":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/SerialInterfaces"},"Status":{"Health":"OK","State":"Enabled"},"UUID":"f36bf50c-1ba9-58c1-a8ab-0409e1a4224a","VirtualMedia":{"@odata.id":"/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1/VirtualMedia"}}`)
	err := mockAddManagertoDB("Managers", "/redfish/v1/Managers/ba06875a-a292-445d-89ef-90e984e806ed.1", managerData, common.InMemory)
	if err != nil {
		t.Fatalf("Error in creating mock Manager :%v", err)
	}
	pluginData := []byte(`{"IP":"localhost","Port":"9091","Username":"admin","Password":"nROsZ0pea8KZQlgHtw644vwIxt1niZ3uPcdmdJVTd48Amh9iYEmEo6Ol3t8u1tM4HtZ7E1zRybWI+WGfh46bJ7WkDbLWqmpX4BBGYX4UwelyZh6Dij68sjvm4SRa68slkdzPJickgC5/+XCV/AGTeeT/bsgGZX39KT98xlf3BQ1hOs31OcRcLzheYO0AndkhLVeV//kaP4w8ITL6RevorFhupTTkN9iMmsinlOS158mPbGC0qnCI82gEtJfL4OzS/QfovdVfZ1ILeVEGC08ohSGdtZ1/b2V/Leu+Lg9O098n2ah8dUXNhzoSZ6QraZKUDJRecesyGYHz1kMrcjP/00eOcYpWoO8HDFydzOWCM4AbmwqkGHLTmCtfy0DqcrfETxlD6Fpbh/J09kg7QtoOneAPi3Ldyv1jhY7sqczVcJZNpotnfDlKY64vrxE9zWdVzBhIP0ncd8TFE0sVhZaVEj7x+vjx1HPq/3BxDqSywPj0F3IKztBzhVoZFGsrn79pMLG1wPsbI5lt72vQnBBvA3CQ5AGLix6EodMDeqR22UcXgNFN8KWFCL7LLT2r9Q1aWLoThI03IuQ/5jz1tBV8a0KADCy5PHCKjOg40XcZx16JJ5mGYonk0WUwvXRMzKcB1JwItYmvZ776+C4t8507TwiZBugv/4o6HA8423BxgLo=","ID":"ILO_v1.0.0","PluginType":"","PreferredAuthType":"BasicAuth","ManagerUUID":"3ccb5c71-0e00-4d14-93bb-8d125c030f27"}`)
	err = mockAddPlugonToDB("Plugin", "ILO_v1.0.0", pluginData, common.OnDisk)
	if err != nil {
		t.Fatalf("Error in creating mock Plugin :%v", err)
	}
	managerDetails := []byte(`{"@odata.context":"/redfish/v1/$metadata#Manager.Manager","@odata.etag":"WAA6D42B0","@odata.id":"/redfish/v1/Managers/3ccb5c71-0e00-4d14-93bb-8d125c030f27","@odata.type":"#Manager.v1_13_0.Manager","Certificates":{"@odata.id":""},"Description":"Plugin Manager","FirmwareVersion":"v1.0.0","Id":"3ccb5c71-0e00-4d14-93bb-8d125c030f27","Links":{"ManagerForChassis":[{"@odata.id":"/redfish/v1/Chassis/ba06875a-a292-445d-89ef-90e984e806ed.1"}],"ManagerForServers":[{"@odata.id":"/redfish/v1/Systems/ba06875a-a292-445d-89ef-90e984e806ed.1"}]},"LogServices":{"@odata.id":"/redfish/v1/Managers/3ccb5c71-0e00-4d14-93bb-8d125c030f27/LogServices"},"ManagerType":"Service","Model":"ILO v1.0.0","Name":"ILO_v1.0.0","PowerState":"On","Status":{"Health":"OK","State":"Enabled"},"UUID":"3ccb5c71-0e00-4d14-93bb-8d125c030f27"}`)
	err = mockAddManagertoDB("Managers", "/redfish/v1/Managers/3ccb5c71-0e00-4d14-93bb-8d125c030f27", managerDetails, common.InMemory)
	if err != nil {
		t.Fatalf("Error in creating mock Manager :%v", err)
	}
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

	create := new(Create)
	create.createPluginClient = func(name string) (plugin.Client, *errors.Error) {
		return &plugin.ClientMock{}, nil
	}
	// Empty Request Body
	req := chassisproto.CreateChassisRequest{}
	response := create.Handle(&req)
	assert.NotNil(t, response, "Request with empty data , Status code should be StatusInternalServerError")
	req = chassisproto.CreateChassisRequest{
		RequestBody: []byte(`{
				"ChassisType": "RackGroup",
				"Description": "My RackGroup",
				"Links": {
				  "ManagedBy": [
					{
					  "@odata.id": "/redfish/v1/Managers/3ccb5c71-0e00-4d14-93bb-8d125c030f27"
					}
				  ]
				},
				"Name": "RG5"
			  }`),
	}
	response = create.Handle(&req)
	assert.NotNil(t, response, "Status code should be StatusNotImplemented")

	// Remove Links from request
	req = chassisproto.CreateChassisRequest{
		RequestBody: []byte(`{
				"ChassisType": "RackGroup",
				"Description": "My RackGroup",
				"Links": {
				  "ManagedBy": [
					
				  ]
				},
				"Name": "RG5"
			  }`),
	}
	response = create.Handle(&req)
	assert.NotNil(t, response, "Status code should be StatusBadRequest")

	// Mocking Db Connection with error
	GetDbConnectFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	req = chassisproto.CreateChassisRequest{
		RequestBody: []byte(`{
				"ChassisType": "RackGroup",
				"Description": "My RackGroup",
				"Links": {
				  "ManagedBy": [
					{
					  "@odata.id": "/redfish/v1/Managers/3ccb5c71-0e00-4d14-93bb-8d125c030f27"
					}
				  ]
				},
				"Name": "RG5"
			  }`),
	}
	response = create.Handle(&req)
	assert.NotNil(t, response, "Can not acquire database connection")

	GetDbConnectFunc = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
	//Mocking GenericSave Func
	GenericSaveFunc = func(body []byte, table, key string) error {
		return &errors.Error{}
	}
	response = create.Handle(&req)
	assert.NotNil(t, response, "Response Status code StatusInternalServerError")

	GenericSaveFunc = func(body []byte, table, key string) error {
		return smodel.GenericSave(body, table, key)
	}
	JSONMarshalFunc = func(v interface{}) ([]byte, error) {
		return nil, &errors.Error{}
	}
	response = create.Handle(&req)
	assert.NotNil(t, response, "Response Status code StatusInternalServerError")
	JSONMarshalFunc = func(v interface{}) ([]byte, error) {
		return json.Marshal(v)
	}
	GetResourceFunc = func(Table, key string) (string, *errors.Error) {
		return "", &errors.Error{}
	}
	response = create.Handle(&req)
	assert.NotNil(t, response, "Response Status code StatusInternalServerError")
	GetResourceFunc = func(Table, key string) (string, *errors.Error) {
		return smodel.GetResource(Table, key)
	}

	JSONUnmarshalFunc1 = func(data []byte, v interface{}) error {
		return &errors.Error{}
	}
	response = create.Handle(&req)
	assert.NotNil(t, response, "Response Status code StatusInternalServerError")
	JSONUnmarshalFunc1 = func(data []byte, v interface{}) error {
		return json.Unmarshal(data, v)
	}

	JSONUnmarshalFunc2 = func(data []byte, v interface{}) error {
		return &errors.Error{}
	}
	response = create.Handle(&req)
	assert.NotNil(t, response, "Response Status code StatusInternalServerError")
	JSONUnmarshalFunc2 = func(data []byte, v interface{}) error {
		return json.Unmarshal(data, v)
	}

	JSONUnmarshalFunc3 = func(data []byte, v interface{}) error {
		return &errors.Error{}
	}
	response = create.Handle(&req)
	assert.NotNil(t, response, "Response Status code StatusInternalServerError")
	JSONUnmarshalFunc3 = func(data []byte, v interface{}) error {
		return json.Unmarshal(data, v)
	}

	JSONUnmarshalFunc4 = func(data []byte, v interface{}) error {
		return &errors.Error{}
	}
	response = create.Handle(&req)
	assert.NotNil(t, response, "Response Status code StatusInternalServerError")
	JSONUnmarshalFunc4 = func(data []byte, v interface{}) error {
		return json.Unmarshal(data, v)
	}

	StrconvUnquote = func(s string) (string, error) {
		return "", &errors.Error{}
	}
	response = create.Handle(&req)
	assert.NotNil(t, response, "Response Status code StatusInternalServerError")
	StrconvUnquote = func(s string) (string, error) {
		return strconv.Unquote(s)
	}
	//Mocking createPluginClient
	create.createPluginClient = func(name string) (plugin.Client, *errors.Error) {
		return &plugin.ClientMock{}, &errors.Error{}
	}
	response = create.Handle(&req)
	assert.NotNil(t, response, "Response Status code StatusInternalServerError")

}
func Test_NewCreateHandler(t *testing.T) {
	create := NewCreateHandler(plugin.NewClientFactory(&config.URLTranslation{}))
	assert.NotNil(t, create, "Should be no error")
}

func mockAddManagertoDB(table, key string, data []byte, dbtype common.DbType) error {
	connPool, connErr := common.GetDBConnection(dbtype)
	if connErr != nil {
		return fmt.Errorf("unable to connect DB: %v", connErr.Error())
	}
	if err := connPool.AddResourceData(table, key, string(data)); err != nil {
		return fmt.Errorf("%v", err.Error())
	}
	return nil
}
func mockAddPlugonToDB(table, key string, data []byte, dbtype common.DbType) error {

	connPool, connErr := common.GetDBConnection(dbtype)
	if connErr != nil {
		return fmt.Errorf("unable to connect DB: %v", connErr.Error())
	}
	var input interface{}

	json.Unmarshal(data, &input)

	if err := connPool.AddResourceData(table, key, input); err != nil {
		return fmt.Errorf("%v", err.Error())
	}
	return nil
}
