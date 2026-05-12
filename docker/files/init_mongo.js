db = db.getSiblingDB('argo_core');

db.authentication.insertMany([
{
    name: 'test_admin',
    email: 'test_admin@example.foo',
    api_key: 'WEB_API_ACCESS_TOKEN'
},
{
    name: 'super_admin',
    email: 'super_admin@example.foo',
    api_key: 'super_admin_key'
},
]);


db.tenants.insertMany([
    {
    id: '42c1152d-e23c-4a19-b51a-b27f1eb7f37f',
    info: {
      name: 'TENANT-TEST',
      email: 'test@example.foo',
      description: 'this is test tenant description',
      image: 'https://test.example.foo/image.png',
      website: 'https://test.example.foo',
      created: '2015-01-01 00:00:00',
      updated: '2015-01-02 00:00:00'
    },
    db_conf: [{
      "store": "ar",
      "server": "localhost",
      "port": 27017,
      "database": "argo_TENANT-TEST",
      "username": "admin",
      "password": ""
    }]      
    ,
    topology: { type: '', feed: '' },
    users: [
      {
          "id": "763723d1-b758-4073-abcf-c60211413b28",
          "name": "admin_testtenant",
          "email": "argo-dev@example.foo",
          "api_key": "admin_testtenant_key",
          "roles": [
              "admin"
          ]
      },
      {
          "id": "f65bf4c5-436f-45dd-8959-92a9bce777b0",
          "name": "admin2_testtenant",
          "email": "argo-dev@example.foo",
          "api_key": "admin2_testtenant_key",
          "roles": [
              "admin"
          ]
      },
      {
          "id": "aaf55727-e1e3-4112-a73c-7dcf27438fcf",
          "name": "editor_testtenant",
          "email": "argo-dev@example.foo",
          "api_key": "editor_testtenant_key",
          "roles": [
              "editor"
          ]
      },
      {
          "id": "22cbd08e-727b-46f8-95d8-a510a626ce2b",
          "name": "admin_viewer_EOSCBEYOND",
          "email": "argo-dev@example.foo",
          "api_key": "viewer_testtenant_key",
          "roles": [
              "viewer"
          ]
      },
    ]
  },
  {
    id: '6b36d6d3-56a3-48a5-93af-aecf3e16a7c6',
    info: {
      name: 'TENANTB',
      email: 'tenantb@example.foo',
      description: 'this is tenant b description',
      image: 'https://tenantb.example.foo/image.png',
      website: 'https://tenantb.example.foo',
      created: '2015-01-01 00:00:00',
      updated: '2015-01-02 00:00:00'
    },
    db_conf: [{
      "store": "ar",
      "server": "localhost",
      "port": 27017,
      "database": "argo_TENANTB",
      "username": "admin",
      "password": ""
    }]      
    ,
    topology: { type: '', feed: '' },
    users: [
      {
          "id": "175bd877-abe3-4e11-bb96-dc4b297ef4d9",
          "name": "admin_tenantb",
          "email": "argo-dev@example.foo",
          "api_key": "admin_tenantb_key",
          "roles": [
              "admin"
          ]
      },
      {
          "id": "9c2a98b7-c695-495d-b625-1210d584f49f",
          "name": "admin2_tenantb",
          "email": "argo-dev@example.foo",
          "api_key": "admin2_tenantb_key",
          "roles": [
              "admin"
          ]
      },
      {
          "id": "2d67f58b-11b1-4ef3-8cd8-5c12f6f6ef2c",
          "name": "editor_tenantb",
          "email": "argo-dev@example.foo",
          "api_key": "editor_tenantb_key",
          "roles": [
              "editor"
          ]
      },
      {
          "id": "1d8e2a23-a7f4-478a-9b31-3ca69840ec36",
          "name": "admin_viewer_EOSCBEYOND",
          "email": "argo-dev@example.foo",
          "api_key": "viewer_tenantb_key",
          "roles": [
              "viewer"
          ]
      },
    ]
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
    resource: 'reports.set_node_report',
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
    resource: 'tenants.node_set',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'tenants.node_unset',
    roles: [ 'super_admin' ]
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
  },
  {
    resource: 'v3.components.access_refresh',
    roles: [ 'component_engine', 'component_monbox', 'component_poem' ]
  },
  {
    resource: 'tenants.update_ready',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'tenants.get_ready',
    roles: [ 'super_admin' ]
  },
  {
    resource: 'tenants.update_node',
    roles: ['super_admin']
  },
  {
    resource: 'reports.set_node_report',
    roles: ['super_admin', 'admin', 'editor']
  },
  {
    resource: 'v4.nodes.summary',
    roles: [ 'super_admin', 'admin', 'editor', 'viewer' ]
  },
  {
    resource: 'v4.nodes.summary.item',
    roles: [ 'super_admin', 'admin', 'editor', 'viewer' ]
  },
  {
    resource: 'v4.nodes.availability',
    roles: [ 'super_admin', 'admin', 'editor', 'viewer' ]
  },
  {
    resource: 'v4.nodes.availability.item',
    roles: [ 'super_admin', 'admin', 'editor', 'viewer' ]
  },
  {
    resource: 'v4.nodes.uptime',
    roles: [ 'super_admin', 'admin', 'editor', 'viewer' ]
  },
  {
    resource: 'v4.nodes.uptime.item',
    roles: [ 'super_admin', 'admin', 'editor', 'viewer' ]
  },
  {
    resource: 'v4.nodes.status',
    roles: [ 'super_admin', 'admin', 'editor', 'viewer' ]
  },
  {
    resource: 'v4.nodes.status.item',
    roles: [ 'super_admin', 'admin', 'editor', 'viewer' ]
  }
])

