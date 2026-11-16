# Table: ovh_cloud_managed_kubernetes_nodepool

A node pool of a Managed Kubernetes Service (MKS) cluster.

The `ovh_cloud_managed_kubernetes_nodepool` table can be used to query information about the node pools of Managed Kubernetes Service clusters. You must specify which cloud project and which cluster in the where clause (`where project_id=xxxx and managed_kubernetes_id=xxxx`).

**Note:** This table reads the OVHcloud API v2 endpoint `/v2/publicCloud/project/{projectId}/managedKubernetes/{managedKubernetesId}/nodepool`.

## Examples

### List the node pools of a cluster

```sql
select
  name,
  flavor,
  availability_zone,
  current_nodes,
  ready_nodes
from
  ovh_cloud_managed_kubernetes_nodepool
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and managed_kubernetes_id='0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b'
```

### List the autoscaled node pools of a cluster

```sql
select
  name,
  desired_nodes,
  autoscaler ->> 'minNodes' as min_nodes,
  autoscaler ->> 'maxNodes' as max_nodes
from
  ovh_cloud_managed_kubernetes_nodepool
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and managed_kubernetes_id='0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b'
  and (autoscaler ->> 'enabled')::boolean
```

### List the node pools with nodes that are not ready or not up to date

```sql
select
  name,
  current_nodes,
  ready_nodes,
  up_to_date_nodes
from
  ovh_cloud_managed_kubernetes_nodepool
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and managed_kubernetes_id='0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b'
  and (ready_nodes < current_nodes or up_to_date_nodes < current_nodes)
```

### List the node pools of every cluster of a cloud project

```sql
select
  k.name as cluster,
  n.name,
  n.flavor,
  n.billing_type,
  n.current_nodes
from
  ovh_cloud_managed_kubernetes as k
  join ovh_cloud_managed_kubernetes_nodepool as n on n.project_id = k.project_id and n.managed_kubernetes_id = k.id
where
  k.project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```
