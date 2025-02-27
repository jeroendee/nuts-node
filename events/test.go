/*
 * Nuts node
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

package events

import (
	"testing"

	"github.com/nuts-foundation/nuts-node/core"
	"github.com/nuts-foundation/nuts-node/test/io"
)

func NewTestManager(t *testing.T) Event {
	config := DefaultConfig()
	testDir := io.TestDirectory(t)

	eventManager := &manager{
		config:  config,
		streams: map[string]Stream{},
	}
	cfg := *core.NewServerConfig()
	cfg.Datadir = testDir
	if err := eventManager.Configure(cfg); err != nil {
		t.Fatal(err)
	}
	if err := eventManager.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		eventManager.Shutdown()
	})
	return eventManager
}
