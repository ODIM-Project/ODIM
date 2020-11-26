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
	"net/http"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

func Test_get_urp_manager(t *testing.T) {
	testApp, _ := createTestApplication()
	httptest.New(t, testApp).
		GET("/ODIM/v1/Managers/"+TEST_CONFIG.RootServiceUUID).
		WithBasicAuth("admin", "Od!m12$4").
		Expect().
		Status(http.StatusOK).
		ContentType("application/json", "UTF-8").
		JSON().Object().
		ValueEqual("@odata.id", "/ODIM/v1/Managers/"+TEST_CONFIG.RootServiceUUID).
		ValueEqual("Name", _PLUGIN_NAME).
		ValueEqual("UUID", TEST_CONFIG.RootServiceUUID).
		ValueEqual("Id", TEST_CONFIG.RootServiceUUID).
		ValueEqual("FirmwareVersion", TEST_CONFIG.FirmwareVersion)
}
