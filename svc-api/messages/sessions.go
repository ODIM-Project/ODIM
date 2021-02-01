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

//Package messages ...
package messages

// SessionCreateRequest ...
type SessionCreateRequest struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"`
}

// SessionCreateResponse ...
type SessionCreateResponse struct {
	StatusCode   int
	Message      string
	SessionID    string
	SessionToken string
	Header       map[string]string
}

// SessionResponse ...
type SessionResponse struct {
	StatusCode int
	Message    string
	Header     map[string]string
	Body       string
}
