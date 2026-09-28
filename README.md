# 🌐 OpenNet

## OpenNet — экспериментальный сетевой инструмент на Go для построения защищённого соединения через подключаемый внешний транспорт без необходимости поднимать собственный VPS.

«⚠️ Status: Early development / experimental
Version: v1.0»

---

✨ Features

- 🔌 SOCKS5 — локальный SOCKS5-прокси на "127.0.0.1:1080".
- 🌐 WebSocket transport — подключение к внешнему Worker через WebSocket.
- ☁️ Cloudflare Worker transport — внешний транспорт без собственного VPS.
- 🔐 Token authentication — авторизация Worker через Bearer-токен.
- 📦 Pluggable architecture — транспорт отделён от локального SOCKS5-сервера.
- 🐧 Linux — нативная CLI-сборка.
- 📱 Termux — клиент OpenNet может запускаться на Android через Termux.
- 🚀 Go — основной код написан на Go.
- 🛠️ Bootstrap — автоматическая подготовка и сборка проекта.

---

🧠 How it works

OpenNet работает как локальный SOCKS5-прокси:

Application
     │
     ▼
127.0.0.1:1080
     │
     ▼
OpenNet
     │
     │ WebSocket
     ▼
Cloudflare Worker
     │
     │ TCP
     ▼
Target server

OpenNet принимает соединение от приложения через SOCKS5, устанавливает WebSocket-соединение с настроенным Worker и передаёт данные через него.

Собственный VPS для Worker-транспорта не требуется.

---

📂 Project Structure

OpenNet/
├── cmd/
│   └── opennet/
│       └── main.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── protocol/
│   │   └── protocol.go
│   ├── transport/
│   │   └── transport.go
│   └── socks5/
│       └── server.go
│
├── data/
│   ├── config.json
│   └── cache/
│
├── worker/
│   ├── src/
│   │   └── index.js
│   ├── wrangler.jsonc
│   └── README.md
│
├── go.mod
├── .gitignore
└── README.md

---

💻 Requirements

Linux

- Go 1.23+
- Git

Проверить Go:

go version

Android / Termux

pkg install golang git

Проверить:

go version

«Деплой Cloudflare Worker рекомендуется выполнять с обычного Linux/macOS/Windows-компьютера. Wrangler может не поддерживать Android/Termux как платформу для установки.»

---

🚀 Installation

Клонировать репозиторий:

git clone https://github.com/qopix/OpenNet.git
cd OpenNet

Собрать:

go mod tidy
go build -o opennet ./cmd/opennet

---

⚙️ Configuration

Открой:

data/config.json

Пример:

{
  "worker_url": "wss://YOUR-WORKER.workers.dev",
  "token": "YOUR_TOKEN",
  "socks_host": "127.0.0.1",
  "socks_port": 1080
}

Параметры

Параметр| Описание
"worker_url"| WebSocket-адрес Worker
"token"| Токен авторизации
"socks_host"| Локальный адрес SOCKS5
"socks_port"| Локальный порт SOCKS5

---

☁️ Worker

Перейди в каталог Worker:

cd worker

Установи Wrangler:

bun add -g wrangler

Авторизуйся:

wrangler login

Создай секрет:

wrangler secret put OPENNET_TOKEN

Затем выполни:

wrangler deploy

После деплоя Cloudflare выдаст адрес Worker.

В "data/config.json" указывается WebSocket-вариант адреса:

wss://YOUR-WORKER.workers.dev

Токен в OpenNet должен совпадать со значением "OPENNET_TOKEN".

---

▶️ Running

Из корня проекта:

./opennet

При успешном запуске:

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

После запуска приложения могут использовать:

SOCKS5
127.0.0.1:1080

Для остановки:

Ctrl+C

---

🔐 Security

OpenNet находится на ранней стадии разработки.

Не публикуйте:

- "OPENNET_TOKEN";
- приватные ключи, если они появятся в будущих версиях;
- локальные конфигурации с секретами;
- другие credentials.

Не открывайте локальный SOCKS5-порт во внешнюю сеть без необходимости.

Для Worker используется секрет:

wrangler secret put OPENNET_TOKEN

Секрет не должен храниться непосредственно в исходном коде или Git-репозитории.

---

🛠️ Roadmap

v1.x

- [x] Local SOCKS5
- [x] WebSocket transport
- [x] Cloudflare Worker transport
- [x] Token authentication
- [x] Basic configuration
- [x] Linux/Termux CLI
- [x] Bootstrap

Future

- [ ] Multiple transport backends
- [ ] Automatic transport selection
- [ ] Better connection recovery
- [ ] Connection multiplexing
- [ ] Persistent configuration
- [ ] Improved cryptographic protocol
- [ ] Android system-wide VPN/TUN mode
- [ ] Local control interface
- [ ] Transport health checks

---

⚠️ Experimental Software

OpenNet is experimental software.

The current v1.0 implementation is intended for development and testing. It should not be considered a production-grade anonymity, privacy, or censorship-resistance system.

Network behavior can depend on the configured transport, network provider, firewall, DNS configuration, and other external factors.

---

📄 License

OpenNet is free and open-source software distributed under the:

GNU General Public License v3.0 (GPLv3)

---

👤 Project

OpenNet

GitHub:

https://github.com/qopix/OpenNet