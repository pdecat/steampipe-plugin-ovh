package ovh

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
)

func tableOvhCloudManagedKubernetes() *plugin.Table {
	return &plugin.Table{
		Name:        "ovh_cloud_managed_kubernetes",
		Description: "A Managed Kubernetes Service (MKS) cluster.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("project_id"),
			Hydrate:    listManagedKubernetes,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"project_id", "id"}),
			Hydrate:    getManagedKubernetes,
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getManagedKubernetesCertificateAuthority,
				// Reading the admin kubeconfig is a privileged call: credentials limited to
				// read-only rules may be refused it, which must only leave the column null.
				IgnoreConfig: &plugin.IgnoreConfig{ShouldIgnoreErrorFunc: shouldIgnoreKubeconfigError},
			},
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
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CurrentState.Name").NullIfZero(),
				Description: "Name of the cluster.",
			},
			{
				Name:        "region",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CurrentState.Region").NullIfZero(),
				Description: "Region of the cluster.",
			},
			{
				Name:        "plan",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CurrentState.Plan").NullIfZero(),
				Description: "Plan of the cluster: FREE or STANDARD.",
			},
			{
				Name:        "version",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CurrentState.Version").NullIfZero(),
				Description: "Kubernetes minor version of the cluster, such as 1.35.",
			},
			{
				Name:        "deployed_version",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CurrentState.DeployedVersion").NullIfZero(),
				Description: "Kubernetes version actually deployed on the cluster, including the OVHcloud patch revision, such as 1.35.2-9.",
			},
			{
				Name:        "resource_status",
				Type:        proto.ColumnType_STRING,
				Description: "Status of the cluster resource: CREATING, CREATION_PLANNED, DELETING, ERROR, OUT_OF_SYNC, READY, SUSPENDED, UNKNOWN or UPDATING.",
			},
			{
				Name:        "api_server",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CurrentState.APIServer"),
				Description: "API server of the cluster: public endpoint, authorized IPs and admission plugins.",
			},
			{
				Name:        "certificate_authority",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getManagedKubernetesCertificateAuthority,
				Transform:   transform.FromValue(),
				Description: "Base64-encoded PEM certificate of the cluster certificate authority, as the certificate-authority-data of a kubeconfig. Read from the admin kubeconfig, of which nothing else is exposed: null when the credentials may not read it.",
			},
			{
				Name:        "authentication",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CurrentState.Authentication"),
				Description: "Authentication of the cluster, such as an OIDC provider, if any.",
			},
			{
				Name:        "auto_upgrade",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CurrentState.AutoUpgrade"),
				Description: "Automatic upgrade settings of the cluster.",
			},
			{
				Name:        "cni",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CurrentState.CNI"),
				Description: "Container network interface of the cluster: CANAL or CILIUM, with its settings.",
			},
			{
				Name:        "ip_allocation_policy",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CurrentState.IPAllocationPolicy"),
				Description: "IP ranges of the pods and services of the cluster.",
			},
			{
				Name:        "network",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CurrentState.Network"),
				Description: "Network of the cluster: PUBLIC or VRACK, with its private network, subnets and default egress.",
			},
			{
				Name:        "current_tasks",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CurrentTasks"),
				Description: "Tasks currently running on the cluster.",
			},
			{
				Name:        "target_spec",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("TargetSpec"),
				Description: "Target specification of the cluster, which its current state converges to.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and timestamp when the cluster was created.",
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_TIMESTAMP,
				Description: "The date and timestamp when the cluster was last updated.",
			},
		},
	}
}

type ManagedKubernetes struct {
	CreatedAt      time.Time                     `json:"createdAt"`
	CurrentState   ManagedKubernetesCurrentState `json:"currentState"`
	CurrentTasks   interface{}                   `json:"currentTasks"`
	ID             string                        `json:"id"`
	ResourceStatus string                        `json:"resourceStatus"`
	TargetSpec     interface{}                   `json:"targetSpec"`
	UpdatedAt      time.Time                     `json:"updatedAt"`
}

type ManagedKubernetesCurrentState struct {
	APIServer          interface{} `json:"apiServer"`
	Authentication     interface{} `json:"authentication"`
	AutoUpgrade        interface{} `json:"autoUpgrade"`
	CNI                interface{} `json:"cni"`
	DeployedVersion    string      `json:"deployedVersion"`
	IPAllocationPolicy interface{} `json:"ipAllocationPolicy"`
	Name               string      `json:"name"`
	Network            interface{} `json:"network"`
	Plan               string      `json:"plan"`
	Region             string      `json:"region"`
	Version            string      `json:"version"`
}

func listManagedKubernetes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_managed_kubernetes.listManagedKubernetes", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()

	err = listAllV2(ctx, client, fmt.Sprintf("/v2/publicCloud/project/%s/managedKubernetes", projectId), func(cluster ManagedKubernetes) {
		d.StreamListItem(ctx, cluster)
	})
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_managed_kubernetes.listManagedKubernetes", "api_error", err)
		return nil, err
	}
	return nil, nil
}

func getManagedKubernetes(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_managed_kubernetes.getManagedKubernetes", "connection_error", err)
		return nil, err
	}
	projectId := d.EqualsQuals["project_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()

	var cluster ManagedKubernetes
	err = client.GetWithContext(ctx, fmt.Sprintf("/v2/publicCloud/project/%s/managedKubernetes/%s", projectId, id), &cluster)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_managed_kubernetes.getManagedKubernetes", "api_error", err, "id", id)
		return nil, err
	}
	return cluster, nil
}
