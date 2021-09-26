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

var CSystemJSON = "{\"@odata.context\":\"/redfish/v1/$metadata#ComputerSystem.ComputerSystem\",\"@odata.id\":\"/redfish/v1/Systems/1/\",\"@odata.type\":\"#ComputerSystem.v1_14_1.ComputerSystem\",\"@odata.etag\":\"W/3356AE4A\", \"Id\":\"1\",\"Description\":\"\",\"Name\":\"Computer System\",\"AssetTag\":\"\",\"BiosVersion\":\"U30 v1.46 (10/02/2018)\",\"HostName\":\"\",\"IndicatorLED\":\"Off\",\"Manufacturer\":\"HPE\",\"Model\":\"ProLiant DL380 Gen10\",\"PartNumber\":\"\",\"PowerRestorePolicy\":\"\",\"PowerState\":\"On\",\"SerialNumber\":\"2M291101JZ\",\"SKU\":\"868704-B21\",\"SubModel\":\"\",\"SystemType\":\"Physical\",\"UUID\":\"37383638-3430-4D32-3239-313130314A5A\",\"HostingRoles\":null,\"PCIeDevices\":null,\"PCIeFunctions\":null,\"Bios\":{\"@odata.id\":\"/redfish/v1/systems/1/bios/\"},\"Boot\":{\"AliasBootOrder\":null,\"BootNext\":\"\",\"BootOptions\":{\"@odata.id\":\"\"}, \"BootOrder\":null,\"BootOrderPropertySelection\":\"\",\"BootSourceOverrideEnabled\":\"Disabled\",\"BootSourceOverrideMode\":\"UEFI\",\"BootSourceOverrideTarget\":\"None\",\"Certificates\":{\"@odata.id\":\"\"},\"UefiTargetBootSourceOverride\":\"None\"},\"EthernetInterfaces\":{\"@odata.id\":\"/redfish/v1/Systems/1/EthernetInterfaces/\"},\"HostedServices\":{\"Oem\":{},\"StorageServices\":{\"@odata.id\":\"\"}},\"HostWatchdogTimer\":{\"FunctionEnabled\":false,\"Oem\":{},\"Status\":{\"Health\":\"\",\"State\":\"\",\"Oem\":{}},\"TimeoutAction\":\"\",\"WarningAction\":\"\"},\"Links\":{\"Chassis\":[{\"@odata.id\":\"/redfish/v1/Chassis/1/\"}],\"ManagedBy\":[{\"@odata.id\":\"/redfish/v1/Managers/1/\"}]},\"LogServices\":{\"@odata.id\":\"/redfish/v1/Systems/1/LogServices/\"},\"Memory\":{\"@odata.id\":\"/redfish/v1/Systems/1/Memory/\"},\"MemoryDomains\":{\"@odata.id\":\"\"},\"MemorySummary\":{\"MemoryMirroring\":\"\",\"TotalSystemMemoryGiB\":384,\"TotalSystemPersistentMemoryGiB\":0,\"Status\":{\"Health\":\"\",\"HealthRollup\":\"OK\",\"State\":\"\",\"Oem\":{}}},\"NetworkInterfaces\":{\"@odata.id\":\"/redfish/v1/Systems/1/NetworkInterfaces/\"},\"Processors\":{\"@odata.id\":\"/redfish/v1/Systems/1/Processors/\"},\"ProcessorSummary\":{\"Count\":2,\"LogicalProcessorCount\":0,\"Model\":\"Intel(R) Xeon(R) Gold 6152 CPU @ 2.10GHz\",\"Metrics\":{\"@odata.id\":\"\"},\"Status\":{\"Health\":\"\",\"HealthRollup\":\"OK\",\"State\":\"\",\"Oem\":{}}},\"Redundancy\":[{\"@odata.id\":\"\"}],\"SecureBoot\":{\"@odata.id\":\"/redfish/v1/Systems/1/SecureBoot/\"},\"SimpleStorage\":{\"@odata.id\":\"\"},\"Status\":{\"Health\":\"OK\",\"State\":\"Starting\",\"Oem\":{}},\"Storage\":{\"@odata.id\":\"/redfish/v1/Systems/1/Storage/\"},\"TrustedModules\":[{\"FirmwareVersion\":\"\",\"FirmwareVersion2\":\"\",\"InterfaceType\":\"\",     \"InterfaceTypeSelection\":\"\",\"Oem\":{},\"Status\":{\"Health\":\"\",\"State\":\"Absent\",\"Oem\":{}}}]}"

func TestComputerSystem_SaveInMemory(t *testing.T) {

	var CSystem, CSystem1 ComputerSystem

	if err := config.SetUpMockConfig(t); err != nil {
		t.Fatalf("fatal: error while trying to collect mock db config: %v", err)
	}

	if err := json.Unmarshal([]byte(CSystemJSON), &CSystem); err != nil {
		t.Fatalf("error while trying to unmarshal the data: %v", err)
		return
	}

	if err := CSystem.SaveInMemory("1234"); err != nil {
		t.Fatalf("error while trying to persist the data: %v", err.Error())
		return
	}

	connPool, err := common.GetDBConnection(common.InMemory)
	if err != nil {
		t.Fatalf("error while db connection pool: %v", err.Error())
		return
	}

	CSystemJSONFromDB, _ := connPool.Read("computersystem", "1234")

	if err := connPool.Delete("computersystem", "1234"); err != nil {
		t.Fatalf("error while deleting data: %v", err.Error())
		return
	}

	if err := json.Unmarshal([]byte(CSystemJSONFromDB), &CSystem1); err != nil {
		t.Fatalf("error while trying to unmarshal the data: %v", err)
		return
	}

	if !reflect.DeepEqual(CSystem1, CSystem) {
		t.Errorf("Got = %#v, Want %#v", CSystem1, CSystem)
	}
}
