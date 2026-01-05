// ##################################################################
// CrunchyCleaner
// Made by: Knuspii, (M)
// Project: https://github.com/knuspii/crunchycleaner
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.
// ##################################################################

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/eiannone/keyboard"
)

// Global constants for UI and Versioning
const (
	CC_VERSION = "2.2"
	COLS       = 62
	LINES      = 30
	YELLOW     = "\033[33m"
	CYAN       = "\033[36m"
	GREEN      = "\033[32m"
	RC         = "\033[0m" // Reset Color
)

var (
	goos              = runtime.GOOS
	getcols, getlines int
	// CLI Flags
	Flagversion = flag.Bool("version", false, "Display version information")
	Flagnoinit  = flag.Bool("no-init", false, "Skip terminal resizing and environment initialization")
	Flagdryrun  = flag.Bool("dry-run", false, "Simulation mode without deleting files (for testing)")
)

// Program represents a target application and its associated cache directories
type Program struct {
	Name    string
	Paths   []string // List of paths (supports wildcards/globbing)
	Checked bool     // Selection state in the menu
}

// --- Helper Functions ---

// clearScreen handles cross-platform terminal clearing
func clearScreen() {
	if goos == "windows" {
		// Windows CMD requires an external call to 'cls'
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		// Unix-like systems use ANSI escape sequences
		fmt.Print("\033[H\033[2J")
	}
}

// cc_exit provides a clean termination of the application
func cc_exit() {
	keyboard.Close()
	fmt.Printf("\nExiting CrunchyCleaner. Goodbye!\n")
	os.Exit(0)
}

