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
	"github.com/stretchr/testify/assert"
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
