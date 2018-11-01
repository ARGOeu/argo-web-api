---
title: 'API documentation | ARGO'
page_title: API - Operations Profile Requests
font_title: fa fa-cogs
description: API Calls for listing existing and creating new threshold profiles
---

# Description
A Threshold profile contains a list of threshold rules. Threshold rules refer to low level monitoring numeric values
that accompany metric data and threshold limits on those values that can deem the status 'WARNING' or 'CRITICAL'

# Threshold format
Each threhsold rule is expressed as string in the following format
`{label}={value}[uom];{warning};{critical};{min};{max}`

- `label` : can contain alphanumeric characters but must always begin with a letter
- `value` : is a float or integer
- `uom`   : is a string unit of measurement (accepted values: `s`,`us`,`ms`,`B`,`KB`,`MB`,`TB`,`%`,`c`)
- `warning`: is a range defining the warning threshold limits
- `critical`: is a range defining the critical threshold limits
- `min`: is a float or integer defining the minimum value
- `max`: is a float or integer defining the maximum value

Note: min,max can be omitted.

Ranges (`{warning}` or `{critical}`) are defined in the following format:
`@{floor}:{ceiling}`
-`@`: optional - negates the range (value should belong outside ranges limits)
- `floor`: integer/float or `~` that defines negative infinity
- `ceiling`: integer/float or empty (defining positive infinity)

# Thresholds rule
Each thresholds rule can contain multiple threshold definitions separated by space
for e.g.
`label01=1s;0:10;11:12 label02=1B;0:200;199:500;0;500`

# Thresholds profile
Each threhsolds profile has a name and contains a list of thresholds rules in the following json format
Each rule must always refer to a metric. It can optionally refer to a host and an endpoint group.
```
{
  "id": "68dbd3d8-c95d-41ea-b13e-7ea3644285e5",
  "name": "example-threshold-profile-101"
  "rules":[
    {
      "metric": "httpd.ResponseTime"
      "thresholds": "response=20ms;0:300;299:1000"
    },
    {
      "host": "webserver01.example.foo"
      "metric": "httpd.ResponseTime"
      "thresholds": "response=20ms;0:200;199:300"
    },
    {
      "endpoint_group": "TEST-SITE-51"
      "metric": "httpd.ResponseTime"
      "thresholds": "response=20ms;0:500;499:1000"
    }
  ]
}
```


# API Calls

