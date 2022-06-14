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

// Package rpc ...
package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/events"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/ODIM-Project/ODIM/svc-events/evresponse"
	"github.com/stretchr/testify/assert"
)

func getMockPluginContactInitializer() *Events {
	connector := &events.ExternalInterfaces{
		External: events.External{
			ContactClient:   evcommon.MockContactClient,
			Auth:            evcommon.MockIsAuthorized,
			CreateTask:      evcommon.MockCreateTask,
			CreateChildTask: evcommon.MockCreateChildTask,
			UpdateTask:      evcommon.MockUpdateTask,
		},
		DB: events.DB{
			GetSessionUserName:               evcommon.MockGetSessionUserName,
			GetTarget:                        evcommon.MockGetTarget,
			GetPluginData:                    evcommon.MockGetPluginData,
			GetFabricData:                    evcommon.MockGetFabricData,
			GetEvtSubscriptions:              evcommon.MockGetEvtSubscriptions,
			GetDeviceSubscriptions:           evcommon.MockGetDeviceSubscriptions,
			SaveEventSubscription:            evcommon.MockSaveEventSubscription,
			UpdateEventSubscription:          evcommon.MockUpdateEventSubscription,
			DeleteDeviceSubscription:         evcommon.MockDeleteDeviceSubscription,
			DeleteEvtSubscription:            evcommon.MockDeleteEvtSubscription,
			UpdateDeviceSubscriptionLocation: evcommon.MockUpdateDeviceSubscriptionLocation,
			GetAllKeysFromTable:              evcommon.MockGetAllKeysFromTable,
			GetAllFabrics:                    evcommon.MockGetAllFabrics,
			GetAllMatchingDetails:            evcommon.MockGetAllMatchingDetails,
			SaveDeviceSubscription:           evcommon.MockSaveDeviceSubscription,
			SaveUndeliveredEvents:            evcommon.MockSaveUndeliveredEvents,
		},
	}
	return &Events{
		Connector: connector,
	}
}

func TestGetEventService(t *testing.T) {
	config.SetUpMockConfig(t)
	var ctx context.Context
	events := getMockPluginContactInitializer()
	req := &eventsproto.EventSubRequest{
		SessionToken: "validToken",
	}

	resp, err := events.GetEventService(ctx, req)
	assert.Nil(t, err, "There should be no error")

	var eventServiceResp evresponse.EventServiceResponse
	json.Unmarshal(resp.Body, &eventServiceResp)

	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "Status code should be StatusOK.")
	assert.True(t, eventServiceResp.ServiceEnabled, "Service should be Enabled ")
	assert.Equal(t, eventServiceResp.Status.State, "Enabled", "serviceState should be Enabled.")
	assert.Equal(t, eventServiceResp.Status.Health, "OK", "Health Status should be OK.")
	assert.Equal(t, eventServiceResp.EventFormatTypes, []string{"Event", "MetricReport"},
		"EventFormatTypes: Possible values are Event and MetricReport")

	req = &eventsproto.EventSubRequest{
		SessionToken: "InValidToken",
	}

	esResp, _ := events.GetEventService(ctx, req)
	assert.Equal(t, int(esResp.StatusCode), http.StatusUnauthorized, "Status code should be StatusUnauthorized.")
}

func TestCreateEventSubscription(t *testing.T) {
	config.SetUpMockConfig(t)
	var ctx context.Context
	events := getMockPluginContactInitializer()
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
		},
	}
	postBody, _ := json.Marshal(&SubscriptionReq)

	// Positive test cases
	req := &eventsproto.EventSubRequest{
		SessionToken: "validToken",
		PostBody:     postBody,
	}

	resp, err := events.CreateEventSubscription(ctx, req)
	assert.Nil(t, err, "There should be no error")

	assert.Equal(t, int(resp.StatusCode), http.StatusAccepted, "Status code should be StatusAccepted.")

	req1 := &eventsproto.EventSubRequest{
		SessionToken: "InValidToken",
	}

	esResp, _ := events.CreateEventSubscription(ctx, req1)
	assert.Equal(t, int(esResp.StatusCode), http.StatusUnauthorized, "Status code should be StatusUnauthorized.")

	req.SessionToken = "token1"
	esRespTest, _ := events.CreateEventSubscription(ctx, req)
	assert.Equal(t, int(esRespTest.StatusCode), http.StatusUnauthorized, "Status code should be StatusUnauthorized.")

	req.SessionToken = "token"

	esRespTest2, _ := events.CreateEventSubscription(ctx, req)
	assert.Equal(t, int(esRespTest2.StatusCode), http.StatusInternalServerError, "Status code should be StatusUnauthorized.")
}

