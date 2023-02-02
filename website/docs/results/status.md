---
id: status_results
title: Status Results
sidebar_position: 2
---

## API calls for retrieving monitoring status results

| Name  | Description | Shortcut |
|-------|-------------|----------|
| GET: List Service Metric Status Timelines | This method may be used to retrieve a specific service metric status timeline (applies on a specific service endpoint).|<a href="#1">Description</a>|
| GET: List Service Endpoint Status Timelines | This method may be used to retrieve a specific service endpoint status timeline (for a specific group and/or for a specific service type inside that group). | <a href="#2">Description</a>|
| GET: List Service  Status Timelines |This method may be used to retrieve a specific service type status timeline (applies for a specific service endpoint group). | <a href="#3">Description</a>|
| GET: List Endpoint Group Status Timelines| This method may be used to retrieve endpoint group status timelines. | <a href="#4">Description</a>|
| GET: Metric Result | This method may be used to retrieve a specific and detailed metric result. | <a href="#5">Description</a>|
| GET: Flat list of Endpoint timelines | This method may be used to retrieve a flat list of all available endpoint status results | <a href="#6">Description</a>|
| GET: Flat list of Metric timelines for a specific metric | This method may be used to retrieve a flat list of all available endpoint metric status timelines filtered by a specific metric | <a href="#7">Description</a>|

<a id="1"></a>

## [GET]: List Service Metric Status Timelines

This method may be used to retrieve a specific service metric status timeline (applies on a specific host endpoint and a specific service flavor).

### Input
##### List All metrics:
```
/status/{report}/{group_type}/{group_name}/services/{service_type}/endpoints/{hostname}/metrics?[start_time]&[end_time]
```
##### List a specific metric:
```
/status/{report}/{endpoint_group_type}/{endpoint_group_name}/services/{service_type}/endpoints/{hostname}/metrics/{metric_name}?[start_time]&[end_time]
```

#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`report`| name of the report used | YES | |
|`group_type`| type of endpoint group| YES |  |
|`group_name`| name of endpoint group| YES |  |
|`service_type`| type of endpoint group| YES |  |
|`hostname`| hostname of service endpoint| YES |  |
|`metric_name`| name of the metric| NO |  |

#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`start_time`| UTC time in W3C format| YES |  |
|`end_time`| UTC time in W3C format| YES |  |

___Notes___:
`group_type` and `group_name` in the specific request refer always to endpoint groups (e.g.. `SITES`).
when `metric_name` is supplied, the request returns results for a specific metric. Else returns results for all available metrics for the specific __endpoint__ (and __report__)

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

##### List All metrics:

###### Example Request:
URL:
```
/status/EGI_CRITICAL/SITES/HG-03-AUTH/services/CREAM-CE/endpoints/cream01.afroditi.gr/metrics?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:59:59Z
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
Response body (XML):

```
<root>
	<group name="HG-03-AUTH" type="SITES">
		<group name="CREAM-CE" type="service">
			<endpoint name="cream01.afroditi.gr">
				<metric name="emi.cream.CREAMCE-JobSubmit">
					<status timestamp="2015-04-30T23:59:00Z" status="OK"></status>
					<status timestamp="2015-05-01T01:00:00Z" status="CRITICAL"></status>
					<status timestamp="2015-05-01T02:00:00Z" status="OK"></status>
					<status timestamp="2015-05-01T23:59:59Z" status="OK"></status>
				</metric>
				<metric name="emi.wn.WN-Bi">
					<status timestamp="2015-04-30T22:59:00Z" status="OK"></status>
					<status timestamp="2015-05-01T02:00:00Z" status="OK"></status>
					<status timestamp="2015-05-01T03:00:00Z" status="OK"></status>
					<status timestamp="2015-05-01T23:59:59Z" status="OK"></status>
				</metric>
			</endpoint>
		</group>
	</group>
