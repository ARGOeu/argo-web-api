---
id: ar_results
title: A/R results
sidebar_position: 1
---

## API Calls

| Name                                                                          | Description                                                                                                                                                                                                                              | Shortcut          |
| ----------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------- |
| GET: List Availability and Reliability results for a group of endpoint groups | This method retrieves the results of a specified group of endpoint group or multiple groups of endpoint groups of a specific type that where computed based on a given report. Results can be retrieved on daily, monthly or custom granularity. | [Description](#1) |
| GET: List Availability and Reliability results for an endpoint group          | This method retrieves the results of a specified endpoint group or multiple endpoint groups of a specific type that where computed based on a given report. Results can be retrieved on daily, monthly or custom granularity.                    | [Description](#2) |
| GET: List Availability and Reliability results for a Service Flavor           | This method retrieves the results of a specified service flavor that where computed based on a given report. Results can be retrieved on daily, monthly or custom granularity.                                                                   | [Description](#3) |
| GET: List Availability and Reliability results for an Endpoint                | This method retrieves the results of a specified service endpoint that where computed based on a given report. Results can be retrieved on daily, monthly or custom granularity.                                                                 | [Description](#4) |
| GET: Flat List of all endpoints Availability and Reliability results                | This method retrieves the results in a flat list of all service endpoints that where computed based on a given report. Results can be retrieved on daily, monthly or custom granularity. Pagination is supported.                                                                 | [Description](#5) |


## List Availabilities and Reliabilities for groups of Endpoint Groups {#1}

The following methods can be used to obtain a tenant's Availability and Reliability result metrics per Group of Endpoint Groups. The api authenticates the tenant using the api-key within the x-api-key header. User can specify time granularity (`monthly`, `daily` or `custom`) for retrieved results and also format using the `Accept` header. Depending on the form of the request the user can request a single group of endpoint groups results or a bulk of group of endpoint groups results filtered by their type.

## [GET] Group of Endpoint groups

### Input

```
/results/{report_name}/{group_type}?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{group_type}/{group_name}?[start_time]&[end_time]&[granularity]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start_time]`  | UTC time in W3C format                                                                          | YES      |
| `[end_time]`    | UTC time in W3C format                                                                          | YES      |
| `[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly`, `daily` or `custom` | NO       | `daily`       |

#### Path Parameters

| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report_name}` | Name of the report that contains all the information about the profile, filter tags, group types etc. | YES      |
| `{group_type}`  | Type of the Group of Endpoint Groups.                                                                 | NO       |
| `{group_name}`  | Name of the Group of Endpoint Groups.                                                                 | NO       |

#### Headers

##### Request

```
x-api-key: "tenant_key_value"
Accept: "application/xml" or "application/json"
```

##### Response

```
Status: 200 OK
```

#### URL

`/api/v2/results/Report_A/GROUP/GROUP_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily`

#### Response Body

```
<root>
    <group name="GROUP_A" type="GROUP">
        <results timestamp="2015-06-22" availability="68.13896116893515" reliability="50.413931144915935"></results>
        <results timestamp="2015-06-23" availability="75.36324059247399" reliability="80.8138510808647"></results>
    </group>
</root>
```


## [GET]: List Availabilities and Reliabilities for Endpoint Groups {#2}

The following methods can be used to obtain a tenant's Availability and Reliability result metrics per Endpoint Group. The api authenticates the tenant using the api-key within the x-api-key header. User can specify time granularity (`monthly`, `daily` or `custom`) for retrieved results and also format using the `Accept` header. Depending on the form of the request the user can request a single endpoint group results or a bulk of endpoint group results filtered by their type and if necessary their "top-level" group.

## [GET] Endpoint Groups

### Input

```
/results/{report_name}/{group_type}/{group_name}/{endpoint_group_type}?[start_time]&[end_time]&[granularity]
or simpler
/results/{report_name}/{endpoint_group_type}?[start_time]&[end_time]&[granularity]
and
/results/{report_name}/{group_type}/{group_name}/{endpoint_group_type}/{endpoint_group_name}?[start_time]&[end_time]&[granularity]
or simpler
/results/{report_name}/{endpoint_group_type}/{endpoint_group_name}?[start_time]&[end_time]&[granularity]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start_time]`  | UTC time in W3C format                                                                          | YES      |
| `[end_time]`    | UTC time in W3C format                                                                          | YES      |
| `[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly`, `daily` or `custom` | NO       | `daily`       |

#### Path Parameters

| Name                    | Description                                                                                                                         | Required | Default value |
| ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report_name}`         | Name of the report that contains all the information about the profile, filter tags, group types etc.                               | YES      |
| `{group_name}`          | Name of the Group of Endpoint Groups. If no name is specified then all Endpoint Groups regardless of top-level group are retrieved. | NO       |
| `{group_type}`          | Type of the Group of Endpoint Groups. If no type is specified then all groups are retrieved.                                        | NO       |
| `{endpoint_group_name}` | Name of the the Endpoint Group. If no name is specified then all groups are retrieved according to the `{endpoint_group_type}`.     | NO       |
| `{endpoint_group_type}` | Type of the the Endpoint Group.                                                                                                     | YES      |

#### Headers

##### Request

```
x-api-key: "tenant_key_value"
Accept: "application/xml" or "application/json"
```

##### Response

```
Status: 200 OK
```

#### URL

`/api/v2/results/Report_A/SITE/ST01?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily`

#### Response Body

```
<root>
    <group name="GROUP_A" type="GROUP">
        <group name="ST01" type="SITE">
            <results timestamp="2015-06-22" availability="66.7" reliability="54.6" unknown="0" uptime="1" downtime="0"></results>
            <results timestamp="2015-06-23" availability="100" reliability="100" unknown="0" uptime="1" downtime="0"></results>
        </group>
    </group>
</root>
```


## [GET]: List Availabilities and Reliabilities for Service Flavors {#3}

The following methods can be used to obtain a tenant's Availability and Reliability result metrics per given Service Flavor(s). The api authenticates the tenant using the api-key within the x-api-key header. The user can specify time granularity (`monthly`, `daily` or `custom`) for retrieved results and also format using the `Accept` header. Depending on the form of the request the user can request a single service flavor results or a bulk of service flavor results.

## [GET] Service Flavors

### Input

```
/results/{report_name}/{group_type}/{group_name}/{endpoint_group_type}/{endpoint_group_name}/services?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{group_type}/{group_name}/{endpoint_group_type}/{endpoint_group_name}/services/{service_flavor_type}?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{endpoint_group_type}/{endpoint_group_name}/services?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{endpoint_group_type}/{endpoint_group_name}/services/{service_flavor_type}?[start_time]&[end_time]&[granularity]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start_time]`  | UTC time in W3C format                                                                          | YES      |
| `[end_time]`    | UTC time in W3C format                                                                          | YES      |
| `[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly`, `daily` or `custom` | NO       | `daily`       |

#### Path Parameters

| Name                    | Description                                                                                                                           | Required | Default value |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report_name}`         | Name of the report that contains all the information about the profile, filter tags, group types etc.                                 | YES      |
| `{group_type}`          | Type of the Group of Endpoint Groups.                                                                                                 | NO       |
| `{group_name}`          | Name of the Group of Endpoint Groups.                                                                                                 | NO       |
| `{endpoint_group_type}` | Type of the the Endpoint Group.                                                                                                       | YES      |
| `{endpoint_group_name}` | Name of the the Endpoint Group.                                                                                                       | YES      |
| `{service_flavor_type}` | Type of the Service Flavor. If no type is given then results for all Service Flavors under the given Endpoint Group will be provided. | NO       |

#### Headers

##### Request

```
x-api-key: "tenant_key_value"
Accept: "application/xml" or "application/json"
```

##### Response

```
Status: 200 OK
```

#### URL

`/api/v2/results/Report_A/SITE/ST01/services?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:23:59Z&granularity=daily`

#### Response Body

```
<root>
  <group name="ST01" type="SITE">
    <group name="SF01" type="service">
      <results timestamp="2015-06-22" availability="98.26389" reliability="98.26389" unknown="0" uptime="0.98264" downtime="0"></results>
      <results timestamp="2015-06-23" availability="54.03509" reliability="81.48148" unknown="0.01042" uptime="0.53472" downtime="0.33333"></results>
    </group>
    <group name="SF02" type="service">
      <results timestamp="2015-06-22" availability="96.875" reliability="96.875" unknown="0" uptime="0.96875" downtime="0"></results>
      <results timestamp="2015-06-23" availability="100" reliability="100" unknown="0" uptime="1" downtime="0"></results>
    </group>
  </group>
</root>
```


## [GET]: List Availabilities and Reliabilities for Endpoints {#4}

The following methods can be used to obtain a tenant's Availability and Reliability result metrics for endpoints under a specific service or group. The api authenticates the tenant using the api-key within the x-api-key header. The user can specify time granularity (`monthly`, `daily` or `custom`) for retrieved results and also format using the `Accept` header. Depending on the form of the request the user can request a single service flavor results or a bulk of endpoint results.

## [GET] Endpoints A/R

### Input

Request endpoint a/r under specific service:

```
/results/{report_name}/{group_type}/{group_name}/{endpoint_group_type}/{endpoint_group_name}/services/{service_name}/endpoints?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{group_type}/{group_name}/{endpoint_group_type}/{endpoint_group_name}/services/{service_name}/endpoints/{endpoint_name}?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{endpoint_group_type}/{endpoint_group_name}/services/{service_name}/endpoints?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{endpoint_group_type}/{endpoint_group_name}/services/{service_name}/endpoints/{endpoint_name}?[start_time]&[end_time]&[granularity]
```

Request endpoint a/r under specific endpoint group:

```
/results/{report_name}/{group_type}/{group_name}/{endpoint_group_type}/{endpoint_group_name}/endpoints?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{group_type}/{group_name}/{endpoint_group_type}/{endpoint_group_name}/endpoints/{endpoint_name}?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{endpoint_group_type}/{endpoint_group_name}/endpoints?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{endpoint_group_type}/{endpoint_group_name}/endpoints/{endpoint_name}?[start_time]&[end_time]&[granularity]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start_time]`  | UTC time in W3C format                                                                          | YES      |
| `[end_time]`    | UTC time in W3C format                                                                          | YES      |
| `[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly`, `daily` or `custom` | NO       | `daily`       |

#### Path Parameters

| Name                    | Description                                                                                           | Required | Default value |
| ----------------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report_name}`         | Name of the report that contains all the information about the profile, filter tags, group types etc. | YES      |
| `{group_type}`          | Type of the Group of Endpoint Groups.                                                                 | NO       |
| `{group_name}`          | Name of the Group of Endpoint Groups.                                                                 | NO       |
| `{endpoint_group_type}` | Type of the the Endpoint Group.                                                                       | YES      |
| `{endpoint_group_name}` | Name of the the Endpoint Group.                                                                       | YES      |
| `{service_name}`        | Name of the specific service.                                                                         | NO       |
| `{endpoint_name}`       | Name of the specific endpoint.                                                                        | NO       |

#### Headers

```
x-api-key: "tenant_key_value"
Accept: "application/xml" or "application/json"
```

##### Response

```
Status: 200 OK
```

## Request endpoint a/r under service: `service_a`

#### URL

`/api/v2/results/Report_A/SITE/ST01/services/service_a/endpoints?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:23:59Z&granularity=daily`

#### Response Body

```
{
   "results": [
     {
       "name": "ST01",
       "type": "SITE",
       "serviceflavors": [
         {
           "name": "service_a",
           "type": "service",
           "endpoints": [
             {
               "name": "e01",
               "type": "endpoint",
               "results": [
                 {
                   "timestamp": "2015-06-22",
                   "availability": "98.26389",
                   "reliability": "98.26389",
                   "unknown": "0",
                   "uptime": "0.98264",
                   "downtime": "0"
                 },
                 {
                   "timestamp": "2015-06-23",
                   "availability": "54.03509",
                   "reliability": "81.48148",
                   "unknown": "0.01042",
                   "uptime": "0.53472",
                   "downtime": "0.33333"
                 }
               ]
             },
             {
               "name": "e02",
               "type": "endpoint",
               "results": [
                 {
                   "timestamp": "2015-06-22",
                   "availability": "96.875",
                   "reliability": "96.875",
                   "unknown": "0",
                   "uptime": "0.96875",
                   "downtime": "0"
                 },
                 {
                   "timestamp": "2015-06-23",
                   "availability": "100",
                   "reliability": "100",
                   "unknown": "0",
                   "uptime": "1",
                   "downtime": "0"
                 }
               ]
             }
           ]
         }
       ]
     }
   ]
 }
```

## Request endpoint a/r under endpoint group: `ST01`

#### URL

`/api/v2/results/Report_A/SITE/ST01/endpoints?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:23:59Z&granularity=daily`

#### Response Body

```
{
   "results": [
     {
       "name": "ST01",
       "type": "SITE",
       "serviceflavors": [
         {
           "name": "service_a",
           "type": "service",
           "endpoints": [
             {
               "name": "e01",
               "type": "endpoint",
               "results": [
                 {
                   "timestamp": "2015-06-22",
                   "availability": "98.26389",
                   "reliability": "98.26389",
                   "unknown": "0",
                   "uptime": "0.98264",
                   "downtime": "0"
                 },
                 {
                   "timestamp": "2015-06-23",
                   "availability": "54.03509",
                   "reliability": "81.48148",
                   "unknown": "0.01042",
                   "uptime": "0.53472",
                   "downtime": "0.33333"
                 }
               ]
             },
             {
               "name": "e02",
               "type": "endpoint",
               "results": [
                 {
                   "timestamp": "2015-06-22",
                   "availability": "96.875",
                   "reliability": "96.875",
                   "unknown": "0",
                   "uptime": "0.96875",
                   "downtime": "0"
                 },
                 {
                   "timestamp": "2015-06-23",
                   "availability": "100",
                   "reliability": "100",
                   "unknown": "0",
                   "uptime": "1",
                   "downtime": "0"
                 }
               ]
             }
           ]
         }
       ]
     }
   ]
 }
```

## Request endpoint a/r for specific endpoint `e01` under endpoint group: `ST01`

#### URL

`/api/v2/results/Report_A/SITE/ST01/services/service_a/endpoints/e01?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:23:59Z&granularity=daily`

#### Response Body

```
{
   "results": [
     {
       "name": "ST01",
       "type": "SITE",
       "serviceflavors": [
         {
           "name": "service_a",
           "type": "service",
           "endpoints": [
             {
               "name": "e01",
               "type": "endpoint",
               "results": [
                 {
                   "timestamp": "2015-06-22",
                   "availability": "98.26389",
                   "reliability": "98.26389",
                   "unknown": "0",
                   "uptime": "0.98264",
                   "downtime": "0"
                 },
                 {
                   "timestamp": "2015-06-23",
                   "availability": "54.03509",
                   "reliability": "81.48148",
                   "unknown": "0.01042",
                   "uptime": "0.53472",
                   "downtime": "0.33333"
                 }
               ]
             }
           ]
         }
       ]
     }
   ]
 }
```

### Extra information for a specific endpoint on endpoint a/r

Some service endpoint a/r have additional information regarding the specific service endpoint such as it's Url, certificate DN etc... If this information is available it will be displayed under each service endpoint along with the a/r results. For example:

```
{
  "results": [
    {
      "name": "ST01",
      "type": "SITE",
      "serviceflavors": [
        {
          "name": "service_a",
          "type": "service",
          "endpoints": [
            {
              "name": "e01",
              "type": "endpoint",
              "info": {
                "Url": "https://foo.example.url"
              },
              "results": [
                {
                  "timestamp": "2015-06-22",
                  "availability": "98.26389",
                  "reliability": "98.26389",
                  "unknown": "0",
                  "uptime": "0.98264",
                  "downtime": "0"
                },
                {
                  "timestamp": "2015-06-23",
                  "availability": "54.03509",
                  "reliability": "81.48148",
                  "unknown": "0.01042",
                  "uptime": "0.53472",
                  "downtime": "0.33333"
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```


## [GET]: Flat List Availabilities and Reliabilities for all service Endpoints {#5}

The following methods can be used to obtain a tenant's flat list of all service endpoints Availability and Reliability result. The api authenticates the tenant using the api-key within the x-api-key header. The user can specify time granularity (`monthly`, `daily` or `custom`) for retrieved results and also format using the `Accept` header. Pagination is also supported by using the optional parameters `pageSize` to define the size of each result page and `nextPageToken` to proceed to the next available page of results.
## [GET] Endpoints A/R

### Input

Request a flat list of all endpoint a/r 

```
/results/{report_name}/endpoints?[start_time]&[end_time]&[granularity]&[pageSize]&[nextPageToken]
```


#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start_time]`  | UTC time in W3C format                                                                          | YES      |
| `[end_time]`    | UTC time in W3C format                                                                          | YES      |
| `[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly`, `daily` or `custom` | NO       | `daily`       |
| `[pageSize]` | How many results to return per request (-1 means return all results) | NO       | -1       |
| `[nextPageToken]` | token to proceed to the next page | NO       |  |

#### Path Parameters

| Name                    | Description                                                                                           | Required | Default value |
| ----------------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report_name}`         | Name of the report that contains all the information about the profile, filter tags, group types etc. | YES      |

#### Headers

```
x-api-key: "tenant_key_value"
Accept: "application/xml" or "application/json"
```

##### Response

```
Status: 200 OK
```

## Request endpoint a/r under service: `service_a`

#### URL

`/api/v2/results/{report_name}/endpoints?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:23:59Z&granularity=daily&pageSize=2`

#### Response Body

```
{
 "endpoints": [
  {
   "name": "host01.example.foo",
   "service": "SERV-A",
   "supergroup": "SITE-A",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T08:47:00Z",
     "value": "WARNING"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  }
 ],
 "nextPageToken": "Mg==",
 "pageSize": 2
}
```

## Request to see next page of results

#### URL

`/api/v2/results/{report_name}/endpoints?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:23:59Z&granularity=daily&pageSize=2&nextPageToken=Mg==`

#### Response Body

```
{
 "endpoints": [
  {
   "name": "host02.example.foo",
   "service": "SERV-B",
   "supergroup": "SITE-A",
   "statuses": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T12:56:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  }
 ],
 "pageSize": 2
}
```
