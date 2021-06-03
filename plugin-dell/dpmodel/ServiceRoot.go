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

//Package dpmodel ...
package dpmodel

// ServiceRoot a Redfish service root
type ServiceRoot struct {
	Name           string
	UUID           string
	ID             string
	RedfishVersion string
	Context        string `json:"@odata.context"`
	Etag           string `json:"@odata.etag,omitempty"`
	Oid            string `json:"@odata.id"`
	Type           string `json:"@odata.type"`
	Systems        Systems
	Chassis        Chassis
	AccountService AccountService
	EventService   EventService
	JSONSchemas    JSONSchemas
	Managers       Managers
	SessionService SessionService
	Links          Links
}

// Systems a Redfish system link
type Systems struct {
	Oid string `json:"@odata.id"`
}

// Chassis a Redfish chassis link
type Chassis struct {
	Oid string `json:"@odata.id"`
}

// AccountService a Redfish account service link
type AccountService struct {
	Oid string `json:"@odata.id"`
}

// EventService a Redfish event service link
type EventService struct {
	Oid string `json:"@odata.id"`
}

// JSONSchemas a Redfish json schemas link
type JSONSchemas struct {
	Oid string `json:"@odata.id"`
}

// Managers a Redfish managers link
type Managers struct {
	Oid string `json:"@odata.id"`
}

// SessionService a Redfish session service link
type SessionService struct {
	Oid string `json:"@odata.id"`
}

// Sessions ...
type Sessions struct {
	Oid string `json:"@odata.id"`
}
