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

##  Modifying Configurations of Events Service
  
Config file of ODIMRA is located at: **odimra/lib-utilities/config/odimra_config.json**  
Refer to the section **Modifying Configurations** in the README.md file to change the configurations of an odimra service.
  
**Specific configurations for Events Service are:**
  
##  Log Location of the Events Service
  
/var/log/ODIMRA/event.log
  
  



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
   "@odata.type":"#EventService.v1_7_0.EventService",
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











## Creating an event subscription

|||
|-----------|-----------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/EventService/Subscriptions` |
|**Description**| This endpoint subscribes a northbound client to events originating from a set of resources \(southbound devices, managers, Resource Aggregator for ODIM itself\) by creating a subscription entry. For use cases, see [Subscription use cases](#event-subscription-use-cases).<br>This operation is performed in the background as a Redfish task. If there is more than one resource that is sending a specific event, the task is further divided into subasks.|
|**Returns** |<ul><li>`Location` URI of the task monitor associated with this operation in the response header.</li><li> Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See "Sample response body \(HTTP 202 status\)".<br>**IMPORTANT:**<br> Note down the task Id. If the task completes with an error, it is required to know which subtask has failed. To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`.</li><li>On success, a `Location` header that contains a link to the newly created subscription and a message in the JSON response body saying that the subscription is created. See "Sample response body \(HTTP 201 status\)".</li></ul>|
|**Response code** |<ul><li>`202 Accepted`</li><li>`201 Created`</li></ul>|
|**Authentication** |Yes|


**Usage**


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
|/redfish/v1/TaskService/Tasks|All tasks scheduled by or being executed by Redfish `TaskService`. By subscribing to Redfish tasks, you can receive task status change notifications on the subscribed destination client.<br> By specifying the task URIs as `OriginResources` and `EventTypes` as `StatusChange`, you can receive notifications automatically when the tasks are complete.<br> To check the status of a specific task manually, perform HTTP `GET` on its task monitor until the task is complete.<br> |
|/redfish/v1/Managers|All manager resources available in Resource Aggregator for ODIM for which the service sends only related events.|


**Event types**

|String|Description|
|------|-----------|
|Alert|A condition exists which requires attention.|
|ResourceAdded|A resource has been added.|
|ResourceRemoved|A resource has been removed.|
|ResourceUpdated|The value of this resource has been updated.|
|StatusChange|The status of this resource has changed.|

**EventFormat type**

|String|Description|
|------|-----------|
|Event|The subscription destination will receive JSON bodies of the Resource Type Event.|

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
   "@odata.type":"#Task.v1_5_1.Task",
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
   "@odata.type":"#Event.v1_6_1.Event ",
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
      "ResourceEvent.1.0.3.ResourceAdded"
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
      "ResourceEvent.1.0.3.ResourceRemoved"
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
      "ResourceEvent.1.0.3.StatusChange"
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
   "@odata.type":"#EventDestination.v1_10_1.EventDestination",
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
   "@odata.type":"#EventDestination.v1_10_1.EventDestination",
   "@odata.id":"/redfish/v1/EventService/Subscriptions/57e22fcc-8b1a-460c-ac1f-b3377e22f1cf",
   "Id":"57e22fcc-8b1a-460c-ac1f-b3377e22f1cf",
   "Name":"Event Subscription",
   "Message":"The resource has been removed successfully.",
   "MessageId":"ResourceEvent.1.0.3.ResourceRemoved",
   "Severity":"OK"
}
```