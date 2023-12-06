// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package create

import (
	"github.com/ironcore-dev/kubectl-ironcore/cmd/create/token"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

func Command(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create configuration / items in a cluster.",
	}

	cmd.AddCommand(
		token.Command(f, streams),
	)

	return cmd
}
