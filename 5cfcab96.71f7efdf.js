(window.webpackJsonp=window.webpackJsonp||[]).push([[16],{72:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return b})),n.d(t,"metadata",(function(){return i})),n.d(t,"rightToc",(function(){return p})),n.d(t,"default",(function(){return o}));var a=n(2),r=n(6),c=(n(0),n(92)),b={id:"issues",title:"issues"},i={unversionedId:"issues",id:"issues",isDocsHomePage:!1,title:"issues",description:"API calls to quicky find endpoints with issues",source:"@site/docs/issues.md",slug:"/issues",permalink:"/argo-web-api/docs/issues",version:"current",sidebar:"someSidebar",previous:{title:"Status Results",permalink:"/argo-web-api/docs/status_results"},next:{title:"trends",permalink:"/argo-web-api/docs/trends"}},p=[{value:"API calls to quicky find endpoints with issues",id:"api-calls-to-quicky-find-endpoints-with-issues",children:[{value:"GET: Endpoint with Issues",id:"get-endpoint-with-issues",children:[]},{value:"Input",id:"input",children:[]},{value:"Response body",id:"response-body",children:[]}]}],l={rightToc:p};function o(e){var t=e.components,n=Object(r.a)(e,["components"]);return Object(c.b)("wrapper",Object(a.a)({},l,n,{components:t,mdxType:"MDXLayout"}),Object(c.b)("h2",{id:"api-calls-to-quicky-find-endpoints-with-issues"},"API calls to quicky find endpoints with issues"),Object(c.b)("h3",{id:"get-endpoint-with-issues"},"[GET]",": Endpoint with Issues"),Object(c.b)("p",null,"This method may be used to retrieve a list of problematic endpoints"),Object(c.b)("h3",{id:"input"},"Input"),Object(c.b)("pre",null,Object(c.b)("code",Object(a.a)({parentName:"pre"},{}),"/issues/{report_name}/endpoints?date=2020-05-01&filter=CRITICAL\n")),Object(c.b)("h4",{id:"path-parameters"},"Path Parameters"),Object(c.b)("table",null,Object(c.b)("thead",{parentName:"table"},Object(c.b)("tr",{parentName:"thead"},Object(c.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Type"),Object(c.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Description"),Object(c.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Required"),Object(c.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Default value"))),Object(c.b)("tbody",{parentName:"table"},Object(c.b)("tr",{parentName:"tbody"},Object(c.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(c.b)("inlineCode",{parentName:"td"},"report_name")),Object(c.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"name of the report"),Object(c.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"YES"),Object(c.b)("td",Object(a.a)({parentName:"tr"},{align:null}))))),Object(c.b)("h4",{id:"url-parameters"},"Url Parameters"),Object(c.b)("table",null,Object(c.b)("thead",{parentName:"table"},Object(c.b)("tr",{parentName:"thead"},Object(c.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Type"),Object(c.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Description"),Object(c.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Required"),Object(c.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Default value"))),Object(c.b)("tbody",{parentName:"table"},Object(c.b)("tr",{parentName:"tbody"},Object(c.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(c.b)("inlineCode",{parentName:"td"},"date")),Object(c.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"Date to view problematic endpoints of"),Object(c.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"NO"),Object(c.b)("td",Object(a.a)({parentName:"tr"},{align:null}))),Object(c.b)("tr",{parentName:"tbody"},Object(c.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(c.b)("inlineCode",{parentName:"td"},"filter")),Object(c.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"Filter (optinally) problematic endpoints by status value"),Object(c.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"NO"),Object(c.b)("td",Object(a.a)({parentName:"tr"},{align:null}))))),Object(c.b)("h4",{id:"headers"},"Headers"),Object(c.b)("pre",null,Object(c.b)("code",Object(a.a)({parentName:"pre"},{}),"x-api-key: shared_key_value\nAccept: application/json or application/xml\n")),Object(c.b)("h4",{id:"response-code"},"Response Code"),Object(c.b)("pre",null,Object(c.b)("code",Object(a.a)({parentName:"pre"},{}),"Status: 200 OK\n")),Object(c.b)("h3",{id:"response-body"},"Response body"),Object(c.b)("h6",{id:"example-request"},"Example Request:"),Object(c.b)("p",null,"URL:"),Object(c.b)("pre",null,Object(c.b)("code",Object(a.a)({parentName:"pre"},{}),"/api/v2/issues/Critica/endpoints?date=2015-05-01\n")),Object(c.b)("p",null,"Headers:"),Object(c.b)("pre",null,Object(c.b)("code",Object(a.a)({parentName:"pre"},{}),"x-api-key: shared_key_value\nAccept: application/json or application/xml\n\n")),Object(c.b)("h6",{id:"example-response"},"Example Response:"),Object(c.b)("p",null,"Code:"),Object(c.b)("pre",null,Object(c.b)("code",Object(a.a)({parentName:"pre"},{}),"Status: 200 OK\n")),Object(c.b)("p",null,"Reponse body:"),Object(c.b)("pre",null,Object(c.b)("code",Object(a.a)({parentName:"pre"},{}),'{\n "status": {\n  "message": "Success",\n  "code": "200"\n },\n "data": [\n  {\n   "timestamp": "2015-05-01T05:00:00Z",\n   "endpoint_group": "SITE-A",\n   "service": "web_portal",\n   "endpoint": web01.example.gr",\n   "status": "WARNING",\n   "info": {\n    "Url": "http://example.foo/path/to/service"\n   }\n  },\n  {\n   "timestamp": "2015-05-01T06:00:00Z",\n   "endpoint_group": "SITE-B",\n   "service": "object-storage",\n   "endpoint": "obj.storage.example.gr",\n   "status": "CRITICAL"\n  }\n ]\n}\n')),Object(c.b)("h6",{id:"example-request-with-property-filtercritical"},"Example Request with property filter=CRITICAL:"),Object(c.b)("p",null,"URL:"),Object(c.b)("pre",null,Object(c.b)("code",Object(a.a)({parentName:"pre"},{}),"/api/v2/issues/Critica/endpoints?date=2015-05-01&filter=CRITICAL\n")),Object(c.b)("p",null,"Headers:"),Object(c.b)("pre",null,Object(c.b)("code",Object(a.a)({parentName:"pre"},{}),"x-api-key: shared_key_value\nAccept: application/json or application/xml\n\n")),Object(c.b)("h6",{id:"example-response-1"},"Example Response:"),Object(c.b)("p",null,"Code:"),Object(c.b)("pre",null,Object(c.b)("code",Object(a.a)({parentName:"pre"},{}),"Status: 200 OK\n")),Object(c.b)("p",null,"Reponse body:"),Object(c.b)("pre",null,Object(c.b)("code",Object(a.a)({parentName:"pre"},{}),'{\n "status": {\n  "message": "Success",\n  "code": "200"\n },\n "data": [\n  {\n   "timestamp": "2015-05-01T06:00:00Z",\n   "endpoint_group": "SITE-B",\n   "service": "object-storage",\n   "endpoint": "obj.storage.example.gr",\n   "status": "CRITICAL"\n  }\n ]\n}\n')))}o.isMDXComponent=!0},92:function(e,t,n){"use strict";n.d(t,"a",(function(){return s})),n.d(t,"b",(function(){return O}));var a=n(0),r=n.n(a);function c(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function b(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?b(Object(n),!0).forEach((function(t){c(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):b(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function p(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},c=Object.keys(e);for(a=0;a<c.length;a++)n=c[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var c=Object.getOwnPropertySymbols(e);for(a=0;a<c.length;a++)n=c[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var l=r.a.createContext({}),o=function(e){var t=r.a.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},s=function(e){var t=o(e.components);return r.a.createElement(l.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return r.a.createElement(r.a.Fragment,{},t)}},d=r.a.forwardRef((function(e,t){var n=e.components,a=e.mdxType,c=e.originalType,b=e.parentName,l=p(e,["components","mdxType","originalType","parentName"]),s=o(n),d=a,O=s["".concat(b,".").concat(d)]||s[d]||u[d]||c;return n?r.a.createElement(O,i(i({ref:t},l),{},{components:n})):r.a.createElement(O,i({ref:t},l))}));function O(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var c=n.length,b=new Array(c);b[0]=d;var i={};for(var p in t)hasOwnProperty.call(t,p)&&(i[p]=t[p]);i.originalType=e,i.mdxType="string"==typeof e?e:a,b[1]=i;for(var l=2;l<c;l++)b[l]=n[l];return r.a.createElement.apply(null,b)}return r.a.createElement.apply(null,n)}d.displayName="MDXCreateElement"}}]);