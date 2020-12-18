---
id: feeds
title: Feeds
---
## API Calls

| Name                                  | Description                                                                     | Shortcut           |
| ------------------------------------- | ------------------------------------------------------------------------------- | ------------------ |
| GET: Feed Topology information   | This method can be used to retrieve a list of feed topology parameters         | [ Description](#1) |
| PUT: Update feed topology info | This method can be used to update feed topology information parameters | [ Description](#2) |
| GET: Feed Weights information   | This method can be used to retrieve a list of feed weights parameters         | [ Description](#3) |
| PUT: Update feed weights info | This method can be used to update feed weights information parameters | [ Description](#4) |

<a id='1'></a>

## [GET]: List Feed topology parameters

This method can be used to retrieve a list of feed topology parameters

### Input

```
GET /feeds/topology
```


### Request headers

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
   "type": "gocdb",
   "feed_url": "https://somewhere.foo.bar/topology/feed",
   "paginated": "true",
   "fetch_type": [
    "item1",
    "item2"
   ],
   "uid_endpoints": "endpointA"
  }
 ]
}
```

<a id='2'></a>

## [PUT]: Update topology feed parameters
This method is used to upadte topology feed parameters

### Input

```
PUT /feeds/topology
```

#### PUT BODY
```json
  {
   "type": "gocdb",
   "feed_url": "https://somewhere.foo.bar/topology/feed",
   "paginated": "true",
   "fetch_type": [
    "item1",
    "item2"
   ],
   "uid_endpoints": "endpointA"
  }
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
  "message": "Feeds resource succesfully updated",
  "code": "200"
 },
 "data": [
  {
   "type": "gocdb",
   "feed_url": "https://somewhere2.foo.bar/topology/feed",
   "paginated": "true",
   "fetch_type": [
    "item4",
    "item5"
   ],
   "uid_endpoints": "endpointA"
  }
 ]
}
```


<a id='3'></a>

## [GET]: List Feed weights parameters

This method can be used to retrieve a list of feed weights parameters

### Input

```
GET /feeds/weights
```


### Request headers

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
      "type": "vapor",
      "feed_url": "https://somewhere.foo.bar/weight/feed",
      "weight_type": "hepspec2006 cpu",
      "group_type": "SITES"
    }
  ]
}
```

<a id='2'></a>

## [PUT]: Update topology feed parameters
This method is used to upadte topology feed parameters

### Input

```
PUT /feeds/topology
```

#### PUT BODY
```json
 {
      "type": "item2",
      "feed_url": "https://somewhere.foo.bar/weight/feed",
      "weight_type": "weight_type2",
      "group_type": "group2"
}
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
    "message": "Feeds resource succesfully updated",
    "code": "200"
  },
  "data": [
    {
      "type": "item2",
      "feed_url": "https://somewhere2.foo.bar/weights/feed",
      "weight_type": "weight_type2",
      "group_type": "group2"
    }
  ]
}
```

