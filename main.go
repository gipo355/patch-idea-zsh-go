package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	flag "github.com/ogier/pflag"
)

type JetBrainsIDE int

const (
	IntelliJ JetBrainsIDE = iota
	PyCharm
	WebStorm
	PhpStorm
	CLion
	Rider
	DataGrip
	RubyMine
	AppCode
)

var ideNames = []string{
	"idea",
	"pycharm",
	"webstorm",
	"phpstorm",
	"clion",
	"rider",
	"datagrip",
	"rubymine",
	"appcode",
}

func (ide JetBrainsIDE) String() string {
	return ideNames[ide]
}

func fromString(input string) (JetBrainsIDE, bool) {
	for i, name := range ideNames {
		if name == input {
			return JetBrainsIDE(i), true
		}
	}
	return 0, false
}

func allIDEs() []JetBrainsIDE {
	return []JetBrainsIDE{
		IntelliJ, PyCharm, WebStorm, PhpStorm, CLion, Rider, DataGrip, RubyMine, AppCode,
	}
}

func main() {
	// Define flags
	helpFlag := flag.BoolP("help", "h", false, "Show help")
	dryRunFlag := flag.BoolP("dry-run", "d", false, "Dry run")
	allIdesFlag := flag.BoolP("all-ides", "a", false, "Select all IDEs")
	allFilesFlag := flag.BoolP("all-files", "y", false, "Select all IDEs")
	repatchFlag := flag.BoolP("repatch", "r", false, "Select all IDEs")
	currentShellFlag := flag.BoolP("current-shell", "c", false, "Use current shell from $SHELL")

	// Parse flags
	flag.Parse()

	// Show help and exit if -h is provided
	if *helpFlag {
		// fmt.Println("Usage: main [-h] [-a] [-c]")
		// fmt.Println("Options:")
		// fmt.Println("  -h  Show help")
		// fmt.Println("  -a  Select all IDEs")
		// fmt.Println("  -c  Use current shell from $SHELL")
		flag.Usage()
		os.Exit(0)
	}

	// Get the .local/share/ directory path
	dirPath := filepath.Join(os.Getenv("HOME"), ".local", "share", "applications")

	// Get the home directory path
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		fmt.Fprintln(os.Stderr, "\x1b[31mHOME environment variable not set.\x1b[0m")
		os.Exit(1)
	}

	// Determine the shell to use
	var shell string
	if *currentShellFlag {
		shell = os.Getenv("SHELL")
		// prevent empty shell var or different from sh/bash/zsh
		if shell == "" || shell != "bash" && shell != "sh" && shell != "zsh" {
			fmt.Fprintln(os.Stderr, "\x1b[31mSHELL environment variable not set or invalid.\x1b[0m")
			os.Exit(1)
		}
	} else {
		fmt.Println("Choose the shell to use (sh/bash/zsh, default is zsh):")
		fmt.Print("> ")
		shellInput := readLine()
		if shellInput == "" {
			shellInput = "zsh"
		}
		shell = shellInput
		if shell != "bash" && shell != "sh" && shell != "zsh" {
			fmt.Fprintln(os.Stderr, "\x1b[31mInvalid shell choice. Please choose either 'bash', 'sh' or 'zsh'.\x1b[0m")
			os.Exit(1)
		}
	}

	// Get the path of the chosen shell
	shellPath, err := exec.LookPath(shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\x1b[31mFailed to find the path for the shell: %s\x1b[0m\n", shell)
		os.Exit(1)
	}

	fmt.Printf("\x1b[32mUsing shell: %s\x1b[0m\n", shellPath)

	// Determine the IDEs to patch
	var selectedIDEs []JetBrainsIDE
	if *allIdesFlag {
		selectedIDEs = allIDEs()
	} else {
		fmt.Println("Choose the JetBrains IDEs to patch (comma-separated numbers, default is all):")
		allIDEs := allIDEs()
		for i, ide := range allIDEs {
			fmt.Printf("%d: %s\n", i+1, ide)
		}
		fmt.Print("> ")
		ideInput := readLine()

		if ideInput == "" {
			selectedIDEs = allIDEs
		} else {
			for _, s := range strings.Split(ideInput, ",") {
				i, err := strconv.Atoi(strings.TrimSpace(s))
				if err == nil && i > 0 && i <= len(allIDEs) {
					selectedIDEs = append(selectedIDEs, allIDEs[i-1])
				}
			}
		}
	}

	// Find all matching files
	files, err := filepath.Glob(filepath.Join(dirPath, "jetbrains-*.desktop"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "\x1b[31mFailed to read the directory: %s\x1b[0m\n", dirPath)
		os.Exit(1)
	}

	var matchingFiles []string
	for _, file := range files {
		for _, ide := range selectedIDEs {
			if strings.Contains(file, ide.String()) {
				matchingFiles = append(matchingFiles, file)
				break
			}
		}
	}

	if len(matchingFiles) == 0 {
		fmt.Fprintln(os.Stderr, "\x1b[31mNo matching JetBrains IDEA desktop files found.\x1b[0m")
		os.Exit(1)
	}

	// List all found files
	fmt.Fprintln(
		// "\x1b[32mFound the following %d JetBrains IDEA desktop files:\x1b[0m",
		os.Stdout,
		fmt.Sprintf(
			"\x1b[32mFound the following %d JetBrains IDEA desktop files:\x1b[0m",
			len(matchingFiles),
		),
	)
	for i, file := range matchingFiles {
		fmt.Printf("%d: %s\n", i+1, file)
	}

	// Ask for confirmation
	var filesToPatch []string
	if !*allFilesFlag {
		fmt.Println(
			"Enter the numbers of the files you want to patch, separated by commas (default is all):",
		)
		fmt.Print("> ")
		input := readLine()

		if input == "" {
			filesToPatch = matchingFiles
		} else {
			for _, s := range strings.Split(input, ",") {
				i, err := strconv.Atoi(strings.TrimSpace(s))
				if err == nil && i > 0 && i <= len(matchingFiles) {
					filesToPatch = append(filesToPatch, matchingFiles[i-1])
				}
			}
		}
	} else {
		filesToPatch = matchingFiles
	}

	if len(filesToPatch) == 0 {
		fmt.Fprintln(os.Stderr, "\x1b[31mNo files selected for patching.\x1b[0m")
		os.Exit(1)
	}

	fmt.Printf("\x1b[32mPatching %d files.\x1b[0m\n", len(filesToPatch))

	// Loop through each file and patch it
	for _, filePath := range filesToPatch {
		// example file:
		// [Desktop Entry]
		// Name=IntelliJ IDEA Ultimate 2024.2.1
		// Exec=/usr/bin/zsh -i -c "/home/wolf/.local/share/JetBrains/Toolbox/apps/intellij-idea-ultimate/bin/idea" %u
		// Version=1.0
		// Type=Application
		// Categories=Development;IDE;
		// Terminal=false
		// Icon=/home/wolf/.local/share/JetBrains/Toolbox/apps/intellij-idea-ultimate/bin/idea.svg
		// Comment=The Leading Java and Kotlin IDE
		// StartupWMClass=jetbrains-idea
		// StartupNotify=true

		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\x1b[31mFailed to read the file: %s\x1b[0m\n", filePath)
			os.Exit(1)
		}

		lines := strings.Split(string(content), "\n")
		alreadyPatched := false

		var oldExecContent string
		var modifiedOldContent []string
		var modifiedContent []string

		for _, line := range lines {

			// prevent repatching if outcome is the same
			if !*repatchFlag && strings.HasPrefix(line, fmt.Sprintf("Exec=%s", shellPath)) {
				fmt.Printf("\x1b[33mx\x1b[0m File %s is already patched. Skipping.\n", filePath)
				alreadyPatched = true
				break
			}

			// swap Exec line and append to modifiedContent if the line starts with Exec=
			if strings.HasPrefix(line, "Exec=") {
				start := strings.Index(line, "\"") + 1
				end := strings.LastIndex(line, "\"")

				if start > 0 && end > start {
					oldExecContent = line[start:end]
				}

				newExecLine := fmt.Sprintf(
					`Exec=%s -i -c "%s" %%u`,
					shellPath,
					oldExecContent,
				)

				modifiedContent = append(modifiedContent, newExecLine)

				continue
			}

			// if line is old content starting with # add to modifiedOldContent, else add to modifiedContent and append to old content with #
			// we will merge those together. the modifiedContent will be the new content wit hthe exec line swapped
			// the old content will keep track of the original content with the old exec line to allow to restore
			if strings.HasPrefix(line, "#") {
				modifiedOldContent = append(modifiedOldContent, line)
			} else {
				modifiedOldContent = append(modifiedOldContent, "# "+line)
				modifiedContent = append(modifiedContent, line)
			}
		}

		// skip if already patched
		if alreadyPatched {
			continue
		}

		currentDate := time.Now().Format("2006-01-02 15:04:05")

		finalOldContent := fmt.Sprintf(
			"# patched on %s\n%s",
			currentDate,
			strings.Join(modifiedOldContent, "\n"),
		)
		finalContent := fmt.Sprintf(
			"%s\n%s",
			strings.Join(modifiedContent, "\n"),
			finalOldContent,
		)

		if !*dryRunFlag {
			err = os.WriteFile(filePath, []byte(finalContent), 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write to the file: %s\n", filePath)
				os.Exit(1)
			}
		}

		fmt.Printf("\x1b[32mv\x1b[0m Patched file: %s\n", filePath)

		if *dryRunFlag {
			println(finalContent)
		}
	}

	if *dryRunFlag {
		println("\nno action taken -- dry run mode on")
	}
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}
