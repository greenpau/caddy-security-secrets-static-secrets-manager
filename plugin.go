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

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	static_secrets_manager "github.com/greenpau/go-authcrunch-secrets-static-secrets-manager"
	"go.uber.org/zap"
)

var (
	pluginPrefix = "security.secrets"
	pluginName   = pluginPrefix + ".static_secrets_manager"

	// Interface guards
	_ caddy.Provisioner     = (*Plugin)(nil)
	_ caddy.Validator       = (*Plugin)(nil)
	_ caddyfile.Unmarshaler = (*Plugin)(nil)
	_ caddy.Module          = (*Plugin)(nil)
)

func init() {
	caddy.RegisterModule(Plugin{})
}

// Config represents provisioned configuration value of statically configured secret.
type Config struct {
	ID     string                 `json:"id,omitempty" xml:"id,omitempty" yaml:"id,omitempty"`
	Secret map[string]interface{} `json:"secret,omitempty" xml:"secret,omitempty" yaml:"secret,omitempty"`
}

// Plugin manages statically configured secrets.
type Plugin struct {
	Name      string          `json:"-"`
	ConfigRaw json.RawMessage `json:"config,omitempty" caddy:"namespace=security.secrets.static_secrets_manager"`
	Config    Config          `json:"-"`
	client    static_secrets_manager.Client
	logger    *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (Plugin) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "security.secrets.static_secrets_manager",
		New: func() caddy.Module { return new(Plugin) },
	}
}

// Provision sets up Handler and loads AwsSecretsManager.
func (p *Plugin) Provision(ctx caddy.Context) error {
	p.Name = pluginName
	p.logger = ctx.Logger(p)

	p.logger.Info(
		"provisioning plugin instance",
		zap.String("plugin_name", p.Name),
	)

	if err := json.Unmarshal(p.ConfigRaw, &p.Config); err != nil {
		p.logger.Error(
			"failed configuring plugin instance",
			zap.String("plugin_name", p.Name),
			zap.Error(err),
		)
		return err
	}

	client, err := static_secrets_manager.NewClient(ctx, p.Config.ID, p.Config.Secret)
	if err != nil {
		p.logger.Error(
			"failed initializing secrets manager client",
			zap.String("plugin_name", p.Name),
			zap.Error(err),
		)
		return err
	}

	p.client = client

	p.logger.Info(
		"provisioned plugin instance",
		zap.String("plugin_name", p.Name),
	)
	return nil
}

// Validate implements caddy.Validator.
func (p *Plugin) Validate() error {
	p.logger.Info(
		"validating plugin instance",
		zap.String("plugin_name", p.Name),
		zap.String("secret_id", p.Config.ID),
	)

	p.logger.Info(
		"validated plugin instance",
		zap.String("plugin_name", p.Name),
		zap.String("secret_id", p.Config.ID),
	)
	return nil
}

// GetConfig returns plugin configuration.
func (p *Plugin) GetConfig(ctx context.Context) map[string]interface{} {
	m := p.client.GetConfig(ctx)
	return m
}
