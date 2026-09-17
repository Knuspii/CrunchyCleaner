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
	"flag"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

// Global constants for UI and Versioning
const (
	CC_VERSION = "2.9"
	COLS       = 56
	LINES      = 32
	GOOS       = runtime.GOOS
	CLEARLINE  = "\r\033[K"
	YELLOW     = "\033[33m"
	CYAN       = "\033[36m"
	GREEN      = "\033[32m"
	RC         = "\033[0m"
)

var (
	// CLI Flags
	Flagversion = flag.Bool("v", false, "Display version information")
	Flagdryrun  = flag.Bool("d", false, "Dry-run mode without deleting files (for testing)")
	Flagauto    = flag.Bool("a", false, "Automate cleaning (select all and start immediately)")
)

var menuHelpLines = [...]string{
	"↑/↓ or W/S to navigate | ENTER to select | C to clean",
}

// Program represents a target application and its associated cache directories
type Program struct {
	Name    string
	Paths   []string // List of paths (supports wildcards/globbing)
	Checked bool     // Selection state in the menu
}

type CacheEntry struct {
	Program
	Size int64
}

type KeyEvent struct {
	Char rune
	Key  keyboard.Key
	Err  error
}

var keyEvents = make(chan KeyEvent, 100)

// ========================= HELPER FUNCTIONS =========================

// Initializes the keyboard and starts a background listener
func startKeyboardListener() {
	if err := keyboard.Open(); err != nil {
		fmt.Printf("Error initializing keyboard: %v\n", err)
		os.Exit(1)
	}

	go func() {
		for {
			char, key, err := keyboard.GetKey()
			if err != nil {
				continue
			}
			// Immediate, global exit on 'q', 'Q' or Ctrl+C
			if char == 'q' || char == 'Q' || key == keyboard.KeyCtrlC {
				cc_exit()
			}
			// Forward other keys to the program for processing
			keyEvents <- KeyEvent{Char: char, Key: key, Err: err}
		}
	}()
}

// flushKeyEvents clears any pending key events in the channel buffer (just for safety)
func flushKeyEvents() {
	for {
		select {
		case <-keyEvents:
		default:
			return
		}
	}
}

// cc_exit provides a clean termination of the application
func cc_exit() {
	// Close keyboard
	keyboard.Close()

	// Enable cursor
	fmt.Print("\033[?25h")

	fmt.Printf("\nExiting CrunchyCleaner...\n")
	os.Exit(0)
}

func pause() {
	flushKeyEvents()
	fmt.Printf("\nPress [ENTER] to continue...")
	for ev := range keyEvents {
		if ev.Key == keyboard.KeyEnter {
			break
		}
	}
}

// line draws a formatted horizontal separator
func line() {
	fmt.Printf("%s#%s~%s\n", YELLOW, strings.Repeat("=", COLS-2), RC)
}

