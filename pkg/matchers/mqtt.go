// Copyright 2020 Herman Slatman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package matchers

import (
	"encoding/json"

	"github.com/caddyserver/caddy/v2"
	"github.com/mholt/caddy-l4/layer4"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(MatchMQTT{})
}

// MatchMQTT is able to match MQTT connections.
type MatchMQTT struct {
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (MatchMQTT) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "layer4.matchers.mqtt",
		New: func() caddy.Module { return new(MatchMQTT) },
	}
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (m *MatchMQTT) UnmarshalJSON(b []byte) error {
	return nil
}

// MarshalJSON satisfies the json.Marshaler interface.
func (m MatchMQTT) MarshalJSON() ([]byte, error) {
	return []byte{}, nil
}

// Provision sets up the matcher.
func (m *MatchMQTT) Provision(ctx caddy.Context) error {

	m.logger = ctx.Logger(m)
	defer m.logger.Sync()

	m.logger.Debug("provisioning mqtt")

	return nil
}

// Match returns true if the conn starts with MQTT ......
func (m MatchMQTT) Match(cx *layer4.Connection) (bool, error) {

	m.logger.Debug("matching mqtt")

	// TODO: check first bytes of connection to see if it's MQTT

	return true, nil // We just assume everything is MQTTS (like a fail-open)
}

// Interface guards
var (
	_ caddy.Module       = (*MatchMQTT)(nil)
	_ caddy.Provisioner  = (*MatchMQTT)(nil)
	_ layer4.ConnMatcher = (*MatchMQTT)(nil)
	_ json.Marshaler     = (*MatchMQTT)(nil)
	_ json.Unmarshaler   = (*MatchMQTT)(nil)
)
