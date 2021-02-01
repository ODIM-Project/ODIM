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

package rest

import (
	"testing"

	"github.com/kataras/iris/v12/httptest"
	"github.com/stretchr/testify/require"
)

func Test_get_not_empty_chassis_collection(t *testing.T) {
	testApp, testRedis := createTestApplication()
	//should be returned
	require.NoError(t, testRedis.Set("Chassis:/ODIM/v1/Chassis/1", ""))
	require.NoError(t, testRedis.Set("Chassis:/ODIM/v1/Chassis/2", ""))
	//should not be returned
	require.NoError(t, testRedis.Set("CONTAINS:Chassis:/ODIM/v1/Chassis/2", ""))
	require.NoError(t, testRedis.Set("CONTAINEDIN:Chassis:/ODIM/v1/Chassis/2", ""))

	httptest.New(t, testApp).
		GET("/ODIM/v1/Chassis/").
		WithBasicAuth("admin", "Od!m12$4").
		Expect().Status(httptest.StatusOK).
		ContentType("application/json", "UTF-8").
		JSON().Object().
		Path(`$.Members..["@odata.id"]`).Array().ContainsOnly("/ODIM/v1/Chassis/1", "/ODIM/v1/Chassis/2")
}
