// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package bootstraptoken

import (
	"fmt"

	"github.com/ironcore-dev/kubectl-ironcore/utils/bootstraptoken"
	"k8s.io/apimachinery/pkg/util/sets"
)

type Type string

const (
	MachinePoolType   Type = "MachinePool"
	VolumePoolType    Type = "VolumePool"
	BucketPoolType    Type = "BucketPool"
	NetworkPluginType Type = "NetworkPlugin"
	APINetletType     Type = "APINetlet"
	MetalnetletType   Type = "Metalnetlet"
)

const (
	MachinePoolBootstrappersGroup   = "system:bootstrappers:compute-ironcore-dev:machinepools"
	VolumePoolBootstrappersGroup    = "system:bootstrappers:storage-ironcore-dev:volumepools"
	BucketPoolBootstrappersGroup    = "system:bootstrappers:storage-ironcore-dev:bucketpools"
	NetworkPluginBootstrappersGroup = "system:bootstrappers:networking-ironcore-dev:networkplugins"
	APINetletBootstrappersGroup     = "system:bootstrappers:apinet-ironcore-dev:apinetlets"
	MetalnetletBootstrappersGroup   = "system:bootstrappers:apinet-ironcore-dev:metalnetlets"
)

var AvailableTypes = sets.New[Type](
	MachinePoolType,
	VolumePoolType,
	BucketPoolType,
	NetworkPluginType,
	APINetletType,
	MetalnetletType,
)

type fields struct {
	Description string
	Usages      []string
	Groups      []string
}

var fieldsByType = map[Type]fields{
	MachinePoolType: {
		Description: "Bootstrap token for registering machine pools.",
		Usages: []string{
			bootstraptoken.UsageSigning,
			bootstraptoken.UsageAuthentication,
		},
		Groups: []string{
			MachinePoolBootstrappersGroup,
		},
	},
	VolumePoolType: {
		Description: "Bootstrap token for registering volume pools.",
		Usages: []string{
			bootstraptoken.UsageSigning,
			bootstraptoken.UsageAuthentication,
		},
		Groups: []string{
			VolumePoolBootstrappersGroup,
		},
	},
	BucketPoolType: {
		Description: "Bootstrap token for registering bucket pools.",
		Usages: []string{
			bootstraptoken.UsageSigning,
			bootstraptoken.UsageAuthentication,
		},
		Groups: []string{
			BucketPoolBootstrappersGroup,
		},
	},
	NetworkPluginType: {
		Description: "Bootstrap token for registering network plugins.",
		Usages: []string{
			bootstraptoken.UsageSigning,
			bootstraptoken.UsageAuthentication,
		},
		Groups: []string{
			NetworkPluginBootstrappersGroup,
		},
	},
	APINetletType: {
		Description: "Bootstrap token for registering apinetlets.",
		Usages: []string{
			bootstraptoken.UsageSigning,
			bootstraptoken.UsageAuthentication,
		},
		Groups: []string{
			APINetletBootstrappersGroup,
		},
	},
	MetalnetletType: {
		Description: "Bootstrap token for registering metalnetlets.",
		Usages: []string{
			bootstraptoken.UsageSigning,
			bootstraptoken.UsageAuthentication,
		},
		Groups: []string{
			MetalnetletBootstrappersGroup,
		},
	},
}

func AddTypeFields(bt *bootstraptoken.BootstrapToken, typ Type) error {
	flds, ok := fieldsByType[typ]
	if !ok {
		return fmt.Errorf("unknown type %q", typ)
	}

	if bt.Description == "" {
		bt.Description = flds.Description
	}

	presentUsages := sets.New(bt.Usages...)
	for _, usage := range flds.Usages {
		if !presentUsages.Has(usage) {
			bt.Usages = append(bt.Usages, usage)
		}
	}

	presentGroups := sets.New(bt.Groups...)
	for _, group := range flds.Groups {
		if !presentGroups.Has(group) {
			bt.Groups = append(bt.Groups, group)
		}
	}
	return nil
}
