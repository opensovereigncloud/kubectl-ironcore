// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"context"
	"fmt"
	"net/http"
	"os"

	computev1alpha1 "github.com/ironcore-dev/ironcore/api/compute/v1alpha1"
	ironcoreclientgo "github.com/ironcore-dev/ironcore/client-go/ironcore"
	ironcoreclientgoscheme "github.com/ironcore-dev/ironcore/client-go/ironcore/scheme"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/util/term"
)

func Command(restClientGetter genericclioptions.RESTClientGetter) *cobra.Command {
	var insecureSkipTLSVerifyBackend bool

	cmd := &cobra.Command{
		Use:   "exec <machine-name>",
		Short: "Exec onto running entities in the cluster.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			name := args[0]
			return Run(ctx, restClientGetter, name, insecureSkipTLSVerifyBackend)
		},
	}

	cmd.Flags().BoolVar(&insecureSkipTLSVerifyBackend, "insecure-skip-tls-verify-backend", insecureSkipTLSVerifyBackend, "Whether to skip tls verification on the machinepoollet exec backend.")

	return cmd
}

func Run(ctx context.Context, restClientGetter genericclioptions.RESTClientGetter, name string, insecureSkipVerifyTLSBackend bool) error {
	cfg, err := restClientGetter.ToRESTConfig()
	if err != nil {
		return fmt.Errorf("error getting rest config: %w", err)
	}

	namespace, _, err := restClientGetter.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return fmt.Errorf("error determining target namespace: %w", err)
	}

	ironcoreClientset, err := ironcoreclientgo.NewForConfig(cfg)
	if err != nil {
		return err
	}

	req := ironcoreClientset.ComputeV1alpha1().RESTClient().
		Post().
		Namespace(namespace).
		Resource("machines").
		Name(name).
		SubResource("exec").
		VersionedParams(&computev1alpha1.MachineExecOptions{InsecureSkipTLSVerifyBackend: insecureSkipVerifyTLSBackend}, ironcoreclientgoscheme.ParameterCodec)

	var sizeQueue remotecommand.TerminalSizeQueue
	tty := term.TTY{
		In:     os.Stdin,
		Out:    os.Stdout,
		Raw:    true,
		TryDev: true,
	}
	if size := tty.GetSize(); size != nil {
		// fake resizing +1 and then back to normal so that attach-detach-reattach will result in the
		// screen being redrawn
		sizePlusOne := *size
		sizePlusOne.Width++
		sizePlusOne.Height++

		// this call spawns a goroutine to monitor/update the terminal size
		sizeQueue = tty.MonitorSize(&sizePlusOne, size)
	}

	exec, err := remotecommand.NewSPDYExecutor(cfg, http.MethodPost, req.URL())
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintln(os.Stderr, "If you don't see a command prompt, try pressing enter.")
	return tty.Safe(func() error {
		return exec.StreamWithContext(ctx, remotecommand.StreamOptions{
			Stdin:             tty.In,
			Stdout:            tty.Out,
			Tty:               true,
			TerminalSizeQueue: sizeQueue,
		})
	})
}
