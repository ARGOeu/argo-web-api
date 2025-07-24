# Changelog

All notable changes in argo-web-api project are documented here

## 1.13.3 - (2023-05-10)

### Added:
- ARGO-4207 provide automatically combined endpoints in topology of combined tenants
- ARGO-4208 Provide automatically combined groups in topology of combined tenants
- ARGO-4209 Provide automatically combined service types in topology of combined tenants

## 1.13.2 - (2023-02-28)

### Added:
- ARGO-4206 Add filter by service type when requesting metric results 

### Fixed:
- ARGO-4226 In recomputations fix support for exclude monitoring source


## 1.13.1 - (2023-02-02)


### Added:
- ARGO-4185 Add title field to service types
- ARGO-4096 Daily issues per group

### Fixed:
- ARGO-4124 Fix regression issue with missing date field in supergroup monthly query


## 1.13.0 - (2022-11-07)

### Added:
- ARGO-4089 support OPTIONS for api v3 calls using resource-ids
- ARGO-4048 add additional tags field list of service types
- ARGO-4078 Implement api/v3 call to serve endpoint ar results by resource id
- ARGO-4026 api v3 call to get status results by resource id

### Changed:
- ARGO-3960 Update to docusaurus v.2.0. Add search plugin. Spellcheck docs

### Fixed:
- ARGO-4087 Fix date validation in api/v3 A/R calls
- ARGO-4084 Fix v3 ar results by resource-id call



## 1.12.0 - (2022-06-09)

### Added:
- ARGO-3874 Add info on the report about what is computed (ar,status,trends)
- ARGO-3750 Implement api v3 call to list status timelines for servicegroups and endpoints
- ARGO-3702 Add list of available services types per tenant
- ARGO-3645 Serve a/r results for top level groups in v3 of api

### Fixed:
- ARGO-3690 Provide proper not found response when requesting details about a metric results that doesn't exist

### Changed: 
- ARGO-3850 Fix status v3 grouping issue for endpoints that belong to multiple groups
- ARGO-3802 Display api v3 status results in sorted order. Add parameter for display latest status results. Make start end time optional


## 1.11.0 - (2022-04-06)

### Added:
- ARGO-3408 Serve threshold flag trough web-api
- ARGO-3532 Serve tag trends for flapping items
- ARGO-3533 Serve status metric trends by tag
- ARGO-3534 Support exclude metrics in recomputations
- ARGO-3458 Serve metric details per report
- ARGO-3377 Add negative operator in topology filtering
- ARGO-3432 Show threshold rule applied in status metric timelines
- ARGO-3433 display threshold result rules in low level metric results
- ARGO-3423 Add list of metric to tags mappings
- Support combined tenant data feeds
- ARGO-3320 Provide info url when available for host through metrics results
- ARGO-3281 create status trends for endpoint groups
- ARGO-3280 create status trends view for services
- ARGO-3279 Create Endpoint status trends view
- ARGO-3280 create status trends view for services
- ARGO-3279 Create Endpoint status trends view

### Fixed:
- ARGO-3657 Fix docs to correctly display metric/tags docs in sidebar
- ARGO-3452 Fix threshold rule applied field value in status metrics
- ARGO-3269 Fix empty downtime response

### Changed:
- changed events with duration_in_minutes in trends section
- ARGO-3581 Remove support for TLS 1.0 and 1.1


## 1.10.0 - 2021-09-02

### Added:
- ARGO-3197 Display status trend results for metrics
- ARGO-3193 Add monthly granularity option over flapping trends
- ARGO-3186 Show flapping trends over range of time
- ARGO-3077 Display group flapping trends
- ARGO-3076 Display service flapping trends
- ARGO-3075 Display endpoint flapping trends 
- ARGO-3074 Display metric flapping trends 
- ARGO-2885 Add notification information to topology items
- ARGO-2724 Introduce topology tags/values method 

### Fixed:
- ARGO-3269 Fix empty downtime response
- ARGO-3241 close mongodb sessions in recomputation handlers
- ARGO-3225 Fix endpoint group type filtering by report
- ARGO-3161 Fix issue with multiple values filtering in tags

### Changed:
- ARGO-3265 Return exact date downtime instead of close to date
- ARGO-2880 Connect specific weight datasets to each report


## 1.9.2 - (2020-12-16)

### Added:
- ARGO-2703 Implement filter recomputations by date and report
- ARGO-2630 Implement weights feed resource
- ARGO-2558 Show list of problematic endpoints
- ARGO-2594 Get user information
- ARGO-2593 List tenant users
- ARGO-2475 Remove user
- ARGO-2476 Refresh User's token
- ARGO-2474 Update user in tenant
- ARGO-2473 Tenant create user
- ARGO-2509 Provide topology feed parameters
- ARGO-2503 Provide postman tests for results

