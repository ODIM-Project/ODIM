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
package update

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/response"
	"github.com/ODIM-Project/ODIM/svc-update/uresponse"
)

func TestGetUpdateService(t *testing.T) {
	successResponse := response.Response{
		OdataType:    "#UpdateService.v1_8_1.UpdateService",
		OdataID:      "/redfish/v1/UpdateService",
		OdataContext: "/redfish/v1/$metadata#UpdateService.UpdateService",
		ID:           "UpdateService",
		Name:         "Update Service",
	}
	successResponse.CreateGenericResponse(response.Success)
	successResponse.Message = ""
	successResponse.MessageID = ""
	successResponse.Severity = ""
	common.SetUpMockConfig()
	tests := []struct {
		name string
		want response.RPC
	}{
		{
			name: "account service enabled",
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Allow":         "GET",
					"Cache-Control": "no-cache",
					"Connection":    "Keep-alive",
					"Content-type":  "application/json; charset=utf-8",
					"Link": "	</redfish/v1/SchemaStore/en/UpdateService.json>; rel=describedby",
					"Transfer-Encoding": "chunked",
					"X-Frame-Options":   "sameorigin",
				},
				Body: uresponse.UpdateService{
					Response: successResponse,
					Status: uresponse.Status{
						State:        "Enabled",
						Health:       "OK",
						HealthRollup: "OK",
					},
					ServiceEnabled: true,
					SoftwareInventory: uresponse.SoftwareInventory{
						OdataID: "/redfish/v1/UpdateService/SoftwareInventory",
					},
					FirmwareInventory: uresponse.FirmwareInventory{
						OdataID: "/redfish/v1/UpdateService/FirmwareInventory",
					},
				},
			},
		},
		{
			name: "account service disabled",
			want: response.RPC{
				StatusCode:    http.StatusOK,
				StatusMessage: response.Success,
				Header: map[string]string{
					"Allow":         "GET",
					"Cache-Control": "no-cache",
					"Connection":    "Keep-alive",
					"Content-type":  "application/json; charset=utf-8",
					"Link": "	</redfish/v1/SchemaStore/en/UpdateService.json>; rel=describedby",
					"Transfer-Encoding": "chunked",
					"X-Frame-Options":   "sameorigin",
				},
				Body: uresponse.UpdateService{
					Response: successResponse,
					Status: uresponse.Status{
						State:        "Disabled",
						Health:       "OK",
						HealthRollup: "OK",
					},
					ServiceEnabled: false,
					SoftwareInventory: uresponse.SoftwareInventory{
						OdataID: "/redfish/v1/UpdateService/SoftwareInventory",
					},
					FirmwareInventory: uresponse.FirmwareInventory{
						OdataID: "/redfish/v1/UpdateService/FirmwareInventory",
					},
				},
			},
		},
	}
	config.Data.EnabledServices = []string{"UpdateService"}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUpdateService()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUpdateService() = %v, want %v", got, tt.want)
			}
		})
		config.Data.EnabledServices = []string{"XXXX"}
	}
}
