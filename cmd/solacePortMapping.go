/*
Copyright © 2022 BEN MANSOUR Mohamed Rafik

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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/rung/go-safecast"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

const group string = "scalable.solace.io/v1alpha1"
const resource string = "solacescalables"

var svcPortsMappingCmd = &cobra.Command{
	Use:   "solmap",
	Short: "Solace service ports mapping with HAProxy service",
	Long: `Show solace scalable service ports mapping with HAProxy service in tabular mode.

	Command example:

	kubectl solmap -n solacescalable -c solaceCrdName`,
	Run: func(cmd *cobra.Command, args []string) {

		on, err := cmd.Flags().GetString("operatorNamespace")
		if err != nil {
			panic(err.Error())
		}
		cn, err := cmd.Flags().GetString("crdName")
		if err != nil {
			panic(err.Error())
		}
		clientset, dynamicset := ClientSet(genericclioptions.NewConfigFlags(true))

		var pubSub = []string{"pub", "sub"}
		solaceScalable, err := dynamicset.Resource(
			schema.GroupVersionResource{
				Group:    group,
				Resource: resource,
			},
		).Namespace(cn).Get(context.TODO(), "solacescalable", metav1.GetOptions{})
		if err != nil {
			panic(err)
		}
		tata, err := json.Marshal((*solaceScalable).Object["spec"])
		if err != nil {
			panic(err)
		}

		haProxy := HaproxyCRD{}
		if err := json.Unmarshal(tata, &haProxy); err != nil {
			panic(err)
		}

		for _, arg := range pubSub {
			var data = []svcMap{}

			tcpCmName := on + "-" + arg + "-tcp-ingress"
			tcpCm, err := clientset.CoreV1().ConfigMaps(on).Get(context.TODO(), tcpCmName, metav1.GetOptions{})
			if err != nil {
				fmt.Printf("\nConfigmap %v does not exist in the %v namespace\n\n", tcpCmName, on)
				panic(err.Error())
			}

			var haProxySvcName = haProxy.Haproxy.Publish.ServiceName
			if arg == "sub" {
				haProxySvcName = haProxy.Haproxy.Subscribe.ServiceName
			}

			haProxySvc, err := clientset.CoreV1().Services(
				haProxy.Haproxy.Namespace,
			).Get(
				context.TODO(),
				haProxySvcName,
				metav1.GetOptions{},
			)
			if err != nil {
				panic(err.Error())
			}

			for _, cmV := range tcpCm.Data {
				//example format:
				//solacescalable/test-botti-1029-amqp-pub:1100
				svc := strings.Split(cmV, "/")

				solacePortSlice := strings.Split(svc[1], ":")
				svcPortSlice := strings.Split(solacePortSlice[0], "-")

				svcName := solacePortSlice[0]
				svcPort := svcPortSlice[2]
				solacePort := solacePortSlice[1]

				for _, svcV := range haProxySvc.Spec.Ports {
					if "tcp-"+svcPort == svcV.Name {
						p, _ := safecast.Atoi32(svcPort)
						sp, _ := safecast.Atoi32(solacePort)

						data = append(
							data,
							svcMap{
								SvcName:       svcName,
								SvcPort:       p,
								SvcSolacePort: sp,
								SvcNodePort:   svcV.NodePort,
							},
						)
					}
				}
			}
			DrawSvcTable(
				[]string{
					"Service Name " + arg,
					"Service Port",
					"Solace Port",
					"HAProxy Port",
				},
				data,
			)
		}
	},
}

func Execute() {
	err := svcPortsMappingCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	svcPortsMappingCmd.Aliases = append(svcPortsMappingCmd.Aliases, "svcmap", "map")
	svcPortsMappingCmd.PersistentFlags().StringP("crdName", "c", "solacescalable", "Name of the CRD")
	svcPortsMappingCmd.PersistentFlags().StringP("operatorNamespace", "n", "default", "operator's namespace")
}

// ClientSet k8s
func ClientSet(configFlags *genericclioptions.ConfigFlags) (*kubernetes.Clientset, dynamic.Interface) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		panic("kube config load error")
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic("gen kube config error")
	}
	dynamicSet, err := dynamic.NewForConfig(config)
	if err != nil {
		panic("gen kube config error")
	}

	return clientSet, dynamicSet
}

func DrawSvcTable(header []string, ports []svcMap) {
	var data = [][]string{}
	for _, v := range ports {
		data = append(data,
			[]string{
				v.SvcName,
				strconv.Itoa(int(v.SvcPort)),
				strconv.Itoa(int(v.SvcSolacePort)),
				strconv.Itoa(int(v.SvcNodePort)),
			},
		)

	}
	renderTable(header, data)
}

func renderTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)

	for _, v := range data {
		table.Append(v)
	}
	table.SetHeader(header)
	table.Render()
}
