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

// Package handle ...
package handle

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	l "github.com/ODIM-Project/ODIM/lib-utilities/logs"
	errResponse "github.com/ODIM-Project/ODIM/lib-utilities/response"
	srv "github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-api/models"
	"github.com/ODIM-Project/ODIM/svc-api/response"
	iris "github.com/kataras/iris/v12"
)

// GetVersion is an API handler method, which build the response body and headers for /redfish API
func GetVersion(ctx iris.Context) {
	defer ctx.Next()
	Version := models.Version{
		V1: "/redfish/v1/",
	}
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	common.SetResponseHeader(ctx, nil)
	ctx.JSON(Version)
}

// ServiceRoot defines getService function
type ServiceRoot struct {
	getService func([]string, string) models.ServiceRoot
}

// InitServiceRoot func returns ServiceRoot
func InitServiceRoot() ServiceRoot {
	return ServiceRoot{
		getService: getService,
	}
}

// getService method takes list of string as parameter and returns Serviceroot struct with assigned values
func getService(microServices []string, uuid string) models.ServiceRoot {
	serviceRoot := models.ServiceRoot{
		OdataType:      "#ServiceRoot.v1_16_0.ServiceRoot",
		ID:             "RootService",
		Name:           "Root Service",
		RedfishVersion: "1.18.1",
		UUID:           uuid, //TODO: persistence of the uuid should be discussed.
		OdataContext:   "/redfish/v1/$metadata#ServiceRoot.ServiceRoot",
		OdataID:        "/redfish/v1/",
		Links: models.Links{
			Sessions: models.Sessions{
				OdataID: "/redfish/v1/SessionService/Sessions"},
		},
		Registries: &models.Service{OdataID: "/redfish/v1/Registries"},
	}
	// To discover the services we need registry
	//Get Service options to retrive the Registry from it.
	for microService := range srv.GetEnabledServiceList() {
		servicePath := "/redfish/v1/" + microService
		switch microService {
		case "AccountService":
			serviceRoot.AccountService = &models.Service{OdataID: servicePath}
		case "EventService":
			serviceRoot.EventService = &models.Service{OdataID: servicePath}
		case "SessionService":
			serviceRoot.SessionService = &models.Service{OdataID: servicePath}
		case "JSONSchemas":
			serviceRoot.JSONSchemas = &models.Service{OdataID: servicePath}
		case "Systems":
			serviceRoot.Systems = &models.Service{OdataID: servicePath}
		case "Chassis":
			serviceRoot.Chassis = &models.Service{OdataID: servicePath}

		case "TaskService":
			serviceRoot.Tasks = &models.Service{OdataID: servicePath}

		case "AggregationService":
			serviceRoot.AggregationService = &models.Service{OdataID: servicePath}
		case "Fabrics":
			serviceRoot.Fabrics = &models.Service{OdataID: servicePath}

		case "Managers":
			serviceRoot.Managers = &models.Service{OdataID: servicePath}

		case "UpdateService":
			serviceRoot.UpdateService = &models.Service{OdataID: servicePath}

		case "TelemetryService":
			serviceRoot.TelemetryService = &models.Service{OdataID: servicePath}

		case "CompositionService":
			serviceRoot.CompositionService = &models.Service{OdataID: servicePath}

		case "LicenseService":
			serviceRoot.LicenseService = &models.Service{OdataID: servicePath}

		}
	}

	return serviceRoot
}

// GetServiceRoot builds response body and headers for /redfish/v1
func (s *ServiceRoot) GetServiceRoot(ctx iris.Context) {
	defer ctx.Next()
	services := config.Data.EnabledServices
	uuid := config.Data.RootServiceUUID
	serviceRoot := s.getService(services, uuid)

	var headers = map[string]string{
		"Allow": "GET",
		"Link":  "</redfish/v1/SchemaStore/en/ServiceRoot.json/>; rel=describedby",
	}
	common.SetResponseHeader(ctx, headers)
	ctx.JSON(serviceRoot)
}

