// (C) Copyright [2020] Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package rpc

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-telemetry/telemetry"
	"github.com/ODIM-Project/ODIM/svc-telemetry/tmodel"
	"github.com/stretchr/testify/assert"
)

func mockIsAuthorized(sessionToken string, privileges, oemPrivileges []string) response.RPC {
	if sessionToken == "InvalidToken" {
		return common.GeneralError(http.StatusUnauthorized, response.NoValidSession, "error while trying to authenticate session", nil, nil)
	}
	return common.GeneralError(http.StatusOK, response.Success, "", nil, nil)
}

func mockContactClient(url, method, token string, odataID string, body interface{}, loginCredential map[string]string) (*http.Response, error) {
	if url == "https://localhost:9091/ODIM/v1/Sessions" {
		body := `{"Token": "12345"}`
		return &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
			Header: http.Header{
				"X-Auth-Token": []string{"12345"},
			},
		}, nil
	} else if url == "https://localhost:9092/ODIM/v1/Sessions" {
		body := `{"Token": ""}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	if url == "https://localhost:9091/ODIM/v1/TelemetryService/MetricReports/CPUUtilCustom1" && token == "12345" {
		body := `{"@odata.id":"/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1","@odata.type":"#MetricReport.v1_0_0.MetricReport","Id":"CPUUtilCustom1","Name":"Metric report of CPU Utilization for 10 minutes with sensing interval of 20 seconds.","MetricReportDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricReportDefinitions/CPUUtilCustom1"},"MetricValues":[{"MetricDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"},"MetricId":"CPUUtil","MetricValue":"0","Timestamp":"2021-06-16T07:59:43Z"},{"MetricDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"},"MetricId":"CPUUtil","MetricValue":"0","Timestamp":"2021-06-16T08:00:04Z"}]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9093/ODIM/v1/TelemetryService/MetricReports/CPUUtilCustom1" {
		body := `{"@odata.id":"/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1","@odata.type":"#MetricReport.v1_0_0.MetricReport","Id":"CPUUtilCustom1","Name":"Metric report of CPU Utilization for 10 minutes with sensing interval of 20 seconds.","MetricReportDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricReportDefinitions/CPUUtilCustom1"},"MetricValues":[{"MetricDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"},"MetricId":"CPUUtil","MetricValue":"0","Timestamp":"2021-06-16T07:59:43Z"},{"MetricDefinition":{"@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"},"MetricId":"CPUUtil","MetricValue":"0","Timestamp":"2021-06-16T08:00:04Z"}]}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9092/ODIM/v1/TelemetryService/MetricReports/CPUUtilCustom1" && token == "23456" {
		body := `{"data": "ODIM/v1/TelemetryService/MetricReports/CPUUtilCustom1"}`
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	} else if url == "https://localhost:9091/ODIM/v1/TelemetryService/MetricReports/CPUUtilCustom1" {
		body := `{"@odata.id": "/ODIM/v1/TelemetryService/MetricReports/CPUUtilCustom1"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
		}, nil
	}
	return nil, fmt.Errorf("InvalidRequest")
}

func mockGetResource(table, key string, dbType common.DbType) (string, *errors.Error) {
	return "body", nil
}

func mockGetAllKeysFromTable(table string, dbType common.DbType) ([]string, error) {
	if table == "Plugin" {
		return []string{"ILO", "GRF"}, nil
	}
	return []string{"/redfish/v1/TelemetryService/Triggers/uuid:1"}, nil
}

func getEncryptedKey(t *testing.T, key []byte) []byte {
	cryptedKey, err := common.EncryptWithPublicKey(key)
	if err != nil {
		t.Fatalf("error: failed to encrypt data: %v", err)
	}
	return cryptedKey
}

func mockGetPluginData(pluginID string) (tmodel.Plugin, *errors.Error) {
	var t *testing.T
	password := getEncryptedKey(t, []byte("$2a$10$OgSUYvuYdI/7dLL5KkYNp.RCXISefftdj.MjbBTr6vWyNwAvht6ci"))
	plugin := tmodel.Plugin{
		IP:                "localhost",
		Port:              "9091",
		Username:          "admin",
		Password:          password,
		ID:                pluginID,
		PreferredAuthType: "XAuthToken",
		PluginType:        "Compute",
	}
	return plugin, nil
}

func mockGetExternalInterface() *telemetry.ExternalInterface {
	return &telemetry.ExternalInterface{
		External: telemetry.External{
			Auth:          mockIsAuthorized,
			ContactClient: mockContactClient,
			GetPluginData: mockGetPluginData,
		},
		DB: telemetry.DB{
			GetAllKeysFromTable: mockGetAllKeysFromTable,
			GetResource:         mockGetResource,
		},
	}
}

