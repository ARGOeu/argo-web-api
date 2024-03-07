---
id: topology_endpoints
title: Topology Endpoints
sidebar_position: 2
---

## API calls for handling topology endpoint resources

| Name                                               | Description                                                            | Shortcut                     |
| -------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------- |
| POST: Create endpoint topology for specific date   | Creates a daily endpoint topology mapping endpoints to endpoint groups | <a href="#1">Description</a> |
| GET: List endpoint topology for specific date      | Lists endpoint topology for a specific date                            | <a href="#2">Description</a> |
| DELETE: delete endpoint topology for specific date | Deletes all endpoint items (topology) for a specific date              | <a href="#3">Description</a> |
| GET: List endpoint topology for specific report    | Lists endpoint topology for a specific report                          | <a href="#4">Description</a> |


## POST: Create endpoint topology for specific date {#1}

Creates a daily endpoint topology mapping endpoints to endpoint groups

### Input

```
POST /topology/endpoints?date=YYYY-MM-DD
```

#### Url Parameters

| Type   | Description            | Required | Default value |
| ------ | ---------------------- | -------- | ------------- |
| `date` | target a specific date | NO       | today's date  |

#### Headers

```
x-api-key: secret_key_value
Accept: application/json
```

### POST BODY

```json
[
    {
        "group": "SITE_A",
        "hostname": "host1.site-a.foo",
        "type": "SITES",
        "service": "a.service.foo",
        "tags": { "scope": "TENANT", "production": "1", "monitored": "1" }
    },
    {
        "group": "SITE_A",
        "hostname": "host2.site-b.foo",
        "type": "SITES",
        "service": "b.service.foo",
        "tags": { "scope": "TENANT", "production": "1", "monitored": "1" }
    },
    {
        "group": "SITE_B",
        "hostname": "host1.site-a.foo",
        "type": "SITES",
        "service": "c.service.foo",
        "tags": { "scope": "TENANT", "production": "1", "monitored": "1" },
        "notifications": {"contacts": ["email01@example.com"], "enabled": true}
    }
]
```

#### Response Code

```
Status: 201 OK Created
```

### Response body

```json
{
    "message": "Topology of 3 endpoints created for date: YYYY-MM-DD",
    "code": "201"
}
```

## 409 Conflict when trying to insert a topology that already exists

When trying to insert a topology for a specific date that already exists the api will answer with the following response:

### Response Code

```
Status: 409 Conflict
```

### Response body

```json
{
    "message": "topology already exists for date: YYYY-MM-DD, please either update it or delete it first!",
    "code": "409"
}
```

User can proceed with either updating the existing endpoint topology OR deleting before trying to create it anew


## GET: List endpoint topology per date {#2}

List endpoint topology for a specific date or the closest available topology to that date. If date is not provided list the latest available endpoint topology.

### Input

##### List All topology statistics

```
GET /topology/endpoints?date=YYYY-MM-DD
```

#### Url Parameters

| Type       | Description                   | Required | Default value |
| ---------- | ----------------------------- | -------- | ------------- |
| `date`     | target a specific date        | NO       | today's date  |
| `group`    | filter by group name          | NO       |               |
| `type`     | filter by group type          | NO       |               |
| `service`  | filter by service             | NO       |               |
| `hostname` | filter by hostname            | NO       |               |
| `tags`     | filter by tag key:value pairs | NO       |               |
| `mode`     | if stating `mode=combined` then if the tenant has data feeds from other tenants their endpoint topology items will be combined in the final results | NO       | empty |

