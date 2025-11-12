// CrunchyCleaner: Core Cleaner Functions

package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

// task defines a cleanup action
// Either a Go-native function (goFunc) or an external command (cmd)
type task struct {
	desc   string       // Short description of the task
	cmd    []string     // Command-line slice to execute
	goFunc func() error // Go function to execute
	path   string       // Folder to clean
}

func cleanFolder(desc, folder string) task {
	return task{
		desc: desc,
		path: folder,
		goFunc: func() error {
			// read all entries (files and subdirectories) inside the folder
			entries, err := os.ReadDir(folder)
			if err != nil {
				// return error if reading the folder fails
				return fmt.Errorf("cleanup failed: %v", err)
			}

			// iterate over each entry and remove it
			for _, entry := range entries {
				// construct full path of the entry
				path := folder + string(os.PathSeparator) + entry.Name()

				// remove the entry and its contents if it is a directory
				os.RemoveAll(path)
			}

			// success: folder contents deleted, folder itself remains
			return nil
		},
	}
}

// cleanup orchestrates all cleanup tasks based on the mode (user/full/etc)
// and optionally a username
func cleanup(mode string, username ...string) {
	printInfo(fmt.Sprintf("Cleanup: (%s mode on %s)", mode, goos))
	askVerbose()
	var tasks []task
	switch mode {
	case "user":
		profile, ok := selectProfile(username)
		if !ok {
			return
		}
		tasks = buildUserTasks(goos, profile)
	default:
		tasks = buildSystemTasks(goos, mode)
	}
	previewTasks(tasks)
	runTasks(tasks)
}

// ==================== PROFILE HANDLING ==================== //

// selectProfile returns the chosen user profile (via CLI arg or interactive input)
func selectProfile(username []string) (string, bool) {
	profiles := getProfiles()
	if len(profiles) == 0 {
		printError("No user profiles found!")
		if !skipPause {
			pause()
		}
		return "", false
	}

	var selectedProfile string
	if len(username) > 0 {
		selectedProfile = username[0]
		if !contains(profiles, selectedProfile) {
			printError("No profile name provided")
			if !skipPause {
				pause()
			}
			os.Exit(1)
		}
	} else {
		fmt.Println("\n\nAvailable profiles:")
		for _, p := range profiles {
			fmt.Printf("  %s\n", p)
		}
		fmt.Printf("Input profile name to clean: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		if !contains(profiles, choice) {
			printError("Invalid profile name provided")
			if !skipPause {
				pause()
			}
			os.Exit(1)
		}
		selectedProfile = choice
	}

	fmt.Printf("%sSelected profile: %s%s\n", YELLOW, selectedProfile, RC)
	return selectedProfile, true
}

// getProfiles fetches user profiles depending on OS
func getProfiles() []string {
	var profiles []string
	var files []os.DirEntry
	var err error

	switch goos {
	case "windows":
		userBase := os.Getenv("SystemDrive") + `\Users`
		files, err = os.ReadDir(userBase)
	default:
		files, err = os.ReadDir("/home")
	}

	if err != nil {
		printError("Reading profiles folder: ")
		fmt.Printf("%s\n", err)
		return nil
	}

	for _, f := range files {
		if f.IsDir() {
			profiles = append(profiles, f.Name())
		}
	}
	return profiles
}

// contains helper: check if slice contains string
func contains(list []string, val string) bool {
	for _, item := range list {
		if item == val {
			return true
		}
	}
	return false
}

// ==================== TASK BUILDERS ==================== //

