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
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/ODIM-Project/ODIM/svc-events/evresponse"
)

// GetEventSubscriptionsDetails collects subscription data against given subscription id
func (e *ExternalInterfaces) GetEventSubscriptionsDetails(ctx context.Context, req *eventsproto.EventRequest) response.RPC {
	var resp response.RPC
	authResp, err := e.Auth(req.SessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("error while trying to authenticate session: status code: %v, status message: %v",
			authResp.StatusCode, authResp.StatusMessage)
		if err != nil {
			errMsg = errMsg + ": " + err.Error()
		}
		l.LogWithFields(ctx).Error(errMsg)
		return authResp
	}
	var subscriptions *evresponse.SubscriptionResponse

	subscriptionDetails, err := e.GetEvtSubscriptions(req.EventSubscriptionID)
	if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
		errorMessage := err.Error()
		return common.GeneralError(http.StatusBadRequest, response.ResourceNotFound,
			errorMessage, []interface{}{"EventSubscription", req.EventSubscriptionID}, nil)
	}
	if len(subscriptionDetails) == 0 {
		errorMessage := fmt.Sprintf("Subscription details not found for ID: %v", req.EventSubscriptionID)
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage,
			[]interface{}{"EventSubscription", req.EventSubscriptionID}, nil)
	}

	for _, evtSubscription := range subscriptionDetails {

		// Since we are searching subscription id with pattern search
		// we need to match the subscription id
		if evtSubscription.SubscriptionID != req.EventSubscriptionID {
			errorMessage := fmt.Sprintf("Subscription details not found for subscription id: %s", req.EventSubscriptionID)
			l.LogWithFields(ctx).Info(errorMessage)
			var msgArgs = []interface{}{"SubscriptionID", req.EventSubscriptionID}
			evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
			return resp
		}
		commonResponse := response.Response{
			OdataType:    common.EventDestinationType,
			ID:           evtSubscription.SubscriptionID,
			Name:         evtSubscription.EventDestination.Name,
			OdataContext: "/redfish/v1/$metadata#EventDestination.EventDestination",
			OdataID:      "/redfish/v1/EventService/Subscriptions/" + evtSubscription.SubscriptionID,
		}

		subscriptions = &evresponse.SubscriptionResponse{
			Response:            commonResponse,
			Destination:         evtSubscription.EventDestination.Destination,
			Protocol:            evtSubscription.EventDestination.Protocol,
			Context:             evtSubscription.EventDestination.Context,
			EventTypes:          evtSubscription.EventDestination.EventTypes,
			SubscriptionType:    evtSubscription.EventDestination.SubscriptionType,
			MessageIds:          evtSubscription.EventDestination.MessageIds,
			ResourceTypes:       evtSubscription.EventDestination.ResourceTypes,
			OriginResources:     evtSubscription.EventDestination.OriginResources,
			DeliveryRetryPolicy: evtSubscription.EventDestination.DeliveryRetryPolicy,
		}
	}
	resp.Body = subscriptions
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}

// GetEventSubscriptionsCollection collects all subscription details
func (e *ExternalInterfaces) GetEventSubscriptionsCollection(ctx context.Context, req *eventsproto.EventRequest) response.RPC {
	var resp response.RPC
	authResp, err := e.Auth(req.SessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("error while trying to authenticate session: status code: %v, status message: %v",
			authResp.StatusCode, authResp.StatusMessage)
		if err != nil {
			errMsg = errMsg + ": " + err.Error()
		}
		l.LogWithFields(ctx).Error(errMsg)
		return authResp
	}
	listMembers := []model.Link{}
	searchKey := "*"

	subscriptionDetails, err := e.GetEvtSubscriptions(searchKey)
	if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
		l.LogWithFields(ctx).Printf("error getting event subscription details : %v", err)
		errorMessage := err.Error()
		return common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection,
			errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil)
	}
	for _, evtSubscription := range subscriptionDetails {
		subscriptionID := evtSubscription.SubscriptionID
		destination := evtSubscription.EventDestination.Destination
		if destination == "" {
			continue
		}
		member := model.Link{
			Oid: "/redfish/v1/EventService/Subscriptions/" + subscriptionID + "/",
		}

		listMembers = append(listMembers, member)
	}

	eventResp := evresponse.ListResponse{
		OdataContext: "/redfish/v1/$metadata#EventDestinationCollection.EventDestinationCollection",
		OdataID:      "/redfish/v1/EventService/Subscriptions",
		OdataType:    "#EventDestinationCollection.EventDestinationCollection",
		Name:         "EventSubscriptions",
		Description:  "Event Subscriptions",
		MembersCount: len(listMembers),
		Members:      listMembers,
	}
	resp.Body = eventResp
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}

// IsAggregateHaveSubscription validate any subscription contain aggregate id, return status
func (e *ExternalInterfaces) IsAggregateHaveSubscription(ctx context.Context, req *eventsproto.EventUpdateRequest) bool {
	authResp, err := e.Auth(req.SessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("error while trying to authenticate session: status code: %v, status message: %v",
			authResp.StatusCode, authResp.StatusMessage)
		if err != nil {
			errMsg = errMsg + ": " + err.Error()
		}
		l.LogWithFields(ctx).Error(errMsg)
		return false
	}
	searchKey := evcommon.GetSearchKey(req.AggregateId, evmodel.SubscriptionIndex)
	subscriptionDetails, err := e.GetEvtSubscriptions(searchKey)
	if err != nil {
		return false
	}
	return len(subscriptionDetails) > 0
}
