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

// Package evmodel have the struct models and DB functionalties
package evmodel

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-persistence-manager/persistencemgr"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/stretchr/testify/assert"
)

func mockSystemResourceData(body []byte, table, key string) error {
	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		return err
	}
	if err = connPool.Create(table, key, string(body)); err != nil {
		return err
	}
	return nil
}

func mockTarget(t *testing.T) {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	target := &Target{
		ManagerAddress: "10.10.0.14",
		Password:       []byte("Password"),
		UserName:       "admin",
		DeviceUUID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
		PluginID:       "GRF",
	}
	const table string = "System"
	//Save data into Database
	if err = connPool.Create(table, target.DeviceUUID, target); err != nil {
		t.Fatalf("error: %v", err)
	}
}

func mockPlugins(t *testing.T) {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		t.Errorf("error while trying to connecting to DB: %v", err.Error())
	}

	password := getEncryptedKey(t, []byte("Password"))
	pluginArr := []Plugin{
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

func mockFabricData(t *testing.T, fabuuid, pluginID string) {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		t.Errorf("error while trying to connecting to DB: %v", err.Error())
	}

	fab := &Fabric{
		FabricUUID: fabuuid,
		PluginID:   pluginID,
	}
	const table string = "Fabric"
	//Save data into Database
	if err = connPool.Create(table, fabuuid, fab); err != nil {
		t.Errorf("error while trying to create: %v", err.Error())
	}
}

func TestGetTarget(t *testing.T) {
	config.SetUpMockConfig(t)
	target := &Target{
		ManagerAddress: "10.10.0.14",
		Password:       []byte("Password"),
		UserName:       "admin",
		DeviceUUID:     "1e61aeb6-0f03-4a35-b266-9c98e08da111",
		PluginID:       "GRF",
	}
	create(target)
	resp, err := GetTarget(target.DeviceUUID)
	if err != nil {
		t.Fatalf("Failed to get the device details")
	}
	assert.Equal(t, resp.ManagerAddress, target.ManagerAddress, "should be same")
	assert.Equal(t, resp.UserName, target.UserName, "should be same")
	assert.Equal(t, resp.PluginID, target.PluginID, "should be same")

	// Negative Test case
	// Invalid device uuid
	resp, err = GetTarget("uuid")
	assert.NotNil(t, err, "Error Should not be nil")
	assert.Nil(t, resp, "resp Should not nil")
}

func create(target *Target) *errors.Error {

	conn, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	//Create a header for data entry
	const table string = "System"
	//Save data into Database
	if err = conn.Create(table, target.DeviceUUID, target); err != nil {
		return err
	}
	return nil
}

func TestGetResource(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	var reqData = `{"@odata.id":"/redfish/v1/Systems/1"}`
	table := "ComputerSystem"
	key := "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"
	mockSystemResourceData([]byte(reqData), table, key)

	resp, err := GetResource(table, key)
	assert.Nil(t, err, "Error Should be nil")
	assert.Equal(t, reqData, resp, "response should be same as reqData")

	resp, err = GetResource("table", key)
	assert.NotNil(t, err, "Error Shouldn't be nil")
}

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
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

