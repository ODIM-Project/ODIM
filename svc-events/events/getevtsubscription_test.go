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
	"errors"
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

func getMockMethods() ExternalInterfaces {
	return ExternalInterfaces{
		External: External{
			ContactClient:   evcommon.MockContactClient,
			Auth:            evcommon.MockIsAuthorized,
			CreateChildTask: evcommon.MockCreateChildTask,
			UpdateTask:      evcommon.MockUpdateTask,
		},
		DB: DB{
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
			SaveUndeliveredEvents:            evcommon.MockSaveUndeliveredEvents,
			SaveDeviceSubscription:           evcommon.MockSaveDeviceSubscription,
			GetUndeliveredEvents:             evcommon.MockGetUndeliveredEvents,
			GetUndeliveredEventsFlag:         evcommon.MockGetUndeliveredEventsFlag,
			SetUndeliveredEventsFlag:         evcommon.MockSetUndeliveredEventsFlag,
			DeleteUndeliveredEventsFlag:      evcommon.MockDeleteUndeliveredEventsFlag,
			DeleteUndeliveredEvents:          evcommon.MockDeleteUndeliveredEvents,
			GetAggregateData:                 evcommon.MockGetAggregateDatacData,
			SaveAggregateSubscription:        evcommon.MockSaveAggregateSubscription,
			GetAggregateHosts:                evcommon.MockGetAggregateHosts,
			UpdateAggregateHosts:             evcommon.MockSaveAggregateSubscription,
			GetAggregateList:                 evcommon.MockGetAggregateHosts,
		},
	}
}
func TestGetEventSubscriptionsCollection(t *testing.T) {
	common.SetUpMockConfig()
	pc := getMockMethods()
	req := &eventsproto.EventRequest{
		SessionToken: "validToken",
	}

	// positive test case
	resp := pc.GetEventSubscriptionsCollection(req)
	data := resp.Body.(evresponse.ListResponse)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")
	assert.Equal(t, 1, data.MembersCount, "MembersCount should be 1")

	// Negative test cases
	// Invalid token
	req1 := &eventsproto.EventRequest{
		SessionToken: "InValidToken",
	}
	resp = pc.GetEventSubscriptionsCollection(req1)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status Code should be StatusUnauthorized")

	pc.DB.GetEvtSubscriptions = func(s string) ([]evmodel.Subscription, error) {
		return []evmodel.Subscription{}, errors.New("")
	}
	resp = pc.GetEventSubscriptionsCollection(req)

}

func TestGetEventSubscription(t *testing.T) {
	common.SetUpMockConfig()
	pc := getMockMethods()

	// positive test case
	req := &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "81de0110-c35a-4859-984c-072d6c5a32d7",
	}
	resp := pc.GetEventSubscriptionsDetails(req)
	data := resp.Body.(*evresponse.SubscriptionResponse)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")
	assert.Equal(t, "81de0110-c35a-4859-984c-072d6c5a32d7", data.Response.ID, "ID should be 1")
	assert.Equal(t, "Subscription", data.Response.Name, "Name should be Subscription")
	assert.Equal(t, "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1", data.OriginResources[0].OdataID, " OdataID should be same /redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1")

	// Negative test cases
	// Invalid token
	req1 := &eventsproto.EventRequest{
		SessionToken: "InValidToken",
	}
	resp = pc.GetEventSubscriptionsDetails(req1)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status Code should be StatusUnauthorized")

	// invalid subscription id
	req = &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "1234",
	}
	resp = pc.GetEventSubscriptionsDetails(req)
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status Code should be StatusNotFound")

	pc.DB.GetEvtSubscriptions = func(s string) ([]evmodel.Subscription, error) {
		return []evmodel.Subscription{}, errors.New("")
	}
	resp = pc.GetEventSubscriptionsDetails(req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusOK")

	pc.DB.GetEvtSubscriptions = func(s string) ([]evmodel.Subscription, error) {
		return []evmodel.Subscription{{UserName: "Admin", SubscriptionID: "test"}}, nil
	}
	resp = pc.GetEventSubscriptionsDetails(req)
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status Code should be StatusOK")

}

func TestExternalInterfaces_IsAggregateHaveSubscription(t *testing.T) {
	config.SetUpMockConfig(t)
	pc := getMockMethods()
	pc.Auth = func(s1 string, s2, s3 []string) response.RPC {
		return response.RPC{
			StatusCode: 400,
		}
	}
	pc.IsAggregateHaveSubscription(&eventsproto.EventUpdateRequest{})
	pc.Auth = func(s1 string, s2, s3 []string) response.RPC {
		return response.RPC{
			StatusCode: 200,
		}
	}
	pc.IsAggregateHaveSubscription(&eventsproto.EventUpdateRequest{})
	pc.DB.GetEvtSubscriptions = func(s string) ([]evmodel.Subscription, error) {
		return []evmodel.Subscription{{UserName: "admin"}}, nil
	}
	pc.IsAggregateHaveSubscription(&eventsproto.EventUpdateRequest{})
}
