package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/steveyegge/beads/internal/utils"
)

func TestGetReadyDatabaseContext_TownRouting(t *testing.T) {
	t.Setenv("BEADS_DIR", "")
	t.Setenv("BEADS_DB", "")

	tmpDir := t.TempDir()

	mayorDir := filepath.Join(tmpDir, "mayor")
	if err := os.MkdirAll(mayorDir, 0750); err != nil {
		t.Fatalf("failed to create mayor dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(mayorDir, "town.json"), []byte(`{"name":"test-town"}`), 0600); err != nil {
		t.Fatalf("failed to create town.json: %v", err)
	}

	townBeadsDir := filepath.Join(tmpDir, ".beads")
	if err := os.MkdirAll(townBeadsDir, 0750); err != nil {
		t.Fatalf("failed to create town beads dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(townBeadsDir, "metadata.json"), []byte(`{}`), 0600); err != nil {
		t.Fatalf("failed to create town metadata.json: %v", err)
	}
	routesContent := `{"prefix":"hq-","path":"."}
{"prefix":"te-","path":"TerraNomadicCity/mayor/rig"}
`
	if err := os.WriteFile(filepath.Join(townBeadsDir, "routes.jsonl"), []byte(routesContent), 0600); err != nil {
		t.Fatalf("failed to write routes.jsonl: %v", err)
	}

	rigDir := filepath.Join(tmpDir, "TerraNomadicCity", "mayor", "rig")
	rigBeadsDir := filepath.Join(rigDir, ".beads")
	if err := os.MkdirAll(rigBeadsDir, 0750); err != nil {
		t.Fatalf("failed to create rig beads dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(rigBeadsDir, "metadata.json"), []byte(`{}`), 0600); err != nil {
		t.Fatalf("failed to create rig metadata.json: %v", err)
	}

	currentDBPath := filepath.Join(rigBeadsDir, "beads.db")
	if err := os.WriteFile(currentDBPath, []byte{}, 0600); err != nil {
		t.Fatalf("failed to create rig db file: %v", err)
	}
	townDBPath := filepath.Join(townBeadsDir, "beads.db")
	if err := os.WriteFile(townDBPath, []byte{}, 0600); err != nil {
		t.Fatalf("failed to create town db file: %v", err)
	}

	t.Chdir(rigDir)

	ctx := getReadyDatabaseContext(currentDBPath)
	if ctx == nil {
		t.Fatal("expected ready database context, got nil")
	}
	if utils.CanonicalizePath(ctx.CurrentDBPath) != utils.CanonicalizePath(currentDBPath) {
		t.Fatalf("current db mismatch: got %s want %s", ctx.CurrentDBPath, currentDBPath)
	}
	if utils.CanonicalizePath(ctx.TownDBPath) != utils.CanonicalizePath(townDBPath) {
		t.Fatalf("town db mismatch: got %s want %s", ctx.TownDBPath, townDBPath)
	}

	if ctx := getReadyDatabaseContext(townDBPath); ctx != nil {
		t.Fatalf("expected nil context when using town db, got %+v", ctx)
	}
}
