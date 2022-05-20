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
package chassis

import (
	"testing"
)

func Test_findAllPlugins(t *testing.T) {
	// common.SetUpMockConfig()
	// mngr := RAManager{
	// 	Name:            "Manager",
	// 	ManagerType:     "Service",
	// 	FirmwareVersion: "1.0",
	// 	ID:              "1",
	// 	UUID:            "1",
	// 	State:           "Enabled",
	// }
	// mocAddManagertoDB(mngr)

	// defer func() {
	// 	err := common.TruncateDB(common.InMemory)
	// 	if err != nil {
	// 		t.Fatalf("error: %v", err)
	// 	}
	// }()

	// delete := NewDeleteHandler(plugin.NewClientFactory(&config.URLTranslation{NorthBoundURL: map[string]string{
	// 	"ODIM": "redfish",
	// },
	// 	SouthBoundURL: map[string]string{
	// 		"redfish": "ODIM",
	// 	}}))

}
