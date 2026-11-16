package ovh

import (
	"encoding/base64"
	"strings"
	"testing"
)

// testCertificateAuthority is a throwaway self-signed certificate (CN=test-ca) whose key was discarded.
const testCertificateAuthority = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUJlVENDQVIrZ0F3SUJBZ0lVZC9VbUhNVmRlWUhTd0MxRUwwVnhoNGZaTlZ3d0NnWUlLb1pJemowRUF3SXcKRWpFUU1BNEdBMVVFQXd3SGRHVnpkQzFqWVRBZUZ3MHlOakV4TVRZeE9EVTFORGxhRncwek5qRXhNVE14T0RVMQpORGxhTUJJeEVEQU9CZ05WQkFNTUIzUmxjM1F0WTJFd1dUQVRCZ2NxaGtqT1BRSUJCZ2dxaGtqT1BRTUJCd05DCkFBVDE1bzRwSnZENW8vZ1lPNFVwOFp0eU5KUWxPL2NSK3Z2aVFpV3ZIR29rRzBybDUwUG9yZkFwZXRVM0hWUVQKT05DOHh0c0E5dkJ4ODFraVdKRkJ0d3N2bzFNd1VUQWRCZ05WSFE0RUZnUVVHbUp1bGJJbHA4cXNpVyt2UWczQQo0K0Fob2F3d0h3WURWUjBqQkJnd0ZvQVVHbUp1bGJJbHA4cXNpVyt2UWczQTQrQWhvYXd3RHdZRFZSMFRBUUgvCkJBVXdBd0VCL3pBS0JnZ3Foa2pPUFFRREFnTklBREJGQWlBS1lLb3pnQ1FsSFYwdExSUGdVVGJCeHJ5WElES0wKNU1WRHJrWlllejNZTGdJaEFQVm8vUC9TbDlId3krUjhYWjgxV1ZldFdRUEo4RlFlUkdkVUpJV1QvSFdXCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"

func testKubeconfig(certificateAuthorityData string) string {
	return `apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: ` + certificateAuthorityData + `
    server: https://abcdef.c1.gra11.k8s.ovh.net
  name: test
contexts:
- context:
    cluster: test
    user: kubernetes-admin-test
  name: kubernetes-admin@test
current-context: kubernetes-admin@test
users:
- name: kubernetes-admin-test
  user:
    client-certificate-data: client-certificate-placeholder
    client-key-data: client-key-placeholder
`
}

func TestCertificateAuthorityFromKubeconfig(t *testing.T) {
	got, err := certificateAuthorityFromKubeconfig(testKubeconfig(testCertificateAuthority))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != testCertificateAuthority {
		t.Fatalf("got %q, want the certificate-authority-data", got)
	}
}

func TestCertificateAuthorityFromKubeconfigRejects(t *testing.T) {
	notACertificate := base64.StdEncoding.EncodeToString([]byte("-----BEGIN PRIVATE KEY-----\nAAAA\n-----END PRIVATE KEY-----\n"))
	for name, content := range map[string]string{
		"not yaml":          "clusters: [unterminated",
		"no clusters":       "apiVersion: v1\nkind: Config\nclusters: []\n",
		"no ca data":        strings.Replace(testKubeconfig(testCertificateAuthority), "certificate-authority-data", "insecure-skip-tls-verify", 1),
		"not base64":        testKubeconfig("not-base64!"),
		"not a certificate": testKubeconfig(notACertificate),
	} {
		t.Run(name, func(t *testing.T) {
			got, err := certificateAuthorityFromKubeconfig(content)
			if err == nil {
				t.Fatalf("expected an error, got %q", got)
			}
			// the kubeconfig holds credentials: an error must never quote it
			if strings.Contains(err.Error(), "placeholder") || strings.Contains(err.Error(), "clusters") {
				t.Fatalf("error quotes the kubeconfig: %v", err)
			}
		})
	}
}
