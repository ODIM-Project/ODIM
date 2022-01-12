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
package handle

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	iris "github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
)

func mockGetEventServiceRPC(req eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error) {
	var response *eventsproto.EventSubResponse
	if req.SessionToken == "ValidToken" {
		response = &eventsproto.EventSubResponse{
			StatusCode: http.StatusOK,
		}
	} else if req.SessionToken == "InValidToken" {
		response = &eventsproto.EventSubResponse{
			StatusCode: http.StatusUnauthorized,
		}
	} else if req.SessionToken == "" {
		response = &eventsproto.EventSubResponse{
			StatusCode: http.StatusUnauthorized,
		}
	} else if req.SessionToken == "token" {
		return &eventsproto.EventSubResponse{}, fmt.Errorf("RPC Error")
	}
	return response, nil
}
func mockCreateEventSubscriptionRPC(req eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error) {
	var response *eventsproto.EventSubResponse
	if req.SessionToken == "ValidToken" {
		response = &eventsproto.EventSubResponse{
			StatusCode: http.StatusCreated,
		}
	} else if req.SessionToken == "InValidToken" {
		response = &eventsproto.EventSubResponse{
			StatusCode: http.StatusUnauthorized,
		}
	} else if req.SessionToken == "" {
		response = &eventsproto.EventSubResponse{
			StatusCode: http.StatusUnauthorized,
		}
	} else if req.SessionToken == "token" {
		return &eventsproto.EventSubResponse{}, fmt.Errorf("RPC Error")
	}
	return response, nil
}

func mockEventSubscriptionRPC(req eventsproto.EventRequest) (*eventsproto.EventSubResponse, error) {
	var response *eventsproto.EventSubResponse
	if req.SessionToken == "ValidToken" && req.EventSubscriptionID == "1A" {
		response = &eventsproto.EventSubResponse{
			StatusCode: http.StatusCreated,
		}
	} else if req.SessionToken == "InValidToken" {
		response = &eventsproto.EventSubResponse{
			StatusCode: http.StatusUnauthorized,
		}
	} else if req.SessionToken == "" {
		response = &eventsproto.EventSubResponse{
			StatusCode: http.StatusUnauthorized,
		}
	} else if req.SessionToken == "token" {
		return &eventsproto.EventSubResponse{}, fmt.Errorf("RPC Error")
	}
	return response, nil
}
func mockGetEventSubscriptionRPC(req eventsproto.EventRequest) (*eventsproto.EventSubResponse, error) {
	var response *eventsproto.EventSubResponse
	if req.SessionToken == "ValidToken" {
		response = &eventsproto.EventSubResponse{
			StatusCode: http.StatusOK,
		}
	} else if req.SessionToken == "InValidToken" {
		response = &eventsproto.EventSubResponse{
			StatusCode: http.StatusUnauthorized,
		}
	} else if req.SessionToken == "" {
		response = &eventsproto.EventSubResponse{
			StatusCode: http.StatusUnauthorized,
		}
	} else if req.SessionToken == "token" {
		return &eventsproto.EventSubResponse{}, fmt.Errorf("RPC Error")
	}
	return response, nil
}

func TestGetEventServiceRPC(t *testing.T) {
	header["Allow"] = []string{"GET"}
	defer delete(header, "Allow")
	var event EventsRPCs
	event.GetEventServiceRPC = mockGetEventServiceRPC
	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1/EventService")
	redfishRoutes.Get("/", event.GetEventService)
	e := httptest.New(t, mockApp)
	// test with valid token
	e.GET(
		"/redfish/v1/EventService",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK).Headers().Equal(header)

	// test with Invalid token
	e.GET(
		"/redfish/v1/EventService",
	).WithHeader("X-Auth-Token", "InValidToken").Expect().Status(http.StatusUnauthorized)

	// test without token
	e.GET(
		"/redfish/v1/EventService",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)
	e.GET(
		"/redfish/v1/EventService",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)

}

func TestCreateEventSubscriptionRPC(t *testing.T) {
	var event EventsRPCs
	event.CreateEventSubscriptionRPC = mockCreateEventSubscriptionRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Post("/EventService/Subscriptions", event.CreateEventSubscription)
	body := map[string]interface{}{
		"Name":                 "EventSubscription",
		"Destination":          "https://localhost/redfish/v1/EventService",
		"EventTypes":           []string{"Alert"},
		"MessageIds":           []string{},
		"ResourceTypes":        []string{},
		"Context":              "ABCDEFGHJLKJ",
		"Protocol":             "Redfish",
		"SubscriptionType":     "RedfishEvent",
		"EventFormatType":      "Event",
		"SubordinateResources": true,
		"OriginResources": []string{
			"/redfish/v1/Systems/7ff3bd97-c41c-5de0-937d-85d390691b73.1",
		},
	}
	e := httptest.New(t, mockApp)
	// test with valid token
	e.POST(
		"/redfish/v1/EventService/Subscriptions",
	).WithHeader("X-Auth-Token", "ValidToken").WithJSON(body).Expect().Status(http.StatusCreated).Headers().Equal(header)

	// test with Invalid token
	e.POST(
		"/redfish/v1/EventService/Subscriptions",
	).WithHeader("X-Auth-Token", "InValidToken").WithJSON(body).Expect().Status(http.StatusUnauthorized).Headers().Equal(header)

	// test without token
	e.POST(
		"/redfish/v1/EventService/Subscriptions",
	).WithHeader("X-Auth-Token", "").WithJSON(body).Expect().Status(http.StatusUnauthorized)

	// test without RequestBody
	e.POST(
		"/redfish/v1/EventService/Subscriptions",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusBadRequest)

	// test for RPC Error
	e.POST(
		"/redfish/v1/EventService/Subscriptions",
	).WithHeader("X-Auth-Token", "token").WithJSON(body).Expect().Status(http.StatusInternalServerError)
}

