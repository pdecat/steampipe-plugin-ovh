# Table: ovh_cloud_kube_nodepool

A node pool of a Managed Kubernetes Service (MKS) cluster.

The `ovh_cloud_kube_nodepool` table can be used to query information about the node pools of Managed Kubernetes Service clusters. You must specify which cloud project and which cluster in the where clause (`where project_id=xxxx and kube_id=xxxx`).

## Examples

### List the node pools of a cluster

```sql
select
  name,
  flavor,
  status,
  current_nodes,
  min_nodes,
  max_nodes
from
  ovh_cloud_kube_nodepool
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and kube_id='0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b'
```

### List the autoscaled node pools of a cluster

```sql
select
  name,
  desired_nodes,
  min_nodes,
  max_nodes,
  autoscaling
from
  ovh_cloud_kube_nodepool
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and kube_id='0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b'
  and autoscale
```

### List the node pools of every cluster of a cloud project

```sql
select
  k.name as cluster,
  n.name,
  n.flavor,
  n.current_nodes,
  n.availability_zones
from
  ovh_cloud_kube as k
  join ovh_cloud_kube_nodepool as n on n.project_id = k.project_id and n.kube_id = k.id
where
  k.project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```
