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
	"errors"
	"reflect"
	"testing"

	taskproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/task"
	"google.golang.org/grpc"
)

func TestDeleteTaskRequest(t *testing.T) {
	type args struct {
		req *taskproto.GetTaskRequest
	}
	tests := []struct {
		name                        string
		args                        args
		ClientFunc                  func(clientName string) (*grpc.ClientConn, error)
		NewGetTaskServiceClientFunc func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient
		want                        *taskproto.TaskResponse
		wantErr                     bool
	}{
		{
			name:                        "Client func error",
			args:                        args{},
			ClientFunc:                  func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewGetTaskServiceClientFunc: func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient { return nil },
			want:                        nil,
			wantErr:                     true,
		},
		{
			name:                        "DeleteTaskRequest error",
			args:                        args{},
			ClientFunc:                  func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewGetTaskServiceClientFunc: func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient { return fakeStruct{} },
			want:                        &taskproto.TaskResponse{StatusCode: 500, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:                     true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewGetTaskServiceClientFunc = tt.NewGetTaskServiceClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteTaskRequest(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTaskRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteTaskRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTaskRequest(t *testing.T) {
	type args struct {
		req *taskproto.GetTaskRequest
	}
	tests := []struct {
		name                        string
		args                        args
		ClientFunc                  func(clientName string) (*grpc.ClientConn, error)
		NewGetTaskServiceClientFunc func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient
		want                        *taskproto.TaskResponse
		wantErr                     bool
	}{
		{
			name:                        "Client func error",
			args:                        args{},
			ClientFunc:                  func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewGetTaskServiceClientFunc: func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient { return nil },
			want:                        nil,
			wantErr:                     true,
		},
		{
			name:                        "GetTaskRequest error",
			args:                        args{},
			ClientFunc:                  func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewGetTaskServiceClientFunc: func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient { return fakeStruct{} },
			want:                        &taskproto.TaskResponse{StatusCode: 500, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:                     true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewGetTaskServiceClientFunc = tt.NewGetTaskServiceClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTaskRequest(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTaskRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTaskRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSubTasks(t *testing.T) {
	type args struct {
		req *taskproto.GetTaskRequest
	}
	tests := []struct {
		name                        string
		args                        args
		ClientFunc                  func(clientName string) (*grpc.ClientConn, error)
		NewGetTaskServiceClientFunc func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient
		want                        *taskproto.TaskResponse
		wantErr                     bool
	}{
		{
			name:                        "Client func error",
			args:                        args{},
			ClientFunc:                  func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewGetTaskServiceClientFunc: func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient { return nil },
			want:                        nil,
			wantErr:                     true,
		},
		{
			name:                        "GetSubTasks error",
			args:                        args{},
			ClientFunc:                  func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewGetTaskServiceClientFunc: func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient { return fakeStruct{} },
			want:                        &taskproto.TaskResponse{StatusCode: 500, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:                     true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewGetTaskServiceClientFunc = tt.NewGetTaskServiceClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSubTasks(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSubTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSubTask(t *testing.T) {
	type args struct {
		req *taskproto.GetTaskRequest
	}
	tests := []struct {
		name                        string
		args                        args
		ClientFunc                  func(clientName string) (*grpc.ClientConn, error)
		NewGetTaskServiceClientFunc func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient
		want                        *taskproto.TaskResponse
		wantErr                     bool
	}{
		{
			name:                        "Client func error",
			args:                        args{},
			ClientFunc:                  func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewGetTaskServiceClientFunc: func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient { return nil },
			want:                        nil,
			wantErr:                     true,
		},
		{
			name:                        "GetSubTask error",
			args:                        args{},
			ClientFunc:                  func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewGetTaskServiceClientFunc: func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient { return fakeStruct{} },
			want:                        &taskproto.TaskResponse{StatusCode: 500, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:                     true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewGetTaskServiceClientFunc = tt.NewGetTaskServiceClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSubTask(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSubTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTaskMonitor(t *testing.T) {
	type args struct {
		req *taskproto.GetTaskRequest
	}
	tests := []struct {
		name                        string
		args                        args
		ClientFunc                  func(clientName string) (*grpc.ClientConn, error)
		NewGetTaskServiceClientFunc func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient
		want                        *taskproto.TaskResponse
		wantErr                     bool
	}{
		{
			name:                        "Client func error",
			args:                        args{},
			ClientFunc:                  func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewGetTaskServiceClientFunc: func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient { return nil },
			want:                        nil,
			wantErr:                     true,
		},
		{
			name:                        "GetTaskMonitor error",
			args:                        args{},
			ClientFunc:                  func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewGetTaskServiceClientFunc: func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient { return fakeStruct{} },
			want:                        &taskproto.TaskResponse{StatusCode: 500, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:                     true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewGetTaskServiceClientFunc = tt.NewGetTaskServiceClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTaskMonitor(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTaskMonitor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTaskMonitor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskCollection(t *testing.T) {
	type args struct {
		req *taskproto.GetTaskRequest
	}
	tests := []struct {
		name                        string
		args                        args
		ClientFunc                  func(clientName string) (*grpc.ClientConn, error)
		NewGetTaskServiceClientFunc func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient
		want                        *taskproto.TaskResponse
		wantErr                     bool
	}{
		{
			name:                        "Client func error",
			args:                        args{},
			ClientFunc:                  func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewGetTaskServiceClientFunc: func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient { return nil },
			want:                        nil,
			wantErr:                     true,
		},
		{
			name:                        "TaskCollection error",
			args:                        args{},
			ClientFunc:                  func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewGetTaskServiceClientFunc: func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient { return fakeStruct{} },
			want:                        &taskproto.TaskResponse{StatusCode: 500, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:                     true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewGetTaskServiceClientFunc = tt.NewGetTaskServiceClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := TaskCollection(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TaskCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTaskService(t *testing.T) {
	type args struct {
		req *taskproto.GetTaskRequest
	}
	tests := []struct {
		name                        string
		args                        args
		ClientFunc                  func(clientName string) (*grpc.ClientConn, error)
		NewGetTaskServiceClientFunc func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient
		want                        *taskproto.TaskResponse
		wantErr                     bool
	}{
		{
			name:                        "Client func error",
			args:                        args{},
			ClientFunc:                  func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewGetTaskServiceClientFunc: func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient { return nil },
			want:                        nil,
			wantErr:                     true,
		},
		{
			name:                        "GetTaskService error",
			args:                        args{},
			ClientFunc:                  func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewGetTaskServiceClientFunc: func(cc *grpc.ClientConn) taskproto.GetTaskServiceClient { return fakeStruct{} },
			want:                        &taskproto.TaskResponse{StatusCode: 500, StatusMessage: "Base.1.13.0.InternalError", Body: []byte("{\"error\":{\"code\":\"Base.1.13.0.GeneralError\",\"message\":\"An error has occurred. See ExtendedInfo for more information.\",\"@Message.ExtendedInfo\":[{\"@odata.type\":\"#Message.v1_1_2.Message\",\"MessageId\":\"Base.1.13.0.InternalError\",\"Message\":\"The request failed due to an internal service error.  The service is still operational.fakeError\",\"Severity\":\"Critical\",\"Resolution\":\"Resubmit the request.  If the problem persists, consider resetting the service.\"}]}}")},
			wantErr:                     true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewGetTaskServiceClientFunc = tt.NewGetTaskServiceClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTaskService(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTaskService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTaskService() = %v, want %v", got, tt.want)
			}
		})
	}
}
