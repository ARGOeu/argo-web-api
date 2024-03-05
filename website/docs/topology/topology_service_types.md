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



## POST: Create list of service types for specific date {#1}

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
    "title": "Service of type A",
    "description": "a short descritpion of service type a"
  },
  {
    "name": "Service_Type_B",
    "title": "Service of type B",
    "description": "a short descritpion of service type b"
  },
  {
    "name": "Service_Type_C",
    "title": "Service of type C",
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


## GET: List service types for specific date {#2}

Lists the available service types per tenant for a specific date

### Input

```
GET /topology/service-types?date=YYYY-MM-DD
```

#### Url Parameters

| Type       | Description                   | Required | Default value |
| ---------- | ----------------------------- | -------- | ------------- |
| `date`     | target a specific date        | NO       | today's date  |
| `mode`     | if stating `mode=combined` then if the tenant has data feeds from other tenants their service lists will be combined in the final results | NO       | empty |


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
      "title": "Service of type A",
      "description": "a short descritpion of service type a"
    },
    {
      "date": "2019-03-03",
      "name": "Service_Type_B",
      "title": "Service of type B",
      "description": "a short descritpion of service type b"
    },
    {
      "date": "2019-03-03",
      "name": "Service_Type_C",
      "title": "Service of type C",
      "description": "a short descritpion of service type c",
      "tags": ["special-service", "beta"]
    }
  ]
}
```

### Combined tenant example:

If the tenant combines data feeds from other tenants (deemed a `combined` tenant) then the optional url property `mode=combined` can
be used to retrieve service types from all other included tenants combined in the final result. Each item retrieved from an included tenant receives an extra `tenant` field to identify its origin

#### Example Request

```
GET /topology/service-types?date=2019-03-03?mode=combined
```

```json
{
  "status": {
    "message": "Success",
    "code": "200"
  },
  "data": [
    {
      "date": "2019-03-03",
      "name": "TenantA-service-type1",
      "title": "Service type1 from tenant A",
      "description": "a short descritpion",
      "tenant": "TenantA"
    },
    {
      "date": "2019-03-03",
      "name": "TenantB-service-type1",
      "title": "Service type2 from tenant B",
      "description": "a short descritpion",
      "tenant": "TenantB"
    }
  ]
}
```


## [DELETE]: Delete list of service types for specific date {#3}

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
