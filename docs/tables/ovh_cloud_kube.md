# Table: ovh_cloud_kube

A Managed Kubernetes Service (MKS) cluster.

The `ovh_cloud_kube` table can be used to query information about Managed Kubernetes Service clusters. You must specify which cloud project in the where clause (`where project_id=xxxx`).

## Examples

### List the clusters of a cloud project

```sql
select
  name,
  region,
  status,
  version,
  plan
from
  ovh_cloud_kube
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
```

### Get a specific cluster

```sql
select
  name,
  region,
  url,
  kube_proxy_mode,
  update_policy
from
  ovh_cloud_kube
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and id='0190a1b2-c3d4-7e5f-8a9b-0c1d2e3f4a5b'
```

### List the clusters that are not up to date

```sql
select
  name,
  version,
  update_policy,
  control_plane_is_up_to_date,
  next_upgrade_versions
from
  ovh_cloud_kube
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and not is_up_to_date
```

### List the clusters attached to a private network

```sql
select
  name,
  private_network_id,
  nodes_subnet_id,
  load_balancers_subnet_id
from
  ovh_cloud_kube
where
  project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
  and private_network_id is not null
```

### Count the nodes of each cluster

```sql
select
  k.name,
  sum(n.current_nodes) as nodes
from
  ovh_cloud_kube as k
  join ovh_cloud_kube_nodepool as n on n.project_id = k.project_id and n.kube_id = k.id
where
  k.project_id='27c5a6d3dfez87893jfd88fdsfmvnqb8'
group by
  k.name
```
