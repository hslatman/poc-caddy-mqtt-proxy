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
	"bufio"
	"bytes"
	"fmt"
	"net"

	"github.com/caddyserver/caddy/v2"
	"github.com/mholt/caddy-l4/layer4"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(MQTTMatcher{})
}

// MQTTMatcher is able to match MQTT connections.
type MQTTMatcher struct {
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (MQTTMatcher) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "layer4.handlers.mqttmatcher", // TODO: need to think about this name a bit more; it's too similar to the matchers
		New: func() caddy.Module { return new(MQTTMatcher) },
	}
}

// Provision sets up the matcher.
func (m *MQTTMatcher) Provision(ctx caddy.Context) error {

	m.logger = ctx.Logger(m)
	defer m.logger.Sync()

	return nil
}

// Handle handles the downstream connection.
func (m *MQTTMatcher) Handle(cx *layer4.Connection, next layer4.Handler) error {

	fmt.Println(cx)
	fmt.Println(next)

	b := newBufferedConn(cx.Conn)
	// p, err := b.Peek(20)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(p)
	// fmt.Println(string(p))

	// cx.Conn = b

	//p := make([]byte, len(mqttPrefix))
	//p = make([]byte, 30)
	//n, err := io.ReadFull(b, p)

	n, err := b.Peek(len(mqttPrefix))
	if err != nil {
		fmt.Println("error")
	}

	fmt.Println(n)
	fmt.Println(string(n))

	isMQTT := bytes.Equal(n, mqttPrefix)

	fmt.Println(isMQTT)

	cx.Conn = b

	// TODO: terminate connection if not MQTT?

	return next.Handle(cx)
}

// 16 46 0 4

//var mqttPrefix = []byte(".MQTT")
var mqttPrefix = []byte{16, 46, 0, 4} // TODO: make this more generic; other MQTT versions etc.

type bufferedConn struct {
	r        *bufio.Reader
	net.Conn // So that most methods are embedded
}

func newBufferedConn(c net.Conn) bufferedConn {
	return bufferedConn{bufio.NewReader(c), c}
}

func newBufferedConnSize(c net.Conn, n int) bufferedConn {
	return bufferedConn{bufio.NewReaderSize(c, n), c}
}

func (b bufferedConn) Peek(n int) ([]byte, error) {
	return b.r.Peek(n)
}

func (b bufferedConn) Read(p []byte) (int, error) {
	return b.r.Read(p)
}

// Interface guards
var (
	_ caddy.Module       = (*MQTTMatcher)(nil)
	_ caddy.Provisioner  = (*MQTTMatcher)(nil)
	_ layer4.NextHandler = (*MQTTMatcher)(nil)
)
