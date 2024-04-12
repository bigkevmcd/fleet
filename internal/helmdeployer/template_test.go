package helmdeployer

import (
	"context"
	"testing"

	"github.com/rancher/fleet/internal/manifest"
	fleet "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	fleetv1 "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestTemplate(t *testing.T) {
	testBundle := &fleetv1.Bundle{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "helm-single-cluster-helm",
			Namespace: "fleet-local",
		},
		Spec: fleetv1.BundleSpec{
			Resources: []fleet.BundleResource{
				{
					Name: "fleet.yaml",
					Content: `
namespace: fleet-helm-example
helm:
releaseName: guestbook
chart: ""
repo: ""
version: ""
force: false
timeoutSeconds: 0
values:
  replicas: 2
`,
				},
				{
					Name: "Chart.yaml",
					Content: `
apiVersion: v2
name: guestbook
description: Sample application
version: 0.0.0
appVersion: 0.0.0
`,
				},
				{
					Name: "templates/frontend-service.yaml",
					Content: `
apiVersion: v1
kind: Service
metadata:
  name: frontend
  labels:
    app: guestbook
    clusterVersion: {{ .Capabilities.KubeVersion.Version }},
spec:
  type: "{{ .Values.serviceType }}"
  ports:
  - port: 80
  selector:
    app: guestbook
    tier: frontend
`,
				},
			},
		},
	}

	objs, err := Template(context.TODO(), "test-bundle", manifest.FromBundle(testBundle), fleetv1.BundleDeploymentOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if l := len(objs); l != 1 {
		t.Errorf("Template() got %d resources, want 1", l)
	}
	svc := testToUnstructured(t, objs[0])
	want := map[string]string{
		"app":                          "guestbook",
		"clusterVersion":               "v1.20.0,",
		"objectset.rio.cattle.io/hash": "896d620caba176df539386eee887bcaa4f248e8c",
	}
	assert.Equal(t, want, svc.GetLabels())
}

func testToUnstructured(t *testing.T, obj runtime.Object) *unstructured.Unstructured {
	rawMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		t.Fatal(err)
	}

	return &unstructured.Unstructured{Object: rawMap}
}