db = db.getSiblingDB('argo_TENANT-TEST');

db.topology_endpoints.ensureIndex({ "date_integer": -1, "id": 1 })
db.topology_groups.ensureIndex({ "date_integer": -1, "id": 1 })
db.topology_service_types.ensureIndex({ "date_integer": -1, "id": 1 })
db.metric_profiles.ensureIndex({ "date_integer": -1, "id": 1 })
db.operations_profiles.ensureIndex({ "date_integer": -1, "id": 1 })
db.aggregation_profiles.ensureIndex({ "date_integer": -1, "id": 1 })

db.topology_service_types.insertMany(
[
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "name": "webportal",
    "title": "Web Portal",
    "description": "Generic service representing a web portal",
    "tags": [
      "topology",
      "web",
      "http"
    ]
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "name": "api",
    "title": "Rest API",
    "description": "Generic service representing a web api",
    "tags": [
      "topology",
      "rest",
      "http",
      "api"
    ]
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "name": "iam",
    "title": "Identity Management",
    "description": "Generic service representing identity management",
    "tags": [
      "topology",
      "iam",
      "oidc"
    ]
  }
]  
)

db.topology_endpoints.insertMany(
  [
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "ESHOP",
    "type": "SERVICEGROUPS",
    "service": "webportal",
    "hostname": "eshop.tenant-test.foo",
    "notifications": {
      "contacts": [
        ""
      ],
      "enabled": true
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANT-TEST"
    }
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "ESHOP",
    "type": "SERVICEGROUPS",
    "service": "webportal",
    "hostname": "eshop2.tenant-test.foo",
    "notifications": {
      "contacts": [
        ""
      ],
      "enabled": true
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANT-TEST"
    }
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "ESHOP",
    "type": "SERVICEGROUPS",
    "service": "api",
    "hostname": "api-eshop.tenant-test.foo",
    "notifications": {
      "contacts": [
        ""
      ],
      "enabled": true
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANT-TEST"
    }
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "HELPDESK",
    "type": "SERVICEGROUPS",
    "service": "webportal",
    "hostname": "kb.tenant-test.foo",
    "notifications": {
      "contacts": [
        ""
      ],
      "enabled": true
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANT-TEST"
    }
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "HELPDESK",
    "type": "SERVICEGROUPS",
    "service": "webportal",
    "hostname": "help.tenant-test.foo",
    "notifications": {
      "contacts": [
        ""
      ],
      "enabled": true
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANT-TEST"
    }
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "HELPDESK",
    "type": "SERVICEGROUPS",
    "service": "api",
    "hostname": "api-eshop.tenant-test.foo",
    "notifications": {
      "contacts": [
        ""
      ],
      "enabled": true
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANT-TEST"
    }
  }
]
)

