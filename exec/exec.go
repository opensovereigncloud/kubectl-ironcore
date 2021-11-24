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
	"time"

	"github.com/onmetal/onmetal-api/apis/compute/v1alpha1"
	"github.com/onmetal/onmetal-console/http"
	"github.com/onmetal/onmetal-console/tty"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const fieldOwner = client.FieldOwner("kubectl-onmetal.onmetal.de/console")

func Run(ctx context.Context, c client.Client, tty tty.TTY, namespace string, name string) error {
	log := controllerruntime.LoggerFrom(ctx)
	log.Info("Applying console")
	console := &v1alpha1.Console{
		TypeMeta: v1.TypeMeta{
			APIVersion: v1alpha1.GroupVersion.String(),
			Kind:       "Console",
		},
		ObjectMeta: v1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Spec: v1alpha1.ConsoleSpec{
			Type: v1alpha1.ConsoleTypeService,
			MachineRef: v12.LocalObjectReference{
				Name: name,
			},
		},
	}
	if err := c.Patch(ctx, console, client.Apply, fieldOwner); err != nil {
		return fmt.Errorf("error applying console: %w", err)
	}
	consoleKey := client.ObjectKeyFromObject(console)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := c.Delete(ctx, &v1alpha1.Console{
			ObjectMeta: v1.ObjectMeta{Namespace: consoleKey.Namespace, Name: consoleKey.Name},
		}); client.IgnoreNotFound(err) != nil {
			log.Error(err, "Error deleting console")
		}
	}()

	if err := wait.PollImmediateUntilWithContext(ctx, 1*time.Second, func(ctx context.Context) (done bool, err error) {
		if err := c.Get(ctx, consoleKey, console); err != nil {
			return false, fmt.Errorf("error getting console: %w", err)
		}

		switch console.Status.State {
		case v1alpha1.ConsoleStateReady:
			return true, nil
		case "", v1alpha1.ConsoleStatePending:
			log.V(1).Info("Console is pending")
			return false, nil
		case v1alpha1.ConsoleStateError:
			return false, fmt.Errorf("console errored")
		default:
			return false, fmt.Errorf("unknown console state %q", console.Status.State)
		}
	}); err != nil {
		return fmt.Errorf("error waiting for console to become ready: %w", err)
	}

	consoleClient, err := http.NewClient(http.ClientOptions{})
	if err != nil {
		return fmt.Errorf("error creating console client: %w", err)
	}

	clearScreen(tty)
	defer clearScreen(tty)
	return consoleClient.Run(ctx, tty, *console.Status.ServiceClientConfig.URL)
}

func clearScreen(t tty.TTY) {
	_, _ = t.Write([]byte("\033[H\033[2J"))
}
