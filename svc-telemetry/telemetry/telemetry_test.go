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

// Package telemetry ...
package telemetry

import (
	"context"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
)

func TestExternalInterface_GetTelemetryService(t *testing.T) {
	tests := []struct {
		name string
		e    *ExternalInterface
		want int
	}{
		{
			name: "Success",
			e:    MockGetExternalInterface(),
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.GetTelemetryService(); int(got.StatusCode) != tt.want {
				t.Errorf("ExternalInterface.GetTelemetryService() = %v, want %v", int(got.StatusCode), tt.want)
			}
		})
	}
}

func TestExternalInterface_GetMetricDefinitionCollection(t *testing.T) {
	type args struct {
		req *teleproto.TelemetryRequest
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
		want int
	}{
		{
			name: "success",
			e:    MockGetExternalInterface(),
			args: args{
				req: &teleproto.TelemetryRequest{
					ResourceID:   "custom1",
					SessionToken: "validToken",
				},
			},
			want: http.StatusOK,
		},
		{
			name: "error",
			e:    MockGetExternalInterface(),
			args: args{
				req: &teleproto.TelemetryRequest{
					URL:          "error",
					ResourceID:   "custom1",
					SessionToken: "validToken",
				},
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.GetMetricDefinitionCollection(tt.args.req); int(got.StatusCode) != tt.want {
				t.Errorf("ExternalInterface.GetMetricDefinitionCollection() = %v, want %v", int(got.StatusCode), tt.want)
			}
		})
	}
}

func TestExternalInterface_GetMetricReportDefinitionCollection(t *testing.T) {
	type args struct {
		req *teleproto.TelemetryRequest
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
		want int
	}{
		{
			name: "success",
			e:    MockGetExternalInterface(),
			args: args{
				req: &teleproto.TelemetryRequest{
					ResourceID:   "custom1",
					SessionToken: "validToken",
				},
			},
			want: http.StatusOK,
		},
		{
			name: "error",
			e:    MockGetExternalInterface(),
			args: args{
				req: &teleproto.TelemetryRequest{
					URL:          "error",
					ResourceID:   "custom1",
					SessionToken: "validToken",
				},
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.GetMetricReportDefinitionCollection(tt.args.req); int(got.StatusCode) != tt.want {
				t.Errorf("ExternalInterface.GetMetricReportDefinitionCollection() = %v, want %v", int(got.StatusCode), tt.want)
			}
		})
	}
}

func TestExternalInterface_GetMetricReportCollection(t *testing.T) {
	type args struct {
		req *teleproto.TelemetryRequest
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
		want int
	}{
		{
			name: "success",
			e:    MockGetExternalInterface(),
			args: args{
				req: &teleproto.TelemetryRequest{
					ResourceID:   "custom1",
					SessionToken: "validToken",
				},
			},
			want: http.StatusOK,
		},
		{
			name: "error",
			e:    MockGetExternalInterface(),
			args: args{
				req: &teleproto.TelemetryRequest{
					URL:          "error",
					ResourceID:   "custom1",
					SessionToken: "validToken",
				},
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.GetMetricReportCollection(tt.args.req); int(got.StatusCode) != tt.want {
				t.Errorf("ExternalInterface.GetMetricReportCollection() = %v, want %v", int(got.StatusCode), tt.want)
			}
		})
	}
}

func TestExternalInterface_GetTriggerCollection(t *testing.T) {
	type args struct {
		req *teleproto.TelemetryRequest
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
		want int
	}{
		{
			name: "success",
			e:    MockGetExternalInterface(),
			args: args{
				req: &teleproto.TelemetryRequest{
					ResourceID:   "custom1",
					SessionToken: "validToken",
				},
			},
			want: http.StatusOK,
		},
		{
			name: "error",
			e:    MockGetExternalInterface(),
			args: args{
				req: &teleproto.TelemetryRequest{
					URL:          "error",
					ResourceID:   "custom1",
					SessionToken: "validToken",
				},
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.GetTriggerCollection(tt.args.req); int(got.StatusCode) != tt.want {
				t.Errorf("ExternalInterface.GetTriggerCollection() = %v, want %v", int(got.StatusCode), tt.want)
			}
		})
	}
}

