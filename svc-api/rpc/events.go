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

//Package rpc ...
package rpc

import (
	"context"
	"fmt"

	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

// DoGetEventService defines the RPC call function for
// the GetEventService from events micro service
func DoGetEventService(req eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error) {
	conn, err := services.ODIMService.Client(services.Events)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	events := eventsproto.NewEventsClient(conn)

	resp, err := events.GetEventService(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoCreateEventSubscription defines the RPC call function for
// the CreateEventSubscription from events micro service
func DoCreateEventSubscription(req eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error) {

	conn, err := services.ODIMService.Client(services.Events)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	events := eventsproto.NewEventsClient(conn)

	resp, err := events.CreateEventSubscription(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoSubmitTestEvent defines the RPC call function for
// the SubmitTestEvent from events micro service
func DoSubmitTestEvent(req eventsproto.EventSubRequest) (*eventsproto.EventSubResponse, error) {

	conn, err := services.ODIMService.Client(services.Events)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	events := eventsproto.NewEventsClient(conn)

	resp, err := events.SubmitTestEvent(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetEventSubscription defines the RPC call function for
// the DoGetEventSubscription from events micro service
func DoGetEventSubscription(req eventsproto.EventRequest) (*eventsproto.EventSubResponse, error) {

	conn, err := services.ODIMService.Client(services.Events)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	events := eventsproto.NewEventsClient(conn)

	resp, err := events.GetEventSubscription(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoDeleteEventSubscription defines the RPC call function for
// the DoDeleteEventSubscription from events micro service
func DoDeleteEventSubscription(req eventsproto.EventRequest) (*eventsproto.EventSubResponse, error) {

	conn, err := services.ODIMService.Client(services.Events)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	events := eventsproto.NewEventsClient(conn)

	resp, err := events.DeleteEventSubscription(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetEventSubscriptionsCollection defines the RPC call function for
// the DoGetEventSubscription from events micro service
func DoGetEventSubscriptionsCollection(req eventsproto.EventRequest) (*eventsproto.EventSubResponse, error) {

	conn, err := services.ODIMService.Client(services.Events)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	events := eventsproto.NewEventsClient(conn)

	resp, err := events.GetEventSubscriptionsCollection(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}
