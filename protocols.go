//  protocols.go - routes/connections per protocol
// 	Copyright (C) 2020  Joshua French
//
// 	This program is free software: you can redistribute it and/or modify
// 	it under the terms of the GNU Affero General Public License as published
// 	by the Free Software Foundation, either version 3 of the License, or
// 	(at your option) any later version.
//
// 	This program is distributed in the hope that it will be useful,
// 	but WITHOUT ANY WARRANTY; without even the implied warranty of
// 	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// 	GNU Affero General Public License for more details.
//
// 	You should have received a copy of the GNU Affero General Public License
// 	along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"io"
	"net"

	"go.bug.st/serial"
)

// VISCARoute represents a connection to a VISCA device
type VISCARoute struct {
	config *RouteConfig
	conn   io.ReadWriteCloser
}

// Defaults
const (
	DefaultPort      = 5678
	DefaultConnProto = "tcp"
)

// NewVISCARoute creates a new route
func NewVISCARoute(config *IfaceConfig) (*VISCARoute, error) {
	var conn io.ReadWriteCloser
	if config.Device != "" {
		// VISCA uses 9600_8N1, which is the default, but we specify here anyway
		mode := &serial.Mode{
			BaudRate: 9600,
			DataBits: 8,
			Parity:   serial.NoParity,
			StopBits: 1,
		}
		port, err := serial.Open(config.Device, mode)
		if err != nil {
			return nil, err
		}
		conn = port
	} else if config.IP != "" {
		port := config.Port
		if port == 0 {
			port = DefaultPort
		}
		proto := config.ConnProto
		if proto == "" {
			proto = DefaultConnProto
		}
		address := fmt.Sprintf("%v:%v", config.IP, port)
		var err error
		conn, err = net.Dial(proto, address)
		if err != nil {
			return nil, err
		}
	}

	return &VISCARoute{
		nil,
		conn,
	}, nil
}

// Send opens a new connection if necessary and sends the message
func (r *VISCARoute) Send(message *AbstractMessage) error {
	return nil
}
