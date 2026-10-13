package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatBytesUsesReadableUnits(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		want  string
	}{
		{name: "megabytes", bytes: 512 * 1024 * 1024, want: "512.00 MB"},
		{name: "gigabytes", bytes: 3 * 1024 * 1024 * 1024, want: "3.00 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatBytes(tt.bytes); got != tt.want {
				t.Fatalf("formatBytes(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestSelectionSummaryTracksSelectedAndTotalSizes(t *testing.T) {
	entries := []CacheEntry{
		{Program: Program{Name: "first", Checked: true}, Size: 1024 * 1024 * 1024},
		{Program: Program{Name: "second"}, Size: 2 * 1024 * 1024 * 1024},
		{Program: Program{Name: "third", Checked: true}, Size: 512 * 1024 * 1024},
	}

	if got, want := selectionSummary(entries), "1.50 GB / 3.50 GB selected"; got != want {
		t.Fatalf("selectionSummary() = %q, want %q", got, want)
	}
}

func TestToggleAllSelectsThenDeselectsEveryEntry(t *testing.T) {
	entries := []CacheEntry{
		{Program: Program{Name: "first", Checked: true}},
		{Program: Program{Name: "second"}},
	}

	toggleAll(entries)
	for _, entry := range entries {
		if !entry.Checked {
			t.Fatal("toggleAll() did not select every entry")
		}
	}

	toggleAll(entries)
	for _, entry := range entries {
		if entry.Checked {
			t.Fatal("toggleAll() did not deselect every entry")
		}
	}
}

func TestCacheSizeDoesNotDoubleCountOverlappingPatterns(t *testing.T) {
	root := t.TempDir()
	cacheDir := filepath.Join(root, "cache")
	if err := os.Mkdir(cacheDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cacheDir, "data.bin"), []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	size, found := getCacheSize([]string{cacheDir, filepath.Join(root, "cache*")})
	if !found {
		t.Fatal("getCacheSize() did not find the cache directory")
	}
	if size != 4 {
		t.Fatalf("getCacheSize() = %d bytes, want 4", size)
	}
}

func TestCacheSizeRejectsRelativePaths(t *testing.T) {
	size, found := getCacheSize([]string{"."})
	if found || size != 0 {
		t.Fatalf("getCacheSize() = (%d, %t), want (0, false) for relative path", size, found)
	}
}

func TestMenuHelpFitsConfiguredWidth(t *testing.T) {
	for _, line := range menuHelpLines {
		if width := len([]rune(line)); width > COLS {
			t.Fatalf("menu help line %q is %d columns wide, limit is %d", line, width, COLS)
		}
	}
	if help := strings.Join(menuHelpLines[:], "\n"); !strings.Contains(help, "A: Select/deselect all") {
		t.Fatalf("menu help does not explain the select-all toggle: %q", help)
	}
}

func TestDeletePathReturnsRemovedBytes(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "cache")
	if err := os.Mkdir(cacheDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cacheDir, "first.bin"), []byte("1234"), 0o644); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(cacheDir, "nested")
	if err := os.Mkdir(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nested, "second.bin"), []byte("123456"), 0o644); err != nil {
		t.Fatal(err)
	}

	if got := deletePath(cacheDir); got != 10 {
		t.Fatalf("deletePath() removed %d bytes, want 10", got)
	}
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("deletePath() left %d entries, want 0", len(entries))
	}
}

func TestScanProgramsDropsEmptyCaches(t *testing.T) {
	root := t.TempDir()
	emptyDir := filepath.Join(root, "empty")
	filledDir := filepath.Join(root, "filled")
	if err := os.Mkdir(emptyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filledDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(filledDir, "cache.bin"), []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	entries := scanPrograms([]Program{
		{Name: "Empty", Paths: []string{emptyDir}},
		{Name: "Filled", Paths: []string{filledDir}},
	})

	if len(entries) != 1 {
		t.Fatalf("scanPrograms() returned %d entries, want 1", len(entries))
	}
	if entries[0].Name != "Filled" || entries[0].Size != 4 {
		t.Fatalf("scanPrograms() returned %#v, want Filled with 4 bytes", entries[0])
	}
}
