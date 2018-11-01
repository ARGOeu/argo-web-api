# Recomputations Requests
 API Calls for listing existing and creating new recomputation requests

Name                                     | Description                                                                            | Shortcut
---------------------------------------- | -------------------------------------------------------------------------------------- | ------------------
GET: List Recomputation Requests         | This method can be used to retrieve a list of current Recomputation requests.          | [ Description](#1)
POST: Create a new recomputation request | This method can be used to insert a new recomputation request onto the Compute Engine. | [ Description](#2)
DELETE: Delete a specific recomputation  | This method can be used to delete a specific recomputation. | [ Description](#3)

<a id='1'></a>

## [GET]: List Recomputation Requests
This method can be used to retrieve a list of current Recomputation requests.

### Input

```
GET /recomputations
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
"root": [
     {
          "requester_name": "Arya Stark",
          "requester_email": "astark@shadowguild.com",
          "reason": "power cuts",
          "start_time": "2015-01-10T12:00:00Z",
          "end_time": "2015-01-30T23:00:00Z",
          "report": "EGI_Critical",
          "exclude": [
           "Gluster"
          ],
          "status": "running",
          "timestamp": "2015-02-01 14:58:40"
     }
 ]
}
```

Xml Response
```xml
<root>
    <Result>
        <requester_name>Arya Stark</requester_name>
        <requester_email>astark@shadowguild.com</requester_email>
        <reason>power cuts</reason>
        <start_time>2015-01-10T12:00:00Z</start_time>
        <end_time>2015-01-30T23:00:00Z</end_time>
        <report>EGI_Critical</report>
        <exclude>Gluster</exclude>
        <status>running</status>
        <timestamp>2015-02-01 14:58:40</timestamp>
    </Result>
</root>
```

<a id='2'></a>

## [POST]: Create a new recomputation request
This method can be used to insert a new recomputation request onto the Compute Engine.

### Input

```
POST /recomputations
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

#### Parameters

Type         | Description                                                                                                                                            | Required | Default value
------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------ | -------- | -------------
`start_time` | UTC time in W3C format                                                                                                                                 | YES      |
`end_time`   | UTC time in W3C format                                                                                                                                 | YES      |
`reason`     | Explain the need for a recomputation                                                                                                                   | YES      |
`report`     | Report for which the recomputation is requested                                                                                                        | YES      |
`exclude`    | Groups to be excluded from recomputation. If more than one group are to be excluded use the parameter as many times as needed within the same API call | NO       |

### Response
Headers: `Status: 201 Created`

<a id='3'></a>

## [DELETE]: Delete a specific recomputation

```
DELETE /recomputations/{ID}
```

#### Request headers

```
x-api-key: shared_key_value
Accept: application/json
```

### Response
`Status 200 OK`