func pause() {
	fmt.Printf("\nPress [ENTER] to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// line draws a formatted horizontal separator
func line() {
	fmt.Printf("%s#%s~%s\n", YELLOW, strings.Repeat("-", COLS-2), RC)
}

// spinner visualizes background tasks to keep the UI responsive
func spinner(text string, done chan bool) {
	frames := []string{"|", "/", "-", "\\"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Print("\r\033[K") // Clear the line when task completes
			return
		default:
			fmt.Printf("\r%s%s%s %s%s %s%s%s ", YELLOW, frames[i%len(frames)], RC, CYAN, text, YELLOW, frames[i%len(frames)], RC)
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}
}

// runCommand wraps exec.Command to provide easy access to combined stdout/stderr
func runCommand(cmd []string) (string, error) {
	if len(cmd) == 0 {
		return "", errors.New("command is empty")
	}
	c := exec.Command(cmd[0], cmd[1:]...)
	outBytes, err := c.CombinedOutput()
	out := strings.TrimSpace(string(outBytes))

	if err != nil {
		return out, fmt.Errorf("command failed: %v", err)
	}
	return out, nil
}

// initApp prepares the terminal environment (Title, Resize, User Info)
func initApp() {
	fmt.Printf("Initializing CrunchyCleaner %s...\n", CC_VERSION)

	// Fetching the current system user for the display
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("Username: unknown\n")
	} else {
		name := usr.Username
		// Strip domain/machine name prefix on Windows
		if goos == "windows" && strings.Contains(name, "\\") {
			parts := strings.Split(name, "\\")
			name = parts[len(parts)-1]
		}
		fmt.Printf("Username: %s\n", name)
	}

	// Set Terminal Title via ANSI sequence
	fmt.Printf("\033]0;CrunchyCleaner %s\007", CC_VERSION)

	// OS-specific Terminal Resizing
	switch goos {
	case "windows":
		runCommand([]string{"cmd", "/C", "title", "CrunchyCleaner"})
		// Use PowerShell to force specific Window and Buffer size
		psCmd := fmt.Sprintf(
			`$w=(Get-Host).UI.RawUI; $s=New-Object System.Management.Automation.Host.Size(%d,%d); $w.WindowSize=$s; $w.BufferSize=$s`,
			COLS, LINES,
		)
		runCommand([]string{"powershell", "-NoProfile", "-Command", psCmd})
	}

	// Generic ANSI resize fallback for modern terminals
	fmt.Printf("\033[8;%d;%dt", LINES, COLS)

	// Verify if the terminal actually resized to match the UI expectations
	sizeErr := error(nil)
	if goos == "windows" {
		out, err := runCommand([]string{
			"powershell",
			"-NoProfile",
			"-Command",
			"$s=$Host.UI.RawUI.WindowSize; Write-Output \"$($s.Width) $($s.Height)\"",
		})
		if err == nil {
			fmt.Sscanf(strings.TrimSpace(out), "%d %d", &getcols, &getlines)
		} else {
			sizeErr = err
		}
	} else {
		out, err := runCommand([]string{"sh", "-c", "stty size < /dev/tty"})
		if err == nil {
			fmt.Sscanf(strings.TrimSpace(out), "%d %d", &getlines, &getcols)
		} else {
			sizeErr = err
		}
	}

	if sizeErr != nil || getcols == 0 || getlines == 0 {
		fmt.Printf("System: Could not detect terminal size.\n")
		time.Sleep(2 * time.Second)
	} else if getcols != COLS || getlines != LINES {
		fmt.Printf("System: Terminal size mismatch (Got %dx%d, Expected %dx%d)\n", getcols, getlines, COLS, LINES)
		time.Sleep(2 * time.Second)
	} else {
		fmt.Printf("System: Terminal size optimized (%dx%d)\n", getcols, getlines)
	}

	time.Sleep(1 * time.Second)
}

// isSafePath checks whether a given path is safe to delete.
// It prevents dangerous operations like deleting root directories
// or very short paths that could wipe critical system locations.
func isSafePath(p string) bool {
	// Clean the path to resolve ".." and remove double slashes
	p = filepath.Clean(p)
	pLower := strings.ToLower(p)

	// Absolute Base Protection: Root Directories
	if p == "/" || p == "\\" {
		return false
	}

	if runtime.GOOS == "windows" {
		// Block "C:", "C:\", "D:\" etc.
		if len(p) <= 3 && strings.Contains(p, ":") {
			return false
		}

		// Critical Windows System Directories
		// Even if globbed incorrectly, these hardcoded paths provide a safety net.
		systemBlacklist := []string{
			"c:\\windows",
			"c:\\windows\\system32",
			"c:\\users",
			"c:\\program files",
			"c:\\program files (x86)",
			"c:\\programdata",
		}
		for _, s := range systemBlacklist {
			if pLower == s {
				return false
			}
		}

		// Prevent deleting the entire User Profile directory
		userProfile := strings.ToLower(os.Getenv("USERPROFILE"))
		if pLower == userProfile {
			return false
		}

	} else {
		// Unix/Linux System Blacklist
		systemBlacklist := []string{
			"/etc", "/bin", "/sbin", "/lib", "/usr", "/boot", "/root", "/home", "/proc", "/sys", "/dev",
		}
		for _, s := range systemBlacklist {
			if pLower == s {
				return false
			}
		}
	}

	// Most cache directories are deep in the file system (e.g., AppData/Local/Temp).
	// We count the number of separators to ensure we aren't deleting a high-level folder.
	// This allows "C:\Users\Name\AppData", but blocks "C:\Users\Name" or "C:\Users".
	parts := strings.Split(strings.Trim(p, string(os.PathSeparator)), string(os.PathSeparator))
	if len(parts) < 3 {
		return false
	}

	return true
}

func getPrograms() []Program {
	if runtime.GOOS == "windows" {
		// Windows
		home, _ := os.UserHomeDir()
		appData := os.Getenv("APPDATA")
		localAppData := os.Getenv("LOCALAPPDATA")
		programFilesX86 := os.Getenv("ProgramFiles(x86)")
		programFiles := os.Getenv("ProgramFiles")
		winDir := os.Getenv("WINDIR")
		return []Program{
			{"System Logs (Admin)", []string{
				filepath.Join(winDir, "Panther"),
				filepath.Join(winDir, "Logs"),
			}, false},
			{"System Temp Folders (Admin)", []string{filepath.Join(winDir, "Temp")}, false},
			{"Update Logs (Admin)", []string{filepath.Join(winDir, "SoftwareDistribution/Download")}, false},
			{"User Temp Folder", []string{filepath.Join(localAppData, "Temp")}, false},
			{"Thumbnail Cache", []string{filepath.Join(localAppData, "Microsoft/Windows/Explorer")}, false},
			{"Firefox Cache", []string{
				filepath.Join(localAppData, "Mozilla/Firefox/Profiles/*/cache2"),
				filepath.Join(localAppData, "Mozilla/Firefox/Profiles/*/jumpListCache"),
				filepath.Join(appData, "Mozilla/Firefox/Profiles/*/shader-cache"),
			}, false},
			{"Chrome Cache", []string{
				filepath.Join(localAppData, "Google/Chrome/User Data/Default/Cache"),
				filepath.Join(localAppData, "Google/Chrome/User Data/Default/Code Cache"),
				filepath.Join(localAppData, "Google/Chrome/User Data/*/Cache"),
				filepath.Join(localAppData, "Google/Chrome/User Data/Default/Media Cache"),
			}, false},
			{"Edge Cache", []string{
				filepath.Join(localAppData, "Microsoft/Edge/User Data/Default/Cache"),
				filepath.Join(localAppData, "Microsoft/Edge/User Data/*/Cache"),
				filepath.Join(localAppData, "Microsoft/Edge/User Data/Default/Media Cache"),
			}, false},
			{"Brave Cache", []string{
				filepath.Join(localAppData, "BraveSoftware/Brave-Browser/User Data/Default/Cache"),
				filepath.Join(localAppData, "BraveSoftware/Brave-Browser/User Data/*/Cache"),
				filepath.Join(localAppData, "BraveSoftware/Brave-Browser/User Data/Default/Media Cache"),
			}, false},
			{"Opera Cache", []string{
				filepath.Join(localAppData, "Opera Software/Opera Stable/Cache"),
				filepath.Join(localAppData, "Opera Software/Opera Stable/Code Cache"),
			}, false},
			{"Thunderbird Cache", []string{
				filepath.Join(localAppData, "Thunderbird/Profiles/*/cache2"),
			}, false},
			{"Steam AppCache", []string{
				filepath.Join(programFilesX86, "Steam/appcache"),
				filepath.Join(programFiles, "Steam/appcache"),
			}, false},
			{"Epic Games Cache", []string{
				filepath.Join(localAppData, "EpicGamesLauncher/Saved/webcache"),
			}, false},
			{"Discord Cache", []string{
				filepath.Join(appData, "discord/Cache"),
				filepath.Join(appData, "discord/Code Cache"),
				filepath.Join(appData, "discord/GPUCache"),
			}, false},
			{"Spotify Storage", []string{filepath.Join(localAppData, "Spotify/Storage")}, false},
			{"VS Code Cache", []string{
				filepath.Join(appData, "Code/Cache"),
				filepath.Join(appData, "Code/CachedData"),
				filepath.Join(appData, "Code/CachedExtensionVSIXs"),
				filepath.Join(appData, "Code/User/workspaceStorage"),
				filepath.Join(appData, "Code/GPUCache"),
			}, false},
			{"DirectX Shader Cache", []string{
				filepath.Join(localAppData, "D3DSCache"),
				filepath.Join(localAppData, "NVIDIA/GLCache"),
			}, false},
			{"Go Build Cache", []string{filepath.Join(localAppData, "go-build")}, false},
			{"Pip Cache", []string{filepath.Join(localAppData, "pip/Cache")}, false},
			{"NPM Cache", []string{filepath.Join(appData, "npm-cache/_cacache")}, false},
			{"Yarn Cache", []string{
				filepath.Join(localAppData, "Yarn/Cache"),
				filepath.Join(appData, "Yarn/Cache"),
			}, false},
			{"Cargo Cache", []string{
				filepath.Join(home, ".cargo/registry/cache"),
				filepath.Join(home, ".cargo/git/db"),
			}, false},
		}
	} else {
		// Linux
		home, _ := os.UserHomeDir()
		return []Program{
			{"System Logs (Root)", []string{"/var/log/*.log"}, false},
			{"System Temp Folders (Root)", []string{"/tmp"}, false},
			{"Thumbnail Cache", []string{filepath.Join(home, ".cache/thumbnails")}, false},
			{"Firefox Cache", []string{filepath.Join(home, ".cache/mozilla/firefox/*/cache2")}, false},
			{"Chromium Cache", []string{
				filepath.Join(home, ".cache/chromium/*/Cache"),
				filepath.Join(home, ".cache/chromium/*/Code Cache"),
			}, false},
			{"Edge Cache", []string{
				filepath.Join(home, ".cache/microsoft-edge/*/Cache"),
				filepath.Join(home, ".cache/microsoft-edge/*/Code Cache"),
			}, false},
			{"Brave Cache", []string{
				filepath.Join(home, ".cache/BraveSoftware/Brave-Browser/*/Cache"),
				filepath.Join(home, ".cache/BraveSoftware/Brave-Browser/*/Code Cache"),
			}, false},
			{"Opera Cache", []string{
				filepath.Join(home, ".cache/opera/Cache"),
				filepath.Join(home, ".config/opera/Cache"),
			}, false},
			{"Thunderbird Cache", []string{
				filepath.Join(home, ".cache/thunderbird/*/cache2"),
			}, false},
			{"Steam Cache", []string{
				filepath.Join(home, ".steam/steam/appcache"),
				filepath.Join(home, ".local/share/Steam/appcache"),
				filepath.Join(home, ".local/share/Steam/config/htmlcache"),
			}, false},
			{"Epic Games (Heroic/Lutris) Cache", []string{
				filepath.Join(home, ".config/heroic/WebCache"),
				filepath.Join(home, ".local/share/lutris/runtime"),
			}, false},
			{"Discord Cache", []string{
				filepath.Join(home, ".config/discord/Cache"),
				filepath.Join(home, ".config/discord/Code Cache"),
				filepath.Join(home, ".config/discord/GPUCache"),
			}, false},
			{"Spotify Cache", []string{filepath.Join(home, ".cache/spotify")}, false},
			{"VS Code Cache", []string{
				filepath.Join(home, ".config/Code/Cache"),
				filepath.Join(home, ".config/Code/CachedData"),
				filepath.Join(home, ".config/Code/User/workspaceStorage"),
				filepath.Join(home, ".config/Code/GPUCache"),
			}, false},
			{"Mesa Shader Cache", []string{filepath.Join(home, ".cache/mesa_shader_cache")}, false},
			{"Go Build Cache", []string{filepath.Join(home, ".cache/go-build")}, false},
			{"Pip Cache", []string{filepath.Join(home, ".cache/pip")}, false},
			{"NPM Cache", []string{filepath.Join(home, ".npm/_cacache")}, false},
			{"Yarn Cache", []string{filepath.Join(home, ".cache/yarn")}, false},
			{"Cargo Cache", []string{filepath.Join(home, ".cargo/registry/cache")}, false},
		}
	}
}

// getDiskMetrics calculates total and free space for the root/system drive
func getDiskMetrics() (freeGB float64, totalStr string, freeStr string) {
	// Default values if something fails
	totalStr = "N/A"
	freeStr = "N/A"

	if goos == "windows" {
		// We use a single PowerShell command to get both values to reduce overhead
		cmdArgs := []string{
			"powershell", "-NoProfile", "-Command",
			"(Get-PSDrive C | Select-Object Used, Free) | ForEach-Object { \"$($_.Used) $($_.Free)\" }",
		}
		out, err := exec.Command(cmdArgs[0], cmdArgs[1:]...).Output()

		if err == nil {
			parts := strings.Fields(string(out))
			if len(parts) >= 2 {
				used, _ := strconv.ParseFloat(parts[0], 64)
				free, _ := strconv.ParseFloat(parts[1], 64)

				totalGB := (used + free) / 1024 / 1024 / 1024
				freeGB = free / 1024 / 1024 / 1024

				totalStr = fmt.Sprintf("%.2f GB", totalGB)
				freeStr = fmt.Sprintf("%.2f GB", freeGB)
			}
		}
	} else {
		// Standard Unix 'df' command
		// -B1 ensures output in bytes for better precision before converting to GB
		out, err := exec.Command("sh", "-c", "df -B1 --output=size,avail / | tail -1").Output()

		if err == nil {
			parts := strings.Fields(string(out))
			if len(parts) >= 2 {
				totalBytes, _ := strconv.ParseFloat(parts[0], 64)
				freeBytes, _ := strconv.ParseFloat(parts[1], 64)

				freeGB = freeBytes / 1024 / 1024 / 1024
				totalStr = fmt.Sprintf("%.2f GB", totalBytes/1024/1024/1024)
				freeStr = fmt.Sprintf("%.2f GB", freeGB)
			}
		}
	}
	return freeGB, totalStr, freeStr
}

// formatMB converts bytes to a string representing Megabytes
func formatMB(bytes int64) string {
	mb := float64(bytes) / 1024 / 1024
	return fmt.Sprintf("%.2f MB", mb)
}

// getDirSize remains the same (calculating in bytes first for precision)
func getDirSize(path string) int64 {
	var size int64
	matches, _ := filepath.Glob(expandHome(path))
	for _, m := range matches {
		filepath.Walk(m, func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() {
				size += info.Size()
			}
			return nil
		})
	}
	return size
}

