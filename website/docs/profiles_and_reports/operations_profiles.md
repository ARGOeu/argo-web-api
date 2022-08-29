---
id: operations_profiles
title: Operation Profiles
sidebar_position: 1
---

## API Calls

Name                                     | Description                                                                            | Shortcut
---------------------------------------- | -------------------------------------------------------------------------------------- | ------------------
GET: List Operations Profile Requests         | This method can be used to retrieve a list of current Operations profiles.          | [ Description](#1)
GET: List a specific  Operations profile         | This method can be used to retrieve a specific  Operations profile based on its id.          | [ Description](#2)
POST: Create a new  Operations profile  | This method can be used to create a new  Operations profile | [ Description](#3)
PUT: Update an Operations profile |This method can be used to update information on an existing  Operations profile | [ Description](#4)
DELETE: Delete an  Operations profile |This method can be used to delete an existing  Operations profile | [ Description](#5)

<a id='1'></a>

## [GET]: List Operations Profiles

This method can be used to retrieve a list of current  Operations profiles. 

### Input

```
GET /operations_profiles
```

#### Optional Query Parameters

Type            | Description                                                                                     | Required
--------------- | ----------------------------------------------------------------------------------------------- | --------
`name`          | Operations profile name to be used as query                                                     | NO      
`date`          | Date to retrieve a historic version of the operation profile. If no date parameter is provided the most current profile will be returned | NO

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
   "date": "2019-11-04",
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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "date": "2019-11-02",
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

## [GET]: List A Specific Operations profile
This method can be used to retrieve specific Operations profile based on its id

### Input

```
GET /operations_profiles/{ID}
```

#### Optional Query Parameters

Type            | Description                                                                                     | Required
--------------- | ----------------------------------------------------------------------------------------------- | --------
`date`          | Date to retrieve a historic version of the operation profile. If no date parameter is provided the most current profile will be returned | NO


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
   "date": "2019-11-04",
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

## [POST]: Create a new Operations Profile
This method can be used to insert a new operations profile

### Input

```
POST /operations_profiles
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### Optional Query Parameters

Type            | Description                                                                                     | Required
--------------- | ----------------------------------------------------------------------------------------------- | --------
`date`          | Date to create a  new historic version of the operation profile. If no date parameter is provided current date will be supplied automatically | NO


#### POST BODY

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

### Response
Headers: `Status: 201 Created`

#### Response body
Json Response

```json
{
 "status": {
  "message": "Operations Profile successfully created",
  "code": "201"
 },
 "data": {
  "id": "{{ID}}",
  "links": {
   "self": "https:///api/v2/operations_profiles/{{ID}}"
  }
 }
}
```

<a id='4'></a>

## [PUT]: Update information on an existing operations profile
This method can be used to update information on an existing operations profile

### Input

```
PUT /operations_profiles/{ID}
```

#### Optional Query Parameters

Type            | Description                                                                                     | Required
--------------- | ----------------------------------------------------------------------------------------------- | --------
`date`          | Date to update a historic version of the operation profile. If no date parameter is provided the current date will be supplied automatically | NO


#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### PUT BODY

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

### Response
Headers: `Status: 200 OK`

#### Response body
Json Response

```json
{
 "status": {
  "message": "Operations Profile successfully updated (new snapshot created)",
  "code": "200"
 }
}
```

<a id='5'></a>

## [DELETE]: Delete an existing aggregation profile
This method can be used to delete an existing aggregation profile

### Input

```
DELETE /operations_profiles/{ID}
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
When submitting or updating a new operations profile, validation checks are performed on json POST/PUT body for the following cases:
 - Check if user has defined more than once a state name in available states list
 - Check if user has defined more than once an operation name in operations list
 - Check if user used an undefined state in operations
 - Check if truth table statements are adequate to handle all cases

When an invalid operations profile is submitted the api responds with a validation error list:

#### Example invalid profile

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

  The above profile definition contains errors like: duplicate states, undefined states and inadequate statements in truth tables. Api response is the following:

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
