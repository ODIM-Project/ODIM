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

	teleproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/telemetry"
	"github.com/ODIM-Project/ODIM/lib-utilities/services"
)

// DoGetTelemetryService defines the RPC call function for
// the GetTelemetryService from telemetry micro service
func DoGetTelemetryService(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	conn, err := services.ODIMService.Client(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	telemetry := teleproto.NewTelemetryClient(conn)

	resp, err := telemetry.GetTelemetryService(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetMetricDefinitionCollection defines the RPC call function for
// the GetMetricDefinitionCollection from telemetry micro service
func DoGetMetricDefinitionCollection(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	conn, err := services.ODIMService.Client(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	telemetry := teleproto.NewTelemetryClient(conn)
	resp, err := telemetry.GetMetricDefinitionCollection(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetMetricReportDefinitionCollection defines the RPC call function for
// the GetMetricReportDefinitionCollection from telemetry micro service
func DoGetMetricReportDefinitionCollection(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	conn, err := services.ODIMService.Client(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	telemetry := teleproto.NewTelemetryClient(conn)

	resp, err := telemetry.GetMetricReportDefinitionCollection(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetMetricReportCollection defines the RPC call function for
// the GetMetricReportCollection from telemetry micro service
func DoGetMetricReportCollection(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	conn, err := services.ODIMService.Client(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	telemetry := teleproto.NewTelemetryClient(conn)

	resp, err := telemetry.GetMetricReportCollection(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetTriggerCollection defines the RPC call function for
// the GetTriggerCollection from telemetry micro service
func DoGetTriggerCollection(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	conn, err := services.ODIMService.Client(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	telemetry := teleproto.NewTelemetryClient(conn)

	resp, err := telemetry.GetTriggerCollection(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetMetricDefinition defines the RPC call function for
// the GetMetricDefinition from telemetry micro service
func DoGetMetricDefinition(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	conn, err := services.ODIMService.Client(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	telemetry := teleproto.NewTelemetryClient(conn)

	resp, err := telemetry.GetMetricDefinition(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetMetricReportDefinition defines the RPC call function for
// the GetMetricReportDefinition from telemetry micro service
func DoGetMetricReportDefinition(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	conn, err := services.ODIMService.Client(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	telemetry := teleproto.NewTelemetryClient(conn)

	resp, err := telemetry.GetMetricReportDefinition(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetMetricReport defines the RPC call function for
// the GetMetricReport from telemetry micro service
func DoGetMetricReport(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	conn, err := services.ODIMService.Client(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	telemetry := teleproto.NewTelemetryClient(conn)

	resp, err := telemetry.GetMetricReport(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoGetTrigger defines the RPC call function for
// the GetTrigger from telemetry micro service
func DoGetTrigger(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	conn, err := services.ODIMService.Client(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	telemetry := teleproto.NewTelemetryClient(conn)

	resp, err := telemetry.GetTrigger(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}

// DoUpdateTrigger defines the RPC call function for
// the UpdateTrigger from telemetry micro service
func DoUpdateTrigger(req teleproto.TelemetryRequest) (*teleproto.TelemetryResponse, error) {

	conn, err := services.ODIMService.Client(services.Telemetry)
	if err != nil {
		return nil, fmt.Errorf("Failed to create client connection: %v", err)
	}
	defer conn.Close()
	telemetry := teleproto.NewTelemetryClient(conn)

	resp, err := telemetry.UpdateTrigger(context.TODO(), &req)
	if err != nil {
		return nil, fmt.Errorf("error: RPC error: %v", err)
	}

	return resp, err
}
