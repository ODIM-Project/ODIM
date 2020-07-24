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
package response

import (
	"testing"
)

func TestCreateGenericResponse(t *testing.T) {

	response := &Response{
		OdataType:    "TestOdataType",
		OdataID:      "/redfish/v1/test/",
		OdataContext: "/redfish/v1/$metadata#Test.Test",
		Description:  "A test response",
		ID:           "Test",
		Name:         "Test",
		Message:      "A test response message",
		Severity:     "OK",
		Resolution:   "Test resolution",
	}

	tests := []struct {
		name        string
		code        string
		message     string
		messageargs []string
		noargs      int
	}{
		{
			name:    Success,
			code:    Success,
			message: "Successfully Completed Request",
		},
		{
			name:    Created,
			code:    Created,
			message: "The resource has been created successfully",
		},
		{
			name:    AccountRemoved,
			code:    AccountRemoved,
			message: "The account was successfully removed.",
		},
		{
			name:    AccountModified,
			code:    AccountModified,
			message: "The account was successfully modified.",
		},
		{
			name:    ResourceRemoved,
			code:    ResourceRemoved,
			message: "The resource has been removed successfully.",
		},
		{
			name:    ResourceCreated,
			code:    ResourceCreated,
			message: "The resource has been created successfully.",
		},
		{
			name:        TaskStarted,
			code:        TaskStarted,
			message:     "The task with id 1234 has started.",
			messageargs: []string{"1234"},
			noargs:      1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response.MessageID = ""
			response.Message = ""
			response.NumberOfArgs = tt.noargs
			response.MessageArgs = tt.messageargs
			response.CreateGenericResponse(tt.code)
			if response.MessageID != tt.code {
				t.Errorf("CreateGenericResponse() = %v, want %v", response.MessageID, tt.code)
			}
			if response.Message != tt.message {
				t.Errorf("CreateGenericResponse() = %v, want %v", response.Message, tt.message)
			}
		})
	}
}
