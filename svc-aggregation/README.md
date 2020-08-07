#  Resource aggregation and management

Resource Aggregator for ODIM allows users to add and group southbound infrastructure into collections for easy management. It exposes Redfish `AggregationService` endpoints to achieve the following:

-   Adding a resource and building its inventory.

-   Resetting one or more resources.

-   Changing the boot path of one or more resources to default settings.

-   Removing a resource from the inventory which is no longer managed.


Using these endpoints, you can add or remove only one resource at a time. You can group the resources into one collection and perform actions in combination on that group.

All aggregation actions are performed as [tasks](#tasks) in Resource Aggregator for ODIM. The actions performed on a group of resources \(resetting or changing the boot order to default settings\) are carried out as a set of subtasks.

<aside class="notice">
To access Redfish `AggregationService` endpoints, you require `ConfigureComponents` privilege. If you access these endpoints without necessary privileges, you will receive an HTTP `403` error.
</aside>
  
  
##  Modifying Configurations of Aggregation Service
  
Config File of ODIMRA is located at: **odimra/lib-utilities/config/odimra_config.json**  
Refer the section **Modifying Configurations** in the README.md to change the configurations of a odimra service
  
**Specific configurations for Aggregation Service are:**
  
  
  
  
##  Log Location of the Aggregation Service
  
/var/log/ODIMRA/aggregation.log
  
  
##  Supported endpoints

|||
|-------|--------------------|
|/redfish/v1/AggregationService|`GET`|
|/redfish/v1/AggregationService/Actions|`GET`|
|/redfish/v1/AggregationService/Actions/AggregationService.Add|`POST`|
|/redfish/v1/AggregationService/Actions/AggregationService.Remove|`POST`|
|/redfish/v1/AggregationService/Actions/AggregationService.Reset|`POST`|
|/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder|`POST`|


## The aggregation service root


```
curl -i --insecure -X GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/AggregationService'


```

> Sample response header 

```
allow:GET
cache-control:no-cache
content-length:1 kilobytes
content-type:application/json; charset=utf-8
date:Mon, 18 Nov 2019 14:52:21 GMT+5h 37m
etag:W/"12345"
link:/v1/SchemaStore/en/AggregationService.json>; rel=describedby
odata-version:4.0
status:200
x-frame-options:sameorigin

```

> Sample response body 

```
{ 
   "@odata.context":"/redfish/v1/$metadata#AggregationService.AggregationService",
   "@odata.etag":"W/\"979B45E7\"",
   "Id":"AggregationService",
   "@odata.id":"/redfish/v1/AggregationService",
   "@odata.type":"#AggregationService.v1_0_0.AggregationService",
   "Name":"AggregationService",
   "Description":"AggregationService",
   "Actions":{ 
      "#AggregationService.Add":{ 
         "target":"/redfish/v1/AggregationService/Actions/AggregationService.Add/",
         "@Redfish.ActionInfo":"/redfish/v1/AggregationService/AddActionInfo"
      },
      "#AggregationService.Remove":{ 
         "target":"/redfish/v1/AggregationService/Actions/AggregationService.Remove/",
         "@Redfish.ActionInfo":"/redfish/v1/AggregationService/RemoveActionInfo"
      },
      "#AggregationService.Reset":{ 
         "target":"/redfish/v1/AggregationService/Actions/AggregationService.Reset/",
         "@Redfish.ActionInfo":"/redfish/v1/AggregationService/ResetActionInfo"
      },
      "#AggregationService.SetDefaultBootOrder":{ 
         "target":"/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder/",
         "@Redfish.ActionInfo":"/redfish/v1/AggregationService/SetDefaultBootOrderActionInfo"
      }
   },
   "ServiceEnabled":true,
   "Status":{ 
      "Health":"OK",
      "HealthRollup":"OK",
      "State":"Enabled"
   }
}
```



|**Method**|`GET` |
|-----|-------------------|
|**URI** |`redfish/v1/AggregationService` |
|**Description** |The URI for the Aggregation service root.|
|**Returns** |Properties for the service and a list of actions you can perform using this service.|
|**Response Code** |`200` on success|

 

 

## Adding a plugin

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{"ManagerAddress":"{BMC_address}:{port}",
  "UserName":"{plugin_userName}",
  "Password":"{plugin_password}",
  "Oem":{"PluginID":"{Redfish_PluginId}",
         "PreferredAuthType":"{Preferred_aunthentication_type}",
         "PluginType":"{plugin_type}"
        }
 }' \
 'https://{odimra_host}:{port}/redfish/v1/AggregationService/Actions/AggregationService.Add'


```

|||
|---------|---------------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/AggregationService/Actions/AggregationService.Add` |
|**Returns** |<ul><li>`Location` URI of the task monitor associated with this operation in the response header.</li><li>Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See "Sample response body \(HTTP 202 status\)".</li><li>When the task is successfully complete, success message in the JSON response body.</li></ul>|
|**Response Code** |<ul><li>`202 Accepted`</li><li>`200 Ok`</li></ul>|
|**Authentication** |Yes|

<br>
**Description**

This action discovers information about a plugin and adds it in the inventory. 



It is performed as a task. To know the progress of this action, perform `GET` on the [task monitor](#task-monitor) returned in the response header \(until the task is complete\).


 

> Sample request body

```
{ 
   ​   "ManagerAddress":"{BMC_address}:45001",
   ​   "UserName":"abc",
   ​   "Password":"abc123",
   ​   "Oem":{ 
      ​      "PluginID":"GRF",
      ​      "PreferredAuthType":"BasicAuth(or)XAuthToken",
      ​      "PluginType":"Compute"      ​
   }   ​
}
```

###  Request parameters

|Parameter|Type|Description|
|---------|----|-----------|
|ManagerAddress|String \(required\)<br> |A valid IP address or hostname and port of the baseboard management controller \(BMC\) where the plugin is installed. The default port for the Generic Redfish Plugin is `45001`.<br>**NOTE:**<br> Ensure that the port is greater than `45000`.|
|UserName|String \(required\)<br> |The plugin username. Example: UserName for the Generic Redfish Plugin - Default admin username is `admin`|
|Password|String \(required\)<br> |The plugin password. Example: Password for the Generic Redfish Plugin - Default admin password is `GRFPlug!n12$4`|
|PluginID|String \(required\)<br> |The id of the plugin you want to add. Example: GRF|
|PreferredAuthType|String \(required\)<br> |Preferred authentication method to connect to the plugin - `BasicAuth` or `XAuthToken`.|
|PluginType|String \(required\)<br> |The string that represents the type of the plugin. For the Generic Redfish Plugin, the type is `Compute`.|


> Sample response header \(HTTP 202 status\)

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/task85de4003-8757-4c7d-942f-55eaf7d6812a
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes

```

> Sample response body \(HTTP 202 status\)

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

> Sample response body \(HTTP 200 status\)

```
{
   "code":"Base.1.6.1.Success",
   "message":"Request completed successfully."
} 
```








## Adding a server

```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{"ManagerAddress":"{BMC_address}","UserName":"{BMC_username}","Password":"{BMC_password}","Oem":{"PluginID":"{Redfish_PluginId}"}}' \
 'https://{odimra_host}:{port}/redfish/v1/AggregationService/Actions/AggregationService.Add'


```



|||
|---------|---------------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/AggregationService/Actions/AggregationService.Add` |
|**Returns** | <ul><li>`Location` URI of the task monitor associated with this operation in the response header.</li><li>Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See "Sample response body \(HTTP 202 status\)".</li><li> On success:<ul><li>A success message in the JSON response body.</li><li>A link \(having `ComputerSystemId` \(`Id`\) specified by Resource Aggregator for ODIM\) to the added server in the `Location` header.</li></ul></li></ul>|
|**Response code** |<ul><li>`202 Accepted`</li><li>`200 Ok`</li></ul>|
|**Authentication** |Yes|

**Description**

This action discovers information about a single server and performs a detailed inventory of it. 


It is performed as a task. To know the progress of this action, perform `GET` on the [task monitor](#task-monitor) returned in the response header \(until the task is complete\). When the task is successfully complete, you will receive `ComputerSystemId` of the added server.


**IMPORTANT:**

 `ComputerSystemId`(`Id`) is unique information about the server that is added. Save it as it is required to perform subsequent actions such as `delete`, `reset`, and `setdefaultbootorder` on this server. 
 
 You can get the `ComputerSystemId`(`Id`) of a specific server later by performing HTTP `GET` on `/redfish/v1/Systems`. See [collection of Computer Systems](#collection-of-computer-systems).


 

> Sample request body

```
{
	"ManagerAddress": "{BMC_address}",
	"UserName": "abc",
	"Password": "abc1234",
	"Oem": {
		"PluginID": "GRF"
	}
}
```

### Request parameters

|Parameter|Type|Description|
|---------|----|-----------|
|ManagerAddress|String \(required\)<br> |A valid IP address or hostname of a baseboard management controller \(BMC\).|
|UserName|String \(required\)<br> |The username of the server BMC administrator account.|
|Password|String \(required\)<br> |The password of the server BMC administrator account.|
|PluginID|String \(required\)<br> |The plugin id of the plugin. Example: "GRF"<br>**NOTE:**<br> Before specifying the plugin Id, ensure that the installed plugin is added in the Resource Aggregator for ODIM inventory. To know how to add a plugin, see [Adding a Plugin](#adding-a-plugin).|

> Sample response header \(HTTP 202 status\)

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/task4aac9e1e-df58-4fff-b781-52373fcb5699
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes

```

> Sample response header \(HTTP 200 status\)

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/redfish/v1/Systems/d8f740ba-5f01-4784-89b3-122fe76af739:1
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Mon,18 May 2020 09:16:09 GMT+5m 15s
Content-Length:73 bytes
```

> Sample response body \(HTTP 202 status\)

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

> Sample response body \(HTTP 200 status\)

```
{
   "code":"Base.1.6.1.Success",
   "message":"Request completed successfully."
}
```



## Resetting servers

```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{
"parameters": {
"ResetCollection": {
"description": "Collection of ResetTargets",
"ResetTarget": [{
"ResetType": "ForceRestart",
"TargetUri": "/redfish/v1/Systems/{ComputerSystemId}",
"Priority": 10,
"Delay": 5
},
{
"ResetType": "ForceOff",
"TargetUri": "/redfish/v1/Systems/{ComputerSystemId}",
"Priority": 9,
"Delay": 0
}
]
}
}
}' \
 'https://{odimra_host}:{port}/redfish/v1/AggregationService/Actions/AggregationService.Reset'


```


|||
|---------|---------------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/AggregationService/Actions/AggregationService.Reset` |
|**Returns** |<ul><li> `Location` URI of the task monitor associated with this operation \(task\) in the response header.</li><li>Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI.<br>**IMPORTANT:**<br>Note down the task Id. If the task completes with an error, it is required to know which subtask has failed. To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`.</li><li>On successful completion of the reset operation, a message in the response body, saying that the reset operation is completed successfully. See "Sample response body \(HTTP 200 status \)".</li></ul>|
|**Response code** |<ul><li>`202 Accepted`</li><li>`200 Ok`</li></ul> |
|**Authentication** |Yes|

**Description**

This action shuts down, powers up, and restarts one or more servers. 
It is performed as a task and is further divided into subtasks to reset each server individually. To know the progress of this action, perform HTTP `GET` on the [task monitor](#task-monitor) returned in the response header \(until the task is complete\).

To get the list of subtask URIs, perform HTTP `GET` on the task URI returned in the JSON response body. See "Sample response body \(HTTP 202 status \)". The JSON response body of each subtask contains a link to the task monitor associated with it. To know the progress of the reset operation \(subtask\) on a specific server, perform HTTP `GET` on the task monitor associated with the respective subtask.



You can perform reset on a group of servers by specifying multiple target URIs in the request.




> Sample request body

```
{
   "parameters":{
      "ResetCollection":{
         "description":"Collection of ResetTargets",
         "ResetTarget":[
            {
               "ResetType":"ForceRestart",
               "TargetUri":"/redfish/v1/Systems/97d08f36-17f5-5918-8082-f5156618f58d:1",
               "Priority":10,
               "Delay":5
            },
            {
               "ResetType":"ForceOff",
               "TargetUri":"/redfish/v1/Systems/24b243cf-f1e3-5318-92d9-2d6737d6b0b9:1",
               "Priority":9,
               "Delay":0
            }
         ]
      }
   }
}

```

### Request parameters

|Parameter|Type|Description|
|---------|----|-----------|
|parameters\[\{|Array|The parameters associated with the reset Action.|
|ResetType|String \(required\)<br> |The type of reset to be performed. For possible values, see "Reset type". If the value is not supported by the target server machine, you will receive an HTTP `400 Bad Request` error.|
|TargetUri|String \(required\)<br> |The URI of the target for `Reset`. Example: `"/redfish/v1/Systems/{ComputerSystemId}"` |
|Priority|Integer\[0-10\] \(optional\)<br> |This parameter is used to indicate which reset action is performed on priority. You can set a priority number in the range of 0-10 \(zero being the least and 10 being the highest\). If this parameter is not specified, a priority of zero will be assigned by default.<br> |
|Delay\}\]|Integer \(seconds\)\[0-3600\] \(optional\)<br> |This parameter is used to defer the reset action by specified time. If this parameter is not specified, a delay of zero will be assigned by default.<br>**NOTE:**<ul><li> If two or more reset actions have equal priority values in a single request, they are performed one after the other in the order of their delay values \(a reset action with zero delay will be performed first\).</li><li>If two or more reset actions have equal priority and delay values, they are performed at the same time.</li></ul>|

### Reset type

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

> Sample response header \(HTTP 202 status\)

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/task85de4103-8757-4c7d-942f-55eaf7d6412a
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes

```

> Sample response body \(HTTP 202 status\)

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

> Sample response body \(subtask\)

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

> Sample response body \(HTTP 200 status\)

```
{ 
   "error":{ 
      "code":"Base.1.6.1.Success",
      "message":"Request completed successfully"
   }
}
```

## Changing the boot order of servers to default settings

```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{ 
   "parameters":{ 
      "ServerCollection":[ 
         "/redfish/v1/Systems/97d08f36-17f5-5918-8082-f5156618f58d:1",
         "/redfish/v1/Systems/76632110-1c75-5a86-9cc2-471325983653:1"
      ]
   }
}' \
 'https://{odimra_host}:{port}/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder'


```


|||
|---------|---------------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/AggregationService/Actions/AggregationService.SetDefaultBootOrder` |
|**Returns** |<ul><li>`Location` URI of the task monitor associated with this operation in the response header. </li><li>Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI.<br>**IMPORTANT:**<br> Note down the task Id. If the task completes with an error, it is required to know which subtask has failed. To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`.<br></li><li> On successful completion of this operation, a message in the response body, saying that the operation is completed successfully. See "Sample response body \(HTTP 200 status\)".</li></ul> |
|**Response code** |<ul><li>`202 Accepted`</li><li>`200 Ok`</li></ul>|
|**Authentication** |Yes|


<br>
**Description**

This action changes the boot order of one or more servers to default settings. 


It is performed as a task and is further divided into subtasks to change the boot order of each server individually. To know the progress of this action, perform HTTP `GET` on the [task monitor](#task-monitor) returned in the response header \(until the task is complete\).

To get the list of subtask URIs, perform HTTP `GET` on the task URI returned in the JSON response body. See "Sample response body \(HTTP 202 status\)". The JSON response body of each subtask contains a link to the task monitor associated with it. To know the progress of `SetDefaultBootOrder` action \(subtask\) on a specific server, perform HTTP `GET` on the task monitor associated with the respective subtask.




You can perform `setDefaultBootOrder` action on a group of servers by specifying multiple server URIs in the request.




> Sample request body

```
{ 
   "parameters":{ 
      "ServerCollection":[ 
         "/redfish/v1/Systems/97d08f36-17f5-5918-8082-f5156618f58d:1",
         "/redfish/v1/Systems/76632110-1c75-5a86-9cc2-471325983653:1"
      ]
   }
}

```

### Request parameters

|Parameter|Type|Description|
|---------|----|-----------|
|parameters\[\{|Array|The parameters associated with the `SetDefaultBootOrder` action.|
|ServerCollection\}\]| \(required\)<br> |Target servers for `SetDefaultBootOrder`.|


> Sample response header \(HTTP 202 status\)

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/task85de4003-8057-4c7d-942f-55eaf7d6412a
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes

```

> Sample response body \(HTTP 202 status\)

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

> Sample response body \(subtask\)

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

> Sample response body \(HTTP 200 status\)

```
{ 
   "error":{ 
      "code":"Base.1.6.1.Success",
      "message":"Request completed successfully"
   }
}
```


## Deleting the server inventory

```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{
"parameters": [
{
"Name": "/redfish/v1/Systems/{ComputerSystemId}"
}
]
}' \
 'https://{odimra_host}:{port}/redfish/v1/AggregationService/Actions/AggregationService.Remove/'


```


|||
|---------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/AggregationService/Actions/AggregationService.Remove` |
|**Returns** |<ul><li>`Location` URI of the task monitor associated with this operation in the response header.</li><li>Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See "Sample response body \(HTTP 202 status\)".</li><li>On successful completion, a message in the response body, saying that the operation is completed successfully. See "Sample response body \(HTTP 200 status\)".</li></ul> |
|**Response code** |<ul><li>`202 Accepted`</li><li>`200 Ok`</li></ul> |
|**Authentication** |Yes|


<br>
**Description**

This action removes the inventory of a specific server and deletes all associated event subscriptions. 


It is performed as a task. To know the progress of this action, perform `GET` on the [task monitor](#task-monitor) returned in the response header \(until the task is complete\).




> Sample request body

```
{
   "parameters":[
      {
         "Name":"/redfish/v1/Systems/97d08f36-17f5-5918-8082-f5156618f58d:1"
      }
   ]
}
```

###  Request parameters

|Parameter|Type|Description|
|---------|----|-----------|
|parameters\[\{|Array|The parameters associated with the `Delete` Action.|
|Name\}\]|String \(required\)<br> |The URI of the target to be removed: `/redfish/v1/Systems/{ComputerSystemId}` |



> Sample response header \(HTTP 202 status\)

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/task85de4003-8757-2c7d-942f-55eaf7d6412a
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes

```

> Sample response body \(HTTP 202 status\)

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

> Sample response body \(HTTP 200 status\)

```
{ 
   "error":{ 
      "code":"ResourceEvent.1.0.2.ResourceRemoved",
      "message":"Request completed successfully"
   }
}
```


## Removing a plugin

```
curl -i POST \
   -H "X-Auth-Token:{X-Auth-Token}" \
   -H "Content-Type:application/json" \
   -d \
'{
"parameters": [{
"Name": "/redfish/v1/Managers/{managerId}"
}
]
}' \
 'https://{odimra_host}:{port\}/redfish/v1/AggregationService/Actions/AggregationService.Remove'


```


|||
|---------|---------------|
|**Method** |`POST` |
|**URI** |`/redfish/v1/AggregationService/Actions/AggregationService.Remove` |
|**Returns** | <ul><li>`Location` URI of the task monitor associated with this operation in the response header.</li><li>Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See "Sample response body \(HTTP 202 status\)".</li><li>On successful completion, a message in the response body, saying that the operation is completed successfully. See "Sample response body \(HTTP 200 status\)".</li></ul>|
|**Response Code** |<ul><li>`202 Accepted`</li><li>`200 OK`</li></ul> |
|**Authentication** |Yes|


<br>
**Description**

This action removes the inventory of a specific plugin. 

It is performed as a task. 
To know the progress of this action, perform `GET` on the [task monitor](#task-monitor) returned in the response header \(until the task is complete\).

<aside class="notice">
Before removing the plugin, ensure that the plugin container is stopped.
</aside>





> Sample request body

```
{
	"parameters": [
		{
	       "Name": "/redfish/v1/Managers/a6ddc4c0-2568-4e16-975d-fa771b0be853"
        }
	]
}
```

###  Request parameters

|Parameter|Type|Description|
|---------|----|-----------|
|parameters\[\{|Array|The parameters associated with the `Delete` Action.|
|Name\}\]|String \(required\)<br> |The URI of the target to be removed: `/redfish/v1/Managers/{managerId}` |



> Sample response header \(HTTP 202 status\)

```
Connection:keep-alive
Content-Type:application/json; charset=utf-8
Location:/taskmon/task85de4003-8757-4c7d-941f-55eaf7d6412a
Odata-Version:4.0
X-Frame-Options:sameorigin
Date:Sun,17 May 2020 14:35:32 GMT+5m 13s
Content-Length:491 bytes

```

> Sample response body \(HTTP 202 status\)

```
{
   "@odata.type":"#Task.v1_4_2.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4003-8757-4c7d-941f-55eaf7d6412a",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task85de4003-8757-4c7d-941f-55eaf7d6412a",
   "Name":"Task task85de4003-8757-4c7d-941f-55eaf7d6412a",
   "Message":"The task with id task85de4003-8757-4c7d-941f-55eaf7d6412a has started.",
   "MessageId":"TaskEvent.1.0.1.TaskStarted",
   "MessageArgs":[
      "task85de4003-8757-4c7d-941f-55eaf7d6412a"
   ],
   "NumberOfArgs":1,
   "Severity":"OK"
}
```

> Sample response body \(HTTP 200 status\)

```
{ 
   "error":{ 
      "code":"ResourceEvent.1.0.2.ResourceRemoved",
      "message":"Request completed successfully"
   }
}
```
