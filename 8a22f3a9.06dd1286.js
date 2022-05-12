(window.webpackJsonp=window.webpackJsonp||[]).push([[22],{78:function(e,t,r){"use strict";r.r(t),r.d(t,"frontMatter",(function(){return l})),r.d(t,"metadata",(function(){return c})),r.d(t,"rightToc",(function(){return b})),r.d(t,"default",(function(){return d}));var n=r(2),a=r(6),o=(r(0),r(93)),i=["components"],l={id:"errors",title:"API Errors"},c={unversionedId:"errors",id:"errors",isDocsHomePage:!1,title:"API Errors",description:"The following error codes exist among the methods of the ARGO Web API:",source:"@site/docs/errors.md",slug:"/errors",permalink:"/argo-web-api/docs/errors",version:"current",sidebar:"someSidebar",previous:{title:"API Validations",permalink:"/argo-web-api/docs/validations"},next:{title:"Availability / Reliability Results (v3)",permalink:"/argo-web-api/docs/v3_ar_results"}},b=[],p={rightToc:b};function d(e){var t=e.components,r=Object(a.a)(e,i);return Object(o.b)("wrapper",Object(n.a)({},p,r,{components:t,mdxType:"MDXLayout"}),Object(o.b)("p",null,"The following error codes exist among the methods of the ARGO Web API:"),Object(o.b)("table",null,Object(o.b)("thead",{parentName:"table"},Object(o.b)("tr",{parentName:"thead"},Object(o.b)("th",{parentName:"tr",align:null},"Error"),Object(o.b)("th",{parentName:"tr",align:null},"HTTP Code"),Object(o.b)("th",{parentName:"tr",align:null},"Description"))),Object(o.b)("tbody",{parentName:"table"},Object(o.b)("tr",{parentName:"tbody"},Object(o.b)("td",{parentName:"tr",align:null},"Bad request"),Object(o.b)("td",{parentName:"tr",align:null},"400"),Object(o.b)("td",{parentName:"tr",align:null},"One or more checks may have failed. More details on the carried out checks can be found ",Object(o.b)("a",{parentName:"td",href:"/argo-web-api/docs/validations"},"here"))),Object(o.b)("tr",{parentName:"tbody"},Object(o.b)("td",{parentName:"tr",align:null},"Wrong start_time"),Object(o.b)("td",{parentName:"tr",align:null},"400"),Object(o.b)("td",{parentName:"tr",align:null},"Use start_time url parameter in zulu format (like ",Object(o.b)("inlineCode",{parentName:"td"},"2006-01-02T15:04:05Z"),") to indicate the query start time")),Object(o.b)("tr",{parentName:"tbody"},Object(o.b)("td",{parentName:"tr",align:null},"Wrong end_time"),Object(o.b)("td",{parentName:"tr",align:null},"400"),Object(o.b)("td",{parentName:"tr",align:null},"Use end_time url parameter in zulu format (like ",Object(o.b)("inlineCode",{parentName:"td"},"2006-01-02T15:04:05Z"),") to indicate the query end time")),Object(o.b)("tr",{parentName:"tbody"},Object(o.b)("td",{parentName:"tr",align:null},"Wrong exec_time"),Object(o.b)("td",{parentName:"tr",align:null},"400"),Object(o.b)("td",{parentName:"tr",align:null},"Use exec_time url parameter in zulu format (like ",Object(o.b)("inlineCode",{parentName:"td"},"2006-01-02T15:04:05Z"),") to indicate the exact probe execution time")),Object(o.b)("tr",{parentName:"tbody"},Object(o.b)("td",{parentName:"tr",align:null},"Wrong granularity"),Object(o.b)("td",{parentName:"tr",align:null},"400"),Object(o.b)("td",{parentName:"tr",align:null},"The parameter value can be either ",Object(o.b)("inlineCode",{parentName:"td"},"daily")," or ",Object(o.b)("inlineCode",{parentName:"td"},"monthly"))),Object(o.b)("tr",{parentName:"tbody"},Object(o.b)("td",{parentName:"tr",align:null},"Unauthorized"),Object(o.b)("td",{parentName:"tr",align:null},"401"),Object(o.b)("td",{parentName:"tr",align:null},"The client needs to provide a correct authentication token using the header ",Object(o.b)("inlineCode",{parentName:"td"},"x-api-key"))),Object(o.b)("tr",{parentName:"tbody"},Object(o.b)("td",{parentName:"tr",align:null},"Forbidden"),Object(o.b)("td",{parentName:"tr",align:null},"403"),Object(o.b)("td",{parentName:"tr",align:null},"Access to the resource is forbidden due to authorization policy enforced")),Object(o.b)("tr",{parentName:"tbody"},Object(o.b)("td",{parentName:"tr",align:null},"Item not found"),Object(o.b)("td",{parentName:"tr",align:null},"404"),Object(o.b)("td",{parentName:"tr",align:null},"Either the path is not found or no results are available for the given query")),Object(o.b)("tr",{parentName:"tbody"},Object(o.b)("td",{parentName:"tr",align:null},"Content not acceptable"),Object(o.b)("td",{parentName:"tr",align:null},"406"),Object(o.b)("td",{parentName:"tr",align:null},"The ",Object(o.b)("inlineCode",{parentName:"td"},"Accept")," header either was not provided or was provided but did not contain any valid content types. Acceptable content types are ",Object(o.b)("inlineCode",{parentName:"td"},"application/xml")," or ",Object(o.b)("inlineCode",{parentName:"td"},"application/json"))))))}d.isMDXComponent=!0},93:function(e,t,r){"use strict";r.d(t,"a",(function(){return d})),r.d(t,"b",(function(){return s}));var n=r(0),a=r.n(n);function o(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function i(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function l(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?i(Object(r),!0).forEach((function(t){o(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):i(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function c(e,t){if(null==e)return{};var r,n,a=function(e,t){if(null==e)return{};var r,n,a={},o=Object.keys(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||(a[r]=e[r]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(a[r]=e[r])}return a}var b=a.a.createContext({}),p=function(e){var t=a.a.useContext(b),r=t;return e&&(r="function"==typeof e?e(t):l(l({},t),e)),r},d=function(e){var t=p(e.components);return a.a.createElement(b.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},m=a.a.forwardRef((function(e,t){var r=e.components,n=e.mdxType,o=e.originalType,i=e.parentName,b=c(e,["components","mdxType","originalType","parentName"]),d=p(r),m=n,s=d["".concat(i,".").concat(m)]||d[m]||u[m]||o;return r?a.a.createElement(s,l(l({ref:t},b),{},{components:r})):a.a.createElement(s,l({ref:t},b))}));function s(e,t){var r=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var o=r.length,i=new Array(o);i[0]=m;var l={};for(var c in t)hasOwnProperty.call(t,c)&&(l[c]=t[c]);l.originalType=e,l.mdxType="string"==typeof e?e:n,i[1]=l;for(var b=2;b<o;b++)i[b]=r[b];return a.a.createElement.apply(null,i)}return a.a.createElement.apply(null,r)}m.displayName="MDXCreateElement"}}]);