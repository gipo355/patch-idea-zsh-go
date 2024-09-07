package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func FindMatchingFiles(dirPath string, selectedIDEs []JetBrainsIDE) []string {
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

	fmt.Printf(
		"\x1b[32mFound the following %d JetBrains IDEA desktop files:\x1b[0m\n",
		len(matchingFiles),
	)
	for i, file := range matchingFiles {
		fmt.Printf("%d: %s\n", i+1, file)
	}
	return matchingFiles
}

func GetFilesToPatch(matchingFiles []string, allFilesFlag bool) []string {
	if allFilesFlag {
		return matchingFiles
	}

	fmt.Println(
		"Enter the numbers of the files you want to patch, separated by commas (default is all):",
	)
	fmt.Print("> ")
	input := ReadLine()

	if input == "" {
		return matchingFiles
	}

	var filesToPatch []string
	for _, s := range strings.Split(input, ",") {
		i, err := strconv.Atoi(strings.TrimSpace(s))
		if err == nil && i > 0 && i <= len(matchingFiles) {
			filesToPatch = append(filesToPatch, matchingFiles[i-1])
		}
	}
	return filesToPatch
}

func PatchFiles(filesToPatch []string, shellPath string, dryRunFlag, repatchFlag bool) {
	if len(filesToPatch) == 0 {
		fmt.Fprintln(os.Stderr, "\x1b[31mNo files selected for patching.\x1b[0m")
		os.Exit(1)
	}

	fmt.Printf("\x1b[32mPatching %d files.\x1b[0m\n", len(filesToPatch))

	for _, filePath := range filesToPatch {
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
			if !repatchFlag && strings.HasPrefix(line, fmt.Sprintf("Exec=%s", shellPath)) {
				fmt.Printf("\x1b[33mx\x1b[0m File %s is already patched. Skipping.\n", filePath)
				alreadyPatched = true
				break
			}

			if strings.HasPrefix(line, "Exec=") {
				start := strings.Index(line, "\"") + 1
				end := strings.LastIndex(line, "\"")

				if start > 0 && end > start {
					oldExecContent = line[start:end]
				}

				newExecLine := fmt.Sprintf(`Exec=%s -i -c "%s" %%u`, shellPath, oldExecContent)
				modifiedContent = append(modifiedContent, newExecLine)
				continue
			}

			if strings.HasPrefix(line, "#") {
				modifiedOldContent = append(modifiedOldContent, line)
			} else {
				modifiedOldContent = append(modifiedOldContent, "# "+line)
				modifiedContent = append(modifiedContent, line)
			}
		}

		if alreadyPatched {
			continue
		}

		currentDate := time.Now().Format("2006-01-02 15:04:05")
		finalOldContent := fmt.Sprintf(
			"# patched on %s\n%s",
			currentDate,
			strings.Join(modifiedOldContent, "\n"),
		)
		finalContent := fmt.Sprintf("%s\n%s", strings.Join(modifiedContent, "\n"), finalOldContent)

		if !dryRunFlag {
			err = os.WriteFile(filePath, []byte(finalContent), 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write to the file: %s\n", filePath)
				os.Exit(1)
			}
		}

		fmt.Printf("\x1b[32mv\x1b[0m Patched file: %s\n", filePath)

		if dryRunFlag {
			println(finalContent)
		}
	}

	if dryRunFlag {
		println("\nno action taken -- dry run mode on")
	}
}
