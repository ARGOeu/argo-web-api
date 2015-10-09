# API - Availabilities & Reliabilities

API calls for retrieving computed Availabilities and Reliabilities

| Name  | Description | Shortcut |
--------|-------------|----------|
| GET: List Availabilities and Reliabilities for Groups | The following methods can be used to obtain Availability and Reliablity metrics per group type elements (i.e. Endpoint Groups, Group of Endpoint Groups etc). Results can be retrieved on daily or monthly granularity.  |<a href="#1"> Description</a>|
| GET: List Availabilities and Reliabilities for Service Flavors | This method can be used to obtain Availability and Reliability metrics for Service Flavors per Site. Results can be retrieved on daily or monthly granularity. | <a href="#2"> Description</a>|

<a id="1"></a>

## [GET]: List Availabilities and Reliabilities for Endpoint Groups, Group of Endpoint Groups

The following methods can be used to obtain Availability and Reliability metrics per group type elements (i.e. Endpoint Groups, Group of Endpoint Groups etc). Results can be retrieved on daily or monthly granularity.

### Endpoint Groups

#### Input

Endpoint Groups

    /endpoint_group_availability?[start_time]&[end_time]&[job]&[granularity]&[format]&[group_name]&[supergroup_name]

##### Parameters

| Type | Description | Required | Default value |
-------|-------------|----------|---------------|
|`start_time`| UTC time in W3C format| YES ||
|`end_time`| UTC time in W3C format| YES| |
|`job`| Name of the job that contains all the information about the profile, filter tags etc. | YES| |
|`format`| Only xml available right now, so the parameter is void thus deactivated for the time being  | NO| `XML` |
|`group_name`| Name of the Endpoint Groups. If no name is specified then all Endpoint Groups are retrieved. |NO| |
|`supergroup_name`| Name of the group that groups the Endpoint Groups. If no name is specified then all groups are retrieved. |NO| |

##### Request headers

    x-api-key: "tenant_key_value"

#### Response

Headers: `Status: 200 OK`

##### Response body for `/endpoint_group_availability` API call

    <root>
      <Job name="Job_A">
        <EndpointGroup name="Site-Name" SuperGroup="SuperGroup-A">
          <Availability timestamp="YYYY-MM-DD" availability="1" reliability="1"/>
          <Availability timestamp="YYYY-MM-DD" availability="1" reliability="1"/>
        </EndpointGroup>
        <EndpointGroup name="Site-Name-Another" SuperGroup="SuperGroup-B">
          <Availability timestamp="YYYY-MM-DD" availability="1" reliability="1"/>
          <Availability timestamp="YYYY-MM-DD" availability="1" reliability="1"/>
        </EndpointGroup>
        .
        .
        .
      </Job>
    </root>

### Group of Endpoint Groups

#### Input

Group of Endpoint Groups

    /group_groups_availability?[start_time]&[end_time]&[job]&[granularity]&[format]&[group_name]

##### Parameters

| Type | Description | Required | Default value |
-------|-------------|----------|---------------|
|`start_time`| UTC time in W3C format| YES ||
|`end_time`| UTC time in W3C format| YES| |
|`job`| Name of the job that contains all the information about the profile, filter tags etc. | YES| |
|`format`| Only xml available right now, so the parameter is void thus deactivated for the time being  | NO| `XML` |
|`group_name`| Name of the group that groups the Endpoint Groups. If no name is specified then all groups are retrieved. |NO| |

##### Request headers

    x-api-key: "tenant_key_value"

#### Response

Headers: `Status: 200 OK`

##### Response body for `/group_groups_availability` API call

    <root>
      <Job name="Job_A">
        <SuperGroup name="GROUP_A">
          <Availability timestamp="2013-08-01" availability="87.64699776723847" reliability="87.64699776723847">
          </Availability>
          <Availability timestamp="2013-08-02" availability="87.63642636198455" reliability="87.63642636198455">
          </Availability>
          <Availability timestamp="2013-08-03" availability="11.307937916910474" reliability="11.307937916910474">
          </Availability>
          .
          .
          .
          <Availability timestamp="2013-08-09" availability="87.69028873148349" reliability="92.9126400880786">
          </Availability>
        </SuperGroup>
      </Job>
    </root>


