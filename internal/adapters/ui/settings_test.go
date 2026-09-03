// Copyright 2025.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ui

import (
	"path/filepath"
	"testing"
)

func TestSettingsManagerMuxModeDefaultsToTmuxCC(t *testing.T) {
	t.Parallel()

	manager := &settingsManager{
		filePath: filepath.Join(t.TempDir(), "settings.json"),
	}

	mode, err := manager.LoadMuxMode("prod")
	if err != nil {
		t.Fatalf("LoadMuxMode() error = %v", err)
	}
	if mode != muxModeTmuxCC {
		t.Fatalf("LoadMuxMode() = %q, want %q", mode, muxModeTmuxCC)
	}
}

func TestSettingsManagerSaveMuxModePersistsByAlias(t *testing.T) {
	t.Parallel()

	manager := &settingsManager{
		filePath: filepath.Join(t.TempDir(), "settings.json"),
	}

	if err := manager.SaveSortMode(SortByLastSeenDesc); err != nil {
		t.Fatalf("SaveSortMode() error = %v", err)
	}
	if err := manager.SaveMuxMode("prod", muxModeOff); err != nil {
		t.Fatalf("SaveMuxMode(prod) error = %v", err)
	}
	if err := manager.SaveMuxMode("staging", muxModeTmuxCC); err != nil {
		t.Fatalf("SaveMuxMode(staging) error = %v", err)
	}

	prodMode, err := manager.LoadMuxMode("prod")
	if err != nil {
		t.Fatalf("LoadMuxMode(prod) error = %v", err)
	}
	if prodMode != muxModeOff {
		t.Fatalf("LoadMuxMode(prod) = %q, want %q", prodMode, muxModeOff)
	}

	stagingMode, err := manager.LoadMuxMode("staging")
	if err != nil {
		t.Fatalf("LoadMuxMode(staging) error = %v", err)
	}
	if stagingMode != muxModeTmuxCC {
		t.Fatalf("LoadMuxMode(staging) = %q, want %q", stagingMode, muxModeTmuxCC)
	}

	mode, err := manager.LoadSortMode()
	if err != nil {
		t.Fatalf("LoadSortMode() error = %v", err)
	}
	if mode != SortByLastSeenDesc {
		t.Fatalf("LoadSortMode() = %v, want %v", mode, SortByLastSeenDesc)
	}
}
