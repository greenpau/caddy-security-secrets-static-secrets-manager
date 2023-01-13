// Copyright 2022 Paul Greenberg greenpau@outlook.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package staticsecretsmanager

import (
	"context"
)

// GetSecret returns a secret in the form of a key-value map.
func (p *Plugin) GetSecret(ctx context.Context) (map[string]interface{}, error) {
	return p.client.GetSecret(ctx)
}

// GetSecretByKey returns a value of key in the secret key-value map.
func (p *Plugin) GetSecretByKey(ctx context.Context, key string) (interface{}, error) {
	return p.client.GetSecretByKey(ctx, key)
}
