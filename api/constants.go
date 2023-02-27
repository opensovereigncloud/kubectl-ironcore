package api

import "sigs.k8s.io/controller-runtime/pkg/client"

const (
	// FieldOwner is the field owner kubectl-onmetal uses.
	FieldOwner = client.FieldOwner("api.onmetal.de/kubectl-onmetal")
)
