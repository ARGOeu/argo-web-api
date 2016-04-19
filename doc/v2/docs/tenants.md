---
title: 'API documentation | ARGO'
page_title: API - Tenant Requests
font_title: fa fa-cogs
description: API Calls for listing existing and creating new tenant
---

# API Calls

Name                                     | Description                                                                            | Shortcut
---------------------------------------- | -------------------------------------------------------------------------------------- | ------------------
GET: List Tenants        | This method can be used to retrieve a list of current tenants          | [ Description](#1)
GET: List a specific tenant         | This method can be used to retrieve a specific metric tenant based on its id.          | [ Description](#2)
POST: Create a new tenant  | This method can be used to create a new tenant | [ Description](#3)
PUT: Update a tenant |This method can be used to update information on an existing tenant | [ Description](#4)
DELETE: Delete a tenant |This method can be used to delete an existing tenant | [ Description](#5)

<a id='1'></a>

## [GET]: List Tenants
This method can be used to retrieve a list of current tenants

### Input

```
GET /admin/tenants
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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "info": {
    "name": "Tenant1",
    "email": "email1@tenant1.com",
    "website": "www.tenant1.com",
    "created": "2015-10-20 02:08:04",
    "updated": "2015-10-20 02:08:04"
   },
   "db_conf": [
    {
     "store": "ar",
     "server": "a.mongodb.org",
     "port": 27017,
     "database": "ar_db",
     "username": "admin",
     "password": "3NCRYPT3D"
    },
    {
     "store": "status",
     "server": "b.mongodb.org",
     "port": 27017,
     "database": "status_db",
     "username": "admin",
     "password": "3NCRYPT3D"
    }
   ],
   "users": [
    {
     "name": "cap",
     "email": "cap@email.com",
     "api_key": "C4PK3Y"
    },
    {
     "name": "thor",
     "email": "thor@email.com",
     "api_key": "TH0RK3Y"
    }
   ]
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "info": {
    "name": "tenant2",
    "email": "tenant2@email.com",
    "website": "www.tenant2.com",
    "created": "2015-10-20 02:08:04",
    "updated": "2015-10-20 02:08:04"
   },
   "db_conf": [
    {
     "store": "ar",
     "server": "a.mongodb.org",
     "port": 27017,
     "database": "ar_db",
     "username": "admin",
     "password": "3NCRYPT3D"
    },
    {
     "store": "status",
     "server": "b.mongodb.org",
     "port": 27017,
     "database": "status_db",
     "username": "admin",
     "password": "3NCRYPT3D"
    }
   ],
   "users": [
    {
     "name": "groot",
     "email": "groot@email.com",
     "api_key": "GR00TK3Y"
    },
    {
     "name": "starlord",
     "email": "starlord@email.com",
     "api_key": "ST4RL0RDK3Y"
    }
   ]
  }
 ]
}
```

<a id='2'></a>

## [GET]: List A Specific tenant
This method can be used to retrieve specific tenant based on its id

### Input

```
GET /admin/tenants/{ID}
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
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",
   "info": {
    "name": "Tenant1",
    "email": "email1@tenant1.com",
    "website": "www.tenant1.com",
    "created": "2015-10-20 02:08:04",
    "updated": "2015-10-20 02:08:04"
   },
   "db_conf": [
    {
     "store": "ar",
     "server": "a.mongodb.org",
     "port": 27017,
     "database": "ar_db",
     "username": "admin",
     "password": "3NCRYPT3D"
    },
    {
     "store": "status",
     "server": "b.mongodb.org",
     "port": 27017,
     "database": "status_db",
     "username": "admin",
     "password": "3NCRYPT3D"
    }
   ],
   "users": [
    {
     "name": "cap",
     "email": "cap@email.com",
     "api_key": "C4PK3Y"
    },
    {
     "name": "thor",
     "email": "thor@email.com",
     "api_key": "TH0RK3Y"
    }
   ]
  }
 ]
}
```

<a id='3'></a>

## [POST]: Create a new Tenant
This method can be used to insert a new tenant

### Input

```
POST /admin/tenants
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### POST BODY

```json
{
  "info": {
    "name": "Tenant1",
    "email": "email1@tenant1.com",
    "website": "www.tenant1.com",
    "created": "2015-10-20 02:08:04",
    "updated": "2015-10-20 02:08:04"
  },
  "db_conf": [
    {
      "store": "ar",
      "server": "a.mongodb.org",
      "port": 27017,
      "database": "ar_db",
      "username": "admin",
      "password": "3NCRYPT3D"
    },
    {
      "store": "status",
      "server": "b.mongodb.org",
      "port": 27017,
      "database": "status_db",
      "username": "admin",
      "password": "3NCRYPT3D"
    }
  ],
  "users": [
    {
      "name": "cap",
      "email": "cap@email.com",
      "api_key": "C4PK3Y"
    },
    {
      "name": "thor",
      "email": "thor@email.com",
      "api_key": "TH0RK3Y"
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
  "message": "Tenant was succesfully created",
  "code": "201"
 },
 "data": {
  "id": "{{ID}}",
  "links": {
   "self": "https:///api/v2/admin/tenants/{{ID}}"
  }
 }
}
```

<a id='4'></a>

## [PUT]: Update information on an existing tenant
This method can be used to update information on an existing tenant

### Input

```
PUT /admin/tenants/{ID}
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### PUT BODY

```json
{
  "info": {
    "name": "Tenant1",
    "email": "email1@tenant1.com",
    "website": "www.tenant1.com",
    "created": "2015-10-20 02:08:04",
    "updated": "2015-10-20 02:08:04"
  },
  "db_conf": [
    {
      "store": "ar",
      "server": "a.mongodb.org",
      "port": 27017,
      "database": "ar_db",
      "username": "admin",
      "password": "3NCRYPT3D"
    },
    {
      "store": "status",
      "server": "b.mongodb.org",
      "port": 27017,
      "database": "status_db",
      "username": "admin",
      "password": "3NCRYPT3D"
    }
  ],
  "users": [
    {
      "name": "cap",
      "email": "cap@email.com",
      "api_key": "C4PK3Y"
    },
    {
      "name": "thor",
      "email": "thor@email.com",
      "api_key": "TH0RK3Y"
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
  "message": "Tenant successfully updated",
  "code": "200"
 }
}
```

<a id='5'></a>

## [DELETE]: Delete an existing tenant
This method can be used to delete an existing tenant

### Input

```
DELETE /admin/tenants/{ID}
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
  "message": "Tenant Successfully Deleted",
  "code": "200"
 }
}
```
