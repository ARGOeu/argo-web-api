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
 <root>
   <host name="www.example.com">
     <metric name="httpd_check">
       <status timestamp="2015-06-20T12:00:00Z" value="CRITICAL">
	 <summary>httpd status is CRITICAL</summary>
	 <message>httpd service seems down. Failed to connect to port 80.</message>
       </status>
     </metric>
   </host>
 </root>
```