### Fixed:
- ARGO-2690 Minor fixes and cleanups in weights resource
- ARGO-2599 Update to latest docusaurus Minor fixes in topology/downtime docs
- ARGO-2572 Fix feeds routing issue

### Changed:
- ARGO-2651 Remove name,id namespacing from daily downtimes
- ARGO-2570 Migrate argo-web-api docs to docusaurus2


## 1.9.1 - (2020-07-08)

### Added:
- ARGO-2305 Flat list of status results by metric type
- ARGO-2228 Add flat list for status endpoint results
- ARGO-2227 Feature: argo-web-api provide flat endpoint list a/r
- ARGO-2250 Display extra information from endpoint topology in endpoint status results
- ARGO-2249 Add description field to metric profiles
- ARGO-2226 Display additional endpoint information when available in endpoint a/r results
- ARGO-2238 Support composite filters in topology

### Fixed: 
- ARGO-2418 Fix threshold profile update name unique check
- ARGO-2274 Fix date input validation checks in historic profiles
- DEVOPS-111 Add clean workspace post build step
- DEVOPS-54 Build each stage in different container
- DEVOPS-67 Jenkinsfile environmental variable is not being interpreted
- DEVOPS-59 Generate argo-web-api swagger docs

## 1.9.0 - (2020-03-26)

### Added:
- ARGO-2187 Use report topology to filter endpoint items 
- ARGO-2188 Use report's filter tags field to filter endpoint topology by tags 
- ARGO-2234 Apply  sort order by id when listing multiple historic profiles
- ARGO-2192 Use report's filter tags field to automatically filter group topology
- ARGO-2181 Use report topology option to filter group topology
- ARGO-2190 Support tags filtering when listing group topology
- ARGO-2186 Support tags filtering in endpoint topology
- ARGO-2189 Support basic filtering when listing group topology
- ARGO-2185 Support basic filtering when listing endpoint topology
- ARGO-2106 Delete group topology
- ARGO-2112 List group topology per day
- ARGO-2110 Insert group topology list per day
- ARGO-2109 Delete endpoint topology for specific date
- ARGO-2108 List Endpoint topology per day
- ARGO-2106 Insert Endpoint Group Topology per day
- ARGO-2161 Add topology feed parameters to tenant configuration
- ARGO-2137 Add historic versioning to downtimes resources
- ARGO-2101 Create Downtime resource ARGO-2102 List all downtime resources 
- ARGO-2103 Delete specific downtime resource 
- ARGO-2104 Update downtime resource
- ARGO-2098 Implement weights resource with CRUD functionality
- ARGO-2099 Add historic versioning to weights resources
- ARGO-2098 Implement weights resource with CRUD functionality
- ARGO-2088 Add historic functionality for thresholds profiles
- ARGO-2001 Add history functionality for aggregation profiles
- ARGO-2002 Add History for operations profiles
- ARGO-2000 Add history functionality for metric profiles
- ARGO-2002 Add History for operations profiles
- ARGO-2098 Implement weights resource with CRUD functionality
- ARGO-2002 Add History for operations profiles
- ARGO-2002 Add history functionality for ops profiles


### Fixed: 
- ARGO-2184 Fix db reference name in package latest unittest
- ARGO-2142 Fix double routing prefixes for weights and downtimes. Fix 404 response when lists empty
- ARGO-2125 Fix regression of returning 404 when profile list were empty

### Changed:
- ARGO-2183 Change topology stats signature path


## 1.8.1 - (2019-11-15)

### Fixed:
- ARGO-2059 Fix latest strict result order by time instead of group name
- ARGO-2057 Fix latest results strict mode to honor limit & filter parameters
- ARGO-2058 Set strict=false by default in latest results call


## 1.8.0 - (2019-11-04)

### Added:
- ARGO-2038 Add version information to binary
- ARGO-2005 add extra information to tenants
- ARGO-2004 add extra information to reports
- ARGO-2003 Add a tenant list for web ui admin users
- ARGO-1997 Change recomputation status through web-api
- ARGO-1996 create recomputation: allow setting up recomputation requester name/email
- ARGO-1964 Return all daily metric data for specific host and date
- ARGO-1747 API Call - Get user by ID
- ARGO-1744 Add UUID for tenant's users

### Changed:
- ARGO-1727 Update the latest api call to be able to only return the latest entry
- ARGO-1983 Add filter param to return metric result list
- ARGO-1958 Fix add end of day point in multiple status timelines


## 1.7.9 - (2019-03-22)

### Added:
- ARGO-1438 Implement tenant general status
- ARGO-1680 serve endpoint a/r results

### Changed:
- ARGO-1455 - Migrate to golang/dep tool


## 1.7.8 - (2018-11-12)

### Added:
- ARGO-1388 Api call update tenant name field

### Changed:
- ARGO-1435 Fix configuration_profile json field in tenant status call
