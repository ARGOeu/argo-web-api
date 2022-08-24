---
id: weights
title: Weights
sidebar_position: 6
---

## API Calls

| Name                                  | Description                                                                     | Shortcut           |
| ------------------------------------- | ------------------------------------------------------------------------------- | ------------------ |
| GET: List Weights resources Request   | This method can be used to retrieve a list of current weight resources.         | [ Description](#1) |
| GET: List a specific Weights resource | This method can be used to retrieve a specific weight resource based on its id. | [ Description](#2) |
| POST: Create a new weight resource    | This method can be used to create a new weight resource                         | [ Description](#3) |
| PUT: Update a weight resource         | This method can be used to update information on an existing weight resource    | [ Description](#4) |
| DELETE: Delete a weight resource      | This method can be used to delete an existing weight resource                   | [ Description](#5) |

<a id='1'></a>

## [GET]: List weight resources

This method can be used to retrieve a list of current weight resources

### Input

```
GET /weights
```

#### Optional Query Parameters

| Type   | Description                                                                                                                             | Required |
| ------ | --------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| `name` | weight resource name to be used as query                                                                                                | NO       |
| `date` | Date to retrieve a historic version of the weights resource. If no date parameter is provided the most current profile will be returned | NO       |

### Request headers

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
            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
            "date": "2019-11-04",
            "name": "Critical",
            "weight_type": "hepsepc",
            "group_type": "SITES",
            "groups": [
                {
                    "name": "SITE-A",
                    "value": 1673
                },
                {
                    "name": "SITE-B",
                    "value": 1234
                },
                {
                    "name": "SITE-C",
                    "value": 523
                },
                {
                    "name": "SITE-D",
                    "value": 2
                }
            ]
        },
        {
            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
            "date": "2019-11-02",
            "name": "NonCritical",
            "weight_type": "hepsepc",
            "group_type": "SERVICEGROUPS",
            "groups": [
                {
                    "name": "SVGROUP-A",
                    "value": 334
                },
                {
                    "name": "SVGROUP-B",
                    "value": 588
                }
            ]
        }
    ]
}
```

<a id='2'></a>

## [GET]: List A Specific weight resource

This method can be used to retrieve specific weight resource based on its id

### Input

```
GET /weights/{ID}
```

#### Optional Query Parameters

| Type   | Description                                                                                                                             | Required |
| ------ | --------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| `date` | Date to retrieve a historic version of the weights resource. If no date parameter is provided the most current profile will be returned | NO       |

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
            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
            "date": "2019-11-02",
            "name": "NonCritical",
            "weight_type": "hepsepc",
            "group_type": "SERVICEGROUPS",
            "groups": [
                {
                    "name": "SVGROUP-A",
                    "value": 334
                },
                {
                    "name": "SVGROUP-B",
                    "value": 588
                }
            ]
        }
    ]
}
```

<a id='3'></a>

## [POST]: Create a new weight resource

This method can be used to insert a new weight resource

### Input

```
POST /weights
```

#### Optional Query Parameters

| Type   | Description                                                                                                                                 | Required |
| ------ | ------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| `date` | Date to create a historic version of the weights resource. If no date parameter is provided the current date will be supplied automatically | NO       |

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### POST BODY

```json
{
    "name": "weight_set3",
    "weight_type": "hepspec2",
    "group_type": "SITES",
    "groups": [
        { "name": "site-a", "value": 336 },
        { "name": "site-b", "value": 343 },
        { "name": "site-c", "value": 553 },
        { "name": "site-d", "value": 435 },
        { "name": "site-e", "value": 3.33 },
        { "name": "site-f", "value": 323.3 }
    ]
}
```

### Response

Headers: `Status: 201 Created`

#### Response body

Json Response

```json
{
    "status": {
        "message": "Weight resource successfully created",
        "code": "201"
    },
    "data": {
        "id": "{{ID}}",
        "links": {
            "self": "https:///api/v2/weights/{{ID}}"
        }
    }
}
```

<a id='4'></a>

## [PUT]: Update information on an existing weight resource

This method can be used to update information on an existing weight resource

### Input

```
PUT /weights/{ID}
```

#### Optional Query Parameters

| Type   | Description                                                                                                                                                                  | Required |
| ------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| `date` | Date to update a historic version of the weights resource. If no date parameter is provided If no date parameter is provided the current date will be supplied automatically | NO       |

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### PUT BODY

```json
{
    "name": "weight_set3",
    "weight_type": "hepspec2",
    "group_type": "SITES",
    "groups": [
        { "name": "site-a", "value": 1336 },
        { "name": "site-b", "value": 1343 }
    ]
}
```

### Response

Headers: `Status: 200 OK`

#### Response body

Json Response

```json
{
    "status": {
        "message": "Weight resource successfully updated",
        "code": "200"
    },
    "data": {
        "id": "{{ID}}",
        "links": {
            "self": "https:///api/v2/weights/{{ID}}"
        }
    }
}
```

<a id='5'></a>

## [DELETE]: Delete an existing weight resource

This method can be used to delete an existing weight resource

### Input

```
DELETE /weights/{ID}
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
        "message": "Weight resource Successfully Deleted",
        "code": "200"
    }
}
```
