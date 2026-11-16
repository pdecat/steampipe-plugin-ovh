# Table: ovh_cloud_managed_kubernetes

A Managed Kubernetes Service (MKS) cluster.

The `ovh_cloud_managed_kubernetes` table can be used to query information about Managed Kubernetes Service clusters. You must specify which cloud project in the where clause (`where project_id=xxxx`).

**Note:** This table reads the OVHcloud API v2 endpoint `/v2/publicCloud/project/{projectId}/managedKubernetes`, which returns the current state of each cluster in a single paginated listing.

## Examples

### List the clusters of a cloud project

```sql
select
  name,
  region,
  plan,
  version,
  deployed_version,
  resource_status
from
  ovh_cloud_managed_kubernetes
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```

### Get a specific cluster

```sql
select
  name,
  region,
  api_server -> 'access' -> 'publicEndpoint' ->> 'url' as endpoint,
  cni ->> 'type' as cni
from
  ovh_cloud_managed_kubernetes
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and id='0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b'
```

### List the clusters whose API server is reachable from any IP

```sql
select
  name,
  region
from
  ovh_cloud_managed_kubernetes
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and coalesce(jsonb_array_length(api_server -> 'access' -> 'authorizedIps' -> 'rules'), 0) = 0
```

### List the clusters without automatic patch upgrades

```sql
select
  name,
  deployed_version
from
  ovh_cloud_managed_kubernetes
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and not coalesce((auto_upgrade -> 'patch' ->> 'enabled')::boolean, false)
```

### List the clusters attached to a vRack private network

```sql
select
  name,
  network -> 'vrack' ->> 'privateNetworkId' as private_network_id,
  network -> 'vrack' ->> 'subnetId' as subnet_id,
  network -> 'vrack' -> 'defaultEgress' ->> 'type' as default_egress
from
  ovh_cloud_managed_kubernetes
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and network ->> 'type' = 'VRACK'
```

### Count the nodes of each cluster

```sql
select
  k.name,
  sum(n.current_nodes) as nodes
from
  ovh_cloud_managed_kubernetes as k
  join ovh_cloud_managed_kubernetes_nodepool as n on n.project_id = k.project_id and n.managed_kubernetes_id = k.id
where
  k.project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
group by
  k.name
```
