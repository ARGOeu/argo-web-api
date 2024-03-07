---
id: topology_groups
title: Topology Groups
sidebar_position: 3
---

## API calls for handling topology group resources

| Name                                            | Description                                                         | Shortcut                     |
| ----------------------------------------------- | ------------------------------------------------------------------- | ---------------------------- |
| POST: Create group topology for specific date   | Creates a daily group topology mapping endpoints to endpoint groups | <a href="#1">Description</a> |
| GET: List group topology for specific date      | Lists group topology for a specific date                            | <a href="#2">Description</a> |
| DELETE: Delete group topology for specific date | Delete group topology items for specific date                       | <a href="#3">Description</a> |
| GET: List group topology for specific report    | Lists group topology for a specific report                          | <a href="#4">Description</a> |


## POST: Create group topology for specific date {#1}

Creates a daily group topology mapping top-level groups to subgroups

### Input

```
POST /topology/groups?date=YYYY-MM-DD
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
        "group": "NGIA",
        "type": "NGIS",
        "subgroup": "SITEA",
        "tags": {
            "scope": "FEDERATION",
            "infrastructure": "Production",
            "certification": "Certified"
        }
    },
    {
        "group": "NGIA",
        "type": "NGIS",
        "subgroup": "SITEB",
        "tags": {
            "scope": "FEDERATION",
            "infrastructure": "Production",
            "certification": "Certified"
        },
        "notifications": {
                "contacts": ["email01@example.com"],
                "enabled": true
        }
    },
    {
        "group": "PROJECTZ",
        "type": "PROJECT",
        "subgroup": "SITEZ",
        "tags": {
            "scope": "FEDERATION",
            "infrastructure": "Production",
            "certification": "Certified"
        }
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
    "message": "Topology of 3 groups created for date: YYYY-MM-DD",
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

User can proceed with either updating the existing topology OR deleting before trying to create it anew


## GET: List group topology for specific date {#2}

Lists group topology items for specific date

### Input

```
GET /topology/groups?date=YYYY-MM-DD
```

#### Url Parameters

| Type       | Description                   | Required | Default value |
| ---------- | ----------------------------- | -------- | ------------- |
| `date`     | target a specific date        | NO       | today's date  |
| `group`    | filter by group name          | NO       |               |
| `type`     | filter by group type          | NO       |               |
| `subgroup` | filter by subgroup            | NO       |               |
| `tags`     | filter by tag key:value pairs | NO       |               |
| `mode`     | if stating `mode=combined` then if the tenant has data feeds from other tenants their group topology items will be combined in the final results | NO       | empty |

_note_ : user can use wildcard \* in filters
_note_ : when using tag filters the query string must follow the pattern: `?tags=key1:value1,key2:value2`
_note_ : You can use `~` as a negative operator in the beginning of a filter value to exclude something: For example you can exclude endpoints with subgroup of value `GROUP_A` by issuing `?subgroup:~SERVICE_A`

#### Headers

```
x-api-key: secret_key_value
Accept: application/json
```

#### Example Request

```
GET /topology/groups?date=2015-07-22
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
            "date": "2015-07-22",
            "group": "NGIA",
            "type": "NGIS",
            "subgroup": "SITEA",
            "tags": {
                "certification": "Certified",
                "infrastructure": "Production"
            }
        },
        {
            "date": "2015-07-22",
            "group": "NGIA",
            "type": "NGIS",
            "subgroup": "SITEB",
            "tags": {
                "certification": "Certified",
                "infrastructure": "Production"
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

If the tenant is configured to receive data from other tenants in its data feeds the combined mode can be used when retrieving topology group items. In this mode the list of items is enriched with items from tenants specified in the data feeds. Items from those tenants receive an extra `tenant` field to signify their origin. To enable the combine mode use the optional url parameter `mode=combined`


#### Example Request
```
GET /topology/groups?date=2015-07-22?mode=combined
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
            "date": "2015-07-22",
            "group": "TENANT1-NGIA",
            "type": "NGIS",
            "subgroup": "SITEX",
            "tags": {
                "certification": "Certified",
                "infrastructure": "Production"
            },
            "tenant":"TENANT1"
        },
        {
            "date": "2015-07-22",
            "group": "TENANT2-NGIA",
            "type": "NGIS",
            "subgroup": "SITEZ",
            "tags": {
                "certification": "Certified",
                "infrastructure": "Production"
            },
            "tenant":"TENANT2"
        }
    ]
}
```


## [DELETE]: Delete group topology for a specific date {#3}

This method can be used to delete all group items contributing to the group topology of a specific date

### Input

```
DELETE /topology/groups?date=YYYY-MM-DD
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
    "message": "Topology of 3 groups deleted for date: 2019-12-12",
    "code": "200"
}
```


## GET: List group topology for specific report {#4}

Lists group topology items for specific report

### Input

```
GET /topology/groups/by_report/{report-name}?date=YYYY-MM-DD
```

#### Url Parameters

| Type          | Description              | Required | Default value |
| ------------- | ------------------------ | -------- | ------------- |
| `report-name` | target a specific report | YES      | none          |
| `date`        | target a specific date   | NO       | today's date  |
| `mode`     | if stating `mode=combined` then if the tenant has data feeds from other tenants their group topology items will be combined in the final results | NO       | empty |

#### Headers

```
x-api-key: secret_key_value
Accept: application/json
```

#### Example Request

```
GET /topology/groups/by_report/Critical?date=2015-07-22
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
            "date": "2015-07-22",
            "group": "NGIA",
            "type": "NGIS",
            "subgroup": "SITEA",
            "tags": {
                "certification": "Certified",
                "infrastructure": "Production"
            }
        },
        {
            "date": "2015-07-22",
            "group": "NGIA",
            "type": "NGIS",
            "subgroup": "SITEB",
            "tags": {
                "certification": "Certified",
                "infrastructure": "Production"
            },
            "notifications": {
                "contacts": ["email01@example.com"],
                "enabled": true
            }
        },
        {
            "date": "2015-07-22",
            "group": "NGIX",
            "type": "NGIS",
            "subgroup": "SITEX",
            "tags": {
                "certification": "Certified",
                "infrastructure": "Production"
            }
        }
    ]
}
```
