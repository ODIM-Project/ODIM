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
)

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
	"SerialInterface":        "SerialInterface",
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
	IP      string `json:"ip"`
	Request []byte `json:"request"`
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
