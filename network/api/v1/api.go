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

package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nuts-foundation/nuts-node/core"
	hash2 "github.com/nuts-foundation/nuts-node/crypto/hash"
	"github.com/nuts-foundation/nuts-node/network"
	"github.com/nuts-foundation/nuts-node/network/log"
)

// Wrapper implements the ServerInterface for the network API.
type Wrapper struct {
	Service network.Transactions
}

func (a *Wrapper) Routes(router core.EchoRouter) {
	RegisterHandlers(router, a)
}

// ListTransactions lists all transactions
func (a Wrapper) ListTransactions(ctx echo.Context) error {
	transactions, err := a.Service.ListTransactions()
	if err != nil {
		log.Logger().Errorf("Error while listing transactions: %v", err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	results := make([]string, len(transactions))
	for i, transaction := range transactions {
		results[i] = string(transaction.Data())
	}
	return ctx.JSON(http.StatusOK, results)
}

// GetTransaction returns a specific transaction
func (a Wrapper) GetTransaction(ctx echo.Context, hashAsString string) error {
	hash, err := hash2.ParseHex(hashAsString)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	transaction, err := a.Service.GetTransaction(hash)
	if err != nil {
		log.Logger().Errorf("Error while retrieving transaction (hash=%s): %v", hash, err)
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	if transaction == nil {
		return ctx.String(http.StatusNotFound, "transaction not found")
	}
	ctx.Response().Header().Set(echo.HeaderContentType, "application/jose")
	ctx.Response().WriteHeader(http.StatusOK)
	_, err = ctx.Response().Writer.Write(transaction.Data())
	return err
}

// GetTransactionPayload returns the payload of a specific transaction
func (a Wrapper) GetTransactionPayload(ctx echo.Context, hashAsString string) error {
	hash, err := hash2.ParseHex(hashAsString)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	data, err := a.Service.GetTransactionPayload(hash)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	if data == nil {
		return ctx.String(http.StatusNotFound, "transaction or contents not found")
	}
	ctx.Response().Header().Set(echo.HeaderContentType, "application/octet-stream")
	ctx.Response().WriteHeader(http.StatusOK)
	_, err = ctx.Response().Writer.Write(data)
	return err
}
