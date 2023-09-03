/*
Copyright AppsCode Inc. and Contributors.

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

package metricshandler

import (
	"context"
	"fmt"
	"strings"

	falcoapi "kubeops.dev/falco-ui-server/apis/falco/v1alpha1"

	"k8s.io/klog/v2"
	"k8s.io/kube-state-metrics/v2/pkg/metric"
	generator "k8s.io/kube-state-metrics/v2/pkg/metric_generator"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type criteria struct {
	node, ns, priority, rule, pod string
}

func (c criteria) toString() string {
	return fmt.Sprintf("N=%s,NS=%s,P=%s,R=%s,POD=%s", c.node, c.ns, c.priority, c.rule, c.pod)
}

func toCriteria(str string) *criteria {
	s := strings.Split(str, ",")
	c := criteria{}
	for _, part := range s {
		kv := strings.Split(part, "=")
		if len(kv) != 2 {
			klog.Errorf("Invalid string found while converting: %s \n", str)
			return nil
		}
		switch kv[0] {
		case "N":
			c.node = kv[1]
		case "NS":
			c.ns = kv[1]
		case "P":
			c.priority = kv[1]
		case "R":
			c.rule = kv[1]
		case "POD":
			c.pod = kv[1]
		default:
		}
	}
	return &c
}

func collectFalcoMetrics(kc client.Client, genTotal generator.FamilyGenerator) (*metric.Family, error) {
	fTotal := genTotal.Generate(nil)

	var feList falcoapi.FalcoEventList
	err := kc.List(context.TODO(), &feList)
	if err != nil {
		return nil, err
	}

	mp := make(map[string]int)

	for i := 0; i < len(feList.Items); i++ {
		fe := feList.Items[i]
		podName := fe.GetLabels()["k8s.pod.name"]
		nsName := fe.GetLabels()["k8s.ns.name"]
		nodeName := fe.GetLabels()["k8s.node.name"]
		if podName == "" || nsName == "" || nodeName == "" {
			continue
		}

		s := criteria{
			pod:      podName,
			ns:       nsName,
			node:     nodeName,
			priority: fe.Spec.Priority,
			rule:     fe.Spec.Rule,
		}.toString()
		mp[s]++
	}

	for l, v := range mp {
		c := toCriteria(l)
		if c == nil {
			continue
		}
		mTotal := metric.Metric{
			LabelKeys: []string{
				"node_name",
				"ns_name",
				"priority",
				"rule",
				"pod_name",
			},
			LabelValues: []string{
				c.node,
				c.ns,
				c.priority,
				c.rule,
				c.pod,
			},
			Value: float64(v),
		}
		fTotal.Metrics = append(fTotal.Metrics, &mTotal)
	}

	return fTotal, nil
}
