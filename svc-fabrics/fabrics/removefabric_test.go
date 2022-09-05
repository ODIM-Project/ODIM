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

package fabrics

import (
	"net/http"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	fabricsproto "github.com/ODIM-Project/ODIM/lib-utilities/proto/fabrics"
	"gotest.tools/assert"
)

func TestRemoveFabric(t *testing.T) {
	Token.Tokens = make(map[string]string)
	common.SetUpMockConfig()
	defer func() {
		err := common.TruncateDB(common.OnDisk)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	}()
	err := mockPluginData(t, "CFM", "XAuthToken", "9091")
	if err != nil {
		t.Fatalf("Error in creating mock DeviceData :%v", err)
	}
	err = mockFabricData("d72dade0-c35a-984c-4859-1108132d72da", "CFM")
	if err != nil {
		t.Fatalf("Error in creating mockFabricData :%v", err)
	}
	req := &fabricsproto.AddFabricRequest{
		OriginResource: "/redfish/v1/Fabrics/d72dade0-c35a-984c-4859-1108132d72da",
		Address:        "10.10.10.10",
	}

	resp := RemoveFabric(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusOK, "should be same")

	resp = RemoveFabric(req)
	assert.Equal(t, int(resp.StatusCode), http.StatusInternalServerError, "should be same")

}
