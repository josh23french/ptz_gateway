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
	"os"

	"github.com/josh23french/visca"
	"github.com/rs/zerolog/log"
)

// FileIface implements the Iface interface for serial VISCA connections
type FileIface struct {
	device       string
	file         *os.File
	scanner      *visca.Scanner
	receiveQueue chan *visca.Packet
	quit         chan struct{}
}

// NewFileIface creates a new FileIface
func NewFileIface(device string) (*FileIface, error) {
	return &FileIface{
		device:       device,
		file:         nil,
		scanner:      nil,
		receiveQueue: nil,
		quit:         make(chan struct{}),
	}, nil
}

// Start the interface
func (i *FileIface) Start() error {
	log.Info().Msgf("Opening serial interface %v...", i.device)

	file, err := os.OpenFile(i.device, os.O_RDWR, os.ModeCharDevice)
	if err != nil {
		return err
	}
	log.Info().Msgf("Opened serial interface %v", i.device)
	i.file = file

	i.quit = make(chan struct{})

	i.scanner = visca.NewScanner(i.file)

	go i.scanner.Scan(i.receiveQueue, i.quit)
	log.Info().Msgf("Started read loop from serial interface %v", i.device)

	return nil
}

// Stop the interface
func (i *FileIface) Stop() {
	if i.file == nil {
		log.Warn().Msg("Never Started")
		return
	}

	// Stop the receive goroutine first
	close(i.quit)

	// Then close the port
	file := i.file
	err := file.Close()
	if err != nil {
		log.Warn().Err(err).Msgf("Error stopping serial interface %v", i.device)
	}
}

// Send a packet
func (i *FileIface) Send(pkt *visca.Packet) error {
	if i.file == nil {
		log.Warn().Msg("not started")
		return ErrIfaceNotStarted
	}
	file := *i.file
	log.Debug().Msgf("Sending packet %v via %v", pkt, i.device)
	written, err := file.Write(pkt.Bytes())
	if err != nil {
		return err
	}

	if written != len(pkt.Bytes()) {
		return ErrIncompletePacketSent
	}

	return nil
}

// SetReceiveQueue for received packets
func (i *FileIface) SetReceiveQueue(q chan *visca.Packet) {
	i.receiveQueue = q
	return
}
