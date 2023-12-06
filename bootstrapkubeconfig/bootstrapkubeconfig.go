// Copyright 2023 IronCore authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bootstrapkubeconfig

import (
	"fmt"

	"github.com/ironcore-dev/kubectl-ironcore/utils/bootstraptoken"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	DefaultBootstrapContext = "bootstrap"
)

// GenerateOptions are options for generating a bootstrap kubeconfig.
type GenerateOptions struct {
	// Context is the context to use for generating the bootstrap kubeconfig.
	// If empty, the current context will be used.
	Context string

	// BootstrapContext is context name of the resulting bootstrap kubeconfig.
	// If empty, DefaultBootstrapContext will be used.
	BootstrapContext string
}

func (o *GenerateOptions) ApplyOptions(opts []func(*GenerateOptions)) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithContext sets the context to use to the given one.
func WithContext(context string) func(*GenerateOptions) {
	return func(options *GenerateOptions) {
		options.Context = context
	}
}

func setGenerateOptionsDefaults(o *GenerateOptions) {
	if o.BootstrapContext == "" {
		o.BootstrapContext = DefaultBootstrapContext
	}
}

// Generate generates a bootstrap clientcmdapi.Config from the given starting config and token.
func Generate(
	startingCfg *clientcmdapi.Config,
	token *bootstraptoken.BootstrapToken,
	opts ...func(*GenerateOptions),
) (*clientcmdapi.Config, error) {
	o := &GenerateOptions{}
	o.ApplyOptions(opts)
	setGenerateOptionsDefaults(o)

	contextName := o.Context
	if contextName == "" {
		contextName = startingCfg.CurrentContext
	}
	if contextName == "" {
		return nil, fmt.Errorf("could not determine context name to use")
	}

	context, ok := startingCfg.Contexts[startingCfg.CurrentContext]
	if !ok {
		return nil, fmt.Errorf("context %q not found", startingCfg.CurrentContext)
	}

	cluster, ok := startingCfg.Clusters[context.Cluster]
	if !ok {
		return nil, fmt.Errorf("cluster %q not found", context.Cluster)
	}

	return &clientcmdapi.Config{
		Preferences: startingCfg.Preferences,
		Clusters: map[string]*clientcmdapi.Cluster{
			o.BootstrapContext: cluster,
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			o.BootstrapContext: {
				Token: fmt.Sprintf("%s.%s", token.ID, token.Secret),
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			o.BootstrapContext: {
				AuthInfo:   o.BootstrapContext,
				Cluster:    o.BootstrapContext,
				Namespace:  context.Namespace,
				Extensions: context.Extensions,
			},
		},
		CurrentContext: o.BootstrapContext,
		Extensions:     startingCfg.Extensions,
	}, nil
}
