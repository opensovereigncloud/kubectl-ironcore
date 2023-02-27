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
	"bufio"
	"crypto/rand"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cluster-bootstrap/token/api"
	"k8s.io/cluster-bootstrap/token/util"
)

const (
	// validBootstrapTokenChars defines the characters a bootstrap token can consist of
	validBootstrapTokenChars = "0123456789abcdefghijklmnopqrstuvwxyz"

	SecretPattern = `\A([a-z0-9]{16})\z`

	UsageSigning        = "signing"
	UsageAuthentication = "authentication"
)

var (
	idRegexp     = regexp.MustCompile(api.BootstrapTokenIDPattern)
	secretRegexp = regexp.MustCompile(SecretPattern)
)

func ValidateID(id string) error {
	if idRegexp.MatchString(id) {
		return nil
	}
	return fmt.Errorf("bootstrap token id %q is invalid (must match string %s)", id, api.BootstrapTokenIDPattern)
}

func IsValidID(id string) bool {
	return ValidateID(id) == nil
}

func ValidateSecret(secret string) error {
	if secretRegexp.MatchString(secret) {
		return nil
	}
	return fmt.Errorf("secret %q is invalid (must match string %s)", secret, SecretPattern)
}

func IsValidSecret(secret string) bool {
	return ValidateSecret(secret) == nil
}

// randBytes returns a random string consisting of the characters in
// validBootstrapTokenChars, with the length customized by the parameter
func randBytes(length int) (string, error) {
	// len(validBootstrapTokenChars) = 36 which doesn't evenly divide
	// the possible values of a byte: 256 mod 36 = 4. Discard any random bytes we
	// read that are >= 252 so the bytes we evenly divide the character set.
	const maxByteValue = 252

	var (
		b     byte
		err   error
		token = make([]byte, length)
	)

	reader := bufio.NewReaderSize(rand.Reader, length*2)
	for i := range token {
		for {
			if b, err = reader.ReadByte(); err != nil {
				return "", err
			}
			if b < maxByteValue {
				break
			}
		}

		token[i] = validBootstrapTokenChars[int(b)%len(validBootstrapTokenChars)]
	}

	return string(token), nil
}

func Generate(template *BootstrapToken) (*BootstrapToken, error) {
	var token BootstrapToken
	if template != nil {
		token = *template
	}

	if token.ID != "" {
		if err := ValidateID(token.ID); err != nil {
			return nil, err
		}
	} else {
		id, err := randBytes(api.BootstrapTokenIDBytes)
		if err != nil {
			return nil, fmt.Errorf("error generating id: %w", err)
		}

		token.ID = id
	}

	if token.Secret != "" {
		if err := ValidateSecret(token.Secret); err != nil {
			return nil, err
		}
	} else {
		secret, err := randBytes(api.BootstrapTokenSecretBytes)
		if err != nil {
			return nil, fmt.Errorf("error generating secret: %w", err)
		}

		token.Secret = secret
	}

	if err := util.ValidateUsages(token.Usages); err != nil {
		return nil, err
	}

	return &token, nil
}

func ToSecret(token *BootstrapToken) *corev1.Secret {
	data := map[string][]byte{
		api.BootstrapTokenIDKey:     []byte(token.ID),
		api.BootstrapTokenSecretKey: []byte(token.Secret),
	}

	if token.Description != "" {
		data[api.BootstrapTokenDescriptionKey] = []byte(token.Description)
	}
	if token.Expires != nil {
		expirationString := token.Expires.UTC().Format(time.RFC3339)
		data[api.BootstrapTokenExpirationKey] = []byte(expirationString)
	}

	for _, usage := range token.Usages {
		data[api.BootstrapTokenUsagePrefix+usage] = []byte("true")
	}

	if len(token.Groups) > 0 {
		data[api.BootstrapTokenExtraGroupsKey] = []byte(strings.Join(token.Groups, ","))
	}

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: metav1.NamespaceSystem,
			Name:      util.BootstrapTokenSecretName(token.ID),
		},
		Type: corev1.SecretTypeBootstrapToken,
		Data: data,
	}
}

func FromSecret(secret *corev1.Secret) (*BootstrapToken, error) {
	tokenID := string(secret.Data[api.BootstrapTokenIDKey])
	if err := ValidateID(tokenID); err != nil {
		return nil, err
	}

	if expectedSecretName := util.BootstrapTokenSecretName(tokenID); secret.Name != expectedSecretName {
		return nil, fmt.Errorf("bootstrap token name is not of the form '%s(token-id)' - actual: %q, expected: %q",
			api.BootstrapTokenSecretPrefix, secret.Name, expectedSecretName)
	}

	tokenSecret := string(secret.Data[api.BootstrapTokenSecretKey])
	if err := ValidateSecret(tokenSecret); err != nil {
		return nil, err
	}

	description := string(secret.Data[api.BootstrapTokenDescriptionKey])

	var expires *time.Time
	if expirationData, ok := secret.Data[api.BootstrapTokenExpirationKey]; ok {
		e, err := time.Parse(time.RFC3339, string(expirationData))
		if err != nil {
			return nil, fmt.Errorf("bootstrap token expiration is invalid: %w", err)
		}

		expires = &e
	}

	var usages []string
	for k, v := range secret.Data {
		if !strings.HasPrefix(k, api.BootstrapTokenUsagePrefix) {
			continue
		}

		// Skip those that don't have this usage set to true
		if string(v) != "true" {
			continue
		}

		usages = append(usages, strings.TrimPrefix(k, api.BootstrapTokenUsagePrefix))
	}

	if usages != nil {
		sort.Strings(usages)
	}

	var groups []string
	groupsString := string(secret.Data[api.BootstrapTokenExtraGroupsKey])
	g := strings.Split(groupsString, ",")
	if len(g) > 0 && len(g[0]) > 0 {
		groups = g
	}

	return &BootstrapToken{
		ID:          tokenID,
		Secret:      tokenSecret,
		Description: description,
		Expires:     expires,
		Usages:      usages,
		Groups:      groups,
	}, nil
}
