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

//Package handle ...
package handle

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/ODIM-Project/ODIM/lib-utilities/common"
	"github.com/ODIM-Project/ODIM/lib-utilities/config"
	errResponse "github.com/ODIM-Project/ODIM/lib-utilities/response"
	srv "github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-api/models"
	"github.com/ODIM-Project/ODIM/svc-api/response"
	iris "github.com/kataras/iris/v12"
)

//GetVersion is an API handler method, which build the response body and headers for /redfish API
func GetVersion(ctx iris.Context) {
	Version := models.Version{
		V1: "/redfish/v1/",
	}
	SetResponseHeaders(ctx, nil)
	ctx.JSON(Version)
}

//ServiceRoot defines getService function
type ServiceRoot struct {
	getService func([]string, string) models.ServiceRoot
}

//InitServiceRoot func returns ServiceRoot
func InitServiceRoot() ServiceRoot {
	return ServiceRoot{
		getService: getService,
	}
}

//getService method takes list of string as parameter and returns Serviceroot struct with assigned values
func getService(microServices []string, uuid string) models.ServiceRoot {
	serviceRoot := models.ServiceRoot{
		OdataType:      "#ServiceRoot.v1_9_0.ServiceRoot",
		ID:             "RootService",
		Name:           "Root Service",
		RedfishVersion: "1.11.1",
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
	options := srv.Service.Options()
	reg := options.Registry
	for _, microService := range microServices {
		servicePath := "/redfish/v1/" + microService
		switch microService {
		case "AccountService":
			serviceNodes, err := reg.GetService(srv.AccountSession)
			if err == nil {
				if len(serviceNodes) != 0 {
					serviceRoot.AccountService = &models.Service{OdataID: servicePath}
				}
			}
		case "EventService":
			serviceNodes, err := reg.GetService(srv.Events)
			if err == nil {
				if len(serviceNodes) != 0 {
					serviceRoot.EventService = &models.Service{OdataID: servicePath}
				}
			}
		case "SessionService":
			serviceNodes, err := reg.GetService(srv.AccountSession)
			if err == nil {
				if len(serviceNodes) != 0 {
					serviceRoot.SessionService = &models.Service{OdataID: servicePath}
				}
			}
		case "JSONSchemas":
			serviceRoot.JSONSchemas = &models.Service{OdataID: servicePath}
		case "Systems":
			serviceNodes, err := reg.GetService(srv.Systems)

			if err == nil {
				if len(serviceNodes) != 0 {
					serviceRoot.Systems = &models.Service{OdataID: servicePath}
				}
			}
		case "Chassis":
			serviceNodes, err := reg.GetService(srv.Systems)
			if err == nil {
				if len(serviceNodes) != 0 {
					serviceRoot.Chassis = &models.Service{OdataID: servicePath}
				}
			}
		case "TaskService":
			serviceNodes, err := reg.GetService(srv.Tasks)
			if err == nil {
				if len(serviceNodes) != 0 {
					serviceRoot.Tasks = &models.Service{OdataID: servicePath}
				}
			}
		case "AggregationService":
			serviceNodes, err := reg.GetService(srv.Aggregator)
			if err == nil {
				if len(serviceNodes) != 0 {
					serviceRoot.AggregationService = &models.Service{OdataID: servicePath}
				}
			}
		case "Fabrics":
			serviceNodes, err := reg.GetService(srv.Fabrics)
			if err == nil {
				if len(serviceNodes) != 0 {
					serviceRoot.Fabrics = &models.Service{OdataID: servicePath}
				}
			}

		case "Managers":
			serviceNodes, err := reg.GetService(srv.Managers)
			if err == nil {
				if len(serviceNodes) != 0 {
					serviceRoot.Managers = &models.Service{OdataID: servicePath}
				}
			}

		case "UpdateService":
			serviceNodes, err := reg.GetService(srv.Update)
			if err == nil {
				if len(serviceNodes) != 0 {
					serviceRoot.UpdateService = &models.Service{OdataID: servicePath}
				}
			}
		case "TelemetryService":
			serviceNodes, err := reg.GetService(srv.Telemetry)
			if err == nil {
				if len(serviceNodes) != 0 {
					serviceRoot.TelemetryService = &models.Service{OdataID: servicePath}
				}
			}
		}
	}

	return serviceRoot
}

//GetServiceRoot builds response body and headers for /redfish/v1
func (s *ServiceRoot) GetServiceRoot(ctx iris.Context) {
	services := config.Data.EnabledServices
	uuid := config.Data.RootServiceUUID
	serviceRoot := s.getService(services, uuid)

	var headers = map[string]string{
		"Allow":             "GET",
		"Cache-Control":     "no-cache",
		"Link":              "</redfish/v1/SchemaStore/en/ServiceRoot.json/>; rel=describedby",
		"Transfer-Encoding": "chunked",
	}
	SetResponseHeaders(ctx, headers)
	ctx.JSON(serviceRoot)
}

//GetOdata builds response body and headers for /redfish/v1/odata
func GetOdata(ctx iris.Context) {
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
	ctx.Gzip(true)
	var odataheaders = map[string]string{
		"Allow":             "GET",
		"Cache-Control":     "no-cache",
		"Transfer-Encoding": "chunked",
	}
	SetResponseHeaders(ctx, odataheaders)
	ctx.JSON(Odata)
}

//GetMetadata build response body and headers for the GET operation on /redfish/v1/$metadata
func GetMetadata(ctx iris.Context) {
	Metadata := models.Metadata{
		Version:   "4.0",
		Xmlnsedmx: "http://docs.oasis-open.org/odata/ns/edmx",
		TopReference: []models.Reference{
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/AccountService_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "AccountService"},
					models.Include{Namespace: "AccountService.v1_6_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/AddressPool_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "AddressPool.v1_1_0"},
					models.Include{Namespace: "AddressPool.v1_1_0a"},
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
					models.Include{Namespace: "Aggregate.v1_0_0"},
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
					models.Include{Namespace: "AggregationService.v1_0_0"},
					models.Include{Namespace: "AggregationService.v1_0_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/AggregationSource_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "AggregationSource"},
					models.Include{Namespace: "AggregationSource.v1_0_0"},
					models.Include{Namespace: "AggregationSource.v1_1_0"},
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
					models.Include{Namespace: "AttributeRegistry.v1_0_0"},
					models.Include{Namespace: "AttributeRegistry.v1_1_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Bios_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Bios"},
					models.Include{Namespace: "Bios.v1_0_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/BootOption_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "BootOption"},
					models.Include{Namespace: "BootOption.v1_0_1"},
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
					models.Include{Namespace: "Chassis.v1_6_0"},
					models.Include{Namespace: "Chassis.v1_10_2"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ChassisCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ChassisCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ComputerSystem_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ComputerSystem"},
					models.Include{Namespace: "ComputerSystem.v1_4_0"},
					models.Include{Namespace: "ComputerSystem.v1_10_0"},
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
					models.Include{Namespace: "ConnectionMethod.v1_0_0"},
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
					models.Include{Namespace: "Drive.v1_8_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Endpoint_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Endpoint.v1_3_1"},
					models.Include{Namespace: "Endpoint.v1_4_0"},
					models.Include{Namespace: "Endpoint.v1_5_0"},
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
					models.Include{Namespace: "EthernetInterface.v1_4_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/EthernetInterfaceCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "EthernetInterfaceCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/EthernetInterfaces_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "EthernetInterfaces"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/EventDestination_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "EventDestination"},
					models.Include{Namespace: "EventDestination.v1_7_0"},
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
					models.Include{Namespace: "EventService.v1_5_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Fabric_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Fabric.v1_0_5"},
					models.Include{Namespace: "Fabric.v1_1_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/FabricCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "FabricCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/FirmwareInventory_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "FirmwareInventory"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/FirmwareInventoryCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "FirmwareInventoryCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/HostInterface_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "HostInterface"},
					models.Include{Namespace: "HostInterface.v1_1_1"},
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
					models.Include{Namespace: "LogEntry.v1_1_0"},
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
					models.Include{Namespace: "LogService.v1_0_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Manager_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Manager"},
					models.Include{Namespace: "Manager.v1_3_3"},
					models.Include{Namespace: "Manager.v1_5_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ManagerAccount_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ManagerAccount"},
					models.Include{Namespace: "ManagerAccount.v1_4_0"},
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
					models.Include{Namespace: "ManagerNetworkProtocol.v1_0_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Memory_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Memory"},
					models.Include{Namespace: "Memory.v1_7_0"},
					models.Include{Namespace: "Memory.v1_7_1"},
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
					models.Include{Namespace: "MessageRegistry.v1_0_0"},
					models.Include{Namespace: "MessageRegistry.v1_1_3"},
					models.Include{Namespace: "MessageRegistry.v1_2_0"},
					models.Include{Namespace: "MessageRegistry.v1_3_0"},
					models.Include{Namespace: "MessageRegistry.v1_3_1"},
					models.Include{Namespace: "MessageRegistry.v1_3_2"},
					models.Include{Namespace: "MessageRegistry.v1_3_3"},
					models.Include{Namespace: "MessageRegistry.v1_4_0"},
					models.Include{Namespace: "MessageRegistry.v1_4_1"},
					models.Include{Namespace: "MessageRegistry.v1_4_2"},
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
					models.Include{Namespace: "MessageRegistryFile.v1_1_0"},
					models.Include{Namespace: "MessageRegistryFile.v1_1_1"},
					models.Include{Namespace: "MessageRegistryFile.v1_1_2"},
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
					models.Include{Namespace: "NetworkAdapter.v1_2_0"},
					models.Include{Namespace: "NetworkAdapter.v1_3_0"},
					models.Include{Namespace: "NetworkAdapter.v1_4_0"},
					models.Include{Namespace: "NetworkAdapter.v1_5_0"},
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
					models.Include{Namespace: "NetworkInterface.v1_1_1"},
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
					models.Include{Namespace: "PCIeDevice.v1_4_0"},
					models.Include{Namespace: "PCIeDevice.v1_5_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Port_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Port"},
					models.Include{Namespace: "Port.v1_1_3"},
					models.Include{Namespace: "Port.v1_2_0"},
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
					models.Include{Namespace: "Power.v1_3_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/PrivilegeRegistry_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "PrivilegeRegistry"},
					models.Include{Namespace: "PrivilegeRegistry.v1_0_0"},
					models.Include{Namespace: "PrivilegeRegistry.v1_1_4"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Processor_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Processor"},
					models.Include{Namespace: "Processor.v1_0_0"},
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
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Resource_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Resource"},
					models.Include{Namespace: "Resource.v1_8_3"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Role_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Role"},
					models.Include{Namespace: "Role.v1_2_4"},
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
					models.Include{Namespace: "SecureBoot.v1_0_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/ServiceRoot_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "ServiceRoot"},
					models.Include{Namespace: "ServiceRoot.v1_1_0"},
					models.Include{Namespace: "ServiceRoot.v1_5_0"},
					models.Include{Namespace: "ServiceRoot.v1_9_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Session_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Session"},
					models.Include{Namespace: "Session.v1_2_1"},
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
					models.Include{Namespace: "SessionService.v1_1_6"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Settings_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Settings"},
					models.Include{Namespace: "Settings.v1_0_0"},
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
					models.Include{Namespace: "Storage.v1_7_1"},
					models.Include{Namespace: "Storage.v1_8_0"},
					models.Include{Namespace: "Storage.v1_8_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/StorageCollection_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "StorageCollection"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Switch_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Switch"},
					models.Include{Namespace: "Switch.v1_2_0"},
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
					models.Include{Namespace: "Task.v1_5_0"},
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
					models.Include{Namespace: "TaskService.v1_1_4"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Thermal_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Thermal"},
					models.Include{Namespace: "Thermal.v1_1_0"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/UpdateService_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "UpdateService"},
					models.Include{Namespace: "UpdateService.v1_1_1"},
					models.Include{Namespace: "UpdateService.v1_6_0"},
					models.Include{Namespace: "UpdateService.v1_8_1"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/VirtualMedia_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "VirtualMedia"},
					models.Include{Namespace: "VirtualMedia.v1_2_0"},
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
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/Zone_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "Zone.v1_3_0"},
					models.Include{Namespace: "Zone.v1_4_0"},
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
					models.Include{Namespace: "SerialInterface.v1_1_7"},
				},
			},
			models.Reference{URI: "http://redfish.dmtf.org/schemas/v1/TelemetryService_v1.xml",
				TopInclude: []models.Include{
					models.Include{Namespace: "TelemetryService.v1_1_7"},
				},
			},
		},
	}
	ctx.Gzip(true)

	var headers = map[string]string{
		"Allow":             "GET",
		"Cache-Control":     "no-cache",
		"Transfer-Encoding": "chunked",
		"Content-type":      "application/xml; charset=utf-8",
	}
	xmlData, _ := xml.Marshal(Metadata)
	SetResponseHeaders(ctx, headers)
	ctx.Write(xmlData)

}

// Registry defines Auth which helps with authorization
type Registry struct {
	Auth func(string, []string, []string) errResponse.RPC
}

//GetRegistryFileCollection is show available collection of registry files.
func (r *Registry) GetRegistryFileCollection(ctx iris.Context) {

	// Authorize the request here
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, errResponse.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error      headers
		ctx.JSON(&response.Body)
		return
	}
	authResp := r.Auth(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Error("error while trying to authorize token")
		ctx.StatusCode(int(authResp.StatusCode))
		SetResponseHeaders(ctx, authResp.Header)
		ctx.JSON(authResp.Body)
		return
	}

	//Get the Registrystore location
	var headers = map[string]string{
		"Allow":             "GET",
		"Cache-Control":     "no-cache",
		"Link":              "</redfish/v1/SchemaStore/en/MessageRegistryFileCollection.json/>; rel=describedby",
		"Transfer-Encoding": "chunked",
	}
	// Get all available file names in the registry store directory in a list
	registryStore := config.Data.RegistryStorePath
	regFiles, err := ioutil.ReadDir(registryStore)
	if err != nil {
		log.Fatal(err.Error())
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
	regFileKeys, err := models.GetAllRegistryFileNamesFromDB("Registries")
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
		OdataID:      "/redfish/v1/Registries/",
		OdataType:    "#MessageRegistryFileCollection.MessageRegistryFileCollection",
		Name:         "Registry File Repository",
		Description:  "Registry Repository",
		MembersCount: len(listMembers),
		Members:      listMembers,
	}
	SetResponseHeaders(ctx, headers)
	ctx.JSON(regCollectionResp)
}

//GetMessageRegistryFileID this is for giving deatiled information about the file and its locations.
func (r *Registry) GetMessageRegistryFileID(ctx iris.Context) {
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
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

	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, errResponse.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error      headers
		ctx.JSON(&response.Body)
		return
	}
	authResp := r.Auth(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Error("error while trying to authorize token")
		ctx.StatusCode(int(authResp.StatusCode))
		SetResponseHeaders(ctx, authResp.Header)
		ctx.JSON(authResp.Body)
		return
	}
	var headers = map[string]string{
		"Allow":             "GET",
		"Content-type":      "application/json; charset=utf-8", //   TODO: add all error headers
		"Cache-Control":     "no-cache",
		"Link":              "</redfish/v1/SchemaStore/en/MessageRegistryFile.json/>; rel=describedby",
		"Transfer-Encoding": "chunked",
	}
	//Get the Registrystore location
	registryStore := config.Data.RegistryStorePath
	regFiles, err := ioutil.ReadDir(registryStore)

	reqRegistryFileName := regFileID + ".json"
	if err != nil {
		log.Error(err.Error())
	}
	// Constuct the registry file names slice
	var regFileNames []string
	for _, regFile := range regFiles {
		regFileNames = append(regFileNames, regFile.Name())
	}
	locationURI := ""
	// Registry file from DB and append
	regFileKeys, err := models.GetAllRegistryFileNamesFromDB("Registries")
	if err != nil {
		// log Critical message but proceed
		log.Error("error: while trying to get the Registry files (Keys/FileNames only)from DB")
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
		responseHeader := map[string]string{
			"Content-type": "application/json; charset=utf-8", //   TODO: add all error headers
		}
		SetResponseHeaders(ctx, responseHeader)
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusNotFound, errResponse.ResourceNotFound, errorMessage, []interface{}{"RegistryFile", regFileID}, nil)
		ctx.StatusCode(http.StatusNotFound) // TODO: add error      headers
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
	SetResponseHeaders(ctx, headers)
	ctx.JSON(resp)
}

//GetMessageRegistryFile this is to retreve the message registry file itself.
func (r *Registry) GetMessageRegistryFile(ctx iris.Context) {
	sessionToken := ctx.Request().Header.Get("X-Auth-Token")
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
	if sessionToken == "" {
		errorMessage := "error: no X-Auth-Token found in request header"
		responseHeader := map[string]string{
			"Content-type": "application/json; charset=utf-8", //   TODO: add all error headers
		}
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusUnauthorized, errResponse.NoValidSession, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusUnauthorized) // TODO: add error      headers
		SetResponseHeaders(ctx, responseHeader)
		ctx.JSON(&response.Body)
		return
	}
	authResp := r.Auth(sessionToken, []string{common.PrivilegeLogin}, []string{})
	if authResp.StatusCode != http.StatusOK {
		log.Error("error while trying to authorize token")
		ctx.StatusCode(int(authResp.StatusCode))
		SetResponseHeaders(ctx, authResp.Header)
		ctx.JSON(authResp.Body)
		return
	}
	var headers = map[string]string{
		"Allow":             "GET",
		"Content-type":      "application/json; charset=utf-8", //   TODO: add all error headers
		"Cache-Control":     "no-cache",
		"Transfer-Encoding": "chunked",
	}
	//Get the Registrystore location

	registryStore := config.Data.RegistryStorePath
	registryStore = strings.TrimSuffix(registryStore, "/")
	regFilePath := registryStore + "/" + regFileID
	// read the file from file system
	content, err := ioutil.ReadFile(regFilePath)
	if err != nil {
		// Check if this file is in DB
		content, err = models.GetRegistryFile("Registries", regFileID)
		if content == nil {
			// file Not found, send 404 error
			log.Error("got error while retreiving fom DB")
			errorMessage := "error: Resource not found"
			responseHeader := map[string]string{
				"Content-type": "application/json; charset=utf-8", //   TODO: add all error headers
			}
			log.Error(errorMessage)
			response := common.GeneralError(http.StatusNotFound, errResponse.ResourceNotFound, errorMessage, []interface{}{"RegistryFile", regFileID}, nil)
			ctx.StatusCode(http.StatusNotFound) // TODO: add error      headers
			SetResponseHeaders(ctx, responseHeader)
			ctx.JSON(&response.Body)
			return
		}
	}
	var data interface{}
	log.Error("Before Unmarshalling Data")
	err = json.Unmarshal(content, &data)
	if err != nil {
		//return fmt.Errorf("error while trying to unmarshal the config data: %v", err)
		log.Error(err.Error())
		errorMessage := "error: Resource not found"
		responseHeader := map[string]string{
			"Content-type": "application/json; charset=utf-8", //   TODO: add all error headers
		}
		log.Error(errorMessage)
		response := common.GeneralError(http.StatusInternalServerError, errResponse.InternalError, errorMessage, nil, nil)
		ctx.StatusCode(http.StatusInternalServerError)
		SetResponseHeaders(ctx, responseHeader)
		ctx.JSON(&response.Body)
		return
	}
	SetResponseHeaders(ctx, headers)
	ctx.JSON(data)

}

// fillMethodNotAllowedErrorResponse fills the status code and status message for MethodNotAllowed responses
func fillMethodNotAllowedErrorResponse(ctx iris.Context) {
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
	ctx.JSON(errArgs.CreateGenericErrorResponse())
	return
}

// AsMethodNotAllowed holds Method to throw 405 Method not allowed on Account Service URLs
func AsMethodNotAllowed(ctx iris.Context) {
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// SsMethodNotAllowed holds builds reponse for the unallowed http operation on Session Service URLs and returns 405 error.
func SsMethodNotAllowed(ctx iris.Context) {
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// SystemsMethodNotAllowed holds builds reponse for the unallowed http operation on Systems URLs and returns 405 error.
func SystemsMethodNotAllowed(ctx iris.Context) {
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
	case "/redfish/v1/Systems/" + systemID + "/Storage/" + storageid + "/Volumes/":
		ctx.ResponseWriter().Header().Set("Allow", "POST, GET")
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

// ManagersMethodNotAllowed holds builds reponse for the unallowed http operation on Managers URLs and returns 405 error.
func ManagersMethodNotAllowed(ctx iris.Context) {
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
	default:
		ctx.ResponseWriter().Header().Set("Allow", "GET")
	}

	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// TsMethodNotAllowed holds builds reponse for the unallowed http operation on Task Service URLs and returns 405 error.
func TsMethodNotAllowed(ctx iris.Context) {
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// ChassisMethodNotAllowed holds builds reponse for the unallowed http operation on Chassis URLs and returns 405 error.
func ChassisMethodNotAllowed(ctx iris.Context) {
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// RegMethodNotAllowed holds builds reponse for the unallowed http operation on Registries URLs and returns 405 error.
func RegMethodNotAllowed(ctx iris.Context) {
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// EvtMethodNotAllowed holds builds reponse for the unallowed http operation on Events URLs and returns 405 error.
func EvtMethodNotAllowed(ctx iris.Context) {
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
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// FabricsMethodNotAllowed holds builds reponse for the unallowed http operation on Fabrics URLs and returns 405 error.
func FabricsMethodNotAllowed(ctx iris.Context) {
	ctx.ResponseWriter().Header().Set("Allow", "GET")
	fillMethodNotAllowedErrorResponse(ctx)
	return
}

// AggregateMethodNotAllowed holds builds reponse for the unallowed http operation on Aggregate URLs and returns 405 error.
func AggregateMethodNotAllowed(ctx iris.Context) {
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
