package firewall

import (
	"path/filepath"
	"runtime"
	"sort"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestExpectedPackages_NonEmpty(t *testing.T) {
	if len(ExpectedPackages()) == 0 {
		t.Fatal("empty")
	}
}

func TestExpectedBinaries_NonEmpty(t *testing.T) {
	if len(ExpectedBinaries()) == 0 {
		t.Fatal("empty")
	}
}

func TestExpectedPackages_Sorted(t *testing.T) {
	if !sort.StringsAreSorted(ExpectedPackages()) {
		t.Fatal("not sorted")
	}
}

func TestExpectedPackages_Unique(t *testing.T) {
	seen := map[string]int{}
	for i, p := range ExpectedPackages() {
		if prev, ok := seen[p]; ok {
			t.Errorf("dup %q at %d and %d", p, prev, i)
		}
		seen[p] = i
	}
}

func TestExpectedPackages_PinnedCount(t *testing.T) {
	const expected = 7 // +eta (iter: per-nationality ETA signal engine)
	if got := len(ExpectedPackages()); got != expected {
		t.Fatalf("got %d, want %d", got, expected)
	}
}

func TestFirewall_EveryExpectedPackageExistsOnDisk(t *testing.T) {
	root := repoRoot(t)
	onDisk, err := ScanInternal(root)
	if err != nil {
		t.Fatalf("ScanInternal: %v", err)
	}
	onDiskSet := map[string]bool{}
	for _, p := range onDisk {
		onDiskSet[p] = true
	}
	for _, expected := range ExpectedPackages() {
		if !onDiskSet[expected] {
			t.Errorf("R145.C drift: expected %q missing from disk", expected)
		}
	}
}

func TestFirewall_EveryOnDiskPackageInExpectedList(t *testing.T) {
	root := repoRoot(t)
	onDisk, err := ScanInternal(root)
	if err != nil {
		t.Fatalf("ScanInternal: %v", err)
	}
	expectedSet := map[string]bool{}
	for _, p := range ExpectedPackages() {
		expectedSet[p] = true
	}
	for _, found := range onDisk {
		if !expectedSet[found] {
			t.Errorf("R145.C drift: %q on disk not in ExpectedPackages", found)
		}
	}
}

func TestFirewall_EveryExpectedBinaryExistsOnDisk(t *testing.T) {
	root := repoRoot(t)
	onDisk, err := ScanCmd(root)
	if err != nil {
		t.Fatalf("ScanCmd: %v", err)
	}
	onDiskSet := map[string]bool{}
	for _, b := range onDisk {
		onDiskSet[b] = true
	}
	for _, expected := range ExpectedBinaries() {
		if !onDiskSet[expected] {
			t.Errorf("R145.C drift: binary %q missing", expected)
		}
	}
}
