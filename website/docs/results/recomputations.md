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

| Parameter                   | Type     | Description                                                                                                                                     | Required |
|----------------------------|----------|-------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| `start_time`               | `string` | UTC timestamp in W3C format (e.g., `2025-05-07T12:00:00Z`)                                                                                      | Yes      |
| `end_time`                 | `string` | UTC timestamp in W3C format (e.g., `2025-05-09T09:00:00Z`)                                                                                      | Yes      |
| `reason`                   | `string` | Explanation of the need for the recomputation                                                                                                   | Yes      |
| `requester_name`           | `string` | Name of the person submitting the recomputation                                                                                                 | Yes      |
| `requester_email`          | `string` | Email of the person submitting the recomputation                                                                                                | Yes      |
| `report`                   | `string` | Report for which the recomputation is requested                                                                                                 | Yes      |
| `exclude_monitoring_source` | `array`  | List of monitoring sources to exclude during recomputation, each with `host`, `start_time`, and `end_time`                                     | No       |
| `exclude`                  | `array`  | List of group names to exclude from recomputation                                                                                               | No       |
| `exclude_metrics`          | `array`  | List of specific metrics (optionally scoped by hostname, service or group) to exclude from recomputation                                       | No       |
| `applied_status_changes`   | `array`  | Manual status overrides for specific topology elements (group, service, endpoint, or metric)                                                   | No       |

### Response
Headers: `Status: 201 Created`

<a id='4'></a>

#### Request Body

The recomputation request body can include any combination of the following optional elements:

- `"exclude_monitoring_source"` — to exclude data from specified monitoring sources during a time window.
- `"exclude"` — to exclude specific groups from availability/reliability calculations.
- `"exclude_metrics"` — to exclude particular metrics, optionally scoped by hostname, service, or group.
- `"applied_status_changes"` — to manually override statuses of monitored topology items.

All these elements are **optional** and can be included individually or together depending on the recomputation needs.

---

#### Example Request Body

```json
{
  "requester_name": "John Doe",
  "requester_email": "johndoe@foo.gr",
  "reason": "Recomputation request including all types",
  "start_time": "2025-05-07T12:00:00Z",
  "end_time": "2025-05-09T09:00:00Z",
  "report": "Report-A",

  "exclude_monitoring_source": [
    {
      "host": "monitoring_node01.example.foo",
      "start_time": "2022-01-10T12:00:00Z",
      "end_time": "2022-01-10T23:00:00Z"
    }
  ],

  "exclude": [
    "Group-1",
    "Group-2"
  ],

  "exclude_metrics": [
    { "metric": "check-1" },
    { "metric": "check-2", "hostname": "host1.example.com" }
  ],

  "applied_status_changes": [
    { "group": "Group-A", "state": "CRITICAL" },
    { "service": "Service-a", "state": "OK" }
  ]
}
```




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



