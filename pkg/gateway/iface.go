//  iface.go - interface code
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
	"errors"

	"github.com/josh23french/visca"
)

// Iface represents an interface
type Iface interface {
	Start() error                       // Start the Iface; returns error if failed be started
	Stop()                              // Stop the Iface, can have no error
	Send(*visca.Packet) error           // Send a packet out the Iface, may be queued
	SetReceiveQueue(chan *visca.Packet) // Tell the Iface where to send received packets
}

// Error constants
var (
	ErrIfaceNotStarted      = errors.New("iface not started")
	ErrIncompletePacketSent = errors.New("incomplete packet sent")
)

// NullIface implements the Iface interface without sending packets anywhere
type NullIface struct {
}

// Start the interface
func (i *NullIface) Start() error {
	return nil
}

// Stop the interface
func (i *NullIface) Stop() {
	return
}

// Send a packet
func (i *NullIface) Send(pkt *visca.Packet) error {
	return nil
}

// SetReceiveQueue for received packets
func (i *NullIface) SetReceiveQueue(q chan *visca.Packet) {
	return
}
