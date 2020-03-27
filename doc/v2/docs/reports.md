---
title: "API documentation | ARGO"
page_title: API - Reports
font_title: fa fa-cogs
description: API Calls for listing existing and creating new Reports
---

# API Calls

| Name                               | Description                                                    | Shortcut           |
| ---------------------------------- | -------------------------------------------------------------- | ------------------ |
| GET: List reports or single report | This method can be used to retrieve a list of existing reports | [ Description](#1) |
| POST: Create a new report          | This method can be used to create a new report.                | [ Description](#2) |
| PUT: Update an existing report     | This method can be used to update an existing report.          | [ Description](#3) |
| DELETE: Delete an existing Report  | This method can be used to delete an existing report.          | [ Description](#4) |

<a id='1'></a>

## [GET]: List Reports

This method can be used to retrieve a list of existing reports or a single report according to its ID.

### Input

#### URL

```
/reports
or
/reports/{id}
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

### Response

Headers: `Status: 200 OK`

#### Response body

Json Response

```json
{
    "status": {
        "message": "Success",
        "code": "200"
    },
    "data": [
        {
            "id": "eba61a9e-22e9-4521-9e47-ecaa4a494364",
            "tenant": "TenantA",
            "disabled": false,
            "info": {
                "name": "Report_A",
                "description": "report aaaaa",
                "created": "2015-9-10 13:43:00",
                "updated": "2015-10-11 13:43:00"
            },
            "topology_schema": {
                "group": {
                    "type": "NGI",
                    "group": {
                        "type": "SITE"
                    }
                }
            },
            "thresholds": {
                "availability": 80.0,
                "reliability": 85.0,
                "uptime": 80.0,
                "unknown": 10.0,
                "downtime": 10.0
            },
            "profiles": [
                {
                    "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
                    "name": "profile1",
                    "type": "metric"
                },
                {
                    "id": "6ac7d684-1f8e-4a02-a502-720e8f11e523",
                    "name": "profile2",
                    "type": "operations"
                },
                {
                    "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50q",
                    "name": "profile3",
                    "type": "aggregation"
                }
            ],
            "filter_tags": [
                {
                    "name": "name1",
                    "value": "value1",
                    "context": ""
                },
                {
                    "name": "name2",
                    "value": "value2",
                    "context": ""
                }
            ]
        }
    ]
}
```

<a id='2'></a>

## [POST]: Create a new report

This method can be used to create a new report

### Input

#### URL

```
POST /reports
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### Request Body

```json
{
    "info": {
        "name": "Report_A",
        "description": "report aaaaa"
    },
    "topology_schema": {
        "group": {
            "type": "ngi",
            "group": {
                "type": "site"
            }
        }
    },
    "thresholds": {
        "availability": 80.0,
        "reliability": 85.0,
        "uptime": 80.0,
        "unknown": 10.0,
        "downtime": 10.0
    },
    "profiles": [
        {
            "id": "422985a7-6386-4964-bc99-5ebd5d7b0aef",
            "type": "metric"
        },
        {
            "id": "1aa74849-2310-4bbc-b63a-8995ac7888ea",
            "type": "aggregation"
        },
        {
            "id": "1eafbdd1-1bbc-4861-b849-65394840762",
            "type": "operations"
        }
    ],
    "filter_tags": [
        {
            "name": "production",
            "value": "N"
        },
        {
            "name": "monitored",
            "value": "Y"
        }
    ]
}
```

### Response

Headers: `Status: 201 Created`

#### Response Body

```json
{
    "status": {
        "message": "Successfully Created Report",
        "code": "201"
    },
    "data": {
        "id": "eba61a9e-22e9-4521-9e47-ecaa4a494360",
        "links": {
            "self": "https://myapi.test.com/api/v2/reports/eba61a9e-22e9-4521-9e47-ecaa4a494360"
        }
    }
}
```

<a id='3'></a>

## [PUT]: Update an existing report

This method can be used to update an existing report. This will replace all the fields in the record so all the old fields that need to be kept must be included in the json of the update request body

### Input

#### URL

```
/reports/{id}
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### Request Body

```json
{
    "disabled": true,
    "weight": "hepspec",
    "info": {
        "name": "newname",
        "description": "newdescription"
    },
    "topology_schema": {
        "group": {
            "type": "ngi",
            "group": {
                "type": "fight"
            }
        }
    },
    "thresholds": {
        "availability": 90.0,
        "reliability": 95.0,
        "uptime": 90.0,
        "unknown": 15.0,
        "downtime": 15.0
    },
    "profiles": [
        {
            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
            "type": "metric",
            "name": "profile1"
        },
        {
            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e523",
            "type": "operations",
            "name": "profile2"
        },
        {
            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50bq",
            "type": "aggregation",
            "name": "profile3"
        }
    ],
    "filter_tags": [
        {
            "name": "production",
            "value": "N"
        },
        {
            "name": "monitored",
            "value": "N"
        }
    ]
}
```

### Response

Headers: `Status: 200 OK`

#### Response Body

```json
{
    "status": {
        "message": "Report was successfully updated",
        "code": "200"
    }
}
```

<a id='4'></a>

## [DELETE]: Delete an existing report

This method can be used to update an existing report

### Input

#### URL

```
DELETE /reports/{id}
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

### Response

Headers: `Status: 200 OK`

#### Response Body

```json
{
    "status": {
        "message": "Report was successfully deleted",
        "code": "200"
    }
}
```

<a id='5'></a>

## Notes on Report Filter tags and topology

As we seen before a report can host a list of filter tags using the following list under the filed `filter_tags`:

```json
{
    "filter_tags": [
        {
            "context": "a context description to define where the filter applies",
            "name": "what to be filter",
            "value": "filter pattern described here"
        }
    ]
}
```

There are special argo contextes that are automatically picked up to filter group and endpoint topology. These contexts are described below:

-   _context:_ `argo.group.filter.fields` - Used to apply filter on basic fields of group topology. Under this context the `name` targets the group field name and the `value` holds the actual field pattern
-   _context:_ `argo.group.filter.tags` - Used to apply filter on tags of group topology. Under this context the `name` targets the group tag name and the `value` holds the actual field pattern
-   _context:_ `argo.endpoint.filter.fields` - Used to apply filter on basic fields of endpoint topology. Under this context the `name` targets the endpoint field name and the `value` holds the actual field pattern
-   _context:_ `argo.endpoint.filter.tags` - Used to apply filter on tags of endpoint topology. Under this context the `name` targets the endpoint tag name and the `value` holds the actual field pattern
