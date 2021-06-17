//(C) Copyright [2020] Hewlett Packard Enterprise Development LP
//(C) Copyright 2020 Intel Corporation
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

//Package router ...
package router

import (
	srv "github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-api/handle"
	"github.com/ODIM-Project/ODIM/svc-api/middleware"
	"github.com/ODIM-Project/ODIM/svc-api/rpc"
	"github.com/kataras/iris/v12"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
)

//Router method to register API handlers.
func Router() *iris.Application {
	r := handle.RoleRPCs{
		GetAllRolesRPC: rpc.GetAllRoles,
		CreateRoleRPC:  rpc.CreateRole,
		GetRoleRPC:     rpc.GetRole,
		UpdateRoleRPC:  rpc.UpdateRole,
		DeleteRoleRPC:  rpc.DeleteRole,
	}
	a := handle.AccountRPCs{
		GetServiceRPC:     rpc.DoGetAccountServiceRequest,
		CreateRPC:         rpc.DoAccountCreationRequest,
		GetAllAccountsRPC: rpc.DoGetAllAccountRequest,
		GetAccountRPC:     rpc.DoGetAccountRequest,
		UpdateRPC:         rpc.DoUpdateAccountRequest,
		DeleteRPC:         rpc.DoAccountDeleteRequest,
	}
	pc := handle.AggregatorRPCs{
		GetAggregationServiceRPC:                rpc.DoGetAggregationService,
		ResetRPC:                                rpc.DoResetRequest,
		SetDefaultBootOrderRPC:                  rpc.DoSetDefaultBootOrderRequest,
		AddAggregationSourceRPC:                 rpc.DoAddAggregationSource,
		GetAllAggregationSourceRPC:              rpc.DoGetAllAggregationSource,
		GetAggregationSourceRPC:                 rpc.DoGetAggregationSource,
		UpdateAggregationSourceRPC:              rpc.DoUpdateAggregationSource,
		DeleteAggregationSourceRPC:              rpc.DoDeleteAggregationSource,
		CreateAggregateRPC:                      rpc.DoCreateAggregate,
		GetAggregateCollectionRPC:               rpc.DoGetAggregateCollection,
		GetAggregateRPC:                         rpc.DoGeteAggregate,
		DeleteAggregateRPC:                      rpc.DoDeleteAggregate,
		AddElementsToAggregateRPC:               rpc.DoAddElementsToAggregate,
		RemoveElementsFromAggregateRPC:          rpc.DoRemoveElementsFromAggregate,
		ResetAggregateElementsRPC:               rpc.DoResetAggregateElements,
		SetDefaultBootOrderAggregateElementsRPC: rpc.DoSetDefaultBootOrderAggregateElements,
		GetAllConnectionMethodsRPC:              rpc.DoGetAllConnectionMethods,
		GetConnectionMethodRPC:                  rpc.DoGetConnectionMethod,
	}

	s := handle.SessionRPCs{
		CreateSessionRPC:        rpc.DoSessionCreationRequest,
		DeleteSessionRPC:        rpc.DeleteSessionRequest,
		GetSessionRPC:           rpc.GetSessionRequest,
		GetAllActiveSessionsRPC: rpc.GetAllActiveSessionRequest,
		GetSessionServiceRPC:    rpc.GetSessionServiceRequest,
	}

	ts := handle.TaskRPCs{
		DeleteTaskRPC:     rpc.DeleteTaskRequest,
		GetTaskRPC:        rpc.GetTaskRequest,
		GetSubTasksRPC:    rpc.GetSubTasks,
		GetSubTaskRPC:     rpc.GetSubTask,
		GetTaskMonitorRPC: rpc.GetTaskMonitor,
		TaskCollectionRPC: rpc.TaskCollection,
		GetTaskServiceRPC: rpc.GetTaskService,
	}

	system := handle.SystemRPCs{
		GetSystemsCollectionRPC:    rpc.GetSystemsCollection,
		GetSystemRPC:               rpc.GetSystemRequestRPC,
		GetSystemResourceRPC:       rpc.GetSystemResource,
		SystemResetRPC:             rpc.ComputerSystemReset,
		SetDefaultBootOrderRPC:     rpc.SetDefaultBootOrder,
		ChangeBiosSettingsRPC:      rpc.ChangeBiosSettings,
		ChangeBootOrderSettingsRPC: rpc.ChangeBootOrderSettings,
		CreateVolumeRPC:            rpc.CreateVolume,
		DeleteVolumeRPC:            rpc.DeleteVolume,
	}

	cha := handle.ChassisRPCs{
		GetChassisCollectionRPC: rpc.GetChassisCollection,
		GetChassisResourceRPC:   rpc.GetChassisResource,
		GetChassisRPC:           rpc.GetChassis,
		CreateChassisRPC:        rpc.CreateChassis,
		DeleteChassisRPC:        rpc.DeleteChassis,
		UpdateChassisRPC:        rpc.UpdateChassis,
	}

	evt := handle.EventsRPCs{
		GetEventServiceRPC:                 rpc.DoGetEventService,
		CreateEventSubscriptionRPC:         rpc.DoCreateEventSubscription,
		SubmitTestEventRPC:                 rpc.DoSubmitTestEvent,
		GetEventSubscriptionRPC:            rpc.DoGetEventSubscription,
		DeleteEventSubscriptionRPC:         rpc.DoDeleteEventSubscription,
		GetEventSubscriptionsCollectionRPC: rpc.DoGetEventSubscriptionsCollection,
	}

	fab := handle.FabricRPCs{
		GetFabricResourceRPC:    rpc.GetFabricResource,
		UpdateFabricResourceRPC: rpc.UpdateFabricResource,
		DeleteFabricResourceRPC: rpc.DeleteFabricResource,
	}

	manager := handle.ManagersRPCs{
		GetManagersCollectionRPC: rpc.GetManagersCollection,
		GetManagersRPC:           rpc.GetManagers,
		GetManagersResourceRPC:   rpc.GetManagersResource,
	}

	update := handle.UpdateRPCs{
		GetUpdateServiceRPC:               rpc.DoGetUpdateService,
		SimpleUpdateRPC:                   rpc.DoSimpleUpdate,
		StartUpdateRPC:                    rpc.DoStartUpdate,
		GetFirmwareInventoryRPC:           rpc.DoGetFirmwareInventory,
		GetFirmwareInventoryCollectionRPC: rpc.DoGetFirmwareInventoryCollection,
		GetSoftwareInventoryRPC:           rpc.DoGetSoftwareInventory,
		GetSoftwareInventoryCollectionRPC: rpc.DoGetSoftwareInventoryCollection,
	}

	registryFile := handle.Registry{
		Auth: srv.IsAuthorized,
	}

	serviceRoot := handle.InitServiceRoot()

	router := iris.New()

	// Parses the URL and performs URL decoding for path
	router.WrapRouter(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rawURI := r.RequestURI
		parsedURI, err := url.Parse(rawURI)
		if err != nil {
			errMessage := "while trying to parse the URL: " + err.Error()
			log.Error(errMessage)
			return
		}
		path := strings.Replace(rawURI, parsedURI.EscapedPath(), parsedURI.Path, -1)
		r.RequestURI = path
		r.URL.Path = parsedURI.Path
		next(w, r)
	})

	taskmon := router.Party("/taskmon")
	taskmon.SetRegisterRule(iris.RouteSkip)
	taskmon.Get("/{TaskID}", ts.GetTaskMonitor)
	taskmon.Any("/{TaskID}", handle.TsMethodNotAllowed)

	redfish := router.Party("/redfish")
	redfish.SetRegisterRule(iris.RouteSkip)
	redfish.Get("/", handle.GetVersion)

	v1 := redfish.Party("/v1")
	v1.SetRegisterRule(iris.RouteSkip)
	v1.Get("/", serviceRoot.GetServiceRoot)
	v1.Get("/odata", handle.GetOdata)
	v1.Get("/$metadata", handle.GetMetadata)

	registry := v1.Party("/Registries")
	registry.SetRegisterRule(iris.RouteSkip)
	registry.Get("/", registryFile.GetRegistryFileCollection)
	registry.Get("/{id}", registryFile.GetMessageRegistryFileID)
	registry.Any("/", handle.RegMethodNotAllowed)
	registry.Any("/{id}", handle.RegMethodNotAllowed)

	session := v1.Party("/SessionService")
	session.SetRegisterRule(iris.RouteSkip)
	session.Get("/", s.GetSessionService)
	session.Get("/Sessions", middleware.SessionDelMiddleware, s.GetAllActiveSessions)
	session.Get("/Sessions/{sessionID}", middleware.SessionDelMiddleware, s.GetSession)
	session.Post("/Sessions", s.CreateSession)
	session.Delete("/Sessions/{sessionID}", middleware.SessionDelMiddleware, s.DeleteSession)
	session.Any("/", handle.SsMethodNotAllowed)

	account := v1.Party("/AccountService", middleware.SessionDelMiddleware)
	account.SetRegisterRule(iris.RouteSkip)
	account.Get("/", a.GetAccountService)
	account.Get("/Accounts", a.GetAllAccounts)
	account.Get("/Accounts/{id}", a.GetAccount)
	account.Post("/Accounts", a.CreateAccount)
	account.Patch("/Accounts/{id}", a.UpdateAccount)
	account.Delete("/Accounts/{id}", a.DeleteAccount)
	account.Any("/", handle.AsMethodNotAllowed)

	role := account.Party("/Roles", middleware.SessionDelMiddleware)
	role.SetRegisterRule(iris.RouteSkip)
	role.Get("/", r.GetAllRoles)
	role.Get("/{id}", r.GetRole)
	role.Post("/", r.CreateRole)
	role.Patch("/{id}", r.UpdateRole)
	role.Delete("/{id}", r.DeleteRole)

	task := v1.Party("/TaskService", middleware.SessionDelMiddleware)
	task.SetRegisterRule(iris.RouteSkip)
	task.Get("/", ts.GetTaskService)
	task.Get("/Tasks", ts.TaskCollection)
	task.Get("/Tasks/{TaskID}", ts.GetTaskStatus)
	task.Get("/Tasks/{TaskID}/SubTasks", ts.GetSubTasks)
	task.Get("/Tasks/{TaskID}/SubTasks/{subTaskID}", ts.GetSubTask)
	task.Delete("/Tasks/{TaskID}", ts.DeleteTask)
	task.Any("/", handle.TsMethodNotAllowed)
	task.Any("/Tasks", handle.TsMethodNotAllowed)
	task.Any("/Tasks/{TaskID}", handle.TsMethodNotAllowed)

	systems := v1.Party("/Systems", middleware.SessionDelMiddleware)
	systems.SetRegisterRule(iris.RouteSkip)
	systems.Get("/", system.GetSystemsCollection)
	systems.Get("/{id}", system.GetSystem)
	systems.Get("/{id}/Processors", system.GetSystemResource)
	systems.Get("/{id}/Processors/{rid}", system.GetSystemResource)
	systems.Get("/{id}/Memory", system.GetSystemResource)
	systems.Get("/{id}/Memory/{rid}", system.GetSystemResource)
	systems.Get("/{id}/NetworkInterfaces", system.GetSystemResource)
	systems.Get("/{id}/NetworkInterfaces/{rid}", system.GetSystemResource)
	systems.Get("/{id}/MemoryDomains", system.GetSystemResource)
	systems.Get("/{id}/EthernetInterfaces", system.GetSystemResource)
	systems.Get("/{id}/EthernetInterfaces/{rid}", system.GetSystemResource)
	systems.Get("/{id}/EthernetInterfaces/{id2}/VLANS", system.GetSystemResource)
	systems.Get("/{id}/EthernetInterfaces/{id2}/VLANS/{rid}", system.GetSystemResource)
	systems.Get("/{id}/SecureBoot", system.GetSystemResource)
	systems.Get("/{id}/BootOptions", system.GetSystemResource)
	systems.Get("/{id}/BootOptions/{rid}", system.GetSystemResource)
	systems.Get("/{id}/LogServices", system.GetSystemResource)
	systems.Get("/{id}/LogServices/{rid}", system.GetSystemResource)
	systems.Get("/{id}/LogServices/{rid}/Entries", system.GetSystemResource)
	systems.Get("/{id}/LogServices/{rid}/Entries/{rid2}", system.GetSystemResource)
	systems.Post("/{id}/LogServices/{rid}/Actions/LogService.ClearLog", system.GetSystemResource)
	systems.Patch("/{id}", system.ChangeBootOrderSettings)
	systems.Get("/{id}/PCIeDevices/{rid}", system.GetSystemResource)
	systems.Any("/{id}/PCIeDevices/{rid}", handle.SystemsMethodNotAllowed)
	systems.Any("/", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/EthernetInterfaces", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/EthernetInterfaces/{rid}", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/SecureBoot", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/MemoryDomains", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/NetworkInterfaces", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/Memory", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/Processors", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/BootOptions", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/BootOptions/{rid}", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/LogServices", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/LogServices/{rid}", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/LogServices/{rid}/Entries", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/LogServices/{rid}/Entries/{rid2}", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/LogServices/{rid}/Actions", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/LogServices/{rid}/Actions/LogService.ClearLog", handle.SystemsMethodNotAllowed)

	systems.Get("/{id}/Bios", system.GetSystemResource)
	systems.Get("/{id}/Bios/Settings", system.GetSystemResource)
	systems.Patch("/{id}/Bios/Settings", system.ChangeBiosSettings)
	systems.Any("/{id}/Bios", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/Processors/{rid}", handle.SystemsMethodNotAllowed)

	storage := v1.Party("/Systems/{id}/Storage", middleware.SessionDelMiddleware)
	storage.SetRegisterRule(iris.RouteSkip)
	storage.Get("/", system.GetSystemResource)
	storage.Get("/{rid}", system.GetSystemResource)
	storage.Get("/{id2}/Drives/{rid}", system.GetSystemResource)
	storage.Get("/{id2}/Volumes", system.GetSystemResource)
	storage.Post("/{id2}/Volumes", system.CreateVolume)
	storage.Delete("/{id2}/Volumes/{rid}", system.DeleteVolume)
	storage.Get("/{id2}/Volumes/{rid}", system.GetSystemResource)
	storage.Any("/", handle.SystemsMethodNotAllowed)
	storage.Any("/{id2}/Drives/{rid}", handle.SystemsMethodNotAllowed)
	storage.Any("/{rid}", handle.SystemsMethodNotAllowed)
	storage.Any("/{id2}/Volumes", handle.SystemsMethodNotAllowed)
	storage.Any("/{id2}/Volumes/{rid}", handle.SystemsMethodNotAllowed)

	systemsAction := systems.Party("/{id}/Actions", middleware.SessionDelMiddleware)
	systemsAction.SetRegisterRule(iris.RouteSkip)
	systemsAction.Post("/ComputerSystem.Reset", system.ComputerSystemReset)
	systemsAction.Post("/ComputerSystem.SetDefaultBootOrder", system.SetDefaultBootOrder)

	aggregation := v1.Party("/AggregationService", middleware.SessionDelMiddleware)
	aggregation.SetRegisterRule(iris.RouteSkip)
	aggregation.Get("/", pc.GetAggregationService)
	aggregation.Post("/Actions/AggregationService.Reset/", pc.Reset)
	aggregation.Any("/Actions/AggregationService.Reset/", handle.AggMethodNotAllowed)
	aggregation.Post("/Actions/AggregationService.SetDefaultBootOrder/", pc.SetDefaultBootOrder)
	aggregation.Any("/Actions/AggregationService.SetDefaultBootOrder/", handle.AggMethodNotAllowed)
	aggregation.Any("/", handle.AggMethodNotAllowed)
	aggregationSource := aggregation.Party("/AggregationSources", middleware.SessionDelMiddleware)
	aggregationSource.Post("/", pc.AddAggregationSource)
	aggregationSource.Get("/", pc.GetAllAggregationSource)
	aggregationSource.Any("/", handle.AggMethodNotAllowed)
	aggregationSource.Get("/{id}", pc.GetAggregationSource)
	aggregationSource.Patch("/{id}", pc.UpdateAggregationSource)
	aggregationSource.Delete("/{id}", pc.DeleteAggregationSource)
	aggregationSource.Any("/{id}", handle.AggMethodNotAllowed)

	connectionMethods := aggregation.Party("/ConnectionMethods", middleware.SessionDelMiddleware)
	connectionMethods.Get("/", pc.GetAllConnectionMethods)
	connectionMethods.Get("/{id}", pc.GetConnectionMethod)
	connectionMethods.Any("/", handle.AggMethodNotAllowed)
	connectionMethods.Any("/{id}", handle.AggMethodNotAllowed)

	aggregates := aggregation.Party("/Aggregates", middleware.SessionDelMiddleware)
	aggregates.Post("/", pc.CreateAggregate)
	aggregates.Get("/", pc.GetAggregateCollection)
	aggregates.Any("/", handle.AggregateMethodNotAllowed)
	aggregates.Get("/{id}", pc.GetAggregate)
	aggregates.Delete("/{id}", pc.DeleteAggregate)
	aggregates.Any("/{id}", handle.AggregateMethodNotAllowed)
	aggregates.Post("/{id}/Actions/Aggregate.AddElements/", pc.AddElementsToAggregate)
	aggregates.Any("/{id}/Actions/Aggregate.AddElements/", handle.AggregateMethodNotAllowed)
	aggregates.Post("/{id}/Actions/Aggregate.RemoveElements/", pc.RemoveElementsFromAggregate)
	aggregates.Any("/{id}/Actions/Aggregate.RemoveElements/", handle.AggregateMethodNotAllowed)
	aggregates.Post("/{id}/Actions/Aggregate.Reset/", pc.ResetAggregateElements)
	aggregates.Any("/{id}/Actions/Aggregate.Reset/", handle.AggregateMethodNotAllowed)
	aggregates.Post("/{id}/Actions/Aggregate.SetDefaultBootOrder/", pc.SetDefaultBootOrderAggregateElements)
	aggregates.Any("/{id}/Actions/Aggregate.SetDefaultBootOrder/", handle.AggregateMethodNotAllowed)

	chassis := v1.Party("/Chassis", middleware.SessionDelMiddleware)
	chassis.SetRegisterRule(iris.RouteSkip)
	chassis.Get("/", cha.GetChassisCollection)
	chassis.Post("/", cha.CreateChassis)
	chassis.Get("/{id}", cha.GetChassis)
	chassis.Patch("/{id}", cha.UpdateChassis)
	chassis.Delete("/{id}", cha.DeleteChassis)
	chassis.Get("/{id}/NetworkAdapters", cha.GetChassisResource)
	chassis.Get("/{id}/NetworkAdapters/{rid}", cha.GetChassisResource)
	chassis.Get("/{id}/NetworkAdapters/{id2}/NetworkDeviceFunctions", cha.GetChassisResource)
	chassis.Get("/{id}/NetworkAdapters/{id2}/NetworkPorts", cha.GetChassisResource)
	chassis.Get("/{id}/NetworkAdapters/{id2}/NetworkDeviceFunctions/{rid}", cha.GetChassisResource)
	chassis.Get("/{id}/NetworkAdapters/{id2}/NetworkPorts/{rid}", cha.GetChassisResource)
	chassis.Any("/", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}/NetworkAdapters", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}/NetworkAdapters/{rid}", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}/NetworkAdapters/{id2}/NetworkDeviceFunctions", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}/NetworkAdapters/{id2}/NetworkPorts", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}/NetworkAdapters/{id2}/NetworkDeviceFunctions/{rid}", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}/NetworkAdapters/{id2}/NetworkPorts/{rid}", handle.ChassisMethodNotAllowed)
	chassis.Get("/{id}/Assembly", cha.GetChassisResource)
	chassis.Any("/{id}/Assembly", handle.ChassisMethodNotAllowed)
	chassis.Get("/{id}/PCIeSlots", cha.GetChassisResource)
	chassis.Get("/{id}/PCIeSlots/{rid}", cha.GetChassisResource)
	chassis.Any("/{id}/PCIeSlots", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}/PCIeSlots/{rid}", handle.ChassisMethodNotAllowed)
	chassis.Get("/{id}/PCIeDevices", cha.GetChassisResource)
	chassis.Get("/{id}/PCIeDevices/{rid}", cha.GetChassisResource)
	chassis.Any("/{id}/PCIeDevices", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}/PCIeDevices/{rid}", handle.ChassisMethodNotAllowed)
	chassis.Get("/{id}/Sensors", cha.GetChassisResource)
	chassis.Get("/{id}/Sensors/{rid}", cha.GetChassisResource)
	chassis.Any("/{id}/Sensors", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}/Sensors/{rid}", handle.ChassisMethodNotAllowed)
	chassis.Get("/{id}/LogServices", cha.GetChassisResource)
	chassis.Get("/{id}/LogServices/{rid}", cha.GetChassisResource)
	chassis.Get("/{id}/LogServices/{rid}/Entries", cha.GetChassisResource)
	chassis.Get("/{id}/LogServices/{rid}/Entries/{rid2}", cha.GetChassisResource)
	// TODO
	// chassis.Post("/{id}/LogServices/{rid}/Actions/LogService.ClearLog", cha.GetChassisResource)
	chassis.Any("/{id}/LogServices", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}/LogServices/{rid}", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}/LogServices/{rid}/Entries", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}/LogServices/{rid}/Entries/{rid2}", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}/LogServices/{rid}/Actions", handle.ChassisMethodNotAllowed)

	chassisPower := chassis.Party("/{id}/Power")
	chassisPower.SetRegisterRule(iris.RouteSkip)
	chassisPower.Get("/", cha.GetChassisResource)
	chassisPower.Get("#PowerControl/{id1}", cha.GetChassisResource)
	chassisPower.Get("#PowerSupplies/{id1}", cha.GetChassisResource)
	chassisPower.Get("#Redundancy/{id1}", cha.GetChassisResource)
	chassisPower.Any("/", handle.ChassisMethodNotAllowed)
	chassisPower.Any("#PowerControl/{id1}", handle.ChassisMethodNotAllowed)
	chassisPower.Any("#PowerSupplies/{id1}", handle.ChassisMethodNotAllowed)
	chassisPower.Any("#Redundancy/{id1}", handle.ChassisMethodNotAllowed)

	chassisThermal := chassis.Party("/{id}/Thermal")
	chassisThermal.SetRegisterRule(iris.RouteSkip)
	chassisThermal.Get("/", cha.GetChassisResource)
	chassisThermal.Get("#Fans/{id1}", cha.GetChassisResource)
	chassisThermal.Get("#Temperatures/{id1}", cha.GetChassisResource)
	chassisThermal.Any("/", handle.ChassisMethodNotAllowed)
	chassisThermal.Any("#Fans/{id1}", handle.ChassisMethodNotAllowed)
	chassisThermal.Any("#Temperatures/{id1}", handle.ChassisMethodNotAllowed)

	events := v1.Party("/EventService", middleware.SessionDelMiddleware)
	events.SetRegisterRule(iris.RouteSkip)
	events.Get("/", evt.GetEventService)
	events.Get("/Subscriptions", evt.GetEventSubscriptionsCollection)
	events.Get("/Subscriptions/{id}", evt.GetEventSubscription)
	events.Post("/Subscriptions", evt.CreateEventSubscription)
	events.Post("/Actions/EventService.SubmitTestEvent", evt.SubmitTestEvent)
	events.Delete("/Subscriptions/{id}", evt.DeleteEventSubscription)
	events.Any("/", handle.EvtMethodNotAllowed)
	events.Any("/Actions", handle.EvtMethodNotAllowed)
	events.Any("/Actions/EventService.SubmitTestEvent", handle.EvtMethodNotAllowed)
	events.Any("/Subscriptions", handle.EvtMethodNotAllowed)

	fabrics := v1.Party("/Fabrics", middleware.SessionDelMiddleware)
	fabrics.SetRegisterRule(iris.RouteSkip)
	fabrics.Get("/", fab.GetFabricCollection)
	fabrics.Get("/{id}", fab.GetFabric)
	fabrics.Get("/{id}/Switches", fab.GetFabricSwitchCollection)
	fabrics.Get("/{id}/Switches/{switchID}", fab.GetFabricSwitch)
	fabrics.Get("/{id}/Switches/{switchID}/Ports", fab.GetSwitchPortCollection)
	fabrics.Get("/{id}/Switches/{switchID}/Ports/{port_uuid}", fab.GetSwitchPort)
	fabrics.Get("/{id}/Zones/", fab.GetFabricZoneCollection)
	fabrics.Get("/{id}/Endpoints/", fab.GetFabricEndPointCollection)
	fabrics.Get("/{id}/AddressPools/", fab.GetFabricAddressPoolCollection)
	fabrics.Get("/{id}/Zones/{zone_uuid}", fab.GetFabricZone)
	fabrics.Get("/{id}/Endpoints/{endpoint_uuid}", fab.GetFabricEndPoints)
	fabrics.Get("/{id}/AddressPools/{addresspool_uuid}", fab.GetFabricAddressPool)
	fabrics.Put("/{id}/Zones/{zone_uuid}", fab.UpdateFabricResource)
	fabrics.Put("/{id}/Endpoints/{endpoint_uuid}", fab.UpdateFabricResource)
	fabrics.Put("/{id}/AddressPools/{addresspool_uuid}", fab.UpdateFabricResource)
	fabrics.Post("/{id}/Zones", fab.UpdateFabricResource)
	fabrics.Post("/{id}/Endpoints", fab.UpdateFabricResource)
	fabrics.Post("/{id}/AddressPools", fab.UpdateFabricResource)
	fabrics.Patch("/{id}/Zones/{zone_uuid}", fab.UpdateFabricResource)
	fabrics.Patch("/{id}/Endpoints/{endpoint_uuid}", fab.UpdateFabricResource)
	fabrics.Patch("/{id}/AddressPools/{addresspool_uuid}", fab.UpdateFabricResource)
	fabrics.Patch("/{id}/Switches/{switchID}/Ports/{port_uuid}", fab.UpdateFabricResource)
	fabrics.Delete("/{id}/Zones/{zone_uuid}", fab.DeleteFabricResource)
	fabrics.Delete("/{id}/Endpoints/{endpoint_uuid}", fab.DeleteFabricResource)
	fabrics.Delete("/{id}/AddressPools/{addresspool_uuid}", fab.DeleteFabricResource)
	fabrics.Any("/", handle.FabricsMethodNotAllowed)
	fabrics.Any("/{id}", handle.FabricsMethodNotAllowed)
	fabrics.Any("/{id}/Switches", handle.FabricsMethodNotAllowed)
	fabrics.Any("/{id}/Switches/{switchID}", handle.FabricsMethodNotAllowed)
	fabrics.Any("/{id}/Switches/{switchID}/Ports", handle.FabricsMethodNotAllowed)
	fabrics.Any("/{id}/Switches/{switchID}/Ports/{port_uuid}", handle.FabricsMethodNotAllowed)
	fabrics.Any("/{id}/Zones/", handle.FabricsMethodNotAllowed)
	fabrics.Any("/{id}/Endpoints/", handle.FabricsMethodNotAllowed)
	fabrics.Any("/{id}/AddressPools/", handle.FabricsMethodNotAllowed)
	fabrics.Any("/{id}/Zones/{zone_uuid}", handle.FabricsMethodNotAllowed)
	fabrics.Any("/{id}/Endpoints/{endpoint_uuid}", handle.FabricsMethodNotAllowed)
	fabrics.Any("/{id}/AddressPools/{addresspool_uuid}", handle.FabricsMethodNotAllowed)

	managers := v1.Party("/Managers", middleware.SessionDelMiddleware)
	managers.SetRegisterRule(iris.RouteSkip)
	managers.Get("/", manager.GetManagersCollection)
	managers.Get("/{id}", manager.GetManager)
	managers.Get("/{id}/EthernetInterfaces", manager.GetManagersResource)
	managers.Get("/{id}/EthernetInterfaces/{rid}", manager.GetManagersResource)
	managers.Get("/{id}/NetworkProtocol", manager.GetManagersResource)
	managers.Get("/{id}/NetworkProtocol/{rid}", manager.GetManagersResource)
	managers.Get("/{id}/HostInterfaces", manager.GetManagersResource)
	managers.Get("/{id}/HostInterfaces/{rid}", manager.GetManagersResource)

	managers.Get("/{id}/SerialInterface", manager.GetManagersResource)
	managers.Get("/{id}/SerialInterface/{rid}", manager.GetManagersResource)
	managers.Get("/{id}/VirtualMedia", manager.GetManagersResource)
	managers.Get("/{id}/VirtualMedia/{rid}", manager.GetManagersResource)
	managers.Get("/{id}/LogServices", manager.GetManagersResource)
	managers.Get("/{id}/LogServices/{rid}", manager.GetManagersResource)
	managers.Get("/{id}/LogServices/{rid}/Entries", manager.GetManagersResource)
	managers.Get("/{id}/LogServices/{rid}/Entries/{rid2}", manager.GetManagersResource)
	managers.Post("/{id}/LogServices/{rid}/Actions/LogService.ClearLog", manager.GetManagersResource)
	managers.Any("/{id}/LogServices", handle.ManagersMethodNotAllowed)
	managers.Any("/{id}/LogServices/{rid}", handle.ManagersMethodNotAllowed)
	managers.Any("/{id}/LogServices/{rid}/Entries", handle.ManagersMethodNotAllowed)
	managers.Any("/{id}/LogServices/{rid}/Entries/{rid2}", handle.ManagersMethodNotAllowed)
	managers.Any("/{id}/LogServices/{rid}/Actions", handle.ManagersMethodNotAllowed)
	managers.Any("/{id}/LogServices/{rid}/Actions/LogService.ClearLog", handle.ManagersMethodNotAllowed)
	managers.Any("/", handle.ManagersMethodNotAllowed)
	managers.Any("/{id}", handle.ManagersMethodNotAllowed)

	updateService := v1.Party("/UpdateService", middleware.SessionDelMiddleware)
	updateService.SetRegisterRule(iris.RouteSkip)
	updateService.Get("/", update.GetUpdateService)
	updateService.Post("/Actions/UpdateService.SimpleUpdate", update.SimpleUpdate)
	updateService.Post("/Actions/UpdateService.StartUpdate", update.StartUpdate)
	updateService.Get("/FirmwareInventory", update.GetFirmwareInventoryCollection)
	updateService.Get("/FirmwareInventory/{firmwareInventory_id}", update.GetFirmwareInventory)
	updateService.Get("/SoftwareInventory", update.GetSoftwareInventoryCollection)
	updateService.Get("/SoftwareInventory/{softwareInventory_id}", update.GetSoftwareInventory)
	return router
}
