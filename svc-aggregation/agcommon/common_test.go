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

package agcommon

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	"github.com/ODIM-Project/ODIM/lib-utilities/errors"
	"github.com/ODIM-Project/ODIM/svc-aggregation/agmodel"
)

func TestAddConnectionMethods(t *testing.T) {
	var e = DBInterface{
		GetAllKeysFromTableInterface: stubGetAllkeys,
		GetConnectionMethodInterface: stubGetConnectionMethod,
		AddConnectionMethodInterface: stubAddConnectionMethod,
		DeleteInterface:              stubDeleteConnectionMethod,
	}
	config.SetUpMockConfig(t)
	err := e.AddConnectionMethods(config.Data.ConnectionMethodConf)
	assert.Nil(t, err, "err should be nil")
}

var connectionMethod = []string{"/redfish/v1/AggregationService/ConnectionMethods/1234567545691234f",
	"/redfish/v1/AggregationService/ConnectionMethods/1234567545691234g",
	"/redfish/v1/AggregationService/ConnectionMethods/1234567545691234h"}

func stubGetAllkeys(tableName string) ([]string, error) {
	return connectionMethod, nil
}

func stubGetConnectionMethod(key string) (agmodel.ConnectionMethod, *errors.Error) {
	if key == "/redfish/v1/AggregationService/ConnectionMethods/1234567545691234f" {
		return agmodel.ConnectionMethod{
			ConnectionMethodType:    "Redfish",
			ConnectionMethodVariant: "Compute:BasicAuth:GRF:1.0.0",
			Links: agmodel.Links{
				AggregationSources: []agmodel.OdataID{},
			},
		}, nil
	}

	if key == "/redfish/v1/AggregationService/ConnectionMethods/1234567545691234g" {
		return agmodel.ConnectionMethod{
			ConnectionMethodType:    "Redfish",
			ConnectionMethodVariant: "Fabric:XAuthToken:FabricPlugin:1.0.0",
			Links: agmodel.Links{
				AggregationSources: []agmodel.OdataID{
					agmodel.OdataID{OdataID: "/redfish/v1/AggregationService/AggregationSources/1234656881231fg1"},
				},
			},
		}, nil
	}
	return agmodel.ConnectionMethod{
		ConnectionMethodType:    "Redfish",
		ConnectionMethodVariant: "Storage:BasicAuth:Stg1:1.0.0",
		Links: agmodel.Links{
			AggregationSources: []agmodel.OdataID{},
		},
	}, nil
}

func stubAddConnectionMethod(data agmodel.ConnectionMethod, key string) *errors.Error {

	return nil
}

func stubDeleteConnectionMethod(table, key string, dbtype common.DbType) *errors.Error {

	return nil

}
