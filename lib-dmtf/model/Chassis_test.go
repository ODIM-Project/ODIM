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

package model

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
)

var chassisJSON = "{\"@odata.context\":\"/redfish/v1/$metadata#Chassis.Chassis\",\"@odata.id\":\"/redfish/v1/Chassis/1/\",\"@odata.type\":\"#Chassis.v1_16_0.Chassis\",\"@odata.etag\":\"W/\\\"ACCE5EFE\\\"\",\"Id\":\"1\",\"Description\":\"\",\"Name\":\"Computer System Chassis\",\"AssetTag\":\"\",\"ChassisType\":\"RackMount\",\"DepthMm\":0,\"EnvironmentalClass\":\"\",\"HeightMm\":0,\"IndicatorLED\":\"On\",\"Manufacturer\":\"HPE\",\"Model\":\"ProLiant DL380 Gen10\",\"PartNumber\":\"\",\"PowerState\":\"\",\"SerialNumber\":\"2M291101JZ\",\"SKU\":\"868704-B21\",\"UUID\":\"\",\"WeightKg\":0,\"WidthMm\":0,\"Links\":{\"ComputerSystems\":[{\"@odata.id\":\"/redfish/v1/Systems/1/\"}],\"ManagedBy\":[{\"@odata.id\":\"/redfish/v1/Managers/1/\"}]},\"Location\":{\"@odata.id\":\"\"},\"LogServices\":{\"@odata.id\":\"\"},\"Assembly\":{\"@odata.id\":\"\"},\"NetworkAdapters\":{\"@odata.id\":\"/redfish/v1/Chassis/1/NetworkAdapters/\"},\"PCIeSlots\":{\"@odata.id\":\"\"},\"PhysicalSecurity\":{\"IntrusionSensor\":\"\",\"IntrusionSensorNumber\":0,\"IntrusionSensorReArm\":\"\"},\"Power\":{\"@odata.id\":\"/redfish/v1/Chassis/1/Power/\"},\"Sensors\":{\"@odata.id\":\"\"},\"Status\":{\"Health\":\"OK\",\"State\":\"Starting\",\"Oem\":{}},\"Thermal\":{\"@odata.id\":\"/redfish/v1/Chassis/1/Thermal/\"}}"

func TestChassis_SaveInMemory(t *testing.T) {
	var Chas, Chas1 Chassis

	if err := config.SetUpMockConfig(t); err != nil {
		t.Fatalf("fatal: error while trying to collect mock db config: %v", err)
	}

	if err := json.Unmarshal([]byte(chassisJSON), &Chas); err != nil {
		t.Fatalf("error while trying to unmarshal the data: %v", err)
		return
	}

	if err := Chas.SaveInMemory("1234"); err != nil {
		t.Fatalf("error while trying to persist the data: %v", err.Error())
		return
	}
	connPool, errs := common.GetDBConnection(common.InMemory)
	if errs != nil {
		t.Fatalf("error while db connection pool: %v", errs.Error())
		return
	}

	chassisJSONFromDB, _ := connPool.Read("chassis", "1234")

	if errs = connPool.Delete("chassis", "1234"); errs != nil {
		t.Fatalf("error while deleting data: %v", errs)
	}

	if err := json.Unmarshal([]byte(chassisJSONFromDB), &Chas1); err != nil {
		t.Fatalf("error while trying to unmarshal the data: %v", err)
		return
	}

	if !reflect.DeepEqual(Chas1, Chas) {
		t.Errorf("Got = %#v, Want %#v", Chas1, Chas)
	}
}