func TestSubmitTestEvent(t *testing.T) {
	config.SetUpMockConfig(t)
	var ctx context.Context
	events := getMockPluginContactInitializer()
	event := map[string]interface{}{
		"MemberID":          "1",
		"EventType":         "Alert",
		"EventID":           "123",
		"Severity":          "OK",
		"Message":           "IndicatorChanged",
		"MessageId":         "IndicatorChanged",
		"OriginOfCondition": "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
	}

	message, err := json.Marshal(event)
	if err != nil {
		t.Errorf("expected err is nil but got : %v", err)
	}
	// Positive test cases
	req := &eventsproto.EventSubRequest{
		SessionToken: "validToken",
		PostBody:     message,
	}

	resp, errTest := events.SubmitTestEvent(ctx, req)
	assert.Nil(t, errTest, "There should be no error")
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetEventSubscriptionsCollection(t *testing.T) {
	config.SetUpMockConfig(t)
	var ctx context.Context
	events := getMockPluginContactInitializer()
	// Positive test cases
	req := &eventsproto.EventRequest{
		SessionToken: "validToken",
	}

	resp, err := events.GetEventSubscriptionsCollection(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "Status code should be StatusOK.")

	var evResp = &evresponse.ListResponse{}
	json.Unmarshal(resp.Body, evResp)
	assert.Equal(t, 1, evResp.MembersCount, "MembersCount should be 1")

}

func TestGetEventSubscriptions(t *testing.T) {
	config.SetUpMockConfig(t)
	var ctx context.Context
	events := getMockPluginContactInitializer()
	// Positive test cases
	req := &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "81de0110-c35a-4859-984c-072d6c5a32d7",
	}

	esResp, err := events.GetEventSubscription(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, int(esResp.StatusCode), http.StatusOK, "Status code should be StatusOK.")

	var evResp = &evresponse.SubscriptionResponse{}
	json.Unmarshal(esResp.Body, evResp)
	assert.Equal(t, "81de0110-c35a-4859-984c-072d6c5a32d7", evResp.Response.ID, "ID should be 81de0110-c35a-4859-984c-072d6c5a32d7")

	req.EventSubscriptionID = "81de0110"
	//resp := &eventsproto.EventSubResponse{}
	resp, _ := events.GetEventSubscription(ctx, req)
	assert.Equal(t, int(resp.StatusCode), http.StatusNotFound, "Status code should be StatusNotFound.")
}

func TestDeleteEventSubscription(t *testing.T) {
	config.SetUpMockConfig(t)
	var ctx context.Context
	events := getMockPluginContactInitializer()
	// Positive test cases
	req := &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "81de0110-c35a-4859-984c-072d6c5a32d7",
	}

	resp, err := events.DeleteEventSubscription(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "Status code should be StatusOK.")

	req.EventSubscriptionID = "81de0110"

	delResp, _ := events.DeleteEventSubscription(ctx, req)
	assert.Equal(t, int(delResp.StatusCode), http.StatusNotFound, "Status code should be StatusNotFound.")
}

func TestDeleteEventSubscriptionwithUUID(t *testing.T) {
	config.SetUpMockConfig(t)
	var ctx context.Context
	events := getMockPluginContactInitializer()
	// Positive test cases
	req := &eventsproto.EventRequest{
		SessionToken: "validToken",
		UUID:         "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
	}

	resp, err := events.DeleteEventSubscription(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, int(resp.StatusCode), http.StatusNoContent, "Status code should be StatusNoContent.")

	req.UUID = "81de0110"

	delResp, _ := events.DeleteEventSubscription(ctx, req)
	assert.Equal(t, int(delResp.StatusCode), http.StatusBadRequest, "Status code should be StatusBadRequest.")
}

func TestCreateDefaultSubscriptions(t *testing.T) {
	config.SetUpMockConfig(t)
	var ctx context.Context
	events := getMockPluginContactInitializer()
	// Positive test cases
	req := &eventsproto.DefaultEventSubRequest{
		SystemID:      []string{"systemid"},
		EventTypes:    []string{"Alert"},
		MessageIDs:    []string{},
		ResourceTypes: []string{},
		Protocol:      "redfish",
	}

	_, err := events.CreateDefaultEventSubscription(ctx, req)
	assert.Nil(t, err, "There should be no error")

}

func TestSubscribeEMB(t *testing.T) {
	var ctx context.Context
	events := getMockPluginContactInitializer()
	evcommon.EMBTopics.TopicsList = make(map[string]bool)
	req := &eventsproto.SubscribeEMBRequest{
		PluginID:     "GRF",
		EMBQueueName: []string{"topic"},
	}

	resp, err := events.SubsribeEMB(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.True(t, resp.Status, "status should be true")
}
