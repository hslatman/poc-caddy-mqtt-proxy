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
	caddy.RegisterModule(MatchMQTTS{})
}

// MatchMQTTS is able to match MQTTS connections.
type MatchMQTTS struct {
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (MatchMQTTS) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "layer4.matchers.mqtts",
		New: func() caddy.Module { return new(MatchMQTTS) },
	}
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (m *MatchMQTTS) UnmarshalJSON(b []byte) error {
	return nil
}

// MarshalJSON satisfies the json.Marshaler interface.
func (m MatchMQTTS) MarshalJSON() ([]byte, error) {
	return []byte{}, nil
}

// Provision sets up the matcher.
func (m *MatchMQTTS) Provision(ctx caddy.Context) error {

	m.logger = ctx.Logger(m)
	defer m.logger.Sync()

	m.logger.Debug("provisioning mqtts")

	// TODO: Add logic similar to the TLS matcher (with the SNI matching) in here?
	// Due to the way the matching now works, I think we can't know if we're actually
	// dealing with (encrypted) MQTT traffic. This is because the connection has not
	// been terminated by Caddy yet (that's done in the TLS handler) and so this
	// information is not available at this point.

	return nil
}

// Match returns true if the conn starts with MQTT ......
func (m MatchMQTTS) Match(cx *layer4.Connection) (bool, error) {

	m.logger.Debug("matching mqtts")

	return true, nil // We just assume everything is MQTTS (like a fail-open)
}

// Interface guards
var (
	_ caddy.Module       = (*MatchMQTTS)(nil)
	_ caddy.Provisioner  = (*MatchMQTTS)(nil)
	_ layer4.ConnMatcher = (*MatchMQTTS)(nil)
	_ json.Marshaler     = (*MatchMQTTS)(nil)
	_ json.Unmarshaler   = (*MatchMQTTS)(nil)
)
