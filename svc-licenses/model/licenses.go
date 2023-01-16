//(C) Copyright [2022] Hewlett Packard Enterprise Development LP
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

package model

import (
	"context"
	"net/http"
)

// PluginContactRequest ...
type PluginContactRequest struct {
	Token            string
	OID              string
	DeviceInfo       interface{}
	BasicAuth        map[string]string
	ContactClient    func(context.Context, string, string, string, string, interface{}, map[string]string) (*http.Response, error)
	PostBody         interface{}
	Plugin           Plugin
	HTTPMethodType   string
	LoginCredentials map[string]string
}

// ResponseStatus ...
type ResponseStatus struct {
	StatusCode    int32
	StatusMessage string
	MsgArgs       []interface{}
}

// Plugin ...
type Plugin struct {
	IP                string
	Port              string
	Username          string
	Password          []byte
	ID                string
	PluginType        string
	PreferredAuthType string
}

// Target ...
type Target struct {
	ManagerAddress string `json:"ManagerAddress"`
	Password       []byte `json:"Password"`
	UserName       string `json:"UserName"`
	PostBody       []byte `json:"PostBody"`
	DeviceUUID     string `json:"DeviceUUID"`
	PluginID       string `json:"PluginID"`
}

// Elements struct used for storing system details
type Elements struct {
	Elements []OdataIDLinks `json:"Elements,omitempty"`
}

// OdataIDLinks struct used for storing odata id
type OdataIDLinks struct {
	OdataID string `json:"@odata.id,omitempty"`
}
