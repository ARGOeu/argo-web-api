---
title: API documentation | ARGO
page_title: API - Recalculation Requests
font_title: 'fa fa-cogs'
description: API Calls for listing and creating new recalculation requests
---


## API Calls

| Name  | Description | Shortcut |
| GET: List Recalculation Requests| This method can be used to retrieve a list of current Recalculation requests. |<a href="#1"> Description</a>|
| POST: Create a new recalculation request | This method can be used to insert a new recalculation request onto the Compute Engine. | <a href="#2"> Description</a>|

<a id='1'></a>

## GET: List Recalculation Requests

This method can be used to retrieve a list of current Recalculation requests.

### Input

    /recomputations

### Response

Headers: `Status: 200 OK`

#### Response body

    <root>
      <Request start_time="2013-12-08T12:03:44Z" end_time="2013-12-10T12:03:44Z" reason="some_reason" ngi_name="Some_NGI" status="pending" timestamp="2014-03-07 12:03:44">
        <Exclude site="site_1"/>
        <Exclude site="site_2"/>
      </Request>
    </root>

<a id='2'></a>

## POST: Create a new recalculation request

This method can be used to insert a new recalculation request onto the Compute Engine.

### Input

    /recomputations

#### Request headers

    x-api-key: shared_key_value
    Content-Type: application/json

#### Parameters

| Type | Description | Required | Default value | 
|`start_time`| UTC time in W3C format | YES | |
|`end_time`| UTC time in W3C format | YES | |
|`reason`| Explain the need for a recalculation | YES | |
|`ngi_name`| NGI for which the recalculation is requested | YES | |
|`exclude_site`| Site to be excluded from recalculation. If more than one sites are to be excluded use the parameter as many times as needed within the same API call | NO | |

### Response

Headers: `Status: 200 OK`