func TestGetPluginData(t *testing.T) {
	config.SetUpMockConfig(t)

	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()

	validPassword := []byte("password")
	invalidPassword := []byte("invalid")
	validPasswordEnc := getEncryptedKey(t, []byte("password"))

	pluginData := Plugin{
		IP:                "localhost",
		Port:              "45001",
		Username:          "admin",
		Password:          validPasswordEnc,
		ID:                "GRF",
		PluginType:        "RF-GENERIC",
		PreferredAuthType: "BasicAuth",
	}
	mockData(t, common.OnDisk, "Plugin", "validPlugin", pluginData)
	pluginData.Password = invalidPassword
	mockData(t, common.OnDisk, "Plugin", "invalidPassword", pluginData)
	mockData(t, common.OnDisk, "Plugin", "invalidPluginData", "pluginData")

	type args struct {
		pluginID string
	}
	tests := []struct {
		name    string
		args    args
		exec    func(*Plugin)
		want    *Plugin
		wantErr bool
	}{
		{
			name: "Positive Case",
			args: args{pluginID: "validPlugin"},
			exec: func(want *Plugin) {
				want.Password = validPassword
			},
			want:    &pluginData,
			wantErr: false,
		},
		{
			name:    "Negative Case - Non-existent plugin",
			args:    args{pluginID: "notFound"},
			exec:    nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Negative Case - Invalid plugin data",
			args:    args{pluginID: "invalidPluginData"},
			exec:    nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Negative Case - Plugin with invalid password",
			args:    args{pluginID: "invalidPassword"},
			exec:    nil,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if tt.exec != nil {
			tt.exec(tt.want)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPluginData(tt.args.pluginID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPluginData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPluginData() = %v, want %v", got, tt.want)
			}
		})
	}
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

	plugins, err := GetAllPlugins()
	assert.Nil(t, err, "Error Should be nil")
	assert.Equal(t, 3, len(plugins), "should be only 3 plugins")
}

func TestGetAllKeysFromTable(t *testing.T) {
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
	var reqData = `{"@odata.id":"/redfish/v1/Systems/1"}`
	table := "ComputerSystem"
	key := "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"
	mockSystemResourceData([]byte(reqData), table, key)

	resp, err := GetAllKeysFromTable(table)
	assert.Nil(t, err, "Error Should be nil")
	assert.Equal(t, 1, len(resp), "response should be same as reqData")

}

func TestGetAllSystems(t *testing.T) {
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
	mockTarget(t)

	resp, err := GetAllSystems()
	assert.Nil(t, err, "Error Should be nil")
	assert.Equal(t, 1, len(resp), "response should be same as reqData")

}

func TestGetSingleSystem(t *testing.T) {
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
	mockTarget(t)

	resp, err := GetSingleSystem("6d4a0a66-7efa-578e-83cf-44dc68d2874e")
	assert.Nil(t, err, "Error Should be nil")
	var system Target
	json.Unmarshal([]byte(resp), &system)
	assert.Equal(t, "10.10.0.14", system.ManagerAddress, "ManagerAddress should be 10.10.0.14")
	assert.Equal(t, "admin", system.UserName, "UserName should be admin")
	assert.Equal(t, "6d4a0a66-7efa-578e-83cf-44dc68d2874e", system.DeviceUUID, "DeviceUUID should be 6d4a0a66-7efa-578e-83cf-44dc68d2874e")
	assert.Equal(t, "GRF", system.PluginID, "PluginID should be GRF")

	// negative test cases
	resp, err = GetSingleSystem("578e0a66-7efa-578e-83cf-44dc68d2874e")
	assert.NotNil(t, err, "Error Should not be nil")
	//assert.Equal(t, "" , resp, "resp Should b empty)
}

func TestGetFabricData(t *testing.T) {
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

	fabuuid := "6d4a0a66-7efa-578e-83cf-44dc68d2874e"
	mockFabricData(t, fabuuid, "CFM")

	fabric, err := GetFabricData(fabuuid)
	assert.Nil(t, err, "Error Should be nil")
	assert.Equal(t, fabuuid, fabric.FabricUUID, "Fabric uuid should be same")
	assert.Equal(t, "CFM", fabric.PluginID, "PluginID should be CFM")

	// Negative Test case
	// Invalid fabric uuid
	_, err = GetFabricData("uuid")
	assert.NotNil(t, err, "Error Should not be nil")
}

func TestGetAllFabrics(t *testing.T) {
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

	fabuuid := "6d4a0a66-7efa-578e-83cf-44dc68d2874e"
	mockFabricData(t, fabuuid, "CFM")

	fabuuid = "44dc0a66-7efa-578e-83cf-44dc68d2874e"
	mockFabricData(t, fabuuid, "CFM")

	fabrics, err := GetAllFabrics()
	assert.Nil(t, err, "Error Should be nil")
	assert.Equal(t, 2, len(fabrics), "there should be 2 fabrics details")
}

