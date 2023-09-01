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
	"fmt"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/stretchr/testify/assert"
)

func TestDeleteEventSubscription(t *testing.T) {
	// Initializing plugin token
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
	resp := pc.DeleteEventSubscriptionsDetails(evcommon.MockContext(), req)
	data := resp.Body.(response.Response)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")
	assert.Equal(t, "81de0110-c35a-4859-984c-072d6c5a32d7", data.ID, "ID should be 81de0110-c35a-4859-984c-072d6c5a32d7")

	pc.DB.GetEvtSubscriptions = func(s string) ([]evmodel.SubscriptionResource, error) { return nil, &errors.Error{} }
	resp = pc.DeleteEventSubscriptionsDetails(evcommon.MockContext(), req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusOK")
	pc = getMockMethods()

	pc.DB.DeleteEvtSubscription = func(s string) error { return &errors.Error{} }
	resp = pc.DeleteEventSubscriptionsDetails(evcommon.MockContext(), req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusOK")
	pc = getMockMethods()
	// positive test case with basic auth type
	req = &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "71de0110-c35a-4859-984c-072d6c5a32d8",
	}
	resp = pc.DeleteEventSubscriptionsDetails(evcommon.MockContext(), req)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")

	// positive test case deletion of collection subscription
	req = &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "71de0110-c35a-4859-984c-072d6c5a3211",
	}
	resp = pc.DeleteEventSubscriptionsDetails(evcommon.MockContext(), req)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")

	// Negative test cases
	// if subscription id is bot present
	req1 := &eventsproto.EventRequest{
		SessionToken:        "validToken",
		EventSubscriptionID: "de018110-4859-984c-c35a-0a32d772d6c5",
	}
	resp = pc.DeleteEventSubscriptionsDetails(evcommon.MockContext(), req1)
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status Code should be StatusNotFound")

	// Invalid token
	req2 := &eventsproto.EventRequest{
		SessionToken: "InValidToken",
	}
	resp = pc.DeleteEventSubscriptionsDetails(evcommon.MockContext(), req2)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status Code should be StatusUnauthorized")
}

