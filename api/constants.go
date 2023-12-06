// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package api

import "sigs.k8s.io/controller-runtime/pkg/client"

const (
	// FieldOwner is the field owner kubectl-ironcore uses.
	FieldOwner = client.FieldOwner("api.ironcore.dev/kubectl-ironcore")
)
