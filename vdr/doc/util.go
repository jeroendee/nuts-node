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

package doc

import (
	"fmt"
	"github.com/nuts-foundation/go-did"
	"github.com/nuts-foundation/go-did/did"
	"github.com/nuts-foundation/nuts-node/vdr/types"
	"strings"
)

const serviceTypeQueryParameter = "type"
const serviceEndpointPath = "serviceEndpoint"

// MakeServiceReference creates a service reference, which can be used as query when looking up services.
func MakeServiceReference(subjectDID did.DID, serviceType string) ssi.URI {
	ref := subjectDID.URI()
	ref.Opaque += "/" + serviceEndpointPath
	ref.Fragment = ""
	ref.RawQuery = fmt.Sprintf("%s=%s", serviceTypeQueryParameter, serviceType)
	return ref
}

// IsServiceReference checks whether the given endpoint string looks like a service reference (e.g. did:nuts:1234/serviceType?type=HelloWorld).
func IsServiceReference(endpoint string) bool {
	return strings.HasPrefix(endpoint, "did:")
}

// ValidateServiceReference checks whether the given URI matches the format for a service reference.
func ValidateServiceReference(endpointURI ssi.URI) error {
	// Parse it as DID URL since DID URLs are rootless and thus opaque (RFC 3986), meaning the path will be part of the URI body, rather than the URI path.
	// For DID URLs the path is parsed properly.
	didEndpointURL, err := did.ParseDIDURL(endpointURI.String())
	if err != nil {
		return types.ErrInvalidServiceQuery{Cause: err}
	}
	if didEndpointURL.Path != serviceEndpointPath {
		// Service reference doesn't refer to `/serviceEndpoint`
		return types.ErrInvalidServiceQuery{Cause: fmt.Errorf("URL path must be '/%s'", serviceEndpointPath)}
	}
	queriedServiceType := endpointURI.Query().Get(serviceTypeQueryParameter)
	typeQueryParameterError := types.ErrInvalidServiceQuery{Cause: fmt.Errorf("URL must contain exactly one '%s' query parameter, with exactly one value", serviceTypeQueryParameter)}
	if len(queriedServiceType) == 0 {
		// Service reference doesn't contain `type` query parameter
		return typeQueryParameterError
	}
	if len(endpointURI.Query()[serviceTypeQueryParameter]) > 1 {
		// Service reference contains more than 1 `type` query parameter
		return typeQueryParameterError
	}
	if len(endpointURI.Query()) > 1 {
		// Service reference contains more than just `type` query parameter
		return typeQueryParameterError
	}
	return nil
}
