"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[8084],{1673:(e,n,i)=>{i.r(n),i.d(n,{assets:()=>a,contentTitle:()=>s,default:()=>o,frontMatter:()=>r,metadata:()=>d,toc:()=>h});var l=i(4848),t=i(8453);const r={id:"v3_ar_results",title:"Availability / Reliability Results (v3)",sidebar_position:1},s=void 0,d={id:"apiv3/v3_ar_results",title:"Availability / Reliability Results (v3)",description:"API Calls",source:"@site/docs/apiv3/v3results.md",sourceDirName:"apiv3",slug:"/apiv3/v3_ar_results",permalink:"/argo-web-api/docs/apiv3/v3_ar_results",draft:!1,unlisted:!1,tags:[],version:"current",sidebarPosition:1,frontMatter:{id:"v3_ar_results",title:"Availability / Reliability Results (v3)",sidebar_position:1},sidebar:"tutorialSidebar",previous:{title:"API Version 3",permalink:"/argo-web-api/docs/category/api-version-3"},next:{title:"Status Results (v3)",permalink:"/argo-web-api/docs/apiv3/v3_status_results"}},a={},h=[{value:"API Calls",id:"api-calls",level:2},{value:"[GET]: List Availability and Reliability results for top level supergroups and included groups",id:"1",level:2},{value:"Input",id:"input",level:3},{value:"Query Parameters",id:"query-parameters",level:4},{value:"Path Parameters",id:"path-parameters",level:4},{value:"Example Request 1: default daily granularity",id:"example-request-1-default-daily-granularity",level:3},{value:"Request",id:"request",level:4},{value:"Method",id:"method",level:5},{value:"Path",id:"path",level:5},{value:"Headers",id:"headers",level:5},{value:"Response",id:"response",level:4},{value:"Code",id:"code",level:5},{value:"Body",id:"body",level:5},{value:"Example Request 2: monthly granularity",id:"example-request-2-monthly-granularity",level:3},{value:"Request",id:"request-1",level:4},{value:"Method",id:"method-1",level:5},{value:"Path",id:"path-1",level:5},{value:"Headers",id:"headers-1",level:5},{value:"Response",id:"response-1",level:4},{value:"Code",id:"code-1",level:5},{value:"Body",id:"body-1",level:5},{value:"Example Request 3: Custom granularity",id:"example-request-3-custom-granularity",level:3},{value:"Request",id:"request-2",level:4},{value:"Method",id:"method-2",level:5},{value:"Path",id:"path-2",level:5},{value:"Headers",id:"headers-2",level:5},{value:"Response",id:"response-2",level:4},{value:"Code",id:"code-2",level:5},{value:"Body",id:"body-2",level:5},{value:"[GET]: List Availability and Reliability results for endpoints with specific resource-id",id:"2",level:2},{value:"Input",id:"input-1",level:3},{value:"Query Parameters",id:"query-parameters-1",level:4},{value:"Path Parameters",id:"path-parameters-1",level:4},{value:"Example Request 1: default daily granularity with specific resource-id",id:"example-request-1-default-daily-granularity-with-specific-resource-id",level:3},{value:"Request",id:"request-3",level:4},{value:"Method",id:"method-3",level:5},{value:"Path",id:"path-3",level:5},{value:"Headers",id:"headers-3",level:5},{value:"Response",id:"response-3",level:4},{value:"Code",id:"code-3",level:5},{value:"Body",id:"body-3",level:5},{value:"Example Request 2: monthly granularity with specific resource-id",id:"example-request-2-monthly-granularity-with-specific-resource-id",level:3},{value:"Request",id:"request-4",level:4},{value:"Method",id:"method-4",level:5},{value:"Path",id:"path-4",level:5},{value:"Headers",id:"headers-4",level:5},{value:"Response",id:"response-4",level:4},{value:"Code",id:"code-4",level:5},{value:"Body",id:"body-4",level:5},{value:"Example Request 3: custom period granularity with specific resource-id",id:"example-request-3-custom-period-granularity-with-specific-resource-id",level:3},{value:"Request",id:"request-5",level:4},{value:"Method",id:"method-5",level:5},{value:"Path",id:"path-5",level:5},{value:"Headers",id:"headers-5",level:5},{value:"Response",id:"response-5",level:4},{value:"Code",id:"code-5",level:5},{value:"Body",id:"body-5",level:5}];function c(e){const n={a:"a",code:"code",em:"em",h2:"h2",h3:"h3",h4:"h4",h5:"h5",p:"p",pre:"pre",table:"table",tbody:"tbody",td:"td",th:"th",thead:"thead",tr:"tr",...(0,t.R)(),...e.components};return(0,l.jsxs)(l.Fragment,{children:[(0,l.jsx)(n.h2,{id:"api-calls",children:"API Calls"}),"\n",(0,l.jsxs)(n.p,{children:[(0,l.jsx)(n.em,{children:"Note"}),": These are v3 api calls implementations found under the path ",(0,l.jsx)(n.code,{children:"/api/v3"})]}),"\n",(0,l.jsxs)(n.table,{children:[(0,l.jsx)(n.thead,{children:(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.th,{children:"Name"}),(0,l.jsx)(n.th,{children:"Description"}),(0,l.jsx)(n.th,{children:"Shortcut"})]})}),(0,l.jsxs)(n.tbody,{children:[(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.td,{children:"GET: List Availability and Reliability results for top level supergroups and included groups"}),(0,l.jsx)(n.td,{children:"This method retrieves the a/r results of all top level supergroups and their included groups"}),(0,l.jsx)(n.td,{children:(0,l.jsx)(n.a,{href:"#1",children:"Description"})})]}),(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.td,{children:"GET: List Availability and Reliability results for specific endpoint using resource id"}),(0,l.jsx)(n.td,{children:"This method retrieves the a/r results of all endoints with specific resource id"}),(0,l.jsx)(n.td,{children:(0,l.jsx)(n.a,{href:"#2",children:"Description"})})]})]})]}),"\n",(0,l.jsx)(n.h2,{id:"1",children:"[GET]: List Availability and Reliability results for top level supergroups and included groups"}),"\n",(0,l.jsxs)(n.p,{children:["The following methods can be used to obtain a tenant's Availability and Reliability result metrics for all top level supergroups and included groups. The api authenticates the tenant using the api-key within the x-api-key header. User can specify time granularity (",(0,l.jsx)(n.code,{children:"monthly"}),", ",(0,l.jsx)(n.code,{children:"daily"})," or ",(0,l.jsx)(n.code,{children:"custom"}),") for retrieved results and also format using the ",(0,l.jsx)(n.code,{children:"Accept"})," header."]}),"\n",(0,l.jsx)(n.h3,{id:"input",children:"Input"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"/results/{report_name}?[start_time]&[end_time]&[granularity]\n"})}),"\n",(0,l.jsx)(n.h4,{id:"query-parameters",children:"Query Parameters"}),"\n",(0,l.jsxs)(n.table,{children:[(0,l.jsx)(n.thead,{children:(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.th,{children:"Type"}),(0,l.jsx)(n.th,{children:"Description"}),(0,l.jsx)(n.th,{children:"Required"}),(0,l.jsx)(n.th,{children:"Default value"})]})}),(0,l.jsxs)(n.tbody,{children:[(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.td,{children:(0,l.jsx)(n.code,{children:"[start_time]"})}),(0,l.jsx)(n.td,{children:"UTC time in W3C format"}),(0,l.jsx)(n.td,{children:"YES"}),(0,l.jsx)(n.td,{})]}),(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.td,{children:(0,l.jsx)(n.code,{children:"[end_time]"})}),(0,l.jsx)(n.td,{children:"UTC time in W3C format"}),(0,l.jsx)(n.td,{children:"YES"}),(0,l.jsx)(n.td,{})]}),(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.td,{children:(0,l.jsx)(n.code,{children:"[granularity]"})}),(0,l.jsxs)(n.td,{children:["Granularity of time that will be used to present data. Possible values are ",(0,l.jsx)(n.code,{children:"monthly"}),",  ",(0,l.jsx)(n.code,{children:"daily"})," or ",(0,l.jsx)(n.code,{children:"custom"})]}),(0,l.jsx)(n.td,{children:"NO"}),(0,l.jsx)(n.td,{children:(0,l.jsx)(n.code,{children:"daily"})})]})]})]}),"\n",(0,l.jsx)(n.h4,{id:"path-parameters",children:"Path Parameters"}),"\n",(0,l.jsxs)(n.table,{children:[(0,l.jsx)(n.thead,{children:(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.th,{children:"Name"}),(0,l.jsx)(n.th,{children:"Description"}),(0,l.jsx)(n.th,{children:"Required"}),(0,l.jsx)(n.th,{children:"Default value"})]})}),(0,l.jsx)(n.tbody,{children:(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.td,{children:(0,l.jsx)(n.code,{children:"{report_name}"})}),(0,l.jsx)(n.td,{children:"Name of the report that contains all the information about the profile, filter tags, group types etc."}),(0,l.jsx)(n.td,{children:"YES"}),(0,l.jsx)(n.td,{})]})})]}),"\n",(0,l.jsx)(n.h3,{id:"example-request-1-default-daily-granularity",children:"Example Request 1: default daily granularity"}),"\n",(0,l.jsx)(n.h4,{id:"request",children:"Request"}),"\n",(0,l.jsx)(n.h5,{id:"method",children:"Method"}),"\n",(0,l.jsx)(n.p,{children:(0,l.jsx)(n.code,{children:"HTTP GET"})}),"\n",(0,l.jsx)(n.h5,{id:"path",children:"Path"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"/api/v3/results/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z \n"})}),"\n",(0,l.jsx)(n.p,{children:"or"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"/api/v3/results/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily`\n"})}),"\n",(0,l.jsx)(n.h5,{id:"headers",children:"Headers"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:'x-api-key: "tenant_key_value"\nAccept: "application/json"\n'})}),"\n",(0,l.jsx)(n.h4,{id:"response",children:"Response"}),"\n",(0,l.jsx)(n.h5,{id:"code",children:"Code"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"Status: 200 OK\n"})}),"\n",(0,l.jsx)(n.h5,{id:"body",children:"Body"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-json",children:'{\n  "results": [\n    {\n      "name": "GROUP_A",\n      "type": "GROUP",\n      "results": [\n        {\n          "date": "2015-06-22",\n          "availability": "68.13896116893515",\n          "reliability": "68.13896116893515"\n        },\n        {\n          "date": "2015-06-23",\n          "availability": "75.36324059247399",\n          "reliability": "75.36324059247399"\n        }\n      ],\n      "groups": [\n        {\n          "name": "ST01",\n          "type": "SITES",\n          "results": [\n            {\n              "date": "2015-06-22",\n              "availability": "66.7",\n              "reliability": "66.7",\n              "unknown": "0",\n              "uptime": "66.7",\n              "downtime": "0"\n            },\n            {\n              "date": "2015-06-23",\n              "availability": "100",\n              "reliability": "100",\n              "unknown": "0",\n              "uptime": "1",\n              "downtime": "0"\n            }\n          ]\n        },\n        {\n          "name": "ST02",\n          "type": "SITES",\n          "results": [\n            {\n              "date": "2015-06-22",\n              "availability": "70",\n              "reliability": "70",\n              "unknown": "0",\n              "uptime": "0.70",\n              "downtime": "0"\n            },\n            {\n              "date": "2015-06-23",\n              "availability": "43.5",\n              "reliability": "43.5",\n              "unknown": "0",\n              "uptime": "0.435",\n              "downtime": "0"\n            }\n          ]\n        }\n      ]\n    }\n  ]\n}\n'})}),"\n",(0,l.jsx)(n.h3,{id:"example-request-2-monthly-granularity",children:"Example Request 2: monthly granularity"}),"\n",(0,l.jsx)(n.h4,{id:"request-1",children:"Request"}),"\n",(0,l.jsx)(n.h5,{id:"method-1",children:"Method"}),"\n",(0,l.jsx)(n.p,{children:(0,l.jsx)(n.code,{children:"HTTP GET"})}),"\n",(0,l.jsx)(n.h5,{id:"path-1",children:"Path"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"/api/v3/results/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=monthly\n"})}),"\n",(0,l.jsx)(n.h5,{id:"headers-1",children:"Headers"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:'x-api-key: "tenant_key_value"\nAccept: "application/json"\n'})}),"\n",(0,l.jsx)(n.h4,{id:"response-1",children:"Response"}),"\n",(0,l.jsx)(n.h5,{id:"code-1",children:"Code"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"Status: 200 OK\n"})}),"\n",(0,l.jsx)(n.h5,{id:"body-1",children:"Body"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-json",children:'{\n  "results": [\n    {\n      "name": "GROUP_A",\n      "type": "GROUP",\n      "results": [\n        {\n          "date": "2015-06",\n          "availability": "99.99999900000002",\n          "reliability": "99.99999900000002"\n        }\n      ],\n      "groups": [\n        {\n          "name": "ST01",\n          "type": "SITES",\n          "results": [\n            {\n              "date": "2015-06",\n              "availability": "99.99999900000002",\n              "reliability": "99.99999900000002",\n              "unknown": "0",\n              "uptime": "1",\n              "downtime": "0"\n            }\n          ]\n        },\n        {\n          "name": "ST02",\n          "type": "SITES",\n          "results": [\n            {\n              "date": "2015-06",\n              "availability": "99.99999900000002",\n              "reliability": "99.99999900000002",\n              "unknown": "0",\n              "uptime": "1",\n              "downtime": "0"\n            }\n          ]\n        }\n      ]\n    }\n  ]\n}\n'})}),"\n",(0,l.jsx)(n.h3,{id:"example-request-3-custom-granularity",children:"Example Request 3: Custom granularity"}),"\n",(0,l.jsxs)(n.p,{children:["This request returns availability/reliability score numbers for the whole custom period defined between ",(0,l.jsx)(n.code,{children:"start_time"})," and ",(0,l.jsx)(n.code,{children:"end_time"}),".\nThis means that for each item the user will receive one availability and reliability result concerning the whole period (instead of multiple daily or monthly results)"]}),"\n",(0,l.jsx)(n.h4,{id:"request-2",children:"Request"}),"\n",(0,l.jsx)(n.h5,{id:"method-2",children:"Method"}),"\n",(0,l.jsx)(n.p,{children:(0,l.jsx)(n.code,{children:"HTTP GET"})}),"\n",(0,l.jsx)(n.h5,{id:"path-2",children:"Path"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"/api/v3/results/Report_A?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=custom\n"})}),"\n",(0,l.jsx)(n.h5,{id:"headers-2",children:"Headers"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:'x-api-key: "tenant_key_value"\nAccept: "application/json"\n'})}),"\n",(0,l.jsx)(n.h4,{id:"response-2",children:"Response"}),"\n",(0,l.jsx)(n.h5,{id:"code-2",children:"Code"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"Status: 200 OK\n"})}),"\n",(0,l.jsx)(n.h5,{id:"body-2",children:"Body"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-json",children:'{\n  "results": [\n    {\n      "name": "GROUP_A",\n      "type": "GROUP",\n      "results": [\n        {\n          "availability": "99.99999900000002",\n          "reliability": "99.99999900000002"\n        }\n      ],\n      "groups": [\n        {\n          "name": "ST01",\n          "type": "SITES",\n          "results": [\n            {\n              "availability": "99.99999900000002",\n              "reliability": "99.99999900000002",\n              "unknown": "0",\n              "uptime": "1",\n              "downtime": "0"\n            }\n          ]\n        },\n        {\n          "name": "ST02",\n          "type": "SITES",\n          "results": [\n            {\n              "availability": "99.99999900000002",\n              "reliability": "99.99999900000002",\n              "unknown": "0",\n              "uptime": "1",\n              "downtime": "0"\n            }\n          ]\n        }\n      ]\n    }\n  ]\n}\n'})}),"\n",(0,l.jsx)(n.h2,{id:"2",children:"[GET]: List Availability and Reliability results for endpoints with specific resource-id"}),"\n",(0,l.jsxs)(n.p,{children:["The following methods can be used to obtain a tenant's Availability and Reliability result for the endpoints that have a specific resource-id. User can specify a period with ",(0,l.jsx)(n.code,{children:"start_time"})," and ",(0,l.jsx)(n.code,{children:"end_time"})," and granularity(",(0,l.jsx)(n.code,{children:"monthly"}),", ",(0,l.jsx)(n.code,{children:"daily"})," or ",(0,l.jsx)(n.code,{children:"custom"}),") for retrieved results. ",(0,l.jsx)(n.code,{children:"Accept"})," header is required."]}),"\n",(0,l.jsx)(n.h3,{id:"input-1",children:"Input"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"/results/{report_name}/id/{resource-id}?[start_time]&[end_time]&[granularity]\n"})}),"\n",(0,l.jsx)(n.h4,{id:"query-parameters-1",children:"Query Parameters"}),"\n",(0,l.jsxs)(n.table,{children:[(0,l.jsx)(n.thead,{children:(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.th,{children:"Type"}),(0,l.jsx)(n.th,{children:"Description"}),(0,l.jsx)(n.th,{children:"Required"}),(0,l.jsx)(n.th,{children:"Default value"})]})}),(0,l.jsxs)(n.tbody,{children:[(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.td,{children:(0,l.jsx)(n.code,{children:"[start_time]"})}),(0,l.jsx)(n.td,{children:"UTC time in W3C format"}),(0,l.jsx)(n.td,{children:"YES"}),(0,l.jsx)(n.td,{})]}),(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.td,{children:(0,l.jsx)(n.code,{children:"[end_time]"})}),(0,l.jsx)(n.td,{children:"UTC time in W3C format"}),(0,l.jsx)(n.td,{children:"YES"}),(0,l.jsx)(n.td,{})]}),(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.td,{children:(0,l.jsx)(n.code,{children:"[granularity]"})}),(0,l.jsxs)(n.td,{children:["Granularity of time that will be used to present data. Possible values are ",(0,l.jsx)(n.code,{children:"monthly"}),", ",(0,l.jsx)(n.code,{children:"daily"})," or ",(0,l.jsx)(n.code,{children:"custom"})]}),(0,l.jsx)(n.td,{children:"NO"}),(0,l.jsx)(n.td,{children:(0,l.jsx)(n.code,{children:"daily"})})]})]})]}),"\n",(0,l.jsx)(n.h4,{id:"path-parameters-1",children:"Path Parameters"}),"\n",(0,l.jsxs)(n.table,{children:[(0,l.jsx)(n.thead,{children:(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.th,{children:"Name"}),(0,l.jsx)(n.th,{children:"Description"}),(0,l.jsx)(n.th,{children:"Required"}),(0,l.jsx)(n.th,{children:"Default value"})]})}),(0,l.jsxs)(n.tbody,{children:[(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.td,{children:(0,l.jsx)(n.code,{children:"{report_name}"})}),(0,l.jsx)(n.td,{children:"Name of the report that contains all the information about the profile, filter tags, group types etc."}),(0,l.jsx)(n.td,{children:"YES"}),(0,l.jsx)(n.td,{})]}),(0,l.jsxs)(n.tr,{children:[(0,l.jsx)(n.td,{children:(0,l.jsx)(n.code,{children:"{id}"})}),(0,l.jsx)(n.td,{children:"The resource id"}),(0,l.jsx)(n.td,{children:"YES"}),(0,l.jsx)(n.td,{})]})]})]}),"\n",(0,l.jsx)(n.h3,{id:"example-request-1-default-daily-granularity-with-specific-resource-id",children:"Example Request 1: default daily granularity with specific resource-id"}),"\n",(0,l.jsx)(n.h4,{id:"request-3",children:"Request"}),"\n",(0,l.jsx)(n.h5,{id:"method-3",children:"Method"}),"\n",(0,l.jsx)(n.p,{children:(0,l.jsx)(n.code,{children:"HTTP GET"})}),"\n",(0,l.jsx)(n.h5,{id:"path-3",children:"Path"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"/api/v3/results/Report_A/id/simple-queue?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z \n"})}),"\n",(0,l.jsx)(n.p,{children:"or"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"/api/v3/results/Report_A/id/simple-queue?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=daily`\n"})}),"\n",(0,l.jsx)(n.h5,{id:"headers-3",children:"Headers"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:'x-api-key: "tenant_key_value"\nAccept: "application/json"\n'})}),"\n",(0,l.jsx)(n.h4,{id:"response-3",children:"Response"}),"\n",(0,l.jsx)(n.h5,{id:"code-3",children:"Code"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"Status: 200 OK\n"})}),"\n",(0,l.jsx)(n.h5,{id:"body-3",children:"Body"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-json",children:'{\n  "id": "simple-queue",\n  "endpoints": [\n    {\n      "name": "host01.example",\n      "service": "service.queue",\n      "group": "Infra-01",\n      "info": {\n        "URL": "http://submit.queue01.example.com"\n      },\n      "results": [\n        {\n          "date": "2015-06-22",\n          "availability": "99.99999900000002",\n          "reliability": "99.99999900000002",\n          "unknown": "0",\n          "uptime": "1",\n          "downtime": "0"\n        },\n        {\n          "date": "2015-06-23",\n          "availability": "99.99999900000002",\n          "reliability": "99.99999900000002",\n          "unknown": "0",\n          "uptime": "1",\n          "downtime": "0"\n        }\n      ]\n    }\n  ]\n}\n'})}),"\n",(0,l.jsx)(n.h3,{id:"example-request-2-monthly-granularity-with-specific-resource-id",children:"Example Request 2: monthly granularity with specific resource-id"}),"\n",(0,l.jsx)(n.h4,{id:"request-4",children:"Request"}),"\n",(0,l.jsx)(n.h5,{id:"method-4",children:"Method"}),"\n",(0,l.jsx)(n.p,{children:(0,l.jsx)(n.code,{children:"HTTP GET"})}),"\n",(0,l.jsx)(n.h5,{id:"path-4",children:"Path"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"/api/v3/results/Report_A/id/simple-queue?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=monthly\n"})}),"\n",(0,l.jsx)(n.h5,{id:"headers-4",children:"Headers"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:'x-api-key: "tenant_key_value"\nAccept: "application/json"\n'})}),"\n",(0,l.jsx)(n.h4,{id:"response-4",children:"Response"}),"\n",(0,l.jsx)(n.h5,{id:"code-4",children:"Code"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"Status: 200 OK\n"})}),"\n",(0,l.jsx)(n.h5,{id:"body-4",children:"Body"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-json",children:'{\n  "id": "simple-queue",\n  "endpoints": [\n    {\n      "name": "host01.example",\n      "service": "service.queue",\n      "group": "Infra-01",\n      "info": {\n        "URL": "http://submit.queue01.example.com"\n      },\n      "results": [\n        {\n          "date": "2015-06",\n          "availability": "99.99999900000002",\n          "reliability": "99.99999900000002",\n          "unknown": "0",\n          "uptime": "1",\n          "downtime": "0"\n        }\n      ]\n    }\n  ]\n}\n'})}),"\n",(0,l.jsx)(n.h3,{id:"example-request-3-custom-period-granularity-with-specific-resource-id",children:"Example Request 3: custom period granularity with specific resource-id"}),"\n",(0,l.jsxs)(n.p,{children:["This request returns availability/reliability score numbers for the whole custom period defined between ",(0,l.jsx)(n.code,{children:"start_time"})," and ",(0,l.jsx)(n.code,{children:"end_time"}),".\nThis means that for each item with the specific resource-id the user will receive one availability and reliability result concerning the whole period (instead of multiple daily or monthly results)"]}),"\n",(0,l.jsx)(n.h4,{id:"request-5",children:"Request"}),"\n",(0,l.jsx)(n.h5,{id:"method-5",children:"Method"}),"\n",(0,l.jsx)(n.p,{children:(0,l.jsx)(n.code,{children:"HTTP GET"})}),"\n",(0,l.jsx)(n.h5,{id:"path-5",children:"Path"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"/api/v3/results/Report_A/id/simple-queue?start_time=2015-06-20T12:00:00Z&end_time=2015-06-26T23:00:00Z&granularity=custom\n"})}),"\n",(0,l.jsx)(n.h5,{id:"headers-5",children:"Headers"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:'x-api-key: "tenant_key_value"\nAccept: "application/json"\n'})}),"\n",(0,l.jsx)(n.h4,{id:"response-5",children:"Response"}),"\n",(0,l.jsx)(n.h5,{id:"code-5",children:"Code"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{children:"Status: 200 OK\n"})}),"\n",(0,l.jsx)(n.h5,{id:"body-5",children:"Body"}),"\n",(0,l.jsx)(n.pre,{children:(0,l.jsx)(n.code,{className:"language-json",children:'{\n  "id": "simple-queue",\n  "endpoints": [\n    {\n      "name": "host01.example",\n      "service": "service.queue",\n      "group": "Infra-01",\n      "info": {\n        "URL": "http://submit.queue01.example.com"\n      },\n      "results": [\n        {\n          "availability": "99.99999900000002",\n          "reliability": "99.99999900000002",\n          "unknown": "0",\n          "uptime": "1",\n          "downtime": "0"\n        }\n      ]\n    }\n  ]\n}\n'})})]})}function o(e={}){const{wrapper:n}={...(0,t.R)(),...e.components};return n?(0,l.jsx)(n,{...e,children:(0,l.jsx)(c,{...e})}):c(e)}},8453:(e,n,i)=>{i.d(n,{R:()=>s,x:()=>d});var l=i(6540);const t={},r=l.createContext(t);function s(e){const n=l.useContext(r);return l.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function d(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(t):e.components||t:s(e.components),l.createElement(r.Provider,{value:n},e.children)}}}]);