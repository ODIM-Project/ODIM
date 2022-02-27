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

//Package handle ...
package handle

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	iris "github.com/kataras/iris/v12"
)

// EventsRPCs defines all the RPC methods in Events service
type EventsRPCs struct {
	GetEventServiceRPC                 func(eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error)
	CreateEventSubscriptionRPC         func(eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error)
	SubmitTestEventRPC                 func(eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error)
	GetEventSubscriptionRPC            func(eventsproto.EventRequest) (*eventsproto.EventSubResponse, error)
	DeleteEventSubscriptionRPC         func(eventsproto.EventRequest) (*eventsproto.EventSubResponse, error)
	GetEventSubscriptionsCollectionRPC func(eventsproto.EventRequest) (*eventsproto.EventSubResponse, error)
}

// GetEventService is the handler to get the Event Service details.
func (e *EventsRPCs) GetEventService(ctx iris.Context) {
	defer ctx.Next()
	req := eventsproto.EventSubRequest{
		SessionToken: ctx.Request().Header.Get("X-Auth-Token"),
	}
	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	resp, err := e.GetEventServiceRPC(req)
	if err != nil {
		log.Error(err.Error())
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// CreateEventSubscription is the handler for creating event subscription
func (e *EventsRPCs) CreateEventSubscription(ctx iris.Context) {
	defer ctx.Next()
	var req eventsproto.EventSubRequest
	// Read Post Body from Request
	var SubscriptionReq interface{}
	err := ctx.ReadJSON(&SubscriptionReq)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the event subscription request body: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}

	req.SessionToken = ctx.Request().Header.Get("X-Auth-Token")

	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	req.PostBody, _ = json.Marshal(&SubscriptionReq)

	resp, err := e.CreateEventSubscriptionRPC(req)
	if err != nil {
		log.Error(err.Error())
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// SubmitTestEvent is the handler to submit test event
func (e *EventsRPCs) SubmitTestEvent(ctx iris.Context) {
	defer ctx.Next()
	var req eventsproto.EventSubRequest
	// Read Post Body from Request
	var SubmitTestEventReq interface{}
	err := ctx.ReadJSON(&SubmitTestEventReq)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the SubmitTestEvent request body: " + err.Error()
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusBadRequest, response.MalformedJSON, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(&response.Body)
		return
	}

	req.SessionToken = ctx.Request().Header.Get("X-Auth-Token")

	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}
	req.PostBody, _ = json.Marshal(&SubmitTestEventReq)

	resp, err := e.SubmitTestEventRPC(req)
	if err != nil {
		log.Error(err.Error())
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetEventSubscription is the handler for getting event subscription
func (e *EventsRPCs) GetEventSubscription(ctx iris.Context) {
	defer ctx.Next()
	var req eventsproto.EventRequest
	req.EventSubscriptionID = ctx.Params().Get("id")
	req.SessionToken = ctx.Request().Header.Get("X-Auth-Token")

	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	resp, err := e.GetEventSubscriptionRPC(req)
	if err != nil {
		log.Error(err.Error())
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response)
		return
	}
	ctx.ResponseWriter().Header().Set("Allow", "GET, DELETE")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// DeleteEventSubscription is the handler for getting event subscription
func (e *EventsRPCs) DeleteEventSubscription(ctx iris.Context) {
	defer ctx.Next()
	var req eventsproto.EventRequest
	req.EventSubscriptionID = ctx.Params().Get("id")
	req.SessionToken = ctx.Request().Header.Get("X-Auth-Token")

	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	resp, err := e.DeleteEventSubscriptionRPC(req)
	if err != nil {
		log.Error(err.Error())
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetEventSubscriptionsCollection is the handler for getting event subscriptions collection
func (e *EventsRPCs) GetEventSubscriptionsCollection(ctx iris.Context) {
	defer ctx.Next()
	var req eventsproto.EventRequest
	req.SessionToken = ctx.Request().Header.Get("X-Auth-Token")

	if req.SessionToken == "" {
		errorMessage := "no X-Auth-Token found in request header"
		response := common.GeneralError(http.StatusUnauthorized, response.NoValidSession, errorMessage, nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(&response.Body)
		return
	}

	resp, err := e.GetEventSubscriptionsCollectionRPC(req)
	if err != nil {
		log.Error(err.Error())
		response := common.GeneralError(http.StatusInternalServerError, response.InternalError, err.Error(), nil, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(&response.Body)
		return
	}

	ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}