db.topology_groups.insertMany(
  [
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "PROJECTA",
    "type": "PROJECT",
    "subgroup": "ESHOP",
    "notifications": {
      "contacts": [
        "eshop-admin@example.foo"
      ]
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANT-TEST"
    }
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "PROJECTA",
    "type": "PROJECT",
    "subgroup": "HELPDESK",
    "notifications": {
      "contacts": [
        "helpdesk-admin@example.foo"
      ]
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANT-TEST"
    }
  }
]
)

db.operations_profiles.insertMany([
{
  "id": "d168c01b-6ddc-4dd6-bc76-6caad1a326ee",
  "date_integer": 20150101,
  "date": "2015-01-01",
  "name": "default_ops",
  "available_states": [
    "OK",
    "WARNING",
    "UNKNOWN",
    "MISSING",
    "CRITICAL",
    "DOWNTIME"
  ],
  "defaults": {
    "down": "DOWNTIME",
    "missing": "MISSING",
    "unknown": "UNKNOWN"
  },
  "operations": [
    {
      "name": "AND",
      "truth_table": [
        {
          "a": "OK",
          "b": "OK",
          "x": "OK"
        },
        {
          "a": "OK",
          "b": "WARNING",
          "x": "WARNING"
        },
        {
          "a": "OK",
          "b": "UNKNOWN",
          "x": "UNKNOWN"
        },
        {
          "a": "OK",
          "b": "MISSING",
          "x": "MISSING"
        },
        {
          "a": "OK",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "OK",
          "b": "DOWNTIME",
          "x": "DOWNTIME"
        },
        {
          "a": "WARNING",
          "b": "WARNING",
          "x": "WARNING"
        },
        {
          "a": "WARNING",
          "b": "UNKNOWN",
          "x": "UNKNOWN"
        },
        {
          "a": "WARNING",
          "b": "MISSING",
          "x": "MISSING"
        },
        {
          "a": "WARNING",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "WARNING",
          "b": "DOWNTIME",
          "x": "DOWNTIME"
        },
        {
          "a": "UNKNOWN",
          "b": "UNKNOWN",
          "x": "UNKNOWN"
        },
        {
          "a": "UNKNOWN",
          "b": "MISSING",
          "x": "MISSING"
        },
        {
          "a": "UNKNOWN",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "UNKNOWN",
          "b": "DOWNTIME",
          "x": "DOWNTIME"
        },
        {
          "a": "MISSING",
          "b": "MISSING",
          "x": "MISSING"
        },
        {
          "a": "MISSING",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "MISSING",
          "b": "DOWNTIME",
          "x": "DOWNTIME"
        },
        {
          "a": "CRITICAL",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "CRITICAL",
          "b": "DOWNTIME",
          "x": "CRITICAL"
        },
        {
          "a": "DOWNTIME",
          "b": "DOWNTIME",
          "x": "DOWNTIME"
        }
      ]
    },
    {
      "name": "OR",
      "truth_table": [
        {
          "a": "OK",
          "b": "OK",
          "x": "OK"
        },
        {
          "a": "OK",
          "b": "WARNING",
          "x": "OK"
        },
        {
          "a": "OK",
          "b": "UNKNOWN",
          "x": "OK"
        },
        {
          "a": "OK",
          "b": "MISSING",
          "x": "OK"
        },
        {
          "a": "OK",
          "b": "CRITICAL",
          "x": "OK"
        },
        {
          "a": "OK",
          "b": "DOWNTIME",
          "x": "OK"
        },
        {
          "a": "WARNING",
          "b": "WARNING",
          "x": "WARNING"
        },
        {
          "a": "WARNING",
          "b": "UNKNOWN",
          "x": "WARNING"
        },
        {
          "a": "WARNING",
          "b": "MISSING",
          "x": "WARNING"
        },
        {
          "a": "WARNING",
          "b": "CRITICAL",
          "x": "WARNING"
        },
        {
          "a": "WARNING",
          "b": "DOWNTIME",
          "x": "WARNING"
        },
        {
          "a": "UNKNOWN",
          "b": "UNKNOWN",
          "x": "UNKNOWN"
        },
        {
          "a": "UNKNOWN",
          "b": "MISSING",
          "x": "UNKNOWN"
        },
        {
          "a": "UNKNOWN",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "UNKNOWN",
          "b": "DOWNTIME",
          "x": "UNKNOWN"
        },
        {
          "a": "MISSING",
          "b": "MISSING",
          "x": "MISSING"
        },
        {
          "a": "MISSING",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "MISSING",
          "b": "DOWNTIME",
          "x": "DOWNTIME"
        },
        {
          "a": "CRITICAL",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "CRITICAL",
          "b": "DOWNTIME",
          "x": "CRITICAL"
        },
        {
          "a": "DOWNTIME",
          "b": "DOWNTIME",
          "x": "DOWNTIME"
        }
      ]
    }
  ]
}
])

