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

package token

import (
	"context"
	"fmt"
	"time"

	"github.com/ironcore-dev/kubectl-ironcore/api"
	"github.com/ironcore-dev/kubectl-ironcore/bootstraptoken"
	utilbootstraptoken "github.com/ironcore-dev/kubectl-ironcore/utils/bootstraptoken"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Flags struct {
	Factory    cmdutil.Factory
	Template   utilbootstraptoken.BootstrapToken
	Type       bootstraptoken.Type
	TTL        time.Duration
	PrintFlags *genericclioptions.PrintFlags
	genericclioptions.IOStreams
}

func NewFlags(f cmdutil.Factory, streams genericclioptions.IOStreams) *Flags {
	printFlags := genericclioptions.NewPrintFlags("created").
		WithTypeSetter(scheme.Scheme)

	return &Flags{
		Factory:    f,
		PrintFlags: printFlags,
		IOStreams:  streams,
	}
}

func (f *Flags) AddFlags(cmd *cobra.Command) {
	cmdutil.AddDryRunFlag(cmd)
	cmd.Flags().StringVar((*string)(&f.Type), "token-type", "", fmt.Sprintf("Token type fields to add. Available types: %v", sets.List(bootstraptoken.AvailableTypes)))
	cmd.Flags().StringVar(&f.Template.ID, "token-id", "", "Token ID to use to generate.")
	cmd.Flags().StringVar(&f.Template.Secret, "token-secret", "", "Token secret to use to generate.")
	cmd.Flags().StringVar(&f.Template.Description, "token-description", "", "Token description to use to generate.")
	cmd.Flags().StringSliceVar(&f.Template.Groups, "token-groups", nil, "Additional token groups.")
	cmd.Flags().StringSliceVar(&f.Template.Usages, "token-usages", nil, "Additional token usages.")
	cmd.Flags().DurationVar(&f.TTL, "token-ttl", 0, "TTL for the token to expire. If unset, token will not expire.")
	f.PrintFlags.AddFlags(cmd)
}

func (f *Flags) ToOptions(cmd *cobra.Command) (*Options, error) {
	template := f.Template
	if f.TTL > 0 {
		template = *f.Template.WithTTL(f.TTL)
	}
	if f.Type != "" {
		if err := bootstraptoken.AddTypeFields(&template, f.Type); err != nil {
			return nil, err
		}
	}

	dryRunStrategy, err := cmdutil.GetDryRunStrategy(cmd)
	if err != nil {
		return nil, err
	}

	cmdutil.PrintFlagsWithDryRunStrategy(f.PrintFlags, dryRunStrategy)
	printer, err := f.PrintFlags.ToPrinter()
	if err != nil {
		return nil, err
	}

	namespace, _, err := f.Factory.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return nil, err
	}

	cfg, err := f.Factory.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	newClient := func() (client.Client, error) {
		return client.New(cfg, client.Options{})
	}

	return &Options{
		DryRun:    dryRunStrategy,
		Printer:   printer,
		Template:  template,
		Namespace: namespace,
		NewClient: newClient,
		IOStreams: f.IOStreams,
	}, nil
}

type Options struct {
	DryRun    cmdutil.DryRunStrategy
	Printer   printers.ResourcePrinter
	Template  utilbootstraptoken.BootstrapToken
	Namespace string
	NewClient func() (client.Client, error)
	genericclioptions.IOStreams
}

func Command(f cmdutil.Factory, streams genericclioptions.IOStreams) *cobra.Command {
	flags := NewFlags(f, streams)

	cmd := &cobra.Command{
		Use:   "token",
		Short: "Create a bootstrap token in a cluster.",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := flags.ToOptions(cmd)
			if err != nil {
				return err
			}

			return Run(cmd.Context(), *opts)
		},
	}

	flags.AddFlags(cmd)

	return cmd
}

func Run(ctx context.Context, opts Options) error {
	t, err := utilbootstraptoken.Generate(&opts.Template)
	if err != nil {
		return fmt.Errorf("error generating token: %w", err)
	}

	secret := utilbootstraptoken.ToSecret(t)

	patchOpts := []client.PatchOption{client.ForceOwnership, api.FieldOwner}
	if opts.DryRun == cmdutil.DryRunServer {
		patchOpts = append(patchOpts, client.DryRunAll)
	}

	if opts.DryRun != cmdutil.DryRunClient {
		c, err := opts.NewClient()
		if err != nil {
			return err
		}
		if err := c.Patch(ctx, secret, client.Apply, patchOpts...); err != nil {
			return err
		}
	}

	if err := opts.Printer.PrintObj(secret, opts.Out); err != nil {
		return fmt.Errorf("error printing object: %w", err)
	}
	return nil
}
