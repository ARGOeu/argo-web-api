---
id: topology_stats
title: Topology Statistics
---

## API calls for retrieving topology statistics per report

| Name                          | Description                                           | Shortcut                     |
| ----------------------------- | ----------------------------------------------------- | ---------------------------- |
| GET: List topology statistics | List number of groups, endpoint groups and services . | <a href="#1">Description</a> |

<a id="1"></a>

## [GET]: List topology statistics

This method may be used to retrieve topology statistics for a specific report. Topology statistics include number of groups, endpoint groups and services included in the report

### Input

##### List All topology statistics

```
/topology/stats/{report}/?[date]
```

#### Path Parameters

| Type     | Description             | Required | Default value |
| -------- | ----------------------- | -------- | ------------- |
| `report` | name of the report used | YES      |               |

#### Url Parameters

| Type   | Description            | Required | Default value |
| ------ | ---------------------- | -------- | ------------- |
| `date` | target a specific data | NO       | today's date  |

#### Headers

```
x-api-key: secret_key_value
Accept: application/json
```

#### Response Code

```
Status: 200 OK
```

### Response body

```json
{
    "status": {
        "message": "application/json",
        "code": "200"
    },
    "data": {
        "group_count": 1,
        "group_type": "type of top-level groups in report",
        "group_list": ["list of top level groups"],
        "endpoint_group_count": 1,
        "endpoint_group_type": "type of endpoint groups in report",
        "endpoint_group_list": ["list of endpoint groups"],
        "service_count": 1,
        "service_list": ["list of available services"]
    }
}
```

###### Example Request:

URL:

```
latest/Report_B/?date=2015-05-01
```

Headers:

```
x-api-key: secret_key_value
Accept: application/json
```

###### Example Response:

Code:

```
Status: 200 OK
```

Response body:

```json
{
    "status": {
        "message": "application/json",
        "code": "200"
    },
    "data": {
        "group_count": 2,
        "group_type": "PROJECTS",
        "group_list": ["PROJECT_A", "PROJECT_B"],
        "endpoint_group_count": 3,
        "endpoint_group_type": "SERVICEGROUPS",
        "endpoint_group_list": ["SGROUP_A", "SGROUP_B", "SGROUP_C"],
        "service_count": 8,
        "service_list": [
            "service_type_1",
            "service_type_2",
            "service_type_3",
            "service_type_4",
            "service_type_5",
            "service_type_6",
            "service_type_7",
            "service_type_8"
        ]
    }
}
```
