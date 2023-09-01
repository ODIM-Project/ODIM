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

// Package events have the functionality of
// - Create Event Subscription
// - Delete Event Subscription
// - Get Event Subscription
// - Post Event Subscription to destination
// and corresponding unit test cases
package events

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestPublishEventsToDestination(t *testing.T) {
	config.SetUpMockConfig(t)
	pc := getMockMethods()
	pc.LoadSubscriptionData(mockContext())

	messages := []common.MessageData{
		{
			OdataType: "#Event",
			Events: []common.Event{
				{
					MemberID:       "1",
					EventType:      "Alert",
					EventID:        "123",
					Severity:       "OK",
					EventTimestamp: "",
					Message:        "IndicatorChanged",
					MessageID:      "IndicatorChanged",
					OriginOfCondition: &common.Link{
						Oid: "/redfish/v1/Systems/1",
					},
				},
			},
		},
	}

	ip := []string{"100.100.100.100", "100.100.100.100", "10.10.1.3", "10.10.1.3"}
	mockCacheData()
	for i, v := range messages {
		var event common.Events
		event.IP = ip[i]
		message, err := json.Marshal(v)
		if err != nil {
			t.Errorf("expected err is nil but got : %v", err)
		}
		event.Request = message
		flag := pc.PublishEventsToDestination(evcommon.MockContext(), event)
		assert.True(t, flag)
	}
	for _, v := range messages {
		var event common.Events
		event.IP = "10.10.10.9"
		message, err := json.Marshal(v)
		if err != nil {
			t.Errorf("expected err is nil but got : %v", err)
		}
		event.Request = message
		flag := pc.PublishEventsToDestination(evcommon.MockContext(), event)
		assert.False(t, flag)
	}
}

func TestPublishEventsWithEmptyOriginOfCondition(t *testing.T) {
	common.SetUpMockConfig()
	pc := getMockMethods()
	pc.LoadSubscriptionData(mockContext())
	mockCacheData()
	message := common.MessageData{
		OdataType: "#Event",
		Events: []common.Event{
			{
				MemberID:          "1",
				EventType:         "Alert",
				EventID:           "123",
				Severity:          "OK",
				EventTimestamp:    "",
				Message:           "IndicatorChanged",
				MessageID:         "IndicatorChanged",
				OriginOfCondition: &common.Link{},
			},
		},
	}

	var event common.Events
	event.IP = "100.100.100.100"
	msg, err := json.Marshal(message)
	if err != nil {
		t.Errorf("expected err is nil but got : %v", err)
	}
	event.Request = msg
	flag := pc.PublishEventsToDestination(evcommon.MockContext(), event)
	assert.True(t, flag)

	// nil events
	flag = pc.PublishEventsToDestination(evcommon.MockContext(), nil)
	assert.False(t, flag)

	event.EventType = "PluginStartUp"
	flag = pc.PublishEventsToDestination(evcommon.MockContext(), event)
	assert.True(t, flag)
	event.EventType = "MetricReport"
	flag = pc.PublishEventsToDestination(evcommon.MockContext(), event)
	assert.False(t, flag)

	event.EventType = "MetricReport"
	flag = pc.PublishEventsToDestination(evcommon.MockContext(), event)
	assert.False(t, flag)

	event.EventType = "Events"
	event.Request = []byte{}
	flag = pc.PublishEventsToDestination(evcommon.MockContext(), event)
	assert.False(t, flag)

}

func TestPublishEventsToDestinationWithMultipleEvents(t *testing.T) {
	config.SetUpMockConfig(t)
	messages := []common.MessageData{
		{
			OdataType: "#Event",
			Events: []common.Event{
				{
					MemberID:       "1",
					EventType:      "Alert",
					EventID:        "123",
					Severity:       "OK",
					EventTimestamp: "",
					Message:        "IndicatorChanged",
					MessageID:      "IndicatorChanged",
					OriginOfCondition: &common.Link{
						Oid: "/redfish/v1/Systems/1",
					},
				},
				{
					MemberID:       "1",
					EventType:      "ResourceAdded",
					EventID:        "1234",
					Severity:       "OK",
					EventTimestamp: "",
					Message:        "IndicatorChanged",
					MessageID:      "IndicatorChanged",
					OriginOfCondition: &common.Link{
						Oid: "/redfish/v1/Systems/1",
					},
				},
			},
		},
	}

	ip := []string{"100.100.100.100", "100.100.100.100", "10.10.1.3", "10.10.1.3"}
	pc := getMockMethods()
	pc.LoadSubscriptionData(mockContext())
	mockCacheData()
	for i, v := range messages {
		var event common.Events
		event.IP = ip[i]
		message, err := json.Marshal(v)
		if err != nil {
			t.Errorf("expected err is nil but got : %v", err)
		}
		event.Request = message
		flag := pc.PublishEventsToDestination(evcommon.MockContext(), event)
		assert.True(t, flag)
	}
}

