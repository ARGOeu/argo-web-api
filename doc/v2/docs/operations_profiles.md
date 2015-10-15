---
title: 'API documentation | ARGO'
page_title: API - Operations Profile Requests
font_title: fa fa-cogs
description: API Calls for listing existing and creating new operations profiles
---

# API Calls

Name                                     | Description                                                                            | Shortcut
---------------------------------------- | -------------------------------------------------------------------------------------- | ------------------
GET: List Operations Profile Requests         | This method can be used to retrieve a list of current Operations profiles.          | [ Description](#1)
GET: List a specific  Operations profile         | This method can be used to retrieve a specific  Operations profile based on its uuid.          | [ Description](#2)
POST: Create a new  Operations profile  | This method can be used to create a new  Operations profile | [ Description](#3)
PUT: Update an Operations profile |This method can be used to update information on an existing  Operations profile | [ Description](#4)
DELETE: Delete an  Operations profile |This method can be used to delete an existing  Operations profile | [ Description](#5)
<a id='1'></a>

# GET: List AOperations Profiles
This method can be used to retrieve a list of current  Operations profiles

## Input

```
GET /operations_profiles
```

#### Optional Query Parameters

Type            | Description                                                                                     | Required
--------------- | ----------------------------------------------------------------------------------------------- | --------
`name`  | Operations profile name to be used as query                                                                          | NO      

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
   "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "name": "ops1",
   "available_states": [
    "A,B,C"
   ],
   "defaults": {
    "down": "B",
    "missing": "A",
    "unknown": "C"
   },
   "operations": [
    {
     "name": "AND",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "B"
      },
      {
       "a": "A",
       "b": "C",
       "x": "C"
      },
      {
       "a": "B",
       "b": "C",
       "x": "C"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "A"
      },
      {
       "a": "A",
       "b": "C",
       "x": "A"
      },
      {
       "a": "B",
       "b": "C",
       "x": "B"
      }
     ]
    }
   ]
  },
  {
   "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "name": "ops2",
   "available_states": [
    "X,Y,Z"
   ],
   "defaults": {
    "down": "Y",
    "missing": "X",
    "unknown": "Z"
   },
   "operations": [
    {
     "name": "AND",
     "truth_table": [
      {
       "a": "X",
       "b": "Y",
       "x": "Y"
      },
      {
       "a": "X",
       "b": "Z",
       "x": "Z"
      },
      {
       "a": "Y",
       "b": "Z",
       "x": "Z"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "X",
       "b": "Y",
       "x": "X"
      },
      {
       "a": "X",
       "b": "Z",
       "x": "X"
      },
      {
       "a": "Y",
       "b": "Z",
       "x": "Y"
      }
     ]
    }
   ]
  }
 ]
}
```

<a id='2'></a>

# GET: List A Specific Operations profile
This method can be used to retrieve specific Operations profile based on its uuid

## Input

```
GET /operations_profiles/{UUID}
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
   "uuid": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "name": "ops1",
   "available_states": [
    "A,B,C"
   ],
   "defaults": {
    "down": "B",
    "missing": "A",
    "unknown": "C"
   },
   "operations": [
    {
     "name": "AND",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "B"
      },
      {
       "a": "A",
       "b": "C",
       "x": "C"
      },
      {
       "a": "B",
       "b": "C",
       "x": "C"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "A"
      },
      {
       "a": "A",
       "b": "C",
       "x": "A"
      },
      {
       "a": "B",
       "b": "C",
       "x": "B"
      }
     ]
    }
   ]
  }
 ]
}
```

<a id='3'></a>

# POST: Create a new Operations Profile
This method can be used to insert a new operations profile

## Input

```
POST /operations_profiles
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
   "name": "tops1",
   "available_states": [
    "A","B","C"
   ],
   "defaults": {
    "down": "B",
    "missing": "A",
    "unknown": "C"
   },
   "operations": [
    {
     "name": "AND",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "B"
      },
      {
       "a": "A",
       "b": "C",
       "x": "C"
      },
      {
       "a": "B",
       "b": "C",
       "x": "C"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "A"
      },
      {
       "a": "A",
       "b": "C",
       "x": "A"
      },
      {
       "a": "B",
       "b": "C",
       "x": "B"
      }
     ]
    }
   ]
  }
