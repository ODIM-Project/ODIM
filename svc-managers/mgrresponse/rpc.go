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

//Package mgrresponse ...
package mgrresponse

// RPC defines the reponse which account-session service returns back as
// part of the RPC call.
//
// StatusCode defines the status code of the requested service operation.
// StatusMessage defines the message regarding the status of the requested operation.
// Header defines the headers required to create a proper response from the api gate way.
// Body defines the actual response of the requested service operation.
type RPC struct {
	Header        map[string]string
	StatusCode    int32
	StatusMessage string
	Body          interface{}
}
