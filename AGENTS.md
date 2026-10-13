# AGENTS.md

Single-module Go repo (`github.com/discohaus/go-rcon`, Go 1.26.6): RCON client library + Bubble Tea CLI.

## Layout

- `pkg/rcon/` — public library. `client.go` (`Client`: per-command dial/close, goroutine-safe), `conn.go` (`Conn`/`Dial`: persistent reusable connection, `SendCommand` mutex-serialized), `packet.go`, `charset.go`, `errors.go`.
- `internal/cli/` — CLI logic: `cli.go` (`NewCli` validation + `ExecuteCommand`), `tea.go` (interactive TUI model), `commands.go`, `executor.go`, `properties.go` (`ReadRconProperties`).
- `cmd/cli/main.go` — cobra entrypoint. Release binary `go-rcon` built from here (see `.goreleaser.yaml`).
- `data/` — gitignored local Minecraft server volume for `make cli-dev`. Never commit.

## Commands (`Makefile` is source of truth; CI mirrors it)

```sh
make all           # lint + test (CI gate)
make test          # go test -v -race ./...
make lint          # golangci-lint run (config: .golangci.yml, v2.13.2 in CI)
make fmt           # go fmt ./...
make tidy          # go mod tidy
go test -race -run TestName ./pkg/rcon/   # single test
go test -race ./pkg/rcon/ ./internal/cli/ # focused package
```

- `make test-coverage` opens HTML in a browser — avoid headless; use `go test -race -coverprofile=coverage.out ./... && go tool cover -func=coverage.out` instead.
- `make cli-dev` starts `itzg/minecraft-server` docker (RCON password `1234`, port 25575, mounts `./data`); then `make cli` runs `go run cmd/cli/main.go -i -H localhost -P 25575 -p 1234`.
- Release: tag `v*.*.*` triggers GoReleaser (`release.yml`); do not hand-build archives.

## Gotchas

- `Dial` requires an `rcon://host:port` URL (scheme validated, port defaults to 25575). CLI builds it as `rcon://host:port` — pass the same form to `NewClient`.
- Charset default differs: `rcon.NewClient` defaults to `CharSetASCII`; CLI flag `--charset/-C` defaults to `latin1`. Use `rcon.WithOptions(rcon.CharSetLatin1)` for `§` color codes, `CharSetUTF8` only for modded servers.
- `Conn` cannot queue requests (Minecraft limitation): one `SendCommand` at a time, termination-packet flow in `conn.go:readPackets`. For concurrent use, use `Client.Send` (dials per command), not a shared `Conn`.
- CLI flag rules (enforced in `main.go`): exactly one of `-i` / `-c` required; `-f` properties file (`rcon.password`, `rcon.port`) only fills in `port`/`password` when those flags were NOT explicitly passed.
- Tests are self-contained unit tests (`*_test.go` next to code); no live server or docker needed for `make test`.
