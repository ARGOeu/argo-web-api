# Topology Endpoints

API calls for handling topology endpoint resources

| Name                                               | Description                                                            | Shortcut                     |
| -------------------------------------------------- | ---------------------------------------------------------------------- | ---------------------------- |
| POST: Create endpoint topology for specific date   | Creates a daily endpoint topology mapping endpoints to endpoint groups | <a href="#1">Description</a> |
| GET: List endpoint topology for specific date      | Lists endpoint topology for a specific date                            | <a href="#2">Description</a> |
| DELETE: delete endpoint topology for specific date | Deletes all endpoint items (topology) for a specific date              | <a href="#3">Description</a> |

<a id="1"></a>

## POST: Create endpoint topology for specific date

Creates a daily endpoint topology mapping endpoints to endpoint groups

### Input

```
POST /topology/endpoints?date=YYYY-MM-DD
```

#### Url Parameters

| Type   | Description            | Required | Default value |
| ------ | ---------------------- | -------- | ------------- |
| `date` | target a specific date | NO       | today's date  |

#### Headers

```
x-api-key: secret_key_value
Accept: application/json
```

### POST BODY

```json
[
    {
        "group": "SITE_A",
        "hostname": "host1.site-a.foo",
        "type": "SITES",
        "service": "a.service.foo",
        "tags": { "scope": "TENANT", "production": "1", "monitored": "1" }
    },
    {
        "group": "SITE_A",
        "hostname": "host2.site-b.foo",
        "type": "SITES",
        "service": "b.service.foo",
        "tags": { "scope": "TENANT", "production": "1", "monitored": "1" }
    },
    {
        "group": "SITE_B",
        "hostname": "host1.site-a.foo",
        "type": "SITES",
        "service": "c.service.foo",
        "tags": { "scope": "TENANT", "production": "1", "monitored": "1" }
    }
]
```

#### Response Code

```
Status: 201 OK Created
```

### Response body

```json
{
    "message": "Topology of 3 endpoints created for date: YYYY-MM-DD",
    "code": "201"
}
```

## 409 Conflict when trying to insert a topology that already exists

When trying to insert a topology for a specific date that already exists the api will answer with the following reponse:

### Response Code

```
Status: 409 Conflict
```

### Response body

```json
{
    "message": "topology already exists for date: YYYY-MM-DD, please either update it or delete it first!",
    "code": "409"
}
```

User can proceed with either updating the existing endpoint topology OR deleting before trying to create it anew

<a id="2"></a>

## GET: List endpoint topology per date

List endpoint topology for a specific date or the closest availabe topology to that date. If date is not provided list the latest available endpoint topology.

### Input

##### List All topology statistics

```
GET /topology/endpoints?date=YYYY-MM-DD
```

#### Url Parameters

| Type       | Description                   | Required | Default value |
| ---------- | ----------------------------- | -------- | ------------- |
| `date`     | target a specific date        | NO       | today's date  |
| `group`    | filter by group name          | NO       |               |
| `type`     | filter by group type          | NO       |               |
| `service`  | filter by service             | NO       |               |
| `hostname` | filter by hostname            | NO       |               |
| `tags`     | filter by tag key:value pairs | NO       |               |

_note_ : user can use wildcard \* in filters
_note_ : when using tags filter the query string must follow the pattern: `?tags=key1:value1,key2:value2`

#### Headers

```
x-api-key: secret_key_value
Accept: application/json
```

#### Response Code

```
Status: 200 OK
```

### Response body

```json
[
    {
        "date": "2019-12-12",
        "group": "SITE_A",
        "hostname": "host1.site-a.foo",
        "type": "SITES",
        "service": "a.service.foo",
        "tags": { "scope": "TENANT", "production": "1", "monitored": "1" }
    },
    {
        "date": "2019-12-12",
        "group": "SITE_A",
        "hostname": "host2.site-b.foo",
        "type": "SITES",
        "service": "b.service.foo",
        "tags": { "scope": "TENANT", "production": "1", "monitored": "1" }
    },
    {
        "date": "2019-12-12",
        "group": "SITE_B",
        "hostname": "host1.site-a.foo",
        "type": "SITES",
        "service": "c.service.foo",
        "tags": { "scope": "TENANT", "production": "1", "monitored": "1" }
    }
]
```

<a id='3'></a>

## [DELETE]: Delete endpoint topology for a specific date

This method can be used to delete all endpoint items contributing to the endpoint topology of a specific date

### Input

```
DELETE /topology/endpoints?date=YYYY-MM-DD
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
    "message": "Topology of 3 endpoints deleted for date: 2019-12-12",
    "code": "200"
}
```
