# Metric results

API call for retrieving detailed metric result.

## [GET]: Metric Result

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

