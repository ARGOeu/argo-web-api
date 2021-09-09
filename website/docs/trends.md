---
id:  trends
title: trends
---

## API calls to find status trends among metrics

This API call displays the top metrics by state (e.g. CRITICAL, WARNING etc) over a period of dates - optionally by monthly aggregation - for a specific report

### [GET]: Daily Status trends in service endpoint metrics
This method may be used to retrieve a list of top status service endpoint metrics. 

### Input

```
/trends/{report_name}/status/metrics?date=2020-05-01
```

#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`report_name`| name of the report| YES |  |

#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`date`| Date to view problematic endpoints of | NO |  |
|`start_date`| define start date to view problematic endpoints over range | NO |  |
|`end_date`| define end date to view problematic endpoints over range | NO |  |
|`top`| integer to define a top number of results displayed | NO |  |
|`granularity`| string value to define if you want monthly granularity in the results - e.g `?granularity=monthly` | NO |  |


#### Headers
```
x-api-key: shared_key_value
Accept: application/json or application/xml
```

#### Response Code
```
Status: 200 OK
```


###### Example Request:
URL:
```
/trends/{report_name}/status/metrics?date=2020-05-01
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
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "metric": "check-2",
          "status": "CRITICAL",
          "events": 40
        }
      ]
    },
    {
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-B",
          "service": "service-A",
          "endpoint": "hosta.example2.foo",
          "metric": "web-check",
          "status": "UNKNOWN",
          "events": 5
        }
      ]
    },
    {
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "metric": "check-1",
          "status": "WARNING",
          "events": 55
        },
        {
          "endpoint_group": "SITE-A",
          "service": "service-B",
          "endpoint": "hostb.example.foo",
          "metric": "web-check",
          "status": "WARNING",
          "events": 12
        }
      ]
    }
  ]
}
```

###### Example Request with Range and top number of results:
URL:
```
/trends/{report_name}/status/metrics?start_date=2020-05-01&end_date=2021-06-15&top=1
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
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "metric": "check-2",
          "status": "CRITICAL",
          "events": 40
        }
      ]
    },
    {
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "metric": "check-1",
          "status": "UNKNOWN",
          "events": 45
        }
      ]
    },
    {
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "metric": "check-1",
          "status": "WARNING",
          "events": 55
        }
      ]
    }
  ]
}
```

###### Example Request with granularity=monthly option enabled:
URL:
```
/trends/{report_name}/status/metrics?start_date=2020-04-01&end_date=2021-05-31&granularity=monthly&top=1
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
      "date": "2015-04",
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "metric": "check-1",
          "status": "CRITICAL",
          "events": 55
        }
      ]
    },
    {
      "date": "2015-04",
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-B",
          "endpoint": "hostb.example.foo",
          "metric": "web-check",
          "status": "UNKNOWN",
          "events": 12
        }
      ]
    },
    {
      "date": "2015-04",
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "metric": "check-2",
          "status": "WARNING",
          "events": 40
        }
      ]
    },
    {
      "date": "2015-05",
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "metric": "check-2",
          "status": "CRITICAL",
          "events": 40
        }
      ]
    },
    {
      "date": "2015-05",
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "metric": "check-1",
          "status": "UNKNOWN",
          "events": 45
        }
      ]
    },
    {
      "date": "2015-05",
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "metric": "check-1",
          "status": "WARNING",
          "events": 55
        }
      ]
    }
  ]
}
```

## API calls to find status trends among endpoints

This API call displays the top endpoints by state (e.g. CRITICAL, WARNING etc) over a period of dates - optionally by monthly aggregation - for a specific report

### [GET]: Daily Status trends in service endpoint metrics
This method may be used to retrieve a list of top status service endpoints. 

### Input

```
/trends/{report_name}/status/endpoints?date=2020-05-01
```

#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`report_name`| name of the report| YES |  |

#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`date`| Date to view problematic endpoints of | NO |  |
|`start_date`| define start date to view problematic endpoints over range | NO |  |
|`end_date`| define end date to view problematic endpoints over range | NO |  |
|`top`| integer to define a top number of results displayed | NO |  |
|`granularity`| string value to define if you want monthly granularity in the results - e.g `?granularity=monthly` | NO |  |


#### Headers
```
x-api-key: shared_key_value
Accept: application/json or application/xml
```

#### Response Code
```
Status: 200 OK
```


