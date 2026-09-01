## Gokus

Gokus is a small terminal Blokus Duo game.

![Demo](./gokus.png)

### Features
- Blokus Duo board and piece logic
- Turn-based 2-player flow
- Piece rotate and flip
- Ghost preview before placing
- Terminal UI built with Bubble Tea + Lip Gloss
- Design based on this nice paper: [Implementing Minimax Search with Alpha-Beta
Pruning in Blokus Duo](https://informatika.stei.itb.ac.id/~rinaldi.munir/Stmik/2024-2025/Makalah2025/Makalah-IF2211-Strategi-Algoritma-2025%20(99).pdf)

### Todo
- [x] SSH server for remote play
- Detect game over and show winner
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

### SSH server

Run the multiplayer SSH server locally:

```bash
go run ./cmd/gokus-ssh
ssh -p 23234 ROOM_NAME@localhost
```

The SSH username is used as the room name. Share the same room name with the
other player.

The server supports these environment variables:

| Variable | Default | Description |
| --- | --- | --- |
| `GOKUS_SSH_ADDRESS` | `localhost:23234` | Address the SSH server listens on |
| `GOKUS_SSH_HOST_KEY_PATH` | `.ssh/gokus_ed25519` | Persistent SSH host key path |

### Docker

Run the published image from GitHub Container Registry:

```bash
docker run --rm -it \
  -p 23234:23234 \
  -v gokus-ssh-data:/data \
  ghcr.io/livensmi1e/gokus:latest
```

Then connect with:

```bash
ssh -p 23234 ROOM_NAME@localhost
```

The `/data` volume preserves the SSH host key between container replacements.
To build and run from source instead, use `docker compose up --build`.

The GitHub Actions container workflow validates pull requests without pushing.
Pushes to `main` publish `main`, `latest`, and `sha-*` tags. Tags such as
`v1.2.3` publish `1.2.3`, `1.2`, `1`, `latest`, and `sha-*` tags to
`ghcr.io/<owner>/<repository>` using the repository's `GITHUB_TOKEN`.

GHCR packages are private by default. If this image should be pullable without
authentication, change the package visibility to public in the package
settings after the first publish.
