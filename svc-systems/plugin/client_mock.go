/*
 * Copyright (c) 2020 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package plugin

import (
	"encoding/json"
	"net/http"

	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/stretchr/testify/mock"
)

// ClientMock helps to mock response APIs
type ClientMock struct {
	mock.Mock
}

// Get mocks response of GET APis
func (c *ClientMock) Get(uri string, opts ...CallOption) response.RPC {
	args := c.Called(uri, opts)

	var r response.RPC
	if arg0 := args.Get(0); arg0 != nil {
		r = arg0.(response.RPC)
	}
	return r
}

// Post mocks response of POST APIs
func (c *ClientMock) Post(uri string, body *json.RawMessage) response.RPC {
	// TODO: Implement this
	return response.RPC{
		StatusCode: http.StatusNotImplemented,
	}
}

// Delete mocks response of Delete APIs
func (c *ClientMock) Delete(uri string) response.RPC {
	// TODO: Implement this
	return response.RPC{
		StatusCode: http.StatusNotImplemented,
	}
}

// Patch mocks response of Patch APIs
func (c *ClientMock) Patch(uri string, body *json.RawMessage) response.RPC {
	if uri == "/redfish/v1/Managers/3ccb5c71-0e00-4d14-93bb-8d125c030f27" {
		return response.RPC{
			StatusCode: http.StatusOK,
		}
	}
	// TODO: Implement this
	return response.RPC{
		StatusCode: http.StatusNotImplemented,
	}
}
