"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[583],{3905:(e,t,n)=>{n.d(t,{Zo:()=>d,kt:()=>m});var a=n(7294);function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function l(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?l(Object(n),!0).forEach((function(t){r(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):l(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},l=Object.keys(e);for(a=0;a<l.length;a++)n=l[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var l=Object.getOwnPropertySymbols(e);for(a=0;a<l.length;a++)n=l[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var o=a.createContext({}),p=function(e){var t=a.useContext(o),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},d=function(e){var t=p(e.components);return a.createElement(o.Provider,{value:t},e.children)},h={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},u=a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,l=e.originalType,o=e.parentName,d=s(e,["components","mdxType","originalType","parentName"]),u=p(n),m=r,c=u["".concat(o,".").concat(m)]||u[m]||h[m]||l;return n?a.createElement(c,i(i({ref:t},d),{},{components:n})):a.createElement(c,i({ref:t},d))}));function m(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var l=n.length,i=new Array(l);i[0]=u;var s={};for(var o in t)hasOwnProperty.call(t,o)&&(s[o]=t[o]);s.originalType=e,s.mdxType="string"==typeof e?e:r,i[1]=s;for(var p=2;p<l;p++)i[p]=n[p];return a.createElement.apply(null,i)}return a.createElement.apply(null,n)}u.displayName="MDXCreateElement"},639:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>o,contentTitle:()=>i,default:()=>h,frontMatter:()=>l,metadata:()=>s,toc:()=>p});var a=n(7462),r=(n(7294),n(3905));const l={id:"threshold_profiles",title:"Threshold Profiles",sidebar_position:4},i=void 0,s={unversionedId:"profiles_and_reports/threshold_profiles",id:"profiles_and_reports/threshold_profiles",title:"Threshold Profiles",description:"Description",source:"@site/docs/profiles_and_reports/threshold_profiles.md",sourceDirName:"profiles_and_reports",slug:"/profiles_and_reports/threshold_profiles",permalink:"/argo-web-api/docs/profiles_and_reports/threshold_profiles",draft:!1,tags:[],version:"current",sidebarPosition:4,frontMatter:{id:"threshold_profiles",title:"Threshold Profiles",sidebar_position:4},sidebar:"tutorialSidebar",previous:{title:"Aggregation Profiles",permalink:"/argo-web-api/docs/profiles_and_reports/aggregation_profiles"},next:{title:"Reports",permalink:"/argo-web-api/docs/profiles_and_reports/reports"}},o={},p=[{value:"Description",id:"description",level:2},{value:"Threshold format",id:"threshold-format",level:2},{value:"Thresholds rule",id:"thresholds-rule",level:2},{value:"Thresholds profile",id:"thresholds-profile",level:2},{value:"API Calls",id:"api-calls",level:2},{value:"GET: List Threshold Profiles",id:"get-list-threshold-profiles",level:2},{value:"Input",id:"input",level:3},{value:"Optional Query Parameters",id:"optional-query-parameters",level:4},{value:"Request headers",id:"request-headers",level:4},{value:"Response",id:"response",level:3},{value:"Response body",id:"response-body",level:4},{value:"GET: List A Specific Thresholds profile",id:"get-list-a-specific-thresholds-profile",level:2},{value:"Input",id:"input-1",level:3},{value:"Request headers",id:"request-headers-1",level:4},{value:"Optional Query Parameters",id:"optional-query-parameters-1",level:4},{value:"Response",id:"response-1",level:3},{value:"Response body",id:"response-body-1",level:4},{value:"POST: Create a new Thresholds Profile",id:"post-create-a-new-thresholds-profile",level:2},{value:"Input",id:"input-2",level:3},{value:"Request headers",id:"request-headers-2",level:4},{value:"POST BODY",id:"post-body",level:4},{value:"Response",id:"response-2",level:3},{value:"Response body",id:"response-body-2",level:4},{value:"PUT: Update information on an existing operations profile",id:"put-update-information-on-an-existing-operations-profile",level:2},{value:"Input",id:"input-3",level:3},{value:"Request headers",id:"request-headers-3",level:4},{value:"PUT BODY",id:"put-body",level:4},{value:"Response",id:"response-3",level:3},{value:"Response body",id:"response-body-3",level:4},{value:"DELETE: Delete an existing aggregation profile",id:"delete-delete-an-existing-aggregation-profile",level:2},{value:"Input",id:"input-4",level:3},{value:"Request headers",id:"request-headers-4",level:4},{value:"Response",id:"response-4",level:3},{value:"Response body",id:"response-body-4",level:4},{value:"Validation Checks",id:"validation-checks",level:2},{value:"Example invalid profile",id:"example-invalid-profile",level:4},{value:"Response",id:"response-5",level:3},{value:"Response body",id:"response-body-5",level:4}],d={toc:p};function h(e){let{components:t,...n}=e;return(0,r.kt)("wrapper",(0,a.Z)({},d,n,{components:t,mdxType:"MDXLayout"}),(0,r.kt)("h2",{id:"description"},"Description"),(0,r.kt)("p",null,"A Threshold profile contains a list of threshold rules. Threshold rules refer to low level monitoring numeric values\nthat accompany metric data and threshold limits on those values that can deem the status 'WARNING' or 'CRITICAL'"),(0,r.kt)("h2",{id:"threshold-format"},"Threshold format"),(0,r.kt)("p",null,"Each threshold rule is expressed as string in the following format\n",(0,r.kt)("inlineCode",{parentName:"p"},"{label}={value}[uom];{warning};{critical};{min};{max}")),(0,r.kt)("ul",null,(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("inlineCode",{parentName:"li"},"label")," : can contain alphanumeric characters but must always begin with a letter"),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("inlineCode",{parentName:"li"},"value")," : is a float or integer"),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("inlineCode",{parentName:"li"},"uom")," : is a string unit of measurement (accepted values: ",(0,r.kt)("inlineCode",{parentName:"li"},"s"),",",(0,r.kt)("inlineCode",{parentName:"li"},"us"),",",(0,r.kt)("inlineCode",{parentName:"li"},"ms"),",",(0,r.kt)("inlineCode",{parentName:"li"},"B"),",",(0,r.kt)("inlineCode",{parentName:"li"},"KB"),",",(0,r.kt)("inlineCode",{parentName:"li"},"MB"),",",(0,r.kt)("inlineCode",{parentName:"li"},"TB"),",",(0,r.kt)("inlineCode",{parentName:"li"},"%"),",",(0,r.kt)("inlineCode",{parentName:"li"},"c"),")"),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("inlineCode",{parentName:"li"},"warning"),": is a range defining the warning threshold limits"),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("inlineCode",{parentName:"li"},"critical"),": is a range defining the critical threshold limits"),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("inlineCode",{parentName:"li"},"min"),": is a float or integer defining the minimum value"),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("inlineCode",{parentName:"li"},"max"),": is a float or integer defining the maximum value")),(0,r.kt)("p",null,"Note: min,max can be omitted."),(0,r.kt)("p",null,"Ranges (",(0,r.kt)("inlineCode",{parentName:"p"},"{warning}")," or ",(0,r.kt)("inlineCode",{parentName:"p"},"{critical}"),") are defined in the following format:\n",(0,r.kt)("inlineCode",{parentName:"p"},"@{floor}:{ceiling}")," -",(0,r.kt)("inlineCode",{parentName:"p"},"@"),": optional - negates the range (value should belong outside ranges limits)"),(0,r.kt)("ul",null,(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("inlineCode",{parentName:"li"},"floor"),": integer/float or ",(0,r.kt)("inlineCode",{parentName:"li"},"~")," that defines negative infinity"),(0,r.kt)("li",{parentName:"ul"},(0,r.kt)("inlineCode",{parentName:"li"},"ceiling"),": integer/float or empty (defining positive infinity)")),(0,r.kt)("h2",{id:"thresholds-rule"},"Thresholds rule"),(0,r.kt)("p",null,"Each thresholds rule can contain multiple threshold definitions separated by space\nfor e.g.\n",(0,r.kt)("inlineCode",{parentName:"p"},"label01=1s;0:10;11:12 label02=1B;0:200;199:500;0;500")),(0,r.kt)("h2",{id:"thresholds-profile"},"Thresholds profile"),(0,r.kt)("p",null,"Each thresholds profile has a name and contains a list of thresholds rules in the following json format\nEach rule must always refer to a metric. It can optionally refer to a host and an endpoint group."),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},'{\n  "id": "68dbd3d8-c95d-41ea-b13e-7ea3644285e5",\n  "name": "example-threshold-profile-101"\n  "rules":[\n    {\n      "metric": "httpd.ResponseTime"\n      "thresholds": "response=20ms;0:300;299:1000"\n    },\n    {\n      "host": "webserver01.example.foo"\n      "metric": "httpd.ResponseTime"\n      "thresholds": "response=20ms;0:200;199:300"\n    },\n    {\n      "endpoint_group": "TEST-SITE-51"\n      "metric": "httpd.ResponseTime"\n      "thresholds": "response=20ms;0:500;499:1000"\n    }\n  ]\n}\n')),(0,r.kt)("h2",{id:"api-calls"},"API Calls"),(0,r.kt)("table",null,(0,r.kt)("thead",{parentName:"table"},(0,r.kt)("tr",{parentName:"thead"},(0,r.kt)("th",{parentName:"tr",align:null},"Name"),(0,r.kt)("th",{parentName:"tr",align:null},"Description"),(0,r.kt)("th",{parentName:"tr",align:null},"Shortcut"))),(0,r.kt)("tbody",{parentName:"table"},(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},"GET: List Thresholds Profile Requests"),(0,r.kt)("td",{parentName:"tr",align:null},"This method can be used to retrieve a list of current Thresholds profiles."),(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("a",{parentName:"td",href:"#1"}," Description"))),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},"GET: List a specific Threshold profile"),(0,r.kt)("td",{parentName:"tr",align:null},"This method can be used to retrieve a specific Threshold profile based on its id."),(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("a",{parentName:"td",href:"#2"}," Description"))),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},"POST: Create a new Threshold profile"),(0,r.kt)("td",{parentName:"tr",align:null},"This method can be used to create a new Threshold profile"),(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("a",{parentName:"td",href:"#3"}," Description"))),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},"PUT: Update an Threshold profile"),(0,r.kt)("td",{parentName:"tr",align:null},"This method can be used to update information on an existing Threshold profile"),(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("a",{parentName:"td",href:"#4"}," Description"))),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},"DELETE: Delete an Threshold profile"),(0,r.kt)("td",{parentName:"tr",align:null},"This method can be used to delete an existing Threshold profile"),(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("a",{parentName:"td",href:"#5"}," Description"))))),(0,r.kt)("a",{id:"1"}),(0,r.kt)("h2",{id:"get-list-threshold-profiles"},"[GET]",": List Threshold Profiles"),(0,r.kt)("p",null,"This method can be used to retrieve a list of current Threshold profiles"),(0,r.kt)("h3",{id:"input"},"Input"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"GET /thresholds_profiles\n")),(0,r.kt)("h4",{id:"optional-query-parameters"},"Optional Query Parameters"),(0,r.kt)("table",null,(0,r.kt)("thead",{parentName:"table"},(0,r.kt)("tr",{parentName:"thead"},(0,r.kt)("th",{parentName:"tr",align:null},"Type"),(0,r.kt)("th",{parentName:"tr",align:null},"Description"),(0,r.kt)("th",{parentName:"tr",align:null},"Required"))),(0,r.kt)("tbody",{parentName:"table"},(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"name")),(0,r.kt)("td",{parentName:"tr",align:null},"thresholds profile name to be used as query"),(0,r.kt)("td",{parentName:"tr",align:null},"NO")),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"date")),(0,r.kt)("td",{parentName:"tr",align:null},"Date to retrieve a historic version of the thresholds profiles. If no date parameter is provided the most current profile will be returned"),(0,r.kt)("td",{parentName:"tr",align:null},"NO")))),(0,r.kt)("h4",{id:"request-headers"},"Request headers"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,r.kt)("h3",{id:"response"},"Response"),(0,r.kt)("p",null,"Headers: ",(0,r.kt)("inlineCode",{parentName:"p"},"Status: 200 OK")),(0,r.kt)("h4",{id:"response-body"},"Response body"),(0,r.kt)("p",null,"Json Response"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "status": {\n        "message": "Success",\n        "code": "200"\n    },\n    "data": [\n        {\n            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",\n            "name": "thr01",\n            "rules": [\n                {\n                    "host": "hostFoo",\n                    "metric": "metricA",\n                    "thresholds": "freshnesss=1s;10;9:;0;25 entries=1;3;2:0;10"\n                }\n            ]\n        },\n        {\n            "id": "6ac7d222-1f8e-4a02-a502-720e8f11e50b",\n            "name": "thr02",\n            "rules": [\n                {\n                    "host": "hostFoo",\n                    "metric": "metricA",\n                    "thresholds": "freshness=1s;10;9:;0;25 entries=1;3;2:0;10"\n                }\n            ]\n        },\n        {\n            "id": "6ac7d555-1f8e-4a02-a502-720e8f11e50b",\n            "name": "thr03",\n            "rules": [\n                {\n                    "host": "hostFoo",\n                    "metric": "metricA",\n                    "thresholds": "freshness=1s;10;9:;0;25 entries=1;3;2:0;10"\n                }\n            ]\n        }\n    ]\n}\n')),(0,r.kt)("a",{id:"2"}),(0,r.kt)("h2",{id:"get-list-a-specific-thresholds-profile"},"[GET]",": List A Specific Thresholds profile"),(0,r.kt)("p",null,"This method can be used to retrieve specific Thresholds profile based on its id"),(0,r.kt)("h3",{id:"input-1"},"Input"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"GET /thresholds_profiles/{ID}\n")),(0,r.kt)("h4",{id:"request-headers-1"},"Request headers"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,r.kt)("h4",{id:"optional-query-parameters-1"},"Optional Query Parameters"),(0,r.kt)("table",null,(0,r.kt)("thead",{parentName:"table"},(0,r.kt)("tr",{parentName:"thead"},(0,r.kt)("th",{parentName:"tr",align:null},"Type"),(0,r.kt)("th",{parentName:"tr",align:null},"Description"),(0,r.kt)("th",{parentName:"tr",align:null},"Required"))),(0,r.kt)("tbody",{parentName:"table"},(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"date")),(0,r.kt)("td",{parentName:"tr",align:null},"Date to retrieve a historic version of the thresholds profile. If no date parameter is provided the most current profile will be returned"),(0,r.kt)("td",{parentName:"tr",align:null},"NO")))),(0,r.kt)("h3",{id:"response-1"},"Response"),(0,r.kt)("p",null,"Headers: ",(0,r.kt)("inlineCode",{parentName:"p"},"Status: 200 OK")),(0,r.kt)("h4",{id:"response-body-1"},"Response body"),(0,r.kt)("p",null,"Json Response"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "status": {\n        "message": "Success",\n        "code": "200"\n    },\n    "data": [\n        {\n            "id": "6ac7d222-1f8e-4a02-a502-720e8f11e50b",\n            "name": "thr02",\n            "rules": [\n                {\n                    "host": "hostFoo",\n                    "metric": "metricA",\n                    "thresholds": "freshness=1s;10;9:;0;25 entries=1;3;2:0;10"\n                }\n            ]\n        }\n    ]\n}\n')),(0,r.kt)("a",{id:"3"}),(0,r.kt)("h2",{id:"post-create-a-new-thresholds-profile"},"[POST]",": Create a new Thresholds Profile"),(0,r.kt)("p",null,"This method can be used to insert a new thresholds profile"),(0,r.kt)("h3",{id:"input-2"},"Input"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"POST /thresholds_profiles\n")),(0,r.kt)("h4",{id:"request-headers-2"},"Request headers"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,r.kt)("table",null,(0,r.kt)("thead",{parentName:"table"},(0,r.kt)("tr",{parentName:"thead"},(0,r.kt)("th",{parentName:"tr",align:null},"Type"),(0,r.kt)("th",{parentName:"tr",align:null},"Description"),(0,r.kt)("th",{parentName:"tr",align:null},"Required"))),(0,r.kt)("tbody",{parentName:"table"},(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"date")),(0,r.kt)("td",{parentName:"tr",align:null},"Date to create a historic version of the thresholds profile. If no date parameter is provided the most current profile will be returned"),(0,r.kt)("td",{parentName:"tr",align:null},"NO")))),(0,r.kt)("h4",{id:"post-body"},"POST BODY"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "name": "thr04",\n    "rules": [\n        {\n            "metric": "metricB",\n            "thresholds": "time=1s;10;9:30;0;30"\n        }\n    ]\n}\n')),(0,r.kt)("h3",{id:"response-2"},"Response"),(0,r.kt)("p",null,"Headers: ",(0,r.kt)("inlineCode",{parentName:"p"},"Status: 201 Created")),(0,r.kt)("h4",{id:"response-body-2"},"Response body"),(0,r.kt)("p",null,"Json Response"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "status": {\n        "message": "Thresholds Profile successfully created",\n        "code": "201"\n    },\n    "data": {\n        "id": "{{ID}}",\n        "links": {\n            "self": "https:///api/v2/thresholds_profiles/{{ID}}"\n        }\n    }\n}\n')),(0,r.kt)("a",{id:"4"}),(0,r.kt)("h2",{id:"put-update-information-on-an-existing-operations-profile"},"[PUT]",": Update information on an existing operations profile"),(0,r.kt)("p",null,"This method can be used to update information on an existing operations profile"),(0,r.kt)("h3",{id:"input-3"},"Input"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"PUT /thresholds_profiles/{ID}\n")),(0,r.kt)("h4",{id:"request-headers-3"},"Request headers"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,r.kt)("table",null,(0,r.kt)("thead",{parentName:"table"},(0,r.kt)("tr",{parentName:"thead"},(0,r.kt)("th",{parentName:"tr",align:null},"Type"),(0,r.kt)("th",{parentName:"tr",align:null},"Description"),(0,r.kt)("th",{parentName:"tr",align:null},"Required"))),(0,r.kt)("tbody",{parentName:"table"},(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"date")),(0,r.kt)("td",{parentName:"tr",align:null},"Date to update a historic version of the thresholds profile. If no date parameter is provided the most current profile will be returned"),(0,r.kt)("td",{parentName:"tr",align:null},"NO")))),(0,r.kt)("h4",{id:"put-body"},"PUT BODY"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "name": "thr04",\n    "rules": [\n        {\n            "metric": "metricB",\n            "thresholds": "time=1s;10;9:30;0;30"\n        }\n    ]\n}\n')),(0,r.kt)("h3",{id:"response-3"},"Response"),(0,r.kt)("p",null,"Headers: ",(0,r.kt)("inlineCode",{parentName:"p"},"Status: 200 OK")),(0,r.kt)("h4",{id:"response-body-3"},"Response body"),(0,r.kt)("p",null,"Json Response"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "status": {\n        "message": "Thresholds Profile successfully updated",\n        "code": "200"\n    }\n}\n')),(0,r.kt)("a",{id:"5"}),(0,r.kt)("h2",{id:"delete-delete-an-existing-aggregation-profile"},"[DELETE]",": Delete an existing aggregation profile"),(0,r.kt)("p",null,"This method can be used to delete an existing aggregation profile"),(0,r.kt)("h3",{id:"input-4"},"Input"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"DELETE /thresholds_profiles/{ID}\n")),(0,r.kt)("h4",{id:"request-headers-4"},"Request headers"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),(0,r.kt)("h3",{id:"response-4"},"Response"),(0,r.kt)("p",null,"Headers: ",(0,r.kt)("inlineCode",{parentName:"p"},"Status: 200 OK")),(0,r.kt)("h4",{id:"response-body-4"},"Response body"),(0,r.kt)("p",null,"Json Response"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "status": {\n        "message": "Operations Profile Successfully Deleted",\n        "code": "200"\n    }\n}\n')),(0,r.kt)("h2",{id:"validation-checks"},"Validation Checks"),(0,r.kt)("p",null,"When submitting or updating a new threshold profile, validation checks are performed on json POST/PUT body for the following cases:"),(0,r.kt)("ul",null,(0,r.kt)("li",{parentName:"ul"},"Check if each thresholds rule is valid according to threshold specification discussed in the first paragraph")),(0,r.kt)("p",null,"When an invalid operations profile is submitted the api responds with a validation error list:"),(0,r.kt)("h4",{id:"example-invalid-profile"},"Example invalid profile"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "name": "test-invalid-01",\n    "rules": [\n        { "thresholds": "bad01=33;33s" },\n        { "thresholds": "good01=33s;33 good02=1s;~:10;9:;-20;30" },\n        { "thresholds": "bad02=33sbad03=1s;~~:10;9:;-20;30" },\n        { "thresholds": "33;33 bad04=33s;33 -20;30" },\n        { "thresholds": "good01=2KB;0:3;2:10;0;20 good02=1c;~:10;9:30;-30;30" }\n    ]\n}\n')),(0,r.kt)("p",null,"Api response is the following:"),(0,r.kt)("h3",{id:"response-5"},"Response"),(0,r.kt)("p",null,"Headers: ",(0,r.kt)("inlineCode",{parentName:"p"},"Status: 422 Unprocessable Entity")),(0,r.kt)("h4",{id:"response-body-5"},"Response body"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "status": {\n        "message": "Validation Error",\n        "code": "422"\n    },\n    "errors": [\n        {\n            "message": "Validation Failed",\n            "code": "422",\n            "details": "Invalid threshold: bad01=33;33s"\n        },\n        {\n            "message": "Validation Failed",\n            "code": "422",\n            "details": "Invalid threshold: bad02=33sbad03=1s;~~:10;9:;-20;30"\n        },\n        {\n            "message": "Validation Failed",\n            "code": "422",\n            "details": "Invalid threshold: 33;33 bad04=33s;33 -20;30"\n        }\n    ]\n}\n')))}h.isMDXComponent=!0}}]);