_note_ : user can use wildcard \* in filters
_note_ : when using tag filters the query string must follow the pattern: `?tags=key1:value1,key2:value2`
_note_ : You can use `~` as a negative operator in the beginning of a filter value to exclude something: For example you can exclude endpoints with service of value `SERVICE_A` by issuing `?service:~SERVICE_A`

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
        "message": "Success",
        "code": "200"
    },
    "data": [
        {
            "date": "2019-12-12",
            "group": "SITE_A",
            "hostname": "host1.site-a.foo",
            "type": "SITES",
            "service": "a.service.foo",
            "tags": {
                "scope": "TENANT",
                "production": "1",
                "monitored": "1"
            }
        },
        {
            "date": "2019-12-12",
            "group": "SITE_A",
            "hostname": "host2.site-b.foo",
            "type": "SITES",
            "service": "b.service.foo",
            "tags": {
                "scope": "TENANT",
                "production": "1",
                "monitored": "1"
            }
        },
        {
            "date": "2019-12-12",
            "group": "SITE_B",
            "hostname": "host1.site-a.foo",
            "type": "SITES",
            "service": "c.service.foo",
            "tags": {
                "scope": "TENANT",
                "production": "1",
                "monitored": "1"
            },
            "notifications": {
                "contacts": ["email01@example.com"],
                "enabled": true
            }
        }
    ]
}
```

### Combined tenant example

If the tenant is configured to receive data from other tenants in its data feeds the combined mode can be used when retrieving topology endpoint items. In this mode the list of items is enriched with items from tenants specified in the data feeds. Items from those tenants receive an extra `tenant` field to signify their origin. To enable the combine mode use the optional url parameter `mode=combined`


#### Example Request
```
GET /topology/endpoints?date=2015-07-22?mode=combined
```

#### Response Code

```
Status: 200 OK
```

#### Response body

```json
{
    "status": {
        "message": "Success",
        "code": "200"
    },
    "data": [
        {
            "date": "2019-12-12",
            "group": "SITE_X",
            "hostname": "host01.site-y.foo",
            "type": "SITES",
            "service": "x.service.foo",
            "tags": {
                "scope": "TENANT",
                "production": "1",
                "monitored": "1"
            },
            "tenant": "TENANT_A"
        },
        {
            "date": "2019-12-12",
            "group": "SITE_Y",
            "hostname": "host01.site-y.foo",
            "type": "SITES",
            "service": "y.service.foo",
            "tags": {
                "scope": "TENANT",
                "production": "1",
                "monitored": "1"
            },
            "tenant": "TENANT_B"
        }
    ]
}
```



## [DELETE]: Delete endpoint topology for a specific date {#3}

This method can be used to delete all endpoint items contributing to the endpoint topology of a specific date

### Input

```
DELETE /topology/endpoints?date=YYYY-MM-DD
```

#### Request headers

```
x-api-key: shared_key_value
Content-Type: application/json
Accept: application/json
```

### Response

Headers: `Status: 200 OK`

#### Response body

Json Response

```json
{
    "message": "Topology of 3 endpoints deleted for date: 2019-12-12",
    "code": "200"
}
```


## GET: List endpoint topology for specific report {#4}

Lists endpoint topology items for specific report

### Input

```
GET /topology/endpoint/by_report/{report-name}?date=YYYY-MM-DD
```

#### Url Parameters

| Type          | Description              | Required | Default value |
| ------------- | ------------------------ | -------- | ------------- |
| `report-name` | target a specific report | YES      | none          |
| `date`        | target a specific date   | NO       | today's date  |
| `mode`     | if stating `mode=combined` then if the tenant has data feeds from other tenants their endpoint topology items will be combined in the final results | NO       | empty |

#### Headers

```
x-api-key: secret_key_value
Accept: application/json
```

#### Example Request

```
GET /topology/endpoints/by_report/Critical?date=2015-07-22
```

#### Response Code

```
Status: 200 OK
```

### Response body

```json
{
    "status": {
        "message": "Success",
        "code": "200"
    },
    "data": [
        {
            "date": "2019-12-12",
            "group": "SITE_A",
            "hostname": "host1.site-a.foo",
            "type": "SITES",
            "service": "a.service.foo",
            "tags": {
                "scope": "TENANT",
                "production": "1",
                "monitored": "1"
            }
        },
        {
            "date": "2019-12-12",
            "group": "SITE_A",
            "hostname": "host2.site-b.foo",
            "type": "SITES",
            "service": "b.service.foo",
            "tags": {
                "scope": "TENANT",
                "production": "1",
                "monitored": "1"
            }
        },
        {
            "date": "2019-12-12",
            "group": "SITE_B",
            "hostname": "host1.site-a.foo",
            "type": "SITES",
            "service": "c.service.foo",
            "tags": {
                "scope": "TENANT",
                "production": "1",
                "monitored": "1"
            },
            "notifications": {
                "contacts": ["email01@example.com"],
                "enabled": true
            }
        }
    ]
}
```
