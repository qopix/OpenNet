🌐 OpenNet

OpenNet — экспериментальный сетевой инструмент на Go с локальным SOCKS5-интерфейсом и подключаемыми транспортами.

«⚠️ Status: Early development / experimental
Version: v1.0»

✨ Features

- 🚀 Written in Go 1.23+
- 🔌 Local SOCKS5 proxy
- 🌐 WebSocket transport
- ☁️ Cloudflare Worker transport
- 🔐 Token-based authentication
- 🧩 Pluggable transport architecture
- 🐧 Linux support
- 📱 Termux / Android CLI support
- ⚙️ JSON configuration
- 📦 Bootstrap-based setup

🏗️ Architecture

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

📂 Project Structure

OpenNet/
├── cmd/
│   └── opennet/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   │
│   ├── protocol/
│   │   └── protocol.go
│   │
│   ├── socks5/
│   │   └── server.go
│   │
│   └── transport/
│       └── transport.go
│
├── data/
│   ├── cache/
│   └── config.json
│
├── worker/
│   ├── src/
│   │   └── index.js
│   ├── wrangler.jsonc
│   └── README.md
│
├── .gitignore
├── go.mod
└── README.md

💻 Requirements

Linux

- Go 1.23+
- Git
- Bun or Node.js

Check Go:

go version

Termux

pkg update
pkg install golang git

Check Go:

go version

«Wrangler deployment should be performed from a desktop Linux, Windows, or macOS system. The Wrangler "workerd" dependency does not currently support Android/Termux natively.»

🚀 Installation

Clone the repository:

git clone https://github.com/qopix/OpenNet.git
cd OpenNet

Build the project:

go mod tidy
go build -o opennet ./cmd/opennet

Run:

./opennet

⚙️ Configuration

Open:

data/config.json

Example:

{
  "worker_url": "wss://YOUR-WORKER.workers.dev",
  "token": "YOUR_TOKEN",
  "socks_host": "127.0.0.1",
  "socks_port": 1080
}

Configuration options

Option| Description
"worker_url"| WebSocket transport endpoint
"token"| Worker authentication token
"socks_host"| Local SOCKS5 address
"socks_port"| Local SOCKS5 port

Default SOCKS5 address:

127.0.0.1:1080

☁️ Cloudflare Worker

OpenNet includes a Cloudflare Worker transport.

Go to the Worker directory:

cd worker

Install Wrangler with Bun:

bun add -g wrangler

Login:

wrangler login

Set the authentication token:

wrangler secret put OPENNET_TOKEN

Deploy the Worker:

wrangler deploy

After deployment, Cloudflare will provide a Worker URL.

Use the WebSocket version of the URL in:

data/config.json

Example:

wss://YOUR-WORKER.workers.dev

The token in "data/config.json" must match the "OPENNET_TOKEN" secret.

▶️ Running OpenNet

From the project root:

./opennet

Expected output:

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

Stop the program with:

Ctrl+C

🔐 Security

OpenNet is experimental software.

Do not commit sensitive information to the repository.

Never publish:

- Authentication tokens
- Private keys
- Credentials
- Local configuration containing secrets

Worker secrets should be stored using:

wrangler secret put OPENNET_TOKEN

The SOCKS5 server is intended to listen on:

127.0.0.1

Avoid exposing the SOCKS5 port to external networks unless you understand the security implications.

🛣️ Roadmap

OpenNet v1.0

- [x] Go client
- [x] Local SOCKS5
- [x] WebSocket transport
- [x] Cloudflare Worker transport
- [x] Token authentication
- [x] JSON configuration
- [x] Linux support
- [x] Termux CLI support

Future

- [ ] Multiple transport backends
- [ ] Automatic transport selection
- [ ] Transport health checks
- [ ] Automatic reconnection
- [ ] Connection multiplexing
- [ ] Improved protocol security
- [ ] Android VPN/TUN mode
- [ ] Local management interface
- [ ] Additional transport implementations

⚠️ Disclaimer

OpenNet is an experimental networking project.

It does not guarantee anonymity, privacy, uninterrupted connectivity, or successful operation on every network.

Network behavior depends on the configured transport and the surrounding network environment.

📜 License

OpenNet is distributed under the GNU General Public License v3.0.

See ""LICENSE"" (LICENSE) for the full license text.

🤝 Contributing

Issues, bug reports, ideas, and pull requests are welcome.

If you find a bug or want to propose a new transport, open an Issue or Pull Request.

⭐ OpenNet

GitHub: https://github.com/qopix/OpenNet