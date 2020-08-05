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

//Package router ...
package router

import (
	srv "github.com/ODIM-Project/ODIM/lib-utilities/services"
	"github.com/ODIM-Project/ODIM/svc-api/handle"
	"github.com/ODIM-Project/ODIM/svc-api/middleware"
	"github.com/ODIM-Project/ODIM/svc-api/rpc"
	iris "github.com/kataras/iris/v12"
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
		GetAggregationServiceRPC:   rpc.DoGetAggregationService,
		AddComputeRPC:              rpc.DoAddComputeRequest,
		DeleteComputeRPC:           rpc.DoDeleteComputeRequest,
		ResetRPC:                   rpc.DoResetRequest,
		SetDefaultBootOrderRPC:     rpc.DoSetDefaultBootOrderRequest,
		AddAggregationSourceRPC:    rpc.DoAddAggregationSource,
		GetAllAggregationSourceRPC: rpc.DoGetAllAggregationSource,
		GetAggregationSourceRPC:    rpc.DoGetAggregationSource,
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
	}

	cha := handle.ChassisRPCs{
		GetChassisCollectionRPC: rpc.GetChassisCollection,
		GetChassisResourceRPC:   rpc.GetChassisResource,
		GetChassisRPC:           rpc.GetChassis,
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
		GetUpdateServiceRPC: rpc.DoGetUpdateService,
	}

	registryFile := handle.Registry{
		Auth: srv.IsAuthorized,
	}

	serviceRoot := handle.InitServiceRoot()

	router := iris.New()

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
	v1.Get("/registries/{id}", registryFile.GetMessageRegistryFile)
	v1.Any("/registries", handle.RegMethodNotAllowed)
	v1.Any("/registries/{id}", handle.RegMethodNotAllowed)

	registry := v1.Party("/Registries")
	registry.SetRegisterRule(iris.RouteSkip)
	registry.Get("/", registryFile.GetRegistryFileCollection)
	registry.Get("/{id}", registryFile.GetMessageRegistryFileID)
	registry.Any("/", handle.RegMethodNotAllowed)
	registry.Any("/{id}", handle.RegMethodNotAllowed)

	session := v1.Party("/SessionService")
	session.SetRegisterRule(iris.RouteSkip)
	session.Get("/", middleware.SessionDelMiddleware, s.GetSessionService)
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
	systems.Get("/{id}/Storage", system.GetSystemResource)
	systems.Get("/{id}/Storage/{rid}", system.GetSystemResource)
	systems.Get("/{id}/Storage/{rid}/Drives", system.GetSystemResource)
	systems.Get("/{id}/Storage/{rid}/Drives/{rid2}", system.GetSystemResource)
	systems.Patch("/{id}", system.ChangeBootOrderSettings)
	systems.Any("/", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/EthernetInterfaces", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/EthernetInterfaces/{rid}", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/SecureBoot", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/MemoryDomains", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/NetworkInterfaces", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/Memory", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/Processors", handle.SystemsMethodNotAllowed)
	systems.Any("/{id}/Storage", handle.SystemsMethodNotAllowed)
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

	systemsAction := systems.Party("/{id}/Actions", middleware.SessionDelMiddleware)
	systemsAction.SetRegisterRule(iris.RouteSkip)
	systemsAction.Post("/ComputerSystem.Reset", system.ComputerSystemReset)
	systemsAction.Post("/ComputerSystem.SetDefaultBootOrder", system.SetDefaultBootOrder)

	aggregation := v1.Party("/AggregationService", middleware.SessionDelMiddleware)
	aggregation.SetRegisterRule(iris.RouteSkip)
	aggregation.Get("/", pc.GetAggregationService)
	aggregation.Post("/Actions/AggregationService.Add/", pc.AddCompute)
	aggregation.Post("/Actions/AggregationService.Remove/", pc.DeleteCompute)
	aggregation.Post("/Actions/AggregationService.Reset/", pc.Reset)
	aggregation.Post("/Actions/AggregationService.SetDefaultBootOrder/", pc.SetDefaultBootOrder)
	aggregation.Any("/", handle.AggMethodNotAllowed)
	aggregation.Any("/Actions/AggregationService.Add/", handle.AggMethodNotAllowed)
	aggregationSource := aggregation.Party("/AggregationSource", middleware.SessionDelMiddleware)
	aggregationSource.Post("/", pc.AddAggregationSource)
	aggregationSource.Get("/", pc.GetAllAggregationSource)
	aggregationSource.Get("/{id}", pc.GetAggregationSource)

	chassis := v1.Party("/Chassis", middleware.SessionDelMiddleware)
	chassis.SetRegisterRule(iris.RouteSkip)
	chassis.Get("/", cha.GetChassisCollection)
	chassis.Get("/{id}", cha.GetChassis)
	chassis.Get("/{id}/NetworkAdapters", cha.GetChassisResource)
	chassis.Any("/", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}", handle.ChassisMethodNotAllowed)
	chassis.Any("/{id}/NetworkAdapters", handle.ChassisMethodNotAllowed)

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

	fabrics := v1.Party("/Fabrics", middleware.SessionDelMiddleware)
	fabrics.SetRegisterRule(iris.RouteSkip)
	fabrics.Get("/", fab.GetFabricResource)
	fabrics.Get("/{id}", fab.GetFabricResource)
	fabrics.Get("/{id}/Switches", fab.GetFabricResource)
	fabrics.Get("/{id}/Switches/{switchID}", fab.GetFabricResource)
	fabrics.Get("/{id}/Switches/{switchID}/Ports", fab.GetFabricResource)
	fabrics.Get("/{id}/Switches/{switchID}/Ports/{port_uuid}", fab.GetFabricResource)
	fabrics.Get("/{id}/Zones/", fab.GetFabricResource)
	fabrics.Get("/{id}/Endpoints/", fab.GetFabricResource)
	fabrics.Get("/{id}/AddressPools/", fab.GetFabricResource)
	fabrics.Get("/{id}/Zones/{zone_uuid}", fab.GetFabricResource)
	fabrics.Get("/{id}/Endpoints/{endpoint_uuid}", fab.GetFabricResource)
	fabrics.Get("/{id}/AddressPools/{addresspool_uuid}", fab.GetFabricResource)
	fabrics.Put("/{id}/Zones/{zone_uuid}", fab.UpdateFabricResource)
	fabrics.Put("/{id}/Endpoints/{endpoint_uuid}", fab.UpdateFabricResource)
	fabrics.Put("/{id}/AddressPools/{addresspool_uuid}", fab.UpdateFabricResource)
	fabrics.Post("/{id}/Zones", fab.UpdateFabricResource)
	fabrics.Post("/{id}/Endpoints", fab.UpdateFabricResource)
	fabrics.Post("/{id}/AddressPools", fab.UpdateFabricResource)
	fabrics.Patch("/{id}/Zones/{zone_uuid}", fab.UpdateFabricResource)
	fabrics.Patch("/{id}/Endpoints/{endpoint_uuid}", fab.UpdateFabricResource)
	fabrics.Patch("/{id}/AddressPools/{addresspool_uuid}", fab.UpdateFabricResource)
	fabrics.Delete("/{id}/Zones/{zone_uuid}", fab.DeleteFabricResource)
	fabrics.Delete("/{id}/Endpoints/{endpoint_uuid}", fab.DeleteFabricResource)
	fabrics.Delete("/{id}/AddressPools/{addresspool_uuid}", fab.DeleteFabricResource)
	fabrics.Any("/", handle.FabricsMethodNotAllowed)

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
	return router
}
