// ##################################
// CrunchyCleaner
// Made by: Knuspii (M)
// ##################################

package main

import (
	"bufio"     // For reading user input
	"fmt"       // Formatted input/output
	"os"        // General OS interactions (exit, files, etc.)
	"os/exec"   // Executing external commands
	"os/signal" // Handling file paths in a cross-platform way
	"runtime"   // Info about the OS / architecture
	"strings"   // String manipulation (Trim, Split, Join, etc.)
	"syscall"   // System calls (for signal handling)
	"time"      // Time-related functions (sleep, timestamp, timeout)

	"github.com/eiannone/keyboard" // For capturing keyboard input (like key presses)
)

const (
	CC_VERSION = "1.6"
	COLS       = 62
	LINES      = 30
	CMDWAIT    = 1 * time.Second // Wait time running a command
	RED        = "\033[31m"
	YELLOW     = "\033[33m"
	GREEN      = "\033[32m"
	BLUE       = "\033[34m"
	CYAN       = "\033[36m"
	RC         = "\033[0m" // Reset ANSI color
)

var (
	origCols, origLines int            // Original terminal size
	verbose             = false        // If true, print all errors
	skipPause           = false        // If true, skip pause
	goos                = runtime.GOOS // Current OS
	reader              = bufio.NewReader(os.Stdin)
	//SPINNERFRAMES  = []rune{'⣾', '⣽', '⣻', '⢿', '⡿', '⣟', '⣯', '⣷'} // Spinner animation frames
	SPINNERFRAMES = []rune{'|', '/', '-', '\\'} // Spinner animation frames
)

func adminCheck() {
	switch goos {
	// WINDOWS
	case "windows":
		// Try running a command that requires admin rights
		cmd := exec.Command("net", "session")
		if err := cmd.Run(); err != nil {
			printError("You need admin privileges")
			os.Exit(1)
		}
	default:
		if os.Geteuid() != 0 { // Check if current user is not root
			printError("You need root privileges")
			os.Exit(1)
		}
	}
}

