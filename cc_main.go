// ##################################################################
// CrunchyCleaner - System & Software Cache Cleaner
// Made by: Knuspii (M)
//
// LICENSE: CC BY-NC 4.0 (Creative Commons Attribution-NonCommercial)
// - You must attribute the author (link to GitHub).
// - Commercial use is strictly prohibited.
// ##################################################################

package main

import (
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
	CC_VERSION = "2.0"
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
	Flagdryrun  = flag.Bool("dry-run", false, "Simulation mode: identifies files without deleting them")
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
	fmt.Printf("\nExiting CrunchyCleaner. Goodbye!\n")
	os.Exit(0)
}

func pause() {
	fmt.Printf("\nPress [ENTER] to continue...")
	fmt.Scanln()
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
			fmt.Printf("\r%s%s%s %s%s %s%s%s", YELLOW, frames[i%len(frames)], RC, CYAN, text, YELLOW, frames[i%len(frames)], RC)
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
	clearScreen()
	fmt.Printf("%sInitializing CrunchyCleaner %s...%s\n", YELLOW, CC_VERSION, RC)

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
	} else if getcols != COLS || getlines != LINES {
		fmt.Printf("System: Terminal size mismatch (Got %dx%d, Expected %dx%d)\n", getcols, getlines, COLS, LINES)
	} else {
		fmt.Printf("System: Terminal size optimized (%dx%d)\n", getcols, getlines)
	}

	time.Sleep(1 * time.Second)
}

// getDiskMetrics calculates total and free space for the root/system drive
func getDiskMetrics() (float64, string, string) {
	var diskTotal, diskFree string
	var freeGB float64

	if goos == "windows" {
		// Query C: drive stats via PowerShell
		sizeOut, _ := exec.Command("powershell", "-Command",
			"(Get-PSDrive -PSProvider FileSystem | Where-Object {$_.Name -eq 'C'}).Used, (Get-PSDrive -PSProvider FileSystem | Where-Object {$_.Name -eq 'C'}).Free").Output()
		parts := strings.Fields(string(sizeOut))
		if len(parts) >= 2 {
			used, _ := strconv.ParseFloat(parts[0], 64)
			free, _ := strconv.ParseFloat(parts[1], 64)
			freeGB = free / 1024 / 1024 / 1024
			totalGB := (used + free) / 1024 / 1024 / 1024
			diskTotal = fmt.Sprintf("%.2f GB", totalGB)
			diskFree = fmt.Sprintf("%.2f GB", freeGB)
		}
	} else {
		// Use standard df command for Unix-like systems
		dfOut, _ := exec.Command("sh", "-c", "df -BG --output=size,avail / | tail -1 | tr -d 'G'").Output()
		parts := strings.Fields(string(dfOut))
		if len(parts) >= 2 {
			totalVal, _ := strconv.ParseFloat(parts[0], 64)
			freeVal, _ := strconv.ParseFloat(parts[1], 64)
			freeGB = freeVal
			diskTotal = fmt.Sprintf("%.2f GB", totalVal)
			diskFree = fmt.Sprintf("%.2f GB", freeVal)
		}
	}
	return freeGB, diskTotal, diskFree
}

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
 |    | __      |  |  | %sCrunchyCleaner - Clear Software Cache%s
 |    ||  |     |  |  | Made by: Knuspii, (M)
 |    ||__|     |  |  | Version: %s
 |____|_________|__|__| Disk-Space: %s / %s%s
