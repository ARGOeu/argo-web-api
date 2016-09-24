// Default Role definitions
//
// Roles are defined in 'roles' collection in argo_core database
// Default roles include the following:
//  - super_admin: only able to edit global tenant information (tenants/users/roles)
//  - admin: able to view/edit all information under a tenant
//  - editor: same as tenant admin (but not able to change privileges - future feature)
//  - viewer: able to only view all information under a tenant (cannot edit)
//
// To run the script issue in mongo shell the following commands:
//   > load('./populate_default_roles.js')
//   true
//   > populate_default_roles()
//   INFO	Opened argo_core db
//   INFO	Polulated default roles in 'roles' collection


function populate_default_roles()
{
  db = db.getSiblingDB('argo_core')
  print("INFO\tOpened argo_core db")
  db.roles.insert([
  {"resource" : "reports.get", "roles" : [ "admin", "editor","viewer"] },
  {"resource" : "reports.list", "roles" : [ "admin", "editor","viewer" ]},
  {"resource" : "reports.create", "roles" : [ "admin", "editor" ] },
  {"resource" : "reports.delete", "roles" : [ "admin", "editor" ] },
  {"resource" : "reports.update", "roles" : [ "admin", "editor" ] },
  {"resource" : "metricProfiles.get", "roles" : [ "admin", "editor","viewer"] },
  {"resource" : "metricProfiles.list", "roles" : [ "admin", "editor","viewer" ]},
  {"resource" : "metricProfiles.create", "roles" : [ "admin", "editor" ] },
  {"resource" : "metricProfiles.delete", "roles" : [ "admin", "editor" ] },
  {"resource" : "metricProfiles.update", "roles" : [ "admin", "editor" ] },
  {"resource" : "operationsProfiles.get", "roles" : [ "admin", "editor","viewer"] },
  {"resource" : "operationsProfiles.list", "roles" : [ "admin", "editor","viewer" ]},
  {"resource" : "operationsProfiles.create", "roles" : [ "admin", "editor" ] },
  {"resource" : "operationsProfiles.delete", "roles" : [ "admin", "editor" ] },
  {"resource" : "operationsProfiles.update", "roles" : [ "admin", "editor" ] },
  {"resource" : "aggregationProfiles.get", "roles" : [ "admin", "editor","viewer"] },
  {"resource" : "aggregationProfiles.list", "roles" : [ "admin", "editor","viewer" ]},
  {"resource" : "aggregationProfiles.create", "roles" : [ "admin", "editor" ] },
  {"resource" : "aggregationProfiles.delete", "roles" : [ "admin", "editor" ] },
  {"resource" : "aggregationProfiles.update", "roles" : [ "admin", "editor" ] },
  {"resource" : "results.get", "roles" : [ "admin", "editor","viewer"] },
  {"resource" : "results.list", "roles" : [ "admin", "editor","viewer" ]},
  {"resource" : "status.get", "roles" : [ "admin", "editor","viewer"] },
  {"resource" : "status.list", "roles" : [ "admin", "editor","viewer" ]},
  {"resource" : "factors.list", "roles" : [ "admin", "editor","viewer"] },
  {"resource" : "tenants.get", "roles" : [ "super_admin"] },
  {"resource" : "tenants.list", "roles" : [ "super_admin" ]},
  {"resource" : "tenants.create", "roles" : [ "super_admin" ] },
  {"resource" : "tenants.delete", "roles" : [ "super_admin" ] },
  {"resource" : "tenants.update", "roles" : [ "super_admin" ] },
  {"resource" : "metricResult.get", "roles" : [ "admin", "editor","viewer" ]},
  {"resource" : "recomputations.list", "roles" : [ "admin","editor"]},
  {"resource" : "recomputations.get", "roles" : [ "admin","editor"]},
  {"resource" : "recomputations.submit", "roles" : [ "admin","editor"]}]);
  print("INFO\tPolulated default roles in \'roles\' collection")
}
