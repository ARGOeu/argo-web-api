(window.webpackJsonp=window.webpackJsonp||[]).push([[18],{74:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return c})),n.d(t,"metadata",(function(){return b})),n.d(t,"rightToc",(function(){return i})),n.d(t,"default",(function(){return p}));var a=n(2),o=n(6),r=(n(0),n(92)),c={id:"topology_endpoints",title:"Topology Endpoints"},b={unversionedId:"topology_endpoints",id:"topology_endpoints",isDocsHomePage:!1,title:"Topology Endpoints",description:"API calls for handling topology endpoint resources",source:"@site/docs/topology_endpoints.md",slug:"/topology_endpoints",permalink:"/argo-web-api/docs/topology_endpoints",version:"current",sidebar:"someSidebar",previous:{title:"Downtimes",permalink:"/argo-web-api/docs/downtimes"},next:{title:"Topology Groups",permalink:"/argo-web-api/docs/topology_groups"}},i=[{value:"API calls for handling topology endpoint resources",id:"api-calls-for-handling-topology-endpoint-resources",children:[]},{value:"POST: Create endpoint topology for specific date",id:"post-create-endpoint-topology-for-specific-date",children:[{value:"Input",id:"input",children:[]},{value:"POST BODY",id:"post-body",children:[]},{value:"Response body",id:"response-body",children:[]}]},{value:"409 Conflict when trying to insert a topology that already exists",id:"409-conflict-when-trying-to-insert-a-topology-that-already-exists",children:[{value:"Response Code",id:"response-code-1",children:[]},{value:"Response body",id:"response-body-1",children:[]}]},{value:"GET: List endpoint topology per date",id:"get-list-endpoint-topology-per-date",children:[{value:"Input",id:"input-1",children:[]},{value:"Response body",id:"response-body-2",children:[]}]},{value:"DELETE: Delete endpoint topology for a specific date",id:"delete-delete-endpoint-topology-for-a-specific-date",children:[{value:"Input",id:"input-2",children:[]},{value:"Response",id:"response",children:[]}]},{value:"GET: List endpoint topology for specific report",id:"get-list-endpoint-topology-for-specific-report",children:[{value:"Input",id:"input-3",children:[]},{value:"Response body",id:"response-body-4",children:[]}]}],l={rightToc:i};function p(e){var t=e.components,n=Object(o.a)(e,["components"]);return Object(r.b)("wrapper",Object(a.a)({},l,n,{components:t,mdxType:"MDXLayout"}),Object(r.b)("h2",{id:"api-calls-for-handling-topology-endpoint-resources"},"API calls for handling topology endpoint resources"),Object(r.b)("table",null,Object(r.b)("thead",{parentName:"table"},Object(r.b)("tr",{parentName:"thead"},Object(r.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Name"),Object(r.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Description"),Object(r.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Shortcut"))),Object(r.b)("tbody",{parentName:"table"},Object(r.b)("tr",{parentName:"tbody"},Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"POST: Create endpoint topology for specific date"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"Creates a daily endpoint topology mapping endpoints to endpoint groups"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(r.b)("a",{href:"#1"},"Description"))),Object(r.b)("tr",{parentName:"tbody"},Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"GET: List endpoint topology for specific date"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"Lists endpoint topology for a specific date"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(r.b)("a",{href:"#2"},"Description"))),Object(r.b)("tr",{parentName:"tbody"},Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"DELETE: delete endpoint topology for specific date"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"Deletes all endpoint items (topology) for a specific date"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(r.b)("a",{href:"#3"},"Description"))),Object(r.b)("tr",{parentName:"tbody"},Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"GET: List endpoint topology for specific report"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"Lists endpoint topology for a specific report"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(r.b)("a",{href:"#4"},"Description"))))),Object(r.b)("a",{id:"1"}),Object(r.b)("h2",{id:"post-create-endpoint-topology-for-specific-date"},"POST: Create endpoint topology for specific date"),Object(r.b)("p",null,"Creates a daily endpoint topology mapping endpoints to endpoint groups"),Object(r.b)("h3",{id:"input"},"Input"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{}),"POST /topology/endpoints?date=YYYY-MM-DD\n")),Object(r.b)("h4",{id:"url-parameters"},"Url Parameters"),Object(r.b)("table",null,Object(r.b)("thead",{parentName:"table"},Object(r.b)("tr",{parentName:"thead"},Object(r.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Type"),Object(r.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Description"),Object(r.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Required"),Object(r.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Default value"))),Object(r.b)("tbody",{parentName:"table"},Object(r.b)("tr",{parentName:"tbody"},Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(r.b)("inlineCode",{parentName:"td"},"date")),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"target a specific date"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"NO"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"today's date")))),Object(r.b)("h4",{id:"headers"},"Headers"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{}),"x-api-key: secret_key_value\nAccept: application/json\n")),Object(r.b)("h3",{id:"post-body"},"POST BODY"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{className:"language-json"}),'[\n    {\n        "group": "SITE_A",\n        "hostname": "host1.site-a.foo",\n        "type": "SITES",\n        "service": "a.service.foo",\n        "tags": { "scope": "TENANT", "production": "1", "monitored": "1" }\n    },\n    {\n        "group": "SITE_A",\n        "hostname": "host2.site-b.foo",\n        "type": "SITES",\n        "service": "b.service.foo",\n        "tags": { "scope": "TENANT", "production": "1", "monitored": "1" }\n    },\n    {\n        "group": "SITE_B",\n        "hostname": "host1.site-a.foo",\n        "type": "SITES",\n        "service": "c.service.foo",\n        "tags": { "scope": "TENANT", "production": "1", "monitored": "1" },\n        "notifications": {"contacts": ["email01@example.com"], "enabled": true}\n    }\n]\n')),Object(r.b)("h4",{id:"response-code"},"Response Code"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{}),"Status: 201 OK Created\n")),Object(r.b)("h3",{id:"response-body"},"Response body"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{className:"language-json"}),'{\n    "message": "Topology of 3 endpoints created for date: YYYY-MM-DD",\n    "code": "201"\n}\n')),Object(r.b)("h2",{id:"409-conflict-when-trying-to-insert-a-topology-that-already-exists"},"409 Conflict when trying to insert a topology that already exists"),Object(r.b)("p",null,"When trying to insert a topology for a specific date that already exists the api will answer with the following reponse:"),Object(r.b)("h3",{id:"response-code-1"},"Response Code"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{}),"Status: 409 Conflict\n")),Object(r.b)("h3",{id:"response-body-1"},"Response body"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{className:"language-json"}),'{\n    "message": "topology already exists for date: YYYY-MM-DD, please either update it or delete it first!",\n    "code": "409"\n}\n')),Object(r.b)("p",null,"User can proceed with either updating the existing endpoint topology OR deleting before trying to create it anew"),Object(r.b)("a",{id:"2"}),Object(r.b)("h2",{id:"get-list-endpoint-topology-per-date"},"GET: List endpoint topology per date"),Object(r.b)("p",null,"List endpoint topology for a specific date or the closest availabe topology to that date. If date is not provided list the latest available endpoint topology."),Object(r.b)("h3",{id:"input-1"},"Input"),Object(r.b)("h5",{id:"list-all-topology-statistics"},"List All topology statistics"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{}),"GET /topology/endpoints?date=YYYY-MM-DD\n")),Object(r.b)("h4",{id:"url-parameters-1"},"Url Parameters"),Object(r.b)("table",null,Object(r.b)("thead",{parentName:"table"},Object(r.b)("tr",{parentName:"thead"},Object(r.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Type"),Object(r.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Description"),Object(r.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Required"),Object(r.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Default value"))),Object(r.b)("tbody",{parentName:"table"},Object(r.b)("tr",{parentName:"tbody"},Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(r.b)("inlineCode",{parentName:"td"},"date")),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"target a specific date"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"NO"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"today's date")),Object(r.b)("tr",{parentName:"tbody"},Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(r.b)("inlineCode",{parentName:"td"},"group")),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"filter by group name"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"NO"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}))),Object(r.b)("tr",{parentName:"tbody"},Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(r.b)("inlineCode",{parentName:"td"},"type")),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"filter by group type"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"NO"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}))),Object(r.b)("tr",{parentName:"tbody"},Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(r.b)("inlineCode",{parentName:"td"},"service")),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"filter by service"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"NO"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}))),Object(r.b)("tr",{parentName:"tbody"},Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(r.b)("inlineCode",{parentName:"td"},"hostname")),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"filter by hostname"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"NO"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}))),Object(r.b)("tr",{parentName:"tbody"},Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(r.b)("inlineCode",{parentName:"td"},"tags")),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"filter by tag key:value pairs"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"NO"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}))))),Object(r.b)("p",null,Object(r.b)("em",{parentName:"p"},"note")," : user can use wildcard ","*"," in filters\n",Object(r.b)("em",{parentName:"p"},"note")," : when using tag filters the query string must follow the pattern: ",Object(r.b)("inlineCode",{parentName:"p"},"?tags=key1:value1,key2:value2"),"\n",Object(r.b)("em",{parentName:"p"},"note")," : You can use ",Object(r.b)("inlineCode",{parentName:"p"},"~")," as a negative operator in the beginning of a filter value to exclude something: For example you can exclude endpoints with service of value ",Object(r.b)("inlineCode",{parentName:"p"},"SERVICE_A")," by issuing ",Object(r.b)("inlineCode",{parentName:"p"},"?service:~SERVICE_A")),Object(r.b)("h4",{id:"headers-1"},"Headers"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{}),"x-api-key: secret_key_value\nAccept: application/json\n")),Object(r.b)("h4",{id:"response-code-2"},"Response Code"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{}),"Status: 200 OK\n")),Object(r.b)("h3",{id:"response-body-2"},"Response body"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{className:"language-json"}),'{\n    "status": {\n        "message": "Success",\n        "code": "200"\n    },\n    "data": [\n        {\n            "date": "2019-12-12",\n            "group": "SITE_A",\n            "hostname": "host1.site-a.foo",\n            "type": "SITES",\n            "service": "a.service.foo",\n            "tags": {\n                "scope": "TENANT",\n                "production": "1",\n                "monitored": "1"\n            }\n        },\n        {\n            "date": "2019-12-12",\n            "group": "SITE_A",\n            "hostname": "host2.site-b.foo",\n            "type": "SITES",\n            "service": "b.service.foo",\n            "tags": {\n                "scope": "TENANT",\n                "production": "1",\n                "monitored": "1"\n            }\n        },\n        {\n            "date": "2019-12-12",\n            "group": "SITE_B",\n            "hostname": "host1.site-a.foo",\n            "type": "SITES",\n            "service": "c.service.foo",\n            "tags": {\n                "scope": "TENANT",\n                "production": "1",\n                "monitored": "1"\n            },\n            "notifications": {\n                "contacts": ["email01@example.com"],\n                "enabled": true\n            }\n        }\n    ]\n}\n')),Object(r.b)("a",{id:"3"}),Object(r.b)("h2",{id:"delete-delete-endpoint-topology-for-a-specific-date"},"[DELETE]",": Delete endpoint topology for a specific date"),Object(r.b)("p",null,"This method can be used to delete all endpoint items contributing to the endpoint topology of a specific date"),Object(r.b)("h3",{id:"input-2"},"Input"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{}),"DELETE /topology/endpoints?date=YYYY-MM-DD\n")),Object(r.b)("h4",{id:"request-headers"},"Request headers"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{}),"x-api-key: shared_key_value\nContent-Type: application/json\nAccept: application/json\n")),Object(r.b)("h3",{id:"response"},"Response"),Object(r.b)("p",null,"Headers: ",Object(r.b)("inlineCode",{parentName:"p"},"Status: 200 OK")),Object(r.b)("h4",{id:"response-body-3"},"Response body"),Object(r.b)("p",null,"Json Response"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{className:"language-json"}),'{\n    "message": "Topology of 3 endpoints deleted for date: 2019-12-12",\n    "code": "200"\n}\n')),Object(r.b)("a",{id:"4"}),Object(r.b)("h2",{id:"get-list-endpoint-topology-for-specific-report"},"GET: List endpoint topology for specific report"),Object(r.b)("p",null,"Lists endpoint topology items for specific report"),Object(r.b)("h3",{id:"input-3"},"Input"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{}),"GET /topology/endpoint/by_report/{report-name}?date=YYYY-MM-DD\n")),Object(r.b)("h4",{id:"url-parameters-2"},"Url Parameters"),Object(r.b)("table",null,Object(r.b)("thead",{parentName:"table"},Object(r.b)("tr",{parentName:"thead"},Object(r.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Type"),Object(r.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Description"),Object(r.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Required"),Object(r.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Default value"))),Object(r.b)("tbody",{parentName:"table"},Object(r.b)("tr",{parentName:"tbody"},Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(r.b)("inlineCode",{parentName:"td"},"report-name")),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"target a specific report"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"YES"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"none")),Object(r.b)("tr",{parentName:"tbody"},Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(r.b)("inlineCode",{parentName:"td"},"date")),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"target a specific date"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"NO"),Object(r.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"today's date")))),Object(r.b)("h4",{id:"headers-2"},"Headers"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{}),"x-api-key: secret_key_value\nAccept: application/json\n")),Object(r.b)("h4",{id:"example-request"},"Example Request"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{}),"GET /topology/endpoints/by_report/Critical?date=2015-07-22\n")),Object(r.b)("h4",{id:"response-code-3"},"Response Code"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{}),"Status: 200 OK\n")),Object(r.b)("h3",{id:"response-body-4"},"Response body"),Object(r.b)("pre",null,Object(r.b)("code",Object(a.a)({parentName:"pre"},{className:"language-json"}),'{\n    "status": {\n        "message": "Success",\n        "code": "200"\n    },\n    "data": [\n        {\n            "date": "2019-12-12",\n            "group": "SITE_A",\n            "hostname": "host1.site-a.foo",\n            "type": "SITES",\n            "service": "a.service.foo",\n            "tags": {\n                "scope": "TENANT",\n                "production": "1",\n                "monitored": "1"\n            }\n        },\n        {\n            "date": "2019-12-12",\n            "group": "SITE_A",\n            "hostname": "host2.site-b.foo",\n            "type": "SITES",\n            "service": "b.service.foo",\n            "tags": {\n                "scope": "TENANT",\n                "production": "1",\n                "monitored": "1"\n            }\n        },\n        {\n            "date": "2019-12-12",\n            "group": "SITE_B",\n            "hostname": "host1.site-a.foo",\n            "type": "SITES",\n            "service": "c.service.foo",\n            "tags": {\n                "scope": "TENANT",\n                "production": "1",\n                "monitored": "1"\n            },\n            "notifications": {\n                "contacts": ["email01@example.com"],\n                "enabled": true\n            }\n        }\n    ]\n}\n')))}p.isMDXComponent=!0},92:function(e,t,n){"use strict";n.d(t,"a",(function(){return d})),n.d(t,"b",(function(){return j}));var a=n(0),o=n.n(a);function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function c(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function b(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?c(Object(n),!0).forEach((function(t){r(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):c(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function i(e,t){if(null==e)return{};var n,a,o=function(e,t){if(null==e)return{};var n,a,o={},r=Object.keys(e);for(a=0;a<r.length;a++)n=r[a],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);for(a=0;a<r.length;a++)n=r[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var l=o.a.createContext({}),p=function(e){var t=o.a.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):b(b({},t),e)),n},d=function(e){var t=p(e.components);return o.a.createElement(l.Provider,{value:t},e.children)},s={inlineCode:"code",wrapper:function(e){var t=e.children;return o.a.createElement(o.a.Fragment,{},t)}},O=o.a.forwardRef((function(e,t){var n=e.components,a=e.mdxType,r=e.originalType,c=e.parentName,l=i(e,["components","mdxType","originalType","parentName"]),d=p(n),O=a,j=d["".concat(c,".").concat(O)]||d[O]||s[O]||r;return n?o.a.createElement(j,b(b({ref:t},l),{},{components:n})):o.a.createElement(j,b({ref:t},l))}));function j(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var r=n.length,c=new Array(r);c[0]=O;var b={};for(var i in t)hasOwnProperty.call(t,i)&&(b[i]=t[i]);b.originalType=e,b.mdxType="string"==typeof e?e:a,c[1]=b;for(var l=2;l<r;l++)c[l]=n[l];return o.a.createElement.apply(null,c)}return o.a.createElement.apply(null,n)}O.displayName="MDXCreateElement"}}]);