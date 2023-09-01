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

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	telemetryproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	iris "github.com/kataras/iris/v12"
)

// TelemetryRPCs used to define the service RPC function
type TelemetryRPCs struct {
	GetTelemetryServiceRPC                 func(context.Context, telemetryproto.TelemetryRequest) (*telemetryproto.TelemetryResponse, error)
	GetMetricDefinitionCollectionRPC       func(context.Context, telemetryproto.TelemetryRequest) (*telemetryproto.TelemetryResponse, error)
	GetMetricReportDefinitionCollectionRPC func(context.Context, telemetryproto.TelemetryRequest) (*telemetryproto.TelemetryResponse, error)
	GetMetricReportCollectionRPC           func(context.Context, telemetryproto.TelemetryRequest) (*telemetryproto.TelemetryResponse, error)
	GetTriggerCollectionRPC                func(context.Context, telemetryproto.TelemetryRequest) (*telemetryproto.TelemetryResponse, error)
	GetMetricDefinitionRPC                 func(context.Context, telemetryproto.TelemetryRequest) (*telemetryproto.TelemetryResponse, error)
	GetMetricReportDefinitionRPC           func(context.Context, telemetryproto.TelemetryRequest) (*telemetryproto.TelemetryResponse, error)
	GetMetricReportRPC                     func(context.Context, telemetryproto.TelemetryRequest) (*telemetryproto.TelemetryResponse, error)
	GetTriggerRPC                          func(context.Context, telemetryproto.TelemetryRequest) (*telemetryproto.TelemetryResponse, error)
	UpdateTriggerRPC                       func(context.Context, telemetryproto.TelemetryRequest) (*telemetryproto.TelemetryResponse, error)
}

const (
	rpcCallFailedErrorMsg    = "error: something went wrong with the RPC calls: "
	invalidAuthTokenErrorMsg = "error: no X-Auth-Token found in request header"
)

// GetTelemetryService is the handler for getting TelemetryService details
func (a *TelemetryRPCs) GetTelemetryService(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := telemetryproto.TelemetryRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
	}
	l.LogWithFields(ctxt).Debug("Incoming request received for getting telemetry service details")
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetTelemetryServiceRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrorMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting telemetry service is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

// GetMetricDefinitionCollection is the handler for getting TelemetryService details
func (a *TelemetryRPCs) GetMetricDefinitionCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := telemetryproto.TelemetryRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting metric definition collection details with request URI %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetMetricDefinitionCollectionRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrorMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting metric definition collection is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

// GetMetricReportDefinitionCollection is the handler for getting TelemetryService details
func (a *TelemetryRPCs) GetMetricReportDefinitionCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := telemetryproto.TelemetryRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting metric report definition collection details with request URI %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetMetricReportDefinitionCollectionRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrorMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting metric report definition collection is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

// GetMetricReportCollection is the handler for getting TelemetryService details
func (a *TelemetryRPCs) GetMetricReportCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := telemetryproto.TelemetryRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting metric report collection details with request URI %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetMetricReportCollectionRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrorMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting metric report collection is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

// GetTriggerCollection is the handler for getting TelemetryService details
func (a *TelemetryRPCs) GetTriggerCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := telemetryproto.TelemetryRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting trigger collection details with request URI %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetTriggerCollectionRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrorMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting trigger collection is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

// GetMetricDefinition is the handler for getting TelemetryService details
func (a *TelemetryRPCs) GetMetricDefinition(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := telemetryproto.TelemetryRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting metric definition details with request URI %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetMetricDefinitionRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrorMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting metric definition is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

// GetMetricReportDefinition is the handler for getting TelemetryService details
func (a *TelemetryRPCs) GetMetricReportDefinition(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := telemetryproto.TelemetryRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting metric report definition details with request URI %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetMetricReportDefinitionRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrorMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting metric report definition is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

// GetMetricReport is the handler for getting TelemetryService details
func (a *TelemetryRPCs) GetMetricReport(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := telemetryproto.TelemetryRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting metric report details with request URI %s", req.URL)
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.GetMetricReportRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrorMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting metric definition collection is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

// GetTrigger is the handler for getting TelemetryService details
func (a *TelemetryRPCs) GetTrigger(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := telemetryproto.TelemetryRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
		URL:          ctx.Request().RequestURI,
	}
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Incoming request received for getting trigger details with request URI %s", req.URL)
	resp, err := a.GetTriggerRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrorMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for getting trigger details is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH")
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}

// UpdateTrigger is the handler for getting TelemetryService details
func (a *TelemetryRPCs) UpdateTrigger(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	req := telemetryproto.TelemetryRequest{
		SessionToken: ctx.Request().Header.Get(AuthTokenHeader),
	}
	l.LogWithFields(ctxt).Debug("Incoming request received for updating trigger")
	if req.SessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
	}
	resp, err := a.UpdateTriggerRPC(ctxt, req)
	if err != nil {
		errorMessage := rpcCallFailedErrorMsg + err.Error()
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	l.LogWithFields(ctxt).Debugf("Outgoing response for updating trigger is %s with status code %d", string(resp.Body), int(resp.StatusCode))
	common.SetResponseHeader(ctx, resp.Header)
	ctx.StatusCode(int(resp.StatusCode))
	ctx.Write(resp.Body)

}
