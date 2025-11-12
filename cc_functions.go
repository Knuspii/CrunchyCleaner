// CrunchyCleaner: Helper Functions

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

func printInfo(msg string) {
	fmt.Printf("%s[INFO]%s %s\n", YELLOW, RC, msg)
}
func printError(msg string) {
	fmt.Printf("%s[ERROR]%s %s\n", RED, RC, msg)
}
func printSuccess(msg string) {
	fmt.Printf("\n%s[INFO]%s %s\n", GREEN, RC, msg)
}
func printTask(msg string) {
	fmt.Printf("%s[*]%s %s\n", GREEN, RC, msg)
}

func line() {
	fmt.Printf("%s#%s~%s\n", YELLOW, strings.Repeat("-", COLS-2), RC)
}
func cmdline() {
	fmt.Printf("%s#%s~%s\n", RED, strings.Repeat("-", COLS-2), RC)
}

func cc_exit() {
	defer keyboard.Close()
	restoreTerm()
	printTask("CrunchyCleaner EXITED")
	os.Exit(0)
}

// pause waits for the user to press ENTER before continuing.
// Uses the "keyboard" library for non-blocking key detection.
func pause() {
	fmt.Printf("\nPress [ENTER] to continue...")

	// Open keyboard input in raw mode (captures individual keypresses)
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close() // make sure the keyboard gets properly closed after use

	for {
		// Read a single key press (character + special key)
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		// If ENTER is pressed, break out of the loop
		if key == keyboard.KeyEnter {
			break
		}
		// If CTRL+C is pressed, exit the program
		if key == keyboard.KeyCtrlC {
			cc_exit()
		}
	}

	// Move the cursor up one line
	fmt.Printf("\033[1A")

	// Clear the current line
	fmt.Printf("\r\033[2K")
}

// asyncSpinner displays a spinning "loading" animation in the terminal.
// It runs asynchronously and stops when the provided context is canceled.
func asyncSpinner(ctx context.Context, text string) {
	i := 0 // Index for spinner frames
	for {
		select {
		// If the context is canceled, stop the spinner and return
		case <-ctx.Done():
			return

		// Default case: continue spinning
		default:
			// Print spinner line:
			// \r        -> Carriage return to overwrite the same line
			// [LOADING] -> Static label
			// YELLOW/RC -> Apply color and reset
			// text      -> Custom text passed to the spinner
			// SPINNERFRAMES[i%len(SPINNERFRAMES)] -> Rotate through spinner characters
			fmt.Printf("\r%s[LOADING]%s %s %s%c%s  ", YELLOW, RC, text, YELLOW, SPINNERFRAMES[i%len(SPINNERFRAMES)], RC)
			time.Sleep(100 * time.Millisecond) // Wait a short time before next frame
			i++                                // Move to the next spinner frame
		}
	}
}

// runCommand executes an external command and returns its combined output (stdout + stderr).
// It takes a slice of strings, where the first element is the command and the rest are arguments.
func runCommand(cmd []string) (string, error) {
	// Check if the command slice is empty
	if len(cmd) == 0 {
		return "", errors.New("command is empty")
	}

	// Create an exec.Command object with the command and its arguments
	c := exec.Command(cmd[0], cmd[1:]...)

	// Run the command and capture both stdout and stderr
	outBytes, err := c.CombinedOutput()

	// Convert output bytes to string and trim whitespace/newlines
	out := strings.TrimSpace(string(outBytes))

	// If the command failed, return the output along with a formatted error
	if err != nil {
		return out, fmt.Errorf("command failed: %v", err)
	}

	// If successful, return the output and nil error
	return out, nil
}

// clearScreen clears the terminal screen for neat output
func clearScreen() {
	var cmd *exec.Cmd
	if goos == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Fallback: ANSI for Clear + Reset Cursor
		fmt.Print("\033[H\033[2J")
	}
}

// getDiskInfo returns total and free disk space as human-readable strings
func getDiskInfo() (string, string) {
	var diskTotal, diskFree string

	if goos == "windows" {
		// Windows: PowerShell query
		sizeOut, _ := exec.Command("powershell", "-Command",
			"(Get-PSDrive -PSProvider FileSystem | Where-Object {$_.Name -eq 'C'}).Used, (Get-PSDrive -PSProvider FileSystem | Where-Object {$_.Name -eq 'C'}).Free").Output()
		parts := strings.Fields(string(sizeOut))
		if len(parts) >= 2 {
			used, _ := strconv.ParseFloat(parts[0], 64)
			free, _ := strconv.ParseFloat(parts[1], 64)
			total := used + free
			diskTotal = fmt.Sprintf("%.0f MB", total/1024/1024)
			diskFree = fmt.Sprintf("%.0f MB", free/1024/1024)
		}
	} else {
		// Linux / Unix: df
		dfOut, _ := exec.Command("sh", "-c", "df -BM --output=size,avail / | tail -1 | tr -d 'M'").Output()
		parts := strings.Fields(string(dfOut))
		if len(parts) >= 2 {
			diskTotal = parts[0] + " MB"
			diskFree = parts[1] + " MB"
		}
	}

	return diskTotal, diskFree
}

// getFreeMB returns free disk space in MB (numeric, for diff calculation)
func getFreeMB() int64 {
	if goos == "windows" {
		out, _ := exec.Command("powershell", "-Command",
			"(Get-PSDrive -PSProvider FileSystem | Where-Object {$_.Name -eq 'C'}).Free").Output()
		free, _ := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)
		return free / 1024 / 1024
	}

	out, _ := exec.Command("sh", "-c", "df --output=avail / | tail -1").Output()
	free, _ := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)
	return free / 1024
}
