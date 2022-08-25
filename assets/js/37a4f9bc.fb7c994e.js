"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[150],{3905:(e,t,r)=>{r.d(t,{Zo:()=>u,kt:()=>g});var n=r(7294);function a(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function o(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function l(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?o(Object(r),!0).forEach((function(t){a(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):o(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function p(e,t){if(null==e)return{};var r,n,a=function(e,t){if(null==e)return{};var r,n,a={},o=Object.keys(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||(a[r]=e[r]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(a[r]=e[r])}return a}var i=n.createContext({}),s=function(e){var t=n.useContext(i),r=t;return e&&(r="function"==typeof e?e(t):l(l({},t),e)),r},u=function(e){var t=s(e.components);return n.createElement(i.Provider,{value:t},e.children)},c={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},d=n.forwardRef((function(e,t){var r=e.components,a=e.mdxType,o=e.originalType,i=e.parentName,u=p(e,["components","mdxType","originalType","parentName"]),d=s(r),g=a,y=d["".concat(i,".").concat(g)]||d[g]||c[g]||o;return r?n.createElement(y,l(l({ref:t},u),{},{components:r})):n.createElement(y,l({ref:t},u))}));function g(e,t){var r=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var o=r.length,l=new Array(o);l[0]=d;var p={};for(var i in t)hasOwnProperty.call(t,i)&&(p[i]=t[i]);p.originalType=e,p.mdxType="string"==typeof e?e:a,l[1]=p;for(var s=2;s<o;s++)l[s]=r[s];return n.createElement.apply(null,l)}return n.createElement.apply(null,r)}d.displayName="MDXCreateElement"},4289:(e,t,r)=>{r.r(t),r.d(t,{assets:()=>i,contentTitle:()=>l,default:()=>c,frontMatter:()=>o,metadata:()=>p,toc:()=>s});var n=r(7462),a=(r(7294),r(3905));const o={id:"topology_stats",title:"Topology Statistics",sidebar_position:4},l=void 0,p={unversionedId:"topology/topology_stats",id:"topology/topology_stats",title:"Topology Statistics",description:"API calls for retrieving topology statistics per report",source:"@site/docs/topology/topology_stats.md",sourceDirName:"topology",slug:"/topology/topology_stats",permalink:"/argo-web-api/docs/topology/topology_stats",draft:!1,tags:[],version:"current",sidebarPosition:4,frontMatter:{id:"topology_stats",title:"Topology Statistics",sidebar_position:4},sidebar:"tutorialSidebar",previous:{title:"Topology Groups",permalink:"/argo-web-api/docs/topology/topology_groups"},next:{title:"Topology Tags & Values",permalink:"/argo-web-api/docs/topology/topology_tags"}},i={},s=[{value:"API calls for retrieving topology statistics per report",id:"api-calls-for-retrieving-topology-statistics-per-report",level:2},{value:"GET: List topology statistics",id:"get-list-topology-statistics",level:2},{value:"Input",id:"input",level:3},{value:"List All topology statistics",id:"list-all-topology-statistics",level:5},{value:"Path Parameters",id:"path-parameters",level:4},{value:"Url Parameters",id:"url-parameters",level:4},{value:"Headers",id:"headers",level:4},{value:"Response Code",id:"response-code",level:4},{value:"Response body",id:"response-body",level:3},{value:"Example Request:",id:"example-request",level:6},{value:"Example Response:",id:"example-response",level:6}],u={toc:s};function c(e){let{components:t,...r}=e;return(0,a.kt)("wrapper",(0,n.Z)({},u,r,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h2",{id:"api-calls-for-retrieving-topology-statistics-per-report"},"API calls for retrieving topology statistics per report"),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:null},"Name"),(0,a.kt)("th",{parentName:"tr",align:null},"Description"),(0,a.kt)("th",{parentName:"tr",align:null},"Shortcut"))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},"GET: List topology statistics"),(0,a.kt)("td",{parentName:"tr",align:null},"List number of groups, endpoint groups and services ."),(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{href:"#1"},"Description"))))),(0,a.kt)("a",{id:"1"}),(0,a.kt)("h2",{id:"get-list-topology-statistics"},"[GET]",": List topology statistics"),(0,a.kt)("p",null,"This method may be used to retrieve topology statistics for a specific report. Topology statistics include number of groups, endpoint groups and services included in the report"),(0,a.kt)("h3",{id:"input"},"Input"),(0,a.kt)("h5",{id:"list-all-topology-statistics"},"List All topology statistics"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"/topology/stats/{report}/?[date]\n")),(0,a.kt)("h4",{id:"path-parameters"},"Path Parameters"),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:null},"Type"),(0,a.kt)("th",{parentName:"tr",align:null},"Description"),(0,a.kt)("th",{parentName:"tr",align:null},"Required"),(0,a.kt)("th",{parentName:"tr",align:null},"Default value"))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("inlineCode",{parentName:"td"},"report")),(0,a.kt)("td",{parentName:"tr",align:null},"name of the report used"),(0,a.kt)("td",{parentName:"tr",align:null},"YES"),(0,a.kt)("td",{parentName:"tr",align:null})))),(0,a.kt)("h4",{id:"url-parameters"},"Url Parameters"),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:null},"Type"),(0,a.kt)("th",{parentName:"tr",align:null},"Description"),(0,a.kt)("th",{parentName:"tr",align:null},"Required"),(0,a.kt)("th",{parentName:"tr",align:null},"Default value"))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("inlineCode",{parentName:"td"},"date")),(0,a.kt)("td",{parentName:"tr",align:null},"target a specific data"),(0,a.kt)("td",{parentName:"tr",align:null},"NO"),(0,a.kt)("td",{parentName:"tr",align:null},"today's date")))),(0,a.kt)("h4",{id:"headers"},"Headers"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"x-api-key: secret_key_value\nAccept: application/json\n")),(0,a.kt)("h4",{id:"response-code"},"Response Code"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"Status: 200 OK\n")),(0,a.kt)("h3",{id:"response-body"},"Response body"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "status": {\n        "message": "application/json",\n        "code": "200"\n    },\n    "data": {\n        "group_count": 1,\n        "group_type": "type of top-level groups in report",\n        "group_list": ["list of top level groups"],\n        "endpoint_group_count": 1,\n        "endpoint_group_type": "type of endpoint groups in report",\n        "endpoint_group_list": ["list of endpoint groups"],\n        "service_count": 1,\n        "service_list": ["list of available services"]\n    }\n}\n')),(0,a.kt)("h6",{id:"example-request"},"Example Request:"),(0,a.kt)("p",null,"URL:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"latest/Report_B/?date=2015-05-01\n")),(0,a.kt)("p",null,"Headers:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"x-api-key: secret_key_value\nAccept: application/json\n")),(0,a.kt)("h6",{id:"example-response"},"Example Response:"),(0,a.kt)("p",null,"Code:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"Status: 200 OK\n")),(0,a.kt)("p",null,"Response body:"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "status": {\n        "message": "application/json",\n        "code": "200"\n    },\n    "data": {\n        "group_count": 2,\n        "group_type": "PROJECTS",\n        "group_list": ["PROJECT_A", "PROJECT_B"],\n        "endpoint_group_count": 3,\n        "endpoint_group_type": "SERVICEGROUPS",\n        "endpoint_group_list": ["SGROUP_A", "SGROUP_B", "SGROUP_C"],\n        "service_count": 8,\n        "service_list": [\n            "service_type_1",\n            "service_type_2",\n            "service_type_3",\n            "service_type_4",\n            "service_type_5",\n            "service_type_6",\n            "service_type_7",\n            "service_type_8"\n        ]\n    }\n}\n')))}c.isMDXComponent=!0}}]);