###### Example Request:
URL:
```
/trends/{report_name}/status/endpoints?date=2020-05-01
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
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "status": "CRITICAL",
          "duration_in_minutes": 40
        }
      ]
    },
    {
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-B",
          "service": "service-A",
          "endpoint": "hosta.example2.foo",
          "status": "UNKNOWN",
          "duration_in_minutes": 5
        }
      ]
    },
    {
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "status": "WARNING",
          "duration_in_minutes": 55
        },
        {
          "endpoint_group": "SITE-A",
          "service": "service-B",
          "endpoint": "hostb.example.foo",
          "status": "WARNING",
          "duration_in_minutes": 12
        }
      ]
    }
  ]
}
```

###### Example Request with Range and top number of results:
URL:
```
/trends/{report_name}/status/endpoints?start_date=2020-05-01&end_date=2021-06-15&top=1
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
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "status": "CRITICAL",
          "duration_in_minutes": 40
        }
      ]
    },
    {
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "status": "UNKNOWN",
          "duration_in_minutes": 45
        }
      ]
    },
    {
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "status": "WARNING",
          "duration_in_minutes": 55
        }
      ]
    }
  ]
}
```

###### Example Request with granularity=monthly option enabled:
URL:
```
/trends/{report_name}/status/endpoints?start_date=2020-04-01&end_date=2021-05-31&granularity=monthly&top=1
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
      "date": "2015-04",
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "status": "CRITICAL",
          "duration_in_minutes": 55
        }
      ]
    },
    {
      "date": "2015-04",
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-B",
          "endpoint": "hostb.example.foo",
          "status": "UNKNOWN",
          "duration_in_minutes": 12
        }
      ]
    },
    {
      "date": "2015-04",
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "status": "WARNING",
          "duration_in_minutes": 40
        }
      ]
    },
    {
      "date": "2015-05",
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "status": "CRITICAL",
          "duration_in_minutes": 40
        }
      ]
    },
    {
      "date": "2015-05",
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "status": "UNKNOWN",
          "duration_in_minutes": 45
        }
      ]
    },
    {
      "date": "2015-05",
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "status": "WARNING",
          "duration_in_minutes": 55
        }
      ]
    }
  ]
}
```

## API calls to find status trends among services

This API call displays the top services by state (e.g. CRITICAL, WARNING etc) over a period of dates - optionally by monthly aggregation - for a specific report

### [GET]: Daily Status trends in service endpoint metrics
This method may be used to retrieve a list of top status service services. 

### Input

```
/trends/{report_name}/status/services?date=2020-05-01
```

#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`report_name`| name of the report| YES |  |

#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`date`| Date to view problematic services of | NO |  |
|`start_date`| define start date to view problematic services over range | NO |  |
|`end_date`| define end date to view problematic servces over range | NO |  |
|`top`| integer to define a top number of results displayed | NO |  |
|`granularity`| string value to define if you want monthly granularity in the results - e.g `?granularity=monthly` | NO |  |


#### Headers
```
x-api-key: shared_key_value
Accept: application/json or application/xml
```

#### Response Code
```
Status: 200 OK
```


###### Example Request:
URL:
```
/trends/{report_name}/status/services?date=2020-05-01
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
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "status": "CRITICAL",
          "duration_in_minutes": 40
        }
      ]
    },
    {
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-B",
          "service": "service-D",
          "status": "UNKNOWN",
          "duration_in_minutes": 5
        }
      ]
    },
    {
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-X",
          "status": "WARNING",
          "duration_in_minutes": 55
        },
        {
          "endpoint_group": "SITE-A",
          "service": "service-B",
          "status": "WARNING",
          "duration_in_minutes": 12
        }
      ]
    }
  ]
}
```

###### Example Request with Range and top number of results:
URL:
```
/trends/{report_name}/status/services?start_date=2020-05-01&end_date=2021-06-15&top=1
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
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "status": "CRITICAL",
          "duration_in_minutes": 40
        }
      ]
    },
    {
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-B",
          "status": "UNKNOWN",
          "duration_in_minutes": 45
        }
      ]
    },
    {
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-C",
          "status": "WARNING",
          "duration_in_minutes": 55
        }
      ]
    }
  ]
}
```

###### Example Request with granularity=monthly option enabled:
URL:
```
/trends/{report_name}/status/services?start_date=2020-04-01&end_date=2021-05-31&granularity=monthly&top=1
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
      "date": "2015-04",
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "status": "CRITICAL",
          "duration_in_minutes": 55
        }
      ]
    },
    {
      "date": "2015-04",
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-B",
          "status": "UNKNOWN",
          "duration_in_minutes": 12
        }
      ]
    },
    {
      "date": "2015-04",
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-C",
          "status": "WARNING",
          "duration_in_minutes": 40
        }
      ]
    },
    {
      "date": "2015-05",
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-D",
          "status": "CRITICAL",
          "duration_in_minutes": 40
        }
      ]
    },
    {
      "date": "2015-05",
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "status": "UNKNOWN",
          "duration_in_minutes": 45
        }
      ]
    },
    {
      "date": "2015-05",
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-B",
          "status": "WARNING",
          "duration_in_minutes": 55
        }
      ]
    }
  ]
}
```

