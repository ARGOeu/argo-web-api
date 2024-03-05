---
id: factors
title: Factors
sidebar_position: 7
---

## API Calls

Name                                     | Description                                                                            | Shortcut
---------------------------------------- | -------------------------------------------------------------------------------------- | ------------------
GET: List Factors Requests         | This method can be used to retrieve a list of factors.          | [ Description](#1)


## GET: List Factors {#1}
This method can be used to retrieve a list of current Factors

## Input

```
GET /factors
```

### Request headers

```
x-api-key: shared_key_value
Accept: application/json or application/xml
```

## Response
Headers: `Status: 200 OK`

### Response body
Json Response
```json
{
 "factors": [
  {
   "site": "CETA-GRID",
   "weight": "5406"
  },
  {
   "site": "CFP-IST",
   "weight": "1019"
  },
  {
   "site": "CIEMAT-LCG2",
   "weight": "14595"
  }
 ]
}
```

XML Response

```xml
<root>
    <Factor site="CETA-GRID" weight="5406"></Factor>
    <Factor site="CFP-IST" weight="1019"></Factor>
    <Factor site="CIEMAT-LCG2" weight="14595"></Factor>
</root>
```