func TestDeleteEventSubscriptionOnDeletedServer(t *testing.T) {
	config.SetUpMockConfig(t)
	// Initializing plugin token
	evcommon.Token.Tokens = map[string]string{
		"ILO": "token",
	}
	pc := getMockMethods()

	// positive test case
	req := &eventsproto.EventRequest{
		SessionToken: "validToken",
		UUID:         "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
	}
	resp := pc.DeleteEventSubscriptions(evcommon.MockContext(), req)
	assert.Equal(t, http.StatusNoContent, int(resp.StatusCode), "Status Code should be StatusNoContent")

	// Negative test cases
	// if UUID is invalid
	req1 := &eventsproto.EventRequest{
		SessionToken: "validToken",
		UUID:         "de018110-4859-984c-c35a-0a32d772d6c5",
	}
	resp = pc.DeleteEventSubscriptions(evcommon.MockContext(), req1)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusNotFound")

	// if UUID is is not present in DB
	req1 = &eventsproto.EventRequest{
		SessionToken: "validToken",
		UUID:         "/redfish/v1/Systems/de018110-4859-984c-c35a-0a32d772d6c5.1",
	}

	resp = pc.DeleteEventSubscriptions(evcommon.MockContext(), req1)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusNotFound")

	req = &eventsproto.EventRequest{
		SessionToken: "validToken",
		UUID:         "/redfish/v1/Systems/abab09db-e7a9-4352-8df0-5e41315a2a4c.1",
	}
	resp = pc.DeleteEventSubscriptions(evcommon.MockContext(), req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusNotFound")

	GetIPFromHostNameFunc = func(fqdn string) (string, error) {
		return "", fmt.Errorf("Not Found")
	}
	resp = pc.DeleteEventSubscriptions(evcommon.MockContext(), req)
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status Code should be ResourceNotFound")

	GetIPFromHostNameFunc = func(fqdn string) (string, error) {
		return common.GetIPFromHostName(fqdn)
	}
	pc.DB.GetEvtSubscriptions = func(s string) ([]evmodel.SubscriptionResource, error) {
		return nil, &errors.Error{}
	}
	resp = pc.DeleteEventSubscriptions(evcommon.MockContext(), req)
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status Code should be ResourceNotFound")

	pc = getMockMethods()

	DecryptWithPrivateKeyFunc = func(ciphertext []byte) ([]byte, error) {
		return nil, &errors.Error{}
	}

	req = &eventsproto.EventRequest{
		SessionToken: "validToken",
		UUID:         "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1",
	}
	resp = pc.DeleteEventSubscriptions(evcommon.MockContext(), req)
	assert.Equal(t, http.StatusInternalServerError, int(resp.StatusCode), "Status Code should be StatusNoContent")

	DecryptWithPrivateKeyFunc = func(ciphertext []byte) ([]byte, error) {
		return common.DecryptWithPrivateKey(ciphertext)
	}
	pc.DB.GetDeviceSubscriptions = func(s string) (*common.DeviceSubscription, error) {
		return nil, &errors.Error{}
	}
	resp = pc.DeleteEventSubscriptions(evcommon.MockContext(), req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusNotFound")

	pc = getMockMethods()
	pc.DB.DeleteEvtSubscription = func(s string) error {
		return &errors.Error{}
	}
	resp = pc.DeleteEventSubscriptions(evcommon.MockContext(), req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusNoContent")

	pc = getMockMethods()
	pc.DB.UpdateEventSubscription = func(s evmodel.SubscriptionResource) error {
		return &errors.Error{}
	}

	req.UUID = "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874d.1"
	resp = pc.DeleteEventSubscriptions(evcommon.MockContext(), req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusNoContent")

}

func TestDeleteEventSubscriptionOnFabrics(t *testing.T) {
	// Initializing plugin token
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
	resp := pc.DeleteEventSubscriptionsDetails(evcommon.MockContext(), req)
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
	// Initializing plugin token
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
	plugin := &common.Plugin{
		IP:                "odim.controller.com",
		Port:              "1234",
		Password:          password,
		Username:          "admin",
		ID:                "ILO",
		PreferredAuthType: "BasicAuth",
		PluginType:        "ILO",
	}
	resp, err := pc.DeleteFabricsSubscription(evcommon.MockContext(), "", plugin)
	assert.Nil(t, err, "error should be nil")
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")
	// Negative
	GetIPFromHostNameFunc = func(fqdn string) (string, error) { return "", fmt.Errorf("Not found") }
	resp, _ = pc.DeleteFabricsSubscription(evcommon.MockContext(), "", plugin)
	assert.Equal(t, http.StatusNotFound, int(resp.StatusCode), "Status Code should be StatusOK")

	GetIPFromHostNameFunc = func(fqdn string) (string, error) { return common.GetIPFromHostName(fqdn) }
	// Negative test cases
	// if subscription id is not present
	plugin.IP = "10.10.10.10"
	resp, err = pc.DeleteFabricsSubscription(evcommon.MockContext(), "", plugin)
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
	target := common.Target{
		ManagerAddress: "100.100.100.100",
		Password:       encryptedData,
		UserName:       "admin",
		DeviceUUID:     "6d4a0a66-7efa-578e-83cf-44dc68d2874e",
		PluginID:       "GRF",
	}

	err = pc.deleteSubscription(evcommon.MockContext(), &target, "/redfish/v1/Systems")
	assert.Nil(t, err, "error should be nil")

	target.PluginID = "non-existent"
	err = pc.deleteSubscription(evcommon.MockContext(), &target, "/redfish/v1/Systems")
	assert.NotNil(t, err, "error should not be nil")
}

func TestExternalInterfaces_DeleteAggregateSubscriptions(t *testing.T) {
	config.SetUpMockConfig(t)
	pc := getMockMethods()
	req := eventsproto.EventUpdateRequest{
		AggregateId: "71de0110-c35a-4859-984c-072d6c5a32d9",
	}
	pc.DeleteAggregateSubscriptions(evcommon.MockContext(), &req, true)
	pc.DB.GetEvtSubscriptions = func(s string) ([]evmodel.SubscriptionResource, error) {

		return []evmodel.SubscriptionResource{
			{
				EventDestination: &model.EventDestination{
					Destination: "https://localhost:9090/events",
					Name:        "Subscription",

					Context:              "context",
					EventTypes:           []string{"Alert"},
					MessageIds:           []string{},
					ResourceTypes:        []string{},
					OriginResources:      []model.Link{model.Link{Oid: "/redfish/v1/Fabrics/123456"}},
					SubordinateResources: true,
				},
				SubscriptionID: "71de0110-c35a-4859-984c-072d6c5a32d9",

				Hosts: []string{"localhost"},
			},
		}, nil
	}
	pc.DeleteAggregateSubscriptions(evcommon.MockContext(), &req, true)

	pc.DB.GetEvtSubscriptions = func(s string) ([]evmodel.SubscriptionResource, error) {

		return []evmodel.SubscriptionResource{
			{
				EventDestination: &model.EventDestination{
					Destination:          "https://localhost:9090/events",
					Name:                 "Subscription",
					Context:              "context",
					EventTypes:           []string{"Alert"},
					MessageIds:           []string{},
					ResourceTypes:        []string{},
					OriginResources:      []model.Link{model.Link{Oid: ""}},
					SubordinateResources: true,
				},
				SubscriptionID: "71de0110-c35a-4859-984c-072d6c5a32d9",
				Hosts:          []string{"localhost"},
			},
		}, nil
	}
	pc.DeleteAggregateSubscriptions(evcommon.MockContext(), &req, true)
}

func TestExternalInterfaces_resubscribeFabricsSubscription(t *testing.T) {
	config.SetUpMockConfig(t)
	pc := getMockMethods()
	event := model.EventDestination{}
	err := pc.resubscribeFabricsSubscription(evcommon.MockContext(), event, "/fabric/scascascsaa", false)
	assert.Nil(t, err)
	err = pc.resubscribeFabricsSubscription(evcommon.MockContext(), event, "/redfish/v1/Fabrics/6d4a0a66-7efa-578e-83cf-44dc68d2874e", false)
	assert.Nil(t, err)
}

func TestExternalInterfaces_subscribe(t *testing.T) {
	config.SetUpMockConfig(t)
	pc := getMockMethods()
	event := model.EventDestination{}

	err := pc.subscribe(evcommon.MockContext(), event, "/redfish/v1/Systems/6d4a0a66-7efa-578e-83cf-44dc68d2874e.1", false, "valid")
	assert.Nil(t, err)

}