func TestExternalInterfaces_checkUndeliveredEvents(t *testing.T) {
	pc := getMockMethods()
	pc.checkUndeliveredEvents("dummy")
	pc.LoadSubscriptionData(mockContext())
	pc.DB.GetUndeliveredEventsFlag = func(s string) (bool, error) { return false, nil }
	pc.DB.GetAllMatchingDetails = func(s1, s2 string, dt common.DbType) ([]string, *errors.Error) {
		return []string{"/dummydestination"}, nil
	}
	SendEventFunc = func(destination string, event []byte) (*http.Response, error) {
		return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("Dummy"))}, &errors.Error{}
	}
	pc.checkUndeliveredEvents("dummy")
	SendEventFunc = func(destination string, event []byte) (*http.Response, error) {
		return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("Dummy"))}, nil
	}
	pc.checkUndeliveredEvents("dummy")
	pc.DB.GetUndeliveredEvents = func(s string) (string, error) { return "", &errors.Error{} }
	pc.checkUndeliveredEvents("dummy")

	SendEventFunc = func(destination string, event []byte) (*http.Response, error) {
		return sendEvent(destination, event)
	}
	pc = getMockMethods()
}

func Test_callPluginStartUp(t *testing.T) {
	pc := getMockMethods()
	config.SetUpMockConfig(t)
	pc.LoadSubscriptionData(mockContext())
	pc.publishMetricReport(evcommon.MockContext(), "")
	pc.DB.GetEvtSubscriptions = func(s string) ([]evmodel.SubscriptionResource, error) {
		return []evmodel.SubscriptionResource{{
			UserName:       "admin",
			SubscriptionID: "test",
			EventDestination: &model.EventDestination{
				Destination: "dummy",
			},
		}}, nil
	}
	pc.publishMetricReport(evcommon.MockContext(), "")
	JSONUnmarshal = func(data []byte, v interface{}) error {
		return &errors.Error{}
	}
	callPluginStartUp(evcommon.MockContext(), common.Events{})

	JSONUnmarshal = func(data []byte, v interface{}) error {
		return nil
	}
	updateSystemPowerState(evcommon.MockContext(), "3bd1f589-117a-4cf9-89f2-da44ee8e012b.1", "/redfish/v1/UpdateService/FirmwareInentory/3bd1f589-117a-4cf9-89f2-da44ee8e012b.1", "ServerPoweredOn")
	updateSystemPowerState(evcommon.MockContext(), "3bd1f589-117a-4cf9-89f2-da44ee8e012b.1", "/redfish/v1/UpdateService/FirmwareInentory/3bd1f589-117a-4cf9-89f2-da44ee8e012b.1", "ServerPoweredOff")
	ServiceDiscoveryFunc = func(clientName string) (*grpc.ClientConn, error) {
		return nil, &errors.Error{}
	}
	updateSystemPowerState(evcommon.MockContext(), "3bd1f589-117a-4cf9-89f2-da44ee8e012b.1", "/redfish/v1/UpdateService/FirmwareInentory/valid.1", "ServerPoweredOff")
	updateSystemPowerState(evcommon.MockContext(), "3bd1f589-117a-4cf9-89f2-da44ee8e012b.1", "/redfish/v1/UpdateService/FirmwareInentory/valid.1", "ServerPoweredOn")

	pc.removeFabricRPCCall(evcommon.MockContext(), "Zones", "test")
	pc.removeFabricRPCCall(evcommon.MockContext(), "Fabric", "test")
	pc.addFabricRPCCall(evcommon.MockContext(), "Zones", "test")
	pc.addFabricRPCCall(evcommon.MockContext(), "Fabric", "test")
	rediscoverSystemInventory(evcommon.MockContext(), "3bd1f589-117a-4cf9-89f2-da44ee8e012b.1", "/redfish/v1/UpdateService/FirmwareInentory/valid.1")

	callPluginStartUp(evcommon.MockContext(), common.Events{})
}

