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
	MachinePoolType Type = "MachinePool"
	VolumePoolType  Type = "VolumePool"
	BucketPoolType  Type = "BucketPool"
)

const (
	MachinePoolBootstrappersGroup = "system:bootstrappers:compute-api-onmetal-de:machinepools"
	VolumePoolBootstrappersGroup  = "system:bootstrappers:storage-api-onmetal-de:volumepools"
	BucketPoolBootstrappersGroup  = "system:bootstrappers:storage-api-onmetal-de:bucketpools"
)

var AvailableTypes = sets.New[Type](
	MachinePoolType,
	VolumePoolType,
	BucketPoolType,
)

type fields struct {
	Usages []string
	Groups []string
}

var fieldsByType = map[Type]fields{
	MachinePoolType: {
		Usages: []string{
			bootstraptoken.UsageSigning,
			bootstraptoken.UsageAuthentication,
		},
		Groups: []string{
			MachinePoolBootstrappersGroup,
		},
	},
	VolumePoolType: {
		Usages: []string{
			bootstraptoken.UsageSigning,
			bootstraptoken.UsageAuthentication,
		},
		Groups: []string{
			VolumePoolBootstrappersGroup,
		},
	},
	BucketPoolType: {
		Usages: []string{
			bootstraptoken.UsageSigning,
			bootstraptoken.UsageAuthentication,
		},
		Groups: []string{
			BucketPoolBootstrappersGroup,
		},
	},
}

func AddTypeFields(bt *bootstraptoken.BootstrapToken, typ Type) error {
	flds, ok := fieldsByType[typ]
	if !ok {
		return fmt.Errorf("unknown type %q", typ)
	}

	presentUsages := sets.New(bt.Usages...)
	presentGroups := sets.New(bt.Groups...)

	for _, usage := range flds.Usages {
		if !presentUsages.Has(usage) {
			bt.Usages = append(bt.Usages, usage)
		}
	}
	for _, group := range flds.Groups {
		if !presentGroups.Has(group) {
			bt.Groups = append(bt.Groups, group)
		}
	}
	return nil
}
