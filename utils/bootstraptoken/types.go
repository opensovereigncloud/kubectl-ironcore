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

import "time"

type BootstrapToken struct {
	ID          string
	Secret      string
	Description string
	Expires     *time.Time
	Usages      []string
	Groups      []string
}

func (t *BootstrapToken) WithTTL(ttl time.Duration) *BootstrapToken {
	expires := time.Now().Add(ttl)
	return &BootstrapToken{
		ID:          t.ID,
		Secret:      t.Secret,
		Description: t.Description,
		Expires:     &expires,
		Usages:      t.Usages,
		Groups:      t.Groups,
	}
}
