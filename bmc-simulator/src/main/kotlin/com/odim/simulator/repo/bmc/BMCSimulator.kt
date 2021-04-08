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

package com.odim.simulator.repo.bmc

import com.odim.simulator.CoreConfig.SERVER_BASIC_CREDENTIALS
import com.odim.simulator.RedfishSimulator
import com.odim.simulator.SimulatorConfig.Config.getConfigProperty
import com.odim.simulator.SimulatorName
import com.odim.simulator.behaviors.Behavior.Companion.behavior
import com.odim.simulator.behaviors.BehaviorResponse.Companion.nonTerminal
import com.odim.simulator.behaviors.BehaviorResponse.Companion.terminal
import com.odim.simulator.http.HttpMethod.DELETE
import com.odim.simulator.http.HttpMethod.GET
import com.odim.simulator.http.HttpMethod.PATCH
import com.odim.simulator.http.HttpMethod.POST
import com.odim.simulator.http.Response.Companion.unauthorized
import com.odim.simulator.repo.behaviors.GetOnChassis
import com.odim.simulator.repo.behaviors.GetOnComputerSystem
import com.odim.simulator.repo.behaviors.GetOnLogEntriesCollection
import com.odim.simulator.repo.bmc.behaviors.ComputerSystemMetrics
import com.odim.simulator.repo.bmc.behaviors.GetOnAccounts
import com.odim.simulator.repo.bmc.behaviors.GetOnManager
import com.odim.simulator.repo.bmc.behaviors.GetOnRegistryMessage
import com.odim.simulator.repo.bmc.behaviors.GetOnStorage
import com.odim.simulator.repo.bmc.behaviors.GetOnTask
import com.odim.simulator.repo.bmc.behaviors.PatchOnBiosSettings
import com.odim.simulator.repo.bmc.behaviors.PatchOnChassis
import com.odim.simulator.repo.bmc.behaviors.PatchOnComputerSystem
import com.odim.simulator.repo.bmc.behaviors.PatchOnEthernetInterface
import com.odim.simulator.repo.bmc.behaviors.PatchOnManager
import com.odim.simulator.repo.bmc.behaviors.PatchOnManagerNetworkProtocol
import com.odim.simulator.repo.bmc.behaviors.PatchOnPower
import com.odim.simulator.repo.bmc.behaviors.PostOnAccounts
import com.odim.simulator.repo.bmc.behaviors.PostOnFirmwareUpdate
import com.odim.simulator.repo.bmc.behaviors.PostOnSimpleUpdate
import com.odim.simulator.repo.bmc.behaviors.ProcessorMetrics
import com.odim.simulator.repo.bmc.behaviors.ResetOnSystem
import com.odim.simulator.repo.bmc.behaviors.SetSuccessfulMessageResponse
import com.odim.simulator.repo.bmc.configurators.AccountServiceConfigurator.Factory.configureAccountService
import com.odim.simulator.repo.bmc.configurators.ChassisConfigurator.Factory.appendPower
import com.odim.simulator.repo.bmc.configurators.ChassisConfigurator.Factory.appendThermal
import com.odim.simulator.repo.bmc.configurators.ChassisConfigurator.Factory.appendVoltageVolt
import com.odim.simulator.repo.bmc.configurators.ChassisConfigurator.Factory.configureChassis
import com.odim.simulator.repo.bmc.configurators.ChassisConfigurator.Factory.configureMainChassis
import com.odim.simulator.repo.bmc.configurators.ManagerConfigurator.Factory.configureEthernetInterfaceForManager
import com.odim.simulator.repo.bmc.configurators.ManagerConfigurator.Factory.configureLogServiceForManager
import com.odim.simulator.repo.bmc.configurators.ManagerConfigurator.Factory.configureManager
import com.odim.simulator.repo.bmc.configurators.MessageRegistryConfigurator.Factory.configureMessageRegistry
import com.odim.simulator.repo.bmc.configurators.SessionServiceConfigurator.Factory.configureSessionService
import com.odim.simulator.repo.bmc.configurators.UpdateServiceConfigurator.Factory.configureUpdateService
import com.odim.simulator.tree.ResourceTree
import com.odim.simulator.tree.ResourceVersion
import com.odim.simulator.tree.structure.ActionElement
import com.odim.simulator.tree.structure.ActionType.RESET
import com.odim.simulator.tree.structure.ActionType.SIMPLE_UPDATE
import com.odim.simulator.tree.structure.ActionType.UPDATE_BIOS
import com.odim.simulator.tree.structure.ActionType.UPDATE_BMC
import com.odim.simulator.tree.structure.ActionType.UPDATE_ME
import com.odim.simulator.tree.structure.ActionType.UPDATE_SDR
import com.odim.simulator.tree.structure.LinkableResourceArray
import com.odim.simulator.tree.structure.Resource
import com.odim.simulator.tree.structure.ResourceCollection
import com.odim.simulator.tree.structure.ResourceCollectionType.LOG_ENTRIES_COLLECTION
import com.odim.simulator.tree.structure.ResourceCollectionType.MANAGER_ACCOUNTS_COLLECTION
import com.odim.simulator.tree.structure.ResourceType.BIOS_SETTINGS
import com.odim.simulator.tree.structure.ResourceType.BOOT_OPTION
import com.odim.simulator.tree.structure.ResourceType.CERTIFICATE_LOCATIONS
import com.odim.simulator.tree.structure.ResourceType.CHASSIS
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM
import com.odim.simulator.tree.structure.ResourceType.COMPUTER_SYSTEM_METRICS
import com.odim.simulator.tree.structure.ResourceType.DRIVE
import com.odim.simulator.tree.structure.ResourceType.ETHERNET_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.FAN
import com.odim.simulator.tree.structure.ResourceType.HOST_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.LOG_ENTRY
import com.odim.simulator.tree.structure.ResourceType.LOG_SERVICE
import com.odim.simulator.tree.structure.ResourceType.MANAGER
import com.odim.simulator.tree.structure.ResourceType.MANAGER_ACCOUNT
import com.odim.simulator.tree.structure.ResourceType.MANAGER_NETWORK_PROTOCOL
import com.odim.simulator.tree.structure.ResourceType.MEMORY
import com.odim.simulator.tree.structure.ResourceType.MESSAGE_REGISTRY_FILE
import com.odim.simulator.tree.structure.ResourceType.NETWORK_ADAPTER
import com.odim.simulator.tree.structure.ResourceType.NETWORK_DEVICE_FUNCTION
import com.odim.simulator.tree.structure.ResourceType.NETWORK_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.NETWORK_PORT
import com.odim.simulator.tree.structure.ResourceType.PCIE_DEVICE
import com.odim.simulator.tree.structure.ResourceType.POWER
import com.odim.simulator.tree.structure.ResourceType.POWER_CONTROL
import com.odim.simulator.tree.structure.ResourceType.POWER_SUPPLY
import com.odim.simulator.tree.structure.ResourceType.PROCESSOR
import com.odim.simulator.tree.structure.ResourceType.PROCESSOR_METRICS
import com.odim.simulator.tree.structure.ResourceType.REDUNDANCY
import com.odim.simulator.tree.structure.ResourceType.ROLE
import com.odim.simulator.tree.structure.ResourceType.SOFTWARE_INVENTORY
import com.odim.simulator.tree.structure.ResourceType.STORAGE
import com.odim.simulator.tree.structure.ResourceType.STORAGE_CONTROLLER
import com.odim.simulator.tree.structure.ResourceType.TASK
import com.odim.simulator.tree.structure.ResourceType.TEMPERATURE
import com.odim.simulator.tree.structure.ResourceType.VIRTUAL_MEDIA
import com.odim.simulator.tree.structure.ResourceType.VLAN_NETWORK_INTERFACE
import com.odim.simulator.tree.structure.ResourceType.VOLTAGE
import com.odim.simulator.tree.structure.ResourceType.VOLUME
import com.odim.simulator.tree.structure.TreeElement
import com.odim.simulator.tree.templates.bmc.BmcVersion.BMC_1_0
import io.javalin.core.security.BasicAuthCredentials
import java.util.UUID.randomUUID