// expandHome resolves the shorthand '~/ ' to the absolute user home directory
func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

// --- Menu UI Logic ---

func showBanner() {
	_, total, free := getDiskMetrics()
	fmt.Printf(`%s  ____________________     .-.
 |  |              |  |    |_|
 |[]|              |[]|    | |
 |  |              |  |    |=|
 |  |              |  |  .=/I\=.
 |  |______________|  | ////V\\\\
 |  |______________|  | |#######|
 |                    | |||||||||
 |     ____________   |
 |    | __      |  |  | %sCrunchyCleaner%s
 |    ||  |     |  |  | Made by: Knuspii, (M)
 |    ||__|     |  |  | Version: %s
 |____|_________|__|__| Disk-Space: %s / %s%s
`, YELLOW, RC, YELLOW, CC_VERSION, free, total, RC)
	line()
}

func logInfo(msg string) {
	fmt.Printf("\r\033[K%s[+] %s%s\n", CYAN, msg, RC)
}

func logWarn(msg string) {
	fmt.Printf("\r\033[K\033[33m[!] %s%s\n", msg, RC)
}

func logOK(msg string) {
	fmt.Printf("\r\033[K\033[32m[✓] %s%s\n", msg, RC)
}

// renderMenu draws the interactive selection list
func renderMenu(existing []Program, idx int, fullRedraw bool) {
	if fullRedraw {
		clearScreen()
		showBanner()
		fmt.Printf("Use ↑/↓ or W/S to navigate | [ENTER] to select | [C] to clean\n")
		fmt.Printf("Folders found: [%d]\n", len(existing))
	}

	// Render each detected program entry
	for i := range existing {
		cursor := "    "
		// Highlight the currently selected entry
		if i == idx {
			cursor = YELLOW + "  >_" + RC
		}
		// Checkbox indicator for selection state
		check := "[ ]"
		if existing[i].Checked {
			check = "[" + GREEN + "X" + RC + "]"
		}
		// Clear the current line and print the menu entry
		fmt.Printf("\r\033[K%s%s %s\n", cursor, check, existing[i].Name)
	}
}