<a id="2"></a>

## [GET]: List Availabilities and Reliabilities for Service Flavors

This method can be used to obtain Availability and Reliability metrics for Service Flavors. Results can be retrieved on daily or monthly granularity.

### Input

    /service_flavor_availability?[start_time]&[end_time]&[job]&[granularity]&[format]&[flavor]&[supergroup]


##### Parameters

| Type | Description | Required | Default value |
-------|-------------|----------|---------------|
|`start_time`| UTC time in W3C format| YES ||
|`end_time`| UTC time in W3C format| YES| |
|`job`| Job (view) according to which A/R has been calculated | YES | |
|`granularity`| Possible values: `daily`, `monthly` | NO | `daily` |
|`format`| Only xml available right now, so the parameter is void thus deactivated for the time being  | NO | `xml` |
|`flavor`| Service Flavor name or list of Service Flavors. If no Service Flavor is specified then all flavors within the supergroup are retrieved. | NO | All flavors within the given supergroup |
|`supergroup`| Name of supergroup. If no supergroup is specified then all available service flavor results are returned | NO | all availabile results |

##### Request headers

    x-api-key: "tenant_key_value"

#### Response

Headers: `Status: 200 OK`

##### Response body

    <root>
      <Job name="EGI_Critical">
        <SuperGroup name="GR-01-AUTH">
          <Flavor Flavor="SRMv2">
            <Availability timestamp="2015-06-20" availability="100" reliability="100"></Availability>2"/>
            <Availability timestamp="2015-06-21" availability="100" reliability="100"></Availability>
            <Availability timestamp="2015-06-22" availability="100" reliability="100"></Availability>
            <Availability timestamp="2015-06-23" availability="100" reliability="100"></Availability>2"/>
            <Availability timestamp="2015-06-24" availability="100" reliability="100"></Availability>
            <Availability timestamp="2015-06-25" availability="100" reliability="100"></Availability>
            <Availability timestamp="2015-06-26" availability="100" reliability="100"></Availability>
          </Flavor>
          <Flavor Flavor="Site-BDII">2"/>
            <Availability timestamp="2015-06-20" availability="100" reliability="100"></Availability>
            <Availability timestamp="2015-06-21" availability="50" reliability="100"></Availability>
            <Availability timestamp="2015-06-22" availability="50" reliability="100"></Availability>
            <Availability timestamp="2015-06-23" availability="50" reliability="60"></Availability>
            <Availability timestamp="2015-06-24" availability="50" reliability="60"></Availability>
            <Availability timestamp="2015-06-25" availability="50" reliability="60"></Availability>
            <Availability timestamp="2015-06-26" availability="100" reliability="100"></Availability>
          </Flavor>
        </SuperGroup>
        <SuperGroup name="HG-03-AUTH">
          <Flavor Flavor="CREAM-CE">
            <Availability timestamp="2015-06-20" availability="100" reliability="100"></Availability>
            <Availability timestamp="2015-06-21" availability="100" reliability="100"></Availability>
            <Availability timestamp="2015-06-22" availability="100" reliability="100"></Availability>
            <Availability timestamp="2015-06-23" availability="100" reliability="100"></Availability>
          </Flavor>
          <Flavor Flavor="SRMv2">
            <Availability timestamp="2015-06-20" availability="100" reliability="100"></Availability>
            <Availability timestamp="2015-06-21" availability="100" reliability="100"></Availability>
            <Availability timestamp="2015-06-22" availability="100" reliability="100"></Availability>
            <Availability timestamp="2015-06-23" availability="100" reliability="100"></Availability>
          </Flavor>
        </SuperGroup>
      </Job>
    </root>