func TestExternalInterfaces_reAttemptEvents(t *testing.T) {
	config.SetUpMockConfig(t)
	pc := getMockMethods()
	// pc.reAttemptEvents(evcommon.MockContext(), "test", "dummy", []byte{})

	pc.DB.GetUndeliveredEvents = func(s string) (string, error) { return "test", nil }
	SendEventFunc = func(destination string, event []byte) (*http.Response, error) {
		return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("Dummy"))}, nil
	}
	pc.DB.DeleteUndeliveredEvents = func(s string) error { return &errors.Error{} }
	// pc.reAttemptEvents(evcommon.MockContext(), "test", "dummy", []byte{})
}

func TestExternalInterfaces_postEvent(t *testing.T) {

	config.SetUpMockConfig(t)
	pc := getMockMethods()

	SendEventFunc = func(destination string, event []byte) (*http.Response, error) {
		return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("Dummy"))}, nil
	}

	// pc.postEvent(evcommon.MockContext(), "dumy", "dummy", []byte{})
	pc.DB.SaveUndeliveredEvents = func(s string, b []byte) error { return &errors.Error{} }
	SendEventFunc = func(destination string, event []byte) (*http.Response, error) {
		return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("Dummy"))}, &errors.Error{}
	}

	// pc.postEvent(evcommon.MockContext(), "dumy", "dummy", []byte{})

	isStringPresentInSlice(evcommon.MockContext(), []string{"data1"}, "", "data2")
	isStringPresentInSlice(evcommon.MockContext(), []string{"data1"}, "data1", "data2")
	isStringPresentInSlice(evcommon.MockContext(), []string{"data1"}, "data3", "data2")
}

func Test_isResourceTypeSubscribed(t *testing.T) {
	// empty origin of condition
	status := isResourceTypeSubscribed(evcommon.MockContext(), []string{}, "", false)
	assert.Equal(t, true, status, "There shoud be no error ")

	resourceType := []string{"AccelerationFunction", "AddressPool", "Assembly"}
	status = isResourceTypeSubscribed(evcommon.MockContext(), resourceType, "/redfish/v1/Systems/uuid:1/processors/1", false)
	assert.Equal(t, false, status, "There shoud be no error ")

	status = isResourceTypeSubscribed(evcommon.MockContext(), resourceType, "/redfish/v1/Systems/uuid:1/processors/1", true)
	assert.Equal(t, false, status, "There shoud be no error ")

	resourceType = []string{"AccelerationFunction", "processor", "Assembly"}
	status = isResourceTypeSubscribed(evcommon.MockContext(), resourceType, "/redfish/v1/Systems/uuid:1/processors/1", false)
	assert.Equal(t, true, status, "There shoud be no error ")

}
func mockCacheData() {
	eventSourceToManagerIDMap = make(map[string]string, 2)
	eventSourceToManagerIDMap["100.100.100.100"] = "6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"
	eventSourceToManagerIDMap["10.10.1.3"] = "11081de0-4859-984c-c35a-6c50732d72da.1"
	subscriptionsCache = make(map[string]model.EventDestination, 1)
	subscriptionsCache["11081de0-4859-984c-c35a-6c50732d7"] = model.EventDestination{
		Destination: "https://10.10.10.10:8080/Destination",
	}
	aggregateIDToSubscriptionsMap = make(map[string]map[string]bool)
	emptyOriginResourceToSubscriptionsMap = make(map[string]bool, 0)
	emptyOriginResourceToSubscriptionsMap["11081de0-4859-984c-c35a-6c50732d7"] = true

	systemToSubscriptionsMap = map[string]map[string]bool{}
	systemToSubscriptionsMap["100.100.100.100"] = map[string]bool{"11081de0-4859-984c-c35a-6c50732d7": true}

}
