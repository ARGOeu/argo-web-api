---
id: v4_groups
title: Group Resoults and Statuses
sidebar_position: 1
---

## Group Results and Statuses

A user can quickly retrive latest results and statuses about the groups of a tenant by using the following calls. If the tenant has one only report this will be selected automatically otherwise the user must specifiy the report by using the `?report=` url parameter. User can query for all groups or a specific one.


## API Calls



| Name                                                                          | Description                                                                                                                                                                                                                              | Shortcut          |
| ----------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------- |
| GET: Group Results | This method retrieves by default the latest daily availability and uptime for all tenant groups or a specific one | [Description](#1) |
| GET: Group Statuses | This method retrieves by default the latest statuses for all tenant groups or a specific one | [Description](#2) |



## [GET]: Group Results {#A}

The following methods can be used to obtain availability and uptime results for a tenant's groups. The api authenticates the tenant using the api-key within the x-api-key header. User can specify time granularity (`monthly`, `daily`) for retrieved results and also format using the `Accept` header. 

The period can be set by either:
- `date` parameter for a specific date
- `start_date` & `end_date` parameters for specific days
- `period` parameter that accept values such as `7d` to go 7 days back or `2w` to go two weeks back or `1m` to go one month back

The report can be specified by using the `report` parameter

### Input

```
/api/v4/results/groups/{group_name}?([start_time]&[end_time]|[start_date]&[end_date]|[date]|[period])&[granularity]&[report]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[date]`        | UTC date in YYYY-MM-DD format                                                                   | NO       |               |
| `[period]`      | Specify the lookback window in days or weeks (e.g., 7d or 2w)                                   | NO       |               |
| `[start_time]`  | UTC time in W3C format                                                                          | NO       |               |
| `[end_time]`    | UTC time in W3C format                                                                          | NO       |               |
| `[start_date]`  | UTC date in YYYY-MM-DD format                                                                   | NO       |               |
| `[start_date]`  | UTC date in YYYY-MM-DD format                                                                   | NO       |               |
| `[report]`      | Name of the report                                                                              | Only if multiple reports exists       |               |    
| `[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly`,  `daily`  | NO       | `daily`       |

- If a user doesn't specify any query parameter the api returns the latest daily results.
- If a user specifies the date parameter the api returns the daily  results for that date
- A user can specify the period of results by using start_date and end_date instead of start_time and end_time
- A user can specify the latest period of results by using the period parameter with example values such as 7d (7 days), 3w (3 weeks)

#### Path Parameters

| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{group_name}`| Target a specific group. If omitted you get a result list for all the groups                | NO       |               |



### Example Request 1: Get default daily results for group ST101

```
/api/v4/results/groups/ST101`
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
  "data": [
    {
      "name": "ST101",
      "results": [
        {
          "date": "2015-06-26",
          "availability": "99",
          "uptime": "1"
        }
      ]
    }
  ]
}
```


### Example Request 2: daily results over a period

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v4/results/groups/ST101?start_date=2015-06-20Z&end_date=2015-06-22T&granularity=daily&report=CORE`
```

or

```
/api/v4/results/groups/ST101?start_time=2015-06-20T12:00:00Z&end_time=2015-06-22T23:00:00Z&granularity=daily&report=CORE`
```

or 

```
/api/v4/results/groups/ST101?period=3d&report=CORE`
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
   "data": [
     {
       "name": "SERVICE101",
       "results": [
         {
           "date": "2015-06-20",
           "availability": "100",
           "uptime": "1"
         },
         {
           "date": "2015-06-21",
           "availability": "100",
           "uptime": "1"
         },
         {
           "date": "2015-06-22",
           "availability": "100",
           "uptime": "1"
         }
       ]
     }
   ]
 }
```

### Example Request 3: monthly granularity

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v4/results/groups/ST101?start_time=2015-07-01T12:00:00Z&end_time=2015-07-31T23:00:00Z&granularity=monthly&report=CORE
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
  "data": [
    {
      "name": "ST101",
      "results": [
        {
          "date": "2015-06",
          "availability": "99.99999900000002",
          "uptime": "1"
        }
      ]
    }
  ]
}  
```

### Example Request 4: Get default daily availability for all the tenant's groups

```
/api/v4/results/groups`
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
  "data": [
    {
      "name": "ST101",
      "results": [
        {
          "date": "2015-06-26",
          "availability": "99",
          "uptime": "0.99"
        }
      ]
    },
    {
      "name": "ST202",
      "results": [
        {
          "date": "2015-06-26",
          "availability": "98",
          "uptime": "0.98"
        }
      ]
    }
  ]
}
```



## [GET]: Status results for the tenant's groups {#2}

The following methods can be used to obtain status results for a tenant's groups. The api authenticates the tenant using the api-key within the x-api-key header. for retrieved results and also format using the `Accept` header. 


The report can be specified by using the `report` parameter

### Input

```
/api/v4/status/groups/{group_name}?[start_time]&[end_time]&[history]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start_time]`  | UTC time in W3C format                                                                          | NO       |               |
| `[end_time]`    | UTC time in W3C format                                                                          | NO       |               |
| `[history]`     | Show full history of status timelines                                                           | NO       |               |        
| `[report]`      | Name of the report                                                                              | Only if multiple reports exists       |               |  
- If a user doesn't specify any query parameter the api returns the status of the services

#### Path Parameters

| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{group_name}`| Target a specific group. If omitted you get a result list for all the groups                | NO       |               |


### Example Request 1: Get latest status of a specific group

```
/api/v4/status/groups/SERVICE101?report=CORE`
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
 "data": [
  {
   "name": "ST101",
   "results": [
    {
     "timestamp": "2015-05-01T17:53:00Z",
     "value": "OK"
    }
   ]
  }
 ]
}
```


### Example Request 2: get status timelines period

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v4/status/groups/ST101?report=CORE&start_date=2015-05-01Z&end_date=2015-05-01T&granularity=daily`
```

or

```
/api/v4/status/groups/ST101?report=CORE&start_time=2015-05-01T12:00:00Z&end_time=2015-05-01T23:00:00Z`
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
 "data": [
  {
   "name": "ST101",
   "results": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T01:00:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T05:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "OK"
    }
   ]
  }
 ]
}
```

### Example Request 3: Get latest status for all tenant's groups

```
/api/v4/status/groups?report=CORE
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
 "data": [
  {
   "name": "ST101",
   "results": [
    {
     "timestamp": "2015-05-01T17:53:00Z",
     "value": "OK"
    }
   ]
  },
  {
   "name": "ST202",
   "results": [
    {
     "timestamp": "2015-05-01T17:53:00Z",
     "value": "CRITICAL"
    }
   ]
  }
 ]
}
```

