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

package dag

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/stretchr/testify/assert"

	"github.com/nuts-foundation/go-stoabs"
	crypto2 "github.com/nuts-foundation/nuts-node/crypto"
	"github.com/nuts-foundation/nuts-node/crypto/hash"
	"github.com/nuts-foundation/nuts-node/vdr/doc"
	"github.com/nuts-foundation/nuts-node/vdr/types"
)

func Test_PrevTransactionVerifier(t *testing.T) {
	ctx := context.Background()
	rootPayload := []byte{0}
	root, _, _ := CreateTestTransactionEx(1, hash.SHA256Sum(rootPayload), nil)

	t.Run("ok - prev is present", func(t *testing.T) {
		testState := createState(t).(*state)
		payload := []byte{0}
		tx, _, _ := CreateTestTransactionEx(1, hash.SHA256Sum(payload), nil, root)
		_ = testState.Add(ctx, root, rootPayload)

		_ = testState.db.Read(ctx, func(dbTx stoabs.ReadTx) error {
			err := NewPrevTransactionsVerifier()(dbTx, tx)
			assert.NoError(t, err)
			return nil
		})
	})

	t.Run("failed - prev not present", func(t *testing.T) {
		testState := createState(t).(*state)
		tx, _, _ := CreateTestTransaction(1, root)

		_ = testState.db.Read(ctx, func(dbTx stoabs.ReadTx) error {
			err := NewPrevTransactionsVerifier()(dbTx, tx)
			assert.Contains(t, err.Error(), "transaction is referring to non-existing previous transaction")
			return nil
		})
	})

	t.Run("error - incorrect lamport clock", func(t *testing.T) {
		testState := createState(t).(*state)
		_ = testState.Add(ctx, root, rootPayload)

		// malformed TX with LC = 2
		unsignedTransaction, _ := NewTransaction(hash.EmptyHash(), "application/did+json", []hash.SHA256Hash{root.Ref()}, nil, 2)
		signer := crypto2.NewTestKey("1")
		signedTransaction, _ := NewTransactionSigner(signer, true).Sign(unsignedTransaction, time.Now())

		_ = testState.db.Read(ctx, func(dbTx stoabs.ReadTx) error {
			err := NewPrevTransactionsVerifier()(dbTx, signedTransaction)
			assert.EqualError(t, err, "transaction has an invalid lamport clock value")
			return nil
		})
	})
}

func TestTransactionSignatureVerifier(t *testing.T) {
	t.Run("embedded JWK, sign -> verify", func(t *testing.T) {
		err := NewTransactionSignatureVerifier(nil)(nil, CreateTestTransactionWithJWK(1))
		assert.NoError(t, err)
	})
	t.Run("embedded JWK, sign -> marshal -> unmarshal -> verify", func(t *testing.T) {
		expected, _ := ParseTransaction(CreateTestTransactionWithJWK(1).Data())
		err := NewTransactionSignatureVerifier(nil)(nil, expected)
		assert.NoError(t, err)
	})
	t.Run("referral with key ID", func(t *testing.T) {
		transaction, _, publicKey := CreateTestTransaction(1)
		expected, _ := ParseTransaction(transaction.Data())
		err := NewTransactionSignatureVerifier(&doc.StaticKeyResolver{Key: publicKey})(nil, expected)
		assert.NoError(t, err)
	})
	t.Run("wrong key", func(t *testing.T) {
		attackerKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		transaction, _, _ := CreateTestTransaction(1)
		expected, _ := ParseTransaction(transaction.Data())
		err := NewTransactionSignatureVerifier(&doc.StaticKeyResolver{Key: attackerKey.Public()})(nil, expected)
		assert.EqualError(t, err, "failed to verify message: failed to verify signature using ecdsa")
	})
	t.Run("key type is incorrect", func(t *testing.T) {
		d, _, _ := CreateTestTransaction(1)
		tx := d.(*transaction)
		tx.signingKey = jwk.NewSymmetricKey()
		err := NewTransactionSignatureVerifier(nil)(nil, tx)
		assert.EqualError(t, err, "failed to verify message: failed to retrieve ecdsa.PublicKey out of []uint8: expected ecdsa.PublicKey or *ecdsa.PublicKey, got []uint8")
	})
	t.Run("unable to derive key from JWK", func(t *testing.T) {
		d, _, _ := CreateTestTransaction(1)
		transaction := d.(*transaction)
		transaction.signingKey = jwk.NewOKPPublicKey()
		err := NewTransactionSignatureVerifier(nil)(nil, transaction)
		assert.EqualError(t, err, "failed to build public key: invalid curve algorithm P-invalid")
	})
	t.Run("unable to resolve key by hash", func(t *testing.T) {
		d := CreateSignedTestTransaction(1, time.Now(), nil, "foo/bar", false)
		ctrl := gomock.NewController(t)
		keyResolver := types.NewMockKeyResolver(ctrl)
		keyResolver.EXPECT().ResolvePublicKey(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed"))

		err := NewTransactionSignatureVerifier(keyResolver)(nil, d)
		if !assert.Error(t, err) {
			return
		}

		assert.Contains(t, err.Error(), "unable to verify transaction signature, can't resolve key by TX ref")
		assert.Contains(t, err.Error(), "failed")
	})
}
