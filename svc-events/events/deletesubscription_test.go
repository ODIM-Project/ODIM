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

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/stretchr/testify/assert"
)

func TestDeleteEventSubscription(t *testing.T) {
	// Intializing plugin token
	evcommon.Token.Tokens = map[string]string{
		"ILO": "token",
	}
	config.SetUpMockConfig(t)

	pc := getMockMethods()

	// positive test case with basic auth type
	req := &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "81de0110-c35a-4859-984c-072d6c5a32d7",
	}
	resp := pc.DeleteEventSubscriptionsDetails(req)
	data := resp.Body.(response.Response)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")
	assert.Equal(t, "81de0110-c35a-4859-984c-072d6c5a32d7", data.ID, "ID should be 81de0110-c35a-4859-984c-072d6c5a32d7")

	// positive test case with basic auth type
	req = &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "71de0110-c35a-4859-984c-072d6c5a32d8",
	}
	resp = pc.DeleteEventSubscriptionsDetails(req)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")

	// positive test case deletion of collection subscription
	req = &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "71de0110-c35a-4859-984c-072d6c5a3211",
	}
	resp = pc.DeleteEventSubscriptionsDetails(req)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")

	// Negative test cases
	// if subscription id is bot present
	req1 := &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "de018110-4859-984c-c35a-0a32d772d6c5",
	}
	resp = pc.DeleteEventSubscriptionsDetails(req1)
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status Code should be StatusNotFound")

	// Invalid token
	req2 := &eventsproto.EventRequest{
		SessionToken: "InValidToken",
	}
	resp = pc.DeleteEventSubscriptionsDetails(req2)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status Code should be StatusUnauthorized")
}

func TestDeleteEventSubscriptionOnDeletServer(t *testing.T) {
	config.SetUpMockConfig(t)
	// Intializing plugin token
	evcommon.Token.Tokens = map[string]string{
		"ILO": "token",
	}
	pc := getMockMethods()

	// positive test case
	req := &eventsproto.EventRequest{
		SessionToken: "validToken",
		UUID:         "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
	}
	resp := pc.DeleteEventSubscriptions(req)
	assert.Equal(t, http.StatusNoContent, int(resp.StatusCode), "Status Code should be StatusNoContent")

	// Negative test cases
	// if UUID is invalid
	req1 := &eventsproto.EventRequest{
		SessionToken: "validToken",
		UUID:         "de018110-4859-984c-c35a-0a32d772d6c5",
	}
	resp = pc.DeleteEventSubscriptions(req1)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusNotFound")

	// if UUID is is not present in DB
	req1 = &eventsproto.EventRequest{
		SessionToken: "validToken",
		UUID:         "/redfish/v1/Systems/de018110-4859-984c-c35a-0a32d772d6c5.1",
	}

	resp = pc.DeleteEventSubscriptions(req1)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusNotFound")

	req = &eventsproto.EventRequest{
		SessionToken: "validToken",
		UUID:         "/redfish/v1/Systems/abab09db-e7a9-4352-8df0-5e41315a2a4c.1",
	}
	resp = pc.DeleteEventSubscriptions(req)
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status Code should be StatusNotFound")

}

func TestDeleteEventSubscriptionOnFabrics(t *testing.T) {
	// Intializing plugin token
	evcommon.Token.Tokens = map[string]string{
		"CFM": "token",
	}
	config.SetUpMockConfig(t)
	pc := getMockMethods()

	// positive test case with basic auth type
	req := &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "71de0110-c35a-4859-984c-072d6c5a32d9",
	}
	resp := pc.DeleteEventSubscriptionsDetails(req)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")
}

func TestIsCollectionOriginResourceURI(t *testing.T) {
	config.SetUpMockConfig(t)
	tests := []struct {
		name string
		arg  string
		want bool
	}{
		{
			name: "Positive: First element in list",
			arg:  "/redfish/v1/Systems",
			want: true,
		},
		{
			name: "Positive: Last element in list",
			arg:  "/redfish/v1/TaskService/Tasks",
			want: true,
		},
		{
			name: "Positive: Middle element in list",
			arg:  "/redfish/v1/Fabrics/",
			want: true,
		},
		{
			name: "Negative: Empty string",
			arg:  "",
			want: false,
		},
		{
			name: "Negative: Non-existent element",
			arg:  "non-existent string",
			want: false,
		},
		{
			name: "Negative: Non-existent element -2",
			arg:  "/redfish/v1/Fabrics1",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCollectionOriginResourceURI(tt.arg); got != tt.want {
				t.Errorf("isCollectionOriginResourceURI got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestDeleteFabricsSubscription(t *testing.T) {
	// Intializing plugin token
	evcommon.Token.Tokens = map[string]string{
		"ILO": "token",
	}
	config.SetUpMockConfig(t)
	pc := getMockMethods()

	password, err := evcommon.GetEncryptedKey([]byte("Password"))
	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	// positive test case with basic auth type
	plugin := &evmodel.Plugin{
		IP:                "odim.domain.com",
		Port:              "1234",
		Password:          password,
		Username:          "admin",
		ID:                "ILO",
		PreferredAuthType: "BasicAuth",
		PluginType:        "ILO",
	}
	resp, err := pc.DeleteFabricsSubscription("", plugin)
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")

	// Negative test cases
	// if subscription id is not present
	plugin.IP = "10.10.10.10"
	resp, err = pc.DeleteFabricsSubscription("", plugin)
	assert.NotNil(t, err, "error should not be nil")
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status Code should be StatusNotFound")
}

func TestDeleteSubscription(t *testing.T) {
	config.SetUpMockConfig(t)
	pc := getMockMethods()

	encryptedData, err := evcommon.GetEncryptedKey([]byte("testData"))
	if err != nil {
		t.Fatalf("%v", err.Error())
	}
	target := evmodel.Target{
		ManagerAddress: "100.100.100.100",
		Password:       encryptedData,
		UserName:       "admin",
		DeviceUUID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
		PluginID:       "GRF",
	}

	err = pc.deleteSubscription(&target, "/redfish/v1/Systems")
	assert.Nil(t, err, "error should be nil")

	target.PluginID = "non-existent"
	err = pc.deleteSubscription(&target, "/redfish/v1/Systems")
	assert.NotNil(t, err, "error should not be nil")
}
