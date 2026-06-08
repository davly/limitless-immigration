// Package firewall implements R145.C FIREWALL-TEST-DISCIPLINE pin for
// limitless-immigration — structural firewall against internal/ + cmd/
// drift.
package firewall

import (
	"os"
	"path/filepath"
	"sort"
)

// ExpectedPackages returns the canonical list of internal/ packages.
// 7 packages: 5 cohort (firewall + honest + legal + manifest + mirrormark)
// + 2 domain (immigration-rules corpus-gate + eta per-nationality signal engine).
func ExpectedPackages() []string {
	return []string{
		"eta",
		"firewall",
		"honest",
		"immigration-rules",
		"legal",
		"manifest",
		"mirrormark",
	}
}

// ExpectedBinaries returns the canonical list of cmd/ binaries.
func ExpectedBinaries() []string {
	return []string{
		"limitless-immigration",
	}
}

// ScanInternal returns the subdirectories under repoRoot/internal/.
func ScanInternal(repoRoot string) ([]string, error) {
	return scanGoSubtree(filepath.Join(repoRoot, "internal"))
}

// ScanCmd returns the subdirectories under repoRoot/cmd/.
func ScanCmd(repoRoot string) ([]string, error) {
	cmdDir := filepath.Join(repoRoot, "cmd")
	entries, err := os.ReadDir(cmdDir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		mainGo := filepath.Join(cmdDir, e.Name(), "main.go")
		if _, err := os.Stat(mainGo); err == nil {
			out = append(out, e.Name())
		}
	}
	sort.Strings(out)
	return out, nil
}

func scanGoSubtree(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		subPath := filepath.Join(root, name)
		hasGo, err := dirHasGoFile(subPath)
		if err != nil {
			return nil, err
		}
		if hasGo {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out, nil
}

func dirHasGoFile(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) == ".go" {
			return true, nil
		}
	}
	return false, nil
}
