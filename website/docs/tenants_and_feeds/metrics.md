---
id: metrics
title: Available Metrics and Tags
sidebar_position: 4
---
## API Calls

| Name                                  | Description                                                                     | Shortcut           |
| ------------------------------------- | ------------------------------------------------------------------------------- | ------------------ |
| GET: List Metrics (Admin)   | This method can be used to retrieve a list of all metrics         | [ Description](#1) |
| PUT: Update Metrics (Admin) | This method can be used to update the list of metrics | [ Description](#2) |
| GET: List Metrics  | This method can be used to retrieve a list of metrics (as a tenant user)        | [ Description](#3) |
| PUT: List Metrics by report | This method can be used to retrieve a list of metrics included in a report (as a tenant user) | [ Description](#4) |



## [GET]: List Metrics (Admin) {#1}

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


## [PUT]: Update Metrics information {#2}
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



## [GET]: List Metrics (as a tenant user) {#3}

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


## [PUT]: List metrics by report (as a tenant user) {#4}
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

Some metric results have additional information regarding the specific service endpoint such as it's Url, certificate DN etc... If this information is available it will be displayed under each service endpoint along with status results. Also some metrics might have a changed status due to a defined threshold rule being applied (see more about [Threshold profiles](/profiles_and_reports/threshold_profiles.md)). Thus they will include additional information such as the original status value (field name: `original_status`), the threshold rule applied (field name: `threshold_rule_applied`) and the actual data (field name: `actual_data`) that the rule has been applied to. For example:

Headers: `Status: 200 OK`

#### Response body

Json Response

```json
{
   "root": [
     {
       "Name": "www.example.com",
       "info": {
                  "Url": "https://example.com/path/to/service/check"
               },
       "Metrics": [
         {
           "Name": "httpd_check",
           "Service": "httpd",
           "Details": [
             {
               "Timestamp": "2015-06-20T12:00:00Z",
               "Value": "OK",
               "Summary": "httpd is ok",
               "Message": "all checks ok"
             },
             {
               "Timestamp": "2015-06-20T23:00:00Z",
               "Value": "OK",
               "Summary": "httpd is ok",
               "Message": "all checks ok"
             }
           ]
         },
         {
           "Name": "httpd_memory",
           "Service": "httpd",
           "Details": [
             {
               "Timestamp": "2015-06-20T06:00:00Z",
               "Value": "OK",
               "Summary": "memcheck ok",
               "Message": "memory under 20%"
             },
             {
               "Timestamp": "2015-06-20T09:00:00Z",
               "Value": "OK",
               "Summary": "memcheck ok",
               "Message": "memory under 20%"
             },
             {
               "Timestamp": "2015-06-20T18:00:00Z",
               "Value": "CRITICAL",
               "Summary": "memcheck ok",
               "Message": "memory under 20%",
               "original_status": "OK",
               "threshold_rule_applied": "reserved_memory=0.1;0.1:0.2;0.2:0.5",
               "actual_data": "reserved_memory=0.15"
             },
           ]
         }
       ]
     }
   ]
 }
```
