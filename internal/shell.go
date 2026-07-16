package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const defaultShell = "zsh"

// supportedShells are the shells the generated Exec line may be wrapped in.
var supportedShells = map[string]bool{"sh": true, "bash": true, "zsh": true}

// IsSupportedShell reports whether name is a shell this tool knows how to wrap.
// $SHELL is an absolute path (e.g. /usr/bin/zsh), so we compare the base name.
func IsSupportedShell(name string) bool {
	return supportedShells[filepath.Base(name)]
}

func DetermineShell(currentShellFlag bool) string {
	var shell string
	if currentShellFlag {
		shell = os.Getenv("SHELL")
		if !IsSupportedShell(shell) {
			fmt.Fprintln(os.Stderr, "\x1b[31mSHELL environment variable not set or invalid.\x1b[0m")
			os.Exit(1)
		}
		shell = filepath.Base(shell)
	} else {
		fmt.Println("Choose the shell to use (sh/bash/zsh, default is zsh):")
		fmt.Print("> ")
		shellInput := ReadLine()
		if shellInput == "" {
			shellInput = defaultShell
		}
		shell = shellInput
		if !IsSupportedShell(shell) {
			fmt.Fprintln(os.Stderr, "\x1b[31mInvalid shell choice. Please choose either 'bash', 'sh' or 'zsh'.\x1b[0m")
			os.Exit(1)
		}
	}
	return shell
}

func GetShellPath(shell string) string {
	shellPath, err := exec.LookPath(shell)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\x1b[31mFailed to find the path for the shell: %s\x1b[0m\n", shell)
		os.Exit(1)
	}
	fmt.Printf("\x1b[32mUsing shell: %s\x1b[0m\n", shellPath)
	return shellPath
}
