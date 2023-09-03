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

package falcosidekick

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"

	"kubeops.dev/falco-ui-server/apis/falco/v1alpha1"
	"kubeops.dev/falco-ui-server/pkg/falcosidekick/metricshandler"
	"kubeops.dev/falco-ui-server/pkg/falcosidekick/types"

	"github.com/google/uuid"
	jsonx "gomodules.xyz/encoding/json"
	core "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
	cu "kmodules.xyz/client-go/client"
	"kmodules.xyz/client-go/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Handler is Falco Sidekick main handler (default).
func Handler(kc client.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			http.Error(w, "Please send a valid request body", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			err := metricshandler.CollectMetrics(kc, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}

		falcopayload, err := newFalcoPayload(r.Body)
		if err != nil || !falcopayload.Check() {
			http.Error(w, "Please send a valid request body", http.StatusBadRequest)
			return
		}
		mustForwardEvent(kc, falcopayload)
	})
}

func newFalcoPayload(payload io.Reader) (types.FalcoPayload, error) {
	var falcopayload types.FalcoPayload

	d := json.NewDecoder(payload)
	d.UseNumber()

	err := d.Decode(&falcopayload)
	if err != nil {
		return types.FalcoPayload{}, err
	}

	if len(config.Customfields) > 0 {
		if falcopayload.OutputFields == nil {
			falcopayload.OutputFields = make(map[string]interface{})
		}
		for key, value := range config.Customfields {
			falcopayload.OutputFields[key] = value
		}
	}

	if falcopayload.Source == "" {
		falcopayload.Source = "syscalls"
	}

	falcopayload.UUID = uuid.New().String()

	var kn, kp string
	for i, j := range falcopayload.OutputFields {
		if j != nil {
			if i == "k8s.ns.name" {
				kn = j.(string)
			}
			if i == "k8s.pod.name" {
				kp = j.(string)
			}
		}
	}

	if len(config.Templatedfields) > 0 {
		if falcopayload.OutputFields == nil {
			falcopayload.OutputFields = make(map[string]interface{})
		}
		for key, value := range config.Templatedfields {
			tmpl, err := template.New("").Parse(value)
			if err != nil {
				log.Printf("[ERROR] : Parsing error for templated field '%v': %v\n", key, err)
				continue
			}
			v := new(bytes.Buffer)
			if err := tmpl.Execute(v, falcopayload.OutputFields); err != nil {
				log.Printf("[ERROR] : Parsing error for templated field '%v': %v\n", key, err)
			}
			falcopayload.OutputFields[key] = v.String()
		}
	}

	promLabels := map[string]string{"rule": falcopayload.Rule, "priority": falcopayload.Priority.String(), "k8s_ns_name": kn, "k8s_pod_name": kp}
	if falcopayload.Hostname != "" {
		promLabels["hostname"] = falcopayload.Hostname
	}

	for key, value := range config.Customfields {
		if regPromLabels.MatchString(key) {
			promLabels[key] = value
		}
	}
	for _, i := range config.Prometheus.ExtraLabelsList {
		promLabels[strings.ReplaceAll(i, ".", "_")] = ""
		for key, value := range falcopayload.OutputFields {
			if key == i && regPromLabels.MatchString(strings.ReplaceAll(key, ".", "_")) {
				switch value.(type) {
				case string:
					promLabels[strings.ReplaceAll(key, ".", "_")] = fmt.Sprintf("%v", value)
				default:
					continue
				}
			}
		}
	}

	if config.BracketReplacer != "" {
		for i, j := range falcopayload.OutputFields {
			if strings.Contains(i, "[") {
				falcopayload.OutputFields[strings.ReplaceAll(strings.ReplaceAll(i, "]", ""), "[", config.BracketReplacer)] = j
				delete(falcopayload.OutputFields, i)
			}
		}
	}

	if config.Debug {
		body, _ := json.Marshal(falcopayload)
		log.Printf("[DEBUG] : Falco's payload : %v\n", string(body))
	}

	return falcopayload, nil
}

func forwardEvent(kc client.Client, payload types.FalcoPayload, evHash uint64) error {
	var nodeName string
	if payload.Hostname != "" {
		var pod core.Pod
		key := client.ObjectKey{
			Namespace: meta.PodNamespace(),
			Name:      payload.Hostname,
		}
		if err := kc.Get(context.TODO(), key, &pod); err == nil {
			nodeName = pod.Spec.NodeName
		}
	}

	obj := &v1alpha1.FalcoEvent{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
			Kind:       v1alpha1.ResourceKindFalcoEvent,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("fe-%d", evHash),
			Labels:      map[string]string{},
			Annotations: nil,
		},
		Spec: v1alpha1.FalcoEventSpec{
			UUID:     payload.UUID,
			Output:   payload.Output,
			Priority: payload.Priority.String(),
			Rule:     payload.Rule,
			Time:     metav1.NewTime(payload.Time),
			// OutputFields: apiextensionsv1.JSON{},
			Source:   payload.Source,
			Tags:     payload.Tags,
			Hostname: payload.Hostname,
		},
	}

	fields, err := jsonx.Marshal(payload.OutputFields)
	if err != nil {
		return err
	}
	obj.Spec.OutputFields = apiextensionsv1.JSON{Raw: fields}

	for k, v := range payload.OutputFields {
		switch k {
		case "k8s.ns.name", "k8s.pod.name":
			val, ok := v.(string)
			if ok {
				obj.Labels[k] = val
			}
		}
	}
	if nodeName != "" {
		obj.Labels["k8s.node.name"] = nodeName
		obj.Spec.Nodename = nodeName
	}

	_, err = cu.CreateOrPatch(context.TODO(), kc, obj, func(in client.Object, createOp bool) client.Object {
		o := in.(*v1alpha1.FalcoEvent)
		o.Labels = obj.Labels
		o.Spec = obj.Spec

		return o
	})
	return err
}

var eventHashes = make(map[uint64]time.Time)

const eventRefreshTTL = 10 * time.Minute

func mustForwardEvent(kc client.Client, payload types.FalcoPayload) {
	hashKey := payload.HashKey()
	lastTime, found := eventHashes[hashKey]

	if !found || time.Since(lastTime) > eventRefreshTTL {
		err := forwardEvent(kc, payload, hashKey)
		if err = client.IgnoreAlreadyExists(err); err != nil {
			klog.ErrorS(err, "failed to write falco event")
		} else {
			eventHashes[hashKey] = payload.Time
		}
	}
}
