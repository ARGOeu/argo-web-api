---
id: tenants
title: Tenants
slug: /
sidebar_position: 1
---

## API Calls

| Name                                  | Description                                                                            | Shortcut           |
| ------------------------------------- | -------------------------------------------------------------------------------------- | ------------------ |
| GET: List Tenants                     | This method can be used to retrieve a list of current tenants                          | [ Description](#1) |
| GET: List a specific tenant           | This method can be used to retrieve a specific metric tenant based on its id.          | [ Description](#2) |
| POST: Create a new tenant             | This method can be used to create a new tenant                                         | [ Description](#3) |
| PUT: Update a tenant                  | This method can be used to update information on an existing tenant                    | [ Description](#4) |
| PUT: Update a tenant's info           | This method can be used to update just the info part on an existing tenant             | [ Description](#4B) |
| PUT: Update a tenant's db conf        | This method can be used to update just the db conf part on an existing tenant          | [ Description](#4C) |
| PUT: Update a tenant's topology       | This method can be used to update just the topology part on an existing tenant         | [ Description](#4D) |
| POST: Set a tenant as node            | This method can be used to set an existing tenant as node                              | [ Description](#4E) |
| POST: Unset a tenant from being a node | This method can be used to unset an existing tenant from being a node                 | [ Description](#4F) |
| DELETE: Delete a tenant               | This method can be used to delete an existing tenant                                   | [ Description](#5) |
| GET: Get a tenant's arg engine status | This method can be used to get status for a specific tenant                            | [ Description](#6) |
| PUT: Update a tenant's engine status  | This method can be used to update argo engine status information for a specific tenant | [ Description](#7) |
| GET: Get a tenant's readiness | This method can be used to get the readiness for a specific tenant                            | [ Description](#7B) |
| PUT: Update a tenant's readiness  | This method can be used to update readiness information (checks) for a specific tenant | [ Description](#7C) |
| POST: Create tenant user  | This method can be used to add a new user to existing tenant| [ Description](#8) |
| PUT: Update tenant user  | This method can be used to update information on an existing user of a specific tenant| [ Description](#9) |
| POST: Renew User's API key | This method can be used to renew user's api key | [ Description](#10) |
| DELETE: Delete Users  | This method can be used to remove and delete a user from a specific tenant| [ Description](#11) |
| GET: List Users  | This method can be used to list all users that belong to a specific tenant| [ Description](#12) |



## [GET]: List Tenants {#1}

This method can be used to retrieve a list of current tenants

**Note**: This method restricts tenant database and user information when the x-api-key token holder is a **restricted** super admin
**Note**: This method shows only tenants that have admin ui users when the x-api-key token holder is a **super_admin_ui**

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

#### Response body for super admin users

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
    "description" : "a simple tenant",
    "image" : "url to image",
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
   "topology": {
    "type": "GOCDB",
    "feed": "gocdb1.example.foo"
   },
   "users": [
    {
     "id": "acb74194-553a-11e9-8647-d663bd873d93",
     "name": "cap",
     "email": "cap@email.com",
     "api_key": "C4PK3Y",
     "roles": [
        "admin"
     ]
    },
    {
    "id": "acb74194-553a-11e9-8647-d663bd873d94",
     "name": "thor",
     "email": "thor@email.com",
     "api_key": "TH0RK3Y",
     "roles": [
        "viewer"
     ]
    }
   ]
  },
  {
   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
   "info": {
    "name": "tenant2",
    "email": "tenant2@email.com",
    "description" : "a simple tenant",
    "image" : "url to image",
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
   "topology": {
    "type": "GOCDB",
    "feed": "gocdb2.example.foo"
   },
   "users": [
    {
    "id": "acb74194-553a-11e9-8647-d663bd873d95",
     "name": "groot",
     "email": "groot@email.com",
     "api_key": "GR00TK3Y",
     "roles": [
         "admin", "admin_ui"
      ]
    },
    {
    "id": "acb74194-553a-11e9-8647-d663bd873d97",
     "name": "starlord",
     "email": "starlord@email.com",
     "api_key": "ST4RL0RDK3Y",
     "roles": [
         "admin"
      ]
    }
   ]
  }
 ]
}
```

#### Response body for restricted super admin users:

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
                "description": "a simple tenant",
                "image": "url to image",
                "website": "www.tenant1.com",
                "created": "2015-10-20 02:08:04",
                "updated": "2015-10-20 02:08:04"
            },
            "topology": {
                "type": "GOCDB",
                "feed": "gocdb1.example.foo"
            }
        },
        {
            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
            "info": {
                "name": "tenant2",
                "email": "tenant2@email.com",
                "description": "a simple tenant",
                "image": "url to image",
                "website": "www.tenant2.com",
                "created": "2015-10-20 02:08:04",
                "updated": "2015-10-20 02:08:04"
            },
            "topology": {
                "type": "GOCDB",
                "feed": "gocdb2.example.foo"
            }
        }
    ]
}
```

#### Response body for super_admin_ui users:

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
            "info": {
                "name": "tenant2",
                "email": "tenant2@email.com",
                "description": "a simple tenant",
                "image": "url to image",
                "website": "www.tenant2.com",
                "created": "2015-10-20 02:08:04",
                "updated": "2015-10-20 02:08:04"
            },
            "topology": {
                "type": "GOCDB",
                "feed": "gocdb2.example.foo"
            },
            "users": [
                {
                    "id": "acb74194-553a-11e9-8647-d663bd873d95",
                    "name": "groot",
                    "email": "groot@email.com",
                    "api_key": "GR00TK3Y",
                    "roles": ["admin", "admin_ui"]
                }
            ]
        }
    ]
}
```


## [GET]: List A Specific tenant {#2}

This method can be used to retrieve specific tenant based on its id

**Note**: This method shows only tenants that have admin ui users when the x-api-key token holder is a **super_admin_ui**

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
            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",
            "info": {
                "name": "tenant2",
                "email": "tenant2@email.com",
                "description": "a simple tenant",
                "image": "url to image",
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
            "topology": {
                "type": "GOCDB",
                "feed": "gocdb1.example.foo"
            },
            "users": [
                {
                    "id": "acb74194-553a-11e9-8647-d663bd873d95",
                    "name": "groot",
                    "email": "groot@email.com",
                    "api_key": "GR00TK3Y",
                    "roles": ["admin", "admin_ui"]
                },
                {
                    "id": "acb74194-553a-11e9-8647-d663bd873d97",
                    "name": "starlord",
                    "email": "starlord@email.com",
                    "api_key": "ST4RL0RDK3Y",
                    "roles": ["admin"]
                }
            ]
        }
    ]
}
```

#### Response body for super_admin_ui users:

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
            "info": {
                "name": "tenant2",
                "email": "tenant2@email.com",
                "description": "a simple tenant",
                "image": "url to image",
                "website": "www.tenant2.com",
                "created": "2015-10-20 02:08:04",
                "updated": "2015-10-20 02:08:04"
            },
            "topology": {
                "type": "GOCDB",
                "feed": "gocdb2.example.foo"
            },
            "users": [
                {
                    "id": "acb74194-553a-11e9-8647-d663bd873d95",
                    "name": "groot",
                    "email": "groot@email.com",
                    "api_key": "GR00TK3Y",
                    "roles": ["admin", "admin_ui"]
                }
            ]
        }
    ]
}
```

## [GET]: List A Specific User

This method can be used to retrieve specific user based on its id

### Input

```
GET /admin/users:byID/{ID}
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
            "id": "acb74194-553a-11e9-8647-d663bd873d93",
            "name": "cap",
            "email": "cap@email.com",
            "api_key": "C4PK3Y",
            "roles": ["admin"]
        }
    ]
}
```

### NOTE

Specifying the filter, `export=flat`, it will return a flat user json object

```json
{
    "id": "acb74194-553a-11e9-8647-d663bd873d93",
    "name": "cap",
    "email": "cap@email.com",
    "api_key": "C4PK3Y",
    "roles": ["admin"]
}
```


## [POST]: Create a new Tenant {#3}

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
        "description": "a simple tenant",
        "image": "url to image",
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
    "topology": {
        "type": "GOCDB",
        "feed": "gocdb.example.foo"
    },
    "users": [
        {
            "name": "cap",
            "email": "cap@email.com",
            "api_key": "C4PK3Y",
            "roles": ["admin"]
        },
        {
            "name": "thor",
            "email": "thor@email.com",
            "api_key": "TH0RK3Y",
            "roles": ["admin"]
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


## [PUT]: Update an existing tenant {#4}

This method can be used to update the whole definition of an existing tenant

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
        "description": "a changed description",
        "image": "a changed url to nwe image",
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
    "topology": {
        "type": "GOCDB",
        "feed": "gocdb.example.foo"
    },
    "users": [
        {
            "name": "cap",
            "email": "cap@email.com",
            "api_key": "C4PK3Y",
            "roles": ["admin"]
        },
        {
            "name": "thor",
            "email": "thor@email.com",
            "api_key": "TH0RK3Y",
            "roles": ["admin"]
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

## [PUT]: Update only the info part of an existing tenant {#4B}

This method can be used to update only the info part of an existing tenant

### Input

```
PUT /admin/tenants/{ID}/info
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
        "name": "Tenant1-updated",
        "email": "email1@tenant1.com",
        "description": "a changed description",
        "image": "a changed url to nwe image",
        "website": "www.tenant1-updated.com",
    }
}
```

### Response

Headers: `Status: 200 OK`

#### Response body

Json Response

```json
{
    "status": {
        "message": "Tenant information successfully updated",
        "code": "200"
    }
}
```

## [PUT]: Update only the db conf part of an existing tenant {#4C}

This method can be used to update only the db conf part of an existing tenant

### Input

```
PUT /admin/tenants/{ID}/db-conf
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### PUT BODY

