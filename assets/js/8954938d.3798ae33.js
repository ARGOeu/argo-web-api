"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[278],{3905:(e,t,n)=>{n.d(t,{Zo:()=>u,kt:()=>c});var a=n(7294);function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function l(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){r(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function p(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},o=Object.keys(e);for(a=0;a<o.length;a++)n=o[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(a=0;a<o.length;a++)n=o[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var i=a.createContext({}),s=function(e){var t=a.useContext(i),n=t;return e&&(n="function"==typeof e?e(t):l(l({},t),e)),n},u=function(e){var t=s(e.components);return a.createElement(i.Provider,{value:t},e.children)},d={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},g=a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,o=e.originalType,i=e.parentName,u=p(e,["components","mdxType","originalType","parentName"]),g=s(n),c=r,k=g["".concat(i,".").concat(c)]||g[c]||d[c]||o;return n?a.createElement(k,l(l({ref:t},u),{},{components:n})):a.createElement(k,l({ref:t},u))}));function c(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var o=n.length,l=new Array(o);l[0]=g;var p={};for(var i in t)hasOwnProperty.call(t,i)&&(p[i]=t[i]);p.originalType=e,p.mdxType="string"==typeof e?e:r,l[1]=p;for(var s=2;s<o;s++)l[s]=n[s];return a.createElement.apply(null,l)}return a.createElement.apply(null,n)}g.displayName="MDXCreateElement"},6241:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>i,contentTitle:()=>l,default:()=>d,frontMatter:()=>o,metadata:()=>p,toc:()=>s});var a=n(7462),r=(n(7294),n(3905));const o={id:"topology_groups",title:"Topology Groups",sidebar_position:3},l=void 0,p={unversionedId:"topology/topology_groups",id:"topology/topology_groups",title:"Topology Groups",description:"API calls for handling topology group resources",source:"@site/docs/topology/topology_groups.md",sourceDirName:"topology",slug:"/topology/topology_groups",permalink:"/argo-web-api/docs/topology/topology_groups",draft:!1,tags:[],version:"current",sidebarPosition:3,frontMatter:{id:"topology_groups",title:"Topology Groups",sidebar_position:3},sidebar:"tutorialSidebar",previous:{title:"Topology Endpoints",permalink:"/argo-web-api/docs/topology/topology_endpoints"},next:{title:"Topology Statistics",permalink:"/argo-web-api/docs/topology/topology_stats"}},i={},s=[{value:"API calls for handling topology group resources",id:"api-calls-for-handling-topology-group-resources",level:2},{value:"POST: Create group topology for specific date",id:"post-create-group-topology-for-specific-date",level:2},{value:"Input",id:"input",level:3},{value:"Url Parameters",id:"url-parameters",level:4},{value:"Headers",id:"headers",level:4},{value:"POST BODY",id:"post-body",level:3},{value:"Response Code",id:"response-code",level:4},{value:"Response body",id:"response-body",level:3},{value:"409 Conflict when trying to insert a topology that already exists",id:"409-conflict-when-trying-to-insert-a-topology-that-already-exists",level:2},{value:"Response Code",id:"response-code-1",level:3},{value:"Response body",id:"response-body-1",level:3},{value:"GET: List group topology for specific date",id:"get-list-group-topology-for-specific-date",level:2},{value:"Input",id:"input-1",level:3},{value:"Url Parameters",id:"url-parameters-1",level:4},{value:"Headers",id:"headers-1",level:4},{value:"Example Request",id:"example-request",level:4},{value:"Response Code",id:"response-code-2",level:4},{value:"Response body",id:"response-body-2",level:3},{value:"DELETE: Delete group topology for a specific date",id:"delete-delete-group-topology-for-a-specific-date",level:2},{value:"Input",id:"input-2",level:3},{value:"Request headers",id:"request-headers",level:4},{value:"Response",id:"response",level:3},{value:"Response body",id:"response-body-3",level:4},{value:"GET: List group topology for specific report",id:"get-list-group-topology-for-specific-report",level:2},{value:"Input",id:"input-3",level:3},{value:"Url Parameters",id:"url-parameters-2",level:4},{value:"Headers",id:"headers-2",level:4},{value:"Example Request",id:"example-request-1",level:4},{value:"Response Code",id:"response-code-3",level:4},{value:"Response body",id:"response-body-4",level:3}],u={toc:s};function d(e){let{components:t,...n}=e;return(0,r.kt)("wrapper",(0,a.Z)({},u,n,{components:t,mdxType:"MDXLayout"}),(0,r.kt)("h2",{id:"api-calls-for-handling-topology-group-resources"},"API calls for handling topology group resources"),(0,r.kt)("table",null,(0,r.kt)("thead",{parentName:"table"},(0,r.kt)("tr",{parentName:"thead"},(0,r.kt)("th",{parentName:"tr",align:null},"Name"),(0,r.kt)("th",{parentName:"tr",align:null},"Description"),(0,r.kt)("th",{parentName:"tr",align:null},"Shortcut"))),(0,r.kt)("tbody",{parentName:"table"},(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},"POST: Create group topology for specific date"),(0,r.kt)("td",{parentName:"tr",align:null},"Creates a daily group topology mapping endpoints to endpoint groups"),(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("a",{href:"#1"},"Description"))),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},"GET: List group topology for specific date"),(0,r.kt)("td",{parentName:"tr",align:null},"Lists group topology for a specific date"),(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("a",{href:"#2"},"Description"))),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},"DELETE: Delete group topology for specific date"),(0,r.kt)("td",{parentName:"tr",align:null},"Delete group topology items for specific date"),(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("a",{href:"#3"},"Description"))),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},"GET: List group topology for specific report"),(0,r.kt)("td",{parentName:"tr",align:null},"Lists group topology for a specific report"),(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("a",{href:"#4"},"Description"))))),(0,r.kt)("a",{id:"1"}),(0,r.kt)("h2",{id:"post-create-group-topology-for-specific-date"},"POST: Create group topology for specific date"),(0,r.kt)("p",null,"Creates a daily group topology mapping top-level groups to subgroups"),(0,r.kt)("h3",{id:"input"},"Input"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"POST /topology/groups?date=YYYY-MM-DD\n")),(0,r.kt)("h4",{id:"url-parameters"},"Url Parameters"),(0,r.kt)("table",null,(0,r.kt)("thead",{parentName:"table"},(0,r.kt)("tr",{parentName:"thead"},(0,r.kt)("th",{parentName:"tr",align:null},"Type"),(0,r.kt)("th",{parentName:"tr",align:null},"Description"),(0,r.kt)("th",{parentName:"tr",align:null},"Required"),(0,r.kt)("th",{parentName:"tr",align:null},"Default value"))),(0,r.kt)("tbody",{parentName:"table"},(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"date")),(0,r.kt)("td",{parentName:"tr",align:null},"target a specific date"),(0,r.kt)("td",{parentName:"tr",align:null},"NO"),(0,r.kt)("td",{parentName:"tr",align:null},"today's date")))),(0,r.kt)("h4",{id:"headers"},"Headers"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"x-api-key: secret_key_value\nAccept: application/json\n")),(0,r.kt)("h3",{id:"post-body"},"POST BODY"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'[\n    {\n        "group": "NGIA",\n        "type": "NGIS",\n        "subgroup": "SITEA",\n        "tags": {\n            "scope": "FEDERATION",\n            "infrastructure": "Production",\n            "certification": "Certified"\n        }\n    },\n    {\n        "group": "NGIA",\n        "type": "NGIS",\n        "subgroup": "SITEB",\n        "tags": {\n            "scope": "FEDERATION",\n            "infrastructure": "Production",\n            "certification": "Certified"\n        },\n        "notifications": {\n                "contacts": ["email01@example.com"],\n                "enabled": true\n        }\n    },\n    {\n        "group": "PROJECTZ",\n        "type": "PROJECT",\n        "subgroup": "SITEZ",\n        "tags": {\n            "scope": "FEDERATION",\n            "infrastructure": "Production",\n            "certification": "Certified"\n        }\n    }\n]\n')),(0,r.kt)("h4",{id:"response-code"},"Response Code"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"Status: 201 OK Created\n")),(0,r.kt)("h3",{id:"response-body"},"Response body"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "message": "Topology of 3 groups created for date: YYYY-MM-DD",\n    "code": "201"\n}\n')),(0,r.kt)("h2",{id:"409-conflict-when-trying-to-insert-a-topology-that-already-exists"},"409 Conflict when trying to insert a topology that already exists"),(0,r.kt)("p",null,"When trying to insert a topology for a specific date that already exists the api will answer with the following response:"),(0,r.kt)("h3",{id:"response-code-1"},"Response Code"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"Status: 409 Conflict\n")),(0,r.kt)("h3",{id:"response-body-1"},"Response body"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "message": "topology already exists for date: YYYY-MM-DD, please either update it or delete it first!",\n    "code": "409"\n}\n')),(0,r.kt)("p",null,"User can proceed with either updating the existing topology OR deleting before trying to create it anew"),(0,r.kt)("a",{id:"2"}),(0,r.kt)("h2",{id:"get-list-group-topology-for-specific-date"},"GET: List group topology for specific date"),(0,r.kt)("p",null,"Lists group topology items for specific date"),(0,r.kt)("h3",{id:"input-1"},"Input"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"GET /topology/groups?date=YYYY-MM-DD\n")),(0,r.kt)("h4",{id:"url-parameters-1"},"Url Parameters"),(0,r.kt)("table",null,(0,r.kt)("thead",{parentName:"table"},(0,r.kt)("tr",{parentName:"thead"},(0,r.kt)("th",{parentName:"tr",align:null},"Type"),(0,r.kt)("th",{parentName:"tr",align:null},"Description"),(0,r.kt)("th",{parentName:"tr",align:null},"Required"),(0,r.kt)("th",{parentName:"tr",align:null},"Default value"))),(0,r.kt)("tbody",{parentName:"table"},(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"date")),(0,r.kt)("td",{parentName:"tr",align:null},"target a specific date"),(0,r.kt)("td",{parentName:"tr",align:null},"NO"),(0,r.kt)("td",{parentName:"tr",align:null},"today's date")),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"group")),(0,r.kt)("td",{parentName:"tr",align:null},"filter by group name"),(0,r.kt)("td",{parentName:"tr",align:null},"NO"),(0,r.kt)("td",{parentName:"tr",align:null})),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"type")),(0,r.kt)("td",{parentName:"tr",align:null},"filter by group type"),(0,r.kt)("td",{parentName:"tr",align:null},"NO"),(0,r.kt)("td",{parentName:"tr",align:null})),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"subgroup")),(0,r.kt)("td",{parentName:"tr",align:null},"filter by subgroup"),(0,r.kt)("td",{parentName:"tr",align:null},"NO"),(0,r.kt)("td",{parentName:"tr",align:null})),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"tags")),(0,r.kt)("td",{parentName:"tr",align:null},"filter by tag key:value pairs"),(0,r.kt)("td",{parentName:"tr",align:null},"NO"),(0,r.kt)("td",{parentName:"tr",align:null})))),(0,r.kt)("p",null,(0,r.kt)("em",{parentName:"p"},"note")," : user can use wildcard ","*"," in filters\n",(0,r.kt)("em",{parentName:"p"},"note")," : when using tag filters the query string must follow the pattern: ",(0,r.kt)("inlineCode",{parentName:"p"},"?tags=key1:value1,key2:value2"),"\n",(0,r.kt)("em",{parentName:"p"},"note")," : You can use ",(0,r.kt)("inlineCode",{parentName:"p"},"~")," as a negative operator in the beginning of a filter value to exclude something: For example you can exclude endpoints with subgroup of value ",(0,r.kt)("inlineCode",{parentName:"p"},"GROUP_A")," by issuing ",(0,r.kt)("inlineCode",{parentName:"p"},"?subgroup:~SERVICE_A")),(0,r.kt)("h4",{id:"headers-1"},"Headers"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"x-api-key: secret_key_value\nAccept: application/json\n")),(0,r.kt)("h4",{id:"example-request"},"Example Request"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"GET /topology/groups?date=2015-07-22\n")),(0,r.kt)("h4",{id:"response-code-2"},"Response Code"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"Status: 200 OK\n")),(0,r.kt)("h3",{id:"response-body-2"},"Response body"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "status": {\n        "message": "Success",\n        "code": "200"\n    },\n    "data": [\n        {\n            "date": "2015-07-22",\n            "group": "NGIA",\n            "type": "NGIS",\n            "subgroup": "SITEA",\n            "tags": {\n                "certification": "Certified",\n                "infrastructure": "Production"\n            }\n        },\n        {\n            "date": "2015-07-22",\n            "group": "NGIA",\n            "type": "NGIS",\n            "subgroup": "SITEB",\n            "tags": {\n                "certification": "Certified",\n                "infrastructure": "Production"\n            },\n            "notifications": {\n                "contacts": ["email01@example.com"],\n                "enabled": true\n            }\n        }\n    ]\n}\n')),(0,r.kt)("a",{id:"3"}),(0,r.kt)("h2",{id:"delete-delete-group-topology-for-a-specific-date"},"[DELETE]",": Delete group topology for a specific date"),(0,r.kt)("p",null,"This method can be used to delete all group items contributing to the group topology of a specific date"),(0,r.kt)("h3",{id:"input-2"},"Input"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"DELETE /topology/groups?date=YYYY-MM-DD\n")),(0,r.kt)("h4",{id:"request-headers"},"Request headers"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nContent-Type: application/json\nAccept: application/json\n")),(0,r.kt)("h3",{id:"response"},"Response"),(0,r.kt)("p",null,"Headers: ",(0,r.kt)("inlineCode",{parentName:"p"},"Status: 200 OK")),(0,r.kt)("h4",{id:"response-body-3"},"Response body"),(0,r.kt)("p",null,"Json Response"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "message": "Topology of 3 groups deleted for date: 2019-12-12",\n    "code": "200"\n}\n')),(0,r.kt)("a",{id:"4"}),(0,r.kt)("h2",{id:"get-list-group-topology-for-specific-report"},"GET: List group topology for specific report"),(0,r.kt)("p",null,"Lists group topology items for specific report"),(0,r.kt)("h3",{id:"input-3"},"Input"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"GET /topology/groups/by_report/{report-name}?date=YYYY-MM-DD\n")),(0,r.kt)("h4",{id:"url-parameters-2"},"Url Parameters"),(0,r.kt)("table",null,(0,r.kt)("thead",{parentName:"table"},(0,r.kt)("tr",{parentName:"thead"},(0,r.kt)("th",{parentName:"tr",align:null},"Type"),(0,r.kt)("th",{parentName:"tr",align:null},"Description"),(0,r.kt)("th",{parentName:"tr",align:null},"Required"),(0,r.kt)("th",{parentName:"tr",align:null},"Default value"))),(0,r.kt)("tbody",{parentName:"table"},(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"report-name")),(0,r.kt)("td",{parentName:"tr",align:null},"target a specific report"),(0,r.kt)("td",{parentName:"tr",align:null},"YES"),(0,r.kt)("td",{parentName:"tr",align:null},"none")),(0,r.kt)("tr",{parentName:"tbody"},(0,r.kt)("td",{parentName:"tr",align:null},(0,r.kt)("inlineCode",{parentName:"td"},"date")),(0,r.kt)("td",{parentName:"tr",align:null},"target a specific date"),(0,r.kt)("td",{parentName:"tr",align:null},"NO"),(0,r.kt)("td",{parentName:"tr",align:null},"today's date")))),(0,r.kt)("h4",{id:"headers-2"},"Headers"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"x-api-key: secret_key_value\nAccept: application/json\n")),(0,r.kt)("h4",{id:"example-request-1"},"Example Request"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"GET /topology/groups/by_report/Critical?date=2015-07-22\n")),(0,r.kt)("h4",{id:"response-code-3"},"Response Code"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre"},"Status: 200 OK\n")),(0,r.kt)("h3",{id:"response-body-4"},"Response body"),(0,r.kt)("pre",null,(0,r.kt)("code",{parentName:"pre",className:"language-json"},'{\n    "status": {\n        "message": "Success",\n        "code": "200"\n    },\n    "data": [\n        {\n            "date": "2015-07-22",\n            "group": "NGIA",\n            "type": "NGIS",\n            "subgroup": "SITEA",\n            "tags": {\n                "certification": "Certified",\n                "infrastructure": "Production"\n            }\n        },\n        {\n            "date": "2015-07-22",\n            "group": "NGIA",\n            "type": "NGIS",\n            "subgroup": "SITEB",\n            "tags": {\n                "certification": "Certified",\n                "infrastructure": "Production"\n            },\n            "notifications": {\n                "contacts": ["email01@example.com"],\n                "enabled": true\n            }\n        },\n        {\n            "date": "2015-07-22",\n            "group": "NGIX",\n            "type": "NGIS",\n            "subgroup": "SITEX",\n            "tags": {\n                "certification": "Certified",\n                "infrastructure": "Production"\n            }\n        }\n    ]\n}\n')))}d.isMDXComponent=!0}}]);