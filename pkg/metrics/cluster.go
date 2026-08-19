/*
Copyright 2022 The Karmada Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	"github.com/karmada-io/karmada/pkg/util"
	utilmetrics "github.com/karmada-io/karmada/pkg/util/metrics"
)

const (
	clusterReadyMetricsName              = "cluster_ready_state"
	clusterTotalNodeNumberMetricsName    = "cluster_node_number"
	clusterReadyNodeNumberMetricsName    = "cluster_ready_node_number"
	clusterMemoryAllocatableMetricsName  = "cluster_memory_allocatable_bytes"
	clusterCPUAllocatableMetricsName     = "cluster_cpu_allocatable_number"
	clusterPodAllocatableMetricsName     = "cluster_pod_allocatable_number"
	clusterMemoryAllocatedMetricsName    = "cluster_memory_allocated_bytes"
	clusterCPUAllocatedMetricsName       = "cluster_cpu_allocated_number"
	clusterPodAllocatedMetricsName       = "cluster_pod_allocated_number"
	clusterSyncStatusDurationMetricsName = "cluster_sync_status_duration_seconds"
	clusterHealthProbeSuccessName        = "cluster_health_probe_success"
	clusterHealthProbeDurationName       = "cluster_health_probe_duration_seconds"
	clusterHealthProbeTotalName          = "cluster_health_probe_total"
	clusterReadySinceName                = "cluster_ready_since_timestamp_seconds"
	clusterConditionLastTransitionName   = "cluster_condition_last_transition_timestamp_seconds"
	clusterHealthTransitionsTotalName    = "cluster_health_transitions_total"

	// Canonical label for Karmada member clusters.
	memberClusterLabel = "member_cluster"

	// HealthStateSuccess indicates the cluster is online and healthy.
	HealthStateSuccess = "success"
	// HealthStateError indicates the cluster is not healthy or not reachable.
	HealthStateError = "error"
)

var (
	// clusterReadyGauge reports if the cluster is ready.
	clusterReadyGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: clusterReadyMetricsName,
		Help: "State of the cluster(1 if ready, 0 otherwise).",
	}, []string{"cluster_name"})

	// clusterTotalNodeNumberGauge reports the number of nodes in the given cluster.
	clusterTotalNodeNumberGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: clusterTotalNodeNumberMetricsName,
		Help: "Number of nodes in the cluster.",
	}, []string{"cluster_name"})

	// clusterReadyNodeNumberGauge reports the number of ready nodes in the given cluster.
	clusterReadyNodeNumberGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: clusterReadyNodeNumberMetricsName,
		Help: "Number of ready nodes in the cluster.",
	}, []string{"cluster_name"})

	// clusterMemoryAllocatableGauge reports the allocatable memory in the given cluster.
	clusterMemoryAllocatableGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: clusterMemoryAllocatableMetricsName,
		Help: "Allocatable cluster memory resource in bytes.",
	}, []string{"cluster_name"})

	// clusterCPUAllocatableGauge reports the allocatable CPU in the given cluster.
	clusterCPUAllocatableGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: clusterCPUAllocatableMetricsName,
		Help: "Number of allocatable CPU in the cluster.",
	}, []string{"cluster_name"})

	// clusterPodAllocatableGauge reports the allocatable Pod number in the given cluster.
	clusterPodAllocatableGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: clusterPodAllocatableMetricsName,
		Help: "Number of allocatable pods in the cluster.",
	}, []string{"cluster_name"})

	// clusterMemoryAllocatedGauge reports the allocated memory in the given cluster.
	clusterMemoryAllocatedGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: clusterMemoryAllocatedMetricsName,
		Help: "Allocated cluster memory resource in bytes.",
	}, []string{"cluster_name"})

	// clusterCPUAllocatedGauge reports the allocated CPU in the given cluster.
	clusterCPUAllocatedGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: clusterCPUAllocatedMetricsName,
		Help: "Number of allocated CPU in the cluster.",
	}, []string{"cluster_name"})

	// clusterPodAllocatedGauge reports the allocated Pod number in the given cluster.
	clusterPodAllocatedGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: clusterPodAllocatedMetricsName,
		Help: "Number of allocated pods in the cluster.",
	}, []string{"cluster_name"})

	// clusterSyncStatusDuration reports the duration of the given cluster syncing status.
	clusterSyncStatusDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: clusterSyncStatusDurationMetricsName,
		Help: "Duration in seconds for syncing the status of the cluster once.",
	}, []string{"cluster_name"})

	// clusterHealthProbeSuccess reports the raw health probe result before threshold adjustment.
	clusterHealthProbeSuccess = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: clusterHealthProbeSuccessName,
		Help: "Result of the last health probe for the member cluster (1 for success, 0 for failure). " +
			"Reflects the raw probe outcome before threshold adjustment.",
	}, []string{memberClusterLabel})

	// clusterHealthProbeDuration reports the duration of the health probe HTTP call only.
	clusterHealthProbeDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    clusterHealthProbeDurationName,
		Help:    "Duration in seconds of the health probe to the member cluster.",
		Buckets: []float64{0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	}, []string{memberClusterLabel})

	// clusterHealthProbeTotal counts probes by result type.
	clusterHealthProbeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: clusterHealthProbeTotalName,
		Help: "Total number of health probes to the member cluster, categorized by result.",
	}, []string{memberClusterLabel, "result"})

	// clusterReadySince records the unix timestamp when the cluster last became Ready.
	// Set to 0 when the cluster is not Ready.
	clusterReadySince = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: clusterReadySinceName,
		Help: "Unix timestamp when the member cluster last became Ready. " +
			"Set to 0 when the cluster is not Ready.",
	}, []string{memberClusterLabel})

	// clusterConditionLastTransition records the unix timestamp of the last condition state transition.
	clusterConditionLastTransition = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: clusterConditionLastTransitionName,
		Help: "Unix timestamp of the last condition state transition for the member cluster.",
	}, []string{memberClusterLabel})

	// clusterHealthTransitionsTotal counts health state transitions.
	clusterHealthTransitionsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: clusterHealthTransitionsTotalName,
		Help: "Total number of health state transitions for the member cluster.",
	}, []string{memberClusterLabel, "from_state", "to_state"})
)

// RecordClusterStatus records the status of the given cluster.
func RecordClusterStatus(cluster *v1alpha1.Cluster) {
	clusterReadyGauge.WithLabelValues(cluster.Name).Set(func() float64 {
		if util.IsClusterReady(&cluster.Status) {
			return 1
		}
		return 0
	}())

	if cluster.Status.NodeSummary != nil {
		clusterTotalNodeNumberGauge.WithLabelValues(cluster.Name).Set(float64(cluster.Status.NodeSummary.TotalNum))
		clusterReadyNodeNumberGauge.WithLabelValues(cluster.Name).Set(float64(cluster.Status.NodeSummary.ReadyNum))
	}

	if cluster.Status.ResourceSummary != nil {
		if cluster.Status.ResourceSummary.Allocatable != nil {
			clusterMemoryAllocatableGauge.WithLabelValues(cluster.Name).Set(cluster.Status.ResourceSummary.Allocatable.Memory().AsApproximateFloat64())
			clusterCPUAllocatableGauge.WithLabelValues(cluster.Name).Set(cluster.Status.ResourceSummary.Allocatable.Cpu().AsApproximateFloat64())
			clusterPodAllocatableGauge.WithLabelValues(cluster.Name).Set(cluster.Status.ResourceSummary.Allocatable.Pods().AsApproximateFloat64())
		}

		if cluster.Status.ResourceSummary.Allocated != nil {
			clusterMemoryAllocatedGauge.WithLabelValues(cluster.Name).Set(cluster.Status.ResourceSummary.Allocated.Memory().AsApproximateFloat64())
			clusterCPUAllocatedGauge.WithLabelValues(cluster.Name).Set(cluster.Status.ResourceSummary.Allocated.Cpu().AsApproximateFloat64())
			clusterPodAllocatedGauge.WithLabelValues(cluster.Name).Set(cluster.Status.ResourceSummary.Allocated.Pods().AsApproximateFloat64())
		}
	}
}

// RecordClusterHealthProbeSuccess records the raw health probe result before threshold adjustment.
func RecordClusterHealthProbeSuccess(clusterName string, online, healthy bool) {
	val := float64(0)
	if online && healthy {
		val = 1
	}
	clusterHealthProbeSuccess.WithLabelValues(clusterName).Set(val)
}

// ProbeResultLabel returns the result label for a health probe.
func ProbeResultLabel(online, healthy bool) string {
	if online && healthy {
		return HealthStateSuccess
	}
	return HealthStateError
}

// RecordClusterHealthProbeTotal increments the probe counter with the appropriate result label.
func RecordClusterHealthProbeTotal(clusterName string, online, healthy bool) {
	clusterHealthProbeTotal.WithLabelValues(clusterName, ProbeResultLabel(online, healthy)).Inc()
}

// RecordClusterHealthProbeDuration records the duration of the health probe HTTP call.
func RecordClusterHealthProbeDuration(clusterName string, startTime time.Time) {
	clusterHealthProbeDuration.WithLabelValues(clusterName).Observe(utilmetrics.DurationInSeconds(startTime))
}

// RecordClusterReadySince updates the ready-since timestamp based on the threshold-adjusted condition.
// On transition to Ready, sets the timestamp to the provided time. On transition away from Ready, sets to 0.
// No-op when the status hasn't changed.
func RecordClusterReadySince(clusterName string, prevStatus, currentStatus metav1.ConditionStatus, timestamp time.Time) {
	if prevStatus == currentStatus {
		return
	}
	if currentStatus == metav1.ConditionTrue {
		clusterReadySince.WithLabelValues(clusterName).Set(float64(timestamp.Unix()))
	} else {
		clusterReadySince.WithLabelValues(clusterName).Set(0)
	}
}

// RecordClusterConditionLastTransition records the timestamp of a condition state transition.
// No-op when the status hasn't changed.
func RecordClusterConditionLastTransition(clusterName string, prevStatus, currentStatus metav1.ConditionStatus, timestamp time.Time) {
	if prevStatus == currentStatus {
		return
	}
	clusterConditionLastTransition.WithLabelValues(clusterName).Set(float64(timestamp.Unix()))
}

// RecordClusterHealthTransition increments the transition counter when the raw probe state changes.
// No-op when the state hasn't changed.
func RecordClusterHealthTransition(clusterName string, fromState, toState metav1.ConditionStatus) {
	if fromState == toState {
		return
	}
	clusterHealthTransitionsTotal.WithLabelValues(clusterName, string(fromState), string(toState)).Inc()
}

// RecordClusterSyncStatusDuration records the duration of the given cluster syncing status
func RecordClusterSyncStatusDuration(cluster *v1alpha1.Cluster, startTime time.Time) {
	clusterSyncStatusDuration.WithLabelValues(cluster.Name).Observe(utilmetrics.DurationInSeconds(startTime))
}

// CleanupMetricsForCluster removes the cluster status metrics after the cluster is deleted.
func CleanupMetricsForCluster(clusterName string) {
	clusterReadyGauge.DeleteLabelValues(clusterName)
	clusterTotalNodeNumberGauge.DeleteLabelValues(clusterName)
	clusterReadyNodeNumberGauge.DeleteLabelValues(clusterName)
	clusterMemoryAllocatableGauge.DeleteLabelValues(clusterName)
	clusterCPUAllocatableGauge.DeleteLabelValues(clusterName)
	clusterPodAllocatableGauge.DeleteLabelValues(clusterName)
	clusterMemoryAllocatedGauge.DeleteLabelValues(clusterName)
	clusterCPUAllocatedGauge.DeleteLabelValues(clusterName)
	clusterPodAllocatedGauge.DeleteLabelValues(clusterName)
	clusterSyncStatusDuration.DeleteLabelValues(clusterName)
	clusterHealthProbeSuccess.DeleteLabelValues(clusterName)
	clusterHealthProbeDuration.DeleteLabelValues(clusterName)
	clusterHealthProbeTotal.DeletePartialMatch(prometheus.Labels{memberClusterLabel: clusterName})
	clusterReadySince.DeleteLabelValues(clusterName)
	clusterConditionLastTransition.DeleteLabelValues(clusterName)
	clusterHealthTransitionsTotal.DeletePartialMatch(prometheus.Labels{memberClusterLabel: clusterName})
}

// ClusterCollectors returns the collectors about clusters.
func ClusterCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		clusterReadyGauge,
		clusterTotalNodeNumberGauge,
		clusterReadyNodeNumberGauge,
		clusterMemoryAllocatableGauge,
		clusterCPUAllocatableGauge,
		clusterPodAllocatableGauge,
		clusterMemoryAllocatedGauge,
		clusterCPUAllocatedGauge,
		clusterPodAllocatedGauge,
		clusterSyncStatusDuration,
		clusterHealthProbeSuccess,
		clusterHealthProbeDuration,
		clusterHealthProbeTotal,
		clusterReadySince,
		clusterConditionLastTransition,
		clusterHealthTransitionsTotal,
	}
}