func TestSubmitTestEventRPC(t *testing.T) {
	var event EventsRPCs
	//	var timeStamp time.Time
	event.SubmitTestEventRPC = mockCreateEventSubscriptionRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Post("/EventService/Actions/EventService.SubmitTestEvent", event.SubmitTestEvent)
	body := map[string]interface{}{
		"EventGroupId":      1,
		"EventId":           "eventId.123",
		"EventTimestamp":    time.Now(),
		"EventType":         "Alert",
		"Message":           "Test Message",
		"MessageArgs":       []string{},
		"MessageId":         "MessageId.123",
		"OriginOfCondition": "/redfish/v1/Systems/7ff3bd97-c41c-5de0-937d-85d390691b73.1",
		"Severity":          "Critical",
	}
	e := httptest.New(t, mockApp)
	// test with valid token
	e.POST(
		"/redfish/v1/EventService/Actions/EventService.SubmitTestEvent",
	).WithHeader("X-Auth-Token", "ValidToken").WithJSON(body).Expect().Status(http.StatusCreated).Headers().Equal(header)

	// test with Invalid token
	e.POST(
		"/redfish/v1/EventService/Actions/EventService.SubmitTestEvent",
	).WithHeader("X-Auth-Token", "InValidToken").WithJSON(body).Expect().Status(http.StatusUnauthorized)

	// test without token
	e.POST(
		"/redfish/v1/EventService/Actions/EventService.SubmitTestEvent",
	).WithHeader("X-Auth-Token", "").WithJSON(body).Expect().Status(http.StatusUnauthorized)

	// test without requestBody
	e.POST(
		"/redfish/v1/EventService/Actions/EventService.SubmitTestEvent",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusBadRequest)

	// test for RPC error
	e.POST(
		"/redfish/v1/EventService/Actions/EventService.SubmitTestEvent",
	).WithHeader("X-Auth-Token", "token").WithJSON(body).Expect().Status(http.StatusInternalServerError)
}

func TestDeleteEventSubscriptionRPC(t *testing.T) {
	var s EventsRPCs
	s.DeleteEventSubscriptionRPC = mockGetEventSubscriptionRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Delete("/EventService/Subscriptions/{id}", s.DeleteEventSubscription)
	e := httptest.New(t, mockApp)

	// test with valid token
	e.DELETE(
		"/redfish/v1/EventService/Subscriptions/1A",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK).Headers().Equal(header)

	// test with invalid token
	e.DELETE(
		"/redfish/v1/EventService/Subscriptions/1A",
	).WithHeader("X-Auth-Token", "InValidToken").Expect().Status(http.StatusUnauthorized)

	// test without token
	e.DELETE(
		"/redfish/v1/EventService/Subscriptions/1A",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)

	// test without token
	e.DELETE(
		"/redfish/v1/EventService/Subscriptions/1A",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetEventSubscriptionRPC(t *testing.T) {
	var s EventsRPCs
	s.GetEventSubscriptionRPC = mockGetEventSubscriptionRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/EventService/Subscriptions/{id}", s.GetEventSubscription)
	e := httptest.New(t, mockApp)
	// test with valid token
	e.GET(
		"/redfish/v1/EventService/Subscriptions/1A",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)

	// test with invalid token
	e.GET(
		"/redfish/v1/EventService/Subscriptions/1A",
	).WithHeader("X-Auth-Token", "InValidToken").Expect().Status(http.StatusUnauthorized)

	// test without token
	e.GET(
		"/redfish/v1/EventService/Subscriptions/1A",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)

	// test for RPC error
	e.GET(
		"/redfish/v1/EventService/Subscriptions/1A",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}

func TestGetEventSubscriptionsCollectionRPC(t *testing.T) {
	var s EventsRPCs
	s.GetEventSubscriptionsCollectionRPC = mockGetEventSubscriptionRPC

	mockApp := iris.New()
	redfishRoutes := mockApp.Party("/redfish/v1")
	redfishRoutes.Get("/EventService/Subscriptions", s.GetEventSubscriptionsCollection)
	e := httptest.New(t, mockApp)
	e.GET(
		"/redfish/v1/EventService/Subscriptions",
	).WithHeader("X-Auth-Token", "ValidToken").Expect().Status(http.StatusOK)

	// test with invalid token
	e.GET(
		"/redfish/v1/EventService/Subscriptions",
	).WithHeader("X-Auth-Token", "InValidToken").Expect().Status(http.StatusUnauthorized)

	// test without token
	e.GET(
		"/redfish/v1/EventService/Subscriptions",
	).WithHeader("X-Auth-Token", "").Expect().Status(http.StatusUnauthorized)

	// test for RPC Error
	e.GET(
		"/redfish/v1/EventService/Subscriptions",
	).WithHeader("X-Auth-Token", "token").Expect().Status(http.StatusInternalServerError)
}