func TestSaveDeviceSubscription(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	var devSubscription = DeviceSubscription{
		EventHostIP:     "10.10.0.1",
		Location:        "https://10.10.10.23/redfish/v1/EventService/Subscriptions/123",
		OriginResources: []string{"/redfish/v1/Systems/uuid.1"},
	}
	if cerr := SaveDeviceSubscription(devSubscription); cerr != nil {
		t.Errorf("Error while making save device suscription: %v\n", cerr.Error())
	}
}

func TestSaveDeviceSubscription_existing_subscription(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	var devSubscription = DeviceSubscription{
		EventHostIP:     "10.10.0.1",
		Location:        "https://10.10.10.23/redfish/v1/EventService/Subscriptions/123",
		OriginResources: []string{"/redfish/v1/Systems/uuid.1"},
	}
	if cerr := SaveDeviceSubscription(devSubscription); cerr != nil {
		t.Errorf("Error while saving device suscription: %v\n", cerr.Error())
	}

	if cerr := SaveDeviceSubscription(devSubscription); cerr == nil {
		t.Errorf("Error while saving device suscription: %v\n", cerr.Error())
	}
}

func TestGetDeviceSubscriptions(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	var devSubscription = DeviceSubscription{
		EventHostIP:     "10.10.0.1",
		Location:        "https://10.10.10.23/redfish/v1/EventService/Subscriptions/123",
		OriginResources: []string{"/redfish/v1/Systems/uuid.1"},
	}
	if cerr := SaveDeviceSubscription(devSubscription); cerr != nil {
		t.Errorf("Error while saving device suscription: %v\n", cerr.Error())
	}

	devSub, err := GetDeviceSubscriptions(devSubscription.EventHostIP)
	if err != nil {
		t.Errorf("Error while getting device suscription: %v\n", err.Error())
	}
	assert.Equal(t, devSubscription.EventHostIP, devSub.EventHostIP, "event host ip should be 10.10.0.1")
	assert.Equal(t, devSubscription.Location, devSub.Location, "Location should be https://10.10.10.23/redfish/v1/EventService/Subscriptions/123")

	if !reflect.DeepEqual(devSubscription.OriginResources, devSub.OriginResources) {
		t.Errorf("Origin Resource are not same")
	}
}

func TestDeleteDeviceSubscription(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	var devSubscription = DeviceSubscription{
		EventHostIP:     "10.10.0.1",
		Location:        "https://10.10.10.23/redfish/v1/EventService/Subscriptions/123",
		OriginResources: []string{"/redfish/v1/Systems/uuid.1"},
	}
	if cerr := SaveDeviceSubscription(devSubscription); cerr != nil {
		t.Errorf("Error while saving device suscription: %v\n", cerr.Error())
	}

	if err := DeleteDeviceSubscription(devSubscription.EventHostIP); err != nil {
		t.Errorf("Error while deleting device suscription: %v\n", err.Error())
	}
}

func TestUpdateDeviceSubscriptionLocation(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	var devSubscription = DeviceSubscription{
		EventHostIP:     "10.10.0.1",
		Location:        "https://10.10.10.23/redfish/v1/EventService/Subscriptions/123",
		OriginResources: []string{"/redfish/v1/Systems/uuid.1"},
	}
	if cerr := SaveDeviceSubscription(devSubscription); cerr != nil {
		t.Errorf("Error while saving device suscription: %v\n", cerr.Error())
	}

	devSubscription.Location = "https://10.10.10.23/redfish/v1/EventService/Subscriptions/12345"
	if err := UpdateDeviceSubscriptionLocation(devSubscription); err != nil {
		t.Errorf("Error while updating device suscription: %v\n", err.Error())
	}

	devSub, err := GetDeviceSubscriptions(devSubscription.EventHostIP)
	if err != nil {
		t.Errorf("Error while getting device suscription: %v\n", err.Error())
	}
	assert.Equal(t, devSubscription.EventHostIP, devSub.EventHostIP, "event host ip should be 10.10.0.1")
	assert.Equal(t, devSubscription.Location, devSub.Location, "Location should be https://10.10.10.23/redfish/v1/EventService/Subscriptions/123")

	if !reflect.DeepEqual(devSubscription.OriginResources, devSub.OriginResources) {
		t.Errorf("Origin Resource are not same")
	}
}