// GetOdata builds response body and headers for /redfish/v1/odata
func GetOdata(ctx iris.Context) {
	defer ctx.Next()
	Odata := models.Odata{
		RedfishCopyright: "Copyright © 2014-2015 Distributed Management Task Force, Inc. (DMTF). All rights reserved.",
		OdataContext:     "/redfish/v1/$metadata",
	}
	var list []string
	list = append(list, "AccountService", "Service", "JsonSchemas", "SessionService")

	for _, service := range list {
		var serviceURL string
		if service == "Service" {
			Odata.Value = append(Odata.Value, &models.Value{Name: service, Kind: "Singleton", URL: "/redfish/v1/"})
		} else if service == "JsonSchemas" {
			Odata.Value = append(Odata.Value, &models.Value{Name: service, Kind: "Singleton", URL: "/redfish/v1/Schemas"})
		} else if service == "Sessions" {
			Odata.Value = append(Odata.Value, &models.Value{Name: service, Kind: "Singleton", URL: "/redfish/v1/SessionService/Sessions/"})
		} else {
			serviceURL = "/redfish/v1/" + service
			Odata.Value = append(Odata.Value, &models.Value{Name: service, Kind: "Singleton", URL: serviceURL})
		}
	}
	var odataheaders = map[string]string{
		"Allow": "GET",
	}
	common.SetResponseHeader(ctx, odataheaders)
	ctx.JSON(Odata)
}

