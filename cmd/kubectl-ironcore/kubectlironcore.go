// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package kubectlironcore

import (
	"os"

	"github.com/ironcore-dev/kubectl-ironcore/cmd/create"
	"github.com/ironcore-dev/kubectl-ironcore/cmd/exec"
	"github.com/ironcore-dev/kubectl-ironcore/cmd/generate"
	"github.com/ironcore-dev/kubectl-ironcore/cmd/options"
	"github.com/ironcore-dev/kubectl-ironcore/cmd/version"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/util/templates"
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
		Use:   "kubectl-ironcore",
		Short: "Command line utility for operating and interacting with ironcore.",
		Run:   runHelp,
	}

	configFlags.AddFlags(cmd.PersistentFlags())

	f := cmdutil.NewFactory(configFlags)

	templates.ActsAsRootCommand(cmd, []string{"options"})

	cmd.AddCommand(
		exec.Command(configFlags),
		create.Command(f, opts.IOStreams),
		generate.Command(clientcmd.NewDefaultPathOptions(), opts.IOStreams),
		options.Command(opts.IOStreams.Out),
		version.Command(opts.IOStreams.Out),
	)

	return cmd
}

func runHelp(cmd *cobra.Command, _ []string) {
	_ = cmd.Help()
}
