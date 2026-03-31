## Gokus

Gokus is a small terminal Blokus Duo game written in Go.

### Features

This project focuses on the core Blokus Duo rules and a simple TUI play experience.
- Blokus Duo board and piece logic
- Turn-based 2-player flow
- Piece rotate and flip
- Ghost preview before placing
- Terminal UI built with Bubble Tea + Lip Gloss

### TODO
- SSH server for remote play
- AI opponent?

### Controls

- Arrow keys: move cursor
- Tab / Shift+Tab: next / previous piece
- Enter: place piece
- r: rotate piece
- f: flip piece
- s: skip turn
- n: new game
- q: quit

### Run
```bash
go run ./cmd/gokus
```

Or with Makefile:

```bash
make build
make run
```

### Test

```bash
go test ./internal/... -v
```
