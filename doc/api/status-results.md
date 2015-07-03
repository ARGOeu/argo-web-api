---
title: API documentation | ARGO
page_title: API - Status results
font_title: 'fa fa-cogs'
description: API calls for retrieving monitoring status results
---

## API Calls

| Name  | Description | Shortcut |
| GET: List Service Metric Status Timelines | This method may be used to retrieve a specific service metric status timeline (applies on a specific host endpoint and a specific service flavor).|<a href="#1"> Description</a>|
| GET: List Service Endpoint Status Timelines | This method may be used to retrieve a specific service endpoint status timeline (applies on a specific service flavor). | <a href="#2"> Description</a>|
| GET: List Service Flavor Status Timelines |This method may be used to retrieve a specific service flavor status timeline (applies for a specific site). | <a href="#3"> Description</a>|
| GET: List Site Status timelines| This method may be used to retrieve a whole site status timeline. | <a href="#4"> Description</a>|


<a id="1"></a>

## GET: List Service Metric Status Timelines

This method may be used to retrieve a specific service metric status timeline (applies on a specific host endpoint and a specific service flavor).

### Input

    /status/metrics/timeline/{group}?[start_time]&[end_time]&[job]&[group_type]

#### Parameters

| Type | Description | Required | Default value |
|`start_time`| UTC time in W3C format| YES |  |
|`end_time`| UTC time in W3C format| YES |  |
|`job`| Job (view) name | YES |  |
|`group_type`| `ngi` or `site` | YES |  |

Depending on the `group_type`, `{group}` is the name of the group (for example `NGI_GRNET` when `group_type=ngi`, `HG-03-AUTH` when `group_type=site`). 

#### Request headers

    x-api-key: "tenant_key_value"


### Response

Headers: `Status: 200 OK`

#### Response body

##### `group_type=ngi`

    <root>
       <job name="ROC_CRITICAL">
         <group name="NGI_GRNET" type="ngi">
           <group name="HG-03-AUTH" type="site">
             <group name="CREAM-CE" type="service_type">
               <host name="cream01.afroditi.gr">
                 <metric name="emi.cream.CREAMCE-JobSubmit">
                   <status timestamp="2015-04-30T23:59:00Z" status="OK"></status>
                   <status timestamp="2015-05-01T01:00:00Z" status="CRITICAL"></status>
                   <status timestamp="2015-05-01T02:00:00Z" status="OK"></status>
                 </metric>
               </host>
             </group>
           </group>
         </group>
       </job>
     </root>


##### `group_type=site`

    <root>
       <job name="ROC_CRITICAL">
         <group name="NGI_GRNET" type="ngi">
           <group name="HG-03-AUTH" type="site">
             <group name="CREAM-CE" type="service_type">
               <host name="cream01.afroditi.gr">
                 <metric name="emi.cream.CREAMCE-JobSubmit">
                   <status timestamp="2015-04-30T23:59:00Z" status="OK"></status>
                   <status timestamp="2015-05-01T01:00:00Z" status="CRITICAL"></status>
                   <status timestamp="2015-05-01T02:00:00Z" status="OK"></status>
                 </metric>
               </host>
             </group>
           </group>
         </group>
       </job>
     </root>
    <root>





<a id="2"></a>

## GET: List Service Endpoint Status Timelines

This method may be used to retrieve a specific service endpoint status timeline (applies on a specific service flavor).

### Input

    /status/endpoints/timeline/{hostname}/{service_flavor}?[start_time]&[end_time]&[job]

#### Parameters

| Type | Description | Required | Default value |
|`start_time`| UTC time in W3C format| YES | |
|`end_time`| UTC time in W3C format| YES| |
|`job`| Job (view) name | YES |  |

#### Request headers

    x-api-key: "tenant_key_value"

### Response

Headers: `Status: 200 OK`

#### Response body


    <root>
       <job name="JOB_A">
         <endpoint hostname="irods01.juelich.de" service="iRods">
           <status timestamp="2015-05-01T00:00:10Z" status="OK"></status>
           <status timestamp="2015-05-01T01:01:00Z" status="CRITICAL"></status>
         </endpoint>
       </job>
     </root>
    <root>




<a id="3"></a>

## GET: List Service Flavor Status Timelines (TBD)

This method may be used to retrieve a specific service flavor status timeline (applies for a specific site).

### Input

    /status/services/timeline/{group}?[start_time]&[end_time]&[vo]&[profile]&[group_type]

#### Parameters

| Type | Description | Required | Default value |
|`start_time`| UTC time in W3C format| YES ||
|`end_time`| UTC time in W3C format| YES| |
|`vo`| vo name | NO | `ops` |
|`profile`| POEM profile name | NO | `ch.cern.sam.ROC_CRITICAL` |
|`group_type`| `site` or `ngi` | NO | `site` |

Depending on the `group_type`, `{group}` is the name of the group (for example `HG-03-AUTH` when `group_type=site` or `NGI_GRNET` when `group_type=ngi`).

### Response

Headers: `Status: 200 OK`

#### Response body

    <root>
      <profile name="A_POEM">
        <flavor name="A_FLAVOR" site="A_SITE-NAME" vo="A_VO" roc="A_ROC" monitoring_host="A_MONHOST">
          <timeline start_time="2014-10-23T00:00:00Z" end_time="2014-10-24T00:00:00Z">
            <status timestamp="2014-10-23T00:12:34Z" value="OK" />
            <status timestamp="2014-10-23T01:12:20Z" value="WARNING" />
            <status timestamp="2014-10-23T02:12:31Z" value="CRITICAL" />
            <status timestamp="2014-10-23T04:12:25Z" value="OK" />
            .
            .
            .
            <status timestamp="2014-10-23T23:17:45Z" value="OK" />
          </timeline>
        </flavor>
      </profile>
    </root>


<a id="4"></a>

## GET: List Metric Results

This method may be used to retrieve a detailed metric result.

### Input

    /status/metrics/msg/{hostname}/{service_flavor}/{metric}?[exec_time]&[job]


#### Parameters

| Type | Description | Required | Default value |
|`exec_time`| UTC time in W3C format| YES |  |
|`job`| Job (view) name | YES |  |

#### Request headers

    x-api-key: "tenant_key_value"

### Response

Headers: `Status: 200 OK`

#### Response body

    <root>
       <job name="ROC_CRITICAL">
         <group name="NGI_GRNET" type="ngi">
           <group name="HG-03-AUTH" type="site">
             <group name="CREAM-CE" type="service_type">
               <host name="cream01.afroditi.gr">
                 <metric name="emi.cream.CREAMCE-JobSubmit">
                   <status timestamp="2015-05-01T01:00:00Z" status="">
                     <summary>Cream status is CRITICAL</summary>
                     <message>Cream job submission test failed!</message>
                   </status>
                 </metric>
               </host>
             </group>
           </group>
         </group>
       </job>
     </root>
    <root>

