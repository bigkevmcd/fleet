package metrics

import (
	"fmt"

	"github.com/rancher/fleet/internal/cmd/controller/summary"
	fleet "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	bundleDeploymentSubsystem = "bundledeployment"
	bundleDeploymentLabels    = []string{
		"name",
		"namespace",
		"cluster_name",
		"cluster_namespace",
		"repo",
		"commit",
		"bundle",
		"bundle_namespace",
		"generation",
		"state",
	}
	BundleDeploymentCollector = CollectorCollection{
		subsystem: bundleDeploymentSubsystem,
		metrics: map[string]prometheus.Collector{
			"state": promauto.NewGaugeVec(
				prometheus.GaugeOpts{
					Namespace: metricPrefix,
					Subsystem: bundleDeploymentSubsystem,
					Name:      "state",
					Help: "Shows the state of this bundle deployment based on the state label. " +
						"A value of 1 is true 0 is false.",
				},
				bundleDeploymentLabels,
			),
		},
		collector: collectBundleDeploymentMetrics,
	}
)

func collectBundleDeploymentMetrics(obj any, metrics map[string]prometheus.Collector) {
	bundleDep, ok := obj.(*fleet.BundleDeployment)
	if !ok {
		panic("unexpected object type")
	}

	currentState := summary.GetDeploymentState(bundleDep)
	labels := prometheus.Labels{
		"name":              bundleDep.Name,
		"namespace":         bundleDep.Namespace,
		"cluster_name":      bundleDep.Labels[fleet.ClusterLabel],
		"cluster_namespace": bundleDep.Labels[fleet.ClusterNamespaceLabel],
		"repo":              bundleDep.Labels[fleet.RepoLabel],
		"commit":            bundleDep.Labels[fleet.CommitLabel],
		"bundle":            bundleDep.Labels[fleet.BundleLabel],
		"bundle_namespace":  bundleDep.Labels[fleet.BundleNamespaceLabel],
		"generation":        fmt.Sprintf("%d", bundleDep.Generation),
		"state":             string(currentState),
	}

	for _, state := range bundleStates {
		labels["state"] = string(state)

		if state == currentState {
			metrics["state"].(*prometheus.GaugeVec).With(labels).Set(1)
		} else {
			metrics["state"].(*prometheus.GaugeVec).With(labels).Set(0)
		}
	}
}
