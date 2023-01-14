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

// Package common ...
package common

// EventConst constant
type EventConst int

const (
	// RedfishEvent constant
	RedfishEvent EventConst = iota
	// MetricReport constant
	MetricReport
)

const (
	// RoleAdmin defines admin role for all service to authorize
	RoleAdmin = "Administrator"
	// RoleMonitor defines monitor role for all service to authorize
	RoleMonitor = "Operator"
	// RoleClient defines client role for all service to authorize
	RoleClient = "ReadOnly"

	// PrivilegeLogin defines the privilege for login
	PrivilegeLogin = "Login"
	// PrivilegeConfigureManager defines the privilege for configuraton manager
	PrivilegeConfigureManager = "ConfigureManager"
	// PrivilegeConfigureUsers defines the privilege for user configuratons
	PrivilegeConfigureUsers = "ConfigureUsers"
	// PrivilegeConfigureSelf defines the privilege for self configuratons
	PrivilegeConfigureSelf = "ConfigureSelf"
	// PrivilegeConfigureComponents defines the privilege for component configuratons
	PrivilegeConfigureComponents = "ConfigureComponents"

	// Below constans are for TaskState to Indicate the state of the task

	//Cancelled - This value shall represent that the operation was cancelled either through
	//a Delete on a Task Monitor or Task Resource or by an internal process.
	Cancelled = "Cancelled"
	//Cancelling - This value shall represent that the operation is in the process of being cancelled.
	Cancelling = "Cancelling"
	//Completed - This value shall represent that the operation is complete and
	//completed sucessfully or with warnings.
	Completed = "Completed"
	//Exception - This value shall represent that the operation is complete and completed with errors.
	Exception = "Exception"
	//Interrupted - This value shall represent that the operation has been interrupted but
	//is expected to restart and is therefore not complete.
	Interrupted = "Interrupted"
	//Killed - This value shall represent that the operation is complete because the task
	//was killed by an operator. Deprecated v1.2+. This value has been
	//deprecated and is being replaced by the value Cancelled which has more
	//determinate semantics.
	Killed = "Killed"
	//New - This value shall represent that this task is newly created but the operation
	//has not yet started.
	New = "New"
	//Pending - This value shall represent that the operation is pending some condition and
	//has not yet begun to execute.
	Pending = "Pending"
	//Running - This value shall represent that the operation is executing.
	Running = "Running"
	//Service - This value shall represent that the operation is now running as a service and
	//expected to continue operation until stopped or killed.
	Service = "Service"
	//Starting - This value shall represent that the operation is starting.
	Starting = "Starting"
	//Stopping - This value shall represent that the operation is stopping but is not yet complete.
	Stopping = "Stopping"
	//Suspended - This value shall represent that the operation has been suspended but is
	//expected to restart and is therefore not complete.
	Suspended = "Suspended"

	//Below constants are for TaskStatus to indicate the completion status
	//of the task

	//Critical - A critical condition exists that requires immediate attention.
	Critical = "Critical"
	//OK - Normal.
	OK = "OK"
	//Warning - A condition exists that requires attention.
	Warning = "Warning"

	// ManagerType constants (as per Manager_1.6.0 of DSP2046_2019.3)

	// ManagerTypeAuxiliaryController - A controller that provides management functions
	// for a particular subsystem or group of devices.
	ManagerTypeAuxiliaryController = "AuxiliaryController"
	// ManagerTypeBMC - A controller that provides management functions for a single computer system.
	ManagerTypeBMC = "BMC"
	// ManagerTypeEnclosureManager - A controller that provides management functions
	// for a chassis or group of devices or systems.
	ManagerTypeEnclosureManager = "EnclosureManager"
	// ManagerTypeManagementController - A controller that primarily monitors or manages the operation of a device or system.
	ManagerTypeManagementController = "ManagementController"
	// ManagerTypeRackManager - A controller that provides management functions for a whole or part of a rack.
	ManagerTypeRackManager = "RackManager"
	// ManagerTypeService - A software-based service that provides management functions.
	ManagerTypeService = "Service"

	// SubscriptionIndex is a index name which required for indexing of event subscriptions
	SubscriptionIndex = "Subscription"
	// DeviceSubscriptionIndex is a index name which required for indexing
	// subscription of device
	DeviceSubscriptionIndex = "DeviceSubscription"

	// ManagerAccountType has schema version to be returned with manager account
	ManagerAccountType = "#ManagerAccount.v1_9_0.ManagerAccount"
	// AccountServiceType has schema version to be returned with accountservice
	AccountServiceType = "#AccountService.v1_11_0.AccountService"
	// RoleType has schema version to be returned with Role
	RoleType = "#Role.v1_3_1.Role"
	// SessionServiceType has schema version to be returned with sessionservice
	SessionServiceType = "#SessionService.v1_1_8.SessionService"
	// SessionType has schema version to be returned with session
	SessionType = "#Session.v1_4_0.Session"
	// EventType has schema version to be returned with event
	EventType = "#Event.v1_7_0.Event"
	// AggregationServiceType has schema version to be returned with Aggregationservice
	AggregationServiceType = "#AggregationService.v1_0_1.AggregationService"
	// TaskType has schema version to be returned with Task
	TaskType = "#Task.v1_6_0.Task"
	// EventDestinationType has schema version to be returned with EventDestination
	EventDestinationType = "#EventDestination.v1_12_0.EventDestination"
	// EventServiceType has schema version to be returned with Event Service Type
	EventServiceType = "#EventService.v1_8_0.EventService"
	// ManagerType has schema version to be returned with Manager
	ManagerType = "#Manager.v1_15_0.Manager"
	// TaskEventType has schema version to be returned with TaskEvent
	TaskEventType = "TaskEvent.1.0.3"
	// UpdateServiceType has schema version to be returned with UpdateService
	UpdateServiceType = "#UpdateService.v1_11_0.UpdateService"
	// SettingsType has schema version to be returned with Settings in update service
	SettingsType = "#Settings.v1_3_3.OperationApplyTimeSupport"
	// TelemetryServiceType has version to be returned with Telemetry Service
	TelemetryServiceType = "#TelemetryService.v1_3_1.TelemetryService"
	//AggregationSourceType has version to be returned with AggregationSource Service
	AggregationSourceType = "#AggregationSource.v1_2_0.AggregationSource"
	//ChassisType has version to be returned with Chassis Service
	ChassisType = "#Chassis.v1_20_0.Chassis"
	// AggregateSubscriptionIndex is a index name which required for indexing
	// subscription of aggregate
	AggregateSubscriptionIndex = "AggregateToHost"
	// Below fields are Process Name for logging
	TransactionID = "transactionid"
	ThreadID      = "threadid"
	ThreadName    = "threadname"
	ActionName    = "actionname"
	ActionID      = "actionid"
	ProcessName   = "processname"
	// Below fields define Service Name
	ManagerService     = "svc-managers"
	AccountService     = "svc-account"
	SystemService      = "svc-systems"
	SessionService     = "svc-account-session"
	ApiService         = "svc-api"
	UpdateService      = "svc-update"
	TaskService        = "svc-task"
	AggregationService = "svc-aggregation"
	// DefaultThreadID to be used for apis
	DefaultThreadID = "0"
	// Invalid Action
	InvalidActionID   = "000"
	InvalidActionName = "MethodNotAllowed"
	// ThreadName
	CheckAuth                              = "Check-Authentication"
	CheckSessionCreation                   = "CheckSessionCreationCredentials"
	CheckSessionTimeout                    = "CheckSessionTimeOut"
	SendRequest                            = "SendRequest"
	StartRequest                           = "StartRequest"
	SimpleUpdate                           = "SimpleUpdate"
	StartUpdate                            = "StartUpdate"
	OverWriteCompletedTaskUtil             = "OverWriteCompletedTaskUtil"
	AsyncTaskDelete                        = "AsyncTaskDelete"
	ResetAggregates                        = "Reset-Aggregates"
	ResetAggregate                         = "Reset-Aggregate"
	SetBootOrder                           = "SettingBootOrder"
	CollectAndSetDefaultBootOrder          = "CollectAndSetDefaultBoorOrder"
	AddAggregationSource                   = "AddingAggregationSource"
	DeleteAggregationSource                = "DeleteAggregationSource"
	SubTaskStatusUpdate                    = "SubTaskStatusUpdate"
	ResetSystem                            = "ResetSystem"
	SetDefaultBootOrderElementsOfAggregate = "SetDefaultBootOrderElementsOfAggregate"
	RediscoverSystemInventory              = "RediscoverSystemInventory"
	CheckPluginStatus                      = "CheckPluginStatus"
	// constants for log
	SessionToken  = "sessiontoken"
	SessionUserID = "sessionuserid"
	SessionRoleID = "sessionroleid"
	StatusCode    = "statuscode"
)