```

## Response
Headers: `Status: 201 Created`

### Response body
Json Response

```json
{
 "status": {
  "message": "Operations Profile successfully created",
  "code": "201"
 },
 "data": {
  "uuid": "{{UUID}}",
  "links": {
   "self": "https:///api/v2/operations_profiles/{{UUID}}"
  }
 }
}
```

<a id='4'></a>

# PUT: Update information on an existing operations profile
This method can be used to update information on an existing operations profile

## Input

```
PUT /operations_profiles/{UUID}
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
	 "name": "tops1",
	 "available_states": [
		"A","B","C"
	 ],
	 "defaults": {
		"down": "B",
		"missing": "A",
		"unknown": "C"
	 },
	 "operations": [
		{
		 "name": "AND",
		 "truth_table": [
			{
			 "a": "A",
			 "b": "B",
			 "x": "B"
			},
			{
			 "a": "A",
			 "b": "C",
			 "x": "C"
			},
			{
			 "a": "B",
			 "b": "C",
			 "x": "C"
			}
		 ]
		},
		{
		 "name": "OR",
		 "truth_table": [
			{
			 "a": "A",
			 "b": "B",
			 "x": "A"
			},
			{
			 "a": "A",
			 "b": "C",
			 "x": "A"
			},
			{
			 "a": "B",
			 "b": "C",
			 "x": "B"
			}
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
  "message": "Operations Profile successfully updated",
  "code": "200"
 }
}
```

<a id='5'></a>

# DELETE: Delete an existing aggregation profile
This method can be used to delete an existing aggregation profile

## Input

```
DELETE /operations_profiles/{UUID}
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
  "message": "Operations Profile Successfully Deleted",
  "code": "200"
 }
}
```

# Validation Checks
When submitting or updating a new operations profile, validation checks are performed on json POST/PUT body for the following cases:
 - Check if user has defined more than once a state name in available states list
 - Check if user has defined more than once an operation name in operations list
 - Check if user used an undefined state in operations
 - Check if truth table statements are adequate to handle all cases

When an invalid operations profile is submitted the api responds with a validation error list:

### Example invalid profile

```json
{
   "name": "ops1",
   "available_states": [
    "A","B","C","C"
   ],
   "defaults": {
    "down": "B",
    "missing": "FOO",
    "unknown": "C"
   },
   "operations": [
    {
     "name": "AND",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "B"
      },
      {
       "a": "A",
       "b": "C",
       "x": "C"
      },
      {
       "a": "B",
       "b": "BAR",
       "x": "C"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "A"
      },
      {
       "a": "A",
       "b": "C",
       "x": "A"
      },
      {
       "a": "B",
       "b": "CAR",
       "x": "B"
      }
     ]
    },
    {
     "name": "OR",
     "truth_table": [
      {
       "a": "A",
       "b": "B",
       "x": "A"
      },
      {
       "a": "A",
       "b": "C",
       "x": "A"
      },
      {
       "a": "B",
       "b": "C",
       "x": "B"
      }
     ]
    }
   ]
  }
  ```

  The above profile definiton contains errors like: duplicate states, undefined states and unadequate statements in truth tables. Api response is the following:

## Response
Headers: `Status: 422 Unprocessable Entity`

### Response body

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
   "details": "State:C is duplicated"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Operation:OR is duplicated"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Default Missing State: FOO not in available States"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "In Operation: AND, statement member b: BAR contains undeclared state"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "In Operation: OR, statement member b: CAR contains undeclared state"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:A in operation: AND"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:B in operation: AND"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:C in operation: AND"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:A in operation: OR"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:B in operation: OR"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:C in operation: OR"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:A in operation: OR"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:B in operation: OR"
  },
  {
   "message": "Validation Failed",
   "code": "422",
   "details": "Not enough mentions of state:C in operation: OR"
  }
 ]
}
 ```
