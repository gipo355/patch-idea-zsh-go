package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gipo355/patch-jetbrains-ide/internal"
	"github.com/spf13/pflag"
)

func main() {
	flags := internal.DefineFlags()

	pflag.Parse()

	if *flags.Help {
		pflag.Usage()
		os.Exit(0)
	}

	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		fmt.Fprintln(os.Stderr, "\x1b[31mHOME environment variable not set.\x1b[0m")
		os.Exit(1)
	}

	if !*flags.WrappersOnly {
		dirPath := filepath.Join(homeDir, ".local", "share", "applications")

		shell := internal.DetermineShell(*flags.CurrentShell)

		shellPath := internal.GetShellPath(shell)

		selectedIDEs := internal.GetSelectedIDEs(*flags.AllIDEs)

		matchingFiles := internal.FindMatchingFiles(dirPath, selectedIDEs)

		filesToPatch := internal.GetFilesToPatch(matchingFiles, *flags.AllFiles)

		internal.PatchFiles(filesToPatch, shellPath, *flags.DryRun, *flags.Repatch)
	}

	if *flags.Wrappers || *flags.WrappersOnly {
		scriptsDir := filepath.Join(homeDir, ".local", "share", "JetBrains", "Toolbox", "scripts")
		targetDir := *flags.WrapperDir
		if targetDir == "" {
			targetDir = filepath.Join(homeDir, ".local", "bin")
		}
		internal.GenerateWrappers(scriptsDir, targetDir, *flags.DryRun, *flags.Repatch)
	}
}
