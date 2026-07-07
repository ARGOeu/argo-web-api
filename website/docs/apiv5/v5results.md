---
id: v5_ar_results
title: Results (v5)
sidebar_position: 1
---

## API Calls

_Note_: These are v5 api calls implementations found under the path `/api/v5`

| Name                                                                          | Description                                                                                                                                                                                                                              | Shortcut          |
| ----------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------- |
| Get results for top level supergroups| This method retrieves the a/r results of all top level supergroups | [Description](#1) |
| Get results for specific supergroup| This method retrieves the a/r results for a specific supergroup | [Description](#2) |
| Get results for top level groups| This method retrieves the a/r results of all top level groups | [Description](#3) |
| Get results for specific group| This method retrieves the a/r results for a specific group | [Description](#4) |



## Get results for top level supergroups {#1}

The following method can be used to obtain a tenant's Availability and Reliability result metrics for all top level supergroups. The api authenticates the tenant using the api-key within the x-api-key header (or using an admin key along with `x-tenant-id` param). User can specify time granularity (`monthly`, `daily` or `custom`) for retrieved results 

### Input


```
HTTP GET /api/v5/results/{report-name}/supergroups?[start-time]&[end-time]&[granularity]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start-time]`  | UTC time in W3C format                                                                          | YES      |
| `[end-time]`    | UTC time in W3C format                                                                          | YES      |
| `[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly`,  `daily` or `custom` | NO       | `daily`       |

#### Path Parameters

| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report-name}` | Name of the report that contains the results | YES      |


### Example Request 1: default daily granularity

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v5/results/Report_A/supergroups?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z 
```
or 
```
/api/v3/results/Report_A/supergroups?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily`
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
      "name": "PROJECT_A",
      "type": "PROJECT",
      "results": [
        {
          "date": "2015-06-22",
          "availability": 68.14,
          "reliability": 68.14
        },
        {
          "date": "2015-06-23",
          "availability": 75.36,
          "reliability": 75.36
        }
      ]
    },
    {
      "name": "PROJECT_B",
      "type": "PROJECT",
      "results": [
        {
          "date": "2015-06-22",
          "availability": 33.13,
          "reliability": 55.13
        },
        {
          "date": "2015-06-23",
          "availability": 23.36,
          "reliability": 88.36
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
/api/v5/results/Report_A/supergroups?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=monthly
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
      "name": "PROJECT_A",
      "type": "PROJECT",
      "results": [
        {
          "date": "2015-06",
          "availability": 99.99,
          "reliability": 99.98
        }
      ]
    },
    {
      "name": "PROJECT_B",
      "type": "PROJECT",
      "results": [
        {
          "date": "2015-06",
          "availability": 45.99,
          "reliability": 33.32
        }
      ]
    }
  ]
}
```

### Example Request 3: Custom granularity
This request returns availability/reliability score numbers for the whole custom period defined between `start-time` and `end-time`. 
This means that for each item the user will receive one availability and reliability result concerning the whole period (instead of multiple daily or monthly results)

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v5/results/Report_A/supergroups?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=custom
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
      "name": "PROJECT_A",
      "type": "PROJECT",
      "results": [
        {
          "date": "2015-06",
          "availability": 99.99,
          "reliability": 99.98
        }
      ]
    },
    {
      "name": "PROJECT_B",
      "type": "PROJECT",
      "results": [
        {
          "date": "2015-06",
          "availability": 33.99,
          "reliability": 33.98
        }
      ]
    }
  ]
}
```

## Get results for specific top level supergroup {#2}

The following method can be used to obtain a tenant's Availability and Reliability result metrics for a specific top level supergroup. The api authenticates the tenant using the api-key within the x-api-key header (or using an admin key along with `x-tenant-id` param). User can specify time granularity (`monthly`, `daily` or `custom`) for retrieved results 

### Input

```
HTTP GET /api/v5/results/{report-name}/supergroups/{supergroup-name}?[start-time]&[end-time]&[granularity]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start-time]`  | UTC time in W3C format                                                                          | YES      |
| `[end-time]`    | UTC time in W3C format                                                                          | YES      |
| `[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly`,  `daily` or `custom` | NO       | `daily`       |

#### Path Parameters

| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report-name}` | Name of the report that contains the results | YES      |
| `{supergroup-name}` | Name of the specific supergroup | YES      |


### Example Request 1: default daily granularity

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v5/results/Report_A/supergroups/PROJECT-A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z 
```
or 
```
/api/v3/results/Report_A/supergroups/PROJECT-A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily`
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
      "name": "PROJECT_A",
      "type": "PROJECT",
      "results": [
        {
          "date": "2015-06-22",
          "availability": 68.14,
          "reliability": 68.14
        },
        {
          "date": "2015-06-23",
          "availability": 75.36,
          "reliability": 75.36
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
/api/v5/results/Report_A/supergroups/PROJECT-A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=monthly
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
      "name": "PROJECT_A",
      "type": "PROJECT",
      "results": [
        {
          "date": "2015-06",
          "availability": 99.99,
          "reliability": 99.98
        }
      ]
    }
  ]
}
```

### Example Request 3: Custom granularity
This request returns availability/reliability score numbers for the whole custom period defined between `start-time` and `end-time`. 
This means that for each item the user will receive one availability and reliability result concerning the whole period (instead of multiple daily or monthly results)

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v5/results/Report_A/supergroups/PROJECT-A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=custom
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
      "name": "PROJECT_A",
      "type": "PROJECT",
      "results": [
        {
          "date": "2015-06",
          "availability": 99.99,
          "reliability": 99.98
        }
      ]
    }
  ]
}
```



