---
id:  trends
title: trends
---

## API calls to find flapping trends among metrics, endpoints, services and groups by Date

Flapping items are the ones that change state too frequently causing a lot of alarms and notifications. State flapping might be the case of wrong configuration. 

### [GET]: Daily Flapping trends in service endpoint metrics
This method may be used to retrieve a list of top flapping service endpoint metrics. 

### Input

```
/trends/{report_name}/flapping/metrics?date=2020-05-01
```

#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`report_name`| name of the report| YES |  |

#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`date`| Date to view problematic endpoints of | NO |  |


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
/trends/{report_name}/flapping/metrics?date=2020-05-01
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
      "endpoint_group": "SITE-A",
      "service": "service-A",
      "endpoint": "hosta.example.foo",
      "metric": "check-1",
      "flapping": 55
    },
    {
      "endpoint_group": "SITE-A",
      "service": "service-A",
      "endpoint": "hosta.example.foo",
      "metric": "check-2",
      "flapping": 40
    },
    {
      "endpoint_group": "SITE-A",
      "service": "service-B",
      "endpoint": "hostb.example.foo",
      "metric": "web-check",
      "flapping": 12
    },
    {
      "endpoint_group": "SITE-B",
      "service": "service-A",
      "endpoint": "hosta.example2.foo",
      "metric": "web-check",
      "flapping": 5
    }
  ]
}
```


### [GET]: Daily Flapping trends in service endpoints 
This method may be used to retrieve a list of top flapping service endpoints

### Input

```
/trends/{report_name}/flapping/endpoints?date=2020-05-01
```

#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`report_name`| name of the report| YES |  |

#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`date`| Date to view problematic endpoints of | NO |  |


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
/trends/{report_name}/flapping/endpoints?date=2020-05-01
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
      "endpoint_group": "SITE-A",
      "service": "service-A",
      "endpoint": "hosta.example.foo",
      "flapping": 55
    },
    {
      "endpoint_group": "SITE-A",
      "service": "service-B",
      "endpoint": "hostb.example.foo",
      "flapping": 12
    },
    {
      "endpoint_group": "SITE-B",
      "service": "service-A",
      "endpoint": "hosta.example2.foo",
      "flapping": 5
    }
  ]
}
```

### [GET]: Daily Flapping trends in services
This method may be used to retrieve a list of top flapping services

### Input

```
/trends/{report_name}/flapping/services?date=2020-05-01
```

#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`report_name`| name of the report| YES |  |

#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`date`| Date to view problematic endpoints of | NO |  |


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
/trends/{report_name}/flapping/services?date=2020-05-01
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
      "endpoint_group": "SITE-A",
      "service": "service-A",
      "flapping": 55
    },
    {
      "endpoint_group": "SITE-A",
      "service": "service-B",
      "flapping": 12
    },
    {
      "endpoint_group": "SITE-B",
      "service": "service-A",
      "flapping": 5
    }
  ]
}
```

### [GET]: Daily Flapping trends in endpoint groups
This method may be used to retrieve a list of top endpoint groups

### Input

```
/trends/{report_name}/flapping/groups?date=2020-05-01
```

#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`report_name`| name of the report| YES |  |

#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`date`| Date to view problematic endpoints of | NO |  |


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
/trends/{report_name}/flapping/groups?date=2020-05-01
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
      "endpoint_group": "SITE-A",
      "flapping": 55
    },
    {
      "endpoint_group": "SITE-B",
      "flapping": 5
    }
  ]
}
```