func TestExternalInterface_GetMetricReportDefinition(t *testing.T) {
	type args struct {
		ctx context.Context
		req *teleproto.TelemetryRequest
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
		want int
	}{
		{
			name: "success",
			e:    MockGetExternalInterface(),
			args: args{
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{
					ResourceID:   "custom1",
					SessionToken: "validToken",
				},
			},
			want: http.StatusOK,
		},
		{
			name: "error",
			e:    MockGetExternalInterface(),
			args: args{
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{
					URL:          "error",
					ResourceID:   "custom1",
					SessionToken: "validToken",
				},
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.GetMetricReportDefinition(tt.args.ctx, tt.args.req); int(got.StatusCode) != tt.want {
				t.Errorf("ExternalInterface.GetMetricReportDefinition() = %v, want %v", int(got.StatusCode), tt.want)
			}
		})
	}
}

func TestExternalInterface_GetMetricReport(t *testing.T) {
	config.SetUpMockConfig(t)
	type args struct {
		ctx context.Context
		req *teleproto.TelemetryRequest
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
		want int
	}{
		{
			name: "success",
			e:    MockGetExternalInterface(),
			args: args{
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{
					ResourceID:   "custom1",
					SessionToken: "validToken",
					URL:          "/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1",
				},
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.GetMetricReport(tt.args.ctx, tt.args.req); int(got.StatusCode) != tt.want {
				t.Errorf("ExternalInterface.GetMetricReport() = %v, want %v", int(got.StatusCode), tt.want)
			}
		})
	}
}

func TestExternalInterface_GetMetricDefinition(t *testing.T) {
	type args struct {
		ctx context.Context
		req *teleproto.TelemetryRequest
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
		want int
	}{
		{
			name: "success",
			e:    MockGetExternalInterface(),
			args: args{
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{
					ResourceID:   "custom1",
					SessionToken: "validToken",
				},
			},
			want: http.StatusOK,
		},
		{
			name: "error",
			e:    MockGetExternalInterface(),
			args: args{
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{
					ResourceID:   "custom1",
					SessionToken: "validToken",
					URL:          "error",
				},
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.GetMetricDefinition(tt.args.ctx, tt.args.req); int(got.StatusCode) != tt.want {
				t.Errorf("ExternalInterface.GetMetricDefinition() = %v, want %v", int(got.StatusCode), tt.want)
			}
		})
	}
}

func TestExternalInterface_GetTrigger(t *testing.T) {
	type args struct {
		ctx context.Context
		req *teleproto.TelemetryRequest
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
		want int
	}{
		{
			name: "success",
			e:    MockGetExternalInterface(),
			args: args{
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{
					ResourceID:   "custom1",
					SessionToken: "validToken",
				},
			},
			want: http.StatusOK,
		},
		{
			name: "error",
			e:    MockGetExternalInterface(),
			args: args{
				ctx: context.Background(),
				req: &teleproto.TelemetryRequest{
					ResourceID:   "custom1",
					SessionToken: "validToken",
					URL:          "error",
				},
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.GetTrigger(tt.args.ctx, tt.args.req); int(got.StatusCode) != tt.want {
				t.Errorf("ExternalInterface.GetTrigger() = %v, want %v", int(got.StatusCode), tt.want)
			}
		})
	}
}

func TestExternalInterface_UpdateTrigger(t *testing.T) {
	type args struct {
		req *teleproto.TelemetryRequest
	}
	tests := []struct {
		name string
		e    *ExternalInterface
		args args
		want int
	}{
		{
			name: "success",
			e:    MockGetExternalInterface(),
			args: args{
				req: &teleproto.TelemetryRequest{
					ResourceID:   "custom1",
					SessionToken: "validToken",
				},
			},
			want: 0, // as update trigger is not implemented. Need to be changed with http.StatusOK once update trigger operation is implemented
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.UpdateTrigger(tt.args.req); int(got.StatusCode) != tt.want {
				t.Errorf("ExternalInterface.UpdateTrigger() = %v, want %v", int(got.StatusCode), tt.want)
			}
		})
	}
}
