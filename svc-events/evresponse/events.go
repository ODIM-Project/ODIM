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

// Package evresponse have error and response struct
// and also have functionality to create error response
package evresponse

import (
	"sync"

	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// SubscriptionResponse is used to return response to end user
type SubscriptionResponse struct {
	response.Response
	Destination      string       `json:"Destination,omitempty"`
	Context          string       `json:"Context,omitempty"`
	Protocol         string       `json:"Protocol,omitempty"`
	EventTypes       []string     `json:"EventTypes,omitempty"`
	SubscriptionType string       `json:"SubscriptionType,omitempty"`
	MessageIds       []string     `json:"MessageIds,omitempty"`
	ResourceTypes    []string     `json:"ResourceTypes,omitempty"`
	OriginResources  []ListMember `json:"OriginResources,omitempty"`
}

// ListResponse define list for odimra
type ListResponse struct {
	OdataContext string       `json:"@odata.context"`
	Etag         string       `json:"@odata.etag,omitempty"`
	OdataID      string       `json:"@odata.id"`
	OdataType    string       `json:"@odata.type"`
	Name         string       `json:"Name,omitempty"`
	Description  string       `json:"Description,omitempty"`
	MembersCount int          `json:"Members@odata.count"`
	Members      []ListMember `json:"Members"`
}

// ListMember containes link to a resource
type ListMember struct {
	OdataID string `json:"@odata.id"`
}

//EventServiceResponse is used to return response
type EventServiceResponse struct {
	OdataContext                      string                        `json:"@odata.context,omitempty"`
	Etag                              string                        `json:"@odata.etag,omitempty"`
	ID                                string                        `json:"Id"`
	OdataID                           string                        `json:"@odata.id"`
	OdataType                         string                        `json:"@odata.type"`
	Name                              string                        `json:"Name"`
	Description                       string                        `json:"Description,omitempty"`
	Actions                           Actions                       `json:"Actions,omitempty"`
	DeliveryRetryAttempts             int                           `json:"DeliveryRetryAttempts"`
	DeliveryRetryIntervalSeconds      int                           `json:"DeliveryRetryIntervalSeconds"`
	EventFormatTypes                  []string                      `json:"EventFormatTypes"`
	EventTypesForSubscription         []string                      `json:"EventTypesForSubscription"` // Deprecated v1.3
	RegistryPrefixes                  []string                      `json:"RegistryPrefixes"`
	ResourceTypes                     []string                      `json:"ResourceTypes"`
	ServerSentEventURI                string                        `json:"ServerSentEventUri,omitempty"`
	ServiceEnabled                    bool                          `json:"ServiceEnabled,omitempty"`
	SSEFilterPropertiesSupported      *SSEFilterPropertiesSupported `json:"SSEFilterPropertiesSupported,omitempty"`
	Status                            Status                        `json:"Status,omitempty"`
	SubordinateResourcesSupported     bool                          `json:"SubordinateResourcesSupported,omitempty"`
	Subscriptions                     Subscriptions                 `json:"Subscriptions,omitempty"`
	Oem                               Oem                           `json:"Oem,omitempty"`
	IncludeOriginOfConditionSupported bool                          `json:"IncludeOriginOfConditionSupported,omitempty"`
	SMTP                              *SMTP                         `json:"SMTP,omitempty"`
}

// SMTP is for SMTP event delivery
type SMTP struct {
	Authentication     string `json:"Authentication,omitempty"`
	ConnectionProtocol string `json:"ConnectionProtocol,omitempty"`
	FromAddress        string `json:"FromAddress,omitempty"`
	Password           string `json:"Password,omitempty"`
	Port               int    `json:"Port,omitempty"`
	ServerAddress      string `json:"ServerAddress,omitempty"`
	ServiceEnabled     bool   `json:"ServiceEnabled,omitempty"`
	Username           string `json:"Username,omitempty"`
}

//SSEFilterPropertiesSupported defines set propertis that are supported in the
//$filter query parameter for the ServerSentEventUri
type SSEFilterPropertiesSupported struct {
	EventFormatType        bool `json:"EventFormatType"`
	EventType              bool `json:"EventType"` //Deprecated v1.3
	MessageID              bool `json:"MessageId"`
	MetricReportDefinition bool `json:"MetricReportDefinition"`
	OriginResource         bool `json:"OriginResource"`
	RegistryPrefix         bool `json:"RegistryPrefix"`
	ResourceType           bool `json:"ResourceType"`
	SubordinateResources   bool `json:"SubordinateResources"`
}

//Subscriptions containes link to a resource
type Subscriptions struct {
	OdataID string `json:"@odata.id"`
}

//Status struct definition
type Status struct {
	Health       string `json:"Health"`
	HealthRollup string `json:"HealthRollup"`
	State        string `json:"State"`
}

//Actions struct definition
type Actions struct {
	SubmitTestEvent Action `json:"#EventService.SubmitTestEvent"`
	Oem             Oem    `json:"Oem"`
}

//Action struct definition
type Action struct {
	Target          string   `json:"target"`
	AllowableValues []string `json:"EventType@Redfish.AllowableValues"`
}

//Oem struct definition placeholder.
type Oem struct {
}

// MutexLock is a struct for mutex lock and Response and hosts
type MutexLock struct {
	Lock     *sync.Mutex
	Hosts    []string
	Response map[string]EventResponse
}

// AddResponse will add the response to the map and host to the hosts slice
func (r *MutexLock) AddResponse(origin, host string, response EventResponse) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	r.Response[origin] = response
	if response.StatusCode == 201 {
		r.Hosts = append(r.Hosts, host)
	}
}

// ReadResponse will read the response from the response
func (r *MutexLock) ReadResponse(subscriptionID string) (response.RPC, []string) {
	var rpcResponse response.RPC
	r.Lock.Lock()
	defer r.Lock.Unlock()
	for _, resp := range r.Response {
		// Sucessfully created subscription
		rpcResponse.StatusCode = int32(resp.StatusCode)
		rpcResponse.Header = map[string]string{
			"Location": "/redfish/v1/EventService/Subscriptions/" + subscriptionID, // TODO make it dynamic
		}
		rpcResponse.Body = resp.Response
	}
	hosts := r.Hosts
	return rpcResponse, hosts
}
