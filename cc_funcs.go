// ##################################################################
// CrunchyCleaner
// Here are all tech functions.
// ##################################################################

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/shirou/gopsutil/v3/disk"
)

// getDiskMetrics uses gopsutil to fetch precise, platform-independent disk data.
func getDiskMetrics() (freeGB float64, totalStr, freeStr string) {
	const GB = 1024 * 1024 * 1024

	// On Linux, "/" is the root. On Windows, gopsutil handles the mapping to the system drive.
	usage, err := disk.Usage("/")
	if err != nil {
		// Fallback for specific Windows environments if "/" fails
		usage, err = disk.Usage("C:")
		if err != nil {
			return 0, "N/A", "N/A"
		}
	}

	freeGB = float64(usage.Free) / GB
	totalStr = fmt.Sprintf("%.2f GB", float64(usage.Total)/GB)
	freeStr = fmt.Sprintf("%.2f GB", freeGB)
	return
}

func formatBytes(bytes int64) string {
	const (
		mb = 1024 * 1024
		gb = 1024 * mb
	)
	if bytes >= gb {
		return fmt.Sprintf("%.2f GB", float64(bytes)/gb)
	}
	return fmt.Sprintf("%.2f MB", float64(bytes)/mb)
}

// getDirSize utilizes WalkDir (introduced in Go 1.16), which is significantly
// faster than the older filepath.Walk because it avoids unnecessary Lstat calls.
func getDirSize(path string) int64 {
	var size int64
	// Resolve glob patterns (e.g., paths containing '*')
	matches := globMatches(path)

	for _, m := range matches {
		size += getPathSize(m)
	}
	return size
}

func getPathSize(path string) int64 {
	var size int64
	_ = filepath.WalkDir(path, func(_ string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			if info, err := d.Info(); err == nil {
				size += info.Size()
			}
		}
		return nil
	})
	return size
}

func globMatches(path string) []string {
	expanded := expandHome(path)
	if !filepath.IsAbs(expanded) {
		return nil
	}
	matches, err := filepath.Glob(expanded)
	if err != nil {
		return nil
	}
	return matches
}

func getCacheSize(paths []string) (int64, bool) {
	seen := make(map[string]struct{})
	var size int64
	for _, path := range paths {
		matches := globMatches(path)
		for _, match := range matches {
			match = filepath.Clean(match)
			if _, exists := seen[match]; exists {
				continue
			}
			seen[match] = struct{}{}
			size += getDirSize(match)
		}
	}
	return size, len(seen) > 0
}

// expandHome resolves the shorthand '~/ ' to the absolute user home directory
func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

// function to scan which programs actually exist on the disk
func scanForExisting() []CacheEntry {
	return scanPrograms(getPrograms())
}

func scanPrograms(programs []Program) []CacheEntry {
	existing := []CacheEntry{}
	for _, p := range programs {
		totalSize, found := getCacheSize(p.Paths)
		if found && totalSize > 0 {
			existing = append(existing, CacheEntry{Program: p, Size: totalSize})
		}
	}
	return existing
}

func selectionSummary(entries []CacheEntry) string {
	var selected, total int64
	for _, entry := range entries {
		total += entry.Size
		if entry.Checked {
			selected += entry.Size
		}
	}
	return fmt.Sprintf("%s / %s selected", formatBytes(selected), formatBytes(total))
}

func toggleAll(entries []CacheEntry) {
	allChecked := true
	for _, entry := range entries {
		if !entry.Checked {
			allChecked = false
			break
		}
	}
	for i := range entries {
		entries[i].Checked = !allChecked
	}
}

func deletePath(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}

	if !info.IsDir() {
		if err := os.Remove(path); err == nil {
			return info.Size()
		}
		return 0
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Print(CLEARLINE)
		logWarn("Cannot read " + path + ": " + err.Error())
		return 0
	}

	var removed int64
	for _, e := range entries {
		full := filepath.Join(path, e.Name())
		size := getPathSize(full)
		if err := os.RemoveAll(full); err != nil {
			fmt.Print(CLEARLINE)
			logWarn("Skipped " + e.Name() + ": " + err.Error())
			continue
		}
		removed += size
	}
	return removed
}