func TestSaveEventSubscription(t *testing.T) {

	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	sub := SubscriptionResource{
		SubscriptionID: "1",
		EventDestination: &model.EventDestination{
			Destination:     "https://10.10.10.23:8080/destination",
			Name:            "Event Subscription",
			EventTypes:      []string{"Alert", "StatusChange"},
			OriginResources: []string{"/redfish/v1/Systems/uuid.1"},
		},
	}
	if cerr := SaveEventSubscription(sub); cerr != nil {
		t.Errorf("Error while making save event subscriptions : %v\n", cerr.Error())
	}
}

func TestSaveEventSubscription_existingData(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	sub := SubscriptionResource{
		SubscriptionID: "123456",
		EventDestination: &model.EventDestination{
			Destination:     "https://10.10.10.23:8080/destination",
			Name:            "Event Subscription",
			EventTypes:      []string{"Alert", "StatusChange"},
			OriginResources: []string{"/redfish/v1/Systems/uuid.1"},
		},
	}
	if cerr := SaveEventSubscription(sub); cerr != nil {
		t.Errorf("Error while making save event subscriptions: %v\n", cerr.Error())
	}
	if cerr := SaveEventSubscription(sub); cerr == nil {
		t.Errorf("Error while making save event subscriptions")
	}
}

func TestGetEvtSubscriptions(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	sub := SubscriptionResource{
		SubscriptionID: "1",
		EventDestination: &model.EventDestination{
			Destination:     "https://10.10.10.23:8080/destination",
			Name:            "Event Subscription",
			EventTypes:      []string{"Alert", "StatusChange"},
			OriginResources: []string{"/redfish/v1/Systems/uuid.1"}},
	}
	if cerr := SaveEventSubscription(sub); cerr != nil {
		t.Errorf("Error while making save event subscriptions: %v\n", cerr.Error())
	}

	evtSub, err := GetEvtSubscriptions("/redfish/v1/Systems/uuid.1")
	if err != nil {
		t.Errorf("Error while getting event subscriptions: %v\n", err.Error())
	}
	assert.Equal(t, sub.SubscriptionID, evtSub[0].SubscriptionID, "SubscriptionID should be 1")
	assert.Equal(t, sub.EventDestination.Destination, evtSub[0].EventDestination.Destination, "Destination should be https://10.10.10.23:8080/destination")
	assert.Equal(t, sub.EventDestination.Name, evtSub[0].EventDestination.Name, "Name should be Event Subscription")
	assert.Equal(t, sub.EventDestination.Destination, evtSub[0].EventDestination.Destination, "Destination should be https://10.10.10.23:8080/destination")
	if !reflect.DeepEqual(sub.EventDestination.EventTypes, evtSub[0].EventDestination.EventTypes) {
		t.Errorf("Event Types are not same")
	}
	if !reflect.DeepEqual(sub.EventDestination.OriginResources, evtSub[0].EventDestination.OriginResources) {
		t.Errorf("OriginResources are not same")
	}
}

func TestDeleteEvtSubscription(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	sub := SubscriptionResource{
		SubscriptionID: "112345",
		EventDestination: &model.EventDestination{
			Destination:     "https://10.10.10.23:8080/destination",
			Name:            "Event Subscription",
			EventTypes:      []string{"Alert", "StatusChange"},
			OriginResources: []string{"/redfish/v1/Systems/uuid.1"},
		},
	}
	if cerr := SaveEventSubscription(sub); cerr != nil {
		t.Errorf("Error while making save event subscriptions: %v\n", cerr.Error())
	}

	if err := DeleteEvtSubscription(sub.SubscriptionID); err != nil {
		t.Errorf("Error while deleting event subscriptions: %v\n", err.Error())
	}
	evtSub, _ := GetEvtSubscriptions("/redfish/v1/Systems/uuid.1")
	assert.Equal(t, 0, len(evtSub), "there should be no data")
}

