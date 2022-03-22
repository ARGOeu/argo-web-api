(window.webpackJsonp=window.webpackJsonp||[]).push([[4],{53:function(e,n,t){"use strict";t.r(n),t.d(n,"frontMatter",(function(){return o})),t.d(n,"metadata",(function(){return b})),t.d(n,"rightToc",(function(){return l})),t.d(n,"default",(function(){return p}));var a=t(2),r=t(6),i=(t(0),t(90)),o={id:"operations_profiles",title:"Operation Profiles"},b={unversionedId:"operations_profiles",id:"operations_profiles",isDocsHomePage:!1,title:"Operation Profiles",description:"API Calls",source:"@site/docs/operations_profiles.md",slug:"/operations_profiles",permalink:"/argo-web-api/docs/operations_profiles",version:"current",sidebar:"someSidebar",previous:{title:"Available Metrics and Tags",permalink:"/argo-web-api/docs/metrics"},next:{title:"Metric Profiles",permalink:"/argo-web-api/docs/metric_profiles"}},l=[{value:"API Calls",id:"api-calls",children:[]},{value:"GET: List Operations Profiles",id:"get-list-operations-profiles",children:[{value:"Input",id:"input",children:[]},{value:"Response",id:"response",children:[]}]},{value:"GET: List A Specific Operations profile",id:"get-list-a-specific-operations-profile",children:[{value:"Input",id:"input-1",children:[]},{value:"Response",id:"response-1",children:[]}]},{value:"POST: Create a new Operations Profile",id:"post-create-a-new-operations-profile",children:[{value:"Input",id:"input-2",children:[]},{value:"Response",id:"response-2",children:[]}]},{value:"PUT: Update information on an existing operations profile",id:"put-update-information-on-an-existing-operations-profile",children:[{value:"Input",id:"input-3",children:[]},{value:"Response",id:"response-3",children:[]}]},{value:"DELETE: Delete an existing aggregation profile",id:"delete-delete-an-existing-aggregation-profile",children:[{value:"Input",id:"input-4",children:[]},{value:"Response",id:"response-4",children:[]}]},{value:"Validation Checks",id:"validation-checks",children:[{value:"Response",id:"response-5",children:[]}]}],s={rightToc:l};function p(e){var n=e.components,t=Object(r.a)(e,["components"]);return Object(i.b)("wrapper",Object(a.a)({},s,t,{components:n,mdxType:"MDXLayout"}),Object(i.b)("h2",{id:"api-calls"},"API Calls"),Object(i.b)("table",null,Object(i.b)("thead",{parentName:"table"},Object(i.b)("tr",{parentName:"thead"},Object(i.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Name"),Object(i.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Description"),Object(i.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Shortcut"))),Object(i.b)("tbody",{parentName:"table"},Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"GET: List Operations Profile Requests"),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"This method can be used to retrieve a list of current Operations profiles."),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(i.b)("a",Object(a.a)({parentName:"td"},{href:"#1"})," Description"))),Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"GET: List a specific  Operations profile"),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"This method can be used to retrieve a specific  Operations profile based on its id."),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(i.b)("a",Object(a.a)({parentName:"td"},{href:"#2"})," Description"))),Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"POST: Create a new  Operations profile"),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"This method can be used to create a new  Operations profile"),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(i.b)("a",Object(a.a)({parentName:"td"},{href:"#3"})," Description"))),Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"PUT: Update an Operations profile"),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"This method can be used to update information on an existing  Operations profile"),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(i.b)("a",Object(a.a)({parentName:"td"},{href:"#4"})," Description"))),Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"DELETE: Delete an  Operations profile"),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"This method can be used to delete an existing  Operations profile"),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(i.b)("a",Object(a.a)({parentName:"td"},{href:"#5"})," Description"))))),Object(i.b)("a",{id:"1"}),Object(i.b)("h2",{id:"get-list-operations-profiles"},"[GET]",": List Operations Profiles"),Object(i.b)("p",null,"This method can be used to retrieve a list of current  Operations profiles. "),Object(i.b)("h3",{id:"input"},"Input"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{}),"GET /operations_profiles\n")),Object(i.b)("h4",{id:"optional-query-parameters"},"Optional Query Parameters"),Object(i.b)("table",null,Object(i.b)("thead",{parentName:"table"},Object(i.b)("tr",{parentName:"thead"},Object(i.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Type"),Object(i.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Description"),Object(i.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Required"))),Object(i.b)("tbody",{parentName:"table"},Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(i.b)("inlineCode",{parentName:"td"},"name")),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"Operations profile name to be used as query"),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"NO")),Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(i.b)("inlineCode",{parentName:"td"},"date")),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"Date to retrieve a historic version of the operation profile. If no date parameter is provided the most current profile will be returned"),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"NO")))),Object(i.b)("h4",{id:"request-headers"},"Request headers"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{}),"x-api-key: shared_key_value\nAccept: application/json\n")),Object(i.b)("h3",{id:"response"},"Response"),Object(i.b)("p",null,"Headers: ",Object(i.b)("inlineCode",{parentName:"p"},"Status: 200 OK")),Object(i.b)("h4",{id:"response-body"},"Response body"),Object(i.b)("p",null,"Json Response"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{className:"language-json"}),'{\n "status": {\n  "message": "Success",\n  "code": "200"\n },\n "data": [\n  {\n   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",\n   "date": "2019-11-04",\n   "name": "ops1",\n   "available_states": [\n    "A,B,C"\n   ],\n   "defaults": {\n    "down": "B",\n    "missing": "A",\n    "unknown": "C"\n   },\n   "operations": [\n    {\n     "name": "AND",\n     "truth_table": [\n      {\n       "a": "A",\n       "b": "B",\n       "x": "B"\n      },\n      {\n       "a": "A",\n       "b": "C",\n       "x": "C"\n      },\n      {\n       "a": "B",\n       "b": "C",\n       "x": "C"\n      }\n     ]\n    },\n    {\n     "name": "OR",\n     "truth_table": [\n      {\n       "a": "A",\n       "b": "B",\n       "x": "A"\n      },\n      {\n       "a": "A",\n       "b": "C",\n       "x": "A"\n      },\n      {\n       "a": "B",\n       "b": "C",\n       "x": "B"\n      }\n     ]\n    }\n   ]\n  },\n  {\n   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",\n   "date": "2019-11-02",\n   "name": "ops2",\n   "available_states": [\n    "X,Y,Z"\n   ],\n   "defaults": {\n    "down": "Y",\n    "missing": "X",\n    "unknown": "Z"\n   },\n   "operations": [\n    {\n     "name": "AND",\n     "truth_table": [\n      {\n       "a": "X",\n       "b": "Y",\n       "x": "Y"\n      },\n      {\n       "a": "X",\n       "b": "Z",\n       "x": "Z"\n      },\n      {\n       "a": "Y",\n       "b": "Z",\n       "x": "Z"\n      }\n     ]\n    },\n    {\n     "name": "OR",\n     "truth_table": [\n      {\n       "a": "X",\n       "b": "Y",\n       "x": "X"\n      },\n      {\n       "a": "X",\n       "b": "Z",\n       "x": "X"\n      },\n      {\n       "a": "Y",\n       "b": "Z",\n       "x": "Y"\n      }\n     ]\n    }\n   ]\n  }\n ]\n}\n')),Object(i.b)("a",{id:"2"}),Object(i.b)("h2",{id:"get-list-a-specific-operations-profile"},"[GET]",": List A Specific Operations profile"),Object(i.b)("p",null,"This method can be used to retrieve specific Operations profile based on its id"),Object(i.b)("h3",{id:"input-1"},"Input"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{}),"GET /operations_profiles/{ID}\n")),Object(i.b)("h4",{id:"optional-query-parameters-1"},"Optional Query Parameters"),Object(i.b)("table",null,Object(i.b)("thead",{parentName:"table"},Object(i.b)("tr",{parentName:"thead"},Object(i.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Type"),Object(i.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Description"),Object(i.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Required"))),Object(i.b)("tbody",{parentName:"table"},Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(i.b)("inlineCode",{parentName:"td"},"date")),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"Date to retrieve a historic version of the operation profile. If no date parameter is provided the most current profile will be returned"),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"NO")))),Object(i.b)("h4",{id:"request-headers-1"},"Request headers"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{}),"x-api-key: shared_key_value\nAccept: application/json\n")),Object(i.b)("h3",{id:"response-1"},"Response"),Object(i.b)("p",null,"Headers: ",Object(i.b)("inlineCode",{parentName:"p"},"Status: 200 OK")),Object(i.b)("h4",{id:"response-body-1"},"Response body"),Object(i.b)("p",null,"Json Response"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{className:"language-json"}),'{\n "status": {\n  "message": "Success",\n  "code": "200"\n },\n "data": [\n  {\n   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",\n   "date": "2019-11-04",\n   "name": "ops1",\n   "available_states": [\n    "A,B,C"\n   ],\n   "defaults": {\n    "down": "B",\n    "missing": "A",\n    "unknown": "C"\n   },\n   "operations": [\n    {\n     "name": "AND",\n     "truth_table": [\n      {\n       "a": "A",\n       "b": "B",\n       "x": "B"\n      },\n      {\n       "a": "A",\n       "b": "C",\n       "x": "C"\n      },\n      {\n       "a": "B",\n       "b": "C",\n       "x": "C"\n      }\n     ]\n    },\n    {\n     "name": "OR",\n     "truth_table": [\n      {\n       "a": "A",\n       "b": "B",\n       "x": "A"\n      },\n      {\n       "a": "A",\n       "b": "C",\n       "x": "A"\n      },\n      {\n       "a": "B",\n       "b": "C",\n       "x": "B"\n      }\n     ]\n    }\n   ]\n  }\n ]\n}\n')),Object(i.b)("a",{id:"3"}),Object(i.b)("h2",{id:"post-create-a-new-operations-profile"},"[POST]",": Create a new Operations Profile"),Object(i.b)("p",null,"This method can be used to insert a new operations profile"),Object(i.b)("h3",{id:"input-2"},"Input"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{}),"POST /operations_profiles\n")),Object(i.b)("h4",{id:"request-headers-2"},"Request headers"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{}),"x-api-key: shared_key_value\nAccept: application/json\n")),Object(i.b)("h4",{id:"optional-query-parameters-2"},"Optional Query Parameters"),Object(i.b)("table",null,Object(i.b)("thead",{parentName:"table"},Object(i.b)("tr",{parentName:"thead"},Object(i.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Type"),Object(i.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Description"),Object(i.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Required"))),Object(i.b)("tbody",{parentName:"table"},Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(i.b)("inlineCode",{parentName:"td"},"date")),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"Date to create a  new historic version of the operation profile. If no date parameter is provided current date will be supplied automatically"),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"NO")))),Object(i.b)("h4",{id:"post-body"},"POST BODY"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{className:"language-json"}),'{\n   "name": "tops1",\n   "available_states": [\n    "A","B","C"\n   ],\n   "defaults": {\n    "down": "B",\n    "missing": "A",\n    "unknown": "C"\n   },\n   "operations": [\n    {\n     "name": "AND",\n     "truth_table": [\n      {\n       "a": "A",\n       "b": "B",\n       "x": "B"\n      },\n      {\n       "a": "A",\n       "b": "C",\n       "x": "C"\n      },\n      {\n       "a": "B",\n       "b": "C",\n       "x": "C"\n      }\n     ]\n    },\n    {\n     "name": "OR",\n     "truth_table": [\n      {\n       "a": "A",\n       "b": "B",\n       "x": "A"\n      },\n      {\n       "a": "A",\n       "b": "C",\n       "x": "A"\n      },\n      {\n       "a": "B",\n       "b": "C",\n       "x": "B"\n      }\n     ]\n    }\n   ]\n  }\n')),Object(i.b)("h3",{id:"response-2"},"Response"),Object(i.b)("p",null,"Headers: ",Object(i.b)("inlineCode",{parentName:"p"},"Status: 201 Created")),Object(i.b)("h4",{id:"response-body-2"},"Response body"),Object(i.b)("p",null,"Json Response"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{className:"language-json"}),'{\n "status": {\n  "message": "Operations Profile successfully created",\n  "code": "201"\n },\n "data": {\n  "id": "{{ID}}",\n  "links": {\n   "self": "https:///api/v2/operations_profiles/{{ID}}"\n  }\n }\n}\n')),Object(i.b)("a",{id:"4"}),Object(i.b)("h2",{id:"put-update-information-on-an-existing-operations-profile"},"[PUT]",": Update information on an existing operations profile"),Object(i.b)("p",null,"This method can be used to update information on an existing operations profile"),Object(i.b)("h3",{id:"input-3"},"Input"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{}),"PUT /operations_profiles/{ID}\n")),Object(i.b)("h4",{id:"optional-query-parameters-3"},"Optional Query Parameters"),Object(i.b)("table",null,Object(i.b)("thead",{parentName:"table"},Object(i.b)("tr",{parentName:"thead"},Object(i.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Type"),Object(i.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Description"),Object(i.b)("th",Object(a.a)({parentName:"tr"},{align:null}),"Required"))),Object(i.b)("tbody",{parentName:"table"},Object(i.b)("tr",{parentName:"tbody"},Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),Object(i.b)("inlineCode",{parentName:"td"},"date")),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"Date to update a historic version of the operation profile. If no date parameter is provided the current date will be supplied automatically"),Object(i.b)("td",Object(a.a)({parentName:"tr"},{align:null}),"NO")))),Object(i.b)("h4",{id:"request-headers-3"},"Request headers"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{}),"x-api-key: shared_key_value\nAccept: application/json\n")),Object(i.b)("h4",{id:"put-body"},"PUT BODY"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{className:"language-json"}),'{\n     "name": "tops1",\n     "available_states": [\n        "A","B","C"\n     ],\n     "defaults": {\n        "down": "B",\n        "missing": "A",\n        "unknown": "C"\n     },\n     "operations": [\n        {\n         "name": "AND",\n         "truth_table": [\n            {\n             "a": "A",\n             "b": "B",\n             "x": "B"\n            },\n            {\n             "a": "A",\n             "b": "C",\n             "x": "C"\n            },\n            {\n             "a": "B",\n             "b": "C",\n             "x": "C"\n            }\n         ]\n        },\n        {\n         "name": "OR",\n         "truth_table": [\n            {\n             "a": "A",\n             "b": "B",\n             "x": "A"\n            },\n            {\n             "a": "A",\n             "b": "C",\n             "x": "A"\n            },\n            {\n             "a": "B",\n             "b": "C",\n             "x": "B"\n            }\n         ]\n        }\n     ]\n    }\n')),Object(i.b)("h3",{id:"response-3"},"Response"),Object(i.b)("p",null,"Headers: ",Object(i.b)("inlineCode",{parentName:"p"},"Status: 200 OK")),Object(i.b)("h4",{id:"response-body-3"},"Response body"),Object(i.b)("p",null,"Json Response"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{className:"language-json"}),'{\n "status": {\n  "message": "Operations Profile successfully updated (new snapshot created)",\n  "code": "200"\n }\n}\n')),Object(i.b)("a",{id:"5"}),Object(i.b)("h2",{id:"delete-delete-an-existing-aggregation-profile"},"[DELETE]",": Delete an existing aggregation profile"),Object(i.b)("p",null,"This method can be used to delete an existing aggregation profile"),Object(i.b)("h3",{id:"input-4"},"Input"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{}),"DELETE /operations_profiles/{ID}\n")),Object(i.b)("h4",{id:"request-headers-4"},"Request headers"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{}),"x-api-key: shared_key_value\nAccept: application/json\n")),Object(i.b)("h3",{id:"response-4"},"Response"),Object(i.b)("p",null,"Headers: ",Object(i.b)("inlineCode",{parentName:"p"},"Status: 200 OK")),Object(i.b)("h4",{id:"response-body-4"},"Response body"),Object(i.b)("p",null,"Json Response"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{className:"language-json"}),'{\n "status": {\n  "message": "Operations Profile Successfully Deleted",\n  "code": "200"\n }\n}\n')),Object(i.b)("h2",{id:"validation-checks"},"Validation Checks"),Object(i.b)("p",null,"When submitting or updating a new operations profile, validation checks are performed on json POST/PUT body for the following cases:"),Object(i.b)("ul",null,Object(i.b)("li",{parentName:"ul"},"Check if user has defined more than once a state name in available states list"),Object(i.b)("li",{parentName:"ul"},"Check if user has defined more than once an operation name in operations list"),Object(i.b)("li",{parentName:"ul"},"Check if user used an undefined state in operations"),Object(i.b)("li",{parentName:"ul"},"Check if truth table statements are adequate to handle all cases")),Object(i.b)("p",null,"When an invalid operations profile is submitted the api responds with a validation error list:"),Object(i.b)("h4",{id:"example-invalid-profile"},"Example invalid profile"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{className:"language-json"}),'{\n   "name": "ops1",\n   "available_states": [\n    "A","B","C","C"\n   ],\n   "defaults": {\n    "down": "B",\n    "missing": "FOO",\n    "unknown": "C"\n   },\n   "operations": [\n    {\n     "name": "AND",\n     "truth_table": [\n      {\n       "a": "A",\n       "b": "B",\n       "x": "B"\n      },\n      {\n       "a": "A",\n       "b": "C",\n       "x": "C"\n      },\n      {\n       "a": "B",\n       "b": "BAR",\n       "x": "C"\n      }\n     ]\n    },\n    {\n     "name": "OR",\n     "truth_table": [\n      {\n       "a": "A",\n       "b": "B",\n       "x": "A"\n      },\n      {\n       "a": "A",\n       "b": "C",\n       "x": "A"\n      },\n      {\n       "a": "B",\n       "b": "CAR",\n       "x": "B"\n      }\n     ]\n    },\n    {\n     "name": "OR",\n     "truth_table": [\n      {\n       "a": "A",\n       "b": "B",\n       "x": "A"\n      },\n      {\n       "a": "A",\n       "b": "C",\n       "x": "A"\n      },\n      {\n       "a": "B",\n       "b": "C",\n       "x": "B"\n      }\n     ]\n    }\n   ]\n  }\n')),Object(i.b)("p",null,"  The above profile definiton contains errors like: duplicate states, undefined states and unadequate statements in truth tables. Api response is the following:"),Object(i.b)("h3",{id:"response-5"},"Response"),Object(i.b)("p",null,"Headers: ",Object(i.b)("inlineCode",{parentName:"p"},"Status: 422 Unprocessable Entity")),Object(i.b)("h4",{id:"response-body-5"},"Response body"),Object(i.b)("pre",null,Object(i.b)("code",Object(a.a)({parentName:"pre"},{className:"language-json"}),'{\n"status": {\n "message": "Validation Error",\n "code": "422"\n},\n"errors": [\n {\n  "message": "Validation Failed",\n  "code": "422",\n  "details": "State:C is duplicated"\n },\n {\n  "message": "Validation Failed",\n  "code": "422",\n  "details": "Operation:OR is duplicated"\n },\n {\n  "message": "Validation Failed",\n  "code": "422",\n  "details": "Default Missing State: FOO not in available States"\n },\n {\n  "message": "Validation Failed",\n  "code": "422",\n  "details": "In Operation: AND, statement member b: BAR contains undeclared state"\n },\n {\n  "message": "Validation Failed",\n  "code": "422",\n  "details": "In Operation: OR, statement member b: CAR contains undeclared state"\n },\n {\n  "message": "Validation Failed",\n  "code": "422",\n  "details": "Not enough mentions of state:A in operation: AND"\n },\n {\n  "message": "Validation Failed",\n  "code": "422",\n  "details": "Not enough mentions of state:B in operation: AND"\n },\n {\n  "message": "Validation Failed",\n  "code": "422",\n  "details": "Not enough mentions of state:C in operation: AND"\n },\n {\n  "message": "Validation Failed",\n  "code": "422",\n  "details": "Not enough mentions of state:A in operation: OR"\n },\n {\n  "message": "Validation Failed",\n  "code": "422",\n  "details": "Not enough mentions of state:B in operation: OR"\n },\n {\n  "message": "Validation Failed",\n  "code": "422",\n  "details": "Not enough mentions of state:C in operation: OR"\n },\n {\n  "message": "Validation Failed",\n  "code": "422",\n  "details": "Not enough mentions of state:A in operation: OR"\n },\n {\n  "message": "Validation Failed",\n  "code": "422",\n  "details": "Not enough mentions of state:B in operation: OR"\n },\n {\n  "message": "Validation Failed",\n  "code": "422",\n  "details": "Not enough mentions of state:C in operation: OR"\n }\n]\n}\n')))}p.isMDXComponent=!0},90:function(e,n,t){"use strict";t.d(n,"a",(function(){return c})),t.d(n,"b",(function(){return u}));var a=t(0),r=t.n(a);function i(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function o(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);n&&(a=a.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,a)}return t}function b(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?o(Object(t),!0).forEach((function(n){i(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):o(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function l(e,n){if(null==e)return{};var t,a,r=function(e,n){if(null==e)return{};var t,a,r={},i=Object.keys(e);for(a=0;a<i.length;a++)t=i[a],n.indexOf(t)>=0||(r[t]=e[t]);return r}(e,n);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(a=0;a<i.length;a++)t=i[a],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(r[t]=e[t])}return r}var s=r.a.createContext({}),p=function(e){var n=r.a.useContext(s),t=n;return e&&(t="function"==typeof e?e(n):b(b({},n),e)),t},c=function(e){var n=p(e.components);return r.a.createElement(s.Provider,{value:n},e.children)},d={inlineCode:"code",wrapper:function(e){var n=e.children;return r.a.createElement(r.a.Fragment,{},n)}},O=r.a.forwardRef((function(e,n){var t=e.components,a=e.mdxType,i=e.originalType,o=e.parentName,s=l(e,["components","mdxType","originalType","parentName"]),c=p(t),O=a,u=c["".concat(o,".").concat(O)]||c[O]||d[O]||i;return t?r.a.createElement(u,b(b({ref:n},s),{},{components:t})):r.a.createElement(u,b({ref:n},s))}));function u(e,n){var t=arguments,a=n&&n.mdxType;if("string"==typeof e||a){var i=t.length,o=new Array(i);o[0]=O;var b={};for(var l in n)hasOwnProperty.call(n,l)&&(b[l]=n[l]);b.originalType=e,b.mdxType="string"==typeof e?e:a,o[1]=b;for(var s=2;s<i;s++)o[s]=t[s];return r.a.createElement.apply(null,o)}return r.a.createElement.apply(null,t)}O.displayName="MDXCreateElement"}}]);