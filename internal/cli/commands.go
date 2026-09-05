package cli

import "strings"

type commandHandler func(*model) commandAction

type command struct {
	Name        string
	Description string
	Handler     commandHandler
}

type commandRegistry struct{ commands []command }

func newCommandRegistry() commandRegistry {
	return commandRegistry{commands: []command{
		{Name: "/help", Description: "show available commands", Handler: func(m *model) commandAction { return commandAction{message: m.registry.help()} }},
		{Name: "/exit", Description: "quit the client", Handler: func(*model) commandAction { return commandAction{quit: true} }},
	}}
}

func (r commandRegistry) find(input string) (command, bool) {
	for _, item := range r.commands {
		if item.Name == input {
			return item, true
		}
	}
	return command{}, false
}

func (r commandRegistry) suggestions(prefix string) []string {
	var result []string
	for _, item := range r.commands {
		if strings.HasPrefix(item.Name, prefix) {
			result = append(result, item.Name)
		}
	}
	return result
}

func (r commandRegistry) help() string {
	lines := []string{"UI commands:"}
	for _, item := range r.commands {
		lines = append(lines, "  "+item.Name+"  "+item.Description)
	}
	return strings.Join(lines, "\n")
}

type parsedInput struct {
	command, argument string
	isUI              bool
}

func parseInput(input string) parsedInput {
	input = strings.TrimSpace(input)
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return parsedInput{}
	}
	if !strings.HasPrefix(parts[0], "/") {
		return parsedInput{command: input}
	}
	result := parsedInput{command: parts[0], isUI: true}
	if result.isUI {
		result.command = strings.ToLower(result.command)
	}
	if len(parts) > 1 {
		result.argument = strings.Join(parts[1:], " ")
	}
	return result
}

type commandAction struct {
	message string
	quit    bool
}
