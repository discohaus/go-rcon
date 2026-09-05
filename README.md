# go-rcon

The Package is very much based on [jltobler's great implementation](https://github.com/jltobler/go-rcon).

He doesn't seem to actively maintain it anymore so we do.

Minecraft RCON client module for connecting to Minecraft server using [RCON](https://minecraft.wiki/RCON) protocol written in Go.

`rcon.Client` features a concurrent-safe `Send` function that creates a separate RCON connection for each command enabling asynchronous execution if desired. 
The ability to reuse a single connection is also available by creating `rcon.Conn` directly and sending commands via `SendCommand` function until the connection is closed.

Fragmented responses over 4096 bytes are also handled and parsed into single response.

## Getting Started

### Installing

`go get` *will always pull the latest tagged release from the main branch.*

```sh
go get github.com/discohaus/go-rcon
```

## CLI

The interactive CLI is available as a prebuilt download on the
[GitHub Releases page](https://github.com/discohaus/go-rcon/releases). Download
the archive for your operating system and architecture, extract it, and run
the `go-rcon` binary.

```sh
go-rcon --host localhost --port 25575 --password 'password' --charset latin1
```

To run the CLI from a checkout instead:

```sh
go run ./cmd/cli --host localhost --port 25575 --password 'password' --charset latin1
```

### CLI flags

| Flag | Short | Default | Description |
| --- | --- | --- | --- |
| `--host` | `-H` | `localhost` | RCON server hostname or IP address |
| `--port` | `-P` | `25575` | RCON server port (`1`-`65535`) |
| `--password` | `-p` | empty | RCON password |
| `--charset` | `-c` | `latin1` | Payload charset: `ascii`, `latin1`, or `utf8` |

Example using short flags:

```sh
./go-rcon -H mc.example.com -P 25575 -p 'password' -c latin1
```

Invalid host, port, and charset values are reported before the TUI starts.

### TUI commands

Server commands are entered without a leading slash and are sent to the RCON
server as written:

```text
time set day
say Server restart in 5 minutes
```

Commands beginning with `/` are reserved for the local TUI and are never sent
to the server:

| Command | Description |
| --- | --- |
| `/help` | Show all available local TUI commands |
| `/exit` | Close the CLI |

Completion is available for local commands. Type `/` to open the palette, type
to filter the command names, press `Tab` to accept the selected suggestion, and
use the arrow keys to navigate the suggestions. Server commands run
asynchronously, so the TUI remains responsive while waiting for a response.

## Go package usage

Import the package into your project.

```go
import "github.com/discohaus/go-rcon/pkg/rcon"
```

Construct a new RCON client which can be used to access the send function.

```go
rconClient := rcon.NewClient("rcon://localhost:25575", "password")
```

Use the Send function to request commands be remotely executed on the Minecraft server.

```go
rconClient, err := rconClient.Send("time set day")
```

Use the `CharSet` type to specify the required character set for the RCON connection. This constraints sent and received data to the specified character set.

```go
rconClient := rcon.NewClient("rcon://localhost:25575", "password", WithOptions(rcon.CharSetASCII))
```
For older Minecraft versions use `CharSetASCII`. This is also the default character set when no options are specified.

For newer Minecraft versions or those that sends Minecraft Color codes such as `§3` use `CharSetLatin1` or for servers understanding UTF8 you can use `CharSetUTF8`.

## Development

### Prerequisites

- Go (defined in `go.mod`)
- [`golangci-lint`](https://golangci-lint.run/)

### Available Commands

A `Makefile` is provided to simplify common development operations:

- **Run tests:**
  ```sh
  make test
  ```
- **Run tests with coverage:**
  ```sh
  make test-coverage
  ```
- **Run linter:**
  ```sh
  make lint
  ```
- **Format code:**
  ```sh
  make fmt
  ```
- **Tidy Go modules:**
  ```sh
  make tidy
  ```
- **Clean build artifacts:**
  ```sh
  make clean
  ```
