db = db.getSiblingDB('argo_core');

db.authentication.insertMany([
{
    name: 'test_admin',
    email: 'test_admin@example.foo',
    api_key: 'WEB_API_ACCESS_TOKEN'
}
]);


db.tenants.insertMany([
    {
    id: 'e1ab046c-8544-47e6-bd8f-e8aa8b83acb3',
    info: {
      name: 'TENANT-TEST',
      email: 'test@gmail.com',
      description: 'this is test tenant description',
      image: 'https://example/image.png',
      website: 'https://test.tenant.org',
      created: '2025-01-01 00:00:00',
      updated: '2025-01-02 00:00:00'
    },
    db_conf: null,
    topology: { type: '', feed: '' },
    users: null
  }
])


db.roles.insertMany([
  {
    resource: 'reports.get',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'reports.delete',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'reports.update',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'metric_profiles.get',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'metric_profiles.create',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'metric_profiles.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'metric_profiles.delete',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'reports.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'metric_profiles.update',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'operations_profiles.create',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'operations_profiles.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'reports.create',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'operations_profiles.delete',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'operations_profiles.update',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'aggregation_profiles.get',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'aggregation_profiles.create',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'aggregation_profiles.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'operations_profiles.get',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'aggregation_profiles.update',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'aggregation_profiles.delete',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'results.get',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'status.get',
    roles: [ 'admin', 'editor', 'viewer', 'status_viewer', 'admin_ui' ]
  },
  {
    resource: 'results.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'status.list',
    roles: [ 'admin', 'editor', 'viewer', 'status_viewer', 'admin_ui' ]
  },
  {
    resource: 'factors.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'tenants.get',
    roles: [ 'super_admin', 'super_admin_ui' ]
  },
  {
    resource: 'tenants.create',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'tenants.update',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'tenants.delete',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'tenants.list',
    roles: [
      'super_admin',
      'super_admin_viewer',
      'super_admin_restricted',
      'super_admin_ui'
    ]
  },
  {
    resource: 'metric_result.get',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'recomputations.get',
    roles: [ 'admin', 'editor', 'admin_ui', 'viewer' ]
  },
  {
    resource: 'recomputations.list',
    roles: [ 'admin', 'editor', 'admin_ui', 'viewer' ]
  },
  {
    resource: 'recomputations.submit',
    roles: [ 'admin', 'editor', 'admin_ui' ]
  },
  {
    resource: 'aggregationProfiles.create',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'aggregationProfiles.get',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'aggregationProfiles.delete',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'aggregationProfiles.update',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'operationsProfiles.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'operationsProfiles.get',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
    {
    resource: 'operationsProfiles.delete',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'metricResult.get',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'operationsProfiles.update',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'operationsProfiles.create',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'metricProfiles.create',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'metricProfiles.delete',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'metricProfiles.update',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'metricProfiles.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'aggregationProfiles.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'metricProfiles.get',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'thresholdsProfiles.get',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'thresholdsProfiles.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'thresholdsProfiles.update',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'thresholdsProfiles.delete',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'thresholdsProfiles.create',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'tenants.get_status',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'latest.get',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'topology.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'recomputations.update',
    roles: [ 'admin', 'editor', 'admin_ui' ]
  },
  {
    resource: 'recomputations.changeStatus',
    roles: [ 'admin', 'editor', 'admin_ui' ]
  },
    {
    resource: 'recomputations.delete',
    roles: [ 'admin', 'editor', 'admin_ui' ]
  },
  {
    resource: 'recomputations.resetStatus',
    roles: [ 'admin', 'editor', 'admin_ui' ]
  },
  {
    resource: 'tenants.update_status',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'weights.update',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'weights.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui', 'admin_ui' ]
  },
  {
    resource: 'weights.get',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'weights.delete',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'downtimes.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui', 'admin_ui' ]
  },
  {
    resource: 'downtimes.get',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'weights.create',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'downtimes.create',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'downtimes.update',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'downtimes.delete',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'downtimes.options',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'topology_endpoints.list',
    roles: [ 'admin', 'editor', 'admin_ui', 'viewer' ]
  },
  {
    resource: 'weights.options',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'topology_endpoints.insert',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'topology_endpoints.delete',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'topology_groups.insert',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'topology_groups.list',
    roles: [ 'admin', 'editor', 'admin_ui', 'viewer' ]
  },
    {
    resource: 'topology_stats.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'topology_groups.delete',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'topology_groups_report.list',
    roles: [ 'admin', 'editor', 'viewer' ]
  },
  {
    resource: 'topology_groups_endpoint.list',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'topology_endpoints_report.list',
    roles: [ 'admin', 'editor', 'viewer' ]
  },
  {
    resource: 'feeds.topo.update',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'feeds.topo.get',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'tenants.create_user',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'tenants.list_users',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'tenants.user_refresh_token',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'tenants.get_user',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'list_endpoints',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'tenants.update_user',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'feeds.weights.update',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'tenants.delete_user',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'issues.list_endpoints',
    roles: [ 'admin', 'editor', 'admin_ui', 'viewer' ]
  },
  {
    resource: 'feeds.weights.get',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'trends.flapping_endpoints',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'trends.flapping_services',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'trends.flapping_endpoint_groups',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
    {
    resource: 'trends.flapping_metrics',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'trends.status_services',
    roles: [ 'admin', 'editor', 'viewer', 'admin-ui', 'admin_ui' ]
  },
  {
    resource: 'trends.status_endpoints',
    roles: [ 'admin', 'editor', 'viewer', 'admin-ui', 'admin_ui' ]
  },
  {
    resource: 'trends.status_endpoint_groups',
    roles: [ 'admin', 'editor', 'viewer', 'admin-ui', 'admin_ui' ]
  },
  {
    resource: 'feeds.data.get',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'feeds.data.update',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'metrics_admin.get',
    roles: [ 'super_admin', 'metrics_admin', 'super_admin_ui' ]
  },
  {
    resource: 'metrics.get',
    roles: [ 'admin', 'editor', 'viewer' ]
  },
  {
    resource: 'metrics_admin.update',
    roles: [ 'super_admin', 'metrics_admin', 'super_admin_ui' ]
  },
  {
    resource: 'metrics_report.get',
    roles: [ 'admin', 'editor', 'viewer' ]
  },
  {
    resource: 'trends.flapping_metrics_tags',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'trends.status_metrics_tags',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'v3.ar.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'topology_service_types.insert',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'trends.status_metrics',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'topology_service_types.delete',
    roles: [ 'admin', 'editor' ]
  },
  {
    resource: 'v3.status.list',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'topology_tags.list',
    roles: [ 'admin', 'editor', 'viewer' ]
  },
  {
    resource: 'v3.status.list-by-id',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'v3.ar.list-by-id',
    roles: [ 'admin', 'editor', 'viewer', 'admin-ui' ]
  },
  {
    resource: 'topology_service_types.list',
    roles: [ 'admin', 'editor', 'viewer' ]
  },
  {
    resource: 'issues.list_group_metrics',
    roles: [ 'admin', 'editor', 'viewer', 'admin_ui' ]
  },
  {
    resource: 'health',
    roles: [ 'viewer', 'editor' ]
  },
  {
    resource: 'consistency.result',
    roles: [ 'consistency-viewer' ]
  },
  {
    resource: 'consistency.auto-check',
    roles: [ 'consistency-check' ]
  },
  {
    resource: 'consistency.ack',
    roles: [ 'consistency-ack' ]
  },
  {
    resource: 'tenants.update_info',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'tenants.update_db_conf',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'tenants.update_topology',
    roles: [ 'super_admin' ]
  }  
])
