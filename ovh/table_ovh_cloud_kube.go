package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func tableOvhCloudKube() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_kube",
		Description: "A Managed Kubernetes Service (MKS) cluster.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listKube,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getKube,
		},
		Columns: []*plugin.Column{
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("project_id"),
				Description: "Project ID.",
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "ID of the cluster.",
			},
			{
				Name:        "name",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Name of the cluster.",
			},
			{
				Name:        "region",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Region of the cluster.",
			},
			{
				Name:        "status",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Status of the cluster, such as INSTALLING, READY, UPDATING, ERROR or USER_ERROR.",
			},
			{
				Name:        "version",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Kubernetes version of the cluster.",
			},
			{
				Name:        "next_upgrade_versions",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("NextUpgradeVersions"),
				Description: "Kubernetes versions available for upgrade.",
			},
			{
				Name:        "plan",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Plan of the cluster, such as free or standard.",
			},
			{
				Name:        "update_policy",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Update policy of the cluster: ALWAYS_UPDATE, MINIMAL_DOWNTIME or NEVER_UPDATE.",
			},
			{
				Name:        "is_up_to_date",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("IsUpToDate"),
				Description: "True if all nodes and the control plane are up to date.",
			},
			{
				Name:        "control_plane_is_up_to_date",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ControlPlaneIsUpToDate"),
				Description: "True if the control plane is up to date.",
			},
			{
				Name:        "url",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Management URL of the cluster, i.e. its API server endpoint.",
			},
			{
				Name:        "nodes_url",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_STRING,
				Description: "URL of the cluster nodes.",
			},
			{
				Name:        "kube_proxy_mode",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_STRING,
				Description: "Mode of kube-proxy: iptables or ipvs.",
			},
			{
				Name:        "private_network_id",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_STRING,
				Description: "OpenStack ID of the private network used by the cluster, if any.",
			},
			{
				Name:        "nodes_subnet_id",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_STRING,
				Description: "OpenStack ID of the subnet used by the cluster nodes, if any.",
			},
			{
				Name:        "load_balancers_subnet_id",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_STRING,
				Description: "OpenStack ID of the subnet used by the cluster load balancers, if any.",
			},
			{
				Name:        "private_network_configuration",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_JSON,
				Description: "Private network configuration of the cluster.",
			},
			{
				Name:        "disable_public_networking",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("DisablePublicNetworking"),
				Description: "True if the public connectivity of the cluster is disabled.",
			},
			{
				Name:        "ip_allocation_policy",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_JSON,
				Description: "IP allocation policy of the pods and services of the cluster.",
			},
			{
				Name:        "customization",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_JSON,
				Description: "Customization of the cluster, such as the admission plugins of its API server.",
			},
			{
				Name:        "audit_logs_subscribed",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("AuditLogsSubscribed"),
				Description: "True if the forwarding of the cluster audit logs has been subscribed.",
			},
			{
				Name:        "created_at",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and timestamp when the cluster was created.",
			},
			{
				Name:        "updated_at",
				Hydrate:     getKubeInfo,
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and timestamp when the cluster was last updated.",
			},
		},
	}
}

type Kube struct {
	AuditLogsSubscribed         bool        `json:"auditLogsSubscribed"`
	ControlPlaneIsUpToDate      bool        `json:"controlPlaneIsUpToDate"`
	CreatedAt                   time.Time   `json:"createdAt"`
	Customization               interface{} `json:"customization"`
	DisablePublicNetworking     *bool       `json:"disablePublicNetworking"`
	ID                          string      `json:"id"`
	IPAllocationPolicy          interface{} `json:"ipAllocationPolicy"`
	IsUpToDate                  bool        `json:"isUpToDate"`
	KubeProxyMode               string      `json:"kubeProxyMode"`
	LoadBalancersSubnetID       string      `json:"loadBalancersSubnetId"`
	Name                        string      `json:"name"`
	NextUpgradeVersions         []string    `json:"nextUpgradeVersions"`
	NodesSubnetID               string      `json:"nodesSubnetId"`
	NodesURL                    string      `json:"nodesUrl"`
	Plan                        string      `json:"plan"`
	PrivateNetworkConfiguration interface{} `json:"privateNetworkConfiguration"`
	PrivateNetworkID            string      `json:"privateNetworkId"`
	Region                      string      `json:"region"`
	Status                      string      `json:"status"`
	UpdatePolicy                string      `json:"updatePolicy"`
	UpdatedAt                   time.Time   `json:"updatedAt"`
	URL                         string      `json:"url"`
	Version                     string      `json:"version"`
}

func getKubeInfo(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	kube := h.Item.(Kube)
	projectId := d.EqualsQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_kube.getKubeInfo", "connection_error", err)
		return nil, err
	}

	err = client.Get(fmt.Sprintf("/cloud/project/%s/kube/%s", projectId, kube.ID), &kube)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_kube.getKubeInfo", "api_error", err, "id", kube.ID)
		return nil, err
	}
	return kube, nil
}

func listKube(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_kube.listKube", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()

	var kubeIds []string
	err = client.Get(fmt.Sprintf("/cloud/project/%s/kube", projectId), &kubeIds)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_kube.listKube", "api_error", err)
		return nil, err
	}
	for _, kubeId := range kubeIds {
		d.StreamListItem(ctx, Kube{ID: kubeId})
	}
	return nil, nil
}

func getKube(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	return Kube{ID: d.EqualsQuals["id"].GetStringValue()}, nil
}
