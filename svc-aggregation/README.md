#  Resource aggregation and management

The resource aggregator allows you to add southbound infrastructure into its inventory, create resource collections, and perform actions in combination on these collections. It exposes Redfish aggregation service endpoints to achieve the following:

-   Adding a resource and building its inventory.

-   Resetting one or more resources.

-   Changing the boot path of one or more resources to default settings.

-   Removing a resource from the inventory which is no longer managed.


All aggregation actions are performed as [tasks](#tasks) in Resource Aggregator for ODIM. The actions performed on a group of resources \(resetting or changing the boot order to default settings\) are carried out as a set of subtasks.

**Supported endpoints**

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/AggregationService|GET|`Login` |
|/redfish/v1/AggregationService/AggregationSources<br> |GET, POST|`Login`, `ConfigureManager` |
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
  
  
##  Modifying configurations of the aggregation service
  
Config file of ODIMRA is located at: **odimra/lib-utilities/config/odimra_config.json**  
Refer to the **Modifying Configurations** section in the README.md file to change the configurations of an ODIMRA service.
  
**Specific configurations for Aggregation Service are:**
  
  
  
  
##  Log location of the aggregation service
  
/var/log/ODIMRA/aggregation.log
  


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
   "@odata.type":"#AggregationService.v1_0_1.AggregationService",
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
   "@odata.type":"#Task.v1_5_1.Task",
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
|<strong>Response Code</strong> |On success, `202 Accepted` On successful completion of the task, `201 Created` <br> |
|<strong>Authentication</strong> |Yes|

**Usage information**

Perform HTTP POST on the mentioned URI with a request body specifying a connection method to use for adding the BMC. To know about connection methods, see [Connection methods](#connection-methods).
				
A Redfish task will be created and you will receive a link to the [task monitor](#viewing-a-task-monitor) associated with it.

To know the progress of this operation, perform HTTP `GET` on the task monitor returned in the response header (until the task is complete).
When the task is successfully complete, you will receive aggregation source Id of the added BMC. Save it as it is required to identify it in the resource inventory later.

After the server is successfully added as an aggregation source, it will also be available as a computer system resource at `/redfish/v1/Systems/` and a manager resource at `/redfish/v1/Managers/`.

To view the list of links to computer system resources, perform HTTP `GET` on `/redfish/v1/Systems/`. Each link contains `ComputerSystemId` of a specific BMC. For more information, see [Collection of computer systems](#(#collection-of-computer-systems)).

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
   "@odata.type":"#Task.v1_5_1.Task",
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
   "@odata.type":"#Task.v1_5_1.Task",
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

>**Sample response body** \(subtask\)

```
{
   "@odata.type":"#SubTask.v1_4_2.SubTask",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4003-8757-4c7d-942f-55eaf7d6412a/SubTasks/task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "@odata.context":"/redfish/v1/$metadata#SubTask.SubTask",
   "Id":"task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "Name":"Task task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "Message":"Successfully Completed Request",
   "MessageId":"Base.1.10.0.Success",
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
      "code":"Base.1.10.0.Success",
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
|<strong>Returns</strong> |- `Location` URI of the task monitor associated with this operation in the response header. See `Location` URI highlighted in bold in "Sample response header \(HTTP 202 status\)".<br>-  Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See the task URI and the task Id highlighted in bold in "Sample response body \(HTTP 202 status\)".<blockquote><br>IMPORTANT:<br>Note down the task Id. If the task completes with an error, it is required to know which subtask has failed. To get the list of subtasks, perform HTTP `GET` on `/redfish/v1/TaskService/Tasks/{taskId}`.</blockquote><br>- On successful completion of this operation, a message in the response body, saying that the operation is completed successfully. See "Sample response body \(HTTP 200 status\)".<br>|
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
   "@odata.type":"#Task.v1_5_1.Task",
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

>**Sample response body** \(subtask\)

```
{
   "@odata.type":"#SubTask.v1_4_2.SubTask",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4003-8757-4c7d-942f-55eaf7d6412a/SubTasks/task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "@odata.context":"/redfish/v1/$metadata#SubTask.SubTask",
   "Id":"task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "Name":"Task task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "Message":"Successfully Completed Request",
   "MessageId":"Base.1.10.0.Success",
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
      "code":"Base.1.10.0.Success",
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
|<strong>Returns</strong> |`Location` URI of the task monitor associated with this operation in the response header. See `Location` URI highlighted in bold in "Sample response header \(HTTP 202 status\)".<br>-   Link to the task and the task Id in the sample response body. To get more information on the task, perform HTTP `GET` on the task URI. See "Sample response body \(HTTP 202 status\)".<br>|
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
   "@odata.type":"#Task.v1_5_1.Task",
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
      "MessageId":"Base.1.10.0.Created",
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
      "MessageId":"Base.1.10.0.Success",
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
   "MessageId":"Base.1.10.0.Success",
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
      "MessageId":"Base.1.10.0.Created",
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
   "@odata.type":"#Task.v1_5_1.Task",
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

>**Sample response body** \(subtask\)

```
{
   "@odata.type":"#SubTask.v1_4_2.SubTask",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591/SubTasks/task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "@odata.context":"/redfish/v1/$metadata#SubTask.SubTask",
   "Id":"task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591",
   "Name":"Task task8cf1ed8b-bb83-431a-9fa6-1f8d349a8591",
   "Message":"Successfully Completed Request",
   "MessageId":"Base.1.10.0.Success",
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
      "code":"Base.1.10.0.Success",
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
   "@odata.type":"#Task.v1_5_1.Task",
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

>**Sample response body** \(subtask\)

```
{
   "@odata.type":"#SubTask.v1_4_2.SubTask",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task85de4003-8757-4c7d-942f-55eaf7d6412a/SubTasks/task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "@odata.context":"/redfish/v1/$metadata#SubTask.SubTask",
   "Id":"task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "Name":"Task task22a98864-5dd8-402b-bfe0-0d61e265391e",
   "Message":"Successfully Completed Request",
   "MessageId":"Base.1.10.0.Success",
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
      "code":"Base.1.10.0.Success",
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


