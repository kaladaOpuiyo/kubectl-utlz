package metrics

import (
	"sort"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

var (
	cpuUsage        int64
	limits          v1.ResourceList
	limitscpu       int64
	limitsmemory    int64
	memoryUsage     int64
	metrics         []PodMetric
	metricsUsage    v1.ResourceList
	requestedcpu    int64
	requestedmemory int64
	requests        v1.ResourceList

	measuredResources = []v1.ResourceName{
		v1.ResourceCPU,
		v1.ResourceMemory,
	}
)

// PodMetric ...
type PodMetric struct {
	CPU             int64
	LimitsCPU       int64
	LimitsMemory    int64
	Memory          int64
	Name            string
	Namespace       string
	RequestedCPU    int64
	RequestedMemory int64
	SortBy          string
}

// NodeMetricsByPod ...
type NodeMetricsByPod struct {
	Clientset        *kubernetes.Clientset
	MetricsClientset *metricsv.Clientset
	Name             string
	Selector         string
	SortBy           string
	View             string
}

// NewNodeMetricsByPod ...
func NewNodeMetricsByPod(clientset *kubernetes.Clientset, metricsClientset *metricsv.Clientset, name, selector, sortBy string, view string) *NodeMetricsByPod {
	return &NodeMetricsByPod{
		Clientset:        clientset,
		MetricsClientset: metricsClientset,
		Name:             name,
		Selector:         selector,
		SortBy:           sortBy,
		View:             view,
	}
}

// GetPodMetrics ...
func (n *NodeMetricsByPod) GetPodMetrics() ([]PodMetric, error) {

	pm, err := n.getPodMetrics()
	if err != nil {
		return nil, err
	}

	return pm, nil
}

func (n *NodeMetricsByPod) getPodMetrics() ([]PodMetric, error) {

	pods, err := n.getPodsByNodeName()
	if err != nil {
		return nil, err
	}

	for _, pod := range pods.Items {

		podMetrics, err := n.podMetrics(pod.GetNamespace(), pod.GetName())

		// Need to handle this better,maybe
		// occurs when a pod is recently removed but still associated with node during the previous request
		// Should check error for specific type and inform user before continuing
		if err != nil {
			continue
		}

		if n.View == "wide" || n.View == "resources" {
			limits, requests = getContainerResources(pod.Spec.Containers)
			limitscpu = limits.Cpu().MilliValue()
			limitsmemory = limits.Memory().Value() / (1024 * 1024)
			requestedcpu = requests.Cpu().MilliValue()
			requestedmemory = requests.Memory().Value() / (1024 * 1024)
		}

		if n.View != "resources" {
			metricsUsage = getContainerMetricsUsage(podMetrics.Containers)
			cpuUsage = metricsUsage.Cpu().MilliValue()
			memoryUsage = metricsUsage.Memory().Value() / (1024 * 1024)
		}

		m := &PodMetric{
			CPU:             cpuUsage,
			LimitsCPU:       limitscpu,
			LimitsMemory:    limitsmemory,
			Memory:          memoryUsage,
			Name:            pod.Name,
			Namespace:       pod.Namespace,
			RequestedCPU:    requestedcpu,
			RequestedMemory: requestedmemory,
			SortBy:          n.SortBy,
		}

		metrics = append(metrics, *m)
	}

	sort.Sort(SortMetrics(metrics))

	return metrics, nil

}

func (n *NodeMetricsByPod) getPodsByNodeName() (*v1.PodList, error) {

	pods, err := n.Clientset.CoreV1().Pods("").List(metav1.ListOptions{
		FieldSelector: n.Selector + n.Name,
	})
	if err != nil {
		return nil, err
	}

	return pods, nil

}

func (n *NodeMetricsByPod) podMetrics(namespace, name string) (*metricsv1beta1.PodMetrics, error) {

	podMetrics, err := n.MetricsClientset.MetricsV1beta1().PodMetricses(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return podMetrics, nil
}

func getContainerMetricsUsage(containers []metricsv1beta1.ContainerMetrics) v1.ResourceList {

	var usage v1.ResourceList
	usu := make(v1.ResourceList)

	for _, c := range containers {
		c.Usage.DeepCopyInto(&usage)
		for _, res := range measuredResources {
			quantity := usu[res]
			quantity.Add(usage[res])
			usu[res] = quantity
		}
	}

	return usu

}

func getContainerResources(containers []v1.Container) (v1.ResourceList, v1.ResourceList) {

	var (
		limits   v1.ResourceList
		requests v1.ResourceList
	)

	lim := make(v1.ResourceList)
	req := make(v1.ResourceList)

	for _, c := range containers {

		c.Resources.Limits.DeepCopyInto(&limits)
		c.Resources.Requests.DeepCopyInto(&requests)

		for _, res := range measuredResources {

			quantityLim := lim[res]
			quantityLim.Add(limits[res])
			lim[res] = quantityLim

			quantityReq := req[res]
			quantityReq.Add(requests[res])
			req[res] = quantityReq

		}
	}

	return lim, req

}
