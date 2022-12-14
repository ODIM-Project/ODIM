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

package rpc

import (
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/stretchr/testify/assert"
)

func Test_auth(t *testing.T) {
	authFunc := func(sessionToken string, privileges, oemPrivileges []string) (response.RPC, error) {
		return response.RPC{}, nil
	}
	callback := func() response.RPC { return response.RPC{} }
	resp := auth(authFunc, "", []string{}, callback)
	assert.Equal(t, http.StatusUnauthorized, int(resp.StatusCode), "status should be StatusUnauthorized")
}
