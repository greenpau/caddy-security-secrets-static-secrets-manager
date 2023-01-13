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
	"encoding/json"
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/google/go-cmp/cmp"
)

// packMapToJSON converts a map to a JSON string.
func packMapToJSON(t *testing.T, m map[string]interface{}) string {
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("failed to marshal %v: %v", m, err)
	}
	return string(b)
}

func TestGetSecret(t *testing.T) {
	testcases := []struct {
		name      string
		cfg       string
		secret    map[string]interface{}
		want      map[string]interface{}
		shouldErr bool
		err       error
	}{
		{
			name: "test get valid secret",
			cfg:  `{"id":"foo","secret":{"foo":"bar","bar":"baz"}}`,
			secret: map[string]interface{}{
				"foo": "bar",
				"bar": "baz",
			},
			want: map[string]interface{}{
				"foo": "bar",
				"bar": "baz",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			p := &Plugin{
				ConfigRaw: json.RawMessage(tc.cfg),
			}
			err := p.Provision(caddy.ActiveContext())
			if err != nil {
				t.Fatalf("unexpected provisioning error: %v", err)
			}

			err = p.Validate()
			if err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}

			got, err := p.GetSecret(context.TODO())
			if err != nil {
				if !tc.shouldErr {
					t.Fatalf("expected success, got: %v", err)
				}
				if diff := cmp.Diff(tc.err.Error(), err.Error()); diff != "" {
					t.Logf("unexpected error: %v", err)
					t.Fatalf("GetSecret() error mismatch (-want +got):\n%s", diff)
				}
				return
			}
			if tc.shouldErr {
				t.Fatalf("unexpected success, want: %v", tc.err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GetSecret() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetSecretByKey(t *testing.T) {
	testcases := []struct {
		name      string
		cfg       string
		secret    map[string]interface{}
		key       string
		want      interface{}
		shouldErr bool
		err       error
	}{
		{
			name: "test get valid secret by key",
			cfg:  `{"id":"foo","secret":{"foo":"bar","bar":"baz"}}`,
			secret: map[string]interface{}{
				"foo": "bar",
				"bar": "baz",
			},
			key:  "foo",
			want: "bar",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			p := &Plugin{
				ConfigRaw: json.RawMessage(tc.cfg),
			}
			err := p.Provision(caddy.ActiveContext())
			if err != nil {
				t.Fatalf("unexpected provisioning error: %v", err)
			}

			err = p.Validate()
			if err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}

			got, err := p.GetSecretByKey(context.TODO(), tc.key)
			if err != nil {
				if !tc.shouldErr {
					t.Fatalf("expected success, got: %v", err)
				}
				if diff := cmp.Diff(tc.err.Error(), err.Error()); diff != "" {
					t.Logf("unexpected error: %v", err)
					t.Fatalf("GetSecretByKey() error mismatch (-want +got):\n%s", diff)
				}
				return
			}
			if tc.shouldErr {
				t.Fatalf("unexpected success, want: %v", tc.err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("GetSecretByKey() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