`, YELLOW, RC, YELLOW, CC_VERSION, free, total, RC)
	line()
}

// expandHome resolves the shorthand '~/ ' to the absolute user home directory
func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

// getPrograms returns a curated list of cache locations for various popular software
func getPrograms() []Program {
	if goos == "windows" {
		appData := os.Getenv("APPDATA")
		localAppData := os.Getenv("LOCALAPPDATA")
		return []Program{
			{"Windows Thumbnails", []string{filepath.Join(localAppData, "Microsoft/Windows/Explorer")}, false},
			{"Firefox Cache", []string{
				filepath.Join(localAppData, "Mozilla/Firefox/Profiles/*/cache2"),
				filepath.Join(localAppData, "Mozilla/Firefox/Profiles/*/jumpListCache"),
				filepath.Join(appData, "Mozilla/Firefox/Profiles/*/shader-cache"),
			}, false},
			{"Chrome Cache", []string{
				filepath.Join(localAppData, "Google/Chrome/User Data/Default/Cache"),
				filepath.Join(localAppData, "Google/Chrome/User Data/ShaderCache"),
			}, false},
			{"Edge Cache", []string{
				filepath.Join(localAppData, "Microsoft/Edge/User Data/Default/Cache"),
				filepath.Join(localAppData, "Microsoft/Edge/User Data/ShaderCache"),
			}, false},
			{"Thunderbird Cache", []string{
				filepath.Join(localAppData, "Thunderbird/Profiles/*/cache2"),
				filepath.Join(localAppData, "Thunderbird/Profiles/*/startupCache"),
			}, false},
			{"Steam AppCache", []string{"C:/Program Files (x86)/Steam/appcache"}, false},
			{"Discord Cache", []string{
				filepath.Join(appData, "discord/Cache"),
				filepath.Join(appData, "discord/Code Cache"),
				filepath.Join(appData, "discord/GPUCache"),
			}, false},
			{"Spotify Storage", []string{filepath.Join(localAppData, "Spotify/Storage")}, false},
			{"VS Code Cache", []string{
				filepath.Join(appData, "Code/Cache"),
				filepath.Join(appData, "Code/CachedData"),
				filepath.Join(appData, "Code/User/workspaceStorage"),
			}, false},
			{"Pip Cache", []string{filepath.Join(localAppData, "pip/Cache")}, false},
			{"Go Build Cache", []string{filepath.Join(localAppData, "go-build")}, false},
			{"NPM Global Cache", []string{filepath.Join(localAppData, "npm-cache")}, false},
		}
	} else {
		// Linux/Unix paths
		home, _ := os.UserHomeDir()
		return []Program{
			{"Thumbnail Cache", []string{filepath.Join(home, ".cache/thumbnails")}, false},
			{"Firefox Cache", []string{filepath.Join(home, ".cache/mozilla/firefox/*/cache2")}, false},
			{"Chrome Cache", []string{filepath.Join(home, ".cache/google-chrome/Default/Cache")}, false},
			{"Edge Cache", []string{filepath.Join(home, ".cache/microsoft-edge/Default/Cache")}, false},
			{"Thunderbird Cache", []string{filepath.Join(home, ".cache/thunderbird/*/cache2")}, false},
			{"Spotify Storage", []string{filepath.Join(home, ".cache/spotify")}, false},
			{"Steam AppCache", []string{filepath.Join(home, ".steam/steam/appcache")}, false},
			{"Discord Cache", []string{filepath.Join(home, ".cache/discord")}, false},
			{"VS Code Cache", []string{filepath.Join(home, ".config/Code/Cache")}, false},
			{"Pip Cache", []string{filepath.Join(home, ".cache/pip")}, false},
			{"Go Build Cache", []string{filepath.Join(home, ".cache/go-build")}, false},
			{"NPM Cache", []string{filepath.Join(home, ".npm")}, false},
		}
	}
}

// --- Menu UI Logic ---

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
		fmt.Printf("Use ↑/↓ or W/S to navigate | SPACE to select | ENTER to clean\n")
		fmt.Printf("Software found: [%d]\n", len(existing))
	}

	for i := range existing {
		cursor := "    "
		if i == idx {
			cursor = YELLOW + "  >_" + RC
		}
		check := "[ ]"
		if existing[i].Checked {
			check = "[" + GREEN + "X" + RC + "]"
		}
		// Clear line and print entry
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

	allPrograms := getPrograms()
	existing := []Program{}
	for _, p := range allPrograms {
		found := false
		for _, path := range p.Paths {
			// Check if any file/folder matches the glob pattern
			matches, _ := filepath.Glob(expandHome(path))
			if len(matches) > 0 {
				found = true
				break
			}
		}
		if found {
			existing = append(existing, p)
		}
	}
	time.Sleep(1 * time.Second)
	done <- true

	// Clean up the scanning line
	fmt.Print("\033[A\033[K")

	if len(existing) == 0 {
		fmt.Printf("\nNo cache directories found on your system.\n")
		pause()
		return
	}

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
		} else if char == ' ' || key == keyboard.KeySpace {
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
		} else if key == keyboard.KeyEnter {
			clearScreen()
			showBanner()
			runCleanup(existing)
		} else if key == keyboard.KeyCtrlC {
			cc_exit()
		}

		if updated {
			// ANSI Escape: Move cursor up by N lines to redraw menu in-place
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
			matches, _ := filepath.Glob(expandHome(path))
			for _, m := range matches {
				if *Flagdryrun {
					// Only log the target in simulation mode
					logInfo(fmt.Sprintf("[SIMULATE] Would delete: %s", m))
				} else {
					// Physical deletion
					err := os.RemoveAll(m)
					if err != nil {
						logWarn("Error cleaning " + p.Name + ": " + err.Error())
					}
				}
			}
		}

		if !*Flagdryrun {
			logOK("Cleaned " + p.Name)
		}
	}
	done <- true

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
	totalCleanedMB := (afterFree - beforeFree) * 1024

	if totalCleanedMB < 0 || *Flagdryrun {
		totalCleanedMB = 0
	}

	line()
	label := "Cleaned"
	if *Flagdryrun {
		label = "Space to be recovered"
	}

	fmt.Printf(" %s: %s%.2f MB%s\n", label, GREEN, totalCleanedMB, RC)
	line()

	fmt.Printf("\nPress [ENTER] to exit...")
	fmt.Scanln()
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
		cc_exit()
	}

	if !*Flagnoinit {
		initApp()
	}

	handleMenu()
}