// GetMetadata build response body and headers for the GET operation on /redfish/v1/$metadata
func GetMetadata(ctx iris.Context) {
	defer ctx.Next()
	Metadata := models.Metadata{
		Version:   "4.0",
		Xmlnsedmx: "http://docs.oasis-open.org/odata/ns/edmx",
		TopReference: []models.Reference{
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/AccountService_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "AccountService"},
					models.Include{Namespace: "AccountService.v1_13_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/AddressPool_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "AddressPool"},
					models.Include{Namespace: "AddressPool.v1_2_2"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/AddressPoolCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "AddressPoolCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Aggregate_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Aggregate"},
					models.Include{Namespace: "Aggregate.v1_0_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/AggregateCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "AggregateCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/AggregationService_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "AggregationService"},
					models.Include{Namespace: "AggregationService.v1_0_2"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/AggregationSource_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "AggregationSource"},
					models.Include{Namespace: "AggregationSource.v1_3_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/AggregationSourceCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "AggregationSourceCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/AttributeRegistry_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "AttributeRegistry"},
					models.Include{Namespace: "AttributeRegistry.v1_3_6"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Bios_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Bios"},
					models.Include{Namespace: "Bios.v1_2_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/BootOption_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "BootOption"},
					models.Include{Namespace: "BootOption.v1_0_4"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/BootOptionCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "BootOptionCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Chassis_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Chassis"},
					models.Include{Namespace: "Chassis.v1_23_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ChassisCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ChassisCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/CompositionService_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "CompositionService"},
					models.Include{Namespace: "CompositionService.v1_2_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ComputerSystem_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ComputerSystem"},
					models.Include{Namespace: "ComputerSystem.v1_20_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ComputerSystemCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ComputerSystemCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ConnectionMethod_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ConnectionMethod"},
					models.Include{Namespace: "ConnectionMethod.v1_1_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ConnectionMethodCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ConnectionMethodCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Drive_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Drive"},
					models.Include{Namespace: "Drive.v1_17_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Endpoint_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Endpoint"},
					models.Include{Namespace: "Endpoint.v1_8_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/EndpointCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "EndpointCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/EthernetInterface_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "EthernetInterface"},
					models.Include{Namespace: "EthernetInterface.v1_10_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/EthernetInterfaceCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "EthernetInterfaceCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/EventDestination_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "EventDestination"},
					models.Include{Namespace: "EventDestination.v1_13_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/EventDestinationCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "EventDestinationCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/EventService_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "EventService"},
					models.Include{Namespace: "EventService.v1_10_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Fabric_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Fabric.v1_0_0"},
					models.Include{Namespace: "Fabric.v1_3_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/FabricCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "FabricCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/HostInterface_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "HostInterface"},
					models.Include{Namespace: "HostInterface.v1_3_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/HostInterfaceCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "HostInterfaceCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/IPAddresses_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "IPAddresses"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/LogEntry_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "LogEntry"},
					models.Include{Namespace: "LogEntry.v1_15_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/LogEntryCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "LogEntryCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/LogService_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "LogService"},
					models.Include{Namespace: "LogService.v1_4_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Manager_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Manager"},
					models.Include{Namespace: "Manager.v1_18_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ManagerAccount_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ManagerAccount"},
					models.Include{Namespace: "ManagerAccount.v1_10_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ManagerAccountCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ManagerAccountCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ManagerCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ManagerCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ManagerNetworkProtocol_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ManagerNetworkProtocol"},
					models.Include{Namespace: "ManagerNetworkProtocol.v1_9_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Memory_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Memory"},
					models.Include{Namespace: "Memory.v1_17_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/MemoryCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "MemoryCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/MemoryDomainCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "MemoryDomainCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Message_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Message"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/MessageRegistry_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "MessageRegistry"},
					models.Include{Namespace: "MessageRegistry.v1_6_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/MessageRegistryCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "MessageRegistryCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/MessageRegistryFile_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "MessageRegistryFile"},
					models.Include{Namespace: "MessageRegistryFile.v1_1_3"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/MessageRegistryFileCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "MessageRegistryFileCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/NetworkAdapter_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "NetworkAdapter"},
					models.Include{Namespace: "NetworkAdapter.v1_9_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/NetworkAdapterCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "NetworkAdapterCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/NetworkInterface_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "NetworkInterface"},
					models.Include{Namespace: "NetworkInterface.v1_2_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/NetworkInterfaceCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "NetworkInterfaceCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/PCIeDevice_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "PCIeDevice"},
					models.Include{Namespace: "PCIeDevice.v1_11_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/PCIeFunction_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "PCIeFunction"},
					models.Include{Namespace: "PCIeFunction.v1_5_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Port_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Port"},
					models.Include{Namespace: "Port.v1_9_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/PortCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "PortCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Power_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Power"},
					models.Include{Namespace: "Power.v1_7_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/PrivilegeRegistry_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "PrivilegeRegistry"},
					models.Include{Namespace: "PrivilegeRegistry.v1_1_4"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Processor_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Processor"},
					models.Include{Namespace: "Processor.v1_18_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ProcessorCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ProcessorCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/RedfishExtensions_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "RedfishExtensions.v1_0_0", Alias: "Redfish"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Redundancy_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Redundancy"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ResourceBlock_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ResourceBlock"},
					models.Include{Namespace: "ResourceBlock.v1_4_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ResourceBlockCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ResourceBlockCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Resource_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Resource"},
					models.Include{Namespace: "Resource.v1_16_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Role_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Role"},
					models.Include{Namespace: "Role.v1_3_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/RoleCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "RoleCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/SecureBoot_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "SecureBoot"},
					models.Include{Namespace: "SecureBoot.v1_1_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ServiceRoot_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ServiceRoot"},
					models.Include{Namespace: "ServiceRoot.v1_16_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Session_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Session"},
					models.Include{Namespace: "Session.v1_6_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/SessionCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "SessionCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/SessionService_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "SessionService"},
					models.Include{Namespace: "SessionService.v1_1_8"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Settings_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Settings"},
					models.Include{Namespace: "Settings.v1_3_5"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/SoftwareInventoryCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "SoftwareInventoryCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Storage_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Storage"},
					models.Include{Namespace: "Storage.v1_15_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/StorageCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "StorageCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/StorageController_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "StorageController"},
					models.Include{Namespace: "StorageController.v1_7_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/StorageControllerCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "StorageControllerCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Switch_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Switch"},
					models.Include{Namespace: "Switch.v1_9_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/SwitchCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "SwitchCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Task_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Task"},
					models.Include{Namespace: "Task.v1_7_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/TaskCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "TaskCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/TaskService_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "TaskService"},
					models.Include{Namespace: "TaskService.v1_2_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Thermal_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Thermal"},
					models.Include{Namespace: "Thermal.v1_7_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/UpdateService_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "UpdateService"},
					models.Include{Namespace: "UpdateService.v1_11_3"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/VirtualMedia_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "VirtualMedia"},
					models.Include{Namespace: "VirtualMedia.v1_6_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/VirtualMediaCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "VirtualMediaCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/VLanNetworkInterface_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "VLanNetworkInterface"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/VolumeCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "VolumeCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Volume_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Volume"},
					models.Include{Namespace: "Volume.v1_9_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Zone_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Zone"},
					models.Include{Namespace: "Zone.v1_6_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ZoneCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ZoneCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/SerialInterfaceCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "SerialInterfaceCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/SerialInterface_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "SerialInterface"},
					models.Include{Namespace: "SerialInterface.v1_1_8"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/PCIeSlots_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "PCIeSlots"},
					models.Include{Namespace: "PCIeSlots.v1_5_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/NetworkPortCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "NetworkPortCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/NetworkPort_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "NetworkPort"},
					models.Include{Namespace: "NetworkPort.v1_4_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/NetworkDeviceFunction_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "NetworkDeviceFunction"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/NetworkDeviceFunctionCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "NetworkDeviceFunctionCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/LogServiceCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "LogServiceCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/TelemetryService_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "TelemetryService"},
					models.Include{Namespace: "TelemetryService.v1_3_2"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/TriggersCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "TriggersCollection"},
				},
			},

			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Triggers_v1.xml",

				TopInclude: []models.Include{
					models.Include{Namespace: "Triggers"},
					models.Include{Namespace: "Triggers.v1_3_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/MetricDefinitionCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "MetricDefinitionCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/MetricDefinition_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "MetricDefinition"},
					models.Include{Namespace: "MetricDefinition.v1_3_2"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/MetricReportDefinitionCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "MetricReportDefinitionCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/MetricReportDefinition_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "MetricReportDefinition"},
					models.Include{Namespace: "MetricReportDefinition.v1_4_3"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/MetricReportCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "MetricReportCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/MetricReport_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "MetricReport"},
					models.Include{Namespace: "MetricReport.v1_5_0"},
				},
			},
			models.Reference{URI: "https://redfish.dmtf.org/schemas/v1/LicenseCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "LicenseCollection"},
				},
			},
			models.Reference{URI: "https://redfish.dmtf.org/schemas/v1/LicenseService_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "LicenseService"},
					models.Include{Namespace: "LicenseService.v1_1_0"},
				},
			},
			models.Reference{URI: "https://redfish.dmtf.org/schemas/v1/License_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "License"},
					models.Include{Namespace: "License.v1_1_1"},
				},
			},
		},
	}

	var headers = map[string]string{
		"Allow":        "GET",
		"Content-type": "application/xml; charset=utf-8",
	}
	xmlData, _ := xml.Marshal(Metadata)
	common.SetResponseHeader(ctx, headers)
	ctx.Write(xmlData)

}

// Registry defines Auth which helps with authorization
type Registry struct {
	Auth func(context.Context, string, []string, []string) (errResponse.RPC, error)
}

// GetRegistryFileCollection is show available collection of registry files.
func (r *Registry) GetRegistryFileCollection(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	// Authorize the request here
	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)
	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
		return
	}
	authResp, err := r.Auth(ctxt, sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		errMsg := "error while trying to authenticate session"
		if err != nil {
			errMsg = errMsg + ": " + err.Error()
		}
		sendAuthErrorResponse(ctxt, ctx, errMsg, authResp)
	}

	//Get the Registrystore location
	var headers = map[string]string{
		"Allow": "GET",
		"Link":  "</redfish/v1/SchemaStore/en/MessageRegistryFileCollection.json/>; rel=describedby",
	}
	// Get all available file names in the registry store directory in a list
	registryStore := config.Data.RegistryStorePath
	regFiles, err := ioutil.ReadDir(registryStore)
	if err != nil {
		l.LogWithFields(ctxt).Fatal(err.Error())
	}
	//Construct the Response body
	var listMembers []response.ListMember
	for _, regFile := range regFiles {
		// checking if any hidden files or non json files, if so we do not show them in the response
		if regFile.Name()[0:1] == "." || !strings.HasSuffix(regFile.Name(), ".json") {
			continue
		}
		regFileID := strings.TrimSuffix(regFile.Name(), ".json")
		member := response.ListMember{
			OdataID: "/redfish/v1/Registries/" + regFileID,
		}
		listMembers = append(listMembers, member)
	}
	// Get Registry file names from db if any
	regFileKeys, err := models.GetAllRegistryFileNamesFromDB(ctxt, "Registries")
	if err != nil {
		// log Critical message but proceed
	}
	for _, regFile := range regFileKeys {
		regFileID := strings.TrimSuffix(regFile, ".json")
		member := response.ListMember{
			OdataID: "/redfish/v1/Registries/" + regFileID,
		}
		listMembers = append(listMembers, member)
	}
	regCollectionResp := response.ListResponse{
		OdataContext: "/redfish/v1/$metadata#MessageRegistryFileCollection.MessageRegistryFileCollection",
		OdataID:      "/redfish/v1/Registries",
		OdataType:    "#MessageRegistryFileCollection.MessageRegistryFileCollection",
		Name:         "Registry File Repository",
		Description:  "Registry Repository",
		MembersCount: len(listMembers),
		Members:      listMembers,
	}
	common.SetResponseHeader(ctx, headers)
	ctx.JSON(regCollectionResp)
}

// GetMessageRegistryFileID this is for giving deatiled information about the file and its locations.
func (r *Registry) GetMessageRegistryFileID(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)
	regFileID := ctx.Params().Get("id")
	if strings.Contains(regFileID, ".json") {
		r.GetMessageRegistryFile(ctx)
		return
	}
	if strings.HasPrefix(regFileID, "#") {
		reqURI := ctx.Request().RequestURI
		// Fetch Registry file ID from request URI
		strArray := strings.Split(reqURI, "/")
		if strings.HasSuffix(reqURI, "/") {
			regFileID = strArray[len(strArray)-2]
		} else {
			regFileID = strArray[len(strArray)-1]
		}
	}
	regFileID = strings.Replace(regFileID, "#", "%23", -1)
	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
		return
	}
	authResp, err := r.Auth(ctxt, sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		errMsg := "error while trying to authenticate session"
		if err != nil {
			errMsg = errMsg + ": " + err.Error()
		}
		sendAuthErrorResponse(ctxt, ctx, errMsg, authResp)
	}
	var headers = map[string]string{
		"Allow": "GET",
		"Link":  "</redfish/v1/SchemaStore/en/MessageRegistryFile.json/>; rel=describedby",
	}
	//Get the Registrystore location
	registryStore := config.Data.RegistryStorePath
	regFiles, err := ioutil.ReadDir(registryStore)

	reqRegistryFileName := regFileID + ".json"
	if err != nil {
		l.LogWithFields(ctxt).Error(err.Error())
	}
	// Constuct the registry file names slice
	var regFileNames []string
	for _, regFile := range regFiles {
		regFileNames = append(regFileNames, regFile.Name())
	}
	locationURI := ""
	// Registry file from DB and append
	regFileKeys, err := models.GetAllRegistryFileNamesFromDB(ctxt, "Registries")
	if err != nil {
		// log Critical message but proceed
		l.LogWithFields(ctxt).Error("error: while trying to get the Registry files (Keys/FileNames only)from DB")
	}
	regFileNames = append(regFileNames, regFileKeys...)
	for _, regFile := range regFileNames {
		if reqRegistryFileName == regFile {
			locationURI = "/redfish/v1/Registries/" + regFile
			break
		}
	}
	if locationURI == "" {
		errorMessage := "error: resource not found"
		l.LogWithFields(ctxt).Error(errorMessage)
		response := common.GeneralError(http.StatusNotFound, errResponse.ResourceNotFound, errorMessage, []interface{}{"RegistryFile", regFileID}, nil)
		common.SetResponseHeader(ctx, response.Header)
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(&response.Body)
		return
	}
	strArr := strings.Split(regFileID, ".")
	registryType := strArr[0]
	//construct the response body
	resp := response.MessageRegistryFileID{
		ID:           regFileID,
		OdataContext: "/redfish/v1/$metadata#MessageRegistryFile.MessageRegistryFile",
		OdataID:      "/redfish/v1/Registries/" + regFileID,
		OdataType:    "#MessageRegistryFile.v1_1_3.MessageRegistryFile",
		Name:         "Registry File Repository",
		Description:  registryType + " Message Registry File Locations",
		Languages:    []string{"en"},
		Registry:     regFileID,
		Location: []response.Location{response.Location{
			Language: "en",
			URI:      locationURI,
		},
		},
	}
	common.SetResponseHeader(ctx, headers)
	ctx.JSON(resp)
}

// GetMessageRegistryFile this is to retreve the message registry file itself.
func (r *Registry) GetMessageRegistryFile(ctx iris.Context) {
	defer ctx.Next()
	ctxt := ctx.Request().Context()
	sessionToken := ctx.Request().Header.Get(AuthTokenHeader)
	regFileID := ctx.Params().Get("id")
	if strings.HasPrefix(regFileID, "#") {
		reqURI := ctx.Request().RequestURI
		// Fetch Registry file ID from request URI
		strArray := strings.Split(reqURI, "/")
		if strings.HasSuffix(reqURI, "/") {
			regFileID = strArray[len(strArray)-2]
		} else {
			regFileID = strArray[len(strArray)-1]
		}
	}
	regFileID = strings.Replace(regFileID, "#", "%23", -1)
	l.LogWithFields(ctxt).Debugf("Retriveing message registry file with file id %s", regFileID)
	if sessionToken == "" {
		errorMessage := invalidAuthTokenErrorMsg
		common.SendInvalidSessionResponse(ctx, errorMessage)
		return
	}
	authResp, err := r.Auth(ctxt, sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		errMsg := "error while trying to authenticate session"
		if err != nil {
			errMsg = errMsg + ": " + err.Error()
		}
		sendAuthErrorResponse(ctxt, ctx, errMsg, authResp)
	}
	var headers = map[string]string{
		"Allow": "GET",
	}
	//Get the Registrystore location

	registryStore := config.Data.RegistryStorePath
	registryStore = strings.TrimSuffix(registryStore, "/")
	regFilePath := registryStore + "/" + regFileID
	// read the file from file system
	content, err := ioutil.ReadFile(regFilePath)
	if err != nil {
		// Check if this file is in DB
		content, err = models.GetRegistryFile(ctxt, "Registries", regFileID)
		if content == nil {
			// file Not found, send 404 error
			l.LogWithFields(ctxt).Error("got error while retreiving fom DB")
			errorMessage := "error: Resource not found"
			l.LogWithFields(ctxt).Error(errorMessage)
			response := common.GeneralError(http.StatusNotFound, errResponse.ResourceNotFound, errorMessage, []interface{}{"RegistryFile", regFileID}, nil)
			common.SetResponseHeader(ctx, response.Header)
			ctx.StatusCode(http.StatusNotFound)
			ctx.JSON(&response.Body)
			return
		}
	}
	var data interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		//return fmt.Errorf("error while trying to unmarshal the config data: %v", err)
		l.LogWithFields(ctxt).Error(err.Error())
		errorMessage := "error: Resource not found"
		l.LogWithFields(ctxt).Error(errorMessage)
		common.SendFailedRPCCallResponse(ctx, errorMessage)
	}
	common.SetResponseHeader(ctx, headers)
	ctx.JSON(data)

}

// fillMethodNotAllowedErrorResponse fills the status code and status message for MethodNotAllowed responses
func fillMethodNotAllowedErrorResponse(ctx iris.Context) {
	defer ctx.Next()
	ctx.StatusCode(http.StatusMethodNotAllowed)
	errArgs := &errResponse.Args{
		Code: errResponse.GeneralError,
		ErrorArgs: []errResponse.ErrArgs{
			errResponse.ErrArgs{
				StatusMessage: errResponse.ActionNotSupported,
				MessageArgs:   []interface{}{ctx.Request().Method},
			},
		},
	}
	common.SetResponseHeader(ctx, nil)
	ctx.JSON(errArgs.CreateGenericErrorResponse())
	return
}

// AsMethodNotAllowed holds Method to throw 405 Method not allowed on Account Service URLs
func AsMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	url := ctx.Request().URL
	path := url.Path
	id := ctx.Params().Get("id")
	switch path {
	case "/redfish/v1/AccountService":
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	case "/redfish/v1/AccountService/Accounts":
		ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	case "/redfish/v1/AccountService/Accounts/" + id:
		ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH, DELETE")
	default:
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	}
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// SsMethodNotAllowed holds builds reponse for the unallowed http operation on Session Service URLs and returns 405 error.
func SsMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	url := ctx.Request().URL
	path := url.Path
	id := ctx.Params().Get("sessionID")
	switch path {
	case "/redfish/v1/SessionService":
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	case "/redfish/v1/SessionService/Sessions":
		ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	case "/redfish/v1/SessionService/Sessions/" + id:
		ctx.ResponseWriter().Header().Set("Allow", "GET, DELETE")
	default:
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	}
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// SystemsMethodNotAllowed holds builds reponse for the unallowed http operation on Systems URLs and returns 405 error.
func SystemsMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	url := ctx.Request().URL
	path := url.Path
	systemID := ctx.Params().Get("id")
	subID := ctx.Params().Get("rid")
	storageid := ctx.Params().Get("id2")
	resourceID := ctx.Params().Get("rid")
	// Extend switch case, when each path, requires different handling
	switch path {
	case "/redfish/v1/Systems/" + systemID:
		ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH")
	case "/redfish/v1/Systems/" + systemID + "/LogServices/" + subID + "Actions":
		ctx.ResponseWriter().Header().Set("Allow", "")
	case "/redfish/v1/Systems/" + systemID + "/LogServices/" + subID + "Actions/LogService.ClearLog":
		ctx.ResponseWriter().Header().Set("Allow", "POST")
	case "/redfish/v1/Systems/" + systemID + "/Storage/" + storageid + "/Volumes":
		ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	case "/redfish/v1/Systems/" + systemID + "/Storage/" + storageid + "/Volumes/" + resourceID:
		ctx.ResponseWriter().Header().Set("Allow", "GET, DELETE")
	default:
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	}

	// Set Allow header for search and filter queries
	if strings.Contains(path, "?") {
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	}

	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// SystemsMethodInvalidURI holds builds reponse for the invalid url operation on Systems URLs and returns 404 error.
func SystemsMethodInvalidURI(ctx iris.Context) {
	defer ctx.Next()
	url := ctx.Request().URL
	ctx.StatusCode(http.StatusNotFound)
	errArgs := &errResponse.Args{
		Code: errResponse.GeneralError,
		ErrorArgs: []errResponse.ErrArgs{
			errResponse.ErrArgs{
				StatusMessage: errResponse.InvalidURI,
				MessageArgs:   []interface{}{url},
			},
		},
	}
	common.SetResponseHeader(ctx, nil)
	ctx.JSON(errArgs.CreateGenericErrorResponse())
	return
}

// CompositionServiceMethodNotAllowed holds builds reponse for the unallowed http operation on Systems URLs and returns 405 error.
func CompositionServiceMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	url := ctx.Request().URL
	path := url.Path
	resourceID := ctx.Params().Get("id")
	// Extend switch case, when each path, requires different handling
	switch path {
	case "/redfish/v1/CompositionService/ResourceBlocks":
		ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	case "/redfish/v1/CompositionServie/ResourceBlocks/" + resourceID:
		ctx.ResponseWriter().Header().Set("Allow", "GET, DELETE")
	case "/redfish/v1/CompositionService/ResourceZones":
		ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	case "/redfish/v1/CompositionService/ResourceZones/" + resourceID:
		ctx.ResponseWriter().Header().Set("Allow", "GET, DELETE")
	default:
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	}

	// Set Allow header for search and filter queries
	if strings.Contains(path, "?") {
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	}

	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// LicenseMethodNotAllowed holds builds reponse for the unallowed http operation on License URLs and returns 405 error.
func LicenseMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	url := ctx.Request().URL
	path := url.Path

	// Extend switch case, when each path, requires different handling
	switch path {
	case "/redfish/v1/LicenseService/Licenses":
		ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	default:
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	}

	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// ManagersMethodNotAllowed holds builds reponse for the unallowed http operation on Managers URLs and returns 405 error.
func ManagersMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	url := ctx.Request().URL
	path := url.Path
	systemID := ctx.Params().Get("id")
	subID := ctx.Params().Get("rid")

	// Extend switch case, when each path, requires different handling
	switch path {
	case "/redfish/v1/Managers/" + systemID + "/LogServices/" + subID + "Actions":
		ctx.ResponseWriter().Header().Set("Allow", "")
	case "/redfish/v1/Managers/" + systemID + "/LogServices/" + subID + "Actions/LogService.ClearLog":
		ctx.ResponseWriter().Header().Set("Allow", "POST")
	case "/redfish/v1/Managers/" + systemID + "/VirtualMedia/" + subID + "/Actions/VirtualMedia.EjectMedia":
		ctx.ResponseWriter().Header().Set("Allow", "POST")
	case "/redfish/v1/Managers/" + systemID + "/VirtualMedia/" + subID + "/Actions/VirtualMedia.InsertMedia":
		ctx.ResponseWriter().Header().Set("Allow", "POST")
	default:
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	}

	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// TsMethodNotAllowed holds builds reponse for the unallowed http operation on Task Service URLs and returns 405 error.
func TsMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// UpdateServiceMethodNotAllowed holds builds reponse for the unallowed http operation on Update Service URLs and returns 405 error.
func UpdateServiceMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// MethodNotAllowed fills status code and status message for MethodNotAllowed responses
func MethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// ChassisMethodNotAllowed holds builds reponse for the unallowed http operation on Chassis URLs and returns 405 error.
func ChassisMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// RegMethodNotAllowed holds builds reponse for the unallowed http operation on Registries URLs and returns 405 error.
func RegMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// EvtMethodNotAllowed holds builds reponse for the unallowed http operation on Events URLs and returns 405 error.
func EvtMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	url := ctx.Request().URL
	path := url.Path

	// Extend switch case, when each path, requires different handling
	switch path {
	case "/redfish/v1/EventService":
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	case "/redfish/v1/EventService/Actions":
		ctx.ResponseWriter().Header().Set("Allow", "")
	case "/redfish/v1/EventService/Actions/EventService.SubmitTestEvent":
		ctx.ResponseWriter().Header().Set("Allow", "POST")
	}
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// AggMethodNotAllowed holds builds reponse for the unallowed http operation on Aggregation Service URLs and returns 405 error.
func AggMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	url := ctx.Request().URL
	path := url.Path
	id := ctx.Params().Get("id")
	switch path {
	case "/redfish/v1/AggregationService":
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	case "/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder":
		ctx.ResponseWriter().Header().Set("Allow", "POST")
	case "/redfish/v1/AggregationService/Actions/AggregationService.Reset":
		ctx.ResponseWriter().Header().Set("Allow", "POST")
	case "/redfish/v1/AggregationService/AggregationSources":
		ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	case "/redfish/v1/AggregationService/AggregationSources/" + id:
		ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH, DELETE")
	case "/redfish/v1/AggregationService/ConnectionMethods":
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	case "/redfish/v1/AggregationService/ConnectionMethods/" + id:
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	default:
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	}
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// FabricsMethodNotAllowed holds builds reponse for the unallowed http operation on Fabrics URLs and returns 405 error.
func FabricsMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// AggregateMethodNotAllowed holds builds reponse for the unallowed http operation on Aggregate URLs and returns 405 error.
func AggregateMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	url := ctx.Request().URL
	path := url.Path
	aggregateID := ctx.Params().Get("id")
	// Extend switch case, when each path, requires different handling
	switch path {
	case "/redfish/v1/AggregationService/Aggregates/":
		ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	case "/redfish/v1/AggregationService/Aggregates/" + aggregateID:
		ctx.ResponseWriter().Header().Set("Allow", "GET, DELETE")
	case "/redfish/v1/AggregationService/Aggregates/" + aggregateID + "Actions/Aggregate.AddElements/":
		ctx.ResponseWriter().Header().Set("Allow", "POST")
	case "/redfish/v1/AggregationService/Aggregates/" + aggregateID + "Actions/Aggregate.RemoveElements/":
		ctx.ResponseWriter().Header().Set("Allow", "POST")
	case "/redfish/v1/AggregationService/Aggregates/" + aggregateID + "Actions/Aggregate.Reset/":
		ctx.ResponseWriter().Header().Set("Allow", "POST")
	case "/redfish/v1/AggregationService/Aggregates/" + aggregateID + "Actions/Aggregate.SetDefaultBootOrder/":
		ctx.ResponseWriter().Header().Set("Allow", "POST")
	}
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// SRMethodNotAllowed holds builds response for the unallowed http operation on service root URLs and returns 405 error.
func SRMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// RoleMethodNotAllowed holds builds response for the unallowed http operation on Role URLs and returns 405 error.
func RoleMethodNotAllowed(ctx iris.Context) {
	defer ctx.Next()
	url := ctx.Request().URL
	path := url.Path
	id := ctx.Params().Get("id")
	switch path {
	case "/redfish/v1/Roles":
		ctx.ResponseWriter().Header().Set("Allow", "GET, POST")
	case "/redfish/v1/Roles/" + id:
		ctx.ResponseWriter().Header().Set("Allow", "GET, PATCH, DELETE")
	}
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// sendAuthErrorResponse writes the response for an RPC call when authentication of session fails
func sendAuthErrorResponse(ctxt context.Context, ctx iris.Context, errorMessage string, authResp errResponse.RPC) {
	l.LogWithFields(ctxt).Error(errorMessage)
	ctx.StatusCode(int(authResp.StatusCode))
	common.SetResponseHeader(ctx, authResp.Header)
	ctx.JSON(authResp.Body)
	return
}