func getAdmin() {
	switch goos {
	// WINDOWS
	case "windows":
		// Try running a command that requires admin privileges
		cmd := exec.Command("net", "session")
		if err := cmd.Run(); err != nil {
			// If it fails, try to restart the program as admin
			printInfo("Restarting as admin...\n")

			elevate := exec.Command("powershell", "-Command", "Start-Process", os.Args[0], "-Verb", "RunAs")

			// Connect the elevated process output to the current console
			elevate.Stdout = os.Stdout
			elevate.Stderr = os.Stderr

			if err := elevate.Run(); err != nil {
				// Failed to elevate privileges
				printError(fmt.Sprintf("Failed to restart as admin: %v", err))
				fmt.Println("CrunchyCleaner might not work correctly without admin rights")
				pause()
			} else {
				// Successfully elevated, exit current process
				os.Exit(0)
			}
		}
	// LINUX
	default:
		if os.Geteuid() != 0 { // Check if current user is not root
			printInfo("Requesting root privileges...\n")

			// Get absolute path to the running executable
			scriptPath, err := exec.LookPath(os.Args[0])
			if err != nil {
				printError(fmt.Sprintf("Could not find executable: %v", err))
				return
			}

			// Prepare command with all original arguments
			args := append([]string{scriptPath}, os.Args[1:]...)
			cmd := exec.Command("sudo", args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin // Needed to read sudo password

			if err := cmd.Run(); err != nil {
				// User canceled or sudo failed
				if strings.Contains(err.Error(), "interrupt") {
					fmt.Println("\nCancelled by user")
					os.Exit(0)
				}
				printError(fmt.Sprintf("Failed to restart with sudo: %v", err))
				fmt.Print("Press [ENTER] to continue anyway: ")
				reader.ReadString('\n')
			} else {
				// Successfully elevated, exit current process
				os.Exit(0)
			}
		}
	}
}

func showInfo() {
	fmt.Printf(`%sCrunchyCleaner Version: %s%s
https://github.com/Knuspii/crunchycleaner

DISCLAIMER:
MADE BY: Knuspii, (M)
Help by the World Wide Web.

A lightweight, cross-platform system cleanup tool.
You use this tool at your own risk.
I do not take any responsibilities.

System Requirements:
- Windows or Linux
- Terminal with ANSI escape code support

This work is licensed under the:
Creative Commons Attribution-
NonCommercial 4.0 International License.

This work uses the following external dependencies:
- github.com/eiannone/keyboard (for cross-platform keyboard input)
`, YELLOW, CC_VERSION, RC)
}

func getTermSize() (cols, lines int, err error) {
	switch goos {
	case "windows":
		// Windows: get with PowerShell
		cmd := exec.Command("powershell", "-Command",
			`$size = $Host.UI.RawUI.WindowSize; Write-Output "$($size.Width) $($size.Height)"`)
		out, err := cmd.Output()
		if err != nil {
			return 0, 0, err
		}
		fmt.Sscanf(string(out), "%d %d", &cols, &lines)
	default:
		// Unix: use stty
		cmd := exec.Command("stty", "size")
		cmd.Stdin = os.Stdin
		out, err := cmd.Output()
		if err != nil {
			return 0, 0, err
		}
		fmt.Sscanf(string(out), "%d %d", &lines, &cols)
	}
	return
}

func setTermSize(cols, lines int) {
	switch goos {
	case "windows":
		// Build PowerShell command to resize terminal window and buffer
		psCmd := fmt.Sprintf(
			`$Host.UI.RawUI.WindowSize = New-Object System.Management.Automation.Host.Size(%d, %d); $Host.UI.RawUI.BufferSize = New-Object System.Management.Automation.Host.Size(%d, 300)`,
			cols, lines, cols,
		)
		_ = exec.Command("powershell", "-Command", psCmd).Run() // execute the command silently
	default:
		// Send ANSI escape code to resize terminal on Linux/macOS
		fmt.Printf("\033[8;%d;%dt", lines, cols)
	}
}

// It saves the original terminal size and sets it to program defaults.
func init_term() {
	// Save original terminal size
	cols, lines, err := getTermSize()
	if err == nil {
		origCols, origLines = cols, lines
	}

	// Resize terminal to program defaults
	setTermSize(COLS, LINES)
}

// Restores the terminal to its original size if it was saved.
func restoreTerm() {
	if origCols > 0 && origLines > 0 {
		setTermSize(origCols, origLines)
	}
}

func showBanner() {
	total, free := getDiskInfo()
	fmt.Printf(`%s
  ____________________     .-.
 |  |              |  |    |_|
 |[]|              |[]|    | |
 |  |              |  |    |=|
 |  |              |  |  .=/I\=.
 |  |              |  | ////V\\\\
 |  |______________|  | |#######|
 |                    | |||||||||
 |     ____________   |
 |    | __      |  |  | %sCrunchyCleaner - Cleanup your system!%s
 |    ||  |     |  |  | Made by: Knuspii, (M)
 |    ||__|     |  |  | Version: %s
 |____|_________|__|__| Disk-Space: %s / %s%s
`, YELLOW, RC, YELLOW, CC_VERSION, free, total, RC)
	line()
}

func usage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("  %scrunchycleaner [option]%s\n\n", CYAN, RC)
	fmt.Printf("Options:\n")
	fmt.Printf("  %sNo option%s   Run with TUI (Text-UI)\n", YELLOW, RC)
	fmt.Printf("  %s-t%s          Run with TUI (Text-UI)\n", YELLOW, RC)
	fmt.Printf("  %s-s%s          Run Safe-Cleanup\n", YELLOW, RC)
	fmt.Printf("  %s-sy%s         Run Safe-Cleanup (non-interactive for scripts)\n", YELLOW, RC)
	fmt.Printf("  %s-f%s          Run Full-Cleanup\n", YELLOW, RC)
	fmt.Printf("  %s-fy%s         Run Full-Cleanup (non-interactive for scripts)\n", YELLOW, RC)
	fmt.Printf("  %s-u <user> %s  Run User-Cleanup\n", YELLOW, RC)
	fmt.Printf("  %s-uy <user>%s  Run User-Cleanup (non-interactive for scripts)\n", YELLOW, RC)
	fmt.Printf("  %s-v%s          Show version\n", YELLOW, RC)
	fmt.Printf("  %s-h%s          Show this help page\n", YELLOW, RC)
}

func normalstartup() {
	adminCheck()
	showBanner()
}

func skipstartup() {
	verbose = true
	skipPause = true
	adminCheck()
	showBanner()
}

