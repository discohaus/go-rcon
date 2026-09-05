package cli

import "testing"

func TestParseInput(t *testing.T) {
	tests := []struct {
		input, command, argument string
		isUI                     bool
	}{
		{"time set day", "time set day", "", false},
		{" /HELP now ", "/help", "now", true},
	}
	for _, test := range tests {
		got := parseInput(test.input)
		if got.command != test.command || got.argument != test.argument || got.isUI != test.isUI {
			t.Errorf("parseInput(%q) = %#v", test.input, got)
		}
	}
}

func TestRegistryAndCompletion(t *testing.T) {
	registry := newCommandRegistry()
	if _, ok := registry.find("/exit"); !ok {
		t.Fatal("/exit is not registered")
	}
	got := registry.suggestions("/h")
	if len(got) != 1 || got[0] != "/help" {
		t.Fatalf("suggestions = %#v", got)
	}
	if _, ok := registry.find("/missing"); ok {
		t.Fatal("unknown command was found")
	}
}
