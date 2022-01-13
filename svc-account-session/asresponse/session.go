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

// Package asresponse ...
package asresponse

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// RedfishSessionResponse will have all the
type RedfishSessionResponse struct {
	Error Error `json:"error"`
}

// Error is the internal structure of RedfishSessionResponse
type Error struct {
	Code          string         `json:"code"`
	Message       string         `json:"message"`
	ExtendedInfos []ExtendedInfo `json:"@Message.ExtendedInfo"`
}

// ExtendedInfo is the extended version of RedfishSessionResponse message
type ExtendedInfo struct {
	MessageID string `json:"MessageId"`
}

// Session struct is used to omit password for display purposes
type Session struct {
	response.Response
	UserName string `json:"UserName"`
}

//SessionService struct definition
type SessionService struct {
	response.Response
	Status         Status   `json:"Status,omitempty"`
	ServiceEnabled bool     `json:"ServiceEnabled,omitempty"`
	SessionTimeout float64  `json:"SessionTimeout,omitempty"`
	Sessions       Sessions `json:"Sessions,omitempty"`
}

//Sessions struct definition
type Sessions struct {
	OdataID string `json:"@odata.id"`
}

//Status struct definition
type Status struct {
	State  string `json:"State"`
	Health string `json:"Health"`
}
