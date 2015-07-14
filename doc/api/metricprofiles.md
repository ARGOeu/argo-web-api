---
title: API documentation | ARGO
page_title: API - POEM Profiles 
font_title: 'fa fa-cogs'
description: API Calls for retrieving available POEM profiles
---

<a id="1"></a>

## GET: List Metric Profiles

This method can be used to retrieve a list of current metric (poem) profiles.

### Input

    /poems

#### Request headers

    x-api-key: "tenant_key_value"

### Response

Headers: `Status: 200 OK`

#### Response body

    <root>
      <Poem profile="ch.cern.sam.CLOUD-MON"/>
      <Poem profile="ch.cern.sam.GLEXEC"/>
      <Poem profile="ch.cern.sam.OPS_MONITOR"/>
      .
      .
    </root>


