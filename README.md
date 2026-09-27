OpenNet

OpenNet — open-source сетевой проект на Go.

Цель проекта — предоставить расширяемую сетевую инфраструктуру с автоматическим выбором доступного транспорта и защищённым соединением.

«Status: Early development / experimental»

Features

- Written in Go
- Cross-platform architecture
- Linux support
- Termux-compatible CLI
- Local web control panel
- Ed25519 cryptographic keys
- Local configuration
- Pluggable transport architecture
- No mandatory VPS
- Open-source

Project structure

OpenNet/
├── cmd/
│   └── opennet/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── keys/
│   │   └── keys.go
│   ├── transport/
│   │   └── transport.go
│   └── web/
│       └── web.go
├── web/
│   └── index.html
├── go.mod
├── .gitignore
└── README.md

Requirements

- Go 1.23 or newer

On Linux:

go version

On Termux:

pkg install golang git
go version

Build

Clone the repository:

git clone https://github.com/qopix/OpenNet.git
cd OpenNet

Build:

go build -o opennet ./cmd/opennet

Run:

./opennet

The local control panel is available at:

http://127.0.0.1:8765

Security

OpenNet is experimental software.

Do not expose the local control panel to the public Internet.

Private keys should be treated as secrets and should never be published in Git repositories.

Development

The project is currently focused on the core architecture:

1. Key management
2. Local configuration
3. Transport abstraction
4. Secure connections
5. Client/server communication
6. Automatic transport selection

License

OpenNet is open source. The project license is GPLv3