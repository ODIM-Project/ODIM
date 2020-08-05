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

package system

import (
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/response"
)

// GetAggregationSourceCollection is to fetch all the AggregationSourceURI uri's and returns with created collection
// of AggregationSource data from odim
func GetAggregationSourceCollection() response.RPC {
	// It need to be removed after the backend is implemented
	return response.RPC{
		StatusCode: http.StatusNotImplemented,
	}
}

// GetAggregationSource is used  to fetch the AggregationSource with given aggregation source uri
//and returns AggregationSource
func GetAggregationSource(reqURI string) response.RPC {
	// It need to be removed after the backend is implemented
	return response.RPC{
		StatusCode: http.StatusNotImplemented,
	}
}
