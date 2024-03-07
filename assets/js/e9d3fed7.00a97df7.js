"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[2562],{2140:(e,s,n)=>{n.r(s),n.d(s,{assets:()=>c,contentTitle:()=>d,default:()=>h,frontMatter:()=>i,metadata:()=>l,toc:()=>a});var t=n(4848),r=n(8453);const i={id:"latest",title:"Latest Metric results",sidebar_position:3},d=void 0,l={id:"results/latest",title:"Latest Metric results",description:"API calls for retrieving latest metric results",source:"@site/docs/results/latest.md",sourceDirName:"results",slug:"/results/latest",permalink:"/argo-web-api/docs/results/latest",draft:!1,unlisted:!1,tags:[],version:"current",sidebarPosition:3,frontMatter:{id:"latest",title:"Latest Metric results",sidebar_position:3},sidebar:"tutorialSidebar",previous:{title:"Status Results",permalink:"/argo-web-api/docs/results/status_results"},next:{title:"Metric Results",permalink:"/argo-web-api/docs/results/metric_results"}},c={},a=[{value:"API calls for retrieving latest metric results",id:"api-calls-for-retrieving-latest-metric-results",level:2},{value:"[GET]: List all latest metric data",id:"1",level:2},{value:"Input",id:"input",level:3},{value:"List All latest metric data",id:"list-all-latest-metric-data",level:5},{value:"Path Parameters",id:"path-parameters",level:4},{value:"Url Parameters",id:"url-parameters",level:4},{value:"Headers",id:"headers",level:4},{value:"Response Code",id:"response-code",level:4},{value:"Response body",id:"response-body",level:3},{value:"Example Request:",id:"example-request",level:6},{value:"Example Response:",id:"example-response",level:6},{value:"[GET]: List All Metric Data for a specific endpoint group",id:"2",level:2},{value:"Input",id:"input-1",level:3},{value:"Path Parameters",id:"path-parameters-1",level:4},{value:"Url Parameters",id:"url-parameters-1",level:4},{value:"Headers",id:"headers-1",level:4},{value:"Response Code",id:"response-code-1",level:4},{value:"Response body",id:"response-body-1",level:3},{value:"Example Request:",id:"example-request-1",level:6},{value:"Example Response:",id:"example-response-1",level:6}];function o(e){const s={code:"code",em:"em",h2:"h2",h3:"h3",h4:"h4",h5:"h5",h6:"h6",p:"p",pre:"pre",strong:"strong",table:"table",tbody:"tbody",td:"td",th:"th",thead:"thead",tr:"tr",...(0,r.R)(),...e.components};return(0,t.jsxs)(t.Fragment,{children:[(0,t.jsx)(s.h2,{id:"api-calls-for-retrieving-latest-metric-results",children:"API calls for retrieving latest metric results"}),"\n",(0,t.jsxs)(s.table,{children:[(0,t.jsx)(s.thead,{children:(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.th,{children:"Name"}),(0,t.jsx)(s.th,{children:"Description"}),(0,t.jsx)(s.th,{children:"Shortcut"})]})}),(0,t.jsxs)(s.tbody,{children:[(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:"GET: List all latest metric data"}),(0,t.jsx)(s.td,{children:"List latest metric data ."}),(0,t.jsx)(s.td,{children:(0,t.jsx)("a",{href:"#1",children:"Description"})})]}),(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:"GET: List latest metric data for Group"}),(0,t.jsx)(s.td,{children:"List latest metric data for a specific endpoint group."}),(0,t.jsx)(s.td,{children:(0,t.jsx)("a",{href:"#2",children:"Description"})})]})]})]}),"\n",(0,t.jsx)(s.h2,{id:"1",children:"[GET]: List all latest metric data"}),"\n",(0,t.jsx)(s.p,{children:"This method may be used to retrieve latest metric data available in a report. User can filer the results by status and limit the amount\nof results returned"}),"\n",(0,t.jsx)(s.h3,{id:"input",children:"Input"}),"\n",(0,t.jsx)(s.h5,{id:"list-all-latest-metric-data",children:"List All latest metric data"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{children:"/latest/{report}/{group_type}?[date]&[filter]&[limit]&[strict]\n"})}),"\n",(0,t.jsx)(s.h4,{id:"path-parameters",children:"Path Parameters"}),"\n",(0,t.jsxs)(s.table,{children:[(0,t.jsx)(s.thead,{children:(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.th,{children:"Type"}),(0,t.jsx)(s.th,{children:"Description"}),(0,t.jsx)(s.th,{children:"Required"}),(0,t.jsx)(s.th,{children:"Default value"})]})}),(0,t.jsxs)(s.tbody,{children:[(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:(0,t.jsx)(s.code,{children:"report"})}),(0,t.jsx)(s.td,{children:"name of the report used"}),(0,t.jsx)(s.td,{children:"YES"}),(0,t.jsx)(s.td,{})]}),(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:(0,t.jsx)(s.code,{children:"group_type"})}),(0,t.jsx)(s.td,{children:"type of endpoint group"}),(0,t.jsx)(s.td,{children:"YES"}),(0,t.jsx)(s.td,{})]}),(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:(0,t.jsx)(s.code,{children:"group_name"})}),(0,t.jsx)(s.td,{children:"name of endpoint group"}),(0,t.jsx)(s.td,{children:"YES"}),(0,t.jsx)(s.td,{})]})]})]}),"\n",(0,t.jsx)(s.h4,{id:"url-parameters",children:"Url Parameters"}),"\n",(0,t.jsxs)(s.table,{children:[(0,t.jsx)(s.thead,{children:(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.th,{children:"Type"}),(0,t.jsx)(s.th,{children:"Description"}),(0,t.jsx)(s.th,{children:"Required"}),(0,t.jsx)(s.th,{children:"Default value"})]})}),(0,t.jsxs)(s.tbody,{children:[(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:(0,t.jsx)(s.code,{children:"date"})}),(0,t.jsx)(s.td,{children:"target a specific data"}),(0,t.jsx)(s.td,{children:"NO"}),(0,t.jsx)(s.td,{children:"today's date"})]}),(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:(0,t.jsx)(s.code,{children:"filter"})}),(0,t.jsxs)(s.td,{children:["filter by status values (",(0,t.jsx)(s.code,{children:"all"}),",",(0,t.jsx)(s.code,{children:"non-ok"}),",",(0,t.jsx)(s.code,{children:"ok"}),",",(0,t.jsx)(s.code,{children:"critical"}),",",(0,t.jsx)(s.code,{children:"warning"}),",",(0,t.jsx)(s.code,{children:"unknown"}),",",(0,t.jsx)(s.code,{children:"missing"}),")"]}),(0,t.jsx)(s.td,{children:"NO"}),(0,t.jsx)(s.td,{children:"all"})]}),(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:(0,t.jsx)(s.code,{children:"limit"})}),(0,t.jsx)(s.td,{children:"limit number of results returned"}),(0,t.jsx)(s.td,{children:"NO"}),(0,t.jsx)(s.td,{children:"500"})]}),(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:(0,t.jsx)(s.code,{children:"strict"})}),(0,t.jsx)(s.td,{children:"strict mode on/off"}),(0,t.jsx)(s.td,{children:"NO"}),(0,t.jsx)(s.td,{children:'"false"'})]})]})]}),"\n",(0,t.jsxs)(s.p,{children:[(0,t.jsx)(s.em,{children:(0,t.jsx)(s.strong,{children:"Notes"})}),":"]}),"\n",(0,t.jsxs)(s.p,{children:[(0,t.jsx)(s.code,{children:"group_type"})," in the specific request refer always to endpoint groups (e.g.. ",(0,t.jsx)(s.code,{children:"SITES"}),")."]}),"\n",(0,t.jsxs)(s.p,{children:[(0,t.jsx)(s.code,{children:"strict"})," when strict mode is ON, the response will contain only the latest grouped by endpoint_group/host/service/metric."]}),"\n",(0,t.jsx)(s.h4,{id:"headers",children:"Headers"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{children:"x-api-key: secret_key_value\nAccept: application/json or application/xml\n"})}),"\n",(0,t.jsx)(s.h4,{id:"response-code",children:"Response Code"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{children:"Status: 200 OK\n"})}),"\n",(0,t.jsx)(s.h3,{id:"response-body",children:"Response body"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{className:"language-json",children:'{\n    "status": {\n        "message": "type of message received",\n        "code": "return code"\n    },\n    "data": {\n        "metric_data": [\n            {\n                "endpoint_group": "name of endpoint group",\n                "service": "name of service",\n                "endpoint": "name of endpoint",\n                "metric": "name of metric",\n                "timestamp": "2018-06-22T11:55:44Z",\n                "status": "OK || WARNING || CRITICAL || MISSING || UNKNOWN",\n                "summary": "summary of the metric message - generated by the monitoring engine (nagios)",\n                "message": "body of nagios generated message"\n            }\n         ]\n      }\n}\n'})}),"\n",(0,t.jsx)(s.h6,{id:"example-request",children:"Example Request:"}),"\n",(0,t.jsx)(s.p,{children:"URL:"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{children:"latest/Report_B/TENANT_SITES?date=2015-05-01T00:00:00Z&filter=non-ok&limit=10\n"})}),"\n",(0,t.jsx)(s.p,{children:"Headers:"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{children:"x-api-key: secret_key_value\nAccept: application/json or application/xml\n\n"})}),"\n",(0,t.jsx)(s.h6,{id:"example-response",children:"Example Response:"}),"\n",(0,t.jsx)(s.p,{children:"Code:"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{children:"Status: 200 OK\n"})}),"\n",(0,t.jsx)(s.p,{children:"Response body:"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{className:"language-json",children:'{\n"status": {\n"message": "application/json",\n"code": "200"\n},\n"data": {\n"metric_data": [\n {\n  "endpoint_group": "EL-01-AUTH",\n  "service": "someService",\n  "endpoint": "someservice.example.gr",\n  "metric": "someService-FileTransfer",\n  "timestamp": "2015-05-01T05:00:00Z",\n  "status": "WARNING",\n  "summary": "someService status is ok",\n  "message": "someService data upload test return value of ok"\n },\n {\n  "endpoint_group": "EL-02-AUTH",\n  "service": "someService",\n  "endpoint": "someservice.example.gr",\n  "metric": "someService-FileTransfer",\n  "timestamp": "2015-05-01T05:00:00Z",\n  "status": "MISSING",\n  "summary": "someService status is ok",\n  "message": "someService data upload test return value of ok"\n },\n {\n  "endpoint_group": "EL-03-AUTH",\n  "service": "someService",\n  "endpoint": "someservice.example.gr",\n  "metric": "someService-FileTransfer",\n  "timestamp": "2015-05-01T05:00:00Z",\n  "status": "CRITICAL",\n  "summary": "someService status is ok",\n  "message": "someService data upload test return value of ok"\n },\n {\n  "endpoint_group": "EL-01-AUTH",\n  "service": "someService",\n  "endpoint": "someservice.example.gr",\n  "metric": "someService-FileTransfer",\n  "timestamp": "2015-05-01T01:00:00Z",\n  "status": "UNKNOWN",\n  "summary": "someService status is CRITICAL",\n  "message": "someService data upload test failed"\n }\n]\n}\n}\n'})}),"\n",(0,t.jsx)(s.h2,{id:"2",children:"[GET]: List All Metric Data for a specific endpoint group"}),"\n",(0,t.jsx)(s.p,{children:"This method may be used to retrieve latest metric data available in a report for a specific endpoint group. User can filer the results by status and limit the amount\nof results returned"}),"\n",(0,t.jsx)(s.h3,{id:"input-1",children:"Input"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{children:"/latest/{report}/{group_type}/{group_name}/?[date]&[filter]&[limit]\n"})}),"\n",(0,t.jsx)(s.h4,{id:"path-parameters-1",children:"Path Parameters"}),"\n",(0,t.jsxs)(s.table,{children:[(0,t.jsx)(s.thead,{children:(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.th,{children:"Type"}),(0,t.jsx)(s.th,{children:"Description"}),(0,t.jsx)(s.th,{children:"Required"}),(0,t.jsx)(s.th,{children:"Default value"})]})}),(0,t.jsxs)(s.tbody,{children:[(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:(0,t.jsx)(s.code,{children:"report"})}),(0,t.jsx)(s.td,{children:"name of the report used"}),(0,t.jsx)(s.td,{children:"YES"}),(0,t.jsx)(s.td,{})]}),(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:(0,t.jsx)(s.code,{children:"group_type"})}),(0,t.jsx)(s.td,{children:"type of endpoint group"}),(0,t.jsx)(s.td,{children:"YES"}),(0,t.jsx)(s.td,{})]}),(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:(0,t.jsx)(s.code,{children:"group_name"})}),(0,t.jsx)(s.td,{children:"name of endpoint group"}),(0,t.jsx)(s.td,{children:"YES"}),(0,t.jsx)(s.td,{})]})]})]}),"\n",(0,t.jsx)(s.h4,{id:"url-parameters-1",children:"Url Parameters"}),"\n",(0,t.jsxs)(s.table,{children:[(0,t.jsx)(s.thead,{children:(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.th,{children:"Type"}),(0,t.jsx)(s.th,{children:"Description"}),(0,t.jsx)(s.th,{children:"Required"}),(0,t.jsx)(s.th,{children:"Default value"})]})}),(0,t.jsxs)(s.tbody,{children:[(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:(0,t.jsx)(s.code,{children:"date"})}),(0,t.jsx)(s.td,{children:"target a specific data"}),(0,t.jsx)(s.td,{children:"NO"}),(0,t.jsx)(s.td,{children:"today's date"})]}),(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:(0,t.jsx)(s.code,{children:"filter"})}),(0,t.jsxs)(s.td,{children:["filter by status values (",(0,t.jsx)(s.code,{children:"all"}),",",(0,t.jsx)(s.code,{children:"non-ok"}),",",(0,t.jsx)(s.code,{children:"ok"}),",",(0,t.jsx)(s.code,{children:"critical"}),",",(0,t.jsx)(s.code,{children:"warning"}),",",(0,t.jsx)(s.code,{children:"unknown"}),",",(0,t.jsx)(s.code,{children:"missing"}),")"]}),(0,t.jsx)(s.td,{children:"NO"}),(0,t.jsx)(s.td,{children:"all"})]}),(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:(0,t.jsx)(s.code,{children:"limit"})}),(0,t.jsx)(s.td,{children:"limit number of results returned"}),(0,t.jsx)(s.td,{children:"NO"}),(0,t.jsx)(s.td,{children:"500"})]}),(0,t.jsxs)(s.tr,{children:[(0,t.jsx)(s.td,{children:(0,t.jsx)(s.code,{children:"strict"})}),(0,t.jsx)(s.td,{children:"strict mode on/off"}),(0,t.jsx)(s.td,{children:"NO"}),(0,t.jsx)(s.td,{children:'"false"'})]})]})]}),"\n",(0,t.jsxs)(s.p,{children:[(0,t.jsx)(s.em,{children:(0,t.jsx)(s.strong,{children:"Notes"})}),":"]}),"\n",(0,t.jsxs)(s.p,{children:[(0,t.jsx)(s.code,{children:"group_type"})," and ",(0,t.jsx)(s.code,{children:"group_name"})," in the specific request refer always to endpoint groups (e.g.. ",(0,t.jsx)(s.code,{children:"SITES"}),")."]}),"\n",(0,t.jsxs)(s.p,{children:[(0,t.jsx)(s.code,{children:"strict"})," when strict mode is ON, the response will contain only the latest grouped by endpoint_group/host/service/metric."]}),"\n",(0,t.jsx)(s.h4,{id:"headers-1",children:"Headers"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{children:"x-api-key: shared_key_value\nAccept: application/json or application/xml\n"})}),"\n",(0,t.jsx)(s.h4,{id:"response-code-1",children:"Response Code"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{children:"Status: 200 OK\n"})}),"\n",(0,t.jsx)(s.h3,{id:"response-body-1",children:"Response body"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{className:"language-json",children:'{\n    "status": {\n        "message": "type of message received",\n        "code": "return code"\n    },\n    "data": {\n        "metric_data": [\n            {\n                "endpoint_group": "name of endpoint group",\n                "service": "name of service",\n                "endpoint": "name of endpoint",\n                "metric": "name of metric",\n                "timestamp": "2018-06-22T11:55:44Z",\n                "status": "OK || WARNING || CRITICAL || MISSING || UNKNOWN",\n                "summary": "summary of the metric message - generated by the monitoring engine (nagios)",\n                "message": "body of nagios generated message"\n            },\n         ]\n      }\n}\n'})}),"\n",(0,t.jsx)(s.h6,{id:"example-request-1",children:"Example Request:"}),"\n",(0,t.jsx)(s.p,{children:"URL:"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{children:"latest/Report_B/TENANT_SITES/EL-01-AUTH?date=2015-05-01T00:00:00Z&filter=non-ok&limit=10\n"})}),"\n",(0,t.jsx)(s.p,{children:"Headers:"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{children:"x-api-key: secret_key_value\nAccept: application/json or application/xml\n\n"})}),"\n",(0,t.jsx)(s.h6,{id:"example-response-1",children:"Example Response:"}),"\n",(0,t.jsx)(s.p,{children:"Code:"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{children:"Status: 200 OK\n"})}),"\n",(0,t.jsx)(s.p,{children:"Response body:"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{className:"language-json",children:'{\n"status": {\n"message": "application/json",\n"code": "200"\n},\n"data": {\n"metric_data": [\n {\n  "endpoint_group": "EL-01-AUTH",\n  "service": "someService",\n  "endpoint": "someservice.example.gr",\n  "metric": "someService-FileTransfer",\n  "timestamp": "2015-05-01T05:00:00Z",\n  "status": "WARNING",\n  "summary": "someService status is warning",\n  "message": "someService data upload test return value of warning"\n },\n {\n  "endpoint_group": "EL-01-AUTH",\n  "service": "someService",\n  "endpoint": "someservice.example.gr",\n  "metric": "someService-FileTransfer",\n  "timestamp": "2015-05-01T01:00:00Z",\n  "status": "UNKNOWN",\n  "summary": "someService status is CRITICAL",\n  "message": "someService data upload test failed"\n }\n]\n}\n}\n'})})]})}function h(e={}){const{wrapper:s}={...(0,r.R)(),...e.components};return s?(0,t.jsx)(s,{...e,children:(0,t.jsx)(o,{...e})}):o(e)}},8453:(e,s,n)=>{n.d(s,{R:()=>d,x:()=>l});var t=n(6540);const r={},i=t.createContext(r);function d(e){const s=t.useContext(i);return t.useMemo((function(){return"function"==typeof e?e(s):{...s,...e}}),[s,e])}function l(e){let s;return s=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:d(e.components),t.createElement(i.Provider,{value:s},e.children)}}}]);