db.metric_profiles.insertMany(
  [
  {
    "id": "bb3cf905-e270-4a05-b053-4234d73b97ba",
    "date": "2015-01-01",
    "date_integer": 20150101,
    "name": "MON_ALL",
    "description": "all checks",
    "services": [
      {
        "service": "webportal",
        "metrics": [
          "http_check",
          "cert_check"
        ]
      },
      {
        "service": "api",
        "metrics": [
          "http_check",
          "cert_check",
          "api_check",
          "auth_check",
          "option_check"
        ]
      },
      {
        "service": "aai",
        "metrics": [
          "http_check",
          "cert_check",
          "redirect_check",
          "login_check"
        ]
      }
    ]
  },
  {
    "id": "5e8712f7-1b38-4f27-8ca4-d0c697760bf3",
    "date": "2015-01-01",
    "date_integer": 20150101,
    "name": "MON_HTTP",
    "description": "just the http checks",
    "services": [
      {
        "service": "webportal",
        "metrics": [
          "http_check",
          "cert_check"
        ]
      },
      {
        "service": "api",
        "metrics": [
          "http_check",
          "cert_check"
        ]
      },
      {
        "service": "aai",
        "metrics": [
          "http_check",
          "cert_check"
        ]
      }
    ]
  }
]
)

db.aggregation_profiles.insertMany(
[
  {
    "id": "b01238e6-597a-4491-b5ec-b30164cbd9a1",
    "date": "2015-01-01",
    "date_integer": 20150101,
    "name": "AGGR_HA",
    "namespace": "tenant",
    "endpoint_group": "SERVICEGROUPS",
    "metric_operation": "AND",
    "profile_operation": "AND",
    "metric_profile": {
      "name": "MON_ALL",
      "id": "bb3cf905-e270-4a05-b053-4234d73b97ba"
    },
    "groups": [
      {
        "name": "eshop",
        "operation": "AND",
        "services": [
          {
            "name": "webportal",
            "operation": "OR"
          },
          {
            "name": "api",
            "operation": "OR"
          },
          {
            "name": "iam",
            "operation": "OR"
          }
        ]
      }
    ]
  }
]
)

