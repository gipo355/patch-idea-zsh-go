package internal_test

import (
	"testing"

	"github.com/gipo355/patch-jetbrains-ide/internal"
)

func TestIdeEnumTablesAligned(t *testing.T) {
	all := internal.AllIDEs()
	if len(all) != len(internal.IdeNames) {
		t.Fatalf("AllIDEs (%d) and IdeNames (%d) lengths differ", len(all), len(internal.IdeNames))
	}
	for i, ide := range all {
		if int(ide) != i {
			t.Errorf("AllIDEs[%d] = %d, want %d (out of order)", i, ide, i)
		}
		if ide.String() != internal.IdeNames[i] {
			t.Errorf("ide %d String() = %q, want %q", i, ide.String(), internal.IdeNames[i])
		}
	}
}

func TestFromString(t *testing.T) {
	for _, name := range internal.IdeNames {
		ide, ok := internal.FromString(name)
		if !ok {
			t.Errorf("FromString(%q) returned ok=false", name)
			continue
		}
		if ide.String() != name {
			t.Errorf("FromString(%q).String() = %q, round-trip failed", name, ide.String())
		}
	}

	if _, ok := internal.FromString("not-an-ide"); ok {
		t.Error("FromString on unknown name should return ok=false")
	}
}
