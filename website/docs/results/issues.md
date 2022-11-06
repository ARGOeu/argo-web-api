---
id:  issues
title: Issues
sidebar_position: 5
---

## API calls to quickly find endpoints with issues

### [GET]: Endpoint with Issues

This method may be used to retrieve a list of problematic endpoints

### Input

```
/issues/{report_name}/endpoints?date=2020-05-01&filter=CRITICAL
```

#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`report_name`| name of the report| YES |  |

#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`date`| Date to view problematic endpoints of | NO |  |
|`filter`| Filter (optionally) problematic endpoints by status value| NO |  |



#### Headers
```
x-api-key: shared_key_value
Accept: application/json or application/xml
```

#### Response Code
```
Status: 200 OK
```

### Response body

###### Example Request:
URL:
```
/api/v2/issues/Critica/endpoints?date=2015-05-01
```
Headers:
```
x-api-key: shared_key_value
Accept: application/json or application/xml

```
###### Example Response:

Code:
```
Status: 200 OK
```
Response body:
```json
{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "timestamp": "2015-05-01T05:00:00Z",
   "endpoint_group": "SITE-A",
   "service": "web_portal",
   "endpoint": "web01.example.gr",
   "status": "WARNING",
   "info": {
    "Url": "http://example.foo/path/to/service"
   }
  },
  {
   "timestamp": "2015-05-01T06:00:00Z",
   "endpoint_group": "SITE-B",
   "service": "object-storage",
   "endpoint": "obj.storage.example.gr",
   "status": "CRITICAL"
  }
 ]
}
```

###### Example Request with property filter=CRITICAL:
URL:
```
/api/v2/issues/Critica/endpoints?date=2015-05-01&filter=CRITICAL
```
Headers:
```
x-api-key: shared_key_value
Accept: application/json or application/xml

```
###### Example Response:

Code:
```
Status: 200 OK
```
Response body:
```json
{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "timestamp": "2015-05-01T06:00:00Z",
   "endpoint_group": "SITE-B",
   "service": "object-storage",
   "endpoint": "obj.storage.example.gr",
   "status": "CRITICAL"
  }
 ]
}
```

### [GET]: Show Metric Issues for a given Group

This method may be used to retrieve a list of problematic metrics that cause issues to a specific Group

### Input

```
/issues/{report_name}/groups/{group_name}/metrics?date=2020-05-01&filter=CRITICAL
```

#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`report_name`| name of the report| YES |  |
|`group_name`| name of the group| YES |  |


#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`date`| Date to view problematic endpoints of | NO |  |
|`filter`| Filter (optionally) problematic endpoints by status value| NO |  |



#### Headers
```
x-api-key: shared_key_value
Accept: application/json or application/xml
```

#### Response Code
```
Status: 200 OK
```

### Response body

###### Example Request:
URL:
```
/api/v2/issues/Critical/groups/SITE-A/metrics?date=2015-05-01
```
Headers:
```
x-api-key: shared_key_value
Accept: application/json or application/xml

```
###### Example Response:

Code:
```
Status: 200 OK
```
Response body:
```json
{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "service": "web_portal",
   "hostname": "web01.example.gr",
   "metric": "http_check",
   "status": "WARNING",
   "info": {
    "Url": "http://example.foo/path/to/service"
   }
  },
  {
   "service": "web_portal",
   "hostname": "web02.example.gr",
   "metric": "http_check",
   "status": "CRITICAL",
   "info": {
    "Url": "http://example2.foo/path/to/service"
   }
  }
 ]
}
```

###### Example Request with property filter=CRITICAL:
URL:
```
/api/v2/issues/Critical/groups/SITE-A/metrics?date=2015-05-01&filter=CRITICAL
```
Headers:
```
x-api-key: shared_key_value
Accept: application/json or application/xml

```
###### Example Response:

Code:
```
Status: 200 OK
```
Response body:
```json
{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "service": "web_portal",
   "hostname": "web02.example.gr",
   "metric": "http_check",
   "status": "CRITICAL",
   "info": {
    "Url": "http://example2.foo/path/to/service"
   }
  }
 ]
}
```

