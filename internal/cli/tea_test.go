package cli

import (
	"errors"
	"testing"
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
