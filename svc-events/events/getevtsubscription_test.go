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
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-events/evresponse"
	"github.com/stretchr/testify/assert"
)

func mockIsAuthorized(sessionToken string, privileges, oemPrivileges []string) response.RPC {
	if sessionToken != "validToken" && sessionToken != "token" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "", nil, nil)
	}
	return common.GeneralError(http.StatusOK, response.Success, "", nil, nil)
}

func TestGetEventSubscriptionsCollection(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	storeTestEventDetails(t)
	pc := PluginContact{
		Auth: mockIsAuthorized,
	}
	req := &eventsproto.EventRequest{
		SessionToken: "validToken",
	}

	// positive test case
	resp := pc.GetEventSubscriptionsCollection(req)
	data := resp.Body.(evresponse.ListResponse)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")
	assert.Equal(t, 7, data.MembersCount, "MembersCount should be 7")

	// Negative test cases
	// Invalid token
	req1 := &eventsproto.EventRequest{
		SessionToken: "InValidToken",
	}
	resp = pc.GetEventSubscriptionsCollection(req1)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status Code should be StatusUnauthorized")

}

func TestGetEventSubscription(t *testing.T) {
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	storeTestEventDetails(t)
	pc := PluginContact{
		Auth: mockIsAuthorized,
	}

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
	assert.Equal(t, "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1", data.OriginResources[0].OdataID, " OdataID should be same /redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e:1")

	// Negative test cases
	// Invalid token
	req1 := &eventsproto.EventRequest{
		SessionToken: "InValidToken",
	}
	resp = pc.GetEventSubscriptionsDetails(req1)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status Code should be StatusUnauthorized")

	common.TruncateDB(common.OnDisk)
	resp = pc.GetEventSubscriptionsDetails(req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

}
