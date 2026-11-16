package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func tableOvhCloudKubeNodePool() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_kube_nodepool",
		Description: "A node pool of a Managed Kubernetes Service (MKS) cluster.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "kube_id"}),
			Hydrate:    listKubeNodePool,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "kube_id", "id"}),
			Hydrate:    getKubeNodePool,
		},
		Columns: []*plugin.Column{
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("project_id"),
				Description: "Project ID.",
			},
			{
				Name:        "kube_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("kube_id"),
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
				Description: "Name of the node pool.",
			},
			{
				Name:        "flavor",
				Type:        proto.ColumnType_STRING,
				Description: "Flavor of the nodes of the node pool.",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Status of the node pool, such as INSTALLING, READY, UPSCALING, DOWNSCALING, ERROR or USER_ERROR.",
			},
			{
				Name:        "size_status",
				Type:        proto.ColumnType_STRING,
				Description: "Size status of the node pool: CAPACITY_OK, OVER_CAPACITY or UNDER_CAPACITY.",
			},
			{
				Name:        "autoscale",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("Autoscale"),
				Description: "True if the node pool is managed by the cluster autoscaler.",
			},
			{
				Name:        "autoscaling",
				Type:        proto.ColumnType_JSON,
				Description: "Autoscaling parameters of the node pool.",
			},
			{
				Name:        "monthly_billed",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("MonthlyBilled"),
				Description: "True if the nodes of the node pool are billed monthly rather than hourly.",
			},
			{
				Name:        "anti_affinity",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AntiAffinity"),
				Description: "True if the nodes of the node pool are spread over distinct hypervisors.",
			},
			{
				Name:        "desired_nodes",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("DesiredNodes"),
				Description: "Number of nodes the node pool should have.",
			},
			{
				Name:        "min_nodes",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MinNodes"),
				Description: "Minimum number of nodes of the node pool.",
			},
			{
				Name:        "max_nodes",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("MaxNodes"),
				Description: "Maximum number of nodes of the node pool.",
			},
			{
				Name:        "current_nodes",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("CurrentNodes"),
				Description: "Number of nodes the node pool currently has.",
			},
			{
				Name:        "available_nodes",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("AvailableNodes"),
				Description: "Number of nodes of the node pool that are ready.",
			},
			{
				Name:        "up_to_date_nodes",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("UpToDateNodes"),
				Description: "Number of nodes of the node pool that are up to date.",
			},
			{
				Name:        "availability_zones",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AvailabilityZones"),
				Description: "Availability zones of the node pool.",
			},
			{
				Name:        "attach_floating_ips",
				Type:        proto.ColumnType_JSON,
				Description: "Floating IP configuration of the nodes of the node pool.",
			},
			{
				Name:        "template",
				Type:        proto.ColumnType_JSON,
				Description: "Template applied to the nodes of the node pool: labels, annotations, finalizers, taints and schedulability.",
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

type KubeNodePool struct {
	AntiAffinity      bool        `json:"antiAffinity"`
	AttachFloatingIps interface{} `json:"attachFloatingIps"`
	Autoscale         bool        `json:"autoscale"`
	Autoscaling       interface{} `json:"autoscaling"`
	AvailabilityZones []string    `json:"availabilityZones"`
	AvailableNodes    int         `json:"availableNodes"`
	CreatedAt         time.Time   `json:"createdAt"`
	CurrentNodes      int         `json:"currentNodes"`
	DesiredNodes      int         `json:"desiredNodes"`
	Flavor            string      `json:"flavor"`
	ID                string      `json:"id"`
	MaxNodes          int         `json:"maxNodes"`
	MinNodes          int         `json:"minNodes"`
	MonthlyBilled     bool        `json:"monthlyBilled"`
	Name              string      `json:"name"`
	SizeStatus        string      `json:"sizeStatus"`
	Status            string      `json:"status"`
	Template          interface{} `json:"template"`
	UpToDateNodes     int         `json:"upToDateNodes"`
	UpdatedAt         time.Time   `json:"updatedAt"`
}

func listKubeNodePool(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_kube_nodepool.listKubeNodePool", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	kubeId := d.EqualsQuals["kube_id"].GetStringValue()

	var nodePools []KubeNodePool
	err = client.Get(fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool", projectId, kubeId), &nodePools)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_kube_nodepool.listKubeNodePool", "api_error", err, "kube_id", kubeId)
		return nil, err
	}
	for _, nodePool := range nodePools {
		d.StreamListItem(ctx, nodePool)
	}
	return nil, nil
}

func getKubeNodePool(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_kube_nodepool.getKubeNodePool", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	kubeId := d.EqualsQuals["kube_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()

	var nodePool KubeNodePool
	err = client.Get(fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", projectId, kubeId, id), &nodePool)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_kube_nodepool.getKubeNodePool", "api_error", err, "kube_id", kubeId, "id", id)
		return nil, err
	}
	return nodePool, nil
}
