---
id: health
title: API Health Information
---

The users can receive information about the health of the argo-web-api service



## GET Health information

__note: Health information can be accessed via both the `/api/v2` and `/api/v3` endpoints. For consistency, the examples provided will utilize `/api/v3`__

A user can get health information about the argo-web-api instance by issuing the following call:

```
GET /api/v3/health
```

### Request headers

```
Accept: application/json
```

### Response
Headers: `Status: 200 OK`

### Response Body

Json Response example:
```json
{
  "status": "OK",
  "timestamp": "2025-07-10T12:30:00Z",
  "message": "No flapping items"
}
```
