
## GET: Metric Result

This method may be used to retrieve a detailed metric result.

### Input

```
/metric_result/{hostname}/{metric_name}?[exec_time]
```

#### Query Parameters

Type            | Description             | Required | Default value
--------------- | ----------------------- | -------- | -------------
`[exec_time]`   | UTC time in W3C format  | YES      |

#### Path Parameters

Name             | Description                                              | Required | Default value
---------------- | -------------------------------------------------------- | -------- | -------------
`{hostname}`     | Name of the endpoint                                     | YES      |
`{metric_name}`  | Name of the metric (probe) for which results are queries | YES      |


#### Headers

##### Request
```
x-api-key: "tenant_key_value"
Accept: "application/xml" or "application/json"
```

##### Response
```
Status: 200 OK
```

#### URL
`/api/v2/metric_result/www.example.com/httpd_check?exec_time=2015-06-20T12:00:00Z`


#### Response body

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