```json
{
    "db_conf": [{
        "store": "ar",
        "server": "mongo-remote.foo",
        "port": 27017,
        "database":"argo_TENANT"
    }]
}
```

### Response

Headers: `Status: 200 OK`

#### Response body

Json Response

```json
{
    "status": {
        "message": "Tenant database configuration successfully updated",
        "code": "200"
    }
}
```

## [PUT]: Update only the topology part of an existing tenant {#4D}

This method can be used to update only the topology part of an existing tenant

### Input

```
PUT /admin/tenants/{ID}/topology
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### PUT BODY

```json
{
    "topology": {
        "type": "csv",
        "feed": "https://example.foo/path/to/csv",
    }
}
```

### Response

Headers: `Status: 200 OK`

#### Response body

Json Response

```json
{
    "status": {
        "message": "Tenant topology configuration successfully updated",
        "code": "200"
    }
}
```

## [POST]: Set an existing tenant as a node {#4E}

Some tenants can represent Nodes and provide results as capabilities.
This method can be used to specify that an existing tenant is also a node
### Input

```
POST /admin/tenants/{ID}/node-set
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
        "message": "Tenant has been set as node",
        "code": "200"
    }
}
```

## [POST]: Unset an existing tenant from being an node {#4F}

Some tenants can represent Nodes and provide results as capabilities.
This method can be used to unset an existing tenant from being a node (thus providing results as capabilities)
### Input

```
POST /admin/tenants/{ID}/node-unset
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
        "message": "Tenant has been unset from being a node",
        "code": "200"
    }
}
```


## [DELETE]: Delete an existing tenant {#5}

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


## [GET]: List A Specific tenant's argo-engine status {#6}

This method can be used to retrieve specific tenant's status based on its id

### Input

```
GET /admin/tenants/{ID}/status
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
                "name": "tenant1",
                "email": "email1@tenant1.com",
                "description": "a simple tenant",
                "image": "url to image",
                "website": "www.tenant1.com",
                "created": "2015-10-20 02:08:04",
                "updated": "2015-10-20 02:08:04"
            },
            "status": {
                "total_status": false,
                "ams": {
                    "metric_data": {
                        "ingestion": false,
                        "publishing": false,
                        "status_streaming": false,
                        "messages_arrived": 0
                    },
                    "sync_data": {
                        "ingestion": false,
                        "publishing": false,
                        "status_streaming": false,
                        "messages_arrived": 0
                    }
                },
                "hdfs": {
                    "metric_data": false
                },
                "engine_config": false,
                "last_check": ""
            }
        }
    ]
}
```


## [PUT]: Update argo-engine status information on an existing tenant {#7}

This method can be used to update status information on an existing tenant

### Input

```
PUT /admin/tenants/{ID}/status
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### PUT BODY

