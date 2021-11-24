// Copyright 2021 OnMetal authors
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

package exec

import (
	"context"
	"fmt"

	"github.com/onmetal/kubectl-onmetal/exec"
	"github.com/onmetal/onmetal-console/tty/os"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/spf13/cobra"
)

func Command(restClientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "exec machine-name",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			name := args[0]
			return Run(ctx, restClientGetter, name)
		},
	}

	return cmd
}

func Run(ctx context.Context, restClientGetter genericclioptions.RESTClientGetter, name string) error {
	cfg, err := restClientGetter.ToRESTConfig()
	if err != nil {
		return fmt.Errorf("error getting rest config: %w", err)
	}

	namespace, _, err := restClientGetter.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return fmt.Errorf("error determining target namespace: %w", err)
	}

	c, err := client.New(cfg, client.Options{})
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	tty, err := os.FromStdStreams()
	if err != nil {
		return fmt.Errorf("error creating tty: %w", err)
	}

	return exec.Run(ctx, c, tty, namespace, name)
}
