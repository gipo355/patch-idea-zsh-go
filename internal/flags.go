package internal

import (
	"github.com/ogier/pflag"
)

func DefineFlags() (helpFlag, dryRunFlag, allIdesFlag, allFilesFlag, repatchFlag, currentShellFlag *bool) {
	helpFlag = pflag.BoolP("help", "h", false, "Show help")
	dryRunFlag = pflag.BoolP("dry-run", "d", false, "Dry run")
	allIdesFlag = pflag.BoolP("all-ides", "a", false, "Select all IDEs")
	allFilesFlag = pflag.BoolP("all-files", "y", false, "Select all files")
	repatchFlag = pflag.BoolP("repatch", "r", false, "Repatch")
	currentShellFlag = pflag.BoolP("current-shell", "c", false, "Use current shell from $SHELL")
	return
}
