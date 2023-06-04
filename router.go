//  router.go - the core ptz_gateway Router
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

	"github.com/josh23french/visca"
)

// Route is anything that supports sending messages
type Route interface {
	Send(*visca.Packet) error
}

// Router processes incoming messages and sends them out another "interface"
type Router struct {
	Interfaces map[int]Iface
	Routes     map[int]Route
	recvQ      chan *visca.Packet
}

// NullRoute sends messages to the void
type NullRoute struct{}

// Send just returns no error and does nothing with your message
func (fr *NullRoute) Send(packet *visca.Packet) error {
	return nil
}

// NewRouter initializes a new Router
func NewRouter() *Router {
	return &Router{
		Routes: make(map[int]Route),
	}
}

// UpsertRoute adds or replaces a route for an address
func (r *Router) UpsertRoute(addr int, route Route) {
	r.Routes[addr] = route
}

// RemoveRoute deletes a route for an address
func (r *Router) RemoveRoute(addr int) {
	delete(r.Routes, addr)
}

// Route receives a message from one Interface and sends it out another
func (r *Router) Route(packet *visca.Packet) error {
	// Look up route
	route, err := r.Lookup(packet.Destination())
	if err != nil {
		return err
	}

	// Pass the whole message out that route
	err = route.Send(packet)
	if err != nil {
		return err
	}

	return nil
}

// Lookup returns the route for an address
func (r *Router) Lookup(addr int) (Route, error) {
	if route, ok := r.Routes[addr]; ok {
		return route, nil
	}
	return nil, fmt.Errorf("no route to host %v", addr)
}
