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
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (c *ClientMock) Get(uri string) response.RPC {
	args := c.Called(uri)

	var r response.RPC
	if arg0 := args.Get(0); arg0 != nil {
		r = arg0.(response.RPC)
	}
	return r
}

func (c *ClientMock) Post(uri string, body *json.RawMessage) response.RPC {
	panic("implement me")
}

func (c *ClientMock) Delete(uri string) response.RPC {
	panic("implement me")
}

func (c *ClientMock) Patch(uri string, body *json.RawMessage) response.RPC {
	panic("implement me")
}
