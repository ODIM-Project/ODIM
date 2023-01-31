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
	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

var (
	NewTelemetryClientFunc = teleproto.NewTelemetryClient
)

// DoGetTelemetryService defines the RPC call function for
// the GetTelemetryService from telemetry micro service
func DoGetTelemetryService(ctx context.Context, req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	telemetry := NewTelemetryClientFunc(conn)

	resp, err := telemetry.GetTelemetryService(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetMetricDefinitionCollection defines the RPC call function for
// the GetMetricDefinitionCollection from telemetry micro service
func DoGetMetricDefinitionCollection(ctx context.Context, req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	telemetry := NewTelemetryClientFunc(conn)
	resp, err := telemetry.GetMetricDefinitionCollection(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetMetricReportDefinitionCollection defines the RPC call function for
// the GetMetricReportDefinitionCollection from telemetry micro service
func DoGetMetricReportDefinitionCollection(ctx context.Context, req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	telemetry := NewTelemetryClientFunc(conn)
	resp, err := telemetry.GetMetricReportDefinitionCollection(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetMetricReportCollection defines the RPC call function for
// the GetMetricReportCollection from telemetry micro service
func DoGetMetricReportCollection(ctx context.Context, req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	telemetry := NewTelemetryClientFunc(conn)

	resp, err := telemetry.GetMetricReportCollection(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetTriggerCollection defines the RPC call function for
// the GetTriggerCollection from telemetry micro service
func DoGetTriggerCollection(ctx context.Context, req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	telemetry := NewTelemetryClientFunc(conn)

	resp, err := telemetry.GetTriggerCollection(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetMetricDefinition defines the RPC call function for
// the GetMetricDefinition from telemetry micro service
func DoGetMetricDefinition(ctx context.Context, req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	telemetry := NewTelemetryClientFunc(conn)

	resp, err := telemetry.GetMetricDefinition(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetMetricReportDefinition defines the RPC call function for
// the GetMetricReportDefinition from telemetry micro service
func DoGetMetricReportDefinition(ctx context.Context, req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	telemetry := NewTelemetryClientFunc(conn)

	resp, err := telemetry.GetMetricReportDefinition(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetMetricReport defines the RPC call function for
// the GetMetricReport from telemetry micro service
func DoGetMetricReport(ctx context.Context, req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	telemetry := NewTelemetryClientFunc(conn)

	resp, err := telemetry.GetMetricReport(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoGetTrigger defines the RPC call function for
// the GetTrigger from telemetry micro service
func DoGetTrigger(ctx context.Context, req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	telemetry := NewTelemetryClientFunc(conn)

	resp, err := telemetry.GetTrigger(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}

// DoUpdateTrigger defines the RPC call function for
// the UpdateTrigger from telemetry micro service
func DoUpdateTrigger(ctx context.Context, req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {
	ctx = common.CreateMetadata(ctx)
	conn, err := ClientFunc(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}

	telemetry := NewTelemetryClientFunc(conn)

	resp, err := telemetry.UpdateTrigger(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}
	defer conn.Close()
	return resp, err
}