db.reports.insertMany(
  [
  {
    "id": "cf010255-cda3-49d8-92d1-926c2c6cf9eb",
    "tenant": "",
    "disabled": false,
    "info": {
      "name": "CORE",
      "description": "Core A/R report",
      "created": "2015-01-01 12:53:27",
      "updated": "2015-01-01 17:56:33"
    },
    "computations": {
      "ar": true,
      "status": true,
      "trends": [
        "flapping",
        "status",
        "tags"
      ]
    },
    "thresholds": {
      "availability": 80,
      "reliability": 90,
      "uptime": 0.800000011920929,
      "unknown": 0.10000000149011612,
      "downtime": 0.10000000149011612
    },
    "topology_schema": {
      "group": {
        "type": "PROJECT",
        "group": {
          "type": "SERVICEGROUPS"
        }
      }
    },
    "profiles": [
      {
        "id": "bb3cf905-e270-4a05-b053-4234d73b97ba",
        "name": "MON_ALL",
        "type": "metric"
      },
      {
        "id": "b01238e6-597a-4491-b5ec-b30164cbd9a1",
        "name": "AGGR_HA",
        "type": "aggregation"
      },
      {
        "id": "d168c01b-6ddc-4dd6-bc76-6caad1a326ee",
        "name": "default_ops",
        "type": "operations"
      }
    ],
    "filter_tags": [],
    "node_report": true
  },
  {
    "id": "c7a6b0d4-4885-46da-9dd1-1f91d0e9142e",
    "tenant": "",
    "disabled": false,
    "info": {
      "name": "Just-Http",
      "description": "Just Http Checks report",
      "created": "2015-01-01 12:53:27",
      "updated": "2015-01-01 17:56:33"
    },
    "computations": {
      "ar": true,
      "status": true,
      "trends": [
        "flapping",
        "status",
        "tags"
      ]
    },
    "thresholds": {
      "availability": 80,
      "reliability": 90,
      "uptime": 0.800000011920929,
      "unknown": 0.10000000149011612,
      "downtime": 0.10000000149011612
    },
    "topology_schema": {
      "group": {
        "type": "PROJECT",
        "group": {
          "type": "SERVICEGROUPS"
        }
      }
    },
    "profiles": [
      {
        "id": "5e8712f7-1b38-4f27-8ca4-d0c697760bf3",
        "name": "MON_HTTP",
        "type": "metric"
      },
      {
        "id": "b01238e6-597a-4491-b5ec-b30164cbd9a1",
        "name": "AGGR_HA",
        "type": "aggregation"
      },
      {
        "id": "d168c01b-6ddc-4dd6-bc76-6caad1a326ee",
        "name": "default_ops",
        "type": "operations"
      }
    ],
    "filter_tags": [
      {
        "name": "monitored",
        "value": "1",
        "context": "argo.group.filter.tags"
      }
    ],
    "node_report": true
  }
]
)


db = db.getSiblingDB('argo_TENANTB');

db.topology_endpoints.ensureIndex({ "date_integer": -1, "id": 1 })
db.topology_groups.ensureIndex({ "date_integer": -1, "id": 1 })
db.topology_service_types.ensureIndex({ "date_integer": -1, "id": 1 })
db.metric_profiles.ensureIndex({ "date_integer": -1, "id": 1 })
db.operations_profiles.ensureIndex({ "date_integer": -1, "id": 1 })
db.aggregation_profiles.ensureIndex({ "date_integer": -1, "id": 1 })


db.topology_service_types.insertMany(
[
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "name": "compute",
    "title": "Cloud Compute",
    "description": "Generic service representing cloud compute",
    "tags": [
      "topology",
      "web",
      "http"
    ]
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "name": "storage",
    "title": "Cloud Storage",
    "description": "Generic service representing a cloud storage",
    "tags": [
      "topology",
      "rest",
      "http",
      "api"
    ]
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "name": "iam",
    "title": "Identity Management",
    "description": "Generic service representing identity management",
    "tags": [
      "topology",
      "iam",
      "oidc"
    ]
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "name": "messaging",
    "title": "Cloud Messaging",
    "description": "Generic service representing cloud messaging",
    "tags": [
      "topology",
      "iam",
      "oidc"
    ]
  }
]  
)

