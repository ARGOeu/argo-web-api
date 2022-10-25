---
id: topology_service_types
title: Topology Service Types
sidebar_position: 1
---

## API calls for handling topology list of service types

| Name                                            | Description                                                         | Shortcut                     |
| ----------------------------------------------- | ------------------------------------------------------------------- | ---------------------------- |
| POST: Create list of service types for specific date   | Creates a daily list of available service types for the specific tenant| <a href="#1">Description</a> |
| GET: List service types for specific date      | Lists the available service types per tenant for a specific date   | <a href="#2">Description</a> |
| DELETE: Delete list of service types for specific date | Delete List of available service types for a specific date | <a href="#3">Description</a> |


<a id="1"></a>

## POST: Create list of service types for specific date 

Creates a daily list of available service types for the specific tenant

### Input

```
POST /topology/service-types?date=YYYY-MM-DD
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
    "name": "Service_Type_A",
    "description": "a short descritpion of service type a"
  },
  {
    "name": "Service_Type_B",
    "description": "a short descritpion of service type b"
  },
  {
    "name": "Service_Type_C",
    "description": "a short descritpion of service type c",
    "tags": ["special-service", "beta"]
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
    "message": "Topology of 3 service types created for date: YYYY-MM-DD",
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
    "message": "Topology list of service types already exists for date: YYYY-MM-DD, please either update it or delete it first!",
    "code": "409"
}
```

User can proceed with either updating the existing topology OR deleting before trying to create it anew

<a id="2"></a>

## GET: List service types for specific date

Lists the available service types per tenant for a specific date

### Input

```
GET /topology/service-types?date=YYYY-MM-DD
```

#### Url Parameters

| Type       | Description                   | Required | Default value |
| ---------- | ----------------------------- | -------- | ------------- |
| `date`     | target a specific date        | NO       | today's date  |


#### Headers

```
x-api-key: secret_key_value
Accept: application/json
```

#### Example Request

```
GET /topology/service-types?date=2019-03-03
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
      "date": "2019-03-03",
      "name": "Service_Type_A",
      "description": "a short descritpion of service type a"
    },
    {
      "date": "2019-03-03",
      "name": "Service_Type_B",
      "description": "a short descritpion of service type b"
    },
    {
      "date": "2019-03-03",
      "name": "Service_Type_C",
      "description": "a short descritpion of service type c",
      "tags": ["special-service", "beta"]
    }
  ]
}
```

<a id='3'></a>

## [DELETE]: Delete list of service types for specific date

This method can be used to delete all service type items contributing to the list of available service types of a specific date

### Input

```
DELETE /topology/service-types?date=YYYY-MM-DD
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
    "message": "Topology of 3 service types deleted for date: 2019-12-12",
    "code": "200"
}
```