Name                                     | Description                                                                            | Shortcut
---------------------------------------- | -------------------------------------------------------------------------------------- | ------------------
GET: List Thresholds Profile Requests         | This method can be used to retrieve a list of current Thresholds profiles.          | [ Description](#1)
GET: List a specific  Threshold profile         | This method can be used to retrieve a specific  Threshold profile based on its id.          | [ Description](#2)
POST: Create a new  Threshold profile  | This method can be used to create a new  Threshold profile | [ Description](#3)
PUT: Update an Threshold profile |This method can be used to update information on an existing  Threshold profile | [ Description](#4)
DELETE: Delete an  Threshold profile |This method can be used to delete an existing Threshold profile | [ Description](#5)

<a id='1'></a>

## [GET]: List Threshold Profiles

This method can be used to retrieve a list of current Threshold profiles

### Input

```
GET /thresholds_profiles
```

#### Optional Query Parameters

Type            | Description                                                                                     | Required
--------------- | ----------------------------------------------------------------------------------------------- | --------
`name`          | thresholds profile name to be used as query                                                     | NO      

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

### Response
Headers: `Status: 200 OK`

#### Response body
Json Response

```json
{
"status": {
"message": "Success",
"code": "200"
},
"data": [
{
 "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
 "name": "thr01",
 "rules": [
  {
   "host": "hostFoo",
   "metric": "metricA",
   "thresholds": "freshnesss=1s;10;9:;0;25 entries=1;3;2:0;10"
  }
 ]
},
{
 "id": "6ac7d222-1f8e-4a02-a502-720e8f11e50b",
 "name": "thr02",
 "rules": [
  {
   "host": "hostFoo",
   "metric": "metricA",
   "thresholds": "freshness=1s;10;9:;0;25 entries=1;3;2:0;10"
  }
 ]
},
{
 "id": "6ac7d555-1f8e-4a02-a502-720e8f11e50b",
 "name": "thr03",
 "rules": [
  {
   "host": "hostFoo",
   "metric": "metricA",
   "thresholds": "freshness=1s;10;9:;0;25 entries=1;3;2:0;10"
  }
 ]
}
]
}
```

<a id='2'></a>

## [GET]: List A Specific Thresholds profile
This method can be used to retrieve specific Thresholds profile based on its id

### Input

```
GET /thresholds_profiles/{ID}
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

### Response
Headers: `Status: 200 OK`

#### Response body
Json Response

```json
{
"status": {
"message": "Success",
"code": "200"
},
"data": [
{
 "id": "6ac7d222-1f8e-4a02-a502-720e8f11e50b",
 "name": "thr02",
 "rules": [
  {
   "host": "hostFoo",
   "metric": "metricA",
   "thresholds": "freshness=1s;10;9:;0;25 entries=1;3;2:0;10"
  }
 ]
}
]
}
```

<a id='3'></a>

## [POST]: Create a new Thresholds Profile
This method can be used to insert a new thresholds profile

### Input

```
POST /thresholds_profiles
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### POST BODY

```json
{
"name" : "thr04",
"rules": [
  {
    "metric": "metricB",
    "thresholds": "time=1s;10;9:30;0;30"
  }
]
}
```

### Response
Headers: `Status: 201 Created`

#### Response body
Json Response

```json
{
"status": {
"message": "Thresholds Profile successfully created",
"code": "201"
},
"data": {
"id": "{{ID}}",
"links": {
 "self": "https:///api/v2/thresholds_profiles/{{ID}}"
}
}
}
```

<a id='4'></a>

## [PUT]: Update information on an existing operations profile
This method can be used to update information on an existing operations profile

### Input

```
PUT /thresholds_profiles/{ID}
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### PUT BODY

```json
{
"name" : "thr04",
"rules": [
  {
    "metric": "metricB",
    "thresholds": "time=1s;10;9:30;0;30"
  }
]
}
```

### Response
Headers: `Status: 200 OK`

#### Response body
Json Response

```json
{
 "status": {
  "message": "Thresholds Profile successfully updated",
  "code": "200"
 }
}
```

<a id='5'></a>

## [DELETE]: Delete an existing aggregation profile
This method can be used to delete an existing aggregation profile

### Input

```
DELETE /thresholds_profiles/{ID}
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```


### Response
Headers: `Status: 200 OK`

#### Response body
Json Response

```json
{
 "status": {
  "message": "Operations Profile Successfully Deleted",
  "code": "200"
 }
}
```

## Validation Checks
When submitting or updating a new threshold profile, validation checks are performed on json POST/PUT body for the following cases:
  - Check if each thresholds rule is valid according to threshold specification discussed in the first paragraph

When an invalid operations profile is submitted the api responds with a validation error list:

#### Example invalid profile

```json
{
    "name":"test-invalid-01",
    "rules":[
      {"thresholds":"bad01=33;33s"},
      {"thresholds":"good01=33s;33 good02=1s;~:10;9:;-20;30"},
      {"thresholds":"bad02=33sbad03=1s;~~:10;9:;-20;30"},
      {"thresholds":"33;33 bad04=33s;33 -20;30"},
      {"thresholds":"good01=2KB;0:3;2:10;0;20 good02=1c;~:10;9:30;-30;30"}
    ]
}
  ```

  Api response is the following:

### Response
Headers: `Status: 422 Unprocessable Entity`

#### Response body

 ```json
 {
  "status": {
   "message": "Validation Error",
   "code": "422"
  },
  "errors": [
   {
    "message": "Validation Failed",
    "code": "422",
    "details": "Invalid threshold: bad01=33;33s"
   },
   {
    "message": "Validation Failed",
    "code": "422",
    "details": "Invalid threshold: bad02=33sbad03=1s;~~:10;9:;-20;30"
   },
   {
    "message": "Validation Failed",
    "code": "422",
    "details": "Invalid threshold: 33;33 bad04=33s;33 -20;30"
   }
  ]
 }
 ```
