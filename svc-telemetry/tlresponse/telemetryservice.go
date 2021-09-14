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

// Package tlresponse ...
package tlresponse

import (
	dmtf "github.com/ODIM-Project/ODIM/lib-dmtf/model"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// Status defines the service status
type Status struct {
	State        string `json:"State"`
	Health       string `json:"Health"`
	HealthRollup string `json:"HealthRollup"`
}

// TelemetryService defines the service properties of update service
type TelemetryService struct {
	response.Response
	Status                       Status     `json:"Status,omitempty"`
	ServiceEnabled               bool       `json:"ServiceEnabled,omitempty"`
	SupportedCollectionFunctions []string   `json:"SupportedCollectionFunctions,omitempty"`
	MinCollectionInterval        string     `json:"MinCollectionInterval,omitempty"`
	MetricDefinitions            *dmtf.Link `json:"MetricDefinitions,omitempty"`
	MetricReportDefinitions      *dmtf.Link `json:"MetricReportDefinitions,omitempty"`
	MetricReports                *dmtf.Link `json:"MetricReports,omitempty"`
	Triggers                     *dmtf.Link `json:"Triggers,omitempty"`
	OEM                          *OEM       `json:"Oem,omitempty"`
}

// OEM defines the ACME defined properties under the service
type OEM struct {
}
