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

// BuildPatchedContent rewrites a JetBrains .desktop file so its Exec line
// launches the IDE through an interactive shell, making the IDE inherit the
// environment (PATH, version-manager shims, etc.) of that shell. The original
// content is preserved, commented out, below a "# patched on <patchedAt>" marker.
//
// It returns the new file content and whether the file was already patched with
// shellPath. When alreadyPatched is true (and repatch is false) the returned
// content is empty and the caller should leave the file untouched.
func BuildPatchedContent(content, shellPath, patchedAt string, repatch bool) (string, bool) {
	lines := strings.Split(content, "\n")

	var oldExecContent string
	var modifiedOldContent []string
	var modifiedContent []string

	for _, line := range lines {
		if !repatch && strings.HasPrefix(line, "Exec="+shellPath) {
			return "", true
		}

		if strings.HasPrefix(line, "Exec=") {
			start := strings.Index(line, "\"") + 1
			end := strings.LastIndex(line, "\"")

			if start > 0 && end > start {
				oldExecContent = line[start:end]
			}

			// The literal double quotes around the IDE path are required by the
			// .desktop Exec format; GLib-based launchers (gio/gtk-launch) refuse
			// to load an Exec line that contains nested/escaped quotes.
			//nolint:gocritic // literal double quotes are required by the .desktop Exec format, not %q escaping
			newExecLine := fmt.Sprintf(`Exec=%s -i -c "%s" %%u`, shellPath, oldExecContent)
			modifiedContent = append(modifiedContent, newExecLine)
			modifiedOldContent = append(modifiedOldContent, "# "+line)
			continue
		}

		if strings.HasPrefix(line, "#") {
			modifiedOldContent = append(modifiedOldContent, line)
		} else {
			modifiedOldContent = append(modifiedOldContent, "# "+line)
			modifiedContent = append(modifiedContent, line)
		}
	}

	finalOldContent := fmt.Sprintf("# patched on %s\n%s", patchedAt, strings.Join(modifiedOldContent, "\n"))
	return fmt.Sprintf("%s\n%s", strings.Join(modifiedContent, "\n"), finalOldContent), false
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

		patchedAt := time.Now().Format("2006-01-02 15:04:05")
		finalContent, alreadyPatched := BuildPatchedContent(string(content), shellPath, patchedAt, repatchFlag)
		if alreadyPatched {
			fmt.Printf("\x1b[33mx\x1b[0m File %s is already patched. Skipping.\n", filePath)
			continue
		}

		if !dryRunFlag {
			//nolint:gosec // .desktop launcher files must stay world-readable (0644)
			if err = os.WriteFile(filePath, []byte(finalContent), 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write to the file: %s\n", filePath)
				os.Exit(1)
			}
		}

		fmt.Printf("\x1b[32mv\x1b[0m Patched file: %s\n", filePath)

		if dryRunFlag {
			fmt.Println(finalContent)
		}
	}

	if dryRunFlag {
		fmt.Println("\nno action taken -- dry run mode on")
	}
}
