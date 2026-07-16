package internal_test

import (
	"strings"
	"testing"

	"github.com/gipo355/patch-jetbrains-ide/internal"
)

const sampleDesktop = `[Desktop Entry]
Name=IntelliJ IDEA
Exec="/path/to/idea" %u
Type=Application`

func TestBuildPatchedContent(t *testing.T) {
	const want = `[Desktop Entry]
Name=IntelliJ IDEA
Exec=/usr/bin/zsh -i -c "/path/to/idea" %u
Type=Application
# patched on 2024-01-01 00:00:00
# [Desktop Entry]
# Name=IntelliJ IDEA
# Exec="/path/to/idea" %u
# Type=Application`

	got, already := internal.BuildPatchedContent(sampleDesktop, "/usr/bin/zsh", "2024-01-01 00:00:00", false)
	if already {
		t.Fatal("BuildPatchedContent reported already-patched for a clean file")
	}
	if got != want {
		t.Errorf("patched content mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
}

func TestBuildPatchedContentPreservesExecFormat(t *testing.T) {
	got, _ := internal.BuildPatchedContent(sampleDesktop, "/usr/bin/zsh", "2024-01-01 00:00:00", false)

	// The IDE path must be wrapped in a single pair of literal double quotes with
	// no escaping: GLib launchers (gio/gtk-launch) reject nested/escaped quotes.
	const wantExec = `Exec=/usr/bin/zsh -i -c "/path/to/idea" %u`
	if !strings.Contains(got, wantExec) {
		t.Errorf("expected patched Exec line %q not found in:\n%s", wantExec, got)
	}
	if strings.Contains(got, `\"`) {
		t.Error("patched output contains escaped quotes, which break GLib-based launchers")
	}
}

func TestBuildPatchedContentIsIdempotent(t *testing.T) {
	patched, _ := internal.BuildPatchedContent(sampleDesktop, "/usr/bin/zsh", "2024-01-01 00:00:00", false)

	_, already := internal.BuildPatchedContent(patched, "/usr/bin/zsh", "2024-01-02 00:00:00", false)
	if !already {
		t.Error("re-patching an already-patched file should report already-patched")
	}
}

func TestBuildPatchedContentRepatch(t *testing.T) {
	patched, _ := internal.BuildPatchedContent(sampleDesktop, "/usr/bin/zsh", "2024-01-01 00:00:00", false)

	_, already := internal.BuildPatchedContent(patched, "/usr/bin/zsh", "2024-01-02 00:00:00", true)
	if already {
		t.Error("repatch=true must re-process even an already-patched file")
	}
}
