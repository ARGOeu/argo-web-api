---
id: adminaccess
title: Accessing Tenant Results as Admin
---

## Tenant User Endpoints

Tenant users can access tenant-specific endpoints using their credentials via the `x-api-key` header.

Some tenant-related routes include the following:
- `/api/v2/operations_profiles`
- `/api/v2/metric_profiles`
- `/api/v2/aggregations_profiles`
- `/api/v2/reports`
- `/api/v2/topology`
- `/api/v2/status`
- `/api/v2/results`
- `/api/v2/issues`
- `/api/v2/trends`

### Authentication behavior
When you authenticate as a tenant user (with `admin`, `editor`, or `viewer` role), your `x-api-key` automatically determines your tenant context. The system selects your tenant during authentication.

To access a different tenant, you must obtain separate credentials for that tenant and use them in your requests.

## Super Admin Access

Super admins can access administrative endpoints under the `/api/v2/admin` path using their super-admin `x-api-key`.

Some super-admin related routes include the following:
- `/api/v2/tenants`
- `/api/v2/tenants/{tenant-id}/users`

These endpoints enable management of tenants, users, and related resources.

### Accessing tenant data as super admin
To view a specific tenant's results, profiles, and other data, use your super-admin credentials with the additional header `x-tenant-id: <tenant-uuid>` (replace `<tenant-uuid>` with the actual tenant ID).