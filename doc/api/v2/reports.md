---
title: 'API documentation | ARGO'
page_title: API - Reports
font_title: fa fa-cogs
description: API Calls for listing existing and creating new Reports
---

# API Calls

Name                              | Description                                                    | Shortcut
--------------------------------- | -------------------------------------------------------------- | ------------------
GET: List reports                 | This method can be used to retrieve a list of existing reports | [ Description](#1)
POST: Create a new report         | This method can be used to create a new report.                | [ Description](#2)
PUT: Update an existing report    | This method can be used to update an existing report.          | [ Description](#3)
DELETE: Delete an existing Report | This method can be used to delete an existing report.          | [ Description](#4)

<a id='1'></a>

# GET: List Recomputation Requests
This method can be used to retrieve a list of existing reports.

## Input

### URL

```
/reports
```

### Request headers

```
x-api-key: shared_key_value
Content-Type: application/json
Accept: application/json
```

## Response
Headers: `Status: 200 OK`

### Response body
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
            "info": {
                "name": "Report A",
                "description": "Description of the report, which is not required",
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
            "profiles": [
                {
                    "name": "profile1",
                    "type": "metric"
                },
                {
                    "name": "profile2",
                    "type": "ops"
                }
            ],
            "filter_tags": [
                {
                    "name": "name1",
                    "value": "value1"
                },
                {
                    "name": "name2",
                    "value": "value2"
                }
            ]
        }
    ]
}
```

Xml Response

```xml
<?xml version="1.0" encoding="UTF-8"?>
<root>
   <status>
      <message>Success</message>
      <code>200</code>
   </status>
   <data>
      <result>
         <id>eba61a9e-22e9-4521-9e47-ecaa4a494364</id>
         <info>
            <name>Report A</name>
            <description>Description of the report, which is not required</description>
            <created>2015-9-10 13:43:00</created>
            <updated>2015-10-11 13:43:00</updated>
         </info>
         <topology_schema>
            <group>
               <type>NGI</type>
               <group>
                  <type>SITE</type>
               </group>
            </group>
         </topology_schema>
         <profile name="profile1" type="metric" />
         <profile name="profile2" type="ops" />
         <tag name="name1" value="value1" />
         <tag name="name2" value="value2" />
      </result>
   </data>
</root>
```

<a id='2'></a>

# POST: Create a new report
This method can be used to create a new report

## Input

### URL


```
/reports
```

### Request headers

```
x-api-key: shared_key_value
Content-Type: application/json
Accept: application/json
```

### Request Body

```json
{
    "info": {
        "name": "Report A",
        "description": "Description of the report, which is not required",
    },
    "topology_schema": {
        "group": {
            "type": "NGI",
            "group": {
                "type": "SITE"
            }
        }
    },
    "profiles": [
        {
            "name": "profile1",
            "type": "metric"
        },
        {
            "name": "profile2",
            "type": "ops"
        }
    ],
    "filter_tags": [
        {
            "name": "name1",
            "value": "value1"
        },
        {
            "name": "name2",
            "value": "value2"
        }
    ]
}
```

## Response
Headers: `Status: 200 OK`

### Response Body

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

# PUT: Update an existing report
This method can be used to update an existing report

## Input

### URL
```
/reports/{uuid}
```

### Request headers

```
x-api-key: shared_key_value
Content-Type: application/json
Accept: application/json
```

### Request Body

```json
{
    "info": {
        "name": "Report A",
        "description": "Description of the report, which is not required",
    },
    "topology_schema": {
        "group": {
            "type": "NGI",
            "group": {
                "type": "SITE"
            }
        }
    },
    "profiles": [
        {
            "name": "profile1",
            "type": "metric"
        },
        {
            "name": "profile2",
            "type": "ops"
        }
    ],
    "filter_tags": [
        {
            "name": "name1",
            "value": "value1"
        },
        {
            "name": "name2",
            "value": "value2"
        }
    ]
}
```

## Response
Headers: `Status: 200 OK`

### Response Body

```json
{
    "status": {
        "message": "Report was successfully updated",
        "code": "200"
    }
}
```

<a id='4'></a>

# DELETE: Delete an existing report
This method can be used to update an existing report

## Input

### URL

```
/reports/{uuid}
```

### Request headers

```
x-api-key: shared_key_value
Content-Type: application/json
Accept: application/json
```

## Response
Headers: `Status: 200 OK`

### Response Body

```json
{
    "status": {
        "message": "Report was successfully deleted",
        "code": "200"
    }
}
```
