package internal_test

import (
	"strings"
	"testing"

	"github.com/gipo355/patch-jetbrains-ide/internal"
)

func TestBuildWrapperContent(t *testing.T) {
	got := internal.BuildWrapperContent("/home/wolf/.local/share/JetBrains/Toolbox/scripts/idea")

	for _, want := range []string{
		"#!/bin/sh\n",
		`nohup "/home/wolf/.local/share/JetBrains/Toolbox/scripts/idea" "$@" >/dev/null 2>&1 &`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("wrapper missing %q\n--- got ---\n%s", want, got)
		}
	}
}

func TestPathOrderWarning(t *testing.T) {
	const target = "/home/wolf/.local/bin"
	const scripts = "/home/wolf/.local/share/JetBrains/Toolbox/scripts"
	sep := ":"

	cases := []struct {
		name    string
		path    string
		wantMsg bool
	}{
		{"target first", target + sep + scripts + sep + "/usr/bin", false},
		{"target present, scripts absent", "/usr/bin" + sep + target, false},
		{"scripts shadows target", scripts + sep + target, true},
		{"target not on path", "/usr/bin" + sep + scripts, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := internal.PathOrderWarning(target, scripts, c.path)
			if (got != "") != c.wantMsg {
				t.Errorf("PathOrderWarning(...) = %q, wantMsg=%v", got, c.wantMsg)
			}
		})
	}
}
