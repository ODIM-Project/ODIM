
# Table of contents

- [Introduction](#introduction)
  * [Resource Aggregator for ODIM logical architecture](#resource-aggregator-for-odim-logical-architecture)
- [API usage and access guidelines](#api-usage-and-access-guidelines)
- [Support for URL Encoding](#support-for-url-encoding)
- [List of supported APIs](#list-of-supported-apis)
  * [Viewing the list of supported Redfish services](#viewing-the-list-of-supported-redfish-services)
- [HTTP request methods, responses, and status codes](#http-request-methods-responses-and-status-codes)
- [Authentication and authorization](#authentication-and-authorization)
  * [Authentication methods for Redfish APIs](#authentication-methods-for-redfish-apis)
  * [Role-based authorization](#role-based-authorization)
- [Sessions](#sessions)
  * [Viewing the session service root](#viewing-the-session-service-root)
  * [Creating a session](#creating-a-session)
  * [Listing sessions](#listing-sessions)
  * [Viewing information about a single session](#viewing-information-about-a-single-session)
  * [Deleting a session](#deleting-a-session)
- [User roles and privileges](#user-roles-and-privileges)
  * [Viewing the account service root](#viewing-the-account-service-root)
  * [Creating a role](#creating-a-role)
  * [Listing roles](#listing-roles)
  * [Viewing information about a role](#viewing-information-about-a-role)
  * [Updating a role](#updating-a-role)
  * [Deleting a role](#deleting-a-role)
- [User accounts](#user-accounts)
  * [Creating a user account](#creating-a-user-account)
    + [Password requirements](#password-requirements)
  * [Listing user accounts](#listing-user-accounts)
  * [Viewing the account details](#viewing-the-account-details)
  * [Updating a user account](#updating-a-user-account)
  * [Deleting a user account](#deleting-a-user-account)
- [Resource aggregation and management](#resource-aggregation-and-management)
  * [Viewing the aggregation service root](#viewing-the-aggregation-service-root)
  * [Connection methods](#connection-methods)
    + [Viewing a collection of connection methods](#viewing-a-collection-of-connection-methods)
    + [Viewing a connection method](#viewing-a-connection-method)
      - [Connection method variants](#connection-method-variants)
  * [Adding a plugin as an aggregation source](#adding-a-plugin-as-an-aggregation-source)
  * [Adding a server as an aggregation source](#adding-a-server-as-an-aggregation-source)
  * [Viewing a collection of aggregation sources](#viewing-a-collection-of-aggregation-sources)
  * [Viewing an aggregation source](#viewing-an-aggregation-source)
  * [Updating an aggregation source](#updating-an-aggregation-source)
  * [Resetting servers](#resetting-servers)
  * [Changing the boot order of servers to default settings](#changing-the-boot-order-of-servers-to-default-settings)
  * [Deleting a resource from the inventory](#deleting-a-resource-from-the-inventory)
  * [Aggregates](#aggregates)
  * [Creating an aggregate](#creating-an-aggregate)
  * [Viewing a list of aggregates](#viewing-a-list-of-aggregates)
  * [Viewing information about a single aggregate](#viewing-information-about-a-single-aggregate)
  * [Deleting an aggregate](#deleting-an-aggregate)
  * [Adding elements to an aggregate](#adding-elements-to-an-aggregate)
  * [Resetting an aggregate of computer systems](#resetting-an-aggregate-of-computer-systems)
  * [Setting boot order of an aggregate to default settings](#setting-boot-order-of-an-aggregate-to-default-settings)
  * [Removing elements from an aggregate](#removing-elements-from-an-aggregate)
- [Resource inventory](#resource-inventory)
  * [Collection of computer systems](#collection-of-computer-systems)
  * [Single computer system](#single-computer-system)
  * [Memory collection](#memory-collection)
  * [Single memory](#single-memory)
  * [Memory domains](#memory-domains)
  * [BIOS](#bios)
  * [Network interfaces](#network-interfaces)
  * [Ethernet interfaces](#ethernet-interfaces)
  * [Single Ethernet interface](#single-ethernet-interface)
  * [PCIeDevice](#pciedevice)
  * [Storage](#storage)
  * [Storage subsystem](#storage-subsystem)
  * [Drives](#drives)
    + [Single drive](#single-drive)
  * [Volumes](#volumes)
    + [Collection of volumes](#collection-of-volumes)
    + [Single volume](#single-volume)
    + [Creating a volume](#creating-a-volume)
    + [Deleting a volume](#deleting-a-volume)
  * [SecureBoot](#secureboot)
  * [Processors](#processors)
  * [Single processor](#single-processor)
  * [Chassis](#chassis)
    + [Collection of chassis](#collection-of-chassis)
    + [Single chassis](#single-chassis)
    + [Thermal metrics](#thermal-metrics)
    + [Collection of network adapters](#collection-of-network-adapters)
    + [Single network adapter](#single-network-adapter)
    + [Power](#power)
    + [Creating a rack group](#creating-a-rack-group)
    + [Creating a rack](#creating-a-rack)
    + [Attaching chassis to a rack](#attaching-chassis-to-a-rack)
    + [Detaching chassis from a rack](#detaching-chassis-from-a-rack)
    + [Deleting a rack](#deleting-a-rack)
    + [Deleting a rack group](#deleting-a-rack-group)
  * [Searching the inventory](#searching-the-inventory)
    + [Request URI parameters](#request-uri-parameters)
- [Actions on a computer system](#actions-on-a-computer-system)
  * [Resetting a computer system](#resetting-a-computer-system)
  * [Changing the boot order of a computer system to default settings](#changing-the-boot-order-of-a-computer-system-to-default-settings)
  * [Changing BIOS settings](#changing-bios-settings)
  * [Changing the boot settings](#changing-the-boot-settings)
- [Managers](#managers)
  * [Collection of managers](#collection-of-managers)
  * [Single manager](#single-manager)
  * [VirtualMedia](#virtualmedia)
    + [Viewing the VirtualMedia collection](#viewing-the-virtualmedia-collection)
    + [Viewing a VirtualMedia Instance](#viewing-a-virtualmedia-instance)
    + [Inserting VirtualMedia](#inserting-virtualmedia)
    + [Ejecting VirtualMedia](#ejecting-virtualmedia)
- [Software and firmware inventory](#software-and-firmware-inventory)
  * [Viewing the update service root](#viewing-the-update-service-root)
  * [Viewing the firmware inventory](#viewing-the-firmware-inventory)
  * [Viewing a specific firmware resource](#viewing-a-specific-firmware-resource)
  * [Viewing the software inventory](#viewing-the-software-inventory)
  * [Viewing a specific software resource](#viewing-a-specific-software-resource)
  * [Actions](#actions)
    + [Simple update](#simple-update)
      - [Request parameters](#request-parameters)
    + [Start update](#start-update)
- [Host to fabric networking](#host-to-fabric-networking)
  * [Collection of fabrics](#collection-of-fabrics)
  * [Single fabric](#single-fabric)
  * [Collection of switches](#collection-of-switches)
  * [Single switch](#single-switch)
  * [Collection of ports](#collection-of-ports)
  * [Single port](#single-port)
  * [Collection of address pools](#collection-of-address-pools)
  * [Single address pool](#single-address-pool)
  * [Collection of endpoints](#collection-of-endpoints)
  * [Single endpoint](#single-endpoint)
  * [Collection of zones](#collection-of-zones)
  * [Single zone](#single-zone)
  * [Creating a zone-specific address pool](#creating-a-zone-specific-address-pool)
  * [Creating an address pool for zone of zones](#creating-an-address-pool-for-zone-of-zones)
  * [Adding a zone of zones](#adding-a-zone-of-zones)
  * [Adding an endpoint](#adding-an-endpoint)
  * [Creating a zone of endpoints](#creating-a-zone-of-endpoints)
  * [Updating a zone](#updating-a-zone)
  * [Deleting a zone](#deleting-a-zone)
  * [Deleting an endpoint](#deleting-an-endpoint)
  * [Deleting an address pool](#deleting-an-address-pool)
- [Tasks](#tasks)
  * [Viewing the task service root](#viewing-the-task-service-root)
  * [Viewing a collection of tasks](#viewing-a-collection-of-tasks)
  * [Viewing information about a specific task](#viewing-information-about-a-specific-task)
  * [Viewing a task monitor](#viewing-a-task-monitor)
  * [Deleting a task](#deleting-a-task)
- [Events](#events)
  * [Viewing the event service root](#viewing-the-event-service-root)
  * [Creating an event subscription](#creating-an-event-subscription)
    + [Creating event subscription with eventformat type “MetricReport”](#creating-event-subscription-with-eventformat-type---metricreport)
  * [Submitting a test event](#submitting-a-test-event)
  * [Event subscription use cases](#event-subscription-use-cases)
    + [Subscribing to resource addition notification](#subscribing-to-resource-addition-notification)
    + [Subscribing to resource removal notification](#subscribing-to-resource-removal-notification)
    + [Subscribing to task status notifications](#subscribing-to-task-status-notifications)
  * [Viewing a collection of event subscriptions](#viewing-a-collection-of-event-subscriptions)
  * [Viewing information about a specific event subscription](#viewing-information-about-a-specific-event-subscription)
  * [Deleting an event subscription](#deleting-an-event-subscription)
- [Message registries](#message-registries)
  * [Viewing a collection of registries](#viewing-a-collection-of-registries)
  * [Viewing a single registry](#viewing-a-single-registry)
  * [Viewing a file in a registry](#viewing-a-file-in-a-registry)
- [Redfish Telemetry](#redfish-telemetry-service)
  * [Viewing the telemetry service root](#viewing-the-telemetry-service-root)
  * [Collection of metric definitions](#collection-of-metric-definitions)
  * [Single metric definition](#single-metric-definition)
  * [Collection of Metric Report Definitions](#collection-of-metric-report-definitions)
  * [Single metric report definition](#single-metric-report-definition)
  * [Collection of metric reports](#collection-of-metric-reports)
  * [Single metric report](#single-metric-report)
  * [Collection of Triggers](#collection-of-triggers)
  * [Single Trigger](#single-trigger)
  * [Updating a trigger](#updating-a-trigger)

# Introduction 

 Welcome to Resource Aggregator for Open Distributed Infrastructure Management!

Resource Aggregator for Open Distributed Infrastructure Management (ODIMRA) is a modular, open framework for
simplified management and orchestration of distributed physical infrastructure. It provides a unified management platform for
converging multivendor hardware equipment. By exposing a standards-based programming interface, it enables easy and
secure management of wide range of multivendor IT infrastructure distributed across multiple data centers.

ODIMRA framework comprises the following two components.

- The resource aggregation function (the resource aggregator)

  The resource aggregation function is the single point of contact between the northbound clients and the
  southbound infrastructure. Its primary function is to build and maintain a central resource inventory. It exposes
  Redfish-compliant APIs to allow northbound infrastructure management systems to:
    - Get a unified view of the southbound compute, local storage, and Ethernet switch fabrics available in the
      resource inventory.
    - Gather crucial configuration information about southbound resources.
    - Manipulate groups of resources in a single action.
    - Listen to similar events from multiple southbound resources.
	
- One or more plugins

  The plugins abstract, translate, and expose southbound resource information to the resource aggregator through
  RESTful APIs. HPE Resource Aggregator for ODIM supports:

    - Generic Redfish plugin for ODIM (GRF): Generic Redfish plugin that can be used as a plugin for any Redfishcompliant
      device.
	- Plugin for unmanaged racks (URP): Plugin that acts as a resource manager for unmanaged racks. 
    - Integration of additional third-party plugins.  	

This guide provides reference information for the northbound APIs exposed by the resource aggregator. These APIs
are designed as per DMTF's [Redfish® Scalable Platforms API (Redfish) specification 1.11.1](https://www.dmtf.org/sites/default/files/standards/documents/DSP0266_1.11.1.pdf) and are fully Redfish-compliant.
Redfish® is an open industry standard specification, API, and schema. It specifies a RESTful interface and uses JSON and OData. The Redfish
standards are designed to deliver simple and secure environment for managing multivendor, converged, and hybrid IT infrastructure.


##  Resource Aggregator for ODIM logical architecture

Resource Aggregator for ODIM framework adopts a layered architecture and has many functional layers.


The following figure shows these functional layers of Resource Aggregator for ODIM deployed in a data center.

![ODIM_architecture](images/arch.png)

- **API layer**


This layer hosts a REST server which is open-source and secure. It learns about the southbound resources from
the plugin layer and exposes the corresponding Redfish data model payloads to the northbound clients. The
northbound clients communicate with this layer through a REST-based protocol that is fully compliant with
DMTF's Redfish® specifications (Schema 2020.3 and Specification 1.11.1).
The API layer sends user requests to the plugins through the aggregation, the event, and the fabric services.

- **Services layer**


The services layer is where all the services are hosted. This layer implements service logic for all use cases
through an extensible domain model (Redfish Data Model). Requests coming from the API layer and the
responses coming from the plugin layer are mapped to the actual end resources in this layer. It maintains the state
for event subscriptions, credentials, and tasks. It also hosts a message bus called the Plug-in Message Bus (PMB).

![Redfish_data_model](images/redfish_data_model.png)

- **Event message bus layer**


This layer hosts a message broker which acts as a communication channel between the upper layers and the
plugin layer. It supports common messaging architecture and real-time streaming. Resource Aggregator for
ODIM uses Kafka as the event message bus.
The services and the event message bus layers host Redis data store.

- **Plugin layer**


This layer connects the actual managed resources to the aggregator layers and is de-coupled from the upper
layers. A plugin abstracts vendor-specific access protocols to a common interface which the aggregator layers use
to communicate with the resources. It uses REST-based communication which is based on OpenAPI Specification
v3.0 to interact with the other layers. It collects events to be exposed to fault management systems and uses the
event message bus to publish events. The messaging mechanism is based on OpenMessaging Specification.
The plugin layer allows developers to create plugins on any tool set of their choice without enforcing any strict
language binding. To know how to develop plugins, refer to [Resource Aggregator for Open Distributed
Infrastructure Management Plugin Developer's Guide](https://github.com/ODIM-Project/ODIM/blob/development/plugin-redfish/README.md).


# API usage and access guidelines

To access the RESTful APIs exposed by the resource aggregator, you need an HTTPS-capable client, such as a web
browser with a REST Client plugin extension, or a Desktop REST Client application, or curl (a popular, free command-line
utility). You can also easily write simple REST clients for RESTful APIs using modern scripting languages.

<blockquote>
Tip: It is good to use a tool, such as curl or any Desktop REST Client application initially to perform requests. Later,
you will want to write your own scripting code to perform requests.
</blockquote>
This guide contains sample request and response payloads. For information on response payload parameters, refer to [Redfish® Scalable Platforms API (Redfish) schema 2020.3](https://www.dmtf.org/sites/default/files/standards/documents/DSP0268_2020.3.pdf).

> **IMPORTANT:**
The response codes and the JSON request and response parameters provided in this guide may vary for systems depending on the vendor, model, and firmware versions.

**HTTP headers**

Include the following HTTP headers:

- `Content-Type:application/json` for all RESTful API operations that include a request body in JSON
    format.
- Authentication header (BasicAuth or XAuthToken) for all RESTful API operations except for HTTP GET on the
    Redfish service root and HTTP POST on sessions.

**Base URL**

Use the following URL in all HTTP requests that you send to the resource aggregator:

`https://{odimra_host}:{port}/`

- {odimra_host} is the fully qualified domain name (FQDN) used for generating certificates while deploying the
    resource aggregator.
	
	>**NOTE:** Ensure that FQDN is provided in the `/etc/hosts` file or in the DNS server.


- {port} is the port where the services of the resource aggregator are running. The default port is 45000. If you
    have changed the default port in the `/etc/odimra_config/odimra_config.json` file, use that as the port.
>**NOTE**: To access the base URL using a REST client, replace `{odimra_host}` with the IP address of the system where the resource aggregator is installed. To use FQDN in place of `{odimra_host}`, add the Resource Aggregator for ODIM server certificate to the browser where the REST client is launched.


**curl usage**

The examples shown in this guide use curl to make HTTP requests.

[curl](https://curl.haxx.se) is a command-line tool which helps you get or send information through URLs using supported protocols. The resource aggregator for ODIM
supports HTTPS.

**curl command options (flags):**

- `--cacert` <file_path> includes a specified X.509 root certificate.
- `-H` passes on custom headers.
- `-X` specifies a custom request method. Use -X for HTTP PATCH, PUT, DELETE operations.
- `-d` posts data to a URI. Use -d for all HTTP operations that include a request body.
- `-i` returns HTTP response headers.
- `-v` fetches verbose.

For a complete list of curl flags, see information provided at [https://curl.haxx.se](https://curl.haxx.se).


>**IMPORTANT:** If you have set proxy configuration, set no_proxy using the following command before running a curl command.<br>
    ```
    export no_proxy="127.0.0.1,localhost,{odimra_host}"
     ```

**Including HTTP certificate**

Without CA certificate, curl fails to verify that HTTP connections are secure and curl commands may fail with the SSL
certificate problem. Provide the root CA certificate to curl for secure SSL communication.


1. If you are running curl commands on the server where the resource aggregator is deployed, provide the `rootCA.crt` file as shown in the curl command:
   ```
   curl -v --cacert {path}/rootCA.crt 'https://{odimra_host}:{port}/redfish/v1'
   ```
    {path} is where you have generated certificates during the deployment of the resource aggregator.

2. If you are running curl commands on a different server, perform the following steps to provide the rootCA.crt file.
    1. Navigate to `~/ODIM/build/cert_generator/certificates` on the server where the resource aggregator is deployed.
    2. Copy the `rootCA.crt` file.
    3. Log in to your server and paste the `rootCA.crt` file in a folder.
    4. Open the `/etc/hosts` file to edit.
    5. Scroll to the end of the file, add the following line, and save:
       `{odim_server_ipv4_address} {FQDN}`
    6. Check if curl is working using the curl command:
        ```
        curl -v --cacert {path}/rootCA.crt 'https://{odimra_host}:{port}/redfish/v1'
        ```

>**NOTE:** To avoid using the `--cacert` flag in every curl command, add `rootCA.crt` in the `ca-certificates.crt` file located in this path:<br> `/etc/ssl/certs/ca-certificates.crt`.

# Support for URL Encoding

The URL encoding mechanism translates the characters in the URLs to a representation that are universally accepted by all web browsers and servers. 

Resource Aggregator for ODIM supports all standard URL encoded characters for all the APIs. When Resource Aggregator for ODIM gets an encoded URL path, the non-ASCII characters in its path are internally translated and sent to the web browsers. In other words, if you replace a character in a URL with its standard encoding notation, Resource Aggregator for ODIM accepts the encoded notation, decodes it to the actual character acceptable by the web browser and sends responses.

**Example**: In the URL`/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1`, if you replace the colon character `:` with its encoding notation `%3A`, or the `/` character with %2F and send the request, Resource Aggregator for ODIM accepts the URL, decodes the encoded notation internally and sends an accurate response.

<blockquote>Tip: You can visit https://www.w3schools.com/tags/ref_urlencode.ASP or browse the Internet to view the standard ASCII Encoding Reference of the URL characters.</blockquote>

# List of supported APIs

Resource Aggregator for ODIM supports the following Redfish APIs:

|Redfish Service Root||
|-------|--------------------|
|/redfish|`GET`|
|/redfish/v1|`GET`|
|/redfish/v1/odata|`GET`|
|/redfish/v1/$metadata|`GET`|

|SessionService||
|-------|--------------------|
|/redfish/v1/SessionService|`GET`|
|/redfish/v1/SessionService/Sessions|`POST`, `GET`|
|redfish/v1/SessionService/Sessions/\{sessionId\}|`GET`, `DELETE`|

|AccountService||
|-------|--------------------|
|/redfish/v1/AccountService|`GET`|
|/redfish/v1/AccountService/Accounts|`POST`, `GET`|
|/redfish/v1/AccountService/Accounts/\{accountId\}|`GET`, `DELETE`, `PATCH`|
|/redfish/v1/AccountService/Roles|`POST`, `GET`|
|/redfish/v1/AccountService/Roles/\{roleId\}|`GET`, `DELETE`, `PATCH`|

|AggregationService||
|-------|--------------------|
|/redfish/v1/AggregationService|`GET`|
|/redfish/v1/AggregationService/AggregationSources<br> |`GET`, `POST`|
|/redfish/v1/AggregationService/AggregationSources/\{aggregationSourceId\}|`GET`, `PATCH`, `DELETE`|
|/redfish/v1/AggregationService/Actions/AggregationService.Reset|`POST`|
|/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder|`POST`|
|/redfish/v1/AggregationService/Aggregates|`GET`, `POST`|
|/redfish/v1/AggregationService/Aggregates/\{aggregateId\}|`GET`, `DELETE`|
|/redfish/v1/AggregationService/Aggregates/\{aggregateId\}/Actions/Aggregate.AddElements|`POST`|
|/redfish/v1/AggregationService/Aggregates/\{aggregateId\}/Aggregate.Reset|`POST`|
|/redfish/v1/AggregationService/Aggregates/\{aggregateId\}/Aggregate.SetDefaultBootOrder|`POST`|
|/redfish/v1/AggregationService/Aggregates/\{aggregateId\}/Actions/Aggregate.RemoveElements|`POST`|
|/redfish/v1/AggregationService/ConnectionMethods|GET|
|/redfish/v1/AggregationService/ConnectionMethods/\{connectionmethodsId\}|GET|

|Systems||
|-------|--------------------|
|/redfish/v1/Systems|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}|`GET`, `PATCH`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Memory|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Memory/\{memoryId\}|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/MemoryDomains|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/NetworkInterfaces|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/EthernetInterfaces|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/EthernetInterfaces/\{Id\}|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Bios|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/SecureBoot|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/PCIeDevices/{PCIeDeviceId}|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage/\{storageSubsystemId\}|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage/\{storageSubsystemId\}/Drives/\{driveId\}|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage/\{storageSubsystemId\}/Volumes|`GET` , `POST`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage/\{storageSubsystemId\}/Volumes/\{volumeId\}|`GET`, `DELETE`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Processors|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Processors/\{Id\}|`GET`|
|/redfish/v1/Systems?filter=\{searchKeys\}%20\{conditionKeys\}%20\{value/regEx\}|`GET`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Bios/Settings<br> |`GET`, `PATCH`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Actions/ComputerSystem.Reset|`POST`|
|/redfish/v1/Systems/\{ComputerSystemId\}/Actions/ComputerSystem.SetDefaultBootOrder|`POST`|

|Chassis||
|-------|--------------------|
|/redfish/v1/Chassis|`GET`, `POST`|
|/redfish/v1/Chassis/\{chassisId\}|`GET`, `PATCH`, `DELETE`|
|/redfish/v1/Chassis/\{chassisId\}/Thermal|`GET`|
|/redfish/v1/Chassis/\{ChassisId\}/Power|`GET`|
|/redfish/v1/Chassis/\{chassisId\}/NetworkAdapters|`GET`|
|/redfish/v1/Chassis/{ChassisId}/NetworkAdapters/{networkadapterId}|GET|`Login` |

|Managers||
|-------|--------------------|
|/redfish/v1/Managers|`GET`|
|/redfish/v1/Managers/\{managerId\}|`GET`|
|/redfish/v1/Managers/\{managerId\}/EthernetInterfaces|`GET`|
|/redfish/v1/Managers/\{managerId\}/HostInterfaces|`GET`|
|/redfish/v1/Managers/\{managerId\}/LogServices|`GET`|
|/redfish/v1/Managers/\{managerId\}/NetworkProtocol|`GET`|
|/redfish/v1/Managers/{ManagerId}/VirtualMedia|`GET`|
|/redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}| `GET`  |
|/redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}/Actions/VirtualMedia.InsertMedia|`POST`|
|/redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}/Actions/VirtualMedia.EjectMedia|`POST`|

|UpdateService||
|-------|--------------------|
|/redfish/v1/UpdateService|`GET`|
|/redfish/v1/UpdateService/FirmwareInventory|`GET`|
|/redfish/v1/UpdateService/FirmwareInventory/\{inventoryId\}|`GET`|
|/redfish/v1/UpdateService/SoftwareInventory|`GET`|
|/redfish/v1/UpdateService/SoftwareInventory/\{inventoryId\}|`GET`|
|/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate|`POST`|
|/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate|`POST`|

|EventService||
|-------|--------------------|
|/redfish/v1/EventService|`GET`|
|/redfish/v1/EventService/Subscriptions|`POST`, `GET`|
|/redfish/v1/EventService/Actions/EventService.SubmitTestEvent|`POST`|
|/redfish/v1/EventService/Subscriptions/\{subscriptionId\}|`GET`, `DELETE`|

|Fabrics||
|-------|--------------------|
|/redfish/v1/Fabrics|`GET`|
|/redfish/v1/Fabrics/\{fabricId\}|`GET`|
|/redfish/v1/Fabrics/\{fabricId\}/Switches|`GET`|
|/redfish/v1/Fabrics/\{fabricId\}/Switches/\{switchId\}|`GET`|
|/redfish/v1/Fabrics/\{fabricId\}/Switches/\{switchId\}/Ports<br> |`GET`|
|/redfish/v1/Fabrics/\{fabricId\} /Switches/\{switchId\}/Ports/\{portid\}<br> |`GET`|
|/redfish/v1/Fabrics/\{fabricId\}/Zones|`GET`, `POST`|
|/redfish/v1/Fabrics/\{fabricId\}/Zones/\{zoneId\}|`GET`, `PATCH`, `DELETE`|
|/redfish/v1/Fabrics/\{fabricId\}/AddressPools|`GET`, `POST`|
|/redfish/v1/Fabrics/\{fabricId\}/AddressPools/\{addresspoolid\}|`GET`, `DELETE`|
|/redfish/v1/Fabrics/\{fabricId\}/Endpoints|`GET`, `POST`|
|/redfish/v1/Fabrics/\{fabricId\}/Endpoints/\{endpointid\}|`GET`, `DELETE`|

|TaskService||
|-------|--------------------|
|/redfish/v1/TaskService|`GET`|
|/redfish/v1/TaskService/Tasks|`GET`|
|/redfish/v1/TaskService/Tasks/\{taskId\}|`GET`, `DELETE`|
| /redfish/v1/TaskService/Tasks/\{taskId\}/SubTasks |`GET`|
| /redfish/v1/TaskService/Tasks/\{taskId\}/SubTasks/ \{subTaskId\} |`GET`|

| TelemetryService                                             |                |
| ------------------------------------------------------------ | -------------- |
| /redfish/v1/TelemetryService                                 | `GET`          |
| /redfish/v1/TelemetryService/MetricDefinitions               | `GET`          |
| /redfish/v1/TelemetryService/MetricDefinitions/{MetricDefinitionId} | `GET`          |
| /redfish/v1//TelemetryService/MetricReportDefinitions        | `GET`,         |
| /redfish/v1/TelemetryService/MetricReportDefinitions/{MetricReportDefinitionId} | `GET`          |
| redfish/v1/TelemetryService/MetricReports                    | `GET`          |
| /redfish/v1/TelemetryService/MetricReports/{MetricReportId}  | `GET`          |
| /redfish/v1/TelemetryService/Triggers                        | `GET`          |
| /redfish/v1/TelemetryService/Triggers/{TriggerId}            | `GET`, `PATCH` |

|Task monitor||
|-------|--------------------|
|/taskmon/\{taskId\}|`GET`|

|Registries||
|-------|--------------------|
|/redfish/v1/Registries|`GET`|
|/redfish/v1/Registries/\{registryId\}|`GET`|
|/redfish/v1/registries/\{registryFileId\}|`GET`|


>**NOTE:**
`ComputerSystemId` is the unique identifier of a system specified by Resource Aggregator for ODIM.
It is represented as `<UUID:n>` in Resource Aggregator for ODIM. `<UUID:n>` is the
universally unique identifier of a system. Example: *ba0a6871-7bc4-5f7a-903d-67f3c205b08c:1*.


## Viewing the list of supported Redfish services

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1` |
|**Description** |This is the URI for the Redfish service root. Perform `GET` on this URI to fetch a list of available Redfish services.|
|**Returns** |All the available services in the service root.|
|**Response Code** |`200 OK` |
|**Authentication** |No|


>**curl command**


```curl
curl -i GET 'https://{odimra_host}:{port}/redfish/v1'
```

>**Sample response header**


```
Allow:GET
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=UTF-8
Link:</redfish/v1/SchemaStore/en/ServiceRoot.json/>; rel=describedby
Odata-Version:4.0
X-Frame-Options:sameorigin
Date":Fri,15 May 2020 13:55:53 GMT+5m 11s
Transfer-Encoding:chunked

```

>**Sample response body**


```
{
   "@odata.context": "/redfish/v1/$metadata#ServiceRoot.ServiceRoot",
   "@odata.id": "/redfish/v1/",
   "@odata.type": "#ServiceRoot.v1_5_0.ServiceRoot",
   "Id": "RootService",
   "Registries": {
      "@odata.id": "/redfish/v1/Registries"
   },
   "SessionService": {
      "@odata.id": "/redfish/v1/SessionService"
   },
   "AccountService": {
      "@odata.id": "/redfish/v1/AccountService"
   },
   "EventService": {
      "@odata.id": "/redfish/v1/EventService"
   },
   "Tasks": {
      "@odata.id": "/redfish/v1/TaskService"
   },
   "TelemetryService": {
      "@odata.id": "/redfish/v1/TelemetryService"
   },
   "AggregationService": {
      "@odata.id": "/redfish/v1/AggregationService"
   },
   "Systems": {
      "@odata.id": "/redfish/v1/Systems"
   },
   "Chassis": {
      "@odata.id": "/redfish/v1/Chassis"
   },
   "Fabrics": {
      "@odata.id": "/redfish/v1/Fabrics"
   },
   "Managers": {
      "@odata.id": "/redfish/v1/Managers"
   },
   "Links": {
      "Sessions": {
         "@odata.id": "/redfish/v1/SessionService/Sessions"
      }
   },
   "Name": "Root Service",
   "Oem":{

   },
   "RedfishVersion": "1.11.1",
   "UUID": "a64fc187-e0e9-4f68-82a8-67a616b84b1d"
}
```



#  HTTP request methods responses and status codes

Following are the Redfish-defined HTTP methods that you can use to implement various actions:

|HTTP Request Methods|Description|
|--------------------|-----------|
|`GET` \[Read Requests\]|Use this method to request a representation of a specified resource \(single resource or collection\).|
|`PATCH` \[Update\]|Use this method to apply partial modifications to a resource.|
|`POST` \[Create\] \[Actions\]|Use this method to create a resource. Submit this request to the resource collection to which you want to add the new resource. You can also use this method to initiate operations on a resource or a collection of resources.|
|`PUT` \[Replace\]|Use this method to replace the property values of a resource completely. It is used to both create and update the state of a resource.|
|`DELETE` \[Delete\]|Use this method to delete a resource.|

Resource Aggregator for ODIM supports the following responses:

|Responses|Description|
|---------|-----------|
|Metadata responses|Describes the resources and types exposed by the service to generic clients.|
|Resource responses|Response in JSON format for an individual resource.|
|Resource collection responses|Response in JSON format for a collection of resources.|
|Error responses|If there is an HTTP error, a high-level JSON response is provided with additional information.|

Following are the HTTP status codes with their descriptions:

| Success code<br> |Description|
|-----------------|-----------|
|200 OK|Successful completion of the request with representation in the body.|
|201 Created|A new resource is successfully created with the `Location` header set to well-defined URI for the newly created resource. The response body may include the representation of the newly created resource.|
|202 Accepted|The request has been accepted for processing but not processed. The `Location` header is set to URI of a task monitor that can be queried later for the status of the operation.|
|204 No Content|The request succeeded, but no content is being returned in the body of the response.|

| Error code<br> |Description|
|-----------------|-----------|
|301 Moved Permanently|The requested resource resides in a different URI given by the `Location` headers.|
|400 Bad Request|The request could not be performed due to missing or invalid information. An extended error message is returned in the response body.|
|401 Unauthorized|Missing or invalid authentication credentials included with a request.|
|403 Forbidden|The server recognizes the credentials to be not having the necessary authorization to perform the operation.|
|404 Not Found|The request specified the URI of a nonexisting resource.|
|405 Method Not Allowed|When the HTTP verb specified in the request \(GET, PATCH, DELETE, and so on\) is not supported for a particular request URI. The response includes `Allow` header that lists the supported methods.|
|409 Conflict|A creation or an update cannot be completed because it would conflict with the current state of the resources supported by the platform.|
|500 Internal Server Error|When the server encounters an unexpected condition that prevents it from fulfilling the request.|
|501 Not Implemented|When the server has not implemented the method for the resource.|
|503 Service Unavailable|When the server is unable to service the request due to temporary overloading or maintenance.|


>**NOTE:**
This guide provides success codes (200, 201, 202, 204) for all referenced API operations. For failed operations, refer to the error codes listed in this section.



# Authentication and authorization

##  Authentication methods for Redfish APIs

 To keep HTTP connections secure, Resource Aggregator for ODIM verifies credentials of HTTP requests. If you perform an unauthenticated HTTP operation on resources except for the following, you will receive an HTTP `401 unauthorized` error.

|Resource|URI|Description|
|--------|---|-----------|
|The Redfish service root|`GET` `/redfish` |The URI for the Redfish service root. It returns the version of Redfish services.|
|List of Redfish services|`GET``/redfish/v1` |It returns a list of available services.|
|$metadata|`GET``/redfish/v1/$metadata` |The Redfish metadata document.|
|OData|`GET``/redfish/v1/odata` |The Redfish OData service document.|
| The `Sessions` resource<br> |`POST``/redfish/v1/SessionService/Sessions` |Creates a Redfish login session.|

To authenticate requests with the Redfish services, implement any one of the following authentication methods:


-   **HTTP BASIC authentication \(BasicAuth\)** 

     To implement HTTP BASIC authentication:

     1. Generate a base64 encoded string of `{valid_username_of_odim_userAccount}:{valid_password_of_odim_userAccount}` using the following command:

         ```
        echo -n '{username}:{password}' | base64 -w0
        ```

        Initially, use the username and the password of the default administrator account. Later, you can create additional [user accounts](#user-accounts) and use their details to implement authentication.

     2. Provide the base64 encoded string in an HTTP `Authorization:Basic` header as shown in the curl command:

         ```
         curl -i --cacert {path}/rootCA.crt GET\
         -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
         'https://{odimra_host}:{port}/redfish/v1/AccountService'
         ```



-   **Redfish session login authentication \(XAuthToken\)** 

      To implement Redfish session login authentication, create a Redfish login [session](#sessions) and obtain an authentication token through session management interface.
    
      Every session that is created has an authentication token called `X-AUTH-TOKEN` associated with it. An `X-AUTH-TOKEN` is returned in the response header from session creation.
    
      To authenticate subsequent requests, provide this token in the `X-AUTH-TOKEN` request header.
    ```
    curl -i --cacert {path}/rootCA.crt GET \
    -H "X-Auth-Token:{X-Auth-Token}" \
     'https://{odimra_host}:{port}/redfish/v1/AccountService'
    ```

      An `X-AUTH-TOKEN` is valid and a session is open only for 30 minutes, unless you continue to send requests to a Redfish service using this token. An idle session is automatically terminated after the time-out interval.

## Role-based authorization

In Resource Aggregator for ODIM, the roles and privileges control which users have what access to resources. If you perform an HTTP operation on a resource without necessary privileges, you will receive an HTTP `403 Forbidden` error.

**Roles**

A role represents a set of operations that a user is allowed to perform. It is determined by a defined set of privileges. You can assign a role to a user at the time of user account creation.

With Resource Aggregator for ODIM, there are two kinds of defined roles:

-   **Redfish predefined roles** 

    Redfish predefined roles come with a predefined set of privileges. These predefined privileges cannot be removed or modified. You may assign additional OEM \(custom\) privileges. Following are the Redfish predefined roles that are available in Resource Aggregator for ODIM by default:

    -   `Administrator` 

    -   `Operator` 

    -   `ReadOnly` 

-   **User-defined roles** 

    User-defined roles are the custom roles that you can create and assign to a user. The privileges of a user-defined role are configurable—you can choose a privilege or a set of privileges to assign to this role at the time of role creation.


<blockquote>
NOTE:

-   Redfish predefined roles cannot be modified.

-   Before assigning a user-defined role to a user at the time of user account creation, ensure that it is created first.
</blockquote>

**Privileges**

A privilege is a permission to perform an operation or a set of operations within a defined management domain.

The following Redfish-specified privileges can be assigned to any user in Resource Aggregator for ODIM:

-    `ConfigureComponents` :
 Users with this privilege can configure components managed by the Redfish services in Resource   Aggregator for ODIM.
This privilege is required to create, update, and delete a resource or a collection of resources exposed by Redfish APIs using HTTP `POST, PATCH, DELETE` operations.

 -    `ConfigureManager` :
 Users with this privilege can configure manager resources.

 -    `ConfigureComponents` :
 Users with this privilege can configure components managed by the services.

 -    `ConfigureSelf` :
 Users with this privilege can change the password for their account.

 -    `ConfigureUsers` :
 Users with this privilege can configure users and their accounts. This privilege is assigned to an `Administrator`.
This privilege is required to create, update, and delete user accounts using HTTP `POST, PATCH, DELETE` operations.

 -    `Login` :
 Users with this privilege can log in to the service and read the resources.
This privilege is required to view any resource or a collection of resources exposed by Redfish APIs using HTTP `GET` operation.

**Mapping of privileges to roles in Resource Aggregator for ODIM**

|Roles|Assigned privileges|
|-----|-------------------|
|Administrator \(Redfish predefined\)|    `Login` <br>   `ConfigureManager` <br>   `ConfigureUsers` <br>    `ConfigureComponents` <br>   `ConfigureSelf` <br>|
|Operator \(Redfish predefined\)|    `Login` <br>   `ConfigureComponents` <br>   `ConfigureSelf` <br>|
|ReadOnly \(Redfish predefined\)|    `Login` <br>   `ConfigureSelf` <br>|


>**NOTE:**
Resource Aggregator for ODIM has a default user account that has all the privileges of an administrator role.


# Sessions

A session represents a window of user's login with a Redfish service and contains details about the user and the user activity. You can run multiple sessions simultaneously.

Resource Aggregator for ODIM offers Redfish `SessionService` interface for creating and managing sessions. It exposes APIs to achieve the following:

-   Fetching the `SessionService` root.

-   Creating a session.

-   Listing active sessions.

-   Deleting a session.


**Supported APIs**

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/SessionService|GET|`Login` |
|/redfish/v1/SessionService/Sessions|POST, GET|`Login`|
|redfish/v1/SessionService/Sessions/\{sessionId\}|GET, DELETE|`Login`, `ConfigureManager`, `ConfigureSelf` |

>**NOTE:**
Before accessing these endpoints, ensure that the user has the required privileges. If you access these endpoints without necessary privileges, you will receive an HTTP `403 Forbidden` error.



## Viewing the session service root

|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/SessionService` |
|**Description** |This endpoint retrieves JSON schema representing the Redfish `SessionService` root.|
|**Returns** |The properties for the Redfish `SessionService` itself and links to the actual list of sessions.|
|**Response Code** |`200 OK` |
|**Authentication** |No|


>**curl command**

```
curl -i GET \
              'https://{odimra_host}:{port}/redfish/v1/SessionService'

```


>**Sample response body**

```
{
   "@odata.type":"#SessionService.v1_1_6.SessionService",
   "@odata.id":"/redfish/v1/SessionService",
   "Id":"Sessions",
   "Name":"Session Service",
   "Status":{
      "State":"Enabled",
      "Health":"OK"
   },
   "ServiceEnabled":true,
   "SessionTimeout":30,
   "Sessions":{
      "@odata.id":"/redfish/v1/SessionService/Sessions"
   }
}
```

##  Creating a session

|||
|---------|---------------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/SessionService/Sessions` |
|**Description** |This operation creates a session to implement authentication. Creating a session allows you to create an `X-AUTH-TOKEN` which is then used to authenticate with other services.<br>**NOTE:**<br>It is a good practice to note down the following:<br><ul><li>The session authentication token returned in the `X-AUTH-TOKEN` header.</li><li>The session Id returned in the `Location` header and the JSON response body.</li></ul><br> You will need the session authentication token to authenticate subsequent requests to the Redfish services and the session Id to log out later.|
|**Returns** |<ul><li> An `X-AUTH-TOKEN` header containing session authentication token.</li><li>`Location` header that contains a link to the newly created session instance.</li><li>The session Id and a message in the JSON response body saying that the session is created.</li></ul> |
|**Response Code** |`201 Created` |
|**Authentication** |No|



>**curl command**

```
curl -i POST \
   -H "Content-Type:application/json" \
   -d \
'{
"UserName": "{username}",
"Password": "{password}"
}' \
 'https://{odimra_host}:{port}/redfish/v1/SessionService/Sessions'

```


>**Sample request body**

```
{
"UserName": "abc",
"Password": "abc123"
}
```



**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|UserName|String \(required\)|Username of the user account for this session. For the first time, use the username of the default administrator account \(admin\). Later, when you create other user accounts, you can use the credentials of those accounts to create a session.<br>**NOTE:**<br> This user must have `Login` privilege.|
|Password|String \(required\)<br> |Password of the user account for this session. For the first time, use the password of the default administrator account. Later, when you create other user accounts, you can use the credentials of those accounts to create a session.<br> |



>**Sample response header**


```
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Link:</redfish/v1/SessionService/Sessions/2d2e8ebc-4e7c-433a-bfd6-74dc420886d0/>; rel=self
Location:{odimra_host}:{port}/redfish/v1/SessionService/Sessions/2d2e8ebc-4e7c-433a-bfd6-74dc420886d0
Odata-Version:4.0
X-Auth-Token:15d0f639-f394-4be7-a8ef-ef9d1df07288
X-Frame-Options:sameorigin
Date:Fri,15 May 2020 14:08:55 GMT+5m 11s
Transfer-Encoding:chunked
```

>**Sample response body**


```
{
	"@odata.type": "#SessionService.v1_1_6.SessionService",
	"@odata.id": "/redfish/v1/SessionService/Sessions/1a547199-0dd3-42de-9b24-1b801d4a1e63",
	"Id": "1a547199-0dd3-42de-9b24-1b801d4a1e63",
	"Name": "Session Service",
	"Message": "The resource has been created successfully",
	"MessageId": "Base.1.6.1.Created",
	"Severity": "OK",
	"UserName": "abc"
}
```




## Listing sessions

|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/SessionService/Sessions` |
|**Description** |This operation lists user sessions.<br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can see a list of all user sessions.<br> Users with `ConfigureSelf` privilege can see only the sessions created by them.|
|**Returns** |Links to user sessions.|
|**Response Code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
               -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
              'https://{odimra_host}:{port}/redfish/v1/SessionService/Sessions'

```



## Viewing information about a single session

|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/SessionService/Sessions/{sessionId}` |
|**Description** |This operation retrieves information about a specific user session.<br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can view information about any user session.<br><br> Users with `ConfigureSelf` privilege can view information about only the sessions created by them.<br>|
|**Returns** |JSON schema representing this session.|
|**Response Code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
                -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/SessionService/Sessions/{sessionId}'

```


>**Sample response body**

```
{
   "@odata.type":"#Session.v1_2_1.Session",
   "@odata.id":"/redfish/v1/SessionService/Sessions/4ee42139-22db-4e2a-97e4-020013248768",
   "Id":"4ee42139-22db-4e2a-97e4-020013248768",
   "Name":"User Session",
   "UserName":"admin"
}
```



## Deleting a session

|||
|---------|---------------|
|**Method** | `DELETE` |
|**URI** |`/redfish/v1/SessionService/Sessions/{sessionId}` |
|**Description** |This operation terminates a specific Redfish session when the user logs out.<br>**NOTE:**<br> Users having the `ConfigureSelf` and `ConfigureComponents` privileges are allowed to delete only those sessions that they created.<br><br> Only a user with all the Redfish-defined privileges \(Redfish-defined `Administrator` role\) is authorized to delete any user session.<br> |
|**Response Code** |`204 No Content` |
|**Authentication** |Yes|

>**curl command**


```
curl -i -X DELETE \
               -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
              'https://{odimra_host}:{port}/redfish/v1/SessionService/Sessions/{sessionId}'

```






#  User roles and privileges

Resource Aggregator for ODIM allows you to create and manage user roles through Redfish APIs.

**Supported APIs**:

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/AccountService|GET|`Login` |
|/redfish/v1/AccountService/Roles|GET, POST|`Login`, `ConfigureManager` |
|/redfish/v1/AccountService/Roles/\{RoleId\}|GET, PATCH, DELETE|`Login`, `ConfigureManager` |


>**NOTE:**
Before accessing these endpoints, ensure that the user has the required privileges. If you access these endpoints without necessary privileges, you will receive an HTTP `403 Forbidden` error.


## Viewing the account service root

|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/AccountService` |
|**Description** |This endpoint fetches JSON schema representing the Redfish `AccountService` root.|
|**Returns** |The properties common to all user accounts and links to the collections of manager accounts and roles.|
|**Response Code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/AccountService'


```

>**Sample response header**

```
Allow:GET
Cache-Control:no-cache
Connection:Keep-alive
Content-Type:application/json; charset=utf-8
Link:</redfish/v1/SchemaStore/en/AccountService.json>; rel=describedby
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Fri,15 May 2020 14:32:09 GMT+5m 12s
Transfer-Encoding:chunked
```


>**Sample response body**

```
{
   "@odata.type":"#AccountService.v1_6_0.AccountService",
   "@odata.id":"/redfish/v1/AccountService",
   "@odata.context":"/redfish/v1/$metadata#AccountService.AccountService",
   "Id":"AccountService",
   "Name":"Account Service",
   "Status":{
      "State":"Enabled",
      "Health":"OK"
   },
   "ServiceEnabled":true,
   "AuthFailureLoggingThreshold":0,
   "MinPasswordLength":12,
   "AccountLockoutThreshold":0,
   "AccountLockoutDuration":0,
   "AccountLockoutCounterResetAfter":0,
   "Accounts":{
      "@odata.id":"/redfish/v1/AccountService/Accounts"
   },
   "Roles":{
      "@odata.id":"/redfish/v1/AccountService/Roles"
   }
}
```



## Creating a role

|||
|---------|---------------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/AccountService/Roles` |
|**Description** |This operation creates a role other than Redfish predefined roles. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can perform this operation.|
|**Returns** |JSON schema representing the newly created role.|
|**Response code** |`201 Created` |
|**Authentication** |Yes|

>**curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{ 
   "RoleId":"CLIENT11",
   "AssignedPrivileges":[ 
      "Login",
      "ConfigureUsers",
      "ConfigureSelf"
   ],
   "OemPrivileges":null 
}' \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Roles'


```



>**Sample request body**

```
{ 
   "RoleId":"CLIENT11",
   "AssignedPrivileges":[ 
      "Login",
      "ConfigureUsers",
      "ConfigureSelf"
   ],
   "OemPrivileges":null
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Id|String \(required, read-only\)<br> |Name for this role. <br>**NOTE:**<br> Id cannot be modified later.|
|AssignedPrivileges|Array \(string \(enum\)\) \(required\)<br> |The Redfish privileges that this role includes. Possible values are:<br>  `ConfigureManager` <br>   `ConfigureSelf` <br>   `ConfigureUsers` <br>   `Login` <br>   `ConfigureComponents` <br>|
|OemPrivileges|Array \(string\) \(required\)<br> |The OEM privileges that this role includes. If you do not want to specify any OEM privileges, use `null` or `[]` as value.|


>**Sample response body**

```
{
   "@odata.type":"#Role.v1_2_4.Role",
   "@odata.id":"/redfish/v1/AccountService/Roles/CLIENT11",
   "Id":"CLIENT11",
   "Name":"User Role",
   "Message":"The resource has been created successfully.",
   "MessageId":"ResourceEvent.1.0.2.ResourceCreated",
   "Severity":"OK",
   "IsPredefined":false,
   "AssignedPrivileges":[
      "Login",
      "ConfigureUsers",
      "ConfigureSelf"
   ],
   "OemPrivileges":null
}
```

## Listing roles

|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/AccountService/Roles` |
|**Description** |This operation lists available user roles. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can perform this operation.|
|**Returns** |Links to user role resources.|
|**Response Code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Roles'

```


>**Sample response body**

```
{ 
   "@odata.type":"#RoleCollection.RoleCollection",
   "@odata.id":"/redfish/v1/AccountService/Roles",
   "Name":"Roles Collection",
   "Members@odata.count":5,
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/AccountService/Roles/Administrator"
      },
      { 
         "@odata.id":"/redfish/v1/AccountService/Roles/Operator"
      },
      { 
         "@odata.id":"/redfish/v1/AccountService/Roles/ReadOnly"
      },
      { 
         "@odata.id":"/redfish/v1/AccountService/Roles/CLIENT13"
      },
      { 
         "@odata.id":"/redfish/v1/AccountService/Roles/CLIENT11"
      }
      
   ]
}
```







## Viewing information about a role


|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/AccountService/Roles/{RoleId}` |
|**Description** |This operation fetches information about a specific user role. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can perform this operation.|
|**Returns** |JSON schema representing this role. The schema has the details such as - Id, name, description, assigned privileges, OEM privileges.|
|**Response Code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Roles/{RoleId}'


```

 >**Sample response body**

```
{
   "@odata.type":"#Role.v1_2_4.Role",
   "@odata.id":"/redfish/v1/AccountService/Roles/CLIENT11",
   "Id":"CLIENT11",
   "Name":"User Role",
   "IsPredefined":false,
   "AssignedPrivileges":[
      "Login",
      "ConfigureUsers",
      "ConfigureSelf"
   ],
   "OemPrivileges":null
}
```





## Updating a role

|||
|---------|---------------|
|**Method** | `PATCH` |
|**URI** |`/redfish/v1/AccountService/Roles/{RoleId}` |
|**Description** |This operation updates privileges of a specific user role - assigned privileges \(Redfish predefined\) and OEM privileges. Id of a role cannot be modified.<br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can perform this operation.|
|**Returns** |JSON schema representing the updated role.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**

```
 curl -i -X PATCH \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
  "AssignedPrivileges": [{Set_Of_Privileges_to_update}],
  "OemPrivileges": []
}' \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Roles/{RoleId}'

```


>**Sample request body**

```
{ 
   "AssignedPrivileges":[ 
      "Login",
      "ConfigureManager",
      "ConfigureUsers"
   ],
   "OemPrivileges": []
}
```

>**Sample response body**

```
{
   "RoleId":"CLIENT11",
   "IsPredefined":false,
   "AssignedPrivileges":[
      "Login",
      "ConfigureManager",
      "ConfigureUsers"
   ],
   "OemPrivileges":null
}
```




## Deleting a role

|||
|---------|---------------|
|**Method** | `DELETE` |
|**URI** |`/redfish/v1/AccountService/Roles/{RoleId}` |
|**Description** |This operation deletes a specific user role. If you attempt to delete a role that is already assigned to a user account, you will receive an HTTP `403 Forbidden` error.<br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can perform this operation.|
|**Response Code** |`204 No Content` |
|**Authentication** |Yes|

>**curl command**

```
curl -i -X DELETE \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Roles/{RoleId}'

```





#  User accounts

Resource Aggregator for ODIM allows users to have accounts to configure their actions and restrictions.

Resource Aggregator for ODIM has an administrator user account by default.

Create other user accounts by defining a username, a password, and a role for each account. The username and the password are used to authenticate with the Redfish services \(using `BasicAuth` or `XAuthToken`\).

Resource Aggregator for ODIM exposes Redfish `AccountsService` APIs to create and manage user accounts. Use these endpoints to perform the following operations:

-   Creating, modifying, and deleting account details.

-   Fetching account details.


**Supported APIs**:

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/AccountService/Accounts|POST, GET|`Login`, `ConfigureUsers` |
|/redfish/v1/AccountService/Accounts/\{accountId\}|GET, DELETE, PATCH|`Login`, `ConfigureUsers`, `ConfigureSelf` |


>**NOTE:**
Before accessing these endpoints, ensure that the user has the required privileges. If you access these endpoints without necessary privileges, you will receive an HTTP `403 Forbidden` error.





## Creating a user account

|||
|-------|--------------------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/AccountService/Accounts` |
|**Description** |This operation creates a user account. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can create other user account.|
|**Returns** |<ul><li>`Location` header that contains a link to the newly created account.</li><li>JSON schema representing the newly created account.</li></ul> |
|**Response Code** |`201 Created` |
|**Authentication** |Yes|

>**curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{"Username":"{username}","Password":"{password}","RoleId":"{roleId}"}
' \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Accounts'


```



>**Sample request body**

```
{ 
   "Username":"monitor32",
   "Password":"Abc1vent2020!",
   "RoleId":"CLIENT11"
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Username|String \(required\)<br> |User name for the user account.|
|Password|String \(required\)<br> |Password for the user account. Before creating a password, see "Password Requirements" .|
|RoleId|String \(required\)<br> |The role for this account. To know more about roles, see [User roles and privileges](#role-based-authorization). Ensure that the `roleId` you want to assign to this user account exists. To check the existing roles, see [Listing Roles](#listing-roles). If you attempt to assign an unavailable role, you will receive an HTTP `400 Bad Request` error.|

### Password requirements

-   Your password must not be same as your username.

-   Your password must be at least 12 characters long and at most 16 characters long.

-   Your password must contain at least one uppercase letter \(A-Z\), one lowercase letter \(a-z\), one digit \(0-9\), and one special character \(~!@\#$%^&\*-+\_|\(\)\{\}:;<\>,.?/\).


>**Sample response header**

```
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Link:</redfish/v1/AccountService/Accounts/monitor32/>; rel=describedby
Location:/redfish/v1/AccountService/Accounts/monitor32/
Odata-Version:4.0
X-Frame-Options:"sameorigin
Date":Fri,15 May 2020 14:36:14 GMT+5m 11s
Transfer-Encoding:chunked
```

>**Sample response body**

```
{
   "@odata.type":"#ManagerAccount.v1_4_0.ManagerAccount",
   "@odata.id":"/redfish/v1/AccountService/Accounts/monitor32",
   "@odata.context":"/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
   "Id":"monitor32",
   "Name":"Account Service",
   "Message":"The resource has been created successfully",
   "MessageId":"Base.1.6.1.Created",
   "Severity":"OK",
   "UserName":"monitor32",
   "RoleId":"CLIENT11",
   "AccountTypes":[
      "Redfish"
   ],
   "Password":null,
   "Links":{
      "Role":{
         "@odata.id":"/redfish/v1/AccountService/Roles/CLIENT11/"
      }
   }
}
```

##  Listing user accounts

|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/AccountService/Accounts` |
|**Description** |This operation retrieves a list of user accounts. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can view a list of user accounts.|
|**Returns** |Links to user accounts.|
|**Response Code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Accounts'


```





##  Viewing the account details

|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/AccountService/Accounts/{accountId}` |
|**Description** |This operation fetches information about a specific user account. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can view information about a user account.|
|**Returns** |JSON schema representing this user account.|
|**Response Code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Accounts/{accountId}'


```

>**Sample response body**

```
{
   "@odata.type":"#ManagerAccount.v1_4_0.ManagerAccount",
   "@odata.id":"/redfish/v1/AccountService/Accounts/monitor32",
   "@odata.context":"/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
   "Id":"monitor32",
   "Name":"Account Service",
   "UserName":"monitor32",
   "RoleId":"CLIENT11",
   "AccountTypes":[
      "Redfish"
   ],
   "Password":null,
   "Links":{
      "Role":{
         "@odata.id":"/redfish/v1/AccountService/Roles/CLIENT11/"
      }
   }
}
```



## Updating a user account

|||
|---------|---------------|
|**Method** | `PATCH` |
|**URI** |`/redfish/v1/AccountService/Accounts/{accountId}` |
|**Description** |This operation updates user account details \(`username`, `password`, and `RoleId`\). To modify account details, add them in the request payload \(as shown in the sample request body\) and perform `PATCH` on the mentioned URI. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can modify other user accounts. Users with `ConfigureSelf` privilege can modify only their own accounts.|
|**Returns** |<ul><li>`Location` header that contains a link to the updated account.</li><li>JSON schema representing the modified account.</li></ul>|
|**Response Code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i -X PATCH \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{ 
   "Password":{new_password}",
   "RoleId":"CLIENT11"
}
' \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Accounts/{accountId}'

```


>**Sample request body**

```
{ 
   "Password":"Testing)9-_?{}",
   "RoleId":"CLIENT11"
}
```

>**Sample response header**

```
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Link:</redfish/v1/AccountService/Accounts/monitor32/>; rel=describedby
Location:/redfish/v1/AccountService/Accounts/monitor32/
Odata-Version:4.0
X-Frame-Options:"sameorigin
Date":Fri,15 May 2020 14:36:14 GMT+5m 11s
Transfer-Encoding:chunked

```

>**Sample response body**

```
{
   "@odata.type":"#ManagerAccount.v1_4_0.ManagerAccount",
   "@odata.id":"/redfish/v1/AccountService/Accounts/monitor32",
   "@odata.context":"/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
   "Id":"monitor32",
   "Name":"Account Service",
   "Message":"The account was successfully modified.",
   "MessageId":"Base.1.6.1.AccountModified",
   "Severity":"OK",
   "UserName":"monitor32",
   "RoleId":"CLIENT11",
   "AccountTypes":[
      "Redfish"
   ],
   "Password":null,
   "Links":{
      "Role":{
         "@odata.id":"/redfish/v1/AccountService/Roles/CLIENT11/"
      }
   }
}
```

## Deleting a user account

|||
|---------|---------------|
|**Method** | `DELETE` |
|**URI** |`/redfish/v1/AccountService/Accounts/{accountId}` |
|**Description** |This operation deletes a user account. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can delete a user account.|
|**Response Code** |`204 No Content` |
|**Authentication** |Yes|

>**curl command**

```
curl  -i -X DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Accounts/{accountId}'

```





#  Resource aggregation and management


The resource aggregator allows you to add southbound infrastructure to its database, create resource collections, and perform actions in combination on these collections. It exposes Redfish aggregation service endpoints to achieve the following:

-   Adding a resource and building its inventory.

-   Resetting one or more resources.

-   Changing the boot path of one or more resources to default settings.

-   Removing a resource from the inventory which is no longer managed.



All aggregation actions are performed as [tasks](#tasks) in Resource Aggregator for ODIM. The actions performed on a group of resources \(resetting or changing the boot order to default settings\) are carried out as a set of subtasks.

**Supported endpoints**

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/AggregationService|GET|`Login` |
| /redfish/v1/AggregationService/AggregationSources<br> |GET, POST|`Login`, `ConfigureManager` |
|/redfish/v1/AggregationService/AggregationSources/\{aggregationSourceId\}|GET, PATCH, DELETE|`Login`, `ConfigureManager` |
|/redfish/v1/AggregationService/Actions/AggregationService.Reset|POST|`ConfigureManager`, `ConfigureComponents` |
|/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder|POST|`ConfigureManager`, `ConfigureComponents` |
|/redfish/v1/AggregationService/Aggregates|GET, POST|`Login`, `ConfigureComponents`, `ConfigureManager` |
|/redfish/v1/AggregationService/Aggregates/\{aggregateId\}|GET, DELETE|`Login`, `ConfigureComponents`, `ConfigureManager` |
|/redfish/v1/AggregationService/Aggregates/\{aggregateId\}/Actions/Aggregate.AddElements|POST|`ConfigureComponents`, `ConfigureManager` |
|/redfish/v1/AggregationService/Aggregates/\{aggregateId\}/Aggregate.Reset|POST|`ConfigureComponents`, `ConfigureManager` |
|/redfish/v1/AggregationService/Aggregates/\{aggregateId\}/Aggregate.SetDefaultBootOrder|POST|`ConfigureComponents`, `ConfigureManager` |
|/redfish/v1/AggregationService/Aggregates/\{aggregateId\}/Actions/Aggregate.RemoveElements|POST|`ConfigureComponents`, `ConfigureManager` |
|/redfish/v1/AggregationService/ConnectionMethods|GET|`Login`|
|/redfish/v1/AggregationService/ConnectionMethods/\{connectionmethodsId\}|GET|`Login`|

>**Note:**
Before accessing these endpoints, ensure that the user has the required privileges. If you access these endpoints without necessary privileges, you will receive an HTTP `403 Forbidden` error.




## Viewing the aggregation service root
|||
|-----|-------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/AggregationService` |
|<strong>Description</strong> |This endpoint retrieves JSON schema representing the aggregation service root.|
|<strong>Returns</strong> |Properties for the service and a list of actions you can perform using this service.|
|<strong>Response Code</strong> |On success, `200 OK` |
|<strong>Authentication</strong> |Yes|

 

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odim_host}:{port}/redfish/v1/AggregationService'


```

>**Sample response header** 

```
Allow:GET
Cache-Control:no-cache
Connection:Keep-alive
Content-Type:application/json; charset=utf-8
Date:Sun,17 May 2020 14:26:49 GMT+5m 14s
Link:</redfish/v1/SchemaStore/en/AggregationService.json>; rel=describedby
Odata-Version:4.0
X-Frame-Options:sameorigin
Transfer-Encoding":chunked

```

>**Sample response body** 

```
{
   "@odata.context":"/redfish/v1/$metadata#AggregationService.AggregationService",
   "Id":"AggregationService",
   "@odata.id":"/redfish/v1/AggregationService",
   "@odata.type":"#AggregationService.v1_0_0.AggregationService",
   "Name":"AggregationService",
   "Description":"AggregationService",
   "Actions":{
      "#AggregationService.Reset":{
         "target":"/redfish/v1/AggregationService/Actions/AggregationService.Reset/",
         "@Redfish.ActionInfo":"/redfish/v1/AggregationService/ResetActionInfo"
      },
      "#AggregationService.SetDefaultBootOrder":{
         "target":"/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder/",
         "@Redfish.ActionInfo":"/redfish/v1/AggregationService/SetDefaultBootOrderActionInfo"
      }
   },
   "Aggregates":{
      "@odata.id":"/redfish/v1/AggregationService/Aggregates"
   },
   "AggregationSources":{
      "@odata.id":"/redfish/v1/AggregationService/AggregationSources"
   },
   "ConnectionMethods":{
      "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods"
   },
   "ServiceEnabled":true,
   "Status":{
      "Health":"OK",
      "HealthRollup":"OK",
      "State":"Enabled"
   }
}
```


## Connection methods

Connection methods indicate protocols, providers, or other methods that are used to communicate with an endpoint. 
The `ConnectionMethod` schema describes these connection methods for the Redfish aggregation service. 

###  Viewing a collection of connection methods


|||
|--------|---------|
|**Method**| `GET` |
|**URI** |`/redfish/v1/AggregationService/ConnectionMethods` |
|**Description** |This operation lists all connection methods associated with the Redfish aggregation service.|
|**Returns** |A list of links to all the available connection method resources.|
|**Response Code** |On success, `200 Ok` |
|**Authentication** |Yes|



>**curl command** 

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/ConnectionMethods'


```

>**Sample response body**

```
{
   ​   "@odata.type":"#ConnectionMethodCollection.ConnectionMethodCollection",
   ​   "@odata.id":"/redfish/v1/AggregationService/ConnectionMethods",
   ​   "@odata.context":"/redfish/v1/$metadata#ConnectionMethodCollection.ConnectionMethodCollection",
   ​   "Name":"Connection Methods",
   ​   "Members@odata.count":3,
   ​   "Members":[
      ​      {
         ​         "@odata.id":"/redfish/v1/AggregationService/ConnectionMethods/c27575d2-052d-4ce9-8be1-978cab002a0f"         ​
      },
      ​      {
         ​         "@odata.id":"/redfish/v1/AggregationService/ConnectionMethods/aa166b6b-a367-40ba-ac2e-402f9a0c818f"         ​
      },
      ​      {
         ​         "@odata.id":"/redfish/v1/AggregationService/ConnectionMethods/7cb9fc3b-8b75-45da-8aad-5ff595968b71"         ​
      }      ​
   ]   ​
}
```

### Viewing a connection method

|||
|--------|---------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/AggregationService/ConnectionMethods/ {connectionmethodsId}` |
|**Description** |This operation retrieves information about a specific connection method.|
|**Returns** |JSON schema representing this connection method.|
|**Response Code** |On success, `200 Ok` |
|**Authentication**|Yes|

>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/ConnectionMethods/{connectionmethodsId}'


```

>**Sample response body**

```
{
   ​   "@odata.type":"#ConnectionMethod.v1_0_0.ConnectionMethod",
   ​   "@odata.id":"/redfish/v1/AggregationService/ConnectionMethods/c27575d2-052d-4ce9-8be1-978cab002a0f",
   ​   "@odata.context":"/redfish/v1/$metadata#ConnectionMethod.v1_0_0.ConnectionMethod",
   ​   "Id":"c27575d2-052d-4ce9-8be1-978cab002a0f",
   ​   "Name":"Connection Method",
   ​   "ConnectionMethodType":"Redfish",
   ​   "ConnectionMethodVariant":"Compute:BasicAuth:GRF_v1.0.0",
   ​   "Links":{
      ​      "AggregationSources":[
         {
            "@odata.id":"/redfish/v1/AggregationService/AggregationSources/839c212d-9ab2-4868-8767-1bdcc0ce862c"
         },
         {
            "@odata.id":"/redfish/v1/AggregationService/AggregationSources/3536bb46-a023-4e3a-ac1a-7528cc18b660"
         }
      ]      ​
   }   ​
}​
```

>**Connection method properties**

|Parameter|Type|Description|
|---------|----|-----------|
|ConnectionMethodType|String| The type of this connection method.<br> For possible property values, see "Connection method types" table.<br> |
|ConnectionMethodVariant|String|The variant of connection method. For more information, see [Connection method variants](#connection-method-variants).|
|Links \{|Object|Links to other resources that are related to this connection method.|
|AggregationSources \[ \{<br> @odata.id<br> \} \]<br> |Array|An array of links to the `AggregationSources` resources that use this connection method.|

 >**Connection method types**

|String|Description|
|------|-----------|
| IPMI15<br> | IPMI 1.5 connection method.<br> |
| IPMI20<br> | IPMI 2.0 connection method.<br> |
| NETCONF<br> | Network Configuration Protocol.<br> |
| OEM<br> | OEM connection method.<br> |
| Redfish<br> | Redfish connection method.<br> |
| SNMP<br> | Simple Network Management Protocol.<br> |

#### Connection method variants

A connection method variant provides details about a plugin and is displayed in the following format:

*`PluginType:PrefferedAuthType:PluginID_Firmwareversion`*. 

It consists of the following parameters:

- **PluginType:**
   The string that represents the type of the plugin.<br>Possible values: Compute, Storage, and Fabric. 
- **PrefferedAuthType:**   
   Preferred authentication method to connect to the plugin - BasicAuth or XAuthToken.  
- **PluginID\_Firmwareversion:**
   The id of the plugin along with the version of the firmware. To know the plugin Ids for all the supported plugins, see "Mapping of plugins and plugin Ids" table.<br>
   Supported values: GRF\_v1.0.0 and URP\_v1.0.0.<br>  


Examples:
1. `Compute:BasicAuth:GRF_v1.0.0`
2. `Compute:BasicAuth:URP_v1.0.0`


>**Mapping of plugins and plugin Ids**

|Plugin Id|Plugin name|
|---------|-----------|
|GRF|Generic Redfish Plugin|
|URP|Unmanaged Rack Plugin|


##  Adding a plugin as an aggregation source

| | |
|-------|------|
|<strong>Method</strong> | `POST` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/AggregationSources` |
|<strong>Description</strong> | This operation creates an aggregation source for a plugin and adds it in the inventory. It is performed in the background as a Redfish task.|
|<strong>Returns</strong> |<ul><li>`Location` URI of the task monitor associated with this operation in the response header. See `Location` URI highlighted in bold in "Sample response header \(HTTP 202 status\)".</li><li>Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See "Sample response body \(HTTP 202 status\)".</li><li>On successful completion:<ul><li>The aggregation source Id, the IP address, the username, and other details of the added plugin in the JSON response body.</li><li> A link \(having the aggregation source Id\) to the added plugin in the `Location` header. See `Location` URI highlighted in bold in "Sample response header \(HTTP 201 status\)".</li></ul></li></ul>  |
|<strong>Response Code</strong> |`202 Accepted` On success, `201 Created`|
|<strong>Authentication</strong> |Yes|

**Usage information**

Perform HTTP POST on the mentioned URI with a request body specifying a connection method to use for adding the plugin. To know about connection methods, see [Connection methods](#connection-methods).
				
A Redfish task will be created and you will receive a link to the [task monitor](#viewing-a-task-monitor) associated with it.
To know the progress of this operation, perform HTTP `GET` on the task monitor returned in the response header (until the task is complete).
		

After the plugin is successfully added as an aggregation source, it will also be available as a manager resource at:

`/redfish/v1/Managers`.

 


**NOTE:**

Only a user with `ConfigureComponents` privilege can add a plugin. If you perform this operation without necessary privileges, you will receive an HTTP `403 Forbidden` error.


>**curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{"HostName":"{plugin_host}:{port}",
  "UserName":"{plugin_userName}",
  "Password":"{plugin_password}", 
  "Links":{
     "ConnectionMethod": {
         "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/{ConnectionMethodId}"
      }
   }
}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/AggregationSources'


```



>**Sample request body for adding the GRF plugin**

```
{
   "HostName":"{plugin_host}:45001",
   "UserName":"admin",
   "Password":"GRFPlug!n12$4",
   "Links":{
      "ConnectionMethod": {
         "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/d172e66c-b4a8-437c-981b-1c07ddfeacaa"
      }
   }
}
```

>**Sample request body for adding URP**

```
{
   "HostName":"{plugin_host}:45003",
   "UserName":"admin",
   "Password":"Od!m12$4",
   "Links":{
      "ConnectionMethod": {
         "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/a171e66c-b4a8-137f-981b-1c07ddfeacbb"
      }
   }
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|HostName|String \(required\)<br> |FQDN of the resource aggregator server and port of a system where the plugin is installed. The default port for the Generic Redfish Plugin is `45001`.<br>The default port for the URP is `45003`.<br> If you are using a different port, ensure that the port is greater than `45000`.<br> IMPORTANT: If you have set the `VerifyPeer` property to false in the plugin `config.json` file \(/etc/plugin\_config/config.json\), you can use IP address of the system where the plugin is installed as `HostName`.<br>|
|UserName|String \(required\)<br> |The plugin username.|
|Password|String \(required\)<br> |The plugin password.|
|Links\{|Object \(required\)<br> |Links to other resources that are related to this resource.|
|ConnectionMethod|Array (required)|Links to the connection method that are used to communicate with this endpoint: `/redfish/v1/AggregationService/AggregationSources`. To know which connection method to use, do the following:<ul><li>Perform HTTP `GET` on: `/redfish/v1/AggregationService/ConnectionMethods`.<br>You will receive a list of  links to available connection methods.</li><li>Perform HTTP `GET` on each link. Check the value of the `ConnectionMethodVariant` property in the JSON response. Choose a connection method having the details of the plugin of your choice.<br>For example, the `ConnectionMethodVariant` property for the GRF plugin displays the following value:<br>`Compute:BasicAuth:GRF_v1.0.0` <br>For more information, see the "connection method properties" table in [Viewing a connection method](#viewing-a-connection-method)</li></ul>|

>**Sample response header \(HTTP 202 status\)**

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/task85de4003-8757-4c7d-942f-55eaf7d6812a
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes

```

>**Sample response header \(HTTP 201 status\)**

```
"cache-control":"no-cache
connection":"keep-alive
content-type":application/json; charset=utf-8
date:"Wed",02 Sep 2020 06:50:43 GMT+7m 2s
link:/v1/AggregationService/AggregationSources/be626e78-7a8a-4b99-afd2-b8ed45ef3d5a:1/>; rel=describedby
location:/redfish/v1/AggregationService/AggregationSources/be626e78-7a8a-4b99-afd2-b8ed45ef3d5a:1
odata-version:4.0
transfer-encoding:"chunked
x-frame-options":"sameorigin"
```

>**Sample response body \(HTTP 202 status\)**

```
{
   "@odata.type":"#Task.v1_4_2.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4003-8757-4c7d-942f-55eaf7d6812a",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task85de4003-8757-4c7d-942f-55eaf7d6812a",
   "Name":"Task task85de4003-8757-4c7d-942f-55eaf7d6812a",
   "Message":"The task with id task85de4003-8757-4c7d-942f-55eaf7d6812a has started.",
   "MessageId":"TaskEvent.1.0.1.TaskStarted",
   "MessageArgs":[
      "task85de4003-8757-4c7d-942f-55eaf7d6812a"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```



>**Sample response body \(HTTP 201 status\)**

```
{
   "@odata.type":"#AggregationSource.v1_0_0.AggregationSource",
   "@odata.id":"/redfish/v1/AggregationService/AggregationSources/be626e78-7a8a-4b99-afd2-b8ed45ef3d5a",
   "@odata.context":"/redfish/v1/$metadata#AggregationSource.AggregationSource",
   "Id":"be626e78-7a8a-4b99-afd2-b8ed45ef3d5a",
   "Name":"Aggregation Source",
   "HostName":"{plugin_host}:45001",
   "UserName":"admin",
   "Links":{
      "ConnectionMethod": {
         "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/d172e66c-b4a8-437c-981b-1c07ddfeacaa"
      }
   }
} 
```









## Adding a server as an aggregation source

| | |
|-------------|---------------------|
|<strong>Method</strong> | `POST` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/AggregationSources` |
|<strong>Description</strong> | This operation creates an aggregation source for a Base Management Controller \(BMC\), discovers information, and performs a detailed inventory of it.<br> The `AggregationSource` schema provides information about a BMC such as the IP address, the username, the password, and more.<br> This operation is performed in the background as a Redfish task.<br> |
|<strong>Returns</strong> |<ul><li>`Location` URI of the task monitor associated with this operation in the response header. See `Location` URI highlighted in bold in "Sample response header \(HTTP 202 status\)".</li><li>Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See "Sample response body \(HTTP 202 status\)".</li><li>On successful completion:<ul><li>The aggregation source Id, the IP address, the username, and other details of the added BMC in the JSON response body.</li><li>A link \(having the aggregation source Id\) to the added BMC in the `Location` header. See `Location` URI highlighted in bold in "Sample response header \(HTTP 201 status\)".</li></ul></li></ul>|
|<strong>Response Code</strong> |On success, `202 Accepted`<br> On successful completion of the task, `201 Created` <br> |
|<strong>Authentication</strong> |Yes|

**Usage information**

Perform HTTP POST on the mentioned URI with a request body specifying a connection method to use for adding the BMC. To know about connection methods, see [Connection methods](#connection-methods).
				
A Redfish task will be created and you will receive a link to the [task monitor](#viewing-a-task-monitor) associated with it.

To know the progress of this operation, perform HTTP `GET` on the task monitor returned in the response header (until the task is complete).
When the task is successfully complete, you will receive aggregation source Id of the added BMC. Save it as it is required to identify it in the resource inventory later.

After the server is successfully added as an aggregation source, it will also be available as a computer system resource at `/redfish/v1/Systems/` and a manager resource at `/redfish/v1/Managers/`.

To view the list of links to computer system resources, perform HTTP `GET` on `/redfish/v1/Systems/`. Each link contains `ComputerSystemId` of a specific BMC. For more information, see [Collection of computer systems](#collection-of-computer-systems).

 `ComputerSystemId` is unique information about the BMC specified by Resource Aggregator for ODIM. It is represented as `<UUID:n>`, where `UUID` is the aggregation source Id of the BMC. Save it as it is required to perform subsequent actions such as `delete, reset`, and `setdefaultbootorder` on this BMC.


**NOTE:**

Only a user with `ConfigureComponents` privilege can add a server. If you perform this operation without necessary privileges, you will receive an HTTP `403 Forbidden` error.


>**curl command**

```
curl -i -X POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{ 
    "HostName": "{BMC_address}", 
    "UserName": "{BMC_UserName}", 
    "Password": "{BMC_Password}", 
    "Links":{     
        "ConnectionMethod": {
         "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/{ConnectionMethodId}"
      }
}
}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/AggregationSources'


```



>**Sample request body**

```
{
   "HostName":"10.24.0.4",
   "UserName":"admin",
   "Password":"{BMC_password}",
   "Links":{
      "ConnectionMethod": {
         "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/d172e66c-b4a8-437c-981b-1c07ddfeacaa"
      }
   }
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|HostName|String \(required\)<br> |A valid IP address or hostname of a baseboard management controller \(BMC\).|
|UserName|String \(required\)<br> |The username of the BMC administrator account.|
|Password|String \(required\)<br> |The password of the BMC administrator account.|
|Links \{|Object \(required\)<br> |Links to other resources that are related to this resource.|
|ConnectionMethod|Array (required)|Links to the connection methods that are used to communicate with this endpoint: `/redfish/v1/AggregationService/AggregationSources`. To know which connection method to use, do the following:<ul><li>Perform HTTP `GET` on: `/redfish/v1/AggregationService/ConnectionMethods`.<br>You will receive a list of  links to available connection methods.</li><li>Perform HTTP `GET` on each link. Check the value of the `ConnectionMethodVariant` property in the JSON response.</li><li>The `ConnectionMethodVariant` property displays the details of a plugin. Choose a connection method having the details of the plugin of your choice.<br> Example: For GRF plugin, the `ConnectionMethodVariant` property displays the following value:<br>`Compute:BasicAuth:GRF:1.0.0`</li></ul>|

>**Sample response header \(HTTP 202 status\)**

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/task4aac9e1e-df58-4fff-b781-52373fcb5699
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes

```

>**Sample response header \(HTTP 201 status\)**

```
"cache-control":"no-cache
connection":"keep-alive
content-type":application/json; charset=utf-8
date:"Wed",02 Sep 2020 06:50:43 GMT+7m 2s
link:/v1/AggregationService/AggregationSources/0102a4b5-03db-40be-ad39-71e3c9f8280e/>; rel=describedby
location:/redfish/v1/AggregationService/AggregationSources/0102a4b5-03db-40be-ad39-71e3c9f8280e
odata-version:4.0
transfer-encoding:"chunked
x-frame-options":"sameorigin"
```

>**Sample response body \(HTTP 202 status\)**

```
{
   "@odata.type":"#Task.v1_4_2.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "Name":"Task task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "Message":"The task with id task4aac9e1e-df58-4fff-b781-52373fcb5699 has started.",
   "MessageId":"TaskEvent.1.0.1.TaskStarted",
   "MessageArgs":[
      "task4aac9e1e-df58-4fff-b781-52373fcb5699"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```



>** Sample response body \(HTTP 201 status\)**
```
 {
   "@odata.type":"#AggregationSource.v1_0_0.AggregationSource",
   "@odata.id":"/redfish/v1/AggregationService/AggregationSources/26562c7b-060b-4fd8-977e-94b1a535f3fb",
   "@odata.context":"/redfish/v1/$metadata#AggregationSource.AggregationSource",
   "Id":"26562c7b-060b-4fd8-977e-94b1a535f3fb",
   "Name":"Aggregation Source",
   "HostName":"10.24.0.4",
   "UserName":"admin",
    "Links":{
      "ConnectionMethod": {
         "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/d172e66c-b4a8-437c-981b-1c07ddfeacaa"
      }
   }
}
```


## Viewing a collection of aggregation sources

| | |
|-------|-------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/AggregationSources` |
|<strong>Description</strong> |This operation lists all aggregation sources available in Resource Aggregator for ODIM.|
|<strong>Returns</strong> |A list of links to all the available aggregation sources.|
|<strong>Response Code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/AggregationSources'


```

>**Sample response body**

```
{
   "@odata.type":"#AggregationSourceCollection.AggregationSourceCollection",
   "@odata.id":"/redfish/v1/AggregationService/AggregationSource",
   "@odata.context":"/redfish/v1/$metadata#AggregationSourceCollection.AggregationSourceCollection",
   "Name":"Aggregation Source",
   "Members@odata.count":2,
   "Members":[
      {
         "@odata.id":"/redfish/v1/AggregationService/AggregationSources/839c212d-9ab2-4868-8767-1bdcc0ce862c"
      },
      {
         "@odata.id":"/redfish/v1/AggregationService/AggregationSources/3536bb46-a023-4e3a-ac1a-7528cc18b660:1"
      }
   ]   
}
```



## Viewing an aggregation source

| | |
|--------|------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/AggregationSources/{AggregationSourceId}` |
|<strong>Description</strong> |This action retrieves information about a specific aggregation source.|
|<strong>Returns</strong> |JSON schema representing this aggregation source.|
|<strong>Response Code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/AggregationSources/{AggregationSourceId}'


```

>**Sample response body**

```
{
   "@odata.type":"#AggregationSource.v1_0_0.AggregationSource",
   "@odata.id":"/redfish/v1/AggregationService/AggregationSources/839c212d-9ab2-4868-8767-1bdcc0ce862c",
   "@odata.context":"/redfish/v1/$metadata#AggregationSource.AggregationSource",
   "Id":"839c212d-9ab2-4868-8767-1bdcc0ce862c",
   "Name":"Aggregation Source",
   "HostName":"10.24.0.4",
   "UserName":"admin",
   "Links":{
      "ConnectionMethod": {
         "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/d172e66c-b4a8-437c-981b-1c07ddfeacaa"
      }
   }
}
```

## Updating an aggregation source

| | |
|------|------|
|<strong>Method</strong> | `PATCH` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/AggregationSources/{AggregationSourceId}` |
|<strong>Description</strong> |This operation updates the details such as the username, the password, and the IP address or hostname of a specific BMC in the resource aggregator inventory. When the username, the password, or the IP address \(or hostname\) of a BMC is changed, you can update those changes in the resource aggregator as well using this operation.<br> |
|<strong>Returns</strong> |Updated JSON schema of this aggregation source.|
|<strong>Response Code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i PATCH \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{

  "HostName": "10.24.0.6",
  "UserName": "admin",
  "Password": "admin1234"

}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/AggregationSources/{AggregationSourceId}'


```

>**Sample request body**

```
{

  "HostName": "10.24.0.4",
  "UserName": "admin",
  "Password": "admin1234"

}
```

>**Sample response body**

```
{
   "@odata.type":"#AggregationSource.v1_0_0.AggregationSource",
   "@odata.id":"/redfish/v1/AggregationService/AggregationSources/839c212d-9ab2-4868-8767-1bdcc0ce862c:1",
   "@odata.context":"/redfish/v1/$metadata#AggregationSource.AggregationSource",
   "Id":"839c212d-9ab2-4868-8767-1bdcc0ce862c:1",
   "Name":"Aggregation Source",
   "HostName":"10.24.0.4",
   "UserName":"admin",
   "Links":{
      "ConnectionMethod": {
         "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/a172e66c-b4a8-437c-981b-1c07ddfeacab"
      }
   } 
}
```








## Resetting servers

|| |
|--------|--------------------|
|<strong>Method</strong> | `POST` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/Actions/AggregationService.Reset` |
|<strong>Description</strong> |This action shuts down, powers up, and restarts one or more servers. It is performed in the background as a Redfish task and is further divided into subtasks to reset each server individually.<br> |
|<strong>Returns</strong> |- `Location` URI of the task monitor associated with this operation \(task\) in the response header. See `Location` URI highlighted in bold in "Sample response header \(HTTP 202 status\)".<br><br>-   Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See the task URI and the task Id highlighted in bold in "Sample response body \(HTTP 202 status\)".<br><blockquote>IMPORTANT: Note down the task Id. If the task completes with an error, it is required to know which subtask has failed. To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`.</blockquote><br>-  On successful completion of the reset operation, a message in the response body, saying that the reset operation is completed successfully. See "Sample response body \(HTTP 200 status\)".|
|<strong>Response code</strong> |On success, `202 Accepted`<br> On successful completion of the task, `200 OK`|
|<strong>Authentication</strong> |Yes|

**Usage information**

To know the progress of this action, perform HTTP `GET` on the [task monitor](#viewing-a-task-monitor) returned in the response header \(until the task is complete\).

To get the list of subtask URIs, perform HTTP `GET` on the task URI returned in the JSON response body. See "Sample response body \(HTTP 202 status\)". The JSON response body of each subtask contains a link to the task monitor associated with it. To know the progress of the reset operation \(subtask\) on a specific server, perform HTTP `GET` on the task monitor associated with the respective subtask. See the link to the task monitor highlighted in bold in "Sample response body \(subtask\)".

You can perform reset on a group of servers by specifying multiple target URIs in the request.


**NOTE:**

Only a user with `ConfigureComponents` privilege can reset servers. If you perform this action without necessary privileges, you will receive an HTTP `403 Forbidden` error.


>**curl command**

```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{
   "BatchSize":2,
   "DelayBetweenBatchesInSeconds":1,
   "ResetType":"ForceRestart",
   "TargetURIs":[
      "/redfish/v1/Systems/{ComputerSystemId1}",
      "/redfish/v1/Systems/{ComputerSystemId2}"
   ]
}' \
 'https://\{odim\_host\}:\{port\}/redfish/v1/AggregationService/Actions/AggregationService.Reset'


```

>**Sample request body**

```
{
   "BatchSize":2,
   "DelayBetweenBatchesInSeconds":1,
   "ResetType":"ForceRestart",
   "TargetURIs":[
      "/redfish/v1/Systems/97d08f36-17f5-5918-8082-f5156618f58d:1",
      "/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b9:1"
   ]
}

```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|BatchSize|Integer \(optional\)<br> |The number of elements to be reset at a time in each batch.|
|DelayBetweenBatchesInSeconds|Integer \(seconds\) \(optional\)<br> |The delay among the batches of elements being reset.|
|ResetType|String \(required\)<br> |The type of reset to be performed. For possible values, see "Reset type". If the value is not supported by the target server machine, you will receive an HTTP `400 Bad Request` error.|
|TargetURIs|Array \(required\)<br> |The URI of the target for `Reset`. Example: `"/redfish/v1/Systems/{ComputerSystemId}"` |

**Reset type**

|String|Description|
|------|-----------|
|ForceOff|Turn off the unit immediately \(non-graceful shutdown\).|
|ForceRestart|Perform an immediate \(non-graceful\) shutdown, followed by a restart.|
|GracefulRestart|Perform a graceful shutdown followed by a restart of the system.|
|GracefulShutdown|Perform a graceful shutdown. Graceful shutdown involves shutdown of the operating system followed by the power off of the physical server.|
|Nmi|Generate a Diagnostic Interrupt \(usually an NMI on x86 systems\) to cease normal operations, perform diagnostic actions and typically halt the system.|
|On|Turn on the unit.|
|PowerCycle|Perform a power cycle of the unit.|
|PushPowerButton|Simulate the pressing of the physical power button on this unit.|

>**Sample response header** \(HTTP 202 status\)

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/task85de4103-8757-4c7d-942f-55eaf7d6412a
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes

```

>**Sample response body** \(HTTP 202 status\)

```
{
   "@odata.type":"#Task.v1_4_2.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4103-8757-4c7d-942f-55eaf7d6412a",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task85de4103-8757-4c7d-942f-55eaf7d6412a",
   "Name":"Task task85de4103-8757-4c7d-942f-55eaf7d6412a",
   "Message":"The task with id task85de4103-8757-4c7d-942f-55eaf7d6412a has started.",
   "MessageId":"TaskEvent.1.0.1.TaskStarted",
   "MessageArgs":[
      "task85de4103-8757-4c7d-942f-55eaf7d6412a"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```

>**Sample response body** \(subtask\)

```
{
   "@odata.type":"#SubTask.v1_4_2.SubTask",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4003-8757-4c7d-942f-55eaf7d6412a/SubTasks/task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "@odata.context":"/redfish/v1/$metadata#SubTask.SubTask",
   "Id":"task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "Name":"Task task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "Message":"Successfully Completed Request",
   "MessageId":"Base.1.6.1.Success",
   "Severity":"OK",
   "Members@odata.count":0,
   "Members":null,
   "TaskState":"Completed",
   "StartTime":"2020-05-13T13:33:59.917329733Z",
   "EndTime":"2020-05-13T13:34:00.320539988Z",
   "TaskStatus":"OK",
   "SubTasks":"",
   "TaskMonitor":"/taskmon/task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "PercentComplete":100,
   "Payload":{
      "HttpHeaders":null,
      "HttpOperation":"POST",
      "JsonBody":"",
      "TargetUri":"/redfish/v1/Systems/97d08f36-17f5-5918-8082-f5156618f58d:1"
   },
   "Messages":null
}
```

>**Sample response body** \(HTTP 200 status\)

```
{ 
   "error":{ 
      "code":"Base.1.6.1.Success",
      "message":"Request completed successfully"
   }
}
```



## Changing the boot order of servers to default settings

| | |
|-----------|------------|
|<strong>Method</strong> | `POST` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder` |
|<strong>Description</strong> |This action changes the boot order of one or more servers to default settings. This operation is performed in the background as a Redfish task and is further divided into subtasks to change the boot order of each server individually.<br> |
|<strong>Returns</strong> |- `Location` URI of the task monitor associated with this operation in the response header. See `Location` URI highlighted in bold in "Sample response header \(HTTP 202 status\)".<br><br>-  Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See the task URI and the task Id highlighted in bold in "Sample response body \(HTTP 202 status\)".<blockquote><br>IMPORTANT:<br>Note down the task Id. If the task completes with an error, it is required to know which subtask has failed. To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`.</blockquote><br>- On successful completion of this operation, a message in the response body, saying that the operation is completed successfully. See "Sample response body \(HTTP 200 status\)".<br>|
|<strong>Response code</strong> |`202 Accepted` On successful completion, `200 OK` <br> |
|<strong>Authentication</strong> |Yes|

**Usage information**

To know the progress of this action, perform HTTP `GET` on the [task monitor](#viewing-a-task-monitor) returned in the response header \(until the task is complete\).

To get the list of subtask URIs, perform HTTP `GET` on the task URI returned in the JSON response body. See "Sample response body \(HTTP 202 status\)". The JSON response body of each subtask contains a link to the task monitor associated with it. To know the progress of `SetDefaultBootOrder` action \(subtask\) on a specific server, perform HTTP `GET` on the task monitor associated with the respective subtask. See the link to the task monitor highlighted in bold in "Sample response body \(subtask\)".

You can perform `setDefaultBootOrder` action on a group of servers by specifying multiple server URIs in the request.


**NOTE:**

Only a user with `ConfigureComponents` privilege can change the boot order of one or more servers to default settings. If you perform this action without necessary privileges, you will receive an HTTP `403 Forbidden` error.


>**curl command**

```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{
   "Systems":[
      {
         "@odata.id":"/redfish/v1/Systems/{ComputerSystemId1}"
      },
      {
         "@odata.id":"/redfish/v1/Systems/{ComputerSystemid2}"
      }
   ]
}' \
 'https://\{odim\_host\}:\{port}/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder'


```

>**Sample request body**

```
{
   "Systems":[
      {
         "@odata.id":"/redfish/v1/Systems/97d08f36-17f5-5918-8082-f5156618f58d:1"
      },
      {
         "@odata.id":"/redfish/v1/Systems/76632110-1c75-5a86-9cc2-471325983653:1"
      }
   ]
}

```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Systems|Array \(required\)<br> |Target servers for `SetDefaultBootOrder`.|

>**Sample response header** \(HTTP 202 status\)

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/task85de4003-8057-4c7d-942f-55eaf7d6412a
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes

```

>**Sample response body** \(HTTP 202 status\)

```
{
   "@odata.type":"#Task.v1_4_2.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4003-8057-4c7d-942f-55eaf7d6412a",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task85de4003-8057-4c7d-942f-55eaf7d6412a",
   "Name":"Task task85de4003-8057-4c7d-942f-55eaf7d6412a",
   "Message":"The task with id task80de4003-8757-4c7d-942f-55eaf7d6412a has started.",
   "MessageId":"TaskEvent.1.0.1.TaskStarted",
   "MessageArgs":[
      "task80de4003-8757-4c7d-942f-55eaf7d6412a"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```

>**Sample response body** \(subtask\)

```
{
   "@odata.type":"#SubTask.v1_4_2.SubTask",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4003-8757-4c7d-942f-55eaf7d6412a/SubTasks/task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "@odata.context":"/redfish/v1/$metadata#SubTask.SubTask",
   "Id":"task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "Name":"Task task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "Message":"Successfully Completed Request",
   "MessageId":"Base.1.6.1.Success",
   "Severity":"OK",
   "Members@odata.count":0,
   "Members":null,
   "TaskState":"Completed",
   "StartTime":"2020-05-13T13:33:59.917329733Z",
   "EndTime":"2020-05-13T13:34:00.320539988Z",
   "TaskStatus":"OK",
   "SubTasks":"",
   "TaskMonitor":"/taskmon/task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "PercentComplete":100,
   "Payload":{
      "HttpHeaders":null,
      "HttpOperation":"POST",
      "JsonBody":"",
      "TargetUri":"/redfish/v1/Systems/97d08f36-17f5-5918-8082-f5156618f58d:1"
   },
   "Messages":null
}
```

>**Sample response body** \(HTTP 200 status\)

```
{ 
   "error":{ 
      "code":"Base.1.6.1.Success",
      "message":"Request completed successfully"
   }
}
```




## Deleting a resource from the inventory

| | |
|--------|--------|
|<strong>Method</strong> | `DELETE` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/AggregationSources/{AggregationSourceId}` |
|<strong>Description</strong> |This operation removes a specific aggregation source \(plugin, BMC, or any manager\) from the inventory. Deleting an aggregation source also deletes all event subscriptions associated with the BMC. This operation is performed in the background as a Redfish task.<br> |
|<strong>Returns</strong> |- `Location` URI of the task monitor associated with this operation in the response header. See `Location` URI highlighted in bold in "Sample response header \(HTTP 202 status\)".<br>-   Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See "Sample response body \(HTTP 202 status\)".<br>|
|<strong>Response Code</strong> |`202 Accepted` On successful completion, `204 No Content` <br> |
|<strong>Authentication</strong> |Yes|

**Usage information**

To know the progress of this action, perform `GET` on the [task monitor](#viewing-a-task-monitor) returned in the response header \(until the task is complete\).


**NOTE:**

Only a user with `ConfigureComponents` privilege can delete a server. If you perform this action without necessary privileges, you will receive an HTTP `403 Forbidden` error.


>**curl command**

```
curl -i DELETE \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/AggregationSources/{AggregationSourceId}'


```

>**Sample response header** \(HTTP 202 status\)

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/task85de4003-8757-2c7d-942f-55eaf7d6412a
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes

```

>**Sample response body** \(HTTP 202 status\)

```
{
   "@odata.type":"#Task.v1_4_2.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4003-8757-2c7d-942f-55eaf7d6412a",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task85de4003-8757-2c7d-942f-55eaf7d6412a",
   "Name":"Task task85de4003-8757-2c7d-942f-55eaf7d6412a",
   "Message":"The task with id task85de4003-8757-2c7d-942f-55eaf7d6412a has started.",
   "MessageId":"TaskEvent.1.0.1.TaskStarted",
   "MessageArgs":[
      "task85de4003-8757-2c7d-942f-55eaf7d6412a"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```




## Aggregates

An aggregate is a user-defined collection of resources.

The aggregate schema provides a mechanism to formally group the southbound resources of your choice. The advantage of creating aggregates is that they are more persistent than the random groupings—the aggregates are available and accessible in the environment of Resource Aggregator for ODIM until you delete them.

The resource aggregator allows you to:

-   Create an aggregate.

-   Populate an aggregate with the resources.

-   Perform actions on all the resources of an aggregate at once.

-   Delete an aggregate.



## Creating an aggregate

|||
|---------|-----------|
|<strong>Method</strong> | `POST` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/Aggregates` |
|<strong>Description</strong> |This operation creates an empty aggregate or an aggregate populated with resources.|
|<strong>Returns</strong> | The `Location` URI of the created aggregate having the aggregate Id. See the `Location` URI highlighted in bold in "Sample response header".<br>-   Link to the new aggregate, its Id, and a message saying that the resource has been created successfully in the JSON response body.<br>|
|<strong>Response Code</strong> |On success, `201 Created` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
      "Elements":[
            "/redfish/v1/Systems/{ComputerSystemId}"      
   ]   
}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Aggregates'


```

>**Sample request body**

```
{
      "Elements":[
            "/redfish/v1/Systems/8da0b6cd-42b7-4fd5-8ccf-97d0f58ae8c1:1"      
   ]   
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Elements|Array \(required\)<br> |An empty array or an array of links to the resources that this aggregate contains. To get the links to the system resources that are available in the resource inventory, perform HTTP `GET` on:<br> `/redfish/v1/Systems/` <br> |

>**Sample response header**

```
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Link:</redfish/v1/AggregationService/Aggregates/c14d91b5-3333-48bb-a7b7-75f74a137d48/>; rel=self
Location:/redfish/v1/AggregationService/Aggregates/c14d91b5-3333-48bb-a7b7-75f74a137d48
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Fri,21 August 2020 14:08:55 GMT+5m 11s
Transfer-Encoding:chunked
```

>**Sample response body**

```
{
      "@odata.type":"#Aggregate.v1_0_0.Aggregate",
      "@odata.id":"/redfish/v1/AggregationService/Aggregates/c14d91b5-3333-48bb-a7b7-75f74a137d48",
      "@odata.context":"/redfish/v1/$metadata#Aggregate.Aggregate",
      "Id":"c14d91b5-3333-48bb-a7b7-75f74a137d48",
      "Name":"Aggregate",
      "Message":"The resource has been created successfully",
      "MessageId":"Base.1.6.1.Created",
      "Severity":"OK",
      "Elements":[
            "/redfish/v1/Systems/8da0b6cd-42b7-4fd5-8ccf-97d0f58ae8c1:1",
            "/redfish/v1/Systems/4da0b6cd-42b7-4fd5-8ccf-97d0f58ae8b1:1"      
   ]   
}
```


## Viewing a list of aggregates

|||
|----------|-----------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/Aggregates` |
|<strong>Description</strong> |This operation lists all aggregates available in Resource Aggregator for ODIM.|
|<strong>Returns</strong> |A list of links to all the available aggregates.|
|<strong>Response Code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Aggregates'


```

>**Sample response body**

```
{
      "@odata.type":"#AggregateCollection.v1_0_0.AggregateCollection",
      "@odata.id":"/redfish/v1/AggregationService/Aggregates",
      "@odata.context":"/redfish/v1/$metadata#AggregateCollection.AggregateCollection",
      "Id":"Aggregate",
      "Name":"Aggregate",
      "Message":"Successfully Completed Request",
      "MessageId":"Base.1.6.1.Success",
      "Severity":"OK",
      "Members@odata.count":1,
      "Members":[
            {
                  "@odata.id":"/redfish/v1/AggregationService/Aggregates/c14d91b5-3333-48bb-a7b7-75f74a137d48"         
      }      
   ]   
}
```


## Viewing information about a single aggregate

|||
|----------|-----------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/Aggregates/{AggregateId}` |
|<strong>Description</strong> |This operation retrieves information about a specific aggregate.|
|<strong>Returns</strong> |JSON schema representing this aggregate.|
|<strong>Response Code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Aggregates/{AggregateId}'


```

>**Sample response body**

```
{
   "@odata.type":"#Aggregate.v1_0_0.Aggregate",
   "@odata.id":"/redfish/v1/AggregationService/Aggregates/c14d91b5-3333-48bb-a7b7-75f74a137d48",
   "@odata.context":"/redfish/v1/$metadata#Aggregate.Aggregate",
   "Id":"c14d91b5-3333-48bb-a7b7-75f74a137d48",
   "Name":"Aggregate",
   "Message":"Successfully Completed Request",
   "MessageId":"Base.1.6.1.Success",
   "Severity":"OK",
   "Elements":[
      "/redfish/v1/Systems/8da0b6cd-42b7-4fd5-8ccf-97d0f58ae8c1:1",
      "/redfish/v1/Systems/4da0b6cd-42b7-4fd5-8ccf-97d0f58ae8b1:1"      
   ]
}
```



## Deleting an aggregate

|||
|--------------|---------|
|<strong>Method</strong> | `DELETE` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/Aggregates/{AggregateId}` |
|<strong>Description</strong> |This operation deletes a specific aggregate.|
|<strong>Response Code</strong> |On success, `204 No Content` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i DELETE \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Aggregates/{AggregateId}'


```


## Adding elements to an aggregate

|||
|----------|-----------|
|<strong>Method</strong> | `POST` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/Aggregates/{AggregateId}/Actions/Aggregate.AddElements` |
|<strong>Description</strong> |This action adds one or more resources to a specific aggregate.|
|<strong>Returns</strong> |JSON schema for this aggregate having links to the added resources.|
|<strong>Response Code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
      "Elements":[
            "/redfish/v1/Systems/{ComputerSystemId1}",
            "/redfish/v1/Systems/{ComputerSystemId2}"     
   ]   
}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Aggregates/{AggregateId}/Actions/Aggregate.AddElements'


```

>**Sample request body**

```
{
      "Elements":[
            "/redfish/v1/Systems/8da0b6cd-42b7-4fd5-8ccf-97d0f58ae8c1:1",
            "/redfish/v1/Systems/7da0b6cd-42b7-4fd5-8ccf-97d0f58ae8e1:1"      
   ]   
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Elements|Array \(required\)<br> |An array of links to the Computer system resources that this aggregate contains.|

>**Sample response body**

```
{
      "@odata.type":"#Aggregate.v1_0_0.Aggregate",
      "@odata.id":"/redfish/v1/AggregationService/Aggregates/c14d91b5-3333-48bb-a7b7-75f74a137d48",
      "@odata.context":"/redfish/v1/$metadata#Aggregate.Aggregate",
      "Id":"c14d91b5-3333-48bb-a7b7-75f74a137d48",
      "Name":"Aggregate",
      "Message":"The resource has been created successfully",
      "MessageId":"Base.1.6.1.Created",
      "Severity":"OK",
      "Elements":[
            "/redfish/v1/Systems/8da0b6cd-42b7-4fd5-8ccf-97d0f58ae8c1:1",
            "/redfish/v1/Systems/4da0b6cd-42b7-4fd5-8ccf-97d0f58ae8b1:1"      
   ]   
}
```


## Resetting an aggregate of computer systems

|||
|--------|-----------|
|<strong>Method</strong> | `POST` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/Aggregates/{AggregateId}/Actions/Aggregate.Reset` |
|<strong>Description</strong> |This action shuts down, powers up, and restarts servers in a specific aggregate. This operation is performed in the background as a Redfish task and is further divided into subtasks to reset each server individually.<br> |
|<strong>Returns</strong> |- `Location` URI of the task monitor associated with this operation \(task\) in the response header. See `Location` URI highlighted in bold in "Sample response header \(HTTP 202 status\)".<br>-   Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See the task URI and the task Id highlighted in bold in "Sample response body \(HTTP 202 status\)".<br><blockquote>IMPORTANT:<br>Note down the task Id. If the task completes with an error, it is required to know which subtask has failed. To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`.<br></blockquote>- On successful completion of the reset operation, a message in the response body, saying that the reset operation is completed successfully. See "Sample response body \(HTTP 200 status\)".<br>|
|<strong>Response Code</strong> |`202 Accepted` On successful completion, `200 OK` <br> |
|<strong>Authentication</strong> |Yes|

**Usage information**

To know the progress of this action, perform HTTP `GET` on the [task monitor](#viewing-a-task-monitor) returned in the response header \(until the task is complete\).

To get the list of subtask URIs, perform HTTP `GET` on the task URI returned in the JSON response body. See "Sample response body \(HTTP 202 status\)". The JSON response body of each subtask contains a link to the task monitor associated with it. To know the progress of the reset operation \(subtask\) on a specific server, perform HTTP `GET` on the task monitor associated with the respective subtask. See the link to the task monitor highlighted in bold in "Sample response body \(subtask\)".


**NOTE:**

Only a user with `ConfigureComponents` privilege can reset servers. If you perform this action without necessary privileges, you will receive an HTTP `403 Forbidden` error.


>**curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
   "BatchSize":2,
   "DelayBetweenBatchesInSeconds":2,
   "ResetType":"ForceRestart"
}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Aggregates/{AggregateId}/Actions/Aggregate.Reset'


```

>**Sample request body**

```
{
   "BatchSize":2,
   "DelayBetweenBatchesInSeconds":2,
   "ResetType":"ForceRestart"
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|BatchSize|Integer \(optional\)<br> |The number of elements to be reset at a time in each batch.|
|DelayBetweenBatchesInSeconds|Integer \(seconds\) \(optional\)<br> |The delay among the batches of elements being reset.|
|ResetType|String \(optional\)<br> |For possible values, refer to "Reset type" table in [Resetting servers](GUID-22EC7FC3-6EF7-4A69-8DE1-385E3786E0C8.md).|

>**Sample response header** \(HTTP 202 status\)

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes

```

>**Sample response body** \(HTTP 202 status\)

```
{
   "@odata.type":"#Task.v1_4_2.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591",
   "Name":"Task task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591",
   "Message":"The task with id task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591 has started.",
   "MessageId":"TaskEvent.1.0.1.TaskStarted",
   "MessageArgs":[
      "task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```

>**Sample response body** \(subtask\)

```
{
   "@odata.type":"#SubTask.v1_4_2.SubTask",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591/SubTasks/task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "@odata.context":"/redfish/v1/$metadata#SubTask.SubTask",
   "Id":"task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591",
   "Name":"Task task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591",
   "Message":"Successfully Completed Request",
   "MessageId":"Base.1.6.1.Success",
   "Severity":"OK",
   "Members@odata.count":0,
   "Members":null,
   "TaskState":"Completed",
   "StartTime":"2020-05-13T13:33:59.917329733Z",
   "EndTime":"2020-05-13T13:34:00.320539988Z",
   "TaskStatus":"OK",
   "SubTasks":"",
   "TaskMonitor":"/taskmon/task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591",
   "PercentComplete":100,
   "Payload":{
      "HttpHeaders":null,
      "HttpOperation":"POST",
      "JsonBody":"",
      "TargetUri":"/redfish/v1/Systems/97d08f36-17f5-5918-8082-f5156618f58d:1"
   },
   "Messages":null
}
```

>**Sample response body** \(HTTP 200 status\)

```
 {
   "error":{
      "code":"Base.1.6.1.Success",
      "message":"Request completed successfully"
   }
}
```




## Setting boot order of an aggregate to default settings

|||
|----------|-----------|
|<strong>Method</strong> | `POST` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/Aggregates/{AggregateId}/Actions/Aggregate.SetDefaultBootOrder` |
|<strong>Description</strong> |This action changes the boot order of all the servers belonging to a specific aggregate to default settings. This operation is performed in the background as a Redfish task and is further divided into subtasks to change the boot order of each server individually.<br> |
|<strong>Returns</strong> |- `Location` URI of the created aggregate having the aggregate Id. See the `Location` URI highlighted in bold in "Sample response header".<br>-   Link to the new aggregate, its Id, and a message saying that the resource has been created successfully in the JSON response body.<br>`Location` URI of the task monitor associated with this operation in the response header. See `Location` URI highlighted in bold in "Sample response header \(HTTP 202 status\)".<br>-   Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See the task URI and the task Id highlighted in bold in "Sample response body \(HTTP 202 status\)".<br><blockquote>IMPORTANT:<br>Note down the task Id. If the task completes with an error, it is required to know which subtask has failed. To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`.<br></blockquote>On successful completion of this operation, a message in the response body, saying that the operation is completed successfully. See "Sample response body \(HTTP 200 status\)".<br>|
|<strong>Response Code</strong> |`202 Accepted` On successful completion, `200 OK` <br> |
|<strong>Authentication</strong> |Yes|

**Usage information**

To know the progress of this action, perform HTTP `GET` on the [task monitor](#viewing-a-task-monitor) returned in the response header \(until the task is complete\).

To get the list of subtask URIs, perform HTTP `GET` on the task URI returned in the JSON response body. See "Sample response body \(HTTP 202 status\)". The JSON response body of each subtask contains a link to the task monitor associated with it. To know the progress of `SetDefaultBootOrder` action \(subtask\) on a specific server, perform HTTP `GET` on the task monitor associated with the respective subtask. See the link to the task monitor highlighted in bold in "Sample response body \(subtask\)".


**NOTE:**

Only a user with `ConfigureComponents` privilege can change the boot order of one or more servers to default settings. If you perform this action without necessary privileges, you will receive an HTTP `403 Forbidden` error.


>**curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Aggregates/{AggregateId}/Actions/Aggregate.SetDefaultBootOrder'


```

>**Sample response header** \(HTTP 202 status\)

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/task85de4003-8057-4c7d-942f-55eaf7d6412a
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes

```

>**Sample response body** \(HTTP 202 status\)

```
{
   "@odata.type":"#Task.v1_4_2.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4003-8057-4c7d-942f-55eaf7d6412a",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task85de4003-8057-4c7d-942f-55eaf7d6412a",
   "Name":"Task task85de4003-8057-4c7d-942f-55eaf7d6412a",
   "Message":"The task with id task80de4003-8757-4c7d-942f-55eaf7d6412a has started.",
   "MessageId":"TaskEvent.1.0.1.TaskStarted",
   "MessageArgs":[
      "task80de4003-8757-4c7d-942f-55eaf7d6412a"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```

>**Sample response body** \(subtask\)

```
{
   "@odata.type":"#SubTask.v1_4_2.SubTask",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4003-8757-4c7d-942f-55eaf7d6412a/SubTasks/task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "@odata.context":"/redfish/v1/$metadata#SubTask.SubTask",
   "Id":"task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "Name":"Task task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "Message":"Successfully Completed Request",
   "MessageId":"Base.1.6.1.Success",
   "Severity":"OK",
   "Members@odata.count":0,
   "Members":null,
   "TaskState":"Completed",
   "StartTime":"2020-05-13T13:33:59.917329733Z",
   "EndTime":"2020-05-13T13:34:00.320539988Z",
   "TaskStatus":"OK",
   "SubTasks":"",
   "TaskMonitor":"/taskmon/task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "PercentComplete":100,
   "Payload":{
      "HttpHeaders":null,
      "HttpOperation":"POST",
      "JsonBody":"",
      "TargetUri":"/redfish/v1/Systems/97d08f36-17f5-5918-8082-f5156618f58d:1"
   },
   "Messages":null
}
```

>**Sample response body** \(HTTP 200 status\)

```
{ 
   "error":{ 
      "code":"Base.1.6.1.Success",
      "message":"Request completed successfully"
   }
}
```




## Removing elements from an aggregate

|||
|--------|---------|
|<strong>Method</strong> | `POST` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/Aggregates/{AggregateId}/Actions/Aggregate.RemoveElements` |
|<strong>Description</strong> |This action removes one or more resources from a specific aggregate.|
|<strong>Returns</strong> |Updated JSON schema representing this aggregate.|
|<strong>Response Code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
      "Elements":[
            "/redfish/v1/Systems/{ComputerSystemId1}",
            "/redfish/v1/Systems/{ComputerSystemId2}"     
   ]   
}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Aggregates/{AggregateId}/Actions/Aggregate.RemoveElements'


```

>**Sample request body**

```
{
      "Elements":[
            "/redfish/v1/Systems/8da0b6cd-42b7-4fd5-8ccf-97d0f58ae8c1:1",
            "/redfish/v1/Systems/7da0b6cd-42b7-4fd5-8ccf-97d0f58ae8e1:1"      
   ]   
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Elements|Array \(required\)<br> |An array of links to the Computer system resources that you want to remove from this aggregate.|

>**Sample response body**

```
{
   "@odata.type":"#Aggregate.v1_0_0.Aggregate",
   "@odata.id":"/redfish/v1/AggregationService/Aggregates/e02faf78-f919-4612-b031-bec7ae59910d",
   "@odata.context":"/redfish/v1/$metadata#Aggregate.Aggregate",
   "Id":"e02faf78-f919-4612-b031-bec7ae59910d",
   "Name":"Aggregate",
   "Message":"The resource has been removed successfully",
   "MessageId":"ResourceRemoved",
   "Severity":"OK",
   "Elements":[

   ]
}
```










#  Resource inventory

Resource Aggregator for ODIM allows you to view the inventory of compute and local storage resources through Redfish `Systems`, `Chassis`, and `Managers` endpoints. 
It also offers the capability to:	
- Search inventory information based on one or more configuration parameters.	
- Manage the resources added in the inventory. 


To discover crucial configuration information about a resource, including chassis, perform `GET` on these endpoints.

**Supported endpoints**

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/Systems|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}|GET, PATCH|`Login`, `ConfigureComponents` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Memory|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Memory/\{memoryId\}|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/MemoryDomains|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/NetworkInterfaces|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/EthernetInterfaces|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/EthernetInterfaces/\{Id\}|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Bios|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/SecureBoot|GET|`Login` |
|/redfish/v1/Systems/{ComputerSystemId}/PCIeDevices/{PCIeDeviceId}|`GET`|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage/\{storageSubsystemId\}|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage/\{storageSubsystemId\}/Drives/\{driveId\}|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage/\{storageSubsystemId\}/Volumes|GET, POST|`Login`, `ConfigureComponents` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Storage/\{storageSubsystemId\}/Volumes/\{volumeId\}|GET, DELETE|`Login`, `ConfigureComponents` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Processors|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Processors/\{Id\}|GET|`Login` |
|/redfish/v1/Systems?$filter=\{searchKeys\}%20\{conditionKeys\}%20\{value\}|GET|`Login` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Bios/Settings<br> |GET, PATCH|`Login`, `ConfigureComponents` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Actions/ComputerSystem.Reset|POST|`ConfigureComponents` |
|/redfish/v1/Systems/\{ComputerSystemId\}/Actions/ComputerSystem.SetDefaultBootOrder|POST|`ConfigureComponents` |

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/Chassis|GET, POST|`Login`, `ConfigureComponents` |
|/redfish/v1/Chassis/\{chassisId\}|GET, PATCH, DELETE|`Login`, `ConfigureComponents`|
|/redfish/v1/Chassis/\{chassisId\}/Thermal|GET|`Login`|
|/redfish/v1/Chassis/\{chassisId\}/NetworkAdapters|GET|`Login` |
|/redfish/v1/Chassis/{ChassisId}/NetworkAdapters/{networkadapterId}|GET|`Login`|

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/Managers|GET|`Login` |
|/redfish/v1/Managers/\{managerId\}|GET|`Login` |
|/redfish/v1/Managers/\{managerId\}/EthernetInterfaces|GET|`Login` |
|/redfish/v1/Managers/\{managerId\}/HostInterfaces|GET|`Login` |
|/redfish/v1/Managers/\{managerId\}/LogServices|GET|`Login` |
|/redfish/v1/Managers/\{managerId\}/NetworkProtocol|GET|`Login` |


>**NOTE:**
To view system, chassis, and manager resources, ensure that you have a minimum privilege of `Login`. If you do not have the necessary privileges, you will receive an HTTP `403 Forbidden` error.







##  Collection of computer systems

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems` |
|**Description** |This operation lists all systems available with Resource Aggregator for ODIM.|
|**Returns** |A collection of links to computer system instances.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Systems'


```

>**Sample response body** 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#ComputerSystemCollection.ComputerSystemCollection",
   "@odata.id":"/redfish/v1/Systems/",
   "@odata.type":"#ComputerSystemCollection.ComputerSystemCollection",
   "Description":"Computer Systems view",
   "Name":"Computer Systems",
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/Systems/ba0a6871-7bc4-5f7a-903d-67f3c205b08c:1"
      },
      { 
         "@odata.id":"/redfish/v1/Systems/7ff3bd97-c41c-5de0-937d-85d390691b73:1"
      }
   ],
   "Members@odata.count":2
}
```







## Single computer system

|||
|---------|-------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}` |
|**Description** |This endpoint fetches information about a specific system.|
|**Returns** |JSON schema representing this computer system instance.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}'

```

>**Sample response body** 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#ComputerSystem.ComputerSystem",
   "@odata.etag":"W/\"8C36EBD2\"",
   "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1",
   "@odata.type":"#ComputerSystem.v1_4_0.ComputerSystem",
   "Id":"e24fb205-6669-4080-b53c-67d4923aa73e:1",
   "Actions":{ 
      "#ComputerSystem.Reset":{ 
         "ResetType@Redfish.AllowableValues":[ 
            "On",
            "ForceOff",
            "GracefulShutdown",
            "ForceRestart",
            "Nmi",
            "PushPowerButton"
         ],
         "target":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/Actions/ComputerSystem.Reset"
      }
   },
   "AssetTag":"",
   "Bios":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/Bios"
   },
   "BiosVersion":"U32 v2.00 (02/02/2019)",
   "Boot":{ 
      "BootOptions":{ 
         "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/BootOptions"
      },
      "BootOrder":[ 
         "Boot000A",
         "Boot000B",
         "Boot000C",
         "Boot0012",
         "Boot000D",
         "Boot000F",
         "Boot000E",
         "Boot0010",
         "Boot0011",
         "Boot0013",
         "Boot0015",
         "Boot0014",
         "Boot0016"
      ],
      "BootSourceOverrideEnabled":"Disabled",
      "BootSourceOverrideMode":"UEFI",
      "BootSourceOverrideTarget":"None",
      "BootSourceOverrideTarget@Redfish.AllowableValues":[ 
         "None",
         "Cd",
         "Hdd",
         "Usb",
         "SDCard",
         "Utilities",
         "Diags",
         "BiosSetup",
         "Pxe",
         "UefiShell",
         "UefiHttp",
         "UefiTarget"
      ],
      "UefiTargetBootSourceOverride":"None",
      "UefiTargetBootSourceOverride@Redfish.AllowableValues":[ 
         "UsbClass(0xFFFF,0xFFFF,0xFF,0xFF,0xFF)",
         "PciRoot(0x0)/Pci(0x14,0x0)/USB(0x13,0x0)",
         "PciRoot(0x3)/Pci(0x0,0x0)/Pci(0x0,0x0)/Scsi(0x0,0x4000)",
         "PciRoot(0x3)/Pci(0x0,0x0)/Pci(0x0,0x0)/Scsi(0x1,0x4000)",
         "PciRoot(0x0)/Pci(0x1C,0x0)/Pci(0x0,0x0)/MAC(8030E02C92B0,0x1)/IPv4(0.0.0.0)/Uri()",
         "PciRoot(0x0)/Pci(0x1C,0x0)/Pci(0x0,0x0)/MAC(8030E02C92B0,0x1)/IPv4(0.0.0.0)",
         "PciRoot(0x0)/Pci(0x1C,0x0)/Pci(0x0,0x0)/MAC(8030E02C92B0,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)/Uri()",
         "PciRoot(0x0)/Pci(0x1C,0x0)/Pci(0x0,0x0)/MAC(8030E02C92B0,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)",
         "HD(2,GPT,E0698C18-D9A0-4F58-93CA-A6AEA6BFC93B,0x96800,0x32000)/\\EFI\\Microsoft\\Boot\\bootmgfw.efi",
         "PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(20677CEEF298,0x1)/IPv4(0.0.0.0)/Uri()",
         "PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(20677CEEF298,0x1)/IPv4(0.0.0.0)",
         "PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(20677CEEF298,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)/Uri()",
         "PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(20677CEEF298,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)"
      ]
   },
   "EthernetInterfaces":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/EthernetInterfaces"
   },
   "HostName":"",
   "IndicatorLED":"Off",
   "Links":{ 
      "ManagedBy":[ 
         { 
            "@odata.id":"/redfish/v1/Managers/e24fb205-6669-4080-b53c-67d4923aa73e:1"
         }
      ],
      "Chassis":[ 
         { 
            "@odata.id":"/redfish/v1/Chassis/e24fb205-6669-4080-b53c-67d4923aa73e:1"
         }
      ]
   },
   "LogServices":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/LogServices"
   },
   "Manufacturer":"HPE",
   "Memory":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/Memory"
   },
   "MemoryDomains":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/MemoryDomains"
   },
   "MemorySummary":{ 
      "Status":{ 
         "HealthRollup":"OK"
      },
      "TotalSystemMemoryGiB":384,
      "TotalSystemPersistentMemoryGiB":0
   },
   "Model":"ProLiant DL360 Gen10",
   "Name":"Computer System",
   "NetworkInterfaces":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/NetworkInterfaces"
   },
   "Oem":{ 
      
         },
         "AggregateHealthStatus":{ 
            "AgentlessManagementService":"Unavailable",
            "BiosOrHardwareHealth":{ 
               "Status":{ 
                  "Health":"OK"
               }
            },
            "FanRedundancy":"Redundant",
            "Fans":{ 
               "Status":{ 
                  "Health":"OK"
               }
            },
            "Memory":{ 
               "Status":{ 
                  "Health":"OK"
               }
            },
            "Network":{ 
               "Status":{ 
                  "Health":"OK"
               }
            },
            "PowerSupplies":{ 
               "PowerSuppliesMismatch":false,
               "Status":{ 
                  "Health":"OK"
               }
            },
            "PowerSupplyRedundancy":"Redundant",
            "Processors":{ 
               "Status":{ 
                  "Health":"OK"
               }
            },
            "SmartStorageBattery":{ 
               "Status":{ 
                  "Health":"OK"
               }
            },
            "Storage":{ 
               "Status":{ 
                  "Health":"OK"
               }
            },
            "Temperatures":{ 
               "Status":{ 
                  "Health":"OK"
               }
            }
         },
         "Bios":{ 
            "Backup":{ 
               "Date":"10/02/2018",
               "Family":"U32",
               "VersionString":"U32 v1.46 (10/02/2018)"
            },
            "Current":{ 
               "Date":"02/02/2019",
               "Family":"U32",
               "VersionString":"U32 v2.00 (02/02/2019)"
            },
            "UefiClass":2
         },
         "CurrentPowerOnTimeSeconds":38039,
         "DeviceDiscoveryComplete":{ 
            "AMSDeviceDiscovery":"NoAMS",
            "DeviceDiscovery":"vMainDeviceDiscoveryComplete",
            "SmartArrayDiscovery":"Complete"
         },
         "ElapsedEraseTimeInMinutes":0,
         "EndOfPostDelaySeconds":null,
         "EstimatedEraseTimeInMinutes":0,
         "IntelligentProvisioningAlwaysOn":true,
         "IntelligentProvisioningIndex":9,
         "IntelligentProvisioningLocation":"System Board",
         "IntelligentProvisioningVersion":"3.20.154",
         "IsColdBooting":false,
         "Links":{ 
            "PCIDevices":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIDevices"
            },
            "PCISlots":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCISlots"
            },
            "NetworkAdapters":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/BaseNetworkAdapters"
            },
            "SmartStorage":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/SmartStorage"
            },
            "USBPorts":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/USBPorts"
            },
            "USBDevices":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/USBDevices"
            },
            "EthernetInterfaces":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/EthernetInterfaces"
            },
            "WorkloadPerformanceAdvisor":{ 
               "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/WorkloadPerformanceAdvisor"
            }
         },
         "PCAPartNumber":"847479-001",
         "PCASerialNumber":"PVZEK0ARHBV392",
         "PostDiscoveryCompleteTimeStamp":"2020-02-23T23:09:45Z",
         "PostDiscoveryMode":null,
         "PostMode":null,
         "PostState":"InPostDiscoveryComplete",
         "PowerAllocationLimit":1000,
         "PowerAutoOn":"Restore",
         "PowerOnDelay":"Minimum",
         "PowerOnMinutes":463,
         "PowerRegulatorMode":"Dynamic",
         "PowerRegulatorModesSupported":[ 
            "OSControl",
            "Dynamic",
            "Max",
            "Min"
         ],
         "SMBIOS":{ 
            "extref":"/smbios"
         },
         "ServerFQDN":"",
         "SmartStorageConfig":[ 
            { 
               "@odata.id":"/redfish/v1/systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/smartstorageconfig"
            }
         ],
         "SystemROMAndiLOEraseComponentStatus":{ 
            "BIOSSettingsEraseStatus":"Idle",
            "iLOSettingsEraseStatus":"Idle"
         },
         "SystemROMAndiLOEraseStatus":"Idle",
         "SystemUsage":{ 
            "AvgCPU0Freq":126,
            "AvgCPU1Freq":0,
            "CPU0Power":62,
            "CPU1Power":54,
            "CPUICUtil":0,
            "CPUUtil":2,
            "IOBusUtil":0,
            "JitterCount":0,
            "MemoryBusUtil":0
         },
         "UserDataEraseComponentStatus":{ 

         },
         "UserDataEraseStatus":"Idle",
         "VirtualProfile":"Inactive"
      }
   },
   "PCIeDevices":[
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/1"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/2"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/3"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/4"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/5"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/6"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/7"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/8"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/9"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/10"
    },
    {
    "@odata.id": "/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/PCIeDevices/11"
    }
   ],
   "PCIeDevices@odata.count": 11,
   "PowerState":"On",
   "ProcessorSummary":{ 
      "Count":2,
      "Model":"Intel(R) Xeon(R) Gold 6152 CPU @ 2.10GHz",
      "Status":{ 
         "HealthRollup":"OK"
      }
   },
   "Processors":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/Processors"
   },
   "SKU":"867959-B21",
   "SecureBoot":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/SecureBoot"
   },
   "SerialNumber":"MXQ91100T9",
   "Status":{ 
      "Health":"OK",
      "HealthRollup":"OK",
      "State":"Starting"
   },
   "Storage":{ 
      "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/Storage"
   },
   "SystemType":"Physical",
   "TrustedModules":[ 
      { 
         "Oem":{ 
            
         },
         "Status":{ 
            "State":"Absent"
         }
      }
   ],
   "UUID":"39373638-3935-584D-5139-313130305439"
```





 



##  Memory collection

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Memory` |
|**Description** |This operation lists all memory devices of a specific server.|
|**Returns** |List of memory resource endpoints.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Memory'


```





## Single memory


|||
|---------|-------|
|**Method** |GET|
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Memory/{memoryId}` |
|**Description** |This endpoint retrieves configuration information of specific memory.|
|**Returns** |JSON schema representing this memory resource.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Memory/{memoryId}'


```

>**Sample response body** 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#Memory.Memory",
   "@odata.etag":"W/\"E6EC3A2C\"",
   "@odata.id":"/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e:1/Memory/proc1dimm1",
   "@odata.type":"#Memory.v1_7_1.Memory",
   "Id":"proc1dimm1",
   "BaseModuleType":"RDIMM",
   "BusWidthBits":72,
   "CacheSizeMiB":0,
   "CapacityMiB":32768,
   "DataWidthBits":64,
   "DeviceLocator":"PROC 1 DIMM 1",
   "ErrorCorrection":"MultiBitECC",
   "LogicalSizeMiB":0,
   "Manufacturer":"HPE",
   "MemoryDeviceType":"DDR4",
   "MemoryLocation":{ 
      "Channel":6,
      "MemoryController":2,
      "Slot":1,
      "Socket":1
   },
   "MemoryMedia":[ 
      "DRAM"
   ],
   "MemoryType":"DRAM",
   "Name":"proc1dimm1",
   "NonVolatileSizeMiB":0,
   "Oem":{ 
      
   },
   "OperatingMemoryModes":[ 
      "Volatile"
   ],
   "OperatingSpeedMhz":2666,
   "PartNumber":"840758-091",
   "PersistentRegionSizeLimitMiB":0,
   "RankCount":2,
   "SecurityCapabilities":{ 

   },
   "Status":{ 
      "Health":"OK",
      "State":"Enabled"
   },
   "VendorID":"52736",
   "VolatileRegionSizeLimitMiB":32768,
   "VolatileSizeMiB":32768
}
```





 




##  Memory domains

|||
|-------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/MemoryDomains` |
|**Description** | This endpoint lists memory domains of a specific system.<br> Memory Domains indicate to the client which Memory \(DIMMs\) can be grouped in Memory Chunks to form interleave sets or otherwise grouped.<br> |
|**Returns** |List of memory domain endpoints.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/MemoryDomains'


```




##  BIOS

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Bios` |
|**Description** | Use this endpoint to discover system-specific information about a BIOS resource and actions for changing to BIOS settings.<br>**NOTE:**<br> Changes to the BIOS typically require a system reset before they take effect.|
|**Returns** |<ul><li>Actions for changing password and resetting BIOS.</li><li>BIOS attributes.</li></ul> |
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Bios'


```





## Network interfaces

|||
|--------|---------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/NetworkInterfaces` |
|**Description** | This endpoint lists network interfaces of a specific system.<br> A network interface contains links to network adapter, network port, and network device function resources.<br> |
|**Returns** |List of network interface endpoints.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/NetworkInterfaces'


```






##  Ethernet interfaces

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/ EthernetInterfaces` |
|**Description** |This endpoint lists Ethernet interfaces or network interface controller \(NIC\) of a specific system.|
|**Returns** |List of Ethernet interface endpoints.|
|**Response Code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces'


```




## Single Ethernet interface

|||
|-----------|----------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/ EthernetInterfaces/{ethernetInterfaceId}` |
|**Description** |This endpoint retrieves information on a single, logical Ethernet interface or network interface controller \(NIC\).|
|**Returns** |JSON schema representing this Ethernet interface.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces/{ethernetInterfaceId}'


```


>**Sample response body** 

```
{
	"@odata.context": "/redfish/v1/$metadata#EthernetInterface.EthernetInterface",
	"@odata.etag": "W/\"5DEAF04A\"",
	"@odata.id": "/redfish/v1/Systems/97d08f36-17f5-5918-8082-f5156618f58d:1/EthernetInterfaces/1",
	"@odata.type": "#EthernetInterface.v1_4_1.EthernetInterface",
	"Id": "1",
	"FullDuplex": true,
	"IPv4Addresses": [],
	"IPv4StaticAddresses": [],
	"IPv6AddressPolicyTable": [],
	"IPv6Addresses": [],
	"IPv6StaticAddresses": [],
	"IPv6StaticDefaultGateways": [],
	"LinkStatus": null,
	"MACAddress": "80:30:e0:32:0a:58",
	"Name": "",
	"NameServers": [],
	"SpeedMbps": null,
	"StaticNameServers": [],
	"Status": {
		"Health": null,
		"State": null
	},
	"UefiDevicePath": "PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)"
}
```



##  PCIeDevice

|||
|---------|-------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/PCIeDevices/{PCIeDeviceId}` |
|**Description** | This operation fetches information about a specific PCIe device.<br> |
|**Returns** |Properties of a PCIe device attached to a computer system such as type, the version of the PCIe specification in use by this device, and more.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/PCIeDevices/{PCIeDeviceId}'


```
> Sample response body

```
{
    "@odata.context": "/redfish/v1/$metadata#PCIeDevice.PCIeDevice",
    "@odata.etag": "W/\"33150E20\"",
    "@odata.id": "/redfish/v1/Systems/1b77fcdd-b6a2-44b4-83f9-cfb4926fcd79:1/PCIeDevices/1",
    "@odata.type": "#PCIeDevice.v1_5_0.PCIeDevice",
    "Id": "1",
    "Name": "HPE Ethernet 1Gb 4-port 331i Adapter - NIC",
    "Oem": {
        "Hpe": {
            "@odata.context": "/redfish/v1/$metadata#HpeServerPciDevice.HpeServerPciDevice",
            "@odata.etag": "W/\"33150E20\"",
            "@odata.id": "/redfish/v1/Systems/1b77fcdd-b6a2-44b4-83f9-cfb4926fcd79:1/PCIDevices/1",
            "@odata.type": "#HpeServerPciDevice.v2_0_0.HpeServerPciDevice",
            "BusNumber": 2,
            "ClassCode": 2,
            "DeviceID": 5719,
            "DeviceInstance": 1,
            "DeviceLocation": "Embedded",
            "DeviceNumber": 0,
            "DeviceSubInstance": 1,
            "DeviceType": "Embedded LOM",
            "FunctionNumber": 0,
            "Id": "1",
            "LocationString": "Embedded LOM 1",
            "Name": "HPE Ethernet 1Gb 4-port 331i Adapter - NIC",
            "SegmentNumber": 0,
            "StructuredName": "NIC.LOM.1.1",
            "SubclassCode": 0,
            "SubsystemDeviceID": 8894,
            "SubsystemVendorID": 4156,
            "UEFIDevicePath": "PciRoot(0x0)/Pci(0x1C,0x0)/Pci(0x0,0x0)",
            "VendorID": 5348
        }
    }
}

```





##  Storage

|||
|---------|-------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Storage` |
|**Description** | This operation lists storage subsystems.<br> A storage subsystem is a set of storage controllers \(physical or virtual\) and the resources such as volumes that can be accessed from that subsystem.<br> |
|**Returns** |Links to storage subsystems.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage'


```





##  Storage subsystem

|||
|---------|-------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}` |
|**Description** | This operation lists resources such as drives and storage controllers in a storage subsystem.<br> |
|**Returns** |Links to the drives and storage controllers of a storage subsystem.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}'


```

> Sample response body 

```
{
    "@odata.context": "/redfish/v1/$metadata#Storage.Storage",
    "@odata.id": "/redfish/v1/Systems/49999b11-3e20-41e8-b6ca-2e466e6d8ccf:1/Storage/ArrayControllers-0",
    "@odata.type": "#Storage.v1_7_1.Storage",
    "Description": "HPE Smart Storage Array Controller View",
    "Drives": [
        {
            "@odata.id": "/redfish/v1/Systems/49999b11-3e20-41e8-b6ca-2e466e6d8ccf:1/Storage/ArrayControllers-0/Drives/0"
        },
        {
            "@odata.id": "/redfish/v1/Systems/49999b11-3e20-41e8-b6ca-2e466e6d8ccf:1/Storage/ArrayControllers-0/Drives/1"
        }
    ],
    "Id": "ArrayController-0",
    "Name": "HpeSmartStorageArrayController",
    "StorageControllers": [
        {
            "@odata.id": "/redfish/v1/Systems/49999b11-3e20-41e8-b6ca-2e466e6d8ccf:1/Storage/ArrayControllers-0#/StorageControllers/0",
            "FirmwareVersion": "1.98",
            "Manufacturer": "HPE",
            "MemberId": "0",
            "Model": "HPE Smart Array P408i-a SR Gen10",
            "Name": "HpeSmartStorageArrayController",
            "PartNumber": "836260-001",
            "PhysicalLocation": {
                "PartLocation": {
                    "LocationOrdinalValue": 0,
                    "LocationType": "Slot",
                    "ServiceLabel": "Slot=0"
                }
            },
            "SerialNumber": "PEYHC0DRHBV3CZ ",
            "Status": {
                "Health": "OK",
                "State": "Enabled"
            }
        }
    ],
    "Volumes": {
        "@odata.id": "/redfish/v1/Systems/49999b11-3e20-41e8-b6ca-2e466e6d8ccf:1/Storage/ArrayControllers-0/Volumes"
    }
}

```



## Drives

The drive schema represents a single physical drive for a system, including links to associated volumes.


###  Single drive

|||
|---------|-------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Drives/{driveId}` |
|**Description** | This operation retrieves information about a specific storage drive.<br> |
|**Returns** |JSON schema representing this drive.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Drives/{driveId}'


```




## Volumes

The volume schema represents a volume, virtual disk, LUN, or other logical storage entity for any system.


### Collection of volumes

| | |
|----------|-----------|
|<strong>Method</strong> |`GET` |
|<strong>URI</strong>  |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes` |
|<strong>Description</strong>  |This endpoint retrieves a collection of volumes in a specific storage subsystem.|
|<strong>Returns</strong> |A list of links to volumes.|
|<strong>Response code</strong> |On success, `200 OK` |
|<strong>Authentication</strong> |Yes|

 

>**curl command** 

```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odim_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes'


```

>**Sample response body** 

```
{
      "@odata.context":"/redfish/v1/$metadata#VolumeCollection.VolumeCollection",
      "@odata.etag":"W/\"AA6D42B0\"",
      "@odata.id":"/redfish/v1/Systems/eb452cf4-306c-4b21-96fb-698a067da407:1/Storage/ArrayControllers-0/Volumes",
      "@odata.type":"#VolumeCollection.VolumeCollection",
      "Description":"Volume Collection view",
      "Members":[
            {
                  "@odata.id":"/redfish/v1/Systems/eb452cf4-306c-4b21-96fb-698a067da407:1/Storage/ArrayControllers-0/Volumes/1"         
      }      
   ],
      "Members@odata.count":1,
      "Name":"Volume Collection"   
}
```



### Single volume


| | |
|----------|-----------|
|<strong>Method</strong> |`GET` |
|<strong>URI</strong>   |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/{volumeId}` |
|<strong>Description</strong>   |This endpoint retrieves information about a specific volume in a storage subsystem.|
|<strong>Returns</strong>  |JSON schema representing this volume.|
|<strong>Response code</strong>  |On success, `200 OK` |
|<strong>Authentication</strong>  |Yes|


>**curl command**


```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odim_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/{volumeId}'


```

>**Sample response body** 

```
{
   "@odata.context":"/redfish/v1/$metadata#Volume.Volume",
   "@odata.etag":"W/\"46916D5D\"",
   "@odata.id":"/redfish/v1/Systems/363bef34-7f89-48ac-8970-ee8955f1b56f:1/Storage/ArrayControllers-0/Volumes/1",
   "@odata.type":"#Volume.v1_4_1.Volume",
   "CapacityBytes":1200209526784,
   "Encrypted":false,
   "Id":"1",
   "Identifiers":[
      {
         "DurableName":"600508B1001C2AFE083D7F9026B2E994",
         "DurableNameFormat":"NAA"
      }
   ],
   "Links":{
      "Drives":[
         {
            "@odata.id":"/redfish/v1/Systems/363bef34-7f89-48ac-8970-ee8955f1b56f:1/Storage/ArrayControllers-0/Drives/0"
         }
      ]
   },
   "Name":"Drive_Volume_Link",
   "RAIDType":"RAID0",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   }
}
```


### Creating a volume

| | |
|----------|-----------|
|<strong>Method</strong> | `POST` |
|<strong>URI</strong>  |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes` |
|<strong>Description</strong>| This operation creates a volume in a specific storage subsystem.|
|<strong>Response code</strong>   |On success, `200 Ok` |
|<strong>Authentication</strong>|Yes|

>**curl command**

```
curl -i -X POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{
   "Name":"Volume_Demo",
   "RAIDType":"RAID1",
   "Drives":[
      {
         "@odata.id":"/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Drives/0"
      },
      {
         "@odata.id":"/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Drives/1"
      }
   ],
   "@Redfish.OperationApplyTime":"OnReset"
}}' \
 'https://{odim_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes'


```

>**Sample request body** 

```
{
   "Name":"Volume_Demo",
   "RAIDType":"RAID1",
   "Drives":[
      {
         "@odata.id":"/redfish/v1/Systems/363bef34-7f89-48ac-8970-ee8955f1b56f:1/Storage/ArrayControllers-0/Drives/0"
      },
      {
         "@odata.id":"/redfish/v1/Systems/363bef34-7f89-48ac-8970-ee8955f1b56f:1/Storage/ArrayControllers-0/Drives/1"
      }
   ],
   "@Redfish.OperationApplyTime":"OnReset"
}}
```

**Request parameters** 

|Parameter|Type|Description|
|---------|----|-----------|
|Name|String \(required\)<br> |Name of the new volume.|
|RAIDType|String \(required\)<br> |The RAID type of the volume you want to create.|
|Drives\[\{|Array \(required\)<br> |An array of links to drive resources to contain the new volume.|
|@odata.id \}\]<br> |String|A link to a drive resource.|
|@Redfish.OperationApplyTimeSupport|Redfish annotation \(optional\)<br> | It enables you to control when the operation is carried out.<br> Supported values: `OnReset` and `Immediate`.<br> `OnReset` indicates that the new volume will be available only after you have successfully reset the system. To know how to reset a system, see [Resetting a computer system](#resetting-a-computer-system).<br>`Immediate` indicates that the created volume will be available in the system immediately after the operation is successfully completed.|

>**Sample response body** 

```
 {
      "error":{
            "@Message.ExtendedInfo":[
                  {
                        "MessageId":"iLO.2.13.SystemResetRequired"            
         }         
      ],
            "code":"iLO.0.10.ExtendedInfo",
            "message":"See @Message.ExtendedInfo for more information."      
   }   
}
```





### Deleting a volume


| | |
|----------|-----------|
|<strong>Method</strong>  | `DELETE` |
|<strong>URI</strong>   |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/{volumeId}` |
|<strong>Description</strong>  | This operation removes a volume in a specific storage subsystem.|
|<strong>Response code</strong>|On success, `204 No Content` |
|<strong>Authentication</strong>  |Yes|

>**curl command**

```
curl -i -X DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
 'https://{odim_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/{volumeId}'


```

>**Sample request body** 

```
{
  
   "@Redfish.OperationApplyTime":"OnReset"
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|@Redfish.OperationApplyTimeSupport|Redfish annotation \(optional\)<br> | It enables you to control when the operation is carried out.<br> Supported values are: `OnReset` and `Immediate`. `OnReset` indicates that the volume will be deleted only after you have successfully reset the system.<br> `Immediate` indicates that the volume will be deleted immediatley after the operation is completed successfully.|




##  SecureBoot

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/SecureBoot` |
|**Description** |Use this endpoint to discover information on `UEFI Secure Boot` and managing the `UEFI Secure Boot` functionality of a specific system.|
|**Returns** | <ul><li>Action for resetting keys.</li><li> `UEFI Secure Boot` properties.<br>**NOTE:**<br> Use URI in the Actions group to discover information about resetting keys.</li></ul>|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/SecureBoot'


```





##  Processors

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Processors` |
|**Description** |This endpoint lists processors of a specific system.|
|**Returns** |List of processor resource endpoints.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Processors'


```






##  Single processor

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Processors/{processoId}` |
|**Description** |This endpoint fetches information about the properties of a processor attached to a specific server.|
|**Returns** |JSON schema representing this processor.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Processors/{processoId}'


```


## Chassis

Chassis represents the physical components of a system—sheet-metal confined spaces, logical zones such as racks, enclosures, chassis and all other containers, and subsystems \(like sensors\).

To view, create, and manage racks or rack groups, ensure that the URP \(Unmanaged Rack Plugin\) is running and is added into the Resource Aggregator for ODIM framework. To know how to add a plugin, see [Adding a plugin as an aggregation source](#adding-a-plugin-as-an-aggregation-source).

>**NOTE:**
URP is installed automatically during the deployment of the resource aggregator.


### Collection of chassis

|||
|-------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis` |
|**Description** | This operation lists chassis instances available with Resource Aggregator for ODIM.<br> Chassis represents the physical components of a system - sheet-metal confined spaces, logical zones such as racks, enclosures, chassis and all other containers, and subsystems \(like sensors\).<br> |
|**Returns** |A collection of links to chassis instances.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis'

```

>**Sample response body** 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#ChassisCollection.ChassisCollection",
   "@odata.id":"/redfish/v1/Chassis/",
   "@odata.type":"#ChassisCollection.ChassisCollection",
   "Description":"Computer System Chassis view",
   "Name":"Computer System Chassis",
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/Chassis/ba0a6871-7bc4-5f7a-903d-67f3c205b08c:1"
      },
      { 
         "@odata.id":"/redfish/v1/Chassis/7ff3bd97-c41c-5de0-937d-85d390691b73:1"
      }
   ],
   "Members@odata.count":2
}
```



 




>**Sample response body** 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#ChassisCollection.ChassisCollection",
   "@odata.id":"/redfish/v1/Chassis/",
   "@odata.type":"#ChassisCollection.ChassisCollection",
   "Description":"Computer System Chassis view",
   "Name":"Computer System Chassis",
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/Chassis/ba0a6871-7bc4-5f7a-903d-67f3c205b08c:1"
      },
      { 
         "@odata.id":"/redfish/v1/Chassis/7ff3bd97-c41c-5de0-937d-85d390691b73:1"
      }
   ],
   "Members@odata.count":2
}
```





### Single chassis

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}` |
|**Description** |This operation fetches information on a specific computer system chassis, rack group, or a rack.|
|**Returns** |JSON schema representing this chassis instance.|
|**Response code** |On success, `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}'


```

>**Sample response body** 

1. **Computer system chassis**

```
{ 
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.etag":"W/\"50540B90\"",
   "@odata.id":"/redfish/v1/Chassis/192083d2-c60a-4318-967b-cb5890c6dfe4:1",
   "@odata.type":"#Chassis.v1_6_0.Chassis",
   "Id":"192083d2-c60a-4318-967b-cb5890c6dfe4:1",
   "ChassisType":"RackMount",
   "Links":{ 
      "ManagedBy":[ 
         { 
            "@odata.id":"/redfish/v1/Managers/141cbba9-1e99-4272-b855-1781730bfe1c:1"
         }
      ],
      "ComputerSystems":[ 
         { 
            "@odata.id":"/redfish/v1/Systems/192083d2-c60a-4318-967b-cb5890c6dfe4:1"
         }
      ]
   },
   "Manufacturer":"HPE",
   "Model":"ProLiant DL380 Gen10",
   "Name":"Computer System Chassis",
   "NetworkAdapters":{ 
      "@odata.id":"/redfish/v1/Chassis/192083d2-c60a-4318-967b-cb5890c6dfe4:1/NetworkAdapters"
   },
   "Oem":{ 
      
         },
         "Firmware":{ 
            "PlatformDefinitionTable":{ 
               "Current":{ 
                  "VersionString":"8.9.0 Build 38"
               }
            },
            "PowerManagementController":{ 
               "Current":{ 
                  "VersionString":"1.0.4"
               }
            },
            "PowerManagementControllerBootloader":{ 
               "Current":{ 
                  "Family":"25",
                  "VersionString":"1.1"
               }
            },
            "SPSFirmwareVersionData":{ 
               "Current":{ 
                  "VersionString":"4.1.4.251"
               }
            },
            "SystemProgrammableLogicDevice":{ 
               "Current":{ 
                  "VersionString":"0x2A"
               }
            }
         },
         "Links":{ 
            "Devices":{ 
               "@odata.id":"/redfish/v1/Chassis/192083d2-c60a-4318-967b-cb5890c6dfe4:1/Devices"
            }
         },
         "MCTPEnabledOnServer":true,
         "SmartStorageBattery":[ 
            { 
               "ChargeLevelPercent":100,
               "FirmwareVersion":"0.70",
               "Index":1,
               "MaximumCapWatts":96,
               "Model":"875241-B21",
               "ProductName":"Smart Storage Battery ",
               "RemainingChargeTimeSeconds":7,
               "SerialNumber":"6WQXL0CB2BV63K",
               "SparePartNumber":"878643-001",
               "Status":{ 
                  "Health":"OK",
                  "State":"Enabled"
               }
            }
         ],
         "SystemMaintenanceSwitches":{ 
            "Sw1":"Off",
            "Sw10":"Off",
            "Sw11":"Off",
            "Sw12":"Off",
            "Sw2":"Off",
            "Sw3":"Off",
            "Sw4":"Off",
            "Sw5":"Off",
            "Sw6":"Off",
            "Sw7":"Off",
            "Sw8":"Off",
            "Sw9":"Off"
         }
      }
   },
   "Power":{ 
      "@odata.id":"/redfish/v1/Chassis/192083d2-c60a-4318-967b-cb5890c6dfe4:1/Power"
   },
   "SKU":"868704-B21",
   "SerialNumber":"2M291101JX",
   "Status":{ 
      "Health":"OK",
      "State":"Disabled"
   },
   "Thermal":{ 
      "@odata.id":"/redfish/v1/Chassis/192083d2-c60a-4318-967b-cb5890c6dfe4:1/Thermal"
   }
}
```

2. **Rack group chassis**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/f4e24c1c-dd2f-5a17-91b7-71620eb070df",
   "@odata.type":"#Chassis.v1_14_0.Chassis",
   "Id":"f4e24c1c-dd2f-5a17-91b7-71620eb070df",
   "Description":"My RackGroup",
   "Name":"RG8",
   "ChassisType":"RackGroup",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/99999999-9999-9999-9999-999999999999"
         }
      ]
   },
   "PowerState":"On",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   }
}
```

3. **Rack chassis**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "@odata.type":"#Chassis.v1_14_0.Chassis",
   "Id":"b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "Description":"rack no 1",
   "Name":"RACK#1",
   "ChassisType":"Rack",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/99999999-9999-9999-9999-999999999999"
         }
      ],
      "ContainedBy":[
         {
            "@odata.id":"/redfish/v1/Chassis/c2459269-011c-58d3-a217-ef914c4c295d"
         }
      ]
   },
   "PowerState":"On",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   }
}
```









###  Thermal metrics

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}/Thermal` |
|**Description** |This operation discovers information on the temperature and cooling of a specific chassis.|
|**Returns** |<ul><li>List of links to Fans</li><li>List of links to Temperatures</li></ul>|
| **Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}/Thermal'


```




### Collection of network adapters

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}/NetworkAdapters|
|**Description** | This endpoint lists network adapters contained in a chassis. A `NetworkAdapter` represents the physical network adapter capable of connecting to a computer network.<br> Some examples include Ethernet, fibre channel, and converged network adapters.|
|**Returns** |Links to network adapter instances available in this chassis.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}/NetworkAdapters'


```


### Single network adapter

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}/NetworkAdapters/{NetworkAdapterId}` |
|**Description** | This endpoint retrieves information on a specific network adapter.|
|**Returns** |JSON schema representing this network adapter.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|



>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}/NetworkAdapters/{NetworkAdapterId}'


```


>**Sample response body**

```

{
   "@odata.context":"/redfish/v1/$metadata#NetworkAdapter.NetworkAdapter",
   "@odata.etag":"W/\"F303ECE9\"",
   "@odata.id":"/redfish/v1/Chassis/a022faa5-107c-496d-874e-89c9f3e2df1c:1/NetworkAdapters/{rid}",
   "@odata.type":"#NetworkAdapter.v1_5_0.NetworkAdapter",
   "Description":"The network adapter resource instances available in this chassis.",
   "Name":"Network Adapter View",
   "Oem":{
      "Hpe":{
         "@odata.context":"/redfish/v1/$metadata#HpeBaseNetworkAdapter.HpeBaseNetworkAdapter",
         "@odata.etag":"W/\"7A9A9CE7\"",
         "@odata.id":"/redfish/v1/Systems/1/BaseNetworkAdapters/1/",
         "@odata.type":"#HpeBaseNetworkAdapter.v2_0_0.HpeBaseNetworkAdapter",
         "Id":"1",
         "FcPorts":[
            
         ],
         "Firmware":{
            "Current":{
               "VersionString":"20.14.54"
            }
         },
         "Name":"HPE Ethernet 1Gb 4-port 331i Adapter - NIC",
         "PhysicalPorts":[
            {
               "FullDuplex":true,
               "IPv4Addresses":[
                  
               ],
               "IPv6Addresses":[
                  
               ],
               "LinkStatus":null,
               "MacAddress":"80:30:e0:2c:92:a4",
               "Name":"",
               "Oem":{
                  "Hpe":{
                     "@odata.context":"/redfish/v1/$metadata#HpeBaseNetworkAdapterExt.HpeBaseNetworkAdapterExt",
                     "@odata.type":"#HpeBaseNetworkAdapterExt.v2_0_0.HpeBaseNetworkAdapterExt",
                     "BadReceives":0,
                     "BadTransmits":0,
                     "GoodReceives":0,
                     "GoodTransmits":0
                  }
               },
               "SpeedMbps":0
            },
            {
               "FullDuplex":true,
               "IPv4Addresses":[
                  
               ],
               "IPv6Addresses":[
                  
               ],
               "LinkStatus":null,
               "MacAddress":"80:30:e0:2c:92:a5",
               "Name":"",
               "Oem":{
                  "Hpe":{
                     "@odata.context":"/redfish/v1/$metadata#HpeBaseNetworkAdapterExt.HpeBaseNetworkAdapterExt",
                     "@odata.type":"#HpeBaseNetworkAdapterExt.v2_0_0.HpeBaseNetworkAdapterExt",
                     "BadReceives":0,
                     "BadTransmits":0,
                     "GoodReceives":0,
                     "GoodTransmits":0
                  }
               },
               "SpeedMbps":0
            },
            {
               "FullDuplex":true,
               "IPv4Addresses":[
                  
               ],
               "IPv6Addresses":[
                  
               ],
               "LinkStatus":null,
               "MacAddress":"80:30:e0:2c:92:a6",
               "Name":"",
               "Oem":{
                  "Hpe":{
                     "@odata.context":"/redfish/v1/$metadata#HpeBaseNetworkAdapterExt.HpeBaseNetworkAdapterExt",
                     "@odata.type":"#HpeBaseNetworkAdapterExt.v2_0_0.HpeBaseNetworkAdapterExt",
                     "BadReceives":0,
                     "BadTransmits":0,
                     "GoodReceives":0,
                     "GoodTransmits":0
                  }
               },
               "SpeedMbps":0
            },
            {
               "FullDuplex":true,
               "IPv4Addresses":[
                  
               ],
               "IPv6Addresses":[
                  
               ],
               "LinkStatus":null,
               "MacAddress":"80:30:e0:2c:92:a7",
               "Name":"",
               "Oem":{
                  "Hpe":{
                     "@odata.context":"/redfish/v1/$metadata#HpeBaseNetworkAdapterExt.HpeBaseNetworkAdapterExt",
                     "@odata.type":"#HpeBaseNetworkAdapterExt.v2_0_0.HpeBaseNetworkAdapterExt",
                     "BadReceives":0,
                     "BadTransmits":0,
                     "GoodReceives":0,
                     "GoodTransmits":0
                  }
               },
               "SpeedMbps":0
            }
         ],
         "Status":{
            "State":null
         },
         "StructuredName":"NIC.LOM.1.1",
         "UEFIDevicePath":"PciRoot(0x0)/Pci(0x1C,0x0)/Pci(0x0,0x0)"
      }
   }
}


```





###  Power

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}/Power` |
|**Description** |This operation retrieves power metrics specific to a server.|
|**Returns** |Information on power consumption and power limiting.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}/Power'


```

### Creating a rack group

|||
|---------|-------|
|Method | `POST` |
|URI |`/redfish/v1/Chassis`|
|Description |This operation creates a rack group.|
|Returns |<ul><li>`Location` header that contains a link to the created rack group \(highlighted in bold in the sample response header\).</li><li>JSON schema representing the created rack group.<br></li></ul>|
|Response code |On success, `201 Created`|
|Authentication |Yes|

 

>**curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
  "ChassisType": "RackGroup",
  "Description": "My RackGroup",
  "Links": {
    "ManagedBy": [
      {
        "@odata.id": "/redfish/v1/Managers/{managerId}"
      }
    ]
  },
  "Name": "RG5"
}
' \
 'https://{odim_host}:{port}/redfish/v1/Chassis'


```

>**Sample request body**

```
{
  "ChassisType": "RackGroup",
  "Description": "My RackGroup",
  "Links": {
    "ManagedBy": [
      {
        "@odata.id": "/redfish/v1/Managers/99999999-9999-9999-9999-999999999999"
      }
    ]
  },
  "Name": "RG5"
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|ChassisType|String \(required\)<br> |The type of chassis. The type to be used to create a rack group is RackGroup.<br> |
|Description|String \(optional\)<br> |Description of this rack group.|
|Links\{|Object \(required\)<br> |Links to the resources that are related to this rack group.|
|ManagedBy \[\{<br> @odata.id<br> \}\]<br> \}<br> |Array \(required\)<br> |An array of links to the manager resources that manage this chassis. The manager resource for racks and rack groups is the URP \(Unmanaged Rack Plugin\) manager. Provide the link to the URP manager.<br> |
|Name|String \(required\)<br> |Name for this rack group.|


>**Sample response header**

```
Connection:keep-alive
Content-Type:application/json; charset=UTF-8
Date:Wed,06 Jan 2021 09:37:43 GMT+15m 26s
**Location:/redfish/v1/Chassis/c2459269-011c-58d3-a217-ef914c4c295d**
Odata-Version:4.0
X-Frame-Options:sameorigin
Content-Length:462 bytes
```

>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/c2459269-011c-58d3-a217-ef914c4c295d",
   "@odata.type":"#Chassis.v1_14_0.Chassis",
   "Id":"c2459269-011c-58d3-a217-ef914c4c295d",
   "Description":"My RackGroup",
   "Name":"RG5",
   "ChassisType":"RackGroup",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/99999999-9999-9999-9999-999999999999"
         }
      ]
   },
   "PowerState":"On",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   }
}
```


### Creating a rack

|||
|---------|-------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/Chassis`|
|**Description**|This operation creates a rack.|
|**Returns** |<ul><li>`Location` header that contains a link to the created rack \(highlighted in bold in the sample response header\).</li><li>JSON schema representing the created rack.<br></li></ul>|
|**Response code** |On success, `201 Created`|
|**Authentication** |Yes|

 

>**curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
  "ChassisType": "Rack",
  "Description": "rack number one",
  "Links": {
    "ManagedBy": [
      {
        "@odata.id": "/redfish/v1/Managers/{managerId}"
      }
    ],
    "ContainedBy": [
      {
	    "@odata.id":"/redfish/v1/Chassis/{chassisId}"
	  }
    ]
  },
  "Name": "RACK#1"
}
' \
 'https://{odim_host}:{port}/redfish/v1/Chassis'


```

>**Sample request body**

```
{
  "ChassisType": "Rack",
  "Description": "rack number one",
  "Links": {
    "ManagedBy": [
      {
        "@odata.id": "/redfish/v1/Managers/675560ae-e903-41d9-bfb2-561951999999"
      }
    ],
    "ContainedBy": [
      {
	    "@odata.id":"/redfish/v1/Chassis/1be678f0-86dd-58ac-ac38-16bf0f6dafee"
	  }
    ]
  },
  "Name": "RACK#1"
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|ChassisType|String \(required\)<br> |The type of chassis. The type to be used to create a rack is Rack.<br> |
|Description|String \(optional\)<br> |Description of this rack.|
|Links\{|Object \(required\)<br> |Links to the resources that are related to this rack.|
|ManagedBy \[\{<br> @odata.id<br> \}\]<br> |Array \(required\)<br> |An array of links to the manager resources that manage this chassis. The manager resource for racks and rack groups is the URP \(Unmanaged Rack Plugin\) manager. Provide the link to the URP manager.<br> |
|ContainedBy \[\{<br> @odata.id<br> \}\]<br> \}<br> |Array \(required\)<br> |An array of links to the rack groups for containing this rack.|
|Name|String \(required\)<br> |Name for this rack group.|


>**Sample response header**

```
Connection:keep-alive
Content-Type:application/json; charset=UTF-8
Date:Wed,06 Jan 2021 09:37:43 GMT+15m 26s
**Location:/redfish/v1/Chassis/b6766cb7-5721-5077-ae0e-3bf3683ad6e2**
Odata-Version:4.0
X-Frame-Options:sameorigin
Content-Length:462 bytes
```

>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "@odata.type":"#Chassis.v1_14_0.Chassis",
   "Id":"b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "Description":"rack no 1",
   "Name":"RACK#1",
   "ChassisType":"Rack",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/99999999-9999-9999-9999-999999999999"
         }
      ],
      "ContainedBy":[
         {
            "@odata.id":"/redfish/v1/Chassis/c2459269-011c-58d3-a217-ef914c4c295d"
         }
      ]
   },
   "PowerState":"On",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   }
}
```



### Attaching chassis to a rack

|||
|---------|-------|
|**Method** | `PATCH` |
|**URI** |`/redfish/v1/Chassis/{rackId}`|
|**Description** |This operation attaches chassis to a specific rack.|
|**Returns** |JSON schema for the modified rack having links to the attached chassis.|
|**Response code** |On success, `200 Ok`|
|**Authentication** |Yes|

 

>**curl command**

```
curl -i PATCH \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
  "Links": {
    "Contains": [
      {
        "@odata.id": "/redfish/v1/Chassis/{chassisId}"
      }
    ]
  }
}
' \
 'https://{odim_host}:{port}/redfish/v1/Chassis/{rackId}'


```

>**Sample request body**

```
{
  "Links": {
    "Contains": [
      {
        "@odata.id": "/redfish/v1/Chassis/46db63a9-2dcb-43b3-bdf2-54ce9c42e9d9:1"
      }
    ]
  }
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Links\{|Object \(required\)<br> |Links to the resources that are related to this rack.|
|Contains \[\{<br> @odata.id<br> \}\]<br> \}<br> |Array \(required\)<br> |An array of links to the computer system chassis resources to be attached to this rack.|




>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "@odata.type":"#Chassis.v1_14_0.Chassis",
   "Id":"b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "Description":"rack no 1",
   "Name":"RACK#1",
   "ChassisType":"Rack",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/99999999-9999-9999-9999-999999999999"
         }
      ],
      "Contains":[
         {
            "@odata.id":"/redfish/v1/Chassis/4159c951-d0d0-4263-858b-0294f5be6377:1"
         }
      ],
      "ContainedBy":[
         {
            "@odata.id":"/redfish/v1/Chassis/c2459269-011c-58d3-a217-ef914c4c295d"
         }
      ]
   },
   "PowerState":"On",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   }
}
```


### Detaching chassis from a rack

|||
|---------|-------|
|**Method** | `PATCH` |
|**URI** |`/redfish/v1/Chassis/{rackId}`|
|**Description** |This operation detaches chassis from a specific rack.|
|**Returns** |JSON schema representing the modified rack.|
|**Response code** |On success, `200 Ok`|
|**Authentication** |Yes|

 

>**curl command**

```
curl -i PATCH \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
  "Links": {
    "Contains": []
  }
}
' \
 'https://{odim_host}:{port}/redfish/v1/Chassis/{rackId}'


```

>**Sample request body**

```
{
  "Links": {
    "Contains": []
  }
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Links\{|Object \(required\)<br> |Links to the resources that are related to this rack.|
|Contains \[\{<br> @odata.id<br> \}\]<br> \}<br> |Array \(required\)<br> |An array of links to the computer system chassis resources to be attached to this rack. To detach chassis from this rack, provide an empty array as value.|




>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "@odata.type":"#Chassis.v1_14_0.Chassis",
   "Id":"b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "Description":"rack no 1",
   "Name":"RACK#1",
   "ChassisType":"Rack",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/99999999-9999-9999-9999-999999999999"
         }
      ],
      "ContainedBy":[
         {
            "@odata.id":"/redfish/v1/Chassis/c2459269-011c-58d3-a217-ef914c4c295d"
         }
      ]
   },
   "PowerState":"On",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   }
}
```

### Deleting a rack

|||
|---------|-------|
|**Method** | `DELETE` |
|**URI** |`/redfish/v1/Chassis/{rackId}`|
|**Description** |This operation deletes a specific rack.<br>**IMPORTANT:**<br> If you try to delete a nonempty rack, you will receive an HTTP `409 Conflict` error. Ensure to detach the chassis attached to a rack before deleting it.<br>|
|**Response code** |On success, `204 No Content`|
|**Authentication** |Yes|

 

>**curl command**

```
curl -i DELETE \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   'https://{odim_host}:{port}/redfish/v1/Chassis/{rackId}'
```


>**Sample request body**

None.

### Deleting a rack group

|||
|---------|-------|
|**Method**| `DELETE` |
|**URI**|`/redfish/v1/Chassis/{rackGroupId}``|
|**Description**|This operation deletes a specific rack group.<br>**IMPORTANT:**<br>If you try to delete a nonempty rack group, you will receive an HTTP `409 Conflict` error. Ensure to remove all the racks contained in a rack group before deleting it.<br>|
|**Response code**|On success, `204 No Content`|
|**Authentication**|Yes|

 

>**curl command**

```
curl -i DELETE \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   'https://{odim_host}:{port}/redfish/v1/Chassis/{rackGroupId}`'
```


>**Sample request body**

None.



##  Searching the inventory

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems?$filter={searchKeys}%20{conditionKeys}%20{value}` |
|**Description** | Use this endpoint to search servers based on a filter - combination of a keyword, a condition, and a value.<br> Two ore more filters can be combined in a single request with the help of logical operands.<br>**NOTE:**<br> Only a user with `Login` privilege can perform this operation.|
|**Returns** |Server endpoints based on the specified filter.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Systems?$filter={searchKeys}%20{conditionKeys}%20{value/regular_expression}%20{logicalOperand}%20{searchKeys}%20{conditionKeys}%20{value}'

```



> Sample usage 

```
curl -i GET \
      -H "X-Auth-Token:{X-Auth-Token}" \
    'http://{odimra_host}:{port}/redfish/v1/Systems?$filter=MemorySummary/TotalSystemMemoryGiB%20eq%20384'

```

### Request URI parameters

-  `{searchkeys}` refers to `ComputerSystem` parameters. Following are the allowed search keys:

       -    `ProcessorSummary/Count` 
        
       -   `ProcessorSummary/Model` 
        
       -   `ProcessorSummary/sockets` 
        
       -    `SystemType` 
        
       -   `MemorySummary/TotalSystemMemoryGiB` 
        
       -   `FirmwareVersion` 
        
       -   `Storage/Drives/Quantity` 
        
       -   `Storage/Drives/Capacity` 
        
       -   `Storage/Drives/Type` 
	
-  `{conditionKeys}` refers to Redfish-specified conditions. Following are the allowed condition keys:

    |Condition Key|Meaning|Supported data type|
    |-------------|-------|-------------------|
    |"eq"|Equal to|All data types|
    |"ne"|Not equal to|All data types|
    |"gt"|Greater than|All numeric data types|
    |"ge"|Greater than or equal to|All numeric data types|
    |"le"|Lesser than or equal to|All numeric data types|
    |"lt"|Lesser than|All numeric data types|

     

- `{value}` refers to the actual value of the search parameter or a regular expression. Allowed regular expressions are as follows:

    `*, ?, ., $,%,^,&, /,!` 

      Examples:

         1. `$filter=TotalSystemMemoryGiB%20eq%20**384**` 
        
         2. `$filter=ProcessorSummary/Model%20eq%20**int\***` 
        
         3. `$filter=Storage/Drives/Type%20eq%20HDD` 

-  `{logicalOperands}` refers to the logical operands that are used to combine two or more filters in a request. Allowed logical operands are `and`, `or`, and `not`.



**Sample filters**

**Simple filter examples:**


1. `$filter=TotalSystemMemoryGiB%20eq%20384`
   
	This filter searches a server having total physical memory of 384 GB.



2. `$filter=ProcessorSummary/Model%20eq%20int*`

    This filter searches a server whose processor model name starts with `int` and ends with any combination of letters/numbers/special characters.

  


**Compound filter example:**


`$filter=(ProcessorSummary/Count%20eq%202%20and%20ProcessorSummary/Model%20eq%20intel)%20and%20(TotalSystemMemoryGiB%20eq%20384)`

This filter searches a server having total physical memory of 384 GB and two Intel processors.





 


>**Sample response body**

```
{ 
   "@odata.context":"/redfish/v1/$metadata#ComputerSystemCollection.ComputerSystemCollection",
   "@odata.id":"/redfish/v1/Systems/",
   "@odata.type":"#ComputerSystemCollection.ComputerSystemCollection",
   "Description":"Computer Systems view",
   "Name":"Computer Systems",
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/Systems/7ff3bd97-c41c-5de0-937d-85d390691b73:1"
      }
   ],
   "Members@odata.count":1
}
```








# Actions on a computer system



##  Resetting a computer system

|||
|---------|-------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Actions/ComputerSystem.Reset` |
|**Description** |This action shuts down, powers up, and restarts a specific system.<br>**NOTE:**<br> To reset an aggregate of systems, use this URI:<br>`/redfish/v1/AggregationService/Actions/AggregationService.Reset` <br> See [Resetting servers](#resetting-servers).|
|**Returns** |Message Id of the actual message in the JSON response body. To get the complete message, look up the specified registry file \(registry file name can be obtained by concatenating `RegistryPrefix` and version number present in the Message Id. Example registry file name: Base.1.4\). See [Message Registries](#message-registries).|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
 curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{
  "ResetType":"ForceRestart"
}
' \
 'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Actions/ComputerSystem.Reset'

```



>**Sample request body**

```
{
  "ResetType":"ForceRestart"
}
```

**Request parameters**

Refer to [Resetting Servers](#resetting-servers) to know about `ResetType.` 



>**Sample response body**

```
{
	"error": {
		"@Message.ExtendedInfo": [{
			"MessageId": "Base.1.4.Success"
		}],
		"code": "iLO.0.10.ExtendedInfo",
		"message": "See @Message.ExtendedInfo for more information."
	}
}
```



##  Changing the boot order of a computer system to default settings

|||
|--------|------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Actions/ComputerSystem.SetDefaultBootOrder` |
|**Description** |This action changes the boot order of a specific system to default settings.<br>**NOTE:**<br> To change the boot order of an aggregate of systems, use this URI:<br> `/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder` <br> See [Changing the Boot Order of Servers to Default Settings](#changing-the-boot-order-of-servers-to-default-settings).|
|**Returns** |Message Id of the actual message in the JSON response body. To get the complete message, look up the specified registry file \(registry file name can be obtained by concatenating `RegistryPrefix` and version number present in the Message Id. Example registry file name: Base.1.4\). See [Message Registries](#message-registries).|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
 curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'' \
 'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Actions/ComputerSystem.SetDefaultBootOrder'

```

>**Sample response body**

```
{
	"error": {
		"@Message.ExtendedInfo": [{
			"MessageId": "Base.1.4.Success"
		}],
		"code": "iLO.0.10.ExtendedInfo",
		"message": "See @Message.ExtendedInfo for more information."
	}
}
```






##  Changing BIOS settings

|||
|-------|-------|
|**Method** |`PATCH` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Bios/Settings` |
|**Description** |This action changes BIOS configuration.<br>**NOTE:**<br> Any change in BIOS configuration will be reflected only after the system resets. To see the change, [reset the computer system](#resetting-a-computer-system).|
|**Returns** |Message Id of the actual message in the JSON response body. To get the complete message, look up the specified registry file \(registry file name can be obtained by concatenating `RegistryPrefix` and version number present in the Message Id\). See [Message registries](#message-registries). For example,`MessageId` in the sample response body is `iLO.2.8.SystemResetRequired`. The registry to look up is `iLO.2.8`.<br> |
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
 curl -i -X PATCH \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{"Attributes": {"BootMode": "LegacyBios"}}' \
 'https://{odimra_host}:{port}/redfish/v1/Systems/{system_id}/Bios/Settings'

```




>**Sample request body**

```
{
	"Attributes": {
		"BootMode": "LegacyBios"
	}
}
```

**Request parameters**

`Attributes` are the list of BIOS attributes specific to the manufacturer or provider. To get a full list of attributes, perform `GET` on:


`https://{odimra_host}:{port}/redfish/v1/Systems/1/Bios/Settings`. 




>**Sample response body**

```
{ 
   "error":{ 
      "@Message.ExtendedInfo":[ 
         { 
            "MessageId":"iLO.2.8.SystemResetRequired"
         }
      ],
      "code":"iLO.0.10.ExtendedInfo",
      "message":"See @Message.ExtendedInfo for more information."
   }
}
```




## Changing the boot settings

|||
|---------|-------|
|**Method** |`PATCH` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}` |
|**Description** |This action changes the boot settings of a specific system such as boot source override target, boot order, and more.<br>**IMPORTANT**<br><ul><li>Ensure that the system is powered off before changing the boot order.</li><li>Power on the system once the operation is successful. The changes will be seen in the system only after a successful reset.</li></ul><br> To know how to power off, power on, or restart a system, see [Resetting a computer system](#resetting-a-computer-system).|
|**Returns** |Message Id of the actual message in the JSON response body. To get the complete message, look up the specified registry file \(registry file name can be obtained by concatenating `RegistryPrefix` and version number present in the Message Id\). See [Message Registries](#message-registries). For example,`MessageId` in the sample response body is `Base.1.0.Success`. The registry to look up is `Base.1.0`.<br> |
|**Response code** |`200 OK`|
|**Authentication** |Yes|

>**curl command**


```
 curl -i -X PATCH \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{ 
   "Boot":{ 
      "BootSourceOverrideTarget":"UefiHttp"
   }
}' \
 'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}}'

```




>**Sample request body**

```
{ 
   "Boot":{ 
      "BootSourceOverrideTarget":"UefiHttp"
   }
}
```

**Request parameters**

To get a full list of boot attributes that you can update, perform `GET` on:


`https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}`.


Check attributes under `Boot` in the JSON response. 

Some of the attributes include:

-   `BootSourceOverrideTarget` 

-   `UefiTargetBootSourceOverride` 

-   `Bootorder`


For possible values, see values listed under `{attribute}.AllowableValues`. 

> Example:

```
BootSourceOverrideTarget@Redfish.AllowableValues":[
"None",
"Cd",
"Hdd",
"Usb",
"SDCard",
"Utilities",
"Diags",
"BiosSetup",
"Pxe",
"UefiShell",
"UefiHttp",
"UefiTarget"
],
"UefiTargetBootSourceOverride@Redfish.AllowableValues":[
"HD(1,GPT,A6AC4D57-9D6D-46C8-8533-E7787450280D,0x800,0x4E2000)/\\EFI\\redhat\\shimx64.efi",
"PciRoot(0x3)/Pci(0x0,0x0)/Pci(0x0,0x0)/Scsi(0x0,0x0)",
"PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(ECEBB89E9928,0x1)/IPv4(0.0.0.0)",
"UsbClass(0xFFFF,0xFFFF,0xFF,0xFF,0xFF)",
"PciRoot(0x0)/Pci(0x14,0x0)/USB(0x13,0x0)",
"PciRoot(0x1)/Pci(0x0,0x0)/Pci(0x0,0x0)/MAC(48DF374FA220,0x1)/IPv4(0.0.0.0)/Uri()",
"PciRoot(0x1)/Pci(0x0,0x0)/Pci(0x0,0x0)/MAC(48DF374FA220,0x1)/IPv4(0.0.0.0)",
"PciRoot(0x1)/Pci(0x0,0x0)/Pci(0x0,0x0)/MAC(48DF374FA220,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)/Uri()",
"PciRoot(0x1)/Pci(0x0,0x0)/Pci(0x0,0x0)/MAC(48DF374FA220,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)",
"PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(ECEBB89E9928,0x1)/IPv4(0.0.0.0)/Uri()",
"PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(ECEBB89E9928,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)/Uri()",
"PciRoot(0x3)/Pci(0x2,0x0)/Pci(0x0,0x0)/MAC(ECEBB89E9928,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)",
"HD(1,GPT,F998DA94-45F7-4877-B907-10EDB6E65B07,0x800,0x64000)/\\EFI\\ubuntu\\shimx64.efi",
"HD(2,GPT,0AF5A707-BBEB-4D8A-8016-D840C2516753,0x40800,0x7700000)/\\EFI\\Microsoft\\Boot\\bootmgfw.efi"
]
}
```


**NOTE:**

If you attempt to update `BootSourceOverrideTarget` to `UefiTarget`, when `UefiTargetBootSourceOverride` is set to `None`, you will receive an HTTP `400 Bad Request` error. Update `UefiTargetBootSourceOverride` before setting `BootSourceOverrideTarget` to `UefiTarget`.



>**Sample response body**

```
{ 
   "error":{ 
      "@Message.ExtendedInfo":[ 
         { 
            "MessageId":"Base.1.0.Success"
         }
      ],
      "code":"iLO.0.10.ExtendedInfo",
      "message":"See @Message.ExtendedInfo for more information."
   }
```










# Managers

Resource Aggregator for ODIM exposes APIs to retrieve information about managers. Examples of managers include:

-   Resource Aggregator for ODIM itself

-   Baseboard Management Controllers \(BMC\)

-   Enclosure Managers

-   Management Controller

-   Other subsystems like plugins


**Supported endpoints**


|||
|-------|--------------------|
|/redfish/v1/Managers|`GET`|
|/redfish/v1/Managers/\{managerId\}|`GET`|
|/redfish/v1/Managers/\{managerId\}/EthernetInterfaces|`GET`|
|/redfish/v1/Managers/\{managerId\}/HostInterfaces|`GET`|
|/redfish/v1/Managers/\{managerId\}/LogServices|`GET`|
|/redfish/v1/Managers/\{managerId\}/NetworkProtocol|`GET`|
|/redfish/v1/Managers/\{managerId\}/VirtualMedia|`GET`|




##  Collection of managers

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Managers` |
|**Description** |A collection of managers.|
|**Returns** |Links to the manager instances. This collection includes a manager for Resource Aggregator for ODIM itself and other managers.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Managers'

```


>**Sample response body** 

```
{
   "@odata.context":"/redfish/v1/$metadata#ManagerCollection.ManagerCollection",
   "@odata.id":"/redfish/v1/Managers",
   "@odata.type":"#ManagerCollection.ManagerCollection",
   "Description":"Managers view",
   "Name":"Managers",
   "Members":[
      {
         "@odata.id":"/redfish/v1/Managers/a64fc187-e0e9-4f68-82a8-67a616b84b1d"
      },
      {
         "@odata.id":"/redfish/v1/Managers/141cbba9-1e99-4272-b855-1781730bfe1c:1"
      },
      {
         "@odata.id":"/redfish/v1/Managers/536cee48-84b2-43dd-b6e2-2459ac0eeac6"
      },
      {
         "@odata.id":"/redfish/v1/Managers/0e778112-4684-433d-9998-ca6f399c031f:1"
      },
      {
         "@odata.id":"/redfish/v1/Managers/a9cf0e1e-c36d-4d5b-9a31-cc07b611c01b"
      },
      {
         "@odata.id":"/redfish/v1/Managers/a6ddc4c0-2568-4e16-975d-fa771b0be853"
      }
   ],
   "Members@odata.count":6
}


```









##  Single manager

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Managers/{managerId}` |
|**Description** |A single manager.|
|**Returns** |Information about a specific management control system or a plugin or Resource Aggregator for ODIM itself. In the JSON schema representing a system \(BMC\) manager, there are links to the managers for:<ul><li>EthernetInterfaces:<br>`/redfish/v1/Managers/{managerId}/EthernetInterfaces`</li><li>HostInterfaces:<br>`/redfish/v1/Managers/{managerId}/HostInterfaces` </li><li>LogServices:<br>`/redfish/v1/Managers/{managerId}/LogServices` </li><li>NetworkProtocol:<br>`/redfish/v1/Managers/{managerId}/NetworkProtocol` <br> To know more about each manager, perform HTTP `GET` on these links.</li></ul>|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Managers/{managerId}'

```



>**Sample response body for a system \(BMC\) manager** 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#Manager.Manager",
   "@odata.etag":"W/\"2D2866FD\"",
   "@odata.id":"/redfish/v1/Managers/88b36c7c-d708-4a4a-8af5-5d58779d377c:1",
   "@odata.type":"#Manager.v1_3_3.Manager",
   "Actions":{ 
      "#Manager.Reset":{ 
         "target":"/redfish/v1/Managers/88b36c7c-d708-4a4a-8af5-5d58779d377c:1/Actions/Manager.Reset"
      }
   },
   "CommandShell":{ 
      "ConnectTypesSupported":[ 
         "SSH",
         "Oem"
      ],
      "MaxConcurrentSessions":9,
      "ServiceEnabled":true
   },
   "EthernetInterfaces":{ 
      "@odata.id":"/redfish/v1/Managers/88b36c7c-d708-4a4a-8af5-5d58779d377c:1/EthernetInterfaces"
   },
   "FirmwareVersion":"iLO 5 v1.40",
   "GraphicalConsole":{ 
      "ConnectTypesSupported":[ 
         "KVMIP"
      ],
      "MaxConcurrentSessions":10,
      "ServiceEnabled":true
   },
   "HostInterfaces":{ 
      "@odata.id":"/redfish/v1/Managers/88b36c7c-d708-4a4a-8af5-5d58779d377c:1/HostInterfaces"
   },
   "Id":"1",
   "Links":{ 
      "ManagerForChassis":[ 
         { 
            "@odata.id":"/redfish/v1/Chassis/88b36c7c-d708-4a4a-8af5-5d58779d377c:1"
         }
      ],
      "ManagerForServers":[ 
         { 
            "@odata.id":"/redfish/v1/Systems/88b36c7c-d708-4a4a-8af5-5d58779d377c:1"
         }
      ],
      "ManagerInChassis":{ 
         "@odata.id":"/redfish/v1/Chassis/88b36c7c-d708-4a4a-8af5-5d58779d377c:1"
      }
   },
   "LogServices":{ 
      "@odata.id":"/redfish/v1/Managers/88b36c7c-d708-4a4a-8af5-5d58779d377c:1/LogServices"
   },
   "ManagerType":"BMC",
   "Name":"Manager",
   "NetworkProtocol":{ 
      "@odata.id":"/redfish/v1/Managers/88b36c7c-d708-4a4a-8af5-5d58779d377c:1/NetworkProtocol"
   },
   "Oem":{},
      
   "SerialConsole":{ 
      "ConnectTypesSupported":[ 
         "SSH",
         "IPMI",
         "Oem"
      ],
      "MaxConcurrentSessions":13,
      "ServiceEnabled":true
   },
   "Status":{ 
      "State":"Absent"
   },
   "UUID":"a964d6a9-a45c-57aa-90ec-08b38850b2f3",
   "VirtualMedia":{ 
      "@odata.id":"/redfish/v1/Managers/88b36c7c-d708-4a4a-8af5-5d58779d377c:1/VirtualMedia"
   }
}
```

>**Sample response body for Resource Aggregator for ODIM manager**

```
{
   "@odata.context":"/redfish/v1/$metadata#Manager.Manager",
   "@odata.id":"/redfish/v1/Managers/a64fc187-e0e9-4f68-82a8-67a616b84b1d",
   "@odata.type":"#Manager.v1_3_3.Manager",
   "Name":"ODIMRA",
   "ManagerType":"Service",
   "Id":"a64fc187-e0e9-4f68-82a8-67a616b84b1d",
   "UUID":"a64fc187-e0e9-4f68-82a8-67a616b84b1d",
   "FirmwareVersion":"1.0",
   "Status":{
      "State":"Enabled"
   }
}
```

>**Sample response body for a plugin manager**

```
{
   "@odata.context":"/redfish/v1/$metadata#Manager.Manager",
   "@odata.etag":"W/\"AA6D42B0\"",
   "@odata.id":"/redfish/v1/Managers/536cee48-84b2-43dd-b6e2-2459ac0eeac6",
   "@odata.type":"#Manager.v1_3_3.Manager",
   "FirmwareVersion":"v1.0.0",
   "Id":"a9cf0e1e-c36d-4d5b-9a31-cc07b611c01b",
   "ManagerType":"Service",
   "Name":"GRF",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   },
   "UUID":"a9cf0e1e-c36d-4d5b-9a31-cc07b611c01b"
}
```

## VirtualMedia

VirtualMedia enables you to connect remote storage media (such as CD-ROM, USB mass storage, ISO image, and floppy disk) to a target server on a network. The target server can access the remote media, read from and write to it as if it were physically connected to the server USB port.

Resource Aggregator for ODIM exposes Redfish `VirtualMedia` APIs to connect the remote storage media to your servers.

**Supported APIs**:

| API URI                                                      | Operation Applicable | Required privileges   |
| ------------------------------------------------------------ | -------------------- | --------------------- |
| /redfish/v1/Managers/{ManagerId}/VirtualMedia                | GET                  | `Login`               |
| /redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID} | GET                  | `Login`               |
| /redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}/Actions/VirtualMedia.InsertMedia | POST                 | `ConfigureComponents` |
| /redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}/Actions/VirtualMedia.EjectMedia | POST                 | `ConfigureComponents` |

>**NOTE:**
>Before accessing these endpoints, ensure that the user has the required privileges. If you access these endpoints without necessary privileges, you will receive an HTTP `403 Forbidden` error.

### Viewing the VirtualMedia collection

| **Method**         | `GET`                                                        |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Managers/{ManagerId}/VirtualMedia`              |
| **Description**    | This operation lists all virtualmedia collections available in Resource Aggregator for ODIM. |
| **Returns**        | A list of links to all the available virtualmedia collections. |
| **Response Code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/Managers/{ManagerId}/VirtualMedia'
```

>**Sample response body**

```
{
  "@odata.context":"/redfish/v1/$metadata#VirtualMediaCollection.VirtualMediaCollection",
   "@odata.id":"/redfish/v1/Managers/1/VirtualMedia/",
   "@odata.type":"#VirtualMediaCollection.VirtualMediaCollection",
   "Description":"Virtual Media Services Settings",
   "Name":"Virtual Media Services",
   "Members":[
      {
         "@odata.id":"/redfish/v1/Managers/1/VirtualMedia/1/"
      },
      {
         "@odata.id":"/redfish/v1/Managers/1/VirtualMedia/2/"
      }
   ],
   "Members@odata.count":2
}
```

### Viewing a VirtualMedia Instance

| <strong>Method</strong>         | `GET`                                                        |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>URI</strong>            | `/redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}` |
| <strong>Description</strong>    | This action retrieves information about a specific virtualmedia instance. |
| <strong>Returns</strong>        | JSON schema representing this virtualmedia instance.         |
| <strong>Response Code</strong>  | On success, `200 Ok`                                         |
| <strong>Authentication</strong> | Yes                                                          |

>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}'
```


>**Sample response body**

```
{
    "error": {
        "@Message.ExtendedInfo": [
            {
                "MessageId": "Base.1.4.Success"
            }
        ],
        "code": "iLO.0.10.ExtendedInfo",
        "message": "See @Message.ExtendedInfo for more information."
    }
}
```

## Inserting VirtualMedia

| **Method**         | `POST`                                                       |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}/Actions/VirtualMedia.InsertMedia` |
| **Description**    | This operation inserts the virtual media on to the manager.  |
| **Returns**        | A message stating the virtual media insertion was successful. |
| **Response Code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

>**curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
	'{
  "Image":"http://<ip address>/<image path>",
  "Inserted":true,
  "WriteProtected":true
}' \
 'https://{odimra_host}:{port}/redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}/Actions/VirtualMedia.InsertMedia'

```

>**Sample response body**

```
{
    "error": {
        "@Message.ExtendedInfo": [
            {
                "MessageId": "Base.1.4.Success"
            }
        ],
        "code": "iLO.0.10.ExtendedInfo",
        "message": "See @Message.ExtendedInfo for more information."
    }
}
 
```

## Ejecting VirtualMedia

| **Method**         | `POST`                                                       |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}/Actions/VirtualMedia.EjectMedia` |
| **Description**    | This operation ejects the virtual media from the manager.    |
| **Returns**        | A message stating the virtual media ejection was successful. |
| **Response Code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

>**curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   
 'https://{odimra_host}:{port}/redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}/Actions/VirtualMedia.EjectMedia'
```

<blockquote>NOTE: No payload is required for this operation. </blockquote>

>**Sample response body**


```
{
    "error": {
        "@Message.ExtendedInfo": [
            {
                "MessageId": "Base.1.4.Success"
            }
        ],
        "code": "iLO.0.10.ExtendedInfo",
        "message": "See @Message.ExtendedInfo for more information."
    }
}
 
```

**Request parameters**

| Parameter      | Type               | Description           |
| -------------- | ------------------ | --------------------- |
| Image          | String (Required)  | Image path            |
| Inserted       | Boolean (Optional) | Default value is true |
| WriteProtected | Boolean (Optional) | Default value is true |

# Software and firmware inventory

The resource aggregator exposes Redfish update service endpoints. Use these endpoints to access and update the software components of a system such as BIOS and firmware. Using these endpoints, you can also upgrade or downgrade firmware of other components such as system drivers and provider software.

The `UpdateService` schema describes the update service and the properties for the service itself. It exposes the firmware and software inventory resources and provides links to access them.


**Supported endpoints**

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/UpdateService|GET|`Login` |
|/redfish/v1/UpdateService/FirmwareInventory|GET|`Login` |
|/redfish/v1/UpdateService/FirmwareInventory/\{inventoryId\}|GET|`Login` |
|/redfish/v1/UpdateService/SoftwareInventory|GET|`Login` |
|/redfish/v1/UpdateService/SoftwareInventory/\{inventoryId\}|GET|`Login` |
|/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate|POST|`ConfigureComponents` |
|/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate|POST|`ConfigureComponents` |

<blockquote>
NOTE:

Before accessing these endpoints, ensure that the user account has the required privileges. If you access these endpoints without necessary privileges, you will receive an HTTP `403 Forbidden` error.

</blockquote>

## Viewing the update service root

| | |
|-----|------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/UpdateService` |
|<strong>Description</strong> |This operation retrieves JSON schema representing the `UpdateService` root.|
|<strong>Returns</strong> |Properties for the service and a list of actions you can perform using this service.|
|<strong>Response code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/UpdateService'


```

>**Sample response body**

```
{
   "@odata.type":"#UpdateService.v1_8_1.UpdateService",
   "@odata.id":"/redfish/v1/UpdateService",
   "@odata.context":"/redfish/v1/$metadata#UpdateService.UpdateService",
   "Id":"UpdateService",
   "Name":"Update Service",
   "Status":{
      "State":"Enabled",
      "Health":"OK",
      "HealthRollup":"OK"
   },
   "ServiceEnabled":true,
   "HttpPushUri":"",
   "FirmwareInventory":{
      "@odata.id":"/redfish/v1/UpdateService/FirmwareInventory"
   },
   "SoftwareInventory":{
      "@odata.id":"/redfish/v1/UpdateService/SoftwareInventory"
   },
   "Action":{
      "#UpdateService.SimpleUpdate":{
         "target":"/redfish/v1/UpdateService/Actions/SimpleUpdate"
      },
      "#UpdateService.StartUpdate":{
         "target":"/redfish/v1/UpdateService/Actions/StartUpdate"
      }
   }
}
```

## Viewing the firmware inventory

| |  |
|-------|---------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/UpdateService/FirmwareInventory` |
|<strong>Description</strong> |This operation lists firmware of all the resources available in Resource Aggregator for ODIM.|
|<strong>Returns</strong> |A collection of links to firmware resources.|
|<strong>Response code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/UpdateService/FirmwareInventory'


```

>**Sample response body**

```
{
   ​   "@odata.context":"/redfish/v1/$metadata#FirmwareInventoryCollection.FirmwareCollection",
   ​   "@odata.id":"/redfish/v1/UpdateService/FirmwareInventory",
   ​   "@odata.type":"#FirmwareInventoryCollection.FirmwareInventoryCollection",
   ​   "Description":"FirmwareInventory view",
   ​   "Name":"FirmwareInventory",
   ​   "Members":​[
      ​      {
         ​         "@odata.id":"/redfish/v1/UpdateService/FirmwareInventory/4c12d2f7-a8e2-430f-bff2-737a80e73803:10"         ​
      },
      ​      {
         ​         "@odata.id":"/redfish/v1/UpdateService/FirmwareInventory/4c12d2f7-a8e2-430f-bff2-737a80e73803:9"         ​
      },
      ​      {
         ​         "@odata.id":"/redfish/v1/UpdateService/FirmwareInventory/4c12d2f7-a8e2-430f-bff2-737a80e73803:6"         ​
      },
      ​      {
         ​         "@odata.id":"/redfish/v1/UpdateService/FirmwareInventory/4c12d2f7-a8e2-430f-bff2-737a80e73803:17"         ​
      },
      ​      {
         ​         "@odata.id":"/redfish/v1/UpdateService/FirmwareInventory/4c12d2f7-a8e2-430f-bff2-737a80e73803:13"         ​
      },
      ​      {
         ​         "@odata.id":"/redfish/v1/UpdateService/FirmwareInventory/4c12d2f7-a8e2-430f-bff2-737a80e73803:5"         ​
      },
      ​      {
         ​         "@odata.id":"/redfish/v1/UpdateService/FirmwareInventory/4c12d2f7-a8e2-430f-bff2-737a80e73803:8"         ​
      },
      ​      {
         ​         "@odata.id":"/redfish/v1/UpdateService/FirmwareInventory/4c12d2f7-a8e2-430f-bff2-737a80e73803:12"         ​
      }      ​
   ],
   ​   "Members@odata.count":8​
}​
```



## Viewing a specific firmware resource

| | |
|-------|-----------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/UpdateService/FirmwareInventory/{inventoryId}` |
|<strong>Description</strong> |This operation retrieves information about a specific firmware resource.|
|<strong>Returns</strong> |JSON schema representing this firmware.|
|<strong>Response code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/UpdateService/FirmwareInventory/{inventoryId}'


```

>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#SoftwareInventory.SoftwareInventory",
   "@odata.etag":"W/\"0539D502\"",
   "@odata.id":"/redfish/v1/UpdateService/FirmwareInventory/3",
   "@odata.type":"#SoftwareInventory.v1_0_0.SoftwareInventory",
   "Description":"PlatformDefinitionTable",
   "Id":"3",
   "Name":"Intelligent Platform Abstraction Data",
   "Oem":{
      "Hpe":{
         "@odata.context":"/redfish/v1/$metadata#HpeiLOSoftwareInventory.HpeiLOSoftwareInventory",
         "@odata.type":"#HpeiLOSoftwareInventory.v2_0_0.HpeiLOSoftwareInventory",
         "DeviceClass":"b8f46d06-85db-465c-94fb-d106e61378ed",
         "DeviceContext":"System Board",
         "Targets":[
            "00000000-0000-0000-0000-000000000204",
            "00000000-0000-0000-0000-000001553332"
         ]
      }
   },
   "Version":"8.5.0 Build 15"
}
```



## Viewing the software inventory

| | |
|-------|-----------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/UpdateService/SoftwareInventory` |
|<strong>Description</strong> |This operation lists software of all the resources available in Resource Aggregator for ODIM.|
|<strong>Returns</strong> |A collection of links to software resources.|
|<strong>Response code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/UpdateService/SoftwareInventory'


```

>**Sample response body**

```
{
   ​   "@odata.context":"/redfish/v1/$metadata#SoftwareInventoryCollection.SoftwareCollection",
   ​   "@odata.id":"/redfish/v1/UpdateService/SoftwareInventory",
   ​   "@odata.type":"#SoftwareInventoryCollection.SoftwareInventoryCollection",
   ​   "Description":"SoftwareInventory view",
   ​   "Name":"SoftwareInventory",
   ​   "Members":null,
   ​   "Members@odata.count":0​
}
```


## Viewing a specific software resource

| | |
|-------|-----------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/UpdateService/SoftwareInventory/{inventoryId}` |
|<strong>Description</strong> |This operation retrieves information about a specific software resource.|
|<strong>Returns</strong> |JSON schema representing this software.|
|<strong>Response code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/UpdateService/SoftwareInventory/{inventoryId}'


```

>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#SoftwareInventory.SoftwareInventory",
   "@odata.etag":"W/\"0539D502\"",
   "@odata.id":"/redfish/v1/UpdateService/SoftwareInventory/3",
   "@odata.type":"#SoftwareInventory.v1_0_0.SoftwareInventory",
   "Description":"PlatformDefinitionTable",
   "Id":"3",
   "Name":"Intelligent Platform Abstraction Data",
   "Oem":{
      "Hpe":{
         "@odata.context":"/redfish/v1/$metadata#HpeiLOSoftwareInventory.HpeiLOSoftwareInventory",
         "@odata.type":"#HpeiLOSoftwareInventory.v2_0_0.HpeiLOSoftwareInventory",
         "DeviceClass":"b8f46d06-85db-465c-94fb-d106e61378ed",
         "DeviceContext":"System Board",
         "Targets":[
            "00000000-0000-0000-0000-000000000204",
            "00000000-0000-0000-0000-000001553332"
         ]
      }
   },
   "Version":"8.5.0 Build 15"
}
```


## Actions

### Simple update

| | |
|-------|-----------|
|<strong>Method</strong> | `POST` |
|<strong>URI</strong> |`/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate` |
|<strong>Description</strong> |This operation creates an update request for updating a software or a firmware component or directly updates a software or a firmware component. The first example in "Sample request body" is used to create an update request and the second one is used to directly update a software or a firmware component of servers.<br>It is performed in the background as a Redfish task. |
|<strong>Response code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|


**Usage information** 
To know the progress of this action, perform HTTP `GET` on the [task monitor](#viewing-a-task-monitor) returned in the response header \(until the task is complete\).


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{
"ImageURI": "<URI_of_the_firmware_image>",
"Password": "<password>",
"Targets": ["/redfish/v1/Systems/{ComputerSystemId}"],
"TransferProtocol": "",
"Username": "<username>"
}' \
 'https://{odim_host}:{port}/redfish/v1/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate'


```

> Sample request body

**Example 1:** 

```
{
  "ImageURI":"http://{IP_address}/ISO/resource.bin",
  "Targets": ["/redfish/v1/Systems/65d01621-4f88-49de-98bc-fcd1419bff3a:1"],
}
```

**Example 2:** 

```
{
  "ImageURI":"http://{IP_address}/ISO/resource.bin",
  "Targets": ["/redfish/v1/Systems/65d01621-4f88-49de-98bc-fcd1419bff3a:1"],
  "@Redfish.OperationApplyTimeSupport": {
            "@odata.type": "#Settings.v1_2_0.OperationApplyTimeSupport",
              "SupportedValues": ["OnStartUpdate"]
            }
}
```

#### Request parameters

|Parameter|Type|Description|
|---------|----|-----------|
|ImageURI|String \(required\)<br> |The URI of the software or firmware image to install. It is the location address of the software or firmware image you want to install.|
|Password|String \(optional\)<br> |The password to access the URI specified by the Image URI parameter.|
|Targets\[\]|Array \(required\)<br> |An array of URIs that indicate where to apply the update image.|
|TransferProtocol|String \(optional\)<br> | The network protocol that the update service uses to retrieve the software or the firmware image file at the URI provided in the `ImageURI` parameter, if the URI does not contain a scheme.<br> For the possible property values, see "Transfer protocol" table.<br> |
|Username|String \(optional\)<br> |The user name to access the URI specified by the Image URI parameter.|
|@Redfish.OperationApplyTimeSupport|Redfish annotation \(optional\)<br> | It enables you to control when the update is carried out.<br> Supported value is: `OnStartUpdate`. It indicates that the update will be carried out only after you perform HTTP POST on:<br> `/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate`.<br> |

|String|Description|
|------|-----------|
|CIFS|Common Internet File System.|
|FTP|File Transfer Protocol.|
|HTTP|Hypertext Transfer Protocol.|
|HTTPS|Hypertext Transfer Protocol Secure.|
| NFS \(v1.3+\)<br> |Network File System.|
| NSF \(deprecated v1.3\)<br> | Network File System.<br> This value has been deprecated in favor of NFS.<br> |
|OEM|A manufacturer-defined protocol.|
|SCP|Secure Copy Protocol.|
| SFTP \(v1.1+\)<br> |Secure File Transfer Protocol.|
|TFTP|Trivial File Transfer Protocol.|


>**Sample response header \(HTTP 202 status\)**

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/task4aac9e1e-df58-4fff-b781-52373fcb5699
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes

```

>**Sample response body \(HTTP 202 status\)**

```
{
   "@odata.type":"#Task.v1_4_2.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "Name":"Task task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "Message":"The task with id task4aac9e1e-df58-4fff-b781-52373fcb5699 has started.",
   "MessageId":"TaskEvent.1.0.1.TaskStarted",
   "MessageArgs":[
      "task4aac9e1e-df58-4fff-b781-52373fcb5699"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```

>**Sample response body \(HTTP 200 status\)**

```
{
   "error":{
      "@Message.ExtendedInfo":[
         {
            "MessageId":"Base.1.4.Success"
         }
      ],
      "code":"iLO.0.10.ExtendedInfo",
      "message":"See @Message.ExtendedInfo for more information."
   }
```


### Start update

| | |
|-------|-----------|
|<strong>Method</strong> | `POST` |
|<strong>URI</strong> |`/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate` |
|<strong>Description</strong> |This operation starts updating software or firmware components for which an update request has been created.<br>It is performed in the background as a Redfish task.<br><blockquote>IMPORTANT:<br>Before performing this operation, ensure that you have created an update request first. To know how to create an update request, see [Simple update](#Simple update).<br></blockquote>|
|<strong>Response code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

**Usage information** 
To know the progress of this action, perform HTTP `GET` on the [task monitor](#viewing-a-task-monitor) returned in the response header \(until the task is complete\).


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
 'https://{odim_host}:{port}/redfish/v1/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate'


```

> Sample request body

None

>**Sample response header \(HTTP 202 status\)**

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/task4aac9e1e-df58-4fff-b781-52373fcb5699
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes

```

>**Sample response body \(HTTP 202 status\)**

```
{
   "@odata.type":"#Task.v1_4_2.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "Name":"Task task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "Message":"The task with id task4aac9e1e-df58-4fff-b781-52373fcb5699 has started.",
   "MessageId":"TaskEvent.1.0.1.TaskStarted",
   "MessageArgs":[
      "task4aac9e1e-df58-4fff-b781-52373fcb5699"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```

>**Sample response body \(HTTP 200 status\)**

```
{
   "error":{
      "@Message.ExtendedInfo":[
         {
            "MessageId":"Base.1.4.Success"
         }
      ],
      "code":"iLO.0.10.ExtendedInfo",
      "message":"See @Message.ExtendedInfo for more information."
```




#  Host to fabric networking

Resource Aggregator for ODIM exposes Redfish APIs to view and manage simple fabrics. A fabric is a network topology consisting of entities such as interconnecting switches, zones, endpoints, and address pools. The Redfish `Fabrics` APIs allow you to create and remove these entities in a fabric.

When creating fabric entities, ensure to create them in the following order:

1.  Zone-specific address pools

2.  Address pools for zone of zones

3.  Zone of zones

4.  Endpoints

5.  Zone of endpoints


When deleting fabric entities, ensure to delete them in the following order:

1.  Zone of endpoints

2.  Endpoints

3.  Zone of zones

4.  Address pools for zone of zones

5.  Zone-specific address pools


<blockquote>
IMPORTANT:

Before using the `Fabrics` APIs, ensure that the fabric manager is installed, its plugin is deployed, and added into the Resource Aggregator for ODIM framework.

</blockquote>


**Supported endpoints**


|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/Fabrics|GET|`Login` |
|/redfish/v1/Fabrics/\{fabricId\}|GET|`Login` |
|/redfish/v1/Fabrics/\{fabricId\}/Switches|GET|`Login` |
|/redfish/v1/Fabrics/\{fabricId\}/Switches/\{switchId\}|GET|`Login` |
| /redfish/v1/Fabrics/\{fabricId\}/Switches/\{switchId\}/Ports<br> |GET|`Login` |
| /redfish/v1/Fabrics/\{fabricId\} /Switches/\{switchId\}/Ports/\{portid\}<br> |GET|`Login` |
|/redfish/v1/Fabrics/\{fabricId\}/Zones|GET, POST|`Login`, `ConfigureComponents` |
|/redfish/v1/Fabrics/\{fabricId\}/Zones/\{zoneId\}|GET, PATCH, DELETE|`Login`, `ConfigureComponents` |
|/redfish/v1/Fabrics/\{fabricId\}/AddressPools|GET, POST|`Login`, `ConfigureComponents` |
|/redfish/v1/Fabrics/\{fabricId\}/AddressPools/\{addresspoolid\}|GET, DELETE|`Login`, `ConfigureComponents` |
|/redfish/v1/Fabrics/\{fabricId\}/Endpoints|GET, POST|`Login`, `ConfigureComponents` |
|/redfish/v1/Fabrics/\{fabricId\}/Endpoints/\{endpointId\}|GET, DELETE|`Login`, `ConfigureComponents` |




##  Collection of fabrics

|||
|---------------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Fabrics` |
|**Description** |A collection of simple fabrics.|
|**Returns** |Links to the fabric instances.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics'

```

>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics",
   "Id":"FabricCollection",
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6"
      }
   ],
   "Members@odata.count":1,
   "Name":"Fabric Collection",
   "RedfishVersion":"1.8.0",
   "@odata.type":"#FabricCollection.FabricCollection"
}
```












## Single fabric

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}` |
|**Description** |Schema representing a specific fabric.|
|**Returns** |Links to various components contained in this fabric instance - address pools, endpoints, switches, and zones.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}'


```



>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5",
   "@odata.type":"#Fabric.v1_0_4.Fabric",
   "AddressPools":{ 
      "@odata.id":"/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/AddressPools"
   },
   "Description":"test",
   "Endpoints":{ 
      "@odata.id":"/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Endpoints"
   },
   "FabricType":"Ethernet",
   "Id":"f4d1578a-d16f-43f2-bb81-cd6db8866db5",
   "Name":"cfm-test",
   "Status":{ 
      "Health":"OK",
      "State":"Enabled"
   },
   "Switches":{ 
      "@odata.id":"/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches"
   },
   "Zones":{ 
      "@odata.id":"/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Zones"
   }
}
```







## Collection of switches

|||
|------------------|----------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Switches` |
|**Description** |A collection of switches located in this fabric.|
|**Returns** |Links to the switch instances.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Switches'

```

>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches",
   "Id":"SwitchCollection",
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches/fb7dc9fd-d0f1-474e-b849-77262f5d73b7"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches/e97a3f0b-cc89-40d8-af3f-9b9bdd793d73"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches/a4ca3161-db95-487d-a930-1b13dc697ed0"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches/bc95a9aa-8447-4b89-a99d-25235f7bae92"
      }
   ],
   "Members@odata.count":4,
   "Name":"Switch Collection",
   "@odata.type":"#SwitchCollection.SwitchCollection"
}
```











## Single switch

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Switches/{switchID}` |
|**Description** |JSON schema representing a particular fabric switch.|
|**Returns** |Details of this switch and links to its ports.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Switches/{switchID}'

```

>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches/fb7dc9fd-d0f1-474e-b849-77262f5d73b7",
   "@odata.type":"#Switch.v1_2_0.Switch",
   "Id":"fb7dc9fd-d0f1-474e-b849-77262f5d73b7",
   "Manufacturer":"Aruba",
   "Model":"Aruba 8325",
   "Name":"Switch_172.10.20.1",
   "PartNumber":"JL636A",
   "Ports":{ 
      "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches/fb7dc9fd-d0f1-474e-b849-77262f5d73b7/Ports"
   },
   "SerialNumber":"TW8BKM302H",
   "Status":{ 
      "Health":"Ok",
      "State":"Enabled"
   },
   "SwitchType":"Ethernet"
}
```




## Collection of ports

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/Fabrics/{fabricID}/Switches/{switchID}/Ports`` |
|**Description** |A collection of ports of this switch.|
|**Returns** |Links to the port instances.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Switches/{switchID}/Ports'

```

>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches/c1b0ac48-e003-4d70-a707-450b128977d9/Ports",
   "Id":"PortCollection",
   "Members":[ 
      "/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches/c1b0ac48-e003-4d70-a707-450b128977d9/Ports/0cb2ff96-b7a7-4627-a7b4-274d915f2524",
      "/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches/c1b0ac48-e003-4d70-a707-450b128977d9/Ports/54096ea1-cfb8-4a6c-b7a3-d6263db729a6",
      "/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches/c1b0ac48-e003-4d70-a707-450b128977d9/Ports/699b8f82-a6bf-47fa-a594-d73c95a8f81e",
      "/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches/c1b0ac48-e003-4d70-a707-450b128977d9/Ports/7dd6b8c6-de72-4499-98dc-568a16e28a88",
      "/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches/c1b0ac48-e003-4d70-a707-450b128977d9/Ports/9d097004-e034-4772-98c5-fa695688cc4d",
      "/redfish/v1/Fabrics/f4d1578a-d16f-43f2-bb81-cd6db8866db5/Switches/c1b0ac48-e003-4d70-a707-450b128977d9/Ports/b82e663b-6d1e-43c9-9d49-63b68b2a5b06"
   ],
   "Members@odata.count":6,
   "Name":"PortCollection",
   "@odata.type":"#PortCollection.PortCollection"
}
	
```








## Single port

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Switches/{switchID}/Ports/{portid}` |
|**Description** |JSON schema representing a specific switch port.|
|**Returns** |Properties of this port.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Switches/{switchID}/Ports/{portid}'

```

>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Switches/a4ca3161-db95-487d-a930-1b13dc697ed0/Ports/80b5f999-25e9-4b37-992c-de2f065ee0e3",
   "@odata.type":"#Port.v1_1_3.Port",
   "CurrentSpeedGbps":0,
   "Description":"single port",
   "Id":"80b5f999-25e9-4b37-992c-de2f065ee0e3",
   "Links":{ 
      "ConnectedPorts":[ 
         { 
            "@odata.id":"/redfish/v1/Systems/768f9da7-56fc-4f13-b6f8-a1cd241e2313:1/EthernetInterfaces/3"
         }
      ],
      "ConnectedSwitches":[ 

      ]
   },
   "MaxSpeedGbps":25,
   "Name":"1/1/3",
   "PortId":"1/1/3",
   "PortProtocol":"Ethernet",
   "PortType":"UpstreamPort",
   "Status":{ 
      "Health":"Ok",
      "State":"Enabled"
   },
   "Width":1
}
```









## Collection of address pools

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/Fabrics/{fabricID}/AddressPools`` |
|**Description** |A collection of address pools.|
|**Returns** |Links to the address pool instances.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/AddressPools'

```


>**Sample response body**

```
{
	"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools",
	"Id": "AddressPool Collection",
	"Members": [{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/54a6d41b-6ed2-460b-90c7-cc5fdd74e6ad"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/f936ba02-fa82-456b-a7d7-3d006228f63c"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/95dec77a-c393-4391-8943-29e3ce03c6ca"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/21d4c00e-0c7c-4af3-af76-fd66df5d5831"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/35698e79-a765-4052-86ed-e290d7b6fd01"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/612aeda9-5cca-4f51-b755-b90008467bad"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/470515c8-6089-4d97-ba4f-e7dabc9d7e6a"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/062f6464-4f0e-4a6b-bb6b-c1857bba1533"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/7b740372-2f88-46d8-af84-7b66fee87695"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/7a98eb5d-99e9-4924-b647-057e3ad772bf"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/9f33f532-0796-42fc-819c-7938a4d6de7c"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/AddressPools/3a821e9d-901a-4469-913b-917351d6eef7"
		}
	],
	"Members@odata.count": 12,
	"Name": "AddressPool Collection",
	"RedfishVersion": "1.8.0",
	"@odata.type": "#AddressPoolCollection.AddressPoolCollection"
}
	
```














## Single address pool

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/Fabrics/{fabricID}/AddressPools/{addresspoolid}`` |
|**Description** |JSON schema representing a specific address pool.|
|**Returns** |Properties of this address pool.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/AddressPools/{addresspoolid}'

```


>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/AddressPools/44c44b52-a784-48e5-9f26-b833d42cf455",
   "@odata.type":"#AddressPool.vxx.AddressPool",
   "BgpEvpn":{ 
      "BgpEvpnEviNumberLowerAddress":200,
      "BgpEvpnEviNumberUpperAddress":220
   },
   "Description":"",
   "ExternalBgp":{ 
      "EbgpAsNumberLowerAddress":65000,
      "EbgpAsNumberUpperAddress":65100
   },
   "IPv4":{ 
      "IPv4FabricLinkLowerAddress":"172.10.20.1",
      "IPv4FabricLinkUpperAddress":"172.10.20.10",
      "IPv4GatewayAddress":"",
      "IPv4HostLowerAddress":"",
      "IPv4HostUpperAddress":"",
      "IPv4LoopbackLowerAddress":"172.10.20.1",
      "IPv4LoopbackUpperAddress":"172.10.20.10",
      "NativeVlan":0,
      "VlanIdentifierLowerAddress":0,
      "VlanIdentifierUpperAddress":0
   },
   "Id":"44c44b52-a784-48e5-9f26-b833d42cf455",
   "Links":{ 
      "Zones":[ 

      ]
   },
   "MultiProtocolIbgp":{ 
      "MPIbgpAsNumberLowerAddress":1,
      "MPIbgpAsNumberUpperAddress":1
   },
   "Name":""
}
	
```






## Collection of endpoints


|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/Fabrics/{fabricID}/Endpoints`` |
|**Description** |A collection of fabric endpoints.|
|**Returns** |Links to the endpoint instances.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Endpoints'

```



>**Sample response body**

```
{
	"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/Endpoints",
 "@odata.type": "#EndpointCollection.EndpointCollection",
	"Members": [{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/Endpoints/f59d59f3-d2ec-4cc1-9255-f35b5b09a31a"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/Endpoints/952f0049-d639-4a00-820a-353a95564d37"
		},
		{
			"@odata.id": "/redfish/v1/Fabrics/edbd1da7-7e2c-4ad0-aa9e-930292619d5f/Endpoints/8a9c27b0-d4d7-4eef-a731-f92aedc49c69"
		}
	],
	"Members@odata.count": 3,
	"Name": "Endpoint Collection"
	
}
	
```





##  Single endpoint

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Endpoints/{endpointId}` |
|**Description** |JSON schema representing a specific fabric endpoint.|
|**Returns** |Details of this endpoint.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Endpoints/{endpointId}'

```

>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/8f6d8828-a21a-464f-abf9-ed062fa08cd9/Endpoints/b21f3e57-e46d-4a8e-92c8-8658edd107cb",
   "@odata.type":"#Endpoint.v1_3_1.Endpoint",
   "Description":"NK Endpoint Collection Description",
   "EndpointProtocol":"Ethernet",
   "Id":"b21f3e57-e46d-4a8e-92c8-8658edd107cb",
   "Links":{ 
      "ConnectedPorts":[ 
         { 
            "@odata.id":"/redfish/v1/Fabrics/8f6d8828-a21a-464f-abf9-ed062fa08cd9/Switches/3f4ac957-90ec-4676-91b9-90f9d78ef56c/Ports/7f708d4f-795d-401d-8bc1-c797fb3ce20b"
         }
      ],
      "Zones":[ 

      ]
   },
   "Name":"NK Endpoint Collection1",
   "Redundancy":[ 
      { 
         "MaxNumSupported":2,
         "MemberId":"Bond0",
         "MinNumNeeded":1,
         "Mode":"",
         "RedundancySet":[ 
            [ 

            ]
         ],
         "Status":{ 
            "Health":"",
            "State":""
         }
      }
   ]
}
```




## Collection of zones

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/Fabrics/{fabricID}/Zones`` |
|**Description** |A collection of fabric zones.|
|**Returns** |Links to the zone instances.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones'

```


>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones",
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/afa007b7-7ea6-4ab3-b5f1-ad37c8aebed7"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/85462c4f-028d-45d6-99d8-73c7889ea263"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/e906a6ab-18ef-4617-a151-420265e7d0f9"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/f310bf40-5163-4cbf-be5b-ac574fe87863"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/5c0b60a0-55f7-43f0-9b23-bfbba9130743"
      },
      { 
         "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/ac424042-7524-4c04-acbd-2d1af0a4832f"
      }
   ],
   "Members@odata.count":6,
   "Name":"Zone Collection",
   "@odata.type":"#ZoneCollection.ZoneCollection"
	
```






## Single zone


|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Zones/{zoneId}` |
|**Description** |JSON schema representing a specific fabric zone.|
|**Returns** |Details of this zone.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones/{zoneId}'

```

>**Sample response body**

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/f310bf40-5163-4cbf-be5b-ac574fe87863",
   "@odata.type":"#Zone.v1_3_0.Zone",
   "DefaultRoutingEnabled":false,
   "Description":"",
   "Id":"f310bf40-5163-4cbf-be5b-ac574fe87863",
   "Links":{ 
      "AddressPools":[ 
         { 
            "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/AddressPools/e8edcc87-81f9-43a9-b1ce-20a895a60014"
         }
      ],
      "ContainedByZones":[ 
         { 
            "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Zones/ac424042-7524-4c04-acbd-2d1af0a4832f"
         }
      ],
      "ContainsZones":[ 

      ],
      "Endpoints":[ 
         { 
            "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Endpoints/a9a41a01-ac4c-460d-923e-98ad9cc7abef"
         }
      ],
      "InvolvedSwitches":[ 

      ]
   },
   "Name":"NK Zone 1",
   "ZoneType":"ZoneOfEndpoints"
}
```









## Creating a zone-specific address pool

|||
|---------------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/AddressPools` |
|**Description** |This operation creates an address pool that can be used by a zone of endpoints.|
|**Returns** |<ul><li>Link to the created address pool in the `Location` header.</li><li>JSON schema representing the created address pool.</li></ul>|
|**Response code** |`201 Created` |
|**Authentication** |Yes|


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
   "Name":"FC 18 vlan_102 - AddressPools",
   "Description":"vlan_102",
   "IPv4":{
      "VlanIdentifierAddressRange":{
         "Lower":102,
         "Upper":102
      }
   },
   "BgpEvpn":{
      "GatewayIPAddressList":[
         "10.18.102.2/24",
         "10.18.102.3/24"
      ],
      "AnycastGatewayIPAddress":"10.18.102.1"
   }
}
'
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/AddressPools'

```







>**Sample request body**

```
{
   "Name":"FC 18 vlan_102 - AddressPools",
   "Description":"vlan_102",
   "IPv4":{
      "VlanIdentifierAddressRange":{
         "Lower":102,
         "Upper":102
      }
   },
   "BgpEvpn":{
      "GatewayIPAddressList":[
         "10.18.102.2/24",
         "10.18.102.3/24"
      ],
      "AnycastGatewayIPAddress":"10.18.102.1"
   }
}
```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Name|String \(optional\)<br> |Name for the address pool.|
|Description|String \(optional\)<br> |Description for the address pool.|
|IPv4\{| \(required\)<br> | |
|VlanIdentifierAddressRange\{| \(optional\)<br> | A single VLAN to assign on the ports or lags.<br> |
|Lower|Integer \(required\)<br> |VLAN lower address.|
|Upper\}\}|Integer \(required\)<br> |VLAN upper address.<br>**NOTE:**<br> `Lower` and `Upper` must have the same value. Ensure that IP range is accurate and it does not overlap with other pools.|
|BgpEvpn\{| \(required\)<br> | |
|GatewayIPAddressList|Array \(required\)<br> | IP pool to assign IPv4 address to the IP interface for VLAN per switch.<br> |
|AnycastGatewayIPAddress\}|String \(required\)<br> | A single active gateway IP address for the IP interface.<br> |
| | | |



>**Sample response header**

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/e2ec196d-4b55-44b3-b928-8273de9fb8bf
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Thu, 14 May 2020 16:18:54 GMT
Transfer-Encoding:chunked

```

>- **Sample response body**

```
{
  "@odata.id": "/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/e2ec196d-4b55-44b3-b928-8273de9fb8bf",
  "@odata.type": "#AddressPool.vxx.AddressPool",
  "BgpEvpn": {
    "AnycastGatewayIPAddress": "10.18.102.1",
    "AnycastGatewayMACAddress": "",
    "GatewayIPAddressList": [
      "10.18.102.2/24",
      "10.18.102.3/24"
    ],
    "RouteDistinguisherList": "",
    "RouteTargetList": [
      
    ]
  },
  "Description": "vlan_102",
  "Ebgp": {
    
  },
  "IPv4": {
    "EbgpAddressRange": {
      "Lower": "",
      "Upper": ""
    },
    "FabricLinkAddressRange": {
      "Lower": "",
      "Upper": ""
    },
    "IbgpAddressRange": {
      "Lower": "",
      "Upper": ""
    },
    "LoopbackAddressRange": {
      "Lower": "",
      "Upper": ""
    },
    "NativeVlan": 0,
    "VlanIdentifierAddressRange": {
      "Lower": 102,
      "Upper": 102
    }
  },
  "Id": "e2ec196d-4b55-44b3-b928-8273de9fb8bf",
  "Links": {
    "Zones": [
      
    ]
  },
  "Name": "FC 18 vlan_102 - AddressPools"
}

```






## Creating an address pool for zone of zones

|||
|---------------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/AddressPools` |
|**Description** |This operation creates an address pool for a zone of zones in a specific fabric.|
|**Returns** |- Link to the created address pool in the `Location` header.<br />- JSON schema representing the created address pool.|
|**Response code** | `201 Created` |
|**Authentication** |Yes|


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
  "Name": "AddressPool for ZoneOfZones - Vlan3002",
  "IPv4": {  
    "VlanIdentifierAddressRange": {
        "Lower": 3002,
        "Upper": 3002
    },
    "IbgpAddressRange": {
              "Lower": "192.12.1.10",
              "Upper": "192.12.1.15"
    },
    "EbgpAddressRange": {
              "Lower": "172.12.1.10",
              "Upper": "172.12.1.15"
    }
  },
  "Ebgp": {
    "AsNumberRange": {
              "Lower": 65120,
              "Upper": 65125
    }
  },
  "BgpEvpn": {
    "RouteDistinguisherList": ["65002:102"],  
    "RouteTargetList": ["65002:102", "65002:102"],
    "GatewayIPAddressList": ["192.168.18.122/31", "192.168.18.123/31"]
  }
}'
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/AddressPools'

```







>**Sample request body**

```
{
  "Name": "AddressPool for ZoneOfZones - Vlan3002",
  "IPv4": {  
    "VlanIdentifierAddressRange": {
        "Lower": 3002,
        "Upper": 3002
    },
    "IbgpAddressRange": {
              "Lower": "192.12.1.10",
              "Upper": "192.12.1.15"
    },
    "EbgpAddressRange": {
              "Lower": "172.12.1.10",
              "Upper": "172.12.1.15"
    }
  },
  "Ebgp": {
    "AsNumberRange": {
              "Lower": 65120,
              "Upper": 65125
    }
  },
  "BgpEvpn": {
    "RouteDistinguisherList": ["65002:102"],  
    "RouteTargetList": ["65002:102", "65002:102"],
    "GatewayIPAddressList": ["192.168.18.122/31", "192.168.18.123/31"]
  }
}

```

**Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Name|String|Name for the address pool.|
|Description|String \(optional\)<br> |Description for the address pool.|
|IPv4\{| \(required\)<br> | |
|VlanIdentifierAddressRange\{| \(required\)<br> | A single VLAN \(virtual LAN\) used for creating the IP interface for the user Virtual Routing and Forwarding \(VRF\).<br> |
|Lower|Integer \(required\)<br> |VLAN lower address|
|Upper\}|Integer \(required\)<br> |VLAN upper address|
|IbgpAddressRange\{| \(required\)<br> | IPv4 address used as the Router Id for the VRF per switch.<br> |
|Lower|String \(required\)<br> |IPv4 lower address|
|Upper\}|String \(required\)<br> |IPv4 upper address|
|EbgpAddressRange\{| \(optional\)<br> |External neighbor IPv4 addresses.|
|Lower|String \(required\)<br> |IPv4 lower address|
|Upper\} \}|String \(required\)<br> |IPv4 upper address|
|Ebgp\{| \(optional\)<br> | |
|AsNumberRange\{| \(optional\)<br> |External neighbor ASN.<br>**NOTE:**<br> `EbgpAddressRange` and `AsNumberRange` values should be a matching sequence and should be of same length.|
|Lower|Integer \(optional\)<br> | |
|Upper\} \}|Integer \(optional\)<br> | |
|BgpEvpn\{| \(required\)<br> | |
|RouteDistinguisherList|Array \(required\)<br> | Single route distinguisher value for the VRF.<br> |
|RouteTargetList|Array \(optional\)<br> | Route targets. By default, the route targets will be configured as both import and export.<br> |
|GatewayIPAddressList\}|Array \(required\)<br> | IP pool to assign IPv4 address to the IP interface used by the VRF per switch.<br> |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Thu, 14 May 2020 16:18:58 GMT
Transfer-Encoding:chunked

```

>**Sample response body**

```
{
   "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d",
   "@odata.type":"#AddressPool.vxx.AddressPool",
   "BgpEvpn":{
      "AnycastGatewayIPAddress":"",
      "AnycastGatewayMACAddress":"",
      "GatewayIPAddressList":[
         "192.168.18.122/31",
         "192.168.18.123/31"
      ],
      "RouteDistinguisherList":[
         "65002:102"
      ],
      "RouteTargetList":[
         "65002:102",
         "65002:102"
      ]
   },
   "Description":"",
   "Ebgp":{
      "AsNumberRange":{
         "Lower":65120,
         "Upper":65125
      }
   },
   "IPv4":{
      "EbgpAddressRange":{
         "Lower":"172.12.1.10",
         "Upper":"172.12.1.15"
      },
      "FabricLinkAddressRange":{
         "Lower":"",
         "Upper":""
      },
      "IbgpAddressRange":{
         "Lower":"192.12.1.10",
         "Upper":"192.12.1.15"
      },
      "LoopbackAddressRange":{
         "Lower":"",
         "Upper":""
      },
      "NativeVlan":0,
      "VlanIdentifierAddressRange":{
         "Lower":3002,
         "Upper":3002
      }
   },
   "Id":"84766158-cbac-4f69-8ed5-fa5f2b331b9d",
   "Links":{
      "Zones":[

      ]
   },
   "Name":"AddressPool for ZoneOfZones - Vlan3002"
}
```


##  Adding a zone of zones



|||
|---------------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Zones` |
|**Description** |This operation creates an empty container zone for all the other zones in a specific fabric. To assign address pools, endpoints, other zones, or switches to this zone, perform HTTP `PATCH` on this zone. See [Updating a Zone](#updating-a-zone).|
|**Returns** |JSON schema representing the created zone.|
|**Response code** |`201 Created` |
|**Authentication** |Yes|


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
   "Name":"Fabric Zone of Zones:",
   "ZoneType":"ZoneOfZones",
   "Links":{
      "AddressPools":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d"
         }
      ]
   }
}
'
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones'

```

>**Sample request body**

```
{
   "Name":"Fabric Zone of Zones:",
   "ZoneType":"ZoneOfZones",
   "Links":{
      "AddressPools":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d"
         }
      ]
   }
}
```

**Request parameters**

|Parameter|Value|Description|
|---------|-----|-----------|
|Name|String \(optional\)<br> |Name for the zone.|
|Description|String \(optional\)<br> |Description for the zone.|
|ZoneType|String|The type of the zone to be created. Options include: `ZoneofZones` and `ZoneofEndpoints`<br> The type of the zone for a zone of zones is `ZoneofZones`.<br> |
|Links\{| \(optional\)<br> | |
|AddressPools|Array \(optional\)<br> | `AddressPool` links supported for the Zone of Zones \(`AddressPool` links created for `ZoneofZones`\).<br> |


>**Sample response header**

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/a2dc8760-ea05-4cab-8f95-866c1c380f98
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Thu, 14 May 2020 16:19:00 GMT
Transfer-Encoding:chunked
```

>**Sample response body**

```
{
   "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/a2dc8760-ea05-4cab-8f95-866c1c380f98",
   "@odata.type":"#Zone.v1_3_0.Zone",
   "Description":"",
   "Id":"a2dc8760-ea05-4cab-8f95-866c1c380f98",
   "Links":{
      "AddressPools":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d"
         }
      ],
      "ContainedByZones":[

      ],
      "ContainsZones":[

      ],
      "Endpoints":[

      ],
      "InvolvedSwitches":[

      ]
   },
   "Name":"Fabric Zone of Zones:",
   "ZoneType":"ZoneOfZones"
}
```


## Adding an endpoint


|||
|---------------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Endpoints` |
|**Description** |This operation creates an endpoint in a specific fabric.|
|**Returns** | <ul><li>Link to the created endpoint in the `Location` header.</li><li>JSON schema representing the created endpoint.</li></ul>|
|**Response code** | `201 Created` |
|**Authentication** |Yes|


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
   "Name":"Host 2 Endpoint 1 Collection",
   "Description":"Host 2 Endpoint 1 Collection Description",
   "Redundancy":[
      {
         "Mode":"Sharing",
         "RedundancySet":[
            {
               "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Switches/81f5ed9a-a4a1-4383-a450-a7f98b792ca2/Ports/29f077b0-e7a5-495f-a3d2-643937f600de"
            },
            {
               "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Switches/b7ba4ece-1716-4d8d-af2c-7aea1682bf91/Ports/62a32f83-c7b1-4cb7-9b47-2f444588d29b"
            }
         ]
      }
   ]
}
'
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Endpoints'

```

>**Sample request body** for a single endpoint

```
{
   "Name":"NK Endpoint Collection",
   "Description":"NK Endpoint Collection Description",
   "Links":{
      "ConnectedPorts":[
         {
            "@odata.id":"/redfish/v1/Fabrics/113a30e3-f312-4221-8f7f-49943c5ff07d/Switches/f4a37f55-be1e-400b-93be-7d7c0afd4cbd/Ports/0d22b201-30d5-43e8-90ab-277c87624c05"
         }
      ]
   }
}
```

>**Sample request body** for a redundant endpoint

```
{
   "Name":"Host 2 Endpoint 1 Collection",
   "Description":"Host 2 Endpoint 1 Collection Description",
   "Redundancy":[
      {
         "Mode":"Sharing",
         "RedundancySet":[
            {
               "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Switches/81f5ed9a-a4a1-4383-a450-a7f98b792ca2/Ports/29f077b0-e7a5-495f-a3d2-643937f600de"
            },
            {
               "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Switches/b7ba4ece-1716-4d8d-af2c-7aea1682bf91/Ports/62a32f83-c7b1-4cb7-9b47-2f444588d29b"
            }
         ]
      }
   ]
}
```

**Request parameters**

|Parameter|Value|Description|
|---------|-----|-----------|
|Name|String \(optional\)<br> |Name for the endpoint.|
|Description|String \(optional\)<br> |Description for the endpoint.|
|Links\{| \(required\)<br> | |
|ConnectedPorts|Array \(required\)<br> | Switch port connected to the switch.<br>  <br> |
|Zones\}|Array \(optional\)<br> | Endpoint is part of `ZoneofEndpoints` and only one zone is permitted in the zones list.<br> |
|Redundancy\[|Array| |
|Mode|String|Redundancy mode.|
|RedundancySet\]|Array| Set of redundancy ports connected to the switches.<br> |

>**Sample response header**

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Endpoints/fe34aff2-e81f-4167-a0c3-9bf5a67e2a97
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Thu, 14 May 2020 16:19:02 GMT
Transfer-Encoding:chunked

```

>**Sample response body**

```
{
   "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Endpoints/fe34aff2-e81f-4167-a0c3-9bf5a67e2a97",
   "@odata.type":"#Endpoint.v1_3_1.Endpoint",
   "Description":"Host 2 Endpoint 1 Collection Description",
   "EndpointProtocol":"Ethernet",
   "Id":"fe34aff2-e81f-4167-a0c3-9bf5a67e2a97",
   "Links":{
      "ConnectedPorts":[

      ],
      "Zones":[

      ]
   },
   "Name":"Host 2 Endpoint 1 Collection",
   "Redundancy":[
      {
         "MaxNumSupported":2,
         "MemberId":"Bond0",
         "MinNumNeeded":1,
         "Mode":"Sharing",
         "RedundancySet":[
            {
               "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Switches/81f5ed9a-a4a1-4383-a450-a7f98b792ca2/Ports/29f077b0-e7a5-495f-a3d2-643937f600de"
            },
            {
               "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Switches/b7ba4ece-1716-4d8d-af2c-7aea1682bf91/Ports/62a32f83-c7b1-4cb7-9b47-2f444588d29b"
            }
         ],
         "Status":{
            "Health":"",
            "State":""
         }
      }
   ]
}
```






## Creating a zone of endpoints

|||
|---------------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Zones` |
|**Description** |This operation creates a zone of endpoints in a specific fabric.<br>**NOTE:**<br> Ensure that the endpoints are created first before assigning them to the zone of endpoints.|
|**Returns** |JSON schema representing the created zone.|
|**Response code** | `201 Created` |
|**Authentication** |Yes|


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
   "Name":"Fabric ZoneofEndpoint",
   "ZoneType":"ZoneOfEndpoints",
   "Links":{
      "ContainedByZones":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/a2dc8760-ea05-4cab-8f95-866c1c380f98"
         }
      ],
      "AddressPools":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/e2ec196d-4b55-44b3-b928-8273de9fb8bf"
         }
      ],
      "Endpoints":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Endpoints/fe34aff2-e81f-4167-a0c3-9bf5a67e2a97"
         }
      ]
   }
}
'
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones'

```






>**Sample request body**

```
{
   "Name":"Fabric ZoneofEndpoint",
   "ZoneType":"ZoneOfEndpoints",
   "Links":{
      "ContainedByZones":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/a2dc8760-ea05-4cab-8f95-866c1c380f98"
         }
      ],
      "AddressPools":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/e2ec196d-4b55-44b3-b928-8273de9fb8bf"
         }
      ],
      "Endpoints":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Endpoints/fe34aff2-e81f-4167-a0c3-9bf5a67e2a97"
         }
      ]
   }
}
```


**Request parameters**

|Parameter|Value|Description|
|---------|-----|-----------|
|Name|String \(optional\)<br> |The name for the zone.|
|Description|String \(optional\)<br> |The description for the zone.|
|DefaultRoutingEnabled|Boolean \(required\)<br> |Set to `false`.|
|ZoneType|String \(required\)<br> |The type of the zone to be created. Options include: `ZoneofZones`and `ZoneofEndpoints`<br> The type of the zone for a zone of endpoints is `ZoneofEndpoints`.<br> |
|Links\{|Object \(optional\)<br> |Contains references to other resources that are related to the zone.|
|ContainedByZones \[\{|Array \(optional\)<br> |Represents an array of `ZoneofZones` for the zone being created \(applicable when creating a zone of endpoints\).|
|@odata.id \}\]|String|Link to a Zone of zones.|
|AddressPools \[\{|Array \(optional\)<br> |Represents an array of address pools linked with the zone \(zone-specific address pools\).|
|@odata.id \}\]|String|Link to an address pool.|
|Endpoints \[\{|Array \(optional\)<br> |Represents an array of endpoints to be included in the zone.|
|@odata.id \}\]|String|Link to an endpoint.|



>**Sample response header**

```
HTTP/1.1 201 Created
Allow: "GET", "PUT", "POST", "PATCH", "DELETE"
Cache-Control: no-cache
Connection: keep-alive
Content-Type: application/json; charset=utf-8
Location: /redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/06d344bb-cce1-4b0c-8414-6f6df1ea373f
Odata-Version: 4.0
X-Frame-Options: sameorigin
Date: Thu, 14 May 2020 16:19:37 GMT
Transfer-Encoding: chunked

```

>**Sample response body**

```
{
   "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/06d344bb-cce1-4b0c-8414-6f6df1ea373f",
   "@odata.type":"#Zone.v1_3_0.Zone",
   "Description":"",
   "Id":"06d344bb-cce1-4b0c-8414-6f6df1ea373f",
   "Links":{
      "AddressPools":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/e2ec196d-4b55-44b3-b928-8273de9fb8bf"
         }
      ],
      "ContainedByZones":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/a2dc8760-ea05-4cab-8f95-866c1c380f98"
         }
      ],
      "ContainsZones":[

      ],
      "Endpoints":[
         {
            "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Endpoints/fe34aff2-e81f-4167-a0c3-9bf5a67e2a97"
         }
      ],
      "InvolvedSwitches":[

      ]
   },
   "Name":"Fabric ZoneofEndpoint",
   "ZoneType":"ZoneOfEndpoints"
}
```



##  Updating a zone

|||
|---------------|---------------|
|**Method** |`PATCH` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Zones/{zoneId}` |
|**Description** |This operation assigns or unassigns a collection of endpoints, address pools, zone of zones, or switches to a zone of endpoints or a collection of address pools to a zone of zones in a specific fabric.|
|**Returns** |JSON schema representing an updated zone.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i -X PATCH \
   -H "X-Auth-Token:{X-Auth-Token}" \
    -d \
'{
	"Links": {
		"Endpoints": [{
			"@odata.id": "/redfish/v1/Fabrics/d76f4c66-aa60-4693-bea1-feac44fb9f81/Endpoints/a9d7f926-fc3c-465f-9724-928ba9becdb2"
		}]
	}
}
'
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones/{zoneId}'

```




>**Sample request body** \(assigning links for a zone of endpoints\)

```
{
	"Links": {
		"Endpoints": [{
			"@odata.id": "/redfish/v1/Fabrics/d76f4c66-aa60-4693-bea1-feac44fb9f81/Endpoints/a9d7f926-fc3c-465f-9724-928ba9becdb2"
		}]
	}
}
```

>**Sample request body** \(unassigning links for a zone of endpoints\)

```
{
	"Links": {
		"Endpoints": []
	}
}
```

>**Sample request body** \(assigning links for a zone of zone\)

```
{
   "Links":{
      "AddressPools":[
        "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d"
      ]
   }
}
```

>**Sample request body** \(unassigning links for a zone of zone\)

```
{
   "Links":{
      "AddressPools":[

      ]
   }
}
```

>**Sample response body** \(assigned links in a zone of endpoints\)

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/143476dc-0ac1-4352-96f3-e0782aeed84a/Zones/57c325f0-eda4-4754-b8da-826d5e266c04",
   "@odata.type":"#Zone.v1_3_0.Zone",
   "Description":"NK Zone Collection Description",
   "Id":"57c325f0-eda4-4754-b8da-826d5e266c04",
   "Links":{ 
      "AddressPools":[ 

      ],
      "ContainedByZones":[ 

      ],
      "ContainsZones":[ 

      ],
      "Endpoints":[ 
         { 
            "@odata.id":"/redfish/v1/Fabrics/77205057-3ef1-4c18-945c-2bf7893ea4a6/Endpoints/e8edcc87-81f9-43a9-b1ce-20a895a60014"
         }
      ],
      "InvolvedSwitches":[ 

      ]
   },
   "Name":"SS Zone Collection default",
   "ZoneType":"ZoneOfEndpoints"
}
```


>**Sample response body** \(unassigned links in a zone of endpoints\)

```
{ 
   "@odata.id":"/redfish/v1/Fabrics/143476dc-0ac1-4352-96f3-e0782aeed84a/Zones/57c325f0-eda4-4754-b8da-826d5e266c04",
   "@odata.type":"#Zone.v1_3_0.Zone",
   "Description":"NK Zone Collection Description",
   "Id":"57c325f0-eda4-4754-b8da-826d5e266c04",
   "Links":{ 
      "AddressPools":[ 

      ],
      "ContainedByZones":[ 

      ],
      "ContainsZones":[ 

      ],
      "Endpoints":[ 

      ],
      "InvolvedSwitches":[ 

      ]
   },
   "Name":"SS Zone Collection default",
   "ZoneType":"ZoneOfEndpoints"
}
```


## Deleting a zone

|||
|---------------|---------------|
|**Method** |`DELETE` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Zones/{zoneId}` |
|**Description** |This operation deletes a zone in a specific fabric.<br>**NOTE:**<br> If you delete a non-empty zone \(a zone which contains links to address pools, other zones, endpoints, or switches\), you will receive an HTTP `400` error. Before attempting to delete, unassign all links in the zone. See [updating a zone](#updating-a-zone).|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i -X DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/Zones/{zoneId}'

```







## Deleting an endpoint

|||
|---------------|---------------|
|**Method** |`DELETE` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/Endpoints/{endpointId}` |
|**Description** |This operation deletes an endpoint in a specific fabric.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odim_hosts}:{port}/redfish/v1/Fabrics/{fabricID}/Endpoints/{endpointId}'

```






## Deleting an address pool

|||
|---------------|---------------|
|**Method** |`DELETE` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}/AddressPools/{addresspoolid}` |
|**Description** |This operation deletes an address pool in a specific fabric.<br>**NOTE:**<br> If you delete an address pool that is being used in any zone, you will receive an HTTP `400` error. Before attempting to delete, ensure that the address pool you want to delete is not present in any zone. To get the list of address pools in a zone, see links to `addresspools` in the sample response body for a [single zone](#single-zone).|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i -X DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/AddressPools/{addresspo'

```


# Tasks

A task represents an operation that takes more time than a user typically wants to wait and is carried out asynchronously.

An example of a task is resetting an aggregate of servers. Resetting all the servers in a group is a time-consuming operation; the user waiting for the result would be blocked from performing other operations. Resource Aggregator for ODIM creates Redfish tasks for such long-duration operations and exposes Redfish `TaskService` APIs and `Task monitor` API. Use these APIs to manage and monitor the tasks until their completion, while performing other operations.

**Supported endpoints**

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/TaskService|GET|`Login` |
|/redfish/v1/TaskService/Tasks|GET|`Login` |
|/redfish/v1/TaskService/Tasks/\{taskId\}|GET, DELETE|`Login`, `ConfigureManager` |
| /redfish/v1/ TaskService/Tasks/\{taskId\}/SubTasks<br> |GET|`Login` |
| /redfish/v1/ TaskService/Tasks/\{taskId\}/SubTasks/ \{subTaskId\}<br> |GET|`Login` |
|/taskmon/\{taskId\}|GET|`Login` |


>**NOTE:**
To view the tasks and the task monitor, ensure that the user has `Login` privilege at the minimum. If you access these endpoints without necessary privileges, you will receive an HTTP `403 Forbidden` error.




##  Viewing the task service root

|||
|-----------|----------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/TaskService` |
|**Description** |This endpoint retrieves JSON schema for the Redfish `TaskService` root.|
|**Returns** |<ul><li> Links to tasks</li><li>Properties of `TaskService`.<br> Following are a few important properties of `TaskService` returned in the JSON response:<br><ul><li>`CompletedTaskOverWritePolicy` : This property indicates the overwrite policy for completed tasks and is set to `oldest` by default - Older completed tasks will be removed automatically.</li><li>`LifeCycleEventOnTaskStateChange`: This property indicates if the task state change event will be sent to the clients who have subscribed to it. It is set to `true` by default.</li></ul></li></ul> |
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/TaskService'


```

 

>**Sample response header** 

```
Allow:GET
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Date:Sun,17 May 2020 15:11:12 GMT+5m 13s
Link:</redfish/v1/SchemaStore/en/TaskService.json>; rel=describedby
Odata-Version:4.0
X-Frame-Options:sameorigin
Transfer-Encoding":chunked
```

 

>**Sample response body** 

```
{
   "@odata.type":"#TaskService.v1_1_4.TaskService",
   "@odata.id":"/redfish/v1/TaskService",
   "@odata.context":"/redfish/v1/$metadata#TaskService.TaskService",
   "Description":"TaskService",
   "Id":"TaskService",
   "Name":"TaskService",
   "CompletedTaskOverWritePolicy":"Oldest",
   "DateTime":"2020-04-17T09:42:04.547136227Z",
   "LifeCycleEventOnTaskStateChange":true,
   "ServiceEnabled":true,
   "Status":{
      "Health":"OK",
      "HealthRollup":"OK",
      "Oem":{

      },
      "State":"Enabled"
   },
   "Tasks":{
      "@odata.id":"/redfish/v1/TaskService/Tasks"
   }
}
```



## Viewing a collection of tasks

|||
|-----------|----------|
|**Method** |**GET** |
|**URI** |`/redfish/v1/TaskService/Tasks` |
|**Description** |This endpoint retrieves a list of tasks scheduled by or being executed by Redfish `TaskService`.<br>**NOTE:**<br>Only an admin or a user with `ConfigureUsers` privilege can view all the running and scheduled tasks in Resource Aggregator for ODIM at any given time. Other users can view tasks created only for their operations with `Login` privilege.<br></blockquote>|
|**Returns** |A list of task endpoints with task Ids.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/TaskService/Tasks'

```

>**Sample response body** 

```
{ 
   "@odata.type":"#TaskCollection.TaskCollection",
   "@odata.id":"/redfish/v1/TaskService/Tasks",
   "@odata.context":"/redfish/v1/$metadata#TaskCollection.TaskCollection",
   "Name":"Task Collection",
   "Members@odata.count":3,
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/TaskService/Tasks/taskc8cf2e2e-6cb2-4e24-8512-247fa5d606b0"
      },
      { 
         "@odata.id":"/redfish/v1/TaskService/Tasks/taskc15aca5a-30a6-4618-adca-c25c889dc409"
      },
      { 
         "@odata.id":"/redfish/v1/TaskService/Tasks/task38d6df20-989f-4c05-ad0e-5939774bab7c"
      }
   ]
}
```



## Viewing information about a specific task

|||
|-----------|----------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/TaskService/Tasks/{TaskID}` |
|**Description** |This endpoint retrieves information about a specific task scheduled by or being executed by Redfish `TaskService`.|
|**Returns** |JSON schema having the details of this task - task Id, name, state of the task, start time and end time of this task, completion percentage, URI of the task monitor associated with this task, subtasks if any. The sample response body given in this section is a JSON response for a task which adds a server.<br> |
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/TaskService/Tasks/{TaskID}'


```


>**Sample response body** 

```
{
   "@odata.type":"#Task.v1_5_0.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task2e4b6684-5c6b-4872-bb64-72cf27f3a78f",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task2e4b6684-5c6b-4872-bb64-72cf27f3a78f",
   "Name":"Task task2e4b6684-5c6b-4872-bb64-72cf27f3a78f",
   "TaskState":"Completed",
   "StartTime":"2021-01-21T07:09:03.366954469Z",
   "EndTime":"2021-01-21T07:10:30.311241695Z",
   "TaskStatus":"OK",
   "SubTasks":"",
   "TaskMonitor":"/taskmon/task2e4b6684-5c6b-4872-bb64-72cf27f3a78f",
   "PercentComplete":100,
   "Payload":{
      "HttpHeaders":[
         "Transfer-Encoding: chunked",
         "Cache-Control: no-cache",
         "Connection: keep-alive",
         "Content-type: application/json; charset=utf-8",
         "Link: </redfish/v1/AggregationService/AggregationSources/7b08ecbd-d23e-4dd5-ad99-58ac2be7576d:1/>; rel=describedby",
         "Location: /redfish/v1/AggregationService/AggregationSources/7b08ecbd-d23e-4dd5-ad99-58ac2be7576d:1",
         "OData-Version: 4.0"
      ],
      "HttpOperation":"POST",
      "JsonBody":"{\"HostName\":\"10.24.0.4\",\"Links\":{\"ConnectionMethod\":{\"@odata.id\":\"/redfish/v1/AggregationService/ConnectionMethods/c31a079c-4b69-4b78-b7d5-41d64bed8ea8\"}},\"Password\":\"HP1nvent\",\"UserName\":\"admin\"}",
      "TargetUri":"/redfish/v1/AggregationService/AggregationSources"
   },
   "Messages":[
      
   ]
}
```




##  Viewing a task monitor

|||
|-----------|----------|
|**Method** | `GET` |
|**URI** |`/taskmon/{TaskID}` |
|**Description** |This endpoint retrieves the task monitor associated with a specific task. A task monitor allows for polling a specific task for its completion. Perform `GET` on a task monitor URI to view the progress of a specific task \(until it is complete\).|
|**Returns** |<ul><li>Details of the task and its progress in the JSON response such as:<br> Link to the task,<br>Id of the task,<br>Task state and status,<br>Percentage of completion,<br>Start time and end time,<br>Link to subtasks \(if any\).<br>To know the status of a subtask, perform `GET` on the respective subtask link.<br>**NOTE:**<br><ul><li>Note down the task Id. If the task completes with an error, it is required to know which subtask has failed.To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`.</li><li>`EndTime` of an ongoing task has `0001-01-01T00:00:00Z` as value, which is equivalent to zero time stamp value. It is updated only after the completion of the task.</li></ul></li><li>On failure, an error message. See "Sample error response".<br> To get the list of subtasks, perform `GET` on the task URI having the Id of the failed task. To know which subtasks have failed, perform `GET` on subtask links individually.</li><li>On successful completion, result of the operation carried out by the task. See "Sample response body \(completed task\)".</li></ul>|
|**Response code** | <ul><li>`202 Accepted` until the task is complete.</li><li>`200 OK`, `201 Created` on success.</li></ul>|
|**Authentication** |Yes|


>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
 'https://{odimra_host}:{port}/taskmon/{TaskID}'

```


>**Sample response header**

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/taskfbd5cdb0-5d33-4ad4-8682-cab90534ba70
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes
```


>**Sample response body** \(ongoing task\)

```
{
   "@odata.type":"#Task.v1_4_2a.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/taskfbd5cdb0-5d33-4ad4-8682-cab90534ba70",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"taskfbd5cdb0-5d33-4ad4-8682-cab90534ba70",
   "Name":"Task taskfbd5cdb0-5d33-4ad4-8682-cab90534ba70",
   "Message":"The task with id taskfbd5cdb0-5d33-4ad4-8682-cab90534ba70 has started.",
   "MessageId":"TaskEvent.1.0.1.TaskStarted",
   "MessageArgs":[
      "taskfbd5cdb0-5d33-4ad4-8682-cab90534ba70"
   ],
   "NumberOfArgs":1,
   "Severity":"OK",
   "TaskState":"Running",
   "StartTime":"2020-04-17T09:39:22.713860589Z",
   "EndTime":"0001-01-01T00:00:00Z",
   "TaskStatus":"OK",
   "SubTasks":"/redfish/v1/Tasks/taskfbd5cdb0-5d33-4ad4-8682-cab90534ba70/SubTasks",
   "TaskMonitor":"/taskmon/taskfbd5cdb0-5d33-4ad4-8682-cab90534ba70",
   "PercentComplete":8,
   "Payload":{
      "HttpHeaders":null,
      "HttpOperation":"",
      "JsonBody":null,
      "TargetUri":""
   },
   "Messages":[

   ]
}
```

>**Sample response body** \(completed task\)

```
{
"code": "Base.1.6.1.Success",
"message": "Request completed successfully."
}
```

>  Sample error response

```
{ 
   "error":{ 
      "code":"Base.1.6.1.GeneralError",
      "message":"one or more of the reset actions failed, check sub tasks for more info."
   }
```



## Deleting a task

|||
|-----------|----------|
|**Method** | `DELETE` |
|**URI** |`/redfish/v1/TaskService/Tasks/{TaskID}` |
|**Description** |This operation deletes a specific task. Deleting a running task aborts the operation being carried out.<br>**NOTE:**<br> Only a user having `ConfigureComponents` privilege is authorized to delete a task. If you do not have the necessary privileges, you will receive an HTTP `403 Forbidden` error.|
|**Returns** |JSON schema representing the deleted task.|
|**Response code** |`204 No Content` |
|**Authentication** |Yes|


>**curl command**


```
curl -i DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/TaskService/Tasks/{TaskID}'

```




# Events

Resource Aggregator for ODIM offers an event interface that allows northbound clients to interact and receive notifications such as alerts and alarms from multiple resources, including Resource Aggregator for ODIM itself. It exposes Redfish `EventService` APIs for managing events.

An event is a way of asynchronously notifying the client of some significant state change or error condition, usually of a time critical nature.

Use these APIs to subscribe a northbound client to southbound events by creating a subscription entry in the service.

**Supported endpoints**

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/EventService|GET|`Login` |
|/redfish/v1/EventService/Subscriptions|GET, POST|`Login`, `ConfigureManager`, `ConfigureComponents` |
|/redfish/v1/EventService/Actions/EventService.SubmitTestEvent|POST|`ConfigureManager` |
|/redfish/v1/EventService/Subscriptions/\{subscriptionId\}|GET, DELETE|`Login`, `ConfigureManager`, `ConfigureSelf` |

>**Note:**
Before accessing these endpoints, ensure that the user has the required privileges. If you access these endpoints without necessary privileges, you will receive an HTTP `403 Forbidden` error.


## Viewing the event service root

|||
|----------|---------|
|**Method** | `GET` |
|**URI** |`redfish/v1/EventService` |
|**Description** |This endpoint retrieves JSON schema for the Redfish `EventService` root.|
|**Returns** |Properties for managing event subscriptions such as allowed event types and a link to the actual collection of subscriptions.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/EventService'


```

>**Sample response header** 

```
Allow:GET
Cache-Control:no-cache
Connection:Keep-alive
Content-Type:application/json; charset=utf-8
Link:/v1/SchemaStore/en/EventService.json>; rel=describedby
Odata-Version:4.0
X-Frame-Options:"sameorigin
Date:Fri,15 May 2020 10:10:15 GMT+5m 11s
Transfer-Encoding:chunked

```

>**Sample response body** 

```
{
   "@odata.context":"/redfish/v1/$metadata#EventService.EventService",
   "Id":"EventService",
   "@odata.id":"/redfish/v1/EventService",
   "@odata.type":"#EventService.v1_5_0.EventService",
   "Name":"EventService",
   "Description":"EventService",
   "Actions":{
      "#EventService.SubmitTestEvent":{
         "target":"/redfish/v1/EventService/Actions/EventService.SubmitTestEvent",
         "EventType@Redfish.AllowableValues":[
            "StatusChange",
            "ResourceUpdated",
            "ResourceAdded",
            "ResourceRemoved",
            "Alert"
         ]
      },
      "Oem":{

      }
   },
   "DeliveryRetryAttempts":3,
   "DeliveryRetryIntervalSeconds":60,
   "EventFormatTypes":[
      "Event"
   ],
   "EventTypesForSubscription":[
      "StatusChange",
      "ResourceUpdated",
      "ResourceAdded",
      "ResourceRemoved",
      "Alert"
   ],
   "RegistryPrefixes":[

   ],
   "ResourceTypes":[
      "ManagerAccount",
      "Switch",
      "EventService",
      "LogService",
      "ManagerNetworkProtocol",
      "ProcessorMetrics",
      "Task",
      "Drive",
      "EthernetInterface",
      "Protocol",
      "Redundancy",
      "Storage",
      "Bios",
      "EventDestination",
      "MessageRegistry",
      "Role",
      "Thermal",
      "VLanNetworkInterface",
      "IPAddresses",
      "NetworkInterface",
      "Volume",
      "Memory",
      "PCIeDevice",
      "PrivilegeRegistry",
      "Privileges",
      "Processor",
      "Sensor",
      "Endpoint",
      "MessageRegistryFile",
      "Session",
      "Assembly",
      "Fabric",
      "Job",
      "Manager",
      "MemoryDomain",
      "Port",
      "Power",
      "ProcessorCollection",
      "Zone",
      "ComputerSystem",
      "MemoryChunks",
      "MemoryMetrics",
      "NetworkDeviceFunction",
      "NetworkPort",
      "PCIeFunction",
      "AddressPool",
      "HostInterface",
      "NetworkAdapter",
      "PhysicalContext",
      "SecureBoot",
      "SerialInterface",
      "AccelerationFunction",
      "Chassis",
      "JobService",
      "Message",
      "Event",
      "LogEntry",
      "PCIeSlots",
      "Resource",
      "BootOption"
   ],
   "ServiceEnabled":true,
   "Status":{
      "Health":"OK",
      "HealthRollup":"OK",
      "State":"Enabled"
   },
   "SubordinateResourcesSupported":true,
   "Subscriptions":{
      "@odata.id":"/redfish/v1/EventService/Subscriptions"
   },
   "Oem":{

   }
}
```



### Creating an event subscription

|||
|-----------|-----------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/EventService/Subscriptions` |
|**Description**| This endpoint subscribes a northbound client to events originating from a set of resources \(southbound devices, managers, Resource Aggregator for ODIM itself\) by creating a subscription entry. For use cases, see [Subscription use cases](#event-subscription-use-cases).<br>This operation is performed in the background as a Redfish task. If there is more than one resource that is sending a specific event, the task is further divided into subasks.|
|**Returns** |<ul><li>`Location` URI of the task monitor associated with this operation in the response header.</li><li> Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See "Sample response body \(HTTP 202 status\)".<br>**IMPORTANT:**<br> Note down the task Id. If the task completes with an error, it is required to know which subtask has failed. To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`.</li><li>On success, a `Location` header that contains a link to the newly created subscription and a message in the JSON response body saying that the subscription is created. See "Sample response body \(HTTP 201 status\)".</li></ul>|
|**Response code** |<ul><li>`202 Accepted`</li><li>`201 Created`</li></ul>|
|**Authentication** |Yes|


To know the progress of this action, perform HTTP `GET` on the [task monitor](#viewing-a-task-monitor) returned in the response header \(until the task is complete\).

To get the list of subtask URIs, perform HTTP `GET` on the task URI returned in the JSON response body. See "Sample response body \(HTTP 202 status\)". The JSON response body of each subtask contains a link to the task monitor associated with it. To know the progress of this operation \(subtask\) on a specific server, perform HTTP `GET` on the task monitor associated with the respective subtask.


>**NOTE:**
Only a user with `ConfigureComponents` privilege is authorized to create event subscriptions. If you perform this action without necessary privileges, you will receive an HTTP`403 Forbidden` error.


>**curl command**

```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json; charset=utf-8" \
   -d \
'{ 
   "Name":"ODIMRA_NBI_client",
   "Destination":"https://{Valid_IP_Address}:{Port}/EventListener",
   "EventTypes":[ 
      "Alert"
   ],
   "MessageIds":[ 

   ],
   "ResourceTypes":[ 
      "ComputerSystem"
   ],
   "Context":"ODIMRA_Event",
   "Protocol":"Redfish",
   "SubscriptionType":"RedfishEvent",
   "EventFormatType":"Event",
   "SubordinateResources":true,
   "OriginResources":[
      { 
        "@odata.id":"/redfish/v1/Systems/{ComputerSystemId}"
      },
      {
        "@odata.id":"/redfish/v1/Systems/{ComputerSystemId}"
      }
   ]
}' \
 'https://{odimra_host}:{port}/redfish/v1/EventService/Subscriptions'

```


>**Sample request body**

```
{ 
   "Name":"ODIMRA_NBI_client",
   "Destination":"https://{valid_destination_IP_Address}:{Port}/EventListener",
   "EventTypes":[ 
      "Alert"
   ],
   "MessageIds":[ 

   ],
   "ResourceTypes":[ 
      "ComputerSystem"
   ],
   "Context":"ODIMRA_Event",
   "Protocol":"Redfish",
   "SubscriptionType":"RedfishEvent",
   "EventFormatType":"Event",
   "SubordinateResources":true,
   "OriginResources":[ 
      { 
        "@odata.id":"/redfish/v1/Systems/{ComputerSystemId}"
      },
      {
        "@odata.id":"/redfish/v1/Systems/{ComputerSystemId}"
      }
   ]
}


```

###  Creating event subscription with eventformat type - MetricReport

If `EventFormatType` is empty, default value will be `Event`.

If `EventTypes` list is empty and `EventFormatType` is `MetricReport`, `EventTypes` will default to `[“MetrciReport”]`.

If both values are empty, `EventFormatType` will default to `Event` and `EventTypes` will default to all supported values apart from `MetricReport`.

>**curl command**

```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json; charset=utf-8" \
   -d \

'{
  "Name": "ODIMRA_NBI_client ",
  "Destination": "https://{Valid_IP_Address}:{Port}/EventListener ",
  "EventTypes": ["MetricReport"],
  "Context": "TelemetryDemo",
  "Protocol": "Redfish",
  "SubscriptionType": "RedfishEvent",
  "EventFormatType": "MetricReport"
 }' \
 'https://{odimra_host}:{port}/redfish/v1/EventService/Subscriptions'
```

>**Sample request body**

```
{ 
   "Name":"ODIMRA_NBI_client",
   "Destination":"https://{valid_destination_IP_Address}:{Port}/EventListener",
   "EventTypes":[ 
      "Alert"
   ],
   "MessageIds":[ 

   ],
   "ResourceTypes":[ 
      "ComputerSystem"
   ],
   "Name": "ODIMRA_NBI_client ",
   "Destination": "https://{Valid_IP_Address}:{Port}/EventListener ",
   "EventTypes": ["MetricReport"],
   "Context": "TelemetryDemo",
   "Protocol": "Redfish",
   "SubscriptionType": "RedfishEvent",
   "EventFormatType": "MetricReport"
}

```

**Request parameters**

|Parameter|Value|Attributes|Description|
|---------|-----|----------|-----------|
|Name|String| \(Optional\)<br> |Name for the subscription.|
|Destination|String|Read-only \(Required on create\)<br> |The URL of the destination event listener that listens to events \(Fault management system or any northbound client\).<br>**NOTE:**<br> `Destination` is unique to a subscription: There can be only one subscription for a destination event listener.<br>To change the parameters of an existing subscription , delete it and then create again with the new parameters and a new destination URL.<br> |
|EventTypes|Array \(string \(enum\)\)|Read-only \(Optional\)<br> |The types of events that are sent to the destination. For possible values, see "Event types" table.|
|ResourceTypes|Array \(string, null\)|Read-only \(Optional\)<br> |The list of resource type values \(Schema names\) that correspond to the `OriginResources`. For possible values, perform `GET` on `redfish/v1/EventService` and check values listed under `ResourceTypes` in the JSON response.<br> Examples: "ComputerSystem", "Storage", "Task"<br> |
|Context|String|Read/write Required \(null\)<br> |A string that is stored with the event destination subscription.|
|MessageIds|Array|Read-only \(Optional\)<br> |The key used to find the message in a Message Registry.|
|Protocol|String \(enum\)|Read-only \(Required on create\)<br> |The protocol type of the event connection. For possible values, see "Protocol" table.|
|SubscriptionType|String \(enum\)|Read-only Required \(null\)<br> |Indicates the subscription type for events. For possible values, see "Subscription type" table.|
|EventFormatType|String \(enum\)|Read-only \(Optional\)<br> |Indicates the content types of the message that this service can send to the event destination. For possible values, see "EventFormat" type table.|
|SubordinateResources|Boolean|Read-only \(null\)|Indicates whether the service supports the `SubordinateResource` property on event subscriptions or not. If it is set to `true`, the service creates subscription for an event originating from the specified `OriginResoures` and also from its subordinate resources. For example, by setting this property to `true`, you can receive specified events from a compute node: `/redfish/v1/Systems/{ComputerSystemId}` and from its subordinate resources such as:<br> `/redfish/v1/Systems/{ComputerSystemId}/Memory`,<br> `/redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces`,<br> `/redfish/v1/Systems/{ComputerSystemId}/Bios`,<br> `/redfish/v1/Systems/{ComputerSystemId}/Storage`|
|OriginResources|Array| Optional \(null\)<br> |Resources for which the service only sends related events. If this property is absent or the array is empty, events originating from any resource will be sent to the subscriber. For possible values, see "Origin resources" table.|

**Origin resources**

|String|Description|
|------|-----------|
|A single resource|A specific resource for which the service sends only related events.|
|A list of resources. Supported collections:<br> |A collection of resources for which the service will send only related events.|
|/redfish/v1/Systems|All computer system resources available in Resource Aggregator for ODIM for which the service sends only related events. By setting `EventType` property in the request payload to `ResourceAdded` or `ResourceRemoved` and `OriginResources` property to `/redfish/v1/Systems`, you can receive notifications when a system is added or removed in Resource Aggregator for ODIM.|
|/redfish/v1/Chassis|All chassis resources available in Resource Aggregator for ODIM for which the service sends only related events.|
|/redfish/v1/Fabrics|All fabric resources available in Resource Aggregator for ODIM for which the service sends only related events.|
|/redfish/v1/Managers|All manager resources available in Resource Aggregator for ODIM for which the service sends only related events.|
|/redfish/v1/TaskService/Tasks|All tasks scheduled by or being executed by Redfish `TaskService`. By subscribing to Redfish tasks, you can receive task status change notifications on the subscribed destination client.<br> By specifying the task URIs as `OriginResources` and `EventTypes` as `StatusChange`, you can receive notifications automatically when the tasks are complete.<br> To check the status of a specific task manually, perform HTTP `GET` on its task monitor until the task is complete.<br> |

**Event types**

|String|Description|
|------|-----------|
|Alert|A condition exists which requires attention.|
|ResourceAdded|A resource has been added.|
|ResourceRemoved|A resource has been removed.|
|ResourceUpdated|The value of this resource has been updated.|
|StatusChange|The status of this resource has changed.|
|MetricReport|Collects resource metrics.|

**EventFormat type**

|String|Description|
|------|-----------|
|Event|The subscription destination will receive JSON bodies of the Resource Type Event.|
|MetricReport|Collects resource metrics.|

**Subscription type**

|String|Description|
|------|-----------|
|RedfishEvent|The subscription follows the Redfish specification for event notifications, which is done by a service sending an HTTP `POST` to the destination URI of the subscriber.|

**Protocol**

|String|Description|
|------|-----------|
|Redfish|The destination follows the Redfish specification for event notifications.|

 

>**Sample response header** \(HTTP 202 status\) 

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/taska9702e20-884c-41e2-bd9c-d779a4dd2e6e
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Fri, 08 Nov 2019 07:49:42 GMT+7m 9s
Content-Length:0 byte

```

 

>**Sample response header** \(HTTP 201 status\) 

```
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/EventService/Subscriptions/76088e1c-4654-4eec-a3f6-60bc33b77cdb
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Thu,14 May 2020 09:48:23 GMT+5m 10s
Transfer-Encoding:chunked
```

>**Sample response body** \(HTTP 202 status\) 

```
{
   "@odata.type":"#Task.v1_4_2.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/taskbab2e46d-2ef9-40e8-a070-4e6c87ef72ad",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"taskbab2e46d-2ef9-40e8-a070-4e6c87ef72ad",
   "Name":"Task taskbab2e46d-2ef9-40e8-a070-4e6c87ef72ad",
   "Message":"The task with id taskbab2e46d-2ef9-40e8-a070-4e6c87ef72ad has started.",
   "MessageId":"TaskEvent.1.0.1.TaskStarted",
   "MessageArgs":[
      "taskbab2e46d-2ef9-40e8-a070-4e6c87ef72ad"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```

>**Sample response body** \(subtask\) 

```
{
   "@odata.type":"#SubTask.v1_4_2.SubTask",
   "@odata.id":"/redfish/v1/TaskService/Tasks/taskbab2e46d-2ef9-40e8-a070-4e6c87ef72a/SubTasks/task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "@odata.context":"/redfish/v1/$metadata#SubTask.SubTask",
   "Id":"taskbab2e46d-2ef9-40e8-a070-4e6c87ef72a",
   "Name":"Task taskbab2e46d-2ef9-40e8-a070-4e6c87ef72a",
   "Message":"Successfully Completed Request",
   "MessageId":"Base.1.6.1.Success",
   "Severity":"OK",
   "Members@odata.count":0,
   "Members":null,
   "TaskState":"Completed",
   "StartTime":"2020-05-13T13:33:59.917329733Z",
   "EndTime":"2020-05-13T13:34:00.320539988Z",
   "TaskStatus":"OK",
   "SubTasks":"",
   "TaskMonitor":"/taskmon/task22a98864-0dd8-402b-bfe0-0d61e265391e",
   "PercentComplete":100,
   "Payload":{
      "HttpHeaders":null,
      "HttpOperation":"POST",
      "JsonBody":"",
      "TargetUri":"/redfish/v1/Systems/97d08f36-17f5-5918-8082-f5156618f58d:1"
   },
   "Messages":null
}
```

>**Sample response body** \(HTTP 201 status\) 

```
{
   "error":{
      "@Message.ExtendedInfo":[
         {
            "MessageId":"Base.1.4.Created"
         }
      ],
      "code":"iLO.0.10.ExtendedInfo",
      "message":"See @Message.ExtendedInfo for more information."
   }
}
```



## Submitting a test event

|||
|-----------|-----------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/EventService/Actions/EventService.SubmitTestEvent` |
|**Description** | Once the subscription is successfully created, you can post a test event to Resource Aggregator for ODIM to check whether you are able to receive events. If the event is successfully posted, you will receive a JSON payload of the event response on the client machine \(destination\) that is listening to events. To know more about this event, look up the message registry using the `MessageId` received in the payload. See "Sample message registry \(Alert.1.0.0\)". For more information on message registries, see [Message registries](#message-registries). |
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json; charset=utf-8" \
   -d \
'{ 
   "EventGroupId":{Group_Id_Integer},
   "EventId":"{Unique_Positive_Integer}",
   "EventTimestamp":"{Event_Time_Stamp}",
   "EventType":"{Event_Type_String}",
   "Message":"{Message_String}",
   "MessageArgs":[ 

   ],
   "MessageId":"{message_id_for_messageRegistry}",
   "OriginOfCondition":"/redfish/v1/Systems/{ComputerSystemId}",
   "Severity":"Critical"
}' \
 'https://{odimra_host}:{port}/redfish/v1/EventService/Actions/EventService.SubmitTestEvent'

```




> Sample event payload 

```
{ 
   "EventGroupId":1,
   "EventId":"132489713478812346",
   "EventTimestamp":"2020-02-17T17:17:42-0600",
   "EventType":"Alert",
   "Message":"The LAN has been disconnected",
   "MessageArgs":[ 
       "EthernetInterface 1",
            "/redfish/v1/Systems/8fbda4f3-f55f-4fe4-8db8-4aec1dc3a7d7:1"
   ],
   "MessageId":"Alert.1.0.LanDisconnect",
   "OriginOfCondition":"/redfish/v1/Systems/8fbda4f3-f55f-4fe4-8db8-4aec1dc3a7d7:1/EthernetInterfaces/1",
   "Severity":"Critical"
}

```

 

**Request parameters** 

|Parameter|Value|Attributes|Description|
|---------|-----|----------|-----------|
|EventGroupId|Integer|Optional|The group Id for the event.|
|EventId|String|Optional|The Id for the event to add. This Id is a string of a unique positive integer. Generate a random positive integer and use it as the Id.|
|EventTimestamp|String|Optional|The date and time stamp for the event to add. When the event is received, it translates as the time the event occurred.|
|EventType|String \(enum\)|Optional|The type for the event to add. For possible property values, see "EventType" in [Creating an event subscription](#creating-an-event-subscription).|
|Message|String|Optional|The human-readable message for the event to add.|
|MessageArgs \[ \]|Array \(string\)|Optional|An array of message arguments for the event to add. The message arguments are substituted for the arguments in the message when looked up in the message registry. It helps in trouble ticketing when there are bad events. For example, `MessageArgs` in "Sample event payload" has the following two substitution variables:<br><ul><li>`EthernetInterface 1`</li><li>`/redfish/v1/Systems/{ComputerSystemId}`</li></ul><br>`Description` and `Message` values in "Sample message registry" are substituted with the above-mentioned variables. They translate to "A LAN Disconnect on `EthernetInterface 1` was detected on system `/redfish/v1/Systems/{ComputerSystemId}.` |
|MessageId|String|Required|The Message Id for the event to add. It is the key used to find the message in a message registry. It has `RegistryPrefix` concatenated with the version, and the unique identifier for the message registry entry. The `RegistryPrefix` concatenated with the version is the name of the message registry. To get the names of available message registries, perform HTTP `GET` on `/redfish/v1/Registries`. The message registry mentioned in the sample request payload is `Alert.1.0`.|
|OriginOfCondition|String|Optional|The URL in the `OriginOfCondition` property of the event to add. It is not a reference object. It is the resource that originated the condition that caused the event to be generated. For possible values, see "Origin resources" in [Creating an event subscription](#creating-an-event-subscription).|
|Severity|String|Optional|The severity for the event to add. For possible values, see "Severity" table.|

**Severity**

|String|Description|
|------|-----------|
|Critical|A critical condition that requires immediate attention.|
|OK|Informational or operating normally.|
|Warning|A condition that requires attention.|

 

>**Sample response header** 

```
Cache-Control:no-cache
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Fri,15 May 2020 07:42:59 GMT+5m 11s
Transfer-Encoding:chunked

```

 

> Sample event response 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#Event.Event",
   "@odata.id":"/redfish/v1/EventService/Events/1",
   "@odata.type":"#Event.v1_1_0.Event ",
   "Id":"1",
   "Name":"Event Array",
   "Context":"ODIMRA_Event",
   "Events":[ 
      { 
         "EventType":"Alert",
         "EventId":"132489713478812346",
         "Severity":"Critical",
         "Message":"The LAN has been disconnected",
         "MessageId":"Alert.1.0.LanDisconnect",
         "MessageArgs":[ 
            "EthernetInterface 1",
            "/redfish/v1/Systems/8fbda4f3-f55f-4fe4-8db8-4aec1dc3a7d7:1"
         ],
         "OriginOfCondition":"/redfish/v1/Systems/8fbda4f3-f55f-4fe4-8db8-4aec1dc3a7d7:1/EthernetInterfaces/1",
         "Context":"Event Subscription"
}
]
}
```

> Sample message registry \(Alert.1.0.0\) 

```
{
"@odata.type": "#MessageRegistry.v1_0_0.MessageRegistry",
"Id": “Alert.1.0.0",
"Name": "Base Message Registry",
"Language": "en",
"Description": "This registry is a sample Redfish alert message registry",
"RegistryPrefix": "Alert",
"RegistryVersion": "1.0.0",
"OwningEntity": "Contoso",
"Messages": {
     "LanDisconnect": {
     "Description": "A LAN Disconnect on %1 was detected on system %2.",
     "Message": "A LAN Disconnect on %1 was detected on system %2.",
     "Severity": "Warning",
     "NumberOfArgs": 2,
     "Resolution": "None"
},
...
```

## Event subscription use cases

### Subscribing to resource addition notification

> Resource addition notification payload

```
{ 
   ​   "Name":"EventSubscription",
   ​   "Destination":"https://{Valid_destination_IP_Address}:{Port}/EventListener",
   ​   "EventTypes":[ 
      "ResourceAdded"
   ],
   ​   "MessageIds":[ 
      "ResourceEvent.1.0.2.ResourceAdded"
   ],
   ​   "ResourceTypes":[ 
       "ComputerSystem",
       "Fabric"
   ],
   ​   "Protocol":"Redfish",
   ​   "Context":"Event Subscription",
   ​   "SubscriptionType":"RedfishEvent",
   ​   "EventFormatType":"Event",
   ​   "SubordinateResources":true,
   ​   "OriginResources":[
            { 
      ​        "@odata.id":"/redfish/v1/Systems"
            },
            {
              "@odata.id":"/redfish/v1/Fabrics" 
            }     ​
   ]   ​
}​
```

To get notified whenever a new resource is added in Resource Aggregator for ODIM, subscribe to `ResourceAdded` event originating from any collection \(Systems, Chassis, Fabrics\). 

To create this subscription, perform HTTP `POST` on `/redfish/v1/EventService/Subscriptions` with the sample request payload:



### Subscribing to resource removal notification


> Resource removal notification payload

 ```
{ 
   ​   "Name":"EventSubscription",
   ​   "Destination":"https://{Valid_destination_IP_Address}:{Port}/EventListener",
   ​   "EventTypes":[ 
      "ResourceRemoved"
   ],
   ​   "MessageIds":[ 
      "ResourceEvent.1.0.2.ResourceRemoved"
   ],
   ​   "ResourceTypes":[ 
      "ComputerSystem",
      "Fabric"
   ],
   ​   "Protocol":"Redfish",
   ​   "Context":"Event Subscription",
   ​   "SubscriptionType":"RedfishEvent",
   ​   "EventFormatType":"Event",
   ​   "SubordinateResources":true,
   ​   "OriginResources":[ 
      ​      { 
      ​        "@odata.id":"/redfish/v1/Systems"
            },
            {
              "@odata.id":"/redfish/v1/Fabrics" 
            }     ​
   ]   ​
}​
 ```

To get notified whenever an existing resource is removed from Resource Aggregator for ODIM, subscribe to `ResourceRemoved` event originating from any collection \(Systems, Chassis, Fabrics\). 

To create this subscription, perform HTTP `POST` on `/redfish/v1/EventService/Subscriptions` with the sample request payload:



### Subscribing to task status notifications


> Task status notification payload

```
{ 
   ​   "Name":"EventSubscription",
   ​   "Destination":"https://{Valid_destination_IP_Address}:{Port}/EventListener",
   ​   "EventTypes":[ 
      "StatusChange"
      
   ],
   ​   "MessageIds":[ 
      "ResourceEvent.1.0.2.StatusChange"
   ],
   ​   "ResourceTypes":[ 
      "Task"
   ],
   ​   "Protocol":"Redfish",
   ​   "Context":"Event Subscription",
   ​   "SubscriptionType":"RedfishEvent",
   ​   "EventFormatType":"Event",
   ​   "SubordinateResources":true,
   ​   "OriginResources":[ 
            {
      ​       "@odata.id":"/redfish/v1/TaskService/Tasks" 
            }     ​
   ]   ​
}​
```

There are two ways of checking the task completion status in Resource Aggregator for ODIM: keep polling a task until its completion or simply subscribe to an event notification for task status change \(to receive changes in the task status asynchronously\).

To get notified of the task completion status, subscribe to `StatusChange` event on `/redfish/v1/TaskService/Tasks`. To create this subscription, perform HTTP `POST` on `/redfish/v1/EventService/Subscriptions` with the sample request payload:




## Viewing a collection of event subscriptions

|||
|-----------|-----------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/EventService/Subscriptions` |
|**Description** |This operation lists all the event subscriptions created by the user.|
|**Returns** |A collection of event subscription links.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/EventService/Subscriptions'

```

 

>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#EventDestinationCollection.EventDestinationCollection",
   "@odata.id":"/redfish/v1/EventService/Subscriptions",
   "@odata.type":"#EventDestinationCollection.EventDestinationCollection",
   "Name":"EventSubscriptions",
   "Description":"Event Subscriptions",
   "Members@odata.count":4,
   "Members":[
      {
         "@odata.id":"/redfish/v1/EventService/Subscriptions/57e22fcc-8b1a-460c-ac1f-b3377e22f1cf/"
      },
      {
         "@odata.id":"/redfish/v1/EventService/Subscriptions/72251989-f5e4-453f-9422-bb36d1d94dec/"
      },
      {
         "@odata.id":"/redfish/v1/EventService/Subscriptions/d43babde-c34e-40b8-8a8c-c57b502b9980/"
      },
      {
         "@odata.id":"/redfish/v1/EventService/Subscriptions/6fe7d515-215c-4aba-8ed5-1faed5a91c87/"
      }
   ]
}
```




## Viewing information about a specific event subscription

|||
|-----------|-----------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/EventService/Subscriptions/{subscriptionId}` |
|**Description** |This operation fetches information about a particular event subscription created by the user.|
|**Returns** |JSON schema having the details of this subscription–subscription Id, destination, event types, origin resource, and so on.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/EventService/Subscriptions/{subscriptionId}'

```

 

>**Sample response body** 

```
{
   "@odata.type":"#EventDestination.v1_7_0.EventDestination",
   "@odata.id":"/redfish/v1/EventService/Subscriptions/57e22fcc-8b1a-460c-ac1f-b3377e22f1cf",
   "@odata.context":"/redfish/v1/$metadata#EventDestination.EventDestination",
   "Id":"57e22fcc-8b1a-460c-ac1f-b3377e22f1cf",
   "Name":"ODIM_NBI_client",
   "Destination":"https://{Valid_IP_Address}:{port}/EventListener",
   "Context":"ODIMRA_Event",
   "Protocol":"Redfish",
   "EventTypes":[
      "Alert"
   ],
   "SubscriptionType":"RedfishEvent",
   "MessageIds":[

   ],
   "ResourceTypes":[
      "ComputerSystem"
   ],
   "OriginResources":[
      "@odata.id":"/redfish/v1/Systems/936f4838-9ce5-4e2a-9e2d-34a45422a389:1"
   ]
}
```


##  Deleting an event subscription

|||
|-----------|-----------|
|**Method** | `DELETE` |
|**URI** |`/redfish/v1/EventService/Subscriptions/{subscriptionId}` |
|**Description** |To unsubscribe from an event, delete the corresponding subscription entry. Perform `DELETE` on this URI to remove an event subscription entry.<br>**NOTE:**<br> Only a user with `ConfigureComponents` privilege is authorized to delete event subscriptions. If you perform this action without necessary privileges, you will receive an HTTP`403 Forbidden` error.|
|**Returns** |A message in the JSON response body saying that the subscription is removed.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i -X DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/EventService/Subscriptions/{subscriptionId}'

```

 

>**Sample response body** 

```
{
   "@odata.type":"#EventDestination.v1_7_0.EventDestination",
   "@odata.id":"/redfish/v1/EventService/Subscriptions/57e22fcc-8b1a-460c-ac1f-b3377e22f1cf",
   "Id":"57e22fcc-8b1a-460c-ac1f-b3377e22f1cf",
   "Name":"Event Subscription",
   "Message":"The resource has been removed successfully.",
   "MessageId":"ResourceEvent.1.0.2.ResourceRemoved",
   "Severity":"OK"
}
```




# Message registries

A `MessageRegistry` represents the properties for a message registry.

A message registry is an array of messages and their attributes organized by `MessageId`. Each entry has

-   Description

-   The message this Id translates to.

-   Severity

-   The number and type of arguments

-   Proposed resolution


The arguments are the substitution variables for the message. The `MessageId` is formed according to the Redfish specification. It consists of the `RegistryPrefix` concatenated with the version and the unique identifier for the message registry entry.

**Supported endpoints**

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/Registries|GET|`Login` |
|/redfish/v1/Registries/\{registryId\}|GET|`Login` |
|/redfish/v1/registries/\{registryFileId\}|GET|`Login` |


>**NOTE:**
To access Redfish message registries, ensure that you have a minimum privilege of `Login`. If you do not have the necessary privileges, you will receive an HTTP `403 Forbidden` error.



##  Viewing a collection of registries

|||
|------|--------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/Registries`` |
|**Description** |This endpoint fetches a collection of Redfish-provided registries and custom registries.|
|**Returns** |Links to the registry instances.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Registries'

```

>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#MessageRegistryFileCollection.MessageRegistryFileCollection",
   "@odata.id":"/redfish/v1/Registries/",
   "@odata.type":"#MessageRegistryFileCollection.MessageRegistryFileCollection",
   "Name":"Registry File Repository",
   "Description":"Registry Repository",
   "Members@odata.count":27,
   "Members":[
      {
         "@odata.id":"/redfish/v1/Registries/Base.1.0.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Base.1.2.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Base.1.3.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Base.1.3.1"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Base.1.4.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Base.1.5.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Base.1.6.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Base.1.6.1"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Composition.1.0.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Redfish_1.0.1_PrivilegeRegistry"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Redfish_1.0.2_PrivilegeRegistry"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Redfish_1.0.3_PrivilegeRegistry"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Redfish_1.0.4_PrivilegeRegistry"
      },
      {
         "@odata.id":"/redfish/v1/Registries/ResourceEvent.1.0.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/ResourceEvent.1.0.1"
      },
      {
         "@odata.id":"/redfish/v1/Registries/ResourceEvent.1.0.2"
      },
      {
         "@odata.id":"/redfish/v1/Registries/TaskEvent.1.0.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/TaskEvent.1.0.1"
      },
      {
         "@odata.id":"/redfish/v1/Registries/BiosAttributeRegistryA40.v1_1_46"
      },
      {
         "@odata.id":"/redfish/v1/Registries/%23SmartStorageMessages.v2_0_1.SmartStorageMessages"
      },
      {
         "@odata.id":"/redfish/v1/Registries/iLO.2.13.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/iLOEvents.2.1.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/BiosAttributeRegistryU32.v1_2_00"
      },
      {
         "@odata.id":"/redfish/v1/Registries/BiosAttributeRegistryU30.v1_2_00"
      }
   ]
}
	
```



##  Viewing a single registry

|||
|------|--------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/Registries/{registryId}`` |
|**Description** |This endpoint fetches information about a single registry.|
|**Returns** |Link to the file inside this registry.|
|**Response code** |On success, `200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Registries/{registryId}'

```

>**Sample response body**

```
{
   "Id":"Base.1.6.1",
   "@odata.context":"/redfish/v1/https://10.24.1.95:45000/redfish/v1/$metadata#MessageRegistryFile.MessageRegistryFile",
   "@odata.id":"/redfish/v1/Registries/Base.1.6.1",
   "@odata.type":"#MessageRegistryFile.v1_1_3.MessageRegistryFile",
   "Name":"Registry File Repository",
   "Description":"Base Message Registry File Locations",
   "Languages":[
      "en"
   ],
   "Location":[
      {
         "Language":"en",
         "Uri":"/redfish/v1/registries/Base.1.6.1.json"
      }
   ],
   "Registry":"Base.1.6.1"
}
```




## Viewing a file in a registry

|||
|------|--------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/registries/{registryFileId}`` |
|**Description** |This endpoint fetches information about a file in a registry.|
|**Returns** |Content of this file.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/registries/{jsonFileId}'

```



# Redfish Telemetry Service

Telemetry refers to the metrics obtained from remote systems for analysis and monitoring. 
The Redfish Telemetry model is designed to obtain characteristics of metrics, send specific metric reports periodically and specify triggers against metrics.

Resource Aggregator for ODIM exposes the Redfish `TelemetryService` APIs to:

- Define characteristics and details of one or more metrics (metadata)
- Generate metric reports at regular intervals, specifying their content and timeframe
- Transmit metric reports with metric readings and any metadata associated with the readings
- Specify trigger thresholds for a list of metrics

**Supported APIs**

| API URI                                                      | Operation Applicable | Required privileges     |
| ------------------------------------------------------------ | -------------------- | ----------------------- |
| /redfish/v1/TelemetryService                                 | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/MetricDefinitions               | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/MetricDefinitions/{MetricDefinitionID} | GET                  | `Login`                 |
| /redfish/v1//TelemetryService/MetricReportDefinitions        | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/MetricReportDefinitions/{MetricReportDefinitionID} | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/MetricReports                   | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/MetricReports/{MetricReportID}  | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/Triggers                        | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/Triggers/{TriggerID}            | GET, PATCH           | `Login`,`ConfigureSelf` |

>**NOTE:**
>Before accessing these endpoints, ensure that the user has the required privileges. If you access these endpoints without necessary privileges, you will receive an HTTP `403 Forbidden` error.

## Viewing the telemetry service root

| **Method**         | `GET`                                                        |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/TelemetryService`                               |
| **Description**    | This operation retrieves JSON schema representing the Redfish `TelemetryService` root. |
| **Returns**        | Properties for the Redfish `TelemetryService` and links to its list of resources. |
| **Response Code**  | `200 OK`                                                     |
| **Authentication** | No                                                           |


>**curl command**

```
curl -i GET \
              'https://{odimra_host}:{port}/redfish/v1/TelemetryService'

```


>**Sample response body**

```
{
   "@odata.id":"/redfish/v1/TelemetryService",
   "@odata.type":"#TelemetryService.v1_3_1.TelemetryService",
   "Id":"TelemetryService",
   "Name":"Telemetry Service",
   "Status":{
      "State":"Enabled",
      "Health":"OK"
   },
   "ServiceEnabled":true,
   "SupportedCollectionFunctions":[
      "Average",
      "Minimum",
      "Maximum"
   ],
   "MetricDefinitions":{
      "@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions"
   },
   "MetricReportDefinitions":{
      "@odata.id":"/redfish/v1/TelemetryService/MetricReportDefinitions"
   },
   "MetricReports":{
      "@odata.id":"/redfish/v1/TelemetryService/MetricReports"
   },
   "Triggers":{
      "@odata.id":"/redfish/v1/TelemetryService/Triggers"
   }
}
```

## Collection of metric definitions

| **Method**         | `GET`                                                        |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/TelemetryService/MetricDefinitions`             |
| **Description**    | This operation lists the metadata information for the metrics collection in Redfish implementation. |
| **Returns**        | JSON schema containing the definition, metadata or the characteristics of the metrics collection |
| **Response Code**  | `200 OK`                                                     |
| **Authentication** | No                                                           |


>**curl command**

```
curl -i GET \
              'https://{odimra_host}:{port}/redfish/v1/TelemetryService/MetricDefinitions/'

```


>**Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#MetricDefinitionCollection.MetricDefinitionCollection",
    "@odata.etag": "W/\"1E796226\"",
    "@odata.id": "/redfish/v1/TelemetryService/MetricDefinitions",
    "@odata.type": "#MetricDefinitionCollection.MetricDefinitionCollection",
    "Description": "Metric Definitions view",
    "Members": [
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricDefinitions/MemoryBusUtil"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricDefinitions/IOBusUtil"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricDefinitions/CPUICUtil"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricDefinitions/JitterCount"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricDefinitions/AvgCPU0Freq"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricDefinitions/CPU0Power"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricDefinitions/AvgCPU1Freq"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricDefinitions/CPU1Power"
        }
    ],
    "Members@odata.count": 9,
    "Name": "Metric Definitions"
}
```

## Single metric definition

| **Method**         | `GET`                                                        |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/TelemetryService/MetricDefinitions/{MetricDefinitionID}` |
| **Description**    | This operation lists the metadata information for a metric in Redfish implementation. |
| **Returns**        | JSON schema containing the definition, metadata or the characteristics of a metric |
| **Response Code**  | `200 OK`                                                     |
| **Authentication** | No                                                           |


>**curl command**

```
curl -i GET \
              'https://{odimra_host}:{port}/redfish/v1/TelemetryService/MetricDefinitions/{MetricDefinitionID}'

```


>**Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#MetricDefinition.MetricDefinition",
    "@odata.etag": "W/\"AB720077\"",
    "@odata.id": "/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil",
    "@odata.type": "#MetricDefinition.v1_0_0.MetricDefinition",
    "Calculable": "NonSummable",
    "CalculationAlgorithm": "Average",
    "Description": "Metric definition for CPU Utilization",
    "Id": "CPUUtil",
    "Implementation": "PhysicalSensor",
    "IsLinear": true,
    "MaxReadingRange": 100,
    "MetricDataType": "Decimal",
    "MetricProperties": [
        "/redfish/v1/Systems/{SystemID}#SystemUsage/CPUUtil"
    ],
    "MetricType": "Numeric",
    "MinReadingRange": 0,
    "Name": "Metric definition for CPU Utilization",
    "Units": "%",
    "Wildcards": [
        {
            "Name": "SystemID",
            "Values": [
                "9616fec9-c76a-4d26-ab53-196d08ce825a:1",
                "ba5cd083-b360-4994-bc30-12b450859b27:1"
            ]
        }
    ]
}
```

## Collection of Metric Report Definitions

| **Method**         | `GET`                                                        |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/TelemetryService/MetricReportDefinitions`       |
| **Description**    | This operation represents a set of metric properties of collected in multiple metric reports. |
| **Returns**        | JSON schema defining the content and time of the metric reports |
| **Response Code**  | `200 OK`                                                     |
| **Authentication** | No                                                           |


>**curl command**

```
curl -i GET \
              'https://{odimra_host}:{port}/redfish/v1/TelemetryService/MetricReportDefinitions/'

```


>**Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#MetricReportDefinitionCollection.MetricReportDefinitionCollection",
    "@odata.etag": "W/\"BFD5C070\"",
    "@odata.id": "/redfish/v1/TelemetryService/MetricReportDefinitions",
    "@odata.type": "#MetricReportDefinitionCollection.MetricReportDefinitionCollection",
    "Description": " MetricReportDefinitions view",
    "Members": [
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricReportDefinitions/CPUUtilCustom1"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricReportDefinitions/CPUUtilCustom2"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricReportDefinitions/CPUUtilCustom3"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricReportDefinitions/CPUUtilCustom4"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricReportDefinitions/MemoryBusUtilCustom1"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricReportDefinitions/MemoryBusUtilCustom2"
        },
    ],
    "Members@odata.count": 6,
    "Name": "MetricReportDefinitions"
}
```

## Single metric report definition 

| **Method**         | `GET`                                                        |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/TelemetryService/MetricReportDefinitions/{MetricReportDefinitionID}` |
| **Description**    | This operation represents metric properties of a single metric report. |
| **Returns**        | JSON schema defining the content and periodicity of the metric report |
| **Response Code**  | `200 OK`                                                     |
| **Authentication** | No                                                           |


>**curl command**

```
curl -i GET \
              'https://{odimra_host}:{port}/redfish/v1/TelemetryService/MetricReportDefinitions/{MetricReportDefinitionID}'

```


>**Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#MetricReportDefinition.MetricReportDefinition",
    "@odata.etag": "W/\"9A613B5C\"",
    "@odata.id": "/redfish/v1/TelemetryService/MetricReportDefinitions/CPUUtilCustom1",
    "@odata.type": "#MetricReportDefinition.v1_0_0.MetricReportDefinition",
    "Description": "Metric report of CPU Utilization for 10 minutes with sensing interval of 20 seconds.",
    "Id": "CPUUtilCustom1",
    "MetricProperties": [
        "/redfish/v1/Systems/{SystemID}#SystemUsage/CPUUtil",
        "SystemUsage/CPUUtil"
    ],
    "MetricReport": {
        "@odata.id": "/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1"
    },
    "MetricReportDefinitionType": "OnRequest",
    "Metrics": [
        {
            "CollectionDuration": "PT20S",
            "CollectionFunction": "Average",
            "CollectionTimeScope": "Interval",
            "MetricId": "CPUUtil"
        }
    ],
    "Name": "Metric report of CPU Utilization for 10 minutes with sensing interval of 20 seconds.",
    "Status": {
        "Health": "OK",
        "State": "Enabled"
    },
    "Wildcards": [
        {
            "Name": "SystemID",
            "Values": [
                "9616fec9-c76a-4d26-ab53-196d08ce825a:1"
            ]
        }
    ]
}
```

## Collection of metric reports

| **Method**         | `GET`                                                        |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/TelemetryService/MetricReports`                 |
| **Description**    | This operation retrieves collection of reports with metric readings and any metadata associated with the readings. |
| **Returns**        | Links of the metric reports.                                 |
| **Response Code**  | `200 OK`                                                     |
| **Authentication** | No                                                           |


>**curl command**

```
curl -i GET \
              'https://{odimra_host}:{port}/redfish/v1/TelemetryService/MetricReports/'

```


>**Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#MetricReportCollection.MetricReportCollection",
    "@odata.etag": "W/\"BFD5C070\"",
    "@odata.id": "/redfish/v1/TelemetryService/MetricReports",
    "@odata.type": "#MetricReportCollection.MetricReportCollection",
    "Description": " Metric Reports view",
    "Members": [
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom1"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom2"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom3"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricReports/MemoryBusUtilCustom1"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricReports/MemoryBusUtilCustom2"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricReports/MemoryBusUtilCustom3"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/MetricReports/MemoryBusUtilCustom4"
        },        
    ],
    "Members@odata.count": 7,
    "Name": "Metric Reports"
}
 
```

## Single metric report 

| **Method**         | `GET`                                                        |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/TelemetryService/MetricReports/{MetricReportID}` |
| **Description**    | This operation retrieves a report with metric readings and any metadata associated with the readings. |
| **Returns**        | Link to the metric report.                                   |
| **Response Code**  | `200 OK`                                                     |
| **Authentication** | No                                                           |


>**curl command**

```
curl -i GET \
              'https://{odimra_host}:{port}/redfish/v1/TelemetryService/MetricReports/{MetricReportID}'

```


>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#MetricReport.MetricReport",
   "@odata.id":"/redfish/v1/TelemetryService/MetricReports/CPUUtilCustom2",
   "@odata.type":"#MetricReport.v1_0_0.MetricReport",
   "Description":"Metric report of CPU Utilization for 60 minutes with sensing interval of 20 seconds.",
   "Id":"CPUUtilCustom2",
   "MetricReportDefinition":{
      "@odata.id":"/redfish/v1/TelemetryService/MetricReportDefinitions/CPUUtilCustom2"
   },
   "MetricValues":[
      {
         "MetricDefinition":{
            "@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"
         },
         "MetricId":"CPUUtil",
         "MetricProperty":"/redfish/v1/Systems/9616fec9-c76a-4d26-ab53-196d08ce825a:1#SystemUsage/CPUUtil",
         "MetricValue":"1",
         "Timestamp":"2021-08-28T13:46:05Z"
      },
      {
         "MetricDefinition":{
            "@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"
         },
         "MetricId":"CPUUtil",
         "MetricProperty":"/redfish/v1/Systems/9616fec9-c76a-4d26-ab53-196d08ce825a:1#SystemUsage/CPUUtil",
         "MetricValue":"1",
         "Timestamp":"2021-08-28T13:46:25Z"
      },
      {
         "MetricDefinition":{
            "@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"
         },
         "MetricId":"CPUUtil",
         "MetricProperty":"/redfish/v1/Systems/9616fec9-c76a-4d26-ab53-196d08ce825a:1#SystemUsage/CPUUtil",
         "MetricValue":"1",
         "Timestamp":"2021-08-28T13:46:45Z"
      },
      {
         "MetricDefinition":{
            "@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"
         },
         "MetricId":"CPUUtil",
         "MetricProperty":"/redfish/v1/Systems/9616fec9-c76a-4d26-ab53-196d08ce825a:1#SystemUsage/CPUUtil",
         "MetricValue":"1",
         "Timestamp":"2021-08-28T13:47:05Z"
      },
      "Name":"Metric report of CPU Utilization for 60 minutes with sensing interval of 20 seconds."
   }
```

<blockquote> NOTE:  After you remove a system and perform a `GET ` operation on the Metric Report Collection, the collection of all individual metric reports is still displayed in the response body. When you perform a `GET` operation on that individual {MetricReportID}, you get a `404-Not Found` error message. After this, when you perform a GET operation on the Metric Report Collection again, the instance of that individual metric report is erased. 
This is an implementation choice in Resource Aggregator for ODIM, because Telemetry service is defined for a collection of BMCs and not for an individual BMC as per the DMTF Redfish specification.</blockquote>

## Collection of triggers

| **Method**         | `GET`                                                        |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/TelemetryService/Triggers`                      |
| **Description**    | This operation retrieves the collection of triggers that apply to multiple metric properties. |
| **Returns**        | Links to a collection of triggers                            |
| **Response Code**  | `200 OK`                                                     |
| **Authentication** | No                                                           |


>**curl command**

```
curl -i GET \
              'https://{odimra_host}:{port}/redfish/v1/TelemetryService/Triggers/'

```


>**Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#TriggersCollection.TriggersCollection",
    "@odata.etag": "W/\"DA402EBA\"",
    "@odata.id": "/redfish/v1/TelemetryService/Triggers",
    "@odata.type": "#TriggersCollection.TriggersCollection",
    "Description": " Triggers view",
    "Members": [
        {
            "@odata.id": "/redfish/v1/TelemetryService/Triggers/CPUUtilTriggers"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/Triggers/MemoryBusUtilTriggers"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/Triggers/IOBusUtilTriggers"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/Triggers/CPUICUtilTriggers"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/Triggers/JitterCountTriggers"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/Triggers/CPU0PowerTriggers"
        },
        {
            "@odata.id": "/redfish/v1/TelemetryService/Triggers/CPU1PowerTriggers"
        }
    ],
    "Members@odata.count": 7,
    "Name": "Triggers"
}
```

## Single trigger 

| **Method**         | `GET`                                                        |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/TelemetryService/Triggers/{TriggersID}`         |
| **Description**    | This endpoint retrieves a trigger that apply to the listed metrics. |
| **Returns**        | Link of a single trigger                                     |
| **Response Code**  | `200 OK`                                                     |
| **Authentication** | No                                                           |


>**curl command**

```
curl -i GET \
              'https://{odimra_host}:{port}/redfish/v1/TelemetryService/Triggers/'

```


>**Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#Triggers.Triggers",
    "@odata.etag": "W/\"BFAAE441\"",
    "@odata.id": "/redfish/v1/TelemetryService/Triggers/CPUUtilTriggers",
    "@odata.type": "#Triggers.v1_0_0.Triggers",
    "Description": "Triggers for CPU Utilization",
    "Id": "CPUUtilTriggers",
    "MetricProperties": [
        "/redfish/v1/Systems/{SystemID}#SystemUsage/CPUUtil"
    ],
    "MetricType": "Numeric",
    "Name": "Triggers for CPU Utilization",
    "NumericThresholds": {
        "LowerCritical": {
            "Activation": "Decreasing",
            "DwellTime": "PT0S",
            "Reading": 0
        },
        "UpperCritical": {
            "Activation": "Increasing",
            "DwellTime": "PT0S",
            "Reading": 0
        }
    },
    "Status": {
        "Health": "OK",
        "State": "Enabled"
    },
    "TriggerActions": [
        "LogToLogService"
    ],
    "Wildcards": [
        {
            "Name": "SystemID",
            "Values": [
                "9616fec9-c76a-4d26-ab53-196d08ce825a:1",
                "ba5cd083-b360-4994-bc30-12b450859b27:1"
            ]
        }
    ]
}
```

## Updating a trigger

| **Method**         | `PATCH`                                              |
| ------------------ | ---------------------------------------------------- |
| **URI**            | `/redfish/v1/TelemetryService/Triggers/{TriggersID}` |
| **Description**    | This operation updates triggers of each metric.      |
| **Returns**        |                                                      |
| **Response Code**  | `200 OK`                                             |
| **Authentication** | No                                                   |


>**curl command**

```
curl -i -X PATCH \
   -H "Authorization:Basic YWRtaW46T2QhbTEyJDQ=" \
   -H "Content-Type:application/json" \
   -d \
'{
  "EventTriggers": ["Alert"]
}' \
 'https://10.24.1.33:30080/redfish/v1/TelemetryService/Triggers/{TriggersID}'

```


>**Sample request body**

```
{
  "EventTriggers": ["Alert"]
}
```



