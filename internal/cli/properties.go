package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// RconFileConfig holds RCON settings parsed from a configuration file (e.g. server.properties).
type RconFileConfig struct {
	Port     *int32
	Password *string
}

// ReadRconProperties reads a properties file and extracts rcon.password and rcon.port if present.
func ReadRconProperties(filePath string) (*RconFileConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open properties file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	config := &RconFileConfig{}
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			key, value, found = strings.Cut(line, ":")
		}
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch key {
		case "rcon.password":
			pwd := value
			config.Password = &pwd
		case "rcon.port":
			portNum, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid rcon.port %q at line %d: must be an integer", value, lineNum)
			}
			if portNum < 1 || portNum > 65535 {
				return nil, fmt.Errorf("invalid rcon.port %d at line %d: must be between 1 and 65535", portNum, lineNum)
			}
			port := int32(portNum)
			config.Port = &port
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading properties file: %w", err)
	}

	return config, nil
}
