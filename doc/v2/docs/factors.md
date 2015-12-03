---
title: 'API documentation | ARGO'
page_title: API - Factors Requests
font_title: fa fa-cogs
description: API Calls for listing existing Factors
---

# API Calls

Name                                     | Description                                                                            | Shortcut
---------------------------------------- | -------------------------------------------------------------------------------------- | ------------------
GET: List Factors Requests         | This method can be used to retrieve a list of factors.          | [ Description](#1)


# GET: List Factors
This method can be used to retrieve a list of current Factors

## Input

```
GET /factors
```

### Request headers

```
x-api-key: shared_key_value
Content-Type: application/json or application/xml
Accept: application/json
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
