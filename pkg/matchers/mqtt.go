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
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/hslatman/poc-caddy-mqtt-proxy/internal/conn"
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

	return nil
}

// Match returns true if the conn starts with MQTT ......
func (m MatchMQTT) Match(cx *layer4.Connection) (bool, error) {

	// TODO: this has not been tested to work fully (yet)!

	wrappedConnection := conn.NewBufferedConn(cx.Conn)

	return m.CheckBufferedConnection(wrappedConnection)
}

func (m *MatchMQTT) CheckBufferedConnection(b conn.BufferedConn) (bool, error) {

	var err error

	// The header/length logic below is largely based on the implementation found in github.com/VolantMQ/volantmq
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

	// Get the remaining (and total) length of the message
	remainingLength, numberOfLengthBytes := binary.Uvarint(header[1:])
	totalLength := 1 + numberOfLengthBytes + int(remainingLength)

	p, err := b.Peek(totalLength)
	if err != nil {
		return false, err
	}

	// we're expecting an MQTT CONNECT packet
	if p[0] != 0x10 { // 00010000; also see: https://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Toc398718029
		return false, fmt.Errorf("unexpected control packet type: 0x%s", hex.EncodeToString(p[0:1]))
	}

	if !bytes.Equal(p[2:8], mqttProtocolName) {
		return false, fmt.Errorf("got wrong protocol name: 0x%s", hex.EncodeToString(p[2:8]))
	}

	revision := p[8]
	if _, ok := mqttRevisions[revision]; !ok {
		return false, fmt.Errorf("got invalid revision: 0x%s", hex.EncodeToString(p[8:9]))
	}

	// TODO: add some configuration and perform checks based on those?
	// TODO: do something with the rest of the bytes?

	return true, err
}

var mqttProtocolName = append([]byte{0x00, 0x04}, []byte("MQTT")...) // https://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Toc398718029

var mqttRevisions = map[byte]string{
	// 1: "",
	// 2: "",
	3: "3.1", // and 3.0, according to https://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html?
	4: "3.1.1",
	5: "5",
}

// Interface guards
var (
	_ caddy.Module       = (*MatchMQTT)(nil)
	_ caddy.Provisioner  = (*MatchMQTT)(nil)
	_ layer4.ConnMatcher = (*MatchMQTT)(nil)
	_ json.Marshaler     = (*MatchMQTT)(nil)
	_ json.Unmarshaler   = (*MatchMQTT)(nil)
)