object FirmwareUpdateMessages {
    const val UPDATE_STARTED_SUFFIX = ", instance:0- Update Started - OEM "
    const val UPDATE_COMPLETED_SUFFIX = ", instance:0- Update Completed Successfully - OEM "
}

@SimulatorName("BMC")
class BMCSimulator(val resourceVersion: ResourceVersion = BMC_1_0,
                   extendingPackages: List<String> = listOf())
    : RedfishSimulator(resourceVersion, mutableListOf("bmc").apply {
    addAll(extendingPackages)
    }) {
    val basicAuthUsername = getUsername()
    val basicAuthPassword = getPassword()
    private val serial = randomUUID().toString()
    private val drive1 = create(DRIVE) { "Id" to "DRV1" }
    private val volume1 = create(VOLUME)
    private val logServiceForManager = configureLogServiceForManager(create(LOG_SERVICE), createResourcesList(25, LOG_ENTRY))
    private val ethInterfaceForManager = configureEthernetInterfaceForManager(create(ETHERNET_INTERFACE))
    val manager = configureManager(create(MANAGER))

    val system = create(COMPUTER_SYSTEM) {
        "Name" to "BMC Simulator System"
        "UUID" to randomUUID().toString()
        "Description" to "Computer system providing compute resources"
        "SerialNumber" to serial
        "IndicatorLED" to "Off"
        "PowerState" to "On"
        "Manufacturer" to "Intel Corporation"
        "Model" to "S2600WFQ"
        "PartNumber" to "H97440-010"
        "AssetTag" to "WP2 BMC"
        "SKU" to "R2208WFTZS"
        "BiosVersion" to "SE5C620.86B.02.01.0008.03192019-20201559"
        "Boot" to {
            "BootSourceOverrideEnabled" to "Disabled"
            "BootSourceOverrideTarget" to "None"
            "BootSourceOverrideMode" to "Legacy"
            "BootSourceOverrideTarget@Redfish.AllowableValues" to array[
                    "None",
                    "Pxe",
                    "Hdd",
                    "Cd",
                    "BiosSetup",
                    "UefiShell",
                    "Usb"
            ]
        }
    }

    private val chassisRackMount = configureMainChassis(create(CHASSIS) { "UUID" to systemUuid})

    private val chassisBaseboard = configureChassis("Baseboard", create(CHASSIS), "CARD001", "1")
            .apply {
                appendThermal(this, createResourcesList(3, TEMPERATURE), create(FAN), create(REDUNDANCY))
                appendPower(this, create(REDUNDANCY), create(POWER_CONTROL), create(POWER_SUPPLY), system, chassisRackMount)
                appendVoltageVolt(this, create(VOLTAGE), system, chassisRackMount)
            }

    private val chassisFrontPanel = configureChassis("FrontPanel", create(CHASSIS), "FP001", "2")
            .apply {
                appendThermal(this, createResourcesList(3, TEMPERATURE), create(FAN), create(REDUNDANCY))
                appendPower(this, create(REDUNDANCY), create(POWER_CONTROL), null, system, chassisRackMount)
                appendVoltageVolt(this, create(VOLTAGE), system, chassisRackMount)
            }

    val processor = createProcessor("cpu_1", "cpu_1")
    val processor2 = createProcessor("cpu_2", "cpu_2")

    val memory = create(MEMORY) {
        "Id" to "Memory1"
        "Name" to "Memory 1"
        "Description" to "System Memory"
        "MemoryType" to "DRAM"
        "MemoryDeviceType" to "DDR4"
        "BaseModuleType" to "RDIMM"
        "CapacityMiB" to 16384
        "DataWidthBits" to 64
        "BusWidthBits" to 72
        "Manufacturer" to "Micron"
        "SerialNumber" to "16A7E5A3"
        "PartNumber" to "18ASF2G72PDZ-2G6B1"
        "AllowedSpeedsMHz" to array[2666]
        "MemoryMedia" to array["DRAM"]
        "RankCount" to 2
        "DeviceLocator" to "CPU1_DIMM_A1"
        "ErrorCorrection" to "MultiBitECC"
        "OperatingSpeedMhz" to 2666
    }

    private val storageController = create(STORAGE_CONTROLLER)
    private val storage = create(STORAGE) { "Id" to "BMCStorage" }
    private val bootOption = create(BOOT_OPTION)
    private val ethernetInterface = create(ETHERNET_INTERFACE)
    private val networkInterface = create(NETWORK_INTERFACE)
    private val vlan = create(VLAN_NETWORK_INTERFACE)
    private val pcieDevice = create(PCIE_DEVICE)
    private val networkAdapter = create(NETWORK_ADAPTER)
    private val networkPort = create(NETWORK_PORT)
    private val networkDeviceFunction = create(NETWORK_DEVICE_FUNCTION)
    private val logServiceForSystem = createLogServiceForSystem()
    private val hostInterface = create(HOST_INTERFACE)
    private val virtualMedia = create(VIRTUAL_MEDIA)

    var powerState
        get() = system.traverse<String>("PowerState")
        set(value) {
            system {
                "PowerState" to value
            }
        }

    val systemUuid get() = system.traverse<String>("UUID")
    val accountService get() = root.traverse<Resource>("AccountService")
    val taskService get() = root.traverse<Resource>("TaskService")
    val bios get() = system.traverse<Resource>("Bios")
    val updateService get() = root.traverse<Resource>("UpdateService")
    val sessionService get() = root.traverse<Resource>("SessionService")

    init {
        root(
                chassisRackMount(
                    networkAdapter(networkDeviceFunction, networkPort)
                ),
                chassisBaseboard,
                chassisFrontPanel,
                system(
                    processor,
                    processor2,
                    memory,
                    storage(drive1, storageController, volume1),
                    logServiceForSystem,
                    bootOption,
                    ethernetInterface(vlan),
                    networkInterface,
                    pcieDevice),
                manager(
                    logServiceForManager,
                    ethInterfaceForManager,
                    hostInterface,
                    virtualMedia)
        )

        configureAccountService(accountService, createResourcesList(5, ROLE), createResourcesList(2, MANAGER_ACCOUNT))
        configureUpdateService(updateService, createResourcesList(5, SOFTWARE_INVENTORY))
        configureSessionService(sessionService)
        configureMessageRegistry(root.traverse("Registries"), createResourcesList(5, MESSAGE_REGISTRY_FILE))
        createCertificateService()
        createTaskService()

        link(chassisRackMount, system)
        link(chassisBaseboard, system)
        link(chassisFrontPanel, system)

        link(chassisRackMount, storage)
        link(storage, drive1)
//        link(system, pcie_device)
        link(manager, system)
        link(manager, chassisRackMount, "ManagerForChassis", "ManagedBy")
        link(manager, chassisBaseboard, "ManagerForChassis", "ManagedBy")
        link(manager, chassisFrontPanel, "ManagerForChassis", "ManagedBy")

        chassisRackMount.traverse<LinkableResourceArray>("Links/Contains").addLink(chassisBaseboard, chassisRackMount)
        chassisRackMount.traverse<LinkableResourceArray>("Links/Contains").addLink(chassisFrontPanel, chassisRackMount)

//        ODIM have problem with recognize ContainedBy property
//        chassisBaseboard.traverse<LinkableResource>("Links/ContainedBy").addLink(chassisRackMount, chassisBaseboard)
//        chassisFrontPanel.traverse<LinkableResource>("Links/ContainedBy").addLink(chassisRackMount, chassisFrontPanel)

        behaviors.prependBehavior(GET, securedBehavior())
        behaviors.prependBehavior(POST, securedBehavior())
        behaviors.prependBehavior(PATCH, securedBehavior())
        behaviors.prependBehavior(DELETE, securedBehavior())

        behaviors.replaceActionBehavior(RESET, ResetOnSystem(logServiceForSystem, bios))
        behaviors.appendBehavior(COMPUTER_SYSTEM, GET, GetOnComputerSystem())
        behaviors.appendBehavior(STORAGE, GET, GetOnStorage())
        behaviors.prependBehavior(CHASSIS, GET, GetOnChassis())
        behaviors.replaceBehavior(CHASSIS, PATCH, PatchOnChassis())
        behaviors.replaceBehavior(POWER, PATCH, PatchOnPower())
        behaviors.appendBehavior(LOG_ENTRIES_COLLECTION, GET, GetOnLogEntriesCollection())
        behaviors.appendBehavior(COMPUTER_SYSTEM_METRICS, GET, ComputerSystemMetrics())
        behaviors.appendBehavior(PROCESSOR_METRICS, GET, ProcessorMetrics())
        behaviors.appendBehavior(MESSAGE_REGISTRY_FILE, GET, GetOnRegistryMessage())
        behaviors.replaceBehavior(MANAGER_ACCOUNTS_COLLECTION, POST, PostOnAccounts())
        behaviors.appendBehavior(MANAGER_ACCOUNT, GET, GetOnAccounts())
        behaviors.appendBehavior(MANAGER, GET, GetOnManager())
        behaviors.replaceBehavior(MANAGER, PATCH, PatchOnManager())
        behaviors.replaceBehavior(MANAGER_NETWORK_PROTOCOL, PATCH, PatchOnManagerNetworkProtocol())
        behaviors.replaceBehavior(ETHERNET_INTERFACE, PATCH, PatchOnEthernetInterface())
        behaviors.appendBehavior(TASK, GET, GetOnTask())
        behaviors.replaceBehavior(BIOS_SETTINGS, PATCH, PatchOnBiosSettings())
        behaviors.appendActionBehavior(UPDATE_BMC, POST, PostOnFirmwareUpdate(logServiceForSystem, manager))
        behaviors.appendActionBehavior(UPDATE_BIOS, POST, PostOnFirmwareUpdate(logServiceForSystem, manager, true))
        behaviors.appendActionBehavior(UPDATE_ME, POST, PostOnFirmwareUpdate(logServiceForSystem, manager, true))
        behaviors.appendActionBehavior(UPDATE_SDR, POST, PostOnFirmwareUpdate(logServiceForSystem, manager))
        behaviors.appendActionBehavior(SIMPLE_UPDATE, POST, PostOnSimpleUpdate())
        behaviors.prependBehavior(COMPUTER_SYSTEM, PATCH, PatchOnComputerSystem())
        behaviors.appendBehavior(COMPUTER_SYSTEM, PATCH, SetSuccessfulMessageResponse())
        behaviors.appendBehavior(CHASSIS, PATCH, SetSuccessfulMessageResponse())

        val biosSettings = bios.traverse<Resource>("@Redfish.Settings/SettingsObject")
        biosSettings {
            "SettingsObject" to {
                "@odata.id" to biosSettings.toLink()
            }
        }
    }

    private fun securedBehavior() = behavior {
        val link = if (item is ActionElement) item.toLink() else (item as TreeElement).toLink()
        val validCredentials = BasicAuthCredentials(basicAuthUsername, basicAuthPassword)
        if (link == "/redfish/v1" || request.basicAuthCredentials?.equals(validCredentials) == true) {
            nonTerminal(response)
        } else {
            val unauthorizedResponse = unauthorized().also {
                it.headers["WWW-Authenticate"] = listOf("Basic realm=\"BMC Simulator\"")
            }
            terminal(unauthorizedResponse)
        }
    }

    private fun createTaskService() {
        taskService(create(TASK))
    }

    private fun createLogServiceForSystem(): Resource {
        val logService = create(LOG_SERVICE) { "Id" to "SEL" }
        val logEntries = createResourcesList(3, LOG_ENTRY)
        (logEntries).map { logService.traverse<ResourceCollection>("Entries").add(it) }

        return logService
    }

    private fun createCertificateService() {
        val certificateService = root.traverse<Resource>("CertificateService")
        val certificateLocation = create(CERTIFICATE_LOCATIONS)
        link(certificateService, certificateLocation)
    }

    private fun createProcessor(id: String, socket: String): Resource = create(PROCESSOR) {
        "Id" to id
        "Socket" to socket
        "ProcessorType" to "CPU"
        "ProcessorArchitecture" to "x86"
        "InstructionSet" to "x86-64"
        "Manufacturer" to "Intel(R) Corporation"
        "Model" to "Intel Xeon processor"
        "MaxSpeedMHz" to 4000
        "TotalCores" to 18
        "TotalThreads" to 36
    }

    private fun getUsername() = getConfigProperty<String>(SERVER_BASIC_CREDENTIALS).split(":").first()
    private fun getPassword() = getConfigProperty<String>(SERVER_BASIC_CREDENTIALS).split(":").last()

}

fun createSelEntry(logService: Resource, tree: ResourceTree, message: String) {
    logService.append(tree.create(LOG_ENTRY) {
        "Message" to message
    })
}

