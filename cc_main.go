// ##################################################################
// CrunchyCleaner
// Made by: Knuspii, (M)
// Project: https://github.com/Knuspii/CrunchyCleaner
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
	CC_VERSION = "2.5"
	COLS       = 62
	LINES      = 32
	GOOS       = runtime.GOOS
	YELLOW     = "\033[33m"
	CYAN       = "\033[36m"
	GREEN      = "\033[32m"
	RC         = "\033[0m" // Reset Color
)

var (
	origCols, origLines int
	// CLI Flags
	Flagversion = flag.Bool("v", false, "Display version information")
	Flagnoinit  = flag.Bool("i", false, "Skip terminal resizing and environment initialization")
	Flagdryrun  = flag.Bool("d", false, "Simulation mode without deleting files (for testing)")
	Flagauto    = flag.Bool("a", false, "Automate cleaning (select all and start immediately)")
)

// Program represents a target application and its associated cache directories
type Program struct {
	Name    string
	Paths   []string // List of paths (supports wildcards/globbing)
	Checked bool     // Selection state in the menu
}

// ========================= HELPER FUNCTIONS =========================

// clearScreen handles cross-platform terminal clearing
func clearScreen() {
	if *Flagnoinit {
		return
	}
	if GOOS == "windows" {
		// Windows CMD requires an external call to 'cls'
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("sh", "-c", "clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

	}
	// Fallback use ANSI escape sequences
	fmt.Print("\033[H\033[2J")
}

// cc_exit provides a clean termination of the application
func cc_exit() {
	// Close keyboard
	keyboard.Close()

	// Enable cursor
	fmt.Print("\033[?25h")

	// Main Screen Buffer
	if !*Flagnoinit {
		fmt.Print("\033[?1049l")
	}

	// Restore size
	if !*Flagnoinit {
		terminalresize(origCols, origLines)
	}

	fmt.Printf("\nExiting CrunchyCleaner...\n")
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

// spinner visualizes background tasks and cleans up properly
func spinner(text string, stop chan bool, ack chan bool) {
	frames := []string{"|", "/", "-", "\\"}
	i := 0
	for {
		select {
		case <-stop:
			// Clear the line and move cursor to start
			fmt.Print("\r\033[K")
			// Signal back to main that we are done
			ack <- true
			return
		default:
			fmt.Printf("\r%s%s%s %s%s %s%s%s ", YELLOW, frames[i%len(frames)], RC, CYAN, text, YELLOW, frames[i%len(frames)], RC)
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}
}

func terminalresize(w int, h int) {
	// OS-specific Terminal Resizing
	if GOOS == "windows" {
		psCmd := fmt.Sprintf(
			`$w=(Get-Host).UI.RawUI; $s=New-Object System.Management.Automation.Host.Size(%d,%d); $w.WindowSize=$s; $w.BufferSize=$s`,
			w, h,
		)
		exec.Command("powershell", "-NoProfile", "-Command", psCmd).Run()
	}

	// Generic ANSI resize fallback for modern terminals
	fmt.Printf("\033[8;%d;%dt", h, w)
}

// initApp prepares the terminal environment (Title, Resize)
func initApp() {
	fmt.Printf("Initializing CrunchyCleaner %s...\n", CC_VERSION)

	// Alternate Screen Buffer
	if !*Flagnoinit {
		fmt.Print("\033[?1049h")
	}

	// Get current terminal size
	if GOOS == "windows" {
		cmd := exec.Command(
			"powershell", "-NoProfile", "-Command",
			"$s=$Host.UI.RawUI.WindowSize; Write-Output \"$($s.Width) $($s.Height)\"",
		)

		out, _ := cmd.Output()
		fmt.Sscanf(strings.TrimSpace(string(out)), "%d %d", &origCols, &origLines)

	} else {
		cmd := exec.Command("sh", "-c", "stty size < /dev/tty")

		out, _ := cmd.Output()
		fmt.Sscanf(strings.TrimSpace(string(out)), "%d %d", &origLines, &origCols)
	}

	// Set Terminal Title via ANSI sequence
	fmt.Printf("\033]0;CrunchyCleaner %s\007", CC_VERSION)

	// Resize terminal
	terminalresize(COLS, LINES)
}

// ========================= PROGRAMS =========================

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
			{"Font Cache (Admin)", []string{filepath.Join(winDir, "ServiceProfiles/LocalService/AppData/Local/FontCache")}, false},
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
			{"Steam Cache", []string{
				filepath.Join(programFilesX86, "Steam/appcache"),
				filepath.Join(programFiles, "Steam/appcache"),
				filepath.Join(localAppData, "Steam/htmlcache"),
			}, false},
			{"Epic Games Cache", []string{filepath.Join(localAppData, "EpicGamesLauncher/Saved/webcache")}, false},
			{"Discord Cache", []string{
				filepath.Join(appData, "discord/Cache"),
				filepath.Join(appData, "discord/Code Cache"),
				filepath.Join(appData, "discord/GPUCache"),
			}, false},
			{"Telegram Cache", []string{filepath.Join(appData, "Telegram Desktop/tdata/user_data/cache")}, false},
			{"Spotify Cache", []string{filepath.Join(localAppData, "Spotify/Storage")}, false},
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
		cache := ".cache/"
		flatpak := ".var/app/"
		return []Program{
			{"System Logs (Root)", []string{"/var/log/*.log"}, false},
			{"System Temp Folders (Root)", []string{"/tmp"}, false},
			{"Thumbnail Cache", []string{filepath.Join(home, cache, "thumbnails")}, false},
			{"Firefox Cache", []string{
				filepath.Join(home, cache, "mozilla/firefox/*/cache2"),
				filepath.Join(home, flatpak, "org.mozilla.firefox/cache/mozilla/firefox/*/cache2"),
			}, false},
			{"Chromium Cache", []string{
				filepath.Join(home, cache, "chromium/*/Cache"),
				filepath.Join(home, cache, "chromium/*/Code Cache"),
				filepath.Join(home, flatpak, "com.google.Chrome/cache/chromium/*/Cache"),
				filepath.Join(home, flatpak, "com.google.Chrome/cache/chromium/*/CodeCache"),
			}, false},
			{"Edge Cache", []string{
				filepath.Join(home, cache, "microsoft-edge/*/Cache"),
				filepath.Join(home, cache, "microsoft-edge/*/Code Cache"),
				filepath.Join(home, flatpak, "com.microsoft.Edge/cache/microsoft-edge/*/Cache"),
				filepath.Join(home, flatpak, "com.microsoft.Edge/cache/microsoft-edge/*/CodeCache"),
			}, false},
			{"Brave Cache", []string{
				filepath.Join(home, cache, "BraveSoftware/Brave-Browser/*/Cache"),
				filepath.Join(home, cache, "BraveSoftware/Brave-Browser/*/Code Cache"),
				filepath.Join(home, flatpak, "com.brave.Browser/cache/Brave-Browser/*/Cache"),
				filepath.Join(home, flatpak, "com.brave.Browser/cache/Brave-Browser/*/Code Cache"),
			}, false},
			{"Opera Cache", []string{
				filepath.Join(home, cache, "opera/Cache"),
				filepath.Join(home, ".config/opera/Cache"),
				filepath.Join(home, flatpak, "com.opera.Opera/cache/opera/Cache"),
				filepath.Join(home, flatpak, ".com.opera.Opera/config/opera/Cache"),
			}, false},
			{"Thunderbird Cache", []string{
				filepath.Join(home, cache, "thunderbird/*/cache2"),
				filepath.Join(home, flatpak, "org.mozilla.Thunderbird/cache/mozilla/Thunderbird/*/cache2"),
			}, false},
			{"Steam Cache", []string{
				filepath.Join(home, ".steam/steam/appcache"),
				filepath.Join(home, ".local/share/Steam/appcache"),
				filepath.Join(home, ".local/share/Steam/config/htmlcache"),
				filepath.Join(home, flatpak, "com.valvesoftware.Steam/steam/steam/appcache"),
				filepath.Join(home, flatpak, "com.valvesoftware.Steam/.local/share/Steam/appcache"),
				filepath.Join(home, flatpak, "com.valvesoftware.Steam/.local/share/Steam/config/htmlcache"),
			}, false},
			{"Epic Games (Heroic/Lutris) Cache", []string{
				filepath.Join(home, ".config/heroic/WebCache"),
				filepath.Join(home, ".local/share/lutris/runtime"),
				filepath.Join(home, flatpak, "com.heroicgameslauncher.hgl/config/heroic/WebCache"),
				filepath.Join(home, flatpak, "com.heroicgameslauncher.hgl/.local/share/lutris/runtime"),
			}, false},
			{"Discord Cache", []string{
				filepath.Join(home, ".config/discord/Cache"),
				filepath.Join(home, ".config/discord/Code Cache"),
				filepath.Join(home, ".config/discord/GPUCache"),
				filepath.Join(home, flatpak, "com.discordapp.Discord/config/discord/Cache"),
				filepath.Join(home, flatpak, "com.discordapp.Discord/config/discord/Code Cache"),
				filepath.Join(home, flatpak, "com.discordapp.Discord/config/discord/GPUCache"),
			}, false},
			{"Telegram Cache", []string{filepath.Join(
				home, ".local/share/TelegramDesktop/tdata/user_data/cache"),
				filepath.Join(home, flatpak, "org.telegram.desktop/data/TelegramDesktop/tdata/user_data/cache"),
			}, false},
			{"Spotify Cache", []string{
				filepath.Join(home, cache, "spotify"),
				filepath.Join(home, flatpak, "com.spotify.Client/cache/spotify"),
			}, false},
			{"VS Code Cache", []string{
				filepath.Join(home, ".config/Code/Cache"),
				filepath.Join(home, ".config/Code/CachedData"),
				filepath.Join(home, ".config/Code/GPUCache"),
				filepath.Join(home, ".config/Code/User/workspaceStorage"),
				filepath.Join(home, flatpak, ".com.visualstudio.code/config/Code/Cache"),
				filepath.Join(home, flatpak, "com.visualstudio.code/config/Code/CachedData"),
				filepath.Join(home, flatpak, "com.visualstudio.code/config/Code/GPUCache"),
				filepath.Join(home, flatpak, "com.visualstudio.code/config/Code/User/workspaceStorage"),
			}, false},
			{"Mesa Shader Cache", []string{filepath.Join(home, cache, "mesa_shader_cache")}, false},
			{"Go Build Cache", []string{filepath.Join(home, cache, "go-build")}, false},
			{"Pip Cache", []string{filepath.Join(home, cache, "pip")}, false},
			{"NPM Cache", []string{filepath.Join(home, ".npm/_cacache")}, false},
			{"Yarn Cache", []string{filepath.Join(home, cache, "yarn")}, false},
			{"Cargo Cache", []string{filepath.Join(home, ".cargo/registry/cache")}, false},
		}
	}
}

// getDiskMetrics calculates total and free space for the root/system drive
func getDiskMetrics() (freeGB float64, totalStr string, freeStr string) {
	// Default values if something fails
	totalStr = "N/A"
	freeStr = "N/A"

	if GOOS == "windows" {
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

// ========================= MENU UI LOGIC =========================

func showBanner() {
	_, total, free := getDiskMetrics()
	fmt.Printf(`%s  ____________________     .-.
 |   |  |       __ |  \    |_|
 |   |  |      |  ||  |    | |
 |   |  |      |__||  |    |=|
 |   |__|__________|  |  .=/I\=.
 |                    | ////V\\\\
 |   ______________   | |#######|
 |  |______________|  | |||||||||
 |  |              |  |
 |  |              |  | %sCrunchyCleaner%s
 |  |              |  | Made by: Knuspii, (M)
 |[]|              |[]| Version: %s
 |__|______________|__| Disk-Space: %s / %s%s
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

// function to scan which programs actually exist on the disk
func scanForExisting() []Program {
	allPrograms := getPrograms()
	existing := []Program{}
	for _, p := range allPrograms {
		found := false
		var totalSize int64
		for _, path := range p.Paths {
			matches, _ := filepath.Glob(expandHome(path))
			if len(matches) > 0 {
				found = true
				totalSize += getDirSize(path)
			}
		}
		if found {
			// We format the name here so it's ready for both UI and Logs
			p.Name = fmt.Sprintf("%-30s %s(%s)%s", p.Name, YELLOW, formatMB(totalSize), RC)
			existing = append(existing, p)
		}
	}
	return existing
}

// handleMenu manages user input for navigation and selection
func handleMenu() {
	clearScreen()
	showBanner()

	// Initial scan of the filesystem to find existing directories
	stop := make(chan bool)
	ack := make(chan bool)
	go spinner("Scanning filesystem", stop, ack)
	time.Sleep(1 * time.Second)

	existing := scanForExisting()

	stop <- true // Tell spinner to stop
	<-ack        // WAIT for spinner to clear the line

	// Abort if nothing was detected
	if len(existing) == 0 {
		fmt.Printf("\nNo cache directories found on your system")
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

func runCleanup(programs []Program) {
	beforeFree, _, _ := getDiskMetrics()

	if *Flagdryrun {
		fmt.Printf("%sNOTE: Dry run active. No files will actually be deleted.%s", YELLOW, RC)
	} else {
		fmt.Printf("You use this tool at your own risk!")
	}
	// Fetching the current system user
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("\nUsername: unknown")
	} else {
		name := usr.Username
		// Strip domain/machine name prefix on Windows
		if GOOS == "windows" && strings.Contains(name, "\\") {
			parts := strings.Split(name, "\\")
			name = parts[len(parts)-1]
		}
		fmt.Printf("\nUsername: %s", name)
	}
	fmt.Printf("\nPress [CTRL+C] to cancel")
	fmt.Printf("\nCleaning caches started...\n")

	stop := make(chan bool)
	ack := make(chan bool)
	go spinner("Cleaning selected caches", stop, ack)
	time.Sleep(3 * time.Second)

	count := 0

	for _, p := range programs {
		if !p.Checked {
			continue
		}
		count++

		for _, path := range p.Paths {
			matches, _ := filepath.Glob(expandHome(path))

			for _, m := range matches {
				if *Flagdryrun {
					logInfo("Would clean: " + m)
					continue
				}
				deletePath(m)
			}
		}

		// Cut the size part
		name := p.Name
		if idx := strings.Index(name, "("); idx != -1 {
			name = strings.TrimSpace(name[:idx])
		}
		logOK(name)
	}

	stop <- true
	<-ack

	if count == 0 {
		fmt.Printf("\nNothing selected")
		return
	}

	if *Flagdryrun {
		logOK("Simulation finished")
	} else {
		logOK("Cleaning finished")
	}

	afterFree, _, _ := getDiskMetrics()
	cleaned := (afterFree - beforeFree) * 1024
	if cleaned < 0 || *Flagdryrun {
		cleaned = 0
	}

	line()
	fmt.Printf("CrunchyCleaner cleaned: %.2f MB\n", cleaned)

	if !*Flagauto {
		keyboard.Close()
		fmt.Printf("\nPress [ENTER] to exit")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
	cc_exit()
}

func deletePath(path string) {
	info, err := os.Stat(path)
	if err != nil {
		return
	}

	if !info.IsDir() {
		_ = os.Remove(path)
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		logWarn("Cannot read " + path + ": " + err.Error())
		return
	}

	for _, e := range entries {
		full := filepath.Join(path, e.Name())
		if err := os.RemoveAll(full); err != nil {
			logWarn("Skipped " + e.Name() + ": " + err.Error())
		}
	}
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

	// AUTOMATION LOGIC
	if *Flagauto {
		showBanner()
		fmt.Printf("%sNOTE: Automation active. Scanning and selecting all caches...%s\n", YELLOW, RC)
		existing := scanForExisting()

		if len(existing) == 0 {
			fmt.Printf("\nNo caches found. Nothing to do.")
			os.Exit(0)
		}

		// Check all found items
		for i := range existing {
			existing[i].Checked = true
		}
		runCleanup(existing)
	} else {
		// Run interactive mode
		handleMenu()
	}
}
