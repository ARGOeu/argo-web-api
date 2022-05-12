(window.webpackJsonp=window.webpackJsonp||[]).push([[11],{67:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return c})),n.d(t,"metadata",(function(){return l})),n.d(t,"rightToc",(function(){return b})),n.d(t,"default",(function(){return p}));var a=n(2),r=n(6),i=(n(0),n(93)),s=["components"],c={id:"metrics",title:"Available Metrics and Tags"},l={unversionedId:"metrics",id:"metrics",isDocsHomePage:!1,title:"Available Metrics and Tags",description:"API Calls",source:"@site/docs/metrics.md",slug:"/metrics",permalink:"/argo-web-api/docs/metrics",version:"current",sidebar:"someSidebar",previous:{title:"weights",permalink:"/argo-web-api/docs/weights"},next:{title:"Operation Profiles",permalink:"/argo-web-api/docs/operations_profiles"}},b=[{value:"API Calls",id:"api-calls",children:[]},{value:"GET: List Metrics (Admin)",id:"get-list-metrics-admin",children:[{value:"Input",id:"input",children:[]},{value:"Request headers",id:"request-headers",children:[]},{value:"Response",id:"response",children:[]}]},{value:"PUT: Update Metrics information",id:"put-update-metrics-information",children:[{value:"Input",id:"input-1",children:[]},{value:"Response",id:"response-1",children:[]}]},{value:"GET: List Metrics (as a tenant user)",id:"get-list-metrics-as-a-tenant-user",children:[{value:"Input",id:"input-2",children:[]},{value:"Request headers",id:"request-headers-2",children:[]},{value:"Response",id:"response-2",children:[]}]},{value:"PUT: List metrics by report (as a tenant user)",id:"put-list-metrics-by-report-as-a-tenant-user",children:[{value:"Input",id:"input-3",children:[]},{value:"Response",id:"response-3",children:[]}]}],o={rightToc:b};function p(e){var t=e.components,n=Object(r.a)(e,s);return Object(i.b)("wrapper",Object(a.a)({},o,n,{components:t,mdxType:"MDXLayout"}),Object(i.b)("h2",{id:"api-calls"},"API Calls"),Object(i.b)("table",null,Object(i.b)("thead",{parentName:"table"},Object(i.b)("tr",{parentName:"thead"},Object(i.b)("th",{parentName:"tr",align:null},"Name"),Object(i.b)("th",{parentName:"tr",align:null},"Description"),Object(i.b)("th",{parentName:"tr",align:null},"Shortcut"))),Object(i.b)("tbody",{parentName:"table"},Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",{parentName:"tr",align:null},"GET: List Metrics (Admin)"),Object(i.b)("td",{parentName:"tr",align:null},"This method can be used to retrieve a list of all metrics"),Object(i.b)("td",{parentName:"tr",align:null},Object(i.b)("a",{parentName:"td",href:"#1"}," Description"))),Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",{parentName:"tr",align:null},"PUT: Update Metrics (Admin)"),Object(i.b)("td",{parentName:"tr",align:null},"This method can be used to update the list of metrics"),Object(i.b)("td",{parentName:"tr",align:null},Object(i.b)("a",{parentName:"td",href:"#2"}," Description"))),Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",{parentName:"tr",align:null},"GET: List Metrics"),Object(i.b)("td",{parentName:"tr",align:null},"This method can be used to retrieve a list of metrics (as a tenant user)"),Object(i.b)("td",{parentName:"tr",align:null},Object(i.b)("a",{parentName:"td",href:"#3"}," Description"))),Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",{parentName:"tr",align:null},"PUT: List Metrics by report"),Object(i.b)("td",{parentName:"tr",align:null},"This method can be used to retrieve a list of metrics included in a report (as a tenant user)"),Object(i.b)("td",{parentName:"tr",align:null},Object(i.b)("a",{parentName:"td",href:"#4"}," Description"))))),Object(i.b)("a",{id:"1"}),Object(i.b)("h2",{id:"get-list-metrics-admin"},"[GET]",": List Metrics (Admin)"),Object(i.b)("p",null,"This method can be used to retrieve a list of all metrics. This is an administrative method. The Metric list is common for all tenants"),Object(i.b)("h3",{id:"input"},"Input"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre"},"GET /admin/metrics\n")),Object(i.b)("h3",{id:"request-headers"},"Request headers"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),Object(i.b)("h3",{id:"response"},"Response"),Object(i.b)("p",null,"Headers: ",Object(i.b)("inlineCode",{parentName:"p"},"Status: 200 OK")),Object(i.b)("h4",{id:"response-body"},"Response body"),Object(i.b)("p",null,"Json Response"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-json"},'{\n  "status": {\n    "message": "Success",\n    "code": "200"\n  },\n  "data": [\n    {\n      "name": "test_metric_1",\n      "tags": [\n        "network",\n        "internal"\n      ]\n    },\n    {\n      "name": "test_metric_2",\n      "tags": [\n        "disk",\n        "agent"\n      ]\n    },\n    {\n      "name": "test_metric_3",\n      "tags": [\n        "aai"\n      ]\n    }\n  ]\n}\n')),Object(i.b)("a",{id:"2"}),Object(i.b)("h2",{id:"put-update-metrics-information"},"[PUT]",": Update Metrics information"),Object(i.b)("p",null,"This method is used to update the list of metrics. This is an administrative method. The list of metrics is common for all tenants"),Object(i.b)("h3",{id:"input-1"},"Input"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre"},"PUT /admin/metrics\n")),Object(i.b)("h4",{id:"put-body"},"PUT BODY"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-json"},'  [\n  {\n    "name": "metric1",\n    "tags": [\n      "tag1",\n      "tag2"\n    ]\n  }\n]\n')),Object(i.b)("h4",{id:"request-headers-1"},"Request headers"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),Object(i.b)("h3",{id:"response-1"},"Response"),Object(i.b)("p",null,"Headers: ",Object(i.b)("inlineCode",{parentName:"p"},"Status: 200 OK")),Object(i.b)("h4",{id:"response-body-1"},"Response body"),Object(i.b)("p",null,"Json Response"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-json"},'{\n  "status": {\n    "message": "Metrics resource succesfully updated",\n    "code": "200"\n  },\n  "data": [\n    {\n      "name": "metric1",\n      "tags": [\n        "tag1",\n        "tag2"\n      ]\n    }\n  ]\n}\n')),Object(i.b)("a",{id:"3"}),Object(i.b)("h2",{id:"get-list-metrics-as-a-tenant-user"},"[GET]",": List Metrics (as a tenant user)"),Object(i.b)("p",null,"This method can be used to retrieve the list of metrics as a tenant user. The list of metrics is common for all tenants but accessible from each tenant."),Object(i.b)("h3",{id:"input-2"},"Input"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre"},"GET /metrics\n")),Object(i.b)("h3",{id:"request-headers-2"},"Request headers"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),Object(i.b)("h3",{id:"response-2"},"Response"),Object(i.b)("p",null,"Headers: ",Object(i.b)("inlineCode",{parentName:"p"},"Status: 200 OK")),Object(i.b)("h4",{id:"response-body-2"},"Response body"),Object(i.b)("p",null,"Json Response"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-json"},'{\n  "status": {\n    "message": "Success",\n    "code": "200"\n  },\n  "data": [\n    {\n      "name": "test_metric_1",\n      "tags": [\n        "network",\n        "internal"\n      ]\n    },\n    {\n      "name": "test_metric_2",\n      "tags": [\n        "disk",\n        "agent"\n      ]\n    },\n    {\n      "name": "test_metric_3",\n      "tags": [\n        "aai"\n      ]\n    }\n  ]\n}\n')),Object(i.b)("a",{id:"4"}),Object(i.b)("h2",{id:"put-list-metrics-by-report-as-a-tenant-user"},"[PUT]",": List metrics by report (as a tenant user)"),Object(i.b)("p",null,"This method is used to retrieve a list of metrics that are included in the metric profile of a specific report."),Object(i.b)("h3",{id:"input-3"},"Input"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre"},"PUT /metrics/by_report/{report_name}\n")),Object(i.b)("h4",{id:"url-parameters"},"Url Parameters"),Object(i.b)("table",null,Object(i.b)("thead",{parentName:"table"},Object(i.b)("tr",{parentName:"thead"},Object(i.b)("th",{parentName:"tr",align:null},"Type"),Object(i.b)("th",{parentName:"tr",align:null},"Description"),Object(i.b)("th",{parentName:"tr",align:null},"Required"),Object(i.b)("th",{parentName:"tr",align:null},"Default value"))),Object(i.b)("tbody",{parentName:"table"},Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",{parentName:"tr",align:null},Object(i.b)("inlineCode",{parentName:"td"},"report_name")),Object(i.b)("td",{parentName:"tr",align:null},"target a specific report"),Object(i.b)("td",{parentName:"tr",align:null},"YES"),Object(i.b)("td",{parentName:"tr",align:null},"none")),Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",{parentName:"tr",align:null},Object(i.b)("inlineCode",{parentName:"td"},"date")),Object(i.b)("td",{parentName:"tr",align:null},"target a specific date"),Object(i.b)("td",{parentName:"tr",align:null},"NO"),Object(i.b)("td",{parentName:"tr",align:null},"today's date")))),Object(i.b)("h4",{id:"request-headers-3"},"Request headers"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json\n")),Object(i.b)("h3",{id:"response-3"},"Response"),Object(i.b)("p",null,"Some metric results have additional information regarding the specific service endpoint such as it's Url, certificat DN etc... If this information is available it will be displayed under each service endpoint along with status results. Also some metrics might have a changed status due to a defined threshold rule being applied (see more about ",Object(i.b)("a",{parentName:"p",href:"/argo-web-api/docs/threshold_profiles"},"Threshold profiles"),"). Thus they will include additional information such as the original status value (field name: ",Object(i.b)("inlineCode",{parentName:"p"},"original_status"),"), the threshold rule applied (field name: ",Object(i.b)("inlineCode",{parentName:"p"},"threshold_rule_applied"),") and the actual data (field name: ",Object(i.b)("inlineCode",{parentName:"p"},"actual_data"),") that the rule has been applied to. For example:"),Object(i.b)("p",null,"Headers: ",Object(i.b)("inlineCode",{parentName:"p"},"Status: 200 OK")),Object(i.b)("h4",{id:"response-body-3"},"Response body"),Object(i.b)("p",null,"Json Response"),Object(i.b)("pre",null,Object(i.b)("code",{parentName:"pre",className:"language-json"},'{\n   "root": [\n     {\n       "Name": "www.example.com",\n       "info": {\n                  "Url": "https://example.com/path/to/service/check"\n               },\n       "Metrics": [\n         {\n           "Name": "httpd_check",\n           "Service": "httpd",\n           "Details": [\n             {\n               "Timestamp": "2015-06-20T12:00:00Z",\n               "Value": "OK",\n               "Summary": "httpd is ok",\n               "Message": "all checks ok"\n             },\n             {\n               "Timestamp": "2015-06-20T23:00:00Z",\n               "Value": "OK",\n               "Summary": "httpd is ok",\n               "Message": "all checks ok"\n             }\n           ]\n         },\n         {\n           "Name": "httpd_memory",\n           "Service": "httpd",\n           "Details": [\n             {\n               "Timestamp": "2015-06-20T06:00:00Z",\n               "Value": "OK",\n               "Summary": "memcheck ok",\n               "Message": "memory under 20%"\n             },\n             {\n               "Timestamp": "2015-06-20T09:00:00Z",\n               "Value": "OK",\n               "Summary": "memcheck ok",\n               "Message": "memory under 20%"\n             },\n             {\n               "Timestamp": "2015-06-20T18:00:00Z",\n               "Value": "CRITICAL",\n               "Summary": "memcheck ok",\n               "Message": "memory under 20%",\n               "original_status": "OK",\n               "threshold_rule_applied": "reserved_memory=0.1;0.1:0.2;0.2:0.5",\n               "actual_data": "reserved_memory=0.15"\n             },\n           ]\n         }\n       ]\n     }\n   ]\n }\n')))}p.isMDXComponent=!0},93:function(e,t,n){"use strict";n.d(t,"a",(function(){return p})),n.d(t,"b",(function(){return m}));var a=n(0),r=n.n(a);function i(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function s(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function c(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?s(Object(n),!0).forEach((function(t){i(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):s(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},i=Object.keys(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var b=r.a.createContext({}),o=function(e){var t=r.a.useContext(b),n=t;return e&&(n="function"==typeof e?e(t):c(c({},t),e)),n},p=function(e){var t=o(e.components);return r.a.createElement(b.Provider,{value:t},e.children)},d={inlineCode:"code",wrapper:function(e){var t=e.children;return r.a.createElement(r.a.Fragment,{},t)}},u=r.a.forwardRef((function(e,t){var n=e.components,a=e.mdxType,i=e.originalType,s=e.parentName,b=l(e,["components","mdxType","originalType","parentName"]),p=o(n),u=a,m=p["".concat(s,".").concat(u)]||p[u]||d[u]||i;return n?r.a.createElement(m,c(c({ref:t},b),{},{components:n})):r.a.createElement(m,c({ref:t},b))}));function m(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var i=n.length,s=new Array(i);s[0]=u;var c={};for(var l in t)hasOwnProperty.call(t,l)&&(c[l]=t[l]);c.originalType=e,c.mdxType="string"==typeof e?e:a,s[1]=c;for(var b=2;b<i;b++)s[b]=n[b];return r.a.createElement.apply(null,s)}return r.a.createElement.apply(null,n)}u.displayName="MDXCreateElement"}}]);