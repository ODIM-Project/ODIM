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
	"errors"
	"reflect"
	"testing"

	eventsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/events"
	"google.golang.org/grpc"
)

func TestDoGetEventService(t *testing.T) {
	type args struct {
		req eventsproto.EventSubRequest
	}
	tests := []struct {
		name                string
		args                args
		ClientFunc          func(clientName string) (*grpc.ClientConn, error)
		NewEventsClientFunc func(cc *grpc.ClientConn) eventsproto.EventsClient
		want                *eventsproto.EventSubResponse
		wantErr             bool
	}{
		{
			name:                "Client func error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewEventsClientFunc: func(cc *grpc.ClientConn) eventsproto.EventsClient { return nil },
			want:                nil,
			wantErr:             true,
		},
		{
			name:                "GetEventService error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewEventsClientFunc: func(cc *grpc.ClientConn) eventsproto.EventsClient { return fakeStruct{} },
			want:                nil,
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewEventsClientFunc = tt.NewEventsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetEventService(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetEventService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetEventService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoCreateEventSubscription(t *testing.T) {
	type args struct {
		req eventsproto.EventSubRequest
	}
	tests := []struct {
		name                string
		args                args
		ClientFunc          func(clientName string) (*grpc.ClientConn, error)
		NewEventsClientFunc func(cc *grpc.ClientConn) eventsproto.EventsClient
		want                *eventsproto.EventSubResponse
		wantErr             bool
	}{
		{
			name:                "Client func error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewEventsClientFunc: func(cc *grpc.ClientConn) eventsproto.EventsClient { return nil },
			want:                nil,
			wantErr:             true,
		},
		{
			name:                "CreateEventSubscription error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewEventsClientFunc: func(cc *grpc.ClientConn) eventsproto.EventsClient { return fakeStruct{} },
			want:                nil,
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewEventsClientFunc = tt.NewEventsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoCreateEventSubscription(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoCreateEventSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoCreateEventSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoSubmitTestEvent(t *testing.T) {
	type args struct {
		req eventsproto.EventSubRequest
	}
	tests := []struct {
		name                string
		args                args
		ClientFunc          func(clientName string) (*grpc.ClientConn, error)
		NewEventsClientFunc func(cc *grpc.ClientConn) eventsproto.EventsClient
		want                *eventsproto.EventSubResponse
		wantErr             bool
	}{
		{
			name:                "Client func error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewEventsClientFunc: func(cc *grpc.ClientConn) eventsproto.EventsClient { return nil },
			want:                nil,
			wantErr:             true,
		},
		{
			name:                "SubmitTestEvent error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewEventsClientFunc: func(cc *grpc.ClientConn) eventsproto.EventsClient { return fakeStruct{} },
			want:                nil,
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewEventsClientFunc = tt.NewEventsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoSubmitTestEvent(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoSubmitTestEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoSubmitTestEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetEventSubscription(t *testing.T) {
	type args struct {
		req eventsproto.EventRequest
	}
	tests := []struct {
		name                string
		args                args
		ClientFunc          func(clientName string) (*grpc.ClientConn, error)
		NewEventsClientFunc func(cc *grpc.ClientConn) eventsproto.EventsClient
		want                *eventsproto.EventSubResponse
		wantErr             bool
	}{
		{
			name:                "Client func error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewEventsClientFunc: func(cc *grpc.ClientConn) eventsproto.EventsClient { return nil },
			want:                nil,
			wantErr:             true,
		},
		{
			name:                "GetEventSubscription error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewEventsClientFunc: func(cc *grpc.ClientConn) eventsproto.EventsClient { return fakeStruct{} },
			want:                nil,
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewEventsClientFunc = tt.NewEventsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetEventSubscription(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetEventSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetEventSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoDeleteEventSubscription(t *testing.T) {
	type args struct {
		req eventsproto.EventRequest
	}
	tests := []struct {
		name                string
		args                args
		ClientFunc          func(clientName string) (*grpc.ClientConn, error)
		NewEventsClientFunc func(cc *grpc.ClientConn) eventsproto.EventsClient
		want                *eventsproto.EventSubResponse
		wantErr             bool
	}{
		{
			name:                "Client func error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewEventsClientFunc: func(cc *grpc.ClientConn) eventsproto.EventsClient { return nil },
			want:                nil,
			wantErr:             true,
		},
		{
			name:                "DeleteEventSubscription error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewEventsClientFunc: func(cc *grpc.ClientConn) eventsproto.EventsClient { return fakeStruct{} },
			want:                nil,
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewEventsClientFunc = tt.NewEventsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoDeleteEventSubscription(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoDeleteEventSubscription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoDeleteEventSubscription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetEventSubscriptionsCollection(t *testing.T) {
	type args struct {
		req eventsproto.EventRequest
	}
	tests := []struct {
		name                string
		args                args
		ClientFunc          func(clientName string) (*grpc.ClientConn, error)
		NewEventsClientFunc func(cc *grpc.ClientConn) eventsproto.EventsClient
		want                *eventsproto.EventSubResponse
		wantErr             bool
	}{
		{
			name:                "Client func error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewEventsClientFunc: func(cc *grpc.ClientConn) eventsproto.EventsClient { return nil },
			want:                nil,
			wantErr:             true,
		},
		{
			name:                "GetEventSubscriptionsCollection error",
			args:                args{},
			ClientFunc:          func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewEventsClientFunc: func(cc *grpc.ClientConn) eventsproto.EventsClient { return fakeStruct{} },
			want:                nil,
			wantErr:             true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewEventsClientFunc = tt.NewEventsClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetEventSubscriptionsCollection(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetEventSubscriptionsCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetEventSubscriptionsCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}
