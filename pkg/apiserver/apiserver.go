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

package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"kubeops.dev/falco-ui-server/apis/falco"
	"kubeops.dev/falco-ui-server/apis/falco/install"
	api "kubeops.dev/falco-ui-server/apis/falco/v1alpha1"
	"kubeops.dev/falco-ui-server/pkg/cleaner"
	"kubeops.dev/falco-ui-server/pkg/falcosidekick"
	"kubeops.dev/falco-ui-server/pkg/falcosidekick/metricshandler"
	festorage "kubeops.dev/falco-ui-server/pkg/registry/falco/falcoevent"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/klog/v2/klogr"
	cu "kmodules.xyz/client-go/client"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var (
	// Scheme defines methods for serializing and deserializing API objects.
	Scheme = runtime.NewScheme()
	// Codecs provides methods for retrieving codecs and serializers for specific
	// versions and content types.
	Codecs = serializer.NewCodecFactory(Scheme)
)

func init() {
	install.Install(Scheme)
	utilruntime.Must(clientgoscheme.AddToScheme(Scheme))
	utilruntime.Must(api.AddToScheme(Scheme))

	// we need to add the options to empty v1
	// TODO fix the server code to avoid this
	metav1.AddToGroupVersion(Scheme, schema.GroupVersion{Version: "v1"})

	// TODO: keep the generic API server from wanting this
	unversioned := schema.GroupVersion{Group: "", Version: "v1"}
	Scheme.AddUnversionedTypes(unversioned,
		&metav1.Status{},
		&metav1.APIVersions{},
		&metav1.APIGroupList{},
		&metav1.APIGroup{},
		&metav1.APIResourceList{},
	)
}

// ExtraConfig holds custom apiserver config
type ExtraConfig struct {
	ClientConfig        *restclient.Config
	KubeClient          kubernetes.Interface
	KubeInformerFactory informers.SharedInformerFactory
	ResyncPeriod        time.Duration
	EventTTLPeriod      time.Duration
}

// Config defines the config for the apiserver
type Config struct {
	GenericConfig *genericapiserver.RecommendedConfig
	ExtraConfig   ExtraConfig
}

// FalcoUIServer contains state for a Kubernetes cluster master/api server.
type FalcoUIServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
	Manager          manager.Manager
}

type completedConfig struct {
	GenericConfig genericapiserver.CompletedConfig
	ExtraConfig   *ExtraConfig
}

// CompletedConfig embeds a private pointer that cannot be instantiated outside of this package.
type CompletedConfig struct {
	*completedConfig
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() CompletedConfig {
	c := completedConfig{
		cfg.GenericConfig.Complete(),
		&cfg.ExtraConfig,
	}
	return CompletedConfig{&c}
}

// New returns a new instance of FalcoUIServer from the given config.
func (c completedConfig) New(ctx context.Context) (*FalcoUIServer, error) {
	genericServer, err := c.GenericConfig.New("falco-ui-server", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	// ctrl.SetLogger(...)
	log.SetLogger(klogr.New()) // nolint:staticcheck
	setupLog := log.Log.WithName("setup")

	cfg := c.ExtraConfig.ClientConfig
	metricsHandlers := map[string]http.Handler{}
	mgr, err := manager.New(cfg, manager.Options{
		Scheme: Scheme,
		Metrics: metricsserver.Options{
			BindAddress:   "",
			ExtraHandlers: metricsHandlers,
		},
		HealthProbeBindAddress: "",
		LeaderElection:         false,
		LeaderElectionID:       "5b87adeb.falco.appscode.com",
		Client: client.Options{
			Cache: &client.CacheOptions{
				DisableFor: []client.Object{
					&api.FalcoEvent{},
				},
			},
		},
		NewClient: cu.NewClient,
		Cache: cache.Options{
			SyncPeriod: &c.ExtraConfig.ResyncPeriod,
			ByObject: map[client.Object]cache.ByObject{
				new(core.Pod): {
					Label: labels.SelectorFromSet(map[string]string{
						"app.kubernetes.io/name": "falco",
					}),
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("unable to start manager, reason: %v", err)
	}
	metricsHandlers["/falcoevents"] = falcosidekick.Handler(mgr.GetClient())
	metricsHandlers["/falcometrics"] = metricshandler.Handler(mgr.GetClient())

	setupLog.Info("setup done!")

	s := &FalcoUIServer{
		GenericAPIServer: genericServer,
		Manager:          mgr,
	}
	{
		apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(falco.GroupName, Scheme, metav1.ParameterCodec, Codecs)

		v1alpha1storage := map[string]rest.Storage{}
		{
			storage, err := festorage.NewStorage(Scheme, c.GenericConfig.RESTOptionsGetter)
			if err != nil {
				return nil, err
			}
			v1alpha1storage[api.ResourceFalcoEvents] = storage
		}
		apiGroupInfo.VersionedResourcesStorageMap["v1alpha1"] = v1alpha1storage

		if err := s.GenericAPIServer.InstallAPIGroup(&apiGroupInfo); err != nil {
			return nil, err
		}
	}
	go cleaner.StartCleaner(mgr.GetClient(), c.ExtraConfig.EventTTLPeriod)
	return s, nil
}
