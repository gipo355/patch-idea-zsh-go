package internal_test

import (
	"testing"

	"github.com/gipo355/patch-jetbrains-ide/internal"
)

func TestIsSupportedShell(t *testing.T) {
	cases := map[string]bool{
		"zsh":           true,
		"bash":          true,
		"sh":            true,
		"/usr/bin/zsh":  true, // $SHELL is an absolute path
		"/bin/bash":     true,
		"fish":          false,
		"/usr/bin/fish": false,
		"":              false,
	}
	for in, want := range cases {
		if got := internal.IsSupportedShell(in); got != want {
			t.Errorf("IsSupportedShell(%q) = %v, want %v", in, got, want)
		}
	}
}
