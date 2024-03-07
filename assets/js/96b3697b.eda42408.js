"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[5207],{5481:(e,n,s)=>{s.r(n),s.d(n,{assets:()=>l,contentTitle:()=>t,default:()=>h,frontMatter:()=>i,metadata:()=>r,toc:()=>o});var a=s(4848),d=s(8453);const i={id:"tenants",title:"Tenants",slug:"/",sidebar_position:1},t=void 0,r={id:"tenants_and_feeds/tenants",title:"Tenants",description:"API Calls",source:"@site/docs/tenants_and_feeds/tenants.md",sourceDirName:"tenants_and_feeds",slug:"/",permalink:"/argo-web-api/docs/",draft:!1,unlisted:!1,tags:[],version:"current",sidebarPosition:1,frontMatter:{id:"tenants",title:"Tenants",slug:"/",sidebar_position:1},sidebar:"tutorialSidebar",previous:{title:"Tenants & Feeds",permalink:"/argo-web-api/docs/category/tenants--feeds"},next:{title:"Feeds",permalink:"/argo-web-api/docs/tenants_and_feeds/feeds"}},l={},o=[{value:"API Calls",id:"api-calls",level:2},{value:"[GET]: List Tenants",id:"1",level:2},{value:"Input",id:"input",level:3},{value:"Request headers",id:"request-headers",level:4},{value:"Response",id:"response",level:3},{value:"Response body for super admin users",id:"response-body-for-super-admin-users",level:4},{value:"Response body for restricted super admin users:",id:"response-body-for-restricted-super-admin-users",level:4},{value:"Response body for super_admin_ui users:",id:"response-body-for-super_admin_ui-users",level:4},{value:"[GET]: List A Specific tenant",id:"2",level:2},{value:"Input",id:"input-1",level:3},{value:"Request headers",id:"request-headers-1",level:4},{value:"Response",id:"response-1",level:3},{value:"Response body",id:"response-body",level:4},{value:"Response body for super_admin_ui users:",id:"response-body-for-super_admin_ui-users-1",level:4},{value:"[GET]: List A Specific User",id:"get-list-a-specific-user",level:2},{value:"Input",id:"input-2",level:3},{value:"Request headers",id:"request-headers-2",level:4},{value:"Response",id:"response-2",level:3},{value:"Response body",id:"response-body-1",level:4},{value:"NOTE",id:"note",level:3},{value:"[POST]: Create a new Tenant",id:"3",level:2},{value:"Input",id:"input-3",level:3},{value:"Request headers",id:"request-headers-3",level:4},{value:"POST BODY",id:"post-body",level:4},{value:"Response",id:"response-3",level:3},{value:"Response body",id:"response-body-2",level:4},{value:"[PUT]: Update information on an existing tenant",id:"4",level:2},{value:"Input",id:"input-4",level:3},{value:"Request headers",id:"request-headers-4",level:4},{value:"PUT BODY",id:"put-body",level:4},{value:"Response",id:"response-4",level:3},{value:"Response body",id:"response-body-3",level:4},{value:"[DELETE]: Delete an existing tenant",id:"5",level:2},{value:"Input",id:"input-5",level:3},{value:"Request headers",id:"request-headers-5",level:4},{value:"Response",id:"response-5",level:3},{value:"Response body",id:"response-body-4",level:4},{value:"[GET]: List A Specific tenant&#39;s argo-engine status",id:"6",level:2},{value:"Input",id:"input-6",level:3},{value:"Request headers",id:"request-headers-6",level:4},{value:"Response",id:"response-6",level:3},{value:"Response body",id:"response-body-5",level:4},{value:"[PUT]: Update argo-engine status information on an existing tenant",id:"7",level:2},{value:"Input",id:"input-7",level:3},{value:"Request headers",id:"request-headers-7",level:4},{value:"PUT BODY",id:"put-body-1",level:4},{value:"Response",id:"response-7",level:3},{value:"Response body",id:"response-body-6",level:4},{value:"[POST]: Create new user",id:"8",level:2},{value:"Input",id:"input-8",level:3},{value:"Request headers",id:"request-headers-8",level:4},{value:"PUT BODY",id:"put-body-2",level:4},{value:"Response",id:"response-8",level:3},{value:"Response body",id:"response-body-7",level:4},{value:"[PUT]: Update user",id:"9",level:2},{value:"Input",id:"input-9",level:3},{value:"Request headers",id:"request-headers-9",level:4},{value:"PUT BODY",id:"put-body-3",level:4},{value:"Response",id:"response-9",level:3},{value:"Response body",id:"response-body-8",level:4},{value:"[POST]: Renew User API key",id:"10",level:2},{value:"Input",id:"input-10",level:3},{value:"Request headers",id:"request-headers-10",level:4},{value:"Response",id:"response-10",level:3},{value:"Response body",id:"response-body-9",level:4},{value:"[DELETE]: Delete User",id:"11",level:2},{value:"Input",id:"input-11",level:3},{value:"Request headers",id:"request-headers-11",level:4},{value:"Response",id:"response-11",level:3},{value:"Response body",id:"response-body-10",level:4},{value:"[GET]: List all available users that belong to a specific tenant",id:"12",level:2},{value:"Input",id:"input-12",level:3},{value:"Request headers",id:"request-headers-12",level:4},{value:"Response",id:"response-12",level:3},{value:"Response body",id:"response-body-11",level:4},{value:"[GET]: Get user",id:"13",level:2},{value:"Input",id:"input-13",level:3},{value:"Request headers",id:"request-headers-13",level:4},{value:"Response",id:"response-13",level:3},{value:"Response body",id:"response-body-12",level:4}];function c(e){const n={a:"a",code:"code",h2:"h2",h3:"h3",h4:"h4",p:"p",pre:"pre",strong:"strong",table:"table",tbody:"tbody",td:"td",th:"th",thead:"thead",tr:"tr",...(0,d.R)(),...e.components};return(0,a.jsxs)(a.Fragment,{children:[(0,a.jsx)(n.h2,{id:"api-calls",children:"API Calls"}),"\n",(0,a.jsxs)(n.table,{children:[(0,a.jsx)(n.thead,{children:(0,a.jsxs)(n.tr,{children:[(0,a.jsx)(n.th,{children:"Name"}),(0,a.jsx)(n.th,{children:"Description"}),(0,a.jsx)(n.th,{children:"Shortcut"})]})}),(0,a.jsxs)(n.tbody,{children:[(0,a.jsxs)(n.tr,{children:[(0,a.jsx)(n.td,{children:"GET: List Tenants"}),(0,a.jsx)(n.td,{children:"This method can be used to retrieve a list of current tenants"}),(0,a.jsx)(n.td,{children:(0,a.jsx)(n.a,{href:"#1",children:" Description"})})]}),(0,a.jsxs)(n.tr,{children:[(0,a.jsx)(n.td,{children:"GET: List a specific tenant"}),(0,a.jsx)(n.td,{children:"This method can be used to retrieve a specific metric tenant based on its id."}),(0,a.jsx)(n.td,{children:(0,a.jsx)(n.a,{href:"#2",children:" Description"})})]}),(0,a.jsxs)(n.tr,{children:[(0,a.jsx)(n.td,{children:"POST: Create a new tenant"}),(0,a.jsx)(n.td,{children:"This method can be used to create a new tenant"}),(0,a.jsx)(n.td,{children:(0,a.jsx)(n.a,{href:"#3",children:" Description"})})]}),(0,a.jsxs)(n.tr,{children:[(0,a.jsx)(n.td,{children:"PUT: Update a tenant"}),(0,a.jsx)(n.td,{children:"This method can be used to update information on an existing tenant"}),(0,a.jsx)(n.td,{children:(0,a.jsx)(n.a,{href:"#4",children:" Description"})})]}),(0,a.jsxs)(n.tr,{children:[(0,a.jsx)(n.td,{children:"DELETE: Delete a tenant"}),(0,a.jsx)(n.td,{children:"This method can be used to delete an existing tenant"}),(0,a.jsx)(n.td,{children:(0,a.jsx)(n.a,{href:"#5",children:" Description"})})]}),(0,a.jsxs)(n.tr,{children:[(0,a.jsx)(n.td,{children:"GET: Get a tenant's arg engine status"}),(0,a.jsx)(n.td,{children:"This method can be used to get status for a specific tenant"}),(0,a.jsx)(n.td,{children:(0,a.jsx)(n.a,{href:"#6",children:" Description"})})]}),(0,a.jsxs)(n.tr,{children:[(0,a.jsx)(n.td,{children:"PUT: Update a tenant's engine status"}),(0,a.jsx)(n.td,{children:"This method can be used to update argo engine status information for a specific tenant"}),(0,a.jsx)(n.td,{children:(0,a.jsx)(n.a,{href:"#7",children:" Description"})})]}),(0,a.jsxs)(n.tr,{children:[(0,a.jsx)(n.td,{children:"POST: Create tenant user"}),(0,a.jsx)(n.td,{children:"This method can be used to add a new user to existing tenant"}),(0,a.jsx)(n.td,{children:(0,a.jsx)(n.a,{href:"#8",children:" Description"})})]}),(0,a.jsxs)(n.tr,{children:[(0,a.jsx)(n.td,{children:"PUT: Update tenant user"}),(0,a.jsx)(n.td,{children:"This method can be used to update information on an existing user of a specific tenant"}),(0,a.jsx)(n.td,{children:(0,a.jsx)(n.a,{href:"#9",children:" Description"})})]}),(0,a.jsxs)(n.tr,{children:[(0,a.jsx)(n.td,{children:"POST: Renew User's API key"}),(0,a.jsx)(n.td,{children:"This method can be used to renew user's api key"}),(0,a.jsx)(n.td,{children:(0,a.jsx)(n.a,{href:"#10",children:" Description"})})]}),(0,a.jsxs)(n.tr,{children:[(0,a.jsx)(n.td,{children:"DELETE: Delete Users"}),(0,a.jsx)(n.td,{children:"This method can be used to remove and delete a user from a specific tenant"}),(0,a.jsx)(n.td,{children:(0,a.jsx)(n.a,{href:"#11",children:" Description"})})]}),(0,a.jsxs)(n.tr,{children:[(0,a.jsx)(n.td,{children:"GET: List Users"}),(0,a.jsx)(n.td,{children:"This method can be used to list all users that belong to a specific tenant"}),(0,a.jsx)(n.td,{children:(0,a.jsx)(n.a,{href:"#12",children:" Description"})})]})]})]}),"\n",(0,a.jsx)(n.h2,{id:"1",children:"[GET]: List Tenants"}),"\n",(0,a.jsx)(n.p,{children:"This method can be used to retrieve a list of current tenants"}),"\n",(0,a.jsxs)(n.p,{children:[(0,a.jsx)(n.strong,{children:"Note"}),": This method restricts tenant database and user information when the x-api-key token holder is a ",(0,a.jsx)(n.strong,{children:"restricted"})," super admin\n",(0,a.jsx)(n.strong,{children:"Note"}),": This method shows only tenants that have admin ui users when the x-api-key token holder is a ",(0,a.jsx)(n.strong,{children:"super_admin_ui"})]}),"\n",(0,a.jsx)(n.h3,{id:"input",children:"Input"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"GET /admin/tenants\n"})}),"\n",(0,a.jsx)(n.h4,{id:"request-headers",children:"Request headers"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,a.jsx)(n.h3,{id:"response",children:"Response"}),"\n",(0,a.jsxs)(n.p,{children:["Headers: ",(0,a.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,a.jsx)(n.h4,{id:"response-body-for-super-admin-users",children:"Response body for super admin users"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n "status": {\n  "message": "Success",\n  "code": "200"\n },\n "data": [\n  {\n   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",\n   "info": {\n    "name": "Tenant1",\n    "email": "email1@tenant1.com",\n    "description" : "a simple tenant",\n    "image" : "url to image",\n    "website": "www.tenant1.com",\n    "created": "2015-10-20 02:08:04",\n    "updated": "2015-10-20 02:08:04"\n   },\n   "db_conf": [\n    {\n     "store": "ar",\n     "server": "a.mongodb.org",\n     "port": 27017,\n     "database": "ar_db",\n     "username": "admin",\n     "password": "3NCRYPT3D"\n    },\n    {\n     "store": "status",\n     "server": "b.mongodb.org",\n     "port": 27017,\n     "database": "status_db",\n     "username": "admin",\n     "password": "3NCRYPT3D"\n    }\n   ],\n   "topology": {\n    "type": "GOCDB",\n    "feed": "gocdb1.example.foo"\n   },\n   "users": [\n    {\n     "id": "acb74194-553a-11e9-8647-d663bd873d93",\n     "name": "cap",\n     "email": "cap@email.com",\n     "api_key": "C4PK3Y",\n     "roles": [\n        "admin"\n     ]\n    },\n    {\n    "id": "acb74194-553a-11e9-8647-d663bd873d94",\n     "name": "thor",\n     "email": "thor@email.com",\n     "api_key": "TH0RK3Y",\n     "roles": [\n        "viewer"\n     ]\n    }\n   ]\n  },\n  {\n   "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",\n   "info": {\n    "name": "tenant2",\n    "email": "tenant2@email.com",\n    "description" : "a simple tenant",\n    "image" : "url to image",\n    "website": "www.tenant2.com",\n    "created": "2015-10-20 02:08:04",\n    "updated": "2015-10-20 02:08:04"\n   },\n   "db_conf": [\n    {\n     "store": "ar",\n     "server": "a.mongodb.org",\n     "port": 27017,\n     "database": "ar_db",\n     "username": "admin",\n     "password": "3NCRYPT3D"\n    },\n    "topology": {\n    "type": "GOCDB",\n    "feed": "gocdb2.example.foo"\n   },\n    {\n     "store": "status",\n     "server": "b.mongodb.org",\n     "port": 27017,\n     "database": "status_db",\n     "username": "admin",\n     "password": "3NCRYPT3D"\n    }\n   ],\n   "users": [\n    {\n    "id": "acb74194-553a-11e9-8647-d663bd873d95",\n     "name": "groot",\n     "email": "groot@email.com",\n     "api_key": "GR00TK3Y",\n     "roles": [\n         "admin", "admin_ui"\n      ]\n    },\n    {\n    "id": "acb74194-553a-11e9-8647-d663bd873d97",\n     "name": "starlord",\n     "email": "starlord@email.com",\n     "api_key": "ST4RL0RDK3Y",\n     "roles": [\n         "admin"\n      ]\n    }\n   ]\n  }\n ]\n}\n'})}),"\n",(0,a.jsx)(n.h4,{id:"response-body-for-restricted-super-admin-users",children:"Response body for restricted super admin users:"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n    "status": {\n        "message": "Success",\n        "code": "200"\n    },\n    "data": [\n        {\n            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",\n            "info": {\n                "name": "Tenant1",\n                "email": "email1@tenant1.com",\n                "description": "a simple tenant",\n                "image": "url to image",\n                "website": "www.tenant1.com",\n                "created": "2015-10-20 02:08:04",\n                "updated": "2015-10-20 02:08:04"\n            },\n            "topology": {\n                "type": "GOCDB",\n                "feed": "gocdb1.example.foo"\n            }\n        },\n        {\n            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",\n            "info": {\n                "name": "tenant2",\n                "email": "tenant2@email.com",\n                "description": "a simple tenant",\n                "image": "url to image",\n                "website": "www.tenant2.com",\n                "created": "2015-10-20 02:08:04",\n                "updated": "2015-10-20 02:08:04"\n            },\n            "topology": {\n                "type": "GOCDB",\n                "feed": "gocdb2.example.foo"\n            }\n        }\n    ]\n}\n'})}),"\n",(0,a.jsx)(n.h4,{id:"response-body-for-super_admin_ui-users",children:"Response body for super_admin_ui users:"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n    "status": {\n        "message": "Success",\n        "code": "200"\n    },\n    "data": [\n        {\n            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",\n            "info": {\n                "name": "tenant2",\n                "email": "tenant2@email.com",\n                "description": "a simple tenant",\n                "image": "url to image",\n                "website": "www.tenant2.com",\n                "created": "2015-10-20 02:08:04",\n                "updated": "2015-10-20 02:08:04"\n            },\n            "topology": {\n                "type": "GOCDB",\n                "feed": "gocdb2.example.foo"\n            },\n            "users": [\n                {\n                    "id": "acb74194-553a-11e9-8647-d663bd873d95",\n                    "name": "groot",\n                    "email": "groot@email.com",\n                    "api_key": "GR00TK3Y",\n                    "roles": ["admin", "admin_ui"]\n                }\n            ]\n        }\n    ]\n}\n'})}),"\n",(0,a.jsx)(n.h2,{id:"2",children:"[GET]: List A Specific tenant"}),"\n",(0,a.jsx)(n.p,{children:"This method can be used to retrieve specific tenant based on its id"}),"\n",(0,a.jsxs)(n.p,{children:[(0,a.jsx)(n.strong,{children:"Note"}),": This method shows only tenants that have admin ui users when the x-api-key token holder is a ",(0,a.jsx)(n.strong,{children:"super_admin_ui"})]}),"\n",(0,a.jsx)(n.h3,{id:"input-1",children:"Input"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"GET /admin/tenants/{ID}\n"})}),"\n",(0,a.jsx)(n.h4,{id:"request-headers-1",children:"Request headers"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,a.jsx)(n.h3,{id:"response-1",children:"Response"}),"\n",(0,a.jsxs)(n.p,{children:["Headers: ",(0,a.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,a.jsx)(n.h4,{id:"response-body",children:"Response body"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n    "status": {\n        "message": "Success",\n        "code": "200"\n    },\n    "data": [\n        {\n            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",\n            "info": {\n                "name": "tenant2",\n                "email": "tenant2@email.com",\n                "description": "a simple tenant",\n                "image": "url to image",\n                "website": "www.tenant2.com",\n                "created": "2015-10-20 02:08:04",\n                "updated": "2015-10-20 02:08:04"\n            },\n            "db_conf": [\n                {\n                    "store": "ar",\n                    "server": "a.mongodb.org",\n                    "port": 27017,\n                    "database": "ar_db",\n                    "username": "admin",\n                    "password": "3NCRYPT3D"\n                },\n                {\n                    "store": "status",\n                    "server": "b.mongodb.org",\n                    "port": 27017,\n                    "database": "status_db",\n                    "username": "admin",\n                    "password": "3NCRYPT3D"\n                }\n            ],\n            "topology": {\n                "type": "GOCDB",\n                "feed": "gocdb1.example.foo"\n            },\n            "users": [\n                {\n                    "id": "acb74194-553a-11e9-8647-d663bd873d95",\n                    "name": "groot",\n                    "email": "groot@email.com",\n                    "api_key": "GR00TK3Y",\n                    "roles": ["admin", "admin_ui"]\n                },\n                {\n                    "id": "acb74194-553a-11e9-8647-d663bd873d97",\n                    "name": "starlord",\n                    "email": "starlord@email.com",\n                    "api_key": "ST4RL0RDK3Y",\n                    "roles": ["admin"]\n                }\n            ]\n        }\n    ]\n}\n'})}),"\n",(0,a.jsx)(n.h4,{id:"response-body-for-super_admin_ui-users-1",children:"Response body for super_admin_ui users:"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n    "status": {\n        "message": "Success",\n        "code": "200"\n    },\n    "data": [\n        {\n            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50c",\n            "info": {\n                "name": "tenant2",\n                "email": "tenant2@email.com",\n                "description": "a simple tenant",\n                "image": "url to image",\n                "website": "www.tenant2.com",\n                "created": "2015-10-20 02:08:04",\n                "updated": "2015-10-20 02:08:04"\n            },\n            "topology": {\n                "type": "GOCDB",\n                "feed": "gocdb2.example.foo"\n            },\n            "users": [\n                {\n                    "id": "acb74194-553a-11e9-8647-d663bd873d95",\n                    "name": "groot",\n                    "email": "groot@email.com",\n                    "api_key": "GR00TK3Y",\n                    "roles": ["admin", "admin_ui"]\n                }\n            ]\n        }\n    ]\n}\n'})}),"\n",(0,a.jsx)(n.h2,{id:"get-list-a-specific-user",children:"[GET]: List A Specific User"}),"\n",(0,a.jsx)(n.p,{children:"This method can be used to retrieve specific user based on its id"}),"\n",(0,a.jsx)(n.h3,{id:"input-2",children:"Input"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"GET /admin/users:byID/{ID}\n"})}),"\n",(0,a.jsx)(n.h4,{id:"request-headers-2",children:"Request headers"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,a.jsx)(n.h3,{id:"response-2",children:"Response"}),"\n",(0,a.jsxs)(n.p,{children:["Headers: ",(0,a.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,a.jsx)(n.h4,{id:"response-body-1",children:"Response body"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n    "status": {\n        "message": "Success",\n        "code": "200"\n    },\n    "data": [\n        {\n            "id": "acb74194-553a-11e9-8647-d663bd873d93",\n            "name": "cap",\n            "email": "cap@email.com",\n            "api_key": "C4PK3Y",\n            "roles": ["admin"]\n        }\n    ]\n}\n'})}),"\n",(0,a.jsx)(n.h3,{id:"note",children:"NOTE"}),"\n",(0,a.jsxs)(n.p,{children:["Specifying the filter, ",(0,a.jsx)(n.code,{children:"export=flat"}),", it will return a flat user json object"]}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n    "id": "acb74194-553a-11e9-8647-d663bd873d93",\n    "name": "cap",\n    "email": "cap@email.com",\n    "api_key": "C4PK3Y",\n    "roles": ["admin"]\n}\n'})}),"\n",(0,a.jsx)(n.h2,{id:"3",children:"[POST]: Create a new Tenant"}),"\n",(0,a.jsx)(n.p,{children:"This method can be used to insert a new tenant"}),"\n",(0,a.jsx)(n.h3,{id:"input-3",children:"Input"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"POST /admin/tenants\n"})}),"\n",(0,a.jsx)(n.h4,{id:"request-headers-3",children:"Request headers"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,a.jsx)(n.h4,{id:"post-body",children:"POST BODY"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n    "info": {\n        "name": "Tenant1",\n        "email": "email1@tenant1.com",\n        "website": "www.tenant1.com",\n        "description": "a simple tenant",\n        "image": "url to image",\n        "created": "2015-10-20 02:08:04",\n        "updated": "2015-10-20 02:08:04"\n    },\n    "db_conf": [\n        {\n            "store": "ar",\n            "server": "a.mongodb.org",\n            "port": 27017,\n            "database": "ar_db",\n            "username": "admin",\n            "password": "3NCRYPT3D"\n        },\n        {\n            "store": "status",\n            "server": "b.mongodb.org",\n            "port": 27017,\n            "database": "status_db",\n            "username": "admin",\n            "password": "3NCRYPT3D"\n        }\n    ],\n    "topology": {\n        "type": "GOCDB",\n        "feed": "gocdb.example.foo"\n    },\n    "users": [\n        {\n            "name": "cap",\n            "email": "cap@email.com",\n            "api_key": "C4PK3Y",\n            "roles": ["admin"]\n        },\n        {\n            "name": "thor",\n            "email": "thor@email.com",\n            "api_key": "TH0RK3Y",\n            "roles": ["admin"]\n        }\n    ]\n}\n'})}),"\n",(0,a.jsx)(n.h3,{id:"response-3",children:"Response"}),"\n",(0,a.jsxs)(n.p,{children:["Headers: ",(0,a.jsx)(n.code,{children:"Status: 201 Created"})]}),"\n",(0,a.jsx)(n.h4,{id:"response-body-2",children:"Response body"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n    "status": {\n        "message": "Tenant was succesfully created",\n        "code": "201"\n    },\n    "data": {\n        "id": "{{ID}}",\n        "links": {\n            "self": "https:///api/v2/admin/tenants/{{ID}}"\n        }\n    }\n}\n'})}),"\n",(0,a.jsx)(n.h2,{id:"4",children:"[PUT]: Update information on an existing tenant"}),"\n",(0,a.jsx)(n.p,{children:"This method can be used to update information on an existing tenant"}),"\n",(0,a.jsx)(n.h3,{id:"input-4",children:"Input"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"PUT /admin/tenants/{ID}\n"})}),"\n",(0,a.jsx)(n.h4,{id:"request-headers-4",children:"Request headers"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,a.jsx)(n.h4,{id:"put-body",children:"PUT BODY"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n    "info": {\n        "name": "Tenant1",\n        "email": "email1@tenant1.com",\n        "description": "a changed description",\n        "image": "a changed url to nwe image",\n        "website": "www.tenant1.com",\n        "created": "2015-10-20 02:08:04",\n        "updated": "2015-10-20 02:08:04"\n    },\n    "db_conf": [\n        {\n            "store": "ar",\n            "server": "a.mongodb.org",\n            "port": 27017,\n            "database": "ar_db",\n            "username": "admin",\n            "password": "3NCRYPT3D"\n        },\n        {\n            "store": "status",\n            "server": "b.mongodb.org",\n            "port": 27017,\n            "database": "status_db",\n            "username": "admin",\n            "password": "3NCRYPT3D"\n        }\n    ],\n    "topology": {\n        "type": "GOCDB",\n        "feed": "gocdb.example.foo"\n    },\n    "users": [\n        {\n            "name": "cap",\n            "email": "cap@email.com",\n            "api_key": "C4PK3Y",\n            "roles": ["admin"]\n        },\n        {\n            "name": "thor",\n            "email": "thor@email.com",\n            "api_key": "TH0RK3Y",\n            "roles": ["admin"]\n        }\n    ]\n}\n'})}),"\n",(0,a.jsx)(n.h3,{id:"response-4",children:"Response"}),"\n",(0,a.jsxs)(n.p,{children:["Headers: ",(0,a.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,a.jsx)(n.h4,{id:"response-body-3",children:"Response body"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n    "status": {\n        "message": "Tenant successfully updated",\n        "code": "200"\n    }\n}\n'})}),"\n",(0,a.jsx)(n.h2,{id:"5",children:"[DELETE]: Delete an existing tenant"}),"\n",(0,a.jsx)(n.p,{children:"This method can be used to delete an existing tenant"}),"\n",(0,a.jsx)(n.h3,{id:"input-5",children:"Input"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"DELETE /admin/tenants/{ID}\n"})}),"\n",(0,a.jsx)(n.h4,{id:"request-headers-5",children:"Request headers"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,a.jsx)(n.h3,{id:"response-5",children:"Response"}),"\n",(0,a.jsxs)(n.p,{children:["Headers: ",(0,a.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,a.jsx)(n.h4,{id:"response-body-4",children:"Response body"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n    "status": {\n        "message": "Tenant Successfully Deleted",\n        "code": "200"\n    }\n}\n'})}),"\n",(0,a.jsx)(n.h2,{id:"6",children:"[GET]: List A Specific tenant's argo-engine status"}),"\n",(0,a.jsx)(n.p,{children:"This method can be used to retrieve specific tenant's status based on its id"}),"\n",(0,a.jsx)(n.h3,{id:"input-6",children:"Input"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"GET /admin/tenants/{ID}/status\n"})}),"\n",(0,a.jsx)(n.h4,{id:"request-headers-6",children:"Request headers"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,a.jsx)(n.h3,{id:"response-6",children:"Response"}),"\n",(0,a.jsxs)(n.p,{children:["Headers: ",(0,a.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,a.jsx)(n.h4,{id:"response-body-5",children:"Response body"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n    "status": {\n        "message": "Success",\n        "code": "200"\n    },\n    "data": [\n        {\n            "id": "6ac7d684-1f8e-4a02-a502-720e8f11e50b",\n            "info": {\n                "name": "tenant1",\n                "email": "email1@tenant1.com",\n                "description": "a simple tenant",\n                "image": "url to image",\n                "website": "www.tenant1.com",\n                "created": "2015-10-20 02:08:04",\n                "updated": "2015-10-20 02:08:04"\n            },\n            "status": {\n                "total_status": false,\n                "ams": {\n                    "metric_data": {\n                        "ingestion": false,\n                        "publishing": false,\n                        "status_streaming": false,\n                        "messages_arrived": 0\n                    },\n                    "sync_data": {\n                        "ingestion": false,\n                        "publishing": false,\n                        "status_streaming": false,\n                        "messages_arrived": 0\n                    }\n                },\n                "hdfs": {\n                    "metric_data": false\n                },\n                "engine_config": false,\n                "last_check": ""\n            }\n        }\n    ]\n}\n'})}),"\n",(0,a.jsx)(n.h2,{id:"7",children:"[PUT]: Update argo-engine status information on an existing tenant"}),"\n",(0,a.jsx)(n.p,{children:"This method can be used to update status information on an existing tenant"}),"\n",(0,a.jsx)(n.h3,{id:"input-7",children:"Input"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"PUT /admin/tenants/{ID}/status\n"})}),"\n",(0,a.jsx)(n.h4,{id:"request-headers-7",children:"Request headers"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,a.jsx)(n.h4,{id:"put-body-1",children:"PUT BODY"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n    "ams": {\n        "metric_data": {\n            "ingestion": true,\n            "publishing": true,\n            "status_streaming": false,\n            "messages_arrived": 100\n        },\n        "sync_data": {\n            "ingestion": true,\n            "publishing": false,\n            "status_streaming": true,\n            "messages_arrived": 200\n        }\n    },\n    "hdfs": {\n        "metric_data": true,\n        "sync_data": {\n            "Critical": {\n                "downtimes": true,\n                "group_endpoints": true,\n                "blank_recompuation": true,\n                "group_groups": true,\n                "weights": true,\n                "operations_profile": true,\n                "metric_profile": true,\n                "aggregation_profile": true\n            }\n        }\n    },\n    "engine_config": true,\n    "last_check": "2018-08-10T12:32:45Z"\n}\n'})}),"\n",(0,a.jsx)(n.h3,{id:"response-7",children:"Response"}),"\n",(0,a.jsxs)(n.p,{children:["Headers: ",(0,a.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,a.jsx)(n.h4,{id:"response-body-6",children:"Response body"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n    "status": {\n        "message": "Tenant successfully updated",\n        "code": "200"\n    }\n}\n'})}),"\n",(0,a.jsx)(n.h2,{id:"8",children:"[POST]: Create new user"}),"\n",(0,a.jsx)(n.p,{children:"This method can be used to create a new user on existing tenant"}),"\n",(0,a.jsx)(n.h3,{id:"input-8",children:"Input"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"POST /admin/tenants/{ID}/users\n"})}),"\n",(0,a.jsx)(n.h4,{id:"request-headers-8",children:"Request headers"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,a.jsx)(n.h4,{id:"put-body-2",children:"PUT BODY"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'  {\n    "name":"new_user",\n    "email":"new_user@email.com",\n    "roles": [\n        "admin"\n    ]\n  }`\n'})}),"\n",(0,a.jsx)(n.h3,{id:"response-8",children:"Response"}),"\n",(0,a.jsxs)(n.p,{children:["Headers: ",(0,a.jsx)(n.code,{children:"Status: 201 OK"})]}),"\n",(0,a.jsx)(n.h4,{id:"response-body-7",children:"Response body"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n "status": {\n  "message": "User was successfully created",\n  "code": "201"\n },\n "data": {\n  "id": "1cb883eb-8b40-428d-bce6-8ec23a9f3ca8",\n  "links": {\n   "self": "https:///api/v2/admin/tenants/6ac7d684-1f8e-4a02-a502-720e8f11e50b/users/1cb883eb-8b40-428d-bce6-8ec23a9f3ca8"\n  }\n }\n}\n'})}),"\n",(0,a.jsx)(n.h2,{id:"9",children:"[PUT]: Update user"}),"\n",(0,a.jsx)(n.p,{children:"This method can be used to update an existing user in a specific tenant"}),"\n",(0,a.jsx)(n.h3,{id:"input-9",children:"Input"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"PUT /admin/tenants/{ID}/users/{USER_ID}\n"})}),"\n",(0,a.jsx)(n.h4,{id:"request-headers-9",children:"Request headers"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,a.jsx)(n.h4,{id:"put-body-3",children:"PUT BODY"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'  {\n    "name":"new_user",\n    "email":"new_user@email.com",\n    "roles": [\n        "admin"\n    ]\n  }`\n'})}),"\n",(0,a.jsx)(n.h3,{id:"response-9",children:"Response"}),"\n",(0,a.jsxs)(n.p,{children:["Headers: ",(0,a.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,a.jsx)(n.h4,{id:"response-body-8",children:"Response body"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n "status": {\n  "message": "User succesfully updated",\n  "code": "200"\n }\n}\n'})}),"\n",(0,a.jsx)(n.h2,{id:"10",children:"[POST]: Renew User API key"}),"\n",(0,a.jsx)(n.p,{children:"This method can be used to renew a user's api access key"}),"\n",(0,a.jsx)(n.h3,{id:"input-10",children:"Input"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"POST /admin/tenants/{ID}/users/{USER_ID}/renew_api_key\n"})}),"\n",(0,a.jsx)(n.h4,{id:"request-headers-10",children:"Request headers"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,a.jsx)(n.h3,{id:"response-10",children:"Response"}),"\n",(0,a.jsxs)(n.p,{children:["Headers: ",(0,a.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,a.jsx)(n.h4,{id:"response-body-9",children:"Response body"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n  "status": {\n    "message": "User api key succesfully renewed",\n    "code": "200"\n  },\n  "data": {\n    "api_key": "s3cr3tT0k3n"\n  }\n}\n'})}),"\n",(0,a.jsx)(n.h2,{id:"11",children:"[DELETE]: Delete User"}),"\n",(0,a.jsx)(n.p,{children:"This method can be used to remove and delete a user from a tenant"}),"\n",(0,a.jsx)(n.h3,{id:"input-11",children:"Input"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"DELETE /admin/tenants/{ID}/users/{USER_ID}\n"})}),"\n",(0,a.jsx)(n.h4,{id:"request-headers-11",children:"Request headers"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,a.jsx)(n.h3,{id:"response-11",children:"Response"}),"\n",(0,a.jsxs)(n.p,{children:["Headers: ",(0,a.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,a.jsx)(n.h4,{id:"response-body-10",children:"Response body"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n "status": {\n  "message": "User succesfully deleted",\n  "code": "200"\n }\n}\n'})}),"\n",(0,a.jsx)(n.h2,{id:"12",children:"[GET]: List all available users that belong to a specific tenant"}),"\n",(0,a.jsx)(n.p,{children:"This method can be used to list all available users that are members of a specific tenant"}),"\n",(0,a.jsx)(n.h3,{id:"input-12",children:"Input"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"GET /admin/tenants/{ID}/users\n"})}),"\n",(0,a.jsx)(n.h4,{id:"request-headers-12",children:"Request headers"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,a.jsx)(n.h3,{id:"response-12",children:"Response"}),"\n",(0,a.jsxs)(n.p,{children:["Headers: ",(0,a.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,a.jsx)(n.h4,{id:"response-body-11",children:"Response body"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n "status": {\n  "message": "Success",\n  "code": "200"\n },\n "data": [\n  {\n   "id": "acb74194-553a-11e9-8647-d663bd873d93",\n   "name": "user_a",\n   "email": "user_a@email.com",\n   "api_key": "user_a_key",\n   "roles": [\n    "admin",\n    "admin_ui"\n   ]\n  },\n  {\n   "id": "acb74432-553a-11e9-8647-d663bd873d93",\n   "name": "user_b",\n   "email": "user_b@email.com",\n   "api_key": "user_b_key",\n   "roles": [\n    "admin"\n   ]\n  }\n ]\n}\n\n'})}),"\n",(0,a.jsx)(n.h2,{id:"13",children:"[GET]: Get user"}),"\n",(0,a.jsx)(n.p,{children:"This method can be used to get information on specific user of a specific tenant"}),"\n",(0,a.jsx)(n.h3,{id:"input-13",children:"Input"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"GET /admin/tenants/{ID}/users/{USER_ID}\n"})}),"\n",(0,a.jsx)(n.h4,{id:"request-headers-13",children:"Request headers"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{children:"x-api-key: shared_key_value\nAccept: application/json\n"})}),"\n",(0,a.jsx)(n.h3,{id:"response-13",children:"Response"}),"\n",(0,a.jsxs)(n.p,{children:["Headers: ",(0,a.jsx)(n.code,{children:"Status: 200 OK"})]}),"\n",(0,a.jsx)(n.h4,{id:"response-body-12",children:"Response body"}),"\n",(0,a.jsx)(n.p,{children:"Json Response"}),"\n",(0,a.jsx)(n.pre,{children:(0,a.jsx)(n.code,{className:"language-json",children:'{\n "status": {\n  "message": "Success",\n  "code": "200"\n },\n "data": [\n  {\n   "id": "acb74432-553a-11e9-8647-d663bd873d93",\n   "name": "user_a",\n   "email": "user_a@email.com",\n   "api_key": "user_a_key",\n   "roles": [\n    "admin"\n   ]\n  }\n ]\n}\n\n'})})]})}function h(e={}){const{wrapper:n}={...(0,d.R)(),...e.components};return n?(0,a.jsx)(n,{...e,children:(0,a.jsx)(c,{...e})}):c(e)}},8453:(e,n,s)=>{s.d(n,{R:()=>t,x:()=>r});var a=s(6540);const d={},i=a.createContext(d);function t(e){const n=a.useContext(i);return a.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function r(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(d):e.components||d:t(e.components),a.createElement(i.Provider,{value:n},e.children)}}}]);