db.topology_endpoints.insertMany(
  [
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "CLOUD-A",
    "type": "SERVICEGROUPS",
    "service": "compute",
    "hostname": "compute1.cloud-a.foo",
    "notifications": {
      "contacts": [
        ""
      ],
      "enabled": true
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANTB"
    }
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "CLOUD-A",
    "type": "SERVICEGROUPS",
    "service": "compute",
    "hostname": "compute2.cloud-a.foo",
    "notifications": {
      "contacts": [
        ""
      ],
      "enabled": true
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANTB"
    }
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "CLOUD-A",
    "type": "SERVICEGROUPS",
    "service": "storage",
    "hostname": "storage.cloud-a.foo",
    "notifications": {
      "contacts": [
        ""
      ],
      "enabled": true
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANTB"
    }
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "CLOUD-A",
    "type": "SERVICEGROUPS",
    "service": "iam",
    "hostname": "iam.cloud-a.foo",
    "notifications": {
      "contacts": [
        "iam-ops@cloud-a.foo"
      ],
      "enabled": true
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANTB",
    }
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "CLOUD-B",
    "type": "SERVICEGROUPS",
    "service": "storage",
    "hostname": "storage.cloud-b.foo",
    "notifications": {
      "contacts": [
        "storage-ops@cloud-b.foo"
      ],
      "enabled": true
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANTB",
      "cloud": "B"
    }
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "CLOUD-B",
    "type": "SERVICEGROUPS",
    "service": "iam",
    "hostname": "iam@cloud-b.foo",
    "notifications": {
      "contacts": [
        ""
      ],
      "enabled": true
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANTB",
      "cloud": "B"
    }
  }
]
)

db.topology_groups.insertMany(
  [
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "CLOUDINFRA",
    "type": "PROJECT",
    "subgroup": "CLOUD-A",
    "notifications": {
      "contacts": [
        "cloud-a@cloudinfra.foo"
      ]
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANTB"
    }
  },
  {
    "date": "2015-01-01",
    "date_integer": 20150101,
    "group": "CLOUDINFRA",
    "type": "PROJECT",
    "subgroup": "CLOUD-B",
    "notifications": {
      "contacts": [
        "cloud-b@cloudinfra.foo"
      ]
    },
    "tags": {
      "monitored": "1",
      "scope": "TENANTB"
    }
  },
]
)

