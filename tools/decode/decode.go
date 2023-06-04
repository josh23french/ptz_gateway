package main

import (
	"bytes"
	"flag"
	"os"

	"github.com/josh23french/ptz_gateway/pkg/gateway"
	"github.com/josh23french/visca"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	var (
		device = flag.String("device", "/dev/cu.usbserial-141440", "device to open")
	)

	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	pc := make(chan *visca.Packet)

	// CONTROLLER
	iface, err := gateway.NewFileIface(*device)
	if err != nil {
		panic(err)
	}
	iface.SetReceiveQueue(pc)
	err = iface.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Couldn't start SerialIface: %v", err.Error())
		panic(err)
	}
	log.Info().Msg("Started Serial Iface...")

	// CAMERA 1
	iface2, err := gateway.NewTCPIface("10.1.2.7:1259")
	if err != nil {
		panic(err)
	}
	iface2.SetReceiveQueue(pc)
	err = iface2.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Couldn't start TCPIface: %v", err.Error())
		panic(err)
	}
	log.Info().Msg("Started TCP Iface for camera 1...")

	// CAMERA 2
	iface3, err := gateway.NewTCPIface("10.1.2.8:5678")
	if err != nil {
		panic(err)
	}
	iface3.SetReceiveQueue(pc)
	err = iface3.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Couldn't start TCPIface for camera 2: %v", err.Error())
		panic(err)
	}
	log.Info().Msg("Started TCP Iface for camera 2...")

	// Synthesize a Network Change packet
	pkt, err := visca.NewPacket(1, 0, []byte{0x38})
	if err != nil {
		log.Fatal().Err(err).Msg("error synthesizing network change packet")
	}
	// Send to controller
	iface.Send(pkt)

	for {
		var err error
		packet := <-pc
		if packet == nil {
			continue
		}
		log.Info().Msgf("Got %v %v packet: From %v To %v Message: %v", packet.Message.Category().String(), packet.Message.Type().String(), packet.Source(), packet.Destination(), packet.Message)
		switch packet.Destination() {
		case 0:
			err = iface.Send(packet)
			if err != nil {
				log.Error().Err(err).Msg("error sending packet")
			}
		case 1:
			err = iface2.Send(packet)
			if err != nil {
				log.Error().Err(err).Msg("error sending packet")
			}
		case 2:
			err = iface3.Send(packet)
			if err != nil {
				log.Error().Err(err).Msg("error sending packet")
			}
		case 8:
			// Broadcast messages

			// 0x30, 0x01 - Network Change / Address Set
			if bytes.Equal(packet.Message, []byte{48, 1}) {
				// Calculate number of cameras we have
				camCount := 2
				pkt, err := visca.NewPacket(0, 8, []byte{0x30, byte(camCount + 1)})
				if err != nil {
					log.Error().Err(err).Msg("error synthesizing response packet")
					continue
				}
				log.Info().Msgf("Synthesized network change packet response: %v", pkt)
				iface.Send(pkt)
			} else {
				log.Warn().Msgf("Unknown broadcast message: %v", packet.Message)
			}
		default:
			log.Warn().Msgf("Unknown destination: %v", packet.Destination())
		}
	}
}
