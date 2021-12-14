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

// Package system ...

package system

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

func mockPublishEventMB(systemID, eventType, collectionType string) {
	return
}

func mockGetResource(Table, key string) (string, *errors.Error) {
	return "", nil
}
func mockChassisData(systemID string) error {
	reqData, _ := json.Marshal(&map[string]interface{}{
		"Id": "1",
	})

	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("Chassis", systemID, string(reqData)); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Chassis", err.Error())
	}
	return nil
}
func mockManagerData(systemID string) error {
	reqData, _ := json.Marshal(&map[string]interface{}{
		"Id": "1",
	})

	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return fmt.Errorf("error while trying to connecting to DB: %v", err.Error())
	}
	if err = connPool.Create("Managers", systemID, string(reqData)); err != nil {
		return fmt.Errorf("error while trying to create new %v resource: %v", "Manager", err.Error())
	}
	return nil
}
func TestExternalInterface_RediscoverResources(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	device1 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "24b243cf-f1e3-5318-92d9-2d6737d6b0b9",
		PluginID:       "GRF",
	}
	device2 := agmodel.Target{
		ManagerAddress: "100.0.0.2",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "7a2c6100-67da-5fd6-ab82-6870d29c7279",
		PluginID:       "GRF",
	}
	device3 := agmodel.Target{
		ManagerAddress: "100.0.0.2",
		Password:       []byte("passwordWithInvalidEncryption"),
		UserName:       "admin",
		DeviceUUID:     "7a2c6100-67da-5fd6-ab82-6870d29c7279",
		PluginID:       "GRF",
	}
	device4 := agmodel.Target{
		ManagerAddress: "100.0.0.2",
		Password:       []byte("someValidPassword"),
		UserName:       "admin",
		DeviceUUID:     "unknown-plugin-uuid",
		PluginID:       "Unknown-Plugin",
	}
	mockPluginData(t, "GRF")
	mockDeviceData("unknown-plugin-uuid", device4)
	mockDeviceData("123443cf-f1e3-5318-92d9-2d6737d65678", device3)
	mockDeviceData("7a2c6100-67da-5fd6-ab82-6870d29c7279", device2)
	mockDeviceData("24b243cf-f1e3-5318-92d9-2d6737d6b0b9", device1)
	tests := []struct {
		name    string
		p       *ExternalInterface
		wantErr error
	}{
		{
			name:    "Positive case: All is well.",
			p:       getMockExternalInterface(),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.RediscoverResources(); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("ExternalInterface.RediscoverResources() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestExternalInterface_RediscoverSystemInventory(t *testing.T) {
	common.MuxLock.Lock()
	config.SetUpMockConfig(t)
	common.MuxLock.Unlock()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	device1 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "24b243cf-f1e3-5318-92d9-2d6737d6b0b9",
		PluginID:       "GRF",
	}
	device2 := agmodel.Target{
		ManagerAddress: "100.0.0.2",
		Password:       []byte("imKp3Q6Cx989b6JSPHnRhritEcXWtaB3zqVBkSwhCenJYfgAYBf9FlAocE"),
		UserName:       "admin",
		DeviceUUID:     "7a2c6100-67da-5fd6-ab82-6870d29c7279",
		PluginID:       "GRF",
	}
	device3 := agmodel.Target{
		ManagerAddress: "100.0.0.3",
		Password:       []byte("passwordWithInvalidEncryption"),
		UserName:       "admin",
		DeviceUUID:     "7a2c6100-67da-5fd6-ab82-6870d29c7279",
		PluginID:       "GRF",
	}
	device4 := agmodel.Target{
		ManagerAddress: "100.0.0.4",
		Password:       []byte("someValidPassword"),
		UserName:       "admin",
		DeviceUUID:     "unknown-plugin-uuid",
		PluginID:       "Unknown-Plugin",
	}
	device5 := agmodel.Target{
		ManagerAddress: "100.0.0.5",
		Password:       []byte("some-password"),
		UserName:       "admin",
		DeviceUUID:     "something",
		PluginID:       "XAuthPlugin",
	}
	mockPluginData(t, "GRF")
	mockPluginData(t, "XAuthPlugin")
	mockDeviceData("something", device5)
	mockDeviceData("unknown-plugin-uuid", device4)
	mockDeviceData("passwordWithInvalidEncryption", device3)
	mockDeviceData("7a2c6100-67da-5fd6-ab82-6870d29c7279", device2)
	mockDeviceData("24b243cf-f1e3-5318-92d9-2d6737d6b0b9", device1)
	mockSystemData("/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b9.1")

	type args struct {
		deviceUUID string
		systemURL  string
		updateFlag bool
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
	}{
		{
			name: "Negative Case: target with invalid encrypted password ",
			e:    getMockExternalInterface(),
			args: args{
				deviceUUID: "passwordWithInvalidEncryption",
				systemURL:  "/redfish/v1/Systems",
				updateFlag: true,
			},
		},
		{
			name: "Negative Case: target with non existing device UUID",
			e:    getMockExternalInterface(),
			args: args{
				deviceUUID: "nonExisting",
				systemURL:  "/redfish/v1/Systems",
				updateFlag: true,
			},
		},
		{
			name: "Negative Case: target with unknown plugin ID",
			e:    getMockExternalInterface(),
			args: args{
				deviceUUID: "unknown-plugin-uuid",
				systemURL:  "/redfish/v1/Systems/unknown-plugin-uuid.1",
				updateFlag: true,
			},
		},
		{
			name: "Positive case: All is well, Need redicovery of the resource",
			e:    getMockExternalInterface(),
			args: args{
				deviceUUID: "7a2c6100-67da-5fd6-ab82-6870d29c7279",
				systemURL:  "/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c7279.1",
				updateFlag: true,
			},
		},
		{
			name: "Positive case: All is well, Need redicovery of storage",
			e:    getMockExternalInterface(),
			args: args{
				deviceUUID: "7a2c6100-67da-5fd6-ab82-6870d29c7279",
				systemURL:  "/redfish/v1/Systems/7a2c6100-67da-5fd6-ab82-6870d29c7279.1/Storage",
				updateFlag: true,
			},
		},
		{
			name: "Positive case: All is well, dont need redicovery of the resource",
			e:    getMockExternalInterface(),
			args: args{
				deviceUUID: "24b243cf-f1e3-5318-92d9-2d6737d6b0b9",
				systemURL:  "/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b9.1",
				updateFlag: true,
			},
		},
		{
			name: "Positive case: All is well, preferred auth type XauthType",
			e:    getMockExternalInterface(),
			args: args{
				deviceUUID: "something",
				systemURL:  "/redfish/v1/Systems/something.1",
				updateFlag: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.RediscoverSystemInventory(tt.args.deviceUUID, tt.args.systemURL, tt.args.updateFlag)
		})
	}
}
func TestExternalInterface_isServerRediscoveryRequired(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		err = common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	device1 := agmodel.Target{
		ManagerAddress: "100.0.0.1",
		Password:       []byte("password"),
		UserName:       "admin",
		DeviceUUID:     "someUUID",
		PluginID:       "GRF",
	}
	device2 := agmodel.Target{
		ManagerAddress: "100.0.0.2",
		Password:       []byte("password"),
		UserName:       "admin",
		DeviceUUID:     "ComputerSystem",
		PluginID:       "GRF",
	}
	device3 := agmodel.Target{
		ManagerAddress: "100.0.0.3",
		Password:       []byte("password"),
		UserName:       "admin",
		DeviceUUID:     "Chassis&System",
		PluginID:       "GRF",
	}
	device4 := agmodel.Target{
		ManagerAddress: "100.0.0.4",
		Password:       []byte("password"),
		UserName:       "admin",
		DeviceUUID:     "Chassis&System&Manager",
		PluginID:       "GRF",
	}
	mockPluginData(t, "GRF")
	mockDeviceData("someUUID", device1)
	mockDeviceData("ComputerSystem", device2)
	mockDeviceData("Chassis&System", device3)
	mockDeviceData("Chassis&System&Manager", device4)
	mockSystemData("/redfish/v1/Systems/ComputerSystem.1")
	mockSystemData("/redfish/v1/Systems/Chassis&System.1")
	mockChassisData("/redfish/v1/Chassis/Chassis&System.1")
	mockSystemData("/redfish/v1/Systems/Chassis&System&Manager.1")
	mockChassisData("/redfish/v1/Chassis/Chassis&System&Manager.1")
	mockManagerData("/redfish/v1/Managers/Chassis&System&Manager.1")
	type args struct {
		deviceUUID string
		systemKey  string
		updateFlag bool
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
		want bool
	}{
		{
			name: "Negative case: NO Inventory in db,  Resource discovery is required.",
			e:    getMockExternalInterface(),
			args: args{
				deviceUUID: "someUUID",
				systemKey:  "/redfish/v1/Systems/1",
				updateFlag: true,
			},
			want: true,
		},
		{
			name: "Negative case: Only ComputerSystem Inventory in db,  Resource discovery is required.",
			e:    getMockExternalInterface(),
			args: args{
				deviceUUID: "ComputerSystem",
				systemKey:  "/redfish/v1/Systems/1",
				updateFlag: true,
			},
			want: true,
		},
		{
			name: "Negative case: Only Chassis and System Inventory in db,  Resource discovery is required.",
			e:    getMockExternalInterface(),
			args: args{
				deviceUUID: "Chassis&System",
				systemKey:  "/redfish/v1/Systems/1",
				updateFlag: true,
			},
			want: true,
		},
		{
			name: "Positive case: ComputerSystem, Chassis and  Manager Inventory in db,  Resource discovery not required.",
			e:    getMockExternalInterface(),
			args: args{
				deviceUUID: "Chassis&System&Manager",
				systemKey:  "/redfish/v1/Systems/1",
				updateFlag: true,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.isServerRediscoveryRequired(tt.args.deviceUUID, tt.args.systemKey); got != tt.want {
				t.Errorf("ExternalInterface.isServerRediscoveryRequired() = %v, want %v", got, tt.want)
			}
		})
	}
}
