---
title: 'API v2.0 documentation | ARGO'
page_title: 'API - Availability & Reliability Results'
font_title: fa fa-cogs
description: API calls for retrieving computed Availability and Reliability Results
---

# API Calls

Name                                                           | Description                                                                                                                                                                                                             | Shortcut
-------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -----------------------------
GET: List Availabilities and Reliabilities for Endpoint Groups | The following methods can be used to obtain Availability and Reliablity metrics per group type elements (i.e. Endpoint Groups, Group of Endpoint Groups etc). Results can be retrieved on daily or monthly granularity. | <a href="#1"> Description</a>

<a id="1"></a>

# GET: List Availabilities and Reliabilities for Endpoint Groups
The following methods can be used to obtain Availability and Reliability metrics per group type elements (i.e. Endpoint Groups, Group of Endpoint Groups etc). Results can be retrieved on daily or monthly granularity.

## Endpoint Groups
### Input
Endpoint Groups

```
/results/{report_name}/{lesser_group_type}?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{lesser_group_type}/{lesser_group_name}?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{group_type}/{lesser_group_type}?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{group_type}/{lesser_group_type}/{lesser_group_name}?[start_time]&[end_time]&[granularity]
```

#### Parameters

Type                  | Description                                                                                               | Required | Default value
--------------------- | --------------------------------------------------------------------------------------------------------- | -------- | -------------
`[start_time]`        | UTC time in W3C format                                                                                    | YES      |
`[end_time]`          | UTC time in W3C format                                                                                    | YES      |
`[report]`            | Name of the report that contains all the information about the profile, filter tags etc.                  | YES      |
`{group_name}`        | Name of the Endpoint Groups. If no name is specified then all Endpoint Groups are retrieved.              | NO       |
`{group_type}`        | Name of the group that groups the Endpoint Groups. If no name is specified then all groups are retrieved. | NO       |
`{lesser_group_type}` | Name of the group that groups the Endpoint Groups. If no name is specified then all groups are retrieved. | NO       |
`{lesser_group_name}` | Name of the group that groups the Endpoint Groups. If no name is specified then all groups are retrieved. | NO       |

#### Request headers

```
x-api-key: "tenant_key_value"
format: "xml or json"
```

### Response
Headers: `Status: 200 OK`

#### Response body for `/api/v2/results/Report_A/SITE/ST01?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily` API call

```
<root>
<group name="GROUP_A" type="GROUP">
<group name="ST01" type="SITE">
  <results timestamp="2015-06-22" availability="66.7" reliability="54.6" unknown="0" uptime="1" downtime="0"></results>
  <results timestamp="2015-06-23" availability="100" reliability="100" unknown="0" uptime="1" downtime="0"></results>
</group>
</group>
</root>
```

## Group of Endpoint groups
