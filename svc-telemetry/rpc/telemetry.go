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

package rpc

import (
	"context"
	"net/http"
	"os"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
)

var podName = os.Getenv("POD_NAME")

// GetTelemetryService is an rpc handler, it gets invoked during GET on TelemetryService API (/redfis/v1/TelemetryService/)
func (a *Telemetry) GetTelemetryService(ctx context.Context, req *teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TelemetryService, podName)
	resp := &teleproto.TelemetryResponse{}
	fillProtoResponse(ctx, resp, a.connector.GetTelemetryService())
	return resp, nil
}

// GetMetricDefinitionCollection an rpc handler which is invoked during GET on MetricDefinition Collection
func (a *Telemetry) GetMetricDefinitionCollection(ctx context.Context, req *teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TelemetryService, podName)
	resp := &teleproto.TelemetryResponse{}
	authResp, err := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, a.connector.GetMetricDefinitionCollection(req))
	return resp, nil
}

// GetMetricReportDefinitionCollection is an rpc handler which is invoked during GET on MetricReportDefinition Collection
func (a *Telemetry) GetMetricReportDefinitionCollection(ctx context.Context, req *teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TelemetryService, podName)
	resp := &teleproto.TelemetryResponse{}
	authResp, err := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, a.connector.GetMetricReportDefinitionCollection(req))
	return resp, nil
}

// GetMetricReportCollection is an rpc handler which is invoked during GET on MetricReport Collection
func (a *Telemetry) GetMetricReportCollection(ctx context.Context, req *teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TelemetryService, podName)
	resp := &teleproto.TelemetryResponse{}
	authResp, err := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, a.connector.GetMetricReportCollection(req))
	return resp, nil
}

// GetTriggerCollection is an rpc handler which is invoked during GET on TriggerCollection
func (a *Telemetry) GetTriggerCollection(ctx context.Context, req *teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TelemetryService, podName)
	resp := &teleproto.TelemetryResponse{}
	authResp, err := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, a.connector.GetTriggerCollection(req))
	return resp, nil
}

// GetMetricDefinition is an rpc handler which is invoked during GET on MetricDefinition
func (a *Telemetry) GetMetricDefinition(ctx context.Context, req *teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TelemetryService, podName)
	resp := &teleproto.TelemetryResponse{}
	authResp, err := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, a.connector.GetMetricDefinition(ctx, req))
	return resp, nil
}

// GetMetricReportDefinition is an rpc handler which is invoked during GET on MetricReportDefinition
func (a *Telemetry) GetMetricReportDefinition(ctx context.Context, req *teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TelemetryService, podName)
	resp := &teleproto.TelemetryResponse{}
	authResp, err := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, a.connector.GetMetricReportDefinition(ctx, req))
	return resp, nil
}

// GetMetricReport is an rpc handler which is invoked during GET on MetricReport
func (a *Telemetry) GetMetricReport(ctx context.Context, req *teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TelemetryService, podName)
	resp := &teleproto.TelemetryResponse{}
	authResp, err := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, a.connector.GetMetricReport(ctx, req))
	return resp, nil
}

// GetTrigger is an rpc handler which is invoked during GET on Triggers
func (a *Telemetry) GetTrigger(ctx context.Context, req *teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TelemetryService, podName)
	resp := &teleproto.TelemetryResponse{}
	authResp, err := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, a.connector.GetTrigger(ctx, req))
	return resp, nil
}

// UpdateTrigger is an rpc handler which is invoked during update on Trigger
func (a *Telemetry) UpdateTrigger(ctx context.Context, req *teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.GetContextData(ctx)
	ctx = common.ModifyContext(ctx, common.TelemetryService, podName)
	resp := &teleproto.TelemetryResponse{}
	authResp, err := a.connector.External.Auth(req.SessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		if err != nil {
			l.LogWithFields(ctx).Errorf("Error while authorizing the session token : %s", err.Error())
		}
		fillProtoResponse(ctx, resp, authResp)
		return resp, nil
	}
	fillProtoResponse(ctx, resp, a.connector.UpdateTrigger(req))
	return resp, nil
}
