---
id: v3_ar_results
title: Availability / Reliability Results (v3)
sidebar_position: 1
---

## API Calls

_Note_: These are v3 api calls implementations found under the path `/api/v3`

| Name                                                                          | Description                                                                                                                                                                                                                              | Shortcut          |
| ----------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------- |
| GET: List Availability and Reliability results for top level supergroups and included groups | This method retrieves the a/r results of all top level supergroups and their included groups | [Description](#1) |
| GET: List Availability and Reliability results for specific endpoint using resource id | This method retrieves the a/r results of all endoints with specific resource id | [Description](#2) |



## [GET]: List Availability and Reliability results for top level supergroups and included groups {#1}

The following methods can be used to obtain a tenant's Availability and Reliability result metrics for all top level supergroups and included groups. The api authenticates the tenant using the api-key within the x-api-key header. User can specify time granularity (`monthly`, `daily` or `custom`) for retrieved results and also format using the `Accept` header. 

### Input

```
/results/{report_name}?[start_time]&[end_time]&[granularity]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start_time]`  | UTC time in W3C format                                                                          | YES      |
| `[end_time]`    | UTC time in W3C format                                                                          | YES      |
| `[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly`,  `daily` or `custom` | NO       | `daily`       |

#### Path Parameters

| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report_name}` | Name of the report that contains all the information about the profile, filter tags, group types etc. | YES      |


### Example Request 1: default daily granularity

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v3/results/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z 
```
or 
```
/api/v3/results/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily`
```

##### Headers

```
x-api-key: "tenant_key_value"
Accept: "application/json"
```

#### Response

##### Code

```
Status: 200 OK
```

##### Body

```json
{
  "results": [
    {
      "name": "GROUP_A",
      "type": "GROUP",
      "results": [
        {
          "date": "2015-06-22",
          "availability": "68.13896116893515",
          "reliability": "68.13896116893515"
        },
        {
          "date": "2015-06-23",
          "availability": "75.36324059247399",
          "reliability": "75.36324059247399"
        }
      ],
      "groups": [
        {
          "name": "ST01",
          "type": "SITES",
          "results": [
            {
              "date": "2015-06-22",
              "availability": "66.7",
              "reliability": "66.7",
              "unknown": "0",
              "uptime": "66.7",
              "downtime": "0"
            },
            {
              "date": "2015-06-23",
              "availability": "100",
              "reliability": "100",
              "unknown": "0",
              "uptime": "1",
              "downtime": "0"
            }
          ]
        },
        {
          "name": "ST02",
          "type": "SITES",
          "results": [
            {
              "date": "2015-06-22",
              "availability": "70",
              "reliability": "70",
              "unknown": "0",
              "uptime": "0.70",
              "downtime": "0"
            },
            {
              "date": "2015-06-23",
              "availability": "43.5",
              "reliability": "43.5",
              "unknown": "0",
              "uptime": "0.435",
              "downtime": "0"
            }
          ]
        }
      ]
    }
  ]
}
```

### Example Request 2: monthly granularity

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v3/results/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=monthly
```
##### Headers

```
x-api-key: "tenant_key_value"
Accept: "application/json"
```

#### Response

##### Code

```
Status: 200 OK
```

##### Body

```json
{
  "results": [
    {
      "name": "GROUP_A",
      "type": "GROUP",
      "results": [
        {
          "date": "2015-06",
          "availability": "99.99999900000002",
          "reliability": "99.99999900000002"
        }
      ],
      "groups": [
        {
          "name": "ST01",
          "type": "SITES",
          "results": [
            {
              "date": "2015-06",
              "availability": "99.99999900000002",
              "reliability": "99.99999900000002",
              "unknown": "0",
              "uptime": "1",
              "downtime": "0"
            }
          ]
        },
        {
          "name": "ST02",
          "type": "SITES",
          "results": [
            {
              "date": "2015-06",
              "availability": "99.99999900000002",
              "reliability": "99.99999900000002",
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
```

### Example Request 3: Custom granularity
This request returns availability/reliability score numbers for the whole custom period defined between `start_time` and `end_time`. 
This means that for each item the user will receive one availability and reliability result concerning the whole period (instead of multiple daily or monthly results)

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v3/results/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=custom
```
##### Headers

```
x-api-key: "tenant_key_value"
Accept: "application/json"
```

#### Response

##### Code

```
Status: 200 OK
```

##### Body

```json
{
  "results": [
    {
      "name": "GROUP_A",
      "type": "GROUP",
      "results": [
        {
          "availability": "99.99999900000002",
          "reliability": "99.99999900000002"
        }
      ],
      "groups": [
        {
          "name": "ST01",
          "type": "SITES",
          "results": [
            {
              "availability": "99.99999900000002",
              "reliability": "99.99999900000002",
              "unknown": "0",
              "uptime": "1",
              "downtime": "0"
            }
          ]
        },
        {
          "name": "ST02",
          "type": "SITES",
          "results": [
            {
              "availability": "99.99999900000002",
              "reliability": "99.99999900000002",
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
```



## [GET]: List Availability and Reliability results for endpoints with specific resource-id {#2}

The following methods can be used to obtain a tenant's Availability and Reliability result for the endpoints that have a specific resource-id. User can specify a period with `start_time` and `end_time` and granularity(`monthly`, `daily` or `custom`) for retrieved results. `Accept` header is required. 

### Input

```
/results/{report_name}/id/{resource-id}?[start_time]&[end_time]&[granularity]
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
| `{id}` | The resource id | YES      |

### Example Request 1: default daily granularity with specific resource-id

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v3/results/Report_A/id/simple-queue?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z 
```
or 
```
/api/v3/results/Report_A/id/simple-queue?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily`
```

##### Headers

```
x-api-key: "tenant_key_value"
Accept: "application/json"
```

#### Response

##### Code

```
Status: 200 OK
```

##### Body

```json
{
  "id": "simple-queue",
  "endpoints": [
    {
      "name": "host01.example",
      "service": "service.queue",
      "group": "Infra-01",
      "info": {
        "URL": "http://submit.queue01.example.com"
      },
      "results": [
        {
          "date": "2015-06-22",
          "availability": "99.99999900000002",
          "reliability": "99.99999900000002",
          "unknown": "0",
          "uptime": "1",
          "downtime": "0"
        },
        {
          "date": "2015-06-23",
          "availability": "99.99999900000002",
          "reliability": "99.99999900000002",
          "unknown": "0",
          "uptime": "1",
          "downtime": "0"
        }
      ]
    }
  ]
}
```

### Example Request 2: monthly granularity with specific resource-id

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v3/results/Report_A/id/simple-queue?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=monthly
```
##### Headers

```
x-api-key: "tenant_key_value"
Accept: "application/json"
```

#### Response

##### Code

```
Status: 200 OK
```

##### Body

```json
{
  "id": "simple-queue",
  "endpoints": [
    {
      "name": "host01.example",
      "service": "service.queue",
      "group": "Infra-01",
      "info": {
        "URL": "http://submit.queue01.example.com"
      },
      "results": [
        {
          "date": "2015-06",
          "availability": "99.99999900000002",
          "reliability": "99.99999900000002",
          "unknown": "0",
          "uptime": "1",
          "downtime": "0"
        }
      ]
    }
  ]
}
```

### Example Request 3: custom period granularity with specific resource-id
This request returns availability/reliability score numbers for the whole custom period defined between `start_time` and `end_time`. 
This means that for each item with the specific resource-id the user will receive one availability and reliability result concerning the whole period (instead of multiple daily or monthly results)

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v3/results/Report_A/id/simple-queue?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=custom
```
##### Headers

```
x-api-key: "tenant_key_value"
Accept: "application/json"
```

#### Response

##### Code

```
Status: 200 OK
```

##### Body

```json
{
  "id": "simple-queue",
  "endpoints": [
    {
      "name": "host01.example",
      "service": "service.queue",
      "group": "Infra-01",
      "info": {
        "URL": "http://submit.queue01.example.com"
      },
      "results": [
        {
          "availability": "99.99999900000002",
          "reliability": "99.99999900000002",
          "unknown": "0",
          "uptime": "1",
          "downtime": "0"
        }
      ]
    }
  ]
}
```
