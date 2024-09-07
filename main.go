package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	helpFlag := flag.Bool("h", false, "Show help")
	allFlag := flag.Bool("a", false, "Select all IDEs")
	currentShellFlag := flag.Bool("c", false, "Use current shell from $SHELL")

	// Parse flags
	flag.Parse()

	// Show help and exit if -h is provided
	if *helpFlag {
		fmt.Println("Usage: main [-h] [-a] [-c]")
		fmt.Println("Options:")
		fmt.Println("  -h  Show help")
		fmt.Println("  -a  Select all IDEs")
		fmt.Println("  -c  Use current shell from $SHELL")
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
		if shell == "" {
			fmt.Fprintln(os.Stderr, "\x1b[31mSHELL environment variable not set.\x1b[0m")
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

	// Create the new Exec line with the home directory path to use
	newExecLine := fmt.Sprintf(
		`Exec=%s -i -c "%s/.local/share/JetBrains/Toolbox/apps/intellij-idea-ultimate/bin/idea" %%u`,
		shellPath,
		homeDir,
	)

	// Determine the IDEs to patch
	var selectedIDEs []JetBrainsIDE
	if *allFlag {
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
	fmt.Println("\x1b[32mFound the following JetBrains IDEA desktop files:\x1b[0m")
	for i, file := range matchingFiles {
		fmt.Printf("%d: %s\n", i+1, file)
	}

	// Ask for confirmation
	fmt.Println(
		"Enter the numbers of the files you want to patch, separated by commas (default is all):",
	)
	fmt.Print("> ")
	input := readLine()

	var filesToPatch []string
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

	if len(filesToPatch) == 0 {
		fmt.Fprintln(os.Stderr, "\x1b[31mNo files selected for patching.\x1b[0m")
		os.Exit(1)
	}

	// Loop through each file and patch it
	for _, filePath := range filesToPatch {
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\x1b[31mFailed to read the file: %s\x1b[0m\n", filePath)
			os.Exit(1)
		}

		lines := strings.Split(string(content), "\n")
		alreadyPatched := false
		for _, line := range lines {
			if strings.HasPrefix(line, fmt.Sprintf("Exec=%s", shellPath)) {
				fmt.Printf("\x1b[33mx\x1b[0m File %s is already patched. Skipping.\n", filePath)
				alreadyPatched = true
				break
			}
		}
		if alreadyPatched {
			continue
		}

		currentDate := time.Now().Format("2006-01-02 15:04:05")

		var modifiedContent []string
		for _, line := range lines {
			if strings.HasPrefix(line, "Exec=") {
				modifiedContent = append(modifiedContent, newExecLine)
			} else {
				modifiedContent = append(modifiedContent, line)
			}
		}

		var modifiedOldContent []string
		for _, line := range lines {
			if strings.HasPrefix(line, "#") {
				modifiedOldContent = append(modifiedOldContent, line)
			} else {
				modifiedOldContent = append(modifiedOldContent, "# "+line)
			}
		}

		finalOldContent := fmt.Sprintf(
			"\n# patched on %s\n%s",
			currentDate,
			strings.Join(modifiedOldContent, "\n"),
		)
		finalContent := fmt.Sprintf(
			"%s\n\n%s",
			strings.Join(modifiedContent, "\n"),
			finalOldContent,
		)

		println(finalContent)

		// err = os.WriteFile(filePath, []byte(finalContent), 0644)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Failed to write to the file: %s\n", filePath)
		// 	os.Exit(1)
		// }

		fmt.Printf("\x1b[32mv\x1b[0m Patched file: %s\n", filePath)
	}
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}
