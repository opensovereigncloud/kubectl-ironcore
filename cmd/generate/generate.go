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

package generate

import (
	bootstrapkubeconfig "github.com/ironcore-dev/kubectl-ironcore/cmd/generate/bootstrap-kubeconfig"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
)

func Command(pathOptions *clientcmd.PathOptions, streams genericclioptions.IOStreams) *cobra.Command {
	if len(pathOptions.ExplicitFileFlag) == 0 {
		pathOptions.ExplicitFileFlag = clientcmd.RecommendedConfigPathFlag
	}

	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate artifacts / configuration commands.",
	}

	// file paths are common to all sub commands
	cmd.PersistentFlags().StringVar(&pathOptions.LoadingRules.ExplicitPath, pathOptions.ExplicitFileFlag, pathOptions.LoadingRules.ExplicitPath, "use a particular kubeconfig file")

	cmd.AddCommand(
		bootstrapkubeconfig.Command(pathOptions, streams),
	)

	return cmd
}
