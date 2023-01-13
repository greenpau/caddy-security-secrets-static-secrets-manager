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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/caddyserver/caddy/v2"
	"github.com/google/go-cmp/cmp"
)

func TestProvisionPlugin(t *testing.T) {
	testcases := []struct {
		name      string
		cfg       string
		want      Config
		shouldErr bool
		err       error
	}{
		{
			name: "test provisioning valid config",
			cfg:  `{"id":"foo","secret":{"foo":"bar","bar":"baz"}}`,
			want: Config{
				ID: "foo",
				Secret: map[string]interface{}{
					"foo": "bar",
					"bar": "baz",
				},
			},
		},
		{
			name:      "test provisioning malformed json config",
			cfg:       `{"id":"foo","secret":{"foo":"bar","bar":"baz"}`,
			shouldErr: true,
			err:       fmt.Errorf("unexpected end of JSON input"),
		},
		{
			name:      "test provisioning config without id",
			cfg:       `{"secret":{"foo":"bar","bar":"baz"}}`,
			shouldErr: true,
			err:       fmt.Errorf("id is empty"),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			p := &Plugin{
				ConfigRaw: json.RawMessage(tc.cfg),
			}
			err := p.Provision(caddy.ActiveContext())
			if err != nil {
				if !tc.shouldErr {
					t.Fatalf("expected success, got: %v", err)
				}
				if diff := cmp.Diff(tc.err.Error(), err.Error()); diff != "" {
					t.Logf("unexpected error: %v", err)
					t.Fatalf("Provision() error mismatch (-want +got):\n%s", diff)
				}
				return
			}
			if tc.shouldErr {
				t.Fatalf("unexpected success, want: %v", tc.err)
			}

			if diff := cmp.Diff(tc.want, p.Config); diff != "" {
				t.Logf("JSON: %s", p.ConfigRaw)
				t.Fatalf("Provision() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestValidatePlugin(t *testing.T) {
	testcases := []struct {
		name      string
		cfg       string
		want      map[string]interface{}
		shouldErr bool
		err       error
	}{
		{
			name: "test validating valid config",
			cfg:  `{"id":"foo","secret":{"foo":"bar","bar":"baz"}}`,
			want: map[string]interface{}{
				"id":       "foo",
				"provider": "static_secrets_manager",
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
				if !tc.shouldErr {
					t.Fatalf("expected success, got: %v", err)
				}
				if diff := cmp.Diff(tc.err.Error(), err.Error()); diff != "" {
					t.Logf("unexpected error: %v", err)
					t.Fatalf("Validate() error mismatch (-want +got):\n%s", diff)
				}
				return
			}
			if tc.shouldErr {
				t.Fatalf("unexpected success, want: %v", tc.err)
			}

			got := p.GetConfig(caddy.ActiveContext())
			if diff := cmp.Diff(tc.want, got); diff != "" {
				// t.Logf("JSON: %s", p.ConfigRaw)
				t.Errorf("Validate() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
