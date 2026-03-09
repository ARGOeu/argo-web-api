---
id: v4_nodes
title: Node Capabilities
sidebar_position: 1
---

## Node capabilities

A tenant can be defined to represent a specific node by adding the following extra information in the body of a tenant:
```json
"node" : {
  "id": "node-id",
  "name": "node-name"
}
```

When a tenant has been configured with node information the node's capabilities become available through the `/api/v4/nodes` calls.
Capabilities include the following:
- availability (get availability results for the node's services)
- uptime (get uptime results for the node's services)
- status (get status results for the node's services)

## API Calls

_Note_: The following calls implement the node capabilities

| Name                                                                          | Description                                                                                                                                                                                                                              | Shortcut          |
| ----------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------- |
| GET: Availability results for the node's services | This method retrieves by default the latest daily availability for all node's services | [Description](#1) |
| GET: Uptime results for the node's services | This method retrieves by default the latest daily availability for all node's services | [Description](#2) |
| GET: Status results for the node's services | This method retrieves by default the latest daily availability for all node's services | [Description](#3) |


## [GET]: Availability results for the node's services {#1}

The following methods can be used to obtain availability results for a node's services. The api authenticates the tenant using the api-key within the x-api-key header. User can specify time granularity (`monthly`, `daily`) for retrieved results and also format using the `Accept` header. 

### Input

```
/nodes/{node_name}/capabilities/availability?([start_time]&[end_time]|[start_date]&[end_date]|[date])&[granularity]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[date]`        | UTC date in YYYY-MM-DD format                                                                   | NO       |               |
| `[start_time]`  | UTC time in W3C format                                                                          | NO       |               |
| `[end_time]`    | UTC time in W3C format                                                                          | NO       |               |
| `[start_date]`  | UTC date in YYYY-MM-DD format                                                                   | NO       |               |
| `[start_date]`  | UTC date in YYYY-MM-DD format                                                                   | NO       |               |
| `[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly`,  `daily`  | NO       | `daily`       |

- If a user doesn't specify any query parameter the api returns the latest daily availability results.
- If a user specifies the date parameter the api retuerns the daily availability results for that date
- A user can specify the period of availability results by using start_date and end_date instead of start_time and end_time

#### Path Parameters

| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{node_name}` | Name of the node | YES      |


### Example Request 1: Get default daily availability

```
/api/v4/nodes/NODE_A/capabilities/availability`
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
      "name": "ST01",
      "results": [
        {
          "date": "2015-06-26",
          "availability": "99"
        }
      ]
    },
    {
      "name": "ST02",
      "results": [
        {
          "date": "2015-06-26",
          "availability": "98"
        }
      ]
    }
  ]
}
```


### Example Request 2: daily availability over a period

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v4/nodes/NODE_A/capabilities/availability?start_date=2015-06-20Z&end_date=2015-06-22T&granularity=daily`
```

or

```
/api/v4/nodes/NODE_A/capabilities/availability?start_time=2015-06-20T12:00:00Z&end_time=2015-06-22T23:00:00Z&granularity=daily`
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
       "name": "ST01",
       "results": [
         {
           "date": "2015-06-20",
           "availability": "100"
         },
         {
           "date": "2015-06-21",
           "availability": "100"
         },
         {
           "date": "2015-06-22",
           "availability": "100"
         }
       ]
     },
     {
       "name": "ST02",
       "results": [
         {
           "date": "2015-06-20",
           "availability": "100"
         },
         {
           "date": "2015-06-21",
           "availability": "98"
         },
         {
           "date": "2015-06-22",
           "availability": "98"
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
/api/v3/results/Report_A?start_time=2015-07-01T12:00:00Z&end_time=2015-07-31T23:00:00Z&granularity=monthly
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
      "name": "ST01",
      "results": [
        {
          "date": "2015-06",
          "availability": "99.99999900000002"
        }
      ]
    },
    {
      "name": "ST02",
      "results": [
        {
          "date": "2015-06",
          "availability": "99.99999900000002"
        }
      ]
    }
  ]
}  
```


## [GET]: Uptime results for the node's services {#2}

The following methods can be used to obtain uptime results for a node's services. The api authenticates the tenant using the api-key within the x-api-key header. User can specify time granularity (`monthly`, `daily`) for retrieved results and also format using the `Accept` header. 

### Input

```
/nodes/{node_name}/capabilities/uptime?([start_time]&[end_time]|[start_date]&[end_date]|[date])&[granularity]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[date]`        | UTC date in YYYY-MM-DD format                                                                   | NO       |               |
| `[start_time]`  | UTC time in W3C format                                                                          | NO       |               |
| `[end_time]`    | UTC time in W3C format                                                                          | NO       |               |
| `[start_date]`  | UTC date in YYYY-MM-DD format                                                                   | NO       |               |
| `[start_date]`  | UTC date in YYYY-MM-DD format                                                                   | NO       |               |
| `[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly`,  `daily`  | NO       | `daily`       |

- If a user doesn't specify any query parameter the api returns the latest daily uptime results.
- If a user specifies the date parameter the api retuerns the daily uptime results for that date
- A user can specify the period of uptime results by using start_date and end_date instead of start_time and end_time

#### Path Parameters

| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{node_name}` | Name of the node | YES      |


### Example Request 1: Get default daily uptime

```
/api/v4/nodes/NODE_A/capabilities/uptime`
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
      "name": "ST01",
      "results": [
        {
          "date": "2015-06-26",
          "uptime": "99"
        }
      ]
    },
    {
      "name": "ST02",
      "results": [
        {
          "date": "2015-06-26",
          "uptime": "98"
        }
      ]
    }
  ]
}
```


### Example Request 2: daily uptime over a period

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v4/nodes/NODE_A/capabilities/uptime?start_date=2015-06-20Z&end_date=2015-06-22T&granularity=daily`
```

or

```
/api/v4/nodes/NODE_A/capabilities/uptime?start_time=2015-06-20T12:00:00Z&end_time=2015-06-22T23:00:00Z&granularity=daily`
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
       "name": "ST01",
       "results": [
         {
           "date": "2015-06-20",
           "uptime": "1"
         },
         {
           "date": "2015-06-21",
           "uptime": "1"
         },
         {
           "date": "2015-06-22",
           "uptime": "1"
         }
       ]
     },
     {
       "name": "ST02",
       "results": [
         {
           "date": "2015-06-20",
           "uptime": "1"
         },
         {
           "date": "2015-06-21",
           "uptime": "1"
         },
         {
           "date": "2015-06-22",
           "uptime": "1"
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
/api/v3/results/Report_A?start_time=2015-07-01T12:00:00Z&end_time=2015-07-31T23:00:00Z&granularity=monthly
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
      "name": "ST01",
      "results": [
        {
          "date": "2015-06",
          "uptime": "1"
        }
      ]
    },
    {
      "name": "ST02",
      "results": [
        {
          "date": "2015-06",
          "uptime": "1"
        }
      ]
    }
  ]
}  
```



## [GET]: Status results for the node's services {#3}

The following methods can be used to obtain status results for a node's services. The api authenticates the tenant using the api-key within the x-api-key header. for retrieved results and also format using the `Accept` header. 

### Input

```
/nodes/{node_name}/capabilities/uptime?[start_time]&[end_time]&[history]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start_time]`  | UTC time in W3C format                                                                          | NO       |               |
| `[end_time]`    | UTC time in W3C format                                                                          | NO       |               |
| `[history]`     | Show full history of status timelines                                                           | NO       |               |        

- If a user doesn't specify any query parameter the api returns the status of the services

#### Path Parameters

| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{node_name}` | Name of the node | YES      |


### Example Request 1: Get latest status

```
/api/v4/nodes/NODE_A/capabilities/status`
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
   "name": "SITEA",
   "results": [
    {
     "timestamp": "2015-05-01T17:53:00Z",
     "value": "OK"
    }
   ]
  },
  {
   "name": "SITEB",
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


### Example Request 2: get status timelines period

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v4/nodes/NODE_A/capabilities/uptime?start_date=2015-05-01Z&end_date=2015-05-01T&granularity=daily`
```

or

```
/api/v4/nodes/NODE_A/capabilities/uptime?start_time=2015-05-01T12:00:00Z&end_time=2015-05-01T23:00:00Z&granularity=daily`
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
   "name": "SITEA",
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
  },
  {
   "name": "SITEB",
   "results": [
    {
     "timestamp": "2015-05-01T00:00:00Z",
     "value": "OK"
    },
    {
     "timestamp": "2015-05-01T03:00:00Z",
     "value": "WARNING"
    },
    {
     "timestamp": "2015-05-01T17:53:00Z",
     "value": "CRITICAL"
    },
    {
     "timestamp": "2015-05-01T23:59:59Z",
     "value": "CRITICAL"
    }
   ]
  }
 ]
}
```

