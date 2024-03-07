---
id: recomputations
title: Recomputation Requests
sidebar_position: 7
---

## API Calls for listing existing and creating new recomputation requests

Name                                     | Description                                                                            | Shortcut
---------------------------------------- | -------------------------------------------------------------------------------------- | ------------------
GET: List Recomputation Requests         | This method can be used to retrieve a list of current Recomputation requests.          | [ Description](#1)
GET: Get a specific recomputation by id | This method can be used to retrieve a specific recomputation by id                     | [ Description](#2)
POST: Create a new recomputation request | This method can be used to insert a new recomputation request onto the Compute Engine. | [ Description](#3)
DELETE: Delete a specific recomputation  | This method can be used to delete a specific recomputation. | [ Description](#4)
POST: change status  | This method can be used to change status of a specific recomputation. | [ Description](#5)
DELETE: Reset status of recomputation  | This method can be used to reset status of a specific recomputation. | [ Description](#6)


## [GET]: List Recomputation Requests {#1}
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


## [GET]: Get specific recomputation request by id {#2}
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



## [POST]: Create a new recomputation request {#3}
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


## [DELETE]: Delete a specific recomputation {#4}

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


## [POST]: Change status of recomputation {#5}

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


## [DELETE]: Reset status of a specific recomputation {#6}

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

# Recomputations that exclude metrics

There is also the ability to run a recomputation and exclude specific metrics. During the recomputation period the metrics that are considered excluded don't take place into any operation or aggregation thus they don't affect their endpoints at all. 

To declare a recomputation that excludes metric you must use the special field "exclude_metrics" in the recomputation and add an array of metrics to be excluded (You can limit the scope also by "group", "service" and "hostname")

For example:

```json
{
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e777",
   "requester_name": "John Doe",
   "requester_email": "johndoe@example.com",
   "reason": "issue with metric checks",
   "start_time": "2022-01-10T12:00:00Z",
   "end_time": "2022-01-10T23:00:00Z",
   "report": "Default",
   "exclude_metrics": [
    {
     "metric": "check-1"
    },
    {
     "metric": "check-2",
     "hostname": "host1.example.com"
    },
    {
     "metric": "check-3",
     "group": "Affected-Site"
    }
   ]
  }
  ```

  If you specify a rule that includes only a `metric` then this type of metric will be excluded globally from all endpoints and groups
  If you specify a rule that includes a `metric` and another field such as `hostname`, `service` or `group` then the rule is scoped accordingly to a specific group or service type or hostname and the metric that belongs there. The field `metric` is mandatory.

  # Recomputations that exclude monitoring sources

There is also the ability to run a recomputation and exclude a monitoring source (e.g. specific monitoring box). This is especially usefull in HA situations where one of the available monitoring sources might have issues for a specific period of time.

For example:

```json
{
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e777",
   "requester_name": "John Doe",
   "requester_email": "johndoe@example.com",
   "reason": "issue with metric checks",
   "start_time": "2022-01-10T12:00:00Z",
   "end_time": "2022-01-10T23:00:00Z",
   "report": "Default",
   "exclude_monitoring_source": [
    {
        "host":"monitoring_node01.example.foo",
        "start_time": "2022-01-10T12:00:00Z",
        "end_time": "2022-01-10T23:00:00Z"
    }
   ]
}
```