package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

type fakeExecutor struct {
	output  string
	err     error
	command string
}

func (f *fakeExecutor) Send(command string) (string, error) {
	f.command = command
	return f.output, f.err
}

func TestModelAsyncSuccessAndBusy(t *testing.T) {
	fake := &fakeExecutor{output: "ok"}
	m := newModel(fake)
	m.textarea.SetValue("say hello")
	_, cmd := m.submit()
	if !m.busy || cmd == nil {
		t.Fatal("submit did not enter busy state")
	}
	result := cmd()
	updated, _ := m.Update(result)
	model := updated.(model)
	if model.busy || fake.command != "say hello" || len(model.messages) != 2 {
		t.Fatalf("unexpected model after send: busy=%v command=%q messages=%v", model.busy, fake.command, model.messages)
	}
}

func TestModelAsyncErrorUpdatesStatus(t *testing.T) {
	fake := &fakeExecutor{err: errors.New("server unavailable")}
	m := newModel(fake)
	m.textarea.SetValue("say hello")
	_, cmd := m.submit()
	updated, _ := m.Update(cmd())
	model := updated.(model)
	if model.status != "error" {
		t.Fatalf("status = %q, want error", model.status)
	}
}

func TestModelShowsNoResponseForEmptyOutput(t *testing.T) {
	m := newModel(&fakeExecutor{output: ""})
	m.textarea.SetValue("list")
	_, cmd := m.submit()
	updated, _ := m.Update(cmd())
	model := updated.(model)

	if len(model.messages) != 2 || !strings.Contains(model.messages[1], "[no response]") {
		t.Fatalf("messages = %v, want no-response message", model.messages)
	}
}

func TestModelUnknownCommandAndExit(t *testing.T) {
	fake := &fakeExecutor{}
	m := newModel(fake)
	m.textarea.SetValue("/wat")
	_, cmd := m.submit()
	if cmd != nil || len(m.messages) != 2 || m.busy {
		t.Fatal("unknown command was not handled locally")
	}
	m.textarea.SetValue("/exit")
	_, cmd = m.submit()
	if cmd == nil {
		t.Fatal("/exit did not return a quit command")
	}
}

func TestModelWrapsMessagesToViewportWidth(t *testing.T) {
	m := newModel(&fakeExecutor{})
	m.resize(10, 20)
	m.addMessage("12345678901234567890", lipgloss.NewStyle())

	if lines := m.viewport.TotalLineCount(); lines != 2 {
		t.Fatalf("wrapped message has %d lines, want 2", lines)
	}

	m.resize(80, 20)
	if lines := m.viewport.TotalLineCount(); lines != 1 {
		t.Fatalf("resized message has %d lines, want 1", lines)
	}
}