// ActionType defines type of action
type ActionType struct {
	ActionID   string
	ActionName string
}

// ActionURL defines type of service, uri and its method
type ActionKey struct {
	Service string
	Uri     string
	Method  string
}

// Actions contain map of URI object and action related to request made by user
var Actions = map[ActionKey]ActionType{
	//Systems URI
	{"Systems", "Systems", "GET"}:                 {"001", "GetSystemsCollection"},
	{"Systems", "Systems/{id}", "GET"}:            {"002", "GetSystem"},
	{"Systems", "Processors", "GET"}:              {"003", "GetProcessorCollection"},
	{"Systems", "Processors/{id}", "GET"}:         {"004", "GetProcessor"},
	{"Systems", "Memory", "GET"}:                  {"005", "GetMemoryCollection"},
	{"Systems", "Memory/{id}", "GET"}:             {"006", "GetMemory"},
	{"Systems", "NetworkInterfaces", "GET"}:       {"007", "GetNetworkInterfacesCollection"},
	{"Systems", "NetworkInterfaces/{id}", "GET"}:  {"008", "GetNetworkInterfaces"},
	{"Systems", "MemoryDomains", "GET"}:           {"009", "GetMemoryDomainsCollection"},
	{"Systems", "EthernetInterfaces", "GET"}:      {"010", "GetEthernetInterfacesCollection"},
	{"Systems", "EthernetInterfaces/{id}", "GET"}: {"011", "GetEthernetInterfaces"},
	{"Systems", "VLANS", "GET"}:                   {"012", "GetVLANSCollection"},
	{"Systems", "VLANS/{id}", "GET"}:              {"013", "GetVLANS"},
	{"Systems", "SecureBoot", "GET"}:              {"014", "GetSecureBoot"},
	{"Systems", "BootOptions", "GET"}:             {"015", "GetBootOptionsCollection"},
	{"Systems", "BootOptions/{id}", "GET"}:        {"016", "GetBootOptions"},
	{"Systems", "LogServices", "GET"}:             {"017", "GetLogServicesCollection"},
	{"Systems", "LogServices/{id}", "GET"}:        {"018", "GetLogServices"},
	{"Systems", "Entries", "GET"}:                 {"019", "GetEntriesCollection"},
	{"Systems", "Entries/{id}", "GET"}:            {"020", "GetEntries"},
	{"Systems", "LogService.ClearLog", "POST"}:    {"021", "ClearLog"},
	{"Systems", "Systems/{id}", "PATCH"}:          {"022", "ChangeBootOrderSettings"},
	{"Systems", "PCIeDevices/{id}", "GET"}:        {"023", "GetPCIeDevices"},
	//Task URI
	{"TaskService", "TaskService", "GET"}:   {"024", "GetTaskService"},
	{"TaskService", "Tasks", "GET"}:         {"025", "TaskCollection"},
	{"TaskService", "Tasks/{id}", "GET"}:    {"026", "GetTaskStatus"},
	{"TaskService", "SubTasks", "GET"}:      {"027", "GetSubTaskCollection"},
	{"TaskService", "SubTasks/{id}", "GET"}: {"028", "GetSubTask"},
	{"TaskService", "Tasks/{id}", "DELETE"}: {"029", "DeleteTask"},
	// Roles URI
	{"Roles", "Roles", "GET"}:         {"030", "GetAllRoles"},
	{"Roles", "Roles/{id}", "GET"}:    {"031", "GetRole"},
	{"Roles", "Roles/{id}", "PATCH"}:  {"032", "UpdateRole"},
	{"Roles", "Roles/{id}", "DELETE"}: {"033", "DeleteRole"},
	// Account Service URI
	{"AccountService", "AccountService", "GET"}:   {"034", "GetAccountService"},
	{"AccountService", "Accounts", "GET"}:         {"035", "GetAllAccounts"},
	{"AccountService", "Accounts/{id}", "GET"}:    {"036", "GetAccount"},
	{"AccountService", "Accounts", "POST"}:        {"037", "CreateAccount"},
	{"AccountService", "Accounts/{id}", "PATCH"}:  {"038", "UpdateAccount"},
	{"AccountService", "Accounts/{id}", "DELETE"}: {"039", "DeleteAccount"},
	// Session Service URI
	{"SessionService", "SessionService", "GET"}:   {"040", "GetSessionService"},
	{"SessionService", "Sessions", "GET"}:         {"041", "GetAllActiveSessions"},
	{"SessionService", "Sessions/{id}", "GET"}:    {"042", "GetSession"},
	{"SessionService", "Sessions", "POST"}:        {"043", "CreateSession"},
	{"SessionService", "Sessions/{id}", "DELETE"}: {"044", "DeleteSession"},
	// Registry URI
	{"Registries", "Registries", "GET"}:      {"045", "GetRegistryFileCollection"},
	{"Registries", "Registries/{id}", "GET"}: {"046", "GetMessageRegistryFileID"},
	// Redfish URI
	{"", "redfish", "GET"}:   {"047", "GetVersion"},
	{"", "v1", "GET"}:        {"048", "GetServiceRoot"},
	{"", "odata", "GET"}:     {"049", "GetOdata"},
	{"", "$metadata", "GET"}: {"050", "GetMetadata"},
	// Taskmon URI
	{"taskmon", "taskmon/{id}", "GET"}: {"051", "GetTaskMonitor"},
	// BIOS URI
	{"Systems", "Bios", "GET"}:       {"052", "GetBiosDetails"},
	{"Systems", "Settings", "GET"}:   {"053", "GetBiosSettings"},
	{"Systems", "Settings", "PATCH"}: {"054", "ChangeBiosSettings"},
	// system Storage URI
	{"Systems", "Storage", "GET"}:               {"055", "GetStorageCollection"},
	{"Systems", "Storage/{id}", "GET"}:          {"056", "GetStorage"},
	{"Systems", "Drives/{id}", "GET"}:           {"057", "GetDrive"},
	{"Systems", "Controllers", "GET"}:           {"058", "GetControllersCollection"},
	{"Systems", "Controllers/{id}", "GET"}:      {"059", "GetSystemController"},
	{"Systems", "Ports", "GET"}:                 {"060", "GetPortsCollection"},
	{"Systems", "Ports/{id}", "GET"}:            {"061", "GetSystemPort"},
	{"Systems", "Volumes", "GET"}:               {"062", "GetVolumeCollection"},
	{"Systems", "Volumes", "POST"}:              {"063", "CreateVolume"},
	{"Systems", "Capabilities", "GET"}:          {"064", "GetVolumeCapabilities"},
	{"Systems", "Volumes/{id}", "DELETE"}:       {"065", "DeleteVolume"},
	{"Systems", "Volumes/{id}", "GET"}:          {"066", "GetVolume"},
	{"Systems", "StoragePools", "GET"}:          {"067", "GetStoragePoolCollection"},
	{"Systems", "StoragePools/{id}", "GET"}:     {"068", "GetStoragePools"},
	{"Systems", "AllocatedVolumes", "GET"}:      {"069", "GetAllocatedVolumesCollection"},
	{"Systems", "AllocatedVolumes/{id}", "GET"}: {"070", "GetAllocatedVolumes"},
	{"Systems", "ProvidingVolumes", "GET"}:      {"071", "GetProvidingVolumesCollection"},
	{"Systems", "ProvidingVolumes/{id}", "GET"}: {"072", "GetProvidingVolumes"},
	{"Systems", "ProvidingDrives", "GET"}:       {"073", "GetProvidingDrivesCollection"},
	{"Systems", "ProvidingDrives/{id}", "GET"}:  {"074", "GetProvidingDrives"},
	// Actions URI
	{"Systems", "ComputerSystem.Reset", "POST"}:               {"075", "ComputerSystemReset"},
	{"Systems", "ComputerSystem.SetDefaultBootOrder", "POST"}: {"076", "SetDefaultBootOrder"},
	// Aggregation URI
	{"AggregationService", "AggregationService", "GET"}:                      {"077", "GetAggregationService"},
	{"AggregationService", "ResetActionInfo", "GET"}:                         {"078", "GetResetActionInfoService"},
	{"AggregationService", "SetDefaultBootOrderActionInfo", "GET"}:           {"079", "GetSetDefaultBootOrderActionInfo"},
	{"AggregationService", "AggregationService.Reset", "POST"}:               {"080", "AggregationServiceReset"},
	{"AggregationService", "AggregationService.SetDefaultBootOrder", "POST"}: {"081", "SetDefaultBootOrder"},
	//AggregationSources URI
	{"AggregationService", "AggregationSources", "POST"}:        {"082", "AddAggregationSource"},
	{"AggregationService", "AggregationSources", "GET"}:         {"083", "GetAllAggregationSource"},
	{"AggregationService", "AggregationSources/{id}", "GET"}:    {"084", "GetAggregationSource"},
	{"AggregationService", "AggregationSources/{id}", "PATCH"}:  {"085", "UpdateAggregationSource"},
	{"AggregationService", "AggregationSources/{id}", "DELETE"}: {"086", "DeleteAggregationSource"},
	// ConnectionMethods URI
	{"AggregationService", "ConnectionMethods", "GET"}:      {"087", "GetAllConnectionMethods"},
	{"AggregationService", "ConnectionMethods/{id}", "GET"}: {"088", "GetConnectionMethod"},
	// Aggregate URI
	{"AggregationService", "Aggregates", "GET"}:                     {"089", "GetAggregateCollection"},
	{"AggregationService", "Aggregates", "POST"}:                    {"090", "CreateAggregate"},
	{"AggregationService", "Aggregates/{id}", "GET"}:                {"100", "GetAggregate"},
	{"AggregationService", "Aggregates/{id}", "DELETE"}:             {"101", "DeleteAggregate"},
	{"AggregationService", "Aggregate.AddElements", "POST"}:         {"102", "AddElementsToAggregate"},
	{"AggregationService", "Aggregate.RemoveElements", "POST"}:      {"103", "RemoveElementsFromAggregate"},
	{"AggregationService", "Aggregate.Reset", "POST"}:               {"104", "ResetAggregateElements"},
	{"AggregationService", "Aggregate.SetDefaultBootOrder", "POST"}: {"105", "SetDefaultBootOrderAggregateElements"},
	// Chassis URI
	{"Chassis", "Chassis", "GET"}:                     {"106", "GetChassisCollection"},
	{"Chassis", "Chassis", "POST"}:                    {"107", "CreateChassis"},
	{"Chassis", "Chassis/{id}", "GET"}:                {"108", "GetChassis"},
	{"Chassis", "Chassis/{id}", "PATCH"}:              {"109", "UpdateChassis"},
	{"Chassis", "Chassis/{id}", "DELETE"}:             {"110", "DeleteChassis"},
	{"Chassis", "NetworkAdapters", "GET"}:             {"111", "GetAllNetworkAdapters"},
	{"Chassis", "NetworkAdapters/{id}", "GET"}:        {"112", "GetNetworkAdapters"},
	{"Chassis", "NetworkDeviceFunctions", "GET"}:      {"113", "GetAllNetworkDeviceFunctions"},
	{"Chassis", "NetworkDeviceFunctions/{id}", "GET"}: {"114", "GetNetworkDeviceFunctions"},
	{"Chassis", "NetworkPorts", "GET"}:                {"115", "GetAllNetworkPorts"},
	{"Chassis", "NetworkPorts/{id}", "GET"}:           {"116", "GetNetworkPorts"},
	{"Chassis", "Assembly", "GET"}:                    {"117", "GetChassisAssembly"},
	{"Chassis", "PCIeSlots", "GET"}:                   {"118", "GetAllPCIeSlots"},
	{"Chassis", "PCIeSlots/{id}", "GET"}:              {"119", "GetPCIeSlots"},
	{"Chassis", "PCIeDevices", "GET"}:                 {"120", "GetAllPCIeDevices"},
	{"Chassis", "PCIeDevices/{id}", "GET"}:            {"121", "GetPCIeDevices"},
	{"Chassis", "PCIeFunctions", "GET"}:               {"122", "GetAllAPCIeFunctions"},
	{"Chassis", "PCIeFunctions/{id}", "GET"}:          {"123", "GetPCIeFunctions"},
	{"Chassis", "Sensors", "GET"}:                     {"124", "GetAllSensors"},
	{"Chassis", "Sensors/{id}", "GET"}:                {"125", "GetSensors"},
	{"Chassis", "LogServices", "GET"}:                 {"126", "GetAllLogServices"},
	{"Chassis", "LogServices/{id}", "GET"}:            {"127", "GetLogServices"},
	{"Chassis", "Entries", "GET"}:                     {"128", "GetAllEntries"},
	{"Chassis", "Entries/{id}", "GET"}:                {"129", "GetEntries"},
	{"Chassis", "Power", "GET"}:                       {"130", "GetChassisPower"},
	{"Chassis", "#PowerControl/{id}", "GET"}:          {"131", "GetPowerControl"},
	{"Chassis", "#PowerSupplies/{id}", "GET"}:         {"132", "GetPowerSupplies"},
	{"Chassis", "#Redundancy/{id}", "GET"}:            {"133", "GetRedundancy"},
	{"Chassis", "Thermal", "GET"}:                     {"134", "GetChassisThermal"},
	{"Chassis", "#Fans/{id}", "GET"}:                  {"135", "GetChassisFans"},
	{"Chassis", "#Temperatures/{id}", "GET"}:          {"136", "GetChassisTemperatures"},
	// EventService URI
	{"EventService", "EventService", "GET"}:                 {"137", "GetEventService"},
	{"EventService", "Subscriptions", "GET"}:                {"138", "GetEventSubscriptionsCollection"},
	{"EventService", "Subscriptions/{id}", "GET"}:           {"139", "GetEventSubscription"},
	{"EventService", "Subscriptions", "POST"}:               {"140", "CreateEventSubscription"},
	{"EventService", "EventService.SubmitTestEvent", "GET"}: {"141", "SubmitTestEvent"},
	{"EventService", "Subscriptions/{id}", "DELETE"}:        {"142", "DeleteEventSubscription"},
	// Fabrics URI
	{"Fabrics", "Fabrics", "GET"}:              {"143", "GetFabricCollection"},
	{"Fabrics", "Fabrics/{id}", "GET"}:         {"144", "GetFabric"},
	{"Fabrics", "Switches", "GET"}:             {"145", "GetFabricSwitchCollection"},
	{"Fabrics", "Switches/{id}", "GET"}:        {"146", "GetFabricSwitch"},
	{"Fabrics", "Ports", "GET"}:                {"147", "GetSwitchPortCollection"},
	{"Fabrics", "Ports/{id}", "GET"}:           {"148", "GetSwitchPort"},
	{"Fabrics", "Zones", "GET"}:                {"149", "GetFabricZoneCollection"},
	{"Fabrics", "Endpoints", "GET"}:            {"150", "GetFabricEndPointCollection"},
	{"Fabrics", "AddressPools", "GET"}:         {"151", "GetFabricAddressPoolCollection"},
	{"Fabrics", "AddressPools/{id}", "GET"}:    {"152", "GetFabricAddressPool"},
	{"Fabrics", "AddressPools/{id}", "PUT"}:    {"153", "VerifyAndUpdateFabricAddressPool"},
	{"Fabrics", "AddressPools", "POST"}:        {"154", "CreateFabricAddressPool"},
	{"Fabrics", "AddressPools/{id}", "PATCH"}:  {"155", "UpdateFabricAddressPool"},
	{"Fabrics", "AddressPools/{id}", "DELETE"}: {"156", "DeleteFabricAddressPool"},
	{"Fabrics", "Zones/{id}", "GET"}:           {"157", "GetFabricZone"},
	{"Fabrics", "Zones/{id}", "PUT"}:           {"158", "VerifyAndUpdateFabricZone"},
	{"Fabrics", "Zones", "POST"}:               {"159", "CreateFabricZone"},
	{"Fabrics", "Zones/{id}", "PATCH"}:         {"160", "UpdateFabricZone"},
	{"Fabrics", "Zones/{id}", "DELETE"}:        {"161", "DeleteFabricZone"},
	{"Fabrics", "Endpoints/{id}", "GET"}:       {"162", "GetFabricEndPoints"},
	{"Fabrics", "Endpoints/{id}", "PUT"}:       {"163", "VerifyAndUpdateFabricEndPoints"},
	{"Fabrics", "Endpoints/{id}", "POST"}:      {"164", "CreateFabricEndPoints"},
	{"Fabrics", "Endpoints/{id}", "PATCH"}:     {"165", "UpdateFabricEndPoints"},
	{"Fabrics", "Endpoints/{id}", "DELETE"}:    {"166", "DeleteFabricEndPoints"},
	{"Fabrics", "Ports/{id}", "Patch"}:         {"167", "UpdateFabricPorts"},
	// Managers URI
	{"Managers", "Managers", "GET"}:                  {"168", "GetManagersCollection"},
	{"Managers", "Managers/{id}", "GET"}:             {"169", "GetManager"},
	{"Managers", "EthernetInterfaces", "GET"}:        {"170", "GetAllEthernetInterfaces"},
	{"Managers", "EthernetInterfaces/{id}", "GET"}:   {"171", "GetEthernetInterfaces"},
	{"Managers", "NetworkProtocol", "GET"}:           {"172", "GetAllNetworkProtocol"},
	{"Managers", "NetworkProtocol/{id}", "GET"}:      {"173", "GetNetworkProtocol"},
	{"Managers", "HostInterfaces", "GET"}:            {"174", "GetAllHostInterfaces"},
	{"Managers", "HostInterfaces/{id}", "GET"}:       {"175", "GetHostInterfaces"},
	{"Managers", "SerialInterfaces", "GET"}:          {"176", "GetAllSerialInterfaces"},
	{"Managers", "SerialInterfaces/{id}", "GET"}:     {"177", "GetSerialInterfaces"},
	{"Managers", "VirtualMedia", "GET"}:              {"178", "GetAllVirtualMedia"},
	{"Managers", "VirtualMedia/{id}", "GET"}:         {"179", "GetVirtualMedia"},
	{"Managers", "VirtualMedia.EjectMedia", "POST"}:  {"180", "VirtualMediaEjectMedia"},
	{"Managers", "VirtualMedia.InsertMedia", "POST"}: {"181", "VirtualMediaInsertMedia"},
	{"Managers", "LogServices", "GET"}:               {"182", "GetAllLogServices"},
	{"Managers", "LogServices/{id}", "GET"}:          {"183", "GetLogServices"},
	{"Managers", "Entries", "GET"}:                   {"184", "GetAllLogServiceEntries"},
	{"Managers", "Entries/{id}", "GET"}:              {"185", "GetLogServiceEntries"},
	{"Managers", "LogService.ClearLog", "POST"}:      {"186", "LogServiceClearLog"},
	{"Managers", "RemoteAccountService", "GET"}:      {"187", "GetRemoteAccountService"},
	{"Managers", "Accounts", "GET"}:                  {"188", "GetAllManagerAccounts"},
	{"Managers", "Accounts/{id}", "GET"}:             {"189", "GetManagerAccounts"},
	{"Managers", "Accounts", "POST"}:                 {"190", "CreateManagerAccounts"},
	{"Managers", "Accounts/{id}", "PATCH"}:           {"191", "UpdateManagerAccounts"},
	{"Managers", "Accounts/{id}", "DELETE"}:          {"192", "DeleteManagerAccounts"},
	{"Managers", "Roles", "GET"}:                     {"193", "GetAllManagersRoles"},
	{"Managers", "Roles/{id}", "GET"}:                {"194", "GetManagersRoles"},
	// Update Service URI
	{"UpdateService", "UpdateService", "GET"}:               {"195", "GetUpdateService"},
	{"UpdateService", "UpdateService.SimpleUpdate", "POST"}: {"196", "UpdateServiceSimpleUpdate"},
	{"UpdateService", "UpdateService.StartUpdate", "POST"}:  {"197", "UpdateServiceStartUpdate"},
	{"UpdateService", "FirmwareInventory", "GET"}:           {"198", "GetFirmwareInventoryCollection"},
	{"UpdateService", "FirmwareInventory/{id}", "GET"}:      {"199", "GetFirmwareInventory"},
	{"UpdateService", "SoftwareInventory", "GET"}:           {"200", "GetSoftwareInventoryCollection"},
	{"UpdateService", "SoftwareInventory/{id}", "GET"}:      {"201", "GetSoftwareInventory"},
	// Telemetry Service URI
	{"TelemetryService", "TelemetryService", "GET"}:             {"202", "GetTelemetryService"},
	{"TelemetryService", "MetricDefinitions", "GET"}:            {"203", "GetMetricDefinitionCollection"},
	{"TelemetryService", "MetricReportDefinitions", "GET"}:      {"204", "GetMetricReportDefinitionCollection"},
	{"TelemetryService", "MetricReports", "GET"}:                {"205", "GetMetricReportCollection"},
	{"TelemetryService", "Triggers", "GET"}:                     {"206", "GetTriggerCollection"},
	{"TelemetryService", "MetricDefinitions/{id}", "GET"}:       {"207", "GetMetricDefinition"},
	{"TelemetryService", "MetricReportDefinitions/{id}", "GET"}: {"208", "GetMetricReportDefinition"},
	{"TelemetryService", "MetricReports/{id}", "GET"}:           {"209", "GetMetricReport"},
	{"TelemetryService", "Triggers/{id}", "GET"}:                {"210", "GetTrigger"},
	{"TelemetryService", "Triggers/{id}", "PATCH"}:              {"211", "UpdateTrigger"},
	//License Service URI
	{"LicenseService", "LicenseService", "GET"}: {"212", "GetLicenseService"},
	{"LicenseService", "Licenses", "GET"}:       {"213", "GetLicenseCollection"},
	{"LicenseService", "Licenses/{id}", "GET"}:  {"214", "GetLicenseResource"},
	{"LicenseService", "Licenses", "POST"}:      {"215", "InstallLicenseService"},
	// 216 and 217 operations are svc-aggregation internal operations pluginhealthcheck and RediscoverSystem
}

