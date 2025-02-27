/*
 * Copyright (C) 2022 Nuts community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package transport

import (
	"github.com/nuts-foundation/go-did/did"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseAddress(t *testing.T) {
	t.Run("valid - with port", func(t *testing.T) {
		addr, err := ParseAddress("grpc://foo:5050")
		assert.NoError(t, err)
		assert.Equal(t, "foo:5050", addr)
	})
	t.Run("valid - without port", func(t *testing.T) {
		addr, err := ParseAddress("grpc://foo")
		assert.NoError(t, err)
		assert.Equal(t, "foo", addr)
	})
	t.Run("invalid - no scheme", func(t *testing.T) {
		addr, err := ParseAddress("foo")
		assert.Empty(t, addr)
		assert.EqualError(t, err, "invalid URL scheme")
	})
	t.Run("invalid - invalid scheme", func(t *testing.T) {
		addr, err := ParseAddress("http://foo")
		assert.Empty(t, addr)
		assert.EqualError(t, err, "invalid URL scheme")
	})
}

func TestPeer_ToFields(t *testing.T) {
	peer := Peer{
		ID:      "abc",
		Address: "def",
		NodeDID: did.MustParseDID("did:abc:123"),
	}

	assert.Len(t, peer.ToFields(), 3)
	assert.Equal(t, "abc", peer.ToFields()["peerID"])
	assert.Equal(t, "def", peer.ToFields()["peerAddr"])
	assert.Equal(t, "did:abc:123", peer.ToFields()["peerDID"])
}
