---
id: consistency
title: API Data Consistency Check
---

Argo Monitoring uses automated checks that run at specific intervals and analyse the consistency of the monitorig data. If for example most of the monitored items appear to be problematic (flapping through states or appear as CRITICAL) that might indicate a network issue affecting Argo Monitoring. The automated checking mechanism will update the argo-web-api with the status of it's findings.There is also a manual mechanism where an argo-monitoring administrator can check the automated consistency warning and respond with an acknowledgement: either by maintaining the problematic state and confirming that there is an issue with the monitoring platform or - in the case that most of the monitoring items are truly experiecing issues idepedent from the state of the monitoring platform - revert the consistency warning and inform the users accordingly.



## GET Consistency information

__note: consistency information is both accessible from `/api/v2` and `/api/v3` paths. In the examples we are going to maintin the `/api/v3`__

A user can get consistency information by issuing the following call:

```
GET /api/v3/consistency
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
  "message": "Flapping items percentage 2% (2/100)"
}
```

## GET Consistency information with verbosity

A user can get consistency information with more details (verbose) by using the url parameter `?verbose` in the previous call:

```
GET /api/v3/consistency?verbose
```

### Request headers

```
Accept: application/json
```

### Response
Headers: `Status: 200 OK`

### Response Body

User can be informed with details about the automated check (when did it run recently, it's status and it's message)

Json Response example:
```json
{
  "status": "CRITICAL",
  "timestamp": "2025-07-10T05:44:00Z",
  "message": "Flapping items percentage 65% (65/100)",
  "auto_check_status": "CRITICAL",
  "auto_check_message": "Flapping items percentage 65% (65/100)",
  "auto_check_timestamp": "2025-07-10T02:30:00Z",
}
```

If there is an admin acknowledgement the user can be informed about it also (in verbose mode)

Json Response example:
```json
{
  "status": "OK",
  "timestamp": "2025-07-10T12:44:00Z",
  "message": "Monitoring is ok - most items are flapping due to issues of their own",
  "auto_check_status": "CRITICAL",
  "auto_check_message": "Flapping items percentage 65% (65/100)",
  "auto_check_timestamp": "2025-07-10T12:30:00Z",
  "ack_check_status": "OK",
  "ack_check_message": "Monitoring is ok - most items are flapping due to issues of their own",
  "ack_check_timestamp": "2025-07-10T06:30:00Z",
  "ack_timeout_hours":6
}
```

So with verbose mode the user can have fully awareness of the consistency status and how it is affected both by the automated check and from the admin acknowledgements. 

## POST Consistency results as automated check

An automated consistency check can post results by having the correct token (with the correct role of `consistency-check`) and by issuing the following request:

```
POST /api/v3/consistency/auto-check
```

### Request headers

```
Accept: application/json
```

### Post body
```json
{
    "status": "OK",
    "message": "All items are ok"
}

```

### Response

The automated check needs only to specify a state (`OK`, `CRITICAL`) and an accompanying message. If the request is succesfull the web-api will respond:

Json Response:
```json
{
  "message": "The Auto Check event was posted succesfully",
  "code": 200
}
```

## POST Acknlowledgement of consistency issue

An user with the correct role (`consistency-ack`) can acknowledge consistency issues raised by the automated mechanism and respond accordingly by issuing the following request:

```
POST /api/v3/consistency/ack
```

### Request headers

```
Accept: application/json
```

### Post body
```json
{
    "status": "OK",
    "message": "False alarm! Monitoring platform is ok"
}
```

Optionally the user can specify an ack timeout in hours (default: `6`) as such:

```json
{
    "status": "OK",
    "message": "False alarm! Monitoring platform is ok",
    "timeout_hours": 2
}
```

### Response

If the request is succesfull the web-api will respond:

Json Response:
```json
{
  "message": "The Ack event was posted succesfully",
  "code": 200
}
```