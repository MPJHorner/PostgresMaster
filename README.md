# PostgresMaster 🐘

A browser-based PostgreSQL client that requires **zero backend infrastructure**. Connect to any Postgres database through a lightweight local proxy, with intelligent SQL autocomplete and a clean, modern interface.

## ✨ Key Features

- **100% Frontend Web App**: Pure static site - no server infrastructure needed
- **Lightweight Local Proxy**: Single ~5MB Go binary bridges browser to Postgres
- **Intelligent Autocomplete**: Schema-aware SQL completion with Monaco Editor
- **Universal Access**: Works with any PostgreSQL database, anywhere
- **Secure by Design**: Credentials never leave your machine

## 🏗️ Architecture

```
Web App (Browser)
    ↓ WebSocket (ws://localhost:8080?secret=xxx)
Go Proxy (Local)
    ↓ Postgres Protocol (TCP + SSL)
Remote Postgres Server
```

## 🚀 Quick Start

### 1. Download Proxy

Download the binary for your platform from [releases](../../releases):

- Windows (x64): `postgres-proxy-windows-amd64.exe`
- macOS (Intel): `postgres-proxy-darwin-amd64`
- macOS (Apple Silicon): `postgres-proxy-darwin-arm64`
- Linux (x64): `postgres-proxy-linux-amd64`

### 2. Connect to Your Database

**Option A: Connection String**
```bash
./postgres-proxy "postgres://user:pass@host:5432/dbname"
```

**Option B: Interactive Mode**
```bash
./postgres-proxy
Host: db.example.com
Port [5432]:
Database: mydb
Username: admin
Password: ****
SSL Mode [prefer]: require

✓ Connected!
→ Open in browser: http://localhost:8080?secret=a1b2c3d4e5f6
```

### 3. Start Querying

Open the URL in your browser and start writing SQL with intelligent autocomplete!

## 🛠️ Tech Stack

### Web Application
- **Framework**: SvelteKit (static adapter)
- **UI**: shadcn-svelte + TailwindCSS
- **Editor**: Monaco Editor with SQL language support
- **Type Safety**: TypeScript

### Local Proxy
- **Language**: Go 1.21+
- **Size**: ~5-10MB single binary
- **Platforms**: Windows, macOS (Intel + ARM), Linux (amd64 + arm64)
- **Libraries**: gorilla/websocket, jackc/pgx

## 📦 Project Structure

```
postgres-client/
├── web/          # SvelteKit frontend application
├── proxy/        # Go WebSocket ↔ Postgres bridge
├── docs/         # Documentation
└── prd.md        # Product requirements & implementation plan
```

## 🧑‍💻 Development

See individual README files in each directory:

- [Web App Development](web/README.md)
- [Proxy Development](proxy/README.md)
- [Architecture Documentation](docs/architecture.md)

## 📝 License

This project is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**.

See [LICENSE](LICENSE) for details.

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 🎯 Project Status

**Current Phase**: Phase 1 - Project Structure Setup ✅

See [prd.md](prd.md) for the complete implementation plan and progress tracking.

---

**Built with ❤️ for the PostgreSQL community**
