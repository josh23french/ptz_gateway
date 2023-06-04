# ptz_gateway

ptz_gateway is a router/translator for various Pan-Tilt-Zoom protocols. It receives commands from a controller and sends them to a specified device using a different medium or protocol.

## Supported Protocols

* VISCA (RS-422/485/232 full-duplex)
* VISCA (TCP)
* VISCA (UDP)

Other protocols can be easily added in the future, including translations from/to VISCA from Pelco-D, for example (adding such a translation layer will be a large undertaking due to needing to have both codecs fully abstracted).

## Installation

```shell
go install github.com/josh23french/ptz_gateway
```

## Running

By default, the software looks for a configuration file of `ptz_gateway.yaml` (or `.yml`) in the current directory, and then in `/etc/`. If no configuration file is found, the software will refuse to start. A configuration file can also be specified on the command line:

```shell
ptz_gateway -config.file /data/studio-c.yml
```

## Configuration

What follows is an example configuration for an RS422 controller attached to ttyS0 and two VISCA-over-IP cameras using TCP and different ports.

When a message comes in from any of these connections, it is forwarded to the destination embedded in the packet's header, and translated to that device's protocol if necessary.

```yaml
routes:
  0:
    protocol: VISCA
    # subproto defaults to serial if provided a device
    device: /dev/ttyS0
  1:
    protocol: VISCA
    connproto: tcp
    ip: 10.1.2.7
    port: 5678
  2:
    protocol: VISCA
    # subproto defaults to tcp if provided an IP
    ip: 10.1.2.8
    port: 1259
```

## Automatic Address Translation

Take for example an IP camera with address 2. IP-based PTZ protocols generally force the camera address to be 1, even if the camera has a different address on the serial bus. So messages require translation.

When the controller refers to camera 2, the gateway must open a connection to the camera over IP and manage sending packets destined to camera 2 over that connection.

Messages to this camera (with `destination=2`) are sent with `destination=1`, and return messages from the camera (with `source=1`) are translated to show they came from `source=2`.

## Design

`Interface`s make connection to cameras or serial buses. These form something similar to the Link Layer. They are responsible for receiving and sending packets on different interfaces: tcp socket, udp socket, serial device, etc. They pass the received packets to the `Router`'s input queue and send any messages on their output queue put there by the `Router`.

The `Router` processes packets in its input queue, determines where the packet is destined, and passes it to the responsible `Interface`'s output queue.

## License

[Affero General Public License 3.0](LICENSE)

## Support Disclaimer

Paid support can be provided upon request.

No other support will be provided.
