// Copyright 2023 OnMetal authors
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

package bootstraptoken

import (
	"fmt"

	"github.com/onmetal/kubectl-onmetal/utils/bootstraptoken"
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
	MachinePoolBootstrappersGroup   = "system:bootstrappers:compute-api-onmetal-de:machinepools"
	VolumePoolBootstrappersGroup    = "system:bootstrappers:storage-api-onmetal-de:volumepools"
	BucketPoolBootstrappersGroup    = "system:bootstrappers:storage-api-onmetal-de:bucketpools"
	NetworkPluginBootstrappersGroup = "system:bootstrappers:networking-api-onmetal-de:networkplugins"
	APINetletBootstrappersGroup     = "system:bootstrappers:apinet-api-onmetal-de:apinetlets"
	MetalnetletBootstrappersGroup   = "system:bootstrappers:apinet-api-onmetal-de:metalnetlets"
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
