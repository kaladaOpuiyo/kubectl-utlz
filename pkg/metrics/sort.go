package metrics

// SortMetrics ...
type SortMetrics []PodMetric

func (ms SortMetrics) Len() int {
	return len(ms)
}

func (ms SortMetrics) Less(i, j int) bool {

	switch ms[0].SortBy {
	case "cpu":
		return ms[i].CPU > ms[j].CPU
	case "memory":
		return ms[i].Memory > ms[j].Memory
	default:
		return false
	}
}

func (ms SortMetrics) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}
