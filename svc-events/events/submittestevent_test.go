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
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/stretchr/testify/assert"
)

// Positive test cases
func TestSubmitTestEvent(t *testing.T) {
	config.SetUpMockConfig(t)
	p := getMockMethods()
	event := map[string]interface{}{
		"MemberID":          "1",
		"EventType":         "Alert",
		"EventID":           "123",
		"EventGroupID":      1,
		"MessageArgs":       []string{"message"},
		"Severity":          "OK",
		"EventTimestamp":    "",
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
	resp := p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusOK")

	// Negative - with invalid request body
	req = &eventsproto.EventSubRequest{
		SessionToken: "validToken",
		PostBody:     []byte{},
	}
	resp = p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// Invalid token
	req = &eventsproto.EventSubRequest{
		SessionToken: "invalidtoken",
		PostBody:     message,
	}
	resp = p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status Code should be StatusUnauthorized")

	// test case for get session user name
	req = &eventsproto.EventSubRequest{
		SessionToken: "InvalidToken",
		PostBody:     message,
	}
	resp = p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status Code should be StatusUnauthorized")

	// test case for invalid user name
	req = &eventsproto.EventSubRequest{
		SessionToken: "validToken",
		PostBody:     message,
	}
	resp = p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status Code should be StatusInternalServerError")

	// with invalid messageId type
	req = &eventsproto.EventSubRequest{
		SessionToken: "validToken",
		PostBody:     []byte(`{"MessageId": 123}`),
	}
	resp = p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// with invalid EventGroupId type
	req = &eventsproto.EventSubRequest{
		SessionToken: "validToken",
		PostBody:     []byte(`{"MessageId": "123", "EventGroupId": "aasdd"}`),
	}
	resp = p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// with invalid EventId type
	req = &eventsproto.EventSubRequest{
		SessionToken: "validToken",
		PostBody:     []byte(`{"MessageId": "123", "EventId": 123}`),
	}
	resp = p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// with invalid EventTimestamp type
	req = &eventsproto.EventSubRequest{
		SessionToken: "validToken",
		PostBody:     []byte(`{"MessageId": "123", "EventTimestamp": 123}`),
	}
	resp = p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// with invalid EventType type
	req = &eventsproto.EventSubRequest{
		SessionToken: "validToken",
		PostBody:     []byte(`{"MessageId": "123", "EventType": 123}`),
	}
	resp = p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// with invalid EventType data
	req = &eventsproto.EventSubRequest{
		SessionToken: "validToken",
		PostBody:     []byte(`{"MessageId": "123", "EventType": "123"}`),
	}
	resp = p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// with invalid Message type
	req = &eventsproto.EventSubRequest{
		SessionToken: "validToken",
		PostBody:     []byte(`{"MessageId": "123", "Message": 123}`),
	}
	resp = p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// with invalid MessageArgs type
	req = &eventsproto.EventSubRequest{
		SessionToken: "validToken",
		PostBody:     []byte(`{"MessageId": "123", "MessageArgs": "123"}`),
	}
	resp = p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// with invalid OriginOfCondition type
	req = &eventsproto.EventSubRequest{
		SessionToken: "validToken",
		PostBody:     []byte(`{"MessageId": "123", "OriginOfCondition": 123}`),
	}
	resp = p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// with invalid Severity type
	req = &eventsproto.EventSubRequest{
		SessionToken: "validToken",
		PostBody:     []byte(`{"MessageId": "123", "Severity": 123}`),
	}
	resp = p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")

	// with invalid Severity value
	req = &eventsproto.EventSubRequest{
		SessionToken: "validToken",
		PostBody:     []byte(`{"MessageId": "123", "Severity": "123"}`),
	}
	resp = p.SubmitTestEvent(req)
	assert.Equal(t, http.StatusBadRequest, int(resp.StatusCode), "Status Code should be StatusBadRequest")
}
