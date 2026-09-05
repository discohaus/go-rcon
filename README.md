# go-rcon

[![Go Reference](https://pkg.go.dev/badge/github.com/discohaus/go-rcon/pkg/rcon.svg)](https://pkg.go.dev/github.com/discohaus/go-rcon/pkg/rcon)

A robust Minecraft RCON client for Go featuring support for asynchronous command execution and automatic reassembly of fragmented responses.

> **Note:** This package is based on the original [implementation by jltobler](https://github.com/jltobler/go-rcon) and is actively maintained and updated by us.

---

## Features

- **Actively Maintained:** Continuation and modernization of the original `jltobler/go-rcon` codebase.
- **Asynchronous Execution:** `rcon.Client` is thread-safe and creates a separate connection per command for asynchronous execution.
- **Connection Reuse:** Optional low-level connection reuse via `rcon.Conn` and `SendCommand` for resource-efficient, persistent connections.
- **Fragmentation Handling:** Automatically merges and parses multi-packet responses exceeding the 4096-byte RCON limit into a single string.
- **Flexible Character Sets:** Full support for `ASCII`, `ISO-8859-1` (`Latin1`), and `UTF-8`.
- **Interactive CLI:** Built-in interactive TUI terminal for direct server management.

---

## Installation

Use `go get` to add the package to your project:

```sh
go get github.com/discohaus/go-rcon/pkg/rcon
```

---

## Usage

### Import

```go
import "github.com/discohaus/go-rcon/pkg/rcon"

```

### Sending Commands (Client Mode)

The `Client` struct is safe for concurrent use across multiple goroutines.

```go
package main

import (
	"fmt"
	"log"

	"github.com/discohaus/go-rcon/pkg/rcon"
)

func main() {
	// Initialize client
	client := rcon.NewClient("rcon://localhost:25575", "password")

	// Send a command
	response, err := client.Send("time set day")
	if err != nil {
		log.Fatalf("Failed to send command: %v", err)
	}

	fmt.Println("Server Response:", response)
}

```

### Configuring Character Sets (CharSet)

To correctly process Minecraft color codes (`§3`) or special characters, specify the required character set using functional options:

```go
// For newer Minecraft versions or color code support (ISO-8859-1)
client := rcon.NewClient(
	"rcon://localhost:25575", 
	"password", 
	rcon.WithOptions(rcon.CharSetLatin1),
)

```

**Available Character Sets:**

* `CharSetASCII`: For legacy Minecraft servers (default when no options are provided).
* `CharSetLatin1`: For modern servers and Minecraft color codes (`§`).
* `CharSetUTF8`: For modded or custom servers supporting native UTF-8 payloads.

---

## Command Line Interface (CLI)

### Downloads & Setup

Prebuilt binaries are available on the [GitHub Releases page](https://github.com/discohaus/go-rcon/releases).

**Run the binary:**

```sh
./go-rcon --host localhost --port 25575 --password 'your_password' --charset latin1

```

**Run directly from source:**

```sh
go run ./cmd/cli --host localhost --port 25575 --password 'your_password' --charset latin1

```

### CLI Flags

| Flag | Short | Default | Description |
| --- | --- | --- | --- |
| `--host` | `-H` | `localhost` | RCON server hostname or IP address |
| `--port` | `-P` | `25575` | RCON server port (`1`–`65535`) |
| `--password` | `-p` | *empty* | RCON password |
| `--charset` | `-c` | `latin1` | Payload character set: `ascii`, `latin1`, or `utf8` |

*Example using short flags:*

```sh
./go-rcon -H localhost -P 25575 -p 'your_password' -c latin1

```

> **Note:** Invalid host, port, and charset values are validated and reported before launching the TUI.

### TUI Commands

* **Server Commands:** Entered without a leading slash (`/`) and sent directly to the RCON server asynchronously.
```text
time set day
say Server restarting in 5 minutes

```


* **Local TUI Commands:** Prefixed with a leading slash (`/`) and **never** sent to the server.

| Command | Description |
| --- | --- |
| `/help` | Show all available local TUI commands |
| `/exit` | Close the CLI application |

*Tip: Press `/` and use `Tab` to navigate local command auto-completion.*

---

## Development

### Prerequisites

* **Go** (version matching `go.mod`)
* **golangci-lint**
* **Make**

### Available Commands

A `Makefile` is provided to streamline common development workflows:

```sh
make all           # Run everything regarding ci pipeline
make test          # Run test suite
make test-coverage # Run tests with coverage reporting
make lint          # Run golangci-lint
make fmt           # Format code
make tidy          # Tidy Go module dependencies
make clean         # Remove build artifacts
make cli           # Run cli
make cli-dev       # Run local mincraft server ready for make cli

```

