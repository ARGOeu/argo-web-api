(window.webpackJsonp=window.webpackJsonp||[]).push([[10],{60:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return i})),n.d(t,"metadata",(function(){return l})),n.d(t,"rightToc",(function(){return s})),n.d(t,"default",(function(){return b}));var r=n(2),a=n(6),o=(n(0),n(93)),c=["components"],i={id:"factors",title:"Factors"},l={unversionedId:"factors",id:"factors",isDocsHomePage:!1,title:"Factors",description:"API Calls",source:"@site/docs/factors.md",slug:"/factors",permalink:"/argo-web-api/docs/factors",version:"current"},s=[{value:"API Calls",id:"api-calls",children:[]},{value:"Input",id:"input",children:[{value:"Request headers",id:"request-headers",children:[]}]},{value:"Response",id:"response",children:[{value:"Response body",id:"response-body",children:[]}]}],p={rightToc:s};function b(e){var t=e.components,n=Object(a.a)(e,c);return Object(o.b)("wrapper",Object(r.a)({},p,n,{components:t,mdxType:"MDXLayout"}),Object(o.b)("h2",{id:"api-calls"},"API Calls"),Object(o.b)("table",null,Object(o.b)("thead",{parentName:"table"},Object(o.b)("tr",{parentName:"thead"},Object(o.b)("th",{parentName:"tr",align:null},"Name"),Object(o.b)("th",{parentName:"tr",align:null},"Description"),Object(o.b)("th",{parentName:"tr",align:null},"Shortcut"))),Object(o.b)("tbody",{parentName:"table"},Object(o.b)("tr",{parentName:"tbody"},Object(o.b)("td",{parentName:"tr",align:null},"GET: List Factors Requests"),Object(o.b)("td",{parentName:"tr",align:null},"This method can be used to retrieve a list of factors."),Object(o.b)("td",{parentName:"tr",align:null},Object(o.b)("a",{parentName:"td",href:"#1"}," Description"))))),Object(o.b)("h1",{id:"get-list-factors"},"GET: List Factors"),Object(o.b)("p",null,"This method can be used to retrieve a list of current Factors"),Object(o.b)("h2",{id:"input"},"Input"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre"},"GET /factors\n")),Object(o.b)("h3",{id:"request-headers"},"Request headers"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json or application/xml\n")),Object(o.b)("h2",{id:"response"},"Response"),Object(o.b)("p",null,"Headers: ",Object(o.b)("inlineCode",{parentName:"p"},"Status: 200 OK")),Object(o.b)("h3",{id:"response-body"},"Response body"),Object(o.b)("p",null,"Json Response"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-json"},'{\n "factors": [\n  {\n   "site": "CETA-GRID",\n   "weight": "5406"\n  },\n  {\n   "site": "CFP-IST",\n   "weight": "1019"\n  },\n  {\n   "site": "CIEMAT-LCG2",\n   "weight": "14595"\n  }\n ]\n}\n')),Object(o.b)("p",null,"XML Response"),Object(o.b)("pre",null,Object(o.b)("code",{parentName:"pre",className:"language-xml"},'<root>\n    <Factor site="CETA-GRID" weight="5406"></Factor>\n    <Factor site="CFP-IST" weight="1019"></Factor>\n    <Factor site="CIEMAT-LCG2" weight="14595"></Factor>\n</root>\n')))}b.isMDXComponent=!0},93:function(e,t,n){"use strict";n.d(t,"a",(function(){return b})),n.d(t,"b",(function(){return f}));var r=n(0),a=n.n(r);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function c(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?c(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):c(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var s=a.a.createContext({}),p=function(e){var t=a.a.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},b=function(e){var t=p(e.components);return a.a.createElement(s.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},d=a.a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,o=e.originalType,c=e.parentName,s=l(e,["components","mdxType","originalType","parentName"]),b=p(n),d=r,f=b["".concat(c,".").concat(d)]||b[d]||u[d]||o;return n?a.a.createElement(f,i(i({ref:t},s),{},{components:n})):a.a.createElement(f,i({ref:t},s))}));function f(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var o=n.length,c=new Array(o);c[0]=d;var i={};for(var l in t)hasOwnProperty.call(t,l)&&(i[l]=t[l]);i.originalType=e,i.mdxType="string"==typeof e?e:r,c[1]=i;for(var s=2;s<o;s++)c[s]=n[s];return a.a.createElement.apply(null,c)}return a.a.createElement.apply(null,n)}d.displayName="MDXCreateElement"}}]);