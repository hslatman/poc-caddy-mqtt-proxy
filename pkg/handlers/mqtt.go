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

package handlers

import (
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/hslatman/poc-caddy-mqtt-proxy/internal/conn"
	"github.com/hslatman/poc-caddy-mqtt-proxy/pkg/matchers"
	"github.com/mholt/caddy-l4/layer4"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(MQTT{})
}

// MQTT is able to match MQTT connections.
type MQTT struct {
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (MQTT) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "layer4.handlers.mqtt", // TODO: need to think about this name a bit more; it's too similar to the matchers
		New: func() caddy.Module { return new(MQTT) },
	}
}

// Provision sets up the matcher.
func (m *MQTT) Provision(ctx caddy.Context) error {

	m.logger = ctx.Logger(m)
	defer m.logger.Sync()

	return nil
}

// Handle handles the downstream connection.
func (m *MQTT) Handle(cx *layer4.Connection, next layer4.Handler) error {

	matcher := &matchers.MatchMQTT{}

	wrappedConnection := conn.NewBufferedConn(cx.Conn)

	isValidMQTTConnection, err := matcher.CheckBufferedConnection(wrappedConnection)
	if err != nil {
		m.logger.Error(err.Error())
	}

	// terminate the connection if it could not be determined to be MQTT
	if !isValidMQTTConnection {
		m.logger.Debug(fmt.Sprintf("terminating connection from %s", cx.Conn.RemoteAddr().String()))
		return cx.Conn.Close()
	}

	m.logger.Debug(fmt.Sprintf("accepted connection from %s", cx.Conn.RemoteAddr().String()))

	cx.Conn = wrappedConnection

	return next.Handle(cx)
}

// Interface guards
var (
	_ caddy.Module       = (*MQTT)(nil)
	_ caddy.Provisioner  = (*MQTT)(nil)
	_ layer4.NextHandler = (*MQTT)(nil)
)
