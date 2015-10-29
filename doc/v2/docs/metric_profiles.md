---
title: 'API documentation | ARGO'
page_title: API - Metric Profile Requests
font_title: fa fa-cogs
description: API Calls for listing existing and creating new metric profiles
---

# API Calls

Name                                     | Description                                                                            | Shortcut
---------------------------------------- | -------------------------------------------------------------------------------------- | ------------------
GET: List Metric Profile Requests         | This method can be used to retrieve a list of current metric profiles.          | [ Description](#1)
GET: List a specific Metric profile         | This method can be used to retrieve a specific metric profile based on its id.          | [ Description](#2)
POST: Create a new metric profile  | This method can be used to create a new metric profile | [ Description](#3)
PUT: Update a metric profile |This method can be used to update information on an existing metric profile | [ Description](#4)
DELETE: Delete a metric profile |This method can be used to delete an existing metric profile | [ Description](#5)
<a id='1'></a>

# GET: List Metric Profiles
This method can be used to retrieve a list of current Metric profiles

## Input

```
GET /metric_profiles
```

#### Optional Query Parameters

Type            | Description                                                                                     | Required
--------------- | ----------------------------------------------------------------------------------------------- | --------
`name`  | metric profile name to be used as query                                                                          | NO      

### Request headers

```
x-api-key: shared_key_value
Content-Type: application/json
Accept: application/json
```

## Response
Headers: `Status: 200 OK`

### Response body
Json Response

```json
{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "name": "ch.cern.SAM.ROC",
   "services": [
    {
     "service": "CREAM-CE",
     "metrics": [
      "emi.cream.CREAMCE-JobSubmit",
      "emi.wn.WN-Bi",
      "emi.wn.WN-Csh",
      "hr.srce.CADist-Check",
      "hr.srce.CREAMCE-CertLifetime",
      "emi.wn.WN-SoftVer"
     ]
    },
    {
     "service": "SRMv2",
     "metrics": [
      "hr.srce.SRM2-CertLifetime",
      "org.sam.SRM-Del",
      "org.sam.SRM-Get",
      "org.sam.SRM-GetSURLs",
      "org.sam.SRM-GetTURLs",
      "org.sam.SRM-Ls",
      "org.sam.SRM-LsDir",
      "org.sam.SRM-Put"
     ]
    }
   ]
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "name": "ch.cern.SAM.ROC_CRITICAL",
   "services": [
    {
     "service": "CREAM-CE",
     "metrics": [
      "emi.cream.CREAMCE-JobSubmit",
      "emi.wn.WN-Bi",
      "emi.wn.WN-Csh",
      "emi.wn.WN-SoftVer"
     ]
    },
    {
     "service": "SRMv2",
     "metrics": [
      "hr.srce.SRM2-CertLifetime",
      "org.sam.SRM-Del",
      "org.sam.SRM-Get",
      "org.sam.SRM-GetSURLs",
      "org.sam.SRM-GetTURLs",
      "org.sam.SRM-Ls",
      "org.sam.SRM-LsDir",
      "org.sam.SRM-Put"
     ]
    }
   ]
  }
 ]
}
```

<a id='2'></a>

# GET: List A Specific Metric profile
This method can be used to retrieve specific metric profile based on its id

## Input

```
GET /metric_profiles/{ID}
```

### Request headers

```
x-api-key: shared_key_value
Content-Type: application/json
Accept: application/json
```

## Response
Headers: `Status: 200 OK`

### Response body
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
   "name": "ch.cern.SAM.ROC_CRITICAL",
   "services": [
    {
     "service": "CREAM-CE",
     "metrics": [
      "emi.cream.CREAMCE-JobSubmit",
      "emi.wn.WN-Bi",
      "emi.wn.WN-Csh",
      "emi.wn.WN-SoftVer"
     ]
    },
    {
     "service": "SRMv2",
     "metrics": [
      "hr.srce.SRM2-CertLifetime",
      "org.sam.SRM-Del",
      "org.sam.SRM-Get",
      "org.sam.SRM-GetSURLs",
      "org.sam.SRM-GetTURLs",
      "org.sam.SRM-Ls",
      "org.sam.SRM-LsDir",
      "org.sam.SRM-Put"
     ]
    }
   ]
  }
 ]
}
```

<a id='3'></a>

# POST: Create a new Metric Profile
This method can be used to insert a new metric profile

## Input

```
POST /metric_profiles
```

### Request headers

```
x-api-key: shared_key_value
Content-Type: application/json
Accept: application/json
```

### POST BODY

```json
{
  "name": "test_profile",
  "services": [
    {
      "service": "Service-A",
      "metrics": [
        "metric.A.1",
        "metric.A.2",
        "metric.A.3",
        "metric.A.4"
      ]
    },
    {
      "service": "Service-B",
      "metrics": [
        "metric.B.1",
        "metric.B.2"
      ]
    }
  ]
}
```

## Response
Headers: `Status: 200 OK`

### Response body
Json Response

```json
{
 "status": {
  "message": "Metric Profile successfully created",
  "code": "201"
 },
 "data": {
  "id": "{{ID}}",
  "links": {
   "self": "https:///api/v2/metric_profiles/{{ID}}"
  }
 }
}
```

<a id='4'></a>

# PUT: Update information on an existing metric profile
This method can be used to update information on an existing metric profile

## Input

```
PUT /metric_profiles/{ID}
```

### Request headers

```
x-api-key: shared_key_value
Content-Type: application/json
Accept: application/json
```

### PUT BODY

```json
{
  "name": "test_profile",
  "services": [
    {
      "service": "Service-A",
      "metrics": [
        "metric.A.1",
        "metric.A.2",
        "metric.A.3",
        "metric.A.4"
      ]
    },
    {
      "service": "Service-B",
      "metrics": [
        "metric.B.1",
        "metric.B.2"
      ]
    }
  ]
}
```

## Response
Headers: `Status: 200 OK`

### Response body
Json Response

```json
{
 "status": {
  "message": "Metric Profile successfully updated",
  "code": "200"
 },
 "data": {
  "id": "{{ID}}",
  "links": {
   "self": "https:///api/v2/metric_profiles/{{ID}}"
  }
 }
}
```

<a id='5'></a>

# DELETE: Delete an existing metric profile
This method can be used to delete an existing metric profile

## Input

```
DELETE /metric_profiles/{ID}
```

### Request headers

```
x-api-key: shared_key_value
Content-Type: application/json
Accept: application/json
```


## Response
Headers: `Status: 200 OK`

### Response body
Json Response

```json
{
 "status": {
  "message": "Metric Profile Successfully Deleted",
  "code": "200"
 }
}
```
