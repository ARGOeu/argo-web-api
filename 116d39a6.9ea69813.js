(window.webpackJsonp=window.webpackJsonp||[]).push([[8],{57:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return i})),n.d(t,"metadata",(function(){return p})),n.d(t,"rightToc",(function(){return b})),n.d(t,"default",(function(){return o}));var a=n(2),r=n(6),l=(n(0),n(93)),c=["components"],i={id:"metric_results",title:"Metric Results"},p={unversionedId:"metric_results",id:"metric_results",isDocsHomePage:!1,title:"Metric Results",description:"API call for retrieving detailed metric result.",source:"@site/docs/metric_results.md",slug:"/metric_results",permalink:"/argo-web-api/docs/metric_results",version:"current",sidebar:"someSidebar",previous:{title:"trends",permalink:"/argo-web-api/docs/trends"},next:{title:"Recomputation Requests",permalink:"/argo-web-api/docs/recomputations"}},b=[{value:"API call for retrieving detailed metric result.",id:"api-call-for-retrieving-detailed-metric-result",children:[{value:"GET: Metric Result",id:"get-metric-result",children:[]},{value:"Input",id:"input",children:[]},{value:"Response body",id:"response-body",children:[]}]},{value:"GET: Multiple Metric Results for a specific host, on a specific day",id:"get-multiple-metric-results-for-a-specific-host-on-a-specific-day",children:[{value:"Input",id:"input-1",children:[]},{value:"Response body",id:"response-body-1",children:[]},{value:"Extra endpoint information on metric results",id:"extra-endpoint-information-on-metric-results",children:[]}]}],m={rightToc:b};function o(e){var t=e.components,n=Object(r.a)(e,c);return Object(l.b)("wrapper",Object(a.a)({},m,n,{components:t,mdxType:"MDXLayout"}),Object(l.b)("h2",{id:"api-call-for-retrieving-detailed-metric-result"},"API call for retrieving detailed metric result."),Object(l.b)("h3",{id:"get-metric-result"},"[GET]",": Metric Result"),Object(l.b)("p",null,"This method may be used to retrieve a specific service metric result."),Object(l.b)("h3",{id:"input"},"Input"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"/metric_result/{hostname}/{metric_name}?[exec_time]\n")),Object(l.b)("h4",{id:"path-parameters"},"Path Parameters"),Object(l.b)("table",null,Object(l.b)("thead",{parentName:"table"},Object(l.b)("tr",{parentName:"thead"},Object(l.b)("th",{parentName:"tr",align:null},"Type"),Object(l.b)("th",{parentName:"tr",align:null},"Description"),Object(l.b)("th",{parentName:"tr",align:null},"Required"),Object(l.b)("th",{parentName:"tr",align:null},"Default value"))),Object(l.b)("tbody",{parentName:"table"},Object(l.b)("tr",{parentName:"tbody"},Object(l.b)("td",{parentName:"tr",align:null},Object(l.b)("inlineCode",{parentName:"td"},"hostname")),Object(l.b)("td",{parentName:"tr",align:null},"hostname of service endpoint"),Object(l.b)("td",{parentName:"tr",align:null},"YES"),Object(l.b)("td",{parentName:"tr",align:null})),Object(l.b)("tr",{parentName:"tbody"},Object(l.b)("td",{parentName:"tr",align:null},Object(l.b)("inlineCode",{parentName:"td"},"metric_name")),Object(l.b)("td",{parentName:"tr",align:null},"name of the metric"),Object(l.b)("td",{parentName:"tr",align:null},"YES"),Object(l.b)("td",{parentName:"tr",align:null})))),Object(l.b)("h4",{id:"url-parameters"},"Url Parameters"),Object(l.b)("table",null,Object(l.b)("thead",{parentName:"table"},Object(l.b)("tr",{parentName:"thead"},Object(l.b)("th",{parentName:"tr",align:null},"Type"),Object(l.b)("th",{parentName:"tr",align:null},"Description"),Object(l.b)("th",{parentName:"tr",align:null},"Required"),Object(l.b)("th",{parentName:"tr",align:null},"Default value"))),Object(l.b)("tbody",{parentName:"table"},Object(l.b)("tr",{parentName:"tbody"},Object(l.b)("td",{parentName:"tr",align:null},Object(l.b)("inlineCode",{parentName:"td"},"exec_time")),Object(l.b)("td",{parentName:"tr",align:null},"The execution date of query in zulu format"),Object(l.b)("td",{parentName:"tr",align:null},"YES"),Object(l.b)("td",{parentName:"tr",align:null})))),Object(l.b)("p",null,Object(l.b)("strong",{parentName:"p"},Object(l.b)("em",{parentName:"strong"},"Notes")),":\n",Object(l.b)("inlineCode",{parentName:"p"},"exec_time")," : The execution date of query in zulu format. In order to get the correct execution time get status results for all metrics (under a given endpoint, service and endpoint group). ( GET /status/{report_name}/{lgroup_type}/{lgroup_name}/services/{service_name}/endpoints/{endpoint_name}/metrics List)"),Object(l.b)("h4",{id:"headers"},"Headers"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json or application/xml\n")),Object(l.b)("h4",{id:"response-code"},"Response Code"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"Status: 200 OK\n")),Object(l.b)("h3",{id:"response-body"},"Response body"),Object(l.b)("h6",{id:"example-request"},"Example Request:"),Object(l.b)("p",null,"URL:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"/api/v2/metric_result/www.example.com/httpd_check?exec_time=2015-06-20T12:00:00Z\n")),Object(l.b)("p",null,"Headers:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json or application/xml\n\n")),Object(l.b)("h6",{id:"example-response"},"Example Response:"),Object(l.b)("p",null,"Code:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"Status: 200 OK\n")),Object(l.b)("p",null,"Reponse body:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},' {\n   "root": [\n     {\n       "Name": "www.example.com",\n       "Metrics": [\n         {\n           "Name": "httpd_check",\n           "Service": "httpd",\n           "Details": [\n             {\n               "Timestamp": "2015-06-20T12:00:00Z",\n               "Value": "OK",\n               "Summary": "httpd is ok",\n               "Message": "all checks ok"\n             }\n           ]\n         }\n       ]\n     }\n   ]\n }\n \n')),Object(l.b)("h2",{id:"get-multiple-metric-results-for-a-specific-host-on-a-specific-day"},"[GET]",": Multiple Metric Results for a specific host, on a specific day"),Object(l.b)("p",null,"This method may be used to retrieve multiple service metric results for a specific host on a specific day"),Object(l.b)("h3",{id:"input-1"},"Input"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"/metric_result/{hostname}?[exec_time]\n")),Object(l.b)("h4",{id:"path-parameters-1"},"Path Parameters"),Object(l.b)("table",null,Object(l.b)("thead",{parentName:"table"},Object(l.b)("tr",{parentName:"thead"},Object(l.b)("th",{parentName:"tr",align:null},"Type"),Object(l.b)("th",{parentName:"tr",align:null},"Description"),Object(l.b)("th",{parentName:"tr",align:null},"Required"),Object(l.b)("th",{parentName:"tr",align:null},"Default value"))),Object(l.b)("tbody",{parentName:"table"},Object(l.b)("tr",{parentName:"tbody"},Object(l.b)("td",{parentName:"tr",align:null},Object(l.b)("inlineCode",{parentName:"td"},"hostname")),Object(l.b)("td",{parentName:"tr",align:null},"hostname of service endpoint"),Object(l.b)("td",{parentName:"tr",align:null},"YES"),Object(l.b)("td",{parentName:"tr",align:null})))),Object(l.b)("h4",{id:"url-parameters-1"},"Url Parameters"),Object(l.b)("table",null,Object(l.b)("thead",{parentName:"table"},Object(l.b)("tr",{parentName:"thead"},Object(l.b)("th",{parentName:"tr",align:null},"Type"),Object(l.b)("th",{parentName:"tr",align:null},"Description"),Object(l.b)("th",{parentName:"tr",align:null},"Required"),Object(l.b)("th",{parentName:"tr",align:null},"Default value"))),Object(l.b)("tbody",{parentName:"table"},Object(l.b)("tr",{parentName:"tbody"},Object(l.b)("td",{parentName:"tr",align:null},Object(l.b)("inlineCode",{parentName:"td"},"exec_time")),Object(l.b)("td",{parentName:"tr",align:null},"The execution date of query in zulu format - timepart is irrelevant (can be 00:00:00Z)"),Object(l.b)("td",{parentName:"tr",align:null},"YES"),Object(l.b)("td",{parentName:"tr",align:null})),Object(l.b)("tr",{parentName:"tbody"},Object(l.b)("td",{parentName:"tr",align:null},Object(l.b)("inlineCode",{parentName:"td"},"filter")),Object(l.b)("td",{parentName:"tr",align:null},"Filter metric results by statuses: non-ok, ok, critical, warning"),Object(l.b)("td",{parentName:"tr",align:null},"NO"),Object(l.b)("td",{parentName:"tr",align:null})))),Object(l.b)("p",null,Object(l.b)("strong",{parentName:"p"},Object(l.b)("em",{parentName:"strong"},"Notes")),":\n",Object(l.b)("inlineCode",{parentName:"p"},"exec_time")," : The specific date of query in zulu format. The time part of the date is irrelevant because all metrics of that day are returned. ( GET /status/{report_name}/{lgroup_type}/{lgroup_name}/services/{service_name}/endpoints/{endpoint_name}/metrics List)"),Object(l.b)("h4",{id:"headers-1"},"Headers"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json or application/xml\n")),Object(l.b)("h4",{id:"response-code-1"},"Response Code"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"Status: 200 OK\n")),Object(l.b)("h3",{id:"response-body-1"},"Response body"),Object(l.b)("h6",{id:"example-request-1"},"Example Request:"),Object(l.b)("p",null,"URL:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"/api/v2/metric_result/www.example.com?exec_time=2015-06-20T00:00:00Z\n")),Object(l.b)("p",null,"Headers:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json or application/xml\n\n")),Object(l.b)("h6",{id:"example-response-1"},"Example Response:"),Object(l.b)("p",null,"Code:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"Status: 200 OK\n")),Object(l.b)("p",null,"Reponse body:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},'{\n   "root": [\n     {\n       "Name": "www.example.com",\n       "Metrics": [\n         {\n           "Name": "httpd_check",\n           "Service": "httpd",\n           "Details": [\n             {\n               "Timestamp": "2015-06-20T12:00:00Z",\n               "Value": "OK",\n               "Summary": "httpd is ok",\n               "Message": "all checks ok"\n             },\n              {\n               "Timestamp": "2015-06-20T18:00:00Z",\n               "Value": "CRITICAL",\n               "Summary": "httpd is critical",\n               "Message": "some checks failed"\n             },\n              {\n               "Timestamp": "2015-06-20T23:00:00Z",\n               "Value": "OK",\n               "Summary": "httpd is ok",\n               "Message": "all checks ok"\n             }\n           ]\n         },\n         {\n           "Name": "httpd_memory",\n           "Service": "httpd",\n           "Details": [\n             {\n               "Timestamp": "2015-06-20T06:00:00Z",\n               "Value": "OK",\n               "Summary": "memcheck ok",\n               "Message": "memory under 20%"\n             },\n             {\n               "Timestamp": "2015-06-20T09:00:00Z",\n               "Value": "OK",\n               "Summary": "memcheck ok",\n               "Message": "memory under 20%"\n             },\n             {\n               "Timestamp": "2015-06-20T18:00:00Z",\n               "Value": "OK",\n               "Summary": "memcheck ok",\n               "Message": "memory under 20%"\n             },\n           ]\n         }\n       ]\n     }\n   ]\n }\n')),Object(l.b)("h6",{id:"example-request-with-filter-parameter-set-to-non-ok"},"Example Request with filter parameter set to ",Object(l.b)("inlineCode",{parentName:"h6"},"non-ok"),":"),Object(l.b)("p",null,"URL:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"/api/v2/metric_result/www.example.com?exec_time=2015-06-20T00:00:00Z&filter=non-ok\n")),Object(l.b)("p",null,"Headers:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json or application/xml\n\n")),Object(l.b)("h6",{id:"example-response-using-fitler-parameter-set-to-non-ok"},"Example Response using fitler parameter set to ",Object(l.b)("inlineCode",{parentName:"h6"},"non-ok"),":"),Object(l.b)("p",null,"Code:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"Status: 200 OK\n")),Object(l.b)("p",null,"Reponse body:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},'{\n   "root": [\n     {\n       "Name": "www.example.com",\n       "Metrics": [\n         {\n           "Name": "httpd_check",\n           "Service": "httpd",\n           "Details": [\n              {\n               "Timestamp": "2015-06-20T18:00:00Z",\n               "Value": "CRITICAL",\n               "Summary": "httpd is critical",\n               "Message": "some checks failed"\n              }\n           ]\n         }\n       ]\n     }\n   ]\n }\n')),Object(l.b)("h6",{id:"example-request-with-filter-parameter-set-to-ok"},"Example Request with filter parameter set to ",Object(l.b)("inlineCode",{parentName:"h6"},"ok"),":"),Object(l.b)("p",null,"URL:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"/api/v2/metric_result/www.example.com?exec_time=2015-06-20T00:00:00Z&filter=ok\n")),Object(l.b)("p",null,"Headers:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json or application/xml\n\n")),Object(l.b)("h6",{id:"example-response-using-fitler-parameter-set-to-ok"},"Example Response using fitler parameter set to ",Object(l.b)("inlineCode",{parentName:"h6"},"ok"),":"),Object(l.b)("p",null,"Code:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},"Status: 200 OK\n")),Object(l.b)("p",null,"Reponse body:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},'{\n   "root": [\n     {\n       "Name": "www.example.com",\n       "Metrics": [\n         {\n           "Name": "httpd_check",\n           "Service": "httpd",\n           "Details": [\n             {\n               "Timestamp": "2015-06-20T12:00:00Z",\n               "Value": "OK",\n               "Summary": "httpd is ok",\n               "Message": "all checks ok"\n             },\n             {\n               "Timestamp": "2015-06-20T23:00:00Z",\n               "Value": "OK",\n               "Summary": "httpd is ok",\n               "Message": "all checks ok"\n             }\n           ]\n         },\n         {\n           "Name": "httpd_memory",\n           "Service": "httpd",\n           "Details": [\n             {\n               "Timestamp": "2015-06-20T06:00:00Z",\n               "Value": "OK",\n               "Summary": "memcheck ok",\n               "Message": "memory under 20%"\n             },\n             {\n               "Timestamp": "2015-06-20T09:00:00Z",\n               "Value": "OK",\n               "Summary": "memcheck ok",\n               "Message": "memory under 20%"\n             },\n             {\n               "Timestamp": "2015-06-20T18:00:00Z",\n               "Value": "OK",\n               "Summary": "memcheck ok",\n               "Message": "memory under 20%"\n             },\n           ]\n         }\n       ]\n     }\n   ]\n }\n')),Object(l.b)("h3",{id:"extra-endpoint-information-on-metric-results"},"Extra endpoint information on metric results"),Object(l.b)("p",null,"Some metric results have additional information regarding the specific service endpoint such as it's Url, certificat DN etc... If this information is available it will be displayed under each service endpoint along with status results. For example:"),Object(l.b)("pre",null,Object(l.b)("code",{parentName:"pre"},'{\n   "root": [\n     {\n       "Name": "www.example.com",\n       "info": {\n                  "Url": "https://example.com/path/to/service/check"\n               },\n       "Metrics": [\n         {\n           "Name": "httpd_check",\n           "Service": "httpd",\n           "Details": [\n             {\n               "Timestamp": "2015-06-20T12:00:00Z",\n               "Value": "OK",\n               "Summary": "httpd is ok",\n               "Message": "all checks ok"\n             },\n             {\n               "Timestamp": "2015-06-20T23:00:00Z",\n               "Value": "OK",\n               "Summary": "httpd is ok",\n               "Message": "all checks ok"\n             }\n           ]\n         },\n         {\n           "Name": "httpd_memory",\n           "Service": "httpd",\n           "Details": [\n             {\n               "Timestamp": "2015-06-20T06:00:00Z",\n               "Value": "OK",\n               "Summary": "memcheck ok",\n               "Message": "memory under 20%"\n             },\n             {\n               "Timestamp": "2015-06-20T09:00:00Z",\n               "Value": "OK",\n               "Summary": "memcheck ok",\n               "Message": "memory under 20%"\n             },\n             {\n               "Timestamp": "2015-06-20T18:00:00Z",\n               "Value": "OK",\n               "Summary": "memcheck ok",\n               "Message": "memory under 20%"\n             },\n           ]\n         }\n       ]\n     }\n   ]\n }\n')))}o.isMDXComponent=!0},93:function(e,t,n){"use strict";n.d(t,"a",(function(){return o})),n.d(t,"b",(function(){return d}));var a=n(0),r=n.n(a);function l(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function c(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?c(Object(n),!0).forEach((function(t){l(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):c(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function p(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},l=Object.keys(e);for(a=0;a<l.length;a++)n=l[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var l=Object.getOwnPropertySymbols(e);for(a=0;a<l.length;a++)n=l[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var b=r.a.createContext({}),m=function(e){var t=r.a.useContext(b),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},o=function(e){var t=m(e.components);return r.a.createElement(b.Provider,{value:t},e.children)},s={inlineCode:"code",wrapper:function(e){var t=e.children;return r.a.createElement(r.a.Fragment,{},t)}},u=r.a.forwardRef((function(e,t){var n=e.components,a=e.mdxType,l=e.originalType,c=e.parentName,b=p(e,["components","mdxType","originalType","parentName"]),o=m(n),u=a,d=o["".concat(c,".").concat(u)]||o[u]||s[u]||l;return n?r.a.createElement(d,i(i({ref:t},b),{},{components:n})):r.a.createElement(d,i({ref:t},b))}));function d(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var l=n.length,c=new Array(l);c[0]=u;var i={};for(var p in t)hasOwnProperty.call(t,p)&&(i[p]=t[p]);i.originalType=e,i.mdxType="string"==typeof e?e:a,c[1]=i;for(var b=2;b<l;b++)c[b]=n[b];return r.a.createElement.apply(null,c)}return r.a.createElement.apply(null,n)}u.displayName="MDXCreateElement"}}]);