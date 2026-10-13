package cli

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestHeaderIconOnlyInSupportedTerminal(t *testing.T) {
	t.Setenv("TERM", "xterm-kitty")
	t.Setenv("TERM_PROGRAM", "")
	t.Setenv("TMUX", "")
	t.Setenv("STY", "")
	withIcon := newModel(&fakeExecutor{})
	if !strings.Contains(withIcon.View(), "\x1b_Ga=T") {
		t.Fatal("Kitty header has no graphics command")
	}

	t.Setenv("TERM", "foot")
	foot := newModel(&fakeExecutor{})
	if !strings.Contains(foot.View(), "\x1bP0;1;0q") || strings.Contains(foot.View(), "\x1b_G") {
		t.Fatal("Foot header has no Sixel image or contains a Kitty command")
	}
	if !strings.HasPrefix(foot.View(), "   ") {
		t.Fatal("Foot header does not reserve space between icon and text")
	}
	painted := ansi.Truncate(strings.Split(foot.View(), "\n")[0], 80, "")
	if !strings.Contains(painted, "\x1b7\r\x1bP0;1;0q") || !strings.HasSuffix(painted, "\x1b\\\x1b8") {
		t.Fatal("Bubble Tea line truncation changed the Sixel placement")
	}

	t.Setenv("TERM", "xterm-256color")
	withoutIcon := newModel(&fakeExecutor{})
	withoutIcon.iconProtocol = noIcon
	if strings.Contains(withoutIcon.View(), "\x1b_G") || strings.Contains(withoutIcon.View(), "\x1bP0;1;0q") {
		t.Fatal("unsupported terminal received graphics command")
	}

	t.Setenv("TERM", "xterm-kitty")
	t.Setenv("TMUX", "1")
	if newModel(&fakeExecutor{}).iconProtocol != noIcon {
		t.Fatal("icon should be disabled inside tmux")
	}
}

func TestFootDetectionWithGenericTerm(t *testing.T) {
	if got := iconProtocolFor("xterm-256color", "", false, true); got != footSixelIcon {
		t.Fatalf("Foot with generic TERM detected as %v", got)
	}
	if got := iconProtocolFor("xterm-256color", "", false, false); got != noIcon {
		t.Fatalf("unsupported terminal detected as %v", got)
	}
	if got := iconProtocolFor("xterm-256color", "", true, true); got != noIcon {
		t.Fatalf("Foot inside tmux detected as %v", got)
	}
}

func TestHeaderIconPayloadAndCellWidth(t *testing.T) {
	const marker = ";"
	start := strings.Index(kittyIconSequence, marker)
	end := strings.Index(kittyIconSequence, "\x1b\\")
	if start < 0 || end < start {
		t.Fatal("invalid Kitty graphics sequence")
	}
	data, err := base64.StdEncoding.DecodeString(kittyIconSequence[start+1 : end])
	if err != nil || !bytes.HasPrefix(data, []byte("\x89PNG\r\n\x1a\n")) {
		t.Fatalf("header icon is not a PNG: %v", err)
	}
	if width := ansi.StringWidth(kittyIconSequence); width != 3 { // two cells plus a gap
		t.Fatalf("header icon occupies %d cells, want 3", width)
	}
	if !strings.Contains(ansi.Truncate(kittyIconSequence+"RCON", 80, ""), "\x1b_Ga=T") {
		t.Fatal("Bubble Tea line truncation removed the graphics command")
	}
	if !strings.HasPrefix(footIconSequence, "\x1bP0;1;0q\"1;1;28;31") || !strings.HasSuffix(footIconSequence, "\x1b\\") {
		t.Fatal("invalid Foot Sixel sequence")
	}
}
