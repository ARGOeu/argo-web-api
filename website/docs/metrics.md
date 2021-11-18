---
id: metric_results
title: Metric Resulst
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

___Notes___:
`exec_time` : The execution date of query in zulu format. In order to get the correct execution time get status results for all metrics (under a given endpoint, service and endpoint group). ( GET /status/{report_name}/{lgroup_type}/{lgroup_name}/services/{service_name}/endpoints/{endpoint_name}/metrics List)

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
Reponse body:
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

#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`exec_time`| The execution date of query in zulu format - timepart is irrelevant (can be 00:00:00Z) | YES |  |
|`filter`| Filter metric results by statuses: non-ok, ok, critical, warning | NO |  |

___Notes___:
`exec_time` : The specific date of query in zulu format. The time part of the date is irrelevant because all metrics of that day are returned. ( GET /status/{report_name}/{lgroup_type}/{lgroup_name}/services/{service_name}/endpoints/{endpoint_name}/metrics List)

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
Reponse body:
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
###### Example Response using fitler parameter set to `non-ok`:

Code:
```
Status: 200 OK
```
Reponse body:
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
###### Example Response using fitler parameter set to `ok`:

Code:
```
Status: 200 OK
```
Reponse body:
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

Some metric results have additional information regarding the specific service endpoint such as it's Url, certificat DN etc... If this information is available it will be displayed under each service endpoint along with status results. Also some metrics might have a changed status due to a defined threshold rule being applied (see more about [Threshold profiles](threshold_profiles.md)). Thus they will include additional information such as the original status value (field name: `original_status`), the threshold rule applied (field name: `threshold_rule_applied`) and the actual data (field name: `actual_data`) that the rule has been applied to. For example:



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