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
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/stretchr/testify/assert"
)

// Positive test cases
func TestCreateEventSubscription(t *testing.T) {
	config.SetUpMockConfig(t)
	p := getMockMethods()
	taskID := "123"
	sessionUserName := "admin"
	SubscriptionReq := map[string]interface{}{
		"Name":                 "EventSubscription",
		"Destination":          "https://localhost:8070/Destination1",
		"EventTypes":           []string{"Alert"},
		"Protocol":             "Redfish",
		"Context":              "Event Subscription",
		"SubscriptionType":     "RedfishEvent",
		"EventFormatType":      "Event",
		"SubordinateResources": true,
		"OriginResources": []evmodel.OdataIDLink{
			{OdataID: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"},
			{OdataID: "/redfish/v1/Systems/11081de0-4859-984c-c35a-6c50732d72da.1"},
		},
	}
	postBody, _ := json.Marshal(&SubscriptionReq)

	// Positive test cases
	req := &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp := p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusCreated, int(resp.StatusCode), "Status Code should be StatusCreated")

	// try to subscrie with already subscribed destinations
	SubscriptionReq["Destination"] = "https://odim.tt.com:9090/events"

	postBody, _ = json.Marshal(&SubscriptionReq)

	req = &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusConflict, int(resp.StatusCode), "Status Code should be StatusCreated")

	// test with different Destinations
	SubscriptionReq["Destination"] = "https://10.10.10.25:8070/Destination2"

	postBody, _ = json.Marshal(&SubscriptionReq)

	req = &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusCreated, int(resp.StatusCode), "Status Code should be StatusCreated")

	// test another subscription with other OriginResources
	SubscriptionReq["OriginResources"] = []evmodel.OdataIDLink{
		{OdataID: "/redfish/v1/Systems/11081de0-4859-984c-c35a-6c50732d72da.1"},
	}
	SubscriptionReq["Destination"] = "https://10.10.10.25:8070/Destination3"
	SubscriptionReq["EventTypes"] = []string{"Alert"}
	postBody, _ = json.Marshal(&SubscriptionReq)

	req = &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusCreated, int(resp.StatusCode), "Status Code should be StatusCreated")

	// test case for collection
	SubscriptionReq["OriginResources"] = []evmodel.OdataIDLink{}
	SubscriptionReq["Destination"] = "https://10.10.10.25:8070/Destination4"
	postBody, _ = json.Marshal(&SubscriptionReq)

	req = &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusCreated, int(resp.StatusCode), "Status Code should be StatusCreated")
}

func TestCreateEventSubscriptionwithHostName(t *testing.T) {
	config.SetUpMockConfig(t)
	p := getMockMethods()
	taskID := "123"
	sessionUserName := "admin"
	SubscriptionReq := map[string]interface{}{
		"Name":                 "EventSubscription",
		"Destination":          "https://10.10.10.24:8070/Destination1",
		"EventTypes":           []string{"Alert"},
		"Protocol":             "Redfish",
		"Context":              "Event Subscription",
		"SubscriptionType":     "RedfishEvent",
		"EventFormatType":      "Event",
		"SubordinateResources": true,
		"OriginResources": []evmodel.OdataIDLink{
			{
				OdataID: "/redfish/v1/Systems/abab09db-e7a9-4352-8df0-5e41315a2a4c.1",
			},
		},
	}
	postBody, _ := json.Marshal(&SubscriptionReq)

	// Positive test cases
	req := &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp := p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusCreated, int(resp.StatusCode), "Status Code should be StatusCreated")

}

