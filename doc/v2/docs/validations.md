

# API Validations

## Parameters validations

Name          | Affected Resources      | Verb   | Shortcut         
------------- | ----------------------- | ------ | -----------------
`start_time`  | `/results`, `/status`   | GET    | [Description](#1)
`end_time`    | `/results`, `/status`   | GET    | [Description](#2)
`exec_time`   | `/metric_result`        | GET    | [Description](#3)
`granularity` | `/results`              | GET    | [Description](#4)


<a id='1'></a>

### `start_time`

The `start_time` query parameter is used under the `/results` and `/status` resources to define the starting time of the query used to match A/R or status results respectively. The value should be given in zulu format like so: `2006-01-02T15:04:05Z`.

In case the parameter is not provided the following response is returned

```json
{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "start_time not set",
   "code": "400",
   "details": "Please use start_time url parameter in zulu format (like 2006-01-02T15:04:05Z) to indicate the query start time"
  }
 ]
}
```

In case the parameter value is malformed (not in zulu expected format) the following response is returned:

```json
{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "start_time parsing error",
   "code": "400",
   "details": "Error parsing date string %s please use zulu format like 2006-01-02T15:04:05Z"
  }
 ]
}
```

<a id='2'></a>

### `end_time`

The `end_time` query parameter is used under the `/results` and `/status` resources to define the ending time of the query used to match A/R or status results respectively. The value should be given in zulu format like so: `2006-01-02T15:04:05Z`.

In case the parameter is not provided the following response is returned

```json
{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "end_time not set",
   "code": "400",
   "details": "Please use end_time url parameter in zulu format (like 2006-01-02T15:04:05Z) to indicate the query end time"
  }
 ]
}
```

In case the parameter value is malformed (not in zulu expected format) the following response is returned:

```json
{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "end_time parsing error",
   "code": "400",
   "details": "Error parsing date string %s please use zulu format like 2006-01-02T15:04:05Z"
  }
 ]
}
```

#### No time span set

In case neither the `start_time` nor the `end_time` parameters are defined the following response is returned by the api:

```json
{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "No time span set",
   "code": "400",
   "details": "Please use start_time and end_time url parameters to set the prefered time span"
  }
 ]
}
```

<a id='3'></a>

### `exec_time`

The `exec_time` query parameter is used under the `/metric_result` resource to define the execution time of the metric result to fetch from the datastore. The value should be given in zulu format like so: `2006-01-02T15:04:05Z`.

In case the parameter is not provided the following response is returned

```json
{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "exec_time not set",
   "code": "400",
   "details": "Please use exec_time url parameter in zulu format (like 2006-01-02T15:04:05Z) to indicate the exact probe execution time"
  }
 ]
}
```

In case the parameter value is malformed (not in zulu expected format) the following response is returned:

```json
{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "exec_time parsing error",
   "code": "400",
   "details": "Error parsing date string %s please use zulu format like 2006-01-02T15:04:05Z"
  }
 ]
}
```

<a id='4'></a>


### `granularity`

The `granularity` query parameter is used *optionally* under the `/results` resource to indicate the granularity level. It's value may be either monthly or daily. If not set by the user `monthly` is used.

In case the parameter value is malformed (neither `daily` nor `monthly`) the following response is returned:

```json
{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "Wrong Granularity",
   "code": "400",
   "details": "%s is not accepted as granularity parameter, please provide either daily, monthly or custom"
  }
 ]
}
```

The Granularity parameter is only relevant for a/r result requests. For status requests granularity is not supported. In order to avoid confusion, if a user provides granularity parameter during status requests the following response is returned:
```json
{
 "status": {
  "message": "Bad Request",
  "code": "400"
 },
 "errors": [
  {
   "message": "Granularity parameter should not be used in status results",
   "code": "400",
   "details": "Granularity parameter is valid only for a/r result requests, not for status results"
  }
 ]
}
```

## Headers validations

Name          | Affected Resources      | Shortcut         
------------- | ----------------------- | -----------------
`Accept`      | All                     | [Description](#5)

<a id='5'></a>

### `Accept`

The `Accept` header is compulsory to use under all api resources. Its value may be either `application/json` or `application/xml`.

In case the Accept header is not provided or is provided but is malformed the following error response is returned by the api:

```json
{
 "status": {
  "message": "Not Acceptable Content Type",
  "code": "406",
  "details": "Accept header provided did not contain any valid content types. Acceptable content types are 'application/xml' and 'application/json'"
 }
}
```
