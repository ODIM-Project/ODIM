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

	aggregatorproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/aggregator"
	"google.golang.org/grpc"
)

func TestDoGetAggregationService(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "GetAggregationServices error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetAggregationService(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetAggregationService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetAggregationService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoResetRequest(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "Aggregator Reset error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoResetRequest(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoResetRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoResetRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoSetDefaultBootOrderRequest(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "SetDefaultBootOrder error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoSetDefaultBootOrderRequest(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoSetDefaultBootOrderRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoSetDefaultBootOrderRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoAddAggregationSource(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "AddAggregationSource error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoAddAggregationSource(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoAddAggregationSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoAddAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetAllAggregationSource(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "GetAllAggregationSource error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetAllAggregationSource(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetAllAggregationSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetAllAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetAggregationSource(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "GetAggregationSource error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetAggregationSource(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetAggregationSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoUpdateAggregationSource(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "UpdateAggregationSource error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoUpdateAggregationSource(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoUpdateAggregationSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoUpdateAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoDeleteAggregationSource(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "DeleteAggregationSource error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoDeleteAggregationSource(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoDeleteAggregationSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoDeleteAggregationSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoCreateAggregate(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "CreateAggregate error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoCreateAggregate(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoCreateAggregate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoCreateAggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetAggregateCollection(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "GetAggregateCollection error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetAggregateCollection(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetAggregateCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetAggregateCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGeteAggregate(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "GetAggregate error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGeteAggregate(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGeteAggregate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGeteAggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoDeleteAggregate(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "DeleteAggregate error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoDeleteAggregate(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoDeleteAggregate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoDeleteAggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoAddElementsToAggregate(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "AddElementsToAggregate error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoAddElementsToAggregate(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoAddElementsToAggregate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoAddElementsToAggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoRemoveElementsFromAggregate(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "RemoveElementsFromAggregate error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoRemoveElementsFromAggregate(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoRemoveElementsFromAggregate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoRemoveElementsFromAggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoResetAggregateElements(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "ResetAggregateElements error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoResetAggregateElements(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoResetAggregateElements() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoResetAggregateElements() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoSetDefaultBootOrderAggregateElements(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "SetDefaultBootOrderAggregateElements error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoSetDefaultBootOrderAggregateElements(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoSetDefaultBootOrderAggregateElements() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoSetDefaultBootOrderAggregateElements() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetAllConnectionMethods(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "GetAllConnectionMethods error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetAllConnectionMethods(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetAllConnectionMethods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetAllConnectionMethods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetConnectionMethod(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "GetConnectionMethod error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetConnectionMethod(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetConnectionMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetConnectionMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetResetActionInfoService(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "GetResetActionInfoService error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetResetActionInfoService(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetResetActionInfoService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetResetActionInfoService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetSetDefaultBootOrderActionInfo(t *testing.T) {
	type args struct {
		req aggregatorproto.AggregatorRequest
	}
	tests := []struct {
		name                    string
		args                    args
		ClientFunc              func(clientName string) (*grpc.ClientConn, error)
		NewAggregatorClientFunc func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient
		want                    *aggregatorproto.AggregatorResponse
		wantErr                 bool
	}{
		{
			name:                    "Client func error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return nil },
			want:                    nil,
			wantErr:                 true,
		},
		{
			name:                    "GetSetDefaultBootOrderActionInfo error",
			args:                    args{},
			ClientFunc:              func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewAggregatorClientFunc: func(cc *grpc.ClientConn) aggregatorproto.AggregatorClient { return fakeStruct{} },
			want:                    nil,
			wantErr:                 true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewAggregatorClientFunc = tt.NewAggregatorClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetSetDefaultBootOrderActionInfo(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetSetDefaultBootOrderActionInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetSetDefaultBootOrderActionInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