func handleargs() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		// Safe-Cleanup
		case "-s":
			normalstartup()
			cleanup("safe")
			os.Exit(0)
		// Safe-Cleanup (non-interactive)
		case "-sy":
			skipstartup()
			cleanup("safe")
			os.Exit(0)
		// Full-Cleanup
		case "-f":
			normalstartup()
			cleanup("full")
			os.Exit(0)
		// Full-Cleanup (non-interactive)
		case "-fy":
			skipstartup()
			cleanup("full")
			os.Exit(0)
		// User-Cleanup
		case "-u":
			adminCheck()
			if len(os.Args) > 2 {
				showBanner()
				cleanup("user", os.Args[2])
			} else {
				printError("No profile name provided")
				os.Exit(1)
			}
			os.Exit(0)
		// User-Cleanup (non-interactive)
		case "-uy":
			skipPause = true
			verbose = true
			adminCheck()
			if len(os.Args) > 2 {
				showBanner()
				cleanup("user", os.Args[2])
			} else {
				printError("No profile name provided")
				os.Exit(1)
			}
			os.Exit(0)
		// Help
		case "-h", "--help":
			showBanner()
			usage()
			os.Exit(0)
		// Version
		case "-v", "--version":
			fmt.Printf("CrunchyCleaner %s\n", CC_VERSION)
			os.Exit(0)
		// TUI
		case "-t":
			// Just continue to TUI
		default:
			fmt.Printf("Unknown option: %s\n", os.Args[1])
			usage()
			os.Exit(1)
		}
	}
}

// MenuItem defines a single menu entry.
// "Name" is the internal identifier used in switch-case logic.
// "Text" is what’s actually shown to the user in the menu.
type MenuItem struct {
	Name string // internal name for switch-case
	Text string // visible text in the menu
}

// All available menu items for the TUI.
var menuItems = []MenuItem{
	{"Full-Clean", "Full-Clean  - Does a Full-Cleanup"},
	{"Safe-Clean", "Safe-Clean  - Does a Safe-Cleanup"},
	{"User-Clean", "User-Clean  - Does a User-Cleanup"},
	{"Info", "Info        - Shows some infos"},
	{"Reset", "Reset       - Reset the TUI"},
	{"Exit", "Exit        - Exit"},
}

// handleMenu controls the interactive TUI (Text User Interface) menu.
// It allows navigation with arrow keys or W/S, and executes actions on ENTER.
func handleMenu() {
	idx := 0 // current selected menu index

	// Open keyboard input in raw mode (captures single key presses)
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close() // ensure cleanup when function exits

	for {
		// Move cursor to line 16 and clear everything below that point
		// (keeps the banner visible at the top)
		fmt.Print("\033[16;0H")
		fmt.Print("\033[J")

		// Draw all menu items
		for i, item := range menuItems {
			if i == idx {
				// Highlight the currently selected item
				fmt.Printf(" %s>%s%s%s\n", CYAN, YELLOW, item.Text, RC)
			} else {
				fmt.Printf("  %s\n", item.Text)
			}
		}

		line() // print separator line
		fmt.Printf("Use ↑/↓ or W/S to navigate, ENTER to select\n")

		// Wait for a key press
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		// Handle navigation and selection
		switch {
		case key == keyboard.KeyArrowUp || char == 'w' || char == 'W':
			// Move selection up
			if idx > 0 {
				idx--
			}

		case key == keyboard.KeyArrowDown || char == 's' || char == 'S':
			// Move selection down
			if idx < len(menuItems)-1 {
				idx++
			}

		case key == keyboard.KeyCtrlC:
			// Exit on CTRL+C
			cc_exit()

		case key == keyboard.KeyEnter:
			// Perform action based on selected menu item
			switch menuItems[idx].Name {

			case "Full-Clean":
				printInfo("Selected: Full-Cleanup")
				cmdline()
				cleanup("full")
				pause()
				cc_exit()

			case "Safe-Clean":
				printInfo("Selected: Safe-Cleanup")
				cmdline()
				cleanup("safe")
				pause()
				cc_exit()

			case "User-Clean":
				printInfo("Selected: User-Cleanup")
				cmdline()
				cleanup("user")
				pause()
				cc_exit()

			case "Info":
				printInfo("Selected: Info")
				cmdline()
				showInfo()
				pause()
				cc_exit()

			case "Reset":
				// Clear screen and reinitialize the TUI
				clearScreen()
				init_term()
				showBanner()

			case "Exit":
				// Restore terminal and exit gracefully
				cc_exit()
			}
		}
	}
}

// main is the entry point of CrunchyCleaner.
// It sets up arguments, handles signals (Ctrl+C), initializes the terminal,
// and starts the main menu.
func main() {
	handleargs()

	// Catch [Ctrl+C] to restore the terminal before exiting
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cc_exit()
	}()

	// Request admin/root privileges if needed
	getAdmin()

	// Initialize terminal (resize, save original size, etc.)
	init_term()
	clearScreen()

	// Show program banner
	showBanner()

	// Start the interactive text menu
	handleMenu()
}