var Types = map[string]string{
	"EthernetInterfaces": "#EthernetInterface.v1_8_0.EthernetInterface",
}

// RediscoverResources contains to get only these resource from the device when
// reset flag is set when device is restarted.
var RediscoverResources = []string{
	"Bios",
	"BootOptions",
	"Storage",
}

// SystemResource contains the Resource name and table name
// this map is basically to fetch the table name against the system resource name,
// so it will be usefull to store the resource data into the particular database table
// and also it will be usefull to retrives the system resource data
var SystemResource = map[string]string{
	"Bios":               "Bios",
	"SecureBoot":         "SecureBoot",
	"Storage":            "StorageCollection",
	"BootOptions":        "BootOptionsCollection",
	"MemoryDomains":      "MemoryDomainsCollection",
	"NetworkInterfaces":  "NetworkInterfacesCollection",
	"Processors":         "ProcessorsCollection",
	"EthernetInterfaces": "EthernetInterfacesCollection",
	"Memory":             "MemoryCollection",
	"VLANS":              "VLANS",
	"LogServices":        "LogServicesCollection",
	"Settings":           "Bios",
	"Volumes":            "VolumesCollection",
	"Drives":             "DrivesCollection",
}

// ChassisResource contains the Resource name and table name
// this map is basically to fetch the table name against the chassis resource name,
// so it will be usefull to store the resource data into the particular database table
// and also it will be usefull to retrives the chassis resource data
var ChassisResource = map[string]string{
	"Power":                  "Power",
	"Thermal":                "Thermal",
	"NetworkAdapters":        "NetworkAdaptersCollection",
	"NetworkPorts":           "NetworkPortsCollection",
	"NetworkDeviceFunctions": "NetworkDeviceFunctionsCollection",
	"Assembly":               "Assembly",
	"PCIeSlots":              "PCIeSlots",
	"PCIeDevices":            "PCIeDevicesCollection",
	"Sensors":                "SensorsCollection",
	"LogServices":            "LogServicesCollection",
}

