// SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

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
