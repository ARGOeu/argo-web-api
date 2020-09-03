---
id: version
title: API Version Information
---

This method can be used to retrieve api version information

## Input

```
GET /version
```

### Request headers

```
Accept: application/json or application/xml
```

## Response
Headers: `Status: 200 OK`

## Response Body

Json Response
```json
{
    "release": "1.7.9",
    "commit": "f9f2e8c5f02fbcc93fe76b0d3cfa5d2089336849",
    "build_time": "2019-11-01T12:51:04Z",
    "golang": "go1.11.5",
    "compiler": "gc",
    "os": "linux",
    "architecture": "amd64"
}
```

XML Response
```xml
<Version>
    <release>1.7.9</release>
    <commit>f9f2e8c5f02fbcc93fe76b0d3cfa5d2089336849</commit>
    <build_time>2019-11-01T12:51:04Z</build_time>
    <golang>go1.11.5</golang>
    <compiler>gc</compiler>
    <os>linux</os>
    <architecture>amd64</architecture>
</Version>
```
