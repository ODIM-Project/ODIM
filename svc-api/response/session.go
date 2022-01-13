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

//Package response ...
package response

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

// Session struct is used to ommit password for display purposes
type Session struct {
	OdataContext string `json:"@odata.context"`
	Etag         string `json:"@odata.etag,omitempty"`
	OdataID      string `json:"@odata.id"`
	OdataType    string `json:"@odata.type"`
	ID           string `json:"ID"`
	Description  string `json:"Description"`
	Name         string `json:"Name"`
	OEM          OEM    `json:"Oem"`
	UserName     string `json:"UserName"`
}

//SessionService struct definition
type SessionService struct {
	OdataType      string   `json:"@odata.type"`
	ID             string   `json:"Id"`
	Name           string   `json:"Name"`
	Description    string   `json:"Description,omitempty"`
	Status         Status   `json:"Status,omitempty"`
	ServiceEnabled bool     `json:"ServiceEnabled,omitempty"`
	SessionTimeout float64  `json:"SessionTimeout,omitempty"`
	Sessions       Sessions `json:"Sessions,omitempty"`
	OdataContext   string   `json:"@odata.context,omitempty"`
	OdataID        string   `json:"@odata.id"`
	Etag           string   `json:"@odata.etag,omitempty"`
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
