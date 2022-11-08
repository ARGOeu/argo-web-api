"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[11],{3905:(e,t,n)=>{n.d(t,{Zo:()=>p,kt:()=>m});var a=n(7294);function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function s(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){r(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},i=Object.keys(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var o=a.createContext({}),d=function(e){var t=a.useContext(o),n=t;return e&&(n="function"==typeof e?e(t):s(s({},t),e)),n},p=function(e){var t=d(e.components);return a.createElement(o.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},c=a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,i=e.originalType,o=e.parentName,p=l(e,["components","mdxType","originalType","parentName"]),c=d(n),m=r,k=c["".concat(o,".").concat(m)]||c[m]||u[m]||i;return n?a.createElement(k,s(s({ref:t},p),{},{components:n})):a.createElement(k,s({ref:t},p))}));function m(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var i=n.length,s=new Array(i);s[0]=c;var l={};for(var o in t)hasOwnProperty.call(t,o)&&(l[o]=t[o]);l.originalType=e,l.mdxType="string"==typeof e?e:r,s[1]=l;for(var d=2;d<i;d++)s[d]=n[d];return a.createElement.apply(null,s)}return a.createElement.apply(null,n)}c.displayName="MDXCreateElement"},5977:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>o,contentTitle:()=>s,default:()=>u,frontMatter:()=>i,metadata:()=>l,toc:()=>d});var a=n(7462),r=(n(7294),n(3905));const i={id:"downtimes",title:"Downtimes",sidebar_position:3},s=void 0,l={unversionedId:"tenants_and_feeds/downtimes",id:"tenants_and_feeds/downtimes",title:"Downtimes",description:"API Calls",source:"@site/docs/tenants_and_feeds/downtimes.md",sourceDirName:"tenants_and_feeds",slug:"/tenants_and_feeds/downtimes",permalink:"/argo-web-api/docs/tenants_and_feeds/downtimes",draft:!1,tags:[],version:"current",sidebarPosition:3,frontMatter:{id:"downtimes",title:"Downtimes",sidebar_position:3},sidebar:"tutorialSidebar",previous:{title:"Feeds",permalink:"/argo-web-api/docs/tenants_and_feeds/feeds"},next:{title:"Available Metrics and Tags",permalink:"/argo-web-api/docs/tenants_and_feeds/metrics"}},o={},d=[{value:"API Calls",id:"api-calls",level:2},{value:"GET: List downtime resources",id:"get-list-downtime-resources",level:2},{value:"Input",id:"input",level:3},{value:"Optional Query Parameters",id:"optional-query-parameters",level:4},{value:"Request headers",id:"request-headers",level:3},{value:"Response",id:"response",level:3},{value:"Response body",id:"response-body",level:4},{value:"Request downtimes and filter by severity and classification example",id:"request-downtimes-and-filter-by-severity-and-classification-example",level:3},{value:"POST: Create a new downtime resource",id:"post-create-a-new-downtime-resource",level:2},{value:"Input",id:"input-1",level:3},{value:"Optional Query Parameters",id:"optional-query-parameters-1",level:4},{value:"Request headers",id:"request-headers-1",level:4},{value:"POST BODY",id:"post-body",level:4},{value:"Response",id:"response-1",level:3},{value:"Response body",id:"response-body-1",level:4},{value:"DELETE: Delete an existing downtime resource",id:"delete-delete-an-existing-downtime-resource",level:2},{value:"Input",id:"input-2",level:3},{value:"Request headers",id:"request-headers-2",level:4},{value:"Response",id:"response-2",level:3},{value:"Response body",id:"response-body-2",level:4}],p={toc:d};function u(e){let{components:t,...n}=e;return(0,r.kt)("wrapper",(0,a.Z)({},p,n,{components:t,mdxType:"MDXLayout"}),(0,r.kt)("h2",{id:"api-calls"},"API Calls"),(0,r.kt)("table",null,(0,r.kt)("thead",{parentName:"table"},(0,r.kt)("tr",{parentName:"thead"},(0,r.kt)("th",{parentName:"tr",align:null},"Name"),(0,r.kt)("th",{parentName:"tr",align:null},"Description"),(0,r.kt)("th",{parentName:"tr",align:null},"Shortcut"))),(0,r.kt)("tbody",{parentName:"table"},(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},"GET: List Downtimes resources Request"),(0,r.kt)("td",{parentName:"tr",align:null},"This method can be used to retrieve a list of current downtime resources per date."),(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("a",{parentName:"td",href:"#1"}," Description"))),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},"POST: Create a new downtime resource"),(0,r.kt)("td",{parentName:"tr",align:null},"This method can be used to create a new downtime resource"),(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("a",{parentName:"td",href:"#2"}," Description"))),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},"DELETE: Delete a downtime resource"),(0,r.kt)("td",{parentName:"tr",align:null},"This method can be used to delete an existing downtime resource"),(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("a",{parentName:"td",href:"#3"}," Description"))))),(0,r.kt)("a",{id:"1"}),(0,r.kt)("h2",{id:"get-list-downtime-resources"},"[GET]",": List downtime resources"),(0,r.kt)("p",null,"This method can be used to retrieve a list of current downtime resources per date"),(0,r.kt)("h3",{id:"input"},"Input"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"GET /downtimes?date=YYYY-MM-DD\n")),(0,r.kt)("h4",{id:"optional-query-parameters"},"Optional Query Parameters"),(0,r.kt)("table",null,(0,r.kt)("thead",{parentName:"table"},(0,r.kt)("tr",{parentName:"thead"},(0,r.kt)("th",{parentName:"tr",align:null},"Type"),(0,r.kt)("th",{parentName:"tr",align:null},"Description"),(0,r.kt)("th",{parentName:"tr",align:null},"Required"))),(0,r.kt)("tbody",{parentName:"table"},(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"date")),(0,r.kt)("td",{parentName:"tr",align:null},"Date to retrieve a historic version of the downtime resource. If no date parameter is provided the most current resource will be returned"),(0,r.kt)("td",{parentName:"tr",align:null},"NO")),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"classification")),(0,r.kt)("td",{parentName:"tr",align:null},"optionally filter downtimes by classification value"),(0,r.kt)("td",{parentName:"tr",align:null},"NO")),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"severity")),(0,r.kt)("td",{parentName:"tr",align:null},"optionally filter downtiumes by severity value"),(0,r.kt)("td",{parentName:"tr",align:null},"NO")))),(0,r.kt)("h3",{id:"request-headers"},"Request headers"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,r.kt)("h3",{id:"response"},"Response"),(0,r.kt)("p",null,"Headers: ",(0,r.kt)("inlineCode",{parentName:"p"},"Status: 200 OK")),(0,r.kt)("h4",{id:"response-body"},"Response body"),(0,r.kt)("p",null,"Json Response"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "status": {\n        "message": "Success",\n        "code": "200"\n    },\n    "data": [\n        {\n            "date": "2019-11-04",\n            "endpoints": [\n                {\n                    "hostname": "host-A",\n                    "service": "service-A",\n                    "start_time": "2019-10-11T04:00:33Z",\n                    "end_time": "2019-10-11T15:33:00Z",\n                    "description": "a simple optional description",\n                    "severity": "optional severity value like critical, warning",\n                    "classification": "optional classification value like outage, scheduled"\n                },\n                {\n                    "hostname": "host-B",\n                    "service": "service-B",\n                    "start_time": "2019-10-11T12:00:33Z",\n                    "end_time": "2019-10-11T12:33:00Z"\n                },\n                {\n                    "hostname": "host-C",\n                    "service": "service-C",\n                    "start_time": "2019-10-11T20:00:33Z",\n                    "end_time": "2019-10-11T22:15:00Z"\n                }\n            ]\n        }\n    ]\n}\n')),(0,r.kt)("h3",{id:"request-downtimes-and-filter-by-severity-and-classification-example"},"Request downtimes and filter by severity and classification example"),(0,r.kt)("p",null,"In the following example we request the downtimes for the date 2022-05-11 that are of outage severity and classified as unscheduled"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"HTTP GET /api/v2/downtimes?date=2022-05-11&severity=outage&classification=outage\n")),(0,r.kt)("p",null,"Response: ",(0,r.kt)("inlineCode",{parentName:"p"},"200 OK"),"\nBody:"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "status": {\n        "message": "Success",\n        "code": "200"\n    },\n    "data": [\n        {\n            "date": "2022-05-11",\n            "endpoints": [\n                {\n                    "hostname": "host-A",\n                    "service": "service-A",\n                    "start_time": "2022-05-11T04:00:33Z",\n                    "end_time": "2022-05-11T15:33:00Z",\n                    "severity": "outage",\n                    "classification": "unscheduled"\n                },\n                {\n                    "hostname": "host-B",\n                    "service": "service-B",\n                    "start_time": "2022-05-11T12:00:33Z",\n                    "end_time": "2022-05-11T12:33:00Z",\n                    "severity": "outage",\n                    "classification": "unscheduled",\n                    "description": "a simple description",\n                }\n            ]\n        }\n    ]\n}\n')),(0,r.kt)("p",null,(0,r.kt)("strong",{parentName:"p"},"note"),": ",(0,r.kt)("inlineCode",{parentName:"p"},"description"),", ",(0,r.kt)("inlineCode",{parentName:"p"},"severity")," and ",(0,r.kt)("inlineCode",{parentName:"p"},"classification")," but quite useful to organise the kind of downtimes declared per day."),(0,r.kt)("a",{id:"2"}),(0,r.kt)("h2",{id:"post-create-a-new-downtime-resource"},"[POST]",": Create a new downtime resource"),(0,r.kt)("p",null,"This method can be used to insert a new downtime resource"),(0,r.kt)("h3",{id:"input-1"},"Input"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"POST /downtimes?date=YYYY-MM-DD\n")),(0,r.kt)("h4",{id:"optional-query-parameters-1"},"Optional Query Parameters"),(0,r.kt)("table",null,(0,r.kt)("thead",{parentName:"table"},(0,r.kt)("tr",{parentName:"thead"},(0,r.kt)("th",{parentName:"tr",align:null},"Type"),(0,r.kt)("th",{parentName:"tr",align:null},"Description"),(0,r.kt)("th",{parentName:"tr",align:null},"Required"))),(0,r.kt)("tbody",{parentName:"table"},(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"date")),(0,r.kt)("td",{parentName:"tr",align:null},"Date to create a new historic version of the downtime resource. If no date parameter is provided current date will be supplied automatically"),(0,r.kt)("td",{parentName:"tr",align:null},"NO")))),(0,r.kt)("h4",{id:"request-headers-1"},"Request headers"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,r.kt)("h4",{id:"post-body"},"POST BODY"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "endpoints": [\n        {\n            "hostname": "host-foo",\n            "service": "service-new-foo",\n            "start_time": "2019-10-11T23:10:00Z",\n            "end_time": "2019-10-11T23:25:00Z",\n            "classification": "unscheduled",\n            "severity": "outage",\n        },\n        {\n            "hostname": "host-bar",\n            "service": "service-new-bar",\n            "start_time": "2019-10-11T23:40:00Z",\n            "end_time": "2019-10-11T23:55:00Z",\n            "classification": "unscheduled",\n            "severity": "outage",\n        }\n    ]\n}\n')),(0,r.kt)("h3",{id:"response-1"},"Response"),(0,r.kt)("p",null,"Headers: ",(0,r.kt)("inlineCode",{parentName:"p"},"Status: 201 Created")),(0,r.kt)("h4",{id:"response-body-1"},"Response body"),(0,r.kt)("p",null,"Json Response"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n "status": {\n  "message": "Downtimes set created for date: 2019-11-29",\n  "code": "201"\n }\n}\n')),(0,r.kt)("a",{id:"3"}),(0,r.kt)("h2",{id:"delete-delete-an-existing-downtime-resource"},"[DELETE]",": Delete an existing downtime resource"),(0,r.kt)("p",null,"This method can be used to delete an existing downtime resource"),(0,r.kt)("h3",{id:"input-2"},"Input"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"DELETE /downtimes?date=YYYY-MM-DD\n")),(0,r.kt)("h4",{id:"request-headers-2"},"Request headers"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,r.kt)("h3",{id:"response-2"},"Response"),(0,r.kt)("p",null,"Headers: ",(0,r.kt)("inlineCode",{parentName:"p"},"Status: 200 OK")),(0,r.kt)("h4",{id:"response-body-2"},"Response body"),(0,r.kt)("p",null,"Json Response"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n "status": {\n  "message": "Downtimes set deleted for date: 2019-10-11",\n  "code": "200"\n }\n}\n')))}u.isMDXComponent=!0}}]);