// Negative test cases
func TestNegativeCasesCreateEventSubscription(t *testing.T) {
	config.SetUpMockConfig(t)
	p := getMockMethods()
	taskID := "123"
	sessionUserName := "admin"
	SubscriptionReq := map[string]interface{}{
		"Name":                 "EventSubscription",
		"Destination":          "https://10.10.10.24:8070/Destination1",
		"EventTypes":           []string{"Alert"},
		"Protocol":             "Redfish",
		"Context":              "Event Subscription",
		"SubscriptionType":     "RedfishEvent",
		"EventFormatType":      "Event",
		"SubordinateResources": true,
		"OriginResources": []evmodel.OdataIDLink{
			{OdataID: "/redfish/v1/Systems/d72dade0-c35a-984c-4859-1108132d72da.1"},
		},
	}

	postBody, _ := json.Marshal(&SubscriptionReq)

	// Bad Request from the plugin
	req := &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp := p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// invalid post body
	req1 := &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     []byte(""),
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req1)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// if Destination is empty
	SubscriptionReq["Destination"] = ""
	postBody, _ = json.Marshal(&SubscriptionReq)

	req2 := &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req2)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// if Destination is invalid
	SubscriptionReq["Destination"] = "destination"
	postBody, _ = json.Marshal(&SubscriptionReq)

	req2 = &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req2)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")
	SubscriptionReq["Destination"] = "https://10.10.10.24:8070/Destination1"

	// if Protocol is empty
	SubscriptionReq["Protocol"] = ""
	postBody, _ = json.Marshal(&SubscriptionReq)

	req = &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// if Protocol is invalid
	SubscriptionReq["Protocol"] = "Protocol"
	postBody, _ = json.Marshal(&SubscriptionReq)

	req = &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")
	SubscriptionReq["Protocol"] = "Redfish"

	// if EventFormatType is Unspported
	SubscriptionReq["EventFormatType"] = "MetricReport"
	postBody, _ = json.Marshal(&SubscriptionReq)

	req = &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// if EventFormatType is invalid
	SubscriptionReq["EventFormatType"] = "EventFormatType"
	postBody, _ = json.Marshal(&SubscriptionReq)

	req = &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")
	SubscriptionReq["EventFormatType"] = "Event"

	// if SubscriptionType is Unsupported
	SubscriptionReq["SubscriptionType"] = "SSE"
	postBody, _ = json.Marshal(&SubscriptionReq)

	req = &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// if SubscriptionType is invalid
	SubscriptionReq["SubscriptionType"] = "SubscriptionType"
	postBody, _ = json.Marshal(&SubscriptionReq)

	req = &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")
	SubscriptionReq["SubscriptionType"] = "RedfishEvent"

	// if ResourceType is unsupported
	SubscriptionReq["ResourceTypes"] = []string{"InvalidResourceType"}
	postBody, _ = json.Marshal(&SubscriptionReq)

	req = &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")
	SubscriptionReq["ResourceTypes"] = []string{}

	postBody, _ = json.Marshal(&SubscriptionReq)

	req = &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")
}

func TestCreateDefaultEventSubscription(t *testing.T) {
	config.SetUpMockConfig(t)
	p := getMockMethods()
	taskID := "123"
	sessionUserName := "admin"
	SubscriptionReq := map[string]interface{}{
		"Name":                 "EventSubscription",
		"Destination":          "https://odim.test24.com:8070/Destination1",
		"EventTypes":           []string{"Alert"},
		"Protocol":             "Redfish",
		"Context":              "Event Subscription",
		"SubscriptionType":     "RedfishEvent",
		"EventFormatType":      "Event",
		"SubordinateResources": true,
		"OriginResources": []evmodel.OdataIDLink{
			{OdataID: "/redfish/v1/Systems"},
		},
	}

	postBody, _ := json.Marshal(&SubscriptionReq)

	// Positive test cases
	req := &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	p.CreateEventSubscription(taskID, sessionUserName, req)

	systemURL := []string{"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"}
	eventTypes := []string{"Alert"}
	messageIDs := []string{}
	resourceTypes := []string{}
	protocol := "redfish"

	resp := p.CreateDefaultEventSubscription(systemURL, eventTypes, messageIDs, resourceTypes, protocol)
	assert.Equal(t, http.StatusCreated, int(resp.StatusCode), "Status Code should be StatusCreated")

}

