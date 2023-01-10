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

// Package rpc ...
package rpc

import (
	"context"
	"fmt"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

var (
	NewEventsClientFunc = eventsproto.NewEventsClient
)

// DoGetEventService defines the RPC call function for
// the GetEventService from events micro service
func DoGetEventService(ctx context.Context, req eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Events)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	events := NewEventsClientFunc(conn)

	resp, err := events.GetEventService(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoCreateEventSubscription defines the RPC call function for
// the CreateEventSubscription from events micro service
func DoCreateEventSubscription(ctx context.Context, req eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Events)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	events := NewEventsClientFunc(conn)

	resp, err := events.CreateEventSubscription(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoSubmitTestEvent defines the RPC call function for
// the SubmitTestEvent from events micro service
func DoSubmitTestEvent(ctx context.Context, req eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Events)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	events := NewEventsClientFunc(conn)

	resp, err := events.SubmitTestEvent(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetEventSubscription defines the RPC call function for
// the DoGetEventSubscription from events micro service
func DoGetEventSubscription(ctx context.Context, req eventsproto.EventRequest) (*eventsproto.EventSubResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Events)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	events := NewEventsClientFunc(conn)

	resp, err := events.GetEventSubscription(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoDeleteEventSubscription defines the RPC call function for
// the DoDeleteEventSubscription from events micro service
func DoDeleteEventSubscription(ctx context.Context, req eventsproto.EventRequest) (*eventsproto.EventSubResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Events)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	events := NewEventsClientFunc(conn)

	resp, err := events.DeleteEventSubscription(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetEventSubscriptionsCollection defines the RPC call function for
// the DoGetEventSubscription from events micro service
func DoGetEventSubscriptionsCollection(ctx context.Context, req eventsproto.EventRequest) (*eventsproto.EventSubResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Events)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	events := NewEventsClientFunc(conn)

	resp, err := events.GetEventSubscriptionsCollection(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}
