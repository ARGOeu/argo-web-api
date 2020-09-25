---
id: downtimes
title: Downtimes
---

## API Calls

| Name                                    | Description                                                                       | Shortcut           |
| --------------------------------------- | --------------------------------------------------------------------------------- | ------------------ |
| GET: List Downtimes resources Request   | This method can be used to retrieve a list of current downtime resources.         | [ Description](#1) |
| GET: List a specific Downtimes resource | This method can be used to retrieve a specific downtime resource based on its id. | [ Description](#2) |
| POST: Create a new downtime resource    | This method can be used to create a new downtime resource                         | [ Description](#3) |
| PUT: Update a downtime resource         | This method can be used to update information on an existing downtime resource    | [ Description](#4) |
| DELETE: Delete a downtime resource      | This method can be used to delete an existing downtime resource                   | [ Description](#5) |

<a id='1'></a>

## [GET]: List downtime resources

This method can be used to retrieve a list of current downtime resources

### Input

```
GET /downtimes
```

#### Optional Query Parameters

| Type   | Description                                                                                                                               | Required |
| ------ | ----------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| `name` | downtime resource name to be used as query                                                                                                | NO       |
| `date` | Date to retrieve a historic version of the downtime resource. If no date parameter is provided the most current resource will be returned | NO       |

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
            "endpoints": [
                {
                    "hostname": "host-A",
                    "service": "service-A",
                    "start_time": "2019-10-11T04:00:33Z",
                    "end_time": "2019-10-11T15:33:00Z"
                },
                {
                    "hostname": "host-B",
                    "service": "service-B",
                    "start_time": "2019-10-11T12:00:33Z",
                    "end_time": "2019-10-11T12:33:00Z"
                },
                {
                    "hostname": "host-C",
                    "service": "service-C",
                    "start_time": "2019-10-11T20:00:33Z",
                    "end_time": "2019-10-11T22:15:00Z"
                }
            ]
        },
        {
            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
            "date": "2019-11-02",
            "name": "NonCritical",
            "endpoints": [
                {
                    "hostname": "host-01",
                    "service": "service-01",
                    "start_time": "2019-10-11T02:00:33Z",
                    "end_time": "2019-10-11T23:33:00Z"
                },
                {
                    "hostname": "host-02",
                    "service": "service-02",
                    "start_time": "2019-10-11T16:00:33Z",
                    "end_time": "2019-10-11T16:45:00Z"
                }
            ]
        }
    ]
}
```

<a id='2'></a>

## [GET]: List A Specific downtime resource

This method can be used to retrieve specific downtime resource based on its id

### Input

```
GET /downtimes/{ID}
```

#### Optional Query Parameters

| Type   | Description                                                                                                                                        | Required |
| ------ | -------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| `date` | Date to retrieve a historic version of the downtime resource. If no date parameter is provided the most current downtime resource will be returned | NO       |

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
            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
            "date": "2019-11-04",
            "name": "Critical",
            "endpoints": [
                {
                    "hostname": "host-A",
                    "service": "service-A",
                    "start_time": "2019-10-11T04:00:33Z",
                    "end_time": "2019-10-11T15:33:00Z"
                },
                {
                    "hostname": "host-B",
                    "service": "service-B",
                    "start_time": "2019-10-11T12:00:33Z",
                    "end_time": "2019-10-11T12:33:00Z"
                },
                {
                    "hostname": "host-C",
                    "service": "service-C",
                    "start_time": "2019-10-11T20:00:33Z",
                    "end_time": "2019-10-11T22:15:00Z"
                }
            ]
        }
    ]
}
```

<a id='3'></a>

## [POST]: Create a new downtime resource

This method can be used to insert a new downtime resource

### Input

```
POST /downtimes
```

#### Optional Query Parameters

| Type   | Description                                                                                                                                  | Required |
| ------ | -------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| `date` | Date to create a new historic version of the downtime resource. If no date parameter is provided current date will be supplied automatically | NO       |

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### POST BODY

```json
{
    "name": "downtimes_set",
    "endpoints": [
        {
            "hostname": "host-foo",
            "service": "service-new-foo",
            "start_time": "2019-10-11T23:10:00Z",
            "end_time": "2019-10-11T23:25:00Z"
        },
        {
            "hostname": "host-bar",
            "service": "service-new-bar",
            "start_time": "2019-10-11T23:40:00Z",
            "end_time": "2019-10-11T23:55:00Z"
        }
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
        "message": "Downtimes resource succesfully created",
        "code": "201"
    },
    "data": {
        "id": "{{id}}",
        "links": {
            "self": "https:///api/v2/downtimes/{{id}}"
        }
    }
}
```

<a id='4'></a>

## [PUT]: Update information on an existing downtime resource

This method can be used to update information on an existing downtime resource

### Input

```
PUT /downtimes/{ID}
```

#### Optional Query Parameters

| Type   | Description                                                                                                                                  | Required |
| ------ | -------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| `date` | Date to update a historic version of the downtime resource. If no date parameter is provided the current date will be supplied automatically | NO       |

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### PUT BODY

```json
{
    "name": "downtimes_set",
    "endpoints": [
        {
            "hostname": "updated-host-foo",
            "service": "service-new-foo",
            "start_time": "2019-10-11T23:10:00Z",
            "end_time": "2019-10-11T23:25:00Z"
        },
        {
            "hostname": "updated-host-bar",
            "service": "service-new-bar",
            "start_time": "2019-10-11T23:40:00Z",
            "end_time": "2019-10-11T23:55:00Z"
        }
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
        "message": "Downtimes resource successfully updated",
        "code": "200"
    },
    "data": {
        "id": "{{ID}}",
        "links": {
            "self": "https:///api/v2/downtimes/{{ID}}"
        }
    }
}
```

<a id='5'></a>

## [DELETE]: Delete an existing downtime resource

This method can be used to delete an existing downtime resource

### Input

```
DELETE /downtimes/{ID}
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
        "message": "Downtimes resource Successfully Deleted",
        "code": "200"
    }
}
```
