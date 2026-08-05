
V5 status · MD
---
id: v5_status
title: Status (v5)
sidebar_position: 2
---
 
## API Calls
 
_Note_: These are v5 api calls implementations found under the path `/api/v5/status`
 
| Name                                                                          | Description                                                                                                                                                                                                                              | Shortcut          |
| ----------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------- |
| Get status for groups | This method retrieves the status timelines of all groups | [Description](#1) |
| Get status for specific group | This method retrieves the status timeline for a specific group | [Description](#2) |
| Get status for service-types of a specific group | This method retrieves the status timelines of all service-types of a specific group | [Description](#3) |
| Get status for a specific service-type of a specific group | This method retrieves the status timeline for a specific service-type of a specific group | [Description](#4) |
| Get status for endpoints or specific endpoint | This method retrieves the status timelines of all endpoints or a specific endpoint | [Description](#5A) |
| Get status for endpoints of a specific group and a specific service-type | This method retrieves the status timelines of all endpoints of a specific group and a specific service-type | [Description](#5) |
| Get status for specific endpoint of a specific group and a specific service-type | This method retrieves the status timeline for a specific endpoint of a specific group and a specific service-type | [Description](#6) |
| Get status for metrics of a specific endpoint | This method retrieves the status timelines of all metrics of a specific endpoint | [Description](#7) |
| Get status for a specific metric of a specific endpoint | This method retrieves the status timeline for a specific metric of a specific endpoint | [Description](#8) |
 
## Status timelines
 
Each item in a status response carries a `statuses` list, which is the timeline of
status changes of that item within the requested period. Each entry has a
`timestamp` and a `value`.
 
The first entry of a timeline may be a _seed_ entry, carrying empty strings for
both `timestamp` and `value`. It represents the state of the item as it was
carried over from before `start-time`, when no status change had yet been
recorded inside the requested window.
 
## Get status for groups {#1}
 
The following method can be used to obtain a tenant's status timelines for all groups. The api authenticates the tenant using the api-key within the x-api-key header (or using an admin key along with `x-tenant-id` param). User can specify the period of the retrieved timelines using `start-time` and `end-time`
 
### Input
 
```
HTTP GET /api/v5/status/{report-name}/groups?[start-time]&[end-time]
```
 
#### Query Parameters
 
| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start-time]`  | UTC time in W3C format                                                                          | YES      |
| `[end-time]`    | UTC time in W3C format                                                                          | YES      |
 
#### Path Parameters
 
| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report-name}` | Name of the report that contains the results | YES      |
 
### Example Request
 
#### Request
 
##### Method
`HTTP GET`
 
##### Path
 
```
/api/v5/status/CORE/groups?start-time=2026-07-23T00:00:00Z&end-time=2026-07-23T23:59:59Z
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
    "groups": [
        {
            "name": "CLOUD-A",
            "type": "SERVICEGROUPS",
            "statuses": [
                {
                    "timestamp": "2026-07-23T09:44:14Z",
                    "value": "WARNING"
                },
                {
                    "timestamp": "2026-07-23T09:45:48Z",
                    "value": "WARNING"
                }
            ]
        },
        {
            "name": "CLOUD-B",
            "type": "SERVICEGROUPS",
            "statuses": [
                {
                    "timestamp": "2026-07-23T00:00:00Z",
                    "value": "OK"
                },
                {
                    "timestamp": "2026-07-23T11:12:03Z",
                    "value": "CRITICAL"
                }
            ]
        }
    ]
}
```
 
## Get status for specific group {#2}
 
The following method can be used to obtain a tenant's status timeline for a specific group. The api authenticates the tenant using the api-key within the x-api-key header (or using an admin key along with `x-tenant-id` param). User can specify the period of the retrieved timelines using `start-time` and `end-time`
 
### Input
 
```
HTTP GET /api/v5/status/{report-name}/groups/{group-name}?[start-time]&[end-time]
```
 
#### Query Parameters
 
| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start-time]`  | UTC time in W3C format                                                                          | YES      |
| `[end-time]`    | UTC time in W3C format                                                                          | YES      |
 
#### Path Parameters
 
| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report-name}` | Name of the report that contains the results | YES      |
| `{group-name}`  | Name of the specific group to target | YES      |
 
### Example Request
 
#### Request
 
##### Method
`HTTP GET`
 
##### Path
 
```
/api/v5/status/CORE/groups/CLOUD-A?start-time=2026-07-23T00:00:00Z&end-time=2026-07-23T23:59:59Z
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
    "groups": [
        {
            "name": "CLOUD-A",
            "type": "SERVICEGROUPS",
            "statuses": [
                {
                    "timestamp": "2026-07-23T09:44:14Z",
                    "value": "WARNING"
                },
                {
                    "timestamp": "2026-07-23T09:45:48Z",
                    "value": "WARNING"
                }
            ]
        }
    ]
}
```
 
## Get status for service-types of a specific group {#3}
 
The following method can be used to obtain a tenant's status timelines for all service-types of a specific group. The api authenticates the tenant using the api-key within the x-api-key header (or using an admin key along with `x-tenant-id` param). User can specify the period of the retrieved timelines using `start-time` and `end-time`
 
### Input
 
```
HTTP GET /api/v5/status/{report-name}/groups/{group-name}/service-types?[start-time]&[end-time]
```
 
#### Query Parameters
 
| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start-time]`  | UTC time in W3C format                                                                          | YES      |
| `[end-time]`    | UTC time in W3C format                                                                          | YES      |
 
#### Path Parameters
 
| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report-name}` | Name of the report that contains the results | YES      |
| `{group-name}` | Name of the group that contains the results | YES      |
 
### Example Request
 
#### Request
 
##### Method
`HTTP GET`
 
##### Path
 
```
/api/v5/status/CORE/groups/CLOUD-A/service-types?start-time=2026-07-23T00:00:00Z&end-time=2026-07-23T23:59:59Z
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
    "groups": [
        {
            "name": "CLOUD-A",
            "type": "SERVICEGROUPS",
            "service-types": [
                {
                    "name": "webportal",
                    "type": "service",
                    "statuses": [
                        {
                            "timestamp": "2026-07-23T09:44:14Z",
                            "value": "WARNING"
                        },
                        {
                            "timestamp": "2026-07-23T09:45:26Z",
                            "value": "WARNING"
                        }
                    ]
                }
            ]
        }
    ]
}
```
 
## Get status for a specific service-type of a specific group {#4}
 
The following method can be used to obtain a tenant's status timeline for a specific service-type of a specific group. The api authenticates the tenant using the api-key within the x-api-key header (or using an admin key along with `x-tenant-id` param). User can specify the period of the retrieved timelines using `start-time` and `end-time`
 
### Input
 
```
HTTP GET /api/v5/status/{report-name}/groups/{group-name}/service-types/{service-type-name}?[start-time]&[end-time]
```
 
#### Query Parameters
 
| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start-time]`  | UTC time in W3C format                                                                          | YES      |
| `[end-time]`    | UTC time in W3C format                                                                          | YES      |
 
#### Path Parameters
 
| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report-name}` | Name of the report that contains the results | YES      |
| `{group-name}`  | Name of the specific group to target | YES      |
| `{service-type-name}`  | Name of the specific service-type to target | YES      |
 
### Example Request
 
#### Request
 
##### Method
`HTTP GET`
 
##### Path
 
```
/api/v5/status/CORE/groups/CLOUD-A/service-types/webportal?start-time=2026-07-23T00:00:00Z&end-time=2026-07-23T23:59:59Z
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
    "groups": [
        {
            "name": "CLOUD-A",
            "type": "SERVICEGROUPS",
            "service-types": [
                {
                    "name": "webportal",
                    "type": "service",
                    "statuses": [
                        {
                            "timestamp": "2026-07-23T09:44:14Z",
                            "value": "WARNING"
                        },
                        {
                            "timestamp": "2026-07-23T09:45:26Z",
                            "value": "WARNING"
                        }
                    ]
                }
            ]
        }
    ]
}
```

 
## Get status for endpoints or a specific endpoint {#5A}
 
The following methods can be used to obtain a tenant's status timelines for all endpoints or a specific endpoint. The api authenticates the tenant using the api-key within the x-api-key header (or using an admin key along with `x-tenant-id` param). User can specify the period of the retrieved timelines using `start-time` and `end-time`
 
### Input
 
```
HTTP GET /api/v5/status/{report-name}/endpoints?[start-time]&[end-time]
```

or 

```
HTTP GET /api/v5/status/{report-name}/endpoints/{endpoint-name}?[start-time]&[end-time]
```
 
#### Query Parameters
 
| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start-time]`  | UTC time in W3C format                                                                          | YES      |
| `[end-time]`    | UTC time in W3C format                                                                          | YES      |
 
#### Path Parameters
 
| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report-name}` | Name of the report that contains the results | YES      |
| `{endpoint-name}`  | Name of the specific endpoint to target | NO      |
 
### Example Request
 
#### Request
 
##### Method
`HTTP GET`
 
##### Path
 
```
/api/v5/status/CORE/endpoints?start-time=2026-07-23T00:00:00Z&end-time=2026-07-23T23:59:59Z
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
    "groups": [
        {
            "name": "CLOUD-A",
            "type": "SERVICEGROUPS",
            "service-types": [
                {
                    "name": "webportal",
                    "type": "service",
                    "endpoints": [
                        {
                            "name": "host1.clouda.foo_ID1",
                            "info": {
                                "ID": "ID1",
                                "URL": "https://host1.clouda.foo"
                            },
                            "statuses": [
                                {
                                    "timestamp": "2026-07-23T09:44:14Z",
                                    "value": "WARNING"
                                },
                                {
                                    "timestamp": "2026-07-23T09:45:09Z",
                                    "value": "WARNING"
                                }
                            ]
                        }
                    ]
                }
            ]
        },
        {
            "name": "CLOUD-B",
            "type": "SERVICEGROUPS",
            "service-types": [
                {
                    "name": "webportal",
                    "type": "service",
                    "endpoints": [
                        {
                            "name": "host1.cloudb.foo_ID2",
                            "info": {
                                "ID": "ID1",
                                "URL": "https://host1.cloudb.foo"
                            },
                            "statuses": [
                                {
                                    "timestamp": "2026-07-23T09:44:14Z",
                                    "value": "OK"
                                },
                                {
                                    "timestamp": "2026-07-23T09:45:09Z",
                                    "value": "OK"
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

### Example Request
 
#### Request
 
##### Method
`HTTP GET`
 
##### Path
 
```
/api/v5/status/CORE/endpoints/host1.clouda.foo_ID1?start-time=2026-07-23T00:00:00Z&end-time=2026-07-23T23:59:59Z
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
    "groups": [
        {
            "name": "CLOUD-A",
            "type": "SERVICEGROUPS",
            "service-types": [
                {
                    "name": "webportal",
                    "type": "service",
                    "endpoints": [
                        {
                            "name": "host1.clouda.foo_ID1",
                            "info": {
                                "ID": "ID1",
                                "URL": "https://host1.clouda.foo"
                            },
                            "statuses": [
                                {
                                    "timestamp": "2026-07-23T09:44:14Z",
                                    "value": "WARNING"
                                },
                                {
                                    "timestamp": "2026-07-23T09:45:09Z",
                                    "value": "WARNING"
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
 
## Get status for endpoints of a specific group and a specific service-type {#5}
 
The following method can be used to obtain a tenant's status timelines for all endpoints of a specific group and a specific service-type. The api authenticates the tenant using the api-key within the x-api-key header (or using an admin key along with `x-tenant-id` param). User can specify the period of the retrieved timelines using `start-time` and `end-time`
 
### Input
 
```
HTTP GET /api/v5/status/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints?[start-time]&[end-time]
```
 
#### Query Parameters
 
| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start-time]`  | UTC time in W3C format                                                                          | YES      |
| `[end-time]`    | UTC time in W3C format                                                                          | YES      |
 
#### Path Parameters
 
| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report-name}` | Name of the report that contains the results | YES      |
| `{group-name}` | Name of the group that contains the results | YES      |
| `{service-type-name}`  | Name of the specific service-type to target | YES      |
 
### Example Request
 
#### Request
 
##### Method
`HTTP GET`
 
##### Path
 
```
/api/v5/status/CORE/groups/CLOUD-A/service-types/webportal/endpoints?start-time=2026-07-23T00:00:00Z&end-time=2026-07-23T23:59:59Z
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
    "groups": [
        {
            "name": "CLOUD-A",
            "type": "SERVICEGROUPS",
            "service-types": [
                {
                    "name": "webportal",
                    "type": "service",
                    "endpoints": [
                        {
                            "name": "host1.clouda.foo_ID1",
                            "info": {
                                "ID": "ID1",
                                "URL": "https://host1.clouda.foo"
                            },
                            "statuses": [
                                {
                                    "timestamp": "2026-07-23T09:44:14Z",
                                    "value": "WARNING"
                                },
                                {
                                    "timestamp": "2026-07-23T09:45:09Z",
                                    "value": "WARNING"
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
 
## Get status for specific endpoint of a specific group and a specific service-type {#6}
 
The following method can be used to obtain a tenant's status timeline for a specific endpoint of a specific group and a specific service-type. The api authenticates the tenant using the api-key within the x-api-key header (or using an admin key along with `x-tenant-id` param). User can specify the period of the retrieved timelines using `start-time` and `end-time`
 
### Input
 
```
HTTP GET /api/v5/status/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints/{endpoint-name}?[start-time]&[end-time]
```
 
#### Query Parameters
 
| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start-time]`  | UTC time in W3C format                                                                          | YES      |
| `[end-time]`    | UTC time in W3C format                                                                          | YES      |
 
#### Path Parameters
 
| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report-name}` | Name of the report that contains the results | YES      |
| `{group-name}`  | Name of the specific group to target | YES      |
| `{service-type-name}`  | Name of the specific service-type to target | YES      |
| `{endpoint-name}`  | Name of the specific endpoint to target | YES      |
 
### Example Request
 
#### Request
 
##### Method
`HTTP GET`
 
##### Path
 
```
/api/v5/status/CORE/groups/CLOUD-A/service-types/webportal/endpoints/host1.clouda.foo_ID1?start-time=2026-07-23T00:00:00Z&end-time=2026-07-23T23:59:59Z
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
    "groups": [
        {
            "name": "CLOUD-A",
            "type": "SERVICEGROUPS",
            "service-types": [
                {
                    "name": "webportal",
                    "type": "service",
                    "endpoints": [
                        {
                            "name": "host1.clouda.foo_ID1",
                            "info": {
                                "ID": "ID1",
                                "URL": "https://host1.clouda.foo"
                            },
                            "statuses": [
                                {
                                    "timestamp": "2026-07-23T09:44:14Z",
                                    "value": "WARNING"
                                },
                                {
                                    "timestamp": "2026-07-23T09:45:09Z",
                                    "value": "WARNING"
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
 
## Get status for metrics of a specific endpoint {#7}
 
The following method can be used to obtain a tenant's status timelines for all metrics of a specific endpoint. The api authenticates the tenant using the api-key within the x-api-key header (or using an admin key along with `x-tenant-id` param). User can specify the period of the retrieved timelines using `start-time` and `end-time`
 
### Input
 
```
HTTP GET /api/v5/status/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints/{endpoint-name}/metrics?[start-time]&[end-time]
```
 
#### Query Parameters
 
| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start-time]`  | UTC time in W3C format                                                                          | YES      |
| `[end-time]`    | UTC time in W3C format                                                                          | YES      |
 
#### Path Parameters
 
| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report-name}` | Name of the report that contains the results | YES      |
| `{group-name}`  | Name of the specific group to target | YES      |
| `{service-type-name}`  | Name of the specific service-type to target | YES      |
| `{endpoint-name}`  | Name of the specific endpoint to target | YES      |
 
### Example Request
 
#### Request
 
##### Method
`HTTP GET`
 
##### Path
 
```
/api/v5/status/CORE/groups/CLOUD-A/service-types/webportal/endpoints/host1.clouda.foo_ID1/metrics?start-time=2026-07-23T00:00:00Z&end-time=2026-07-23T23:59:59Z
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
    "groups": [
        {
            "name": "CLOUD-A",
            "type": "SERVICEGROUPS",
            "service-types": [
                {
                    "name": "webportal",
                    "type": "service",
                    "endpoints": [
                        {
                            "name": "host1.clouda.foo_ID1",
                            "info": {
                                "ID": "ID1",
                                "URL": "https://host1.clouda.foo"
                            },
                            "metrics": [
                                {
                                    "name": "generic.http.connect",
                                    "statuses": [
                                        {
                                            "timestamp": "",
                                            "value": ""
                                        },
                                        {
                                            "timestamp": "2026-07-23T09:44:14Z",
                                            "value": "WARNING"
                                        }
                                    ]
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
 
## Get status for a specific metric of a specific endpoint {#8}
 
The following method can be used to obtain a tenant's status timeline for a specific metric of a specific endpoint. The api authenticates the tenant using the api-key within the x-api-key header (or using an admin key along with `x-tenant-id` param). User can specify the period of the retrieved timelines using `start-time` and `end-time`
 
### Input
 
```
HTTP GET /api/v5/status/{report-name}/groups/{group-name}/service-types/{service-type-name}/endpoints/{endpoint-name}/metrics/{metric-name}?[start-time]&[end-time]
```
 
#### Query Parameters
 
| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start-time]`  | UTC time in W3C format                                                                          | YES      |
| `[end-time]`    | UTC time in W3C format                                                                          | YES      |
 
#### Path Parameters
 
| Name            | Description                                                                                           | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report-name}` | Name of the report that contains the results | YES      |
| `{group-name}`  | Name of the specific group to target | YES      |
| `{service-type-name}`  | Name of the specific service-type to target | YES      |
| `{endpoint-name}`  | Name of the specific endpoint to target | YES      |
| `{metric-name}`  | Name of the specific metric to target | YES      |
 
### Example Request
 
#### Request
 
##### Method
`HTTP GET`
 
##### Path
 
```
/api/v5/status/CORE/groups/CLOUD-A/service-types/webportal/endpoints/host1.clouda.foo_ID1/metrics/generic.http.connect?start-time=2026-07-23T00:00:00Z&end-time=2026-07-23T23:59:59Z
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
    "groups": [
        {
            "name": "CLOUD-A",
            "type": "SERVICEGROUPS",
            "service-types": [
                {
                    "name": "webportal",
                    "type": "service",
                    "endpoints": [
                        {
                            "name": "host1.clouda.foo_ID1",
                            "info": {
                                "ID": "ID1",
                                "URL": "https://host1.clouda.foo"
                            },
                            "metrics": [
                                {
                                    "name": "generic.http.connect",
                                    "statuses": [
                                        {
                                            "timestamp": "",
                                            "value": ""
                                        },
                                        {
                                            "timestamp": "2026-07-23T09:44:14Z",
                                            "value": "WARNING"
                                        }
                                    ]
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