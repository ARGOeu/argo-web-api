"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[7548],{5409:(e,t,s)=>{s.r(t),s.d(t,{assets:()=>d,contentTitle:()=>i,default:()=>a,frontMatter:()=>r,metadata:()=>l,toc:()=>p});var o=s(4848),n=s(8453);const r={id:"topology_stats",title:"Topology Statistics",sidebar_position:4},i=void 0,l={id:"topology/topology_stats",title:"Topology Statistics",description:"API calls for retrieving topology statistics per report",source:"@site/docs/topology/topology_stats.md",sourceDirName:"topology",slug:"/topology/topology_stats",permalink:"/argo-web-api/docs/topology/topology_stats",draft:!1,unlisted:!1,tags:[],version:"current",sidebarPosition:4,frontMatter:{id:"topology_stats",title:"Topology Statistics",sidebar_position:4},sidebar:"tutorialSidebar",previous:{title:"Topology Groups",permalink:"/argo-web-api/docs/topology/topology_groups"},next:{title:"Topology Tags & Values",permalink:"/argo-web-api/docs/topology/topology_tags"}},d={},p=[{value:"API calls for retrieving topology statistics per report",id:"api-calls-for-retrieving-topology-statistics-per-report",level:2},{value:"[GET]: List topology statistics",id:"1",level:2},{value:"Input",id:"input",level:3},{value:"List All topology statistics",id:"list-all-topology-statistics",level:5},{value:"Path Parameters",id:"path-parameters",level:4},{value:"Url Parameters",id:"url-parameters",level:4},{value:"Headers",id:"headers",level:4},{value:"Response Code",id:"response-code",level:4},{value:"Response body",id:"response-body",level:3},{value:"Example Request:",id:"example-request",level:6},{value:"Example Response:",id:"example-response",level:6}];function c(e){const t={code:"code",h2:"h2",h3:"h3",h4:"h4",h5:"h5",h6:"h6",p:"p",pre:"pre",table:"table",tbody:"tbody",td:"td",th:"th",thead:"thead",tr:"tr",...(0,n.R)(),...e.components};return(0,o.jsxs)(o.Fragment,{children:[(0,o.jsx)(t.h2,{id:"api-calls-for-retrieving-topology-statistics-per-report",children:"API calls for retrieving topology statistics per report"}),"\n",(0,o.jsxs)(t.table,{children:[(0,o.jsx)(t.thead,{children:(0,o.jsxs)(t.tr,{children:[(0,o.jsx)(t.th,{children:"Name"}),(0,o.jsx)(t.th,{children:"Description"}),(0,o.jsx)(t.th,{children:"Shortcut"})]})}),(0,o.jsx)(t.tbody,{children:(0,o.jsxs)(t.tr,{children:[(0,o.jsx)(t.td,{children:"GET: List topology statistics"}),(0,o.jsx)(t.td,{children:"List number of groups, endpoint groups and services ."}),(0,o.jsx)(t.td,{children:(0,o.jsx)("a",{href:"#1",children:"Description"})})]})})]}),"\n",(0,o.jsx)(t.h2,{id:"1",children:"[GET]: List topology statistics"}),"\n",(0,o.jsx)(t.p,{children:"This method may be used to retrieve topology statistics for a specific report. Topology statistics include number of groups, endpoint groups and services included in the report"}),"\n",(0,o.jsx)(t.h3,{id:"input",children:"Input"}),"\n",(0,o.jsx)(t.h5,{id:"list-all-topology-statistics",children:"List All topology statistics"}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{children:"/topology/stats/{report}/?[date]\n"})}),"\n",(0,o.jsx)(t.h4,{id:"path-parameters",children:"Path Parameters"}),"\n",(0,o.jsxs)(t.table,{children:[(0,o.jsx)(t.thead,{children:(0,o.jsxs)(t.tr,{children:[(0,o.jsx)(t.th,{children:"Type"}),(0,o.jsx)(t.th,{children:"Description"}),(0,o.jsx)(t.th,{children:"Required"}),(0,o.jsx)(t.th,{children:"Default value"})]})}),(0,o.jsx)(t.tbody,{children:(0,o.jsxs)(t.tr,{children:[(0,o.jsx)(t.td,{children:(0,o.jsx)(t.code,{children:"report"})}),(0,o.jsx)(t.td,{children:"name of the report used"}),(0,o.jsx)(t.td,{children:"YES"}),(0,o.jsx)(t.td,{})]})})]}),"\n",(0,o.jsx)(t.h4,{id:"url-parameters",children:"Url Parameters"}),"\n",(0,o.jsxs)(t.table,{children:[(0,o.jsx)(t.thead,{children:(0,o.jsxs)(t.tr,{children:[(0,o.jsx)(t.th,{children:"Type"}),(0,o.jsx)(t.th,{children:"Description"}),(0,o.jsx)(t.th,{children:"Required"}),(0,o.jsx)(t.th,{children:"Default value"})]})}),(0,o.jsx)(t.tbody,{children:(0,o.jsxs)(t.tr,{children:[(0,o.jsx)(t.td,{children:(0,o.jsx)(t.code,{children:"date"})}),(0,o.jsx)(t.td,{children:"target a specific data"}),(0,o.jsx)(t.td,{children:"NO"}),(0,o.jsx)(t.td,{children:"today's date"})]})})]}),"\n",(0,o.jsx)(t.h4,{id:"headers",children:"Headers"}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{children:"x-api-key: secret_key_value\nAccept: application/json\n"})}),"\n",(0,o.jsx)(t.h4,{id:"response-code",children:"Response Code"}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{children:"Status: 200 OK\n"})}),"\n",(0,o.jsx)(t.h3,{id:"response-body",children:"Response body"}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-json",children:'{\n    "status": {\n        "message": "application/json",\n        "code": "200"\n    },\n    "data": {\n        "group_count": 1,\n        "group_type": "type of top-level groups in report",\n        "group_list": ["list of top level groups"],\n        "endpoint_group_count": 1,\n        "endpoint_group_type": "type of endpoint groups in report",\n        "endpoint_group_list": ["list of endpoint groups"],\n        "service_count": 1,\n        "service_list": ["list of available services"]\n    }\n}\n'})}),"\n",(0,o.jsx)(t.h6,{id:"example-request",children:"Example Request:"}),"\n",(0,o.jsx)(t.p,{children:"URL:"}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{children:"latest/Report_B/?date=2015-05-01\n"})}),"\n",(0,o.jsx)(t.p,{children:"Headers:"}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{children:"x-api-key: secret_key_value\nAccept: application/json\n"})}),"\n",(0,o.jsx)(t.h6,{id:"example-response",children:"Example Response:"}),"\n",(0,o.jsx)(t.p,{children:"Code:"}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{children:"Status: 200 OK\n"})}),"\n",(0,o.jsx)(t.p,{children:"Response body:"}),"\n",(0,o.jsx)(t.pre,{children:(0,o.jsx)(t.code,{className:"language-json",children:'{\n    "status": {\n        "message": "application/json",\n        "code": "200"\n    },\n    "data": {\n        "group_count": 2,\n        "group_type": "PROJECTS",\n        "group_list": ["PROJECT_A", "PROJECT_B"],\n        "endpoint_group_count": 3,\n        "endpoint_group_type": "SERVICEGROUPS",\n        "endpoint_group_list": ["SGROUP_A", "SGROUP_B", "SGROUP_C"],\n        "service_count": 8,\n        "service_list": [\n            "service_type_1",\n            "service_type_2",\n            "service_type_3",\n            "service_type_4",\n            "service_type_5",\n            "service_type_6",\n            "service_type_7",\n            "service_type_8"\n        ]\n    }\n}\n'})})]})}function a(e={}){const{wrapper:t}={...(0,n.R)(),...e.components};return t?(0,o.jsx)(t,{...e,children:(0,o.jsx)(c,{...e})}):c(e)}},8453:(e,t,s)=>{s.d(t,{R:()=>i,x:()=>l});var o=s(6540);const n={},r=o.createContext(n);function i(e){const t=o.useContext(r);return o.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function l(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(n):e.components||n:i(e.components),o.createElement(r.Provider,{value:t},e.children)}}}]);