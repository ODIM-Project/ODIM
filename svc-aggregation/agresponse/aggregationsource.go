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

package agresponse

import (
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// AggregationSourceResponse defines the response for AggregationSource
type AggregationSourceResponse struct {
	response.Response
	HostName string      `json:"HostName"`
	UserName string      `json:"UserName"`
	Links    interface{} `json:"Links"`
}

// AggregateResponse defines the response for aggregate
type AggregateResponse struct {
	response.Response
	Elements []string `json:"Elements"`
}
