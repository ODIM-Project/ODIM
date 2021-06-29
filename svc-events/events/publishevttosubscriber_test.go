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
	"encoding/json"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/stretchr/testify/assert"
)

func storeTestEventDetails(t *testing.T) {
	subarr := []evmodel.Subscription{
		// if SubordinateResources true
		{
			UserName:             "admin",
			SubscriptionID:       "81de0110-c35a-4859-984c-072d6c5a32d7",
			Destination:          "https://10.24.1.15:9090/events",
			Name:                 "Subscription",
			Location:             "https://10.24.1.2/EventService/Subscriptions/1",
			Context:              "context",
			EventTypes:           []string{"Alert"},
			MessageIds:           []string{"IndicatorChanged"},
			ResourceTypes:        []string{"ComputerSystem"},
			OriginResources:      []string{"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"},
			Hosts:                []string{"10.4.1.2"},
			SubordinateResources: true,
		},
		{
			UserName:             "admin",
			SubscriptionID:       "11081de0-4859-984c-c35a-6c50732d72da",
			Destination:          "https://10.24.1.15:9090/events",
			Name:                 "Subscription",
			Location:             "https://10.24.1.2/EventService/Subscriptions/1",
			Context:              "context",
			EventTypes:           []string{"Alert", "StatusChange"},
			MessageIds:           []string{"IndicatorChanged", "StateChanged"},
			ResourceTypes:        []string{"ComputerSystem"},
			OriginResources:      []string{"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1", "/redfish/v1/Systems/11081de0-4859-984c-c35a-6c50732d72da:1"},
			Hosts:                []string{"10.4.1.2", "10.4.1.3"},
			SubordinateResources: true,
		},
		// if SubordinateResources false
		{
			UserName:             "admin",
			SubscriptionID:       "71de0110-c35a-4859-984c-072d6c5a32d8",
			Destination:          "https://10.24.1.16:9090/events",
			Name:                 "Subscription",
			Location:             "https://10.24.1.3/EventService/Subscriptions/1",
			Context:              "context",
			EventTypes:           []string{"Alert"},
			MessageIds:           []string{},
			ResourceTypes:        []string{},
			OriginResources:      []string{"/redfish/v1/Systems/11081de0-4859-984c-c35a-6c50732d72da:1"},
			Hosts:                []string{"10.4.1.3"},
			SubordinateResources: false,
		},
		{
			SubscriptionID:       "71de0110-c35a-4859-984c-072d6c5a32d9",
			Destination:          "https://10.24.1.16:9090/events",
			Name:                 "Subscription",
			Location:             "/ODIM/v1/Subscriptions/12345",
			Context:              "context",
			EventTypes:           []string{"Alert"},
			MessageIds:           []string{},
			ResourceTypes:        []string{},
			OriginResources:      []string{"/redfish/v1/Fabrics/123456"},
			Hosts:                []string{"10.4.1.6"},
			SubordinateResources: true,
		},
		{
			SubscriptionID:       "71de0110-c35a-4859-984c-072d6c5a3210",
			Destination:          "https://10.24.1.16:9090/events",
			Name:                 "Subscription",
			Location:             "/ODIM/v1/Subscriptions/12345",
			Context:              "context",
			EventTypes:           []string{"Alert", "ResourceAdded"},
			MessageIds:           []string{},
			ResourceTypes:        []string{},
			OriginResources:      []string{"/redfish/v1/Fabrics/123456"},
			Hosts:                []string{"10.4.1.6"},
			SubordinateResources: true,
		},
		{
			SubscriptionID:       "5a321010-c35a-4859-984c-072d6c",
			Destination:          "https://10.24.1.16:9090/events",
			Name:                 "Subscription",
			Location:             "/ODIM/v1/Subscriptions/123",
			Context:              "context",
			EventTypes:           []string{"Alert", "ResourceAdded"},
			MessageIds:           []string{},
			ResourceTypes:        []string{},
			OriginResources:      []string{"/redfish/v1/Fabrics/123"},
			Hosts:                []string{"10.4.1.5"},
			SubordinateResources: true,
		},
		{
			SubscriptionID: "71de0110-c35a-4859-984c-072d6c5a3211",
			Destination:    "https://localhost:1234/eventsListener",
			Name:           "Subscription",
			Location:       "/ODIM/v1/Subscriptions/12345",
			Context:        "context",
			EventTypes:     []string{"Alert"},
			MessageIds:     []string{},
			ResourceTypes:  []string{},
			OriginResources: []string{"/redfish/v1/Systems",
				"/redfish/v1/Chassis",
				"/redfish/v1/Fabrics",
				"/redfish/v1/Managers",
				"/redfish/v1/TaskService/Tasks"},
			Hosts:                []string{"localhost"},
			SubordinateResources: true,
		},
		{
			UserName:             "admin",
			SubscriptionID:       "81de0110-c35a-4859-984c-072d6c5a32d8",
			Destination:          "https://10.24.1.9:9090/events",
			Name:                 "Subscription",
			Location:             "https://10.24.1.2/EventService/Subscriptions/1",
			Context:              "context",
			EventTypes:           []string{"Alert"},
			MessageIds:           []string{"IndicatorChanged"},
			ResourceTypes:        []string{"ComputerSystem"},
			OriginResources:      []string{"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"},
			Hosts:                []string{"10.4.1.2"},
			SubordinateResources: true,
		},
	}

	for _, sub := range subarr {
		if cerr := evmodel.SaveEventSubscription(sub); cerr != nil {
			t.Fatalf("Error while making save event subscriptions : %v\n", cerr.Error())
		}
	}

	devSubArr := []evmodel.DeviceSubscription{
		{
			Location:        "https://10.24.1.2/EventService/Subscriptions/1",
			EventHostIP:     "10.4.1.2",
			OriginResources: []string{"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"},
		},
		{
			Location:        "https://10.24.1.3/EventService/Subscriptions/1",
			EventHostIP:     "10.4.1.3",
			OriginResources: []string{"/redfish/v1/Systems/11081de0-4859-984c-c35a-6c50732d72da:1"},
		},
		{
			Location:        "/ODIM/v1/Subscriptions/12345",
			EventHostIP:     "10.4.1.6",
			OriginResources: []string{"/redfish/v1/Fabrics/123456"},
		},
		{
			Location:        "/ODIM/v1/Subscriptions/123",
			EventHostIP:     "10.4.1.5",
			OriginResources: []string{"/redfish/v1/Fabrics/123"},
		},
		{
			Location:        "/ODIM/v1/Subscriptions/12345",
			EventHostIP:     "localhost",
			OriginResources: []string{""},
		},
		{
			Location:        "https://10.24.1.9/EventService/Subscriptions/1",
			EventHostIP:     "10.4.1.9",
			OriginResources: []string{"/redfish/v1/Systems/11081de0-4859-984c-c35a-6c50732d72ea:1"},
		},
	}
	for _, devSub := range devSubArr {
		if cerr := evmodel.SaveDeviceSubscription(devSub); cerr != nil {
			t.Fatalf("Error while making save device subscriptions : %v\n", cerr.Error())
		}
	}

}
func TestPublishEventsToDestiantion(t *testing.T) {
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
	storeTestEventDetails(t)
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

	ip := []string{"10.4.1.2", "10.4.1.2", "10.4.1.3", "10.4.1.3"}
	for i, v := range messages {
		var event common.Events
		event.IP = ip[i]
		message, err := json.Marshal(v)
		if err != nil {
			t.Errorf("expected err is nil but got : %v", err)
		}
		event.Request = message
		flag := PublishEventsToDestination(event)
		assert.True(t, flag)
	}
	for _, v := range messages {
		var event common.Events
		event.IP = "10.24.1.9"
		message, err := json.Marshal(v)
		if err != nil {
			t.Errorf("expected err is nil but got : %v", err)
		}
		event.Request = message
		flag := PublishEventsToDestination(event)
		assert.False(t, flag)
	}
}

func TestPublishEventsWithEmptyOriginOfCondition(t *testing.T) {
	common.SetUpMockConfig()
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
	event.IP = "10.4.1.2"
	msg, err := json.Marshal(message)
	if err != nil {
		t.Errorf("expected err is nil but got : %v", err)
	}
	event.Request = msg
	flag := PublishEventsToDestination(event)
	assert.False(t, flag)

}

func TestPublishEventsToDestiantionWithMultipleEvents(t *testing.T) {
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
	storeTestEventDetails(t)
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

	ip := []string{"10.4.1.2", "10.4.1.2", "10.4.1.3", "10.4.1.3"}
	for i, v := range messages {
		var event common.Events
		event.IP = ip[i]
		message, err := json.Marshal(v)
		if err != nil {
			t.Errorf("expected err is nil but got : %v", err)
		}
		event.Request = message
		flag := PublishEventsToDestination(event)
		assert.True(t, flag)
	}
}