// buildUserTasks builds cleanup tasks for a given user profile
func buildUserTasks(goos, profile string) []task {
	switch goos {
	case "windows":
		userPath := os.Getenv("SystemDrive") + `\Users\` + profile
		return []task{
			cleanFolder("Windows Explorer Cache", userPath+`\AppData\Local\Microsoft\Windows\Explorer`),
			cleanFolder("Local Crash Dumps", userPath+`\AppData\Local\CrashDumps`),
			cleanFolder("Temp Folder", userPath+`\AppData\Local\Temp`),
		}
	default:
		userPath := "/home/" + profile
		return []task{
			cleanFolder("Cache Folder", userPath+"/.cache"),
			cleanFolder("Thumbnails", userPath+"/.thumbnails"),
		}
	}
}

// buildSystemTasks builds cleanup tasks depending on OS and mode
func buildSystemTasks(goos, mode string) []task {
	switch goos {
	case "windows":
		return buildWindowsTasks(mode)
	default:
		return buildLinuxTasks(mode)
	}
}

// ==================== WINDOWS TASKS ==================== //
func buildWindowsTasks(mode string) []task {
	systemRoot := os.Getenv("SystemRoot")
	programData := os.Getenv("ProgramData")

	if mode == "full" {
		return []task{
			cleanFolder("Windows Event Logs", systemRoot+`\System32\winevt\Logs`),
			cleanFolder("Windows Update Logs", programData+`\Microsoft\Windows\WindowsUpdate\Logs`),
			cleanFolder("Windows Logs", systemRoot+`\Logs`),
			cleanFolder("Temp Folder", systemRoot+`\Temp`),
			cleanFolder("Prefetch Folder", systemRoot+`\Prefetch`),
			cleanFolder("Windows WER", programData+`\Microsoft\Windows\WER`),
			cleanFolder("WDI Log Files", systemRoot+`\System32\WDI\LogFiles`),
			cleanFolder("Defender CacheManager", programData+`\Microsoft\Windows Defender\Scans\History\CacheManager`),
			cleanFolder("CBS Logs", systemRoot+`\Logs\CBS`),
			cleanFolder("Delivery Optimization", systemRoot+`\SoftwareDistribution\DeliveryOptimization`),
			cleanFolder("Windows.old", `C:\Windows.old`),
			cleanFolder("SoftwareDistribution Download", systemRoot+`\SoftwareDistribution\Download`),
			cleanFolder("Driver FileRepository", systemRoot+`\System32\DriverStore\FileRepository`),
		}
	}
	return []task{
		cleanFolder("Temp Folder", systemRoot+`\Temp`),
		cleanFolder("Windows Update Logs", programData+`\Microsoft\Windows\WindowsUpdate\Logs`),
		cleanFolder("Defender CacheManager", programData+`\Microsoft\Windows Defender\Scans\History\CacheManager`),
		cleanFolder("Delivery Optimization", systemRoot+`\SoftwareDistribution\DeliveryOptimization`),
	}
}

// ==================== LINUX TASKS ==================== //
func buildLinuxTasks(mode string) []task {
	if mode == "full" {
		return []task{
			cleanFolder("Temp", "/tmp"),
			cleanFolder("Var Temp", "/var/tmp"),
			cleanFolder("Var Cache", "/var/cache"),
			cleanFolder("Systemd Coredump", "/var/lib/systemd/coredump"),
			cleanFolder("Var Crash", "/var/crash"),
			{desc: "All System Logs (>10 days)", cmd: []string{"sh", "-c", "find /var/log -type f -mtime +10 -exec rm -f {} +"}},
			{desc: "fc-cache", cmd: []string{"fc-cache", "-fr"}},
			{desc: "Systemd-Tmpfiles", cmd: []string{"systemd-tmpfiles", "--clean"}},
			{desc: "Apt Cache", cmd: []string{"apt-get", "clean"}},
			{desc: "Flatpak Cache", cmd: []string{"flatpak", "uninstall", "--unused", "-y"}},
			{desc: "Pacman Cache", cmd: []string{"pacman", "-Scc", "--noconfirm"}},
			{desc: "Nix Garbage Collector", cmd: []string{"nix-collect-garbage", "-d"}},
			{desc: "Yarn Cache", cmd: []string{"yarn", "cache", "clean"}},
			{desc: "DNF Cache", cmd: []string{"dnf", "clean", "all"}},
			{desc: "Npm Cache", cmd: []string{"npm", "cache", "clean", "--force"}},
			{desc: "Composer Cache", cmd: []string{"composer", "clear-cache"}},
			{desc: "Pip Cache", cmd: []string{"pip", "cache", "purge"}},
			{desc: "Go Module Cache", cmd: []string{"go", "clean", "-modcache"}},
			{desc: "Rust Cargo Cache", cmd: []string{"cargo", "clean"}},
			{desc: "Docker System Prune", cmd: []string{"docker", "system", "prune", "-af"}},
			{desc: "Podman System Prune", cmd: []string{"podman", "system", "prune", "-af"}},
		}
	}
	return []task{
		cleanFolder("Temp", "/tmp"),
		{desc: "Journal Logs (>100 days)", cmd: []string{"journalctl", "--vacuum-time=100d"}},
		{desc: "fc-cache", cmd: []string{"fc-cache", "-fr"}},
		{desc: "Apt Cache", cmd: []string{"apt-get", "clean"}},
		{desc: "Flatpak Cache", cmd: []string{"flatpak", "uninstall", "--unused", "-y"}},
		{desc: "Pacman Cache", cmd: []string{"pacman", "-Scc", "--noconfirm"}},
		{desc: "DNF Cache", cmd: []string{"dnf", "clean", "all"}},
	}
}

// ==================== EXECUTION ==================== //

// askVerbose asks the user if verbose logging should be enabled (keyboard version)
func askVerbose() {
	if skipPause {
		return
	}

	fmt.Printf("Enable verbose logging?\n")
	fmt.Printf("%s[Y]es%s/%s[N]o%s: ", GREEN, RC, RED, RC)

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		switch {
		case char == 'y' || char == 'Y':
			verbose = true
			fmt.Printf("YES")
			return
		case char == 'n' || char == 'N' || key == keyboard.KeyEsc:
			verbose = false
			fmt.Printf("NO")
			return
		case key == keyboard.KeyCtrlC:
			// Exit on CTRL+C
			cc_exit()
		}
	}
}

// previewTasks prints all tasks before execution
func previewTasks(tasks []task) {
	fmt.Printf("\n\n")
	printInfo("The following cleanup tasks will be executed:")

	for _, t := range tasks {
		var detail string
		if t.goFunc != nil {
			detail = t.path // automatically shows the folder path
		} else if len(t.cmd) > 0 {
			detail = strings.Join(t.cmd, " ") // join command slice into a single string
		}

		// print task description and details
		fmt.Printf("%sCleaning: %s%s\n  └─ %s\n", CYAN, t.desc, RC, detail)
	}

	printInfo("The above cleanup tasks will be executed")
	if verbose {
		printInfo("Verbose mode: showing all task details and errors")
	} else {
		printInfo("Verbose mode OFF: errors will be hidden")
	}
	printInfo("Press [CTRL+C] to cancel")
	printInfo("!!! You use this tool at your own risk !!!")
	if !skipPause {
		pause()
	}
}

// runTasks executes all cleanup tasks with spinner and logging
func runTasks(tasks []task) {
	fmt.Printf("\n%s#############################################%s\n", RED, RC)
	printInfo("*** CrunchyCleaner Cleanup STARTED ***")
	time.Sleep(2 * time.Second)

	startFree := getFreeMB()
	// iterate over all tasks
	for _, t := range tasks {
		// create a context to control the async spinner
		ctx, cancel := context.WithCancel(context.Background())
		go asyncSpinner(ctx, "Cleaning: "+t.desc)
		time.Sleep(CMDWAIT)

		var err error
		// execute task: either Go-native function or shell command
		if t.goFunc != nil {
			err = t.goFunc()
		} else if len(t.cmd) > 0 {
			_, err = runCommand(t.cmd)
		}
		cancel()
		fmt.Printf("\r\033[2K") // clear spinner line

		if verbose && err != nil {
			fmt.Printf("%sCleaning: %s%s FINISHED\n  └─ %s\n", CYAN, t.desc, RC, err)
		} else {
			fmt.Printf("%sCleaning: %s%s FINISHED\n", CYAN, t.desc, RC)
		}
	}

	// calculate freed disk space
	endFree := getFreeMB()
	diff := endFree - startFree
	if diff < 0 {
		diff = 0
	}
	printInfo(fmt.Sprintf("Cleaned approx: %s%d MB%s disk space", YELLOW, diff, RC))
	printSuccess("*** CrunchyCleaner Cleanup FINISHED ***")
}
