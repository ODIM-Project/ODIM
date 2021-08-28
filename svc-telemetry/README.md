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
   "@odata.type":"#TelemetryService.v1_2_0.TelemetryService",
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
   "@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/PowerConsumedWatts",
   "@odata.type":"#MetricDefinition.v1_0_3.MetricDefinition",
   "Id":"PowerConsumedWatts",
   "Name":"Power Consumed Watts Metric Definition",
   "MetricType":"Numeric",
   "Implementation":"PhysicalSensor",
   "PhysicalContext":"PowerSupply",
   "MetricDataType":"Decimal",
   "Units":"W",
   "Precision":4,
   "Accuracy":1,
   "Calibration":2,
   "MinReadingRange":0,
   "MaxReadingRange":50,
   "SensingInterval":"PT1S",
   "TimestampAccuracy":"PT1S",
   "Wildcards":[
      {
         "Name":"ChassisID",
         "Values":[
            "1",
            "2",
            "3"
         ]
      }
   ],
   "MetricProperties":[
      "/redfish/v1/Chassis/{ChassisID}/Power#/PowerControl/0/PowerConsumedWatts",
      "/redfish/v1/Chassis/{ChassisID}/Power#/PowerControl/1/PowerConsumedWatts"
   ]
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
   "@odata.type":"#MetricDefinition.v1_2_0.MetricDefinition",
   "Id":"PowerConsumedWatts",
   "Name":"Power Consumed Watts Metric Definition",
   "MetricType":"Numeric",
   "Implementation":"PhysicalSensor",
   "PhysicalContext":"PowerSupply",
   "MetricDataType":"Decimal",
   "Units":"W",
   "Precision":4,
   "Accuracy":1,
   "Calibration":2,
   "MinReadingRange":0,
   "MaxReadingRange":50,
   "SensingInterval":"PT1S",
   "TimestampAccuracy":"PT1S",
   "Wildcards":[
      {
         "Name":"ChassisID",
         "Values":[
            "1"
         ]
      }
   ],
   "MetricProperties":[
      "/redfish/v1/Chassis/{ChassisID}/Power#/PowerControl/0/PowerConsumedWatts"
   ],
   "@odata.id":"/redfish/v1/TelemetryService/MetricDefinitions/PowerConsumedWatts"
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
   "@odata.id":"/redfish/v1/TelemetryService/MetricReportDefinitions/PowerMetrics",
   "@odata.type":"#MetricReportDefinition.v1_3_0.MetricReportDefinition",
   "Id":"PowerMetrics",
   "Name":"Transmit and Log Power Metrics",
   "MetricReportDefinitionType":"Periodic",
   "MetricReportDefinitionEnabled":true,
   "Schedule":{
      "RecurrenceInterval":"PT0.1S"
   },
   "ReportActions":[
      "RedfishEvent",
      "LogToMetricReportsCollection"
   ],
   "ReportUpdates":"Overwrite",
   "MetricReport":{
      "@odata.id":"/redfish/v1/TelemetryService/MetricReports/PowerMetrics"
   },
   "Status":{
      "State":"Enabled"
   },
   "Wildcards":[
      {
         "Name":"PWild",
         "Values":[
            "0",
            "1"
         ]
      },
      {
         "Name":"TWild",
         "Values":[
            "Tray_1",
            "Tray_2",
            "Tray_3"
         ]
      }
   ],
   "MetricProperties":[
      "/redfish/v1/Chassis/{TWild}/Power#/PowerControl/{PWild}/PowerMetrics/
AverageConsumedWatts",
      "/redfish/v1/Chassis/{TWild}/Power#/PowerControl/{PWild}/PowerMetrics/
MinConsumedWatts",
      "/redfish/v1/Chassis/{TWild}/Power#/PowerControl/{PWild}/PowerMetrics/
MaxConsumedWatts"
   ]
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
   "@odata.type":"#MetricReport.v1_4_2.MetricReport",
   "Id":"AvgPlatformPowerUsage",
   "Name":"Average Platform Power Usage metric report",
   "ReportSequence":"127",
   "MetricReportDefinition":{
      "@odata.id":"/redfish/v1/TelemetryService/MetricReportDefinitions/AvgPlatformPowerUsage"
   },
   "MetricValues":[
      {
         "MetricId":"AverageConsumedWatts",
         "MetricValue":"100",
         "Timestamp":"2016-11-08T12:25:00-05:00",
         "MetricProperty":"/redfish/v1/Chassis/Tray_1/Power#/0/PowerConsumedWatts"
      },
      {
         "MetricId":"AverageConsumedWatts",
         "MetricValue":"94",
         "Timestamp":"2016-11-08T13:25:00-05:00",
         "MetricProperty":"/redfish/v1/Chassis/Tray_1/Power#/0/PowerConsumedWatts"
      },
      {
         "MetricId":"AverageConsumedWatts",
         "MetricValue":"100",
         "Timestamp":"2016-11-08T14:25:00-05:00",
         "MetricProperty":"/redfish/v1/Chassis/Tray_1/Power#/0/PowerConsumedWatts"
      }
   ],
   "@odata.id":"/redfish/v1/TelemetryService/MetricReports/AvgPlatformPowerUsage"
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
   "@odata.id":"/redfish/v1/TelemetryService/MetricReports/PlatformPowerUsage",
   "@odata.type":"#MetricReport.v1_3_0.MetricReport",
   "Id":"PlatformPowerUsage",
   "Name":"PlatformPowerUsage",
   "MetricReportDefinition":{
      "@odata.id":"/redfish/v1/TelemetryService/MetricReportDefinitions/
PlatformPowerUsage"
   },
   "MetricValues":[
      {
         "Timestamp":"2016-11-08T12:25:00-05:00",
         "MetricValue":"103",
         "MetricProperty":"/redfish/v1/Chassis/Tray_1/Power#/PowerControl/0/
PlatformConsumedWatts"
      },
      {
         "Timestamp":"2016-11-08T12:25:00-05:00",
         "MetricValue":"103",
         "MetricProperty":"/redfish/v1/Chassis/Tray_2/Power#/PowerControl/0/
PlatformConsumedWatts"
      },
      {
         "Timestamp":"2016-11-08T13:25:00-05:00",
         "MetricValue":"106",
         "MetricProperty":"/redfish/v1/Chassis/Tray_1/Power#/PowerControl/0/
PlatformConsumedWatts"
      },
      {
         "Timestamp":"2016-11-08T13:25:00-05:00",
         "MetricValue":"106",
         "MetricProperty":"/redfish/v1/Chassis/Tray_2/Power#/PowerControl/0/
PlatformConsumedWatts"
      },
      {
         "Timestamp":"2016-11-08T14:25:00-05:00",
         "MetricValue":"107",
         "MetricProperty":"/redfish/v1/Chassis/Tray_1/Power#/PowerControl/0/
PlatformConsumedWatts"
      },
      {
         "Timestamp":"2016-11-08T14:25:00-05:00",
         "MetricValue":"107",
         "MetricProperty":"/redfish/v1/Chassis/Tray_2/Power#/PowerControl/0/
PlatformConsumedWatts"
      }
   ]
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
   "@odata.type":"#MetricReport.v1_4_2.MetricReport",
   "Id":"AvgPlatformPowerUsage",
   "Name":"Average Platform Power Usage metric report",
   "ReportSequence":"127",
   "MetricReportDefinition":{
      "@odata.id":"/redfish/v1/TelemetryService/MetricReportDefinitions/AvgPlatformPowerUsage"
   },
   "MetricValues":[
      {
         "MetricId":"AverageConsumedWatts",
         "MetricValue":"100",
         "Timestamp":"2016-11-08T12:25:00-05:00",
         "MetricProperty":"/redfish/v1/Chassis/Tray_1/Power#/0/PowerConsumedWatts"
      },
      {
         "MetricId":"AverageConsumedWatts",
         "MetricValue":"94",
         "Timestamp":"2016-11-08T13:25:00-05:00",
         "MetricProperty":"/redfish/v1/Chassis/Tray_1/Power#/0/PowerConsumedWatts"
      },
      {
         "MetricId":"AverageConsumedWatts",
         "MetricValue":"100",
         "Timestamp":"2016-11-08T14:25:00-05:00",
         "MetricProperty":"/redfish/v1/Chassis/Tray_1/Power#/0/PowerConsumedWatts"
      }
   ],
   "@odata.id":"/redfish/v1/TelemetryService/MetricReports/AvgPlatformPowerUsage"
}
```

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
   "@odata.id":"/redfish/v1/TelemetryService/Triggers/PlatformPowerCapTriggers",
   "@odata.type":"#Triggers.v1_1_1.Triggers",
   "Id":"PlatformPowerCapTriggers",
   "Name":"Triggers for platform power consumed",
   "MetricType":"Numeric",
   "TriggerActions":[
      "RedfishEvent"
   ],
   "NumericThresholds":{
      "UpperCritical":{
         "Reading":50,
         "Activation":"Increasing",
         "DwellTime":"PT0.001S"
      },
      "UpperWarning":{
         "Reading":48.1,
         "Activation":"Increasing",
         "DwellTime":"PT0.004S"
      }
   },
   "MetricProperties":[
      "/redfish/v1/Chassis/1/Power#/PowerControl/0/PowerConsumedWatts"
   ]
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
   "@odata.type":"#Triggers.v1_1_4.Triggers",
   "Id":"PlatformPowerCapTriggers",
   "Name":"Triggers for platform power consumed",
   "MetricType":"Numeric",
   "TriggerActions":[
      "RedfishEvent"
   ],
   "NumericThresholds":{
      "UpperCritical":{
         "Reading":50,
         "Activation":"Increasing",
         "DwellTime":"PT0.001S"
      },
      "UpperWarning":{
         "Reading":48.1,
         "Activation":"Increasing",
         "DwellTime":"PT0.004S"
      }
   },
   "MetricProperties":[
      "/redfish/v1/Chassis/1/Power#/PowerControl/0/PowerConsumedWatts"
   ],
   "@odata.id":"/redfish/v1/TelemetryService/Triggers/{TriggerID}"
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



