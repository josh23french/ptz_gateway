//  tcp.go - TCP interface
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

package gateway

import (
	"net"
	"time"

	"github.com/josh23french/visca"
	"github.com/rs/zerolog/log"
)

// TCPIface implements the Iface interface without sending packets anywhere
type TCPIface struct {
	hostPort     string
	conn         net.Conn
	scanner      *visca.Scanner
	receiveQueue chan *visca.Packet
	quit         chan struct{}
}

// NewTCPIface creates a new TCPIface
func NewTCPIface(hostPort string) (*TCPIface, error) {
	return &TCPIface{
		hostPort:     hostPort,
		conn:         nil,
		scanner:      nil,
		receiveQueue: nil,
		quit:         make(chan struct{}),
	}, nil
}

// Start the interface
func (i *TCPIface) Start() error {
	conn, err := net.DialTimeout("tcp", i.hostPort, time.Second)
	if err != nil {
		return err
	}

	i.conn = conn

	i.scanner = visca.NewScanner(i.conn)

	go i.scanner.Scan(i.receiveQueue, i.quit)
	log.Info().Msgf("Started read loop from tcp interface %v", i.hostPort)

	return nil
}

// Stop the interface
func (i *TCPIface) Stop() {
	if i.conn == nil {
		log.Warn().Msg("Never Started")
		return
	}

	// Stop the receive goroutine first
	close(i.quit)
	err := i.conn.Close()
	if err != nil {
		log.Warn().Err(err).Msgf("Error stopping tcp interface %v", i.hostPort)
	}
	return
}

// Send a packet
func (i *TCPIface) Send(pkt *visca.Packet) error {
	if i.conn == nil {
		log.Warn().Msg("not started")
		return ErrIfaceNotStarted
	}
	log.Debug().Msgf("Sending packet %v to %v", pkt, i.hostPort)
	written, err := i.conn.Write(pkt.Bytes())
	if err != nil {
		return err
	}

	if written != len(pkt.Bytes()) {
		return ErrIncompletePacketSent
	}

	return nil
}

// SetReceiveQueue for received packets
func (i *TCPIface) SetReceiveQueue(q chan *visca.Packet) {
	i.receiveQueue = q
	return
}
