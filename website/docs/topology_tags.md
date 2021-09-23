---
id: topology_tags
title: Topology Tags & Values
---

## API calls for all available tags and distinct values present in topology items

| Name                          | Description                                           | Shortcut                     |
| ----------------------------- | ----------------------------------------------------- | ---------------------------- |
| GET: List topology tags | List available tags and distinct values present in topology items. | <a href="#1">Description</a> |

<a id="1"></a>

## [GET]: List topology tags

This method may be used to retrieve available topology tags and their distinct values available 

### Input

##### List All topology tag values

```
/topology/tags?[date]
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

#### Response Code

```
Status: 200 OK
```

### Response body

```json
{
    "status": {
        "message": "application/json",
        "code": "200"
    },
    "data": [
        {
            "name": "endpoints",
            "values": [
                {
                    "name": "tag1",
                    "values": ["value1", "value2", "value3"]
                },
                {
                    "name": "tag2",
                    "values": ["value1", "value2"]
                }
            ]
        },
        {
            "name": "groups",
            "values": [
                {
                    "name": "tag3",
                    "values": ["value1", "value2", "value3"]
                },
                {
                    "name": "tag4",
                    "values": ["value1"]
                }
            ]
        }
    ]
}
```

###### Example Request:

URL:

```
/topology/tags?date=2016-01-01
```


Headers:

```
x-api-key: secret_key_value
Accept: application/json
```

###### Example Response:

Code:

```
Status: 200 OK
```

Response body:

```json
{
 "status": {
  "message": "Success",
  "code": "200"
 },
 "data": [
  {
   "name": "endpoints",
   "values": [
    {
     "name": "scope",
     "values": [
      "GROUPB",
      "GROUPA",
      "GROUPE",
      "GROUPD",
      "GROUPC"
     ]
    },
    {
     "name": "production",
     "values": [
      "0",
      "1"
     ]
    },
    {
     "name": "monitored",
     "values": [
      "0",
      "1"
     ]
    }
   ]
  },
  {
   "name": "groups",
   "values": [
    {
     "name": "certification",
     "values": [
      "Certified",
      "Uncertified"
     ]
    },
    {
     "name": "infrastructure",
     "values": [
      "production",
      "devel",
      "devtest"
     ]
    }
   ]
  }
 ]
}
```
