package internal

import (
	"fmt"
	"os"
	"os/exec"
)

func DetermineShell(currentShellFlag bool) string {
	var shell string
	if currentShellFlag {
		shell = os.Getenv("SHELL")
		if shell == "" || shell != "bash" && shell != "sh" && shell != "zsh" {
			fmt.Fprintln(os.Stderr, "\x1b[31mSHELL environment variable not set or invalid.\x1b[0m")
			os.Exit(1)
		}
	} else {
		fmt.Println("Choose the shell to use (sh/bash/zsh, default is zsh):")
		fmt.Print("> ")
		shellInput := ReadLine()
		if shellInput == "" {
			shellInput = "zsh"
		}
		shell = shellInput
		if shell != "bash" && shell != "sh" && shell != "zsh" {
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