func TestFabricEventSubscription(t *testing.T) {
	if config.Data.URLTranslation == nil {
		config.SetUpMockConfig(t)
	}

	p := getMockMethods()
	taskID := "123"
	sessionUserName := "admin"
	SubscriptionReq := map[string]interface{}{
		"Name":                 "EventSubscription",
		"Destination":          "https://odim.test24.com:8070/Destination1",
		"EventTypes":           []string{"Alert"},
		"Protocol":             "Redfish",
		"Context":              "Event Subscription",
		"SubscriptionType":     "RedfishEvent",
		"EventFormatType":      "Event",
		"SubordinateResources": true,
		"OriginResources": []evmodel.OdataIDLink{
			{OdataID: "/redfish/v1/Fabrics/6d4a0a66-7efa-578e-83cf-44dc68d2874e"},
		},
	}

	postBody, _ := json.Marshal(&SubscriptionReq)

	// Positive test cases
	req := &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp := p.CreateEventSubscription(taskID, sessionUserName, req)
	assert.Equal(t, http.StatusCreated, int(resp.StatusCode), "Status Code should be StatusCreated")

	// Negative test cases

	// Invalid Fabric ID
	SubscriptionReq["OriginResources"] = []evmodel.OdataIDLink{
		{OdataID: "/redfish/v1/Fabrics/11081de0-4859-984c-c35a-6c50732d72da"},
	}
	SubscriptionReq["Destination"] = "https://10.10.10.24:8070/Destination2"
	postBody, _ = json.Marshal(&SubscriptionReq)
	req1 := &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}
	resp = p.CreateEventSubscription(taskID, sessionUserName, req1)
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status Code should be StatusCreated")

	// Invalid Plugin ID
	resp = p.CreateEventSubscription(taskID, sessionUserName, req1)
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status Code should be StatusCreated")

	// Unauthorized token
	SubscriptionReq["OriginResources"] = []evmodel.OdataIDLink{
		{OdataID: "/redfish/v1/Fabrics/48591de0-4859-1108-c35a-6c50110872da"},
	}
	SubscriptionReq["Destination"] = "https://10.10.10.24:8070/Destination4"
	postBody, _ = json.Marshal(&SubscriptionReq)
	req2 := &eventsproto.EventSubRequest{
		SessionToken: "token",
		PostBody:     postBody,
	}

	resp = p.CreateEventSubscription(taskID, sessionUserName, req2)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status Code should be StatusCreated")
}

func TestRmDupEleStrSlc(t *testing.T) {
	tests := []struct {
		name  string
		arg1  []string
		arg2  int
		want1 []string
		want2 int
	}{
		{
			name:  "Empty string slice",
			arg1:  []string{},
			arg2:  0,
			want1: []string{},
			want2: 0,
		},
		{
			name:  "String slice with single element",
			arg1:  []string{"1"},
			arg2:  1,
			want1: []string{"1"},
			want2: 1,
		},
		{
			name:  "String slice with no duplicate elements",
			arg1:  []string{"1", "2", "3", "4", "5"},
			arg2:  5,
			want1: []string{"1", "2", "3", "4", "5"},
			want2: 5,
		},
		{
			name:  "String slice with duplicate elements",
			arg1:  []string{"1", "2", "3", "4", "5", "5", "4", "3", "2", "1"},
			arg2:  10,
			want1: []string{"1", "2", "3", "4", "5"},
			want2: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			removeDuplicatesFromSlice(&tt.arg1, &tt.arg2)
			if tt.arg2 != tt.want2 && !reflect.DeepEqual(tt.arg1, tt.want1) {
				t.Errorf("TestRmDupEleStrSlc got1 = %v, want1 = %v", tt.arg1, tt.want2)
				t.Errorf("TestRmDupEleStrSlc got2 = %v, want2 = %v", tt.arg2, tt.want2)
			}
		})
	}
}

func TestCheckCollectionSubscription(t *testing.T) {
	config.SetUpMockConfig(t)
	p := getMockMethods()
	originResources := "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1"
	protocol := "Redfish"
	p.checkCollectionSubscription(originResources, protocol)
	devSub, _ := p.GetDeviceSubscriptions("*" + originResources)
	assert.Equal(t, "https://100.100.100.100/ODIM/v1/Subscriptions/1", devSub.Location, "Location should be https://100.100.100.100/ODIM/v1/Subscriptions/12")
	assert.Equal(t, "100.100.100.100", devSub.EventHostIP, "EventHostIP should be 100.100.100.100")
}
