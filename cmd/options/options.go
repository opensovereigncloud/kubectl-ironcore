// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"io"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"
)

func Command(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "options",
		Short: "Print the list of flags inherited by all commands",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Usage()
		},
	}

	// The `options` command needs write its output to the `out` stream
	// (typically stdout). Without calling SetOutput here, the Usage()
	// function call will fall back to stderr.
	//
	// See https://github.com/kubernetes/kubernetes/pull/46394 for details.
	cmd.SetOut(out)
	cmd.SetErr(out)

	templates.UseOptionsTemplates(cmd)
	return cmd
}
