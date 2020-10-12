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
// - Post TestEvent (SubmitTestEvent)
// and corresponding unit test cases
package events

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
)

// SubmitTestEvent is a helper method to handle the submit test event request.
func (p *PluginContact) SubmitTestEvent(req *eventsproto.EventSubRequest) response.RPC {
	var resp response.RPC
	authStatusCode, authStatusMessage := p.Auth(
		req.SessionToken,
		[]string{
			common.PrivilegeConfigureComponents,
		},
		[]string{},
	)
	if authStatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("error while trying to authenticate session for submit test events: status code: %v, status message: %v", authStatusCode, authStatusMessage)
		log.Printf(errMsg)
		return common.GeneralError(authStatusCode, authStatusMessage, errMsg, nil, nil)
	}
	// First get the UserName from SessionToken
	sessionUserName, err := p.GetSessionUserName(req.SessionToken)
	if err != nil {
		// handle the error case with appropriate response body
		errMsg := "error while trying to authenticate session: " + err.Error()
		log.Printf(errMsg)
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errMsg, nil, nil)
	}

	testEvent, statusMessage, errMsg, msgArgs := validAndGenSubTestReq(req.PostBody)
	if statusMessage != response.Success {
		log.Printf(errMsg)
		return common.GeneralError(http.StatusBadRequest, statusMessage, errMsg, msgArgs, nil)
	}

	// parsing the event
	var eventObj interface{}
	err = json.Unmarshal(req.PostBody, &eventObj)
	if err != nil {
		errMsg := "unable to parse the event request" + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}

	// Validating the request JSON properties for case sensitive
	invalidProperties, err := common.RequestParamsCaseValidator(req.PostBody, eventObj)
	if err != nil {
		errMsg := "error while validating request parameters: " + err.Error()
		log.Println(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	} else if invalidProperties != "" {
		errorMessage := "error: one or more properties given in the request body are not valid, ensure properties are listed in uppercamelcase "
		log.Println(errorMessage)
		resp := common.GeneralError(http.StatusBadRequest, response.PropertyUnknown, errorMessage, []interface{}{invalidProperties}, nil)
		return resp
	}

	// Find out all the subscription destinations of the requesting user
	subscriptions, err := evmodel.GetEvtSubscriptions(sessionUserName)
	if err != nil {
		// Internall error
		errMsg := "error while trying to find the event destination"
		log.Printf(errMsg)
		return common.GeneralError(http.StatusInternalServerError, response.InternalError, errMsg, nil, nil)
	}
	// we need common.MessageData to find the correct destination to send test event
	var message common.MessageData
	message.Events = append(message.Events, *testEvent)
	messageBytes, _ := json.Marshal(message)
	for _, sub := range subscriptions {

		for _, origin := range sub.OriginResources {
			if sub.Destination != "" {
				if filterEventsToBeForwarded(sub, messageBytes, []string{origin}) {
					log.Printf("Destination: %v\n", sub.Destination)
					go postEvent(sub.Destination, messageBytes)
				}
			}
		}
	}

	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	resp.Body = response.ErrorClass{
		Code:    resp.StatusMessage,
		Message: "Request completed successfully.",
	}
	return resp

}