func TestTelemetry_GetTelemetryService(t *testing.T) {
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *teleproto.TelemetryRequest
	}
	tests := []struct {
		name    string
		a       *Telemetry
		args    args
		wantErr bool
	}{
		{
			name: "positive GetTelemetryService",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{SessionToken: "validToken"},
			},
			wantErr: false,
		},
		{
			name: "auth fail",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{SessionToken: "invalidToken"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := tt.a.GetTelemetryService(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("Telemetry.GetTelemetryService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTelemetry_GetMetricDefinitionCollection(t *testing.T) {
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *teleproto.TelemetryRequest
	}
	tests := []struct {
		name       string
		a          *Telemetry
		args       args
		StatusCode int
	}{
		{
			name: "Request with valid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "validToken",
				},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "invalidToken",
				},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, err := tt.a.GetMetricDefinitionCollection(tt.args.ctx, tt.args.req); err != nil {
				t.Errorf("Telemetry.GetMetricDefinitionCollection() got = %v, want %v", resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestTelemetry_GetMetricReportDefinitionCollection(t *testing.T) {
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *teleproto.TelemetryRequest
	}
	tests := []struct {
		name       string
		a          *Telemetry
		args       args
		StatusCode int
	}{
		{
			name: "Request with valid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "validToken",
				},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "invalidToken",
				},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, err := tt.a.GetMetricReportDefinitionCollection(tt.args.ctx, tt.args.req); err != nil {
				t.Errorf("Telemetry.GetMetricReportDefinitionCollection() got = %v, want %v", resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestTelemetry_GetMetricReportCollection(t *testing.T) {
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *teleproto.TelemetryRequest
	}
	tests := []struct {
		name       string
		a          *Telemetry
		args       args
		StatusCode int
	}{
		{
			name: "Request with valid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "validToken",
				},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "invalidToken",
				},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, err := tt.a.GetMetricReportCollection(tt.args.ctx, tt.args.req); err != nil {
				t.Errorf("Telemetry.GetMetricReportCollection() got = %v, want %v", resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestTelemetry_GetTriggerCollection(t *testing.T) {
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	type args struct {
		ctx context.Context
		req *teleproto.TelemetryRequest
	}
	tests := []struct {
		name       string
		a          *Telemetry
		args       args
		StatusCode int
	}{
		{
			name: "Request with valid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "validToken",
				},
			}, StatusCode: 200,
		},
		{
			name: "Request with invalid token",
			a:    telemetry,
			args: args{
				req: &teleproto.TelemetryRequest{
					SessionToken: "invalidToken",
				},
			}, StatusCode: 401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if resp, err := tt.a.GetTriggerCollection(tt.args.ctx, tt.args.req); err != nil {
				t.Errorf("Telemetry.GetTriggerCollection() got = %v, want %v", resp.StatusCode, tt.StatusCode)
			}
		})
	}
}

func TestGetMetricDefinitionwithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "InvalidToken",
	}
	resp, err := telemetry.GetMetricDefinition(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetMetricDefinitionwithValidtoken(t *testing.T) {
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "validToken",
	}
	resp, err := telemetry.GetMetricDefinition(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetMetricReportDefinitionwithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "InvalidToken",
	}
	resp, err := telemetry.GetMetricReportDefinition(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetMetricReportDefinitionwithValidtoken(t *testing.T) {
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "validToken",
	}
	resp, err := telemetry.GetMetricReportDefinition(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetMetricReportwithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		SessionToken: "InvalidToken",
	}
	resp, err := telemetry.GetMetricReport(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetMetricReportwithValidtoken(t *testing.T) {
	config.SetUpMockConfig(t)
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		SessionToken: "validToken",
		URL:          "/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1",
	}
	resp, err := telemetry.GetMetricReport(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetTriggerwithInValidtoken(t *testing.T) {
	common.SetUpMockConfig()
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "InvalidToken",
	}
	resp, err := telemetry.GetTrigger(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "Status code should be StatusOK.")
}

func TestGetTriggerwithValidtoken(t *testing.T) {
	var ctx context.Context
	telemetry := new(Telemetry)
	telemetry.connector = mockGetExternalInterface()
	req := &teleproto.TelemetryRequest{
		ResourceID:   "custom1",
		SessionToken: "validToken",
	}
	resp, err := telemetry.GetTrigger(ctx, req)
	assert.Nil(t, err, "There should be no error")
	assert.Equal(t, http.StatusOK, int(resp.StatusCode), "Status code should be StatusOK.")
}
