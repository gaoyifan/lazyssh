package ui

import (
	"path/filepath"
	"testing"
)

func TestSettingsManagerTmuxCCDefaultsToEnabled(t *testing.T) {
	t.Parallel()

	manager := &settingsManager{
		filePath: filepath.Join(t.TempDir(), "settings.json"),
	}

	enabled, err := manager.LoadTmuxCCEnabled("prod")
	if err != nil {
		t.Fatalf("LoadTmuxCCEnabled() error = %v", err)
	}
	if !enabled {
		t.Fatalf("LoadTmuxCCEnabled() = false, want true by default")
	}
}

func TestSettingsManagerSaveTmuxCCEnabledPersistsByAlias(t *testing.T) {
	t.Parallel()

	manager := &settingsManager{
		filePath: filepath.Join(t.TempDir(), "settings.json"),
	}

	if err := manager.SaveSortMode(SortByLastSeenDesc); err != nil {
		t.Fatalf("SaveSortMode() error = %v", err)
	}
	if err := manager.SaveTmuxCCEnabled("prod", false); err != nil {
		t.Fatalf("SaveTmuxCCEnabled(prod) error = %v", err)
	}
	if err := manager.SaveTmuxCCEnabled("staging", true); err != nil {
		t.Fatalf("SaveTmuxCCEnabled(staging) error = %v", err)
	}

	prodEnabled, err := manager.LoadTmuxCCEnabled("prod")
	if err != nil {
		t.Fatalf("LoadTmuxCCEnabled(prod) error = %v", err)
	}
	if prodEnabled {
		t.Fatalf("LoadTmuxCCEnabled(prod) = true, want false")
	}

	stagingEnabled, err := manager.LoadTmuxCCEnabled("staging")
	if err != nil {
		t.Fatalf("LoadTmuxCCEnabled(staging) error = %v", err)
	}
	if !stagingEnabled {
		t.Fatalf("LoadTmuxCCEnabled(staging) = false, want true")
	}

	mode, err := manager.LoadSortMode()
	if err != nil {
		t.Fatalf("LoadSortMode() error = %v", err)
	}
	if mode != SortByLastSeenDesc {
		t.Fatalf("LoadSortMode() = %v, want %v", mode, SortByLastSeenDesc)
	}
}
