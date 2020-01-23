#Topology Endpoints

API calls for handling topology endpoint resources

| Name                                             | Description                                                            | Shortcut                     |
| ------------------------------------------------ | ---------------------------------------------------------------------- | ---------------------------- |
| POST: Create endpoint topology for specific date | Creates a daily endpoint topology mapping endpoints to endpoint groups | <a href="#1">Description</a> |

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
| `date` | target a specific data | NO       | today's date  |

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
