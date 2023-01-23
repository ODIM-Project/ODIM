

# Table of contents

- [Resource Aggregator for Open Distributed Infrastructure Management](#resource-aggregator-for-open-distributed-infrastructure-management)
  * [Resource Aggregator for ODIM logical architecture](#resource-aggregator-for-odim-logical-architecture)
- [API usage and access guidelines](#api-usage-and-access-guidelines)
  - [HTTP headers](#http-headers)
  - [Base URL](#base-url)
  - [Curl command](#curl-command)
  - [Curl command options (flags)](#curl-command-options)
  - [Including HTTP certificate](#including-http-certificate)
  - [HTTP request methods](#http-request-methods)
  - [Responses](#Responses)
  - [Common response header properties](#common-response-header-properties)
  - [Status codes](#status-codes)
- [IPV6 support](#ipv6-support)
- [Support for URL Encoding](#support-for-url-encoding)
- [List of supported APIs](#list-of-supported-apis)
  * [Viewing the list of supported Redfish services](#viewing-the-list-of-supported-redfish-services)
  * [Modifying configurations for services](#Modifying-configurations-for-services)
- [Rate limits](#rate-limits)
- [Authentication and authorization](#authentication-and-authorization)
  * [Authentication methods for Redfish APIs](#authentication-methods-for-redfish-apis)
  * [Role-based authorization](#role-based-authorization)
      + [Roles](#roles)
      + [Priviliges](#privileges)
- [Sessions](#sessions)
  * [Viewing the SessionService root](#viewing-the-sessionservice-root)
  * [Creating a session](#creating-a-session)
  * [Viewing a list of sessions](#viewing-a-list-of-sessions)
  * [Viewing information about a session](#viewing-information-about-a-session)
  * [Deleting a session](#deleting-a-session)
- [User roles and privileges](#user-roles-and-privileges)
  * [Viewing the AccountService root](#viewing-the-accountservice-root)
  * [Viewing a list of roles](#viewing-a-list-of-roles)
  * [Viewing information about a role](#viewing-information-about-a-role)
- [User accounts](#user-accounts)
  * [Creating a user account](#creating-a-user-account)
    + [Password requirements](#password-requirements)
  * [Viewing a list of user accounts](#viewing-a-list-of-user-accounts)
  * [Viewing information about an account](#viewing-information-about-an-account)
  * [Updating a user account](#updating-a-user-account)
  * [Deleting a user account](#deleting-a-user-account)
- [Resource aggregation and management](#resource-aggregation-and-management)
  * [Viewing the AggregationService root](#viewing-the-aggregationservice-root)
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
  * [StoragePools](#StoragePools)
    * [Collection of StoragePools](#StoragePools-collection)
    * [Single StoragePool](#Single-StoragePool)
    * [Collection of AllocatedVolumes](#AllocatedVolumes-Collection)
    * [Single AllocatedVolume](#single-AllocatedVolume)
    * [Collection of ProvidingDrives](#ProvidingDrives-Collection)
    * [Single ProvidingDrive](#single-ProvidingDrive)
  * [Storage subsystem](#storage-subsystem)
  * [Drives](#drives)
    + [Single drive](#single-drive)
  * [Volumes](#volumes)
    + [Collection of volumes](#collection-of-volumes)
    + [Viewing volume capabilities](#viewing-volume-capabilities)
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
  * [Remote BMC accounts and roles](#remote-bmc-accounts-and-roles)
    * [Viewing the RemoteAccountService root](#viewing-the-remoteaccountservice-root)
    * [Collection of BMC user accounts](#collection-of-bmc-user-accounts)
    * [Single BMC user account](#single-bmc-user-account)
    * [Creating a BMC account](#creating-a-bmc-account)
    * [Updating a BMC account](#updating-a-bmc-account)
    * [Deleting a BMC account](#deleting-a-bmc-account)
    * [Collection of BMC roles](#collection-of-bmc-roles)
    * [Single role](#single-role)
- [Software and firmware inventory](#software-and-firmware-inventory)
  * [Viewing the UpdateService root](#viewing-the-updateservice-root)
  * [Viewing the firmware inventory](#viewing-the-firmware-inventory)
  * [Viewing a specific firmware resource](#viewing-a-specific-firmware-resource)
  * [Viewing the software inventory](#viewing-the-software-inventory)
  * [Viewing a specific software resource](#viewing-a-specific-software-resource)
  * [Actions](#actions)
    + [Simple update](#simple-update)
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
  * [Viewing the TaskService root](#viewing-the-taskservice-root)
  * [Viewing a collection of tasks](#viewing-a-collection-of-tasks)
  * [Viewing information about a specific task](#viewing-information-about-a-specific-task)
  * [Viewing a task monitor](#viewing-a-task-monitor)
  * [Deleting a task](#deleting-a-task)
- [Events](#events)
  * [Viewing the event service root](#viewing-the-eventservice-root)
  * [Creating an event subscription](#creating-an-event-subscription)
    + [Sample event](#sample-event)
    + [Creating event subscription with eventformat type “MetricReport”](#creating-event-subscription-with-eventformat-type---metricreport)
  * [Submitting a test event](#submitting-a-test-event)
  * [Event subscription use cases](#event-subscription-use-cases)
    + [Subscribing to resource addition notification](#subscribing-to-resource-addition-notification)
    + [Subscribing to resource removal notification](#subscribing-to-resource-removal-notification)
    + [Subscribing to task status notifications](#subscribing-to-task-status-notifications)
  * [Viewing a collection of event subscriptions](#viewing-a-collection-of-event-subscriptions)
  * [Viewing information about a specific event subscription](#viewing-information-about-a-specific-event-subscription)
  * [Deleting an event subscription](#deleting-an-event-subscription)
  * [Undelivered events](#undelivered-events)
- [Message registries](#message-registries)
  * [Viewing a collection of registries](#viewing-a-collection-of-registries)
  * [Viewing a single registry](#viewing-a-single-registry)
  * [Viewing a file in a registry](#viewing-a-file-in-a-registry)
- [Redfish Telemetry Service](#redfish-telemetry-service)
  * [Viewing the TelemetryService root](#viewing-the-telemetryservice-root)
  * [Collection of metric definitions](#collection-of-metric-definitions)
  * [Single metric definition](#single-metric-definition)
  * [Collection of Metric Report Definitions](#collection-of-metric-report-definitions)
  * [Single metric report definition](#single-metric-report-definition)
  * [Collection of metric reports](#collection-of-metric-reports)
  * [Single metric report](#single-metric-report)
  * [Collection of Triggers](#collection-of-triggers)
  * [Single Trigger](#single-trigger)
  * [Updating a trigger](#updating-a-trigger)
- [License Service](#license-service)
  - [Viewing the LicenseService root](#viewing-the-licenseservice-root)
  - [Viewing the license collection](#viewing-the-license-collection)
  - [Viewing information about a license](#viewing-information-about-a-license)
  - [Installing a license](#installing-a-license)
- [Audit logs](#audit-logs)
- [Security logs](#security-logs)
- [Application logs](#Application-logs)

# Resource Aggregator for Open Distributed Infrastructure Management

Resource Aggregator for Open Distributed Infrastructure Management (Resource Aggregator for ODIM) is a modular, open framework for simplified management and orchestration of distributed physical infrastructure. It provides a unified management platform for converging multivendor hardware equipment. By exposing a standards-based programming interface, it enables easy and secure management of a wide range of multivendor IT infrastructure distributed across multiple data centers.

Resource Aggregator for ODIM framework comprises the following two components.

- The resource aggregation function (the resource aggregator)

  The resource aggregator is the single point of contact between the northbound clients and the southbound infrastructure. The primary function of the resource aggregator is to build and maintain a central resource inventory. It exposes Redfish-compliant APIs to allow northbound infrastructure management systems to:

    - Get a unified view of the southbound compute, local storage, and Ethernet switch fabrics available in the resource inventory
    - Gather crucial configuration information about southbound resources
    - Manipulate groups of resources in a single action
    - Listen to similar events from multiple southbound resources
  
- One or more plugins

  The plugins abstract, translate, and expose southbound resource information to the resource aggregator through RESTful APIs. Resource Aggregator for ODIM supports:
    - Generic Redfish (GRF) plugin for ODIM—Plugin that can be used for any Redfish-compliant device
    - Plugin for unmanaged racks (URP)—Plugin that acts as a resource manager for unmanaged racks
    - Integration of additional third-party plugins—Dell, Lenovo and Cisco ACI plugins


This guide provides reference information for the northbound APIs exposed by the resource aggregator. These APIs are designed as per DMTF's *[Redfish® Scalable Platforms API (Redfish) specification 1.15.1](https://www.dmtf.org/sites/default/files/standards/documents/DSP0266_1.15.1.pdf)* and are Redfish-compliant.

The Redfish® standard is a suite of specifications that deliver an industry standard protocol providing a RESTful interface for the simple and secure management of servers, storage, networking, multivendor, converged and hybrid IT infrastructure. Redfish uses JSON and OData.


##  Resource Aggregator for ODIM logical architecture

Resource Aggregator for ODIM framework adopts a layered architecture and has many functional layers. The architecture diagram shows these functional layers of Resource Aggregator for ODIM deployed in a data center.

![ODIM_architecture](images/arch.png)

**API layer**

This layer hosts a REST server which is open-source and secure. It learns about the southbound resources from the plugin layer and exposes the corresponding Redfish data model payloads to the northbound clients. The northbound clients communicate with this layer through a REST-based protocol that is compliant with DMTF's Redfish® specifications (Schema 2022.1 and Specification 1.15.1).
The API layer sends user requests to the plugins through the aggregation, the event, and the fabric services.

**Services layer**

This layer hosts all the services. The layer implements service logic for all use cases through an extensible domain model (Redfish Data Model). All resource information is stored in this data model and is used to service the API requests coming from the API layer. Any responses from the plugin layer might update the domain model. It maintains the state for event subscriptions, credentials, and tasks.


![Redfish_data_model](images/redfish_data_model.png)

**Event message bus layer**

This layer hosts a message broker which acts as a communication channel between the plugin layer and the upper layers. The layer supports common messaging architecture to forward events received from the plugin layer to the upper layers. During the runtime, Resource Aggregator for ODIM uses either Kafka or the RedisStreams service as the event message bus.  The services and the event message bus layers host Redis data store.

**Plugin layer**

Plugins abstract vendor-specific access protocols to a common interface which the aggregator layers use to communicate with the resources. The plugin layer connects the actual managed resources to the aggregator layers and is decoupled from the upper layers. The layer uses REST-based communication to interact with the other layers. It collects events to be exposed to fault management systems and uses the event message bus to publish events. 
The plugin layer allows developers to create plugins on the tool set of their choice without enforcing any strict language binding. To know how to develop plugins, see *[Resource Aggregator for Open Distributed Infrastructure Management Plugin Developer's Guide](https://github.com/ODIM-Project/ODIM/blob/development/plugin-redfish/README.md)*.


# API usage and access guidelines

> **PREREQUISITE**: Ensure that you have the required privileges to access all the services to avoid encountering the HTTP `403 Forbidden` error.

This guide contains sample request and response payloads. For information on response payload parameters, see *[Redfish® Scalable Platforms API (Redfish) schema 2022.1](https://www.dmtf.org/sites/default/files/standards/documents/DSP2046_2022.1.pdf)*.

To access the RESTful APIs exposed by the resource aggregator, you need an HTTPS-capable client, such as a web browser with a REST Client plugin extension, or a Desktop REST Client application, or curl (a popular, free command-line utility). 

> **TIP**: It is good to use a tool, such as curl or any Desktop REST Client application to send requests.


> **IMPORTANT:** The response codes, JSON request and response parameters provided in this guide might vary for systems depending on the vendor, model, and firmware versions.

## **HTTP headers**

HTTP headers include the following:

- `"Content-type":"application/json; charset=utf-8"` for all RESTful API operations that include a request body in JSON format.
- Authentication header (`BasicAuth` or `XAuthToken`) for all RESTful API operations except the HTTP `GET` operation on the Redfish service root and the HTTP `POST` operation on sessions.

## **Base URL**

Use the following base URL in all your HTTP requests:

`https://{odimra_host}:{port}/`

- {odimra_host} is the fully qualified domain name (FQDN) used for generating certificates while deploying the resource aggregator.

	>**NOTE:** Ensure that FQDN is provided in the `/etc/hosts` file or in the DNS server.


- {port} is the port where the services of the resource aggregator are running. The default port is 45000. If you have changed the default port in the `/etc/odimra_config/odimra_config.json` file, use that as the port in the base URL.
>**NOTE**: To access the base URL using a REST client, replace `{odimra_host}` with the IP address of the system where the resource aggregator is installed. To use FQDN in place of `{odimra_host}`, add the Resource Aggregator for ODIM server certificate to the browser where the REST client is launched.

## curl

*[curl](https://curl.haxx.se)* is a command-line tool which helps you get or send information through URLs using supported protocols. Resource Aggregator for ODIM supports HTTPS protocol. Examples in this document use curl commands to make HTTP requests.

>**IMPORTANT:** If you have set proxy configuration, set `no_proxy` using the following command before you run a curl command:
>
>```
>export no_proxy="127.0.0.1,localhost,{odimra_host}"
>```

## curl command options

- `--cacert` <file_path> includes a specified X.509 root certificate.
- `-H` passes on custom headers.
- `-X` specifies a custom request method. Use `-X` for HTTP `PATCH`, `PUT`, and `DELETE` operations.
- `-d` posts data to a URI. Use `-d` for all HTTP operations that include a request body.
- `-i` returns HTTP response headers.
- `-v` fetches verbose.

For a complete list of curl flags, see *[https://curl.haxx.se](https://curl.haxx.se)*.

## Including HTTP certificate

Without CA certificate, curl fails to verify that HTTP connections are secure and the curl commands might fail with the SSL certificate problem. Provide the root CA certificate in curl for secure SSL communication.

- To run curl commands on the server on which you deployed Resource Aggregator for ODIM, provide the `rootCA.crt` file by running the following command:

  ```
  curl -v --cacert {path}/rootCA.crt 'https://{odimra_host}:{port}/redfish/v1'
  ```

   {path} is where you have generated certificates during the Resource Aggregator for ODIM deployment.

- To run curl commands on a different server, perform the following steps to provide the rootCA.crt file.

   1. Navigate to `~/ODIM/build/cert_generator/certificates` on the server where you have deployed Resource Aggregator for ODIM.

   2. Copy the `rootCA.crt` file.
   3. Log in to your server and paste the `rootCA.crt` file in a folder.
   4. Open the `/etc/hosts` file to edit it.
   5. Scroll to the end of the file, add the following line, and save:
      `{odim_server_ipv4_address} {FQDN}`
   6. Check if curl is working by running the curl command:
       ```
       curl -v --cacert {path}/rootCA.crt 'https://{odimra_host}:{port}/redfish/v1'
       ```

   >**NOTE:** To avoid using the `--cacert` flag in every curl command, add `rootCA.crt` in the `ca-certificates.crt` file located in this path:<br> `/etc/ssl/certs/ca-certificates.crt`.

## HTTP request methods

Use the listed Redfish-defined HTTP methods to implement various actions.

| HTTP Request Method       | Description                                                  |
| ------------------------- | ------------------------------------------------------------ |
| `GET` [Read Requests]     | Use this method to request a representation of a specified resource (single resource or collection). |
| `PATCH` [Update]          | Use this method to apply partial modifications to a resource. |
| `POST` [Create] [Actions] | Use this method to create a resource. Submit this request to the resource collection to which you want to add the new resource. You can also use this method to initiate operations on a resource or a collection of resources. |
| `PUT` [Replace]           | Use this method to replace the property values of a resource completely. It is used to both create and update the state of a resource. |
| `DELETE` [Delete]         | Use this method to delete a resource.                        |

## Responses 

Resource Aggregator for ODIM supports the listed responses:

| Response                     | Description                                                  |
| ---------------------------- | ------------------------------------------------------------ |
| Metadata response            | Describes the resources and types exposed by the service to generic clients |
| Resource response            | Response in JSON format for an individual resource           |
| Resource collection response | Response in JSON format for a collection of resources        |
| Error response               | If there is an HTTP error, a JSON response is returned with additional information |

## Common response header properties

The listed properties are common across all response headers, and are omitted from the samples in this document. 

```
"Connection": "keep-alive",
"OData-Version": "4.0",
"X-Frame-Options": "sameorigin",
"X-Content-Type-Options":"nosniff",
"Content-type":"application/json; charset=utf-8",
"Cache-Control":"no-cache, no-store, must-revalidate",
"Transfer-Encoding":"chunked",
```

## Status codes

The HTTP status codes include the success codes and the error codes and their respective descriptions for all API operations.

| Success code<br> | Description                                                  |
| ---------------- | ------------------------------------------------------------ |
| 200 OK           | The request completes successfully with the representation in the body. |
| 201 Created      | A new resource is successfully created with the `Location` header set to well-defined URI for the newly created resource. The response body might include the representation of the newly created resource. |
| 202 Accepted     | The request has been accepted for processing but not processed. The `Location` header is set to URI of a task monitor that can be queried later for the status of the operation. |
| 204 No Content   | The request succeeds, but no content is returned in the response body. |

| Error code<br>            | Description                                                  |
| ------------------------- | ------------------------------------------------------------ |
| 301 Moved Permanently     | The requested resource resides in a different URI given by the `Location` headers. |
| 400 Bad Request           | The request cannot be performed due to missing or invalid information. An extended error message is returned in the response body. |
| 401 Unauthorized          | The request has missing or invalid authentication credentials. |
| 403 Forbidden             | The server recognizes that the credentials do not have the necessary authorization to perform the operation. |
| 404 Not Found             | The request specifies the URI of a non-existing resource.    |
| 405 Method Not Allowed    | The HTTP method specified in the request is not supported for a particular request URI. The response includes `Allow` header that lists the supported methods. |
| 409 Conflict              | A resource creation or an update is incomplete because it conflicts with the current state of the resources supported by the platform. |
| 500 Internal Server Error | The server encounters an unexpected condition that prevents it from fulfilling the request. |
| 501 Not Implemented       | The server has not implemented the method for the resource.  |
| 503 Service Unavailable   | The server is unable to service the request due to temporary overloading or maintenance. |

# IPv6 support

Resource Aggregator for ODIM supports IPv6 address to send API service requests. 

The default value for `nwPreferences` parameter in the Resource Aggregator for ODIM deployment configuration file (`kube_deploy_nodes.yaml`) is `ipv4`. To send the requests using the IPv6 addresses, set the `nwPreferences` parameter value to `dualStack`. This allows you to send requests using both the IPv4 and IPv6 addresses. 

### Sample APIs with IPv6 address

- curl command to view a collection of systems

  ```
  curl -i GET \
  -H "X-Auth-Token:{X-Auth-Token}" \
  'https://{IPv6 address}:{port}/redfish/v1/Systems'
  ```

- curl command to add a server

  ```
  curl -i -X POST \
  -H "Authorization:Basic YWRtaW46T2QhbTEyJDQ=" \
  -H "Content-Type:application/json" \
  -d \
  '{
    "HostName":"xxx.xxx.xxx.xxx",
    "UserName":"admin",
    "Password":"<your_password>",
    "Links":{
       "ConnectionMethod":{
       "@odata.id":"/redfish/v1/AggregationService/ConnectionMethods/e9fec4a3-a9f7-4d4e-b65f-8d9316e7f0d9"
    }
   }
  }' \
  'https://{IPv6 address}:{port}/redfish/v1/AggregationService/AggregationSources'
  ```

- curl command to view the connection method

  ```
  curl -i -X GET \
  -H "Authorization:Basic YWRtaW46T2QhbTEyJDQ=" \
  'https://{IPv6 address}:{port}/redfish/v1/AggregationService/ConnectionMethods/'
  ```


# Support for URL Encoding

The URL encoding mechanism translates the characters in the URLs to a representation that are universally accepted by all web browsers and servers. 

Resource Aggregator for ODIM supports all standard encoded characters for all the API URLs. When Resource Aggregator for ODIM gets an encoded URL path, the non-ASCII characters in the URL are internally translated and sent to the web browsers. 

Replace a character in a URL with its standard encoding notation. Resource Aggregator for ODIM accepts the encoded notation, decodes it to the actual character acceptable by the web browsers and sends an accurate response.

**For example**: In the URL`/redfish/v1/Systems/e24fb205-6669-4080-b53c-67d4923aa73e.1`, if you replace the  `/` character with its encoded notation %2F and send the request, Resource Aggregator for ODIM accepts and decodes the encoded notation internally and sends a response.

> **Tip**: You can visit *https://www.w3schools.com/tags/ref_urlencode.ASP* or browse the Internet to view the standard ASCII Encoding Reference of the URL characters.

# List of supported APIs

Resource Aggregator for ODIM supports the listed Redfish APIs:

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
|redfish/v1/SessionService/Sessions/{sessionId}|`GET`, `DELETE`|

|AccountService||
|-------|--------------------|
|/redfish/v1/AccountService|`GET`|
|/redfish/v1/AccountService/Accounts|`POST`, `GET`|
|/redfish/v1/AccountService/Accounts/{accountId}|`GET`, `DELETE`, `PATCH`|
|/redfish/v1/AccountService/Roles|`POST`, `GET`|
|/redfish/v1/AccountService/Roles/{roleId}|`GET`, `DELETE`, `PATCH`|

|AggregationService||
|-------|--------------------|
|/redfish/v1/AggregationService|`GET`|
|/redfish/v1/AggregationService/AggregationSources<br> |`GET`, `POST`|
|/redfish/v1/AggregationService/AggregationSources/{aggregationSourceId}|`GET`, `PATCH`, `DELETE`|
|/redfish/v1/AggregationService/Actions/AggregationService.Reset|`POST`|
|/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder|`POST`|
|/redfish/v1/AggregationService/Aggregates|`GET`, `POST`|
|/redfish/v1/AggregationService/Aggregates/{aggregateId}|`GET`, `DELETE`|
|/redfish/v1/AggregationService/Aggregates/{aggregateId}/Actions/Aggregate.AddElements|`POST`|
|/redfish/v1/AggregationService/Aggregates/{aggregateId}/Actions/Aggregate.Reset|`POST`|
|/redfish/v1/AggregationService/Aggregates/{aggregateId}/Actions/Aggregate.SetDefaultBootOrder|`POST`|
|/redfish/v1/AggregationService/Aggregates/{aggregateId}/Actions/Aggregate.RemoveElements|`POST`|
|/redfish/v1/AggregationService/ConnectionMethods|`GET`|
|/redfish/v1/AggregationService/ConnectionMethods/{connectionmethodsId}|`GET`|

|Systems||
|-------|--------------------|
|/redfish/v1/Systems|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}|`GET`, `PATCH`|
|/redfish/v1/Systems/{ComputerSystemId}/Memory|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/Memory/{memoryId}|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/MemoryDomains|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/NetworkInterfaces|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces/{id}|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/Bios|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/SecureBoot|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/PCIeDevices/{PCIeDeviceId}|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/Storage|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Drives/{driveId}|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes|`GET` , `POST`|
|/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/Capabilities|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/{volumeId}|`GET`, `DELETE`|
|/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageControllerId}/StoragePools|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageControllerId}/StoragePools/{storagepool_Id}|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/AllocatedVolumes|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/AllocatedVolumes/{allocatedvolumes_Id}|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/CapacitySources/{capacitysources_Id}/ProvidingDrives|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/CapacitySources/{capacitysources_Id}/ProvidingDrives/{providingdrives_id}|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/Processors|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/Processors/{id}|`GET`|
|/redfish/v1/Systems?filter={searchKeys*}%20{conditionKeys}%20{value/regEx}|`GET`|
|/redfish/v1/Systems/{ComputerSystemId}/Bios/Settings<br> |`GET`, `PATCH`|
|/redfish/v1/Systems/{ComputerSystemId}/Actions/ComputerSystem.Reset|`POST`|
|/redfish/v1/Systems/{ComputerSystemId}/Actions/ComputerSystem.SetDefaultBootOrder|`POST`|

|Chassis||
|-------|--------------------|
|/redfish/v1/Chassis|`GET`, `POST`|
|/redfish/v1/Chassis/{chassisId}|`GET`, `PATCH`, `DELETE`|
|/redfish/v1/Chassis/{chassisId}/Thermal|`GET`|
|/redfish/v1/Chassis/{chassisId}/Power|`GET`|
|/redfish/v1/Chassis/{chassisId}/NetworkAdapters|`GET`|
|/redfish/v1/Chassis/{ChassisId}/NetworkAdapters/{networkadapterId}|`GET`|

|Managers||
|-------|--------------------|
|/redfish/v1/Managers|`GET`|
|/redfish/v1/Managers/{managerId}|`GET`|
|/redfish/v1/Managers/{managerId}/EthernetInterfaces|`GET`|
|/redfish/v1/Managers/{managerId}/HostInterfaces|`GET`|
|/redfish/v1/Managers/{managerId}/LogServices|`GET`|
|/redfish/v1/Managers/{managerId}/NetworkProtocol|`GET`|
|/redfish/v1/Managers/{ManagerId}/VirtualMedia|`GET`|
|/redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}| `GET`  |
|/redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}/Actions/VirtualMedia.InsertMedia|`POST`|
|/redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}/Actions/VirtualMedia.EjectMedia|`POST`|
|/redfish/v1/Managers/{ManagerId}/RemoteAccountService|`GET`|
|/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Accounts|`GET`, `POST`|
|/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Accounts/{AccountID}|`GET`, `PATCH`, `DELETE`|
|/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Roles|`GET`|
|/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Roles/{Roleid}|`GET`|

|UpdateService||
|-------|--------------------|
|/redfish/v1/UpdateService|`GET`|
|/redfish/v1/UpdateService/FirmwareInventory|`GET`|
|/redfish/v1/UpdateService/FirmwareInventory/{inventoryId}|`GET`|
|/redfish/v1/UpdateService/SoftwareInventory|`GET`|
|/redfish/v1/UpdateService/SoftwareInventory/{inventoryId}|`GET`|
|/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate|`POST`|
|/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate|`POST`|

|EventService||
|-------|--------------------|
|/redfish/v1/EventService|`GET`|
|/redfish/v1/EventService/Subscriptions|`POST`, `GET`|
|/redfish/v1/EventService/Actions/EventService.SubmitTestEvent|`POST`|
|/redfish/v1/EventService/Subscriptions/{subscriptionId}|`GET`, `DELETE`|

|LicenseService||
|-------|--------------------|
|/redfish/v1/LicenseService|`GET`|
|/redfish/v1/LicenseService/Licenses/|`GET`,`POST`|
|/redfish/v1/LicenseService/Licenses/{LicenseId}|`GET`|

|Fabrics||
|-------|--------------------|
|/redfish/v1/Fabrics|`GET`|
|/redfish/v1/Fabrics/{fabricId}|`GET`|
|/redfish/v1/Fabrics/{fabricId}/Switches|`GET`|
|/redfish/v1/Fabrics/{fabricId}/Switches/{switchId}|`GET`|
|/redfish/v1/Fabrics/{fabricId}/Switches/{switchId}/Ports<br> |`GET`|
|/redfish/v1/Fabrics/{fabricId} /Switches/{switchId}/Ports/{portid}<br> |`GET`|
|/redfish/v1/Fabrics/{fabricId}/Zones|`GET`, `POST`|
|/redfish/v1/Fabrics/{fabricId}/Zones/{zoneId}|`GET`, `PATCH`, `DELETE`|
|/redfish/v1/Fabrics/{fabricId}/AddressPools|`GET`, `POST`|
|/redfish/v1/Fabrics/{fabricId}/AddressPools/{addresspoolid}|`GET`, `DELETE`|
|/redfish/v1/Fabrics/{fabricId}/Endpoints|`GET`, `POST`|
|/redfish/v1/Fabrics/{fabricId}/Endpoints/{endpointId}|`GET`, `DELETE`|

|TaskService||
|-------|--------------------|
|/redfish/v1/TaskService|`GET`|
|/redfish/v1/TaskService/Tasks|`GET`|
|/redfish/v1/TaskService/Tasks/{taskId}|`GET`, `DELETE`|
| /redfish/v1/TaskService/Tasks/{taskId}/SubTasks |`GET`|
| /redfish/v1/TaskService/Tasks/{taskId}/SubTasks/ {subTaskId} |`GET`|

| TelemetryService                                             |                |
| ------------------------------------------------------------ | -------------- |
| /redfish/v1/TelemetryService                                 | `GET`          |
| /redfish/v1/TelemetryService/MetricDefinitions               | `GET`          |
| /redfish/v1/TelemetryService/MetricDefinitions/{MetricDefinitionId} | `GET`          |
| /redfish/v1/TelemetryService/MetricReportDefinitions         | `GET`          |
| /redfish/v1/TelemetryService/MetricReportDefinitions/{MetricReportDefinitionId} | `GET`          |
| redfish/v1/TelemetryService/MetricReports                    | `GET`          |
| /redfish/v1/TelemetryService/MetricReports/{MetricReportId}  | `GET`          |
| /redfish/v1/TelemetryService/Triggers                        | `GET`          |
| /redfish/v1/TelemetryService/Triggers/{TriggerId}            | `GET`, `PATCH` |

|Task monitor||
|-------|--------------------|
|/taskmon/{taskId}|`GET`|

|Registries||
|-------|--------------------|
|/redfish/v1/Registries|`GET`|
|/redfish/v1/Registries/{registryId}|`GET`|
|/redfish/v1/registries/{registryFileId}|`GET`|


## Viewing the list of supported Redfish services

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1` |
|**Description** |This is the URI for the Redfish service root. Perform `GET` on this URI to fetch a list of available Redfish services.|
|**Returns** |All available services in the service root.|
|**Response Code** |`200 OK` |
|**Authentication** |No|


>**curl command**


```curl
curl -i GET 'https://{odimra_host}:{port}/redfish/v1'
```

>**Sample response header**


```
Allow:GET
Link:</redfish/v1/SchemaStore/en/ServiceRoot.json/>; rel=describedby
Date":Fri,15 May 2022 13:55:53 GMT+5m 11s
```

>**Sample response body**


```
{
   "@odata.context": "/redfish/v1/$metadata#ServiceRoot.ServiceRoot",
   "@odata.id": "/redfish/v1/",
   "@odata.type": "#ServiceRoot.v1_11_0.ServiceRoot",
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
   "LicenseService": {
      "@odata.id": "/redfish/v1/LicenseService"
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
   "RedfishVersion": "1.15.1",
   "UUID": "a64fc187-e0e9-4f68-82a8-67a616b84b1d"
}
```

## Modifying configurations for services

You can modify the existing configurations of all Resource Aggregator for ODIM services by editing the configuration file at:

```
odimra/lib-utilities/config/odimra_config.json
```

If Resource Aggregator for ODIM is deployed already, run the following command to apply the latest configurations.

```
python3 odim-controller.py --config /home/${USER}/ODIM/odim-controller/scripts/kube_deploy_nodes.yaml --upgrade odimra-config
```



# Rate limits

It is important to protect the shared services from excessive use to maintain service availability. Rate limits are used to control the rate of requests being sent or received in a network to prevent the frequency of an operation from exceeding specific limits.

Resource Aggregator for ODIM supports the following rate limits:

- Specify time (in milliseconds) to limit multiple resource requests. 
  These resources include the log service entries that take more retrieval time from the Baseboard Management Controller (BMC) servers.
- Limit the number of concurrent API requests being sent per session.
- Limit the number of active sessions per user.

Specify values for `resourceRateLimit`, `requestLimitPerSession`, and `sessionLimitPerUser` in the `kube_deploy_nodes.yaml` deployment configuration file [optional]. By default, the values for these parameters are blank, meaning there is no limit on these numbers, unless specified.

> **Samples**

- **`resourceRateLimit`**: Specify values for the parameter in the following format:

  ```
  resourceRateLimit:
  - /redfish/v1/Systems/{id}/LogServices/SL/Entries:10000
  - /redfish/v1/Systems/{id}/LogServices/IML/Entries:8000
  - /redfish/v1/Managers/{id}/LogServices/IEL/Entries:7000
  ```

  In case of multiple requests for resources, the `503` error code is returned for the specified time (in milliseconds). The response header for this request consists of a property `Retry-after`, which displays time in seconds. After this time, requests are processed with the `200` status code.

  > **Sample response body**

  ```
  {
     "error":{
        "code":"Base.1.13.0.GeneralError",
        "message":"An error has occurred. See ExtendedInfo for more information.",
        "@Message.ExtendedInfo":[
           {
              "@odata.type":"#Message.v1_1_2.Message",
              "MessageId":"Base.1.13.0.GeneralError",
              "Message":"too many requests, retry after some time",
              "Severity":"Critical",
              "Resolution":"Retry after some time"
           }
        ]
     }
  }
  ```

  > **Sample response header**

  ```
  Retry-After: 1
  ```

  > **NOTE:** The value for 'Retry-After' property is in seconds.

- **`requestLimitPerSession`**: Specify the number of concurrent API requests that can be sent per session. If you specify `15` as the value for this parameter, 15 API requests are processed with `200` status code and the remaining concurrent requests triggered from your session return the `503` error code.

  > **Sample response body**

  ```
  {
     "error":{
        "code":"Base.1.11.0.GeneralError",
        "message":"An error has occurred. See ExtendedInfo for more information.",
        "@Message.ExtendedInfo":[
           {
              "@odata.type":"#Message.v1_1_2.Message",
              "MessageId":"Base.1.13.0.GeneralError",
              "Message":"A general error has occurred. See Resolution for information on how to resolve the error, or @Message.ExtendedInfo if Resolution is not provided.",
              "Severity":"Critical",
              "Resolution":"None"
           }
        ]
     }
  }
  ```

- **`sessionLimitPerUser`**: Specify the number of active sessions a user can have. If you specify `10` as the value for this parameter, 10 sessions can be created for a particular user, which return the `201` status code. Beyond this, the `503` error code is returned.

  > **Sample response body**

  ```
  {
     "error":{
        "code":"Base.1.11.0.GeneralError",
        "message":"An error has occurred. See ExtendedInfo for more information.",
        "@Message.ExtendedInfo":[
           {
              "@odata.type":"#Message.v1_1_2.Message",
              "MessageId":"Base.1.11.0.SessionLimitExceeded",
              "Message":"The session establishment failed due to the number of
  simultaneous sessions exceeding the limit of the implementation.",
              "Severity":"Critical",
              "Resolution":"Reduce the number of other sessions before trying to establish the session or increase the limit of simultaneous sessions, if
  supported."
           }
        ]
     }
  }
  ```

  

# Authentication and authorization

##  Authentication methods for Redfish APIs

To keep the HTTP connections secure, Resource Aggregator for ODIM verifies credentials of HTTP requests. If you perform an unauthenticated HTTP operation on resources except the listed ones, you get an HTTP `401 unauthorized` error.

|Resource|URI|Description|
|--------|---|-----------|
|The Redfish service root|`GET` `/redfish` |The URI for the Redfish service root. It returns the version of Redfish services.|
|List of Redfish services|`GET` `/redfish/v1` |It returns a list of available services.|
|$metadata|`GET` `/redfish/v1/$metadata` |The Redfish metadata document.|
|OData|`GET` `/redfish/v1/odata` |The Redfish OData service document.|
| The `Sessions` resource<br> |`POST` `/redfish/v1/SessionService/Sessions` |Creates a Redfish login session.|

To authenticate requests with Redfish services, implement one of the following authentication methods:


-   **HTTP BASIC authentication \(BasicAuth\)** 

    To implement HTTP BASIC authentication:

     1. Generate a `base64` encoded string of `{valid_username_of_odim_userAccount}:{valid_password_of_odim_userAccount}` using the following command:

         ```
        echo -n '{username}:{password}' | base64 -w0
        ```

        Initially, use the username and the password of the default administrator account. Later, you can create additional *[user accounts](#user-accounts)* and use their details to implement authentication.

     2. Provide the base64 encoded string in an HTTP `Authorization:Basic` header as shown in the curl command:

         ```
         curl -i --cacert {path}/rootCA.crt GET\
         -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
         'https://{odimra_host}:{port}/redfish/v1/AccountService'
         ```

-   **Redfish session login authentication (XAuthToken)** 

    1. To implement Redfish session login authentication, create a Redfish login *[session](#sessions)* and obtain an authentication token through session management interface.
       Every session created has an authentication token called `X-AUTH-TOKEN` that is returned in the response header.
    
    2. To authenticate subsequent requests, provide the token in the `X-AUTH-TOKEN` request header.
    
       ```
       curl -i --cacert {path}/rootCA.crt GET \
       -H "X-Auth-Token:{X-Auth-Token}" \
        'https://{odimra_host}:{port}/redfish/v1/AccountService'
       ```

       An `X-AUTH-TOKEN` is valid and the session is available for only 30 minutes, unless you continue to send requests to a Redfish service using this token. An idle session is automatically terminated after the time-out interval.

## Role-based authorization

In Resource Aggregator for ODIM, the roles and privileges control users' access to specific resources. If you perform an HTTP operation on a resource without the required privileges, you encounter an HTTP `403 Forbidden` error.

### **Roles**

Role represents a set of operations that a user is allowed to perform with a defined set of privileges. You can assign a role to a user while creating the user account.

With Resource Aggregator for ODIM, there are two types of defined roles:

-   **Redfish predefined roles** 

    Redfish predefined roles have predefined set of privileges. These privileges cannot be removed or modified. You may assign additional OEM \(custom\) privileges. The following are the default Redfish predefined roles that are available in Resource Aggregator for ODIM:

    -   `Administrator` 

    -   `Operator` 

    -   `ReadOnly` 

-   **User-defined roles** 

    > **PREREQUISITE**: Ensure that a user role is created before assigning it to a user account.
    
    User-defined roles are the custom roles that you can create and assign to a user. The privileges of a user-defined role are configurable. You can select a privilege or a set of privileges to assign to this role while creating the user account.

### **Privileges**

Privilege is a permission to perform an operation or a set of operations within a defined management domain.

The following Redfish-specified privileges can be assigned to the users in Resource Aggregator for ODIM:

-    `ConfigureComponents`—Users with this privilege can configure components managed by the Redfish services in Resource Aggregator for ODIM. This privilege is required to create, update, and delete a resource or a collection of resources exposed by Redfish APIs using HTTP `POST`, `PATCH`, and `DELETE` operations.

 -    `ConfigureManager`—Users with this privilege can configure manager resources.

 -    `ConfigureComponents`—Users with this privilege can configure components managed by the services.

 -    `ConfigureSelf`—Users with this privilege can change the password for their account.

 -    `ConfigureUsers`—Users with this privilege can configure users and their accounts. This privilege is assigned to an `Administrator`. This privilege is required to create, update, and delete user accounts using HTTP `POST`, `PATCH`, and `DELETE` operations.

 -    `Login`—Users with this privilege can log in to the service and read the resources.
This privilege is required to view any resource or a collection of resources exposed by Redfish APIs using HTTP `GET` operation.

#### **Mapping of privileges to roles**

|Roles|Assigned privileges|
|-----|-------------------|
|Administrator (Redfish predefined)| `Login` <br>`ConfigureManager` <br>`ConfigureUsers` <br>`ConfigureComponents` <br>`ConfigureSelf` <br> |
|Operator (Redfish predefined)| `Login` <br>`ConfigureComponents` <br>`ConfigureSelf` <br> |
|ReadOnly (Redfish predefined)| `Login` <br>`ConfigureSelf` <br> |


>**NOTE:** Resource Aggregator for ODIM has a default user account that has all the privileges of an administrator role.


# Sessions

A session represents a window of a user login with a Redfish service. Sessions contain details on the user and the user activities. You can run multiple sessions simultaneously.

Resource Aggregator for ODIM allows you to view, create, and manage user sessions through Redfish APIs.


**Supported APIs**

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/SessionService|`GET`|`Login` |
|/redfish/v1/SessionService/Sessions|`POST`, `GET`|`Login`|
|redfish/v1/SessionService/Sessions/{sessionId}|`GET`, `DELETE`|`Login`, `ConfigureManager`, `ConfigureSelf` |

## Viewing the SessionService root

|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/SessionService` |
|**Description** |This endpoint retrieves JSON schema representing the Redfish `SessionService` root.|
|**Returns** |The properties for the Redfish `SessionService` and the links to the actual list of sessions.|
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
   "@odata.type":"#SessionService.v1_1_8.SessionService",
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
|**Description** |This operation creates a session to implement authentication. Creating a session allows you to create an `X-AUTH-TOKEN` which is then used to authenticate with other services.<br>**NOTE:** It is a good practice to make a note of the following:<br><ul><li>The session authentication token returned in the `X-AUTH-TOKEN` header.</li><li>The session id returned in the `Location` header and the JSON response body.</li></ul><br>You need the session authentication token to authenticate to subsequent requests to the Redfish services and the session id to log out later.|
|**Returns** |<ul><li> An `X-AUTH-TOKEN` header containing session authentication token.</li><li>`Location` header that contains a link to the new session instance.</li><li>The session id and a message in the JSON response body denoting a session creation.</li></ul> |
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
"UserName": "{username}",
"Password": "{password}"
}
```

> **Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|UserName|String (required)|Username of the user account for the session. For the first time, use the username of the default administrator account (admin). Subsequently, when you create other user accounts, you can use the credentials of these accounts to create a session.<br>**NOTE:** This user must have `Login` privilege.|
|Password|String (required)<br> |Password of the user account for the session. For the first time, use the password of the default administrator account. Subsequently, when you create other user accounts, you can use the credentials of these accounts to create a session. |

>**Sample response header**


```
Link:</redfish/v1/SessionService/Sessions/2d2e8ebc-4e7c-433a-bfd6-74dc420886d0/>; rel=self
Location:{odimra_host}:{port}/redfish/v1/SessionService/Sessions/2d2e8ebc-4e7c-433a-bfd6-74dc420886d0
X-Auth-Token:15d0f639-f394-4be7-a8ef-ef9d1df07288
Date:Fri,15 May 2020 14:08:55 GMT+5m 11s
```

>**Sample response body**


```
{
	"@odata.type": "#SessionService.v1_1_8.SessionService",
	"@odata.id": "/redfish/v1/SessionService/Sessions/1a547199-0dd3-42de-9b24-1b801d4a1e63",
	"Id": "1a547199-0dd3-42de-9b24-1b801d4a1e63",
	"Name": "Session Service",
	"Message": "The resource has been created successfully",
	"MessageId": "Base.1.13.0.Created",
	"Severity": "OK",
	"UserName": "{username}"
}
```


## Viewing a list of sessions

|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/SessionService/Sessions` |
|**Description** |This operation lists user sessions.<br>**NOTE:** Only a user with `ConfigureUsers` privilege can view a list of all user sessions.<br>Users with `ConfigureSelf` privilege can view the sessions created only by them.|
|**Returns** |Links to the list of user sessions|
|**Response Code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
               -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
              'https://{odimra_host}:{port}/redfish/v1/SessionService/Sessions'
```

>**Sample response body**

```
{
   "@odata.type":"#SessionCollection.SessionCollection",
   "@odata.id":"/redfish/v1/SessionService/Sessions",
   "@odata.context":"/redfish/v1/$metadata#SessionCollection.SessionCollection",
   "Name":"Session Service",
   "Members@odata.count":3,
   "Members":[
      {
         "@odata.id":"/redfish/v1/SessionService/Sessions/947d9c57-1e0e-4816-9830-c4a2f97d8991",
         "@odata.id":"/redfish/v1/SessionService/Sessions/547d4c57-1e0c-2816-3830-c4a2f97d1981",
         "@odata.id":"/redfish/v1/SessionService/Sessions/747c3c57-1w2e-4656-9550-c4a2f97d8938"
         
      }
   ]
}
```

## Viewing information about a session

|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/SessionService/Sessions/{sessionId}` |
|**Description** |This operation retrieves information about a specific user session.<br>**NOTE:** Only a user with `ConfigureUsers` privilege can view information about any user session.<br>Users with `ConfigureSelf` privilege can view information about the sessions created only by them.|
|**Returns** |JSON schema representing the session|
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
   "@odata.type":"#Session.v1_4_0.Session",
   "@odata.id":"/redfish/v1/SessionService/Sessions/4ee42139-22db-4e2a-97e4-020013248768",
   "Id":"4ee42139-22db-4e2a-97e4-020013248768",
   "Name":"User Session",
   "UserName":"admin"
   "CreatedTime": "2022-06-30T06:32:59Z"
}
```

## Deleting a session

|||
|---------|---------------|
|**Method** | `DELETE` |
|**URI** |`/redfish/v1/SessionService/Sessions/{sessionId}` |
|**Description** |This operation terminates a specific Redfish session when the user logs out.<br>**NOTE**: Users having the `ConfigureSelf` and `ConfigureComponents` privileges can delete only the sessions they created.<br>Only a user with all the Redfish-defined privileges \(Redfish-defined `Administrator` role\) is authorized to delete any user session. |
|**Response Code** |`204 No Content` |
|**Authentication** |Yes|

>**curl command**


```
curl -i -X DELETE \
               -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
              'https://{odimra_host}:{port}/redfish/v1/SessionService/Sessions/{sessionId}'
```


#  User roles and privileges

Resource Aggregator for ODIM allows you to view, create, and manage user roles through Redfish APIs.

**PREREQUISITE:** Only a user with `ConfigureUsers` privilege can perform the operations of user roles and privileges.

**Supported APIs**:

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/AccountService|`GET`|`Login` |
|/redfish/v1/AccountService/Roles|`GET`|`Login` |
|/redfish/v1/AccountService/Roles/{roleId}|`GET`|`Login` |


## Viewing the AccountService root

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
Link:</redfish/v1/SchemaStore/en/AccountService.json>; rel=describedby
Date:Fri,15 May 2020 14:32:09 GMT+5m 12s
```


>**Sample response body**

```
{
   "@odata.type":"#AccountService.v1_11_0.AccountService",
   "@odata.id":"/redfish/v1/AccountService",
   "@odata.context":"/redfish/v1/$metadata#AccountService.AccountService",
   "Id":"AccountService",
   "Name":"Account Service",
   "Status":{
      "State":"Enabled",
      "Health":"OK"
   },
   "ServiceEnabled":true,
   "MinPasswordLength":12,
   "Accounts":{
      "@odata.id":"/redfish/v1/AccountService/Accounts"
   },
   "Roles":{
      "@odata.id":"/redfish/v1/AccountService/Roles"
   }
}
```

## Viewing a list of roles

|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/AccountService/Roles` |
|**Description** |This operation lists available user roles.|
|**Returns** |Links to user role resources|
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
   "Members@odata.count":3,
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
   ]
}
```

## Viewing information about a role


|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/AccountService/Roles/{RoleId}` |
|**Description** |This operation fetches information about a specific user role.|
|**Returns** |JSON schema representing this role. The schema has the details such as id, name, description, assigned privileges, and OEM privileges.|
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
   "@odata.type":"#Role.v1_3_1.Role",
   "@odata.id":"/redfish/v1/AccountService/Roles/ReadOnly",
   "Id":"ReadOnly",
   "Name":"User Role",
   "IsPredefined":true,
   "AssignedPrivileges":[
      "ConfigureSelf",
      "Login"
   ]
}
```



#  User accounts

Resource Aggregator for ODIM allows users to have accounts to configure their actions and restrictions.

Resource Aggregator for ODIM has an administrator user account by default. Create other user accounts by defining a username, a password, and a role for each account. The username and the password are used to authenticate with the Redfish services (using `BasicAuth` or `XAuthToken`).

Resource Aggregator for ODIM exposes Redfish `AccountsService` APIs to view, create and manage user accounts. 


**Supported APIs**:

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/AccountService/Accounts|`POST`, `GET`|`Login`, `ConfigureUsers` |
|/redfish/v1/AccountService/Accounts/{accountId}|`GET`, `DELETE`, `PATCH`|`Login`, `ConfigureUsers`, `ConfigureSelf` |

## Creating a user account

|||
|-------|--------------------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/AccountService/Accounts` |
|**Description** |This operation creates a user account. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can create other user accounts.|
|**Returns** |<ul><li>`Location` header that contains a link to the new account</li><li>JSON schema representing the new account</li></ul> |
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
   "UserName":"{username}",
   "Password":"{password}",
   "RoleId":"{roleId}"
}
```

> **Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Username|String (required)<br> |User name for the user account.|
|Password|String (required)<br> |Password for the user account. Before creating a password, see the *[Password Requirements](#password-requirements)* section.|
|RoleId|String (required)<br> |Role for this account. To know more about roles, see *[User roles and privileges](#role-based-authorization)*. Ensure the `roleId` you want to assign to this user account exists. To check the existing roles, see *[Listing Roles](#listing-roles)*. If you attempt to assign an unavailable role, an HTTP `400 Bad Request` error is displayed.|


### Password requirements

-   Your password must not be same as your username.

-   Your password must be at least 12 characters long and at most 16 characters long.

-   Your password must contain at least one uppercase letter (A-Z), one lowercase letter (a-z), one digit (0-9), and one special character (~!@\#$%^&\*-+\_|(){}:;<\>,.?/).


>**Sample response header**

```
Link:</redfish/v1/AccountService/Accounts/monitor32/>; rel=describedby
Location:/redfish/v1/AccountService/Accounts/monitor32/
Date":Fri,15 May 2020 14:36:14 GMT+5m 11s
```

>**Sample response body**

```
{
   "@odata.type":"#ManagerAccount.v1_9_0.ManagerAccount",
   "@odata.id":"/redfish/v1/AccountService/Accounts/{accountId}",
   "@odata.context":"/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
   "Id":"{accountId}",
   "Name":"Account Service",
   "Message":"The resource has been created successfully",
   "MessageId":"Base.1.13.0.Created",
   "Severity":"OK",
   "UserName":"{Username}",
   "RoleId":"ReadOnly",
   "AccountTypes":[
      "Redfish"
   ],
   "Password":null,
   "Links":{
      "Role":{
         "@odata.id":"/redfish/v1/AccountService/Roles/ReadOnly"
      }
   }
}
```

##  Viewing a list of user accounts

|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/AccountService/Accounts` |
|**Description** |This operation retrieves a list of user accounts.|
|**Returns** |Links to user accounts.|
|**Response Code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Accounts'
```

>**Sample response body**

```
{
   "@odata.type":"#ManagerAccountCollection.ManagerAccountCollection",
   "@odata.id":"/redfish/v1/AccountService/Accounts",
   "@odata.context":"/redfish/v1/$metadata#ManagerAccountCollection.ManagerAccountCollection",
   "Name":"Account Service",
   "Members@odata.count":1,
   "Members":[
      {
         "@odata.id":"/redfish/v1/AccountService/Accounts/admin"
      }
   ]
}
```




##  Viewing information about an account

|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/AccountService/Accounts/{accountId}` |
|**Description** |This operation fetches information about a specific user account. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can view information about a user account.|
|**Returns** |JSON schema representing the user account.|
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
   "@odata.type":"#ManagerAccount.v1_9_0.ManagerAccount",
   "@odata.id":"/redfish/v1/AccountService/Accounts/{accountId}",
   "@odata.context":"/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
   "Id":"{accountId}",
   "Name":"Account Service",
   "UserName":"{Username}",
   "RoleId":"ReadOnly",
   "AccountTypes":[
      "Redfish"
   ],
   "Password":null,
   "Links":{
      "Role":{
         "@odata.id":"/redfish/v1/AccountService/Roles/ReadOnly"
      }
   }
}
```

## Updating a user account

|||
|---------|---------------|
|**Method** | `PATCH` |
|**URI** |`/redfish/v1/AccountService/Accounts/{accountId}` |
|**Description** |This operation updates user account details (`password`, and `RoleId`). To modify account details, add them in the request payload (as shown in the sample request body) and perform `PATCH` on the mentioned URI. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can modify other user accounts. Users with `ConfigureSelf` privilege can modify only their own accounts.|
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
   "RoleId":"{roleId}"
}
' \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Accounts/{accountId}'
```


>**Sample request body**

```
{ 
   "Password":"{new_password}",
   "RoleId":"{roleId}"
}
```

>**Sample response header**

```
Link:</redfish/v1/AccountService/Accounts/monitor32/>; rel=describedby
Location:/redfish/v1/AccountService/Accounts/monitor32/
Date":Fri,15 May 2020 14:36:14 GMT+5m 11s
```

>**Sample response body**

```
{
   "@odata.type":"#ManagerAccount.v1_9_0.ManagerAccount",
   "@odata.id":"/redfish/v1/AccountService/Accounts/{accountId}",
   "@odata.context":"/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
   "Id":"{accountId}",
   "Name":"Account Service",
   "Message":"The account was successfully modified.",
   "MessageId":"Base.1.13.0.AccountModified",
   "Severity":"OK",
   "UserName":"{Username}",
   "RoleId":"ReadOnly",
   "AccountTypes":[
      "Redfish"
   ],
   "Password":null,
   "Links":{
      "Role":{
         "@odata.id":"/redfish/v1/AccountService/Roles/ReadOnly"
      }
   }
}
```

## Deleting a user account

|||
|---------|---------------|
|**Method** | `DELETE` |
|**URI** |`/redfish/v1/AccountService/Accounts/{accountId}` |
|**Description** |This operation deletes a user account.|
|**Response Code** |`204 No Content` |
|**Authentication** |Yes|

>**curl command**

```
curl  -i -X DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Accounts/{accountId}'
```



#  Resource aggregation and management


The resource aggregator allows you to add southbound infrastructure to its database, create resource collections, and perform actions in combination on these collections. It exposes the Redfish `AggregationService` APIs to achieve the following:

-   Adding a resource and building its inventory

-   Resetting one or more resources

-   Changing the boot path of one or more resources to default settings

-   Removing a resource from the inventory, which is no longer managed


All aggregation actions are performed as *[tasks](#tasks)* in Resource Aggregator for ODIM. The actions performed on a group of resources (resetting or changing the boot order to default settings) are carried out as a set of subtasks.

**Supported endpoints**

|API URI|Supported operations|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/AggregationService|`GET`|`Login` |
| /redfish/v1/AggregationService/AggregationSources<br> |`GET`, `POST`|`Login`, `ConfigureManager` |
|/redfish/v1/AggregationService/AggregationSources/{aggregationSourceId}|`GET`, `PATCH`, `DELETE`|`Login`, `ConfigureManager` |
|/redfish/v1/AggregationService/Actions/AggregationService.Reset|`POST`|`ConfigureManager`, `ConfigureComponents` |
|/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder|`POST`|`ConfigureManager`, `ConfigureComponents` |
|/redfish/v1/AggregationService/Aggregates|`GET`, `POST`|`Login`, `ConfigureComponents`, `ConfigureManager` |
|/redfish/v1/AggregationService/Aggregates/{aggregateId}|`GET`, `DELETE`|`Login`, `ConfigureComponents`, `ConfigureManager` |
|/redfish/v1/AggregationService/Aggregates/{aggregateId}/Actions/Aggregate.AddElements|`POST`|`ConfigureComponents`, `ConfigureManager` |
|/redfish/v1/AggregationService/Aggregates/{aggregateId}/Actions/Aggregate.Reset|`POST`|`ConfigureComponents`, `ConfigureManager` |
|/redfish/v1/AggregationService/Aggregates/{aggregateId}/Actions/Aggregate.SetDefaultBootOrder|`POST`|`ConfigureComponents`, `ConfigureManager` |
|/redfish/v1/AggregationService/Aggregates/{aggregateId}/Actions/Aggregate.RemoveElements|`POST`|`ConfigureComponents`, `ConfigureManager` |
|/redfish/v1/AggregationService/ConnectionMethods|`GET`|`Login`|
|/redfish/v1/AggregationService/ConnectionMethods/{connectionmethodsId}|`GET`|`Login`|

## Viewing the AggregationService root
|||
|-----|-------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/AggregationService` |
|<strong>Description</strong> |This endpoint retrieves JSON schema representing the aggregation service root.|
|<strong>Returns</strong> |Properties for the service and a list of actions you can perform using this service|
|<strong>Response Code</strong> |On success, `200 OK` |
|<strong>Authentication</strong> |Yes|

 **curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odim_host}:{port}/redfish/v1/AggregationService'
```

>**Sample response header** 

```
Allow:GET
Date:Sun,17 May 2020 14:26:49 GMT+5m 14s
Link:</redfish/v1/SchemaStore/en/AggregationService.json>; rel=describedby
```

>**Sample response body** 

```
{
   "@odata.context":"/redfish/v1/$metadata#AggregationService.AggregationService",
   "Id":"AggregationService",
   "@odata.id":"/redfish/v1/AggregationService",
   "@odata.type":"#AggregationService.v1_0_1.AggregationService",
   "Name":"AggregationService",
   "Description":"AggregationService",
   "Actions":{
      "#AggregationService.Reset":{
            "target": "/redfish/v1/AggregationService/Actions/AggregationService.Reset/",
            "@Redfish.ActionInfo": "/redfish/v1/AggregationService/ResetActionInfo"
      },
      "#AggregationService.SetDefaultBootOrder":{
            "target": "/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder/",
            "@Redfish.ActionInfo": "/redfish/v1/AggregationService/SetDefaultBootOrderActionInfo"
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
|**Returns** |A list of links to all the available connection method resources|
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
   "@odata.type":"#ConnectionMethodCollection.ConnectionMethodCollection",
   "@odata.id":"/redfish/v1/AggregationService/ConnectionMethods",
   "@odata.context":"/redfish/v1/$metadata#ConnectionMethodCollection.ConnectionMethodCollection",
   "Name":"Connection Methods",
   "Members@odata.count":3,
   "Members":[
      {
         "@odata.id":"/redfish/v1/AggregationService/ConnectionMethods/18312728-d687-4cd3-b7e6-27a1cbd3b2e3"
      },
      {
         "@odata.id":"/redfish/v1/AggregationService/ConnectionMethods/a1b31c57-dcaa-4d7c-b405-9244a24b502c"
      },
      {
         "@odata.id":"/redfish/v1/AggregationService/ConnectionMethods/3077f07c-6496-4503-a3c2-b02108a54000"
      }
   ]
}
```

### Viewing a connection method

|||
|--------|---------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/AggregationService/ConnectionMethods/{connectionmethodsId}` |
|**Description** |This operation retrieves information about a specific connection method.|
|**Returns** |JSON schema representing this connection method|
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
      "@odata.type":"#ConnectionMethod.v1_0_0.ConnectionMethod",
      "@odata.id":"/redfish/v1/AggregationService/ConnectionMethods/c27575d2-052d-4ce9-8be1-978cab002a0f",
      "@odata.context":"/redfish/v1/$metadata#ConnectionMethod.v1_0_0.ConnectionMethod",
      "Id":"c27575d2-052d-4ce9-8be1-978cab002a0f",
      "Name":"Connection Method",
      "ConnectionMethodType":"Redfish",
      "ConnectionMethodVariant":"Compute:BasicAuth:GRF_v1.0.0",
      "Links":{
            "AggregationSources":[
         {
            "@odata.id":"/redfish/v1/AggregationService/AggregationSources/839c212d-9ab2-4868-8767-1bdcc0ce862c"
         },
         {
            "@odata.id":"/redfish/v1/AggregationService/AggregationSources/3536bb46-a023-4e3a-ac1a-7528cc18b660"
         }
      ]      
   }   
}
```

>**Connection method properties**

|Parameter|Type|Description|
|---------|----|-----------|
|ConnectionMethodType|String| The type of this connection method.<br> For possible property values, see "Connection method types" table.<br> |
|ConnectionMethodVariant|String|The variant of connection method. For more information, see *[Connection method variants](#connection-method-variants)*.|
|Links {|Object|Links to other resources that are related to this connection method.|
|AggregationSources [ {<br> @odata.id<br> } ] |Array|An array of links to the `AggregationSources` resources that use this connection method.|

 >**Connection method types**

|String|Description|
|------|-----------|
| IPMI15<br> | IPMI 1.5 connection method |
| IPMI20<br> | IPMI 2.0 connection method |
| NETCONF<br> | Configuration Protocol |
| OEM<br> | OEM connection method |
| Redfish<br> | Redfish connection method |
| SNMP<br> | Simple Network Management Protocol |

#### Connection method variants

A connection method variant provides details about a plugin and is displayed in the following format:

```
PluginType:PreferredAuthType:PluginID_Firmwareversion
```

It consists of the following parameters:

- **PluginType**
   The string that represents the type of the plugin.<br>Possible values: Compute, Storage, and Fabric. 
- **PreferredAuthType:**   
   Preferred authentication method to connect to the plugin - BasicAuth or XAuthToken.  
- **PluginID_Firmwareversion:**
   The id of the plugin along with the version of the firmware. To know the plugin ids for the supported plugins, see *Mapping of plugins and plugin Ids* table.<br>
   Supported values: `GRF_v1.0.0` and `URP_v1.0.0`<br>

**Examples**:

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
|<strong>Returns</strong> |<ul><li>`Location` URI of the task monitor associated with this operation in the response header. See `Location` URI in *Sample response header (HTTP 202 status)*.</li><li>Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See *Sample response body (HTTP 202 status)*.</li><li>On successful completion:<ul><li>The aggregation source Id, the IP address, the username, and other details of the added plugin in the JSON response body.</li><li> A link (having the aggregation source id) to the added plugin in the `Location` header. See `Location` URI in *Sample response header (HTTP 201 status)*.</li></ul></li></ul>  |
|<strong>Response Code</strong> |`202 Accepted` On success, `201 Created`|
|<strong>Authentication</strong> |Yes|

**Usage information**

Perform HTTP `POST` on the mentioned URI with a request body specifying a connection method to use for adding the plugin. To know about connection methods, see *[Connection methods](#connection-methods)*.
A Redfish task is created and you receive a link to the *[task monitor](#viewing-a-task-monitor)* associated with it.
To know the progress of this operation, perform HTTP `GET` on the task monitor returned in the response header (until the task is complete).

After the plugin is successfully added as an aggregation source, it is also be available as a manager resource at `/redfish/v1/Managers`.


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
   "HostName":"{plugin_host}:45007",
   "UserName":"admin",
   "Password":"Plug!n12$4",
   "Links":{
      "ConnectionMethod": {
         "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/a171e66c-b4a8-137f-981b-1c07ddfeacbb"
      }
   }
}
```

> **Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|HostName|String (required)<br> |FQDN of the resource aggregator server and port of a system where the plugin is installed. The default port for the Generic Redfish Plugin is `45001`.<br>The default port for the URP is `45003`.<br> If you are using a different port, ensure that the port is greater than `45000`.<br> **IMPORTANT**: If you have set the `VerifyPeer` property to false in `/etc/plugin_config/config.json`, you can use IP address of the system where the plugin is installed as `HostName`.<br>|
|UserName|String (required)<br> |The plugin username.|
|Password|String (required)<br> |The plugin password.|
|Links{|Object (required)<br> |Links to other resources that are related to this resource.|
|ConnectionMethod|Array (required)|Links to the connection method that are used to communicate with this endpoint: `/redfish/v1/AggregationService/AggregationSources`. To know which connection method to use, do the following:<ul><li>Perform HTTP `GET` on: `/redfish/v1/AggregationService/ConnectionMethods`.<br>You will receive a list of  links to available connection methods.</li><li>Perform HTTP `GET` on each link. Check the value of the `ConnectionMethodVariant` property in the JSON response. Choose a connection method having the details of the plugin of your choice.<br>For example, the `ConnectionMethodVariant` property for the GRF plugin displays the following value:<br>`Compute:BasicAuth:GRF_v1.0.0` <br>For more information, see the "connection method properties" table in *[Viewing a connection method](#viewing-a-connection-method)*</li></ul>|

>**Sample response header (HTTP 202 status)**

```
Location:/taskmon/task85de4003-8757-4c7d-942f-55eaf7d6812a
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes
```

>**Sample response header (HTTP 201 status)**

```
date:"Wed",02 Sep 2020 06:50:43 GMT+7m 2s
link:/v1/AggregationService/AggregationSources/be626e78-7a8a-4b99-afd2-b8ed45ef3d5a.1/>; rel=describedby
location:/redfish/v1/AggregationService/AggregationSources/be626e78-7a8a-4b99-afd2-b8ed45ef3d5a.1
```

>**Sample response body (HTTP 202 status)**

```
{
   "@odata.type":"#Task.v1_6_0.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4003-8757-4c7d-942f-55eaf7d6812a",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task85de4003-8757-4c7d-942f-55eaf7d6812a",
   "Name":"Task task85de4003-8757-4c7d-942f-55eaf7d6812a",
   "Message":"The task with id task85de4003-8757-4c7d-942f-55eaf7d6812a has started.",
   "MessageId":"TaskEvent.1.0.3.TaskStarted",
   "MessageArgs":[
      "task85de4003-8757-4c7d-942f-55eaf7d6812a"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```


>**Sample response body (HTTP 201 status)**

```
{
   "@odata.type":"#AggregationSource.v1_2_0.AggregationSource",
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

> **PREREQUISITE**: Generate and import certificate for the server.

### Generating and importing certificate

1. Obtain the CSR cert from the Base Management Controller (BMC) server.

2. To get the required information, read the CSR using the following command:

   ```
   openssl req -text -noout –-in BMC.csr
   ```
   
3.  Create a file called `cert.conf` and copy the following content to it.

   ```
   [req]
   default_bits = <Key Length 3072>
   encrypt_key = no
   default_md = <Digest Algorithm sha256/sha512>
   prompt = no
   utf8 = yes
   distinguished_name = req_distinguished_name
   req_extensions = v3_req
   
   [req_distinguished_name]
   C = <Country>
   ST = <State>
   L = <Location>
   O = <Organization>
   OU = <Organization Unit>
   CN = <Common Name>
   
   [v3_req]
   subjectKeyIdentifier = hash
   authorityKeyIdentifier = keyid:always,issuer:always
   keyUsage = critical, nonRepudiation, digitalSignature,
   keyEncipherment
   extendedKeyUsage = clientAuth, serverAuth
   subjectAltName = @alt_names
   
   [alt_names]
   DNS.1 = <Server DNS 1>
   IP.1 = <Server IP 1>
   ```
   
4. In this file, update the CSR details obtained in step 2 and ensure that the DNS and IP addresses (IPv4 or IPv6) of the BMC are configured in `[alt_names]`.

5. Run the following command to generate the certificate.

   ```
   openssl x509 -req -days 365 --in BMC.csr \
   -CA rootCA.crt -CAkey rootCA.key -CAcreateserial -out\
   BMC.crt -extensions v3_req -extfile cert.conf
   ```

   > **NOTE**: Copy `rootCA.key` and `rootCA.cert` from `<OdimCertspath>`. 
   >  `<odimCertsPath>` is the path specified for the `odimCertsPath` parameter in the `kube_deploy_nodes.yaml` file.


7. To verify that all the parameters passed as input are present in the generated certificate, run the following command:

   ```
   openssl x509 -text -noout --in BMC.crt
   ```
   
8. Open the generated certificate to copy its content by running the following command:
   
   ```
   cat BMC.crt
   ```
   
   The content of the certificate file is displayed.
   
9. Import `BMC.crt` in the BMC server.

   |                                 |                                                              |
   | ------------------------------- | ------------------------------------------------------------ |
   | <strong>Method</strong>         | `POST`                                                       |
   | <strong>URI</strong>            | `/redfish/v1/AggregationService/AggregationSources`          |
   | <strong>Description</strong>    | This operation creates an aggregation source for a BMC, discovers information, and performs a detailed inventory of it.<br> The `AggregationSource` schema provides information about a BMC such as its IP address, username, password, and so on.<br> This operation is performed in the background as a Redfish task.<br> |
   | <strong>Returns</strong>        | <ul><li>`Location` URI of the task monitor associated with this operation in the response header. See `Location` URI in *Sample response header (HTTP 202 status)*.</li><li>Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See *Sample response body (HTTP 202 status)*.</li><li>On successful completion:<ul><li>The aggregation source id, the IP address, the username, and other details of the added BMC in the JSON response body.</li><li>A link (having the aggregation source id) to the added BMC in the `Location` header. See `Location` URI in *Sample response header (HTTP 201 status)*.</li></ul></li></ul> |
   | <strong>Response Code</strong>  | On success, `202 Accepted`<br>On successful completion of the task, `201 Created` <br> |
   | <strong>Authentication</strong> | Yes                                                          |

**Usage information**

1. Perform HTTP `POST` on the mentioned URI with a request body specifying a connection method to use for adding the BMC. To know about connection methods, see *[Connection methods](#connection-methods)*.			
   A Redfish task is created and you will receive a link to the *[task monitor](#viewing-a-task-monitor)* associated with it.

2. To know the progress of this operation, perform HTTP `GET` on the task monitor returned in the response header (until the task is complete).
   When the task is successfully complete, you will receive aggregation source id of the added BMC. 
3. Save it because you need to identify it in the resource inventory later.

After the server is successfully added as an aggregation source, it will also be available as a computer system resource at `/redfish/v1/Systems/` and a manager resource at `/redfish/v1/Managers/`.

<blockquote>NOTE: Along with the UUID of the server, check the BMC address to ensure the server isn't already present.</blockquote>

To view the list of links to computer system resources, perform HTTP `GET` on `/redfish/v1/Systems/`. Each link contains `ComputerSystemId` of a specific BMC. For more information, see *[Collection of computer systems](#collection-of-computer-systems)*.

 `ComputerSystemId` is unique information about the BMC specified by Resource Aggregator for ODIM. It is represented as `<UUID:n>`, where `UUID` is the aggregation source id of the BMC. Save it as it is required to perform subsequent actions such as `delete, reset`, and `setdefaultbootorder` on this BMC.


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
   "HostName":"{IPv4_address}",
   "UserName":"admin",
   "Password":"{BMC_password}",
   "Links":{
      "ConnectionMethod": {
         "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/d172e66c-b4a8-437c-981b-1c07ddfeacaa"
      }
   }
}
```

>**Sample request body**

```
{
   "HostName":"{IPv6_address}",
   "UserName":"admin",
   "Password":"{BMC_password}",
   "Links":{
      "ConnectionMethod":{
         "@odata.id":"/redfish/v1/AggregationService/ConnectionMethods/ae910e17-4953-4495-95a9-436bf35fe8e4"
      }
   }
}
```

> **Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|HostName|String (required)<br> |A valid IPv4 or IPv6 address, or hostname of the BMC.|
|UserName|String (required)<br> |The username of the BMC administrator account.|
|Password|String (required)<br> |The password of the BMC administrator account.|
|Links {|Object (required)<br> |Links to other resources that are related to this resource.|
|ConnectionMethod|Array (required)|Links to the connection methods that are used to communicate with this endpoint: `/redfish/v1/AggregationService/AggregationSources`. To know which connection method to use, do the following:<ul><li>Perform HTTP `GET` on: `/redfish/v1/AggregationService/ConnectionMethods`.<br>You will receive a list of  links to available connection methods.</li><li>Perform HTTP `GET` on each link. Check the value of the `ConnectionMethodVariant` property in the JSON response.</li><li>The `ConnectionMethodVariant` property displays the details of a plugin. Choose a connection method having the details of the plugin of your choice.<br> Example: For GRF plugin, the `ConnectionMethodVariant` property displays the following value:<br>`Compute:BasicAuth:GRF:1.0.0`</li></ul>|

>**Sample response header (HTTP 202 status)**

```
Location:/taskmon/task4aac9e1e-df58-4fff-b781-52373fcb5699
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes
```

>**Sample response header (HTTP 201 status)**

```
date:"Wed",02 Sep 2020 06:50:43 GMT+7m 2s
link:/v1/AggregationService/AggregationSources/0102a4b5-03db-40be-ad39-71e3c9f8280e/>; rel=describedby
location:/redfish/v1/AggregationService/AggregationSources/0102a4b5-03db-40be-ad39-71e3c9f8280e
```

>**Sample response body (HTTP 202 status)**

```
{
   "@odata.type":"#Task.v1_6_0.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "Name":"Task task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "Message":"The task with id task4aac9e1e-df58-4fff-b781-52373fcb5699 has started.",
   "MessageId":"TaskEvent.1.0.3.TaskStarted",
   "MessageArgs":[
      "task4aac9e1e-df58-4fff-b781-52373fcb5699"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```

>**Sample response body (HTTP 201 status)**
```
 {
   "@odata.type":"#AggregationSource.v1_2_0.AggregationSource",
   "@odata.id":"/redfish/v1/AggregationService/AggregationSources/26562c7b-060b-4fd8-977e-94b1a535f3fb",
   "@odata.context":"/redfish/v1/$metadata#AggregationSource.AggregationSource",
   "Id":"26562c7b-060b-4fd8-977e-94b1a535f3fb",
   "Name":"Aggregation Source",
   "HostName":"{IPv4_address}",
   "UserName":"admin",
    "Links":{
      "ConnectionMethod": {
         "@odata.id": "/redfish/v1/AggregationService/ConnectionMethods/d172e66c-b4a8-437c-981b-1c07ddfeacaa"
      }
   }
}
```

>**Sample response body (HTTP 201 status) - IPv6 address**

```
 {
   "@odata.type":"#AggregationSource.v1_2_0.AggregationSource",
   "@odata.id":"/redfish/v1/AggregationService/AggregationSources/26562c7b-060b-4fd8-977e-94b1a535f3fb",
   "@odata.context":"/redfish/v1/$metadata#AggregationSource.AggregationSource",
   "Id":"26562c7b-060b-4fd8-977e-94b1a535f3fb",
   "Name":"Aggregation Source",
   "HostName":"{IPv6_address}",
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
|<strong>Returns</strong> |Links of the available aggregation sources|
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
   "@odata.id":"/redfish/v1/AggregationService/AggregationSources",
   "@odata.context":"/redfish/v1/$metadata#AggregationSourceCollection.AggregationSourceCollection",
   "Name":"Aggregation Source",
   "Members@odata.count":2,
   "Members":[
      {
         "@odata.id":"/redfish/v1/AggregationService/AggregationSources/839c212d-9ab2-4868-8767-1bdcc0ce862c"
      },
      {
         "@odata.id":"/redfish/v1/AggregationService/AggregationSources/3536bb46-a023-4e3a-ac1a-7528cc18b660.1"
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
|<strong>Returns</strong> |JSON schema representing this aggregation source|
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
   "@odata.type":"#AggregationSource.v1_2_0.AggregationSource",
   "@odata.id":"/redfish/v1/AggregationService/AggregationSources/839c212d-9ab2-4868-8767-1bdcc0ce862c",
   "@odata.context":"/redfish/v1/$metadata#AggregationSource.AggregationSource",
   "Id":"839c212d-9ab2-4868-8767-1bdcc0ce862c",
   "Name":"Aggregation Source",
   "HostName":"{IPv4_address}",
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
|<strong>Description</strong> |This operation updates the username, password, and IP address or hostname of a specific BMC in the resource aggregator inventory.<br> |
|<strong>Returns</strong> |Updated JSON schema of this aggregation source|
|<strong>Response Code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i PATCH \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
  "HostName": "{IPv4_address}",
  "UserName": "{username}",
  "Password": "{password}"
}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/AggregationSources/{AggregationSourceId}'
```

>**Sample request body**

```
{
  "HostName": "{IPv4_address}",
  "UserName": "{username}",
  "Password": "{password}"
}
```

>**Sample response body**

```
{
   "@odata.type":"#AggregationSource.v1_2_0.AggregationSource",
   "@odata.id":"/redfish/v1/AggregationService/AggregationSources/839c212d-9ab2-4868-8767-1bdcc0ce862c.1",
   "@odata.context":"/redfish/v1/$metadata#AggregationSource.AggregationSource",
   "Id":"839c212d-9ab2-4868-8767-1bdcc0ce862c.1",
   "Name":"Aggregation Source",

   "HostName":"{IPv4_address}",
   "UserName":"{username}",
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
|<strong>Returns</strong> |- `Location` URI of the task monitor associated with this operation (task) in the response header. See `Location` URI in *Sample response header (HTTP 202 status)*.<br><br>-   Link to the task and the task id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See the task URI and the task id  in *Sample response body (HTTP 202 status)*.<br>**IMPORTANT**: Make a note of the task id. If the task completes with an error, it is required to know which subtask has failed. To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`. <br>-  Upon the completion of the reset operation, you receive a success message in the response body. See *Sample response body (HTTP 200 status)*.|
|<strong>Response code</strong> |On success, `202 Accepted`.<br> On successful completion of the task, `200 OK`.|
|<strong>Authentication</strong> |Yes|

**Usage information**

To know the progress of this action, perform HTTP `GET` on the *[task monitor](#viewing-a-task-monitor)* returned in the response header (until the task is complete).

To get the list of subtask URIs, perform HTTP `GET` on the task URI returned in the JSON response body. See *Sample response body (HTTP 202 status)*. The JSON response body of each subtask contains a link to the task monitor associated with it. To know the progress of the reset operation (subtask) on a specific server, perform HTTP `GET` on the task monitor associated with the respective subtask. See the link to the task monitor  in "Sample response body (subtask)".

You can perform reset on a group of servers by specifying multiple target URIs in the request.


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
      "/redfish/v1/Systems/{ComputerSystemId}",
      "/redfish/v1/Systems/{ComputerSystemId2}"
   ]
}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Actions/AggregationService.Reset'


```

>**Sample request body**

```
{
   "BatchSize":2,
   "DelayBetweenBatchesInSeconds":1,
   "ResetType":"ForceRestart",
   "TargetURIs":[
      "/redfish/v1/Systems/97d08f36-17f5-5918-8082-f5156618f58d.1",
      "/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b9.1"
   ]
}
```

> **Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|BatchSize|Integer (optional)<br> |The number of elements to be reset at a time in each batch.|
|DelayBetweenBatchesInSeconds|Integer (seconds) (optional)<br> |The delay among the batches of elements being reset.|
|ResetType|String (required)<br> |The type of reset to be performed. For possible values, see *Reset type*. If the value is not supported by the target server, you receive an HTTP `400 Bad Request` error.|
|TargetURIs|Array (required)<br> |The URI of the target for `Reset`. Example: `"/redfish/v1/Systems/{ComputerSystemId}"` |

> **Reset type**

|String|Description|
|------|-----------|
|ForceOff|Turn off the unit immediately (non-graceful shutdown).|
|ForceRestart|Perform an immediate (non-graceful) shutdown, followed by a restart of the system.|
|GracefulRestart|Perform a graceful shutdown followed by a restart of the system.|
|GracefulShutdown|Perform a graceful shutdown. Graceful shutdown involves shutdown of the operating system followed by the power off of the physical server.|
|Nmi|Generate a Diagnostic Interrupt (usually an NMI on x86 systems) to cease normal operations, perform diagnostic actions, and halt the system.|
|On|Turn on the unit.|
|PowerCycle|Perform a power cycle of the unit.|
|PushPowerButton|Simulate the pressing of the physical power button on this unit.|

>**Sample response header** (HTTP 202 status)

```
Location:/taskmon/task85de4103-8757-4c7d-942f-55eaf7d6412a
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes
```

>**Sample response body** (HTTP 202 status)

```
{
   "@odata.type":"#Task.v1_6_0.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4103-8757-4c7d-942f-55eaf7d6412a",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task85de4103-8757-4c7d-942f-55eaf7d6412a",
   "Name":"Task task85de4103-8757-4c7d-942f-55eaf7d6412a",
   "Message":"The task with id task85de4103-8757-4c7d-942f-55eaf7d6412a has started.",
   "MessageId":"TaskEvent.1.0.3.TaskStarted",
   "MessageArgs":[
      "task85de4103-8757-4c7d-942f-55eaf7d6412a"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```

>**Sample response body** (subtask)

```
{
"@odata.type": "#Task.v1_6_0.Task",
"@odata.id": "/redfish/v1/TaskService/Tasks/task2da1ea5d-5604-49e2-9795-694909f99e15",
"@odata.context": "/redfish/v1/$metadata#Task.Task",
"Id": "task2da1ea5d-5604-49e2-9795-694909f99e15",
"Name": "Task task2da1ea5d-5604-49e2-9795-694909f99e15",
"TaskState": "Exception",
"StartTime": "2022-08-18T10:41:44.629353222Z",
"EndTime": "2022-08-18T10:41:46.220216368Z",
"TaskStatus": "Critical",
"SubTasks":{
      "@odata.id": "/redfish/v1/TaskService/Tasks/task2da1ea5d-5604-49e2-9795-694909f99e15/SubTasks"
},
"TaskMonitor": "/taskmon/task2da1ea5d-5604-49e2-9795-694909f99e15",
"PercentComplete": 100,
"Payload":{
"HttpHeaders":[],
"HttpOperation": "POST",
"JsonBody": "{\"BatchSize\":2,\"DelayBetweenBatchesInSeconds\":1,\"Password\":\"null\",\"ResetType\":\"ForceRestart\",\"TargetURIs\":[\"/redfish/v1/Systems/35270772-a2a0-4b1c-811f-600528eac2d9.1\"]}",
"TargetUri": "/redfish/v1/AggregationService/Actions/AggregationService.Reset/"
},
"Oem":{}
}
```

>**Sample response body** \(HTTP 200 status\)

```
{ 
   "error":{ 
      "code":"Base.1.13.0.Success",
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
|<strong>Returns</strong> |- `Location` URI of the task monitor associated with this operation in the response header. See `Location` URI in *Sample response header (HTTP 202 status)*.<br><br>-  Link to the task and the task id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See the task URI and the task Id  in *Sample response body (HTTP 202 status)*.<br>IMPORTANT:<br>Make a note of the task id. If the task completes with an error, it is required to know which subtask has failed. To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`.</blockquote><br>- On successful completion of this operation, a message in the response body, saying that the operation is completed successfully. See *Sample response body (HTTP 200 status)*.<br>|
|<strong>Response code</strong> |`202 Accepted`. On successful completion, `200 OK` <br> |
|<strong>Authentication</strong> |Yes|

**Usage information**

To know the progress of this action, perform HTTP `GET` on the *[task monitor](#viewing-a-task-monitor)* returned in the response header (until the task is complete).

To get the list of subtask URIs, perform HTTP `GET` on the task URI returned in the JSON response body. See *Sample response body (HTTP 202 status)*. The JSON response body of each subtask contains a link to the task monitor associated with it. To know the progress of `SetDefaultBootOrder` action (subtask) on a specific server, perform HTTP `GET` on the task monitor associated with the respective subtask. See the link to the task monitor  in *Sample response body (subtask)*.

You can perform `setDefaultBootOrder` action on a group of servers by specifying multiple server URIs in the request.


>**curl command**

```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{
   "Systems":[
      {
         "@odata.id":"/redfish/v1/Systems/{ComputerSystemId}"
      },
      {
         "@odata.id":"/redfish/v1/Systems/{ComputerSystemid2}"
      }
   ]
}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder'
```

>**Sample request body**

```
{
   "Systems":[
      {
         "@odata.id":"/redfish/v1/Systems/97d08f36-17f5-5918-8082-f5156618f58d.1"
      },
      {
         "@odata.id":"/redfish/v1/Systems/76632110-1c75-5a86-9cc2-471325983653.1"
      }
   ]
}
```

> **Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Systems|Array (required)<br> |Target servers for `SetDefaultBootOrder`.|

>**Sample response header** (HTTP 202 status)

```
Location:/taskmon/task85de4003-8057-4c7d-942f-55eaf7d6412a
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes
```

>**Sample response body** (HTTP 202 status)

```
{
   "@odata.type":"#Task.v1_6_0.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4003-8057-4c7d-942f-55eaf7d6412a",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task85de4003-8057-4c7d-942f-55eaf7d6412a",
   "Name":"Task task85de4003-8057-4c7d-942f-55eaf7d6412a",
   "Message":"The task with id task80de4003-8757-4c7d-942f-55eaf7d6412a has started.",
   "MessageId":"TaskEvent.1.0.3.TaskStarted",
   "MessageArgs":[
      "task80de4003-8757-4c7d-942f-55eaf7d6412a"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```

>**Sample response body** (subtask)

```
{
    "@odata.type": "#Task.v1_6_0.Task",
    "@odata.id": "/redfish/v1/TaskService/Tasks/taskabd8c681-a484-44fe-8ec4-e4929a44d1f2",
    "@odata.context": "/redfish/v1/$metadata#Task.Task",
    "Id": "taskabd8c681-a484-44fe-8ec4-e4929a44d1f2",
    "Name": "Task taskabd8c681-a484-44fe-8ec4-e4929a44d1f2",
    "TaskState": "Exception",
    "StartTime": "2022-02-25T14:50:02.00265165Z",
    "EndTime": "2022-02-25T14:50:02.987585968Z",
    "TaskStatus": "Critical",
    "SubTasks":{
      "@odata.id": "/redfish/v1/TaskService/Tasks/task2da1ea5d-5604-49e2-9795-694909f99e15/SubTasks"
},

    "TaskMonitor": "/taskmon/taskabd8c681-a484-44fe-8ec4-e4929a44d1f2",
    "PercentComplete": 100,
    "Payload": {
        "HttpHeaders": [
        ],
        "HttpOperation": "POST",
        "JsonBody": "{\"Systems\":[{\"@odata.id\":\"/redfish/v1/Systems/a84005b0-928a-4a8d-9994-335bbe15a915.1\"},{\"@odata.id\":\"/redfish/v1/Systems/921b734b-b35e-4387-a71a-14a1fd0bdc69.1\"}]}",
        "TargetUri": "/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder"
    },
    "Oem": {
    }
}
```

>**Sample response body** \(HTTP 200 status\)

```
{ 
   "error":{ 
      "code":"Base.1.13.0.Success",
      "message":"Request completed successfully"
   }
}
```


## Deleting a resource from the inventory

| |                                                              |
|--------|--------|
|<strong>Method</strong> | `DELETE` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/AggregationSources/{AggregationSourceId}` |
|<strong>Description</strong> |This operation removes a specific aggregation source (plugin, BMC, or any manager) from the inventory. Deleting an aggregation source also deletes all event subscriptions associated with the BMC. This operation is performed in the background as a Redfish task.<br> |
|<strong>Returns</strong> |- `Location` URI of the task monitor associated with this operation in the response header. See `Location` URI in *Sample response header (HTTP 202 status)*.<br>-   Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See *Sample response body (HTTP 202 status)*.<br>|
|<strong>Response Code</strong> |`202 Accepted` On successful completion, `204 No Content` <br> |
|<strong>Authentication</strong> |Yes|

**Usage information**

To know the progress of this action, perform `GET` on the *[task monitor](#viewing-a-task-monitor)* returned in the response header (until the task is complete).


>**curl command**

```
curl -i DELETE \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/AggregationSources/{AggregationSourceId}'
```

>**Sample response header** (HTTP 202 status)

```
Location:/taskmon/task85de4003-8757-2c7d-942f-55eaf7d6412a
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes
```

>**Sample response body** (HTTP 202 status)

```
{
   "@odata.type":"#Task.v1_6_0.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4003-8757-2c7d-942f-55eaf7d6412a",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task85de4003-8757-2c7d-942f-55eaf7d6412a",
   "Name":"Task task85de4003-8757-2c7d-942f-55eaf7d6412a",
   "Message":"The task with id task85de4003-8757-2c7d-942f-55eaf7d6412a has started.",
   "MessageId":"TaskEvent.1.0.3.TaskStarted",
   "MessageArgs":[
      "task85de4003-8757-2c7d-942f-55eaf7d6412a"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```


## Aggregates

An aggregate is a user-defined collection of resources.

The aggregate schema provides a mechanism to formally group the southbound resources of your choice. The advantage of creating aggregates is that they are more persistent than the random groupings. The aggregates are available and accessible in the environment of Resource Aggregator for ODIM until you delete them.

The resource aggregator allows you to perform the following tasks:

-   Create an aggregate

-   Populate an aggregate with the resources

-   Perform actions on all the resources of an aggregate at once

-   Delete an aggregate

## Creating an aggregate

|||
|---------|-----------|
|<strong>Method</strong> | `POST` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/Aggregates` |
|<strong>Description</strong> |This operation creates an empty aggregate or an aggregate populated with resources.|
|<strong>Returns</strong> | The `Location` URI of the created aggregate having the aggregate Id. See the `Location` URI in "Sample response header".<br>-   Link to the new aggregate, its Id, and a message saying that the resource has been created successfully in the JSON response body. |
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
            {"@odata.id": "/redfish/v1/Systems/{ComputerSystemId}"      
   ]   
}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Aggregates'
```

>**Sample request body**

```
{
      "Elements":[

             {
               "@odata.id": "/redfish/v1/Systems/{ComputerSystemId}"
             }      
   ]   
}
```

> **Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Elements|Array of objects (required)<br> |An empty array or an array of links to the object resources that this aggregate contains. To get the links to the system resources that are available in the resource inventory, perform an HTTP `GET` on `/redfish/v1/Systems/`. |

>**Sample response header**

```
Link:</redfish/v1/AggregationService/Aggregates/c14d91b5-3333-48bb-a7b7-75f74a137d48/>; rel=self
Location:/redfish/v1/AggregationService/Aggregates/c14d91b5-3333-48bb-a7b7-75f74a137d48
Date:Fri,21 August 2020 14:08:55 GMT+5m 11s
```

>**Sample response body**

```
{
      "@odata.type":"#Aggregate.v1_0_1.Aggregate",
      "@odata.id":"/redfish/v1/AggregationService/Aggregates/c14d91b5-3333-48bb-a7b7-75f74a137d48",
      "@odata.context":"/redfish/v1/$metadata#Aggregate.Aggregate",
      "Id":"c14d91b5-3333-48bb-a7b7-75f74a137d48",
      "Name":"Aggregate",
      "Elements":[
            "@odata.id": "/redfish/v1/Systems/{ComputerSystemId}"      
   ]   
}
```


## Viewing a list of aggregates

|||
|----------|-----------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/AggregationService/Aggregates` |
|<strong>Description</strong> |This operation lists all aggregates available in Resource Aggregator for ODIM.|
|<strong>Returns</strong> |Links of all the available aggregates|
|<strong>Response Code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

>**curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Aggregates'
```
> **Sample response body**

```
{
    "@odata.type": "#AggregateCollection.AggregateCollection",
    "@odata.id": "/redfish/v1/AggregationService/Aggregates",
    "@odata.context": "/redfish/v1/$metadata#AggregateCollection.AggregateCollection",
    "Description": "Aggregate collection view",
    "Name": "Aggregate",
    "Members@odata.count": 1,
    "Members": [
        {
            "@odata.id": "/redfish/v1/AggregationService/Aggregates/30e04950-df9c-4e4d-8ff1-1f5ffae9c7cb"
        }
    ]
}
```

## Viewing information about a single aggregate

|                                 |                                                              |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>Method</strong>         | `GET`                                                        |
| <strong>URI</strong>            | `/redfish/v1/AggregationService/Aggregates/{AggregateId}`    |
| <strong>Description</strong>    | This operation retrieves information about a specific aggregate. |
| <strong>Returns</strong>        | JSON schema representing this aggregate                      |
| <strong>Response Code</strong>  | On success, `200 Ok`                                         |
| <strong>Authentication</strong> | Yes                                                          |

> **curl command**

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Aggregates/{AggregateId}'

```
> **Sample response body**

```
{
    "@odata.type": "#Aggregate.v1_0_1.Aggregate",
    "@odata.id": "/redfish/v1/AggregationService/Aggregates/30e04950-df9c-4e4d-8ff1-1f5ffae9c7cb",
    "@odata.context": "/redfish/v1/$metadata#Aggregate.Aggregate",
    "Id": "30e04950-df9c-4e4d-8ff1-1f5ffae9c7cb",
    "Name": "Aggregate",
    "ElementsCount": 1,
    "Elements": [
        {
            "@odata.id": "/redfish/v1/Systems/766b0eca-ad76-46d5-afb4-b5d6b3650c0e.1"
        }
    ],
    "Actions": {
        "#Aggregate.Reset": {
            "target": "/redfish/v1/AggregationService/Aggregates/30e04950-df9c-4e4d-8ff1-1f5ffae9c7cb/Actions/Aggregate.Reset"
        },
        "#Aggregate.SetDefaultBootOrder": {
            "target": "/redfish/v1/AggregationService/Aggregates/30e04950-df9c-4e4d-8ff1-1f5ffae9c7cb/Actions/Aggregate.SetDefaultBootOrder"
        },
        "#Aggregate.AddElements": {
            "target": "/redfish/v1/AggregationService/Aggregates/30e04950-df9c-4e4d-8ff1-1f5ffae9c7cb/Actions/Aggregate.AddElements"
        },
        "#Aggregate.RemoveElements": {
            "target": "/redfish/v1/AggregationService/Aggregates/30e04950-df9c-4e4d-8ff1-1f5ffae9c7cb/Actions/Aggregate.RemoveElements"
        }
    }
}
```
## Deleting an aggregate

|                                 |                                                           |
| ------------------------------- | --------------------------------------------------------- |
| <strong>Method</strong>         | `DELETE`                                                  |
| <strong>URI</strong>            | `/redfish/v1/AggregationService/Aggregates/{AggregateId}` |
| <strong>Description</strong>    | This operation deletes a specific aggregate.              |
| <strong>Response Code</strong>  | On success, `204 No Content`                              |
| <strong>Authentication</strong> | Yes                                                       |

> **curl command**

```
curl -i DELETE \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Aggregates/{AggregateId}'
```
## Adding elements to an aggregate

|                                 |                                                              |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>Method</strong>         | `POST`                                                       |
| <strong>URI</strong>            | `/redfish/v1/AggregationService/Aggregates/{AggregateId}/Actions/Aggregate.AddElements` |
| <strong>Description</strong>    | This action adds one or more resources to a specific aggregate. |
| <strong>Returns</strong>        | JSON schema for this aggregate having links to the added resources |
| <strong>Response Code</strong>  | On success, `200 Ok`                                         |
| <strong>Authentication</strong> | Yes                                                          |

> **curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
      "Elements":[
            {
              "@odata.id": "/redfish/v1/Systems/{ComputerSystemId}" 
            }
   ]   
}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Aggregates/{AggregateId}/Actions/Aggregate.AddElements'
```

> **Sample request body**

```
{
      "Elements":[
            {
              "@odata.id": "/redfish/v1/Systems/{ComputerSystemId}" 
            }
   ]   
}
```
> **Request parameters**

| Parameter | Type                            | Description                                                  |
| --------- | ------------------------------- | ------------------------------------------------------------ |
| Elements  | Array of objects (required)<br> | An array of object links to the Computer system resources the aggregate contains |

> **Sample response body**

```
{
      "@odata.type":"#Aggregate.v1_0_1.Aggregate",
      "@odata.id":"/redfish/v1/AggregationService/Aggregates/c14d91b5-3333-48bb-a7b7-75f74a137d48",
      "@odata.context":"/redfish/v1/$metadata#Aggregate.Aggregate",
      "Id":"c14d91b5-3333-48bb-a7b7-75f74a137d48",
      "Name":"Aggregate",
      "Message":"Successfully Completed Request",
      "MessageId":"Base.1.13.0.Created",
      "Severity":"OK",
      "Elements":[
            {
              "@odata.id": "/redfish/v1/Systems/{ComputerSystemId}" 
            }     
   ]   
}
```
## Resetting an aggregate of computer systems

|                                 |                                                              |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>Method</strong>         | `POST`                                                       |
| <strong>URI</strong>            | `/redfish/v1/AggregationService/Aggregates/{AggregateId}/Actions/Aggregate.Reset` |
| <strong>Description</strong>    | This action shuts down, powers up, and restarts servers in a specific aggregate. This operation is performed in the background as a Redfish task and is further divided into subtasks to reset each server individually.<br> |
| <strong>Returns</strong>        | - `Location` URI of the task monitor associated with this operation (task) in the response header. See `Location` URI in *Sample response header (HTTP 202 status)*.<br>- Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See the task URI and the task Id  in *Sample response body (HTTP 202 status)*.<br>**IMPORTANT**: Make a note of the task id. If the task completes with an error, it is required to know which subtask has failed. To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`.<br>- Upon the completion of the reset operation, you receive a success message in the response body. See *Sample response body (HTTP 200 status)*. |
| <strong>Response Code</strong>  | `202 Accepted` On successful completion, `200 OK` <br>       |
| <strong>Authentication</strong> | Yes                                                          |

**Usage information**

To know the progress of this action, perform HTTP `GET` on the *[task monitor](#viewing-a-task-monitor)* returned in the response header (until the task is complete).

To get the list of subtask URIs, perform HTTP `GET` on the task URI returned in the JSON response body. See *Sample response body (HTTP 202 status)*. The JSON response body of each subtask contains a link to the task monitor associated with it. To know the progress of the reset operation (subtask) on a specific server, perform HTTP `GET` on the task monitor associated with the respective subtask. See the link to the task monitor  in *Sample response body (subtask)*.

> **curl command**

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
> **Sample request body**


```
{
   "BatchSize":2,
   "DelayBetweenBatchesInSeconds":2,
   "ResetType":"ForceRestart"
}

```
> **Request parameters**

| Parameter                    | Type                             | Description                                                  |
| ---------------------------- | -------------------------------- | ------------------------------------------------------------ |
| BatchSize                    | Integer (optional)<br>           | The number of elements to be reset at a time in each batch   |
| DelayBetweenBatchesInSeconds | Integer (seconds) (optional)<br> | The delay among the batches of elements being reset          |
| ResetType                    | String (optional)<br>            | For possible values, see *Reset type* table in *[Resetting servers](#resetting-servers)*. |

> **Sample response header** (HTTP 202 status)

```
Location:/taskmon/task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes
```

> **Sample response body** (HTTP 202 status)

```
{
   "@odata.type":"#Task.v1_6_0.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591",
   "Name":"Task task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591",
   "Message":"The task with id task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591 has started.",
   "MessageId":"TaskEvent.1.0.3.TaskStarted",
   "MessageArgs":[
      "task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```
> **Sample response body** (subtask)

```
{
    "@odata.type": "#Task.v1_6_0.Task",
    "@odata.id": "/redfish/v1/TaskService/Tasks/taskd862139f-c664-4cb2-b771-3e702bde40e3",
    "@odata.context": "/redfish/v1/$metadata#Task.Task",
    "Id": "taskd862139f-c664-4cb2-b771-3e702bde40e3",
    "Name": "Task taskd862139f-c664-4cb2-b771-3e702bde40e3",
    "TaskState": "Completed",
    "StartTime": "2022-02-25T13:07:05.938018291Z",
    "EndTime": "2022-02-25T13:07:08.108846323Z",
    "TaskStatus": "OK",
    "SubTasks":{
      "@odata.id": "/redfish/v1/TaskService/Tasks/task2da1ea5d-5604-49e2-9795-694909f99e15/SubTasks"
},
    "TaskMonitor": "/taskmon/taskd862139f-c664-4cb2-b771-3e702bde40e3",
    "PercentComplete": 100,
    "Payload": {
        "HttpHeaders": [
        ],
        "HttpOperation": "POST",
        "JsonBody": "{\"BatchSize\":2,\"DelayBetweenBatchesInSeconds\":2,\"Password\":\"null\",\"ResetType\":\"ForceRestart\"}",
        "TargetUri": "/redfish/v1/AggregationService/Aggregates/ca3f2462-15b5-4eb6-80c1-89f99ac36b12/Actions/Aggregate.Reset"
    },
    "Oem": {
    }
}
```
> **Sample response body** (HTTP 200 status)

```
{
   "error":{
      "code":"Base.1.13.0.Success",
      "message":"Request completed successfully"
   }
}
```
 ## Setting boot order of an aggregate to default settings

|                                 |                                                              |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>Method</strong>         | `POST`                                                       |
| <strong>URI</strong>            | `/redfish/v1/AggregationService/Aggregates/{AggregateId}/Actions/Aggregate.SetDefaultBootOrder` |
| <strong>Description</strong>    | This action changes the boot order of all the servers belonging to a specific aggregate to default settings. This operation is performed in the background as a Redfish task and is further divided into subtasks to change the boot order of each server individually.<br> |
| <strong>Returns</strong>        | - `Location` URI of the created aggregate having the aggregate id. See the `Location` URI in *Sample response header*.<br>-   Link to the new aggregate, its id, and a success message in the JSON response body.<br>`Location` URI of the task monitor associated with this operation in the response header. See `Location` URI in *Sample response header (HTTP 202 status)*.<br>-   Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See the task URI and the task Id  in *Sample response header (HTTP 202 status)*.<br>**IMPORTANT**: Make a note of the task id. If the task completes with an error, it is required to know which subtask has failed. To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`.<br>Upon the completion of the operation, you receive a success message in the response body. See *Sample response body (HTTP 200 status)*.<br> |
| <strong>Response Code</strong>  | `202 Accepted`. On successful completion, `200 OK` <br>      |
| <strong>Authentication</strong> | Yes                                                          |

**Usage information**

To know the progress of this action, perform HTTP `GET` on the *[task monitor](#viewing-a-task-monitor)* returned in the response header (until the task is complete).

To get the list of subtask URIs, perform HTTP `GET` on the task URI returned in the JSON response body. See *Sample response body (HTTP 202 status)*. The JSON response body of each subtask contains a link to the task monitor associated with it. To know the progress of `SetDefaultBootOrder` action (subtask) on a specific server, perform HTTP `GET` on the task monitor associated with the respective subtask. See the link to the task monitor  in *Sample response body (subtask)*".

> **curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Aggregates/{AggregateId}/Actions/Aggregate.SetDefaultBootOrder'
```
> **Sample response header** (HTTP 202 status)

```
Location:/taskmon/task85de4003-8057-4c7d-942f-55eaf7d6412a
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes
```
> **Sample response body** (HTTP 202 status)

```
{
   "@odata.type":"#Task.v1_6_0.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4003-8057-4c7d-942f-55eaf7d6412a",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task85de4003-8057-4c7d-942f-55eaf7d6412a",
   "Name":"Task task85de4003-8057-4c7d-942f-55eaf7d6412a",
   "Message":"The task with id task80de4003-8757-4c7d-942f-55eaf7d6412a has started.",
   "MessageId":"TaskEvent.1.0.3.TaskStarted",
   "MessageArgs":[
      "task80de4003-8757-4c7d-942f-55eaf7d6412a"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```
> **Sample response body** (subtask)

```
{
    "@odata.type": "#Task.v1_6_0.Task",
    "@odata.id": "/redfish/v1/TaskService/Tasks/task94f9af7a-fbe4-4846-94c9-9d5f7b949e40",
    "@odata.context": "/redfish/v1/$metadata#Task.Task",
    "Id": "task94f9af7a-fbe4-4846-94c9-9d5f7b949e40",
    "Name": "Task task94f9af7a-fbe4-4846-94c9-9d5f7b949e40",
    "TaskState": "Exception",
    "StartTime": "2022-02-25T13:27:29.518305955Z",
    "EndTime": "2022-02-25T13:27:30.007624346Z",
    "TaskStatus": "Critical",
    "SubTasks":{
      "@odata.id": "/redfish/v1/TaskService/Tasks/task2da1ea5d-5604-49e2-9795-694909f99e15/SubTasks"
},
    "TaskMonitor": "/taskmon/task94f9af7a-fbe4-4846-94c9-9d5f7b949e40",
    "PercentComplete": 100,
    "Payload": {
        "HttpHeaders": [
        ],
        "HttpOperation": "POST",
        "JsonBody": "{\"Password\":\"null\",\"SessionToken\":\"1e2ce744-8bcb-4e97-9eb1-2b419b1e7a2c\",\"URL\":\"/redfish/v1/AggregationService/Aggregates/ca3f2462-15b5-4eb6-80c1-89f99ac36b12/Actions/Aggregate.SetDefaultBootOrder\"}",
        "TargetUri": "/redfish/v1/AggregationService/Aggregates/ca3f2462-15b5-4eb6-80c1-89f99ac36b12/Actions/Aggregate.SetDefaultBootOrder"
    },
    "Oem": {
    }
}
```
> **Sample response body** (HTTP 200 status)

```
{ 
   "error":{ 
      "code":"Base.1.13.0.Success",
      "message":"Request completed successfully"
   }
}
```
## Removing elements from an aggregate

|                                 |                                                              |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>Method</strong>         | `POST`                                                       |
| <strong>URI</strong>            | `/redfish/v1/AggregationService/Aggregates/{AggregateId}/Actions/Aggregate.RemoveElements` |
| <strong>Description</strong>    | This action removes one or more resources from a specific aggregate. |
| <strong>Returns</strong>        | Updated JSON schema representing this aggregate.             |
| <strong>Response Code</strong>  | On success, `200 Ok`                                         |
| <strong>Authentication</strong> | Yes                                                          |

> **curl command**

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{
      "Elements":[
             {
               "@odata.id": "/redfish/v1/Systems/{ComputerSystemId}" 
             }
   ]   
}' \
 'https://{odim_host}:{port}/redfish/v1/AggregationService/Aggregates/{AggregateId}/Actions/Aggregate.RemoveElements'
```
> **Sample request body**

```
{
      "Elements":[
        {
          "@odata.id": "/redfish/v1/Systems/{ComputerSystemId}"
        }
   ] 
}
```
> **Request parameters**

| Parameter | Type                            | Description                                                  |
| --------- | ------------------------------- | ------------------------------------------------------------ |
| Elements  | Array of objects (required)<br> | An array of object links of the Computer system resources that you want to remove from this aggregate |

> **Sample response body**



```
{
   "@odata.type":"#Aggregate.v1_0_1.Aggregate",
   "@odata.id":"/redfish/v1/AggregationService/Aggregates/e02faf78-f919-4612-b031-bec7ae59910d",
   "@odata.context":"/redfish/v1/$metadata#Aggregate.Aggregate",
   "Id":"e02faf78-f919-4612-b031-bec7ae59910d",
   "Name":"Aggregate",
   "Message": "Successfully Completed Request",
   "MessageId": "Base.1.13.0.Success",
   "Severity":"OK",
   "Elements":[
   ]
}
```
#  Resource inventory

Resource Aggregator for ODIM allows you to view the inventory of compute and local storage resources through Redfish `Systems`, `Chassis`, and `Managers` endpoints. 
It also offers the capability to perform the following tasks:	

- Search inventory information based on one or more configuration parameters
- Manage resources added in the inventory


To discover crucial configuration information about a resource, including chassis, perform `GET` on these endpoints.

**Supported endpoints**

| API URI                                                      | Supported operations | Required privileges            |
| ------------------------------------------------------------ | -------------------- | ------------------------------ |
| /redfish/v1/Systems                                          | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}                       | `GET`, `PATCH`       | `Login`, `ConfigureComponents` |
| /redfish/v1/Systems/{ComputerSystemId}/Memory                | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/Memory/{memoryId}     | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/MemoryDomains         | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/NetworkInterfaces     | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces    | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces/{id} | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/Bios                  | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/SecureBoot            | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/PCIeDevices/{PCIeDeviceId} | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/Storage               | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id} | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/AllocatedVolumes | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/AllocatedVolumes/{allocatedvolumes_Id} | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/CapacitySources/{capacitysources_Id}/ProvidingDrives | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/CapacitySources/{capacitysources_Id}/ProvidingDrives/{providingdrives_id} | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId} | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Drives/{driveId} | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes | `GET`, `POST`        | `Login`, `ConfigureComponents` |
| /redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/Capabilities | `GET`                |                                |
| /redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/{volumeId} | `GET`, `DELETE`      | `Login`, `ConfigureComponents` |
| /redfish/v1/Systems/{ComputerSystemId}/Processors            | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/Processors/{id}       | `GET`                | `Login`                        |
| /redfish/v1/Systems?$filter={searchKeys}%20{conditionKeys}%20{value} | `GET`                | `Login`                        |
| /redfish/v1/Systems/{ComputerSystemId}/Bios/Settings<br>     | `GET`, `PATCH`       | `Login`, `ConfigureComponents` |
| /redfish/v1/Systems/{ComputerSystemId}/Actions/ComputerSystem.Reset | `POST`               | `ConfigureComponents`          |
| /redfish/v1/Systems/{ComputerSystemId}/Actions/ComputerSystem.SetDefaultBootOrder | `POST`               | `ConfigureComponents`          |

| API URI                                                      | Operation Applicable     | Required privileges            |
| ------------------------------------------------------------ | ------------------------ | ------------------------------ |
| /redfish/v1/Chassis                                          | `GET`, `POST`            | `Login`, `ConfigureComponents` |
| /redfish/v1/Chassis/{chassisId}                              | `GET`, `PATCH`, `DELETE` | `Login`, `ConfigureComponents` |
| /redfish/v1/Chassis/{chassisId}/Thermal                      | `GET`                    | `Login`                        |
| /redfish/v1/Chassis/{chassisId}/NetworkAdapters              | `GET`                    | `Login`                        |
| /redfish/v1/Chassis/{ChassisId}/NetworkAdapters/{networkadapterId} | `GET`                    | `Login`                        |

| API URI                                             | Supported operations | Required privileges |
| --------------------------------------------------- | -------------------- | ------------------- |
| /redfish/v1/Managers                                | `GET`                | `Login`             |
| /redfish/v1/Managers/{managerId}                    | `GET`                | `Login`             |
| /redfish/v1/Managers/{managerId}/EthernetInterfaces | `GET`                | `Login`             |
| /redfish/v1/Managers/{managerId}/HostInterfaces     | `GET`                | `Login`             |
| /redfish/v1/Managers/{managerId}/LogServices        | `GET`                | `Login`             |
| /redfish/v1/Managers/{managerId}/NetworkProtocol    | `GET`                | `Login`             |


##  Collection of computer systems

Each computer system has a `ComputerSystemId`, a unique identifier of a system specified by Resource Aggregator for ODIM. It is represented as `<UUID.n>` in Resource Aggregator for ODIM. `<UUID.n>` is the universally unique identifier o f a system. 
**Example**: *ba0a6871-7bc4-5f7a-903d-67f3c205b08c.1*.

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/Systems`                                        |
| **Description**    | This operation lists all systems available with Resource Aggregator for ODIM. |
| **Returns**        | A collection of links to computer system instances           |
| **Response code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Systems'
```

> **Sample response body** 


```
{ 
   "@odata.context":"/redfish/v1/$metadata#ComputerSystemCollection.ComputerSystemCollection",
   "@odata.id":"/redfish/v1/Systems/",
   "@odata.type":"#ComputerSystemCollection.ComputerSystemCollection",
   "Description":"Computer Systems view",
   "Name":"Computer Systems",
   "Members":[ 
      { 
         "@odata.id":"/redfish/v1/Systems/ba0a6871-7bc4-5f7a-903d-67f3c205b08c.1"
      },
      { 
         "@odata.id":"/redfish/v1/Systems/7ff3bd97-c41c-5de0-937d-85d390691b73.1"
      }
   ],
   "Members@odata.count":2
}
```
## Single computer system

|                    |                                                            |
| ------------------ | ---------------------------------------------------------- |
| **Method**         | `GET`                                                      |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}`                   |
| **Description**    | This endpoint fetches information about a specific system. |
| **Returns**        | JSON schema representing this computer system instance     |
| **Response code**  | `200 OK`                                                   |
| **Authentication** | Yes                                                        |

> **curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}'
```
> **Sample response body** 

```
{
    "@odata.context": "/redfish/v1/$metadata#ComputerSystem.ComputerSystem",
    "@odata.etag": "W/\"BB5DA93F\"",
    "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1",
    "@odata.type": "#ComputerSystem.v1_18_0.ComputerSystem",
    "Actions": {
        "#ComputerSystem.Reset": {
            "ResetType@Redfish.AllowableValues": [
                "On",
                "ForceOff",
                "GracefulShutdown",
                "ForceRestart",
                "Nmi",
                "PushPowerButton",
                "GracefulRestart"
            ],
            "target": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/Actions/ComputerSystem.Reset"
        }
    },
    "AssetTag": "",
    "Bios": {
        "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/Bios"
    },
    "BiosVersion": "A40 v1.46 (07/10/2019)",
    "Boot": {
        "BootOptions": {
            "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/BootOptions"
        },
        "BootOrder": [
            "Boot0024:Unknown.Unknown.200.1",
            "Boot0017:NIC.FlexLOM.1.1.IPv4",
            "Boot000A:Generic.USB.1.1",
            "Boot000C:HD.SD.1.2",
            "Boot0011:HD.EmbRAID.1.3",
            "Boot0012:HD.EmbRAID.1.4",
            "Boot0013:HD.EmbRAID.1.5",
            "Boot0016:NIC.FlexLOM.1.1.Httpv4",
            "Boot001A:NIC.LOM.1.1.Httpv4",
            "Boot001B:NIC.LOM.1.1.IPv4",
            "Boot0018:NIC.LOM.1.1.Httpv6",
            "Boot0019:NIC.LOM.1.1.IPv6",
            "Boot0014:NIC.FlexLOM.1.1.Httpv6",
            "Boot0015:NIC.FlexLOM.1.1.IPv6",
            "Boot0021:NIC.Slot.1.1.Httpv4",
            "Boot0022:NIC.Slot.1.1.IPv4",
            "Boot001D:NIC.Slot.2.1.Httpv4",
            "Boot001E:NIC.Slot.2.1.IPv4",
            "Boot001F:NIC.Slot.1.1.Httpv6",
            "Boot0020:NIC.Slot.1.1.IPv6",
            "Boot000B:NIC.Slot.2.1.Httpv6",
            "Boot001C:NIC.Slot.2.1.IPv6",
            "Boot0009:HD.EmbRAID.1.6",
            "Boot000E:HD.EmbRAID.1.7",
            "Boot000F:HD.EmbRAID.1.8",
            "Boot000D:HD.EmbRAID.1.2"
        ],
        "BootSourceOverrideEnabled": "Disabled",
        "BootSourceOverrideMode": "UEFI",
        "BootSourceOverrideTarget": "None",
        "BootSourceOverrideTarget@Redfish.AllowableValues": [
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
        "UefiTargetBootSourceOverride": "None",
        "UefiTargetBootSourceOverride@Redfish.AllowableValues": [
            "HD(1,GPT,D8898303-6CD4-43FA-BDA0-66F8967EEA78,0x800,0x64000)/\\EFI\\red\\grubx64.efi",
            "PciRoot(0x0)/Pci(0x1,0x1)/Pci(0x0,0x0)/MAC(48DF377EF730,0x1)/IPv4(0.0.0.0)",
            "UsbClass(0xFFFF,0xFFFF,0xFF,0xFF,0xFF)",
            "PciRoot(0x0)/Pci(0x7,0x1)/Pci(0x0,0x3)/USB(0x3,0x0)/USB(0x0,0x0)",
            "PciRoot(0x1)/Pci(0x1,0x1)/Pci(0x0,0x0)/Scsi(0x4,0x4000)",
            "PciRoot(0x1)/Pci(0x1,0x1)/Pci(0x0,0x0)/Scsi(0x5,0x4000)",
            "PciRoot(0x1)/Pci(0x1,0x1)/Pci(0x0,0x0)/Scsi(0x6,0x4000)",
            "PciRoot(0x0)/Pci(0x1,0x1)/Pci(0x0,0x0)/MAC(48DF377EF730,0x1)/IPv4(0.0.0.0)/Uri()",
            "PciRoot(0x0)/Pci(0x1,0x2)/Pci(0x0,0x0)/MAC(08F1EA8EE70C,0x1)/IPv4(0.0.0.0)/Uri()",
            "PciRoot(0x0)/Pci(0x1,0x2)/Pci(0x0,0x0)/MAC(08F1EA8EE70C,0x1)/IPv4(0.0.0.0)",
            "PciRoot(0x0)/Pci(0x1,0x2)/Pci(0x0,0x0)/MAC(08F1EA8EE70C,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)/Uri()",
            "PciRoot(0x0)/Pci(0x1,0x2)/Pci(0x0,0x0)/MAC(08F1EA8EE70C,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)",
            "PciRoot(0x0)/Pci(0x1,0x1)/Pci(0x0,0x0)/MAC(48DF377EF730,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)/Uri()",
            "PciRoot(0x0)/Pci(0x1,0x1)/Pci(0x0,0x0)/MAC(48DF377EF730,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)",
            "PciRoot(0x2)/Pci(0x3,0x1)/Pci(0x0,0x0)/MAC(040973D10340,0x1)/IPv4(0.0.0.0)/Uri()",
            "PciRoot(0x2)/Pci(0x3,0x1)/Pci(0x0,0x0)/MAC(040973D10340,0x1)/IPv4(0.0.0.0)",
            "PciRoot(0x3)/Pci(0x3,0x1)/Pci(0x0,0x0)/MAC(48DF374763E4,0x1)/IPv4(0.0.0.0)/Uri()",
            "PciRoot(0x3)/Pci(0x3,0x1)/Pci(0x0,0x0)/MAC(48DF374763E4,0x1)/IPv4(0.0.0.0)",
            "PciRoot(0x2)/Pci(0x3,0x1)/Pci(0x0,0x0)/MAC(040973D10340,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)/Uri()",
            "PciRoot(0x2)/Pci(0x3,0x1)/Pci(0x0,0x0)/MAC(040973D10340,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)",
            "PciRoot(0x3)/Pci(0x3,0x1)/Pci(0x0,0x0)/MAC(48DF374763E4,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)/Uri()",
            "PciRoot(0x3)/Pci(0x3,0x1)/Pci(0x0,0x0)/MAC(48DF374763E4,0x1)/IPv6(0000:0000:0000:0000:0000:0000:0000:0000)",
            "PciRoot(0x1)/Pci(0x1,0x1)/Pci(0x0,0x0)/Scsi(0x3,0x4000)",
            "PciRoot(0x1)/Pci(0x1,0x1)/Pci(0x0,0x0)/Scsi(0x1,0x4000)",
            "PciRoot(0x1)/Pci(0x1,0x1)/Pci(0x0,0x0)/Scsi(0x0,0x4000)",
            "PciRoot(0x1)/Pci(0x1,0x1)/Pci(0x0,0x0)/Scsi(0x0,0x0)"
        ]
    },
    "EthernetInterfaces": {
        "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/EthernetInterfaces"
    },
    "Id": "add8f39d-aea7-4eea-aa24-fc1764c33040.1",
    "IndicatorLED": "Off",
    "Links": {
        "Chassis": [
            {
                "@odata.id": "/redfish/v1/Chassis/add8f39d-aea7-4eea-aa24-fc1764c33040.1"
            }
        ],
        "ManagedBy": [
            {
                "@odata.id": "/redfish/v1/Managers/add8f39d-aea7-4eea-aa24-fc1764c33040.1"
            }
        ]
    },
    "LogServices": {
        "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/LogServices"
    },
    "Manufacturer": "HPE",
    "Memory": {
        "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/Memory"
    },
    "MemoryDomains": {
        "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/MemoryDomains"
    },
    "MemorySummary": {
        "Status": {
            "HealthRollup": "OK"
        },
        "TotalSystemMemoryGiB": 512,
        "TotalSystemPersistentMemoryGiB": 0
    },
    "Model": "ProLiant DL385 Gen10",
    "Name": "Computer System",
    "NetworkInterfaces": {
        "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/NetworkInterfaces"
    },
    "Oem": {
        "Hpe": {
            "@odata.context": "/redfish/v1/$metadata#HpeComputerSystemExt.HpeComputerSystemExt",
            "@odata.type": "#HpeComputerSystemExt.v2_9_0.HpeComputerSystemExt",
            "Actions": {
                "#HpeComputerSystemExt.PowerButton": {
                    "PushType@Redfish.AllowableValues": [
                        "Press",
                        "PressAndHold"
                    ],
                    "target": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/Actions/Oem/Hpe/HpeComputerSystemExt.PowerButton"
                },
                "#HpeComputerSystemExt.SecureSystemErase": {
                    "target": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/Actions/Oem/Hpe/HpeComputerSystemExt.SecureSystemErase"
                },
                "#HpeComputerSystemExt.SystemReset": {
                    "ResetType@Redfish.AllowableValues": [
                        "ColdBoot",
                        "AuxCycle"
                    ],
                    "target": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/Actions/Oem/Hpe/HpeComputerSystemExt.SystemReset"
                }
            },
            "AggregateHealthStatus": {
                "AgentlessManagementService": "Unavailable",
                "BiosOrHardwareHealth": {
                    "Status": {
                        "Health": "OK"
                    }
                },
                "FanRedundancy": "Redundant",
                "Fans": {
                    "Status": {
                        "Health": "OK"
                    }
                },
                "Memory": {
                    "Status": {
                        "Health": "OK"
                    }
                },
                "Network": {
                    "Status": {
                        "Health": "OK"
                    }
                },
                "PowerSupplies": {
                    "PowerSuppliesMismatch": false,
                    "Status": {
                        "Health": "OK"
                    }
                },
                "PowerSupplyRedundancy": "Redundant",
                "Processors": {
                    "Status": {
                        "Health": "OK"
                    }
                },
                "SmartStorageBattery": {
                    "Status": {
                        "Health": "OK"
                    }
                },
                "Storage": {
                    "Status": {
                        "Health": "OK"
                    }
                },
                "Temperatures": {
                    "Status": {
                        "Health": "OK"
                    }
                }
            },
            "Bios": {
                "Backup": {
                    "Date": "06/24/2019",
                    "Family": "A40",
                    "VersionString": "A40 v1.44 (06/24/2019)"
                },
                "Current": {
                    "Date": "07/10/2019",
                    "Family": "A40",
                    "VersionString": "A40 v1.46 (07/10/2019)"
                },
                "UefiClass": 2
            },
            "CriticalTempRemainOff": false,
            "CurrentPowerOnTimeSeconds": null,
            "DeviceDiscoveryComplete": {
                "AMSDeviceDiscovery": "NoAMS",
                "DeviceDiscovery": "vMainDeviceDiscoveryComplete",
                "SmartArrayDiscovery": "Complete"
            },
            "ElapsedEraseTimeInMinutes": 0,
            "EndOfPostDelaySeconds": null,
            "EstimatedEraseTimeInMinutes": 0,
            "IntelligentProvisioningAlwaysOn": true,
            "IntelligentProvisioningIndex": 8,
            "IntelligentProvisioningLocation": "System Board",
            "IntelligentProvisioningVersion": "3.30.213",
            "IsColdBooting": false,
            "Links": {
                "EthernetInterfaces": {
                    "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/EthernetInterfaces"
                },
                "NetworkAdapters": {
                    "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/BaseNetworkAdapters"
                },
                "PCISlots": {
                    "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/PCISlots"
                },
                "PCIeDevices": {
                    "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/PCIeDevices"
                },
                "SecureEraseReportService": {
                    "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/SecureEraseReportService"
                },
                "SmartStorage": {
                    "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/SmartStorage"
                },
                "USBDevices": {
                    "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/USBDevices"
                },
                "USBPorts": {
                    "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/USBPorts"
                },
                "WorkloadPerformanceAdvisor": {
                    "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/WorkloadPerformanceAdvisor"
                }
            },
            "PCAPartNumber": "866342-001",
            "PCASerialNumber": "PWCDH%%LMBT086",
            "PostDiscoveryCompleteTimeStamp": null,
            "PostDiscoveryMode": null,
            "PostMode": null,
            "PostState": "InPostDiscoveryComplete",
            "PowerAllocationLimit": 1600,
            "PowerAutoOn": "Restore",
            "PowerOnDelay": "Minimum",
            "PowerOnMinutes": 270768,
            "PowerRegulatorMode": "Dynamic",
            "PowerRegulatorModesSupported": [
            ],
            "SMBIOS": {
                "extref": "/smBios"
            },
            "ServerFQDN": "",
            "SmartStorageConfig": [
                {
                    "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/smartstorageconfig"
                }
            ],
            "SystemROMAndiLOEraseComponentStatus": {
                "BIOSSettingsEraseStatus": "Idle",
                "iLOSettingsEraseStatus": "Idle"
            },
            "SystemROMAndiLOEraseStatus": "Idle",
            "UserDataEraseComponentStatus": {
            },
            "UserDataEraseStatus": "Idle",
            "VirtualProfile": "Inactive"
        }
    },
    "PCIeDevices": [
        {
            "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/PCIeDevices/1"
        },
        {
            "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/PCIeDevices/2"
        },
        {
            "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/PCIeDevices/3"
        }
    ],
    "PCIeDevices@odata.count": 3,
    "PowerState": "On",
    "ProcessorSummary": {
        "Count": 2,
        "Model": "AMD EPYC 7601 32-Core Processor                ",
        "Status": {
            "HealthRollup": "OK"
        }
    },
    "Processors": {
        "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/Processors"
    },
    "SKU": "878612-B21",
    "SecureBoot": {
        "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/SecureBoot"
    },
    "SerialNumber": "2M29120289",
    "Status": {
        "Health": "OK",
        "HealthRollup": "OK",
        "State": "Enabled"
    },
    "Storage": {
        "@odata.id": "/redfish/v1/Systems/add8f39d-aea7-4eea-aa24-fc1764c33040.1/Storage"
    },
    "SystemType": "Physical",
    "TrustedModules": [
        {
            "Oem": {
                "Hpe": {
                    "@odata.context": "/redfish/v1/$metadata#HpeTrustedModuleExt.HpeTrustedModuleExt",
                    "@odata.type": "#HpeTrustedModuleExt.v2_0_0.HpeTrustedModuleExt"
                }
            },
            "Status": {
                "State": "Absent"
            }
        }
    ],
    "UUID": "36383738-3231-4D32-3239-313230323839",
    "VirtualMedia": {
        "@odata.id": "/redfish/v1/Managers/add8f39d-aea7-4eea-aa24-fc1764c33040.1/VirtualMedia"
    }
}
```
##  Memory collection

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}/Memory`              |
| **Description**    | This operation lists all memory devices of a specific server. |
| **Returns**        | List of memory resource endpoints                            |
| **Response code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Memory'

```
> **Sample response body** 

```
{
    "@odata.context": "/redfish/v1/$metadata#MemoryCollection.MemoryCollection",
    "@odata.etag": "W/\"09417F5F\"",
    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Memory",
    "@odata.type": "#MemoryCollection.MemoryCollection",
    "Description": "Memory DIMM Collection",
    "Members": [
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Memory/proc1dimm1"
        },
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Memory/proc1dimm2"
        },
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Memory/proc1dimm3"
        },
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Memory/proc1dimm4"
        }
    ],
    "Members@odata.count": 4,
    "Name": "Memory DIMM Collection",
    "Oem": {
        "Hpe": {
            "@odata.context": "/redfish/v1/$metadata#HpeAdvancedMemoryProtection.HpeAdvancedMemoryProtection",
            "@odata.type": "#HpeAdvancedMemoryProtection.v2_0_0.HpeAdvancedMemoryProtection",
            "AmpModeActive": "A3DC",
            "AmpModeStatus": "DegradedA3DC",
            "AmpModeSupported": [
                "AdvancedECC",
                "OnlineSpareRank",
                "IntrasocketMirroring",
                "A3DC"
            ],
            "MemoryList": [
                {
                    "BoardCpuNumber": 1,
                    "BoardNumberOfSockets": 12,
                    "BoardOperationalFrequency": 2666,
                    "BoardOperationalVoltage": 1200,
                    "BoardTotalMemorySize": 196608
                },
                {
                    "BoardCpuNumber": 2,
                    "BoardNumberOfSockets": 12,
                    "BoardOperationalFrequency": 2666,
                    "BoardOperationalVoltage": 1200,
                    "BoardTotalMemorySize": 196608
                }
            ]
        }
    }
}
```
## Single memory

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | GET                                                          |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}/Memory/{memoryId}`   |
| **Description**    | This endpoint retrieves configuration information of specific memory. |
| **Returns**        | JSON schema representing this memory resource.               |
| **Response code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Memory/{memoryId}'
```
> **Sample response body** 

```
{
    "@odata.context": "/redfish/v1/$metadata#Memory.Memory",
    "@odata.etag": "W/\"E6EC3A2C\"",
    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Memory/proc1dimm1",
    "@odata.type": "#Memory.v1_7_1.Memory",
    "BaseModuleType": "RDIMM",
    "BusWidthBits": 72,
    "CacheSizeMiB": 0,
    "CapacityMiB": 32768,
    "DataWidthBits": 64,
    "DeviceLocator": "PROC 1 DIMM 1",
    "ErrorCorrection": "MultiBitECC",
    "Id": "proc1dimm1",
    "LogicalSizeMiB": 0,
    "Manufacturer": "HPE",
    "MemoryDeviceType": "DDR4",
    "MemoryLocation": {
        "Channel": 6,
        "MemoryController": 2,
        "Slot": 1,
        "Socket": 1
    },
    "MemoryMedia": [
        "DRAM"
    ],
    "MemoryType": "DRAM",
    "Name": "proc1dimm1",
    "NonVolatileSizeMiB": 0,
    "Oem": {
        "Hpe": {
            "@odata.context": "/redfish/v1/$metadata#HpeMemoryExt.HpeMemoryExt",
            "@odata.type": "#HpeMemoryExt.v2_5_0.HpeMemoryExt",
            "Attributes": [
                "HpeSmartMemory"
            ],
            "BaseModuleType": "RDIMM",
            "DIMMManufacturingDate": "1828",
            "DIMMStatus": "GoodInUse",
            "MaxOperatingSpeedMTs": 2666,
            "MinimumVoltageVoltsX10": 12,
            "VendorName": "Samsung"
        }
    },
    "OperatingMemoryModes": [
        "Volatile"
    ],
    "OperatingSpeedMhz": 2666,
    "PartNumber": "M393A4K40CB2-CTD   ",
    "PersistentRegionSizeLimitMiB": 0,
    "RankCount": 2,
    "SecurityCapabilities": {
    },
    "SerialNumber": "39F51030",
    "Status": {
        "Health": "OK",
        "State": "Enabled"
    },
    "VendorID": "52736",
    "VolatileRegionSizeLimitMiB": 32768,
    "VolatileSizeMiB": 32768
}
```
##  Memory domains

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}/MemoryDomains`       |
| **Description**    | This endpoint lists memory domains of a specific system.<br>Memory Domains indicate to the client which Memory (DIMMs) can be grouped in Memory Chunks to form interleave sets, or otherwise grouped.<br> |
| **Returns**        | List of memory domain endpoints                              |
| **Response code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/MemoryDomains'
```
> **Sample response body** 

```
{
    "@odata.context": "/redfish/v1/$metadata#MemoryDomainCollection.MemoryDomainCollection",
    "@odata.etag": "W/\"75983E8D\"",
    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/MemoryDomains",
    "@odata.type": "#MemoryDomainCollection.MemoryDomainCollection",
    "Description": "Memory Domains Collection",
    "Members": [
    ],
    "Members@odata.count": 0,
    "Name": "Memory Domains Collection"
}
```
##  BIOS

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}/Bios`                |
| **Description**    | Use this endpoint to discover system-specific information about a BIOS resource and actions for changing to BIOS settings.<br>**NOTE:** Changes to the BIOS typically require a system reset before they take effect. |
| **Returns**        | <ul><li>Actions for changing password and resetting BIOS</li><li>BIOS attributes</li></ul> |
| **Response code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Bios'
```
> **Sample response body** 

```
{
    "@Redfish.Settings": {
        "@odata.type": "#Settings.v1_0_0.Settings",
        "ETag": "5D44558E",
        "Messages": [
            {
                "MessageId": "Base.1.0.Success"
            }
        ],
        "SettingsObject": {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Bios/Settings/"
        },
        "Time": "2022-03-09T11:56:28+00:00"
    },
    "@odata.context": "/redfish/v1/$metadata#Bios.Bios",
    "@odata.etag": "W/\"A8CFEAA3407A6E6E6E2970FA2E980355\"",
    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Bios/",
    "@odata.type": "#Bios.v1_0_0.Bios",
    "Actions": {
        "#Bios.ChangePassword": {
            "target": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Bios/Settings/Actions/Bios.ChangePasswords/"
        },
        "#Bios.ResetBios": {
            "target": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Bios/Settings/Actions/Bios.ResetBios/"
        }
    },
    "AttributeRegistry": "BiosAttributeRegistryU32.v1_2_22",
    "Attributes": {
        "AcpiHpet": "Disabled",
        "AcpiRootBridgePxm": "Enabled",
        "AcpiSlit": "Enabled",
        "AdjSecPrefetch": "Enabled",
        "AdminEmail": "admin2@someorg.com",
        "AdminName": "admin11",
        "AdminOtherInfo": "",
        "AdminPhone": "",
        "AdvCrashDumpMode": "Disabled",
        "AdvancedMemProtection": "FastFaultTolerantADDDC",
        "AsrStatus": "Enabled",
        "AsrTimeoutMinutes": "Timeout10",
        "AssetTagProtection": "Unlocked",
        "AutoPowerOn": "RestoreLastState",
        "BootMode": "Uefi",
        "BootOrderPolicy": "RetryIndefinitely",
        "ChannelInterleaving": "Enabled",
        "CollabPowerControl": "Enabled",
        "ConsistentDevNaming": "LomsAndSlots",
        "CustomPostMessage": "",
        "DaylightSavingsTime": "Disabled",
        "DcuIpPrefetcher": "Enabled",
        "DcuStreamPrefetcher": "Enabled",
        "Dhcpv4": "Enabled",
        "DirectToUpi": "Auto",
        "DynamicPowerCapping": "Disabled",
        "EmbNicEnable": "Auto",
        "EmbNicLinkSpeed": "Auto",
        "EmbNicPCIeOptionROM": "Enabled",
        "EmbSas1Aspm": "Disabled",
        "EmbSas1Boot": "TwentyFourTargets",
        "EmbSas1Enable": "Auto",
        "EmbSas1LinkSpeed": "Auto",
        "EmbSas1PcieOptionROM": "Enabled",
        "EmbSata1Aspm": "Disabled",
        "EmbSata2Aspm": "Disabled",
        "EmbVideoConnection": "Auto",
        "EmbeddedDiagnostics": "Enabled",
        "EmbeddedSata": "Ahci",
        "EmbeddedSerialPort": "Com2Irq3",
        "EmbeddedUefiShell": "Enabled",
        "EmsConsole": "Disabled",
        "EnabledCoresPerProc": 0,
        "EnergyEfficientTurbo": "Enabled",
        "EnergyPerfBias": "BalancedPerf",
        "EraseUserDefaults": "No",
        "ExtendedAmbientTemp": "Disabled",
        "ExtendedMemTest": "Disabled",
        "F11BootMenu": "Enabled",
        "FCScanPolicy": "CardConfig",
        "FanFailPolicy": "Shutdown",
        "FanInstallReq": "EnableMessaging",
        "FlexLom1Aspm": "Disabled",
        "FlexLom1Enable": "Auto",
        "FlexLom1LinkSpeed": "Auto",
        "FlexLom1PCIeOptionROM": "Enabled",
        "HttpSupport": "Auto",
        "HwPrefetcher": "Enabled",
        "IODCConfiguration": "Auto",
        "IntelDmiLinkFreq": "Auto",
        "IntelNicDmaChannels": "Enabled",
        "IntelPerfMonitoring": "Disabled",
        "IntelProcVtd": "Enabled",
        "IntelUpiFreq": "Auto",
        "IntelUpiLinkEn": "Auto",
        "IntelUpiPowerManagement": "Enabled",
        "IntelligentProvisioning": "Enabled",
        "InternalSDCardSlot": "Enabled",
        "Ipv4Address": "0.0.0.0",
        "Ipv4Gateway": "0.0.0.0",
        "Ipv4PrimaryDNS": "0.0.0.0",
        "Ipv4SecondaryDNS": "0.0.0.0",
        "Ipv4SubnetMask": "0.0.0.0",
        "Ipv6Address": "::",
        "Ipv6ConfigPolicy": "Automatic",
        "Ipv6Duid": "Auto",
        "Ipv6Gateway": "::",
        "Ipv6PrimaryDNS": "::",
        "Ipv6SecondaryDNS": "::",
        "LLCDeadLineAllocation": "Enabled",
        "LlcPrefetch": "Disabled",
        "LocalRemoteThreshold": "Auto",
        "MaxMemBusFreqMHz": "Auto",
        "MaxPcieSpeed": "PerPortCtrl",
        "MemClearWarmReset": "Disabled",
        "MemFastTraining": "Enabled",
        "MemMirrorMode": "Full",
        "MemPatrolScrubbing": "Enabled",
        "MemRefreshRate": "Refreshx1",
        "MemoryControllerInterleaving": "Auto",
        "MemoryRemap": "NoAction",
        "MinProcIdlePkgState": "C6Retention",
        "MinProcIdlePower": "C6",
        "MixedPowerSupplyReporting": "Enabled",
        "NetworkBootRetry": "Enabled",
        "NetworkBootRetryCount": 20,
        "NicBoot1": "NetworkBoot",
        "NicBoot2": "Disabled",
        "NicBoot3": "Disabled",
        "NicBoot4": "Disabled",
        "NicBoot5": "NetworkBoot",
        "NicBoot6": "Disabled",
        "NodeInterleaving": "Disabled",
        "NumaGroupSizeOpt": "Flat",
        "NvmeOptionRom": "Enabled",
        "OpportunisticSelfRefresh": "Disabled",
        "PciPeerToPeerSerialization": "Disabled",
        "PciResourcePadding": "Normal",
        "PciSlot1Bifurcation": "Auto",
        "PciSlot2Bifurcation": "Auto",
        "PciSlot3Bifurcation": "Auto",
        "PersistentMemBackupPowerPolicy": "WaitForBackupPower",
        "PostBootProgress": "Disabled",
        "PostDiscoveryMode": "Auto",
        "PostF1Prompt": "Delayed20Sec",
        "PostVideoSupport": "DisplayAll",
        "PowerButton": "Enabled",
        "PowerOnDelay": "NoDelay",
        "PowerRegulator": "DynamicPowerSavings",
        "PreBootNetwork": "Auto",
        "PrebootNetworkEnvPolicy": "Auto",
        "PrebootNetworkProxy": "",
        "ProcAes": "Enabled",
        "ProcHyperthreading": "Enabled",
        "ProcTurbo": "Enabled",
        "ProcVirtualization": "Enabled",
        "ProcX2Apic": "Enabled",
        "ProcessorConfigTDPLevel": "Normal",
        "ProcessorJitterControl": "Disabled",
        "ProcessorJitterControlFrequency": 0,
        "ProcessorJitterControlOptimization": "ZeroLatency",
        "ProductId": "867959-B21",
        "RedundantPowerSupply": "BalancedMode",
        "RemovableFlashBootSeq": "ExternalKeysFirst",
        "RestoreDefaults": "No",
        "RestoreManufacturingDefaults": "No",
        "RomSelection": "CurrentRom",
        "SataSecureErase": "Disabled",
        "SaveUserDefaults": "No",
        "SecStartBackupImage": "Disabled",
        "SecureBootStatus": "Disabled",
        "SerialConsoleBaudRate": "BaudRate115200",
        "SerialConsoleEmulation": "Vt100Plus",
        "SerialConsolePort": "Auto",
        "SerialNumber": "MXQ91100T6",
        "ServerAssetTag": "",
        "ServerConfigLockStatus": "Disabled",
        "ServerName": "SRVMXQ91100T6",
        "ServerOtherInfo": "",
        "ServerPrimaryOs": "",
        "ServiceEmail": "",
        "ServiceName": "",
        "ServiceOtherInfo": "",
        "ServicePhone": "",
        "SetupBrowserSelection": "Auto",
        "Slot1MctpBroadcastSupport": "Enabled",
        "Slot2MctpBroadcastSupport": "Enabled",
        "Slot3MctpBroadcastSupport": "Enabled",
        "Sriov": "Enabled",
        "StaleAtoS": "Disabled",
        "SubNumaClustering": "Disabled",
        "ThermalConfig": "OptimalCooling",
        "ThermalShutdown": "Enabled",
        "TimeFormat": "Utc",
        "TimeZone": "Unspecified",
        "TpmChipId": "None",
        "TpmFips": "FipsMode",
        "TpmState": "NotPresent",
        "TpmType": "NoTpm",
        "UefiOptimizedBoot": "Enabled",
        "UefiSerialDebugLevel": "Disabled",
        "UefiShellBootOrder": "Disabled",
        "UefiShellScriptVerification": "Disabled",
        "UefiShellStartup": "Disabled",
        "UefiShellStartupLocation": "Auto",
        "UefiShellStartupUrl": "",
        "UefiShellStartupUrlFromDhcp": "Disabled",
        "UncoreFreqScaling": "Auto",
        "UpiPrefetcher": "Enabled",
        "UrlBootFile": "",
        "UrlBootFile2": "",
        "UrlBootFile3": "",
        "UrlBootFile4": "",
        "UsbBoot": "Enabled",
        "UsbControl": "UsbEnabled",
        "UserDefaultsState": "Disabled",
        "UtilityLang": "English",
        "VirtualInstallDisk": "Disabled",
        "VirtualSerialPort": "Com1Irq4",
        "VlanControl": "Disabled",
        "VlanId": 0,
        "VlanPriority": 0,
        "WakeOnLan": "Enabled",
        "WorkloadProfile": "GeneralPowerEfficientCompute",
        "XptPrefetcher": "Auto",
        "iSCSIPolicy": "SoftwareInitiator"
    },
    "Id": "Bios",
    "Name": "BIOS Current Settings",
    "Oem": {
        "Hpe": {
            "@odata.type": "#HpeBiosExt.v2_0_0.HpeBiosExt",
            "Links": {
                "BaseConfigs": {
                    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Bios/baseconfigs/"
                },
                "Boot": {
                    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Bios/boot/"
                },
                "KmsConfig": {
                    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Bios/kmsconfig/"
                },
                "Mappings": {
                    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Bios/mappings/"
                },
                "ServerConfigLock": {
                    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Bios/serverconfiglock/"
                },
                "TlsConfig": {
                    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Bios/tlsconfig/"
                },
                "iScsi": {
                    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Bios/iscsi/"
                }
            },
            "SettingsObject": {
                "UnmodifiedETag": "W/\"E1E562A3BB8E1C1C1CBDC6F3070B67B2\""
            }
        }
    }
}
```
## Network interfaces

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}/NetworkInterfaces`   |
| **Description**    | This endpoint lists network interfaces of a specific system.<br> A network interface contains links to network adapter, network port, and network device function resources. |
| **Returns**        | List of network interface endpoints                          |
| **Response code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/NetworkInterfaces'
```
> **Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#NetworkInterfaceCollection.NetworkInterfaceCollection",
    "@odata.etag": "W/\"C321D970\"",
    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/NetworkInterfaces",
    "@odata.type": "#NetworkInterfaceCollection.NetworkInterfaceCollection",
    "Description": "The collection of network interfaces available in this system.",
    "Members": [
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/NetworkInterfaces/DC07A000"
        }
    ],
    "Members@odata.count": 1,
    "Name": "Network Interface Collection",
    "Oem": {
        "Hpe": {
            "@odata.context": "/redfish/v1/$metadata#HpeNetworkInterfaceStatus.HpeNetworkInterfaceStatus",
            "@odata.type": "#HpeNetworkInterfaceStatus.v1_0_0.HpeNetworkInterfaceStatus",
            "MemberContents": "AllDevices"
        }
    }
}
```
##  Ethernet interfaces

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}/ EthernetInterfaces` |
| **Description**    | This endpoint lists Ethernet interfaces or network interface controllers (NICs) of a specific system. |
| **Returns**        | List of Ethernet interface endpoints                         |
| **Response Code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces'
```
> **Sample response body**


```
{
    "@odata.context": "/redfish/v1/$metadata#EthernetInterfaceCollection.EthernetInterfaceCollection",
    "@odata.etag": "W/\"D5EC731D\"",
    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/EthernetInterfaces",
    "@odata.type": "#EthernetInterfaceCollection.EthernetInterfaceCollection",
    "Description": "Collection of System Ethernet Interfaces",
    "Members": [
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/EthernetInterfaces/1"
        },
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/EthernetInterfaces/2"
        },
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/EthernetInterfaces/3"
        },
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/EthernetInterfaces/4"
        },
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/EthernetInterfaces/5"
        },
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/EthernetInterfaces/6"
        }
    ],
    "Members@odata.count": 6,
    "Name": "System Ethernet Interfaces"
}
```
## Single Ethernet interface

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}/ EthernetInterfaces/{ethernetInterfaceId}` |
| **Description**    | This endpoint retrieves information on a single, logical Ethernet interface or network interface controller (NIC). |
| **Returns**        | JSON schema representing this Ethernet interface             |
| **Response code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces/{ethernetInterfaceId}'
```
> **Sample response body** 

```
{
    "@odata.context": "/redfish/v1/$metadata#EthernetInterface.EthernetInterface",
    "@odata.etag": "W/\"A04B8EF5\"",
    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/EthernetInterfaces/1",
    "@odata.type": "#EthernetInterface.v1_8_0.EthernetInterface",
    "FQDN": null,
    "FullDuplex": false,
    "HostName": null,
    "IPv4Addresses": [
    ],
    "IPv4StaticAddresses": [
    ],
    "IPv6AddressPolicyTable": [
    ],
    "IPv6Addresses": [
    ],
    "IPv6StaticAddresses": [
    ],
    "IPv6StaticDefaultGateways": [
    ],
    "Id": "1",
    "InterfaceEnabled": null,
    "LinkStatus": null,
    "MACAddress": "20:67:7c:e9:f6:40",
    "Name": "",
    "NameServers": [
    ],
    "SpeedMbps": null,
    "StaticNameServers": [
    ],
    "Status": {
        "Health": null,
        "State": null
    },
    "UefiDevicePath": "PciRoot(0x0)/Pci(0x1C,0x0)/Pci(0x0,0x0)"
}
```
##  PCIeDevice

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}/PCIeDevices/{PCIeDeviceId}` |
| **Description**    | This operation fetches information about a specific PCIe device.<br> |
| **Returns**        | Properties of a PCIe device attached to a computer system such as type, version of the PCIe specification in use by this device and so on. |
| **Response code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/PCIeDevices/{PCIeDeviceId}'
```
> **Sample response body**


```
{
    "@odata.context": "/redfish/v1/$metadata#PCIeDevice.PCIeDevice",
    "@odata.etag": "W/\"33150E20\"",
    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/PCIeDevices/1",
    "@odata.type": "#PCIeDevice.v1_9_0.PCIeDevice",
    "Id": "1",
    "Name": "HPE Ethernet 1Gb 4-port 331i Adapter - NIC",
    "Oem": {
        "Hpe": {
            "@odata.context": "/redfish/v1/$metadata#HpeServerPciDevice.HpeServerPciDevice",
            "@odata.etag": "W/\"33150E20\"",
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/PCIDevices/1",
            "@odata.type": "#HpeServerPciDevice.v2_0_0.HpeServerPciDevice",
            "Bifurcated": "BifurcationNotSupported",
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

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}/Storage`             |
| **Description**    | This operation lists storage subsystems.<br> A storage subsystem is a set of storage controllers (physical or virtual) and the resources such as volumes that can be accessed from that subsystem.<br> |
| **Returns**        | Links to storage subsystems                                  |
| **Response code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage'
```
> **Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#StorageCollection.StorageCollection",
    "@odata.etag": "W/\"AA6D42B0\"",
    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Storage",
    "@odata.type": "#StorageCollection.StorageCollection",
    "Description": "Storage Collection view",
    "Members": [
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Storage/ArrayControllers-0"
        }
    ],
    "Members@odata.count": 1,
    "Name": "Storage Collection"
}
```
## StoragePools 

The StoragePools schema represents storage pools, allocated volumes, and drives.

### StoragePools Collection

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools` |
| **Description**    | This operation returns a collection of StoragePool resource instances. |
| **Response code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools'
```

> **Sample response body**

```
{
	"@odata.etag": "\"2a840a57e9592422136\"",
	"@odata.id": "/redfish/v1/Systems/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1/Storage/RAID_Slot4/StoragePools",
	"@odata.type": "#StoragePoolCollection.StoragePoolCollection",
	"Description": "A collection of StoragePool resource instances.",
	"Members": [{
		"@odata.id": "/redfish/v1/Systems/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1/Storage/RAID_Slot4/StoragePools/Pool_1_27"
	}],
	"Members@odata.count": 1,
	"Name": "StoragePoolCollection"
}
```

### Single StoragePool

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}` |
| **Description**    | This operation represents a single StoragePool instance.     |
| **Response code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```curl -i GET \
 curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}'
```

> **Sample response body**

```
{
	"@odata.etag": "\"7a55c980802224f456c\"",
	"@odata.id": "/redfish/v1/Systems/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1/Storage/RAID_Slot4/StoragePools/Pool_1_27",
	"@odata.type": "#StoragePool.v1_5_0.StoragePool",
	"AllocatedVolumes": {
		"@odata.id": "/redfish/v1/Systems/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1/Storage/RAID_Slot4/StoragePools/Pool_1_27/AllocatedVolumes"
	},
	"Capacity": {
		"Data": {
			"AllocatedBytes": 998999326720,
			"ConsumedBytes": 998999326720
		},
		"Metadata": {},
		"Snapshot": {}
	},
	"CapacitySources": [{
		"@odata.etag": "\"31ecba89b48925a6bf3\"",
		"@odata.id": "/redfish/v1/Systems/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1/Storage/RAID_Slot4/StoragePools/Pool_1_27/CapacitySources/1",
		"@odata.type": "#Capacity.v1_1_3.CapacitySource",
		"Description": "The resource is used to represent a capacity for a Redfish implementation.",
		"Id": "1",
		"Name": "CapacitySources_1",
		"ProvidingDrives": {
			"@odata.id": "/redfish/v1/Systems/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1/Storage/RAID_Slot4/StoragePools/Pool_1_27/CapacitySources/1/ProvidingDrives"
		}
	}],
	"CapacitySources@odata.count": 1,
	"Description": "The resource is used to represent a storage pool for a Redfish implementation.",
	"Id": "Pool_1_27",
	"Name": "Pool_1_27",
	"Status": {
		"State": "Enabled"
	},
	"SupportedRAIDTypes": ["RAID0"]
}
```

### AllocatedVolumes Collection

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/AllocatedVolumes` |
| **Description**    | This operation returns a collection of volume resource instances. |
| **Response code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```curl -i GET \
 curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/AllocatedVolumes'
```

> **Sample response body**

```
{
   "@odata.etag":"\"2cbdffec21f02963c79\"",
   "@odata.id":"/redfish/v1/Systems/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1/Storage/RAID_Slot4/StoragePools/Pool_1_27/AllocatedVolumes",
   "@odata.type":"#VolumeCollection.VolumeCollection",
   "Description":"A collection of volume resource instances.",
   "Members":[
      {
         "@odata.id":"/redfish/v1/Systems/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1/Storage/RAID_Slot4/StoragePools/Pool_1_27/AllocatedVolumes/27"
      }
   ],
   "Members@odata.count":1,
   "Name":"VolumeCollection"
}
```

### Single AllocatedVolume

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/AllocatedVolumes/{allocatedvolumes_Id}` |
| **Description**    | This operation represents a single volume instance.          |
| **Response code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```curl -i GET \
 curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/AllocatedVolumes/{allocatedvolumes_Id}'
```

> **Sample response body**

```
{
   "@odata.etag":"\"915af5f726a127f66c2\"",
   "@odata.id":"/redfish/v1/Systems/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1/Storage/RAID_Slot4/StoragePools/Pool_1_27/AllocatedVolumes/27",
   "@odata.type":"#Volume.v1_6_2.Volume",
   "AccessCapabilities":[
   ],
   "Actions":{
      "#Volume.Initialize":{
         "InitializeType@Redfish.AllowableValues":[
            "Fast"
         ],
         "target":"/redfish/v1/Systems/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1/Storage/RAID_Slot4/StoragePools/Pool_1_27/AllocatedVolumes/27/Actions/Volume.Initialize",
         "title":"Initialize"
      }
   },
   "BlockSizeBytes":512,
   "Capacity":{
      "Data":{
         
      },
      "Metadata":{
         
      },
      "Snapshot":{
         
      }
   },
   "CapacityBytes":998999326720,
   "Description":"This resource is used to represent a volume for a Redfish implementation.",
   "DisplayName":"VD_1",
   "Id":"27",
   "Links":{
      "Drives":[
         {
            "@odata.id":"/redfish/v1/Systems/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1/Storage/RAID_Slot4/Drives/Disk.0"
         }
      ]
   },
   "Name":"VD_1",
   "Oem":{
      "Lenovo":{
         "@odata.type":"#LenovoStorageVolume.v1_0_0.LenovoStorageVolume",
         "AccessPolicy":"",
         "Bootable":true,
         "DriveCachePolicy":"",
         "…"
      }
   },
   "RAIDType":"RAID0",
   "ReadCachePolicy":null,
   "ReadCachePolicy@Redfish.AllowableValues":[
      "Off",
      "ReadAhead"
   ],
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   },
   "StripSizeBytes":0,
   "WriteCachePolicy":null,
   "WriteCachePolicy@Redfish.AllowableValues":[
      "WriteThrough",
      "UnprotectedWriteBack",
      "ProtectedWriteBack"
   ]
}
```

### ProvidingDrives Collection

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/CapacitySources/{capacitysources_Id}/ProvidingDrives` |
| **Description**    | This operation returns a collection of drives.               |
| **Response code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```curl -i GET \
 curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/CapacitySources/{capacitysources_Id}/ProvidingDrives'

```

> **Sample response body**

```
{
   "@odata.etag":"\"252e43d3e5862ae3d30\"",
   "@odata.id":"/redfish/v1/Systems/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1/Storage/RAID_Slot4/StoragePools/Pool_1_27/CapacitySources/1/ProvidingDrives",
   "@odata.type":"#DriveCollection.DriveCollection",
   "Members":[
      {
         "@odata.id":"/redfish/v1/Systems/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1/Storage/RAID_Slot4/Drives/Disk.0"
      }
   ],
   "Members@odata.count":1,
   "Name":"DriveCollection"
}
```

### Single ProvidingDrive

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/CapacitySources/{capacitysources_Id}/ProvidingDrives/{providingdrives_id}` |
| **Description**    | This operation represents a single drive instance.           |
| **Response code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

> **curl command**

```curl -i GET \
 curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{StorageControllerId}/StoragePools/{storagepool_Id}/CapacitySources/{capacitysources_Id}/ProvidingDrives/{providingdrives_id}'
```

> **Sample response body**

```
{
   "@odata.etag":"\"a81c2a4f4e972e18b6306\"",
   "@odata.id":"/redfish/v1/Systems/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1/Storage/RAID_Slot4/Drives/Disk.0",
   "@odata.type":"#Drive.v1_13_1.Drive",
   "AssetTag":"",
   "BlockSizeBytes":512,
   "CapableSpeedGbs":6,
   "CapacityBytes":1000204886016,
   "Description":"This resource is used to represent a drive for a Redfish implementation.",
   "EncryptionAbility":"None",
   "EncryptionStatus":"Unencrypted",
   "FailurePredicted":false,
   "HotspareType":"None",
   "Id":"Disk.0",
   "Identifiers":[
      {
         "DurableName":"",
         "DurableNameFormat":"UUID"
      }
   ],
   "Links":{
      "Chassis":{
         "@odata.id":"/redfish/v1/Chassis/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1"
      },
      "PCIeFunctions":[
         
      ],
      "Volumes":[
         {
            "@odata.id":"/redfish/v1/Systems/8b9da958-52d7-4f33-a01a-74b6ab4d3886.1/Storage/RAID_Slot4/Volumes/27"
         }
      ]
   },
   "Manufacturer":"Seagate",
   "MediaType":"HDD",
   "Model":"ST1000NX0423",
   "Name":"1.00TB 7.2K 6Gbps SATA 2.5 HDD",
   "NegotiatedSpeedGbs":6,
   "Oem":{
      "Lenovo":{
         "@odata.type":"#LenovoDrive.v1_0_0.LenovoDrive",
         "DriveStatus":"Online",
         "Temperature":27
      }
   },
   "PartNumber":"D7A01874",
   "PhysicalLocation":{
      "Info":"Slot 0",
      "Info@Redfish.Deprecated":"The property is deprecated. Please use PartLocation instead.",
      "InfoFormat":"Slot Number",
      "InfoFormat@Redfish.Deprecated":"The property is deprecated. Please use PartLocation instead.",
      "PartLocation":{
         "LocationOrdinalValue":0,
         "LocationType":"Bay",
         "ServiceLabel":"Drive 0"
      }
   },
   "PredictedMediaLifeLeftPercent":null,
   "Protocol":"SATA",
   "Revision":"LEK9",
   "RotationSpeedRPM":7200,
   "SKU":"00YK025",
   "SerialNumber":"W473AMB3",
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   },
   "StatusIndicator":null
}
```

##  Storage subsystem

|||
|---------|-------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}` |
|**Description** | This operation lists resources such as drives and storage controllers in a storage subsystem. |
|**Returns** |Links to the drives and storage controllers of a storage subsystem|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}'
```

> **Sample response body** 

```
{
    "@odata.context": "/redfish/v1/$metadata#Storage.Storage",
    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Storage/ArrayControllers-0",
    "@odata.type": "#Storage.v1_13_0.Storage",
    "Description": "HPE Smart Storage Array Controller View",
    "Drives": [
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Storage/ArrayControllers-0/Drives/0"
        },
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Storage/ArrayControllers-0/Drives/1"
        },
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Storage/ArrayControllers-0/Drives/2"
        }
    ],
    "Id": "ArrayControllers-0",
    "Name": "Hpe Smart Storage Array Controller",
    "StorageControllers": [
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Storage/ArrayControllers-0#/StorageControllers/0",
            "@odata.type": "#Storage.v1_13_0.Storage",
            "FirmwareVersion": "2.65",
            "Location": {
                "PartLocation": {
                    "LocationOrdinalValue": 0,
                    "LocationType": "Slot",
                    "ServiceLabel": "Slot=0"
                }
            },
            "Manufacturer": "HPE",
            "MemberId": "0",
            "Model": "HPE Smart Array P408i-a SR Gen10",
            "Name": "Hpe Smart Storage Array Controller",
            "PartNumber": "836260-001",
            "SerialNumber": "PEYHC0DRHBV947 ",
            "Status": {
                "Health": "OK",
                "State": "Enabled"
            }
        }
    ],
    "Volumes": {
        "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Storage/ArrayControllers-0/Volumes"
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
|**Returns** |JSON schema representing this drive|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Drives/{driveId}'
```

> **Sample response body** 

```
{
    "@odata.context": "/redfish/v1/$metadata#Drive.Drive",
    "@odata.etag": "W/\"990C0D8A\"",
    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Storage/ArrayControllers-0/Drives/0",
    "@odata.type": "#Drive.v1_15_0.Drive",
    "BlockSizeBytes": 512,
    "CapacityBytes": 1200000000000,
    "Description": "HPE Smart Storage Disk Drive View",
    "Id": "0",
    "Links": {
        "Volumes": [
            {
                "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Storage/ArrayControllers-0/Volumes/1"
            }
        ]
    },
    "MediaType": "HDD",
    "Model": "EG001200JWJNQ",
    "Name": "HpeStorageDiskDrive",
    "PhysicalLocation": {
        "PartLocation": {
            "LocationOrdinalValue": 1,
            "LocationType": "Bay",
            "ServiceLabel": "Port=1I:Box=1:Bay=1:LegacyBootPriority=None"
        }
    },
    "Revision": "HPD3",
    "RotationSpeedRPM": 10500,
    "SerialNumber": "WFK25Z6F",
    "Status": {
        "Health": "OK",
        "State": "Enabled"
    }
}
```



## Volumes

The volume schema represents a volume, virtual disk, LUN, or other logical storage entity for a system.


### Collection of volumes

| | |
|----------|-----------|
|<strong>Method</strong> |`GET` |
|<strong>URI</strong>  |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes` |
|<strong>Description</strong>  |This endpoint retrieves a collection of volumes in a specific storage subsystem.|
|<strong>Returns</strong> |A list of links to volumes|
|<strong>Response code</strong> |On success, `200 OK` |
|<strong>Authentication</strong> |Yes|

 **curl command** 

```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odim_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes'
```

>**Sample response body** 

```
{
    "@Redfish.CollectionCapabilities": {
        "@odata.type": "#CollectionCapabilities.v1_4_0.CollectionCapabilities",
        "Capabilities": [
            {
                "CapabilitiesObject": {
                    "@odata.id": "/redfish/v1/Systems/64992250-2a1a-41c6-82c6-b046140d615d.1/Storage/ArrayControllers-0/Volumes/Capabilities"
                },
                "Links": {
                    "TargetCollection": {
                        "@odata.id": "/redfish/v1/Systems/64992250-2a1a-41c6-82c6-b046140d615d.1/Storage/ArrayControllers-0/Volumes"
                    }
                },
                "UseCase": "VolumeCreation"
            }
        ]
    },
    "@odata.context": "/redfish/v1/$metadata#VolumeCollection.VolumeCollection",
    "@odata.etag": "W/\"AA6D42B0\"",
    "@odata.id": "/redfish/v1/Systems/64992250-2a1a-41c6-82c6-b046140d615d.1/Storage/ArrayControllers-0/Volumes",
    "@odata.type": "#VolumeCollection.VolumeCollection",
    "Description": "Volume Collection view",
    "Members": [
        {
            "@odata.id": "/redfish/v1/Systems/64992250-2a1a-41c6-82c6-b046140d615d.1/Storage/ArrayControllers-0/Volumes/1"
        },
        {
            "@odata.id": "/redfish/v1/Systems/64992250-2a1a-41c6-82c6-b046140d615d.1/Storage/ArrayControllers-0/Volumes/2"
        }
    ],
    "Members@odata.count": 2,
    "Name": "Volume Collection"
}
```

### Viewing volume capabilities

|                                 |                                                              |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>Method</strong>         | `GET`                                                        |
| <strong>URI</strong>            | `/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/Capabilities` |
| <strong>Description</strong>    | This operation displays all allowed property values you can use while creating a volume. |
| <strong>Returns</strong>        | JSON schema representing this volume                         |
| <strong>Response code</strong>  | On success, `200 OK`                                         |
| <strong>Authentication</strong> | Yes                                                          |


>**curl command**


```
curl -i -X GET \
   -H "Authorization:Basic YWRtaW46T2QhbTEyJDQ=" \
   -H "Content-Type:application/json" \
 'https://{odim_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/Capabilities'
```

>**Sample response body** 

```
{
    "@odata.id": "/redfish/v1/Systems/45201b16-5305-49f0-846b-4597e982f6f8.1/Storage/DE00C000/Volumes/Capabilities",
    "@odata.type": "#Volume.v1_6_2.Volume",
    "Id": "Capabilities",
    "Links": {
        "Drives@Redfish.RequiredOnCreate": true
    },
    "Links@Redfish.RequiredOnCreate": true,
    "Name": "Capabilities for the volume collection",
    "RAIDType@Redfish.AllowableValues": [
        "RAID0",
        "RAID1",
        "RAID10",
        "RAID5",
        "RAID50",
        "RAID6",
        "RAID60",
        "RAID1Triple",
        "RAID10Triple"
    ],
    "RAIDType@Redfish.RequiredOnCreate": true
}
```

### Single volume


| | |
|----------|-----------|
|<strong>Method</strong> |`GET` |
|<strong>URI</strong>   |`/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes/{volumeId}` |
|<strong>Description</strong>   |This endpoint retrieves information about a specific volume in a storage subsystem.|
|<strong>Returns</strong>  |JSON schema representing this volume|
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
   "@odata.id":"/redfish/v1/Systems/363bef34-7f89-48ac-8970-ee8955f1b56f.1/Storage/ArrayControllers-0/Volumes/1",
   "@odata.type":"#Volume.v1_6_2.Volume",
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
            "@odata.id":"/redfish/v1/Systems/363bef34-7f89-48ac-8970-ee8955f1b56f.1/Storage/ArrayControllers-0/Drives/0"
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
   "RAIDType":"RAID1",
   "Links":{
     "Drives":[
      {
         "@odata.id":"/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Drives/0"
      },
      {
         "@odata.id":"/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Drives/1"
      }
   ]
 }, 
   "@Redfish.OperationApplyTime":"OnReset"
}' \
 'https://{odim_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Storage/{storageSubsystemId}/Volumes'
```

>**Sample request body** 

```
{
   "RAIDType":"RAID1",
   "Links":{
     "Drives":[
      {
         "@odata.id":"/redfish/v1/Systems/363bef34-7f89-48ac-8970-ee8955f1b56f.1/Storage/ArrayControllers-0/Drives/0"
      },
      {
         "@odata.id":"/redfish/v1/Systems/363bef34-7f89-48ac-8970-ee8955f1b56f.1/Storage/ArrayControllers-0/Drives/1"
      }
   ]
  },
   "DisplayName":"<volume_name>"
   "@Redfish.OperationApplyTime":"OnReset"
}
```

> **Request parameters** 

|Parameter|Type|Description|
|---------|----|-----------|
|RAIDType|String (required)<br>|The RAID type of the volume you want to create.|
|Links {|Object (required)|Links to individual drives.|
|Drives[{|Array (required)<br> |An array of links to drive resources to contain the new volume.|
|@odata.id }]}<br> |String|A link to a drive resource.|
|DisplayName |String|Name of the volume (optional).|
|@Redfish.OperationApplyTime|Redfish annotation (optional)<br> | It enables you to control when the operation is carried out.<br> Supported values: `OnReset` and `Immediate`.<br> `OnReset` indicates that the new volume is available only after you successfully reset the system. To know how to reset a system, see [Resetting a computer system](#resetting-a-computer-system).<br>`Immediate` indicates that the created volume is available in the system immediately after the operation is successfully complete. |

>**Sample response body** 

```
 {
      "error":{
            "@Message.ExtendedInfo":[
                  {
                        "MessageId": "Base.1.13.Success"            
         }         
      ],
            "code":"iLO.0.10.ExtendedInfo",
            "message":"See @Message.ExtendedInfo for more information."      
   }   
}
```

> **NOTE**: Reset your system only if prompted in your response message id. After the system reset, the new volume is available. In case of successful message id in the response, system reset is not required.

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

> **Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|@Redfish.OperationApplyTime|Redfish annotation (optional)<br> | It enables you to control when the operation is carried out.<br> Supported values are: `OnReset` and `Immediate`. `OnReset` indicates that the volume is deleted only after you successfully reset the system.<br> `Immediate` indicates that the volume is deleted immediately after the operation is successfully complete. |


##  SecureBoot

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/SecureBoot` |
|**Description** |Use this endpoint to discover information on `UEFI Secure Boot` and manage the `UEFI Secure Boot` functionality of a specific system.|
|**Returns** | <ul><li>Action for resetting keys</li><li> `UEFI Secure Boot` properties<br>**NOTE:** Use URI in the *Actions* group to discover information about resetting keys.</li></ul> |
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/SecureBoot'
```

> **Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#SecureBoot.SecureBoot",
    "@odata.etag": "W/\"4A4CB737\"",
    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/SecureBoot",
    "@odata.type": "#SecureBoot.v1_0_0.SecureBoot",
    "Actions": {
        "#SecureBoot.ResetKeys": {
            "target": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/SecureBoot/Actions/SecureBoot.ResetKeys"
        }
    },
    "Id": "SecureBoot",
    "Name": "SecureBoot",
    "SecureBootCurrentBoot": "Disabled",
    "SecureBootEnable": false,
    "SecureBootMode": "UserMode"
}
```

##  Processors

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Processors` |
|**Description** |This endpoint lists processors of a specific system.|
|**Returns** |List of processor resource endpoints|
|**Response code** |`200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Processors'
```

> **Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#ProcessorCollection.ProcessorCollection",
    "@odata.etag": "W/\"570254F2\"",
    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Processors",
    "@odata.type": "#ProcessorCollection.ProcessorCollection",
    "Description": "Processors view",
    "Members": [
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Processors/1"
        },
        {
            "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Processors/2"
        }
    ],
    "Members@odata.count": 2,
    "Name": "Processors Collection"
}
```

### Single processor

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems/{ComputerSystemId}/Processors/{processoId}` |
|**Description** |This endpoint fetches information about the properties of a processor attached to a specific server.|
|**Returns** |JSON schema representing this processor|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
         -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}/Processors/{processoId}'
```

> **Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#Processor.Processor",
    "@odata.etag": "W/\"18ABF8BD\"",
    "@odata.id": "/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Processors/1",
    "@odata.type": "#Processor.v1_7_2.Processor",
    "Id": "1",
    "InstructionSet": "x86-64",
    "Manufacturer": "Intel(R) Corporation",
    "MaxSpeedMHz": 4000,
    "Model": "Intel(R) Xeon(R) Gold 6152 CPU @ 2.10GHz",
    "Name": "Processors",
    "Oem": {
        "Hpe": {
            "@odata.context": "/redfish/v1/$metadata#HpeProcessorExt.HpeProcessorExt",
            "@odata.type": "#HpeProcessorExt.v2_0_0.HpeProcessorExt",
            "AssetTag": "UNKNOWN",
            "Cache": [
                {
                    "Associativity": "8waySetAssociative",
                    "CacheSpeedns": 0,
                    "CurrentSRAMType": [
                        "Synchronous"
                    ],
                    "EccType": "SingleBitECC",
                    "InstalledSizeKB": 1408,
                    "Location": "Internal",
                    "MaximumSizeKB": 1408,
                    "Name": "L1-Cache",
                    "Policy": "WriteBack",
                    "Socketed": false,
                    "SupportedSRAMType": [
                        "Synchronous"
                    ],
                    "SystemCacheType": "Unified"
                },
                {
                    "Associativity": "16waySetAssociative",
                    "CacheSpeedns": 0,
                    "CurrentSRAMType": [
                        "Synchronous"
                    ],
                    "EccType": "SingleBitECC",
                    "InstalledSizeKB": 22528,
                    "Location": "Internal",
                    "MaximumSizeKB": 22528,
                    "Name": "L2-Cache",
                    "Policy": "Varies",
                    "Socketed": false,
                    "SupportedSRAMType": [
                        "Synchronous"
                    ],
                    "SystemCacheType": "Unified"
                },
                {
                    "Associativity": "FullyAssociative",
                    "CacheSpeedns": 0,
                    "CurrentSRAMType": [
                        "Synchronous"
                    ],
                    "EccType": "SingleBitECC",
                    "InstalledSizeKB": 30976,
                    "Location": "Internal",
                    "MaximumSizeKB": 30976,
                    "Name": "L3-Cache",
                    "Policy": "Varies",
                    "Socketed": false,
                    "SupportedSRAMType": [
                        "Synchronous"
                    ],
                    "SystemCacheType": "Unified"
                }
            ],
            "Characteristics": [
                "64Bit",
                "MultiCore",
                "HwThread",
                "ExecuteProtection",
                "EnhancedVirtualization",
                "PowerPerfControl"
            ],
            "ConfigStatus": {
                "Populated": true,
                "State": "Enabled"
            },
            "CoresEnabled": 22,
            "ExternalClockMHz": 100,
            "MicrocodePatches": [
                {
                    "CpuId": "0x00050654",
                    "Date": "2019-09-05T00:00:00Z",
                    "PatchId": "0x02000065"
                },
                {
                    "CpuId": "0x00050655",
                    "Date": "2018-10-08T00:00:00Z",
                    "PatchId": "0x0300000F"
                },
                {
                    "CpuId": "0x00050656",
                    "Date": "2019-09-05T00:00:00Z",
                    "PatchId": "0x0400002C"
                },
                {
                    "CpuId": "0x00050657",
                    "Date": "2019-09-05T00:00:00Z",
                    "PatchId": "0x0500002C"
                }
            ],
            "PartNumber": "",
            "RatedSpeedMHz": 2100,
            "SerialNumber": "",
            "VoltageVoltsX10": 16
        }
    },
    "PartNumber": "",
    "ProcessorArchitecture": "x86",
    "ProcessorId": {
        "EffectiveFamily": "179",
        "EffectiveModel": "5",
        "IdentificationRegisters": "0x06540005fbffbfeb",
        "MicrocodeInfo": null,
        "Step": "4",
        "VendorId": "Intel(R) Corporation"
    },
    "ProcessorType": "CPU",
    "SerialNumber": "",
    "Socket": "Proc 1",
    "Status": {
        "Health": "OK",
        "State": "Enabled"
    },
    "TotalCores": 22,
    "TotalThreads": 44
}
```



## Chassis

Chassis represents the physical components of a system—sheet-metal confined spaces, logical zones such as racks, enclosures, chassis and all other containers, and subsystems (like sensors).

To view, create, and manage racks or rack groups, ensure that the URP is running and is added into the Resource Aggregator for ODIM framework. To know how to add a plugin, see *[Adding a plugin as an aggregation source](#adding-a-plugin-as-an-aggregation-source)*.

>**NOTE:** URP is automatically installed during the Resource Aggregator for ODIM deployment.


### Collection of chassis

|||
|-------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis` |
|**Description** | This operation lists chassis instances available with Resource Aggregator for ODIM. |
|**Returns** |A collection of links to chassis instances|
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
         "@odata.id":"/redfish/v1/Chassis/ba0a6871-7bc4-5f7a-903d-67f3c205b08c.1"
      },
      { 
         "@odata.id":"/redfish/v1/Chassis/7ff3bd97-c41c-5de0-937d-85d390691b73.1"
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
         "@odata.id":"/redfish/v1/Chassis/ba0a6871-7bc4-5f7a-903d-67f3c205b08c.1"
      },
      { 
         "@odata.id":"/redfish/v1/Chassis/7ff3bd97-c41c-5de0-937d-85d390691b73.1"
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
|**Returns** |JSON schema representing this chassis instance|
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
   "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1",
   "@odata.type":"#Chassis.v1_20_0.Chassis",
   "@odata.etag":"W/\"59209823\"",
   "Id":"b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1",
   "Name":"Computer System Chassis",
   "AssetTag":null,
   "ChassisType":"RackMount",
   "IndicatorLED":"Off",
   "Manufacturer":"HPE",
   "Model":"ProLiant DL360 Gen10",
   "PartNumber":null,
   "PowerState":"On",
   "SerialNumber":"MXQ91100T6",
   "SKU":"867959-B21",
   "Links":{
      "ComputerSystems":[
         {
            "@odata.id":"/redfish/v1/Systems/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1"
         }
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1"
         }
      ]
   },
   "NetworkAdapters":{
      "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/NetworkAdapters"
   },
   "PCIeSlots":{
      "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/PCIeSlots"
   },
   "Power":{
      "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Power"
   },
   "Status":{
      "Health":"OK",
      "State":"Starting"
   },
   "Thermal":{
      "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Thermal"
   },
   "Oem":{
      "Hpe":{
         "@odata.context":"/redfish/v1/$metadata#HpeServerChassis.HpeServerChassis",
         "@odata.type":"#HpeServerChassis.v2_3_1.HpeServerChassis",
         "Actions":{
            "#HpeServerChassis.DisableMCTPOnServer":{
               "target":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Actions/Oem/Hpe/HpeServerChassis.DisableMCTPOnServer"
            },
            "#HpeServerChassis.FactoryResetMCTP":{
               "target":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Actions/Oem/Hpe/HpeServerChassis.FactoryResetMCTP"
            }
         },
         "ElConfigOverride":false,
         "Firmware":{
            "PlatformDefinitionTable":{
               "Current":{
                  "VersionString":"9.8.0 Build 15"
               }
            },
            "PowerManagementController":{
               "Current":{
                  "VersionString":"1.0.7"
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
                  "VersionString":"4.1.4.601"
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
               "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Devices"
            }
         },
         "SmartStorageBattery":[
            {
               "ChargeLevelPercent":99,
               "FirmwareVersion":"0.70",
               "Index":1,
               "MaximumCapWatts":96,
               "Model":"875241-B21",
               "ProductName":"HPE Smart Storage Battery ",
               "RemainingChargeTimeSeconds":37,
               "SerialNumber":"6WQXL0CB2BX63Z",
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
   "PCIeDevices":{
      "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/PCIeDevices"
   },
   "ThermalManagedByParent.omitempty":false
}
```

2. **Rack group chassis**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/22804541-c439-5d2a-81d5-23d23e0ebe38",
   "@odata.type":"#Chassis.v1_20_0.Chassis",
   "Id":"22804541-c439-5d2a-81d5-23d23e0ebe38",
   "Description":"My RackGroup",
   "Name":"RG2",
   "ChassisType":"RackGroup",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/b44b87c0-00de-4184-ad2b-cdd4da52a805"
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
   "@odata.id":"/redfish/v1/Chassis/f03fed09-dd75-5585-ad81-75cd4ae6266a",
   "@odata.type":"#Chassis.v1_20_0.Chassis",
   "Id":"f03fed09-dd75-5585-ad81-75cd4ae6266a",
   "Description":"My RackGroup",
   "Name":"RG_2",
   "ChassisType":"RackGroup",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/b44b87c0-00de-4184-ad2b-cdd4da52a805"
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

> **Sample response body**

```
{
    "@odata.context": "/redfish/v1/$metadata#Thermal.Thermal",
    "@odata.etag": "W/\"B51E22EA\"",
    "@odata.id": "/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Thermal",
    "@odata.type": "#Thermal.v1_6_2.Thermal",
    "Fans": [
        {
            "@odata.id": "/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Thermal#Fans/0",
            "MemberId": "0",
            "Name": "Fan 1",
            "Oem": {
                "Hpe": {
                    "@odata.context": "/redfish/v1/$metadata#HpeServerFan.HpeServerFan",
                    "@odata.type": "#HpeServerFan.v2_0_0.HpeServerFan",
                    "HotPluggable": true,
                    "Location": "System",
                    "Redundant": true
                }
            },
            "Reading": 30,
            "ReadingUnits": "Percent",
            "Status": {
                "Health": "OK",
                "State": "Enabled"
            }
        },
        {
            "@odata.id": "/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Thermal#Fans/1",
            "MemberId": "1",
            "Name": "Fan 2",
            "Oem": {
                "Hpe": {
                    "@odata.context": "/redfish/v1/$metadata#HpeServerFan.HpeServerFan",
                    "@odata.type": "#HpeServerFan.v2_0_0.HpeServerFan",
                    "HotPluggable": true,
                    "Location": "System",
                    "Redundant": true
                }
            },
            "Reading": 30,
            "ReadingUnits": "Percent",
            "Status": {
                "Health": "OK",
                "State": "Enabled"
            }
        }
    ],
    "Id": "Thermal",
    "Name": "Thermal",
    "Oem": {
        "Hpe": {
            "@odata.context": "/redfish/v1/$metadata#HpeThermalExt.HpeThermalExt",
            "@odata.type": "#HpeThermalExt.v2_0_0.HpeThermalExt",
            "Actions": {
            },
            "FanPercentMinimum": 0,
            "ThermalConfiguration": "OptimalCooling"
        }
    },
    "Temperatures": [
        {
            "@odata.id": "/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Thermal#Temperatures/1",
            "MemberId": "1",
            "Name": "02-CPU 1",
            "Oem": {
                "Hpe": {
                    "@odata.context": "/redfish/v1/$metadata#HpeSeaOfSensors.HpeSeaOfSensors",
                    "@odata.type": "#HpeSeaOfSensors.v2_0_0.HpeSeaOfSensors",
                    "LocationXmm": 11,
                    "LocationYmm": 5
                }
            },
            "PhysicalContext": "CPU",
            "ReadingCelsius": 40,
            "SensorNumber": 2,
            "Status": {
                "Health": "OK",
                "State": "Enabled"
            },
            "UpperThresholdCritical": 70,
            "UpperThresholdFatal": null
        },
        {
            "@odata.id": "/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Thermal#Temperatures/2",
            "MemberId": "2",
            "Name": "03-CPU 2",
            "Oem": {
                "Hpe": {
                    "@odata.context": "/redfish/v1/$metadata#HpeSeaOfSensors.HpeSeaOfSensors",
                    "@odata.type": "#HpeSeaOfSensors.v2_0_0.HpeSeaOfSensors",
                    "LocationXmm": 4,
                    "LocationYmm": 5
                }
            },
            "PhysicalContext": "CPU",
            "ReadingCelsius": 40,
            "SensorNumber": 3,
            "Status": {
                "Health": "OK",
                "State": "Enabled"
            },
            "UpperThresholdCritical": 70,
            "UpperThresholdFatal": null
        }
    ]
}
```



### Collection of network adapters

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}/NetworkAdapters`|
|**Description** | This endpoint lists network adapters contained in a chassis. A `NetworkAdapter` represents the physical network adapter capable of connecting to a computer network.<br> Some examples include Ethernet, fibre channel, and converged network adapters.|
|**Returns** |Links to network adapter instances available in this chassis|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}/NetworkAdapters'
```

> **Sample response body**

```
{   "@odata.context":"/redfish/v1/$metadata#NetworkAdapterCollection.NetworkAdapterCollection",
   "@odata.etag":"W/\"C321D970\"",
   "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/NetworkAdapters",
   "@odata.type":"#NetworkAdapterCollection.NetworkAdapterCollection",
   "Description":"The collection of network adapter resource instances available in this chassis.",
   "Members":[
      {
         "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/NetworkAdapters/DC07A000"
      }
   ],
   "Members@odata.count":1,
   "Name":"NetworkAdapterCollection",
   "Oem":{
      "Hpe":{
         "@odata.context":"/redfish/v1/$metadata#HpeNetworkAdapterStatus.HpeNetworkAdapterStatus",
         "@odata.type":"#HpeNetworkAdapterStatus.v1_0_0.HpeNetworkAdapterStatus",
         "MemberContents":"AllDevices"
      }
   }
}
```

### Single network adapter

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}/NetworkAdapters/{NetworkAdapterId}` |
|**Description** | This endpoint retrieves information on a specific network adapter.|
|**Returns** |JSON schema representing this network adapter|
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
   "@odata.etag":"W/\"DA153FEC\"",
   "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/NetworkAdapters/DC07A000",
   "@odata.type":"#NetworkAdapter.v1_4_0.NetworkAdapter",
   "Actions":{
      "#NetworkAdapter.ResetSettingsToDefault":{
         "target":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/NetworkAdapters/DC07A000/Actions/NetworkAdapter.ResetSettingsToDefault",
         "title":"Reset network adapter configuration to factory default values."
      }
   },
   "Controllers":[
      {
         "ControllerCapabilities":{
            "DataCenterBridging":{
               "Capable":true
            },
            "NPAR":{
               "NparCapable":false,
               "NparEnabled":false
            },
            "NPIV":{
               "MaxDeviceLogins":128,
               "MaxPortLogins":64
            },
            "NetworkDeviceFunctionCount":8,
            "NetworkPortCount":2,
            "VirtualizationOffload":{
               "SRIOV":{
                  "SRIOVVEPACapable":false
               },
               "VirtualFunction":{
                  "DeviceMaxCount":128,
                  "MinAssignmentGroupSize":8,
                  "NetworkPortMaxCount":64
               }
            }
         },
         "FirmwarePackageVersion":"07.18.27.00",
         "Location":{
            "PartLocation":{
               "LocationOrdinalValue":0,
               "LocationType":"Slot",
               "ServiceLabel":"Embedded ALOM"
            }
         }
      }
   ],
   "Description":"Device capabilities and characteristics with active configuration status",
   "Id":"DC07A000",
   "Manufacturer":"Hewlett Packard Enterprise",
   "Model":"HP FlexFabric 10Gb 2-port 534FLR-SFP+ Adapter",
   "Name":"HP FlexFabric 10Gb 2port 534FLR-SFP+ Adapter",
   "NetworkDeviceFunctions":{
      "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/NetworkAdapters/DC07A000/NetworkDeviceFunctions"
   },
   "NetworkPorts":{
      "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/NetworkAdapters/DC07A000/NetworkPorts"
   },
   "Oem":{
      "Hpe":{
         "@odata.context":"/redfish/v1/$metadata#HpeNetworkAdapter.HpeNetworkAdapter",
         "@odata.type":"#HpeNetworkAdapter.v1_3_0.HpeNetworkAdapter",
         "Actions":{
            "#HpeNetworkAdapter.FlushConfigurationToNVM":{
               "target":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/NetworkAdapters/DC07A000/Actions/Oem/Hpe/HpeNetworkAdapter.FlushConfigurationToNVM",
               "title":"Force a save of current network adapter configuration to non-volatile storage."
            },
            "#NetworkAdapter.FlushConfigurationToNVM":{
               "target":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/NetworkAdapters/DC07A000/Actions/Oem/Hpe/NetworkAdapter.FlushConfigurationToNVM",
               "title":"NOTE: Deprecated, will be removed in a future release. Replaced by HpeNetworkAdapter.FlushConfigurationToNVM. Force a save of current network adapter configuration to non-volatile storage."
            }
         },
         "CLPVersion":"",
         "Controllers":[
            {
               "DeviceLimitationsBitmap":0,
               "EdgeVirtualBridging":{
                  "ChannelDescriptionTLVCapable":true,
                  "ChannelLinkControlTLVCapable":true
               },
               "EmbeddedLLDPFunctions":{
                  "Enabled":true,
                  "Optional":true
               },
               "FunctionTypeLimits":[
                  {
                     "ConstraintDescription":"RES1",
                     "FCoEResourcesConsumed":1,
                     "TotalSharedResourcesAvailable":1,
                     "iSCSIResourcesConsumed":1
                  }
               ],
               "FunctionTypes":[
                  "Ethernet",
                  "iSCSI",
                  "FCoE"
               ],
               "MostRecentConfigurationChangeSource":"None",
               "RDMASupport":[
                  "None"
               ],
               "UnderlyingDataSource":"DCi"
            }
         ],
         "FactoryDefaultsActuationBehavior":"Immediate",
         "PCAVersion":"700749-001",
         "RedfishConfiguration":"Disabled"
      }
   },
   "PartNumber":"0",
   "SKU":"534FLR",
   "SerialNumber":"CN7842V67N"
}
```

###  Power

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Chassis/{ChassisId}/Power` |
|**Description** |This operation retrieves power metrics specific to a server.|
|**Returns** |Information on power consumption and power limiting|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Chassis/{ChassisId}/Power'
```

> **Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#Power.Power",
   "@odata.etag":"W/\"ADB9FA3D\"",
   "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Power",
   "@odata.type":"#Power.v1_3_0.Power",
   "Id":"Power",
   "Name":"PowerMetrics",
   "Oem":{
      "Hpe":{
         "@odata.context":"/redfish/v1/$metadata#HpePowerMetricsExt.HpePowerMetricsExt",
         "@odata.type":"#HpePowerMetricsExt.v2_3_0.HpePowerMetricsExt",
         "BrownoutRecoveryEnabled":true,
         "HasCpuPowerMetering":true,
         "HasDimmPowerMetering":true,
         "HasGpuPowerMetering":false,
         "HasPowerMetering":true,
         "HighEfficiencyMode":"Balanced",
         "Links":{
            "FastPowerMeter":{
               "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Power/FastPowerMeter"
            },
            "FederatedGroupCapping":{
               "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Power/FederatedGroupCapping"
            },
            "PowerMeter":{
               "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Power/PowerMeter"
            },
            "SlowPowerMeter":{
               "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Power/SlowPowerMeter"
            }
         },
         "MinimumSafelyAchievableCap":null,
         "MinimumSafelyAchievableCapValid":false,
         "SNMPPowerThresholdAlert":{
            "DurationInMin":0,
            "ThresholdWatts":0,
            "Trigger":"Disabled"
         }
      }
   },
   "PowerControl":[
      {
         "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Power#PowerControl/0",
         "MemberId":"0",
         "PowerCapacityWatts":1000,
         "PowerConsumedWatts":202,
         "PowerLimit":{
            "LimitException":null,
            "LimitInWatts":null
         },
         "PowerMetrics":{
            "AverageConsumedWatts":211,
            "IntervalInMin":20,
            "MaxConsumedWatts":341,
            "MinConsumedWatts":202
         }
      }
   ],
   "PowerSupplies":[
      {
         "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Power#PowerSupplies/0",
         "FirmwareVersion":"1.00",
         "LastPowerOutputWatts":104,
         "LineInputVoltage":211,
         "LineInputVoltageType":"ACHighLine",
         "Manufacturer":"LTEON",
         "MemberId":"0",
         "Model":"865408-B21",
         "Name":"HpeServerPowerSupply",
         "Oem":{
            "Hpe":{
               "@odata.context":"/redfish/v1/$metadata#HpeServerPowerSupply.HpeServerPowerSupply",
               "@odata.type":"#HpeServerPowerSupply.v2_0_0.HpeServerPowerSupply",
               "AveragePowerOutputWatts":104,
               "BayNumber":1,
               "HotplugCapable":true,
               "MaxPowerOutputWatts":120,
               "Mismatched":false,
               "PowerSupplyStatus":{
                  "State":"Ok"
               },
               "iPDUCapable":false
            }
         },
         "PowerCapacityWatts":500,
         "PowerSupplyType":"AC",
         "SerialNumber":"5WBXK0ELLB96DW",
         "SparePartNumber":"866729-001",
         "Status":{
            "Health":"OK",
            "State":"Enabled"
         }
      },
      {
         "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Power#PowerSupplies/1",
         "FirmwareVersion":"1.00",
         "LastPowerOutputWatts":98,
         "LineInputVoltage":210,
         "LineInputVoltageType":"ACHighLine",
         "Manufacturer":"LTEON",
         "MemberId":"1",
         "Model":"865408-B21",
         "Name":"HpeServerPowerSupply",
         "Oem":{
            "Hpe":{
               "@odata.context":"/redfish/v1/$metadata#HpeServerPowerSupply.HpeServerPowerSupply",
               "@odata.type":"#HpeServerPowerSupply.v2_0_0.HpeServerPowerSupply",
               "AveragePowerOutputWatts":98,
               "BayNumber":2,
               "HotplugCapable":true,
               "MaxPowerOutputWatts":101,
               "Mismatched":false,
               "PowerSupplyStatus":{
                  "State":"Ok"
               },
               "iPDUCapable":false
            }
         },
         "PowerCapacityWatts":500,
         "PowerSupplyType":"AC",
         "SerialNumber":"5WBXK0ELLB96AP",
         "SparePartNumber":"866729-001",
         "Status":{
            "Health":"OK",
            "State":"Enabled"
         }
      }
   ],
   "Redundancy":[
      {
         "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Power#Redundancy/0",
         "MaxNumSupported":2,
         "MemberId":"0",
         "MinNumNeeded":2,
         "Mode":"Failover",
         "Name":"PowerSupply Redundancy Group 1",
         "RedundancySet":[
            {
               "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Power#PowerSupplies/0"
            },
            {
               "@odata.id":"/redfish/v1/Chassis/b1ae6e44-ca60-4b72-87ce-f1c5d59a094d.1/Power#PowerSupplies/1"
            }
         ],
         "Status":{
            "Health":"OK",
            "State":"Enabled"
         }
      }
   ]
}
```

### Creating a rack group

|||
|---------|-------|
|Method | `POST` |
|URI |`/redfish/v1/Chassis`|
|Description |This operation creates a rack group.|
|Returns |<ul><li>`Location` header that contains a link to the created rack group ( in the sample response header)</li><li>JSON schema representing the created rack group<br></li></ul>|
|Response code |On success, `201 Created`|
|Authentication |Yes|

 **curl command**

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

> **Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|ChassisType|String (required)<br> |The type of chassis. The type to be used to create a rack group is RackGroup.<br> |
|Description|String (optional)<br> |Description of this rack group.|
|Links{|Object (required)<br> |Links to the resources that are related to this rack group.|
|ManagedBy [{<br> @odata.id<br> }]<br> }<br> |Array (required)<br> |An array of links to the manager resources that manage this chassis. The manager resource for racks and rack groups is the URP manager. Provide the link to the URP manager.<br> |
|Name|String (required)<br> |Name for this rack group.|


>**Sample response header**

```
Date:Wed,06 Jan 2021 09:37:43 GMT+15m 26s
**Location:/redfish/v1/Chassis/c2459269-011c-58d3-a217-ef914c4c295d**
Content-Length:462 bytes
```

>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/22804541-c439-5d2a-81d5-23d23e0ebe38",
   "@odata.type":"#Chassis.v1_20_0.Chassis",
   "Id":"22804541-c439-5d2a-81d5-23d23e0ebe38",
   "Description":"My RackGroup",
   "Name":"RG2",
   "ChassisType":"RackGroup",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/b44b87c0-00de-4184-ad2b-cdd4da52a805"
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
|**Returns** |<ul><li>`Location` header that contains a link to the created rack ( in the sample response header)</li><li>JSON schema representing the created rack<br></li></ul>|
|**Response code** |On success, `201 Created`|
|**Authentication** |Yes|

 **curl command**

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

> **Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|ChassisType|String (required)<br> |The type of chassis. The type to be used to create a rack is Rack.<br> |
|Description|String (optional)<br> |Description of this rack.|
|Links{|Object (required)<br> |Links to the resources that are related to this rack.|
|ManagedBy [{<br> @odata.id<br> }]<br> |Array (required)<br> |An array of links to the manager resources that manage this chassis. The manager resource for racks and rack groups is the URP manager. Provide the link to the URP manager.<br> |
|ContainedBy [{<br> @odata.id<br> }]<br> }<br> |Array (required)<br> |An array of links to the rack groups for containing this rack.|
|Name|String (required)<br> |Name for this rack group.|


>**Sample response header**

```
Date:Wed,06 Jan 2021 09:37:43 GMT+15m 26s
**Location:/redfish/v1/Chassis/b6766cb7-5721-5077-ae0e-3bf3683ad6e2**
Content-Length:462 bytes
```

>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "@odata.type":"#Chassis.v1_20_0.Chassis",
   "Id":"b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "Description":"rack no 1",
   "Name":"RACK#1",
   "ChassisType":"Rack",
   "Links":{
      "ComputerSystems":[
         
      ],
      "ManagedBy":[
         {
            "@odata.id":"/redfish/v1/Managers/b44b87c0-00de-4184-ad2b-cdd4da52a805"
         }
      ],
      "ContainedBy":[
         {
            "@odata.id":"/redfish/v1/Chassis/22804541-c439-5d2a-81d5-23d23e0ebe38"
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
|**Returns** |JSON schema for the modified rack having links to the attached chassis|
|**Response code** |On success, `200 Ok`|
|**Authentication** |Yes|

 **curl command**

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
        "@odata.id": "/redfish/v1/Chassis/46db63a9-2dcb-43b3-bdf2-54ce9c42e9d9.1"
      }
    ]
  }
}
```

> **Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Links{|Object (required)<br> |Links to the resources that are related to this rack|
|Contains [{<br> @odata.id<br> }]<br> }<br> |Array (required)<br> |An array of links to the computer system chassis resources to be attached to this rack|


>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "@odata.type":"#Chassis.v1_20_0.Chassis",
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
            "@odata.id":"/redfish/v1/Chassis/4159c951-d0d0-4263-858b-0294f5be6377.1"
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
|**Returns** |JSON schema representing the modified rack|
|**Response code** |On success, `200 Ok`|
|**Authentication** |Yes|

 **curl command**

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

> **Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Links{|Object (required)<br> |Links to the resources that are related to this rack.|
|Contains [{<br> @odata.id<br> }]<br> }<br> |Array (required)<br> |An array of links to the computer system chassis resources to be attached to this rack. To detach chassis from this rack, provide an empty array as value.|


>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#Chassis.Chassis",
   "@odata.id":"/redfish/v1/Chassis/b6766cb7-5721-5077-ae0e-3bf3683ad6e2",
   "@odata.type":"#Chassis.v1_20_0.Chassis",
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
|**Description** |This operation deletes a specific rack.<br>**IMPORTANT:** If you try to delete a non-empty rack, you will receive an HTTP `409 Conflict` error. Ensure to detach the chassis attached to a rack before deleting the rack.<br>|
|**Response code** |On success, `204 No Content`|
|**Authentication** |Yes|

 **curl command**

```
curl -i DELETE \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   'https://{odim_host}:{port}/redfish/v1/Chassis/{rackId}'
```

### Deleting a rack group

|||
|---------|-------|
|**Method**| `DELETE` |
|**URI**|`/redfish/v1/Chassis/{rackGroupId}`|
|**Description**|This operation deletes a specific rack group.<br>**IMPORTANT:** If you try to delete a non-empty rack group, you will receive an HTTP `409 Conflict` error. Ensure to remove all the racks contained in a rack group before deleting the rack group.|
|**Response code**|On success, `204 No Content`|
|**Authentication**|Yes|

 **curl command**

```
curl -i DELETE \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   'https://{odim_host}:{port}/redfish/v1/Chassis/{rackGroupId}'
```



##   Searching the inventory

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Systems?$filter={searchKeys}%20{conditionKeys}%20{value}` |
|**Description** | Use this endpoint to search servers based on filters - combination of a keyword, condition, and a value.<br> Two ore more filters can be combined in a single request with the help of logical operands.<br>**NOTE:** Only a user with `Login` privilege can perform this operation. |
|**Returns** |Server endpoints based on the specified filter|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Systems?$filter={searchKeys}%20{conditionKeys}%20{value/regular_expression}%20{logicalOperand}%20{searchKeys}%20{conditionKeys}%20{value}'

```

> **Sample usage** 

```
curl -i GET \
      -H "X-Auth-Token:{X-Auth-Token}" \
    'http://{odimra_host}:{port}/redfish/v1/Systems?$filter=MemorySummary/TotalSystemMemoryGiB%20eq%20384'
```

### Request URI parameters

-  `{searchkeys}` refers to `ComputerSystem` parameters. Following are the allowed search keys:

   - `ProcessorSummary/Count` 
   
   -   `ProcessorSummary/Model` 
   
   -   `ProcessorSummary/sockets` 
   
   -   `SystemType` 
   
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

        $filter=TotalSystemMemoryGiB%20eq%20**384**
        
        $filter=ProcessorSummary/Model%20eq%20**int\***
        
        $filter=Storage/Drives/Type%20eq%20HDD

-  `{logicalOperands}` refers to the logical operands that are used to combine two or more filters in a request. Allowed logical operands are `and`, `or`, and `not`.


#### **Sample filters**

- `$filter=TotalSystemMemoryGiB%20eq%20384`
  This filter searches a server having total physical memory of 384 GB.

- `$filter=ProcessorSummary/Model%20eq%20int*`
  This filter searches a server whose processor model name starts with `int` and ends with any combination of letters, numbers and/or special characters.

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
         "@odata.id":"/redfish/v1/Systems/7ff3bd97-c41c-5de0-937d-85d390691b73.1"
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
|**Description** |This action shuts down, powers up, and restarts a specific system.<br>**NOTE:** To reset an aggregate of systems, use the following URI:<br>`/redfish/v1/AggregationService/Actions/AggregationService.Reset` <br> See *[Resetting servers](#resetting-servers)*.|
|**Returns** |A Redfish task in the response header and you receive a link to the task monitor associated with it. To know the progress of this operation, perform an `HTTP GET` on the task monitor (until the task is complete).|
|**Response code** | `202 Accepted`. On successful completion, `200 OK`. |
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

> **Request parameters**

See *[Resetting Servers](#resetting-servers)* to know about `ResetType.` 

>**Sample response body**

```
{
    "error": {
        "@Message.ExtendedInfo": [{
            "MessageId": "Base.1.13.Success"
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
|**Description** |This action changes the boot order of a specific system to default settings.<br>**NOTE:**<br> To change the boot order of an aggregate of systems, use the following URI:<br> `/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder` <br> See *[Changing the Boot Order of Servers to Default Settings](#changing-the-boot-order-of-servers-to-default-settings)*.|
|**Returns** |Message Id of the actual message in the JSON response body. To get the complete message, look up the specified registry file. Registry file name can be obtained by concatenating `RegistryPrefix` and version number present in the Message id. <br />**Example registry file name**: Base.1.4. See *[Message Registries](#message-registries)*.|
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
			"MessageId": "Base.1.13.0.Success"
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
|**Description** |This action changes BIOS configuration.<br>**NOTE:** Any change in BIOS configuration is reflected only after the system reset. To see the change, *[reset the computer system](#resetting-a-computer-system)*.|
|**Returns** |Message Id of the actual message in the JSON response body. To get the complete message, look up the specified registry file. Registry file name can be obtained by concatenating `RegistryPrefix` and version number present in the Message ID. See *[Message registries](#message-registries)*. <br />For example:`MessageId` in the sample response body is `iLO.2.8.SystemResetRequired`. The registry to look up is `iLO.2.8`.<br> |
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

> **Request parameters**

`Attributes` are the list of BIOS attributes specific to the manufacturer or provider. To get a full list of attributes, perform `GET` on `https://{odimra_host}:{port}/redfish/v1/Systems/1/Bios/Settings`. 

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
|**Description** |This action changes the boot settings of a specific system such as boot source override target, boot order, and more.<br>**IMPORTANT**<br><ul><li>Ensure that the system is powered off before changing the boot order.</li><li>Power on the system once the operation is successful. Changes are seen in the system only after a successful reset.</li></ul><br> To know how to power off, power on, or restart a system, see *[Resetting a computer system](#resetting-a-computer-system)*.|
|**Returns** |Message Id of the actual message in the JSON response body. To get the complete message, look up the specified registry file. Registry file name can be obtained by concatenating `RegistryPrefix` and version number present in the Message Id. See *[Message Registries](#message-registries)*. <br />For example,`MessageId` in the sample response body is `Base.1.13.0.Success`. The registry to look up is `Base.1.13.0`.<br> |
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

> **Request parameters**

To get a full list of boot attributes that you can update, perform `GET` on:


`https://{odimra_host}:{port}/redfish/v1/Systems/{ComputerSystemId}`.


Check attributes under `Boot` in the JSON response. Some of the attributes include:

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

**NOTE:** If you attempt to update `BootSourceOverrideTarget` to `UefiTarget`, when `UefiTargetBootSourceOverride` is set to `None`, you encounter an HTTP `400 Bad Request` error. Update `UefiTargetBootSourceOverride` before setting `BootSourceOverrideTarget` to `UefiTarget`.

>**Sample response body**

```
{ 
   "error":{ 
      "@Message.ExtendedInfo":[ 
         { 
            "MessageId":"Base.1.13.0.Success"
         }
      ],
      "code":"iLO.0.10.ExtendedInfo",
      "message":"See @Message.ExtendedInfo for more information."
   }
```




# Managers

Resource Aggregator for ODIM exposes APIs to retrieve information about managers that include:

-   Resource Aggregator for ODIM

-   BMCs

-   Enclosure Managers

-   Management Controller

-   Other subsystems like plugins


**Supported endpoints**


|||
|-------|--------------------|
|/redfish/v1/Managers|`GET`|
|/redfish/v1/Managers/{managerId}|`GET`|
|/redfish/v1/Managers/{managerId}/EthernetInterfaces|`GET`|
|/redfish/v1/Managers/{managerId}/HostInterfaces|`GET`|
|/redfish/v1/Managers/{managerId}/LogServices|`GET`|
|/redfish/v1/Managers/{managerId}/NetworkProtocol|`GET`|
|/redfish/v1/Managers/{managerId}/VirtualMedia|`GET`|




##  Collection of managers

|||
|---------|-------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Managers` |
|**Description** |A collection of managers.|
|**Returns** |Links to the manager instances. This collection includes a manager for Resource Aggregator for ODIM and other managers.|
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
         "@odata.id":"/redfish/v1/Managers/141cbba9-1e99-4272-b855-1781730bfe1c.1"
      },
      {
         "@odata.id":"/redfish/v1/Managers/536cee48-84b2-43dd-b6e2-2459ac0eeac6"
      },
      {
         "@odata.id":"/redfish/v1/Managers/0e778112-4684-433d-9998-ca6f399c031f.1"
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
|**Returns** |Information about a specific management control system or a plugin or Resource Aggregator for ODIM. In the JSON schema representing a system (BMC) manager, you can view links to the managers for:<br /><ul><li>EthernetInterfaces: `/redfish/v1/Managers/{managerId}/EthernetInterfaces`</li><br /><li>HostInterfaces: `/redfish/v1/Managers/{managerId}/HostInterfaces` </li><li><br />LogServices: `/redfish/v1/Managers/{managerId}/LogServices` </li><br /><li>NetworkProtocol: `/redfish/v1/Managers/{managerId}/NetworkProtocol`<br /> **NOTE**: To know more about each manager, perform HTTP `GET` on these links.</li></ul>|
|**Response code** | `200 OK` |
|**Authentication** |Yes|


>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Managers/{managerId}'

```

>**Sample response body for a system (BMC) manager** 

```
{
    "@odata.context": "/redfish/v1/$metadata#Manager.Manager",
    "@odata.etag": "W/\"887BAFE1\",
    "@odata.id: "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1",
    "@odata.type": "#Manager.v1_15_0.Manager",
    "Actions": {
        "#Manager.Reset": {
            ResetType@Redfish.AllowableValues: [
                "ForceRestart",
               "GracefulRestart"
            ],
            "target": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/Actions/Manager.Reset"
        }
    },
    "CommandShell": {
        "ConnectTypesSupported": [
            "SSH",
            "Oem"
        ],
        "MaxConcurrentSessions": 9,
        "ServiceEnabled": true
    },
    "DateTime": "2022-08-22T05:23:29Z",
    "DateTimeLocalOffset": "+00:00",
    "Description": "BMC Manager",
    "EthernetInterfaces": {
        "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/EthernetInterfaces"
    },
    "FirmwareVersion": "iLO 5 v2.70",
    "GraphicalConsole": {
        "ConnectTypesSupported": [
            "KVMIP"
        ],
        "MaxConcurrentSessions": 10,
        "ServiceEnabled": true
    },
    "HostInterfaces": {
        "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/HostInterfaces"
    },
    "Id": "1",
    "Links": {
        "ManagerForChassis": [
            {
                "@odata.id": "/redfish/v1/Chassis/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1"
            }
        ],
        "ManagerForServers": [
            {
                "@odata.id": "/redfish/v1/Systems/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1"
            }
        ],
        "ManagerInChassis": {
            "@odata.id": "/redfish/v1/Chassis/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1"
        }
    },
    "LogServices": {
        "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/LogServices"
    },
    "ManagerType": "BMC",
    "Model": "iLO 5",
    "Name": "Manager",
    "NetworkProtocol": {
        "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/NetworkProtocol"
    },
    "Oem": {
        "Hpe": {
            "@odata.context": "/redfish/v1/$metadata#HpeiLO.HpeiLO",
            "@odata.type": "#HpeiLO.v2_8_1.HpeiLO",
            "Actions": {
                "#HpeiLO.ClearHotKeys": {
                    "target": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/Actions/Oem/Hpe/HpeiLO.ClearHotKeys"
                },
                "#HpeiLO.ClearRestApiState": {
                    "target": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/Actions/Oem/Hpe/HpeiLO.ClearRestApiState"
                },
                "#HpeiLO.DisableCloudConnect": {
                    "target": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/Actions/Oem/Hpe/HpeiLO.DisableCloudConnect"
                },
                "#HpeiLO.DisableiLOFunctionality": {
                    "target": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/Actions/Oem/Hpe/HpeiLO.DisableiLOFunctionality"
                },
                "#HpeiLO.EnableCloudConnect": {
                    "target": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/Actions/Oem/Hpe/HpeiLO.EnableCloudConnect"
                },
                "#HpeiLO.RequestFirmwareAndOsRecovery": {
                    "target": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/Actions/Oem/Hpe/HpeiLO.RequestFirmwareAndOsRecovery"
                },
                "#HpeiLO.ResetToFactoryDefaults": {
                    ResetType@Redfish.AllowableValues: [
                        "Default"
                    ],
                    "target": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/Actions/Oem/Hpe/HpeiLO.ResetToFactoryDefaults"
                }
            },
            "ClearRestApiStatus": "DataPresent",
            "CloudConnect": {
                "ActivationKey": "",
                "CloudConnectStatus": "NotEnabled"
            },
            "ConfigurationLimitations": "None",
            "ConfigurationSettings": "Current",
            "FederationConfig": {
                "IPv6MulticastScope": "Site",
                "MulticastAnnouncementInterval": 600,
                "MulticastDiscovery": "Enabled",
                "MulticastTimeToLive": 5,
                "iLOFederationManagement": "Enabled"
            },
            "Firmware": {
                "Current": {
                    "Date": "May 16 2022",
                    "DebugBuild": false,
                    "MajorVersion": 2,
                    "MinorVersion": 70,
                    "VersionString": "iLO 5 v2.70"
                }
            },
            "FrontPanelUSB": {
                "State": "Ready"
            },
            "IdleConnectionTimeoutMinutes": 120,
            "IntegratedRemoteConsole": {
                "HotKeys": [
                    {
                        "KeySequence": [
                            "NONE",
                            "NONE",
                            "NONE",
                            "NONE",
                            "NONE"
                        ],
                        "Name": "Ctrl-T"
                    },
                    {
                        "KeySequence": [
                            "NONE",
                            "NONE",
                            "NONE",
                            "NONE",
                            "NONE"
                        ],
                        "Name": "Ctrl-U"
                    },
                    {
                        "KeySequence": [
                            "NONE",
                            "NONE",
                            "NONE",
                            "NONE",
                            "NONE"
                        ],
                        "Name": "Ctrl-V"
                    },
                    {
                        "KeySequence": [
                            "NONE",
                            "NONE",
                            "NONE",
                            "NONE",
                            "NONE"
                        ],
                        "Name": "Ctrl-W"
                    },
                    {
                        "KeySequence": [
                            "NONE",
                            "NONE",
                            "NONE",
                            "NONE",
                            "NONE"
                        ],
                        "Name": "Ctrl-X"
                    },
                    {
                        "KeySequence": [
                            "NONE",
                            "NONE",
                            "NONE",
                            "NONE",
                            "NONE"
                        ],
                        "Name": "Ctrl-Y"
                    }
                ],
                "LockKey": {
                    "CustomKeySequence": [
                        "NONE",
                        "NONE",
                        "NONE",
                        "NONE",
                        "NONE"
                    ],
                    "LockOption": "Disabled"
                },
                "TrustedCertificateRequired": false
            },
            "License": {
                "LicenseKey": "XXXXX-XXXXX-XXXXX-XXXXX-7BK6M",
                "LicenseString": "iLO Advanced limited-distribution test",
                "LicenseType": "Internal"
            },
            "Links": {
                "ActiveHealthSystem": {
                    "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/ActiveHealthSystem"
                },
                "BackupRestoreService": {
                    "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/BackupRestoreService"
                },
                "DateTimeService": {
                    "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/DateTime"
                },
                "EmbeddedMediaService": {
                    "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/EmbeddedMedia"
                },
                "FederationDispatch": {
                    "extref": "/dispatch"
                },
                "FederationGroups": {
                    "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/FederationGroups"
                },
                "FederationPeers": {
                    "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/FederationPeers"
                },
                "GUIService": {
                    "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/GUIService"
                },
                "LicenseService": {
                    "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/LicenseService"
                },
                "RemoteSupport": {
                    "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/RemoteSupportService"
                },
                "SNMPService": {
                    "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/SnmpService"
                },
                "SecurityService": {
                    "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/SecurityService"
                },
                "Thumbnail": {
                    "extref": "/images/thumbnail.bmp"
                },
                "VSPLogLocation": {
                    "extref": "/sol.log.gz"
                }
            },
            "PersistentMouseKeyboardEnabled": false,
            "PhysicalMonitorHealthStatusEnabled": true,
            "RIBCLEnabled": true,
            "RemoteConsoleThumbnailEnabled": true,
            "RequireHostAuthentication": false,
            "RequiredLoginForiLORBSU": false,
            "SerialCLISpeed": 9600,
            "SerialCLIStatus": "EnabledAuthReq",
            "SerialCLIUART": "Present",
            "VSPDlLoggingEnabled": false,
            "VSPLogDownloadEnabled": false,
            "VideoPresenceDetectOverride": true,
            "VideoPresenceDetectOverrideSupported": true,
            "VirtualNICEnabled": false,
            "WebGuiEnabled": true,
            "iLOFunctionalityEnabled": true,
            "iLOFunctionalityRequired": false,
            "iLOIPduringPOSTEnabled": true,
            "iLORBSUEnabled": true,
            "iLOSelfTestResults": [
                {
                    "Notes": "",
                    "SelfTestName": "NVRAMData",
                    "Status": "OK"
                },
                {
                    "Notes": "Controller firmware revision  2.11.00  ",
                    "SelfTestName": "EmbeddedFlash",
                    "Status": "OK"
                },
                {
                    "Notes": "",
                    "SelfTestName": "EEPROM",
                    "Status": "OK"
                },
                {
                    "Notes": "",
                    "SelfTestName": "HostRom",
                    "Status": "OK"
                },
                {
                    "Notes": "",
                    "SelfTestName": "SupportedHost",
                    "Status": "OK"
                },
                {
                    "Notes": "Version 1.0.8",
                    "SelfTestName": "PowerManagementController",
                    "Status": "Informational"
                },
                {
                    "Notes": "ProLiant DL380 Gen10 Plus System Programmable Logic Device 0x15",
                    "SelfTestName": "CPLDPAL0",
                    "Status": "Informational"
                },
                {
                    "Notes": "",
                    "SelfTestName": "ASICFuses",
                    "Status": "OK"
                }
            ],
            "iLOServicePort": {
                "MassStorageAuthenticationRequired": false,
                "USBEthernetAdaptersEnabled": true,
                "USBFlashDriveEnabled": true,
                "iLOServicePortEnabled": true
            }
        }
    },
    "PowerState": "On",
    "RemoteAccountService": {
        "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/RemoteAccountService"
    },
    "SerialConsole": {
        "ConnectTypesSupported": [
            "SSH",
            "IPMI",
            "Oem"
        ],
        "MaxConcurrentSessions": 13,
        "ServiceEnabled": true
    },
    "SerialInterfaces": {
        "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/SerialInterfaces"
    },
    "Status": {
        "Health": "OK",
        "State": "Enabled"
    },
    "UUID": "2a4f0469-b023-58f4-9f89-bfd1f67ed70e",
    "VirtualMedia": {
        "@odata.id": "/redfish/v1/Managers/7859c05c-8ed4-4f2d-bef5-ce8b7d2528fc.1/VirtualMedia"
    }
}
```

>**Sample response body for Resource Aggregator for ODIM manager**

```
{
   "@odata.context":"/redfish/v1/$metadata#Manager.Manager",
   "@odata.id":"/redfish/v1/Managers/1df3248f-5ddd-4b62-868d-74f33c4a89d0",
   "@odata.type":"#Manager.v1_15_0.Manager",
   "Name":"odimra",
   "ManagerType":"Service",
   "Id":"1df3248f-5ddd-4b62-868d-74f33c4a89d0",
   "UUID":"1df3248f-5ddd-4b62-868d-74f33c4a89d0",
   "FirmwareVersion":"1.0",
   "Status":{
      "State":"Enabled",
      "Health":"OK"
   },
   "LogServices":{
      "@odata.id":"/redfish/v1/Managers/1df3248f-5ddd-4b62-868d-74f33c4a89d0/LogServices"
   },
   "Links":{
      "ManagerForChassis":[
         {
            "@odata.id":"/redfish/v1/Chassis/573bbf22-6b28-48ce-9e22-2a55c9d1adde.1"
         }
      ],
      "ManagerForServers":[
         {
            "@odata.id":"/redfish/v1/Systems/573bbf22-6b28-48ce-9e22-2a55c9d1adde.1"
         }
      ],
      "ManagerForManagers":[
         {
            "@odata.id":"/redfish/v1/Managers/573bbf22-6b28-48ce-9e22-2a55c9d1adde.1"
         },
         {
            "@odata.id":"/redfish/v1/Managers/386710f8-3a38-4938-a986-5f1048f487fd"
         }
      ]
   },
   "DateTime":"2022-04-07T10:27:40Z",
   "Model":"ODIMRA 1.0",
   "PowerState":"On",
   "SerialConsole":{
      
   },
   "Description":"Odimra Manager",
   "DateTimeLocalOffset":"+00:00"
}
```

>**Sample response body for a plugin manager**

```
{
   "@odata.context":"/redfish/v1/$metadata#Manager.Manager",
   "@odata.etag":"W/\"AA6D42B0\"",
   "@odata.id":"/redfish/v1/Managers/ac04517b-b582-4501-b1a9-7158149cda10",
   "@odata.type":"#Manager.v1_15_0.Manager",
   "DateTime":"2022-02-22 09:52:43.651476316 +0000 UTC",
   "DateTimeLocalOffset":"+00:00",
   "Description":"Plugin Manager",
   "FirmwareVersion":"v1.0.0",
   "Id":"ac04517b-b582-4501-b1a9-7158149cda10",
   "Links":{
      "ManagerForChassis":[
         {
            "@odata.id":"/redfish/v1/Chassis/b91d2658-0a5f-4478-bd11-3e494687afc5.1"
         }
      ],
      "ManagerForServers":[
         {
            "@odata.id":"/redfish/v1/Systems/b91d2658-0a5f-4478-bd11-3e494687afc5.1"
         }
      ]
   },
   "LogServices":{
      "@odata.id":"/redfish/v1/Managers/ac04517b-b582-4501-b1a9-7158149cda10/LogServices"
   },
   "ManagerType":"Service",
   "Name":"ILO",
   "SerialConsole":{
      
   },
   "Status":{
      "Health":"OK",
      "State":"Enabled"
   },
   "UUID":"ac04517b-b582-4501-b1a9-7158149cda10"
}
```



## VirtualMedia

The `VirtualMedia` resource enables you to connect remote storage media (such as CD-ROM, USB mass storage, ISO image, and floppy disk) to a target server on a network. The target server can access the remote media, read from and write to it as if it were physically connected to the server USB port.

Resource Aggregator for ODIM exposes Redfish `VirtualMedia` APIs to connect the remote storage media to your servers.

**Supported APIs**:

| API URI                                                      | Supported operations | Required privileges   |
| ------------------------------------------------------------ | -------------------- | --------------------- |
| /redfish/v1/Managers/{ManagerId}/VirtualMedia                | `GET`                | `Login`               |
| /redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID} | `GET`                | `Login`               |
| /redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}/Actions/VirtualMedia.InsertMedia | `POST`               | `ConfigureComponents` |
| /redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}/Actions/VirtualMedia.EjectMedia | `POST`               | `ConfigureComponents` |

### Viewing the VirtualMedia collection

| **Method**         | `GET`                                                        |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Managers/{ManagerId}/VirtualMedia`              |
| **Description**    | This operation lists all virtualmedia collections available in Resource Aggregator for ODIM. |
| **Returns**        | A list of links to all the available virtualmedia collections |
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
  "@odata.etag": "W/\"570254F2\"",
   "@odata.id":"/redfish/v1/Managers/1/VirtualMedia/",
   "@odata.type":"#VirtualMediaCollection.VirtualMediaCollection",
   "Description":"Virtual Media Services Settings",
   "Name":"Virtual Media Services",
   "Members":[
      {
         "@odata.id":"/redfish/v1/Managers/1/VirtualMedia/1"
      },
      {
         "@odata.id":"/redfish/v1/Managers/1/VirtualMedia/2"
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
| <strong>Returns</strong>        | JSON schema representing this virtualmedia instance          |
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
    "@odata.context": "/redfish/v1/$metadata#VirtualMedia.VirtualMedia",
    "@odata.etag": "W/\"3B0F66BA\"",
    "@odata.id": "/redfish/v1/Managers/c55ea6a6-a501-44a5-b159-3579c67cb81e.1/VirtualMedia/1",
    "@odata.type": "#VirtualMedia.v1_2_0.VirtualMedia",
    "Actions": {
        "#VirtualMedia.EjectMedia": {
            "target": "/redfish/v1/Managers/c55ea6a6-a501-44a5-b159-3579c67cb81e.1/VirtualMedia/1/Actions/VirtualMedia.EjectMedia"
        },
        "#VirtualMedia.InsertMedia": {
            "target": "/redfish/v1/Managers/c55ea6a6-a501-44a5-b159-3579c67cb81e.1/VirtualMedia/1/Actions/VirtualMedia.InsertMedia"
        }
    },
    "ConnectedVia": "NotConnected",
    "Description": "Virtual Removable Media",
    "Id": "1",
    "Image": "",
    "Inserted": false,
    "MediaTypes": [
        "Floppy",
        "USBStick"
    ],
    "Name": "VirtualMedia",
    "Oem": {
        "Hpe": {
            "@odata.context": "/redfish/v1/$metadata#HpeiLOVirtualMedia.HpeiLOVirtualMedia",
            "@odata.type": "#HpeiLOVirtualMedia.v2_2_0.HpeiLOVirtualMedia",
            "Actions": {
                "#HpeiLOVirtualMedia.EjectVirtualMedia": {
                    "target": "/redfish/v1/Managers/c55ea6a6-a501-44a5-b159-3579c67cb81e.1/VirtualMedia/1/Actions/Oem/Hpe/HpeiLOVirtualMedia.EjectVirtualMedia"
                },
                "#HpeiLOVirtualMedia.InsertVirtualMedia": {
                    "target": "/redfish/v1/Managers/c55ea6a6-a501-44a5-b159-3579c67cb81e.1/VirtualMedia/1/Actions/Oem/Hpe/HpeiLOVirtualMedia.InsertVirtualMedia"
                }
            }
        }
    },
    "WriteProtected": true
}
```

### Inserting VirtualMedia

| **Method**         | `POST`                                                       |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}/Actions/VirtualMedia.InsertMedia` |
| **Description**    | This operation inserts the virtual media on to the manager.  |
| **Returns**        | A message stating the virtual media insertion was successful |
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
    "Error": {
        "@Message.ExtendedInfo": [
            {
                "Message": "Successfully performed virtual media actions",
                "MessageArgs": [
               ],
                "MessageId": "Base.1.11.0.Success"
            }
        ],
        "Code": "Base.1.11.0.Success",
        "Message": "See @Message.ExtendedInfo for more information."
    }
} 

```

### Ejecting VirtualMedia

| **Method**         | `POST`                                                       |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/Managers/{ManagerId}/VirtualMedia/{VirtualMediaID}/Actions/VirtualMedia.EjectMedia` |
| **Description**    | This operation ejects the virtual media from the manager.    |
| **Returns**        | A message stating the virtual media ejection was successful  |
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
    "Error": {
        "@Message.ExtendedInfo": [
            {
                "Message": "Successfully performed virtual media actions",
                "MessageArgs": [
                ],
                "MessageId": "Base.1.11.0.Success"
            }
        ],
        "Code": "Base.1.11.0.Success",
        "Message": "See @Message.ExtendedInfo for more information."
    }
}
```

> **Request parameters**

| Parameter      | Type               | Description           |
| -------------- | ------------------ | --------------------- |
| Image          | String (Required)  | Image path            |
| Inserted       | Boolean (Optional) | Default value is true |
| WriteProtected | Boolean (Optional) | Default value is true |

## Remote BMC accounts and roles

Resource Aggregator for ODIM exposes `RemoteAccountService` APIs to manage BMC accounts and roles. 

**Supported APIs**:

| API URI                                                      | Supported operations     | Required privileges            |
| ------------------------------------------------------------ | ------------------------ | ------------------------------ |
| /redfish/v1/Managers/{ManagerId}/RemoteAccountService        | `GET`                    | `Login`                        |
| /redfish/v1/Managers/{ManagerId}/RemoteAccountService/Accounts | `GET`, `POST`            | `Login`, `ConfigureComponents` |
| /redfish/v1/Managers/{ManagerId}/RemoteAccountService/Accounts/{AccountId} | `GET`, `PATCH`, `DELETE` | `Login`, `ConfigureComponents` |
| /redfish/v1/Managers/{ManagerId}/RemoteAccountService/Roles  | `GET`                    | `Login`                        |
| /redfish/v1/Managers/{ManagerId}/RemoteAccountService/Roles/{Roleid} | `GET`                    | `Login`                        |

### Viewing the RemoteAccountService root

|                                 |                                                              |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>Method</strong>         | `GET`                                                        |
| <strong>URI</strong>            | `/redfish/v1/Managers/{ManagerID}/RemoteAccountService`      |
| <strong>Description</strong>    | This operation retrieves JSON schema representing the Redfish `RemoteAccountService` root. |
| <strong>Returns</strong>        | The properties common to all remote BMC accounts and links to the collections of BMC accounts and roles |
| <strong>Response code</strong>  | On success, `200 Ok`                                         |
| <strong>Authentication</strong> | Yes                                                          |

> **curl command**

```
curl -i GET \
-H "X-Auth-Token:{X-Auth-Token}" \
'https://{odim_host}:{port}/redfish/v1/Managers/{ManagerID}/RemoteAccountService'
```


> **Sample response header**

```
Allow: GET
Cache-Control: no-cache, no-store, must-revalidate
Connection: keep-alive
Content-Type: application/json; charset=utf-8
Odata-Version: 4.0
X-Content-Type-Options: nosniff
X-Frame-Options: sameorigin
Date: Tue, 12 Apr 2022 06:15:41 GMT-1s
Transfer-Encoding: chunked
```

> **Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#AccountService.AccountService",
   "@odata.etag":"W/\"8F1B1B4B\"",
   "@odata.id":"/redfish/v1/Managers/bdcc9c30-d062-4239-9a1a-3dc87b4913c7.1/RemoteAccountService",
   "@odata.type":"#AccountService.v1_5_0.AccountService",
   "Id":"AccountService",
   "Name":"Account Service",
   "Description":"iLO User Accounts",
   "Status":{
      
   },
   "Accounts":{
      "@odata.id":"/redfish/v1/Managers/bdcc9c30-d062-4239-9a1a-3dc87b4913c7.1/RemoteAccountService/Accounts"
   },
   "Roles":{
      "@odata.id":"/redfish/v1/Managers/bdcc9c30-d062-4239-9a1a-3dc87b4913c7.1/RemoteAccountService/Roles"
   },
   "MinPasswordLength":8,
   "LocalAccountAuth":"Enabled"
}
```

### Collection of BMC user accounts

|                                 |                                                              |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>Method</strong>         | `GET`                                                        |
| <strong>URI</strong>            | `/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Accounts` |
| <strong>Description</strong>    | A collection of BMC user accounts                            |
| <strong>Returns</strong>        | Links to the BMC user account instances                      |
| <strong>Response code</strong>  | On success, `200 Ok`                                         |
| <strong>Authentication</strong> | Yes                                                          |

> **curl command**

```
curl -i GET \
-H "X-Auth-Token:{X-Auth-Token}" \
'https://{odim_host}:{port}/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Accounts'
```

> **Sample response body**

```
{
"@odata.context":"/redfish/v1/$metadata#ManagerAccountCollection.ManagerAccountCollection",
   "@odata.etag":"W/\"21C260DB\"",
   "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-a3b8293774ba.1/RemoteAccountService/Accounts",
   "@odata.type":"#ManagerAccountCollection.ManagerAccountCollection",
   "Description":"iLO User Accounts",
   "Members":[
      {
         "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-a3b8293774ba.1/RemoteAccountService/Accounts/1"
      },
      {
         "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-a3b8293774ba.1/RemoteAccountService/Accounts/8"
      },
      {
         "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-a3b8293774ba.1/RemoteAccountService/Accounts/7"
      },
      {
         "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-a3b8293774ba.1/RemoteAccountService/Accounts/5"
      },
      {
         "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-a3b8293774ba.1/RemoteAccountService/Accounts/6"
      },
      {
         "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-a3b8293774ba.1/RemoteAccountService/Accounts/2"
      }
   ],
   "Members@odata.count":6,
   "Name":"Accounts"
}
```

### Single BMC user account

|                                 |                                                              |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>Method</strong>         | `GET`                                                        |
| <strong>URI</strong>            | `/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Accounts/{AccountID}` |
| <strong>Description</strong>    | This operation retrieves information about a single BMC user account. |
| <strong>Returns</strong>        | JSON schema representing this user account                   |
| <strong>Response code</strong>  | On success, `200 Ok`                                         |
| <strong>Authentication</strong> | Yes                                                          |

> **curl command**

```
curl -i GET \
-H "X-Auth-Token:{X-Auth-Token}" \
'https://{odim_host}:{port}/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Accounts/{AccountID}'
```

> **Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
   "@odata.etag":"W/\"226E6C7B\"",
   "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-a3b8293774ba.1/RemoteAccountService/Accounts/1",
   "@odata.type":"#ManagerAccount.v1_3_0.ManagerAccount",
   "Id":"1",
   "Name":"User Account",
   "Description":"iLO User Account",
   "UserName":"Administrator",
   "RoleId":"Administrator",
   "Links":{
      "Role":{
         "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-a3b8293774ba.1/RemoteAccountService/Roles/Administrator"
      }
   },
}
```

### Creating a BMC account

|                                 |                                                              |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>Method</strong>         | `POST`                                                       |
| <strong>URI</strong>            | `/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Accounts` |
| <strong>Description</strong>    | This operation creates a BMC user account.                   |
| <strong>Returns</strong>        | JSON schema representing the created user account            |
| <strong>Response code</strong>  | On success, `201 Created`                                    |
| <strong>Authentication</strong> | Yes                                                          |

> **curl command**

```
curl -i POST \
 -H "X-Auth-Token:{X-Auth-Token}" \
 -d \
'{
  "UserName":"{username}",
  "Password":"{password}",
  "RoleId":"Administrator"
}' \
'https://{odim_host}:{port}/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Accounts
```

> **Sample request body**

```
{
  "UserName":"{username}",
  "Password":"{password}",
  "RoleId":"Administrator"
}
```

> **Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
   "@odata.etag":"W/\"A2973884\"",
   "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-a3b8293774ba.1/RemoteAccountService/Accounts/13",
   "@odata.type":"#ManagerAccount.v1_3_0.ManagerAccount",
   "Id":"13",
   "Name":"User Account",
   "Description":"iLO User Account",
   "UserName":"{username}",
   "RoleId":"{roleId}",
   "Links":{
      "Role":{
         "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-a3b8293774ba.1/RemoteAccountService/Roles/Administrator"
      }
   },
}
```

### Updating a BMC account

|                                 |                                                              |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>Method</strong>         | `PATCH`                                                      |
| <strong>URI</strong>            | `/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Accounts/{AccountID}` |
| <strong>Description</strong>    | This operation updates a BMC user account.                   |
| <strong>Returns</strong>        | JSON schema representing the updated user account            |
| <strong>Response code</strong>  | On success, `200 Ok`                                         |
| <strong>Authentication</strong> | Yes                                                          |

> **curl command**

```
curl -i PATCH \
 -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 -d \
'{
  "Password":"{password}",
  "RoleId":"{roleId}"
}' \
'https://{odim_host}:{port}/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Accounts/{AccountID}
```


> **Sample request body**

```
{
"RoleId":"{roleId}",
"Password": "{password}"
}
```

> **Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
   "@odata.etag":"W/\"FC5BE4C2\"",
   "@odata.id":"/redfish/v1/Managers/4dbb506c-b0b6-4da3-87f0-9c70e37bf7b5.1/RemoteAccountService/Accounts/16",
   "@odata.type":"#ManagerAccount.v1_3_0.ManagerAccount",
   "Id":"16",
   "Name":"User Account",
   "Description":"BMC User Account",
   "UserName":"{username}",
   "RoleId":"{roleId}",
   "Links":{
      "Role":{
         "@odata.id":"/redfish/v1/Managers/{ManagerID}/RemoteAccountService/Roles/Administrator"
      }
   }
}
```

### Deleting a BMC account

|                                 |                                                              |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>Method</strong>         | `DELETE`                                                     |
| <strong>URI</strong>            | `/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Accounts/{AccountID}` |
| <strong>Description</strong>    | This operation deletes a BMC user account.                   |
| <strong>Returns</strong>        | JSON schema representing this user account                   |
| <strong>Response code</strong>  | On success, `204 No Content`                                 |
| <strong>Authentication</strong> | Yes                                                          |

> **curl command**

```
curl -i -X DELETE \
               -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
              'https://{odim_host}:{port}/redfish/v1/Managers/{ManagerID}/RemoteAccountService/Accounts/{AccountID}'
```

### Collection of BMC roles

|                                 |                                                              |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>Method</strong>         | `GET`                                                        |
| <strong>URI</strong>            | `/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Roles` |
| <strong>Description</strong>    | This operation retrieves information on the collection of user roles. |
| <strong>Returns</strong>        | Links to the user role instances                             |
| <strong>Response code</strong>  | On success, `200 Ok`                                         |
| <strong>Authentication</strong> | Yes                                                          |

> **curl command**

```
curl -i GET \
    -H "X-Auth-Token:{X-Auth-Token}" \
'https://{odim_host}:{port}/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Roles'
```

> **Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#RoleCollection.RoleCollection",
   "@odata.etag":"W/\"08A22FCA\"",
   "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-a3b8293774ba.1/
RemoteAccountService/Roles",
   "@odata.type":"#RoleCollection.RoleCollection",
   "Description":"iLO Roles Collection",
   "Members":[
      {
         "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-
a3b8293774ba.1/RemoteAccountService/Roles/Administrator"
      },
      {
         "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-
a3b8293774ba.1/RemoteAccountService/Roles/Operator"
      },
      {
         "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-
a3b8293774ba.1/RemoteAccountService/Roles/ReadOnly"
      },
      {
         "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-
a3b8293774ba.1/RemoteAccountService/Roles/dirgroupb3d8954f6ebbe735764e9f7c"
      },
      {
         "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-
a3b8293774ba.1/RemoteAccountService/Roles/dirgroup9d4546a03a03bb977c03086a"
      }
   ],
   "Members@odata.count":5,
   "Name":"Roles"
}
```

### Single role

|                                 |                                                              |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>Method</strong>         | `GET`                                                        |
| <strong>URI</strong>            | `/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Roles/{roleid}` |
| <strong>Description</strong>    | This operation retrieves information about a single user role. |
| <strong>Returns</strong>        | JSON schema representing this user role                      |
| <strong>Response code</strong>  | On success, `200 Ok`                                         |
| <strong>Authentication</strong> | Yes                                                          |

> **curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
'https://{odim_host}:{port}/redfish/v1/Managers/{ManagerId}/RemoteAccountService/Roles/{roleid}'
```

> **Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#Role.Role",
   "@odata.etag":"W/\"B60B0A30\"",
   "@odata.id":"/redfish/v1/Managers/4c7d1c54-4aea-4197-9892-a3b8293774ba.1/RemoteAccountService/Roles/Administrator",
   "@odata.type":"#Role.v1_2_1.Role",
   "Id":"Administrator",
   "Name":"User Role",
   "Description":"iLO User Role",
   "AssignedPrivileges":[
      "Login",
      "ConfigureManager",
      "ConfigureUsers",
      "ConfigureSelf",
      "ConfigureComponents"
   ],
   "IsPredefined":true,
   "RoleId":"{roleId}"
}
```



# Software and firmware inventory

The resource aggregator exposes Redfish update service endpoints. Use these endpoints to access and update the software components of a system such as BIOS and firmware. Using these endpoints, you can also upgrade or downgrade firmware of other components such as system drivers and provider software.

The `UpdateService` schema describes the update service and the properties for the service. It exposes the firmware and software inventory resources and provides links to access them.


**Supported endpoints**

|API URI|Supported operations|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/UpdateService|`GET`|`Login` |
|/redfish/v1/UpdateService/FirmwareInventory|`GET`|`Login` |
|/redfish/v1/UpdateService/FirmwareInventory/{inventoryId}|`GET`|`Login` |
|/redfish/v1/UpdateService/SoftwareInventory|`GET`|`Login` |
|/redfish/v1/UpdateService/SoftwareInventory/{inventoryId}|`GET`|`Login` |
|/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate|`POST`|`ConfigureComponents` |
|/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate|`POST`|`ConfigureComponents` |



## Viewing the UpdateService root

| | |
|-----|------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/UpdateService` |
|<strong>Description</strong> |This operation retrieves JSON schema representing the `UpdateService` root.|
|<strong>Returns</strong> |Properties for the service and a list of actions you can perform using this service|
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
    "@odata.type": "#UpdateService.v1_11_0.UpdateService",
    "@odata.id": "/redfish/v1/UpdateService",
    "@odata.context": "/redfish/v1/$metadata#UpdateService.UpdateService",
    "Id": "UpdateService",
    "Name": "Update Service",
    "Status": {
        "State": "Enabled",
        "Health": "OK",
        "HealthRollup": "OK"
    },
    "ServiceEnabled": true,
    "HttpPushUri": "",
    "FirmwareInventory": {
        "@odata.id": "/redfish/v1/UpdateService/FirmwareInventory"
    },
    "SoftwareInventory": {
        "@odata.id": "/redfish/v1/UpdateService/SoftwareInventory"
    },
    "Actions": {
        "#UpdateService.SimpleUpdate": {
            "target": "/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate",
            "@Redfish.OperationApplyTime": {
                "@odata.type": "#Settings.v1_3_3.OperationApplyTimeSupport",
                "SupportedValues": [
                    "OnStartUpdateRequest"
                ]
            }
        },
        "#UpdateService.StartUpdate": {
            "target": "/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate"
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
|<strong>Returns</strong> |A collection of links to firmware resources|
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
    "@odata.context": "/redfish/v1/$metadata#FirmwareInventoryCollection.FirmwareCollection",
    "@odata.id": "/redfish/v1/UpdateService/FirmwareInventory",
    "@odata.type": "#SoftwareInventoryCollection.SoftwareInventoryCollection",
    "Description": "FirmwareInventory view",
    "Name": "FirmwareInventory",
    "Members": [
        {
            "@odata.id": "/redfish/v1/UpdateService/FirmwareInventory/1c117017-37b7-4beb-b205-97ee73627d6c.7"
        },
        {
            "@odata.id": "/redfish/v1/UpdateService/FirmwareInventory/1c117017-37b7-4beb-b205-97ee73627d6c.4"
        },
        {
            "@odata.id": "/redfish/v1/UpdateService/FirmwareInventory/1c117017-37b7-4beb-b205-97ee73627d6c.6"
        },
        {
            "@odata.id": "/redfish/v1/UpdateService/FirmwareInventory/1c117017-37b7-4beb-b205-97ee73627d6c.11"
        },
        {
            "@odata.id": "/redfish/v1/UpdateService/FirmwareInventory/1c117017-37b7-4beb-b205-97ee73627d6c.13"
        },
        {
            "@odata.id": "/redfish/v1/UpdateService/FirmwareInventory/1c117017-37b7-4beb-b205-97ee73627d6c.10"
        },
        {
            "@odata.id": "/redfish/v1/UpdateService/FirmwareInventory/1c117017-37b7-4beb-b205-97ee73627d6c.9"
        },
        {
            "@odata.id": "/redfish/v1/UpdateService/FirmwareInventory/1c117017-37b7-4beb-b205-97ee73627d6c.12"
        }
    ],
    "Members@odata.count": 8
}
```

## Viewing a specific firmware resource

| | |
|-------|-----------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/UpdateService/FirmwareInventory/{inventoryId}` |
|<strong>Description</strong> |This operation retrieves information about a specific firmware resource.|
|<strong>Returns</strong> |JSON schema representing this firmware|
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
   "@odata.type":"#SoftwareInventory.v1_5_0.SoftwareInventory",
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
   "Status": {
        "Health": "OK",
        "State": "Enabled"
    },
    "Updateable": true,
   "Version": "1.0.0.20"
}
```

## Viewing the software inventory

| | |
|-------|-----------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/UpdateService/SoftwareInventory` |
|<strong>Description</strong> |This operation lists software of all the resources available in Resource Aggregator for ODIM.|
|<strong>Returns</strong> |A collection of links to software resources|
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
{  "@odata.context":"/redfish/v1/$metadata#SoftwareInventoryCollection.SoftwareCollection",
   "@odata.id":"/redfish/v1/UpdateService/SoftwareInventory",
   "@odata.type":"#SoftwareInventoryCollection.SoftwareInventoryCollection",
   "Description":"SoftwareInventory view",
   "Name":"SoftwareInventory",
   "Members":[],
   "Members@odata.count":0
}
```


## Viewing a specific software resource

| | |
|-------|-----------|
|<strong>Method</strong> | `GET` |
|<strong>URI</strong> |`/redfish/v1/UpdateService/SoftwareInventory/{inventoryId}` |
|<strong>Description</strong> |This operation retrieves information about a specific software resource.|
|<strong>Returns</strong> |JSON schema representing this software|
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
   "@odata.type":"#SoftwareInventory.v1_5_0.SoftwareInventory",
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
To know the progress of this action, perform HTTP `GET` on the *[task monitor](#viewing-a-task-monitor)* returned in the response header (until the task is complete).

> **curl command**

```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{
"ImageURI": "<URI_of_the_firmware_image>",
"Targets": ["/redfish/v1/Systems/{ComputerSystemId}"],
"@Redfish.OperationApplyTime": "OnStartUpdateRequest"
}' \
 'https://{odim_host}:{port}/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate'
```

> **Sample request body**

**Example 1:** 

```
{
  "ImageURI":"http://{IP_address}/ISO/resource.bin",
  "Targets": ["/redfish/v1/Systems/65d01621-4f88-49de-98bc-fcd1419bff3a.1"]
}
```

**Example 2:** 

```
{
  "ImageURI":"http://{IP_address}/ISO/resource.bin",
  "Targets": ["/redfish/v1/Systems/65d01621-4f88-49de-98bc-fcd1419bff3a.1"],
  "@Redfish.OperationApplyTime": "OnStartUpdateRequest"
}
```

> **Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|ImageURI|String (required)<br> |The URI of the software or firmware image to install. It is the location address of the software or firmware image you want to install.|
|Password|String (optional)<br> |The password to access the URI specified by the Image URI parameter.|
|Targets[]|Array (required)<br> |An array of URIs that indicate where to apply the update image.|
|TransferProtocol|String (optional)<br> | The network protocol that the update service uses to retrieve the software or the firmware image file at the URI provided in the `ImageURI` parameter, if the URI does not contain a scheme.<br> For the possible property values, see *Transfer protocol* table.<br> |
|Username|String (optional)<br> |The user name to access the URI specified by the Image URI parameter.|
|@Redfish.OperationApplyTime|Redfish annotation (optional)<br> | It enables you to control when the update is carried out.<br> Supported value is: `OnStartUpdate`. It indicates that the update will be carried out only after you perform HTTP POST on:<br> `/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate`.<br> |

|String|Description|
|------|-----------|
|CIFS|Common Internet File System.|
|FTP|File Transfer Protocol.|
|HTTP|Hypertext Transfer Protocol.|
|HTTPS|Hypertext Transfer Protocol Secure.|
| NFS (v1.3+)<br> |Network File System.|
| NSF (deprecated v1.3)<br> | Network File System.<br>This value has been deprecated in favor of NFS.<br> |
|OEM|A manufacturer-defined protocol.|
|SCP|Secure Copy Protocol.|
| SFTP (v1.1+)<br> |Secure File Transfer Protocol.|
|TFTP|Trivial File Transfer Protocol.|


>**Sample response header (HTTP 202 status)**

```
Location:/taskmon/task4aac9e1e-df58-4fff-b781-52373fcb5699
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes
```

>**Sample response body (HTTP 202 status)**

```
{
   "@odata.type":"#Task.v1_6_0.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "Name":"Task task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "Message":"The task with id task4aac9e1e-df58-4fff-b781-52373fcb5699 has started.",
   "MessageId":"TaskEvent.1.0.3.TaskStarted",
   "MessageArgs":[
      "task4aac9e1e-df58-4fff-b781-52373fcb5699"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```

>**Sample response body (HTTP 200 status)**

```
{
   "error":{
      "@Message.ExtendedInfo":[
         {
            "MessageId":"Base.1.13.0.Success"
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
|<strong>Description</strong> |This operation starts updating software or firmware components for which an update request has been created.<br>It is performed in the background as a Redfish task.<br>**IMPORTANT**: Before performing this operation, ensure that you have created an update request first. To know how to create an update request, see *[Simple update](#Simple update)*.|
|<strong>Response code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|

**Usage information** 
To know the progress of this action, perform HTTP `GET` on the *[task monitor](#viewing-a-task-monitor)* returned in the response header (until the task is complete).


```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
 'https://{odim_host}:{port}/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate'
```


> Sample request body

None

>**Sample response header (HTTP 202 status)**

```
Location:/taskmon/task4aac9e1e-df58-4fff-b781-52373fcb5699
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes
```

>**Sample response body (HTTP 202 status)**

```
{
   "@odata.type":"#Task.v1_6_0.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "Name":"Task task4aac9e1e-df58-4fff-b781-52373fcb5699",
   "Message":"The task with id task4aac9e1e-df58-4fff-b781-52373fcb5699 has started.",
   "MessageId":"TaskEvent.1.0.3.TaskStarted",
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
            "MessageId":"Base.1.13.0.Success"
         }
      ],
      "code":"iLO.0.10.ExtendedInfo",
      "message":"See @Message.ExtendedInfo for more information."
 }
}
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

**IMPORTANT**: 

- Before using the `Fabrics` APIs, ensure that the fabric manager is installed, its plugin is deployed, and added into the Resource Aggregator for ODIM framework. 
- The fabric is removed from Resource Aggregator for ODIM when the delete fabric event is received from the fabric plugin.


**Supported endpoints**


|API URI|Supported operations|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/Fabrics|`GET`|`Login` |
|/redfish/v1/Fabrics/{fabricId}|`GET`|`Login` |
|/redfish/v1/Fabrics/{fabricId}/Switches|`GET`|`Login` |
|/redfish/v1/Fabrics/{fabricId}/Switches/{switchId}|`GET`|`Login` |
| /redfish/v1/Fabrics/{fabricId}/Switches/{switchId}/Ports<br> |`GET`|`Login` |
| /redfish/v1/Fabrics/{fabricId} /Switches/{switchId}/Ports/{portid}<br> |`GET`|`Login` |
|/redfish/v1/Fabrics/{fabricId}/Zones|`GET`, `POST`|`Login`, `ConfigureComponents` |
|/redfish/v1/Fabrics/{fabricId}/Zones/{zoneid}|`GET`, `PATCH`, `DELETE`|`Login`, `ConfigureComponents` |
|/redfish/v1/Fabrics/{fabricId}/AddressPools|`GET`, `POST`|`Login`, `ConfigureComponents` |
|/redfish/v1/Fabrics/{fabricId}/AddressPools/{addresspoolid}|`GET`, `DELETE`|`Login`, `ConfigureComponents` |
|/redfish/v1/Fabrics/{fabricId}/Endpoints|`GET`, `POST`|`Login`, `ConfigureComponents` |
|/redfish/v1/Fabrics/{fabricId}/Endpoints/{endpointId}|`GET`, `DELETE`|`Login`, `ConfigureComponents` |


##  Collection of fabrics

|||
|---------------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/Fabrics` |
|**Description** |This operation retrieves a collection of simple fabrics.|
|**Returns** |Links to the fabric instances|
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
   "RedfishVersion":"1.15.1",
   "@odata.type":"#FabricCollection.FabricCollection"
}
```


## Single fabric

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |`/redfish/v1/Fabrics/{fabricID}` |
|**Description** |This operation retrieves a schema representing a specific fabric.|
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
   "@odata.type":"#Fabric.v1_2_2.Fabric",
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
|**Description** |This operation retrieves a collection of switches located in this fabric.|
|**Returns** |Links to the switch instances|
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
|**Description** |This operation retrieves JSON schema representing a particular fabric switch.|
|**Returns** |Details of this switch and links to its ports|
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
   "@odata.type":"#Switch.v1_6_0.Switch",
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
|**Description** |This operation retrieves a collection of ports of this switch.|
|**Returns** |Links to the port instances|
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
|**Description** |This operation retrieves a JSON schema representing a specific switch port.|
|**Returns** |Properties of this port|
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
   "@odata.type":"#Port.v1_5_0.Port",
   "CurrentSpeedGbps":0,
   "Description":"single port",
   "Id":"80b5f999-25e9-4b37-992c-de2f065ee0e3",
   "Links":{ 
      "ConnectedPorts":[ 
         { 
            "@odata.id":"/redfish/v1/Systems/768f9da7-56fc-4f13-b6f8-a1cd241e2313.1/EthernetInterfaces/3"
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
|**Description** |This operation retrieves a collection of address pools.|
|**Returns** |Links to the address pool instances|
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
	"RedfishVersion": "1.15.1",
	"@odata.type": "#AddressPoolCollection.AddressPoolCollection"
}	
```


## Single address pool

|||
|---------------|---------------|
|**Method** |`GET` |
|**URI** |``/redfish/v1/Fabrics/{fabricID}/AddressPools/{addresspoolid}`` |
|**Description** |This operation retrieves JSON schema representing a specific address pool.|
|**Returns** |Properties of this address pool|
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
      "IPv4FabricLinkLowerAddress":"xxx.xxx.xxx.8",
      "IPv4FabricLinkUpperAddress":"xxx.xxx.xxx.18",
      "IPv4GatewayAddress":"",
      "IPv4HostLowerAddress":"",
      "IPv4HostUpperAddress":"",
      "IPv4LoopbackLowerAddress":"xxx.xxx.xxx.28",
      "IPv4LoopbackUpperAddress":"xxx.xxx.xxx.38",
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
|**Description** |This operation retrieves a collection of fabric endpoints.|
|**Returns** |Links to the endpoint instances|
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
|**Description** |This operation retrieves JSON schema representing a specific fabric endpoint.|
|**Returns** |Details of this endpoint|
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
   "@odata.type":"#Endpoint.v1_6_1.Endpoint",
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
|**Description** |This operation retrieves a collection of fabric zones.|
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
|**Description** |This operation retrieves JSON schema representing a specific fabric zone.|
|**Returns** |Details of this zone|
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
   "@odata.type":"#Zone.v1_6_1.Zone",
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
|**Returns** |<ul><li>Link to the created address pool in the `Location` header</li><li>JSON schema representing the created address pool</li></ul>|
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
         "xxx.xxx.xxx.8/24",
         "xxx.xxx.xxx.9/24"  
      ],
      "AnycastGatewayIPAddress":"xxx.xxx.xxx.10"
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
         "xxx.xxx.xxx.8/24",
         "xxx.xxx.xxx.9/24" 
      ],
      "AnycastGatewayIPAddress":"xxx.xxx.xxx.10"
   }
}
```

> **Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Name|String (optional)<br> |Name for the address pool.|
|Description|String (optional)<br> |Description for the address pool.|
|IPv4{| (required)<br> | |
|VlanIdentifierAddressRange{| (optional)<br> | A single VLAN to assign on the ports or lags.<br> |
|Lower|Integer (required)<br> |VLAN lower address.|
|Upper}}|Integer (required)<br> |VLAN upper address.<br>**NOTE:** `Lower` and `Upper` must have the same value. Ensure that IP range is accurate and it does not overlap with other pools.|
|BgpEvpn{| (required)<br> | |
|GatewayIPAddressList|Array (required)<br> | IP pool to assign IPv4 address to the IP interface for VLAN per switch.<br> |
|AnycastGatewayIPAddress}|String (required)<br> | A single active gateway IP address for the IP interface.<br> |


>**Sample response header**

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Location:/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/e2ec196d-4b55-44b3-b928-8273de9fb8bf
Date:Thu, 14 May 2020 16:18:54 GMT
```

>**Sample response body**

```
{
  "@odata.id": "/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/e2ec196d-4b55-44b3-b928-8273de9fb8bf",
  "@odata.type": "#AddressPool.v1_2_0.AddressPool",
  "BgpEvpn": {
    "AnycastGatewayIPAddress": "xxx.xxx.xxx.10",
    "AnycastGatewayMACAddress": "",
    "GatewayIPAddressList": [
      "xxx.xxx.xxx.8/24",
      "xxx.xxx.xxx.9/24" 
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
|**Returns** |- Link to the created address pool in the `Location` header<br />- JSON schema representing the created address pool|
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
              "Lower": "xxx.xxx.xxx.8",
              "Upper": "xxx.xxx.xxx.14"
    },
    "EbgpAddressRange": {
              "Lower": "xxx.xxx.xxx.15",
              "Upper": "xxx.xxx.xxx.20"
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
    "GatewayIPAddressList": ["xxx.xxx.xxx.21/31", "xxx.xxx.xxx.22/31"]
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
              "Lower": "xxx.xxx.xxx.8",
              "Upper": "xxx.xxx.xxx.14"
    },
    "EbgpAddressRange": {
             "Lower": "xxx.xxx.xxx.15",
             "Upper": "xxx.xxx.xxx.20"
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
    "GatewayIPAddressList": ["xxx.xxx.xxx.21/31", "xxx.xxx.xxx.22/31"]
  }
}
```

> **Request parameters**

|Parameter|Type|Description|
|---------|----|-----------|
|Name|String|Name for the address pool.|
|Description|String (optional)<br> |Description for the address pool.|
|IPv4{| (required)<br> | |
|VlanIdentifierAddressRange{| (required)<br> | A single VLAN (virtual LAN) used for creating the IP interface for the user Virtual Routing and Forwarding (VRF).<br> |
|Lower|Integer (required)<br> |VLAN lower address|
|Upper}|Integer (required)<br> |VLAN upper address|
|IbgpAddressRange{| (required)<br> | IPv4 address used as the Router Id for the VRF per switch.<br> |
|Lower|String (required)<br> |IPv4 lower address|
|Upper}|String (required)<br> |IPv4 upper address|
|EbgpAddressRange{| (optional)<br> |External neighbor IPv4 addresses.|
|Lower|String (required)<br> |IPv4 lower address|
|Upper} }|String (required)<br> |IPv4 upper address|
|Ebgp{| (optional)<br> | |
|AsNumberRange{| (optional)<br> |External neighbor ASN.<br>**NOTE:** `EbgpAddressRange` and `AsNumberRange` values should be a matching sequence and should be of same length.|
|Lower|Integer (optional)<br> | |
|Upper} }|Integer (optional)<br> | |
|BgpEvpn{| (required)<br> | |
|RouteDistinguisherList|Array (required)<br> | Single route distinguisher value for the VRF.<br> |
|RouteTargetList|Array (optional)<br> | Route targets. By default, the route targets will be configured as both import and export.<br> |
|GatewayIPAddressList}|Array (required)<br> | IP pool to assign IPv4 address to the IP interface used by the VRF per switch.<br> |

>**Sample response header** 

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Location:/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d
Date:Thu, 14 May 2020 16:18:58 GMT
```

>**Sample response body**

```
{
   "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/AddressPools/84766158-cbac-4f69-8ed5-fa5f2b331b9d",
   "@odata.type":"#AddressPool.v1_2_0.AddressPool",
   "BgpEvpn":{
      "AnycastGatewayIPAddress":"",
      "AnycastGatewayMACAddress":"",
      "GatewayIPAddressList":[
         "xxx.xxx.xxx.8/31",
         "xxx.xxx.xxx.9/31"
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
         "Lower": "xxx.xxx.xxx.15",
         "Upper": "xxx.xxx.xxx.20"
      },
      "FabricLinkAddressRange":{
         "Lower":"",
         "Upper":""
      },
      "IbgpAddressRange":{
         "Lower": "xxx.xxx.xxx.8",
         "Upper": "xxx.xxx.xxx.14"
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
|**Description** |This operation creates an empty container zone for all the other zones in a specific fabric. To assign address pools, endpoints, other zones, or switches to this zone, perform an HTTP `PATCH` on this zone. See [Updating a Zone](#updating-a-zone).|
|**Returns** |JSON schema representing the created zone|
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

> **Request parameters**

|Parameter|Value|Description|
|---------|-----|-----------|
|Name|String (optional)<br> |Name for the zone.|
|Description|String (optional)<br> |Description for the zone.|
|ZoneType|String|The type of the zone to be created. Options include: `ZoneofZones` and `ZoneofEndpoints`<br> The type of the zone for a zone of zones is `ZoneofZones`. |
|Links{| (optional)<br> | |
|AddressPools|Array (optional)<br> | `AddressPool` links supported for the Zone of Zones (`AddressPool` links created for `ZoneofZones`).<br> |


>**Sample response header**

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Location:/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/a2dc8760-ea05-4cab-8f95-866c1c380f98
Date:Thu, 14 May 2020 16:19:00 GMT
```

>**Sample response body**

```
{
   "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/a2dc8760-ea05-4cab-8f95-866c1c380f98",
   "@odata.type":"#Zone.v1_6_1.Zone",
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
|**Returns** | <ul><li>Link to the created endpoint in the `Location` header</li><li>JSON schema representing the created endpoint</li></ul> |
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

> **Request parameters**

|Parameter|Value|Description|
|---------|-----|-----------|
|Name|String (optional)<br> |Name for the endpoint.|
|Description|String (optional)<br> |Description for the endpoint.|
|Links{| (required)<br> | |
|ConnectedPorts|Array (required)<br> | Switch port connected to the switch.<br>  <br> |
|Zones}|Array (optional)<br> | Endpoint is part of `ZoneofEndpoints`. Only one zone is permitted in the zones list.<br> |
|Redundancy[|Array| |
|Mode|String|Redundancy mode.|
|RedundancySet\]|Array| Set of redundancy ports connected to the switches.<br> |

>**Sample response header**

```
HTTP/1.1 201 Created
Allow:"GET", "PUT", "POST", "PATCH", "DELETE"
Location:/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Endpoints/fe34aff2-e81f-4167-a0c3-9bf5a67e2a97
Date:Thu, 14 May 2020 16:19:02 GMT
```

>**Sample response body**

```
{
   "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Endpoints/fe34aff2-e81f-4167-a0c3-9bf5a67e2a97",
   "@odata.type":"#Endpoint.v1_6_1.Endpoint",
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
|**Description** |This operation creates a zone of endpoints in a specific fabric.<br>**NOTE:** Ensure that the endpoints are created first before assigning them to the zone of endpoints.|
|**Returns** |JSON schema representing the created zone|
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

> **Request parameters**

|Parameter|Value|Description|
|---------|-----|-----------|
|Name|String (optional)<br> |The name for the zone.|
|Description|String (optional)<br> |The description for the zone.|
|DefaultRoutingEnabled|Boolean (required)<br> |Set to `false`.|
|ZoneType|String (required)<br> |The type of the zone to be created. Options include: `ZoneofZones`and `ZoneofEndpoints`<br>The type of the zone for a zone of endpoints is `ZoneofEndpoints`.<br> |
|Links{|Object (optional)<br> |Contains references to other resources that are related to the zone.|
|ContainedByZones [{|Array (optional)<br> |Represents an array of `ZoneofZones` for the zone being created (applicable when creating a zone of endpoints).|
|@odata.id }]|String|Link to a Zone of zones.|
|AddressPools [{|Array (optional)<br> |Represents an array of address pools linked with the zone \(zone-specific address pools\).|
|@odata.id }]|String|Link to an address pool.|
|Endpoints [{|Array (optional)<br> |Represents an array of endpoints to be included in the zone.|
|@odata.id }]|String|Link to an endpoint.|

>**Sample response header**

```
HTTP/1.1 201 Created
Allow: "GET", "PUT", "POST", "PATCH", "DELETE"
Location: /redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/06d344bb-cce1-4b0c-8414-6f6df1ea373f
Date: Thu, 14 May 2020 16:19:37 GMT
```

>**Sample response body**

```
{
   "@odata.id":"/redfish/v1/Fabrics/995c85a6-3de7-477f-af6f-b52de671abd5/Zones/06d344bb-cce1-4b0c-8414-6f6df1ea373f",
   "@odata.type":"#Zone.v1_6_1.Zone",
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
|**Description** |This operation assigns or unassigns a collection of endpoints, address pools, zone of zones, switches to a zone of endpoints, or a collection of address pools to a zone of zones in a specific fabric.|
|**Returns** |JSON schema representing an updated zone|
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
   "@odata.type":"#Zone.v1_6_1.Zone",
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
   "@odata.type":"#Zone.v1_6_1.Zone",
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
|**Description** |This operation deletes a zone in a specific fabric.<br>**NOTE:**<br> If you delete a non-empty zone (a zone which contains links to address pools, other zones, endpoints, or switches), you encounter an HTTP `400` error. Before attempting to delete, unassign all links in the zone. See *[updating a zone](#updating-a-zone)*.|
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
|**Description** |This operation deletes an address pool in a specific fabric.<br>**NOTE:**<br> If you delete an address pool that is being used in any zone, you encounter an HTTP `400` error. Before attempting to delete, ensure that the address pool you want to delete is not present in any zone. To get the list of address pools in a zone, see links to `addresspools` in the sample response body for a *[single zone](#single-zone)*.|
|**Response code** | `200 OK` |
|**Authentication** |Yes|

>**curl command**


```
curl -i -X DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/Fabrics/{fabricID}/AddressPools/{addresspoolid}'
```


# Tasks

A task represents an operation that takes more time than a user typically wants to wait and is carried out asynchronously.

An example of a task is resetting an aggregate of servers. Resetting all the servers in a group is a time-consuming operation. Users waiting for the result are blocked from performing other operations. Resource Aggregator for ODIM creates Redfish tasks for such long-duration operations and exposes Redfish `TaskService` APIs and `Task monitor` API. Use these APIs to manage and monitor the tasks until their completion, while performing other operations.

**Supported endpoints**

|API URI|Supported operations|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/TaskService|`GET`|`Login` |
|/redfish/v1/TaskService/Tasks|`GET`|`Login` |
|/redfish/v1/TaskService/Tasks/{taskId}|`GET`, `DELETE`|`Login`, `ConfigureManager` |
| /redfish/v1/ TaskService/Tasks/{taskId}/SubTasks<br> |`GET`|`Login` |
| /redfish/v1/ TaskService/Tasks/{taskId}/SubTasks/ {subTaskId}<br> |`GET`|`Login` |
|/taskmon/{taskId}|`GET`|`Login` |



##  Viewing the TaskService root

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

 **Sample response header** 

```
Allow:GET
Date:Sun,17 May 2020 15:11:12 GMT+5m 13s
Link:</redfish/v1/SchemaStore/en/TaskService.json>; rel=describedby
```

 **Sample response body** 

```
{
   "@odata.type":"#TaskService.v1_2_0.TaskService",
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
    "Tasks": {
        "@odata.id": "/redfish/v1/TaskService/Tasks"
    },
    "Oem": {
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
|**Returns** |JSON schema having the details of this task - task id, name, state of the task, start time and end time of this task, completion percentage, URI of the task monitor associated with this task, and subtasks if any. The sample response body given in this section is a JSON response for a task which adds a server.<br> |
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
   "@odata.type":"#Task.v1_6_0.Task",
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
         "Link: </redfish/v1/AggregationService/AggregationSources/7b08ecbd-d23e-4dd5-ad99-58ac2be7576d.1/>; rel=describedby",
         "Location: /redfish/v1/AggregationService/AggregationSources/7b08ecbd-d23e-4dd5-ad99-58ac2be7576d.1",
         "OData-Version: 4.0"
      ],
      "HttpOperation":"POST",
      "JsonBody": "{\"Context\":\"\",\"DeliveryRetryPolicy\":\"RetryForever1\",\"Destination\":\"https://node.odim.com:8080/Destination\",\"EventFormatType\":\"Event\",\"EventTypes\":[],\"MessageIds\":[],\"Name\":\"Bruce\",\"OriginResources\":[],\"Protocol\":\"Redfish\",\"ResourceTypes\":[],\"SubordinateResources\":true,\"SubscriptionType\":\"RedfishEvent\"}",
        "TargetUri": "/redfish/v1/EventService/Subscriptions"
    },
    "Oem": {
    }
}
```


##  Viewing a task monitor

|||
|-----------|----------|
|**Method** | `GET` |
|**URI** |`/taskmon/{TaskID}` |
|**Description** |This endpoint retrieves the task monitor associated with a specific task. A task monitor allows for polling a specific task for its completion. Perform `GET` on a task monitor URI to view the progress of a specific task (until it is complete).|
|**Returns** |Details of the task and its progress in the JSON response such as:<br>- Link to the task<br />- Task id<br />- Task state and status<br>- Percentage of completion<br>- Start time and end time<br>- Link to subtasks, if any<br>To know the status of a subtask, perform `GET` on the respective subtask link.<br>**NOTE:** <ul><li>Note down the task Id. If the task completes with an error, it is required to know which subtask has failed. To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`.</li><li>`EndTime` of an ongoing task has `0001-01-01T00:00:00Z` as value, which is equivalent to zero time stamp value. It is updated only after the completion of the task.</li></ul></li><li>On failure, an error message. See *Sample error response*.<br> To get the list of subtasks, perform `GET` on the task URI having the id of the failed task. To know which subtasks have failed, perform `GET` on subtask links individually.</li><li>On successful completion, result of the operation carried out by the task. See *Sample response body (completed task)*.</li></ul>|
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
Location:/taskmon/taskfbd5cdb0-5d33-4ad4-8682-cab90534ba70
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes
```


>**Sample response body** \(ongoing task\)

```
{
   "@odata.type":"#Task.v1_6_0.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/taskfbd5cdb0-5d33-4ad4-8682-cab90534ba70",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"taskfbd5cdb0-5d33-4ad4-8682-cab90534ba70",
   "Name":"Task taskfbd5cdb0-5d33-4ad4-8682-cab90534ba70",
   "Message":"The task with id taskfbd5cdb0-5d33-4ad4-8682-cab90534ba70 has started.",
   "MessageId":"TaskEvent.1.0.3.TaskStarted",
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
   "Oem": {
   }
}
```

>**Sample response body** \(completed task\)

```
{
"code": "Base.1.13.0.Success",
"message": "Request completed successfully."
}
```

>  **Sample error response**

```
{ 
   "error":{ 
      "code":"Base.1.13.0.GeneralError",
      "message":"one or more of the reset actions failed, check sub tasks for more info."
   }
```

## Deleting a task

|||
|-----------|----------|
|**Method** | `DELETE` |
|**URI** |`/redfish/v1/TaskService/Tasks/{TaskID}` |
|**Description** |This operation deletes a specific task. Deleting a running task aborts the operation being carried out.|
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

An event asynchronously notifies the client of some significant state change or error condition, usually of a time critical nature. Use the following APIs to subscribe a northbound client to southbound events by creating a subscription entry in the service.

**Supported endpoints**

|API URI|Supported operations|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/EventService|`GET`|`Login` |
|/redfish/v1/EventService/Subscriptions|`GET`, `POST`|`Login`, `ConfigureManager`, `ConfigureComponents` |
|/redfish/v1/EventService/Actions/EventService.SubmitTestEvent|`POST`|`ConfigureManager` |
|/redfish/v1/EventService/Subscriptions/{subscriptionId}|`GET`, `DELETE`|`Login`, `ConfigureManager`, `ConfigureSelf` |



## Viewing the EventService root

|||
|----------|---------|
|**Method** | `GET` |
|**URI** |`redfish/v1/EventService` |
|**Description** |This endpoint retrieves JSON schema for the Redfish `EventService` root.|
|**Returns** |Properties for managing event subscriptions such as allowed event types and a link to the actual collection of subscriptions|
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
Link:/v1/SchemaStore/en/EventService.json>; rel=describedby
Date:Fri,15 May 2020 10:10:15 GMT+5m 11s
```

>**Sample response body** 

```
{
   "@odata.context":"/redfish/v1/$metadata#EventService.EventService",
   "Id":"EventService",
   "@odata.id":"/redfish/v1/EventService",
   "@odata.type":"#EventService.v1_7_2.EventService",
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
      "MetricReport"
   ],
   "EventTypesForSubscription":[
      "StatusChange",
      "ResourceUpdated",
      "ResourceAdded",
      "ResourceRemoved",
      "Alert"
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
|**Description**| This endpoint subscribes a northbound client to events originating from a set of resources (southbound devices, managers, Resource Aggregator for ODIM itself\) by creating a subscription entry. For use cases, see *[Subscription use cases](#event-subscription-use-cases)*.<br>This operation is performed in the background as a Redfish task. If there is more than one resource that is sending a specific event, the task is further divided into subtasks. |
|**Returns** |<ul><li>`Location` URI of the task monitor associated with this operation in the response header.</li><li> Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See *Sample response body (HTTP 202 status)*.<br>**IMPORTANT:**<br> Make a note of the task Id. If the task completes with an error, it is required to know which subtask has failed. To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`.</li><li>On success, a `Location` header that contains a link to the newly created subscription and a message in the JSON response body saying that the subscription is created. See *Sample response body (HTTP 201 status)*.</li></ul>|
|**Response code** |<ul><li>`202 Accepted`</li><li>`201 Created`</li></ul>|
|**Authentication** |Yes|


To know the progress of this action, perform HTTP `GET` on the *[task monitor](#viewing-a-task-monitor)* returned in the response header (until the task is complete).

To get the list of subtask URIs, perform HTTP `GET` on the task URI returned in the JSON response body. See *Sample response body (HTTP 202 status)*. The JSON response body of each subtask contains a link to the task monitor associated with it. To know the progress of this operation (subtask) on a specific server, perform HTTP `GET` on the task monitor associated with the respective subtask.


>**curl command**

```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json; charset=utf-8" \
   -d \
'{ 
   "Name":"ODIMRA_NBI_client",
   "Destination":"https://{Valid_destination__IP_Address}:{Port}/EventListener",
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
   "DeliveryRetryPolicy": "RetryForever"
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
   ],
   "DeliveryRetryPolicy": "RetryForever"

}
```

> **Request parameters**

| Parameter            | Value                 | Attributes                         | Description                                                  |
| -------------------- | --------------------- | ---------------------------------- | ------------------------------------------------------------ |
| Name                 | String                | (optional)<br>                     | Name for the subscription.                                   |
| Destination          | String                | Read-only (Required on create)<br> | The URL of the destination event listener that listens to events (Fault management system or any northbound client).<br>**NOTE:** <br />Destinations with both IPv4 and IPv6 addresses are supported.<br />`Destination` is unique to a subscription. There can be only one subscription for a destination event listener.<br>To change the parameters of an existing subscription , delete it and then create again with the new parameters and a new destination URL.<br> |
| EventTypes           | Array (string (enum)) | Read-only (optional)<br>           | The types of events that are sent to the destination. For possible values, see *Event types* table. |
| ResourceTypes        | Array (string, null)  | Read-only (optional)<br>           | The list of resource type values (Schema names) that correspond to the `OriginResources`.  Examples: "Systems", "Chassis", "Tasks"<br>For possible values, perform `GET` on `redfish/v1/EventService` and check values listed under `ResourceTypes` in the JSON response.<br/> |
| Context              | String                | Read/write Required (null)<br>     | A string that is stored with the event destination subscription. |
| MessageIds           | Array                 | Read-only (optional)<br>           | The key used to find the message in a Message Registry.      |
| Protocol             | String (enum)         | Read-only (Required on create)<br> | The protocol type of the event connection. For possible values, see *Protocol* table. |
| SubscriptionType     | String (enum)         | Read-only Required (null)<br>      | Indicates the subscription type for events. For possible values, see *Subscription type* table. |
| EventFormatType      | String (enum)         | Read-only (optional)<br>           | Indicates the content types of the message that this service can send to the event destination. For possible values, see *EventFormat type* table. |
| SubordinateResources | Boolean               | Read-only (null)                   | Indicates whether the service supports the `SubordinateResource` property on event subscriptions or not. If it is set to `true`, the service creates subscription for an event originating from the specified `OriginResoures` and also from its subordinate resources. For example, by setting this property to `true`, you can receive specified events from a compute node: `/redfish/v1/Systems/{ComputerSystemId}` and from its subordinate resources such as:<br> `/redfish/v1/Systems/{ComputerSystemId}/Memory`<br> `/redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces`<br> `/redfish/v1/Systems/{ComputerSystemId}/Bios`<br> `/redfish/v1/Systems/{ComputerSystemId}/Storage` |
| OriginResources      | Array                 | Optional (null)<br>                | Resources for which the service sends related events. If this property is absent or the array is empty, events originating from any resource is sent to the subscriber. For possible values, see *[Origin resources](#origin-resources)* table. |
| DeliveryRetryPolicy  | String                | Optional                           | This property shall indicate the subscription delivery retry policy for events where the subscription type is `RedfishEvent`. Supported value is `RetryForever`, which implies that the attempts at delivery of future events shall continue regardless of the number of retries. |


> **Sample event**

~~~
{
  "@odata.context": "/redfish/v1/$metadata#Event.Event",
  "@odata.type": "#Event.v1_7_0.Event",
  "Events": [
    {
      "EventId": "aa378d6b-d612-e146-4d0c-6a58eb43179b",
      "EventTimestamp": "2022-07-05T08:54:42Z",
      "EventType": "Alert",
      "MemberId": "0",
      "Message": "",
      "MessageArgs": [
        "Off"
      ],
      "MessageId": "iLOEvents.2.3.IndicatorLEDStateChanged",
      "OriginOfCondition": {
        "@odata.id": "/redfish/v1/Systems/799bdb08-8bb6-4067-b2d5-3ffd87341b1a.1/"
      },
      "Severity": "OK"
    }
  ],
  "Name": "Events"
}
~~~

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

> **Request parameters**

|Parameter|Value|Attributes|Description|
|---------|-----|----------|-----------|
|Name|String| (optional)<br> |Name for the subscription.|
|Destination|String|Read-only (Required on create)<br> |The URL of the destination event listener that listens to events (Fault management system or any northbound client).<br/>**NOTE:** <br />Destinations with both IPv4 and IPv6 addresses are supported.<br />`Destination` is unique to a subscription. There can be only one subscription for a destination event listener.<br/>To change the parameters of an existing subscription , delete it and then create again with the new parameters and a new destination URL.<br/> |
|EventTypes|Array (string (enum))|Read-only (optional)<br> |The types of events that are sent to the destination. For possible values, see *Event types* table.|
|ResourceTypes|Array (string, null)|Read-only (optional)<br> |The list of resource type values (Schema names) that correspond to the `OriginResources`. For possible values, perform `GET` on `redfish/v1/EventService` and check values listed under `ResourceTypes` in the JSON response.<br> Examples: "ComputerSystem", "Storage", "Task"<br> |
|Context|String|Read/write Required (null)<br> |A string that is stored with the event destination subscription.|
|MessageIds|Array|Read-only (optional)<br> |The key used to find the message in a Message Registry.|
|Protocol|String (enum)|Read-only (Required on create)<br> |The protocol type of the event connection. For possible values, see *Protocol* table.|
|SubscriptionType|String (enum)|Read-only Required (null)<br> |Indicates the subscription type for events. For possible values, see *Subscription type* table.|
|EventFormatType|String (enum)|Read-only (optional)<br> |Indicates the content types of the message that this service can send to the event destination. For possible values, see *EventFormat type* table.|
|SubordinateResources|Boolean|Read-only (null)|Indicates whether the service supports the `SubordinateResource` property on event subscriptions or not. If it is set to `true`, the service creates subscription for an event originating from the specified `OriginResoures` and also from its subordinate resources. For example, by setting this property to `true`, you can receive specified events from a compute node: `/redfish/v1/Systems/{ComputerSystemId}` and from its subordinate resources such as:<br> `/redfish/v1/Systems/{ComputerSystemId}/Memory`<br> `/redfish/v1/Systems/{ComputerSystemId}/EthernetInterfaces`<br> `/redfish/v1/Systems/{ComputerSystemId}/Bios`<br> `/redfish/v1/Systems/{ComputerSystemId}/Storage`|
|OriginResources|Array| Optional (null)<br> |Resources for which the service only sends related events. If this property is absent or the array is empty, events originating from any resource is sent to the subscriber. For possible values, see *Origin resources* table.|

##### **Origin resources**

|String|Description|
|------|-----------|
|A single resource|A specific resource for which the service sends only related events.|
|A list of resources. Supported collections:<br> |A collection of resources for which the service will send only related events.|
|/redfish/v1/Systems|All computer system resources available in Resource Aggregator for ODIM for which the service sends only related events. By setting `EventType` property in the request payload to `ResourceAdded` or `ResourceRemoved` and `OriginResources` property to `/redfish/v1/Systems`, you can receive notifications when a system is added or removed in Resource Aggregator for ODIM.|
|/redfish/v1/Chassis|All chassis resources available in Resource Aggregator for ODIM for which the service sends only related events.|
|/redfish/v1/Fabrics|All fabric resources available in Resource Aggregator for ODIM for which the service sends only related events.|
|/redfish/v1/Managers|All manager resources available in Resource Aggregator for ODIM for which the service sends only related events.|
|/redfish/v1/TaskService/Tasks|All tasks scheduled by or being executed by Redfish `TaskService`. By subscribing to Redfish tasks, you can receive task status change notifications on the subscribed destination client.<br> By specifying the task URIs as `OriginResources` and `EventTypes` as `StatusChange`, you can receive notifications automatically when the tasks are complete.<br> To check the status of a specific task manually, perform HTTP `GET` on its task monitor until the task is complete.<br> |
| /redfish/v1/Aggregates/{AggregateId}         |Individual aggregate available in Resource Aggregator for ODIM for which the service sends only related events. |

**Event types**

|String|Description|
|------|-----------|
|Alert|A condition exists which requires attention|
|ResourceAdded|A resource has been added|
|ResourceRemoved|A resource has been removed|
|ResourceUpdated|The value of this resource has been updated|
|StatusChange|The status of this resource has changed|
|MetricReport|Collects resource metrics|

**EventFormat type**

|String|Description|
|------|-----------|
|Event|The subscription destination will receive JSON bodies of the Resource Type Event|
|MetricReport|Collects resource metrics|

**Subscription type**

|String|Description|
|------|-----------|
|RedfishEvent|The subscription follows the Redfish specification for event notifications, which is done by a service sending an HTTP `POST` to the destination URI of the subscriber.|

**Protocol**

|String|Description|
|------|-----------|
|Redfish|The destination follows the Redfish specification for event notifications.|

 

>**Sample response header** (HTTP 202 status) 

```
Location:/taskmon/taska9702e20-884c-41e2-bd9c-d779a4dd2e6e
Date:Fri, 08 Nov 2019 07:49:42 GMT+7m 9s
Content-Length:0 byte
```

>**Sample response header** (HTTP 201 status) 

```
Location:/redfish/v1/EventService/Subscriptions/76088e1c-4654-4eec-a3f6-60bc33b77cdb
Date:Thu,14 May 2020 09:48:23 GMT+5m 10s
```

>**Sample response body** (HTTP 202 status) 

```
{
   "@odata.type":"#Task.v1_6_0.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/taskbab2e46d-2ef9-40e8-a070-4e6c87ef72ad",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"taskbab2e46d-2ef9-40e8-a070-4e6c87ef72ad",
   "Name":"Task taskbab2e46d-2ef9-40e8-a070-4e6c87ef72ad",
   "Message":"The task with id taskbab2e46d-2ef9-40e8-a070-4e6c87ef72ad has started.",
   "MessageId":"TaskEvent.1.0.3.TaskStarted",
   "MessageArgs":[
      "taskbab2e46d-2ef9-40e8-a070-4e6c87ef72ad"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```

>**Sample response body** (subtask) 


```
{
    "@odata.type": "#Task.v1_6_0.Task",
    "@odata.id": "/redfish/v1/TaskService/Tasks/taskd862139f-c664-4cb2-b771-3e702bde40e3",
    "@odata.context": "/redfish/v1/$metadata#Task.Task",
    "Id": "taskd862139f-c664-4cb2-b771-3e702bde40e3",
    "Name": "Task taskd862139f-c664-4cb2-b771-3e702bde40e3",
    "TaskState": "Completed",
    "StartTime": "2022-02-25T13:07:05.938018291Z",
    "EndTime": "2022-02-25T13:07:08.108846323Z",
    "TaskStatus": "OK",
    "SubTasks": "/redfish/v1/TaskService/Tasks/taskd862139f-c664-4cb2-b771-3e702bde40e3/SubTasks",
    "TaskMonitor": "/taskmon/taskd862139f-c664-4cb2-b771-3e702bde40e3",
    "PercentComplete": 100,
    "Payload": {
        "HttpHeaders": [
        ],
        "HttpOperation": "POST",
        "JsonBody": "{\"BatchSize\":2,\"DelayBetweenBatchesInSeconds\":2,\"ResetType\":\"ForceRestart\"}",
        "TargetUri": "/redfish/v1/AggregationService/Aggregates/ca3f2462-15b5-4eb6-80c1-89f99ac36b12/Actions/Aggregate.Reset"
    },
    "Oem": {
    }
}
```

>**Sample response body** (HTTP 201 status) 

```
{
   "error":{
      "@Message.ExtendedInfo":[
         {
            "MessageId":"Base.1.13.0.Created"
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
|**Description** | Once the subscription is successfully created, you can post a test event to Resource Aggregator for ODIM to check whether you are able to receive events. If the event is successfully posted, you will receive a JSON payload of the event response on the client machine (destination) that is listening to events. To know more about this event, look up the message registry using the `MessageId` received in the payload. See *Sample message registry (Alert.1.0.0)*. For more information on message registries, see *[Message registries](#message-registries)*. |
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



> **Sample event payload** 

```
{ 
   "EventGroupId":1,
   "EventId":"132489713478812346",
   "EventTimestamp":"2020-02-17T17:17:42-0600",
   "EventType":"Alert",
   "Message":"The LAN has been disconnected",
   "MessageArgs":[ 
       "EthernetInterface 1",
            "/redfish/v1/Systems/8fbda4f3-f55f-4fe4-8db8-4aec1dc3a7d7.1"
   ],
   "MessageId":"Alert.1.0.LanDisconnect",
   "OriginOfCondition":"/redfish/v1/Systems/8fbda4f3-f55f-4fe4-8db8-4aec1dc3a7d7.1/EthernetInterfaces/1",
   "Severity":"Critical"
}
```

> **Request parameters** 

|Parameter|Value|Attributes|Description|
|---------|-----|----------|-----------|
|EventGroupId|Integer|Optional|The group Id for the event.|
|EventId|String|Optional|The Id for the event to add. This id is a string of a unique positive integer. Generate a random positive integer and use it as the Id.|
|EventTimestamp|String|Optional|The date and time stamp for the event to add. When the event is received, it translates as the time the event occurred.|
|EventType|String (enum)|Optional|The type for the event to add. For possible property values, see *EventType* in *[Creating an event subscription](#creating-an-event-subscription)*.|
|Message|String|Optional|The human-readable message for the event to add.|
|MessageArgs [ ]|Array (string)|Optional|An array of message arguments for the event to add. The message arguments are substituted for the arguments in the message when looked up in the message registry. It helps in trouble ticketing when there are bad events. For example, `MessageArgs` in *Sample event payload* has the following two substitution variables:<br><ul><li>`EthernetInterface 1`</li><li>`/redfish/v1/Systems/{ComputerSystemId}`</li></ul><br>`Description` and `Message` values in "Sample message registry" are substituted with the above-mentioned variables. They translate to "A LAN Disconnect on `EthernetInterface 1` was detected on system `/redfish/v1/Systems/{ComputerSystemId}.` |
|MessageId|String|Required|The Message Id for the event to add. It is the key used to find the message in a message registry. It has `RegistryPrefix` concatenated with the version, and the unique identifier for the message registry entry. The `RegistryPrefix` concatenated with the version is the name of the message registry. To get the names of available message registries, perform HTTP `GET` on `/redfish/v1/Registries`. The message registry mentioned in the sample request payload is `Alert.1.0`.|
|OriginOfCondition|String|Optional|The URL in the `OriginOfCondition` property of the event to add. It is not a reference object. It is the resource that originated the condition that caused the event to be generated. For possible values, see *Origin resources* in *[Creating an event subscription](#creating-an-event-subscription)*.|
|Severity|String|Optional|The severity for the event to add. For possible values, see *Severity* table.|

> **Severity**

|String|Description|
|------|-----------|
|Critical|A critical condition that requires immediate attention|
|OK|Informational or operating normally|
|Warning|A condition that requires attention|

> **Sample response header** 

```
Date:Fri,15 May 2020 07:42:59 GMT+5m 11s
```

> **Sample event response** 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#Event.Event",
   "@odata.id":"/redfish/v1/EventService/Events/1",
   "@odata.type":"#Event.v1_7_0.Event",
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
            "/redfish/v1/Systems/8fbda4f3-f55f-4fe4-8db8-4aec1dc3a7d7.1"
         ],
         "OriginOfCondition":"/redfish/v1/Systems/8fbda4f3-f55f-4fe4-8db8-4aec1dc3a7d7.1/EthernetInterfaces/1",
         "Context":"Event Subscription"
       }
   ]
}
```

> Sample message registry \(Alert.1.0.0\) 

```
{
"@odata.type": "#MessageRegistry.v1_4_2.MessageRegistry",
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
```

## Event subscription use cases

### Subscribing to resource addition notification

> **Resource addition notification payload**

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


> **Resource removal notification payload**

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


> **Task status notification payload**

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

 **Sample response body**

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
|**Returns** |JSON schema having the details of this subscription–subscription id, destination, event types, origin resource, and so on|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/EventService/Subscriptions/{subscriptionId}'
```

 **Sample response body** 

```
{
   "@odata.type":"#EventDestination.v1_11_0.EventDestination",
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
      {
      "@odata.id":"/redfish/v1/Systems/936f4838-9ce5-4e2a-9e2d-34a45422a389.1"
      }
   ]
}
```


##  Deleting an event subscription

|||
|-----------|-----------|
|**Method** | `DELETE` |
|**URI** |`/redfish/v1/EventService/Subscriptions/{subscriptionId}` |
|**Description** |To unsubscribe from an event, delete the corresponding subscription entry. Perform `DELETE` on this URI to remove an event subscription entry.|
|**Returns** |A message in the JSON response body about the subscription removal.|
|**Response code** |`200 OK` |
|**Authentication** |Yes|

>**curl command**

```
curl -i -X DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/EventService/Subscriptions/{subscriptionId}'
```

 **Sample response body** 

```
{
   "@odata.type":"#EventDestination.v1_11_0.EventDestination",
   "@odata.id":"/redfish/v1/EventService/Subscriptions/57e22fcc-8b1a-460c-ac1f-b3377e22f1cf",
   "Id":"57e22fcc-8b1a-460c-ac1f-b3377e22f1cf",
   "Name":"Event Subscription",
   "Message":"The resource has been removed successfully.",
   "MessageId":"ResourceEvent.1.0.2.ResourceRemoved",
   "Severity":"OK"
}
```

## Undelivered events

In instances where your subscribed destination is unavailable to listen to the events for a certain period, the events are saved in the product database as undelivered events. By default, Resource Aggregator for ODIM tries to repost the undelivered events three times in the interval of every 60 seconds. 

Eventually, when the destination becomes available for the new events to be published, the undelivered events are published to the destination and are deleted from the database.

You can configure the number of reposting instances and the required time interval by editing the values for `DeliveryRetryAttempts` and `DeliveryRetryIntervalSeconds` properties.




# Message registries

The`MessageRegistry` endpoint represents the properties for a message registry.

A message registry is an array of messages and their attributes organized by `MessageId`. Each entry has:

-   Description

-   Message this id translates to

-   Severity

-   Number and type of arguments

-   Proposed resolution


The arguments are the substitution variables for the message. The `MessageId` is formed according to the Redfish specification. It consists of the `RegistryPrefix` concatenated with the version and the unique identifier for the message registry entry.

**Supported endpoints**

|API URI|Supported operations|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/Registries|`GET`|`Login` |
|/redfish/v1/Registries/{registryId}|`GET`|`Login` |
|/redfish/v1/registries/{registryFileId}|`GET`|`Login` |


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
   "@odata.id":"/redfish/v1/Registries",
   "@odata.type":"#MessageRegistryFileCollection.MessageRegistryFileCollection",
   "Name":"Registry File Repository",
   "Description":"Registry Repository",
   "Members@odata.count":14,
   "Members":[
      {
         "@odata.id":"/redfish/v1/Registries/Base.1.13.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Composition.1.1.1"
      },
      {
         "@odata.id":"/redfish/v1/Registries/EthernetFabric.1.0.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Fabric.1.0.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/JobEvent.1.0.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/License.1.0.1"
      },
      {
         "@odata.id":"/redfish/v1/Registries/LogService.1.0.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/NetworkDevice.1.0.1"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Redfish_1.3.0_PrivilegeRegistry"
      },
      {
         "@odata.id":"/redfish/v1/Registries/ResourceEvent.1.0.3"
      },
      {
         "@odata.id":"/redfish/v1/Registries/ResourceEvent.1.2.1"
      },
      {
         "@odata.id":"/redfish/v1/Registries/StorageDevice.1.1.0"
      },
      {
         "@odata.id":"/redfish/v1/Registries/TaskEvent.1.0.3"
      },
      {
         "@odata.id":"/redfish/v1/Registries/Update.1.0.1"
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
   "@Redfish.Copyright":"Copyright 2014-2022 DMTF. All rights reserved.",
   "@Redfish.License":"Creative Commons Attribution 4.0 License.  For full text see link: https://creativecommons.org/licenses/by/4.0/",
   "@odata.type":"#MessageRegistry.v1_5_0.MessageRegistry",
   "Description":"This registry defines the base messages for Redfish",
   "Id":"Base.1.13.0",
   "Language":"en",
   "Messages":{
      "Name":"Base Message Registry",
      "OwningEntity":"DMTF",
      "RegistryPrefix":"Base",
      "RegistryVersion":"1.13.0"
   }
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

>**Sample response body**

```
{
   "Id":"Base.1.13.0",
   "@odata.context":"/redfish/v1/$metadata#MessageRegistryFile.MessageRegistryFile",
   "@odata.id":"/redfish/v1/Registries/Base.1.13.0",
   "@odata.type":"#MessageRegistryFile.v1_1_3.MessageRegistryFile",
   "Name":"Registry File Repository",
   "Description":"Base Message Registry File Locations",
   "Languages":[
      "en"
   ],
   "Location":[
      {
         "Language":"en",
         "Uri":"/redfish/v1/Registries/Base.1.13.0.json"
      }
   ],
   "Registry":"Base.1.13.0"
}
```



# Redfish Telemetry Service

Telemetry refers to the metrics obtained from remote systems for analysis and monitoring. 
The Redfish Telemetry model is designed to obtain characteristics of metrics, send specific metric reports periodically and specify triggers against metrics.

Resource Aggregator for ODIM exposes the Redfish `TelemetryService` APIs to perform the following tasks:

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
| /redfish/v1/TelemetryService/MetricReportDefinitions         | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/MetricReportDefinitions/{MetricReportDefinitionID} | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/MetricReports                   | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/MetricReports/{MetricReportID}  | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/Triggers                        | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/Triggers/{TriggerID}            | GET, PATCH           | `Login`,`ConfigureSelf` |

## Viewing the TelemetryService root

| **Method**         | `GET`                                                        |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/TelemetryService`                               |
| **Description**    | This operation retrieves JSON schema representing the Redfish `TelemetryService` root. |
| **Returns**        | Properties for the Redfish `TelemetryService` and links to its list of resources |
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
   "@odata.context": "/redfish/v1/$metadata#TelemetryService.TelemetryService",
   "@odata.type":"#TelemetryService.v1_3_1.TelemetryService",
   "Id":"TelemetryService",
   "Name":"Telemetry Service",
   "Status":{
      "State":"Enabled",
      "Health":"OK",
      "HealthRollup":"OK"
   },
   "ServiceEnabled":true,
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
    "@odata.type": "#MetricDefinition.v1_2_0.MetricDefinition",
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
                "9616fec9-c76a-4d26-ab53-196d08ce825a.1",
                "ba5cd083-b360-4994-bc30-12b450859b27.1"
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
    "@odata.type": "#MetricReportDefinition.v1_4_1.MetricReportDefinition",
    "Description": "Metric report of CPU Utilization for 10 minutes with sensing interval of 20 seconds.",
    "Id": "CPUUtilCustom1",
    "MetricProperties": [
        "/redfish/v1/Systems/{SystemID}#SystemUsage/CPUUtil",
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
                "9616fec9-c76a-4d26-ab53-196d08ce825a.1"
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
| **Returns**        | Links of the metric reports                                  |
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
| **Returns**        | Link to the metric report                                    |
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
         "MetricProperty":"/redfish/v1/Systems/9616fec9-c76a-4d26-ab53-196d08ce825a.1#SystemUsage/CPUUtil",
         "MetricValue":"1",
         "Timestamp":"2021-08-28T13:46:05Z"
      },
      {
         "MetricDefinition":{
            "@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"
         },
         "MetricId":"CPUUtil",
         "MetricProperty":"/redfish/v1/Systems/9616fec9-c76a-4d26-ab53-196d08ce825a.1#SystemUsage/CPUUtil",
         "MetricValue":"1",
         "Timestamp":"2021-08-28T13:46:25Z"
      },
      {
         "MetricDefinition":{
            "@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"
         },
         "MetricId":"CPUUtil",
         "MetricProperty":"/redfish/v1/Systems/9616fec9-c76a-4d26-ab53-196d08ce825a.1#SystemUsage/CPUUtil",
         "MetricValue":"1",
         "Timestamp":"2021-08-28T13:46:45Z"
      },
      {
         "MetricDefinition":{
            "@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"
         },
         "MetricId":"CPUUtil",
         "MetricProperty":"/redfish/v1/Systems/9616fec9-c76a-4d26-ab53-196d08ce825a.1#SystemUsage/CPUUtil",
         "MetricValue":"1",
         "Timestamp":"2021-08-28T13:47:05Z"
      },
      "Name":"Metric report of CPU Utilization for 60 minutes with sensing interval of 20 seconds."
   }
```

> **NOTE**:  After you remove a system and perform a `GET ` operation on the Metric Report Collection, the collection of all individual metric reports is still displayed in the response body. When you perform a `GET` operation on that individual {MetricReportID}, you get a `404-Not Found` error message. After this, when you perform a GET operation on the Metric Report Collection again, the instance of that individual metric report is erased. 
> This is an implementation choice in Resource Aggregator for ODIM, because Telemetry service is defined for a collection of BMCs and not for an individual BMC as per the DMTF Redfish specification.


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
                "9616fec9-c76a-4d26-ab53-196d08ce825a.1",
                "ba5cd083-b360-4994-bc30-12b450859b27.1"
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
 'https://{odimra_host}:{port}/redfish/v1/TelemetryService/Triggers/{TriggersID}'

```


>**Sample request body**

```
{
  "EventTriggers": ["Alert"]
}
```



# License Service

Resource Aggregator for ODIM offers `LicenseService` APIs to view and install licenses on multiple BMC servers.

**Supported APIs**

| API URI                                         | Supported operations | Required privileges         |
| ----------------------------------------------- | -------------------- | --------------------------- |
| /redfish/v1/LicenseService                      | `GET`                | `Login`                     |
| /redfish/v1/LicenseService/Licenses/            | `GET`, `POST`        | `Login`, `ConfigureManager` |
| /redfish/v1/LicenseService/Licenses/{LicenseID} | `GET`                | `Login`                     |

## Viewing the LicenseService root

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/LicenseService`                                 |
| **Description**    | This endpoint fetches JSON schema representing the Redfish `LicenseService` root. |
| **Returns**        | Properties for viewing the service and links to the actual collections of manager licenses |
| **Response Code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/LicenseService'
```


>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#LicenseService.LicenseService",
   "@odata.id":"/redfish/v1/LicenseService",
   "@odata.type":"#LicenseService.v1_0_0.LicenseService",
   "Description":"License Service",
   "Id":"LicenseService",
   "Name":"License Service",
   "Licenses":{
      "@odata.id":"/redfish/v1/LicenseService/Licenses"
   },
   "ServiceEnabled":true
}
```

## Viewing the license collection

|                    |                                                              |
| ------------------ | ------------------------------------------------------------ |
| **Method**         | `GET`                                                        |
| **URI**            | `/redfish/v1/LicenseService/Licenses/`                       |
| **Description**    | This endpoint fetches JSON schema representing the available License collections. |
| **Returns**        | Links to the licenses collection                             |
| **Response Code**  | `200 OK`                                                     |
| **Authentication** | Yes                                                          |

>**curl command**

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/LicenseService/Licenses/'
```


>**Sample response body**

```
{
   "@odata.context":"/redfish/v1/$metadata#LicenseCollection.LicenseCollection",
   "@odata.id":"/redfish/v1/LicenseService/Licenses",
   "@odata.type":"#LicenseCollection.v1_0_0.LicenseCollection",
   "Description":"License Collection",
   "Name":"License Collection",
   "Members":[
      {
         "@odata.id":"/redfish/v1/LicenseService/Licenses/8dd3fb4d-0429-4262-989f-906df092aefd.1.1"
      }
   ],
   "Members@odata.count":1
}
```

## Viewing information about a license

|                                 |                                                              |
| ------------------------------- | ------------------------------------------------------------ |
| <strong>Method</strong>         | `GET`                                                        |
| <strong>URI</strong>            | /redfish/v1/LicenseService/Licenses/{LicenseID}              |
| <strong>Description</strong>    | This endpoint retrieves information about a specific license of a server. |
| <strong>Returns</strong>        | JSON schema representing the single license                  |
| <strong>Response code</strong>  | On success, `200 OK`                                         |
| <strong>Authentication</strong> | Yes                                                          |


>**curl command**


```
curl -i GET \
             -H "X-Auth-Token:{X-Auth-Token}" \
              'https://{odim_host}:{port}/redfish/v1/LicenseService/Licenses/{LicenseID}'
```

>**Sample response body** 

```
{
   "@odata.context":"/redfish/v1/$metadata#License.License",
   "@odata.id":"/redfish/v1/LicenseService/Licenses/8dd3fb4d-0429-4262-989f-906df092aefd.1.1",
   "@odata.type":"#License.v1_0_0.License",
   "Id":"8dd3fb4d-0429-4262-989f-906df092aefd.1.1",
   "Name":"iLO License",
   "Description":"iLO License View",
   "ExpirationDate": "Activated until 14 Jan 2023",
   "InstallDate": "16 May 2022",
   "LicenseType":"Perpetual",
   "SerialNumber":"CN704614C4",
}
```



## Installing a license

|                    |                                                      |
| ------------------ | ---------------------------------------------------- |
| **Method**         | `POST`                                               |
| **URI**            | `/redfish/v1/LicenseService/Licenses`                |
| **Description**    | This endpoint installs a license on the BMC servers. |
| **Returns**        | No content                                           |
| **Response Code**  | `204 No Content`                                     |
| **Authentication** | Yes                                                  |

>**curl command**

```
curl -i -X POST \
  -H "Authorization:Basic YWRtaW46T2QhbTEyJDQ=" \
  -H "Content-Type:application/json" \
  -d \
'{
  "LicenseString": "xxxxx-xxxxx-xxxxx-xxxxx-xxxxx",
  "Links": {
  "AuthorizedDevices": [{
"@odata.id": "/redfish/v1/Systems/78869dd2-d2e2-4a49-854f-d495f873f199.1"
     }
  ]
 }
}
' \
'https://{odim_host}:{port}/redfish/v1/LicenseService/Licenses'

```


>**Sample request body**

```
{
  "LicenseString": "xxxxx-xxxxx-xxxxx-xxxxx-xxxxx",
  "Links": {
  "AuthorizedDevices": [{
  "@odata.id": "/redfish/v1/Systems/78869dd2-d2e2-4a49-854f-d495f873f199.1"
   }
  ]
 }
}
```

**Request parameter**

| Parameter     | Type   | Description                              |
| ------------- | ------ | ---------------------------------------- |
| LicenseString | string | The base64-encoded string of the license |



# Audit logs

Audit logs provide information on each API and are stored in the `api.log` file in `odimra` logs. Each log consists of a priority value, date and time of the log, hostname from which the APIs are sent, user account and role details, API request method and resource, response body, response code, and the message.

**Sample logs**

```
<110> 2009-11-10T23:00:00Z xxx.xxx.xxx.xxx [account@1 user="admin" roleID="Administrator"][request@1 method="GET" resource="/redfish/v1/Systems" requestBody=""][response@1 responseCode=200] Operation Successful
```

```
<107> 2009-11-10T23:00:00Z xxx.xxx.xxx.xxx [account@1 user="admin" roleID="Administrator"][request@1 method="GET" resource="/redfish/v1/Systems" requestBody=""][response@1 responseCode=404] Operation failed
```

> **Note**: <110> and <107> are priority values. <110> is the audit information log and <107> is the audit error log.



# Security logs

Security logs provide information on the successful and failed user authentication and authorization attempts. The logs are stored in `api.log` and `account_session.log` file in `odimra` logs. Each log consists of a priority value, date and time of the log, user account and role details, and the message.

**Sample logs**

```
<86> 2022-01-28T04:44:09Z [account@1 user="admin" roleID="Administrator"] Authentication/Authorization successful for session token 388281e8-4a45-45e5-862b-6b1ccfd6e6a3
```

```
<84> 2022-01-28T04:43:39Z [account@1 user="admin1" roleID="null"] Authentication failed, Invalid username or password
```

<blockquote> Note: <86> and <84> are priority values. <86> is security information log and <84> is the warning log.</blockquote>



# Application logs

Application logs provide information on all operations performed during specific times.

**Sample log**

```
<11>1 2022-12-21T07:39:49Z odim-host  svc-managers  managers-b5465c4df-fj4ss_9  GetManager [process@1 processName="managers-b5465c4df-fj4ss" transactionID="b3b66d09-f844-41bd-8b8d-957addf09b20" actionID="169" actionName="GetManager" threadID="0" threadName="svc-managers"] unable to get managers details: no data with the key /redfish/v1/Managers/386710f8-3a38-4938-a986-5f1048f487fdf found
```

The following table lists the properties, values and description of the application logs:

| Property      | Value (given in sample log)                                  | Description                                                  |
| ------------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| pri           | <11>                                                         | Priority value. <br/>For the list of all priority values, see *Log levels* section in the *Resource Aggregator for ODIM Getting Started Readme*. |
| version       | 1                                                            | Version number of the syslog protocol specification.         |
| timestamp     | 2022-12-21T07:39:49Z                                         | Date and time of the log.                                    |
| hostName      | odim-host                                                    | Name of the host from which the syslog messages are sent.    |
| appName       | svc-managers                                                 | Application that originated the message.                     |
| procID        | managers                                                     | Process name or process ID associated with the syslog system. |
| msgID         | GetManager                                                   | Identifies the type of message.                              |
| processName   | managers-b5465c4df-fj4ss                                     | System that originally sent the message.                     |
| transactionID | b3b66d09-f844-41bd-8b8d-957addf09b20                         | Unique UUID for each API request.                            |
| actionID      | 169                                                          | Unique ID for each actionName.                               |
| actionName    | GetManager                                                   | HTTP operation performed.                                    |
| threadID      | 0                                                            | Unique ID of the current running thread.                     |
| threadName    | svc-managers                                                 | Name of the current running thread.                          |
| message       | no data with the key /redfish/v1/Managers/386710f8-3a38-4938-a986-5f1048f487fdf found | Logged result message.                                       |