func TestUpdateEvtSubscription(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()

	sub := SubscriptionResource{
		SubscriptionID: "112345",
		EventDestination: &model.EventDestination{
			Destination:     "https://10.10.10.23:8080/destination",
			Name:            "Event Subscription",
			EventTypes:      []string{"Alert", "StatusChange"},
			OriginResources: []string{"/redfish/v1/Systems/uuid.1"},
		},
	}
	if cerr := SaveEventSubscription(sub); cerr != nil {
		t.Errorf("Error while making save event subscriptions: %v\n", cerr.Error())
	}

	sub.EventDestination.Destination = "https://10.10.10.23:8080/destination1"
	if err := UpdateEventSubscription(sub); err != nil {
		t.Errorf("Error while updating event subscriptions: %v\n", err.Error())
	}

	evtSub, err := GetEvtSubscriptions(sub.SubscriptionID)
	if err != nil {
		t.Errorf("Error while getting event subscriptions: %v\n", err.Error())
	}
	assert.Equal(t, "https://10.10.10.23:8080/destination1", evtSub[0].EventDestination.Destination, "Destination should be https://10.10.10.23:8080/destination1")

}
func TestSaveUndeliveredEvents(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	eventByte := []byte(`event`)
	if cerr := SaveUndeliveredEvents("destination", eventByte); cerr != nil {
		t.Errorf("Error while making save undelivered events : %v\n", cerr.Error())
	}
}

func TestGetUndeliveredEvents(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	eventByte := []byte(`event`)
	if cerr := SaveUndeliveredEvents("destination", eventByte); cerr != nil {
		t.Errorf("Error while making save undelivered events : %v\n", cerr.Error())
	}

	eventData, err := GetUndeliveredEvents("destination")
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, string(eventData), eventData, "there should be event data")
}

func TestDeleteUndeliveredEvents(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	eventByte := []byte(`event`)
	if cerr := SaveUndeliveredEvents("destination", eventByte); cerr != nil {
		t.Errorf("Error while making save undelivered events : %v\n", cerr.Error())
	}

	err := DeleteUndeliveredEvents("destination")
	assert.Nil(t, err, "error should be nil")

	_, err = GetUndeliveredEvents("destination")
	assert.NotNil(t, err, "error should not be nil")
}

func TestSetUndeliveredEventsFlag(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	if cerr := SetUndeliveredEventsFlag("destination"); cerr != nil {
		t.Errorf("Error while making set undelivered events flag: %v\n", cerr.Error())
	}

}

func TestGetUndeliveredEventsFlag(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	if cerr := SetUndeliveredEventsFlag("destination"); cerr != nil {
		t.Errorf("Error while making set undelivered events flag: %v\n", cerr.Error())
	}

	flag, err := GetUndeliveredEventsFlag("destination")
	assert.Nil(t, err, "error should be nil")
	assert.True(t, flag, "flag should be true")
}

func TestDeleteUndeliveredEventsFlag(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	if cerr := SetUndeliveredEventsFlag("destination"); cerr != nil {
		t.Errorf("Error while making set undelivered events flag: %v\n", cerr.Error())
	}
	err := DeleteUndeliveredEventsFlag("destination")
	assert.Nil(t, err, "error should be nil")

	flag, err := GetUndeliveredEventsFlag("destination")
	assert.NotNil(t, err, "error should be nil")
	assert.False(t, flag, "flag should be false")

	GetDbConnection = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) { return nil, &errors.Error{} }
	_, err = GetUndeliveredEventsFlag("destination")
	assert.NotNil(t, err, "error should be nil")
	GetDbConnection = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}

}

