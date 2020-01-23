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

function populate_default_roles() {
    db = db.getSiblingDB("argo_core");
    print("INFO\tOpened argo_core db");
    db.roles.insert([
        {
            resource: "topology_stats.list",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        {
            resource: "latest.get",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        {
            resource: "reports.get",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        {
            resource: "reports.list",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        { resource: "reports.create", roles: ["admin", "editor"] },
        { resource: "reports.delete", roles: ["admin", "editor"] },
        { resource: "reports.update", roles: ["admin", "editor"] },
        {
            resource: "metricProfiles.get",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        {
            resource: "metricProfiles.list",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        { resource: "metricProfiles.create", roles: ["admin", "editor"] },
        { resource: "metricProfiles.delete", roles: ["admin", "editor"] },
        { resource: "metricProfiles.update", roles: ["admin", "editor"] },
        {
            resource: "operationsProfiles.get",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        {
            resource: "operationsProfiles.list",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        { resource: "operationsProfiles.create", roles: ["admin", "editor"] },
        { resource: "operationsProfiles.delete", roles: ["admin", "editor"] },
        { resource: "operationsProfiles.update", roles: ["admin", "editor"] },
        {
            resource: "aggregationProfiles.get",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        {
            resource: "aggregationProfiles.list",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        { resource: "aggregationProfiles.create", roles: ["admin", "editor"] },
        { resource: "aggregationProfiles.delete", roles: ["admin", "editor"] },
        { resource: "aggregationProfiles.update", roles: ["admin", "editor"] },
        {
            resource: "thresholdsProfiles.get",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        {
            resource: "thresholdsProfiles.list",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        { resource: "thresholdsProfiles.update", roles: ["admin", "editor"] },
        { resource: "thresholdsProfiles.create", roles: ["admin", "editor"] },
        { resource: "thresholdsProfiles.delete", roles: ["admin", "editor"] },
        {
            resource: "weights.get",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        {
            resource: "weights.list",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        { resource: "weights.create", roles: ["admin", "editor"] },
        { resource: "weights.delete", roles: ["admin", "editor"] },
        { resource: "weights.update", roles: ["admin", "editor"] },
        {
            resource: "results.get",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        {
            resource: "results.list",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        {
            resource: "status.get",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        {
            resource: "status.list",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        {
            resource: "factors.list",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        { resource: "tenants.get", roles: ["super_admin"] },
        {
            resource: "tenants.list",
            roles: ["super_admin", "super_admin_restricted", "super_admin_ui"]
        },
        { resource: "tenants.create", roles: ["super_admin"] },
        { resource: "tenants.delete", roles: ["super_admin"] },
        { resource: "tenants.update", roles: ["super_admin"] },
        {
            resource: "tenants.user_by_id",
            roles: ["super_admin", "super_admin_restricted", "super_admin_ui"]
        },
        {
            resource: "tenants.get_status",
            roles: ["super_admin", "super_admin_restricted", "super_admin_ui"]
        },
        { resource: "tenants.update_status", roles: ["super_admin"] },
        {
            resource: "metricResult.get",
            roles: ["admin", "editor", "viewer", "admin_ui"]
        },
        {
            resource: "recomputations.list",
            roles: ["admin", "editor", "admin_ui"]
        },
        {
            resource: "recomputations.get",
            roles: ["admin", "editor", "admin_ui"]
        },
        {
            resource: "recomputations.submit",
            roles: ["admin", "editor", "admin_ui"]
        },
        {
            resource: "recomputations.delete",
            roles: ["admin", "editor", "admin_ui"]
        },
        {
            resource: "recomputations.update",
            roles: ["admin", "editor", "admin_ui"]
        },
        {
            resource: "topology_endpoints.insert",
            roles: ["admin", "editor"]
        },
        {
            resource: "topology_endpoints.list",
            roles: ["admin", "editor"]
        },
        {
            resource: "topology_endpoints.delete",
            roles: ["admin", "editor"]
        },
        {
            resource: "topology_groups.insert",
            roles: ["admin", "editor"]
        },
        {
            resource: "topology_groups.list",
            roles: ["admin", "editor"]
        },
        {
            resource: "topology_groups.delete",
            roles: ["admin", "editor"]
        }
    ]);
    print("INFO\tPolulated default roles in 'roles' collection");
}
