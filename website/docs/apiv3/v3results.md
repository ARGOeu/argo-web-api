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

<a id="1"></a>

# [GET]: List Availability and Reliability results for top level supergroups and included groups

The following methods can be used to obtain a tenant's Availability and Reliability result metrics for all top level supergroups and included groups. The api authenticates the tenant using the api-key within the x-api-key header. User can specify time granularity (`monthly` or `daily`) for retrieved results and also format using the `Accept` header. 

### Input

```
/results/{report_name}?[start_time]&[end_time]&[granularity]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start_time]`  | UTC time in W3C format                                                                          | YES      |
| `[end_time]`    | UTC time in W3C format                                                                          | YES      |
| `[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly` or `daily` | NO       | `daily`       |

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
/api/v2/results/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z 
```
or 
```
/api/v2/results/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily`
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
          "reliability": "50.413931144915935"
        },
        {
          "date": "2015-06-23",
          "availability": "75.36324059247399",
          "reliability": "80.8138510808647"
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
              "reliability": "54.6",
              "unknown": "0",
              "uptime": "1",
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
              "reliability": "45",
              "unknown": "0",
              "uptime": "1",
              "downtime": "0"
            },
            {
              "date": "2015-06-23",
              "availability": "43.5",
              "reliability": "56",
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

### Example Request 2: monthly granularity

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v2/results/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=monthly
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
          "availability": "71.75110088070457",
          "reliability": "65.61389111289031"
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

<a id="2"></a>

# [GET]: List Availability and Reliability results for endpoints with specific resource-id

The following methods can be used to obtain a tenant's Availability and Reliability result for the endpoints that have a specific resource-id. User can specify a period with `start_time` and `end_time` and granularity(`monthly` or `daily`) for retrieved results. `Accept` header is required. 

### Input

```
/results/{report_name}/id/{resource-id}?[start_time]&[end_time]&[granularity]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start_time]`  | UTC time in W3C format                                                                          | YES      |
| `[end_time]`    | UTC time in W3C format                                                                          | YES      |
| `[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly` or `daily` | NO       | `daily`       |

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
/api/v2/results/Report_A/id/simple-queue?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z 
```
or 
```
/api/v2/results/Report_A/id/simple-queue?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily`
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
  ]
}
```

### Example Request 2: monthly granularity with specific resource-id

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v2/results/Report_A/id/simple-queue?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=monthly
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
  ]
}
```
