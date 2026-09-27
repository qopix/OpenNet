# 🌐 OpenNet

**OpenNet** — это open-source сетевой проект на языке Go, предназначенный для создания расширяемой сетевой инфраструктуры с автоматическим выбором доступного транспорта и защищённым соединением.

> ⚠️ **Status:** Early development / experimental

---

### ✨ Features

* **Language:** Written entirely in Go (1.23+)
* **Cross-platform:** Native Linux support and Termux-compatible CLI
* **Control Panel:** Local web-based management panel
* **Security:** Cryptographic keys based on Ed25519 architecture
* **Flexibility:** Pluggable transport architecture with local configuration
* **Independence:** No mandatory VPS or centralized servers required

---

### 📂 Project Structure

```text
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
```

---

### 💻 Requirements

* **Go 1.23** or newer

#### On Linux:
```bash
go version
```

#### On Termux:
```bash
pkg install golang git
go version
```

---

### 🚀 Build & Run

#### 1. Clone the repository:
```bash
git clone https://github.com/qopix/OpenNet.git
cd OpenNet
```

#### 2. Build the project:
```bash
go build -o opennet ./cmd/opennet
```

#### 3. Run the executable:
```bash
./opennet
```

После запуска локальная панель управления будет доступна по адресу:
👉 `http://127.0.0.1:8765`

---

### 🔒 Security Note

* OpenNet is experimental software.
* **Do not expose** the local control panel to the public Internet.
* Private keys should be treated as secrets and should **never** be published in Git repositories.

---

### 🛠️ Development Roadmap

Текущий фокус разработки сосредоточен на проектировании и реализации базовой архитектуры:
1. 🔑 Key management
2. ⚙️ Local configuration
3. 🔄 Transport abstraction
4. 🛡️ Secure connections
5. 📡 Client/server communication
6. 🔀 Automatic transport selection

---

### 📄 License

OpenNet is open-source software licensed under the **GPLv3 License**.
