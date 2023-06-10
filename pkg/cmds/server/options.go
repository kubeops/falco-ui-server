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

package server

import (
	"time"

	"kubeops.dev/falco-ui-server/pkg/apiserver"

	"github.com/spf13/pflag"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

type ExtraOptions struct {
	QPS          float64
	Burst        int
	ResyncPeriod time.Duration

	EventTTLPeriod time.Duration
}

func NewExtraOptions() *ExtraOptions {
	return &ExtraOptions{
		ResyncPeriod:   10 * time.Minute,
		QPS:            1e6,
		Burst:          1e6,
		EventTTLPeriod: time.Hour * 12,
	}
}

func (s *ExtraOptions) AddFlags(fs *pflag.FlagSet) {
	fs.DurationVar(&s.ResyncPeriod, "resync-period", s.ResyncPeriod, "If non-zero, will re-list this often. Otherwise, re-list will be delayed aslong as possible (until the upstream source closes the watch or times out.")
	fs.Float64Var(&s.QPS, "qps", s.QPS, "The maximum QPS to the master from this client")
	fs.IntVar(&s.Burst, "burst", s.Burst, "The maximum burst for throttle")

	fs.DurationVar(&s.EventTTLPeriod, "event-ttl", s.EventTTLPeriod, "Events older than this period will be garbage collected")
}

func (s *ExtraOptions) ApplyTo(cfg *apiserver.ExtraConfig) error {
	cfg.ClientConfig.QPS = float32(s.QPS)
	cfg.ClientConfig.Burst = s.Burst
	cfg.ResyncPeriod = s.ResyncPeriod
	cfg.EventTTLPeriod = s.EventTTLPeriod

	var err error
	if cfg.KubeClient, err = kubernetes.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	cfg.KubeInformerFactory = informers.NewSharedInformerFactory(cfg.KubeClient, cfg.ResyncPeriod)

	return nil
}

func (s *ExtraOptions) Validate() []error {
	return nil
}
