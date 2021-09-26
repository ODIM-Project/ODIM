# Redfish Telemetry

Telemetry refers to the metrics obtained from remote systems for analysis and monitoring. 
The Redfish Telemetry model is designed to obtain characteristics of metrics, send specific metric reports periodically and specify triggers against metrics.

Resource Aggregator for ODIM exposes the Redfish `TelemetryService` APIs to:

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
| /redfish/v1//TelemetryService/MetricReportDefinitions        | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/MetricReportDefinitions/{MetricReportDefinitionID} | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/MetricReports                   | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/MetricReports/{MetricReportID}  | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/Triggers                        | GET                  | `Login`                 |
| /redfish/v1/TelemetryService/Triggers/{TriggerID}            | GET, PATCH           | `Login`,`ConfigureSelf` |

>**NOTE:**
>Before accessing these endpoints, ensure that the user has the required privileges. If you access these endpoints without necessary privileges, you will receive an HTTP `403 Forbidden` error.

## Viewing the telemetry service root

| **Method**         | `GET`                                                        |
| ------------------ | ------------------------------------------------------------ |
| **URI**            | `/redfish/v1/TelemetryService`                               |
| **Description**    | This operation retrieves JSON schema representing the Redfish `TelemetryService` root. |
| **Returns**        | Properties for the Redfish `TelemetryService` and links to its list of resources. |
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
   "@odata.type":"#TelemetryService.v1_3_1.TelemetryService",
   "Id":"TelemetryService",
   "Name":"Telemetry Service",
   "Status":{
      "State":"Enabled",
      "Health":"OK"
   },
   "ServiceEnabled":true,
   "SupportedCollectionFunctions":[
      "Average",
      "Minimum",
      "Maximum"
   ],
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
    "@odata.type": "#MetricDefinition.v1_0_0.MetricDefinition",
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
                "9616fec9-c76a-4d26-ab53-196d08ce825a:1",
                "ba5cd083-b360-4994-bc30-12b450859b27:1"
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
    "@odata.type": "#MetricReportDefinition.v1_0_0.MetricReportDefinition",
    "Description": "Metric report of CPU Utilization for 10 minutes with sensing interval of 20 seconds.",
    "Id": "CPUUtilCustom1",
    "MetricProperties": [
        "/redfish/v1/Systems/{SystemID}#SystemUsage/CPUUtil",
        "SystemUsage/CPUUtil"
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
                "9616fec9-c76a-4d26-ab53-196d08ce825a:1"
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
| **Returns**        | Links of the metric reports.                                 |
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
| **Returns**        | Link to the metric report.                                   |
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
         "MetricProperty":"/redfish/v1/Systems/9616fec9-c76a-4d26-ab53-196d08ce825a:1#SystemUsage/CPUUtil",
         "MetricValue":"1",
         "Timestamp":"2021-08-28T13:46:05Z"
      },
      {
         "MetricDefinition":{
            "@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"
         },
         "MetricId":"CPUUtil",
         "MetricProperty":"/redfish/v1/Systems/9616fec9-c76a-4d26-ab53-196d08ce825a:1#SystemUsage/CPUUtil",
         "MetricValue":"1",
         "Timestamp":"2021-08-28T13:46:25Z"
      },
      {
         "MetricDefinition":{
            "@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"
         },
         "MetricId":"CPUUtil",
         "MetricProperty":"/redfish/v1/Systems/9616fec9-c76a-4d26-ab53-196d08ce825a:1#SystemUsage/CPUUtil",
         "MetricValue":"1",
         "Timestamp":"2021-08-28T13:46:45Z"
      },
      {
         "MetricDefinition":{
            "@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/CPUUtil"
         },
         "MetricId":"CPUUtil",
         "MetricProperty":"/redfish/v1/Systems/9616fec9-c76a-4d26-ab53-196d08ce825a:1#SystemUsage/CPUUtil",
         "MetricValue":"1",
         "Timestamp":"2021-08-28T13:47:05Z"
      },
      "Name":"Metric report of CPU Utilization for 60 minutes with sensing interval of 20 seconds."
   }
```

<blockquote> NOTE:  After you remove a system and perform a `GET ` operation on the Metric Report Collection, the collection of all individual metric reports is still displayed in the response body. When you perform a `GET` operation on that individual {MetricReportID}, you get a `404-Not Found` error message. After this, when you perform a GET operation on the Metric Report Collection again, the instance of that individual metric report is erased. 
This is an implementation choice in Resource Aggregator for ODIM, because Telemetry service is defined for a collection of BMCs and not for an individual BMC as per the DMTF Redfish specification.</blockquote>

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
                "9616fec9-c76a-4d26-ab53-196d08ce825a:1",
                "ba5cd083-b360-4994-bc30-12b450859b27:1"
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
| **Returns**        |                                                      |
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
 'https://10.24.1.33:30080/redfish/v1/TelemetryService/Triggers/{TriggersID}'

```


>**Sample request body**

```
{
  "EventTriggers": ["Alert"]
}
```



