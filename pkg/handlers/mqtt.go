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
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"

	"github.com/caddyserver/caddy/v2"
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

	fmt.Println(cx)
	fmt.Println(next)

	b := newBufferedConn(cx.Conn)

	isValidMQTTConnection, err := checkConnection(b)
	if err != nil {
		m.logger.Warn(err.Error())
	}

	if !isValidMQTTConnection {
		// TODO: terminate connection if not MQTT?
	}

	m.logger.Debug(fmt.Sprintf("%t", isValidMQTTConnection))

	cx.Conn = b

	return next.Handle(cx)
}

var mqttPrefix = []byte{0, 4, 77, 81, 84, 84} // TODO: make this more descriptive

var mqttRevisions = map[byte]string{
	// 1: "",
	// 2: "",
	3: "3.1", // and 3.0, according to https://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html?
	4: "3.1.1",
	5: "5",
}

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

func checkConnection(b bufferedConn) (bool, error) {

	var err error

	// The logic below is largely based on the implementation found in github.com/VolantMQ/volantmq
	// Also see: https://github.com/VolantMQ/volantmq/blob/3e00ff16a9a086c5b40c45edd77fec98ab530bc5/connection/reader.go#L114

	var header []byte
	peekCount := 2
	for {
		if peekCount > 5 {
			// break out if we try to peek too far for what's allowed in MQTT
			return false, fmt.Errorf("error in protocol")
		}

		if header, err = b.Peek(peekCount); err != nil {
			// break out on error when peeking
			return false, err
		}

		// when sufficient bytes have been read, check if we should continue
		// peeking or have found the correct length of the message
		if header[peekCount-1] >= 0x80 {
			peekCount++
		} else {
			break
		}
	}

	fmt.Println(header)

	// Get the remaining (and total) length of the message
	remainingLength, numberOfLengthBytes := binary.Uvarint(header[1:])
	totalLength := 1 + numberOfLengthBytes + int(remainingLength)

	p, err := b.Peek(totalLength)
	if err != nil {
		return false, err
	}

	fmt.Println(p)

	// we're expecting a CONNECT
	if p[0] != 0x10 { // 00010000; also see: https://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Toc398718029
		return false, fmt.Errorf("unexpected control packet type: 0x%s", hex.EncodeToString(p[0:1]))
	}

	if !bytes.Equal(p[2:8], mqttPrefix) {
		return false, fmt.Errorf("got wrong protocol name: 0x%s", hex.EncodeToString(p[2:8]))
	}

	revision := p[8]
	if _, ok := mqttRevisions[revision]; !ok {
		return false, fmt.Errorf("got invalid revision: 0x%s", hex.EncodeToString(p[8:9]))
	}

	// TODO: add some configuration and perform checks based on those?
	// TODO: do some stuff with the rest of the bytes?

	return true, err
}

// Interface guards
var (
	_ caddy.Module       = (*MQTT)(nil)
	_ caddy.Provisioner  = (*MQTT)(nil)
	_ layer4.NextHandler = (*MQTT)(nil)
)
