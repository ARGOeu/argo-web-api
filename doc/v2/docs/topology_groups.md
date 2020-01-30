# Topology Groups

API calls for handling topology group resources

| Name                                            | Description                                                         | Shortcut                     |
| ----------------------------------------------- | ------------------------------------------------------------------- | ---------------------------- |
| POST: Create group topology for specific date   | Creates a daily group topology mapping endpoints to endpoint groups | <a href="#1">Description</a> |
| GET: List group topology for specific date      | Lists group topology for a specific date                            | <a href="#2">Description</a> |
| DELETE: Delete group topology for specific date | Delete group topology items for specific date                       | <a href="#3">Description</a> |

<a id="1"></a>

## POST: Create group topology for specific date

Creates a daily group topology mapping top-level groups to subgroups

### Input

```
POST /topology/groups?date=YYYY-MM-DD
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
        "group": "NGIA",
        "type": "NGIS",
        "service": "SITEA",
        "tags": {
            "scope": "FEDERATION",
            "infrastructure": "Production",
            "certification": "Certified"
        }
    },
    {
        "group": "NGIA",
        "type": "NGIS",
        "service": "SITEB",
        "tags": {
            "scope": "FEDERATION",
            "infrastructure": "Production",
            "certification": "Certified"
        }
    },
    {
        "group": "NGIZ",
        "type": "NGIS",
        "service": "SITEZ",
        "tags": {
            "scope": "FEDERATION",
            "infrastructure": "Production",
            "certification": "Certified"
        }
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
    "message": "Topology of 3 groups created for date: YYYY-MM-DD",
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

User can proceed with either updating the existing topology OR deleting before trying to create it anew

<a id="2"></a>

## POST: List group topology for specific date

Lists group topology items for specific date

### Input

```
GET /topology/groups?date=YYYY-MM-DD
```

#### Url Parameters

| Type       | Description                   | Required | Default value |
| ---------- | ----------------------------- | -------- | ------------- |
| `date`     | target a specific date        | NO       | today's date  |
| `group`    | filter by group name          | NO       |               |
| `type`     | filter by group type          | NO       |               |
| `subgroup` | filter by subgroup            | NO       |               |
| `tags`     | filter by tag key:value pairs | NO       |               |

_note_ : user can use wildcard \* in filters
_note_ : when using tag filters the query string must follow the pattern: `?tags=key1:value1,key2:value2`

_note_ : user can use wildcard \* in filters

#### Headers

```
x-api-key: secret_key_value
Accept: application/json
```

#### Example Request

```
GET /topology/groups?date=2015-07-22
```

#### Response Code

```
Status: 200 OK
```

### Response body

```json
{
    "status": {
        "message": "Success",
        "code": "200"
    },
    "data": [
        {
            "date": "2015-07-22",
            "group": "NGIA",
            "type": "NGIS",
            "subgroup": "SITEA",
            "tags": {
                "certification": "Certified",
                "infrastructure": "Production"
            }
        },
        {
            "date": "2015-07-22",
            "group": "NGIA",
            "type": "NGIS",
            "subgroup": "SITEB",
            "tags": {
                "certification": "Certified",
                "infrastructure": "Production"
            }
        },
        {
            "date": "2015-07-22",
            "group": "NGIX",
            "type": "NGIS",
            "subgroup": "SITEX",
            "tags": {
                "certification": "Certified",
                "infrastructure": "Production"
            }
        }
    ]
}
```

<a id='3'></a>

## [DELETE]: Delete group topology for a specific date

This method can be used to delete all group items contributing to the group topology of a specific date

### Input

```
DELETE /topology/groups?date=YYYY-MM-DD
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
    "message": "Topology of 3 groups deleted for date: 2019-12-12",
    "code": "200"
}
```
