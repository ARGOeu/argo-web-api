---
id:  issues
title: issues
---

## API calls to quicky find endpoints with issues

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
|`filter`| Filter (optinally) problematic endpoints by status value| NO |  |



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
Reponse body:
```
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
   "endpoint": web01.example.gr",
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
Reponse body:
```
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

