---
id: metric_results
title: Metric Results
sidebar_position: 4
---

## API call for retrieving detailed metric result.

### [GET]: Metric Result

This method may be used to retrieve a specific service metric result.

### Input

```
/metric_result/{hostname}/{metric_name}?[exec_time]
```

#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`hostname`| hostname of service endpoint| YES |  |
|`metric_name`| name of the metric| YES |  |

#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`exec_time`| The execution date of query in zulu format| YES |  |
|`service`| Filter by service type | NO | |

___Notes___:
`exec_time` : The execution date of query in zulu format. In order to get the correct execution time get status results for all metrics (under a given endpoint, service and endpoint group). ( GET /status/\{report_name\}/\{lgroup_type\}/\{lgroup_name\}/services/\{service_name\}/endpoints/\{endpoint_name\}/metrics List)

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
/api/v2/metric_result/www.example.com/httpd_check?exec_time=2015-06-20T12:00:00Z
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
```
 {
   "root": [
     {
       "Name": "www.example.com",
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
             }
           ]
         }
       ]
     }
   ]
 }
 
```

###### Example Request with service filter:
URL:
```
/api/v2/metric_result/www.example.com/json_check?exec_time=2018-06-20T12:00:00Z&service=object-storage
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
```
 {
   "root": [
     {
       "Name": "www.example.com",
       "Metrics": [
         {
           "Name": "write-obj",
           "Service": "object-storage",
           "Details": [
             {
               "Timestamp": "2015-06-20T12:00:00Z",
               "Value": "OK",
               "Summary": "file written",
               "Message": "all checks ok"
             }
           ]
         }
       ]
     }
   ]
 }
 
```

## [GET]: Multiple Metric Results for a specific host, on a specific day

This method may be used to retrieve multiple service metric results for a specific host on a specific day

### Input

```
/metric_result/{hostname}?[exec_time]
```

#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`hostname`| hostname of service endpoint| YES |  |
|`service`| Filter by service type | NO | |
#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`exec_time`| The execution date of query in zulu format - time part is irrelevant (can be 00:00:00Z) | YES |  |
|`filter`| Filter metric results by statuses: non-ok, ok, critical, warning | NO |  |

___Notes___:
`exec_time` : The specific date of query in zulu format. The time part of the date is irrelevant because all metrics of that day are returned. ( GET /status/\{report_name\}/\{lgroup_type\}/\{lgroup_name\}/services/\{service_name\}/endpoints/\{endpoint_name\}/metrics List)

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
/api/v2/metric_result/www.example.com?exec_time=2015-06-20T00:00:00Z
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
```
{
   "root": [
     {
       "Name": "www.example.com",
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
               "Timestamp": "2015-06-20T18:00:00Z",
               "Value": "CRITICAL",
               "Summary": "httpd is critical",
               "Message": "some checks failed"
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
               "Value": "OK",
               "Summary": "memcheck ok",
               "Message": "memory under 20%"
             },
           ]
         }
       ]
     }
   ]
 }
```

###### Example Request with filter by service type:
URL:
```
/api/v2/metric_result/www.example.com?exec_time=2020-03-03T00:00:00Z&service=squid
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
```
{
   "root": [
     {
       "Name": "www.example.com",
       "Metrics": [
         {
           "Name": "squid_check",
           "Service": "squid",
           "Details": [
             {
               "Timestamp": "2020-03-03T12:00:00Z",
               "Value": "OK",
               "Summary": "squid is ok",
               "Message": "all checks ok"
             },
              {
               "Timestamp": "2020-03-03T18:00:00Z",
               "Value": "CRITICAL",
               "Summary": "squid is critical",
               "Message": "some checks failed"
             },
              {
               "Timestamp": "2020-03-03T23:00:00Z",
               "Value": "OK",
               "Summary": "squid is ok",
               "Message": "all checks ok"
             }
           ]
         }
       ]
     }
   ]
 }
```

###### Example Request with filter parameter set to `non-ok`:
URL:
```
/api/v2/metric_result/www.example.com?exec_time=2015-06-20T00:00:00Z&filter=non-ok
```
Headers:
```
x-api-key: shared_key_value
Accept: application/json or application/xml

```
###### Example Response using filter parameter set to `non-ok`:

Code:
```
Status: 200 OK
```
Response body:
```
{
   "root": [
     {
       "Name": "www.example.com",
       "Metrics": [
         {
           "Name": "httpd_check",
           "Service": "httpd",
           "Details": [
              {
               "Timestamp": "2015-06-20T18:00:00Z",
               "Value": "CRITICAL",
               "Summary": "httpd is critical",
               "Message": "some checks failed"
              }
           ]
         }
       ]
     }
   ]
 }
```

###### Example Request with filter parameter set to `ok`:
URL:
```
/api/v2/metric_result/www.example.com?exec_time=2015-06-20T00:00:00Z&filter=ok
```
Headers:
```
x-api-key: shared_key_value
Accept: application/json or application/xml

```
###### Example Response using filter parameter set to `ok`:

Code:
```
Status: 200 OK
```
Response body:
```
{
   "root": [
     {
       "Name": "www.example.com",
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
               "Value": "OK",
               "Summary": "memcheck ok",
               "Message": "memory under 20%"
             },
           ]
         }
       ]
     }
   ]
 }
```

### Extra endpoint information on metric results

Some metric results have additional information regarding the specific service endpoint such as it's Url, certificate DN etc... If this information is available it will be displayed under each service endpoint along with status results. For example:



```
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
               "Value": "OK",
               "Summary": "memcheck ok",
               "Message": "memory under 20%"
             },
           ]
         }
       ]
     }
   ]
 }
```