db.operations_profiles.insertMany([
{
  "id": "b131b9ac-07ef-442b-a20a-e18faf534f26",
  "date_integer": 20150101,
  "date": "2015-01-01",
  "name": "default_ops",
  "available_states": [
    "OK",
    "WARNING",
    "UNKNOWN",
    "MISSING",
    "CRITICAL",
    "DOWNTIME"
  ],
  "defaults": {
    "down": "DOWNTIME",
    "missing": "MISSING",
    "unknown": "UNKNOWN"
  },
  "operations": [
    {
      "name": "AND",
      "truth_table": [
        {
          "a": "OK",
          "b": "OK",
          "x": "OK"
        },
        {
          "a": "OK",
          "b": "WARNING",
          "x": "WARNING"
        },
        {
          "a": "OK",
          "b": "UNKNOWN",
          "x": "UNKNOWN"
        },
        {
          "a": "OK",
          "b": "MISSING",
          "x": "MISSING"
        },
        {
          "a": "OK",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "OK",
          "b": "DOWNTIME",
          "x": "DOWNTIME"
        },
        {
          "a": "WARNING",
          "b": "WARNING",
          "x": "WARNING"
        },
        {
          "a": "WARNING",
          "b": "UNKNOWN",
          "x": "UNKNOWN"
        },
        {
          "a": "WARNING",
          "b": "MISSING",
          "x": "MISSING"
        },
        {
          "a": "WARNING",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "WARNING",
          "b": "DOWNTIME",
          "x": "DOWNTIME"
        },
        {
          "a": "UNKNOWN",
          "b": "UNKNOWN",
          "x": "UNKNOWN"
        },
        {
          "a": "UNKNOWN",
          "b": "MISSING",
          "x": "MISSING"
        },
        {
          "a": "UNKNOWN",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "UNKNOWN",
          "b": "DOWNTIME",
          "x": "DOWNTIME"
        },
        {
          "a": "MISSING",
          "b": "MISSING",
          "x": "MISSING"
        },
        {
          "a": "MISSING",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "MISSING",
          "b": "DOWNTIME",
          "x": "DOWNTIME"
        },
        {
          "a": "CRITICAL",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "CRITICAL",
          "b": "DOWNTIME",
          "x": "CRITICAL"
        },
        {
          "a": "DOWNTIME",
          "b": "DOWNTIME",
          "x": "DOWNTIME"
        }
      ]
    },
    {
      "name": "OR",
      "truth_table": [
        {
          "a": "OK",
          "b": "OK",
          "x": "OK"
        },
        {
          "a": "OK",
          "b": "WARNING",
          "x": "OK"
        },
        {
          "a": "OK",
          "b": "UNKNOWN",
          "x": "OK"
        },
        {
          "a": "OK",
          "b": "MISSING",
          "x": "OK"
        },
        {
          "a": "OK",
          "b": "CRITICAL",
          "x": "OK"
        },
        {
          "a": "OK",
          "b": "DOWNTIME",
          "x": "OK"
        },
        {
          "a": "WARNING",
          "b": "WARNING",
          "x": "WARNING"
        },
        {
          "a": "WARNING",
          "b": "UNKNOWN",
          "x": "WARNING"
        },
        {
          "a": "WARNING",
          "b": "MISSING",
          "x": "WARNING"
        },
        {
          "a": "WARNING",
          "b": "CRITICAL",
          "x": "WARNING"
        },
        {
          "a": "WARNING",
          "b": "DOWNTIME",
          "x": "WARNING"
        },
        {
          "a": "UNKNOWN",
          "b": "UNKNOWN",
          "x": "UNKNOWN"
        },
        {
          "a": "UNKNOWN",
          "b": "MISSING",
          "x": "UNKNOWN"
        },
        {
          "a": "UNKNOWN",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "UNKNOWN",
          "b": "DOWNTIME",
          "x": "UNKNOWN"
        },
        {
          "a": "MISSING",
          "b": "MISSING",
          "x": "MISSING"
        },
        {
          "a": "MISSING",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "MISSING",
          "b": "DOWNTIME",
          "x": "DOWNTIME"
        },
        {
          "a": "CRITICAL",
          "b": "CRITICAL",
          "x": "CRITICAL"
        },
        {
          "a": "CRITICAL",
          "b": "DOWNTIME",
          "x": "CRITICAL"
        },
        {
          "a": "DOWNTIME",
          "b": "DOWNTIME",
          "x": "DOWNTIME"
        }
      ]
    }
  ]
}
])

db.metric_profiles.insertMany(
  [
  {
    "id": "90191168-7690-4967-8769-11ccddb51b1e",
    "date": "2015-01-01",
    "date_integer": 20150101,
    "name": "MON_CLOUD",
    "description": "cloud checks",
    "services": [
      {
        "service": "compute",
        "metrics": [
          "create_vm_check",
          "resize_vm_check"
        ]
      },
      {
        "service": "storage",
        "metrics": [
          "list_files_check",
          "upload_file_check",
          "remove_file_check",
          "rename_file_check",
        ]
      },
      {
        "service": "aai",
        "metrics": [
          "http_check",
          "cert_check",
          "redirect_check",
          "login_check"
        ]
      },
       {
        "service": "messaging",
        "metrics": [
          "publish_check",
          "consume_check",
          "empty_queue_check",
          "create_topic_check",
          "delete_topic_check"
        ]
      }
    ]
  },
  {
    "id": "ce137987-ad50-49a9-8a37-77f4cde64b90",
    "date": "2015-01-01",
    "date_integer": 20150101,
    "name": "MON_CLOUD_LIGHT",
    "description": "just basic checks",
     "services": [
      {
        "service": "compute",
        "metrics": [
          "create_vm_check",
        ]
      },
      {
        "service": "storage",
        "metrics": [
          "list_files_check",
        ]
      },
      {
        "service": "aai",
        "metrics": [
          "http_check",
          "cert_check",
        ]
      },
        {
        "service": "messaging",
        "metrics": [
          "publish_check",
          "consume_check",
        ]
      }
    ]
  }
]
)

