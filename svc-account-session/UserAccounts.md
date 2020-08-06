#  User accounts

Resource Aggregator for ODIM allows users to have accounts to configure their actions and restrictions.

Resource Aggregator for ODIM has an administrator user account by default.

Create other user accounts by defining a username, a password, and a role for each account. The username and the password are used to authenticate with the Redfish services \(using `BasicAuth` or `XAuthToken`\).

Resource Aggregator for ODIM exposes Redfish `AccountsService` APIs to create and manage user accounts. Use these endpoints to perform the following operations:

-   Creating, modifying, and deleting account details.

-   Fetching account details.

  
  
##  Modifying Configurations of Session and Account Service
  
Config File of ODIMRA is located at: **odimra/lib-utilities/config/odimra_config.json**  
Refer the section **Modifying Configurations** in the README.md to change the configurations of a odimra service
  
**Specific configurations for Session and Account Service are:**
  
##  Log Location of the Session and Account Service
  
/var/log/ODIMRA/account_session.log
    
  

##  Supported APIs

|||
|-------|--------------------|
|/redfish/v1/AccountService|`GET`|
|/redfish/v1/AccountService/Accounts|`POST`, `GET`|
|/redfish/v1/AccountService/Accounts/\{accountId\}|`GET`, `DELETE`, `PATCH`|




## AccountService root

```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/AccountService'


```

> Sample response header

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


> Sample response body

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


|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/AccountService` |
|**Description** |The schema representing the Redfish `AccountService` root.|
|**Returns** |The properties common to all user accounts and links to the collections of manager accounts and roles.|
|**Response Code** | `200 OK` |
|**Authentication** |Yes|

 










## Creating a user account

```
curl -i POST \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
   -H "Content-Type:application/json" \
   -d \
'{"UserName":"{username}","Password":"{password}","RoleId":"{roleId}"}
' \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Accounts'


```


|||
|-------|--------------------|
|**Method** | `POST` |
|**URI** |`/redfish/v1/AccountService/Accounts` |
|**Description** |This operation creates a user account. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can create a user account.|
|**Returns** |<ul><li>`Location` header that contains a link to the newly created account.</li><li>JSON schema representing the newly created account.</li></ul> |
|**Response Code** |`201 Created` |
|**Authentication** |Yes|

 





> Sample request body

```
{ 
   "UserName":"monitor32",
   "Password":"Abc1vent2020!",
   "RoleId":"CLIENT11"
}
```

###  Request URI parameters

|Parameter|Type|Description|
|---------|----|-----------|
|UserName|String \(required\)<br> |User name for the user account.|
|Password|String \(required\)<br> |Password for the user account. Before creating a password, see "Password Requirements" .|
|RoleId|String \(required\)<br> |The role for this account. To know more about roles, see [User roles and privileges](#user-roles-and-privileges). Ensure that the `roleId` you want to assign to this user account exists. To check the existing roles, see [Listing Roles](#listing-roles). If you attempt to assign an unavailable role, you will receive an HTTP `400 Bad Request` error.|

### Password requirements

-   Your password must not be same as your username.

-   Your password must be at least 12 characters long and at most 16 characters long.

-   Your password must contain at least one uppercase letter \(A-Z\), one lowercase letter \(a-z\), one digit \(0-9\), and one special character \(~!@\#$%^&\*-+\_|\(\)\{\}:;<\>,.?/\).


> Sample response header

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

> Sample response body

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


```
curl -i GET \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Accounts'


```


|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/AccountService/Accounts` |
|**Description** |This operation retrieves a list of user accounts. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can view a list of user accounts.|
|**Returns** |Links to user accounts.|
|**Response Code** |`200 OK` |
|**Authentication** |Yes|


##  Viewing the account details

```
curl -i GET \
   -H 'Authorization:Basic {base64_encoded_string_of_[username:password]}' \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Accounts/{accountId}'


```

> Sample response body

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


|||
|---------|---------------|
|**Method** | `GET` |
|**URI** |`/redfish/v1/AccountService/Accounts/{accountId}` |
|**Description** |This operation fetches information about a specific user account. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can view information about a user account.|
|**Returns** |JSON schema representing this user account.|
|**Response Code** |`200 OK` |
|**Authentication** |Yes|

 







## Updating a user account

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


|||
|---------|---------------|
|**Method** | `PATCH` |
|**URI** |`/redfish/v1/AccountService/Accounts/{accountId}` |
|**Description** |This operation updates user account details \(`username`, `password`, and `RoleId`\). To modify account details, add them in the request payload \(as shown in the sample request body\) and perform `PATCH` on the mentioned URI. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can modify user accounts. Users with `ConfigureSelf` privilege can modify only their own accounts.|
|**Returns** |<ul><li>`Location` header that contains a link to the updated account.</li><li>JSON schema representing the modified account.</li></ul>|
|**Response Code** |`200 OK` |
|**Authentication** |Yes|

 




> Sample request body

```
{ 
   "Password":"Testing)9-_?{}",
   "RoleId":"CLIENT11"
}
```

> Sample response header

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

> Sample response body

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

```
curl  -i -X DELETE \
   -H "X-Auth-Token:{X-Auth-Token}" \
 'https://{odimra_host}:{port}/redfish/v1/AccountService/Accounts/{accountId}'

```


|||
|---------|---------------|
|**Method** | `DELETE` |
|**URI** |`/redfish/v1/AccountService/Accounts/{accountId}` |
|**Description** |This operation deletes a user account. <br>**NOTE:**<br> Only a user with `ConfigureUsers` privilege can delete a user account.|
|**Response Code** |`204 No Content` |
|**Authentication** |Yes|