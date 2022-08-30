"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[11],{3905:(e,t,n)=>{n.d(t,{Zo:()=>p,kt:()=>c});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function s(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function o(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?s(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):s(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function i(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},s=Object.keys(e);for(r=0;r<s.length;r++)n=s[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var s=Object.getOwnPropertySymbols(e);for(r=0;r<s.length;r++)n=s[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var l=r.createContext({}),d=function(e){var t=r.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):o(o({},t),e)),n},p=function(e){var t=d(e.components);return r.createElement(l.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},m=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,s=e.originalType,l=e.parentName,p=i(e,["components","mdxType","originalType","parentName"]),m=d(n),c=a,k=m["".concat(l,".").concat(c)]||m[c]||u[c]||s;return n?r.createElement(k,o(o({ref:t},p),{},{components:n})):r.createElement(k,o({ref:t},p))}));function c(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var s=n.length,o=new Array(s);o[0]=m;var i={};for(var l in t)hasOwnProperty.call(t,l)&&(i[l]=t[l]);i.originalType=e,i.mdxType="string"==typeof e?e:a,o[1]=i;for(var d=2;d<s;d++)o[d]=n[d];return r.createElement.apply(null,o)}return r.createElement.apply(null,n)}m.displayName="MDXCreateElement"},5977:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>l,contentTitle:()=>o,default:()=>u,frontMatter:()=>s,metadata:()=>i,toc:()=>d});var r=n(7462),a=(n(7294),n(3905));const s={id:"downtimes",title:"Downtimes",sidebar_position:3},o=void 0,i={unversionedId:"tenants_and_feeds/downtimes",id:"tenants_and_feeds/downtimes",title:"Downtimes",description:"API Calls",source:"@site/docs/tenants_and_feeds/downtimes.md",sourceDirName:"tenants_and_feeds",slug:"/tenants_and_feeds/downtimes",permalink:"/argo-web-api/docs/tenants_and_feeds/downtimes",draft:!1,tags:[],version:"current",sidebarPosition:3,frontMatter:{id:"downtimes",title:"Downtimes",sidebar_position:3},sidebar:"tutorialSidebar",previous:{title:"Feeds",permalink:"/argo-web-api/docs/tenants_and_feeds/feeds"},next:{title:"Available Metrics and Tags",permalink:"/argo-web-api/docs/tenants_and_feeds/metrics"}},l={},d=[{value:"API Calls",id:"api-calls",level:2},{value:"GET: List downtime resources",id:"get-list-downtime-resources",level:2},{value:"Input",id:"input",level:3},{value:"Optional Query Parameters",id:"optional-query-parameters",level:4},{value:"Request headers",id:"request-headers",level:3},{value:"Response",id:"response",level:3},{value:"Response body",id:"response-body",level:4},{value:"POST: Create a new downtime resource",id:"post-create-a-new-downtime-resource",level:2},{value:"Input",id:"input-1",level:3},{value:"Optional Query Parameters",id:"optional-query-parameters-1",level:4},{value:"Request headers",id:"request-headers-1",level:4},{value:"POST BODY",id:"post-body",level:4},{value:"Response",id:"response-1",level:3},{value:"Response body",id:"response-body-1",level:4},{value:"DELETE: Delete an existing downtime resource",id:"delete-delete-an-existing-downtime-resource",level:2},{value:"Input",id:"input-2",level:3},{value:"Request headers",id:"request-headers-2",level:4},{value:"Response",id:"response-2",level:3},{value:"Response body",id:"response-body-2",level:4}],p={toc:d};function u(e){let{components:t,...n}=e;return(0,a.kt)("wrapper",(0,r.Z)({},p,n,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h2",{id:"api-calls"},"API Calls"),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:null},"Name"),(0,a.kt)("th",{parentName:"tr",align:null},"Description"),(0,a.kt)("th",{parentName:"tr",align:null},"Shortcut"))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},"GET: List Downtimes resources Request"),(0,a.kt)("td",{parentName:"tr",align:null},"This method can be used to retrieve a list of current downtime resources per date."),(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"#1"}," Description"))),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},"POST: Create a new downtime resource"),(0,a.kt)("td",{parentName:"tr",align:null},"This method can be used to create a new downtime resource"),(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"#2"}," Description"))),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},"DELETE: Delete a downtime resource"),(0,a.kt)("td",{parentName:"tr",align:null},"This method can be used to delete an existing downtime resource"),(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"#3"}," Description"))))),(0,a.kt)("a",{id:"1"}),(0,a.kt)("h2",{id:"get-list-downtime-resources"},"[GET]",": List downtime resources"),(0,a.kt)("p",null,"This method can be used to retrieve a list of current downtime resources per date"),(0,a.kt)("h3",{id:"input"},"Input"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"GET /downtimes?date=YYYY-MM-DD\n")),(0,a.kt)("h4",{id:"optional-query-parameters"},"Optional Query Parameters"),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:null},"Type"),(0,a.kt)("th",{parentName:"tr",align:null},"Description"),(0,a.kt)("th",{parentName:"tr",align:null},"Required"))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("inlineCode",{parentName:"td"},"date")),(0,a.kt)("td",{parentName:"tr",align:null},"Date to retrieve a historic version of the downtime resource. If no date parameter is provided the most current resource will be returned"),(0,a.kt)("td",{parentName:"tr",align:null},"NO")))),(0,a.kt)("h3",{id:"request-headers"},"Request headers"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,a.kt)("h3",{id:"response"},"Response"),(0,a.kt)("p",null,"Headers: ",(0,a.kt)("inlineCode",{parentName:"p"},"Status: 200 OK")),(0,a.kt)("h4",{id:"response-body"},"Response body"),(0,a.kt)("p",null,"Json Response"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "status": {\n        "message": "Success",\n        "code": "200"\n    },\n    "data": [\n        {\n            "date": "2019-11-04",\n            "endpoints": [\n                {\n                    "hostname": "host-A",\n                    "service": "service-A",\n                    "start_time": "2019-10-11T04:00:33Z",\n                    "end_time": "2019-10-11T15:33:00Z"\n                },\n                {\n                    "hostname": "host-B",\n                    "service": "service-B",\n                    "start_time": "2019-10-11T12:00:33Z",\n                    "end_time": "2019-10-11T12:33:00Z"\n                },\n                {\n                    "hostname": "host-C",\n                    "service": "service-C",\n                    "start_time": "2019-10-11T20:00:33Z",\n                    "end_time": "2019-10-11T22:15:00Z"\n                }\n            ]\n        }\n    ]\n}\n')),(0,a.kt)("a",{id:"2"}),(0,a.kt)("h2",{id:"post-create-a-new-downtime-resource"},"[POST]",": Create a new downtime resource"),(0,a.kt)("p",null,"This method can be used to insert a new downtime resource"),(0,a.kt)("h3",{id:"input-1"},"Input"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"POST /downtimes?date=YYYY-MM-DD\n")),(0,a.kt)("h4",{id:"optional-query-parameters-1"},"Optional Query Parameters"),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:null},"Type"),(0,a.kt)("th",{parentName:"tr",align:null},"Description"),(0,a.kt)("th",{parentName:"tr",align:null},"Required"))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("inlineCode",{parentName:"td"},"date")),(0,a.kt)("td",{parentName:"tr",align:null},"Date to create a new historic version of the downtime resource. If no date parameter is provided current date will be supplied automatically"),(0,a.kt)("td",{parentName:"tr",align:null},"NO")))),(0,a.kt)("h4",{id:"request-headers-1"},"Request headers"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,a.kt)("h4",{id:"post-body"},"POST BODY"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "endpoints": [\n        {\n            "hostname": "host-foo",\n            "service": "service-new-foo",\n            "start_time": "2019-10-11T23:10:00Z",\n            "end_time": "2019-10-11T23:25:00Z"\n        },\n        {\n            "hostname": "host-bar",\n            "service": "service-new-bar",\n            "start_time": "2019-10-11T23:40:00Z",\n            "end_time": "2019-10-11T23:55:00Z"\n        }\n    ]\n}\n')),(0,a.kt)("h3",{id:"response-1"},"Response"),(0,a.kt)("p",null,"Headers: ",(0,a.kt)("inlineCode",{parentName:"p"},"Status: 201 Created")),(0,a.kt)("h4",{id:"response-body-1"},"Response body"),(0,a.kt)("p",null,"Json Response"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-json"},'{\n "status": {\n  "message": "Downtimes set created for date: 2019-11-29",\n  "code": "201"\n }\n}\n')),(0,a.kt)("a",{id:"3"}),(0,a.kt)("h2",{id:"delete-delete-an-existing-downtime-resource"},"[DELETE]",": Delete an existing downtime resource"),(0,a.kt)("p",null,"This method can be used to delete an existing downtime resource"),(0,a.kt)("h3",{id:"input-2"},"Input"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"DELETE /downtimes?date=YYYY-MM-DD\n")),(0,a.kt)("h4",{id:"request-headers-2"},"Request headers"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,a.kt)("h3",{id:"response-2"},"Response"),(0,a.kt)("p",null,"Headers: ",(0,a.kt)("inlineCode",{parentName:"p"},"Status: 200 OK")),(0,a.kt)("h4",{id:"response-body-2"},"Response body"),(0,a.kt)("p",null,"Json Response"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-json"},'{\n "status": {\n  "message": "Downtimes set deleted for date: 2019-10-11",\n  "code": "200"\n }\n}\n')))}u.isMDXComponent=!0}}]);