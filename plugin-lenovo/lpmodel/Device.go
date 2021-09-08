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

//Package lpmodel ...
package lpmodel

//Device struct definition
type Device struct {
	Host     string `json:"ManagerAddress"`
	Username string `json:"UserName"`
	Password []byte `json:"Password"`
	PostBody []byte `json:"PostBody"`
	Location string `json:"Location"`
}

//EvtSubPost ...
type EvtSubPost struct {
	Destination   string        `json:"Destination"`
	EventTypes    []string      `json:"EventTypes,omitempty"`
	MessageIds    []string      `json:"MessageIds,omitempty"`
	ResourceTypes []string      `json:"ResourceTypes,omitempty"`
	HTTPHeaders   []HTTPHeaders `json:"HttpHeaders"`
	Context       string        `json:"Context"`
	Protocol      string        `json:"Protocol"`
}

//HTTPHeaders ...
type HTTPHeaders struct {
	ContentType string `json:"Content-Type"`
}

// Startup struct recieve request on Startup call
type Startup struct {
	Location   string   `json:"Location"`
	EventTypes []string `json:"EventTypes,omitempty"`
	Device     Device   `json:"Device"`
}
