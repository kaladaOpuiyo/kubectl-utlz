package util

import (
	"fmt"
	"io"

	"github.com/kaladaOpuiyo/kubectl-utlz/pkg/metrics"
)

var (

	// BaseColumns ...
	BaseColumns = []string{"NAMESPACE", "NAME"}

	// UsageColumns ...
	UsageColumns = []string{"CPU(cores)", "MEMORY(bytes)"}

	//PodLimitsRequestColumns ...
	PodLimitsRequestColumns = []string{"CPU REQUESTED(cores)", "MEMORY REQUESTED(bytes)", "CPU LIMITS(cores)", "MEMORY LIMITS(bytes)"}

	//View Options
	resourceView = "resources"
	wideView     = "wide"
)

// PrintColumnNames ...
func printColumnNames(out io.Writer, columns []string) {
	for _, name := range columns {
		printValue(out, name)
	}
	fmt.Fprint(out, "\n")
}

// PrintValue ...
func printValue(out io.Writer, value interface{}) {
	fmt.Fprintf(out, "%v\t", value)
}

// Print ...
func Print(out io.Writer, podMetrics []metrics.PodMetric, view string) {

	if podMetrics != nil {

		var columns []string

		switch view {
		case wideView:
			columns = append(BaseColumns, UsageColumns...)
			columns = append(columns, PodLimitsRequestColumns...)
		case resourceView:
			columns = append(BaseColumns, PodLimitsRequestColumns...)
		default:
			columns = append(BaseColumns, UsageColumns...)
		}

		printColumnNames(out, columns)

		for _, p := range podMetrics {

			switch view {
			case wideView:
				fmt.Fprintf(out, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t",
					p.Namespace, p.Name,
					fmt.Sprintf("%vm", p.CPU), fmt.Sprintf("%vMi", p.Memory),
					fmt.Sprintf("%vm", p.RequestedCPU), fmt.Sprintf("%vMi", p.RequestedMemory),
					fmt.Sprintf("%vm", p.LimitsCPU), fmt.Sprintf("%vMi", p.LimitsMemory))
			case resourceView:
				fmt.Fprintf(out, "%v\t%v\t%v\t%v\t%v\t%v\t",
					p.Namespace, p.Name,
					fmt.Sprintf("%vm", p.RequestedCPU), fmt.Sprintf("%vMi", p.RequestedMemory),
					fmt.Sprintf("%vm", p.LimitsCPU), fmt.Sprintf("%vMi", p.LimitsMemory))
			default:
				fmt.Fprintf(out, "%v\t%v\t%v\t%v\t",
					p.Namespace, p.Name,
					fmt.Sprintf("%vm", p.CPU), fmt.Sprintf("%vMi", p.Memory))

			}

			fmt.Fprint(out, "\n")

		}
	} else {
		fmt.Fprint(out, "No metrics found for Node provided ¯\\_(ツ)_/¯ \n")

	}

}
