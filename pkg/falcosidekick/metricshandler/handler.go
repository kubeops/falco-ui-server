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
	"io"
	"net/http"

	"kubeops.dev/ui-server/pkg/metricsstore"

	"k8s.io/kube-state-metrics/v2/pkg/metric"
	generator "k8s.io/kube-state-metrics/v2/pkg/metric_generator"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	falcoMetricPrefix = "falco_appscode_com_"
)

func Handler(kc client.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "Please send a valid request body", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			err := CollectMetrics(kc, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}
	})
}

func CollectMetrics(kc client.Client, w io.Writer) error {
	generators := getFamilyGenerators()
	if len(generators) == 0 {
		_, err := w.Write([]byte("OK"))
		return err
	}

	// Generate the headers for the resources metrics
	headers := generator.ExtractMetricFamilyHeaders(generators)
	store := metricsstore.NewMetricsStore(headers)

	if all, err := collectFalcoMetrics(kc, generators[0]); err != nil {
		return err
	} else {
		store.Add(all)
	}

	return store.WriteAll(w)
}

func getFamilyGenerators() []generator.FamilyGenerator {
	fn := func(obj any) *metric.Family { return new(metric.Family) }
	generators := make([]generator.FamilyGenerator, 0, 1)
	generators = append(generators, generator.FamilyGenerator{
		Name:              falcoMetricPrefix + "events",
		Help:              "All statistics",
		Type:              metric.Gauge,
		DeprecatedVersion: "",
		GenerateFunc:      fn,
	})
	return generators
}