func validAndGenSubTestReq(reqBody []byte) (*common.Event, string, string, []interface{}) {
	var testEvent common.Event
	var req map[string]interface{}
	json.Unmarshal(reqBody, &req)
	if val, ok := req["MessageId"]; ok {
		switch v := val.(type) {
		case string:
			testEvent.MessageID = v
		default:
			return nil, response.PropertyValueTypeError, "error: required parameter MessageId must be of type string", []interface{}{fmt.Sprintf("%v", v), "MessageId"}
		}
	} else {
		return nil, response.PropertyMissing, "error: MessageId is a required parameter", []interface{}{"MessageId"}
	}

	if val, ok := req["EventGroupId"]; ok {
		switch v := val.(type) {
		case int:
			testEvent.EventGroupID = v
		case float64:
			testEvent.EventGroupID = int(v)
		default:
			return nil, response.PropertyValueTypeError, "error: optional parameter EventGroupId must be of type integer", []interface{}{fmt.Sprintf("%v", v), "EventGroupId"}
		}
	}

	if val, ok := req["EventId"]; ok {
		switch v := val.(type) {
		case string:
			testEvent.EventID = v
		default:
			return nil, response.PropertyValueTypeError, "error: optional parameter EventId must be of type string", []interface{}{fmt.Sprintf("%v", v), "EventId"}
		}
	}

	if val, ok := req["EventTimestamp"]; ok {
		switch v := val.(type) {
		case string:
			testEvent.EventTimestamp = v
		default:
			return nil, response.PropertyValueTypeError, "error: optional parameter EventTimestamp must be of type string", []interface{}{fmt.Sprintf("%v", v), "EventTimestamp"}
		}
	}

	if val, ok := req["EventType"]; ok {
		switch v := val.(type) {
		case string:
			if ok = validEventType(v); ok {
				testEvent.EventType = v
			} else {
				return nil, response.PropertyValueNotInList, "error: optional parameter EventType must have allowed value", []interface{}{fmt.Sprintf("%v", v), "EventType"}
			}
		default:
			return nil, response.PropertyValueTypeError, "error: optional parameter EventType must be of type string", []interface{}{fmt.Sprintf("%v", v), "EventType"}
		}
	}

	if val, ok := req["Message"]; ok {
		switch v := val.(type) {
		case string:
			testEvent.Message = v
		default:
			return nil, response.PropertyValueTypeError, "error: optional parameter Message must be of type string", []interface{}{fmt.Sprintf("%v", v), "Message"}
		}
	}

	if val, ok := req["MessageArgs"]; ok {
		switch v := val.(type) {
		case []string:
			testEvent.MessageArgs = v
		case []interface{}:
			msg, _ := json.Marshal(v)
			var msgArgs []string
			json.Unmarshal(msg, &msgArgs)
			testEvent.MessageArgs = msgArgs
		default:
			return nil, response.PropertyValueTypeError, "error: optional parameter MessageArgs must be of type array(string)", []interface{}{fmt.Sprintf("%v", v), "MessageArgs"}
		}
	}

	if val, ok := req["OriginOfCondition"]; ok {
		switch v := val.(type) {
		case string:
			testEvent.OriginOfCondition = &common.Link{
				Oid: v,
			}
		default:
			return nil, response.PropertyValueTypeError, "error: optional parameter OriginOfCondition must be of type string", []interface{}{fmt.Sprintf("%v", v), "OriginOfCondition"}
		}
	}

	if val, ok := req["Severity"]; ok {
		switch v := val.(type) {
		case string:
			if ok = validSeverity(v); ok {
				testEvent.Severity = v
			} else {
				return nil, response.PropertyValueNotInList, "error: optional parameter Severity must have allowed value", []interface{}{fmt.Sprintf("%v", v), "Severity"}
			}
		default:
			return nil, response.PropertyValueTypeError, "error: optional parameter Severity must be of type string", []interface{}{fmt.Sprintf("%v", v), "Severity"}
		}
	}

	return &testEvent, response.Success, "", nil
}

func validEventType(got string) bool {
	events := getAllowedEventTypes()
	for _, event := range events {
		if event == got {
			return true
		}
	}
	return false
}

func validSeverity(got string) bool {
	severities := getAllowedSeverities()
	for _, severity := range severities {
		if severity == got {
			return true
		}
	}
	return false
}

func getAllowedEventTypes() []string {
	return []string{
		"Alert",
		"MetricReport",
		"Other",
		"ResourceAdded",
		"ResourceRemoved",
		"ResourceUpdated",
		"StatusChange",
	}
}

func getAllowedSeverities() []string {
	return []string{
		"Critical",
		"OK",
		"Warning",
	}
}
