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

// Package agresponse ...
package agresponse

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// ResetResponse ...
type ResetResponse struct {
	ResetType string      `json:"ResetType"`
	TargetURI string      `json:"TargetUri"`
	Response  interface{} `json:"Response"`
}

//AggregationServiceResponse is used to give back the response
type AggregationServiceResponse struct {
	OdataContext       string  `json:"@odata.context"`
	Etag               string  `json:"@odata.etag,omitempty"`
	ID                 string  `json:"Id"`
	OdataID            string  `json:"@odata.id"`
	OdataType          string  `json:"@odata.type"`
	Name               string  `json:"Name"`
	Description        string  `json:"Description"`
	Actions            Actions `json:"Actions"`
	Aggregates         OdataID `json:"Aggregates"`
	AggregationSources OdataID `json:"AggregationSources"`
	ConnectionMethods  OdataID `json:"ConnectionMethods"`
	ServiceEnabled     bool    `json:"ServiceEnabled"`
	Status             Status  `json:"Status"`
}

//Actions struct definition
type Actions struct {
	Reset               Action `json:"#AggregationService.Reset"`
	SetDefaultBootOrder Action `json:"#AggregationService.SetDefaultBootOrder"`
}

//Status struct definition
type Status struct {
	Health       string `json:"Health"`
	HealthRollup string `json:"HealthRollup"`
	State        string `json:"State"`
}

//Action struct definition
type Action struct {
	Target string `json:"target"`
}

//OdataID struct definition for @odata.id
type OdataID struct {
	OdataID string `json:"@odata.id"`
}

// ConnectionMethodResponse defines the response for connection method
type ConnectionMethodResponse struct {
	response.Response
	ConnectionMethodType    string      `json:"ConnectionMethodType"`
	ConnectionMethodVariant string      `json:"ConnectionMethodVariant"`
	Links                   interface{} `json:"Links"`
}