## API calls to find status trends among endpoint groups

This API call displays the top endpoint groups by state (e.g. CRITICAL, WARNING etc) over a period of dates - optionally by monthly aggregation - for a specific report

### [GET]: Daily Status trends in service endpoint groups
This method may be used to retrieve a list of top status endpoint groups. 

### Input

```
/trends/{report_name}/status/groups?date=2020-05-01
```

#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`report_name`| name of the report| YES |  |

#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`date`| Date to view problematic endpoint groups of | NO |  |
|`start_date`| define start date to view problematic endpoint groups over range | NO |  |
|`end_date`| define end date to view problematic endpoint groups over range | NO |  |
|`top`| integer to define a top number of results displayed | NO |  |
|`granularity`| string value to define if you want monthly granularity in the results - e.g `?granularity=monthly` | NO |  |


#### Headers
```
x-api-key: shared_key_value
Accept: application/json or application/xml
```

#### Response Code
```
Status: 200 OK
```


###### Example Request:
URL:
```
/trends/{report_name}/status/groups?date=2020-05-01
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
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "status": "CRITICAL",
          "duration_in_minutes": 40
        }
      ]
    },
    {
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-B",
          "status": "UNKNOWN",
          "duration_in_minutes": 5
        }
      ]
    },
    {
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-C",
          "status": "WARNING",
          "duration_in_minutes": 55
        },
        {
          "endpoint_group": "SITE-D",
          "status": "WARNING",
          "duration_in_minutes": 12
        }
      ]
    }
  ]
}
```

###### Example Request with Range and top number of results:
URL:
```
/trends/{report_name}/status/groups?start_date=2020-05-01&end_date=2021-06-15&top=1
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
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "status": "CRITICAL",
          "duration_in_minutes": 40
        }
      ]
    },
    {
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-B",
          "status": "UNKNOWN",
          "duration_in_minutes": 45
        }
      ]
    },
    {
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-C",
          "status": "WARNING",
          "duration_in_minutes": 55
        }
      ]
    }
  ]
}
```

###### Example Request with granularity=monthly option enabled:
URL:
```
/trends/{report_name}/status/groups?start_date=2020-04-01&end_date=2021-05-31&granularity=monthly&top=1
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
      "date": "2015-04",
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "status": "CRITICAL",
          "duration_in_minutes": 55
        }
      ]
    },
    {
      "date": "2015-04",
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-B",
          "status": "UNKNOWN",
          "duration_in_minutes": 12
        }
      ]
    },
    {
      "date": "2015-04",
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-C",
          "status": "WARNING",
          "duration_in_minutes": 40
        }
      ]
    },
    {
      "date": "2015-05",
      "status": "CRITICAL",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "status": "CRITICAL",
          "duration_in_minutes": 40
        }
      ]
    },
    {
      "date": "2015-05",
      "status": "UNKNOWN",
      "top": [
        {
          "endpoint_group": "SITE-B",
          "status": "UNKNOWN",
          "duration_in_minutes": 45
        }
      ]
    },
    {
      "date": "2015-05",
      "status": "WARNING",
      "top": [
        {
          "endpoint_group": "SITE-C",
          "status": "WARNING",
          "duration_in_minutes": 55
        }
      ]
    }
  ]
}
```

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
|`start_date`| define start date to view problematic endpoints over range | NO |  |
|`end_date`| define end date to view problematic endpoints over range | NO |  |
|`top`| integer to define a top number of results displayed | NO |  |
|`granularity`| string value to define if you want monthly granularity in the results - e.g `?granularity=monthly` | NO |  |


#### Headers
```
x-api-key: shared_key_value
Accept: application/json or application/xml
```

#### Response Code
```
Status: 200 OK
```


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

###### Example Request with Range and top number of results:
URL:
```
/trends/{report_name}/flapping/metrics?start_date=2020-05-01&end_date=2021-06-15&top=3
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
      "flapping": 255
    },
    {
      "endpoint_group": "SITE-A",
      "service": "service-A",
      "endpoint": "hosta.example.foo",
      "metric": "check-2",
      "flapping": 340
    },
    {
      "endpoint_group": "SITE-A",
      "service": "service-B",
      "endpoint": "hostb.example.foo",
      "metric": "web-check",
      "flapping": 112
    }
  ]
}
```