// ManagersResource contains the Resource name and table name
// this map is basically to fetch the table name against the manager resource name,
// so it will be usefull to store the resource data into the particular database table
// and also it will be usefull to retrives the manager resource data
var ManagersResource = map[string]string{
	"NetworkProtocol":    "NetworkProtocol",
	"EthernetInterfaces": "EthernetInterfacesCollection",
	"HostInterfaces":     "HostInterfacesCollection",
	"VirtualMedia":       "VirtualMediaCollection",
	"LogServices":        "LogServicesCollection",
	"SerialInterfaces":   "SerialInterfaceCollection",
	"Entries":            "EntriesCollection",
}

// ResourceTypes specifies the map  of valid resource types that can be used for an event subscription
var ResourceTypes = map[string]string{
	"AccelerationFunction":   "AccelerationFunction",
	"AddressPool":            "AddressPool",
	"Assembly":               "Assembly",
	"Bios":                   "Bios",
	"BootOption":             "BootOptions",
	"Chassis":                "Chassis",
	"ComputerSystem":         "Systems",
	"Drive":                  "Drive",
	"Endpoint":               "Endpoint",
	"EthernetInterface":      "EthernetInterfaces",
	"Event":                  "Event",
	"EventDestination":       "EventDestination",
	"EventService":           "EventService",
	"Fabric":                 "Fabric",
	"HostInterface":          "HostInterfaces",
	"IPAddresses":            "IPAddresses",
	"Job":                    "Job",
	"JobService":             "JobService",
	"LogEntry":               "LogEntry",
	"LogService":             "LogServices",
	"Manager":                "Manager",
	"ManagerAccount":         "ManagerAccount",
	"ManagerNetworkProtocol": "ManagerNetworkProtocol",
	"Memory":                 "Memory",
	"MemoryChunks":           "MemoryChunks",
	"MemoryDomain":           "MemoryDomains",
	"MemoryMetrics":          "MemoryMetrics",
	"Message":                "Message",
	"MessageRegistry":        "MessageRegistry",
	"MessageRegistryFile":    "MessageRegistryFile",
	"NetworkAdapter":         "NetworkAdapters",
	"NetworkDeviceFunction":  "NetworkDeviceFunction",
	"NetworkInterface":       "NetworkInterfaces",
	"NetworkPort":            "NetworkPort",
	"PCIeDevice":             "PCIeDevices",
	"PCIeFunction":           "PCIeFunction",
	"PCIeSlots":              "PCIeSlots",
	"PhysicalContext":        "PhysicalContext",
	"Port":                   "Port",
	"Power":                  "Power",
	"PrivilegeRegistry":      "PrivilegeRegistry",
	"Privileges":             "Privileges",
	"Processor":              "Processors",
	"ProcessorCollection":    "ProcessorCollection",
	"ProcessorMetrics":       "ProcessorMetrics",
	"Protocol":               "Protocol",
	"Redundancy":             "Redundancy",
	"Resource":               "Resource",
	"Role":                   "Role",
	"SecureBoot":             "SecureBoot",
	"Sensor":                 "Sensor",
	"SerialInterfaces":       "SerialInterfaces",
	"Session":                "Session",
	"Storage":                "Storage",
	"Switch":                 "Switch",
	"Task":                   "Task",
	"Thermal":                "Thermal",
	"VLanNetworkInterface":   "VLanNetworkInterface",
	"Volume":                 "Volume",
	"Zone":                   "Zone",
}

