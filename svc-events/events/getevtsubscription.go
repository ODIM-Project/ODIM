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
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-events/evcommon"
	"github.com/ODIM-Project/ODIM/svc-events/evmodel"
	"github.com/ODIM-Project/ODIM/svc-events/evresponse"
)

// GetEventSubscriptionsDetails collects subscription data against given subscription id
func (p *PluginContact) GetEventSubscriptionsDetails(req *eventsproto.EventRequest) response.RPC {
	var resp response.RPC
	authResp := p.Auth(req.SessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Printf("error while trying to authenticate session: status code: %v, status message: %v", authResp.StatusCode, authResp.StatusMessage)
		return authResp
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	var subscriptions *evresponse.SubscriptionResponse

	subscriptionDetails, err := evmodel.GetEvtSubscriptions(req.EventSubscriptionID)
	if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
		log.Printf("error getting eventsubscription details : %v", err)
		errorMessage := err.Error()
		return common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, errorMessage, []interface{}{"EventSubscription", req.EventSubscriptionID}, nil)
	}
	if len(subscriptionDetails) < 1 {
		log.Printf("Subscription details not found for ID: %v", req.EventSubscriptionID)
		errorMessage := fmt.Sprintf("Subscription details not found for ID: %v", req.EventSubscriptionID)
		return common.GeneralError(http.StatusNotFound, response.ResourceNotFound, errorMessage, []interface{}{"EventSubscription", req.EventSubscriptionID}, nil)
	}

	for _, evtSubscription := range subscriptionDetails {

		// Since we are searching subscription id with pattern search
		// we need to match the subscripton id
		if evtSubscription.SubscriptionID != req.EventSubscriptionID {
			errorMessage := fmt.Sprintf("Subscription details not found for subscription id: %s", req.EventSubscriptionID)
			log.Println(errorMessage)
			var msgArgs = []interface{}{"SubscriptionID", req.EventSubscriptionID}
			evcommon.GenErrorResponse(errorMessage, response.ResourceNotFound, http.StatusNotFound, msgArgs, &resp)
			return resp
		}
		// if requested subscription id not matched with stored subscription id then return error
		if req.EventSubscriptionID != evtSubscription.SubscriptionID {
			log.Printf("Subscription details not found for ID: %v", req.EventSubscriptionID)
			errorMessage := fmt.Sprintf("Subscription details not found for ID: %v", req.EventSubscriptionID)
			return common.GeneralError(http.StatusBadRequest, response.ResourceNotFound, errorMessage, []interface{}{"EventSubscription", req.EventSubscriptionID}, nil)
		}
		commonResponse := response.Response{
			OdataType:    common.EventDestinationType,
			ID:           evtSubscription.SubscriptionID,
			Name:         evtSubscription.Name,
			OdataContext: "/redfish/v1/$metadata#EventDestination.EventDestination",
			OdataID:      "/redfish/v1/EventService/Subscriptions/" + evtSubscription.SubscriptionID,
		}

		subscriptions = &evresponse.SubscriptionResponse{
			Response:         commonResponse,
			Destination:      evtSubscription.Destination,
			Protocol:         evtSubscription.Protocol,
			Context:          evtSubscription.Context,
			EventTypes:       evtSubscription.EventTypes,
			SubscriptionType: evtSubscription.SubscriptionType,
			MessageIds:       evtSubscription.MessageIds,
			ResourceTypes:    evtSubscription.ResourceTypes,
			OriginResources:  updateOriginResourceswithOdataID(evtSubscription.OriginResources),
		}
	}
	resp.Body = subscriptions
	resp.StatusCode = http.StatusOK
	resp.StatusMessage = response.Success
	return resp
}

func updateOriginResourceswithOdataID(originResources []string) []evresponse.ListMember {
	var originRes []evresponse.ListMember
	for _, origin := range originResources {
		originRes = append(originRes, evresponse.ListMember{OdataID: origin})
	}
	return originRes
}

// GetEventSubscriptionsCollection collects all subscription details
func (p *PluginContact) GetEventSubscriptionsCollection(req *eventsproto.EventRequest) response.RPC {
	var resp response.RPC
	authResp := p.Auth(req.SessionToken, []string{common.PrivilegeConfigureComponents}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Printf("error while trying to authenticate session: status code: %v, status message: %v", authResp.StatusCode, authResp.StatusMessage)
		return authResp
	}
	resp.Header = map[string]string{
		"Cache-Control":     "no-cache",
		"Connection":        "keep-alive",
		"Content-type":      "application/json; charset=utf-8",
		"Transfer-Encoding": "chunked",
		"OData-Version":     "4.0",
	}
	listMembers := []evresponse.ListMember{}
	searchKey := "*"

	subscriptionDetails, err := evmodel.GetEvtSubscriptions(searchKey)
	if err != nil && !strings.Contains(err.Error(), "No data found for the key") {
		log.Printf("error getting eventsubscription details : %v", err)
		errorMessage := err.Error()
		return common.GeneralError(http.StatusServiceUnavailable, response.CouldNotEstablishConnection, errorMessage, []interface{}{config.Data.DBConf.InMemoryHost + ":" + config.Data.DBConf.InMemoryPort}, nil)
	}
	for _, evtSubscription := range subscriptionDetails {
		subscriptionID := evtSubscription.SubscriptionID
		destination := evtSubscription.Destination
		if destination == "" {
			continue
		}
		member := evresponse.ListMember{
			OdataID: "/redfish/v1/EventService/Subscriptions/" + subscriptionID + "/",
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
