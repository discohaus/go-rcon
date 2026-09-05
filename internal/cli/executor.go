package cli

import "github.com/discohaus/go-rcon/pkg/rcon"

type executor interface {
	Send(command string) (string, error)
}

type rconExecutor struct{ client *rcon.Client }

func (e rconExecutor) Send(command string) (string, error) { return e.client.Send(command) }
