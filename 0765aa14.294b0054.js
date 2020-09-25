(window.webpackJsonp=window.webpackJsonp||[]).push([[5],{54:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return p})),n.d(t,"metadata",(function(){return c})),n.d(t,"rightToc",(function(){return i})),n.d(t,"default",(function(){return l}));var r=n(2),a=n(6),o=(n(0),n(86)),p={id:"feeds",title:"Feeds"},c={unversionedId:"feeds",id:"feeds",isDocsHomePage:!1,title:"Feeds",description:"API Calls",source:"@site/docs/feeds.md",slug:"/feeds",permalink:"/argo-web-api/docs/feeds",version:"current",sidebar:"someSidebar",previous:{title:"Tenants",permalink:"/argo-web-api/docs/"},next:{title:"Downtimes",permalink:"/argo-web-api/docs/downtimes"}},i=[{value:"API Calls",id:"api-calls",children:[]},{value:"GET: List Feed topology parameters",id:"get-list-feed-topology-parameters",children:[{value:"Input",id:"input",children:[]},{value:"Request headers",id:"request-headers",children:[]},{value:"Response",id:"response",children:[]}]},{value:"PUT: Update topology feed parameters",id:"put-update-topology-feed-parameters",children:[{value:"Input",id:"input-1",children:[]},{value:"Response",id:"response-1",children:[]}]}],b={rightToc:i};function l(e){var t=e.components,n=Object(a.a)(e,["components"]);return Object(o.b)("wrapper",Object(r.a)({},b,n,{components:t,mdxType:"MDXLayout"}),Object(o.b)("h2",{id:"api-calls"},"API Calls"),Object(o.b)("table",null,Object(o.b)("thead",{parentName:"table"},Object(o.b)("tr",{parentName:"thead"},Object(o.b)("th",Object(r.a)({parentName:"tr"},{align:null}),"Name"),Object(o.b)("th",Object(r.a)({parentName:"tr"},{align:null}),"Description"),Object(o.b)("th",Object(r.a)({parentName:"tr"},{align:null}),"Shortcut"))),Object(o.b)("tbody",{parentName:"table"},Object(o.b)("tr",{parentName:"tbody"},Object(o.b)("td",Object(r.a)({parentName:"tr"},{align:null}),"GET: Feed Topology information"),Object(o.b)("td",Object(r.a)({parentName:"tr"},{align:null}),"This method can be used to retrieve a list of feed topology parameters"),Object(o.b)("td",Object(r.a)({parentName:"tr"},{align:null}),Object(o.b)("a",Object(r.a)({parentName:"td"},{href:"#1"})," Description"))),Object(o.b)("tr",{parentName:"tbody"},Object(o.b)("td",Object(r.a)({parentName:"tr"},{align:null}),"PUT: Update feed topology info"),Object(o.b)("td",Object(r.a)({parentName:"tr"},{align:null}),"This method can be used to update feed topology information parameters"),Object(o.b)("td",Object(r.a)({parentName:"tr"},{align:null}),Object(o.b)("a",Object(r.a)({parentName:"td"},{href:"#2"})," Description"))))),Object(o.b)("a",{id:"1"}),Object(o.b)("h2",{id:"get-list-feed-topology-parameters"},"[GET]",": List Feed topology parameters"),Object(o.b)("p",null,"This method can be used to retrieve a list of feed topology parameters"),Object(o.b)("h3",{id:"input"},"Input"),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{}),"GET /feeds/topology\n")),Object(o.b)("h3",{id:"request-headers"},"Request headers"),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{}),"x-api-key: shared_key_value\nAccept: application/json\n")),Object(o.b)("h3",{id:"response"},"Response"),Object(o.b)("p",null,"Headers: ",Object(o.b)("inlineCode",{parentName:"p"},"Status: 200 OK")),Object(o.b)("h4",{id:"response-body"},"Response body"),Object(o.b)("p",null,"Json Response"),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{className:"language-json"}),'{\n "status": {\n  "message": "Success",\n  "code": "200"\n },\n "data": [\n  {\n   "type": "gocdb",\n   "feed_url": "https://somewhere.foo.bar/topology/feed",\n   "paginated": "true",\n   "fetch_type": [\n    "item1",\n    "item2"\n   ],\n   "uid_endpoints": "endpointA"\n  }\n ]\n}\n')),Object(o.b)("a",{id:"2"}),Object(o.b)("h2",{id:"put-update-topology-feed-parameters"},"[PUT]",": Update topology feed parameters"),Object(o.b)("p",null,"This method is used to upadte topology feed parameters"),Object(o.b)("h3",{id:"input-1"},"Input"),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{}),"PUT /feeds/topology\n")),Object(o.b)("h4",{id:"put-body"},"PUT BODY"),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{className:"language-json"}),'  {\n   "type": "gocdb",\n   "feed_url": "https://somewhere.foo.bar/topology/feed",\n   "paginated": "true",\n   "fetch_type": [\n    "item1",\n    "item2"\n   ],\n   "uid_endpoints": "endpointA"\n  }\n')),Object(o.b)("h4",{id:"request-headers-1"},"Request headers"),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{}),"x-api-key: shared_key_value\nAccept: application/json\n")),Object(o.b)("h3",{id:"response-1"},"Response"),Object(o.b)("p",null,"Headers: ",Object(o.b)("inlineCode",{parentName:"p"},"Status: 200 OK")),Object(o.b)("h4",{id:"response-body-1"},"Response body"),Object(o.b)("p",null,"Json Response"),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{className:"language-json"}),'{\n "status": {\n  "message": "Feeds resource succesfully updated",\n  "code": "200"\n },\n "data": [\n  {\n   "type": "gocdb",\n   "feed_url": "https://somewhere2.foo.bar/topology/feed",\n   "paginated": "true",\n   "fetch_type": [\n    "item4",\n    "item5"\n   ],\n   "uid_endpoints": "endpointA"\n  }\n ]\n}\n')))}l.isMDXComponent=!0},86:function(e,t,n){"use strict";n.d(t,"a",(function(){return s})),n.d(t,"b",(function(){return O}));var r=n(0),a=n.n(r);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function p(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function c(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?p(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):p(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function i(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var b=a.a.createContext({}),l=function(e){var t=a.a.useContext(b),n=t;return e&&(n="function"==typeof e?e(t):c(c({},t),e)),n},s=function(e){var t=l(e.components);return a.a.createElement(b.Provider,{value:t},e.children)},d={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},u=a.a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,o=e.originalType,p=e.parentName,b=i(e,["components","mdxType","originalType","parentName"]),s=l(n),u=r,O=s["".concat(p,".").concat(u)]||s[u]||d[u]||o;return n?a.a.createElement(O,c(c({ref:t},b),{},{components:n})):a.a.createElement(O,c({ref:t},b))}));function O(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var o=n.length,p=new Array(o);p[0]=u;var c={};for(var i in t)hasOwnProperty.call(t,i)&&(c[i]=t[i]);c.originalType=e,c.mdxType="string"==typeof e?e:r,p[1]=c;for(var b=2;b<o;b++)p[b]=n[b];return a.a.createElement.apply(null,p)}return a.a.createElement.apply(null,n)}u.displayName="MDXCreateElement"}}]);