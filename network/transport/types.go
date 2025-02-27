/*
 * Copyright (C) 2021 Nuts community
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
	"errors"
	"fmt"
	"github.com/nuts-foundation/nuts-node/core"
	"github.com/sirupsen/logrus"
	"net/url"
	"time"

	"github.com/nuts-foundation/go-did/did"
)

// PeerID defines a peer's unique identifier.
type PeerID string

// String returns the PeerID as string.
func (p PeerID) String() string {
	return string(p)
}

// ParseAddress parses the given input string to a gRPC target address. The input must include the protocol scheme (e.g. grpc://).
func ParseAddress(input string) (string, error) {
	parsed, err := url.Parse(input)
	if err != nil {
		return "", err
	}
	if parsed.Scheme != "grpc" {
		return "", errors.New("invalid URL scheme")
	}
	return parsed.Host, nil
}

// Peer holds the properties of a remote node we're connected to
type Peer struct {
	// ID holds the unique identificator of the peer
	ID PeerID
	// Address holds the remote address of the node we're actually connected to
	Address string
	// NodeDID holds the DID that the peer uses to identify its node on the network.
	// It is only set when properly authenticated.
	NodeDID did.DID
	// AcceptUnauthenticated indicates if a connection may be made with this Peer even if the NodeDID could not be authenticated.
	AcceptUnauthenticated bool
}

// ToFields returns the peer as a map of fields, to be used when logging the peer details.
func (p Peer) ToFields() logrus.Fields {
	return map[string]interface{}{
		core.LogFieldPeerID:      p.ID.String(),
		core.LogFieldPeerAddr:    p.Address,
		core.LogFieldPeerNodeDID: p.NodeDID.String(),
	}
}

// String returns the peer as string.
func (p Peer) String() string {
	if p.NodeDID.Empty() {
		return fmt.Sprintf("%s@%s", p.ID, p.Address)
	}
	return fmt.Sprintf("%s(%s)@%s", p.ID, p.NodeDID.String(), p.Address)
}

// Diagnostics contains information that is shared to this node's peers on request.
type Diagnostics struct {
	// Uptime the uptime (time since the node started) in seconds.
	Uptime time.Duration `json:"uptime"`
	// Peers contains the peer IDs of the node's peers.
	Peers []PeerID `json:"peers"`
	// NumberOfTransactions contains the total number of transactions on the node's DAG.
	NumberOfTransactions uint32 `json:"transactionNum"`
	// SoftwareVersion contains an indication of the software version of the node. It's recommended to use a (Git) commit ID that uniquely resolves to a code revision, alternatively a semantic version could be used (e.g. 1.2.5).
	SoftwareVersion string `json:"softwareVersion"`
	// SoftwareID contains an indication of the vendor of the software of the node. For open source implementations it's recommended to specify URL to the public, open source repository.
	// Proprietary implementations could specify the product's or vendor's name.
	SoftwareID string `json:"softwareID"`
}

// ConnectorStats holds statistics of an outbound connector.
type ConnectorStats struct {
	// Address holds the target address the connector is connecting to.
	Address string
	// Attempts holds the number of times the node tried to connect to the peer.
	Attempts uint32
	// LastAttempt holds the time of the last connection attempt.
	LastAttempt time.Time
}

// NutsCommServiceType holds the DID document service type that specifies the Nuts network service address of the Nuts node.
const NutsCommServiceType = "NutsComm"