</root>
```

Response body (JSON):
```
{
  "groups": [
    {
      "name": "HG-03-AUTH",
      "type": "SITES",
      "services": [
        {
          "name": "CREAM-CE",
          "type": "service",
          "endpoints": [
            {
              "name": "cream01.afroditi.gr",
              "metrics": [
                {
                  "name": "emi.cream.CREAMCE-JobSubmit",
                  "statuses": [
                    {
                      "timestamp": "2015-04-30T23:59:00Z",
                      "value": "OK"
                    },
                    {
                      "timestamp": "2015-05-01T01:00:00Z",
                      "value": "CRITICAL"
                    },
                    {
                      "timestamp": "2015-05-01T02:00:00Z",
                      "value": "OK"
                    },
                    {
                      "timestamp": "2015-05-01T23:59:59Z",
                      "value": "OK"
                    }
                  ]
                },
                {
                  "name": "emi.wn.WN-Bi",
                  "statuses": [
                    {
                      "timestamp": "2015-04-30T22:59:00Z",
                      "value": "OK"
                    },
                    {
                      "timestamp": "2015-05-02T00:00:00Z",
                      "value": "OK"
                    },
                    {
                      "timestamp": "2015-05-03T01:00:00Z",
                      "value": "OK"
                    },
                    {
                      "timestamp": "2015-05-01T23:59:59Z",
                      "value": "OK"
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

#### List specific metric
(`metric_name=emi.cream.CREAM-CE-JobSubmit`):

##### Example Request:
URL:
```
/status/EGI_CRITICAL/SITES/HG-03-AUTH/services/CREAM-CE/endpoints/cream01.afroditi.gr/metrics/emi.cream.CREAMCE-JobSubmit?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:59:59Z
```
Headers:
```
x-api-key:"INSERTTENANTKEYHERE"
Accept:"application/xml"

```
##### Example Response:
Code:
```
Status: 200 OK
```
Response body (XML):

```
<root>
	<group name="HG-03-AUTH" type="SITES">
		<group name="CREAM-CE" type="service">
			<endpoint name="cream01.afroditi.gr">
				<metric name="emi.cream.CREAMCE-JobSubmit">
					<status timestamp="2015-04-30T23:59:00Z" status="OK"></status>
					<status timestamp="2015-05-01T01:00:00Z" status="CRITICAL"></status>
					<status timestamp="2015-05-01T02:00:00Z" status="OK"></status>
					<status timestamp="2015-05-01T23:59:59Z" status="OK"></status>
				</metric>
			</endpoint>
		</group>
	</group>
</root>
```

Response body (JSON):
```
{
  "groups": [
    {
      "name": "HG-03-AUTH",
      "type": "SITES",
      "services": [
        {
          "name": "CREAM-CE",
          "type": "service",
          "endpoints": [
            {
              "name": "cream01.afroditi.gr",
              "metrics": [
                {
                  "name": "emi.cream.CREAMCE-JobSubmit",
                  "statuses": [
                    {
                      "timestamp": "2015-04-30T23:59:00Z",
                      "value": "OK"
                    },
                    {
                      "timestamp": "2015-05-01T01:00:00Z",
                      "value": "CRITICAL"
                    },
                    {
                      "timestamp": "2015-05-02T01:00:00Z",
                      "value": "OK"
                    },
                    {
                      "timestamp": "2015-05-01T23:59:59Z",
                      "value": "OK"
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

### Extra endpoint information on metric status results

Some service metric status results have additional information regarding the specific service endpoint such as it's Url, certificate DN etc... If this information is available it will be displayed under each service endpoint along with status results. For example:



```
{
  "groups": [
    {
      "name": "HG-03-AUTH",
      "type": "SITES",
      "services": [
        {
          "name": "CREAM-CE",
          "type": "service",
          "endpoints": [
            {
              "name": "cream01.afroditi.gr",
              "info": {
                  "Url": "https://cream01.afroditi.gr/path/to/service"
               },
              "metrics": [
                {
                  "name": "emi.cream.CREAMCE-JobSubmit",
                  "statuses": [
                    {
                      "timestamp": "2015-04-30T23:59:00Z",
                      "value": "OK"
                    },
                    {
                      "timestamp": "2015-05-01T01:00:00Z",
                      "value": "CRITICAL"
                    },
                    {
                      "timestamp": "2015-05-02T01:00:00Z",
                      "value": "OK"
                    },
                    {
                      "timestamp": "2015-05-01T23:59:59Z",
                      "value": "OK"
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

### Threshold rule information in status metric timelines

By using the url parameter `view=details` the argo-web-api will enrich the status timeline results with additional information in case a threshold rule has been applied to the results. For example:

```json
{
  "groups": [
    {
      "name": "HG-03-AUTH",
      "type": "SITES",
      "services": [
        {
          "name": "CREAM-CE",
          "type": "service",
          "endpoints": [
            {
              "name": "cream01.afroditi.gr",
              "info": {
                  "Url": "https://cream01.afroditi.gr/path/to/service"
               },
              "metrics": [
                {
                  "name": "emi.cream.CREAMCE-JobSubmit",
                  "statuses": [
                    {
                      "timestamp": "2015-04-30T23:59:00Z",
                      "value": "OK"
                    },
                    {
                      "timestamp": "2015-05-01T01:00:00Z",
                      "value": "CRITICAL",
                      "actual_data": "latency=15s",
                      "threshold_rule_applied": "latency=15s",
                      "original_status": "OK"
                    },
                    {
                      "timestamp": "2015-05-02T01:00:00Z",
                      "value": "OK"
                    },
                    {
                      "timestamp": "2015-05-01T23:59:59Z",
                      "value": "OK"
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```


<a id="2"></a>


## [GET]: List Service Endpoint Status Timelines

This method may be used to retrieve a specific service endpoint status timeline (applies on a specific service type).

### Input
##### List All endpoints:
```
/status/{report}/{group_type}/{group_name}/services/{service_type}/endpoints?[start_time]&[end_time]
```
##### List a specific endpoint:
```
/status/{report}/{endpoint_group_type}/{endpoint_group_name}/services/{service_type}/endpoints/{hostname}?[start_time]&[end_time]
```
##### List All endpoints in a specific group:
```
/status/{report}/{group_type}/{group_name}/endpoints?[start_time]&[end_time]
```

#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`report`| name of the report used | YES | |
|`group_type`| type of endpoint group| YES |  |
|`group_name`| name of endpoint group| YES |  |
|`service_type`| type of endpoint group| YES |  |
|`hostname`| hostname of service endpoint| NO |  |

#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`start_time`| UTC time in W3C format| YES |  |
|`end_time`| UTC time in W3C format| YES |  |

___Notes___:
`group_type` and `group_name` in the specific request refer always to endpoint groups (e.g.. `SITES`).
when `hostname` is supplied, the request returns results for a specific endpoint. Else returns results for all available metrics for the specific __endpoint__ (and __report__)

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

##### List All metrics:

###### Example Request:
URL:
```
/status/EGI_CRITICAL/SITES/HG-03-AUTH/services/CREAM-CE/endpoints?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:59:59Z
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
Response body (XML):
```
<root>
	<group name="HG-03-AUTH" type="SITES">
		<group name="CREAM-CE" type="service">
			<endpoint name="cream01.afroditi.gr">
				<status timestamp="2015-04-30T23:59:00Z" status="OK"></status>
				<status timestamp="2015-05-01T01:00:00Z" status="CRITICAL"></status>
				<status timestamp="2015-05-01T02:00:00Z" status="OK"></status>
				<status timestamp="2015-05-01T23:59:59Z" status="OK"></status>
			</endpoint>
			<endpoint name="cream02.afroditi.gr">
				<status timestamp="2015-04-30T23:59:00Z" status="OK"></status>
				<status timestamp="2015-05-01T01:00:00Z" status="CRITICAL"></status>
				<status timestamp="2015-05-01T02:00:00Z" status="OK"></status>
				<status timestamp="2015-05-01T23:59:59Z" status="OK"></status>
			</endpoint>
		</group>
	</group>
</root>
```

Response body (JSON):
```
{
  "groups": [
    {
      "name": "HG-03-AUTH",
      "type": "SITES",
      "services": [
        {
          "name": "CREAM-CE",
          "type": "service",
          "endpoints": [
            {
              "name": "cream01.afroditi.gr",
              "statuses": [
                {
                  "timestamp": "2015-04-30T23:59:00Z",
                  "value": "OK"
                },
                {
                  "timestamp": "2015-05-01T01:00:00Z",
                  "value": "CRITICAL"
                },
                {
                  "timestamp": "2015-05-01T02:00:00Z",
                  "value": "OK"
                },
                {
                  "timestamp": "2015-05-01T23:59:59Z",
                  "value": "OK"
                }
              ]
            },
                        {
              "name": "cream02.afroditi.gr",
              "statuses": [
                {
                  "timestamp": "2015-04-30T23:59:00Z",
                  "value": "OK"
                },
                {
                  "timestamp": "2015-05-01T01:00:00Z",
                  "value": "CRITICAL"
                },
                {
                  "timestamp": "2015-05-01T02:00:00Z",
                  "value": "OK"
                },
                {
                  "timestamp": "2015-05-01T23:59:59Z",
                  "value": "OK"
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```


#### List specific endpoint
(`hostname=cream01.afroditi.gr`):

##### Example Request:
URL:
```
/status/EGI_CRITICAL/SITES/HG-03-AUTH/services/CREAM-CE/endpoints/cream01.afroditi.gr?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:59:59Z
```
Headers:
```
x-api-key: shared_key_value
Accept: application/json or application/xml

```
##### Example Response:
Code:
```
Status: 200 OK
```
Response body (XML):

```
<root>
	<group name="HG-03-AUTH" type="SITES">
		<group name="CREAM-CE" type="service">
			<endpoint name="cream01.afroditi.gr">
				<status timestamp="2015-04-30T23:59:00Z" status="OK"></status>
				<status timestamp="2015-05-01T01:00:00Z" status="CRITICAL"></status>
				<status timestamp="2015-05-01T02:00:00Z" status="OK"></status>
				<status timestamp="2015-05-01T23:59:59Z" status="OK"></status>
			</endpoint>
		</group>
	</group>
</root>
```

Response body (JSON):
```
{
  "groups": [
    {
      "name": "HG-03-AUTH",
      "type": "SITES",
      "services": [
        {
          "name": "CREAM-CE",
          "type": "service",
          "endpoints": [
            {
              "name": "cream01.afroditi.gr",
              "statuses": [
                {
                  "timestamp": "2015-04-30T23:59:00Z",
                  "value": "OK"
                },
                {
                  "timestamp": "2015-05-01T01:00:00Z",
                  "value": "CRITICAL"
                },
                {
                  "timestamp": "2015-05-01T02:00:00Z",
                  "value": "OK"
                },
                {
                  "timestamp": "2015-05-01T23:59:59Z",
                  "value": "OK"
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

#### List all endpoints included in a specific group
For a specific group, return all endpoints included - grouped by service type

##### Example Request:
URL:
```
/status/EGI_CRITICAL/SITES/HG-03-AUTH/endpoints?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:59:59Z
```
Headers:
```
x-api-key: shared_key_value
Accept: application/json or application/xml

```
##### Example Response:
Code:
```
Status: 200 OK
```

Response body (JSON):
```json
{
 "groups": [
  {
   "name": "HG-03-AUTH",
   "type": "SITES",
   "services": [
    {
     "name": "CREAM-CE",
     "type": "service",
     "endpoints": [
      {
       "name": "cream01.afroditi.gr",
       "info": {
        "Url": "http://example.foo/path/to/service"
       },
       "statuses": [
        {
         "timestamp": "2015-05-01T00:00:00Z",
         "value": "OK"
        },
        {
         "timestamp": "2015-05-01T01:00:00Z",
         "value": "CRITICAL"
        },
        {
         "timestamp": "2015-05-01T05:00:00Z",
         "value": "OK"
        },
        {
         "timestamp": "2015-05-01T23:59:59Z",
         "value": "OK"
        }
       ]
      },
      {
       "name": "cream02.afroditi.gr",
       "statuses": [
        {
         "timestamp": "2015-05-01T00:00:00Z",
         "value": "OK"
        },
        {
         "timestamp": "2015-05-01T08:47:00Z",
         "value": "WARNING"
        },
        {
         "timestamp": "2015-05-01T12:00:00Z",
         "value": "OK"
        },
        {
         "timestamp": "2015-05-01T23:59:59Z",
         "value": "OK"
        }
       ]
      },
      {
       "name": "cream03.afroditi.gr",
       "statuses": [
        {
         "timestamp": "2015-05-01T00:00:00Z",
         "value": "OK"
        },
        {
         "timestamp": "2015-05-01T04:40:00Z",
         "value": "UNKNOWN"
        },
        {
         "timestamp": "2015-05-01T06:00:00Z",
         "value": "CRITICAL"
        },
        {
         "timestamp": "2015-05-01T23:59:59Z",
         "value": "CRITICAL"
        }
       ]
      }
     ]
    },
    {
     "name": "STORAGE",
     "type": "service",
     "endpoints": [
      {
       "name": "storage.example.foo",
       "statuses": [
        {
         "timestamp": "2015-05-01T00:00:00Z",
         "value": "OK"
        },
        {
         "timestamp": "2015-05-01T06:40:00Z",
         "value": "WARNING"
        },
        {
         "timestamp": "2015-05-01T09:00:00Z",
         "value": "CRITICAL"
        },
        {
         "timestamp": "2015-05-01T23:59:59Z",
         "value": "CRITICAL"
        }
       ]
      }
     ]
    }
   ]
  }
 ]
}
```

### Extra information for a specific endpoint on endpoint status results

Some service endpoint status results have additional information regarding the specific service endpoint such as it's Url, certificate DN etc... If this information is available it will be displayed under each service endpoint along with status results. For example:



```
{
  "groups": [
    {
      "name": "HG-03-AUTH",
      "type": "SITES",
      "services": [
        {
          "name": "CREAM-CE",
          "type": "service",
          "endpoints": [
            {
              "name": "cream01.afroditi.gr",
              "info": {
                  "Url": "https://cream01.afroditi.gr/path/to/service"
               },
              "statuses": [
                {
                  "timestamp": "2015-04-30T23:59:00Z",
                  "value": "OK"
                },
                {
                  "timestamp": "2015-05-01T01:00:00Z",
                  "value": "CRITICAL"
                },
                {
                  "timestamp": "2015-05-01T02:00:00Z",
                  "value": "OK"
                },
                {
                  "timestamp": "2015-05-01T23:59:59Z",
                  "value": "OK"
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

### Threshold rule information in status endpoint timelines

By using the url parameter `view=details` the argo-web-api will enrich the status timeline results with additional information in case a threshold rule has been applied to the results. For example:

```
{
  "groups": [
    {
      "name": "HG-03-AUTH",
      "type": "SITES",
      "services": [
        {
          "name": "CREAM-CE",
          "type": "service",
          "endpoints": [
            {
              "name": "cream01.afroditi.gr",
              "info": {
                  "Url": "https://cream01.afroditi.gr/path/to/service"
               },
              "statuses": [
                {
                  "timestamp": "2015-04-30T23:59:00Z",
                  "value": "OK"
                },
                {
                  "timestamp": "2015-05-01T01:00:00Z",
                  "value": "CRITICAL",
                  "affected_by_threshold_rule": true
                },
                {
                  "timestamp": "2015-05-01T02:00:00Z",
                  "value": "OK"
                },
                {
                  "timestamp": "2015-05-01T23:59:59Z",
                  "value": "OK"
                }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

<a id="3"></a>


## [GET]: List Service Status Timelines

This method may be used to retrieve a specific service flavor status timeline (applies for a specific service endpoint group).

### Input
##### List All service types:
```
/status/{report}/{group_type}/{group_name}/services[start_time]&[end_time]
```
##### List a specific service type:
```
/status/{report}/{group_type}/{group_name}/services/{service_type}[start_time]&[end_time]
```
#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`report`| name of the report used | YES | |
|`group_type`| type of endpoint group| YES |  |
|`group_name`| name of endpoint group| YES |  |
|`service_type`| type of endpoint group| NO |  |

#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`start_time`| UTC time in W3C format| YES |  |
|`end_time`| UTC time in W3C format| YES |  |

___Notes___:
`group_type` and `group_name` in the specific request refer always to endpoint groups (e.g.. `SITES`).
when `service_name` is supplied, the request returns results for a specific service type. Else returns results for all available service types for the specific __endpoint_group__ (and __report__)

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

##### List All service types:

###### Example Request:
URL:
```
/status/EGI_CRITICAL/SITES/HG-03-AUTH/services?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:59:59Z
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
Response body (XML):
```
<root>
	<group name="HG-03-AUTH" type="SITES">
		<group name="CREAM-CE" type="service">
			<status timestamp="2015-04-30T23:59:00Z" status="OK"></status>
			<status timestamp="2015-05-01T01:00:00Z" status="CRITICAL"></status>
			<status timestamp="2015-05-01T02:00:00Z" status="OK"></status>
			<status timestamp="2015-05-01T23:59:59Z" status="OK"></status>
		</group>
		<group name="SRMv2" type="service">
			<status timestamp="2015-04-30T23:59:00Z" status="OK"></status>
			<status timestamp="2015-05-01T01:00:00Z" status="CRITICAL"></status>
			<status timestamp="2015-05-01T02:00:00Z" status="OK"></status>
			<status timestamp="2015-05-01T23:59:59Z" status="OK"></status>
		</group>
	</group>
</root>
```

Response body (JSON):
```
{
  "groups": [
    {
      "name": "HG-03-AUTH",
      "type": "SITES",
      "services": [
        {
          "name": "CREAM-CE",
          "type": "service",
          "statuses": [
            {
              "timestamp": "2015-04-30T23:59:00Z",
              "value": "OK"
            },
            {
              "timestamp": "2015-05-01T01:00:00Z",
              "value": "CRITICAL"
            },
            {
              "timestamp": "2015-05-01T02:00:00Z",
              "value": "OK"
            },
            {
              "timestamp": "2015-05-01T23:59:59Z",
              "value": "OK"
            }
          ]
        },
        {
          "name": "SRMv2",
          "type": "service",
          "statuses": [
            {
              "timestamp": "2015-04-30T23:59:00Z",
              "value": "OK"
            },
            {
              "timestamp": "2015-05-01T01:00:00Z",
              "value": "CRITICAL"
            },
            {
              "timestamp": "2015-05-01T02:00:00Z",
              "value": "OK"
            },
            {
              "timestamp": "2015-05-01T23:59:59Z",
              "value": "OK"
            }
          ]
        }
      ]
    }
  ]
}
```


##### List specific service type
(`service_type=CREAM-CE`):

###### Example Request:
URL:
```
/status/EGI_CRITICAL/SITES/HG-03-AUTH/services/CREAM-CE?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:59:59Z
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
Response body (XML):

```
<root>
	<group name="HG-03-AUTH" type="SITES">
		<group name="CREAM-CE" type="service">
			<status timestamp="2015-04-30T23:59:00Z" status="OK"></status>
			<status timestamp="2015-05-01T01:00:00Z" status="CRITICAL"></status>
			<status timestamp="2015-05-01T02:00:00Z" status="OK"></status>
			<status timestamp="2015-05-01T23:59:59Z" status="OK"></status>
		</group>
	</group>
</root>
```

Response body (JSON):
```
{
  "groups": [
    {
      "name": "HG-03-AUTH",
      "type": "SITES",
      "services": [
        {
          "name": "CREAM-CE",
          "type": "service",
          "statuses": [
            {
              "timestamp": "2015-04-30T23:59:00Z",
              "value": "OK"
            },
            {
              "timestamp": "2015-05-01T01:00:00Z",
              "value": "CRITICAL"
            },
            {
              "timestamp": "2015-05-01T02:00:00Z",
              "value": "OK"
            },
            {
              "timestamp": "2015-05-01T23:59:59Z",
              "value": "OK"
            }
          ]
        }
      ]
    }
  ]
}
```

### Threshold rule information in status service timelines

By using the url parameter `view=details` the argo-web-api will enrich the status timeline results with additional information in case a threshold rule has been applied to the results. For example:


```
{
  "groups": [
    {
      "name": "HG-03-AUTH",
      "type": "SITES",
      "services": [
        {
          "name": "CREAM-CE",
          "type": "service",
          "statuses": [
            {
              "timestamp": "2015-04-30T23:59:00Z",
              "value": "OK"
            },
            {
              "timestamp": "2015-05-01T01:00:00Z",
              "value": "CRITICAL",
              "affected_by_threshold_rule": true
            },
            {
              "timestamp": "2015-05-01T02:00:00Z",
              "value": "OK"
            },
            {
              "timestamp": "2015-05-01T23:59:59Z",
              "value": "OK"
            }
          ]
        }
      ]
    }
  ]
}
```


<a id="4"></a>




## [GET]: List Endpoint Group Status Timelines

This method may be used to retrieve status timelines for endpoint groups.

### Input
##### List All endpoint groups of specific type:
```
/status/{report}/{group_type}[start_time]&[end_time]
```
##### List a specific endpoint group of specific type:
```
/status/{report}/{group_type}/{group_name}[start_time]&[end_time]
```
#### Path Parameters
| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`report`| name of the report used | YES | |
|`group_type`| type of endpoint group| YES |  |
|`group_name`| name of endpoint group| NO |  |

#### Url Parameters

| Type | Description | Required | Default value |
|------|-------------|----------|---------------|
|`start_time`| UTC time in W3C format| YES |  |
|`end_time`| UTC time in W3C format| YES |  |

___Notes___:
`group_type` and `group_name` in the specific request refer always to endpoint groups (e.g.. `SITES`).
when `group_name` is supplied, the request returns results for a specific endpoint group. Else returns results for all available endpoint groups of the specific __group_type__

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

##### List All endpoint groups:

###### Example Request:
URL:
```
/status/EGI_CRITICAL/SITES?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:59:59Z
```
Headers:
```
x-api-key:"INSERTTENANTKEYHERE"
Accept:"application/xml"
```
###### Example Response:
Code:
```
Status: 200 OK
```
Response body (XML):
```
<root>
	<group name="HG-03-AUTH" type="SITES">
		<status timestamp="2015-05-01T00:00:00Z" status="CRITICAL"></status>
		<status timestamp="2015-05-01T01:00:00Z" status="WARNING"></status>
		<status timestamp="2015-05-01T05:00:00Z" status="OK"></status>
		<status timestamp="2015-05-01T23:59:59Z" status="OK"></status>
	</group>
	<group name="HG-01-AUTH" type="SITES">
		<status timestamp="2015-05-01T00:00:00Z" status="CRITICAL"></status>
		<status timestamp="2015-05-01T02:00:00Z" status="UNKNOWN"></status>
		<status timestamp="2015-05-01T05:00:00Z" status="OK"></status>
		<status timestamp="2015-05-01T23:59:59Z" status="OK"></status>
	</group>
</root>
```

Response body (JSON):
```
{
  "groups": [
    {
      "name": "HG-03-AUTH",
      "type": "SITES",
      "statuses": [
        {
          "timestamp": "2015-05-01T00:00:00Z",
          "value": "CRITICAL"
        },
        {
          "timestamp": "2015-05-01T01:00:00Z",
          "value": "WARNING"
        },
        {
          "timestamp": "2015-05-01T05:00:00Z",
          "value": "OK"
        },
        {
          "timestamp": "2015-05-01T23:59:59Z",
          "value": "OK"
        }
      ]
    },
    {
      "name": "HG-01-AUTH",
      "type": "SITES",
      "statuses": [
        {
          "timestamp": "2015-05-01T00:00:00Z",
          "value": "CRITICAL"
        },
        {
          "timestamp": "2015-05-01T02:00:00Z",
          "value": "UNKNOWN"
        },
        {
          "timestamp": "2015-05-01T05:00:00Z",
          "value": "OK"
        },
        {
          "timestamp": "2015-05-01T23:59:59Z",
          "value": "OK"
        }
      ]
    }
  ]
}
```


##### List specific endpoint group
(`group_name=HG-03-AUTH`):

###### Example Request:
URL:
```
/status/EGI_CRITICAL/SITES/HG-03-AUTH?start_time=2015-05-01T00:00:00Z&end_time=2015-05-01T23:59:59Z
```
Headers:
```
x-api-key:"INSERTTENANTKEYHERE"
Accept:"application/xml"

```
###### Example Response:
Code:
```
Status: 200 OK
```
Response body (XML):

```
<root>
	<group name="HG-03-AUTH" type="SITES">
		<status timestamp="2015-05-01T00:00:00Z" status="CRITICAL"></status>
		<status timestamp="2015-05-01T02:00:00Z" status="WARNING"></status>
		<status timestamp="2015-05-01T05:00:00Z" status="OK"></status>
		<status timestamp="2015-05-01T23:59:59Z" status="OK"></status>
	</group>
</root>
```

Response body (JSON):
```
{
  "groups": [
    {
      "name": "HG-03-AUTH",
      "type": "SITES",
      "statuses": [
        {
          "timestamp": "2015-05-01T00:00:00Z",
          "value": "CRITICAL"
        },
        {
          "timestamp": "2015-05-01T02:00:00Z",
          "value": "WARNING"
        },
        {
          "timestamp": "2015-05-01T05:00:00Z",
          "value": "OK"
        },
        {
          "timestamp": "2015-05-01T23:59:59Z",
          "value": "OK"
        }
      ]
    }
  ]
}
```

### Threshold rule information in status endpoint group timelines

By using the url parameter `view=details` the argo-web-api will enrich the status timeline results with additional information in case a threshold rule has been applied to the results. For example:


```
{
  "groups": [
    {
      "name": "HG-03-AUTH",
      "type": "SITES",
      "statuses": [
        {
          "timestamp": "2015-05-01T00:00:00Z",
          "value": "CRITICAL"
        },
        {
          "timestamp": "2015-05-01T02:00:00Z",
          "value": "WARNING",
          "affected_by_threshold_rule": true
        },
        {
          "timestamp": "2015-05-01T05:00:00Z",
          "value": "OK"
        },
        {
          "timestamp": "2015-05-01T23:59:59Z",
          "value": "OK"
        }
      ]
    }
  ]
}
```

<a id="5"></a>

## [GET]: Metric Result

This method may be used to retrieve a detailed metric result.

### Input

```
/metric_result/{hostname}/{metric_name}?[exec_time]
```

#### Path Parameters

Name             | Description                                              | Required | Default value
---------------- | -------------------------------------------------------- | -------- | -------------
`{hostname}`     | Name of the endpoint                                     | YES      |
`{metric_name}`  | Name of the metric (probe) for which results are queries | YES      |

#### URL Parameters

Type            | Description             | Required | Default value
--------------- | ----------------------- | -------- | -------------
`[exec_time]`   | UTC time in W3C format  | YES      |


#### Headers

```
x-api-key: shared_key_value
Accept: application/json or application/xml
```


#### Response Code
```
Status: 200 OK
```



#### Response body
##### Example Request:
URL:
```
/api/v2/metric_result/www.example.com/httpd_check?exec_time=2015-06-20T12:00:00Z
```
Headers:
```
x-api-key:"INSERTTENANTKEYHERE"
Accept:"application/xml"
```
##### Example Response:
Code:
```
Status: 200 OK
```
Response body (XML):
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

Response body (JSON):
```
{
  "root": [
    {
      "Name": "www.example.com",
      "Metrics": [
        {
          "Name": "httpd_check",
          "Details": [
            {
              "Timestamp": "2015-06-20T12:00:00Z",
              "Value": "CRITICAL",
              "Summary": "httpd status is CRITICAL",
              "Message": "httpd service seems down. Failed to connect to port 80."
            }
          ]
        }
      ]
    }
  ]
}
```


<a id="6"></a>

# [GET]: Flat List of status timelines for all available service Endpoints

The following methods can be used to obtain a tenant's flat list of all service endpoints status timelines. The api authenticates the tenant using the api-key within the x-api-key header. Pagination is also supported by using the optional parameters `pageSize` to define the size of each result page and `nextPageToken` to proceed to the next available page of results.
## [GET] Endpoints Status timelines

### Input

Request a flat list of all endpoint status timelines

```
/status/{report_name}/endpoints?[start_time]&[end_time]&[granularity]&[pageSize]&[nextPageToken]
```


#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start_time]`  | UTC time in W3C format                                                                          | YES      |
| `[end_time]`    | UTC time in W3C format                                                                          | YES      |
| `[pageSize]` | How many results to return per request (-1 means return all results) | NO       | -1       |
| `[nextPageToken]` | token to proceed to the next page | NO       |  |

#### Path Parameters

| Name                    | Description                                                                                           | Required | Default value |
| ----------------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report_name}`         | Name of the report that contains all the information about the profile, filter tags, group types etc. | YES      |

#### Headers

```
x-api-key: "tenant_key_value"
Accept: "application/xml" or "application/json"
```

##### Response

```
Status: 200 OK
```

## Request endpoint a/r under service: `service_a`

#### URL

`/api/v2/status/{report_name}/endpoints?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:23:59Z&pageSize=2`

#### Response Body

```
{
  "results": [
    {
      "name": "e01",
      "type": "endpoint",
      "service": "service_a",
      "supergroup": "ST01",
      "results": [
        {
          "timestamp": "2015-06-22",
          "availability": "98.26389",
          "reliability": "98.26389",
          "unknown": "0",
          "uptime": "0.98264",
          "downtime": "0"
        },
        {
          "timestamp": "2015-06-23",
          "availability": "54.03509",
          "reliability": "81.48148",
          "unknown": "0.01042",
          "uptime": "0.53472",
          "downtime": "0.33333"
        }
      ]
    },
    {
      "name": "e02",
      "type": "endpoint",
      "service": "service_a",
      "supergroup": "ST01",
      "results": [
        {
          "timestamp": "2015-06-22",
          "availability": "96.875",
          "reliability": "96.875",
          "unknown": "0",
          "uptime": "0.96875",
          "downtime": "0"
        }
      ]
    }
  ],
  "pageSize": 2,
  "nextPageToken": "Mg=="
}
```

## Request to see next page of results

#### URL

`/api/v2/status/{report_name}/endpoints?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:23:59Z&granularity=daily&pageSize=2&nextPageToken=Mg==`

#### Response Body

```
{
  "results": [
    {
      "name": "e02",
      "type": "endpoint",
      "service": "service_a",
      "supergroup": "ST01",
      "results": [
        {
          "timestamp": "2015-06-22",
          "availability": "96.875",
          "reliability": "96.875",
          "unknown": "0",
          "uptime": "0.96875",
          "downtime": "0"
        }
      ]
    }
  ],
  "pageSize": 2,
}
```


<a id="7"></a>

# [GET]: Flat List of status timelines for all available endpoint metrics filter by a specific metric

The following method can be used to obtain a tenant's flat list of all available service endpoints metric status timelines filtered by a specific metric. The api authenticates the tenant using the api-key within the x-api-key header. Pagination is also supported by using the optional parameters `pageSize` to define the size of each result page and `nextPageToken` to proceed to the next available page of results.
## [GET] Endpoints Status timelines

### Input

Request a flat list of all endpoint metric status timelines filtered by a specific metric

```
/status/{report_name}/metrics/{metric_name}?[start_time]&[end_time]&[granularity]&[pageSize]&[nextPageToken]
```


#### Query Parameters

| Type            | Description                                                                                     | Required | Default value |
| --------------- | ----------------------------------------------------------------------------------------------- | -------- | ------------- |
| `[start_time]`  | UTC time in W3C format                                                                          | YES      |
| `[end_time]`    | UTC time in W3C format                                                                          | YES      |
| `[pageSize]` | How many results to return per request (-1 means return all results) | NO       | -1       |
| `[nextPageToken]` | token to proceed to the next page | NO       |  |

#### Path Parameters

| Name                    | Description                                                                                           | Required | Default value |
| ----------------------- | ----------------------------------------------------------------------------------------------------- | -------- | ------------- |
| `{report_name}`         | Name of the report that contains all the information about the profile, filter tags, group types etc. | YES      |
| `{metric_name}`         | Name of the metric to filter results by | YES      |

#### Headers

```
x-api-key: "tenant_key_value"
Accept: "application/xml" or "application/json"
```

##### Response

```
Status: 200 OK
```

## Request endpoint a/r under service: `service_a`

#### URL

`/api/v2/status/{report_name}/metrics/someService-FileTransfer?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:23:59Z&pageSize=2`

#### Response Body

```
{
  "results": [
    {
      "type":"endpoint_metric",
      "name": "someservice.example.gr",
      "service": "someService",
      "supergroup": "GROUPA",
      "metric": "someService-FileTransfer",
      "statuses": [
        {
          "timestamp": "2015-04-30T23:59:00Z",
          "value": "OK"
        },
        {
          "timestamp": "2015-05-01T00:00:00Z",
          "value": "OK"
        },
        {
          "timestamp": "2015-05-01T01:00:00Z",
          "value": "CRITICAL"
        }
      ]
    }
  ],
  "nextPageToken": "Mg==",
  "pageSize": 2
}
```

## Request to see next page of results

#### URL

`/api/v2/status/{report_name}/metrics/someService-FileTransfer?start_time=2015-06-22T00:00:00Z&end_time=2015-06-23T23:23:59Z&granularity=daily&pageSize=2&nextPageToken=Mg==`

#### Response Body

```
{
  "results": [
    {
      "type":"endpoint_metric",
      "name": "someservice.example.gr",
      "service": "someService",
      "supergroup": "GROUPA",
      "metric": "someService-FileTransfer",
      "statuses": [
        {
          "timestamp": "2015-05-01T05:00:00Z",
          "value": "OK"
        }
      ]
    },
    {
      "type":"endpoint_metric",
      "name": "someservice2.example.gr",
      "service": "someService",
      "supergroup": "GROUPA",
      "metric": "someService-FileTransfer",
      "statuses": [
        {
          "timestamp": "2015-04-30T23:59:00Z",
          "value": "OK"
        },
        {
          "timestamp": "2015-05-01T00:00:00Z",
          "value": "OK"
        }
      ]
    }
  ],
  "nextPageToken": "NA==",
  "pageSize": 2
}
```

### Threshold rule information in flat status endpoint group timelines

By using the url parameter `view=details` the argo-web-api will enrich the status timeline results with additional information in case a threshold rule has been applied to the results. For example:



```
{
  "results": [
    {
      "type":"endpoint_metric",
      "name": "someservice.example.gr",
      "service": "someService",
      "supergroup": "GROUPA",
      "metric": "someService-FileTransfer",
      "statuses": [
        {
          "timestamp": "2015-05-01T05:00:00Z",
          "value": "OK"
        }
      ]
    },
    {
      "type":"endpoint_metric",
      "name": "someservice2.example.gr",
      "service": "someService",
      "supergroup": "GROUPA",
      "metric": "someService-FileTransfer",
      "statuses": [
        {
          "timestamp": "2015-04-30T23:59:00Z",
          "value": "OK",
          "affected_by_threshold_rule": true
        },
        {
          "timestamp": "2015-05-01T00:00:00Z",
          "value": "OK"
        }
      ]
    }
  ],
  "nextPageToken": "NA==",
  "pageSize": 2
}
```


