#Topology Groups

API calls for handling topology group resources

| Name                                             | Description                                                         | Shortcut                     |
| ------------------------------------------------ | ------------------------------------------------------------------- | ---------------------------- |
| POST: Create endpoint topology for specific date | Creates a daily group topology mapping endpoints to endpoint groups | <a href="#1">Description</a> |

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