## Get results for groups {#3}

The following method can be used to obtain a tenant's Availability and Reliability result metrics for all groups. The api authenticates the tenant using the api-key within the x-api-key header (or using an admin key along with `x-tenant-id` param). User can specify time granularity (`monthly`, `daily` or `custom`) for retrieved results 

### Input

```
HTTP GET /api/v5/results/{report-name}/groups?[start-time]&[end-time]&[granularity]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start-time]`  | UTC time in W3C format                                                                          | YES      |
| `[end-time]`    | UTC time in W3C format                                                                          | YES      |
| `[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly`,  `daily` or `custom` | NO       | `daily`       |

#### Path Parameters

| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report-name}` | Name of the report that contains the results | YES      |


### Example Request 1: default daily granularity

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v5/results/Report_A/groups?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z 
```
or 
```
/api/v3/results/Report_A/groups?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily`
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
      "name": "PROJECT_A",
      "type": "PROJECT",
      "groups": [
        {
          "name": "GROUP-A",
          "type": "SERVICEGROUP",
          "results": [
            {
              "date": "2015-06-22",
              "availability": 68.14,
              "reliability": 68.14,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.68
            },
            {
              "date": "2015-06-23",
              "availability": 75.36,
              "reliability": 75.36,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.75
            }
          ]
        },
        {
          "name": "GROUP-B",
          "type": "SERVICEGROUP",
          "results": [
            {
              "date": "2015-06-22",
              "availability": 68.14,
              "reliability": 68.14,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.68
            },
            {
              "date": "2015-06-23",
              "availability": 75.36,
              "reliability": 75.36,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.75
            }
          ]
        }
      ]
    },
    {
      "name": "PROJECT_B",
      "type": "PROJECT",
      "groups": [
        {
          "name": "GROUP-X",
          "type": "SERVICEGROUP",
          "results": [
            {
              "date": "2015-06-22",
              "availability": 68.14,
              "reliability": 68.14,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.68
            },
            {
              "date": "2015-06-23",
              "availability": 75.36,
              "reliability": 75.36,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.75
            }
          ]
        },
        {
          "name": "GROUP-Y",
          "type": "SERVICEGROUP",
          "results": [
            {
              "date": "2015-06-22",
              "availability": 68.14,
              "reliability": 68.14,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.68
            },
            {
              "date": "2015-06-23",
              "availability": 75.36,
              "reliability": 75.36,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.75
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
/api/v5/results/Report_A/groups?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=monthly
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
      "name": "PROJECT_A",
      "type": "PROJECT",
      "groups": [
        {
          "name": "GROUP-A",
          "type": "SERVICEGROUP",
          "results": [
            {
              "date": "2015-06",
              "availability": 70,
              "reliability": 70,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.7
            }
          ]
        },
        {
          "name": "GROUP-B",
          "type": "SERVICEGROUP",
          "results": [
            {
              "date": "2015-06",
              "availability": 70,
              "reliability": 70,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.7
            }
          ]
        }
      ]
    },
    {
      "name": "PROJECT_B",
      "type": "PROJECT",
      "groups": [
        {
          "name": "GROUP-X",
          "type": "SERVICEGROUP",
          "results": [
            {
              "date": "2015-06",
              "availability": 70,
              "reliability": 70,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.7
            }
          ]
        },
        {
          "name": "GROUP-Y",
          "type": "SERVICEGROUP",
          "results": [
            {
              "date": "2015-06",
              "availability": 80,
              "reliability": 80,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.8
            }
          ]
        }
      ]
    }
  ]
}
```

### Example Request 3: Custom granularity
This request returns availability/reliability score numbers for the whole custom period defined between `start-time` and `end-time`. 
This means that for each item the user will receive one availability and reliability result concerning the whole period (instead of multiple daily or monthly results)

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v5/results/Report_A/groups?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=custom
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
      "name": "PROJECT_A",
      "type": "PROJECT",
      "groups": [
        {
          "name": "GROUP-A",
          "type": "SERVICEGROUP",
          "results": [
            {
              "date": "2015-06",
              "availability": 70,
              "reliability": 70,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.7
            }
          ]
        },
        {
          "name": "GROUP-B",
          "type": "SERVICEGROUP",
          "results": [
            {
              "date": "2015-06",
              "availability": 70,
              "reliability": 70,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.7
            }
          ]
        }
      ]
    },
    {
      "name": "PROJECT_B",
      "type": "PROJECT",
      "groups": [
        {
          "name": "GROUP-X",
          "type": "SERVICEGROUP",
          "results": [
            {
              "date": "2015-06",
              "availability": 90,
              "reliability": 90,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.9
            }
          ]
        },
        {
          "name": "GROUP-Y",
          "type": "SERVICEGROUP",
          "results": [
            {
              "date": "2015-06",
              "availability": 33,
              "reliability": 33,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.33
            }
          ]
        }
      ]
    }
  ]
}
```

## Get results for specific group {#4}

The following method can be used to obtain a tenant's Availability and Reliability result metrics for specific group. The api authenticates the tenant using the api-key within the x-api-key header (or using an admin key along with `x-tenant-id` param). User can specify time granularity (`monthly`, `daily` or `custom`) for retrieved results 

### Input

```
HTTP GET /api/v5/results/{report-name}/groups/{group-name}?[start-time]&[end-time]&[granularity]
```

#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start-time]`  | UTC time in W3C format                                                                          | YES      |
| `[end-time]`    | UTC time in W3C format                                                                          | YES      |
| `[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly`,  `daily` or `custom` | NO       | `daily`       |

#### Path Parameters

| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report-name}` | Name of the report that contains the results | YES      |
| `{group-name}`  | Name of the specific group to target | YES      |


### Example Request 1: default daily granularity

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v5/results/Report_A/groups/GROUP-A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z 
```
or 
```
/api/v3/results/Report_A/groups/GROUP-A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily`
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
      "name": "PROJECT_A",
      "type": "PROJECT",
      "groups": [
        {
          "name": "GROUP-A",
          "type": "SERVICEGROUP",
          "results": [
            {
              "date": "2015-06-22",
              "availability": 68.14,
              "reliability": 68.14,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.68
            },
            {
              "date": "2015-06-23",
              "availability": 75.36,
              "reliability": 75.36,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.75
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
/api/v5/results/Report_A/groups/GROUP-A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=monthly
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
      "name": "PROJECT_A",
      "type": "PROJECT",
      "groups": [
        {
          "name": "GROUP-A",
          "type": "SERVICEGROUP",
          "results": [
            {
              "date": "2015-06",
              "availability": 70,
              "reliability": 70,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.7
            }
          ]
        }
      ]
    }
  ]
}
```

### Example Request 3: Custom granularity
This request returns availability/reliability score numbers for the whole custom period defined between `start-time` and `end-time`. 
This means that for each item the user will receive one availability and reliability result concerning the whole period (instead of multiple daily or monthly results)

#### Request

##### Method
`HTTP GET`

##### Path

```
/api/v5/results/Report_A/groups/GROUP-A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=custom
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
      "name": "PROJECT_A",
      "type": "PROJECT",
      "groups": [
        {
          "name": "GROUP-A",
          "type": "SERVICEGROUP",
          "results": [
            {
              "date": "2015-06",
              "availability": 70,
              "reliability": 70,
              "unknown": 0,
              "downtime": 0,
              "uptime": 0.7
            }
          ]
        }
      ]
    }
  ]
}
```
