// Package cli contains the main CLI logic for the go-rcon tool
package cli

import (
	"errors"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/discohaus/go-rcon/pkg/rcon"
)

// Cli is the main struct for the go-rcon CLI tool
type Cli struct {
	rconClient *rcon.Client
}

// Run starts the TUI program for the go-rcon CLI tool.
func (c *Cli) Run() error {
	p := tea.NewProgram(initialModel(c.rconClient), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error in TUI: %w", err)
	}
	return nil
}

// NewCli creates a new Cli instance with the given host, port, password, and character set.
func NewCli(host *string, port *int32, password *string, charSet *string) (*Cli, error) {

	var parsedCharSet rcon.CharSet
	switch strings.ToLower(*charSet) {
	case "latin1":
		parsedCharSet = rcon.CharSetLatin1
	case "ascii":
		parsedCharSet = rcon.CharSetASCII
	case "utf8":
		parsedCharSet = rcon.CharSetUTF8
	default:
		parsedCharSet = rcon.CharSetASCII
	}

	if host == nil {
		return nil, errors.New("host required")
	}

	if port == nil {
		return nil, errors.New("port required")
	}

	var parsedPassword string
	if password != nil {
		parsedPassword = *password
	}
	rconClient := rcon.NewClient(fmt.Sprintf("rcon://%s:%d", *host, int(*port)), parsedPassword, rcon.WithOptions(parsedCharSet))
	if err := rconClient.CheckConnection(); err != nil {
		return nil, err
	}
	cli := &Cli{
		rconClient: rconClient,
	}

	return cli, nil
}
