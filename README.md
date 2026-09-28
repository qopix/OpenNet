# OpenNet

**Experimental networking tool written in Go** with a local SOCKS5 interface and pluggable transports.

> **Status:** early development / experimental  
> **Version:** v1.0

---

## Features

- Written in Go 1.23+
- Local SOCKS5 proxy
- WebSocket transport
- Cloudflare Worker transport
- Token-based authentication
- Pluggable transport architecture
- Linux support
- Termux / Android CLI support
- JSON configuration
- Bootstrap-based setup

---

## Architecture

```text
Application
     │
     │ SOCKS5
     ▼
┌─────────────┐
│   OpenNet   │
│  Go Client  │
└──────┬──────┘
       │
       │ WebSocket
       ▼
┌─────────────┐
│  Cloudflare │
│    Worker   │
└──────┬──────┘
       │
       │ TCP
       ▼
 Target Server
```

---

## Project Structure

```text
OpenNet/
├── cmd/
│   └── opennet/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── protocol/
│   │   └── protocol.go
│   ├── socks5/
│   │   └── server.go
│   └── transport/
│       └── transport.go
├── data/
│   ├── cache/
│   └── config.json
├── worker/
│   ├── src/
│   │   └── index.js
│   ├── wrangler.jsonc
│   └── README.md
├── .gitignore
├── go.mod
└── README.md
```

---

## Requirements

### Linux

- Go 1.23+
- Git
- Bun or Node.js

Check Go:

```bash
go version
```

### Termux

```bash
pkg update
pkg install golang git
go version
```

> Wrangler deployment should be performed from a desktop Linux, Windows, or macOS system.  
> The Wrangler `workerd` dependency does not currently support Android/Termux natively.

---

## Installation

Clone the repository:

```bash
git clone https://github.com/qopix/OpenNet.git
cd OpenNet
```

Build:

```bash
go mod tidy
go build -o opennet ./cmd/opennet
```

Run:

```bash
./opennet
```

---

## Configuration

Edit `data/config.json`:

```json
{
  "worker_url": "wss://YOUR-WORKER.workers.dev",
  "token": "YOUR_TOKEN",
  "socks_host": "127.0.0.1",
  "socks_port": 1080
}
```

| Option        | Description                      |
|---------------|----------------------------------|
| `worker_url`  | WebSocket transport endpoint     |
| `token`       | Worker authentication token      |
| `socks_host`  | Local SOCKS5 address             |
| `socks_port`  | Local SOCKS5 port                |

Default SOCKS5 address: `127.0.0.1:1080`

---

## Cloudflare Worker

OpenNet includes a Cloudflare Worker transport.

```bash
cd worker
bun add -g wrangler
wrangler login
wrangler secret put OPENNET_TOKEN
wrangler deploy
```

After deployment, Cloudflare provides a Worker URL.  
Use the WebSocket form of that URL in `data/config.json`:

```text
wss://YOUR-WORKER.workers.dev
```

The token in `data/config.json` **must** match the `OPENNET_TOKEN` secret.

---

## Running OpenNet

From the project root:

```bash
./opennet
```

Expected output:

```text
=================================
 OpenNet v1.0
=================================
[OpenNet] Configuration loaded.
[OpenNet] Starting transport...
[OpenNet] Starting SOCKS5...
[OpenNet] Shield: ON
[OpenNet] SOCKS5: 127.0.0.1:1080
OpenNet is running.
Press Ctrl+C to stop.
```

Stop with `Ctrl+C`.

---

## Security

OpenNet is experimental software.

**Do not commit** sensitive information to the repository.

Never publish:

- Authentication tokens
- Private keys
- Credentials
- Local configuration that contains secrets

Store Worker secrets with:

```bash
wrangler secret put OPENNET_TOKEN
```

The SOCKS5 server is intended to listen on `127.0.0.1`.  
Do not expose the SOCKS5 port to external networks unless you understand the security implications.

---

## Roadmap

### OpenNet v1.0

- [x] Go client
- [x] Local SOCKS5
- [x] WebSocket transport
- [x] Cloudflare Worker transport
- [x] Token authentication
- [x] JSON configuration
- [x] Linux support
- [x] Termux CLI support

### Future

- [ ] Multiple transport backends
- [ ] Automatic transport selection
- [ ] Transport health checks
- [ ] Automatic reconnection
- [ ] Connection multiplexing
- [ ] Improved protocol security
- [ ] Android VPN/TUN mode
- [ ] Local management interface
- [ ] Additional transport implementations

---

## Disclaimer

OpenNet is an experimental networking project.

It does **not** guarantee anonymity, privacy, uninterrupted connectivity, or successful operation on every network.

Network behavior depends on the configured transport and the surrounding network environment.

---

## License

OpenNet is distributed under the [GNU General Public License v3.0](LICENSE).

---

## Contributing

Issues, bug reports, ideas, and pull requests are welcome.

If you find a bug or want to propose a new transport, open an Issue or a Pull Request.

---

**GitHub:** [https://github.com/qopix/OpenNet](https://github.com/qopix/OpenNet)
