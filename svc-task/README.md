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
  
  
##  Modifying Configurations of Tasks Service
  
Config file of ODIMRA is located at: **odimra/lib-utilities/config/odimra_config.json**  
Refer to the section **Modifying Configurations** in the README.md file to change the configurations of an odimra service.
  
**Specific configurations for Tasks Service are:**
  
##  Log Location of the Tasks Service
  
/var/log/ODIMRA/task.log
    
  







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
   "@odata.type":"#Task.v1_4_2a.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task6e3cdbd8-65ca-4842-9437-0f29d5c6bce3",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task6e3cdbd8-65ca-4842-9437-0f29d5c6bce3",
   "Name":"Task task6e3cdbd8-65ca-4842-9437-0f29d5c6bce3",
   "TaskState":"Completed",
   "StartTime":"2020-04-17T06:52:15.632439043Z",
   "EndTime":"2020-04-17T06:52:15.947796761Z",
   "TaskStatus":"OK",
   "SubTasks":"",
   "TaskMonitor":"/taskmon/task6e3cdbd8-65ca-4842-9437-0f29d5c6bce3",
   "PercentComplete":100,
   "Payload":{
      "HttpHeaders":{
         "Cache-Control":"no-cache",
         "Connection":"keep-alive",
         "Content-type":"application/json; charset=utf-8",
         "Location":"/redfish/v1/Systems/2412b8c8-3b02-40cd-926e-8e85cb406e63:1",
         "OData-Version":"4.0",
         "Transfer-Encoding":"chunked"
      },
      "HttpOperation":"POST",
      "JsonBody":{
         "code":"Base.1.6.1.Success",
         "message":"Request completed successfully."
      },
      "TargetUri":"/redfish/v1/Systems/2412b8c8-3b02-40cd-926e-8e85cb406e63:1"
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







 

##  Deleting a task

|||
|-----------|----------|
|**Method** | `DELETE` |
|**URI** |`/redfish/v1/TaskService/Tasks/{TaskID}` |
|**Description** |This operation deletes a specific task. Deleting a running task aborts the operation being carried out.<br>**NOTE:**<br> Only a user having `ConfigureComponents` privilege is authorized to delete a task. If you do not have the necessary privileges, you will receive an HTTP `403 Forbidden` error.|
|**Returns** |JSON schema representing the deleted task.|
|**Response code** |`202 Accepted` |
|**Authentication** |Yes|


>**curl command**


```
curl -i DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/TaskService/Tasks/{TaskID}'

```

>**Sample response body** 

```
{
   "@odata.type":"#Task.v1_4_2a.Task",
   "@odata.id":"/redfish/v1/TaskService/Tasks/task52589fad-22fb-4505-86fd-cf845d500a33",
   "@odata.context":"/redfish/v1/$metadata#Task.Task",
   "Id":"task52589fad-22fb-4505-86fd-cf845d500a33",
   "Name":"Task task52589fad-22fb-4505-86fd-cf845d500a33",
   "Message":"The task with id task52589fad-22fb-4505-86fd-cf845d500a33 has started.",
   "MessageId":"TaskEvent.1.0.1.TaskStarted",
   "MessageArgs":[
      "task52589fad-22fb-4505-86fd-cf845d500a33"
   ],
   "NumberOfArgs":1,
   "Severity":"OK",
   "TaskState":"Completed",
   "StartTime":"2020-04-13T10:29:21.406470494Z",
   "EndTime":"2020-04-13T10:29:21.927909927Z",
   "TaskStatus":"OK",
   "SubTasks":"",
   "TaskMonitor":"/taskmon/task52589fad-22fb-4505-86fd-cf845d500a33",
   "PercentComplete":100,
   "Payload":{
      "HttpHeaders":{
         "Content-type":"application/json; charset=utf-8",
         "Location":"/redfish/v1/Managers/a6ddc4c0-2568-4e16-975d-fa771b0be853"
      },
      "HttpOperation":"POST",
      "JsonBody":{
         "code":"Base.1.6.1.Success",
         "message":"Request completed successfully."
      },
      "TargetUri":"/redfish/v1/AggregationService/Actions/AggregationService.Add"
   },
   "Messages":[

   ]
}
```