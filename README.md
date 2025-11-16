# PostgresMaster 🐘

[![Tests](https://github.com/MPJHorner/PostgresMaster/actions/workflows/test.yml/badge.svg)](https://github.com/MPJHorner/PostgresMaster/actions/workflows/test.yml)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/MPJHorner/PostgresMaster)](https://goreportcard.com/report/github.com/MPJHorner/PostgresMaster)

> **A browser-based PostgreSQL client that requires zero backend infrastructure.**

Connect to any Postgres database through a lightweight local proxy, with intelligent SQL autocomplete and a clean, modern interface. No server setup, no configuration files, no hassle.

## 📸 Screenshots

> **Coming Soon**: Screenshots and demo GIFs will be added once the web application is deployed.

## ✨ Key Features

### 🌐 Zero Infrastructure
- **100% Frontend Web App**: Pure static site - no server infrastructure, no backend to maintain
- **Deploy Anywhere**: Host on Cloudflare Pages, Netlify, GitHub Pages, or any static hosting
- **No Cloud Lock-in**: Your data and queries never touch our servers

### 🚀 Powerful SQL Editor
- **Monaco Editor**: The same editor that powers VS Code, right in your browser
- **Intelligent Autocomplete**: Context-aware SQL completion with:
  - PostgreSQL keywords and functions
  - Table and column names from your schema
  - Data types and constraints
- **Multi-line Queries**: Write complex queries with proper formatting
- **Keyboard Shortcuts**: Execute queries with `Ctrl+Enter`

### 🔒 Security First
- **Local-Only Proxy**: Lightweight ~5MB Go binary runs on your machine
- **Credentials Stay Local**: Database credentials never leave your computer
- **Secure WebSocket**: Authenticated connection between browser and local proxy
- **SSL/TLS Support**: Connect to Postgres databases with SSL enabled

### 📊 Query Results
- **Tabular Display**: Clean, readable results with column headers and types
- **Type-Aware Formatting**: Proper handling of NULL, JSON, timestamps, and more
- **Execution Metadata**: View row count and execution time for each query
- **Query History**: Track recent queries during your session (in-memory)

### 🎯 Developer Experience
- **Cross-Platform**: Binaries for Windows, macOS (Intel + ARM), and Linux (amd64 + arm64)
- **No Installation**: Just download and run - no dependencies required
- **Connection Flexibility**: Use connection strings or interactive mode
- **Retry Logic**: Automatic retry with exponential backoff for flaky connections
- **Graceful Shutdown**: Clean connection cleanup on exit

## 🏗️ Architecture

```
Web App (Browser)
    ↓ WebSocket (ws://localhost:8080?secret=xxx)
Go Proxy (Local)
    ↓ Postgres Protocol (TCP + SSL)
Remote Postgres Server
```

## 🚀 Quick Start

Get up and running in under 60 seconds:

### 1. Download Proxy

Download the binary for your platform from [releases](../../releases):

| Platform | Download |
|----------|----------|
| **Windows (x64)** | `postgres-proxy-windows-amd64.exe` |
| **macOS (Intel)** | `postgres-proxy-darwin-amd64` |
| **macOS (Apple Silicon)** | `postgres-proxy-darwin-arm64` |
| **Linux (x64)** | `postgres-proxy-linux-amd64` |
| **Linux (ARM64)** | `postgres-proxy-linux-arm64` |

**macOS/Linux**: Make the binary executable:
```bash
chmod +x postgres-proxy-*
```

### 2. Connect to Your Database

**Option A: Connection String** (Fastest)
```bash
./postgres-proxy "postgres://user:pass@host:5432/dbname"
```

**Option B: Interactive Mode** (Most Secure - password not in shell history)
```bash
./postgres-proxy
Host: db.example.com
Port [5432]:
Database: mydb
Username: admin
Password: ****
SSL Mode [prefer]: require

Connecting to postgres://db.example.com:5432/mydb...
✓ Connected!

WebSocket server listening on :8080
→ Open in browser: http://localhost:8080?secret=a1b2c3d4e5f6
```

### 3. Start Querying

1. **Copy the URL** from the proxy output (includes the secret)
2. **Open in your browser** - the web app will auto-connect
3. **Start writing SQL** with intelligent autocomplete!

```sql
-- Try it out!
SELECT * FROM pg_tables WHERE schemaname = 'public';
```

**Tip**: Press `Ctrl+Enter` to execute your query.

## 💡 What Makes This Different?

Unlike traditional database clients:

- **No Installation Required**: Just download and run - works on any OS
- **No Cloud Service**: Your credentials and queries stay on your machine
- **No Backend to Maintain**: Static web app means zero server costs
- **No Complex Setup**: Works with any Postgres database in seconds
- **No Vendor Lock-in**: Open source (AGPL-3.0), deploy anywhere

Perfect for:
- 🧪 Quick database exploration
- 🔍 Testing queries on remote databases
- 📚 Learning PostgreSQL
- 🚀 Lightweight alternative to pgAdmin or DBeaver
- 🔒 Secure access to production databases without exposing credentials

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

## 📚 Documentation

### User Documentation
- **[Installation Guide](docs/installation.md)**: Detailed installation instructions for all platforms
- **[Usage Guide](docs/usage.md)**: Query examples, autocomplete tips, and keyboard shortcuts
- **[FAQ](docs/installation.md#faq)**: Common questions and troubleshooting

### Developer Documentation
- **[Web App Development](web/README.md)**: Frontend development setup and guidelines
- **[Proxy Development](proxy/README.md)**: Go proxy development and testing
- **[Architecture Overview](docs/architecture.md)**: System design and technical details
- **[Contributing Guide](CONTRIBUTING.md)**: How to contribute to the project

### Quick Links
- **[PRD & Implementation Plan](prd.md)**: Complete product requirements and development roadmap
- **[Releases](../../releases)**: Download binaries and view changelog

## 🗺️ Roadmap

### ✅ Current Version (v0.1.0 - MVP)
- ✅ Go proxy with WebSocket server
- ✅ Static web app with SvelteKit
- ✅ Monaco SQL editor with autocomplete
- ✅ Query execution and results display
- ✅ Schema introspection
- ✅ Query history (in-memory)
- ✅ Cross-platform binaries

### 🚧 Next Version (v0.2.0)
- [ ] Query history persistence (LocalStorage)
- [ ] Export results to CSV/JSON
- [ ] Multiple query tabs
- [ ] Dark/light theme toggle
- [ ] Query templates/snippets
- [ ] Result pagination for large datasets
- [ ] Keyboard shortcuts panel
- [ ] Schema sidebar with tree view

### 🔮 Future (v0.3.0+)
- [ ] Query visualization (charts/graphs)
- [ ] EXPLAIN plan visualization
- [ ] Query formatting and linting
- [ ] Saved queries
- [ ] Connection profiles
- [ ] Auto-update mechanism for proxy

See the full roadmap in [prd.md](prd.md).

## 🎯 Project Status

**Current Phase**: Phase 5 - Testing & Polish 🚧

| Phase | Status | Progress |
|-------|--------|----------|
| Phase 1: Go Proxy Foundation | ✅ Complete | 100% |
| Phase 2: Web App Foundation | ✅ Complete | 100% |
| Phase 3: SQL Editor & Autocomplete | ✅ Complete | 100% |
| Phase 4: Query Execution & Results | ✅ Complete | 100% |
| Phase 5: Testing & Polish | 🚧 In Progress | 85% |
| Phase 6: Deployment & Release | ⏳ Pending | 0% |

See [prd.md](prd.md) for the complete implementation plan and detailed progress tracking.

## 🤝 Contributing

We welcome contributions! Here's how you can help:

- 🐛 **Report Bugs**: Open an issue with details and steps to reproduce
- 💡 **Suggest Features**: Share your ideas in the discussions
- 🔧 **Submit PRs**: Check our [Contributing Guide](CONTRIBUTING.md) for guidelines
- 📖 **Improve Docs**: Help make our documentation better
- ⭐ **Star the Repo**: Show your support!

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## 📝 License

This project is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**.

**What this means**:
- ✅ You can use, modify, and distribute this software freely
- ✅ You can use it commercially
- ⚠️ If you modify and deploy this software (including as a web service), you **must**:
  - Make your source code available under AGPL-3.0
  - Provide attribution to the original project
  - Disclose your changes

See [LICENSE](LICENSE) for the full legal text.

## 🙏 Credits & Attribution

### Built With
- **[SvelteKit](https://kit.svelte.dev/)**: Web application framework
- **[Monaco Editor](https://microsoft.github.io/monaco-editor/)**: The code editor that powers VS Code
- **[shadcn-svelte](https://www.shadcn-svelte.com/)**: Beautiful UI components for Svelte
- **[TailwindCSS](https://tailwindcss.com/)**: Utility-first CSS framework
- **[Go](https://go.dev/)**: Programming language for the proxy
- **[gorilla/websocket](https://github.com/gorilla/websocket)**: WebSocket implementation
- **[jackc/pgx](https://github.com/jackc/pgx)**: PostgreSQL driver and toolkit

### Inspired By
- **[pgAdmin](https://www.pgadmin.org/)**: The industry-standard Postgres GUI
- **[DBeaver](https://dbeaver.io/)**: Universal database tool
- **[Datasette](https://datasette.io/)**: Inspiration for the static-first approach

### Special Thanks
- The PostgreSQL community for building an amazing database
- All open source contributors who made this project possible
- Early testers and feedback providers

## 📞 Support & Community

- **Issues**: [GitHub Issues](../../issues)
- **Discussions**: [GitHub Discussions](../../discussions)
- **Website**: Coming soon

---

<div align="center">

**Built with ❤️ for the PostgreSQL community**

[⭐ Star on GitHub](../../stargazers) • [🐛 Report Bug](../../issues) • [💡 Request Feature](../../issues)

</div>
