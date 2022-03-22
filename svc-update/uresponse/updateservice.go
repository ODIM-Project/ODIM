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

// Package uresponse ...
package uresponse

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// Status defines the service status
type Status struct {
	State        string `json:"State"`
	Health       string `json:"Health"`
	HealthRollup string `json:"HealthRollup"`
}

// FirmwareInventory defines the links to BMC firmware inventory
type FirmwareInventory struct {
	OdataID string `json:"@odata.id"`
}

// SoftwareInventory defines the links to BMC software inventory
type SoftwareInventory struct {
	OdataID string `json:"@odata.id"`
}

// ClientCertificates defines client identity certificates provided to the server referenced by the Destination property
type ClientCertificates struct {
	OdataID string `json:"@odata.id"`
}

// RemoteServerCertificates define server certificates for the server referenced by the ImageURI property in SimpleUpdate
type RemoteServerCertificates struct {
	OdataID string `json:"@odata.id"`
}

// UpdateServiceSimpleUpdate defines Target information for the upgrade
type UpdateServiceSimpleUpdate struct {
	Target                           string                           `json:"target"`
	ActionInfo                       string                           `json:"@Redfish.ActionInfo,omitempty"`
	RedfishOperationApplyTimeSupport RedfishOperationApplyTimeSupport `json:"@Redfish.OperationApplyTimeSupport,omitempty"`
}

// RedfishOperationApplyTimeSupport struct defines the apply time for the action in place
type RedfishOperationApplyTimeSupport struct {
	OdataType       string   `json:"@odata.type,omitempty"`
	SupportedValues []string `json:"SupportedValues,omitempty"`
}

// UpdateServiceStartUpdate defines Target information for the upgrade
type UpdateServiceStartUpdate struct {
	Target     string `json:"target"`
	ActionInfo string `json:"@Redfish.ActionInfo,omitempty"`
}

// Actions defines the links to the actions available under the service
type Actions struct {
	UpdateServiceSimpleUpdate UpdateServiceSimpleUpdate `json:"#UpdateService.SimpleUpdate"`
	UpdateServiceStartUpdate  UpdateServiceStartUpdate  `json:"#UpdateService.StartUpdate"`
}

// UpdateService defines the service properties of update service
type UpdateService struct {
	response.Response
	Status                        Status                    `json:"Status"`
	ServiceEnabled                bool                      `json:"ServiceEnabled"`
	HTTPPushURI                   string                    `json:"HttpPushUri"`
	FirmwareInventory             FirmwareInventory         `json:"FirmwareInventory"`
	SoftwareInventory             SoftwareInventory         `json:"SoftwareInventory"`
	Actions                       Actions                   `json:"Actions"`
	OEM                           *OEM                      `json:"Oem,omitempty"`
	ClientCertificates            *ClientCertificates       `json:"ClientCertificates,omitempty"`
	HTTPPushURIOptions            *HTTPPushURIOptions       `json:"HttpPushUriOptions,omitempty"`
	HTTPPushURIiOptionsBusy       bool                      `json:"HttpPushUriOptionsBusy,omitempty"`
	HTTPPushURITargets            []string                  `json:"HttpPushUriTargets,omitempty"`
	HTTPPushURITargetsBusy        bool                      `json:"HttpPushUriTargetsBusy,omitempty"`
	MaxImageSizeBytes             int                       `json:"MaxImageSizeBytes,omitempty"`
	MultipartHTTPPushURI          string                    `json:"MultipartHttpPushUri,omitempty"`
	VerifyRemoteServerCertificate bool                      `json:"VerifyRemoteServerCertificate,omitempty"`
	RemoteServerCertificates      *RemoteServerCertificates `json:"RemoteServerCertificates,omitempty"`
}

// HTTPPushURIOptions defines the options for HttpPushUri-provided software updates
type HTTPPushURIOptions struct {
	HTTPPushURIApplyTime *HTTPPushURIApplyTime `json:"HTTPPushURIApplyTime,omitempty"`
}

// HTTPPushURIApplyTime defines the settings for when toapply HttpPushUri-provided firmware.
type HTTPPushURIApplyTime struct {
	ApplyTime                          string `json:"ApplyTime,omitempty"`
	MaintenanceWindowDurationInSeconds int    `json:"MaintenanceWindowDurationInSeconds,omitempty"`
	MaintenanceWindowStartTime         string `json:"MaintenanceWindowStartTime,omitempty"`
}

// OEM defines the ACME defined properties under the service
type OEM struct {
}
