"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[637],{3905:(e,t,n)=>{n.d(t,{Zo:()=>c,kt:()=>f});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function s(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function i(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var l=r.createContext({}),p=function(e){var t=r.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):s(s({},t),e)),n},c=function(e){var t=p(e.components);return r.createElement(l.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},d=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,o=e.originalType,l=e.parentName,c=i(e,["components","mdxType","originalType","parentName"]),d=p(n),f=a,m=d["".concat(l,".").concat(f)]||d[f]||u[f]||o;return n?r.createElement(m,s(s({ref:t},c),{},{components:n})):r.createElement(m,s({ref:t},c))}));function f(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var o=n.length,s=new Array(o);s[0]=d;var i={};for(var l in t)hasOwnProperty.call(t,l)&&(i[l]=t[l]);i.originalType=e,i.mdxType="string"==typeof e?e:a,s[1]=i;for(var p=2;p<o;p++)s[p]=n[p];return r.createElement.apply(null,s)}return r.createElement.apply(null,n)}d.displayName="MDXCreateElement"},7883:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>l,contentTitle:()=>s,default:()=>u,frontMatter:()=>o,metadata:()=>i,toc:()=>p});var r=n(7462),a=(n(7294),n(3905));const o={id:"factors",title:"Factors",sidebar_position:7},s=void 0,i={unversionedId:"tenants_and_feeds/factors",id:"tenants_and_feeds/factors",title:"Factors",description:"API Calls",source:"@site/docs/tenants_and_feeds/factors.md",sourceDirName:"tenants_and_feeds",slug:"/tenants_and_feeds/factors",permalink:"/argo-web-api/docs/tenants_and_feeds/factors",draft:!1,tags:[],version:"current",sidebarPosition:7,frontMatter:{id:"factors",title:"Factors",sidebar_position:7},sidebar:"tutorialSidebar",previous:{title:"Weights",permalink:"/argo-web-api/docs/tenants_and_feeds/weights"},next:{title:"Topology",permalink:"/argo-web-api/docs/category/topology"}},l={},p=[{value:"API Calls",id:"api-calls",level:2},{value:"Input",id:"input",level:2},{value:"Request headers",id:"request-headers",level:3},{value:"Response",id:"response",level:2},{value:"Response body",id:"response-body",level:3}],c={toc:p};function u(e){let{components:t,...n}=e;return(0,a.kt)("wrapper",(0,r.Z)({},c,n,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h2",{id:"api-calls"},"API Calls"),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:null},"Name"),(0,a.kt)("th",{parentName:"tr",align:null},"Description"),(0,a.kt)("th",{parentName:"tr",align:null},"Shortcut"))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},"GET: List Factors Requests"),(0,a.kt)("td",{parentName:"tr",align:null},"This method can be used to retrieve a list of factors."),(0,a.kt)("td",{parentName:"tr",align:null},(0,a.kt)("a",{parentName:"td",href:"#1"}," Description"))))),(0,a.kt)("h1",{id:"get-list-factors"},"GET: List Factors"),(0,a.kt)("p",null,"This method can be used to retrieve a list of current Factors"),(0,a.kt)("h2",{id:"input"},"Input"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"GET /factors\n")),(0,a.kt)("h3",{id:"request-headers"},"Request headers"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"x-api-key: shared_key_value\nAccept: application/json or application/xml\n")),(0,a.kt)("h2",{id:"response"},"Response"),(0,a.kt)("p",null,"Headers: ",(0,a.kt)("inlineCode",{parentName:"p"},"Status: 200 OK")),(0,a.kt)("h3",{id:"response-body"},"Response body"),(0,a.kt)("p",null,"Json Response"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-json"},'{\n "factors": [\n  {\n   "site": "CETA-GRID",\n   "weight": "5406"\n  },\n  {\n   "site": "CFP-IST",\n   "weight": "1019"\n  },\n  {\n   "site": "CIEMAT-LCG2",\n   "weight": "14595"\n  }\n ]\n}\n')),(0,a.kt)("p",null,"XML Response"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-xml"},'<root>\n    <Factor site="CETA-GRID" weight="5406"></Factor>\n    <Factor site="CFP-IST" weight="1019"></Factor>\n    <Factor site="CIEMAT-LCG2" weight="14595"></Factor>\n</root>\n')))}u.isMDXComponent=!0}}]);