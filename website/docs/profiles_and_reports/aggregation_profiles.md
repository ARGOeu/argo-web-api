---
id: aggregation_profiles
title: Aggregation Profiles
sidebar_position: 3
---

## API Calls

Name                                     | Description                                                                            | Shortcut
---------------------------------------- | -------------------------------------------------------------------------------------- | ------------------
GET: List Aggregation Profile Requests         | This method can be used to retrieve a list of current aggregation profiles.          | [ Description](#1)
GET: List a specific  aggregation profile         | This method can be used to retrieve a specific  aggregation profile based on its id.          | [ Description](#2)
POST: Create a new  aggregation profile  | This method can be used to create a new  aggregation profile | [ Description](#3)
PUT: Update an aggregation profile |This method can be used to update information on an existing  aggregation profile | [ Description](#4)
DELETE: Delete an  aggregation profile |This method can be used to delete an existing  aggregation profile | [ Description](#5)
<a id='1'></a>

## [GET]: List Aggregation Profiles
This method can be used to retrieve a list of current  aggregation profiles

### Input

```
GET /aggregation_profiles
```

#### Optional Query Parameters

Type            | Description                                                                                     | Required
--------------- | ----------------------------------------------------------------------------------------------- | --------
`name`          | aggregation profile name to be used as query                                                    | NO      
`date`          | Date to retrieve a historic version of the aggregations profiles. If no date parameter is provided the most current profile will be returned | NO

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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "date": "2019-10-10",
   "name": "cloud",
   "namespace": "test",
   "endpoint_group": "sites",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "roc.critical",
    "id": "5637d684-1f8e-4a02-a502-720e8f11e432"
   },
   "groups": [
    {
     "name": "compute",
     "operation": "OR",
     "services": [
      {
       "name": "SERVICEA",
       "operation": "AND"
      },
      {
       "name": "SERVICEB",
       "operation": "AND"
      }
     ]
    },
    {
     "name": "images",
     "operation": "OR",
     "services": [
      {
       "name": "SERVICEC",
       "operation": "AND"
      },
      {
       "name": "SERVICED",
       "operation": "AND"
      }
     ]
    }
   ]
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "date": "2019-11-11",
   "name": "critical",
   "namespace": "test",
   "endpoint_group": "sites",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "roc.critical",
    "id": "5637d684-1f8e-4a02-a502-720e8f11e432"
   },
   "groups": [
    {
     "name": "compute",
     "operation": "OR",
     "services": [
      {
       "name": "CREAM-CE",
       "operation": "AND"
      },
      {
       "name": "ARC-CE",
       "operation": "AND"
      }
     ]
    },
    {
     "name": "storage",
     "operation": "OR",
     "services": [
      {
       "name": "SRMv2",
       "operation": "AND"
      },
      {
       "name": "SRM",
       "operation": "AND"
      }
     ]
    }
   ]
  }
 ]
}
```

<a id='2'></a>

## [GET]: List A Specific Aggregation profile
This method can be used to retrieve specific aggregation profile based on its id

### Input

```
GET /aggregation_profiles/{ID}
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```


#### Optional Query Parameters

Type            | Description                                                                                     | Required
--------------- | ----------------------------------------------------------------------------------------------- | --------
`date`          | Date to list a historic version of the aggregation profile. If no date parameter is provided current date will be supplied automatically | NO

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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "date": "2019-10-10",
   "name": "cloud",
   "namespace": "test",
   "endpoint_group": "sites",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "roc.critical",
    "id": "5637d684-1f8e-4a02-a502-720e8f11e432"
   },
   "groups": [
    {
     "name": "compute",
     "operation": "OR",
     "services": [
      {
       "name": "SERVICEA",
       "operation": "AND"
      },
      {
       "name": "SERVICEB",
       "operation": "AND"
      }
     ]
    },
    {
     "name": "images",
     "operation": "OR",
     "services": [
      {
       "name": "SERVICEC",
       "operation": "AND"
      },
      {
       "name": "SERVICED",
       "operation": "AND"
      }
     ]
    }
   ]
  }
 ]
}
```

<a id='3'></a>

## [POST]: Create a new Aggregation Profile
This method can be used to insert a new aggregation profile

### Input

```
POST /aggregation_profiles
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```


#### Optional Query Parameters

Type            | Description                                                                                     | Required
--------------- | ----------------------------------------------------------------------------------------------- | --------
`date`          | Date to create a  new historic version of the aggregation profile. If no date parameter is provided current date will be supplied automatically | NO

#### POST BODY

```json
{
   "name": "new_aggregation_profile",
   "date": "2019-12-12",
   "namespace": "test",
   "endpoint_group": "sites",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "test.metric.profile",
    "id": "5637d684-1f8e-4a02-a502-720e8f11e432"
   },
   "groups": [
    {
     "name": "compute",
     "operation": "OR",
     "services": [
      {
       "name": "SERVICEA",
       "operation": "AND"
      },
      {
       "name": "SERVICEB",
       "operation": "AND"
      }
     ]
    },
    {
     "name": "images",
     "operation": "OR",
     "services": [
      {
       "name": "SERVICEC",
       "operation": "AND"
      },
      {
       "name": "SERVICED",
       "operation": "AND"
      }
     ]
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
  "message": "Aggregation Profile successfully created",
  "code": "201"
 },
 "data": {
  "id": "{{ID}}",
  "links": {
   "self": "https:///api/v2/aggregation_profiles/{{ID}}"
  }
 }
}
```

<a id='4'></a>

## [PUT]: Update information on an existing aggregation profile
This method can be used to update information on an existing aggregation profile

### Input

```
PUT /aggregation_profiles/{ID}
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```


#### Optional Query Parameters

Type            | Description                                                                                     | Required
--------------- | ----------------------------------------------------------------------------------------------- | --------
`date`          | Date to update a  new historic version of the aggregation profile. If no date parameter is provided current date will be supplied automatically | NO

#### PUT BODY

```json
{
   "name": "updated_profile",
   "namespace": "test",
   "endpoint_group": "sites",
   "metric_operation": "AND",
   "profile_operation": "AND",
   "metric_profile": {
    "name": "test.metric.profile",
    "id": "5637d684-1f8e-4a02-a502-720e8f11e432"
   },
   "groups": [
    {
     "name": "compute",
     "operation": "OR",
     "services": [
      {
       "name": "SERVICEA",
       "operation": "AND"
      },
      {
       "name": "SERVICEB",
       "operation": "AND"
      }
     ]
    },
    {
     "name": "images",
     "operation": "OR",
     "services": [
      {
       "name": "SERVICEC",
       "operation": "AND"
      },
      {
       "name": "SERVICED",
       "operation": "AND"
      }
     ]
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
  "message": "Aggregation Profile successfully updated",
  "code": "200"
 },
 "data": {
  "id": "{{ID}}",
  "links": {
   "self": "https:///api/v2/aggregation_profiles/{{ID}}"
  }
 }
}
```

<a id='5'></a>

## [DELETE]: Delete an existing aggregation profile
This method can be used to delete an existing aggregation profile

### Input

```
DELETE /aggregation_profiles/{ID}
```

#### Request headers

```
x-api-key: shared_key_value
Content-Type: application/json
Accept: application/json
```


### Response
Headers: `Status: 200 OK`

#### Response body
Json Response

```json
{
 "status": {
  "message": "Aggregation Profile Successfully Deleted",
  "code": "200"
 }
}
```
