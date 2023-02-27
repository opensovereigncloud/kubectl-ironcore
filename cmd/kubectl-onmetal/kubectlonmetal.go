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

package kubectlonmetal

import (
	"os"

	"github.com/onmetal/kubectl-onmetal/cmd/create"
	"github.com/onmetal/kubectl-onmetal/cmd/exec"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

type Options struct {
	genericclioptions.IOStreams
}

func DefaultCommand() *cobra.Command {
	return Command(Options{
		IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
	})
}

func Command(opts Options) *cobra.Command {
	var (
		configFlags = genericclioptions.NewConfigFlags(true)
	)

	cmd := &cobra.Command{
		Use: "kubectl-onmetal",
	}

	configFlags.AddFlags(cmd.PersistentFlags())

	f := cmdutil.NewFactory(configFlags)

	cmd.AddCommand(
		exec.Command(configFlags),
		create.Command(f, opts.IOStreams),
	)

	return cmd
}
