---
title: API documentation | ARGO
page_title: API - Availability Profiles 
font_title: 'fa fa-cogs'
description: API Calls for managing Availability Profiles
---

## API Calls

| Name  | Description | Shortcut |
| GET: List Availability Profiles | This method can be used to retrieve a list of current Availability profiles used by the Compute Engine. |<a href="#1"> Description</a>|
| POST: Create a new availability profile |This method can be used to insert a new Availability Profile onto the Compute Engine. | <a href="#2"> Description</a>|
| PUT: Update an existing Availability Profile |This method can be used to update an existing Availability Profile onto the Compute Engine.|<a href="#3"> Description</a>|
| DELETE: Delete an availability profile  | This method can be used to delete an existing Availability Profile.| <a href="#4"> Description</a>|



<a id="1"></a>

## GET: List Availability Profiles

This method can be used to retrieve a list of current Availability profiles used by the Compute Engine.

### Input

    /AP?[name]&[namespace]

#### Parameters

| Type | Description | Required | Default value |
|`name`| Profile name (both name and namespace are needed) | YES| |
|`namespace`| Profile namespace (both name and namespace are needed) | YES| |


### Response

Headers: `Status: 200 OK`

#### Response body

    <root>
      <profile id="some_id" name="name" namespace="namespace" poems="poem_name">
        <AND>
          <OR>
            <Group service_flavor="A_FLAVOR"/>
            <Group service_flavor="B_FLAVOR"/>
            <Group service_flavor="C_FLAVOR"/>
            <Group service_flavor="D_FLAVOR"/>
            <Group service_flavor="E_FLAVOR"/>
          </OR>
          <OR>
            <Group service_flavor="F_FLAVOR"/>
            <Group service_flavor="G_FLAVOR"/>
          </OR>
          <OR>
            <Group service_flavor="H_FLAVOR"/>
          </OR>
        </AND>
      </profile>
      .
      .
    </root>


<a id="2"></a>

## POST: Create a new availability profile

This method can be used to insert a new Availability Profile onto the Compute Engine.

### Input

    /AP

#### Parameters (Json formatted string)

    { 
      "name" : <name>, 
      "namespace" :<namespace>, 
      "poems" : [<POEM profile>],
      "groups" : [ [<flavor1>, <flavor2>,... ],...,[<flavorN>] ]
    }

| Type | Description | Required | Default value |
|`name`| Profile name (both name and namespace are needed) | YES| |
|`namespace`| Profile namespace (both name and namespace are needed) | YES| |
|`poems`| Low level profile matching checks with service flavors | YES| |
|`groups`|  Group of service flavors groups  | YES| |

#### Request headers

    x-api-key: "shared_key_value"
    Content-Type: application/json

#### Example 

This example creates a new availability profile `ap-test-2` for the `someRI` tenant. The fictional service flavors `Torque`, `Slurm` and `Moab` are logically **ORed** in the calculation of this service group as they more or less represent a similar functionality (lets name this **the compute resource scheduling capability of the RI**). Likewise the fictional `iRods` and `GridFTP` service flavors are logically **ORed** in the calculation as they respresent in this context a **storage capability** of the given RI. Finally a valuable resource that is taken in consideration is the fictional IS service flavor which provides the **information service capability** for the given RI. There is only one flavor of that kind on the whole RI Tier hence it makes sense to leave it alone in the A/R calculation. The overall A/R metrics are deduced by grouping these three capabilities ( compute resource scheduling, storage and information service) with a logical **AND** operation. The check (probe) results that are taken into consideration within this calculation are defined within the poem profile, which is extracted via a connector from the POEM endpoint of the RI. 

    {  
      "name":"ap-test-2",
      "namespace":"someRI",
      "groups":[  
        [  
          "Torque",
          "Slurm",
          "Moab"
        ],
        [  
          "iRods",
          "GridFTP"
        ],
        [  
          "IS"
        ]
      ],
      "poems":[  
        "someRI.root.CRITICAL"
      ]
    }

### Response

Headers: `Status: 200 OK`

#### Response body

##### Successful creation of profile

    <root>
       <Message>Availability Profile record successfully created</Message>
    </root>

##### Profile exists

    <root>
       <Message>An availability profile with that name already exists</Message>
    </root>

<a id="3"></a>

## PUT: Update an existing Availability Profile

This method can be used to update an existing Availability Profile onto the Compute Engine.

### Input

    /AP/{id}

#### Parameters ( Json formatted string) 

    { 
      "name" : <name>, 
      "namespace" :<namespace>, 
      "poems" : [<POEM profile>],
      "groups" : [ [<flavor1>, <flavor2>,... ],...,[<flavorN>] ]
    }


##### Input parameters

| Type | Description | Required | Default value |
| id | The availability profile id | YES | | 


##### JSON parameters

| Type | Description | Required | Default value |
|`name`| Profile name (both name and namespace are needed) | YES| |
|`namespace`| Profile namespace (both name and namespace are needed) | YES| |
|`poems`| Low level profile matching checks with service flavors | YES| |
|`groups`|  Group of service flavors groups  | YES| |

Check the POST: Create a new availability profile example. 

#### Request headers

    x-api-key: "shared_key_value"
    Content-Type: application/json

### Response

Headers: `Status: 200 OK`

#### Response body

##### Successful update

    <root>
       <Message>Update successful</Message>
    </root>

##### Profile does not exist

    <root>
       <Message>No profile matching the requested id</Message>
    </root>


<a id="4"></a>

## DELETE: Delete an availability profile

This method can be used to delete an existing Availability Profile.

### Input

    /AP/{id}

#### Parameters

| Type | Description | Required | Default value |
| id | The availability profile id | YES | | 

#### Request headers

    x-api-key: "shared_key_value"

### Response

Headers: `Status: 200 OK`

#### Response body

##### Successful delete

    <root>
       <Message>Delete successful</Message>
    </root>

##### Profile does not exist

    <root>
       <Message>No profile matching the requested id</Message>
    </root>


