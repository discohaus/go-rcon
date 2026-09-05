// Package cli contains the main CLI logic for the go-rcon tool
package cli

import (
	"fmt"
	"net"
	"strconv"
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
	if host == nil || strings.TrimSpace(*host) == "" {
		return nil, fmt.Errorf("host required")
	}
	if net.ParseIP(strings.TrimSpace(*host)) == nil && strings.ContainsAny(*host, " \t\n\r") {
		return nil, fmt.Errorf("invalid host %q", *host)
	}
	if port == nil {
		return nil, fmt.Errorf("port required")
	}
	if *port < 1 || *port > 65535 {
		return nil, fmt.Errorf("port must be between 1 and 65535, got %d", *port)
	}
	if charSet == nil || strings.TrimSpace(*charSet) == "" {
		return nil, fmt.Errorf("charset required")
	}

	var parsedCharSet rcon.CharSet
	switch strings.ToLower(strings.TrimSpace(*charSet)) {
	case "latin1":
		parsedCharSet = rcon.CharSetLatin1
	case "ascii":
		parsedCharSet = rcon.CharSetASCII
	case "utf8":
		parsedCharSet = rcon.CharSetUTF8
	default:
		return nil, fmt.Errorf("unsupported charset %q (choose ascii, latin1, or utf8)", *charSet)
	}

	var parsedPassword string
	if password != nil {
		parsedPassword = *password
	}
	rconClient := rcon.NewClient(fmt.Sprintf("rcon://%s:%s", strings.TrimSpace(*host), strconv.FormatInt(int64(*port), 10)), parsedPassword, rcon.WithOptions(parsedCharSet))
	if err := rconClient.CheckConnection(); err != nil {
		return nil, fmt.Errorf("RCON connection failed — is the server running?\n%w", err)
	}
	cli := &Cli{
		rconClient: rconClient,
	}

	return cli, nil
}