// handleMenu manages user input for navigation and selection
func handleMenu() {
	clearScreen()
	showBanner()

	// Initial scan of the filesystem to find existing directories
	done := make(chan bool)
	go spinner("Scanning filesystem", done)

	// Retrieve all known programs and filter only those
	// whose cache paths actually exist on the system
	allPrograms := getPrograms()
	existing := []Program{}
	for _, p := range allPrograms {
		var totalSize int64
		found := false
		for _, path := range p.Paths {
			matches, _ := filepath.Glob(expandHome(path))
			if len(matches) > 0 {
				found = true
				totalSize += getDirSize(path) // Calculate size for this program
			}
		}
		if found {
			// We modify the name slightly to include the size
			p.Name = fmt.Sprintf("%-30s %s(%s)%s", p.Name, YELLOW, formatMB(totalSize), RC)
			existing = append(existing, p)
		}
	}
	time.Sleep(1 * time.Second)
	done <- true

	// Abort if nothing was detected
	if len(existing) == 0 {
		fmt.Printf("No cache directories found on your system.\n")
		pause()
		return
	}

	// Enable raw keyboard input mode
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	idx := 0
	renderMenu(existing, idx, true)
	// Main Input Loop
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			break
		}

		updated := false

		// Navigation and selection controls
		if key == keyboard.KeyArrowUp || char == 'w' || char == 'W' {
			if idx > 0 {
				idx--
				updated = true
			}
		} else if key == keyboard.KeyArrowDown || char == 's' || char == 'S' {
			if idx < len(existing)-1 {
				idx++
				updated = true
			}
		} else if char == ' ' || key == keyboard.KeyEnter || key == keyboard.KeySpace {
			existing[idx].Checked = !existing[idx].Checked
			updated = true
		} else if char == 'a' || char == 'A' {
			// Toggle "Select All" logic
			allChecked := true
			for _, p := range existing {
				if !p.Checked {
					allChecked = false
					break
				}
			}
			for i := range existing {
				existing[i].Checked = !allChecked
			}
			updated = true
		} else if char == 'c' || char == 'C' {
			clearScreen()
			showBanner()
			runCleanup(existing)
		} else if key == keyboard.KeyCtrlC {
			cc_exit()
		}

		// Redraw menu entries in-place if state changed
		if updated {
			// Move cursor up to the start of the menu list
			fmt.Printf("\033[%dA", len(existing))
			renderMenu(existing, idx, false)
		}
	}
}

