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
|/redfish/v1/AccountService|GET|`Login` |
|/redfish/v1/AccountService/Roles|GET, POST|`Login`, `ConfigureManager` |
|/redfish/v1/AccountService/Roles/\{RoleId\}|GET, PATCH, DELETE|`Login`, `ConfigureManager` |
|/redfish/v1/AccountService/Accounts|POST, GET|`Login`, `ConfigureUsers` |
|/redfish/v1/AccountService/Accounts/\{accountId\}|GET, DELETE, PATCH|`Login`, `ConfigureUsers`, `ConfigureSelf` |



>**NOTE:**
Before accessing these endpoints, ensure that the user has the required privileges. If you access these endpoints without necessary privileges, you will receive an HTTP `403 Forbidden` error.

  
  
##  Modifying Configurations of Session and Account Service
  
Config file of ODIMRA is located at: **odimra/lib-utilities/config/odimra_config.json**  
Refer to the section **Modifying Configurations** in the README.md file to change the configurations of an odimra service.

  
**Specific configurations for Session and Account Service are:**
  
##  Log Location of the Session and Account Service
  
/var/log/ODIMRA/account_session.log
    
  
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
   "@odata.type":"#AccountService.v1_9_0.AccountService",
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
   "@odata.type":"#Role.v1_3_1.Role",
   "@odata.id":"/redfish/v1/AccountService/Roles/CLIENT11",
   "Id":"CLIENT11",
   "Name":"User Role",
   "Message":"The resource has been created successfully.",
   "MessageId":"ResourceEvent.1.0.3.ResourceCreated",
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
   "@odata.type":"#Role.v1_3_1.Role",
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
   "@odata.type":"#ManagerAccount.v1_8_0.ManagerAccount",
   "@odata.id":"/redfish/v1/AccountService/Accounts/monitor32",
   "@odata.context":"/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
   "Id":"monitor32",
   "Name":"Account Service",
   "Message":"The resource has been created successfully",
   "MessageId":"Base.1.10.0.Created",
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
   "@odata.type":"#ManagerAccount.v1_8_0.ManagerAccount",
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
   "@odata.type":"#ManagerAccount.v1_8_0.ManagerAccount",
   "@odata.id":"/redfish/v1/AccountService/Accounts/monitor32",
   "@odata.context":"/redfish/v1/$metadata#ManagerAccount.ManagerAccount",
   "Id":"monitor32",
   "Name":"Account Service",
   "Message":"The account was successfully modified.",
   "MessageId":"Base.1.10.0.AccountModified",
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