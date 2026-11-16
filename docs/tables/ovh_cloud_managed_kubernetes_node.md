# Table: ovh_cloud_managed_kubernetes_node

A node of a Managed Kubernetes Service (MKS) cluster, with the ID of the public cloud instance running it.

The `ovh_cloud_managed_kubernetes_node` table can be used to query information about the nodes of Managed Kubernetes Service clusters. You must specify which cloud project and which cluster in the where clause (`where project_id=xxxx and managed_kubernetes_id=xxxx`). The node pool can optionally be specified too (`and nodepool_id=xxxx`), otherwise every node pool of the cluster is queried.

**Note:** This table reads the OVHcloud API v2 endpoint `/v2/publicCloud/project/{projectId}/managedKubernetes/{managedKubernetesId}/nodepool/{nodePoolId}/node`.

**Note:** The nodes of a cluster run on public cloud instances of the same project, which `ovh_cloud_instance` lists like any other instance. The `instance_id` column is the way to tell them apart.

## Examples

### List the nodes of a cluster

```sql
select
  name,
  nodepool_id,
  flavor,
  version,
  up_to_date
from
  ovh_cloud_managed_kubernetes_node
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and managed_kubernetes_id='0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b'
```

### List the nodes of a node pool

```sql
select
  name,
  instance_id,
  deployed_at
from
  ovh_cloud_managed_kubernetes_node
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and managed_kubernetes_id='0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b'
  and nodepool_id='0190a1b3-d4e5-7f60-9a1b-2c3d4e5f6a7b'
```

### List the instances of a cloud project that are Kubernetes nodes

```sql
select
  i.name,
  i.flavor_id,
  k.name as cluster
from
  ovh_cloud_instance as i
  join ovh_cloud_managed_kubernetes_node as n on n.instance_id = i.id
  join ovh_cloud_managed_kubernetes as k on k.project_id = n.project_id and k.id = n.managed_kubernetes_id
where
  i.project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and n.project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and k.project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```
