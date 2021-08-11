# Sessions

A session represents a window of user's login with a Redfish service and contains details about the user and the user activity. You can run multiple sessions simultaneously.

Resource Aggregator for ODIM offers Redfish `SessionService` interface for creating and managing sessions. It exposes APIs to achieve the following:

-   Fetching the `SessionService` root

-   Creating a session

-   Listing active sessions

-   Deleting a session


**Supported APIs**

|API URI|Operation Applicable|Required privileges|
|-------|--------------------|-------------------|
|/redfish/v1/SessionService|GET|`Login` |
|/redfish/v1/SessionService/Sessions|POST, GET|`Login`,|
|redfish/v1/SessionService/Sessions/\{sessionId\}|GET, DELETE|`Login`, `ConfigureManager`, `ConfigureSelf` |

>**NOTE:**
Before accessing these endpoints, ensure that the user has the required privileges. If you access these endpoints without necessary privileges, you will receive an HTTP `403 Forbidden` error.

##  Modifying Configurations of Session and Account Service
  
Config file of ODIMRA is located at: **odimra/lib-utilities/config/odimra_config.json**  
Refer to the section **Modifying Configurations** in the README.md file to change the configurations of an odimra service.
  
**Specific configurations for Session and Account Service are:**
  
##  Log Location of the Session and Account Service
  
/var/log/ODIMRA/account_session.log
    
  




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
	"MessageId": "Base.1.10.0.Created",
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
   "@odata.type":"#Session.v1_3_0.Session",
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