---
id: downtimes
title: Downtimes
---

## API Calls

| Name                                    | Description                                                                       | Shortcut           |
| --------------------------------------- | --------------------------------------------------------------------------------- | ------------------ |
| GET: List Downtimes resources Request   | This method can be used to retrieve a list of current downtime resources per date.         | [ Description](#1) |
| POST: Create a new downtime resource    | This method can be used to create a new downtime resource                         | [ Description](#2) |
| DELETE: Delete a downtime resource      | This method can be used to delete an existing downtime resource                   | [ Description](#3) |

<a id='1'></a>

## [GET]: List downtime resources

This method can be used to retrieve a list of current downtime resources per date

### Input

```
GET /downtimes?date=YYYY-MM-DD
```

#### Optional Query Parameters

| Type   | Description                                                                                                                               | Required |
| ------ | ----------------------------------------------------------------------------------------------------------------------------------------- | --------                                                                                           | NO       |
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
            "date": "2019-11-04",
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



<a id='2'></a>

## [POST]: Create a new downtime resource

This method can be used to insert a new downtime resource

### Input

```
POST /downtimes?date=YYYY-MM-DD
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
  "message": "Downtimes set created for date: 2019-11-29",
  "code": "201"
 }
}
```

<a id='3'></a>

## [DELETE]: Delete an existing downtime resource

This method can be used to delete an existing downtime resource

### Input

```
DELETE /downtimes?date=YYYY-MM-DD
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
  "message": "Downtimes set deleted for date: 2019-10-11",
  "code": "200"
 }
}
```
