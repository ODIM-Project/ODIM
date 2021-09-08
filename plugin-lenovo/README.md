# Lenovo-plugin
  
Lenovo-plugin communicates with redfish compliant BMC.
This is an independent module which provides two primary communication channels:  
1.  An API mechanism that is used to exchange control data  
2.  An Event Message Bus (EMB) that is used to exchange event and notifications.


This guide provides a set of guidelines for developing API and EMB functions to work within the Resource Aggregator for ODIM™ environment. It ensures consistency around API semantics for all plugins.

To ensure continued adoption of open technologies, the APIs for the plugins are based on the [OpenAPI specification](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md). Messaging is based on the now-evolving [OpenMessaging specifications](https://github.com/openmessaging/specification/blob/master/domain_architecture.md) under Linux Foundation.




## API accessibility

The plugin layer uses JSON as the primary data format for communication. Standardizing on a well-known data-interchange format ensures consistency among plugins and simplifies the task for plugin developers. The API service uses [HATEOAS \(Hypermedia as the Engine of Application State\)](https://restfulapi.net/hateoas/) principles to link resources using the `href` key.

The API service under the plugin layer uses token-based authentication for securing the platform. The token-based authentication is applicable to:

-   The authentication information flowing from the aggregator to the plugin where the aggregator is authenticated.

-   The data flowing from the plugin to the aggregator where the aggregator is authenticated.


The plugin currently uses credentials of the client for authenticating the same.

Data on the wire is encrypted using TLS and is not sent out as cleartext. For this, the plugin exposes a CA signed certificate for the clients to authenticate itself. The plugins communicate primarily with the aggregator. To gather resource information, they can also communicate with another plugin through the north-bound APIs provided by the aggregator. For plugin-to-plugin communication, the aggregator defines a plugin role to set and allows permissions for plugins to communicate with other plugins.

API operations must adhere to the standard Restful API rules—Ensure that the API operations are not idempotent and concurrent. APIs can, in selective cases, implement capabilities to use subresources, filtering, sorting, and other value additions effectively. Return codes are fully in compliance with HTTP. A core objective of the plugin layer is to be able to perform many different operations using the primary HTTP operations —GET, PUT, POST, and DELETE.

The primary media type for plugin content is `Application/json`. The future releases may have other media types and custom media types.


## Naming conventions

-   `PascalCase` is used for all naming requirements within the plugin. It includes the plugin, the resources it contains, functions it implements, and the variables it defines.

         <aside class="notice">
	     NOTE: `PascalCase` is a naming convention where the first letter of each word is capitalized. It aligns with [DMTF standards for Redfish](http://redfish.dmtf.org/schemas/DSP0266_1.6.1.html), which is another key specification for ODIM™.
	     </aside>

-   The names of all resources must be ‘nouns’.

-   The names of the operations on each of the resources must be ‘verbs’.


## Plugin authentication

The plugin uses account credentials previously created in the aggregator Account Service to obtain a token from the Session Service to authenticate itself with the aggregator. For ongoing communication, the plugin and the aggregator use credentials of the resource aggregator account that is configured in the plugin during installation.

To authenticate the requestor and authorize the request to plugin, implement the following authentication methods:

-   **HTTP BASIC authentication (BasicAuth)** 

    Basic authentication is a simple authentication scheme built into the HTTP protocol. The client sends HTTP requests with the Authorization header that contains the word "Basic" followed by a space and a base64-encoded string of username: password. For example, to authorize as demo/p@55w0rd, the Authorization header to be used is:

    `Authorization: Basic ZGVtbzpwQDU1dzByZA==` 

-   **Session login authentication (XAuthToken)** 

    Session-based authentication allows users to obtain a token by entering their username and password which allows them to fetch a specific resource—without using their username and password.


## Plugin APIs

The plugin layer uses an API framework to perform control-related tasks such as adding a resource, modifying a session, and providing a heartbeat status response to the aggregator.

Each plugin is typically started as a system service and hence requires an explicit method of authentication of the resource aggregator during start-up. The plugin layer relies on the aggregator credentials over an HTTPS connection to authenticate the aggregator. The plugin layer uses the session service to authenticate itself to the aggregator. The current implementation defines both a specific token-based authentication and basic authentication mechanisms for plugin/aggregator security.

The plugin layer’s primary role is to act as a translator between the aggregator and the resource. It receives aggregator-native messages and translates them to resource-native messages one way and then receives resource-native messages that must be translated to aggregator-native messages in the other.

In the context of ODIM™, the aggregator receives Redfish messages from its north-bound clients that it passes as a payload to the plugin’s API server. The plugin translates the payload to a resource-native mechanism and performs an operation on the resource. Similarly, when a plugin receives a response from the resource in the resource-native format, the plugin responds to the aggregator with a Redfish payload put in an ODIM API response. By providing a common set of endpoints that do not pertain to either the north-bound or south-bound protocols, the plugin layer is able to perform tasks on any identified protocol without modifying its existing interfaces.


### Certificate TLS communication

Certificate is required for enabling secure communication from and to a plugin for the following scenarios:

-   where a plugin acts as server:

-   ODIM connects to the plugin to perform control operation.

-   Device to plugin \(for event listener\).

-   When plugin acts as a client to BMC or devices supporting HTTPS.

<aside class="notice">
NOTE: Recommended TLS version is 1.2.
</aside>

### Mandatory and optional functions implemented by plugins

The plugin layer forms the southern-most boundary of Resource Aggregator for ODIM's architecture. Plugins are started by the aggregator and they use HTTPS encrypted communication for security reasons.

The plugin's primary responsibility is to interface with the resource on behalf of the aggregator. There are two key components to any plugin:

1.  Control data

    Control data describes any messages sent by the administrator to enact on a certain resource. This includes tasks such as adding a resource, discovering a resource, setting up events, and retrieving resource information among others. Control data exchange is synchronous and is initiated by the administrator or by another entity through the north-bound APIs of the aggregator.

2.  Event data

    Event data describes any messages sent by the resource based on a certain previous event notification request. Resource-specific events such as component failures, telemetry, and log files among others are examples of Event data. Event data exchange is asynchronous \(a push operation\) and is initiated by the resource. This communication happens from the plugin to the aggregator only and not the opposite way.

### Plugin API service

#### Mandatory functions

Each plugin implements API services conforming to specific standards targeted at addressing control data transfer.

1.  Plugin Life cycle
    -   Start-up handshake

        The start-up handshake occurs after the plugin's service has been started up by the aggregator based on incoming request. The start-up handshake method exchanges state information from the aggregator to the plugin which contains information on currently managed resources, their last-known configuration, and a list of events subscriptions for compute plugins. Fabric plugins which maintain their own state need not use this data.

    -   Status

        The status method provides a way for the aggregator to verify if the service is still up by providing, at a minimal level, a heartbeat response while having the option to be able to provide any other relevant information on the plugin's status.

2.  Action on Resources and Collections

    Servers are an ideal example for action on resources. As an example, when the aggregator is about to take action on servers, it is typically performed on the actual resource instead of the entire collection since the collection has been updated by the resident Redfish implementation on the iLO.

    Switches are an ideal example for action on collections. As an example, when a new switch is added to the fabric, the fabric manager has to take an action on the fabric collection by adding a new resource.

    -   Discover

        Discover tasks the plugin to look for resources that fit a certain profile and report back to the aggregator. The aggregator might then choose to take an action on the returned resources or collections. The Discover method is typically a trigger from the aggregator initiated locally by the aggregator or another north-bound entity to perform the task. Auto-discovery of resources is a desirable feature, but an optional feature at this time.

        Lookup of resources, however, is not a plugin task and is implemented in the aggregator. The aggregator gets all the necessary information from the plugin and performs the filtering itself thereby reducing the need for additional processing on the plugin front keeping the tasks simpler.

    -   Add

        Add tasks the plugin to add a certain resource or collection to the aggregator as a managed resource.

    -   Remove

        Remove results in a certain resource or collection removed from the list of managed resources.

    -   Verify

        Verify checks if this resource or collection can currently be managed by the plugin. This usually results in the plugin performing an idempotent task on the resource to determine what version of the specification the resource implements and if the plugin can manage that resource.

    -   Configure - resource specific

        Configure is a suite of API calls that allow the aggregator to configure the underlying resource. This varies from one plugin to another and will be described in detail for each plugin.

    -   Subscribe - resource specific

        Subscribe is a suite of API calls that sets up event notifications of a certain type for resources. Subscriptions and event notifications vary from one plugin to another and will be described in detail for each plugin. Plugins that work with resources that are fully Redfish compliant and do not expect a significantly high number of event notifications can therefore enable just the API-based mechanism for requests and responses.

    -   Request status

        Plugins can optionally support a method that responds with a report on the status of an ongoing job. In situations when an immediate response cannot be relayed back, the plugin can send a `202` message indicating that it has accepted a request. This response can optionally include an endpoint that the aggregator can query for status of this ongoing request.



#### Optional functions

1.  Redfish OEM object processing

    Resource characteristics that are currently not covered in a Redfish approved schema \(either approved spec or draft spec\) are typically added to the OEM block within each resource type. A plugin uses the OEM object to receive and respond to changes that are unique to the resource.

2.  Redfish Draft Schemas

    Support for schema elements that are currently in draft status is an optional method that plugins can support. This allows an implementation of Resource Aggregator for ODIM to validate newer schemas that are relevant to the customer's use-case. A corresponding optional element to process the new schema might be provided by the aggregator for this draft schema. As an example, Redfish has a draft schema to enable `syslog` methods to be sent to a collector and this can be enabled by a server plugin. An aggregator may choose to recognize draft schemas.

3.  Auto Discovery

    Plugins can use well-defined auto discovery mechanisms, such as SSDP, to detect and report events on finding new resources added to the deployment.


### Message bus services

#### Mandatory functions

 

1.  Events on resources:

    The aggregator should be able to post a message requesting for a subscription to certain events on a resource. The type of events varies based on the resource in question and will be described in detail for each individual plugin.

2.  Event synchronizer:

    If the resource does not support a Redfish-aligned, REST-based event notification system, the plugin implements an event-oriented synchronizer that receives the event notifications in their native format \(SNMP, logs, and so on\) and responds to the aggregator or an external entity with compliant message types.


#### Optional functions

1.  Telemetry on resources:

    The aggregator can request to collect resource telemetry information that will be set up by the plugins. The type of telemetry information varies from plugin to plugin.

2.  Multicast status:

    The aggregator can send a multicast request to all plugins requesting status information for all active plugins to determine the health of the system.

<aside class="notice">
NOTE: All events are sent to the aggregator over the message bus.
</aside>


## Plugin message bus

Plugin message bus is a mandatory function.

### Message payloads

The current version of ODIM™ uses the JSON schema of the resource as the primary payload across the platform. Events originating from a southbound resource reaches the listener process on the plugin. The plugin then translates the notification into a Redfish Event schema in a JSON object model and also adds server address to event structure JSON before publishing it on the EMB. The listener on the aggregator layer receives this JSON object before sending it to the appropriate northbound listener who originally subscribed to this notification.

Event structure to post on message bus:

```
{
"ip" : <device_ip>,
"request" : <UTF-8 encoded event data>
}
```

### Encoding and decoding

The current version of ODIM™ uses basic Unicode `UTF-8` encoding and decoding between producers and consumers. The `UTF-8` encoding does not provide significant performance advantages as other schemes. It has been chosen to provide the easiest entry point for various plugins and their associated resources to be part of the ODIM ecosystem.



## Plugin API server information

This section describes the server information for the entity hosting the API server. The server is hosted at:

`https://<IPAddressOfServer><:Port>/ODIM/v1/`.

Each plugin adds its name and the URI for its resources under the root described above.

**Resources** 

The resources section lists all the resources currently supported by the plugin. If an action is translatable as a CRUD function, it resides in the resources section with the necessary fields in the schema updated for that action.

**Functions** 

The functions section lists the supported functions which are implemented as an action, their parameters, the format of the request body, if required and the response for each one of them.

For the ODIM™ project, "ODIM" serves as the Service Root. Hence, the API server is at `https://<IPAddressOfServer><:Port>/ODIM/v1/` on all instances that host the plugin. In scenarios where multiple plugins must run on the same server each plugin uses a unique port number. Each plugin binds itself to a port on the server that is specified in the configuration file of the plugin. This information enables the plugin to bind to the right port on start-up. It also ensures that the aggregator can connect to the right plugin using the right port.

Each API service is hosted under HTTPS to ensure secure access to the resources. Additional protocols beyond HTTPS may be implemented in the future. Resources and corresponding operations are available for each plugin in two broad categories. Some resources are common for most plugins while others are specialized resources that are unique to a certain plugin. Similarly, resources might have generalized operations that are common for all plugins while there are other specialized operations that are unique to certain resources.

The following table contains different types of operations:

### Resource related APIs

#### Resource: /ODIM/v1/<PluginResource\>/ \(Optional\)


|||
|-------|--------|
|Operation| `GET {URI for ResourceID}` |
|Description|Gets a resource based on query string|
|Parameters|Response Body with resource information specific to plugin|
|Payload|Application/json|
|Response| `200 (Success) with parameters, 403 (Forbidden), 404 (Not Found)` |

|||
|-------|--------|
|Operation| `POST` |
|Description|Adds a new resource to the collection|
|Parameters|Request Body with desired resource state specific to plugin; Response body with status on success|
|Payload|Application/json|
|Response| `200 (Success), 403 (Forbidden), 404 (Not Found)` |


|||
|-------|--------|
|Operation| `DELETE {ResourceId}` |
|Description|Removes an existing resource from the collection|
|Parameters|None|
|Payload|None|
|Response| `200 (Success), 403 (Forbidden), 404 (Not Found)` |


### Subscription-related APIs

When adding a new subscription, follow these guidelines:

-   When creating a subscription, make plugin listener as the subscription destination on the resource.

-   Check if subscription is present with the plugin listener as the destination in the resource. Given it is present, check if it is matching the current request. If they do not match, remove the old subscription and create a subscription with the subscription request details.


#### Resource: /ODIM/v1/Subscription/ \(Mandatory\)


|||
|-------|--------|
|Operation| `GET {Resource Id}` |
|Description|Gets a current subscription based on a query string. If not, gets all subscriptions.|
|Parameters|Response body with object array including all sessions for event notifications|
|Payload|Application/json|
|Response| `200 (Success), 404 (Not Found)` |

|||
|-------|--------|
|Operation| `POST` |
|Description|Adds a new subscription|
|Parameters|Request body with new session information|
|Payload|Application/json|
|Response| `200 (Success), 202 (Accepted), 403 (Forbidden)` |


|||
|-------|--------|
|Operation| `DELETE {Resource Id}` |
|Description|Removes an existing subscription|
|Parameters|None|
|Payload|None|
|Response| `200 (Success), 403 (Forbidden), 404 (Not Found)` |


### Plugin control APIs

#### Resource: /ODIM/v1/validate/ \(Mandatory\)


|||
|-------|--------|
|Operation| `POST` |
|URI| `/ODIM/v1/validate/` |
|Description|Check the server credentials.|
|Response Code| `200 (Success), 401 (Unauthorized)` |
|Authentication|Yes|
|Payload|Application/json|

**Request Body** 
```
{
   "ManagerAddress":"<hostaddress>",
   "UserName":"<user_name>",
   "Password":"<password>"
}
```



**Response** 
```
{
   "ServerIP":"<server_IP>",
   "Username":"<user_name>",
   "device_UUID":"<device uuid>"
}
```



#### Resource: /ODIM/v1/Sessions/ \(Mandatory\)


|||
|-------|--------|
|Operation| `POST` |
|URI| `/ODIM/v1/Sessions/` |
|Description|Create a session on the plugin.|
|Payload|Application/json|
|Response Header|X-Auth-Token: 15d0f639-f394-4be7-a8ef-ef9d1df07288|
|Response Code| `201 (Created), 401 (unauthorized), 400 (Bad Request)` |


**Request Body** 
```
{
   "UserName":"admin",
   "Password":"admin"
}
```



#### Resource: /ODIM/v1/Status/ \(Mandatory\)


|||
|-------|--------|
|Operation| `GET` |
|URI| `/ODIM/v1/Status/` |
|Description|Gets a representation of the status of the plugin and associated information such as time alive or on a pending request.|
|Response Code| `200 (Success), 401 (unauthorized)` |
|Authentication|Yes|

**Response** 
```
{
   "_comment":"Plugin Status Response",
   "Name":"Common Lenovo Plugin Status",
   "Version":"v0.1",
   "Status":{
      "Available":"yes",
      "Uptime":"2020-06-18T00:10:07-06:00",
      "TimeStamp":"2020-06-22T03:29:37-06:00"
   },
   "EventMessageBus":{
      "EmbType":"Kafka",
      "EmbQueue":[
         {
            "EmbQueueName":"REDFISH-EVENTS-TOPIC",
            "EmbQueueDesc":"Queue for redfish events"
         }
      ]
   }
}
```


#### Resource: /ODIM/v1/Startup/ \(Mandatory\)


|||
|-------|--------|
|Operation| `POST` |
|URI| `/ODIM/v1/Startup/` |
|Description| Posts a new representation to the StartUp resource. `POST` allows the system to keep track of state information sent through each start-up and potential rollbacks of plugin.<br> The value given in the Location parameter for each resource should be used to do a `GET` request to verify if the subscription is present. If the location is not present in the device or the subscription details are different compared to that specified in this request, we must delete the subscription and resubscribe with new details. The attributes in the following body is for Redfish-based subscriptions. Other device types like Fabric and Storage schema will have other attributes. Also keep in mind the differences in implementation across device types and the optional/mandatory attributes.<br> If the subscription is altered, the ID in the URI sent back will be updated accordingly.<br> |
|Payload|Application/json|
|Response Code| `200 (Success), 401 (unauthorized)` |
|Authentication|Yes|


**Request Body** 
```
[
   {
      "Location":"https://<hostaddress>/redfish/v1/EventService/Subscriptions/1",
      "EventTypes":[
         "Alert",
         "StatusChange"
      ],
      "MessageIds":[

      ],
      "OriginResources":[

      ],
      "RegistryPrefixes":[

      ],
      "ResourceTypes":[

      ],
      "SubordinateResources":[

      ]      "Device":{
         "ManagerAddress":"<hostaddress>",
         "UserName":"admin",
         "Password":"admin"
      }
   }
]
```



**Response** 
```
{
   "<hostaddress>":"https://<hostaddress>/redfish/v1/EventService/Subscriptions/2"
}
```


#### Resource: /ODIM/v1/Managers \(Mandatory\)


|||
|-------|--------|
|Operation| `GET` |
|URI| `/ODIM/v1/Managers/` |
|Description|Get on managers.|
|Payload|Application/json|
|Response Code| `200 (Success), 401 (unauthorized)` |
|Authentication|Yes|

**Response** 
```
{
   "@odata.context":"/ODIM/v1/$metadata#ManagerCollection.ManagerCollection",
   "@odata.etag":"W/\"AA6D42B0\"",
   "@odata.id":"/ODIM/v1/Managers",
   "@odata.type":"#ManagerCollection.ManagerCollection",
   "Name":"Managers",
   "Members":{
      "@odata.id":"/ODIM/v1/Managers/<uuid>"
   },
   "Description":"Manager collection",
   "Members@odata.count":1
}
```

|||
|-------|--------|
|Operation| `GET` |
|URI| `/ODIM/v1/Managers/<ManagerID>/` |
|Description|Get on manager.|
|Payload|Application/json|
|Response Code| `200 (Success), 401 (unauthorized)` |
|Authentication|Yes|

**Response** 
```
{
   "@odata.context":"/ODIM/v1/$metadata#Manager.Manager",
   "@odata.etag":"W/\"AA6D42B0\"",
   "@odata.id":"/ODIM/v1/Managers/<Id>",
   "@odata.type":"#Manager.v1_3_3.Manager",
   "Name":"<plugin name>",
   "ManagerType":"Service",
   "Id":"<ManagerID>",
   "UUID":"<uuid>",
   "FirmwareVersion":1,
   "Status":{
      "State":"Enabled",
      "Health":"OK"
   }
}
```








## The plugin service details

The Plugin service is an in-memory process started as a docker instance as part of the overall host start-up process. This service hosts the API server, event synchronizer, load balancers, worker threads, EMB publishers and, subscribers among other entities as the implementation decides.

The plugin service has the capability to schedule individual, short-lived instances that perform specific functionality as required by the north-bound entity. The current plugin service is hosted centrally. In future, it may be possible to deploy individual instances of the plugin service across distributed sites thereby allowing the plugin layer to scale.

Information on parameters needed by the plugin service on start-up are available from the plugin configuration file for each instance. The files are configured as JSON files for the aggregator to read and take action. It includes information about:

-   Credentials (password is hashed)

-   Plugin ID

-   Firmware version host IP address and port

-   TLS configuration that specifies TLS version and cipher suites to be used.

-   Certificate paths

-   Message bus configuration required to publish events

-   Session timeout configurations

-   Rules for converting south bound messages to ODIM format \(optionally\).


## Deployment guidelines

The plugin layer, and all the components of ODIM™ are built in a deployment-agnostic manner. ODIM™ and all its components, including the plugin layer, must be deployed in an environment regardless of any underlying virtualization mechanism – KVM, ESX, or containers or its absence thereof. Any plugin expecting a tighter dependency on the underlying infrastructure must identify it as part of its specifications. The open-source version has packages to run on Docker. This is not a mandatory requirement. Individual projects may choose to have their own deployment platforms and strategies.

The plugin component must provide access to the source code and build/deploy instructions on the source repository publically hosted on GitHub. Deployers can use this information to deploy ODIM™ and its components within their existing framework as a virtual machine or a container or a bare metal service. The individual services required by a plugin \(For example, API server\) are part of the build instructions provided by the plugin layer.

As an example, for a containerized version, the deployment looks as follows:

Three artifacts are required by the deployment tool. A GitHub repository (owned by the plugin developers), a Docker file.

-   The plugin code and artifacts are available from GitHub. The Docker file indicates how to build a containerized image of the plugin and its associated processes.

-   The deployer provides an operating environment that uses Docker Containers for virtualization.

-   The repository provides scripts and documentation to deploy the Docker images.



## Pseudo code for API implementation

<aside class="notice">
NOTE:

-   Translate protocol name in the URL from ODIM to the device protocol. The device might support Redfish or any other protocol. Plugin will do the translation required when you do any operation on the resource. Example:

    `/ODIM/v1/Systems to /redfish/v1/Systems/` 

-   When any request comes to plugin, do the following:

         -  Check if header has auth token or basic auth. If the header does not have auth token or basic auth, return an unauthorized error.

         -  Check if auth token is valid or basic auth has valid user name and password. If auth token is not valid or basic auth does not have valid user name and password, return an unauthorized error.


</aside>

**The main function pseudo code**

```
main () {
  Check if plugin is not running as root user
  Else 
        Log an error and exit
  
  Initialize the endpoints in the router
 }

```

Following are the endpoints:

###  Session

**URL:** `/ODIM/v1/Sessions/` 

**Method:** `POST` 

**Pseudo code:** 

```
Func CreateSession (context) {
	  Read user_name and password from the request context
	  Validate the credentials with the configured values
	  If success return newly created token in the header and status code as StatusCreated.
	  Else return an unauthorized error
  }
  
  URL: /ODIM/v1/validate/
  Endpoint: POST
  Pseudo Code:
  Func Validate (context) {
	Read input json from request context 
	Verify the given credential from the input request by invoking GET on /redfish/v1/Systems
	If success return json object {ServerIP : <server_IP>, Username : <user_name>, device_UUID : <device uuid>}
    Else return response which comes from the resource
  }

```

### Status

**URL:** `/ODIM/v1/Status/` 

**Method:** `GET` 

**Pseudo code:** 

```
Func GetPluginStatus (context) {
    Build a response of status of the plugin and message bus type and topic name.
	Return a response
  }

```


### Startup

**URL:** `/ODIM/v1/Startup/` 

**Method:** `POST` 

**Pseudo code:** 

```
  Func GetPluginStartup (context) {
    Read input json from request context (input json will have collection of resources to check and subscribe for requested event types)
	For each resource 
	    GET on the subscription location
	    If it's success then
		   Delete subscription from the resource
		   Re-subscribe with the event types provided in the request parameter "EventTypes"
        Else 
           Subscribe with the event types provided in the request parameter "EventTypes"
    Return response containing key value pair of server_address and subscription location		   
  }

```


### Subscriptions

**URL:**`/ODIM/v1/Subscriptions/` 

**Method:** `POST` 

**Pseudo code:** 

```
Func CreateEventSubscription (context) { 
	Read input json from request context
	DeleteMatchingSubscription on the resource 
		GET on /redfish/v1/EventService/Subscriptions on the resource
		Check each subscription, if the plugin listener as the destination in the resource
		If it is present then 
		   Delete Subscription from the resource with subscription location
	Make plugin listener as the subscription destination in the subscription request
	Subscribe Events on the resource with the subscription request
	Return response having subscription location in the response header
  }

```

 

**URL:** `/ODIM/v1/Subscriptions/` 

**Method:** `GET` 

**Pseudo code:** 

```
Func GetEventSubscription (context) { 
	Read input json from request context
	GET on subscription uri on the resource
	If it is GET on collection subscription then return collection subscription
    Else return a requested subscription details
  }

```

 

**URL:**`/ODIM/v1/Subscriptions/` 

**Method:** `DELETE` 

**Pseudo code:** 

```
  Func DeleteEventSubscription (context) { 
	Read input json from request context
	Delete subscription from the requested subscription location
	Return a resource response to ODIM
  }

```


### Managers

**URL:** `/ODIM/v1/Managers/` 

**Method:** `GET` 

**Pseudo code:** 

```
Func GetManagerCollection (context) {
    Check if ManagerAddress in the request is empty
	If it is empty then build manager collection response having plugin root service id as manager id in the @odata.id
		Return manager collection response 
	Else 
	   GET on managers on the resource
	   Return response comes from the resource
  }

```

 

**URL:** `/ODIM/v1/Managers/<manager_id>/` 

**Method:**`GET` 

**Pseudo code:** 

```
Func GetManagerCollection (context) {
    Check if ManagerAddress in the request is empty
	If it is empty then build manager response 
		Return manager response 
	Else 
	   GET on manager on the resource
	   Return response comes from the resource
  }   

```



### Systems

**URL:** `/ODIM/v1/Systems/` 


<aside class="notice">
NOTE: All the system and its child URIs use same pseudo code.
</aside>

**Method:** `GET` 

**Pseudo code:** 

```
  Func GetResource (context) {
    Read input json from request context
	GET on the request uri on the resource
	Return a response comes from resource
  }
}

```

