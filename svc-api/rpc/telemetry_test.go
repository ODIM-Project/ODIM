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

	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	"google.golang.org/grpc"
)

func TestDoGetTelemetryService(t *testing.T) {
	type args struct {
		req teleproto.TelemetryRequest
	}
	tests := []struct {
		name                   string
		args                   args
		ClientFunc             func(clientName string) (*grpc.ClientConn, error)
		NewTelemetryClientFunc func(cc *grpc.ClientConn) teleproto.TelemetryClient
		want                   *teleproto.TelemetryResponse
		wantErr                bool
	}{
		{
			name:                   "Client func error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return nil },
			want:                   nil,
			wantErr:                true,
		},
		{
			name:                   "DoGetTelemetryService error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return fakeStruct{} },
			want:                   nil,
			wantErr:                true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewTelemetryClientFunc = tt.NewTelemetryClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetTelemetryService(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetTelemetryService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetTelemetryService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetMetricDefinitionCollection(t *testing.T) {
	type args struct {
		req teleproto.TelemetryRequest
	}
	tests := []struct {
		name                   string
		args                   args
		ClientFunc             func(clientName string) (*grpc.ClientConn, error)
		NewTelemetryClientFunc func(cc *grpc.ClientConn) teleproto.TelemetryClient
		want                   *teleproto.TelemetryResponse
		wantErr                bool
	}{
		{
			name:                   "Client func error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return nil },
			want:                   nil,
			wantErr:                true,
		},
		{
			name:                   "DoGetMetricDefinitionCollection error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return fakeStruct{} },
			want:                   nil,
			wantErr:                true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewTelemetryClientFunc = tt.NewTelemetryClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetMetricDefinitionCollection(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetMetricDefinitionCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetMetricDefinitionCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetMetricReportDefinitionCollection(t *testing.T) {
	type args struct {
		req teleproto.TelemetryRequest
	}
	tests := []struct {
		name                   string
		args                   args
		ClientFunc             func(clientName string) (*grpc.ClientConn, error)
		NewTelemetryClientFunc func(cc *grpc.ClientConn) teleproto.TelemetryClient
		want                   *teleproto.TelemetryResponse
		wantErr                bool
	}{
		{
			name:                   "Client func error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return nil },
			want:                   nil,
			wantErr:                true,
		},
		{
			name:                   "DoGetMetricReportDefinitionCollection error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return fakeStruct{} },
			want:                   nil,
			wantErr:                true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewTelemetryClientFunc = tt.NewTelemetryClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetMetricReportDefinitionCollection(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetMetricReportDefinitionCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetMetricReportDefinitionCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetMetricReportCollection(t *testing.T) {
	type args struct {
		req teleproto.TelemetryRequest
	}
	tests := []struct {
		name                   string
		args                   args
		ClientFunc             func(clientName string) (*grpc.ClientConn, error)
		NewTelemetryClientFunc func(cc *grpc.ClientConn) teleproto.TelemetryClient
		want                   *teleproto.TelemetryResponse
		wantErr                bool
	}{
		{
			name:                   "Client func error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return nil },
			want:                   nil,
			wantErr:                true,
		},
		{
			name:                   "DoGetMetricReportCollection error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return fakeStruct{} },
			want:                   nil,
			wantErr:                true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewTelemetryClientFunc = tt.NewTelemetryClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetMetricReportCollection(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetMetricReportCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetMetricReportCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetTriggerCollection(t *testing.T) {
	type args struct {
		req teleproto.TelemetryRequest
	}
	tests := []struct {
		name                   string
		args                   args
		ClientFunc             func(clientName string) (*grpc.ClientConn, error)
		NewTelemetryClientFunc func(cc *grpc.ClientConn) teleproto.TelemetryClient
		want                   *teleproto.TelemetryResponse
		wantErr                bool
	}{
		{
			name:                   "Client func error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return nil },
			want:                   nil,
			wantErr:                true,
		},
		{
			name:                   "DoGetTriggerCollection error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return fakeStruct{} },
			want:                   nil,
			wantErr:                true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewTelemetryClientFunc = tt.NewTelemetryClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetTriggerCollection(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetTriggerCollection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetTriggerCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetMetricDefinition(t *testing.T) {
	type args struct {
		req teleproto.TelemetryRequest
	}
	tests := []struct {
		name                   string
		args                   args
		ClientFunc             func(clientName string) (*grpc.ClientConn, error)
		NewTelemetryClientFunc func(cc *grpc.ClientConn) teleproto.TelemetryClient
		want                   *teleproto.TelemetryResponse
		wantErr                bool
	}{
		{
			name:                   "Client func error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return nil },
			want:                   nil,
			wantErr:                true,
		},
		{
			name:                   "DoGetMetricDefinition error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return fakeStruct{} },
			want:                   nil,
			wantErr:                true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewTelemetryClientFunc = tt.NewTelemetryClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetMetricDefinition(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetMetricDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetMetricDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetMetricReportDefinition(t *testing.T) {
	type args struct {
		req teleproto.TelemetryRequest
	}
	tests := []struct {
		name                   string
		args                   args
		ClientFunc             func(clientName string) (*grpc.ClientConn, error)
		NewTelemetryClientFunc func(cc *grpc.ClientConn) teleproto.TelemetryClient
		want                   *teleproto.TelemetryResponse
		wantErr                bool
	}{
		{
			name:                   "Client func error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return nil },
			want:                   nil,
			wantErr:                true,
		},
		{
			name:                   "DoGetMetricReportDefinition error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return fakeStruct{} },
			want:                   nil,
			wantErr:                true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewTelemetryClientFunc = tt.NewTelemetryClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetMetricReportDefinition(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetMetricReportDefinition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetMetricReportDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetMetricReport(t *testing.T) {
	type args struct {
		req teleproto.TelemetryRequest
	}
	tests := []struct {
		name                   string
		args                   args
		ClientFunc             func(clientName string) (*grpc.ClientConn, error)
		NewTelemetryClientFunc func(cc *grpc.ClientConn) teleproto.TelemetryClient
		want                   *teleproto.TelemetryResponse
		wantErr                bool
	}{
		{
			name:                   "Client func error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return nil },
			want:                   nil,
			wantErr:                true,
		},
		{
			name:                   "DoGetMetricReport error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return fakeStruct{} },
			want:                   nil,
			wantErr:                true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewTelemetryClientFunc = tt.NewTelemetryClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetMetricReport(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetMetricReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetMetricReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoGetTrigger(t *testing.T) {
	type args struct {
		req teleproto.TelemetryRequest
	}
	tests := []struct {
		name                   string
		args                   args
		ClientFunc             func(clientName string) (*grpc.ClientConn, error)
		NewTelemetryClientFunc func(cc *grpc.ClientConn) teleproto.TelemetryClient
		want                   *teleproto.TelemetryResponse
		wantErr                bool
	}{
		{
			name:                   "Client func error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return nil },
			want:                   nil,
			wantErr:                true,
		},
		{
			name:                   "DoGetTrigger error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return fakeStruct{} },
			want:                   nil,
			wantErr:                true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewTelemetryClientFunc = tt.NewTelemetryClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoGetTrigger(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoGetTrigger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoGetTrigger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoUpdateTrigger(t *testing.T) {
	type args struct {
		req teleproto.TelemetryRequest
	}
	tests := []struct {
		name                   string
		args                   args
		ClientFunc             func(clientName string) (*grpc.ClientConn, error)
		NewTelemetryClientFunc func(cc *grpc.ClientConn) teleproto.TelemetryClient
		want                   *teleproto.TelemetryResponse
		wantErr                bool
	}{
		{
			name:                   "Client func error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, errors.New("fakeError") },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return nil },
			want:                   nil,
			wantErr:                true,
		},
		{
			name:                   "DoUpdateTrigger error",
			args:                   args{},
			ClientFunc:             func(clientName string) (*grpc.ClientConn, error) { return nil, nil },
			NewTelemetryClientFunc: func(cc *grpc.ClientConn) teleproto.TelemetryClient { return fakeStruct{} },
			want:                   nil,
			wantErr:                true,
		},
	}
	for _, tt := range tests {
		ClientFunc = tt.ClientFunc
		NewTelemetryClientFunc = tt.NewTelemetryClientFunc
		t.Run(tt.name, func(t *testing.T) {
			got, err := DoUpdateTrigger(context.Background(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("DoUpdateTrigger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DoUpdateTrigger() = %v, want %v", got, tt.want)
			}
		})
	}
}
