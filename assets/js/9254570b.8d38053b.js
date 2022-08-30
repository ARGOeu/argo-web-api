"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[810],{3905:(e,t,n)=>{n.d(t,{Zo:()=>d,kt:()=>u});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function l(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},i=Object.keys(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var o=r.createContext({}),p=function(e){var t=r.useContext(o),n=t;return e&&(n="function"==typeof e?e(t):l(l({},t),e)),n},d=function(e){var t=p(e.components);return r.createElement(o.Provider,{value:t},e.children)},c={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},m=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,i=e.originalType,o=e.parentName,d=s(e,["components","mdxType","originalType","parentName"]),m=p(n),u=a,k=m["".concat(o,".").concat(u)]||m[u]||c[u]||i;return n?r.createElement(k,l(l({ref:t},d),{},{components:n})):r.createElement(k,l({ref:t},d))}));function u(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var i=n.length,l=new Array(i);l[0]=m;var s={};for(var o in t)hasOwnProperty.call(t,o)&&(s[o]=t[o]);s.originalType=e,s.mdxType="string"==typeof e?e:a,l[1]=s;for(var p=2;p<i;p++)l[p]=n[p];return r.createElement.apply(null,l)}return r.createElement.apply(null,n)}m.displayName="MDXCreateElement"},7844:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>o,contentTitle:()=>l,default:()=>c,frontMatter:()=>i,metadata:()=>s,toc:()=>p});var r=n(7462),a=(n(7294),n(3905));const i={id:"metric_profiles",title:"Metric Profiles",sidebar_position:2},l=void 0,s={unversionedId:"profiles_and_reports/metric_profiles",id:"profiles_and_reports/metric_profiles",title:"Metric Profiles",description:"API Calls",source:"@site/docs/profiles_and_reports/metric_profiles.md",sourceDirName:"profiles_and_reports",slug:"/profiles_and_reports/metric_profiles",permalink:"/argo-web-api/docs/profiles_and_reports/metric_profiles",draft:!1,tags:[],version:"current",sidebarPosition:2,frontMatter:{id:"metric_profiles",title:"Metric Profiles",sidebar_position:2},sidebar:"tutorialSidebar",previous:{title:"Operation Profiles",permalink:"/argo-web-api/docs/profiles_and_reports/operations_profiles"},next:{title:"Aggregation Profiles",permalink:"/argo-web-api/docs/profiles_and_reports/aggregation_profiles"}},o={},p=[{value:"API Calls",id:"api-calls",level:2},{value:"GET: List Metric Profiles",id:"get-list-metric-profiles",level:2},{value:"Input",id:"input",level:3},{value:"Optional Query Parameters",id:"optional-query-parameters",level:4},{value:"Request headers",id:"request-headers",level:3},{value:"Response",id:"response",level:3},{value:"Response body",id:"response-body",level:4},{value:"GET: List A Specific Metric profile",id:"get-list-a-specific-metric-profile",level:2},{value:"Input",id:"input-1",level:3},{value:"Optional Query Parameters",id:"optional-query-parameters-1",level:4},{value:"Request headers",id:"request-headers-1",level:4},{value:"Response",id:"response-1",level:3},{value:"Response body",id:"response-body-1",level:4},{value:"POST: Create a new Metric Profile",id:"post-create-a-new-metric-profile",level:2},{value:"Input",id:"input-2",level:3},{value:"Request headers",id:"request-headers-2",level:4},{value:"Optional Query Parameters",id:"optional-query-parameters-2",level:4},{value:"POST BODY",id:"post-body",level:4},{value:"Response",id:"response-2",level:3},{value:"Response body",id:"response-body-2",level:4},{value:"PUT: Update information on an existing metric profile",id:"put-update-information-on-an-existing-metric-profile",level:2},{value:"Input",id:"input-3",level:3},{value:"Request headers",id:"request-headers-3",level:4},{value:"Optional Query Parameters",id:"optional-query-parameters-3",level:4},{value:"PUT BODY",id:"put-body",level:4},{value:"Response",id:"response-3",level:3},{value:"Response body",id:"response-body-3",level:4},{value:"DELETE: Delete an existing metric profile",id:"delete-delete-an-existing-metric-profile",level:2},{value:"Input",id:"input-4",level:3},{value:"Request headers",id:"request-headers-4",level:4},{value:"Response",id:"response-4",level:3},{value:"Response body",id:"response-body-4",level:4}],d={toc:p};function c(e){let{components:t,...n}=e;return(0,a.kt)("wrapper",(0,r.Z)({},d,n,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h2",{id:"api-calls"},"API Calls"),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:null},"Name"),(0,a.kt)("th",{parentName:"tr",align:null},"Description"),(0,a.kt)("th",{parentName:"tr",align:null},"Shortcut"))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},"GET: List Metric Profile Requests"),(0,a.kt)("td",{parentName:"tr",align:null},"This method can be used to retrieve a list of current metric profiles."),(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"#1"}," Description"))),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},"GET: List a specific Metric profile"),(0,a.kt)("td",{parentName:"tr",align:null},"This method can be used to retrieve a specific metric profile based on its id."),(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"#2"}," Description"))),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},"POST: Create a new metric profile"),(0,a.kt)("td",{parentName:"tr",align:null},"This method can be used to create a new metric profile"),(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"#3"}," Description"))),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},"PUT: Update a metric profile"),(0,a.kt)("td",{parentName:"tr",align:null},"This method can be used to update information on an existing metric profile"),(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"#4"}," Description"))),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},"DELETE: Delete a metric profile"),(0,a.kt)("td",{parentName:"tr",align:null},"This method can be used to delete an existing metric profile"),(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"#5"}," Description"))))),(0,a.kt)("a",{id:"1"}),(0,a.kt)("h2",{id:"get-list-metric-profiles"},"[GET]",": List Metric Profiles"),(0,a.kt)("p",null,"This method can be used to retrieve a list of current Metric profiles"),(0,a.kt)("h3",{id:"input"},"Input"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"GET /metric_profiles\n")),(0,a.kt)("h4",{id:"optional-query-parameters"},"Optional Query Parameters"),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:null},"Type"),(0,a.kt)("th",{parentName:"tr",align:null},"Description"),(0,a.kt)("th",{parentName:"tr",align:null},"Required"))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("inlineCode",{parentName:"td"},"name")),(0,a.kt)("td",{parentName:"tr",align:null},"metric profile name to be used as query"),(0,a.kt)("td",{parentName:"tr",align:null},"NO")),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("inlineCode",{parentName:"td"},"date")),(0,a.kt)("td",{parentName:"tr",align:null},"Date to retrieve a historic version of the metric profile. If no date parameter is provided the most current profile will be returned"),(0,a.kt)("td",{parentName:"tr",align:null},"NO")))),(0,a.kt)("h3",{id:"request-headers"},"Request headers"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,a.kt)("h3",{id:"response"},"Response"),(0,a.kt)("p",null,"Headers: ",(0,a.kt)("inlineCode",{parentName:"p"},"Status: 200 OK")),(0,a.kt)("h4",{id:"response-body"},"Response body"),(0,a.kt)("p",null,"Json Response"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-json"},'{\n "status": {\n  "message": "Success",\n  "code": "200"\n },\n "data": [\n  {\n   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",\n   "date": "2019-10-10",\n   "name": "ch.cern.SAM.ROC",\n   "description": "default profile",\n   "services": [\n    {\n     "service": "CREAM-CE",\n     "metrics": [\n      "emi.cream.CREAMCE-JobSubmit",\n      "emi.wn.WN-Bi",\n      "emi.wn.WN-Csh",\n      "hr.srce.CADist-Check",\n      "hr.srce.CREAMCE-CertLifetime",\n      "emi.wn.WN-SoftVer"\n     ]\n    },\n    {\n     "service": "SRMv2",\n     "metrics": [\n      "hr.srce.SRM2-CertLifetime",\n      "org.sam.SRM-Del",\n      "org.sam.SRM-Get",\n      "org.sam.SRM-GetSURLs",\n      "org.sam.SRM-GetTURLs",\n      "org.sam.SRM-Ls",\n      "org.sam.SRM-LsDir",\n      "org.sam.SRM-Put"\n     ]\n    }\n   ]\n  },\n  {\n   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",\n   "date" : "2019-11-01",\n   "name": "ch.cern.SAM.ROC_CRITICAL",\n   "description": "",\n   "services": [\n    {\n     "service": "CREAM-CE",\n     "metrics": [\n      "emi.cream.CREAMCE-JobSubmit",\n      "emi.wn.WN-Bi",\n      "emi.wn.WN-Csh",\n      "emi.wn.WN-SoftVer"\n     ]\n    },\n    {\n     "service": "SRMv2",\n     "metrics": [\n      "hr.srce.SRM2-CertLifetime",\n      "org.sam.SRM-Del",\n      "org.sam.SRM-Get",\n      "org.sam.SRM-GetSURLs",\n      "org.sam.SRM-GetTURLs",\n      "org.sam.SRM-Ls",\n      "org.sam.SRM-LsDir",\n      "org.sam.SRM-Put"\n     ]\n    }\n   ]\n  }\n ]\n}\n')),(0,a.kt)("a",{id:"2"}),(0,a.kt)("h2",{id:"get-list-a-specific-metric-profile"},"[GET]",": List A Specific Metric profile"),(0,a.kt)("p",null,"This method can be used to retrieve specific metric profile based on its id"),(0,a.kt)("h3",{id:"input-1"},"Input"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"GET /metric_profiles/{ID}\n")),(0,a.kt)("h4",{id:"optional-query-parameters-1"},"Optional Query Parameters"),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:null},"Type"),(0,a.kt)("th",{parentName:"tr",align:null},"Description"),(0,a.kt)("th",{parentName:"tr",align:null},"Required"))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("inlineCode",{parentName:"td"},"date")),(0,a.kt)("td",{parentName:"tr",align:null},"Date to retrieve a historic version of the metric profile. If no date parameter is provided the most current profile will be returned"),(0,a.kt)("td",{parentName:"tr",align:null},"NO")))),(0,a.kt)("h4",{id:"request-headers-1"},"Request headers"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,a.kt)("h3",{id:"response-1"},"Response"),(0,a.kt)("p",null,"Headers: ",(0,a.kt)("inlineCode",{parentName:"p"},"Status: 200 OK")),(0,a.kt)("h4",{id:"response-body-1"},"Response body"),(0,a.kt)("p",null,"Json Response"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-json"},'{\n "status": {\n  "message": "Success",\n  "code": "200"\n },\n "data": [\n  {\n   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",\n   "date" : "2019-11-01",\n   "name": "ch.cern.SAM.ROC_CRITICAL",\n   "description": "a critical profile",\n   "services": [\n    {\n     "service": "CREAM-CE",\n     "metrics": [\n      "emi.cream.CREAMCE-JobSubmit",\n      "emi.wn.WN-Bi",\n      "emi.wn.WN-Csh",\n      "emi.wn.WN-SoftVer"\n     ]\n    },\n    {\n     "service": "SRMv2",\n     "metrics": [\n      "hr.srce.SRM2-CertLifetime",\n      "org.sam.SRM-Del",\n      "org.sam.SRM-Get",\n      "org.sam.SRM-GetSURLs",\n      "org.sam.SRM-GetTURLs",\n      "org.sam.SRM-Ls",\n      "org.sam.SRM-LsDir",\n      "org.sam.SRM-Put"\n     ]\n    }\n   ]\n  }\n ]\n}\n')),(0,a.kt)("a",{id:"3"}),(0,a.kt)("h2",{id:"post-create-a-new-metric-profile"},"[POST]",": Create a new Metric Profile"),(0,a.kt)("p",null,"This method can be used to insert a new metric profile"),(0,a.kt)("h3",{id:"input-2"},"Input"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"POST /metric_profiles\n")),(0,a.kt)("h4",{id:"request-headers-2"},"Request headers"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,a.kt)("h4",{id:"optional-query-parameters-2"},"Optional Query Parameters"),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:null},"Type"),(0,a.kt)("th",{parentName:"tr",align:null},"Description"),(0,a.kt)("th",{parentName:"tr",align:null},"Required"))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("inlineCode",{parentName:"td"},"date")),(0,a.kt)("td",{parentName:"tr",align:null},"Date to create a  new historic version of the metric profile. If no date parameter is provided current date will be supplied automatically"),(0,a.kt)("td",{parentName:"tr",align:null},"NO")))),(0,a.kt)("h4",{id:"post-body"},"POST BODY"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-json"},'{\n  "name": "test_profile",\n  "description": "a profile just for testing",\n  "services": [\n    {\n      "service": "Service-A",\n      "metrics": [\n        "metric.A.1",\n        "metric.A.2",\n        "metric.A.3",\n        "metric.A.4"\n      ]\n    },\n    {\n      "service": "Service-B",\n      "metrics": [\n        "metric.B.1",\n        "metric.B.2"\n      ]\n    }\n  ]\n}\n')),(0,a.kt)("h3",{id:"response-2"},"Response"),(0,a.kt)("p",null,"Headers: ",(0,a.kt)("inlineCode",{parentName:"p"},"Status: 201 Created")),(0,a.kt)("h4",{id:"response-body-2"},"Response body"),(0,a.kt)("p",null,"Json Response"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-json"},'{\n "status": {\n  "message": "Metric Profile successfully created",\n  "code": "201"\n },\n "data": {\n  "id": "{{ID}}",\n  "links": {\n   "self": "https:///api/v2/metric_profiles/{{ID}}"\n  }\n }\n}\n')),(0,a.kt)("a",{id:"4"}),(0,a.kt)("h2",{id:"put-update-information-on-an-existing-metric-profile"},"[PUT]",": Update information on an existing metric profile"),(0,a.kt)("p",null,"This method can be used to update information on an existing metric profile"),(0,a.kt)("h3",{id:"input-3"},"Input"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"PUT /metric_profiles/{ID}\n")),(0,a.kt)("h4",{id:"request-headers-3"},"Request headers"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,a.kt)("h4",{id:"optional-query-parameters-3"},"Optional Query Parameters"),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:null},"Type"),(0,a.kt)("th",{parentName:"tr",align:null},"Description"),(0,a.kt)("th",{parentName:"tr",align:null},"Required"))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("inlineCode",{parentName:"td"},"date")),(0,a.kt)("td",{parentName:"tr",align:null},"Date to update a  new historic version of the operation profile. If no date parameter is provided current date will be supplied automatically"),(0,a.kt)("td",{parentName:"tr",align:null},"NO")))),(0,a.kt)("h4",{id:"put-body"},"PUT BODY"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-json"},'{\n  "name": "test_profile",\n  "description": "this profile is just for tests",\n  "services": [\n    {\n      "service": "Service-A",\n      "metrics": [\n        "metric.A.1",\n        "metric.A.2",\n        "metric.A.3",\n        "metric.A.4"\n      ]\n    },\n    {\n      "service": "Service-B",\n      "metrics": [\n        "metric.B.1",\n        "metric.B.2"\n      ]\n    }\n  ]\n}\n')),(0,a.kt)("h3",{id:"response-3"},"Response"),(0,a.kt)("p",null,"Headers: ",(0,a.kt)("inlineCode",{parentName:"p"},"Status: 200 OK")),(0,a.kt)("h4",{id:"response-body-3"},"Response body"),(0,a.kt)("p",null,"Json Response"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-json"},'{\n "status": {\n  "message": "Metric Profile successfully updated",\n  "code": "200"\n },\n "data": {\n  "id": "{{ID}}",\n  "links": {\n   "self": "https:///api/v2/metric_profiles/{{ID}}"\n  }\n }\n}\n')),(0,a.kt)("a",{id:"5"}),(0,a.kt)("h2",{id:"delete-delete-an-existing-metric-profile"},"[DELETE]",": Delete an existing metric profile"),(0,a.kt)("p",null,"This method can be used to delete an existing metric profile"),(0,a.kt)("h3",{id:"input-4"},"Input"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"DELETE /metric_profiles/{ID}\n")),(0,a.kt)("h4",{id:"request-headers-4"},"Request headers"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,a.kt)("h3",{id:"response-4"},"Response"),(0,a.kt)("p",null,"Headers: ",(0,a.kt)("inlineCode",{parentName:"p"},"Status: 200 OK")),(0,a.kt)("h4",{id:"response-body-4"},"Response body"),(0,a.kt)("p",null,"Json Response"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-json"},'{\n "status": {\n  "message": "Metric Profile Successfully Deleted",\n  "code": "200"\n }\n}\n')))}c.isMDXComponent=!0}}]);