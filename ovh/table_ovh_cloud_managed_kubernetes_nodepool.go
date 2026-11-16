package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func tableOvhCloudManagedKubernetesNodePool() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_managed_kubernetes_nodepool",
		Description: "A node pool of a Managed Kubernetes Service (MKS) cluster.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "managed_kubernetes_id"}),
			Hydrate:    listManagedKubernetesNodePool,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "managed_kubernetes_id", "id"}),
			Hydrate:    getManagedKubernetesNodePool,
		},
		Columns: []*plugin.Column{
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("project_id"),
				Description: "Project ID.",
			},
			{
				Name:        "managed_kubernetes_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("managed_kubernetes_id"),
				Description: "ID of the cluster the node pool belongs to.",
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the node pool.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CurrentState.Name").NullIfZero(),
				Description: "Name of the node pool.",
			},
			{
				Name:        "resource_status",
				Type:        proto.ColumnType_STRING,
				Description: "Status of the node pool resource: CREATING, CREATION_PLANNED, DELETING, ERROR, OUT_OF_SYNC, READY, SUSPENDED, UNKNOWN or UPDATING.",
			},
			{
				Name:        "flavor",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CurrentState.NodeConfig.Flavor").NullIfZero(),
				Description: "Flavor of the nodes of the node pool.",
			},
			{
				Name:        "availability_zone",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CurrentState.AvailabilityZone").NullIfZero(),
				Description: "Availability zone of the node pool, in a multi-zone region.",
			},
			{
				Name:        "billing_type",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CurrentState.BillingType").NullIfZero(),
				Description: "Billing type of the nodes of the node pool: HOURLY or MONTHLY.",
			},
			{
				Name:        "desired_nodes",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CurrentState.DesiredNodes"),
				Description: "Number of nodes the node pool should have.",
			},
			{
				Name:        "current_nodes",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CurrentState.Status.CurrentNodes"),
				Description: "Number of nodes the node pool currently has.",
			},
			{
				Name:        "ready_nodes",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CurrentState.Status.ReadyNodes"),
				Description: "Number of nodes of the node pool that are ready.",
			},
			{
				Name:        "up_to_date_nodes",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CurrentState.Status.UpToDateNodes"),
				Description: "Number of nodes of the node pool that are up to date.",
			},
			{
				Name:        "autoscaler",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CurrentState.Autoscaler"),
				Description: "Autoscaler settings of the node pool: enabled, minimum and maximum number of nodes, scale down parameters.",
			},
			{
				Name:        "anti_affinity",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CurrentState.AntiAffinity"),
				Description: "Anti-affinity settings of the node pool, which spread its nodes over distinct hypervisors.",
			},
			{
				Name:        "attach_floating_ip",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CurrentState.NodeConfig.AttachFloatingIP"),
				Description: "Floating IP settings of the nodes of the node pool.",
			},
			{
				Name:        "rollout",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CurrentState.Rollout"),
				Description: "Rollout strategy of the node pool upgrades.",
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CurrentState.Template"),
				Description: "Template applied to the nodes of the node pool: labels, annotations, finalizers, taints and schedulability.",
			},
			{
				Name:        "current_tasks",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CurrentTasks"),
				Description: "Tasks currently running on the node pool.",
			},
			{
				Name:        "target_spec",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TargetSpec"),
				Description: "Target specification of the node pool, which its current state converges to.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and timestamp when the node pool was created.",
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and timestamp when the node pool was last updated.",
			},
		},
	}
}

type ManagedKubernetesNodePool struct {
	CreatedAt      time.Time                             `json:"createdAt"`
	CurrentState   ManagedKubernetesNodePoolCurrentState `json:"currentState"`
	CurrentTasks   interface{}                           `json:"currentTasks"`
	ID             string                                `json:"id"`
	ResourceStatus string                                `json:"resourceStatus"`
	TargetSpec     interface{}                           `json:"targetSpec"`
	UpdatedAt      time.Time                             `json:"updatedAt"`
}

// The node counts are pointers so that an absent count stays NULL while a real 0, such as a
// node pool scaled down to zero by the autoscaler, is kept as 0.
type ManagedKubernetesNodePoolCurrentState struct {
	AntiAffinity     interface{}                     `json:"antiAffinity"`
	Autoscaler       interface{}                     `json:"autoscaler"`
	AvailabilityZone string                          `json:"availabilityZone"`
	BillingType      string                          `json:"billingType"`
	DesiredNodes     *int64                          `json:"desiredNodes"`
	Name             string                          `json:"name"`
	NodeConfig       ManagedKubernetesNodeConfig     `json:"nodeConfig"`
	Rollout          interface{}                     `json:"rollout"`
	Status           ManagedKubernetesNodePoolStatus `json:"status"`
	Template         interface{}                     `json:"template"`
}

type ManagedKubernetesNodeConfig struct {
	AttachFloatingIP interface{} `json:"attachFloatingIp"`
	Flavor           string      `json:"flavor"`
}

type ManagedKubernetesNodePoolStatus struct {
	CurrentNodes  *int64 `json:"currentNodes"`
	ReadyNodes    *int64 `json:"readyNodes"`
	UpToDateNodes *int64 `json:"upToDateNodes"`
}

func listManagedKubernetesNodePool(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_managed_kubernetes_nodepool.listManagedKubernetesNodePool", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	clusterId := d.EqualsQuals["managed_kubernetes_id"].GetStringValue()

	path := fmt.Sprintf("/v2/publicCloud/project/%s/managedKubernetes/%s/nodepool", projectId, clusterId)
	err = listAllV2(ctx, client, path, func(nodePool ManagedKubernetesNodePool) {
		d.StreamListItem(ctx, nodePool)
	})
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_managed_kubernetes_nodepool.listManagedKubernetesNodePool", "api_error", err, "managed_kubernetes_id", clusterId)
		return nil, err
	}
	return nil, nil
}

func getManagedKubernetesNodePool(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_managed_kubernetes_nodepool.getManagedKubernetesNodePool", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	clusterId := d.EqualsQuals["managed_kubernetes_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()

	var nodePool ManagedKubernetesNodePool
	path := fmt.Sprintf("/v2/publicCloud/project/%s/managedKubernetes/%s/nodepool/%s", projectId, clusterId, id)
	err = client.GetWithContext(ctx, path, &nodePool)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_managed_kubernetes_nodepool.getManagedKubernetesNodePool", "api_error", err, "managed_kubernetes_id", clusterId, "id", id)
		return nil, err
	}
	return nodePool, nil
}
