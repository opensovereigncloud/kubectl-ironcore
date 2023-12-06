// Copyright 2021 IronCore authors
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

package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-logr/zapr"
	kubectlironcore "github.com/ironcore-dev/kubectl-ironcore/cmd/kubectl-ironcore"
	"go.uber.org/zap"
)

func main() {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func() { _ = zapLog.Sync() }()

	setupLog := zapr.NewLogger(zapLog)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := kubectlironcore.DefaultCommand().ExecuteContext(ctx); err != nil {
		setupLog.Error(err, "Error running command")
		os.Exit(1)
	}
}