// runCleanup executes the deletion logic (or simulation if dry-run is active)
func runCleanup(programs []Program) {
	beforeFree, _, _ := getDiskMetrics()

	statusMsg := "Cleaning selected caches"
	if *Flagdryrun {
		statusMsg = "[DRY RUN] Simulating cleanup"
	}

	done := make(chan bool)
	go spinner(statusMsg, done)

	fmt.Printf("Cleaning caches started...\n")

	if *Flagdryrun {
		fmt.Printf("%sNOTE: Dry run active. No files will actually be deleted.%s\n", YELLOW, RC)
	} else {
		fmt.Printf("You use this tool at your own risk!\n")
	}

	fmt.Printf("Press [CTRL+C] to abort\n")
	time.Sleep(1 * time.Second)

	count := 0
	for _, p := range programs {
		if !p.Checked {
			continue
		}
		count++

		for _, path := range p.Paths {
			// Expand wildcards (e.g., /Profiles/*/cache2)
			matches, _ := filepath.Glob(expandHome(path))

			for _, m := range matches {
				if *Flagdryrun {
					logInfo(fmt.Sprintf("Would empty directory: %s", m))
				} else {
					if !isSafePath(m) {
						logWarn("Skipped unsafe path: " + m)
						continue
					}

					// Get file information to check if it's a directory
					info, err := os.Stat(m)
					if err != nil {
						continue
					}

					if !info.IsDir() {
						// If it's just a file, delete it directly
						err := os.Remove(m)
						if err != nil {
							logWarn("Could not delete file " + m + ": " + err.Error())
						}
					} else {
						// If it's a directory, read its contents
						entries, err := os.ReadDir(m)
						if err != nil {
							logWarn("Could not read directory " + m + ": " + err.Error())
							continue
						}

						// Loop through all files and subfolders inside the directory
						for _, entry := range entries {
							fullPath := filepath.Join(m, entry.Name())

							// Delete the entry (file or subfolder)
							err := os.RemoveAll(fullPath)
							if err != nil {
								// Just log the error message provided by the OS
								// This covers "Permission Denied", "In Use", etc.
								logWarn(fmt.Sprintf("Skipped %s: %v", entry.Name(), err))
							}
						}
					}
				}
			}
		}

		if !*Flagdryrun {
			logOK("Cleaned " + p.Name)
		}
	}

	done <- true
	fmt.Print("\r\033[K") // Clear the spinner line

	if *Flagdryrun {
		logOK("Simulation finished. No files were removed.")
	} else {
		logOK("Cleaning finished")
	}

	if count == 0 {
		fmt.Printf("Nothing selected to clean.\n")
	}

	// Calculate and display space savings
	afterFree, _, _ := getDiskMetrics()

	// Convert GB difference back to MB for display
	totalCleanedMB := (afterFree - beforeFree) * 1024

	if totalCleanedMB < 0 || *Flagdryrun {
		totalCleanedMB = 0
	}

	line()
	label := "Cleaned"
	if *Flagdryrun {
		label = "NOTHING CLEANED (Dry Run)"
	}

	fmt.Printf(" %s: %s%.2f MB%s\n", label, GREEN, totalCleanedMB, RC)
	line()

	fmt.Printf("\nPress [ENTER] to exit...")
	keyboard.Close()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cc_exit()
}

func main() {
	flag.Parse()

	// Capture OS Interrupts (like Ctrl+C) for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cc_exit()
	}()

	if *Flagversion {
		fmt.Printf("CrunchyCleaner %s\n", CC_VERSION)
		os.Exit(0)
	}

	if !*Flagnoinit {
		initApp()
	}

	handleMenu()
}
