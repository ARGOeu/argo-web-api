---
id: metrics
title: Metrics
---
## API Calls

| Name                                  | Description                                                                     | Shortcut           |
| ------------------------------------- | ------------------------------------------------------------------------------- | ------------------ |
| GET: List Metrics (Admin)   | This method can be used to retrieve a list of all metrics         | [ Description](#1) |
| PUT: Update Metrics (Admin) | This method can be used to update the list of metrics | [ Description](#2) |
| GET: List Metrics  | This method can be used to retrieve a list of metrics (as a tenant user)        | [ Description](#3) |
| PUT: List Metrics by report | This method can be used to retrieve a list of metrics included in a report (as a tenant user) | [ Description](#4) |


<a id='1'></a>

## [GET]: List Metrics (Admin)

This method can be used to retrieve a list of all metrics. This is an administrative method. The Metric list is common for all tenants

### Input

```
GET /admin/metrics
```


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
      "name": "test_metric_1",
      "tags": [
        "network",
        "internal"
      ]
    },
    {
      "name": "test_metric_2",
      "tags": [
        "disk",
        "agent"
      ]
    },
    {
      "name": "test_metric_3",
      "tags": [
        "aai"
      ]
    }
  ]
}
```

<a id='2'></a>

## [PUT]: Update Metrics information
This method is used to update the list of metrics. This is an administrative method. The list of metrics is common for all tenants

### Input

```
PUT /admin/metrics
```

#### PUT BODY
```json
  [
  {
    "name": "metric1",
    "tags": [
      "tag1",
      "tag2"
    ]
  }
]
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
    "message": "Metrics resource succesfully updated",
    "code": "200"
  },
  "data": [
    {
      "name": "metric1",
      "tags": [
        "tag1",
        "tag2"
      ]
    }
  ]
}
```


<a id='3'></a>

## [GET]: List Metrics (as a tenant user)

This method can be used to retrieve the list of metrics as a tenant user. The list of metrics is common for all tenants but accessible from each tenant.

### Input

```
GET /metrics
```


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
      "name": "test_metric_1",
      "tags": [
        "network",
        "internal"
      ]
    },
    {
      "name": "test_metric_2",
      "tags": [
        "disk",
        "agent"
      ]
    },
    {
      "name": "test_metric_3",
      "tags": [
        "aai"
      ]
    }
  ]
}
```


<a id='4'></a>

## [PUT]: List metrics by report (as a tenant user)
This method is used to retrieve a list of metrics that are included in the metric profile of a specific report.

### Input

```
PUT /metrics/by_report/{report_name}
```

#### Url Parameters

| Type          | Description              | Required | Default value |
| ------------- | ------------------------ | -------- | ------------- |
| `report_name` | target a specific report | YES      | none          |
| `date`        | target a specific date   | NO       | today's date  |


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
```json
{
  "status": {
    "message": "Success",
    "code": "200"
  },
  "data": [
    {
      "name": "test_metric_1",
      "tags": [
        "network",
        "internal"
      ]
    }
  ]
}
```
