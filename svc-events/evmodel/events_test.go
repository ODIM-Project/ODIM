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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
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
	sub := Subscription{
		SubscriptionID:  "1",
		Destination:     "https://10.10.10.23:8080/destination",
		Name:            "Event Subscription",
		EventTypes:      []string{"Alert", "StatusChange"},
		OriginResources: []string{"/redfish/v1/Systems/uuid.1"},
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
	sub := Subscription{
		SubscriptionID:  "123456",
		Destination:     "https://10.10.10.23:8080/destination",
		Name:            "Event Subscription",
		EventTypes:      []string{"Alert", "StatusChange"},
		OriginResources: []string{"/redfish/v1/Systems/uuid.1"},
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

	sub := Subscription{
		SubscriptionID:  "1",
		Destination:     "https://10.10.10.23:8080/destination",
		Name:            "Event Subscription",
		EventTypes:      []string{"Alert", "StatusChange"},
		OriginResources: []string{"/redfish/v1/Systems/uuid.1"},
	}
	if cerr := SaveEventSubscription(sub); cerr != nil {
		t.Errorf("Error while making save event subscriptions: %v\n", cerr.Error())
	}

	evtSub, err := GetEvtSubscriptions("/redfish/v1/Systems/uuid.1")
	if err != nil {
		t.Errorf("Error while getting event subscriptions: %v\n", err.Error())
	}
	assert.Equal(t, sub.SubscriptionID, evtSub[0].SubscriptionID, "SubscriptionID should be 1")
	assert.Equal(t, sub.Destination, evtSub[0].Destination, "Destination should be https://10.10.10.23:8080/destination")
	assert.Equal(t, sub.Name, evtSub[0].Name, "Name should be Event Subscription")
	assert.Equal(t, sub.Destination, evtSub[0].Destination, "Destination should be https://10.10.10.23:8080/destination")
	if !reflect.DeepEqual(sub.EventTypes, evtSub[0].EventTypes) {
		t.Errorf("Event Types are not same")
	}
	if !reflect.DeepEqual(sub.OriginResources, evtSub[0].OriginResources) {
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

	sub := Subscription{
		SubscriptionID:  "112345",
		Destination:     "https://10.10.10.23:8080/destination",
		Name:            "Event Subscription",
		EventTypes:      []string{"Alert", "StatusChange"},
		OriginResources: []string{"/redfish/v1/Systems/uuid.1"},
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

	sub := Subscription{
		SubscriptionID:  "112345",
		Destination:     "https://10.10.10.23:8080/destination",
		Name:            "Event Subscription",
		EventTypes:      []string{"Alert", "StatusChange"},
		OriginResources: []string{"/redfish/v1/Systems/uuid.1"},
	}
	if cerr := SaveEventSubscription(sub); cerr != nil {
		t.Errorf("Error while making save event subscriptions: %v\n", cerr.Error())
	}

	sub.Destination = "https://10.10.10.23:8080/destination1"
	if err := UpdateEventSubscription(sub); err != nil {
		t.Errorf("Error while updating event subscriptions: %v\n", err.Error())
	}

	evtSub, err := GetEvtSubscriptions(sub.SubscriptionID)
	if err != nil {
		t.Errorf("Error while getting event subscriptions: %v\n", err.Error())
	}
	assert.Equal(t, "https://10.10.10.23:8080/destination1", evtSub[0].Destination, "Destination should be https://10.10.10.23:8080/destination1")

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
}
