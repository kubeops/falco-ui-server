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
	"log"
	"os"
	"regexp"
	"strings"

	"kubeops.dev/falco-ui-server/pkg/falcosidekick/types"
)

var (
	config *types.Configuration

	regPromLabels *regexp.Regexp
)

func init() {
	// detect unit testing and skip init.
	// see: https://github.com/alecthomas/kingpin/issues/187
	testing := (strings.HasSuffix(os.Args[0], ".test") ||
		strings.HasSuffix(os.Args[0], "__debug_bin"))
	if testing {
		return
	}

	regPromLabels, _ = regexp.Compile("^[a-zA-Z_:][a-zA-Z0-9_:]*$")

	config = getConfig()
}

func getConfig() *types.Configuration {
	c := &types.Configuration{
		Customfields:    make(map[string]string),
		BracketReplacer: "",
		Debug:           false,
	}

	// v.GetStringMapString("Customfields")

	if value, present := os.LookupEnv("CUSTOMFIELDS"); present {
		customfields := strings.Split(value, ",")
		for _, label := range customfields {
			tagkeys := strings.Split(label, ":")
			if len(tagkeys) == 2 {
				if strings.HasPrefix(tagkeys[1], "%") {
					if s := os.Getenv(tagkeys[1][1:]); s != "" {
						c.Customfields[tagkeys[0]] = s
					} else {
						log.Printf("[ERROR] : Can't find env var %v for custom fields", tagkeys[1][1:])
					}
				} else {
					c.Customfields[tagkeys[0]] = tagkeys[1]
				}
			}
		}
	}

	if c.Prometheus.ExtraLabels != "" {
		c.Prometheus.ExtraLabelsList = strings.Split(strings.ReplaceAll(c.Prometheus.ExtraLabels, " ", ""), ",")
	}
	return c
}