db.aggregation_profiles.insertMany(
[
  {
    "id": "e7ef618d-b2e3-454b-85f2-1dd0a2849ae2",
    "date": "2015-01-01",
    "date_integer": 20150101,
    "name": "AGGR_HA",
    "namespace": "tenant",
    "endpoint_group": "SERVICEGROUPS",
    "metric_operation": "AND",
    "profile_operation": "AND",
    "metric_profile": {
      "name": "MON_CLOUD",
      "id": "90191168-7690-4967-8769-11ccddb51b1e"
    },
    "groups": [
      {
        "name": "cloud",
        "operation": "AND",
        "services": [
          {
            "name": "compute",
            "operation": "OR"
          },
          {
            "name": "storage",
            "operation": "OR"
          },
          {
            "name": "iam",
            "operation": "OR"
          },
           {
            "name": "messaging",
            "operation": "OR"
          }
        ]
      }
    ]
  }
]
)

db.reports.insertMany(
  [
  {
    "id": "16b2b932-1cf6-42dc-8ce2-1e29bc6879b8",
    "tenant": "",
    "disabled": false,
    "info": {
      "name": "CORE",
      "description": "Core A/R report",
      "created": "2015-01-01 12:53:27",
      "updated": "2015-01-01 17:56:33"
    },
    "computations": {
      "ar": true,
      "status": true,
      "trends": [
        "flapping",
        "status",
        "tags"
      ]
    },
    "thresholds": {
      "availability": 80,
      "reliability": 90,
      "uptime": 0.800000011920929,
      "unknown": 0.10000000149011612,
      "downtime": 0.10000000149011612
    },
    "topology_schema": {
      "group": {
        "type": "PROJECT",
        "group": {
          "type": "SERVICEGROUPS"
        }
      }
    },
    "profiles": [
      {
        "id": "90191168-7690-4967-8769-11ccddb51b1e",
        "name": "MON_CLOUD",
        "type": "metric"
      },
      {
        "id": "e7ef618d-b2e3-454b-85f2-1dd0a2849ae2",
        "name": "AGGR_HA",
        "type": "aggregation"
      },
      {
        "id": "b131b9ac-07ef-442b-a20a-e18faf534f26",
        "name": "default_ops",
        "type": "operations"
      }
    ],
    "filter_tags": [],
    "node_report": true
  },
  {
    "id": "c7a6b0d4-4885-46da-9dd1-1f91d0e9142e",
    "tenant": "",
    "disabled": false,
    "info": {
      "name": "CLOUD-B",
      "description": "Report only for cloud b",
      "created": "2015-01-01 12:53:27",
      "updated": "2015-01-01 17:56:33"
    },
    "computations": {
      "ar": true,
      "status": true,
      "trends": [
        "flapping",
        "status",
        "tags"
      ]
    },
    "thresholds": {
      "availability": 80,
      "reliability": 90,
      "uptime": 0.800000011920929,
      "unknown": 0.10000000149011612,
      "downtime": 0.10000000149011612
    },
    "topology_schema": {
      "group": {
        "type": "PROJECT",
        "group": {
          "type": "SERVICEGROUPS"
        }
      }
    },
    "profiles": [
      {
        "id": "90191168-7690-4967-8769-11ccddb51b1e",
        "name": "MON_CLOUD",
        "type": "metric"
      },
      {
        "id": "e7ef618d-b2e3-454b-85f2-1dd0a2849ae2",
        "name": "AGGR_HA",
        "type": "aggregation"
      },
      {
        "id": "b131b9ac-07ef-442b-a20a-e18faf534f26",
        "name": "default_ops",
        "type": "operations"
      }
    ],
    "filter_tags": [
      {
          "name": "subgroup",
          "value": "CLOUD-B",
          "context": "argo.group.filter.fields"
      },
    ],
    "node_report": true
  }
]
)