###### Example Request with granularity=monthly option enabled:
URL:
```
/trends/{report_name}/flapping/metrics?start_date=2020-04-01&end_date=2021-05-31&granularity=monthly&top=3
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
      "date": "2015-04",
      "top": [
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
          "endpoint_group": "SITE-XB",
          "service": "service-XA",
          "endpoint": "hosta.examplex2.foo",
          "metric": "web-check",
          "flapping": 25
        }
      ]
    },
    {
      "date": "2015-05",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "metric": "check-1",
          "flapping": 100
        },
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "metric": "check-2",
          "flapping": 72
        },
        {
          "endpoint_group": "SITE-A",
          "service": "service-B",
          "endpoint": "hostb.example.foo",
          "metric": "web-check",
          "flapping": 20
        }
      ]
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
|`start_date`| define start date to view problematic endpoints over range | NO |  |
|`end_date`| define end date to view problematic endpoints over range | NO |  |
|`top`| integer to define a top number of results displayed | NO |  |
|`granularity`| string value to define if you want monthly granularity in the results - e.g `?granularity=monthly` | NO |  |


#### Headers
```
x-api-key: shared_key_value
Accept: application/json or application/xml
```

#### Response Code
```
Status: 200 OK
```


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

###### Example Request with date range and top number of results:
URL:
```
/trends/{report_name}/flapping/endpoints?start_date=2020-05-01&end_date=2020-05-15&top=2
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
      "flapping": 83
    },
    {
      "endpoint_group": "SITE-A",
      "service": "service-B",
      "endpoint": "hostb.example.foo",
      "flapping": 53
    }
  ]
}
```

###### Example Request with granularity=monthly option enabled
URL:
```
/trends/{report_name}/flapping/endpoints?start_date=2020-04-01&end_date=2020-05-31&top=2&granularity=monthly
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
      "date": "2015-04",
      "top": [
        {
          "endpoint_group": "SITE-XB",
          "service": "service-XA",
          "endpoint": "hosta.exampleX2.foo",
          "flapping": 35
        },
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "flapping": 25
        }
      ]
    },
    {
      "date": "2015-05",
      "top": [
        {
          "endpoint_group": "SITE-A",
          "service": "service-A",
          "endpoint": "hosta.example.foo",
          "flapping": 103
        },
        {
          "endpoint_group": "SITE-A",
          "service": "service-B",
          "endpoint": "hostb.example.foo",
          "flapping": 19
        }
      ]
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
|`start_date`| define start date to view problematic endpoints over range | NO |  |
|`end_date`| define end date to view problematic endpoints over range | NO |  |
|`top`| integer to define a top number of results displayed | NO |  |
|`granularity`| string value to define if you want monthly granularity in the results - e.g `?granularity=monthly` | NO |  |


#### Headers
```
x-api-key: shared_key_value
Accept: application/json or application/xml
```

#### Response Code
```
Status: 200 OK
```


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

###### Example Request with date range and top number of results:
URL:
```
/trends/{report_name}/flapping/services?start_date=2020-05-01&end_date=2020-07-05&top=1
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
      "flapping": 955
    }
  ]
}
```

###### Example Request with granularity=monthly option enabled:
URL:
```
/trends/{report_name}/flapping/services?start_date=2020-04-01&end_date=2020-05-31&top=1&granularity=monthly
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
   "date": "2015-04",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "flapping": 25
    }
   ]
  },
  {
   "date": "2015-05",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "service": "service-A",
     "flapping": 98
    }
   ]
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
|`start_date`| define start date to view problematic endpoints over range | NO |  |
|`end_date`| define end date to view problematic endpoints over range | NO |  |
|`top`| integer to define a top number of results displayed | NO |  |
|`granularity`| string value to define if you want monthly granularity in the results - e.g `?granularity=monthly` | NO |  |


#### Headers
```
x-api-key: shared_key_value
Accept: application/json or application/xml
```

#### Response Code
```
Status: 200 OK
```

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

###### Example Request with date range and top number of results:
URL:
```
/trends/{report_name}/flapping/groups?start_date=2020-05-01&end_date=2020-05-03&top=1
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
      "flapping": 75
    }
  ]
}
```

###### Example Request with granularity=monthly option enabled
URL:
```
/trends/{report_name}/flapping/groups?start_date=2020-04-01&end_date=2020-05-31&top=1&granularity=monthly
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
   "date": "2015-04",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "flapping": 35
    }
   ]
  },
  {
   "date": "2015-05",
   "top": [
    {
     "endpoint_group": "SITE-A",
     "flapping": 66
    }
   ]
  }
 ]
}
```