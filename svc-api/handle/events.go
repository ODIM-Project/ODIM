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

// Package handle ...
package handle

import (
	"context"
	"encoding/json"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	iris "github.com/kataras/iris/v12"
)

// EventsRPCs defines all the RPC methods in Events service
type EventsRPCs struct {
	GetEventServiceRPC                 func(context.Context, eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error)
	CreateEventSubscriptionRPC         func(context.Context, eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error)
	SubmitTestEventRPC                 func(context.Context, eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error)
	GetEventSubscriptionRPC            func(context.Context, eventsproto.EventRequest) (*eventsproto.EventSubResponse, error)
	DeleteEventSubscriptionRPC         func(context.Context, eventsproto.EventRequest) (*eventsproto.EventSubResponse, error)
	GetEventSubscriptionsCollectionRPC func(context.Context, eventsproto.EventRequest) (*eventsproto.EventSubResponse, error)
}

// GetEventService is the handler to get the Event Service details.
func (e *EventsRPCs) GetEventService(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	l.LogWithFields(ctxt).Debug("Incoming request received for the Get Event service")
	req := eventsproto.EventSubRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
	}
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := e.GetEventServiceRPC(ctxt, req)
	if err != nil {
		l.LogWithFields(ctxt).Error(err.Error())
		common.SendFailedRPCCallResponse(ctx, err.Error())
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for Getting Event service is %s and response status %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// CreateEventSubscription is the handler for creating event subscription
func (e *EventsRPCs) CreateEventSubscription(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req eventsproto.EventSubRequest
	// Read Post Body from Request
	var SubscriptionReq interface{}
	err := ctx.ReadJSON(&SubscriptionReq)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the event subscription request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}
	req.SessionToken = ctx.Request().Header.Get(AuthTokenHeader)

	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	req.PostBody, _ = json.Marshal(&SubscriptionReq)
	l.LogWithFields(ctxt).Debugf("Incoming request received for creating event subscription with request body %s", string(req.PostBody))
	resp, err := e.CreateEventSubscriptionRPC(ctxt, req)
	if err != nil {
		l.LogWithFields(ctxt).Error(err.Error())
		common.SendFailedRPCCallResponse(ctx, err.Error())
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for creating event subscription is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// SubmitTestEvent is the handler to submit test event
func (e *EventsRPCs) SubmitTestEvent(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req eventsproto.EventSubRequest
	// Read Post Body from Request
	var SubmitTestEventReq interface{}
	err := ctx.ReadJSON(&SubmitTestEventReq)
	if err != nil {
		errorMessage := "error while trying to get JSON body from the SubmitTestEvent request body: " + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendMalformedJSONRequestErrResponse(ctx, errorMessage)
	}

	req.SessionToken = ctx.Request().Header.Get(AuthTokenHeader)

	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	req.PostBody, _ = json.Marshal(&SubmitTestEventReq)
	l.LogWithFields(ctxt).Debugf("Incoming request received for submit test event with request body %s", string(req.PostBody))
	resp, err := e.SubmitTestEventRPC(ctxt, req)
	if err != nil {
		l.LogWithFields(ctxt).Error(err.Error())
		common.SendFailedRPCCallResponse(ctx, err.Error())
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for submit test event is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetEventSubscription is the handler for getting event subscription
func (e *EventsRPCs) GetEventSubscription(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req eventsproto.EventRequest
	req.EventSubscriptionID = ctx.Params().Get("id")
	req.SessionToken = ctx.Request().Header.Get(AuthTokenHeader)
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting event subscription with id %s", req.EventSubscriptionID)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := e.GetEventSubscriptionRPC(ctxt, req)
	if err != nil {
		l.LogWithFields(ctxt).Error(err.Error())
		common.SendFailedRPCCallResponse(ctx, err.Error())
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting event subscription is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET, DELETE")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// DeleteEventSubscription is the handler for getting event subscription
func (e *EventsRPCs) DeleteEventSubscription(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req eventsproto.EventRequest
	req.EventSubscriptionID = ctx.Params().Get("id")
	req.SessionToken = ctx.Request().Header.Get(AuthTokenHeader)
	l.LogWithFields(ctxt).Debugf("Incoming request received for deleting event subscription with id %s", req.EventSubscriptionID)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := e.DeleteEventSubscriptionRPC(ctxt, req)
	if err != nil {
		l.LogWithFields(ctxt).Error(err.Error())
		common.SendFailedRPCCallResponse(ctx, err.Error())
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for deleting event subscription is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}

// GetEventSubscriptionsCollection is the handler for getting event subscriptions collection
func (e *EventsRPCs) GetEventSubscriptionsCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	var req eventsproto.EventRequest
	req.SessionToken = ctx.Request().Header.Get(AuthTokenHeader)
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting all the  event subscription collections")
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}

	resp, err := e.GetEventSubscriptionsCollectionRPC(ctxt, req)
	if err != nil {
		l.LogWithFields(ctxt).Error(err.Error())
		common.SendFailedRPCCallResponse(ctx, err.Error())
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting event subscription collections is %s with response code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)
}