// Events contains the data with IP sent fro mplugin to PMB
type Events struct {
	IP        string `json:"ip"`
	Request   []byte `json:"request"`
	EventType string `json:"eventType"`
}

// MessageData contains information of Events and message details including arguments
// it will be used to pass to gob encoding/decoding which will register the type.
// it will be send as byte stream on the wire to/from kafka
type MessageData struct {
	OdataType string  `json:"@odata.type"`
	Name      string  `json:"Name"`
	Context   string  `json:"@odata.context"`
	Events    []Event `json:"Events"`
}

// Event contains the details of the event subscribed from PMB
type Event struct {
	MemberID          string      `json:"MemberId,omitempty"`
	EventType         string      `json:"EventType"`
	EventGroupID      int         `json:"EventGroupId,omitempty"`
	EventID           string      `json:"EventId"`
	Severity          string      `json:"Severity"`
	EventTimestamp    string      `json:"EventTimestamp"`
	Message           string      `json:"Message"`
	MessageArgs       []string    `json:"MessageArgs,omitempty"`
	MessageID         string      `json:"MessageId"`
	Oem               interface{} `json:"Oem,omitempty"`
	OriginOfCondition *Link       `json:"OriginOfCondition,omitempty"`
}

// Link  property shall contain a link to the resource or object that originated the condition that caused the event to be generated
type Link struct {
	Oid string `json:"@odata.id"`
}

// DeviceSubscription is a model to store the subscription details of a device
type DeviceSubscription struct {
	EventHostIP     string   `json:"EventHostIP,omitempty"`
	OriginResources []string `json:"OriginResources"`
	Location        string   `json:"location,omitempty"`
}

// URIWithNoAuth contains the list of URI's which does not require authentication
var URIWithNoAuth = []string{
	"/redfish/v1",
	"/redfish/v1/$metadata",
	"/redfish/v1/odata",
	"/redfish/v1/SessionService",
	"/redfish/v1/SessionService/Sessions",
}

var SessionURI = "/redfish/v1/SessionService/Sessions"
