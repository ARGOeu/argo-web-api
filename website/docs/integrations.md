---
id: integrations
title: Integrations - Components
---

Calls for  components that integrate with argo-web-api and access tenant data.

## Refresh Access Token for a Component Integration

Component administrators can request to refresh access keys for their assigned components within a specific tenant.

```
POST /api/v3/integrations/components/{COMPONENT}/by-tenant-name/{TENANT}/refresh
```

Where `{COMPONENT}` the name of the component account supported (one of the following):
- engine
- poem-admin
- poem-viewer
- monbox
- probe

Where `{TENANT}` the name of the tenant

### Request headers

```
Accept: application/json
x-api-key: component-admin-s3cr3t
```

### Response
Headers: `Status: 200 OK`

### Response Body

Json Response example:
```json
{
    "status": {
        "message": "component api key succesfully renewed",
        "code": "200"
    },
    "data": {
        "api_key": "..."
    }
}
```

## Errors

Component administrators have access only to their assigned component type. For example, a `component_poem` admin can only access `poem` components.

Attempting to access an unauthorized component type returns:

```json
{
    "status": {
        "message": "Access to the resource is Forbidden",
        "code": "403"
    }
}
```

If a component access account has not yet been set up for a specific tenant, then the response will be:

```json
{
    "status": {
        "message": "Component poem-viewer not configured for tenant FOO",
        "code": "404"
    }
}

```

### Types of components

Types of component accounts supported:

- engine
- poem-admin
- poem-viewer
- monbox
- probe