```json
{
    "ams": {
        "metric_data": {
            "ingestion": true,
            "publishing": true,
            "status_streaming": false,
            "messages_arrived": 100
        },
        "sync_data": {
            "ingestion": true,
            "publishing": false,
            "status_streaming": true,
            "messages_arrived": 200
        }
    },
    "hdfs": {
        "metric_data": true,
        "sync_data": {
            "Critical": {
                "downtimes": true,
                "group_endpoints": true,
                "blank_recompuation": true,
                "group_groups": true,
                "weights": true,
                "operations_profile": true,
                "metric_profile": true,
                "aggregation_profile": true
            }
        }
    },
    "engine_config": true,
    "last_check": "2018-08-10T12:32:45Z"
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


## [GET]: List A Specific tenant's readiness {#7B}

This method can be used to retrieve specific tenant's status based on its id

### Input

```
GET /admin/tenants/{ID}/ready
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
    "data": {
        "id": "3f9e3380-43bc-4fcd-b4c4-d0331b34e1a3",
        "name": "TENANT-FOO",
        "ready": false,
        "data": {
            "ready": true,
            "message": "AMS has data. Hdfs has data"
        },
        "topology": {
            "ready": true,
            "message": "Groups, endpoints and service-types set"
        },
        "reports": {
            "ready": false,
            "message": "Tenant doesn't have reports"
        },
        "last_check": "2026-02-06T00:00:00Z"
    }
}
```


## [PUT]: Update argo-engine readines information on an existing tenant {#7C}

This method can be used to update readiness information (checks) on an existing tenant

### Input

```
PUT /admin/tenants/{ID}/ready
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### PUT BODY

