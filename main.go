// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
