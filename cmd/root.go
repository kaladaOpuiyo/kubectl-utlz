/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"fmt"
	"os"

	"github.com/kaladaOpuiyo/kubectl-utlz/pkg/client"
	"github.com/kaladaOpuiyo/kubectl-utlz/pkg/metrics"
	"github.com/kaladaOpuiyo/kubectl-utlz/pkg/util"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
)

var (
	sortBy string
	view   string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubectl-utlz [node name]",
	Short: "kubectl plugin to generate cpu/mem for pods scheduled on a node",
	Long:  `kubectl plugin to generate cpu/mem for pods scheduled on a node`,

	Run: func(cmd *cobra.Command, args []string) {

		streams := genericclioptions.IOStreams{
			Out: os.Stdout,
		}

		w := printers.GetNewTabWriter(streams.Out)
		defer w.Flush()

		// Cheap way to check this value is set
		var nodeName string

		if len(os.Args) > 1 {
			nodeName = os.Args[1]
		} else {
			fmt.Fprint(w, "Please provide node name\n")
			os.Exit(1)

		}

		configFlags := genericclioptions.NewConfigFlags(true)
		config, err := configFlags.ToRESTConfig()
		if err != nil {
			fmt.Fprint(w, err)
			os.Exit(1)
		}

		clientset, metricsClientset, err := client.NewMeticsClientSets(config)
		if err != nil {
			fmt.Fprint(w, err)
			os.Exit(1)
		}

		nodeMetricsByPod := metrics.NewNodeMetricsByPod(clientset, metricsClientset, nodeName, sortBy, view)

		podMetrics, err := nodeMetricsByPod.GetPodMetrics()
		if err != nil {
			fmt.Fprint(w, err)
			os.Exit(1)
		}

		util.Print(w, podMetrics, view)

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&sortBy, "sort-by", "", "sort by (cpu) or (memory)")
	rootCmd.PersistentFlags().StringVar(&view, "view", "", "view (resources) or resources + usage (wide) ")
}