func TestGetAggregateList(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		hostIP string
	}
	tests := []struct {
		name            string
		args            args
		want            []string
		GetDbConnection func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error)
	}{
		{
			name: "Invalid Db Connections ",
			args: args{
				hostIP: "",
			},
			want: []string{},
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return nil, &errors.Error{}
			},
		},
		{
			name: "Valid Db Connections ",
			args: args{
				hostIP: "",
			},
			want: []string{},
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		},
		{
			name: "Valid Aggregate",
			args: args{
				hostIP: "10.10.10.10",
			},
			want: []string{"3bd1f589-117a-4cf9-89f2-da44ee8e012b"},
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Valid Aggregate" {
				mockAggregateData(t, common.OnDisk, AggregateSubscriptionIndex, "3bd1f589-117a-4cf9-89f2-da44ee8e012b", []string{"10.10.10.10"})
			}
			GetDbConnection = tt.GetDbConnection
			got, _ := GetAggregateList(tt.args.hostIP)
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetAggregateList() = %v, want %v", got, tt.want)
			}
		})
	}
}
func mockAggregateData(t *testing.T, dbType common.DbType, table, id string, data []string) {
	connPool, err := common.GetDBConnection(dbType)
	if err != nil {
		t.Fatalf("error: mockAggregateData() failed to DB connection: %v", err)
	}
	if err1 := connPool.CreateAggregateHostIndex(table, id, data); err != nil {
		t.Fatalf("error: mockAggregateData() failed to create entry %s-%s: %v", table, id, err1)
	}
}

