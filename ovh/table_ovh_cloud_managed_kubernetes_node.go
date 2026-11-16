package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func tableOvhCloudManagedKubernetesNode() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_managed_kubernetes_node",
		Description: "A node of a Managed Kubernetes Service (MKS) cluster, with the ID of the public cloud instance running it.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "project_id", Require: plugin.Required},
				{Name: "managed_kubernetes_id", Require: plugin.Required},
				{Name: "nodepool_id", Require: plugin.Optional},
			},
			Hydrate: listManagedKubernetesNode,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "managed_kubernetes_id", "nodepool_id", "id"}),
			Hydrate:    getManagedKubernetesNode,
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
				Description: "ID of the cluster the node belongs to.",
			},
			{
				Name:        "nodepool_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CurrentState.NodePoolID").NullIfZero(),
				Description: "ID of the node pool the node belongs to.",
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the node.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CurrentState.Name").NullIfZero(),
				Description: "Name of the node, which is also the name of its instance.",
			},
			{
				Name:        "instance_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CurrentState.InstanceID").NullIfZero(),
				Description: "ID of the public cloud instance running the node, as listed by ovh_cloud_instance.",
			},
			{
				Name:        "flavor",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CurrentState.Flavor").NullIfZero(),
				Description: "Flavor of the node.",
			},
			{
				Name:        "version",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CurrentState.Version").NullIfZero(),
				Description: "Kubernetes version of the node.",
			},
			{
				Name:        "up_to_date",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("CurrentState.UpToDate"),
				Description: "True if the node is up to date.",
			},
			{
				Name:        "resource_status",
				Type:        proto.ColumnType_STRING,
				Description: "Status of the node resource: CREATING, CREATION_PLANNED, DELETING, ERROR, OUT_OF_SYNC, READY, SUSPENDED, UNKNOWN or UPDATING.",
			},
			{
				Name:        "deployed_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CurrentState.DeployedAt"),
				Description: "The date and timestamp when the node was deployed.",
			},
			{
				Name:        "deleted_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("CurrentState.DeletedAt"),
				Description: "The date and timestamp when the node was deleted, if it is being deleted.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and timestamp when the node was created.",
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and timestamp when the node was last updated.",
			},
		},
	}
}

type ManagedKubernetesNode struct {
	CreatedAt      time.Time                         `json:"createdAt"`
	CurrentState   ManagedKubernetesNodeCurrentState `json:"currentState"`
	ID             string                            `json:"id"`
	ResourceStatus string                            `json:"resourceStatus"`
	UpdatedAt      time.Time                         `json:"updatedAt"`
}

type ManagedKubernetesNodeCurrentState struct {
	DeletedAt  *time.Time `json:"deletedAt"`
	DeployedAt *time.Time `json:"deployedAt"`
	Flavor     string     `json:"flavor"`
	InstanceID string     `json:"instanceId"`
	Name       string     `json:"name"`
	NodePoolID string     `json:"nodePoolId"`
	UpToDate   *bool      `json:"upToDate"`
	Version    string     `json:"version"`
}

// listManagedKubernetesNode lists the nodes of one node pool when nodepool_id is given, and
// otherwise walks every node pool of the cluster, so that all the nodes of a cluster, and the
// instances running them, can be listed from the cluster ID alone.
func listManagedKubernetesNode(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_managed_kubernetes_node.listManagedKubernetesNode", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	clusterId := d.EqualsQuals["managed_kubernetes_id"].GetStringValue()
	clusterPath := fmt.Sprintf("/v2/publicCloud/project/%s/managedKubernetes/%s", projectId, clusterId)

	var nodePoolIds []string
	if nodePoolId := d.EqualsQualString("nodepool_id"); nodePoolId != "" {
		nodePoolIds = []string{nodePoolId}
	} else {
		err = listAllV2(ctx, client, clusterPath+"/nodepool", func(nodePool ManagedKubernetesNodePool) {
			nodePoolIds = append(nodePoolIds, nodePool.ID)
		})
		if err != nil {
			plugin.Logger(ctx).Error("ovh_cloud_managed_kubernetes_node.listManagedKubernetesNode", "api_error", err, "managed_kubernetes_id", clusterId)
			return nil, err
		}
	}

	for _, nodePoolId := range nodePoolIds {
		err = listAllV2(ctx, client, fmt.Sprintf("%s/nodepool/%s/node", clusterPath, nodePoolId), func(node ManagedKubernetesNode) {
			d.StreamListItem(ctx, node)
		})
		if err != nil {
			plugin.Logger(ctx).Error("ovh_cloud_managed_kubernetes_node.listManagedKubernetesNode", "api_error", err, "managed_kubernetes_id", clusterId, "nodepool_id", nodePoolId)
			return nil, err
		}
	}
	return nil, nil
}

func getManagedKubernetesNode(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_managed_kubernetes_node.getManagedKubernetesNode", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	clusterId := d.EqualsQuals["managed_kubernetes_id"].GetStringValue()
	nodePoolId := d.EqualsQuals["nodepool_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()

	var node ManagedKubernetesNode
	path := fmt.Sprintf("/v2/publicCloud/project/%s/managedKubernetes/%s/nodepool/%s/node/%s", projectId, clusterId, nodePoolId, id)
	err = client.GetWithContext(ctx, path, &node)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_managed_kubernetes_node.getManagedKubernetesNode", "api_error", err, "managed_kubernetes_id", clusterId, "nodepool_id", nodePoolId, "id", id)
		return nil, err
	}
	return node, nil
}
