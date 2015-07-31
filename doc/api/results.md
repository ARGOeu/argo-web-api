# API Calls

Name                                                                 | Description                                                                                                                                                                                                           | Shortcut
-------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -----------------------------
GET: List Availability and Reliability results for an endpoint group | This method retrieves the results of a specified endpoint group or multiple endpoint groups of a specific type that where computed based on a given report. Results can be retrieved on daily or monthly granularity. | [Description](#1)

<a id="1"></a>

# GET: List Availabilities and Reliabilities for Endpoint Groups

The following methods can be used to obtain a tenant's Availability and Reliability result metrics per Endpoint Group. The api authenticates the tenant using the api-key within the x-api-key header. User can specify time granularity (`monthly` or `daily`) for retrieved results and also format using the `Accept` header. Depending on the form of the request the user can request a single endpoint group results or a bulk of endpoint group results filtered by their type and if necessary their "top-level" group.

## Endpoint Groups

### Input

```
/results/{report_name}/{endpoint_group_type}?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{endpoint_group_type}/{endpoint_group_name}?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{group_type}/{group_name}/{endpoint_group_type}?[start_time]&[end_time]&[granularity]
or
/results/{report_name}/{group_type}/{group_name}/{endpoint_group_type}/{endpoint_group_name}?[start_time]&[end_time]&[granularity]
```

#### Query Parameters

Type            | Description                                                                                     | Required | Default value
--------------- | ----------------------------------------------------------------------------------------------- | -------- | -------------
`[start_time]`  | UTC time in W3C format                                                                          | YES      |
`[end_time]`    | UTC time in W3C format                                                                          | YES      |
`[granularity]` | Granularity of time that will be used to present data. Possible values are `monthly` or `daily` | NO       | `daily`

#### Path Parameters

Name                  | Description                                                                                                                         | Required | Default value
--------------------- | ----------------------------------------------------------------------------------------------------------------------------------- | -------- | -------------
`{report_name}`       | Name of the report that contains all the information about the profile, filter tags, group types etc.                               | YES      |
`{group_name}`        | Name of the Group of Endpoint Groups. If no name is specified then all Endpoint Groups regardless of top-level group are retrieved. | NO       |
`{group_type}`        | Type of the Group of Endpoint Groups. If no type is specified then all groups are retrieved.                                        | NO       |
`{endpoint_group_name}` | Name of the the Endpoint Group. If no name is specified then all groups are retrieved according to the `{endpoint_group_type}`.       | NO       |
`{endpoint_group_type}` | Type of the the Endpoint Group.                                                                                                     | YES      |

#### Headers

##### Request
```
x-api-key: "tenant_key_value"
Accept: "application/xml" or "application/json"
```

##### Response
```
Status: 200 OK
```


#### URL
`/api/v2/results/Report_A/SITE/ST01?start_time=2015-06-20T12:00:00Z&end_time=2015-06-23T23:00:00Z&granularity=daily`


#### Response Body

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