func TestGetAggregateHosts(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		aggregateID string
	}
	tests := []struct {
		name            string
		args            args
		want            []string
		GetDbConnection func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error)
	}{
		{
			name: "Invalid Db Connections ",
			args: args{
				aggregateID: "",
			},
			want: []string{},
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return nil, &errors.Error{}
			},
		},
		{
			name: "Valid Db Connections ",
			args: args{
				aggregateID: "",
			},
			want: []string{},
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		},
		{
			name: "Valid Aggregate",
			args: args{
				aggregateID: "3bd1f589-117a-4cf9-89f2-da44ee8e012b",
			},
			want: []string{"10.10.10.10"},
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetDbConnection = tt.GetDbConnection
			if tt.name == "Valid Aggregate" {
				mockAggregateData(t, common.OnDisk, AggregateSubscriptionIndex, "3bd1f589-117a-4cf9-89f2-da44ee8e012b", []string{"10.10.10.10"})
			}
			got, _ := GetAggregateHosts(tt.args.aggregateID)
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("GetAggregateHosts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateAggregateHosts(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	mockAggregateData(t, common.OnDisk, AggregateSubscriptionIndex, "3bd1f589-117a-4cf9-89f2-da44ee8e012b", []string{"10.10.10.10"})
	type args struct {
		aggregateID string
		hostIP      []string
	}
	tests := []struct {
		name            string
		args            args
		GetDbConnection func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error)
		wantErr         bool
	}{
		{
			name: "Invalid Db Connections ",
			args: args{
				aggregateID: "",
			},
			wantErr: true,
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return nil, &errors.Error{}
			},
		},
		{
			name: "Valid Db Connections ",
			args: args{
				aggregateID: "",
			},
			wantErr: true,
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		},
		{
			name: "Valid Aggregate",
			args: args{
				aggregateID: "3bd1f589-117a-4cf9-89f2-da44ee8e012b",
				hostIP:      []string{"20.20.20"},
			},
			wantErr: false,
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetDbConnection = tt.GetDbConnection
			if err := UpdateAggregateHosts(tt.args.aggregateID, tt.args.hostIP); (err != nil) != tt.wantErr {
				t.Errorf("UpdateAggregateHosts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSaveAggregateSubscription(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		aggregateID string
		hostIP      []string
	}
	tests := []struct {
		name            string
		args            args
		GetDbConnection func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error)
		wantErr         bool
	}{
		{
			name: "Invalid Db Connections ",
			args: args{
				aggregateID: "",
			},
			wantErr: true,
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return nil, &errors.Error{}
			},
		},
		{
			name: "Valid Db Connections ",
			args: args{
				aggregateID: "",
			},
			wantErr: false,
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		},
		{
			name: "Invalid data",
			args: args{
				aggregateID: "",
				hostIP:      []string{""},
			},
			wantErr: false,
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		},
		{
			name: "Valid Aggregate",
			args: args{
				aggregateID: "3bd1f589-117a-4cf9-89f2-da44ee8e012b",
				hostIP:      []string{"20.20.20"},
			},
			wantErr: false,
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetDbConnection = tt.GetDbConnection
			if err := SaveAggregateSubscription(tt.args.aggregateID, tt.args.hostIP); (err != nil) != tt.wantErr {
				t.Errorf("SaveAggregateSubscription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAllMatchingDetails(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		table   string
		pattern string
		dbtype  common.DbType
	}
	tests := []struct {
		name            string
		args            args
		GetDbConnection func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error)
		want            []string
	}{
		{
			name: "Positive Data ",
			args: args{
				table:   AggregateSubscriptionIndex,
				pattern: "*",
				dbtype:  common.OnDisk,
			},
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		}, {
			name: "Invalid Db Connections ",
			args: args{},
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return nil, &errors.Error{}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetDbConnection = tt.GetDbConnection
			got, _ := GetAllMatchingDetails(tt.args.table, tt.args.pattern, tt.args.dbtype)
			if len(got) != len(tt.want) {
				t.Errorf("GetAllMatchingDetails() got = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestGetAggregateData(t *testing.T) {
	defer func() {
		common.TruncateDB(common.OnDisk)
		common.TruncateDB(common.InMemory)
	}()
	type args struct {
		aggreagetKey string
	}
	var aggregate = Aggregate{
		Elements: []OdataIDLink{
			OdataIDLink{OdataID: "dummy"},
		},
	}
	mockData(t, common.OnDisk, "Aggregate", "3bd1f589-117a-4cf9-89f2-da44ee8e012b", aggregate)
	mockData(t, common.OnDisk, "Aggregate", "3bd1f589-117a-4cf9-89f2-da44ee8e012c", "")
	tests := []struct {
		name            string
		args            args
		want            Aggregate
		GetDbConnection func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error)
		wantErr         bool
	}{
		{
			name:    "Invalid DB connection",
			args:    args{},
			wantErr: true,
			want: Aggregate{
				Elements: []OdataIDLink{{OdataID: ""}},
			},
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return nil, &errors.Error{}
			},
		},
		{
			name:    "invalid data ",
			args:    args{aggreagetKey: ""},
			wantErr: true,
			want: Aggregate{
				Elements: []OdataIDLink{{OdataID: ""}},
			},
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		},
		{
			name:    "Valid data",
			args:    args{aggreagetKey: "3bd1f589-117a-4cf9-89f2-da44ee8e012b"},
			wantErr: false,
			want: Aggregate{
				Elements: []OdataIDLink{{OdataID: ""}},
			},
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		},
		{
			name:    "Valid key",
			args:    args{aggreagetKey: "3bd1f589-117a-4cf9-89f2-da44ee8e012c"},
			wantErr: true,
			want: Aggregate{
				Elements: []OdataIDLink{{OdataID: ""}},
			},
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return common.GetDBConnection(dbFlag)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetDbConnection = tt.GetDbConnection
			_, err := GetAggregateData(tt.args.aggreagetKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAggregateData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
func TestInvalidDbConnection(t *testing.T) {
	GetDbConnection = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return nil, &errors.Error{}
	}
	cerr := SetUndeliveredEventsFlag("destination")
	assert.NotNil(t, cerr)
	_, cerr = GetUndeliveredEvents("destination")
	assert.NotNil(t, cerr)
	cerr = DeleteUndeliveredEvents("destination")
	assert.NotNil(t, cerr)
	_, err := GetResource("", "")
	assert.NotNil(t, "there should be an error ", err)
	_, err1 := GetTarget("")
	assert.NotNil(t, "there should be an error ", err1)
	_, err = GetPluginData("")
	assert.NotNil(t, "there should be an error ", err)
	_, err = GetAllPlugins()
	assert.NotNil(t, "there should be an error ", err)
	_, err1 = GetAllKeysFromTable("")
	assert.NotNil(t, "there should be an error ", err1)
	_, err1 = GetAllSystems()
	assert.NotNil(t, "there should be an error ", err1)
	_, err1 = GetSingleSystem("")
	assert.NotNil(t, "there should be an error ", err1)
	_, err1 = GetFabricData("")
	assert.NotNil(t, "there should be an error ", err1)
	_, err1 = GetAllFabrics()
	assert.NotNil(t, "there should be an error ", err1)
	_, err1 = GetDeviceSubscriptions("")
	assert.NotNil(t, "there should be an error ", err1)
	err1 = UpdateDeviceSubscriptionLocation(DeviceSubscription{})
	assert.NotNil(t, "there should be an error ", err1)
	err1 = SaveDeviceSubscription(DeviceSubscription{})
	assert.NotNil(t, "there should be an error ", err1)
	err1 = DeleteDeviceSubscription("")
	assert.NotNil(t, "there should be an error ", err1)
	err1 = SaveEventSubscription(SubscriptionResource{})
	assert.NotNil(t, "there should be an error ", err1)
	_, err1 = GetEvtSubscriptions("")
	assert.NotNil(t, "there should be an error ", err1)
	err1 = DeleteEvtSubscription("")
	assert.NotNil(t, "there should be an error ", err1)
	err1 = UpdateEventSubscription(SubscriptionResource{})
	assert.NotNil(t, "there should be an error ", err1)
	err1 = SaveUndeliveredEvents("", []byte{})
	assert.NotNil(t, "there should be an error ", err1)
	GetDbConnection = func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
		return common.GetDBConnection(dbFlag)
	}
}

func Test_getSliceFromString(t *testing.T) {
	type args struct {
		sliceString string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSliceFromString(tt.args.sliceString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSliceFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSliceFromString(t *testing.T) {
	type args struct {
		sliceString string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Positive Test case",
			args: args{
				sliceString: "[SystemCollection ManagerCollection]",
			},
			want: []string{"SystemCollection", "ManagerCollection"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSliceFromString(tt.args.sliceString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSliceFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllEvtSubscriptions(t *testing.T) {
	config.SetUpMockConfig(t)
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name:    "Positive Test case",
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAllEvtSubscriptions()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllEvtSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestGetAllDeviceSubscriptions(t *testing.T) {
	config.SetUpMockConfig(t)
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name:    "Positive Test case",
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAllDeviceSubscriptions()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllDeviceSubscriptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestGetAllAggregates(t *testing.T) {
	config.SetUpMockConfig(t)
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name:    "Positive Test case",
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAllAggregates()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllAggregates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestGetAggregate(t *testing.T) {
	config.SetUpMockConfig(t)
	mockAggregateList()
	type args struct {
		aggregateURI string
	}
	tests := []struct {
		name            string
		args            args
		want            Aggregate
		GetDbConnection func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error)
	}{
		{
			name: "Valid aggregate URL",
			args: args{
				aggregateURI: "/redfish/v1/AggregationService/Aggregates/b98ab95b-9187-442a-817f-b9ec60046575",
			},
			GetDbConnection: func(dbFlag common.DbType) (*persistencemgr.ConnPool, *errors.Error) {
				return nil, &errors.Error{}
			},
			want: Aggregate{Elements: []OdataIDLink{
				{
					OdataID: "/redfish/v1/Systems/e2616735-aa1f-49d9-9e03-bb1823b3100e.1",
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := GetAggregate(tt.args.aggregateURI)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAggregate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
func mockAggregateList() error {
	connPool, err := common.GetDBConnection(common.OnDisk)
	if err != nil {
		return err
	}
	aggregate := Aggregate{
		Elements: []OdataIDLink{
			{
				OdataID: "/redfish/v1/Systems/e2616735-aa1f-49d9-9e03-bb1823b3100e.1",
			},
		},
	}
	err = connPool.Create("Aggregate", "/redfish/v1/AggregationService/Aggregates/b98ab95b-9187-442a-817f-b9ec60046575", aggregate)
	if err != nil {
		return fmt.Errorf("error while trying to save Aggregate %v", err.Error())
	}
	return nil
}
