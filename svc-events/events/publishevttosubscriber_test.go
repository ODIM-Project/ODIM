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

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestPublishEventsToDestiantion(t *testing.T) {
	config.SetUpMockConfig(t)
	messages := []common.MessageData{
		{
			OdataType: "#Event",
			Events: []common.Event{
				common.Event{
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
	pc := getMockMethods()
	for i, v := range messages {
		var event common.Events
		event.IP = ip[i]
		message, err := json.Marshal(v)
		if err != nil {
			t.Errorf("expected err is nil but got : %v", err)
		}
		event.Request = message
		flag := pc.PublishEventsToDestination(event)
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
		flag := pc.PublishEventsToDestination(event)
		assert.False(t, flag)
	}
}

func TestPublishEventsWithEmptyOriginOfCondition(t *testing.T) {
	common.SetUpMockConfig()
	message := common.MessageData{
		OdataType: "#Event",
		Events: []common.Event{
			common.Event{
				MemberID:       "1",
				EventType:      "Alert",
				EventID:        "123",
				Severity:       "OK",
				EventTimestamp: "",
				Message:        "IndicatorChanged",
				MessageID:      "IndicatorChanged",
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
	pc := getMockMethods()
	flag := pc.PublishEventsToDestination(event)
	assert.False(t, flag)

	// nil events
	flag = pc.PublishEventsToDestination(nil)
	assert.False(t, flag)

	event.EventType = "PluginStartUp"
	flag = pc.PublishEventsToDestination(event)
	assert.True(t, flag)
	event.EventType = "MetricReport"
	flag = pc.PublishEventsToDestination(event)
	assert.False(t, flag)

	event.EventType = "MetricReport"
	flag = pc.PublishEventsToDestination(event)
	assert.False(t, flag)

	event.EventType = "Events"
	event.Request = []byte{}
	flag = pc.PublishEventsToDestination(event)
	assert.False(t, flag)

}

func TestPublishEventsToDestiantionWithMultipleEvents(t *testing.T) {
	config.SetUpMockConfig(t)
	messages := []common.MessageData{
		{
			OdataType: "#Event",
			Events: []common.Event{
				common.Event{
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
				common.Event{
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
	for i, v := range messages {
		var event common.Events
		event.IP = ip[i]
		message, err := json.Marshal(v)
		if err != nil {
			t.Errorf("expected err is nil but got : %v", err)
		}
		event.Request = message
		flag := pc.PublishEventsToDestination(event)
		assert.True(t, flag)
	}
}

func TestExternalInterfaces_getCollectionSubscriptionInfoForOID(t *testing.T) {
	pc := getMockMethods()
	pc.getCollectionSubscriptionInfoForOID("Chassis", "")
	pc.getCollectionSubscriptionInfoForOID("Managers", "")
	pc.getCollectionSubscriptionInfoForOID("Fabrics", "")
	pc.getCollectionSubscriptionInfoForOID("", "")
}

func TestExternalInterfaces_checkUndeliveredEvents(t *testing.T) {
	pc := getMockMethods()
	pc.checkUndeliveredEvents("dummy")
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
	pc.publishMetricReport("")
	pc.DB.GetEvtSubscriptions = func(s string) ([]evmodel.Subscription, error) {
		return []evmodel.Subscription{{
			UserName:       "admin",
			SubscriptionID: "test",
			Destination:    "dummy",
		}}, nil
	}
	pc.publishMetricReport("")
	JSONUnmarshal = func(data []byte, v interface{}) error {
		return &errors.Error{}
	}
	callPluginStartUp(common.Events{})

	JSONUnmarshal = func(data []byte, v interface{}) error {
		return nil
	}
	updateSystemPowerState("3bd1f589-117a-4cf9-89f2-da44ee8e012b.1", "/redfish/v1/UpdateService/FirmwareInentory/3bd1f589-117a-4cf9-89f2-da44ee8e012b.1", "ServerPoweredOn")
	updateSystemPowerState("3bd1f589-117a-4cf9-89f2-da44ee8e012b.1", "/redfish/v1/UpdateService/FirmwareInentory/3bd1f589-117a-4cf9-89f2-da44ee8e012b.1", "ServerPoweredOff")
	ServiceDiscoveryFunc = func(clientName string) (*grpc.ClientConn, error) {
		return nil, &errors.Error{}
	}
	updateSystemPowerState("3bd1f589-117a-4cf9-89f2-da44ee8e012b.1", "/redfish/v1/UpdateService/FirmwareInentory/valid.1", "ServerPoweredOff")
	updateSystemPowerState("3bd1f589-117a-4cf9-89f2-da44ee8e012b.1", "/redfish/v1/UpdateService/FirmwareInentory/valid.1", "ServerPoweredOn")

	pc.removeFabricRPCCall("Zones", "test")
	pc.removeFabricRPCCall("Fabric", "test")
	pc.addFabricRPCCall("Zones", "test")
	pc.addFabricRPCCall("Fabric", "test")
	rediscoverSystemInventory("3bd1f589-117a-4cf9-89f2-da44ee8e012b.1", "/redfish/v1/UpdateService/FirmwareInentory/valid.1")

	callPluginStartUp(common.Events{})
}

func TestExternalInterfaces_reAttemptEvents(t *testing.T) {
	config.SetUpMockConfig(t)
	pc := getMockMethods()
	pc.reAttemptEvents("test", "dummy", []byte{})

	pc.DB.GetUndeliveredEvents = func(s string) (string, error) { return "test", nil }
	SendEventFunc = func(destination string, event []byte) (*http.Response, error) {
		return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("Dummy"))}, nil
	}
	pc.DB.DeleteUndeliveredEvents = func(s string) error { return &errors.Error{} }
	pc.reAttemptEvents("test", "dummy", []byte{})
}

func TestExternalInterfaces_postEvent(t *testing.T) {

	config.SetUpMockConfig(t)
	pc := getMockMethods()

	SendEventFunc = func(destination string, event []byte) (*http.Response, error) {
		return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("Dummy"))}, nil
	}

	pc.postEvent("dumy", "dummy", []byte{})
	pc.DB.SaveUndeliveredEvents = func(s string, b []byte) error { return &errors.Error{} }
	SendEventFunc = func(destination string, event []byte) (*http.Response, error) {
		return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString("Dummy"))}, &errors.Error{}
	}

	pc.postEvent("dumy", "dummy", []byte{})

	isStringPresentInSlice([]string{"data1"}, "", "data2")
	isStringPresentInSlice([]string{"data1"}, "data1", "data2")
	isStringPresentInSlice([]string{"data1"}, "data3", "data2")
	pc.addFabric(common.MessageData{}, "dummy")

}

func Test_isResourceTypeSubscribed(t *testing.T) {
	// empty origin of condition
	status := isResourceTypeSubscribed([]string{}, "", false)
	assert.Equal(t, true, status, "There shoud be no error ")

	resourceType := []string{"AccelerationFunction", "AddressPool", "Assembly"}
	status = isResourceTypeSubscribed(resourceType, "/redfish/v1/Systems/uuid:1/processors/1", false)
	assert.Equal(t, false, status, "There shoud be no error ")

	status = isResourceTypeSubscribed(resourceType, "/redfish/v1/Systems/uuid:1/processors/1", true)
	assert.Equal(t, false, status, "There shoud be no error ")

	resourceType = []string{"AccelerationFunction", "processor", "Assembly"}
	status = isResourceTypeSubscribed(resourceType, "/redfish/v1/Systems/uuid:1/processors/1", false)
	assert.Equal(t, true, status, "There shoud be no error ")

}
