/*
 * Copyright (c) Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.odim.simulator.tree.structure

import com.odim.simulator.tree.structure.ResourceType.ACCELERATION_FUNCTION
import com.odim.simulator.tree.structure.ResourceType.ADDRESS_POOL
import com.odim.simulator.tree.structure.ResourceType.AGGREGATE
import com.odim.simulator.tree.structure.ResourceType.AGGREGATION_SOURCE
import com.odim.simulator.tree.structure.ResourceType.BOOT_OPTION
import com.odim.simulator.tree.structure.ResourceType.CERTIFICATE
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.CLASS_OF_SERVICE
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import com.odim.simulator.tree.structure.ResourceType.CONNECTION
import com.odim.simulator.tree.structure.ResourceType.CONNECTION_METHOD
import com.odim.simulator.tree.structure.ResourceType.CONSISTENCY_GROUP
import com.odim.simulator.tree.structure.ResourceType.DRIVE
import com.odim.simulator.tree.structure.ResourceType.ENDPOINT
import com.odim.simulator.tree.structure.ResourceType.ENDPOINT_GROUP
import com.odim.simulator.tree.structure.ResourceType.ETHERNET_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.EVENT_DESTINATION
import com.odim.simulator.tree.structure.ResourceType.EXTERNAL_ACCOUNT_PROVIDER
import com.odim.simulator.tree.structure.ResourceType.FABRIC
import com.odim.simulator.tree.structure.ResourceType.FACILITY
import com.odim.simulator.tree.structure.ResourceType.FILE_SHARE
import com.odim.simulator.tree.structure.ResourceType.FILE_SYSTEM
import com.odim.simulator.tree.structure.ResourceType.HOST_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.JOB
import com.odim.simulator.tree.structure.ResourceType.JSON_SCHEMA_FILE
import com.odim.simulator.tree.structure.ResourceType.LINE_OF_SERVICE
import com.odim.simulator.tree.structure.ResourceType.LOG_ENTRY
import com.odim.simulator.tree.structure.ResourceType.LOG_SERVICE
import com.odim.simulator.tree.structure.ResourceType.MANAGER
import com.odim.simulator.tree.structure.ResourceType.MANAGER_ACCOUNT
import com.odim.simulator.tree.structure.ResourceType.MEDIA_CONTROLLER
import com.odim.simulator.tree.structure.ResourceType.MEMORY
import com.odim.simulator.tree.structure.ResourceType.MEMORY_CHUNKS
import com.odim.simulator.tree.structure.ResourceType.MEMORY_DOMAIN
import com.odim.simulator.tree.structure.ResourceType.MESSAGE_REGISTRY_FILE
import com.odim.simulator.tree.structure.ResourceType.METRIC_DEFINITION
import com.odim.simulator.tree.structure.ResourceType.METRIC_REPORT
import com.odim.simulator.tree.structure.ResourceType.METRIC_REPORT_DEFINITION
import com.odim.simulator.tree.structure.ResourceType.NETWORK_ADAPTER
import com.odim.simulator.tree.structure.ResourceType.NETWORK_DEVICE_FUNCTION
import com.odim.simulator.tree.structure.ResourceType.NETWORK_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.NETWORK_PORT
import com.odim.simulator.tree.structure.ResourceType.OPERATING_CONFIG
import com.odim.simulator.tree.structure.ResourceType.OUTLET
import com.odim.simulator.tree.structure.ResourceType.PCIE_DEVICE
import com.odim.simulator.tree.structure.ResourceType.PCIE_FUNCTION
import com.odim.simulator.tree.structure.ResourceType.PORT
import com.odim.simulator.tree.structure.ResourceType.POWER_DISTRIBUTION
import com.odim.simulator.tree.structure.ResourceType.POWER_DOMAIN
import com.odim.simulator.tree.structure.ResourceType.PROCESSOR
import com.odim.simulator.tree.structure.ResourceType.RESOURCE_BLOCK
import com.odim.simulator.tree.structure.ResourceType.ROLE
import com.odim.simulator.tree.structure.ResourceType.ROUTE_ENTRY
import com.odim.simulator.tree.structure.ResourceType.ROUTE_SET_ENTRY
import com.odim.simulator.tree.structure.ResourceType.SECURE_BOOT_DATABASE
import com.odim.simulator.tree.structure.ResourceType.SENSOR
import com.odim.simulator.tree.structure.ResourceType.SERIAL_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.SESSION
import com.odim.simulator.tree.structure.ResourceType.SIGNATURE
import com.odim.simulator.tree.structure.ResourceType.SIMPLE_STORAGE
import com.odim.simulator.tree.structure.ResourceType.SOFTWARE_INVENTORY
import com.odim.simulator.tree.structure.ResourceType.STORAGE
import com.odim.simulator.tree.structure.ResourceType.STORAGE_CONTROLLER
import com.odim.simulator.tree.structure.ResourceType.STORAGE_GROUP
import com.odim.simulator.tree.structure.ResourceType.STORAGE_POOL
import com.odim.simulator.tree.structure.ResourceType.STORAGE_SERVICE
import com.odim.simulator.tree.structure.ResourceType.SWITCH
import com.odim.simulator.tree.structure.ResourceType.TASK
import com.odim.simulator.tree.structure.ResourceType.TRIGGERS
import com.odim.simulator.tree.structure.ResourceType.VCAT_ENTRY
import com.odim.simulator.tree.structure.ResourceType.VIRTUAL_MEDIA
import com.odim.simulator.tree.structure.ResourceType.VLAN_NETWORK_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.VOLUME
import com.odim.simulator.tree.structure.ResourceType.ZONE
import com.odim.simulator.tree.structure.TypeKind.EMBEDDED
import com.odim.simulator.tree.structure.TypeKind.SIMPLE
import com.odim.simulator.tree.structure.TypeKind.SINGLETON

typealias ResourceObject = MutableMap<String, Any?>
typealias StringArray = MutableList<String>

fun ResourceObject.removeProperties(vararg paths: String) = paths.forEach(::removeProperty)

@Suppress("UNCHECKED_CAST")
fun ResourceObject.removeProperty(path: String) {
    val pathChunked = path.split("/")

    var map = this
    for (part in pathChunked) {
        when (map[part]) {
            is Map<*, *> -> map = map[part] as ResourceObject
            else -> {
                map.remove(part)
                return
            }
        }
    }
    map.remove(pathChunked.last())
}

fun ResourceObject.removeAction(actionType: ActionType) {
    (get("Actions") as? Actions)?.removeAction(actionType)
}

fun ResourceObject.addToProperty(path: String, additive: ResourceObject) {
    val pathChunked = path.split("/")

    var map = this
    for (part in pathChunked) {
        if (map[part] is Map<*, *>) {
            map = map[part] as ResourceObject
        }
    }
    map.putAll(additive)
}

interface Type {
    fun oDataType(): String
    fun jsonName(): String
    fun isEmbedded(): Boolean
    fun isSingleton(): Boolean
}

interface ResourceCollectionTypeBase : Type {
    fun of(): ResourceTypeBase
    override fun isEmbedded() = false
    override fun isSingleton() = false
}

interface ResourceTypeBase : Type {
    fun actionsNamespace(): String
}

enum class TypeKind {
    SIMPLE,
    SINGLETON,
    EMBEDDED
}

interface ExtensibleEmbeddedObjectType

enum class ResourceCollectionType(
        private val odataTypeValue: String,
        private val of: ResourceType,
        private val jsonName: String
) : ResourceCollectionTypeBase {
    COMPUTER_SYSTEMS_COLLECTION("ComputerSystemCollection", COMPUTER_SYSTEM, "Systems"),
    MANAGERS_COLLECTION("ManagerCollection", MANAGER, "Managers"),
    PROCESSORS_COLLECTION("ProcessorCollection", PROCESSOR, "Processors"),
    MEMORY_COLLECTION("MemoryCollection", MEMORY, "Memory"),
    ETHERNET_INTERFACES_COLLECTION("EthernetInterfaceCollection", ETHERNET_INTERFACE, "EthernetInterfaces"),
    CHASSIS_COLLECTION("ChassisCollection", CHASSIS, "Chassis"),
    PCIE_DEVICES_COLLECTION("PCIeDevicesCollection", PCIE_DEVICE, "PCIeDevices"),
    STORAGE_SERVICES_COLLECTION("StorageServiceCollection", STORAGE_SERVICE, "StorageServices"),
    RESOURCE_BLOCKS_COLLECTION("ResourceBlockCollection", RESOURCE_BLOCK, "ResourceBlocks"),
    SWITCHES_COLLECTION("SwitchCollection", SWITCH, "Switches"),
    PORTS_COLLECTION("PortCollection", PORT, "Ports"),
    NETWORK_ADAPTERS_COLLECTION("NetworkAdapterCollection", NETWORK_ADAPTER, "NetworkAdapters"),
    NETWORK_DEVICE_FUNCTIONS_COLLECTION("NetworkDeviceFunctionCollection", NETWORK_DEVICE_FUNCTION, "NetworkDeviceFunctions"),
    PCIE_FUNCTIONS_COLLECTION("PcieFunctionCollection", PCIE_FUNCTION, "PCIeFunctions"),
    DRIVES_COLLECTION("DriveCollection", DRIVE, "Drives"),
    SESSIONS_COLLECTION("SessionCollection", SESSION, "Sessions"),
    MANAGER_ACCOUNTS_COLLECTION("ManagerAccountCollection", MANAGER_ACCOUNT, "Accounts"),
    ROLES_COLLECTION("RoleCollection", ROLE, "Roles"),
    EXTERNAL_ACCOUNT_PROVIDERS_COLLECTION("ExternalAccountProviderCollection", EXTERNAL_ACCOUNT_PROVIDER,
            "AdditionalExternalAccountProviders"),
    MESSAGE_REGISTRY_FILES_COLLECTION("MessageRegistryFileCollection", MESSAGE_REGISTRY_FILE, "Registries"),
    JSON_SCHEMA_FILES_COLLECTION("JsonSchemaFileCollection", JSON_SCHEMA_FILE, "JsonSchemas"),
    EVENT_DESTINATIONS_COLLECTION("EventDestinationCollection", EVENT_DESTINATION, "Subscriptions"),
    FABRICS_COLLECTION("FabricCollection", FABRIC, "Fabrics"),
    ZONES_COLLECTION("ZoneCollection", ZONE, "Zones"),
    ENDPOINTS_COLLECTION("EndpointCollection", ENDPOINT, "Endpoints"),
    SOFTWARE_INVENTORIES_COLLECTION("SoftwareInventoryCollection", SOFTWARE_INVENTORY, "SoftwareInventory"),
    SERIAL_INTERFACES_COLLECTION("SerialInterfaceCollection", SERIAL_INTERFACE, "SerialInterfaces"),
    LOG_SERVICES_COLLECTION("LogServiceCollection", LOG_SERVICE, "LogServices"),
    LOG_ENTRIES_COLLECTION("LogEntryCollection", LOG_ENTRY, "Entries"),
    VIRTUAL_MEDIAS_COLLECTION("VirtualMediaCollection", VIRTUAL_MEDIA, "VirtualMedia"),
    HOST_INTERFACES_COLLECTION("HostInterfaceCollection", HOST_INTERFACE, "HostInterfaces"),
    BOOT_OPTIONS_COLLECTION("BootOptionCollection", BOOT_OPTION, "BootOptions"),
    SIMPLE_STORAGES_COLLECTION("SimpleStorageCollection", SIMPLE_STORAGE, "SimpleStorage"),
    STORAGES_COLLECTION("StorageCollection", STORAGE, "Storage"),
    VOLUMES_COLLECTION("VolumeCollection", VOLUME, "Volumes"),
    MEMORY_DOMAINS_COLLECTION("MemoryDomainCollection", MEMORY_DOMAIN, "MemoryDomains"),
    MEMORY_CHUNKS_COLLECTION("MemoryChunksCollection", MEMORY_CHUNKS, "MemoryChunks"),
    NETWORK_INTERFACES_COLLECTION("NetworkInterfaceCollection", NETWORK_INTERFACE, "NetworkInterfaces"),
    NETWORK_PORTS_COLLECTION("NetworkPortCollection", NETWORK_PORT, "NetworkPorts"),
    STORAGE_GROUPS_COLLECTION("StorageGroupCollection", STORAGE_GROUP, "StorageGroups"),
    ENDPOINT_GROUPS_COLLECTION("EndpointGroupCollection", ENDPOINT_GROUP, "EndpointGroups"),
    CLASS_OF_SERVICES_COLLECTION("ClassOfServiceCollection", CLASS_OF_SERVICE, "ClassesOfService"),
    FILE_SYSTEMS_COLLECTION("FileSystemCollection", FILE_SYSTEM, "FileSystems"),
    STORAGE_POOLS_COLLECTION("StoragePoolCollection", STORAGE_POOL, "StoragePools"),
    ADDRESS_POOLS_COLLECTION("AddressPoolCollection", ADDRESS_POOL, "AddressPools"),
    FILE_SHARES_COLLECTION("FileShareCollection", FILE_SHARE, "ExportedShares"),
    TASKS_COLLECTION("TaskCollection", TASK, "Tasks"),
    ACCELERATION_FUNCTIONS_COLLETION("AccelerationFunctionCollection", ACCELERATION_FUNCTION, "AccelerationFunctions"),
    METRIC_DEFINITIONS_COLLECTION("MetricDefinitionCollection", METRIC_DEFINITION, "MetricDefinitions"),
    METRIC_REPORT_DEFINITIONS_COLLECTION("MetricReportDefinitionCollection", METRIC_REPORT_DEFINITION,
            "MetricReportDefinitions"),
    METRIC_REPORTS_COLLECTION("MetricReportCollection", METRIC_REPORT, "MetricReports"),
    TRIGGERS_COLLECTION("TriggersCollection", TRIGGERS, "Triggers"),
    HOSTED_STORAGE_SERVICES("HostedStorageServices", STORAGE_SERVICE, "StorageServices"),
    CERTIFICATES_COLLECTION("CertificateCollection", CERTIFICATE, "Certificates"),
    SENSORS_COLLECTION("SensorCollection", SENSOR, "Sensors"),
    VLAN_NETWORK_INTERFACES_COLLECTION("VLanNetworkInterfaceCollection", VLAN_NETWORK_INTERFACE, "VLANS"),
    STORAGE_SYSTEMS_COLLECTION("StorageSystemCollection", COMPUTER_SYSTEM, "StorageSystems"),
    JOBS_COLLECTION("JobCollection", JOB, "Jobs"),
    AGGREGATES_COLLECTION("AggregateCollection", AGGREGATE, "Aggregates"),
    AGGREGATION_SOURCES_COLLECTION("AggregationSourceCollection", AGGREGATION_SOURCE, "AggregationSources"),
    CONNECTION_METHODS_COLLECTION("ConnectionMethodCollection", CONNECTION_METHOD, "ConnectionMethods"),
    MEDIA_CONTROLLERS_COLLECTION("MediaControllerCollection", MEDIA_CONTROLLER, "MediaControllers"),
    FACILITIES_COLLECTION("FacilityCollection", FACILITY, "Facility"),
    CIRCUITS_COLLECTION("CircuitCollection", FACILITY, "Circuits"),
    OUTLETS_COLLECTION("OutletCollection", OUTLET, "Outlets"),
    FABRIC_ADAPTERS_COLLECTION("FabricAdapterCollection", OUTLET, "FabricAdapters"),
    CONNECTIONS_COLLECTION("ConnectionCollection", CONNECTION, "Connections"),
    CONSISTENCY_GROUPS_COLLECTION("ConsistencyGroupCollection", CONSISTENCY_GROUP, "ConsistencyGroups"),
    POWER_DISTRIBUTIONS_COLLECTION("PowerDistributionCollection", POWER_DISTRIBUTION, "PowerDistributiona"),
    POWER_DOMAINS_COLLECTION("PowerDomainCollection", POWER_DOMAIN, "PowerDomains"),
    LINES_OF_SERVICE_COLLECTION("LineOfServiceCollection", LINE_OF_SERVICE, "LinesOfService"),
    OPERATING_CONFIGS_COLLECTION("OperatingConfigCollection", OPERATING_CONFIG, "OperatingConfigs"),
    OUTLET_GROUPS_COLLECTION("OutletGroupCollection", OPERATING_CONFIG, "OutletGroups"),
    ROUTE_ENTRIES_COLLECTION("RouteEntryCollection", ROUTE_ENTRY, "RouteEntries"),
    ROUTE_SET_ENTRIES_COLLECTION("RouteSetEntryCollection", ROUTE_SET_ENTRY, "RouteSetEntries"),
    SECURE_BOOT_DATABASES_COLLECTION("SecureBootDatabaseCollection", SECURE_BOOT_DATABASE, "SecureBootDatabases"),
    SIGNATURES_COLLECTION("SecureBootDatabaseCollection", SIGNATURE, "SecureBootDatabases"),
    STORAGE_CONTROLLERS_COLLECTION("StorageControllerCollection", STORAGE_CONTROLLER, "StorageControllers"),
    VCAT_ENTRIES_COLLECTION("VCATEntry", VCAT_ENTRY, "VCATEntries"),

    //Added
    FIRMWARE_INVENTORIES_COLLECTION("FirmwareInventoryCollection", SOFTWARE_INVENTORY, "FirmwareInventory");

    override fun oDataType() = this.odataTypeValue

    override fun jsonName() = this.jsonName

    override fun of(): ResourceTypeBase = of
}

enum class ResourceType(
        private val kind: TypeKind,
        private val odataTypeValue: String = "",
        private val jsonName: String = ""
) : ResourceTypeBase {
    ANY(SIMPLE),
    SERVICE_ROOT(SIMPLE, "ServiceRoot"),
    COMPUTER_SYSTEM(SIMPLE, "ComputerSystem"),
    MANAGER(SIMPLE, "Manager"),
    MANIFEST(SIMPLE, "Manifest"),
    PROCESSOR(SIMPLE, "Processor"),
    MEMORY(SIMPLE, "Memory"),
    ETHERNET_INTERFACE(SIMPLE, "EthernetInterface"),
    CHASSIS(SIMPLE, "Chassis"),
    PCIE_DEVICE(SIMPLE, "PCIeDevice"),
    STORAGE_SERVICE(SIMPLE, "StorageService"),
    RESOURCE_BLOCK(SIMPLE, "ResourceBlock"),
    SWITCH(SIMPLE, "Switch"),
    PORT(SIMPLE, "Port"),
    PCIE_FUNCTION(SIMPLE, "PCIeFunction"),
    NETWORK_ADAPTER(SIMPLE, "NetworkAdapter"),
    DRIVE(SIMPLE, "Drive"),
    NETWORK_DEVICE_FUNCTION(SIMPLE, "NetworkDeviceFunction"),
    SESSION(SIMPLE, "Session"),
    MANAGER_ACCOUNT(SIMPLE, "ManagerAccount"),
    ROLE(SIMPLE, "Role"),
    PRIVILEGE_REGISTRY(SIMPLE, "PrivilegeRegistry"),
    EXTERNAL_ACCOUNT_PROVIDER(SIMPLE, "ExternalAccountProvider"),
    MESSAGE_REGISTRY_FILE(SIMPLE, "MessageRegistryFile"),
    JSON_SCHEMA_FILE(SIMPLE, "JsonSchemaFile"),
    EVENT_DESTINATION(SIMPLE, "EventDestination"),
    FABRIC(SIMPLE, "Fabric"),
    ZONE(SIMPLE, "Zone"),
    ENDPOINT(SIMPLE, "Endpoint"),
    SOFTWARE_INVENTORY(SIMPLE, "SoftwareInventory"),
    SERIAL_INTERFACE(SIMPLE, "SerialInterface"),
    LOG_SERVICE(SIMPLE, "LogService"),
    LOG_ENTRY(SIMPLE, "LogEntry"),
    VIRTUAL_MEDIA(SIMPLE, "VirtualMedia"),
    HOST_INTERFACE(SIMPLE, "HostInterface"),
    BOOT_OPTION(SIMPLE, "BootOption"),
    SIMPLE_STORAGE(SIMPLE, "SimpleStorage"),
    STORAGE(SIMPLE, "Storage"),
    VOLUME(SIMPLE, "Volume"),
    MEMORY_DOMAIN(SIMPLE, "MemoryDomain"),
    MEMORY_CHUNKS(SIMPLE, "MemoryChunks"),
    NETWORK_INTERFACE(SIMPLE, "NetworkInterface"),
    NETWORK_PORT(SIMPLE, "NetworkPort"),
    ETHERNET_SWITCH(SIMPLE, "EthernetSwitch"),
    ETHERNET_SWITCH_PORT(SIMPLE, "EthernetSwitchPort"),
    ETHERNET_SWITCH_ACL(SIMPLE, "EthernetSwitchACL"),
    ETHERNET_SWITCH_ACL_RULE(SIMPLE, "EthernetSwitchACLRule"),
    ETHERNET_SWITCH_STATIC_MAC(SIMPLE, "EthernetSwitchStaticMAC"),
    STORAGE_GROUP(SIMPLE, "StorageGroup"),
    ENDPOINT_GROUP(SIMPLE, "EndpointGroup"),
    CLASS_OF_SERVICE(SIMPLE, "ClassOfService"),
    FILE_SYSTEM(SIMPLE, "FileSystem"),
    STORAGE_POOL(SIMPLE, "StoragePool"),
    ADDRESS_POOL(SIMPLE, "AddressPool"),
    FILE_SHARE(SIMPLE, "FileShare"),
    CAPACITY_SOURCE(SIMPLE, "CapacitySource"),
    DATA_PROTECTION_LINE_OF_SERVICE(SIMPLE, "DataProtectionLineOfService"),
    DATA_SECURITY_LINE_OF_SERVICE(SIMPLE, "DataSecurityLineOfService"),
    DATA_STORAGE_LINE_OF_SERVICE(SIMPLE, "DataStorageLineOfService"),
    IO_PERFORMANCE_LINE_OF_SERVICE(SIMPLE, "IOPerformanceLineOfService"),
    IO_CONNECTIVITY_LINE_OF_SERVICE(SIMPLE, "IOConnectivityLineOfService"),
    TASK(SIMPLE, "Task"),
    ACCELERATION_FUNCTION(SIMPLE, "AccelerationFunction"),
    METRIC_DEFINITION(SIMPLE, "MetricDefinition"),
    METRIC_REPORT_DEFINITION(SIMPLE, "MetricReportDefinition"),
    METRIC_REPORT(SIMPLE, "MetricReport"),
    TRIGGERS(SIMPLE, "Triggers"),
    CERTIFICATE(SIMPLE, "Certificate"),
    SENSOR(SIMPLE, "Sensor"),
    PCIE_SLOTS(SIMPLE, "PCIeSlots"),
    EVENT(SIMPLE, "Event"),
    STORAGE_REPLICA_INFO(SIMPLE, "StorageReplicaInfo"),
    SPARE_RESOURCE_SET(SIMPLE, "SpareResourceSet"),
    VLAN_NETWORK_INTERFACE(SIMPLE, "VLanNetworkInterface"),
    CERTIFICATE_LOCATIONS(SIMPLE, "CertificateLocations"),
    JOB(SIMPLE, "Job"),
    DEVICE(SIMPLE, "Device"),
    AGGREGATE(SIMPLE, "Aggregate"),
    AGGREGATION_SOURCE(SIMPLE, "AggregationSource"),
    CONNECTION_METHOD(SIMPLE, "ConnectionMethod"),
    MEDIA_CONTROLLER(SIMPLE, "MediaController"),
    FACILITY(SIMPLE, "Facility"),
    CIRCUIT(SIMPLE, "Circuit"),
    OUTLET(SIMPLE, "Outlet"),
    FABRIC_ADAPTER(SIMPLE, "FabricAdapter"),
    CONNECTION(SIMPLE, "Connection"),
    CONSISTENCY_GROUP(SIMPLE, "ConsistencyGroup"),
    POWER_DISTRIBUTION(SIMPLE, "PowerDistribution"),
    POWER_DOMAIN(SIMPLE, "PowerDomain"),
    FEATURES_REGISTRY(SIMPLE, "FeaturesRegistry"),
    LINE_OF_SERVICE(SIMPLE, "LineOfService"),
    OPERATING_CONFIG(SIMPLE, "OperatingConfig"),
    OUTLET_GROUP(SIMPLE, "OutletGroup"),
    POWER_EQUIPMENT(SIMPLE, "PowerEquipment"),
    ROUTE_ENTRY(SIMPLE, "RouteEntry"),
    ROUTE_SET_ENTRY(SIMPLE, "RouteSetEntry"),
    SECURE_BOOT_DATABASE(SIMPLE, "SecureBootDatabase"),
    SIGNATURE(SIMPLE, "Signature"),
    VCAT_ENTRY(SIMPLE, "VCATEntry"),

    COMPUTER_SYSTEM_METRICS(SINGLETON, "ComputerSystemMetrics", "Metrics"),
    MEMORY_METRICS(SINGLETON, "MemoryMetrics", "Metrics"),
    PROCESSOR_METRICS(SINGLETON, "ProcessorMetrics", "Metrics"),
    PORT_METRICS(SINGLETON, "PortMetrics", "Metrics"),
    ETHERNET_SWITCH_METRICS(SINGLETON, "EthernetSwitchMetrics", "Metrics"),
    TASK_SERVICE(SINGLETON, "TaskService", "TaskService"),
    SESSION_SERVICE(SINGLETON, "SessionService", "SessionService"),
    ACCOUNT_SERVICE(SINGLETON, "AccountService", "AccountService"),
    EVENT_SERVICE(SINGLETON, "EventService", "EventService"),
    UPDATE_SERVICE(SINGLETON, "UpdateService", "UpdateService"),
    COMPOSITION_SERVICE(SINGLETON, "CompositionService", "CompositionService"),
    MANAGER_NETWORK_PROTOCOL(SINGLETON, "ManagerNetworkProtocol", "NetworkProtocol"),
    POWER(SINGLETON, "Power", "Power"),
    THERMAL(SINGLETON, "Thermal", "Thermal"),
    BIOS(SINGLETON, "Bios", "Bios"),
    SECURE_BOOT(SINGLETON, "SecureBoot", "SecureBoot"),
    ASSEMBLY(SINGLETON, "Assembly", "Assembly"),
    DATA_PROTECTION_LOS_CAPABILITIES(SINGLETON, "DataProtectionLoSCapabilities", "DataProtectionLoSCapabilities"),
    DATA_SECURITY_LOS_CAPABILITIES(SINGLETON, "DataSecurityLoSCapabilities", "DataSecurityLoSCapabilities"),
    DATA_STORAGE_LOS_CAPABILITIES(SINGLETON, "DataStorageLoSCapabilities", "DataStorageLoSCapabilities"),
    IO_CONNECTIVITY_LOS_CAPABILITIES(SINGLETON, "IOConnectivityLoSCapabilities", "IOConnectivityLoSCapabilities"),
    IO_PERFORMANCE_LOS_CAPABILITIES(SINGLETON, "IOPerformanceLoSCapabilities", "IOPerformanceLoSCapabilities"),
    TELEMETRY_SERVICE(SINGLETON, "TelemetryService", "TelemetryService"),
    JOB_SERVICE(SINGLETON, "JobService", "JobService"),
    CERTIFICATE_SERVICE(SINGLETON, "CertificateService", "CertificateService"),
    BIOS_SETTINGS(SINGLETON, "Settings", "Settings"),
    AGGREGATION_SERVICE(SINGLETON, "AggregationService", "AggregationSevice"),
    POWER_DISTRIBUTION_METRICS(SINGLETON, "PowerDistributionMetrics", "Metrics"),

    POWER_CONTROL(EMBEDDED, jsonName = "PowerControl"),
    VOLTAGE(EMBEDDED, jsonName = "Voltages"),
    POWER_SUPPLY(EMBEDDED, jsonName = "PowerSupplies"),
    TEMPERATURE(EMBEDDED, jsonName = "Temperatures"),
    FAN(EMBEDDED, jsonName = "Fans"),
    STORAGE_CONTROLLER(EMBEDDED, jsonName = "StorageControllers"),
    ASSEMBLY_DATA(EMBEDDED, jsonName = "Assemblies"),
    REDUNDANCY(EMBEDDED, jsonName = "Redundancy"),
    EVENT_RECORD(EMBEDDED, jsonName = "EventRecord");

    override fun oDataType() = this.odataTypeValue

    override fun jsonName() = this.jsonName

    override fun isEmbedded() = this.kind == EMBEDDED

    override fun isSingleton() = this.kind == SINGLETON

    override fun actionsNamespace() = this.odataTypeValue
}

enum class ActionType(val actionName: String) {
    ADD_ELEMENTS("AddElements"),
    ADD_ENDPOINT("AddEndpoint"),
    ADD_RESOURCE_BLOCK("AddResourceBlock"),
    ASSIGN_REPLICA_TARGET("AssignReplicaTarget"),
    BREAKER_CONTROL("BreakerControl"),
    CHANGE_PASSWORD("ChangePassword"),
    CHANGE_RAID_LAYOUT("ChangeRAIDLayout"),
    CHECK_CONSISTENCY("CheckConsistency"),
    CLEAR_CURRENT_PERIOD("ClearCurrentPeriod"),
    CLEAR_LOG("ClearLog"),
    COLLECT_DIAGNOSTIC_DATA("CollectDiagnosticData"),
    CREATE_REPLICA_TARGET("CreateReplicaTarget"),
    DISABLE_PASSPHRASE("DisablePassphrase"),
    EJECT_MEDIA("EjectMedia"),
    EXPOSE_VOLUMES("ExposeVolumes"),
    FORCE_ENABLE("ForceEnable"),
    FORCE_FAILOVER("ForceFailover"),
    GENERATE_CSR("GenerateCSR"),
    HIDE_VOLUMES("HideVolumes"),
    INITIALIZE("Initialize"),
    INSERT_MEDIA("InsertMedia"),
    MODIFY_REDUNDANCY_SET("ModifyRedundancySet"),
    OVERWRITE_UNIT("OverwriteUnit"),
    POWER_CONTROL("PowerControl"),
    POWER_SUPPLY_RESET("PowerSupplyReset"),
    REKEY("Rekey"),
    REMOVE_ELEMENTS("RemoveElements"),
    REMOVE_ENDPOINT("RemoveEndpoint"),
    REMOVE_REPLICA_RELATIONSHIP("RemoveReplicaRelationship"),
    REMOVE_RESOURCE_BLOCK("RemoveResourceBlock"),
    RENEW("Renew"),
    REPLACE_CERTIFICATE("ReplaceCertificate"),
    RESET("Reset"),
    RESET_BIOS("ResetBios"),
    RESET_KEYS("ResetKeys"),
    RESET_METRICS("ResetMetrics"),
    RESET_SETTINGS_TO_DEFAULT("ResetSettingsToDefault"),
    RESET_TO_DEFAULTS("ResetToDefaults"),
    RESUME_REPLICATION("ResumeReplication"),
    RESUME_SUBSCRIPTION("ResumeSubscription"),
    REVERSE_REPLICATION_RELATIONSHIP("ReverseReplicationRelationship"),
    SECURE_ERASE("SecureErase"),
    SECURE_ERASE_UNIT("SecureEraseUnit"),
    SET_DEFAULT_BOOT_ORDER("SetDefaultBootOrder"),
    SET_ENCRYPTION_KEY("SetEncryptionKey"),
    SET_PASSPHRASE("SetPassphrase"),
    SIMPLE_UPDATE("SimpleUpdate"),
    SPLIT_REPLICATION("SplitReplication"),
    START_UPDATE("StartUpdate"),
    SUBMIT_TEST_EVENT("SubmitTestEvent"),
    SUBMIT_TEST_METRIC_REPORT("SubmitTestMetricReport"),
    SUSPEND_REPLICATION("SuspendReplication"),
    TRANSFER_CONTROL("TransferControl"),
    UNLOCK_UNIT("UnlockUnit"),
    UPDATE_BIOS("UpdateBIOS"),
    UPDATE_BIOS_BACKUP("UpdateBIOS.Backup"),
    UPDATE_BMC("UpdateBMC"),
    UPDATE_ME("UpdateME"),
    UPDATE_SDR("UpdateSDR")
}

@Suppress("LargeClass")
enum class EmbeddedObjectType : ExtensibleEmbeddedObjectType {
    EMPTY,
    ACCOUNT_SERVICE_AUTHENTICATION,
    ACCOUNT_SERVICE_LDAP_SEARCH_SETTINGS,
    ACCOUNT_SERVICE_LDAP_SERVICE,
    ACCOUNT_SERVICE_ROLE_MAPPING,
    ADDRESS_POOL_ETHERNET,
    ADDRESS_POOL_GEN_Z,
    ALARM_TRIPS,
    ANA_CHARACTERISTICS,
    AS_NUMBER_RANGE,
    ATTRIBUTES,
    BASE_SPEED_PRIORITY_SETTINGS,
    BFD_SINGLE_HOP_ONLY,
    BGP_EVPN,
    BGP_NEIGHBOR,
    BGP_ROUTE,
    BOOT,
    BOOT_PROGRESS,
    BOOT_TARGETS,
    CACHE_METRICS,
    CALCULATION_PARAMS_TYPE,
    CAPABILITY,
    CAPACITY,
    CAPACITY_INFO,
    CERTIFICATE_IDENTIFIER,
    CHAP_INFORMATION,
    CIRCUIT_CURRENT_SENSORS,
    CIRCUIT_VOLTAGE_SENSORS,
    COLLECTION_CAPABILITIES,
    COMMAND_SHELL,
    COMMON_BGP_PROPERTIES,
    COMPOSITION_STATUS,
    CONFIGURED_NETWORK_LINK,
    CONNECTED_ENTITY,
    CONTACT_INFO,
    CONTROLLERS,
    CONTROLLER_CAPABILITIES,
    CORE_METRICS,
    CREDENTIAL_BOOTSTRAPPING,
    CURRENT_PERIOD,
    C_STATE_RESIDENCY,
    DATA_CENTER_BRIDGING,
    DEEP_OPERATIONS,
    DEVICE,
    DHCHAP_INFORMATION,
    DHCP,
    DHCPV4_CONFIGURATION,
    DHCPV6_CONFIGURATION,
    DISCRETE_TRIGGER,
    EBGP,
    ENDPOINT_GEN_Z,
    END_GRP_LIFETIME,
    ENERGY_SENSORS,
    ENGINE_ID,
    ENUMERATION_MEMBER,
    ESI_NUMBER_RANGE,
    ETHERNET,
    ETHERNET_INTERFACE,
    ETHERNET_PROPERTIES,
    EVENT_DESTINATION_SYSLOG_FILTER,
    EVI_NUMBER_RANGE,
    EXPAND,
    EXTERNAL_ACCOUNT_PROVIDER,
    EXTERNAL_ACCOUNT_PROVIDER_AUTHENTICATION,
    EXTERNAL_ACCOUNT_PROVIDER_LDAP_SEARCH_SETTINGS,
    EXTERNAL_ACCOUNT_PROVIDER_LDAP_SERVICE,
    EXTERNAL_ACCOUNT_PROVIDER_ROLE_MAPPING,
    FABRIC_ADAPTER_GEN_Z,
    FEATURES_REGISTRY_PROPERTY,
    FIBRE_CHANNEL,
    FIBRE_CHANNEL_PROPERTIES,
    FPGA,
    FPGA_RECONFIGURATION_SLOT,
    GCID,
    GENERATE_CSR_RESPONSE,
    GRACEFUL_RESTART,
    GRAPHICAL_CONSOLE,
    HEALTH_DATA,
    HOSTED_SERVICES,
    HOST_GRAPHICAL_CONSOLE,
    HOST_SERIAL_CONSOLE,
    HTTPS_PROTOCOL,
    HTTP_HEADER_PROPERTY,
    HTTP_PUSH_URI_APPLY_TIME,
    HTTP_PUSH_URI_OPTIONS,
    IMPORTED_SHARE,
    INFINI_BAND,
    INPUT_RANGE,
    INTERLEAVE_SET,
    IO_STATISTICS,
    IO_WORKLOAD,
    IO_WORKLOAD_COMPONENT,
    IP_TRANSPORT_DETAILS,
    IP_V4,
    IP_V4_ADDRESS,
    IP_V4_ADDRESS_RANGE,
    IP_V6_ADDRESS,
    IP_V6_ADDRESS_POLICY_ENTRY,
    IP_V6_GATEWAY_STATIC_ADDRESS,
    IP_V6_STATIC_ADDRESS,
    ISCSI_BOOT,
    JOB_PAYLOAD,
    JOB_SERVICE_CAPABILITIES,
    JSON_SCHEMA_FILE_LOCATION,
    LIFE_TIME,
    LINK_CONFIGURATION,
    LOG_SERVICE_SYSLOG_FILTER,
    MAINTENANCE_WINDOW,
    MANAGER_SERVICE,
    MANIFEST_ENTRY,
    MAPPED_VOLUME,
    MAPPING,
    MAX_PREFIX,
    MEMORY_LOCATION,
    MEMORY_SET,
    MEMORY_SUMMARY,
    MESSAGE,
    MESSAGE_REGISTRY_FILE_LOCATION,
    METRIC,
    METRIC_DEFINITION_WILDCARD,
    METRIC_REPORT_DEFINITION_WILDCARD,
    METRIC_REPORT_METRIC_VALUE,
    METRIC_VALUE,
    MULTIPLE_PATHS,
    NAMESPACE_FEATURES,
    NETWORK_DEVICE_FUNCTION_ETHERNET,
    NET_DEV_FUNC_MAX_B_W_ALLOC,
    NET_DEV_FUNC_MIN_B_W_ALLOC,
    NIC_PARTITIONING,
    NPIV,
    NTP_PROTOCOL,
    NVME_CONTROLLER_ATTRIBUTES,
    NVME_CONTROLLER_PROPERTIES,
    NVME_ENDURANCE_GROUP_PROPERTIES,
    NVME_NAMESPACE_PROPERTIES,
    NVME_SET_PROPERTIES,
    NVME_SMART_CRITICAL_WARNINGS,
    OPERATION,
    OPERATIONS,
    OPERATION_APPLY_TIME_SUPPORT,
    OPERATION_MAP,
    OPERATION_PRIVILEGE,
    OUTLET_CURRENT_SENSORS,
    OUTLET_VOLTAGE_SENSORS,
    PART_LOCATION,
    PCIE_INTERFACE,
    PCIE_SLOT,
    PCI_ID,
    PHYSICAL_SECURITY,
    PLACEMENT,
    PORT_GEN_Z,
    PORT_METRICS_GEN_Z,
    POSTAL_ADDRESS,
    POWER_LIMIT,
    POWER_MANAGEMENT_POLICY,
    POWER_METRIC,
    POWER_SENSORS,
    PREFERRED_APPLY_TIME,
    PROCESSOR_ID,
    PROCESSOR_INTERFACE,
    PROCESSOR_MEMORY,
    PROCESSOR_SUMMARY,
    PROPERTY_PATTERN,
    PROTOCOL,
    PROTOCOL_FEATURES_SUPPORTED,
    REGION_SET,
    REKEY_RESPONSE,
    RENEW_RESPONSE,
    REPLICA_INFO,
    REPLICA_REQUEST,
    RESOURCE_BLOCK_LIMITS,
    RESOURCE_IDENTIFIER,
    RESOURCE_LOCATION,
    REVISION_TYPE,
    ROUTE_DISTINGUISHER_RANGE,
    ROUTE_TARGET_RANGE,
    SCHEDULE,
    SECURITY_CAPABILITIES,
    SENSOR_THRESHOLD,
    SENSOR_THRESHOLDS,
    SERIAL_CONSOLE,
    SERIAL_CONSOLE_PROTOCOL,
    SETTINGS,
    SMTP,
    SNMP_COMMUNITY,
    SNMP_PROTOCOL,
    SNMP_SETTINGS,
    SNMP_USER_INFO,
    SRIOV,
    SSD_PROTOCOL,
    SSE_FILTER_PROPERTIES_SUPPORTED,
    STATELESS_ADDRESS_AUTO_CONFIGURATION,
    STATUS,
    STORAGE_CACHE_SUMMARY,
    STORAGE_CONTROLLER_CACHE_SUMMARY,
    STORAGE_CONTROLLER_RATES,
    STORAGE_RATES,
    SUPPORTED_FEATURE,
    SUPPORTED_LINK_CAPABILITIES,
    SYSTEM_CPU_PERFORMANCE_CONFIGURATION,
    TARGET_PRIVILEGE_MAP,
    TASK_MESSAGE,
    TASK_PAYLOAD,
    TELEMETRY_SERVICE_METRIC_VALUE,
    TRANSFER_CONFIGURATION,
    TRANSFER_CRITERIA,
    TRIGGERS_THRESHOLD,
    TRIGGERS_THRESHOLDS,
    TRIGGERS_WILDCARD,
    TRUSTED_MODULES,
    TURBO_PROFILE_DATAPOINT,
    UPDATE_PARAMETERS,
    VCA_TABLE_ENTRY,
    VIRTUALIZATION_OFFLOAD,
    VIRTUAL_FUNCTION,
    VIRTUAL_MEDIA_CONFIG,
    VLAN,
    VLAN_IDENTIFIER_ADDRESS_RANGE,
    VOLUME_INFO,
    WATCHDOG_TIMER
}
