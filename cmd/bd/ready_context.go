package main

import (
	"fmt"
	"path/filepath"

	"github.com/steveyegge/beads/internal/beads"
	"github.com/steveyegge/beads/internal/configfile"
	"github.com/steveyegge/beads/internal/routing"
	"github.com/steveyegge/beads/internal/utils"
)

type readyDatabaseContext struct {
	CurrentDBPath string
	TownDBPath    string
}

func getReadyDatabaseContext(currentDBPath string) *readyDatabaseContext {
	if currentDBPath == "" {
		return nil
	}

	currentBeadsDir := beads.FindBeadsDir()
	if currentBeadsDir == "" {
		return nil
	}

	townBeadsDir, _, err := routing.ResolveBeadsDirForRig("hq", currentBeadsDir)
	if err != nil {
		return nil
	}

	townDBPath := databasePathForBeadsDir(townBeadsDir)
	if townDBPath == "" {
		return nil
	}

	if samePath(currentDBPath, townDBPath) {
		return nil
	}

	return &readyDatabaseContext{
		CurrentDBPath: currentDBPath,
		TownDBPath:    townDBPath,
	}
}

func databasePathForBeadsDir(beadsDir string) string {
	if beadsDir == "" {
		return ""
	}
	cfg, err := configfile.Load(beadsDir)
	if err == nil && cfg != nil {
		return cfg.DatabasePath(beadsDir)
	}
	return filepath.Join(beadsDir, "beads.db")
}

func samePath(left, right string) bool {
	return utils.CanonicalizePath(left) == utils.CanonicalizePath(right)
}

func maybePrintReadyDatabaseContext(currentDBPath string) {
	if jsonOutput {
		return
	}

	ctx := getReadyDatabaseContext(currentDBPath)
	if ctx == nil {
		return
	}

	fmt.Printf("Database: %s\n", ctx.CurrentDBPath)
	fmt.Printf("Town database: %s\n", ctx.TownDBPath)
	fmt.Printf("Hint: run bd --db %s ready to query town ready work.\n\n", ctx.TownDBPath)
}