```json
{
  "data": {
    "ready": true,
    "message": "AMS has data. Hdfs has data"
  },
  "topology": {
    "ready": true,
    "message": "Groups, endpoints and service-types set"
  },
  "reports": {
    "ready": true,
    "message": "Tenant has at least one report"
  },
  "last_check": "2026-02-06T01:00:00Z"
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


## [POST]: Create new user {#8}

This method can be used to create a new user on existing tenant

### Input

```
POST /admin/tenants/{ID}/users
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### PUT BODY

```json
  {
    "name":"new_user",
    "email":"new_user@email.com",
    "roles": [
        "admin"
    ]
  }`
```

### Response

Headers: `Status: 201 OK`

#### Response body

Json Response

```json
{
 "status": {
  "message": "User was successfully created",
  "code": "201"
 },
 "data": {
  "id": "1cb883eb-8b40-428d-bce6-8ec23a9f3ca8",
  "links": {
   "self": "https:///api/v2/admin/tenants/6ac7d684-1f8e-4a02-a502-720e8f11e50b/users/1cb883eb-8b40-428d-bce6-8ec23a9f3ca8"
  }
 }
}
```

__Note__: If a user account is meant for a specific component integration then an optional field `"component"` can be specified when creating the user account such as in the following example:

```json
{
    "name":"poem_viwer_account",
    "email":"ops@email.foo",
    "roles": [
        "viewer"
    ],
    "component": "poem-viewer"
 }`

```


## [PUT]: Update user {#9}

This method can be used to update an existing user in a specific tenant

### Input

```
PUT /admin/tenants/{ID}/users/{USER_ID}
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### PUT BODY

```json
  {
    "name":"new_user",
    "email":"new_user@email.com",
    "roles": [
        "admin"
    ]
  }`
```

### Response

Headers: `Status: 200 OK`

#### Response body

Json Response

```json
{
 "status": {
  "message": "User succesfully updated",
  "code": "200"
 }
}
```

__Note__: If a user account is meant for a specific component integration then an optional field `"component"` can be specified when updating the user account such as in the following example:

```json
{
    "name":"poem_viwer_account",
    "email":"ops@email.foo",
    "roles": [
        "viewer"
    ],
    "component": "poem-viewer"
 }`

```



## [POST]: Renew User API key {#10}

This method can be used to renew a user's api access key

### Input

```
POST /admin/tenants/{ID}/users/{USER_ID}/renew_api_key
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
    "message": "User api key succesfully renewed",
    "code": "200"
  },
  "data": {
    "api_key": "s3cr3tT0k3n"
  }
}
```


## [DELETE]: Delete User {#11}

This method can be used to remove and delete a user from a tenant

### Input

```
DELETE /admin/tenants/{ID}/users/{USER_ID}
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
  "message": "User succesfully deleted",
  "code": "200"
 }
}
```


## [GET]: List all available users that belong to a specific tenant {#12}
This method can be used to list all available users that are members of a specific tenant

### Input

```
GET /admin/tenants/{ID}/users
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
   "id": "acb74194-553a-11e9-8647-d663bd873d93",
   "name": "user_a",
   "email": "user_a@email.com",
   "api_key": "user_a_key",
   "roles": [
    "admin",
    "admin_ui"
   ]
  },
  {
   "id": "acb74432-553a-11e9-8647-d663bd873d93",
   "name": "user_b",
   "email": "user_b@email.com",
   "api_key": "user_b_key",
   "roles": [
    "admin"
   ]
  }
 ]
}

```


## [GET]: Get user {#13}
This method can be used to get information on specific user of a specific tenant

### Input

```
GET /admin/tenants/{ID}/users/{USER_ID}
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
   "id": "acb74432-553a-11e9-8647-d663bd873d93",
   "name": "user_a",
   "email": "user_a@email.com",
   "api_key": "user_a_key",
   "roles": [
    "admin"
   ]
  }
 ]
}

```

















