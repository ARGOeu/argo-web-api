"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[390],{9054:(e,n,s)=>{s.r(n),s.d(n,{assets:()=>c,contentTitle:()=>d,default:()=>h,frontMatter:()=>r,metadata:()=>o,toc:()=>a});var t=s(4848),i=s(8453);const r={id:"recomputations",title:"Recomputation Requests",sidebar_position:7},d=void 0,o={id:"results/recomputations",title:"Recomputation Requests",description:"API Calls for listing existing and creating new recomputation requests",source:"@site/docs/results/recomputations.md",sourceDirName:"results",slug:"/results/recomputations",permalink:"/argo-web-api/docs/results/recomputations",draft:!1,unlisted:!1,tags:[],version:"current",sidebarPosition:7,frontMatter:{id:"recomputations",title:"Recomputation Requests",sidebar_position:7},sidebar:"tutorialSidebar",previous:{title:"Trends",permalink:"/argo-web-api/docs/results/trends"},next:{title:"Validation & Errors",permalink:"/argo-web-api/docs/category/validation--errors"}},c={},a=[{value:"API Calls for listing existing and creating new recomputation requests",id:"api-calls-for-listing-existing-and-creating-new-recomputation-requests",level:2},{value:"[GET]: List Recomputation Requests",id:"1",level:2},{value:"Input",id:"input",level:3},{value:"Optional Query Parameters",id:"optional-query-parameters",level:4},{value:"Request headers",id:"request-headers",level:4},{value:"Response",id:"response",level:4},{value:"Response body",id:"response-body",level:4},{value:"Example Request #2",id:"example-request-2",level:3},{value:"Response",id:"response-1",level:4},{value:"Response body",id:"response-body-1",level:4},{value:"Example Request #3",id:"example-request-3",level:3},{value:"Response",id:"response-2",level:4},{value:"Response body",id:"response-body-2",level:4},{value:"[GET]: Get specific recomputation request by id",id:"2",level:2},{value:"Input",id:"input-1",level:3},{value:"Request headers",id:"request-headers-1",level:4},{value:"Response",id:"response-3",level:4},{value:"Response body",id:"response-body-3",level:4},{value:"[POST]: Create a new recomputation request",id:"3",level:2},{value:"Input",id:"input-2",level:3},{value:"Request headers",id:"request-headers-2",level:4},{value:"Parameters",id:"parameters",level:4},{value:"Response",id:"response-4",level:3},{value:"[DELETE]: Delete a specific recomputation",id:"4",level:2},{value:"Request headers",id:"request-headers-3",level:4},{value:"Response",id:"response-5",level:3},{value:"[POST]: Change status of recomputation",id:"5",level:2},{value:"Request headers",id:"request-headers-4",level:4},{value:"POST body",id:"post-body",level:3},{value:"Response",id:"response-6",level:3},{value:"Response body",id:"response-body-4",level:4},{value:"[DELETE]: Reset status of a specific recomputation",id:"6",level:2},{value:"Request headers",id:"request-headers-5",level:4},{value:"Response",id:"response-7",level:3},{value:"Response body",id:"response-body-5",level:4}];function l(e){const n={a:"a",code:"code",h1:"h1",h2:"h2",h3:"h3",h4:"h4",li:"li",p:"p",pre:"pre",table:"table",tbody:"tbody",td:"td",th:"th",thead:"thead",tr:"tr",ul:"ul",...(0,i.R)(),...e.components};return(0,t.jsxs)(t.Fragment,{children:[(0,t.jsx)(n.h2,{id:"api-calls-for-listing-existing-and-creating-new-recomputation-requests",children:"API Calls for listing existing and creating new recomputation requests"}),"\n",(0,t.jsxs)(n.table,{children:[(0,t.jsx)(n.thead,{children:(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.th,{children:"Name"}),(0,t.jsx)(n.th,{children:"Description"}),(0,t.jsx)(n.th,{children:"Shortcut"})]})}),(0,t.jsxs)(n.tbody,{children:[(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.td,{children:"GET: List Recomputation Requests"}),(0,t.jsx)(n.td,{children:"This method can be used to retrieve a list of current Recomputation requests."}),(0,t.jsx)(n.td,{children:(0,t.jsx)(n.a,{href:"#1",children:" Description"})})]}),(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.td,{children:"GET: Get a specific recomputation by id"}),(0,t.jsx)(n.td,{children:"This method can be used to retrieve a specific recomputation by id"}),(0,t.jsx)(n.td,{children:(0,t.jsx)(n.a,{href:"#2",children:" Description"})})]}),(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.td,{children:"POST: Create a new recomputation request"}),(0,t.jsx)(n.td,{children:"This method can be used to insert a new recomputation request onto the Compute Engine."}),(0,t.jsx)(n.td,{children:(0,t.jsx)(n.a,{href:"#3",children:" Description"})})]}),(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.td,{children:"DELETE: Delete a specific recomputation"}),(0,t.jsx)(n.td,{children:"This method can be used to delete a specific recomputation."}),(0,t.jsx)(n.td,{children:(0,t.jsx)(n.a,{href:"#4",children:" Description"})})]}),(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.td,{children:"POST: change status"}),(0,t.jsx)(n.td,{children:"This method can be used to change status of a specific recomputation."}),(0,t.jsx)(n.td,{children:(0,t.jsx)(n.a,{href:"#5",children:" Description"})})]}),(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.td,{children:"DELETE: Reset status of recomputation"}),(0,t.jsx)(n.td,{children:"This method can be used to reset status of a specific recomputation."}),(0,t.jsx)(n.td,{children:(0,t.jsx)(n.a,{href:"#6",children:" Description"})})]})]})]}),"\n",(0,t.jsx)(n.h2,{id:"1",children:"[GET]: List Recomputation Requests"}),"\n",(0,t.jsx)(n.p,{children:"This method can be used to retrieve a list of current Recomputation requests."}),"\n",(0,t.jsx)(n.h3,{id:"input",children:"Input"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{children:"GET /recomputations\n"})}),"\n",(0,t.jsx)(n.h4,{id:"optional-query-parameters",children:"Optional Query Parameters"}),"\n",(0,t.jsxs)(n.table,{children:[(0,t.jsx)(n.thead,{children:(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.th,{children:"Type"}),(0,t.jsx)(n.th,{children:"Description"}),(0,t.jsx)(n.th,{children:"Required"})]})}),(0,t.jsxs)(n.tbody,{children:[(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.td,{children:(0,t.jsx)(n.code,{children:"report"})}),(0,t.jsx)(n.td,{children:"Filter recomputations by report name"}),(0,t.jsx)(n.td,{children:"NO"})]}),(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.td,{children:(0,t.jsx)(n.code,{children:"date"})}),(0,t.jsx)(n.td,{children:"Specific date to retrieve all relevant recomputations that their period include this date"}),(0,t.jsx)(n.td,{children:"NO"})]})]})]}),"\n",(0,t.jsx)(n.h4,{id:"request-headers",children:"Request headers"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,t.jsx)(n.h4,{id:"response",children:"Response"}),"\n",(0,t.jsxs)(n.p,{children:["Headers: ",(0,t.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,t.jsx)(n.h4,{id:"response-body",children:"Response body"}),"\n",(0,t.jsx)(n.p,{children:"Json Response"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-json",children:'{\n"root": [\n     {\n          "id": "56db43ee-f331-46ca-b0fd-4555b4aa1cfc",\n          "requester_name": "John Doe",\n          "requester_email": "JohnDoe@foo.com",\n          "reason": "power cuts",\n          "start_time": "2015-01-10T12:00:00Z",\n          "end_time": "2015-01-30T23:00:00Z",\n          "report": "Critical",\n          "exclude": [\n           "Gluster"\n          ],\n          "status": "running",\n          "timestamp": "2015-02-01T14:58:40",\n          "history": [\n              { \n                  "status": "pending", \n                  "timestamp" : "2015-02-01T14:58:40"\n              },\n              { \n                  "status": "approved", \n                  "timestamp" : "2015-02-02T08:58:40"\n              },\n              { \n                  "status": "running", \n                  "timestamp" : "2015-02-02T09:10:40"\n              },\n\n          ]\n     },\n     {\n          "id": "f68b43ee-f331-46ca-b0fd-4555b4aa1cfc",\n          "requester_name": "John Doe",\n          "requester_email": "JohnDoe@foo.com",\n          "reason": "power cuts",\n          "start_time": "2015-03-10T12:00:00Z",\n          "end_time": "2015-03-30T23:00:00Z",\n          "report": "OPS-Critical",\n          "exclude": [\n           "Gluster"\n          ],\n          "status": "running",\n          "timestamp": "2015-02-01T14:58:40",\n          "history": [\n              { \n                  "status": "pending", \n                  "timestamp" : "2015-04-01T14:58:40"\n              },\n              { \n                  "status": "approved", \n                  "timestamp" : "2015-04-02T08:58:40"\n              },\n              { \n                  "status": "running", \n                  "timestamp" : "2015-04-02T09:10:40"\n              },\n\n          ]\n     }\n ]\n}\n'})}),"\n",(0,t.jsx)(n.h3,{id:"example-request-2",children:"Example Request #2"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{children:"GET /recomputations?date=2015-03-15\n"})}),"\n",(0,t.jsx)(n.h4,{id:"response-1",children:"Response"}),"\n",(0,t.jsxs)(n.p,{children:["Headers: ",(0,t.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,t.jsx)(n.h4,{id:"response-body-1",children:"Response body"}),"\n",(0,t.jsx)(n.p,{children:"Json Response"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-json",children:'{\n"root": [\n     {\n          "id": "f68b43ee-f331-46ca-b0fd-4555b4aa1cfc",\n          "requester_name": "John Doe",\n          "requester_email": "JohnDoe@foo.com",\n          "reason": "power cuts",\n          "start_time": "2015-03-10T12:00:00Z",\n          "end_time": "2015-03-30T23:00:00Z",\n          "report": "OPS-Critical",\n          "exclude": [\n           "Gluster"\n          ],\n          "status": "running",\n          "timestamp": "2015-02-01T14:58:40",\n          "history": [\n              { \n                  "status": "pending", \n                  "timestamp" : "2015-04-01T14:58:40"\n              },\n              { \n                  "status": "approved", \n                  "timestamp" : "2015-04-02T08:58:40"\n              },\n              { \n                  "status": "running", \n                  "timestamp" : "2015-04-02T09:10:40"\n              },\n\n          ]\n     }\n ]\n}\n'})}),"\n",(0,t.jsx)(n.h3,{id:"example-request-3",children:"Example Request #3"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{children:"GET /recomputations?report=OPS-Critical\n"})}),"\n",(0,t.jsx)(n.h4,{id:"response-2",children:"Response"}),"\n",(0,t.jsxs)(n.p,{children:["Headers: ",(0,t.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,t.jsx)(n.h4,{id:"response-body-2",children:"Response body"}),"\n",(0,t.jsx)(n.p,{children:"Json Response"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-json",children:'{\n"root": [\n     {\n          "id": "f68b43ee-f331-46ca-b0fd-4555b4aa1cfc",\n          "requester_name": "John Doe",\n          "requester_email": "JohnDoe@foo.com",\n          "reason": "power cuts",\n          "start_time": "2015-03-10T12:00:00Z",\n          "end_time": "2015-03-30T23:00:00Z",\n          "report": "OPS-Critical",\n          "exclude": [\n           "Gluster"\n          ],\n          "status": "running",\n          "timestamp": "2015-02-01T14:58:40",\n          "history": [\n              { \n                  "status": "pending", \n                  "timestamp" : "2015-04-01T14:58:40"\n              },\n              { \n                  "status": "approved", \n                  "timestamp" : "2015-04-02T08:58:40"\n              },\n              { \n                  "status": "running", \n                  "timestamp" : "2015-04-02T09:10:40"\n              },\n\n          ]\n     }\n ]\n}\n'})}),"\n",(0,t.jsx)(n.h2,{id:"2",children:"[GET]: Get specific recomputation request by id"}),"\n",(0,t.jsx)(n.p,{children:"This method can be used to retrieve a specific recomputation request by its id"}),"\n",(0,t.jsx)(n.h3,{id:"input-1",children:"Input"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{children:"GET /recomputations/{ID}\n"})}),"\n",(0,t.jsx)(n.h4,{id:"request-headers-1",children:"Request headers"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,t.jsx)(n.h4,{id:"response-3",children:"Response"}),"\n",(0,t.jsxs)(n.p,{children:["Headers: ",(0,t.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,t.jsx)(n.h4,{id:"response-body-3",children:"Response body"}),"\n",(0,t.jsx)(n.p,{children:"Json Response"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-json",children:'{\n"root": [\n     {\n          "id": "56db43ee-f331-46ca-b0fd-4555b4aa1cfc",\n          "requester_name": "John Doe",\n          "requester_email": "JohnDoe@foo.com",\n          "reason": "power cuts",\n          "start_time": "2015-01-10T12:00:00Z",\n          "end_time": "2015-01-30T23:00:00Z",\n          "report": "Critical",\n          "exclude": [\n           "Gluster"\n          ],\n          "status": "running",\n          "timestamp": "2015-02-01T14:58:40",\n          "history": [\n              { \n                  "status": "pending", \n                  "timestamp" : "2015-02-01T14:58:40"\n              },\n              { \n                  "status": "approved", \n                  "timestamp" : "2015-02-02T08:58:40"\n              },\n              { \n                  "status": "running", \n                  "timestamp" : "2015-02-02T09:10:40"\n              },\n\n          ]\n     }\n'})}),"\n",(0,t.jsx)(n.h2,{id:"3",children:"[POST]: Create a new recomputation request"}),"\n",(0,t.jsx)(n.p,{children:"This method can be used to insert a new recomputation request onto the Compute Engine."}),"\n",(0,t.jsx)(n.h3,{id:"input-2",children:"Input"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{children:"POST /recomputations\n"})}),"\n",(0,t.jsx)(n.h4,{id:"request-headers-2",children:"Request headers"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,t.jsx)(n.h4,{id:"parameters",children:"Parameters"}),"\n",(0,t.jsxs)(n.table,{children:[(0,t.jsx)(n.thead,{children:(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.th,{children:"Type"}),(0,t.jsx)(n.th,{children:"Description"}),(0,t.jsx)(n.th,{children:"Required"}),(0,t.jsx)(n.th,{children:"Default value"})]})}),(0,t.jsxs)(n.tbody,{children:[(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.td,{children:(0,t.jsx)(n.code,{children:"start_time"})}),(0,t.jsx)(n.td,{children:"UTC time in W3C format"}),(0,t.jsx)(n.td,{children:"YES"}),(0,t.jsx)(n.td,{})]}),(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.td,{children:(0,t.jsx)(n.code,{children:"end_time"})}),(0,t.jsx)(n.td,{children:"UTC time in W3C format"}),(0,t.jsx)(n.td,{children:"YES"}),(0,t.jsx)(n.td,{})]}),(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.td,{children:(0,t.jsx)(n.code,{children:"reason"})}),(0,t.jsx)(n.td,{children:"Explain the need for a recomputation"}),(0,t.jsx)(n.td,{children:"YES"}),(0,t.jsx)(n.td,{})]}),(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.td,{children:(0,t.jsx)(n.code,{children:"requester_name"})}),(0,t.jsx)(n.td,{children:"The name of the person submitting the recomputation"}),(0,t.jsx)(n.td,{children:"YES"}),(0,t.jsx)(n.td,{})]}),(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.td,{children:(0,t.jsx)(n.code,{children:"requester_email"})}),(0,t.jsx)(n.td,{children:"The email of the person submitting the recomputation"}),(0,t.jsx)(n.td,{children:"YES"}),(0,t.jsx)(n.td,{})]}),(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.td,{children:(0,t.jsx)(n.code,{children:"report"})}),(0,t.jsx)(n.td,{children:"Report for which the recomputation is requested"}),(0,t.jsx)(n.td,{children:"YES"}),(0,t.jsx)(n.td,{})]}),(0,t.jsxs)(n.tr,{children:[(0,t.jsx)(n.td,{children:(0,t.jsx)(n.code,{children:"exclude"})}),(0,t.jsx)(n.td,{children:"Groups to be excluded from recomputation. If more than one group are to be excluded use the parameter as many times as needed within the same API call"}),(0,t.jsx)(n.td,{children:"NO"}),(0,t.jsx)(n.td,{})]})]})]}),"\n",(0,t.jsx)(n.h3,{id:"response-4",children:"Response"}),"\n",(0,t.jsxs)(n.p,{children:["Headers: ",(0,t.jsx)(n.code,{children:"Status: 201 Created"})]}),"\n",(0,t.jsx)(n.h2,{id:"4",children:"[DELETE]: Delete a specific recomputation"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{children:"DELETE /recomputations/{ID}\n"})}),"\n",(0,t.jsx)(n.h4,{id:"request-headers-3",children:"Request headers"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,t.jsx)(n.h3,{id:"response-5",children:"Response"}),"\n",(0,t.jsx)(n.p,{children:(0,t.jsx)(n.code,{children:"Status 200 OK"})}),"\n",(0,t.jsx)(n.h2,{id:"5",children:"[POST]: Change status of recomputation"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{children:"POST /recomputations/{ID}/status\n"})}),"\n",(0,t.jsx)(n.h4,{id:"request-headers-4",children:"Request headers"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,t.jsx)(n.h3,{id:"post-body",children:"POST body"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-json",children:'{\n  "status" : "approved"\n}\n'})}),"\n",(0,t.jsx)(n.p,{children:"Eligible recomputation status values:"}),"\n",(0,t.jsxs)(n.ul,{children:["\n",(0,t.jsx)(n.li,{children:"pending"}),"\n",(0,t.jsx)(n.li,{children:"approved"}),"\n",(0,t.jsx)(n.li,{children:"rejected"}),"\n",(0,t.jsx)(n.li,{children:"running"}),"\n",(0,t.jsx)(n.li,{children:"done"}),"\n"]}),"\n",(0,t.jsxs)(n.p,{children:["If recomputation status input not in eligible values the api will respond with status code ",(0,t.jsx)(n.code,{children:"404"}),":",(0,t.jsx)(n.code,{children:"conflict"})]}),"\n",(0,t.jsx)(n.h3,{id:"response-6",children:"Response"}),"\n",(0,t.jsx)(n.p,{children:(0,t.jsx)(n.code,{children:"Status 200 OK"})}),"\n",(0,t.jsx)(n.h4,{id:"response-body-4",children:"Response body"}),"\n",(0,t.jsx)(n.p,{children:"Json Response"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-json",children:'{\n "status": {\n  "message": "Recomputation status updated successfully to: approved",\n  "code": "200"\n }\n}\n'})}),"\n",(0,t.jsx)(n.h2,{id:"6",children:"[DELETE]: Reset status of a specific recomputation"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{children:"DELETE /recomputations/{ID}/status\n"})}),"\n",(0,t.jsx)(n.h4,{id:"request-headers-5",children:"Request headers"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,t.jsx)(n.h3,{id:"response-7",children:"Response"}),"\n",(0,t.jsx)(n.p,{children:(0,t.jsx)(n.code,{children:"Status 200 OK"})}),"\n",(0,t.jsx)(n.h4,{id:"response-body-5",children:"Response body"}),"\n",(0,t.jsx)(n.p,{children:"Json Response"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-json",children:'{\n "status": {\n  "message": "Recomputation status reset to: pending",\n  "code": "200"\n }\n}\n'})}),"\n",(0,t.jsx)(n.h1,{id:"recomputations-that-exclude-metrics",children:"Recomputations that exclude metrics"}),"\n",(0,t.jsx)(n.p,{children:"There is also the ability to run a recomputation and exclude specific metrics. During the recomputation period the metrics that are considered excluded don't take place into any operation or aggregation thus they don't affect their endpoints at all."}),"\n",(0,t.jsx)(n.p,{children:'To declare a recomputation that excludes metric you must use the special field "exclude_metrics" in the recomputation and add an array of metrics to be excluded (You can limit the scope also by "group", "service" and "hostname")'}),"\n",(0,t.jsx)(n.p,{children:"For example:"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-json",children:'{\n   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e777",\n   "requester_name": "John Doe",\n   "requester_email": "johndoe@example.com",\n   "reason": "issue with metric checks",\n   "start_time": "2022-01-10T12:00:00Z",\n   "end_time": "2022-01-10T23:00:00Z",\n   "report": "Default",\n   "exclude_metrics": [\n    {\n     "metric": "check-1"\n    },\n    {\n     "metric": "check-2",\n     "hostname": "host1.example.com"\n    },\n    {\n     "metric": "check-3",\n     "group": "Affected-Site"\n    }\n   ]\n  }\n'})}),"\n",(0,t.jsxs)(n.p,{children:["If you specify a rule that includes only a ",(0,t.jsx)(n.code,{children:"metric"})," then this type of metric will be excluded globally from all endpoints and groups\nIf you specify a rule that includes a ",(0,t.jsx)(n.code,{children:"metric"})," and another field such as ",(0,t.jsx)(n.code,{children:"hostname"}),", ",(0,t.jsx)(n.code,{children:"service"})," or ",(0,t.jsx)(n.code,{children:"group"})," then the rule is scoped accordingly to a specific group or service type or hostname and the metric that belongs there. The field ",(0,t.jsx)(n.code,{children:"metric"})," is mandatory."]}),"\n",(0,t.jsx)(n.h1,{id:"recomputations-that-exclude-monitoring-sources",children:"Recomputations that exclude monitoring sources"}),"\n",(0,t.jsx)(n.p,{children:"There is also the ability to run a recomputation and exclude a monitoring source (e.g. specific monitoring box). This is especially usefull in HA situations where one of the available monitoring sources might have issues for a specific period of time."}),"\n",(0,t.jsx)(n.p,{children:"For example:"}),"\n",(0,t.jsx)(n.pre,{children:(0,t.jsx)(n.code,{className:"language-json",children:'{\n   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e777",\n   "requester_name": "John Doe",\n   "requester_email": "johndoe@example.com",\n   "reason": "issue with metric checks",\n   "start_time": "2022-01-10T12:00:00Z",\n   "end_time": "2022-01-10T23:00:00Z",\n   "report": "Default",\n   "exclude_monitoring_source": [\n    {\n        "host":"monitoring_node01.example.foo",\n        "start_time": "2022-01-10T12:00:00Z",\n        "end_time": "2022-01-10T23:00:00Z"\n    }\n   ]\n}\n'})})]})}function h(e={}){const{wrapper:n}={...(0,i.R)(),...e.components};return n?(0,t.jsx)(n,{...e,children:(0,t.jsx)(l,{...e})}):l(e)}},8453:(e,n,s)=>{s.d(n,{R:()=>d,x:()=>o});var t=s(6540);const i={},r=t.createContext(i);function d(e){const n=t.useContext(r);return t.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function o(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(i):e.components||i:d(e.components),t.createElement(r.Provider,{value:n},e.children)}}}]);