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
	"testing"

	"github.com/josh23french/visca"
	"github.com/stretchr/testify/assert"
)

func TestInsertLookupDeleteRoute(t *testing.T) {
	routeInserted := &NullRoute{}

	addr := 1
	router := NewRouter()
	router.UpsertRoute(addr, routeInserted)
	routeLookedUp, err := router.Lookup(addr)

	assert.Nil(t, err, "should have no error looking up route just inserted")
	assert.Equal(t, routeInserted, routeLookedUp, "route inserted should equal the one looked up from the same address")

	router.RemoveRoute(addr)
	routeLookedUp, err = router.Lookup(addr)
	assert.NotNil(t, err, "should have error looking up route we just deleted")
}

func TestNullRoute(t *testing.T) {
	route := &NullRoute{}
	pkt, err := visca.NewPacket(0, 1, []byte("invalid"))
	assert.Nil(t, err, "should have no error creating invalid packet")
	err = route.Send(pkt)
	assert.Nil(t, err, "send should never have an error, even if the message is invalid")
}
