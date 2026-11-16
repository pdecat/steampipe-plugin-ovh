package ovh

import (
	"context"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"

	"github.com/ovh/go-ovh/ovh"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"gopkg.in/yaml.v2"
)

// managedKubernetesKubeconfig is the response of the admin kubeconfig endpoint. Its content holds the
// cluster admin credentials: it must never be returned, logged or wrapped into an error, only the
// certificate authority is read out of it.
type managedKubernetesKubeconfig struct {
	Content string `json:"content"`
}

// getManagedKubernetesCertificateAuthority returns the certificate-authority-data of a cluster. The
// OVHcloud API exposes the cluster certificate authority nowhere but in the admin kubeconfig, so the
// kubeconfig is fetched and everything but the certificate authority is dropped.
func getManagedKubernetesCertificateAuthority(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	cluster := h.Item.(ManagedKubernetes)
	projectId := d.EqualsQuals["project_id"].GetStringValue()

	client, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_managed_kubernetes.getManagedKubernetesCertificateAuthority", "connection_error", err)
		return nil, err
	}

	var kubeconfig managedKubernetesKubeconfig
	path := fmt.Sprintf("/v2/publicCloud/project/%s/managedKubernetes/%s/kubeconfig", projectId, cluster.ID)
	if err = client.GetWithContext(ctx, path, &kubeconfig); err != nil {
		plugin.Logger(ctx).Error("ovh_cloud_managed_kubernetes.getManagedKubernetesCertificateAuthority", "api_error", err, "id", cluster.ID)
		return nil, err
	}

	certificateAuthority, err := certificateAuthorityFromKubeconfig(kubeconfig.Content)
	kubeconfig.Content = ""
	if err != nil {
		plugin.Logger(ctx).Warn("ovh_cloud_managed_kubernetes.getManagedKubernetesCertificateAuthority", "error", err, "id", cluster.ID)
		return nil, nil
	}
	return certificateAuthority, nil
}

// certificateAuthorityFromKubeconfig returns the base64-encoded certificate-authority-data of the first
// cluster of a kubeconfig, after checking it decodes to a PEM certificate. Its errors never quote the
// kubeconfig, which holds credentials.
func certificateAuthorityFromKubeconfig(content string) (string, error) {
	var kubeconfig struct {
		Clusters []struct {
			Cluster struct {
				CertificateAuthorityData string `yaml:"certificate-authority-data"`
			} `yaml:"cluster"`
		} `yaml:"clusters"`
	}
	if yaml.Unmarshal([]byte(content), &kubeconfig) != nil {
		return "", errors.New("the kubeconfig is not valid YAML")
	}
	if len(kubeconfig.Clusters) == 0 || kubeconfig.Clusters[0].Cluster.CertificateAuthorityData == "" {
		return "", errors.New("the kubeconfig has no certificate-authority-data")
	}
	data := kubeconfig.Clusters[0].Cluster.CertificateAuthorityData
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.New("the certificate-authority-data is not valid base64")
	}
	if block, _ := pem.Decode(decoded); block == nil || block.Type != "CERTIFICATE" {
		return "", errors.New("the certificate-authority-data is not a PEM certificate")
	}
	return data, nil
}

// shouldIgnoreKubeconfigError reports whether reading the admin kubeconfig was refused (403) or found
// nothing (404): credentials limited to read-only rules may not read it, which must only leave the
// certificate_authority column null instead of failing the whole query.
func shouldIgnoreKubeconfigError(_ context.Context, _ *plugin.QueryData, _ *plugin.HydrateData, err error) bool {
	var apiError *ovh.APIError
	if errors.As(err, &apiError) {
		return apiError.Code == http.StatusForbidden || apiError.Code == http.StatusNotFound
	}
	return false
}
