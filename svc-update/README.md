# Software and firmware update

The resource aggregator exposes Redfish update service endpoints. Use these endpoints to access and update the software components of a system such as BIOS and firmware. Using these endpoints, you can also update firmware of other components such as system drivers and provider software.

The `UpdateService` schema describes the update service and the properties for the service itself. It exposes the firmware and software inventory resources and provides links to access them.



|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/UpdateService|`GET`|`Login` |
|/redfish/v1/UpdateService/FirmwareInventory|`GET`|`Login` |
|/redfish/v1/UpdateService/FirmwareInventory/\{inventoryId\}|`GET`|`Login` |
|/redfish/v1/UpdateService/SoftwareInventory|`GET`|`Login` |
|/redfish/v1/UpdateService/SoftwareInventory/\{inventoryId\}|`GET`|`Login` |
|/redfish/v1/UpdateService/Actions/UpdateService.SimpleUpdate|`POST`|`ConfigureComponents` |
|/redfish/v1/UpdateService/Actions/UpdateService.StartUpdate|`POST`|`ConfigureComponents` |

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
|<strong>Response Code</strong> |On success, `200 Ok` |
|<strong>Authentication</strong> |Yes|



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
   "@odata.type":"#SoftwareInventory.v1_4_0.SoftwareInventory",
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
   "@odata.type":"#SoftwareInventory.v1_4_0.SoftwareInventory",
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

 
**Usage**

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
            "@odata.type": "#Settings.v1_3_3.OperationApplyTimeSupport",
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
|NFS \(v1.3+\)<br> |Network File System.|
|NSF \(deprecated v1.3\)<br> | Network File System.<br> This value has been deprecated in favor of NFS.<br> |
|OEM|A manufacturer-defined protocol.|
|SCP|Secure Copy Protocol.|
|SFTP \(v1.1+\)<br> |Secure File Transfer Protocol.|
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

**Usage** 

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


