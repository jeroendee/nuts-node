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

package holder

import (
	ssi "github.com/nuts-foundation/go-did"
	"github.com/nuts-foundation/go-did/did"
	"github.com/nuts-foundation/go-did/vc"
	"github.com/nuts-foundation/nuts-node/vcr/signature/proof"
)

// VerifiableCredentialLDContextV1 holds the URI of the JSON-LD context for Verifiable Credentials.
var VerifiableCredentialLDContextV1 = ssi.MustParseURI("https://www.w3.org/2018/credentials/v1")

// VerifiablePresentationLDType holds the JSON-LD type for Verifiable Presentations.
var VerifiablePresentationLDType = ssi.MustParseURI("VerifiablePresentation")

// Holder holds logic for Presenting credentials
// A Holder is a role an entity might perform by:
// * Possessing one or more verifiable credentials
// * Generating verifiable presentations from them
// Example holders include students, employees, and customers.
type Holder interface {
	// BuildVP builds and signs a Verifiable Presentation using the given Verifiable Credentials.
	// The assertion key used for signing it is taken from signerDID's DID document.
	// If signerDID is not provided, it will be derived from the credentials credentialSubject.id fields. But only if all provided credentials have the same (singular) credentialSubject.id field.
	BuildVP(credentials []vc.VerifiableCredential, proofOptions proof.ProofOptions, signerDID *did.DID, validateVC bool) (*vc.VerifiablePresentation, error)
}
