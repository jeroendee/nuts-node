/*
 * Nuts node
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

package vcr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/nuts-foundation/go-did/vc"
	"github.com/nuts-foundation/nuts-node/crypto"
	"github.com/nuts-foundation/nuts-node/crypto/hash"
	"github.com/nuts-foundation/nuts-node/crypto/util"
	"github.com/nuts-foundation/nuts-node/events"
	"github.com/nuts-foundation/nuts-node/jsonld"
	"github.com/nuts-foundation/nuts-node/network"
	"github.com/nuts-foundation/nuts-node/network/dag"
	"github.com/nuts-foundation/nuts-node/test"
	"github.com/nuts-foundation/nuts-node/vcr/credential"
	"github.com/nuts-foundation/nuts-node/vcr/types"
	"github.com/nuts-foundation/nuts-node/vcr/verifier"
	"github.com/stretchr/testify/assert"
)

func TestNewAmbassador(t *testing.T) {
	a := NewAmbassador(nil, nil, nil, nil)

	assert.NotNil(t, a)
}

func TestAmbassador_Configure(t *testing.T) {
	t.Run("calls network.subscribe", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		nMock := network.NewMockTransactions(ctrl)

		a := NewAmbassador(nMock, nil, nil, nil)
		nMock.EXPECT().WithPersistency().Times(2)
		nMock.EXPECT().Subscribe("vcr_vcs", gomock.Any(), gomock.Any())
		nMock.EXPECT().Subscribe("vcr_revocations", gomock.Any(), gomock.Any())

		a.Configure()
	})
}

func TestAmbassador_Start(t *testing.T) {
	t.Run("error on stream subscription", func(t *testing.T) {
		ctx := newMockContext(t)
		mockEvent := events.NewMockEvent(ctx.ctrl)
		ctx.vcr.ambassador.(*ambassador).eventManager = mockEvent
		mockPool := events.NewMockConnectionPool(ctx.ctrl)
		mockConnection := events.NewMockConn(ctx.ctrl)
		mockEvent.EXPECT().Pool().Return(mockPool)
		mockPool.EXPECT().Acquire(gomock.Any()).Return(mockConnection, nil, nil)
		mockConnection.EXPECT().JetStream().Return(nil, errors.New("b00m!"))

		err := ctx.vcr.ambassador.Start()

		assert.EqualError(t, err, "failed to subscribe to REPROCESS event stream: b00m!")
	})

	t.Run("error on nats connection acquire", func(t *testing.T) {
		ctx := newMockContext(t)
		mockEvent := events.NewMockEvent(ctx.ctrl)
		ctx.vcr.ambassador.(*ambassador).eventManager = mockEvent
		mockPool := events.NewMockConnectionPool(ctx.ctrl)
		mockEvent.EXPECT().Pool().Return(mockPool)
		mockPool.EXPECT().Acquire(gomock.Any()).Return(nil, nil, errors.New("b00m!"))

		err := ctx.vcr.ambassador.Start()

		assert.EqualError(t, err, "failed to subscribe to REPROCESS event stream: b00m!")
	})
}

func TestAmbassador_handleReprocessEvent(t *testing.T) {
	ctx := newMockContext(t)
	mockWriter := NewMockWriter(ctx.ctrl)
	ctx.vcr.ambassador.(*ambassador).writer = mockWriter

	// load VC
	vc := vc.VerifiableCredential{}
	vcJSON, _ := os.ReadFile("test/vc.json")
	json.Unmarshal(vcJSON, &vc)

	// load key
	pem, _ := os.ReadFile("test/private.pem")
	signer, _ := util.PemToPrivateKey(pem)
	key := crypto.TestKey{
		PrivateKey: signer,
		Kid:        fmt.Sprintf("%s#1", vc.Issuer.String()),
	}

	// trust otherwise Resolve wont work
	ctx.vcr.Trust(vc.Type[0], vc.Issuer)
	ctx.vcr.Trust(vc.Type[1], vc.Issuer)

	// mocks
	ctx.keyResolver.EXPECT().ResolveSigningKey(gomock.Any(), gomock.Any()).Return(signer.Public(), nil)

	// Publish a VC
	payload, _ := json.Marshal(vc)
	unsignedTransaction, err := dag.NewTransaction(hash.SHA256Sum(payload), "application/vc+json", nil, nil, uint32(0))
	signedTransaction, err := dag.NewTransactionSigner(key, true).Sign(unsignedTransaction, time.Now())
	twp := events.TransactionWithPayload{
		Transaction: signedTransaction,
		Payload:     payload,
	}
	twpBytes, _ := json.Marshal(twp)

	_, js, _ := ctx.vcr.eventManager.Pool().Acquire(context.Background())
	_, err = js.Publish("REPROCESS.application/vc+json", twpBytes)

	if !assert.NoError(t, err) {
		return
	}

	test.WaitFor(t, func() (bool, error) {
		_, err := ctx.vcr.Resolve(*vc.ID, nil)
		return err == nil, nil
	}, time.Second, "timeout while waiting for event to be processed")
}

func TestAmbassador_vcCallback(t *testing.T) {
	payload := []byte(jsonld.TestCredential)
	tx, _ := dag.NewTransaction(hash.EmptyHash(), types.VcDocumentType, nil, nil, 0)
	stx := tx.(dag.Transaction)
	validAt := stx.SigningTime()

	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		wMock := NewMockWriter(ctrl)

		target := vc.VerifiableCredential{}
		a := NewAmbassador(nil, wMock, nil, nil).(*ambassador)
		wMock.EXPECT().StoreCredential(gomock.Any(), &validAt).DoAndReturn(func(f interface{}, g interface{}) error {
			target = f.(vc.VerifiableCredential)
			return nil
		})

		err := a.vcCallback(stx, payload)

		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, "did:nuts:B8PUHs2AUHbFF1xLLK4eZjgErEcMXHxs68FteY7NDtCY#123", target.ID.String())
	})

	t.Run("error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		wMock := NewMockWriter(ctrl)

		a := NewAmbassador(nil, wMock, nil, nil).(*ambassador)
		wMock.EXPECT().StoreCredential(gomock.Any(), &validAt).Return(errors.New("b00m!"))

		err := a.vcCallback(stx, payload)

		assert.Error(t, err)
	})

	t.Run("error - invalid payload", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		wMock := NewMockWriter(ctrl)

		a := NewAmbassador(nil, wMock, nil, nil).(*ambassador)

		err := a.vcCallback(stx, []byte("{"))

		assert.Error(t, err)
	})
}

func TestAmbassador_handleNetworkVCs(t *testing.T) {
	tx, _ := dag.NewTransaction(hash.EmptyHash(), types.VcDocumentType, nil, nil, 0)
	stx := tx.(dag.Transaction)

	t.Run("error - invalid payload", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		wMock := NewMockWriter(ctrl)

		a := NewAmbassador(nil, wMock, nil, nil).(*ambassador)

		value, err := a.handleNetworkVCs(dag.Event{
			Transaction: stx,
			Payload:     []byte("{"),
		})

		assert.False(t, value)
		assert.Error(t, err)
	})
}

func Test_ambassador_jsonLDRevocationCallback(t *testing.T) {
	payload, _ := os.ReadFile("test/ld-revocation.json")
	tx, _ := dag.NewTransaction(hash.EmptyHash(), types.RevocationLDDocumentType, nil, nil, 0)
	stx := tx.(dag.Transaction)

	t.Run("ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		revocation := credential.Revocation{}
		assert.NoError(t, json.Unmarshal(payload, &revocation))

		mockVerifier := verifier.NewMockVerifier(ctrl)
		mockVerifier.EXPECT().RegisterRevocation(revocation)
		a := NewAmbassador(nil, nil, mockVerifier, nil).(*ambassador)

		err := a.jsonLDRevocationCallback(stx, payload)
		assert.NoError(t, err)
	})

	t.Run("error - invalid payload", func(t *testing.T) {
		a := NewAmbassador(nil, nil, nil, nil).(*ambassador)

		err := a.jsonLDRevocationCallback(stx, []byte("b00m"))
		assert.EqualError(t, err, "revocation processing failed: invalid character 'b' looking for beginning of value")
	})

	t.Run("error - storing fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockVerifier := verifier.NewMockVerifier(ctrl)
		mockVerifier.EXPECT().RegisterRevocation(gomock.Any()).Return(errors.New("foo"))
		a := NewAmbassador(nil, nil, mockVerifier, nil).(*ambassador)

		err := a.jsonLDRevocationCallback(stx, payload)
		assert.EqualError(t, err, "foo")
	})
}