// spinner visualizes background tasks and cleans up properly
func spinner(text string, stop chan bool, ack chan bool) {
	frames := []string{"|", "/", "-", "\\"}
	i := 0
	for {
		select {
		case <-stop:
			// Clear the line and move cursor to start
			fmt.Print(CLEARLINE)
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
 |__|______________|__| Disk: %s / %s%s
`, YELLOW, RC, YELLOW, CC_VERSION, free, total, RC)
	line()
}

func logInfo(msg string) { fmt.Printf("%s[+] %s%s\n", CYAN, msg, RC) }
func logOK(msg string)   { fmt.Printf("%s[✓] %s%s\n", GREEN, msg, RC) }
func logWarn(msg string) { fmt.Printf("%s[!] %s%s\n", YELLOW, msg, RC) }

// renderMenu draws the interactive selection list
func renderMenu(existing []CacheEntry, idx int, fullRedraw bool) {
	if fullRedraw {
		showBanner()
	}

	fmt.Printf("%s%s | Folders found: %d\n", CLEARLINE, selectionSummary(existing), len(existing))
	for _, helpLine := range menuHelpLines {
		fmt.Printf("%s%s\n", CLEARLINE, helpLine)
	}

	// Render each detected program entry
	for i := range existing {
		cursor := "   "
		// Highlight the currently selected entry
		if i == idx {
			cursor = YELLOW + " >_" + RC
		}
		// Checkbox indicator for selection state
		check := "[ ]"
		if existing[i].Checked {
			check = "[" + GREEN + "X" + RC + "]"
		}
		// Clear the current line and print the menu entry
		fmt.Printf("%s%s%s %-30s %s(%s)%s\n", CLEARLINE, cursor, check, existing[i].Name, YELLOW, formatBytes(existing[i].Size), RC)
	}
}

// handleMenu manages user input for navigation and selection
func handleMenu() {
	// Channels used to synchronize and control the lifecycle of the background spinner goroutine
	stop := make(chan bool)
	ack := make(chan bool)

	// Spin up the visual loader in a separate thread so it doesn't block filesystem scanning
	go spinner("Starting CrunchyCleaner & Scanning filesystem", stop, ack)
	time.Sleep(1 * time.Second) // Give the user a brief moment to see the loading spinner

	// Perform the actual scan on the local machine to find targeted cache paths
	existing := scanForExisting()

	// Terminate the background spinner and wait for its clean shutdown signal
	stop <- true
	<-ack

	idx := 0
	// Perform the initial, full redraw of the menu screen
	renderMenu(existing, idx, true)

	// If no matching cache directories are found on the machine, display a warning and exit
	if len(existing) <= 0 {
		fmt.Printf("\nNo cache directories found on your system...")
		pause()
		cc_exit()
	}

	// Main Input Loop
	// This loops infinitely, waiting for and processing incoming key events
	for {
		// Read the next key event from the global background listener channel
		ev := <-keyEvents
		char, key, err := ev.Char, ev.Key, ev.Err
		if err != nil {
			break
		}

		// Track whether the menu selection or checkbox state has changed
		updated := false

		// Navigate selection upwards (using Arrow Up or 'W'/'w')
		if key == keyboard.KeyArrowUp || char == 'w' || char == 'W' {
			if idx > 0 {
				idx--
				updated = true
			}
			// Navigate selection downwards (using Arrow Down or 'S'/'s')
		} else if key == keyboard.KeyArrowDown || char == 's' || char == 'S' {
			if idx < len(existing)-1 {
				idx++
				updated = true
			}
			// Toggle selection state of the highlighted list item (Spacebar or Enter)
		} else if char == ' ' || key == keyboard.KeyEnter || key == keyboard.KeySpace {
			existing[idx].Checked = !existing[idx].Checked
			updated = true
			// Toggle 'Select All' / 'Deselect All' logic when pressing 'A'/'a'
		} else if char == 'a' || char == 'A' {
			toggleAll(existing)
			updated = true
			// Trigger the cleanup sequence for all checked items (using 'C'/'c')
		} else if char == 'c' || char == 'C' {
			runCleanup(existing)
		}

		// Redraw the menu dynamically if the UI state changed
		if updated {
			// Move terminal cursor back up to the start of the menu using ANSI escape codes.
			// This prevents screen flickering by avoiding a complete terminal screen clear.
			fmt.Printf("\033[%dA", len(existing)+1+len(menuHelpLines))
			renderMenu(existing, idx, false)
		}
	}
}

func runCleanup(programs []CacheEntry) {
	if *Flagdryrun {
		fmt.Printf("\n%sNOTE: Dry run active. No files will actually be deleted.%s", YELLOW, RC)
	} else {
		fmt.Printf("\nYou use this tool at your own risk!")
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
	fmt.Printf("\nPress Q to cancel")
	fmt.Printf("\nCleaning caches started...\n")

	stop := make(chan bool)
	ack := make(chan bool)
	go spinner("Cleaning selected caches", stop, ack)
	time.Sleep(3 * time.Second)

	count := 0
	var cleaned int64

	for _, p := range programs {
		if !p.Checked {
			continue
		}
		count++

		for _, path := range p.Paths {
			matches := globMatches(path)

			for _, m := range matches {
				if *Flagdryrun {
					fmt.Print(CLEARLINE)
					logInfo("Would clean: " + m)
					continue
				}
				cleaned += deletePath(m)
			}
		}

		fmt.Print(CLEARLINE)
		logOK(p.Name)
	}

	stop <- true
	<-ack

	if count == 0 {
		fmt.Printf("\nNothing selected")
		time.Sleep(3 * time.Second)
		cc_exit()
	}

	if *Flagdryrun {
		logOK("Dry-Run finished")
	} else {
		logOK("Cleaning finished")
	}

	line()
	if *Flagdryrun {
		fmt.Printf("CrunchyCleaner cleaned: NOTHING (DRY-RUN)\n")
	} else {
		fmt.Printf("CrunchyCleaner cleaned: %s%s%s\n", YELLOW, formatBytes(cleaned), RC)
	}

	if !*Flagauto {
		for {
			fmt.Printf("\nPress Q to exit")
			ev := <-keyEvents
			char, err := ev.Key, ev.Err
			if err != nil {
				break
			}

			if char == 'q' || char == 'Q' {
				continue
			}
		}
	}
	cc_exit()
}

func main() {
	flag.Parse()

	if *Flagversion {
		fmt.Printf("CrunchyCleaner %s\n", CC_VERSION)
		return
	}

	startKeyboardListener()

	// AUTOMATION LOGIC
	if *Flagauto {
		showBanner()
		fmt.Printf("%sNOTE: Automation active. Scanning and selecting all caches...%s\n", YELLOW, RC)
		existing := scanForExisting()

		// Check all found items
		for i := range existing {
			existing[i].Checked = true
		}
		runCleanup(existing)
	}

	// Run interactive mode
	handleMenu()
}
