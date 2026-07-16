package internal

import (
	"github.com/spf13/pflag"
)

// Flags holds pointers to the parsed command-line flags. The pointers are
// populated by pflag.Parse(), which the caller must run after DefineFlags.
type Flags struct {
	Help         *bool
	DryRun       *bool
	AllIDEs      *bool
	AllFiles     *bool
	Repatch      *bool
	CurrentShell *bool
	Wrappers     *bool
	WrappersOnly *bool
	WrapperDir   *string
}

func DefineFlags() *Flags {
	return &Flags{
		Help:         pflag.BoolP("help", "h", false, "Show help"),
		DryRun:       pflag.BoolP("dry-run", "d", false, "Dry run"),
		AllIDEs:      pflag.BoolP("all-ides", "a", false, "Select all IDEs"),
		AllFiles:     pflag.BoolP("all-files", "y", false, "Select all files"),
		Repatch:      pflag.BoolP("repatch", "r", false, "Repatch (overwrite existing patches/wrappers)"),
		CurrentShell: pflag.BoolP("current-shell", "c", false, "Use current shell from $SHELL"),
		Wrappers:     pflag.BoolP("wrappers", "w", false, "Also generate detaching terminal wrappers"),
		WrappersOnly: pflag.Bool("wrappers-only", false, "Only generate wrappers, skip desktop patching"),
		WrapperDir:   pflag.String("wrapper-dir", "", "Target dir for wrappers (default ~/.local/bin)"),
	}
}
