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

// Package dphandler ..
package dphandler

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

func TestRedfishEvents(t *testing.T) {
	reqBody := []byte(`{"foo": "bar"}`)
	req, err := http.NewRequest("POST", "/redfish/events", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a mock HTTP response writer
	res := httptest.NewRecorder()

	// Call the RedfishEvents function with the mock request and response
	RedfishEvents(res, req)

	req.Header.Set("X-Forwarded-For", "414")

	// Call the RedfishEvents function with the mock request and response
	RedfishEvents(res, req)

	// Check that the response status code is correct
	if status := res.Code; status != http.StatusOK {
		t.Errorf("RedfishEvents returned wrong status code: got %v, want %v", status, http.StatusOK)
	}
	req, err = http.NewRequest("POST", "/redfish/events", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("X-Forwarded-For", "414")

	// Call the RedfishEvents function with the mock request and response
	RedfishEvents(res, req)

}
