---
id: recomputations
title: Recomputation Requests
---

## API Calls for listing existing and creating new recomputation requests

Name                                     | Description                                                                            | Shortcut
---------------------------------------- | -------------------------------------------------------------------------------------- | ------------------
GET: List Recomputation Requests         | This method can be used to retrieve a list of current Recomputation requests.          | [ Description](#1)
GET: Get a specific recomputatiomn by id | This method can be used to retrieve a specific recomputation by id                     | [ Description](#2)
POST: Create a new recomputation request | This method can be used to insert a new recomputation request onto the Compute Engine. | [ Description](#3)
DELETE: Delete a specific recomputation  | This method can be used to delete a specific recomputation. | [ Description](#4)
POST: change status  | This method can be used to change status of a specific recomputation. | [ Description](#5)
DELETE: Reset status of recomputation  | This method can be used to reset status of a specific recomputation. | [ Description](#6)


<a id='1'></a>

## [GET]: List Recomputation Requests
This method can be used to retrieve a list of current Recomputation requests.

### Input

```
GET /recomputations
```

#### Optional Query Parameters

| Type     | Description                                                                                                                             | Required |
| -------- | --------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| `report` | Filter recomputations by report name                                                                                                    | NO       |
| `date`   | Specific date to retrieve all relevant recomputations that their period include this date                                               | NO       |



#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### Response
Headers: `Status: 200 OK`

#### Response body
Json Response

```json
{
"root": [
     {
          "id": "56db43ee-f331-46ca-b0fd-4555b4aa1cfc",
          "requester_name": "John Doe",
          "requester_email": "JohnDoe@foo.com",
          "reason": "power cuts",
          "start_time": "2015-01-10T12:00:00Z",
          "end_time": "2015-01-30T23:00:00Z",
          "report": "Critical",
          "exclude": [
           "Gluster"
          ],
          "status": "running",
          "timestamp": "2015-02-01T14:58:40",
          "history": [
              { 
                  "status": "pending", 
                  "timestamp" : "2015-02-01T14:58:40"
              },
              { 
                  "status": "approved", 
                  "timestamp" : "2015-02-02T08:58:40"
              },
              { 
                  "status": "running", 
                  "timestamp" : "2015-02-02T09:10:40"
              },

          ]
     },
     {
          "id": "f68b43ee-f331-46ca-b0fd-4555b4aa1cfc",
          "requester_name": "John Doe",
          "requester_email": "JohnDoe@foo.com",
          "reason": "power cuts",
          "start_time": "2015-03-10T12:00:00Z",
          "end_time": "2015-03-30T23:00:00Z",
          "report": "OPS-Critical",
          "exclude": [
           "Gluster"
          ],
          "status": "running",
          "timestamp": "2015-02-01T14:58:40",
          "history": [
              { 
                  "status": "pending", 
                  "timestamp" : "2015-04-01T14:58:40"
              },
              { 
                  "status": "approved", 
                  "timestamp" : "2015-04-02T08:58:40"
              },
              { 
                  "status": "running", 
                  "timestamp" : "2015-04-02T09:10:40"
              },

          ]
     }
 ]
}
```

### Example Request #2 
```
GET /recomputations?date=2015-03-15
```

#### Response
Headers: `Status: 200 OK`

#### Response body
Json Response

```json
{
"root": [
     {
          "id": "f68b43ee-f331-46ca-b0fd-4555b4aa1cfc",
          "requester_name": "John Doe",
          "requester_email": "JohnDoe@foo.com",
          "reason": "power cuts",
          "start_time": "2015-03-10T12:00:00Z",
          "end_time": "2015-03-30T23:00:00Z",
          "report": "OPS-Critical",
          "exclude": [
           "Gluster"
          ],
          "status": "running",
          "timestamp": "2015-02-01T14:58:40",
          "history": [
              { 
                  "status": "pending", 
                  "timestamp" : "2015-04-01T14:58:40"
              },
              { 
                  "status": "approved", 
                  "timestamp" : "2015-04-02T08:58:40"
              },
              { 
                  "status": "running", 
                  "timestamp" : "2015-04-02T09:10:40"
              },

          ]
     }
 ]
}
```

### Example Request #3 
```
GET /recomputations?report=OPS-Critical
```

#### Response
Headers: `Status: 200 OK`

#### Response body
Json Response

```json
{
"root": [
     {
          "id": "f68b43ee-f331-46ca-b0fd-4555b4aa1cfc",
          "requester_name": "John Doe",
          "requester_email": "JohnDoe@foo.com",
          "reason": "power cuts",
          "start_time": "2015-03-10T12:00:00Z",
          "end_time": "2015-03-30T23:00:00Z",
          "report": "OPS-Critical",
          "exclude": [
           "Gluster"
          ],
          "status": "running",
          "timestamp": "2015-02-01T14:58:40",
          "history": [
              { 
                  "status": "pending", 
                  "timestamp" : "2015-04-01T14:58:40"
              },
              { 
                  "status": "approved", 
                  "timestamp" : "2015-04-02T08:58:40"
              },
              { 
                  "status": "running", 
                  "timestamp" : "2015-04-02T09:10:40"
              },

          ]
     }
 ]
}
```

<a id='2'></a>

## [GET]: Get specific recomputation request by id
This method can be used to retrieve a specific recomputation request by its id

### Input

```
GET /recomputations/{ID}
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### Response
Headers: `Status: 200 OK`

#### Response body
Json Response

```json
{
"root": [
     {
          "id": "56db43ee-f331-46ca-b0fd-4555b4aa1cfc",
          "requester_name": "John Doe",
          "requester_email": "JohnDoe@foo.com",
          "reason": "power cuts",
          "start_time": "2015-01-10T12:00:00Z",
          "end_time": "2015-01-30T23:00:00Z",
          "report": "Critical",
          "exclude": [
           "Gluster"
          ],
          "status": "running",
          "timestamp": "2015-02-01T14:58:40",
          "history": [
              { 
                  "status": "pending", 
                  "timestamp" : "2015-02-01T14:58:40"
              },
              { 
                  "status": "approved", 
                  "timestamp" : "2015-02-02T08:58:40"
              },
              { 
                  "status": "running", 
                  "timestamp" : "2015-02-02T09:10:40"
              },

          ]
     }
```


<a id='3'></a>

## [POST]: Create a new recomputation request
This method can be used to insert a new recomputation request onto the Compute Engine.

### Input

```
POST /recomputations
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### Parameters

Type         | Description                                                                                                                                            | Required | Default value
------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------ | -------- | -------------
`start_time` | UTC time in W3C format                                                                                                                                 | YES      |
`end_time`   | UTC time in W3C format                                                                                                                                 | YES      |
`reason`     | Explain the need for a recomputation                                                                                                                   | YES      |
`requester_name`     | The name of the person submitting the recomputation                                                                                            | YES      |
`requester_email`    | The email of the person submitting the recomputation                                                                                           | YES      |
`report`     | Report for which the recomputation is requested                                                                                                        | YES      |
`exclude`    | Groups to be excluded from recomputation. If more than one group are to be excluded use the parameter as many times as needed within the same API call | NO       |

### Response
Headers: `Status: 201 Created`

<a id='4'></a>

## [DELETE]: Delete a specific recomputation

```
DELETE /recomputations/{ID}
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

### Response
`Status 200 OK`

<a id='5'></a>

## [POST]: Change status of recomputation

```
POST /recomputations/{ID}/status
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

### POST body
```json
{
  "status" : "approved"
}
```

Eligible recomputation status values:

- pending
- approved
- rejected
- running
- done

If recomputation status input not in eligible values the api will respond with status code `404`:`conflict`

### Response
`Status 200 OK`

#### Response body
Json Response

```json
{
 "status": {
  "message": "Recomputation status updated successfully to: approved",
  "code": "200"
 }
}
```


<a id='6'></a>

## [DELETE]: Reset status of a specific recomputation

```
DELETE /recomputations/{ID}/status
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

### Response
`Status 200 OK`

#### Response body
Json Response

```json
{
 "status": {
  "message": "Recomputation status reset to: pending",
  "code": "200"
 }
}
```