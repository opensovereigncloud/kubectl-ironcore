package bootstrapkubeconfig

import (
	"bytes"
	"fmt"
	"io"
	"os"

	utilbootstraptoken "github.com/onmetal/kubectl-onmetal/utils/bootstraptoken"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Flags struct {
	Filename     string
	NoFlatten    bool
	ConfigAccess clientcmd.ConfigAccess
	genericclioptions.IOStreams
}

func NewFlags(configAccess clientcmd.ConfigAccess, streams genericclioptions.IOStreams) *Flags {
	return &Flags{
		ConfigAccess: configAccess,
		IOStreams:    streams,
	}
}

func (f *Flags) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.Filename, "filename", "f", "", "File to read for bootstrap token secret. Specify '-' for using stdin.")
	cmd.Flags().BoolVar(&f.NoFlatten, "no-flatten", false, "Whether to skip flattening of the resulting kubeconfig.")
}

func (f *Flags) ToOptions() (*Options, error) {
	if f.Filename == "" {
		return nil, fmt.Errorf("must specify filename")
	}

	return &Options{
		Filename:     f.Filename,
		ConfigAccess: f.ConfigAccess,
		IOStreams:    f.IOStreams,
		NoFlatten:    f.NoFlatten,
	}, nil
}

type Options struct {
	Filename     string
	ConfigAccess clientcmd.ConfigAccess
	Context      string
	genericclioptions.IOStreams
	NoFlatten bool
}

func Command(configAccess clientcmd.ConfigAccess, streams genericclioptions.IOStreams) *cobra.Command {
	flags := NewFlags(configAccess, streams)

	cmd := &cobra.Command{
		Use:   "bootstrap-kubeconfig",
		Short: "Generate a bootstrap-kubeconfig from a bootstrap-token secret and a kubeconfig.",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts, err := flags.ToOptions()
			if err != nil {
				return err
			}

			return Run(*opts)
		},
	}

	flags.AddFlags(cmd)

	return cmd
}

func readFileOrStdin(filename string, streams genericclioptions.IOStreams) ([]byte, error) {
	if filename == "-" {
		return io.ReadAll(streams.In)
	}
	return os.ReadFile(filename)
}

func Run(opts Options) error {
	secretData, err := readFileOrStdin(opts.Filename, opts.IOStreams)
	if err != nil {
		return fmt.Errorf("error reading secret: %w", err)
	}

	secret := &corev1.Secret{}
	if err := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(secretData), 4096).Decode(secret); err != nil {
		return fmt.Errorf("error decoding secret: %w", err)
	}

	token, err := utilbootstraptoken.FromSecret(secret)
	if err != nil {
		return fmt.Errorf("error decoding bootstrap token from secret: %w", err)
	}

	startingCfg, err := opts.ConfigAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	contextName := startingCfg.CurrentContext
	if opts.Context != "" {
		contextName = opts.Context
	}
	if contextName == "" {
		return fmt.Errorf("must specify context or current context must exist in kubeconfig")
	}

	context, ok := startingCfg.Contexts[startingCfg.CurrentContext]
	if !ok {
		return fmt.Errorf("context %s not found", startingCfg.CurrentContext)
	}

	cluster, ok := startingCfg.Clusters[context.Cluster]
	if !ok {
		return fmt.Errorf("cluster %s not found", context.Cluster)
	}

	apiCfg := &clientcmdapi.Config{
		Preferences: startingCfg.Preferences,
		Clusters: map[string]*clientcmdapi.Cluster{
			"bootstrap": cluster,
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			"bootstrap": {
				Token: fmt.Sprintf("%s.%s", token.ID, token.Secret),
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			"bootstrap": {
				AuthInfo:   "bootstrap",
				Cluster:    "bootstrap",
				Namespace:  context.Namespace,
				Extensions: context.Extensions,
			},
		},
		CurrentContext: "bootstrap",
		Extensions:     startingCfg.Extensions,
	}
	if !opts.NoFlatten {
		if err := clientcmdapi.FlattenConfig(apiCfg); err != nil {
			return err
		}
	}

	apiCfgData, err := clientcmd.Write(*apiCfg)
	if err != nil {
		return err
	}

	_, _ = io.Copy(opts.Out, bytes.NewReader(apiCfgData))
	return nil
}
