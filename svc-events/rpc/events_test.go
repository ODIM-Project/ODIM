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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/ODIM-Project/ODIM/svc-events/evresponse"
	"github.com/stretchr/testify/assert"
)

func mockContactClient(url, method, token string, odataID string, body interface{}, credentials map[string]string) (*http.Response, error) {
	fmt.Println(url)
	if url == "https://localhost:1234/ODIM/v1/EventSubscriptions" {
		body := `{"MessageId": "` + response.Success + `"}`
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:1234/ODIM/v1/Subscriptions" {
		body := `{"MessageId": "` + response.Success + `"}`
		return &http.Response{
			StatusCode: http.StatusCreated,
			Header: map[string][]string{
				"location": {"https://localhost:1234/ODIM/v1/Subscriptions/12"},
			},
			Body: ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:1234/ODIM/v1/Sessions" {
		body := `{"MessageId": "` + response.Success + `"}`

		r := &http.Response{
			StatusCode: http.StatusCreated,
			Header: map[string][]string{
				"X-Auth-Token": {"token"},
			},
			Body: ioutil.NopCloser(bytes.NewBufferString(body)),
		}
		return r, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func storeTestEventDetails(t *testing.T) {
	subarr := []evmodel.Subscription{
		// if SubordinateResources true
		{
			UserName:             "admin",
			SubscriptionID:       "81de0110-c35a-4859-984c-072d6c5a32d7",
			Destination:          "https://10.24.1.15:9090/events",
			Name:                 "Subscription",
			Context:              "context",
			EventTypes:           []string{"Alert"},
			MessageIds:           []string{"IndicatorChanged"},
			ResourceTypes:        []string{"#Event"},
			OriginResource:       "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
			OriginResources:      []string{"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"},
			Hosts:                []string{"10.4.1.2"},
			SubordinateResources: true,
		},
	}
	for _, sub := range subarr {
		err := evmodel.SaveEventSubscription(sub)
		if err != nil {
			t.Errorf("error: %v", err)
		}
	}
	var updatedDeviceSubscription = evmodel.DeviceSubscription{
		Location:        "https://10.24.1.2/EventService/Subscriptions/1",
		EventHostIP:     "10.4.1.2",
		OriginResources: []string{"/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"},
	}
	err := evmodel.SaveDeviceSubscription(updatedDeviceSubscription)
	if err != nil {
		t.Errorf("error: %v", err)
	}

}
func mockIsAuthorized(sessionToken string, privileges, oemPrivileges []string) response.RPC {
	if sessionToken != "validToken" && sessionToken != "token" && sessionToken != "token1" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "", nil, nil)
	}
	return common.GeneralError(http.StatusOK, response.Success, "", nil, nil)
}

func getMockedSessionUserName(sessionToken string) (string, error) {
	if sessionToken == "token" {
		return "non-admin", nil
	} else if sessionToken == "token1" {
		return "", fmt.Errorf("no details")
	}
	return "admin", nil
}

func mockCreateChildTask(sessionID, taskid string) (string, error) {
	return "123456", nil
}

func mockUpdateTask(task common.TaskData) error {
	return nil
}

func mockCreateTask(sessionusername string) (string, error) {
	if sessionusername == "non-admin" {
		return "", fmt.Errorf("no task details")
	}
	return "/redfish/v1/tasks/123", nil
}
func TestGetEventService(t *testing.T) {
	config.SetUpMockConfig(t)
	defer func() {
		err := common.TruncateDB(common.InMemory)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	var ctx context.Context
	events := new(Events)
	events.IsAuthorizedRPC = mockIsAuthorized
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

	req = &eventsproto.EventSubRequest{
		SessionToken: "InValidToken",
	}

	esResp, _ := events.GetEventService(ctx, req)
	assert.Equal(t, int(esResp.StatusCode), http.StatusUnauthorized, "Status code should be StatusUnauthorized.")
}

func TestCreateEventSubscription(t *testing.T) {
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
	var ctx context.Context
	events := new(Events)
	events.IsAuthorizedRPC = mockIsAuthorized
	events.GetSessionUserNameRPC = getMockedSessionUserName
	events.ContactClientRPC = mockContactClient
	events.CreateChildTaskRPC = mockCreateChildTask
	events.UpdateTaskRPC = mockUpdateTask
	events.CreateTaskRPC = mockCreateTask
	SubscriptionReq := map[string]interface{}{
		"Name":                 "EventSubscription",
		"Destination":          "https://10.24.1.24:8070/Destination1",
		"EventTypes":           []string{"Alert"},
		"Protocol":             "Redfish",
		"Context":              "Event Subscription",
		"SubscriptionType":     "RedfishEvent",
		"EventFormatType":      "Event",
		"SubordinateResources": true,
		"OriginResources": []evmodel.OdataIDLink{
			{OdataID: "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1"},
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
	var ctx context.Context
	events := new(Events)
	events.IsAuthorizedRPC = mockIsAuthorized
	events.GetSessionUserNameRPC = getMockedSessionUserName
	events.ContactClientRPC = mockContactClient
	event := map[string]interface{}{
		"MemberID":          "1",
		"EventType":         "Alert",
		"EventID":           "123",
		"Severity":          "OK",
		"Message":           "IndicatorChanged",
		"MessageId":         "IndicatorChanged",
		"OriginOfCondition": "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
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
	var ctx context.Context
	events := new(Events)
	events.IsAuthorizedRPC = mockIsAuthorized
	events.ContactClientRPC = mockContactClient
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
	var ctx context.Context
	events := new(Events)
	events.IsAuthorizedRPC = mockIsAuthorized
	events.ContactClientRPC = mockContactClient
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
	assert.Equal(t, int(resp.StatusCode), http.StatusNotFound, "Status code should be StatusBadRequest.")
}

func TestDeleteEventSubscription(t *testing.T) {
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
	var ctx context.Context
	events := new(Events)
	events.IsAuthorizedRPC = mockIsAuthorized
	events.ContactClientRPC = mockContactClient
	// Positive test cases
	req := &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "81de0110-c35a-4859-984c-072d6c5a32d7",
	}

	resp, err := events.DeleteEventSubscription(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, int(resp.StatusCode), http.StatusBadRequest, "Status code should be StatusNotFound.")

	req.EventSubscriptionID = "81de0110"

	delResp, _ := events.DeleteEventSubscription(ctx, req)
	assert.Equal(t, int(delResp.StatusCode), http.StatusNotFound, "Status code should be StatusNotFound.")
}

func TestDeleteEventSubscriptionwithUUID(t *testing.T) {
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
	var ctx context.Context
	events := new(Events)
	events.IsAuthorizedRPC = mockIsAuthorized
	events.ContactClientRPC = mockContactClient
	// Positive test cases
	req := &eventsproto.EventRequest{
		SessionToken: "validToken",
		UUID:         "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1",
	}

	resp, err := events.DeleteEventSubscription(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, int(resp.StatusCode), http.StatusBadRequest, "Status code should be StatusBadRequest.")

	req.UUID = "81de0110"

	delResp, _ := events.DeleteEventSubscription(ctx, req)
	assert.Equal(t, int(delResp.StatusCode), http.StatusBadRequest, "Status code should be StatusBadRequest.")
}

func TestCreateDefaultSubscriptions(t *testing.T) {
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

	var ctx context.Context
	events := new(Events)
	events.IsAuthorizedRPC = mockIsAuthorized
	events.ContactClientRPC = mockContactClient
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
	events := new(Events)
	evcommon.EMBTopics.TopicsList = make(map[string]bool)
	req := &eventsproto.SubscribeEMBRequest{
		PluginID:     "GRF",
		EMBQueueName: []string{"topic"},
	}

	resp, err := events.SubsribeEMB(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.True(t, resp.Status, "status should be true")
}
