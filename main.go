package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ogier/pflag"

	"github.com/gipo355/patch-idea-zsh-go/internal"
)

func main() {
	helpFlag, dryRunFlag, allIdesFlag, allFilesFlag, repatchFlag, currentShellFlag := internal.DefineFlags()

	pflag.Parse()

	if *helpFlag {
		pflag.Usage()
		os.Exit(0)
	}

	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		fmt.Fprintln(os.Stderr, "\x1b[31mHOME environment variable not set.\x1b[0m")
		os.Exit(1)
	}

	dirPath := filepath.Join(homeDir, ".local", "share", "applications")

	shell := internal.DetermineShell(*currentShellFlag)

	shellPath := internal.GetShellPath(shell)

	selectedIDEs := internal.GetSelectedIDEs(*allIdesFlag)

	matchingFiles := internal.FindMatchingFiles(dirPath, selectedIDEs)

	filesToPatch := internal.GetFilesToPatch(matchingFiles, *allFilesFlag)

	internal.PatchFiles(filesToPatch, shellPath, *dryRunFlag, *repatchFlag)
}
