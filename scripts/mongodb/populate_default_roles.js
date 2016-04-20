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
  {"resource" : "metric_profiles.get", "roles" : [ "admin", "editor","viewer"] },
  {"resource" : "metric_profiles.list", "roles" : [ "admin", "editor","viewer" ]},
  {"resource" : "metric_profiles.create", "roles" : [ "admin", "editor" ] },
  {"resource" : "metric_profiles.delete", "roles" : [ "admin", "editor" ] },
  {"resource" : "metric_profiles.update", "roles" : [ "admin", "editor" ] },
  {"resource" : "operations_profiles.get", "roles" : [ "admin", "editor","viewer"] },
  {"resource" : "operations_profiles.list", "roles" : [ "admin", "editor","viewer" ]},
  {"resource" : "operations_profiles.create", "roles" : [ "admin", "editor" ] },
  {"resource" : "operations_profiles.delete", "roles" : [ "admin", "editor" ] },
  {"resource" : "operations_profiles.update", "roles" : [ "admin", "editor" ] },
  {"resource" : "aggregation_profiles.get", "roles" : [ "admin", "editor","viewer"] },
  {"resource" : "aggregation_profiles.list", "roles" : [ "admin", "editor","viewer" ]},
  {"resource" : "aggregation_profiles.create", "roles" : [ "admin", "editor" ] },
  {"resource" : "aggregation_profiles.delete", "roles" : [ "admin", "editor" ] },
  {"resource" : "aggregation_profiles.update", "roles" : [ "admin", "editor" ] },
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
  {"resource" : "metric_result.get", "roles" : [ "admin", "editor","viewer" ]},
  {"resource" : "recomputations.list", "roles" : [ "admin","editor"]},
  {"resource" : "recomputations.get", "roles" : [ "admin","editor"]},
  {"resource" : "recomputations.submit", "roles" : [ "admin","editor"]}]);
  print("INFO\tPolulated default roles in \'roles\' collection")
}
